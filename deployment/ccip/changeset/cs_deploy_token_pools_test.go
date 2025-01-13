package changeset

import (
	"math/big"
	"testing"
	"time"

	"go.uber.org/zap/zapcore"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/stretchr/testify/require"
)

const localTokenDecimals = 18
const testTokenSymbol TokenSymbol = "TEST"

// createSymmetricRateLimits is a utility to quickly create a rate limiter config with equal inbound and outbound values
func createSymmetricRateLimits(rate int64, capacity int64) RateLimiterConfig {
	return RateLimiterConfig{
		Inbound: token_pool.RateLimiterConfig{
			IsEnabled: rate != 0 || capacity != 0,
			Rate:      big.NewInt(rate),
			Capacity:  big.NewInt(capacity),
		},
		Outbound: token_pool.RateLimiterConfig{
			IsEnabled: rate != 0 || capacity != 0,
			Rate:      big.NewInt(rate),
			Capacity:  big.NewInt(capacity),
		},
	}
}

// setup2ChainEnvironment preps the environment for token pool deployment testing
func setup2ChainEnvironment(t *testing.T) (deployment.Environment, uint64, uint64, map[uint64]*deployment.ContractDeploy[*burn_mint_erc677.BurnMintERC677], map[uint64]*proposalutils.TimelockExecutionContracts) {
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 2,
	})
	selectors := e.AllChainSelectors()

	addressBook := deployment.NewMemoryAddressBook()
	var prereqCfg []DeployPrerequisiteConfigPerChain
	for _, selector := range selectors {
		prereqCfg = append(prereqCfg, DeployPrerequisiteConfigPerChain{
			ChainSelector: selector,
		})
	}

	mcmsCfg := make(map[uint64]commontypes.MCMSWithTimelockConfig)
	for _, selector := range selectors {
		mcmsCfg[selector] = proposalutils.SingleGroupTimelockConfig(t)
	}

	// Deploy one burn-mint token per chain to use in the tests
	tokens := make(map[uint64]*deployment.ContractDeploy[*burn_mint_erc677.BurnMintERC677])
	for _, selector := range selectors {
		token, err := deployment.DeployContract(e.Logger, e.Chains[selector], addressBook,
			func(chain deployment.Chain) deployment.ContractDeploy[*burn_mint_erc677.BurnMintERC677] {
				tokenAddress, tx, token, err := burn_mint_erc677.DeployBurnMintERC677(
					e.Chains[selector].DeployerKey,
					e.Chains[selector].Client,
					string(testTokenSymbol),
					string(testTokenSymbol),
					localTokenDecimals,
					big.NewInt(0).Mul(big.NewInt(1e9), big.NewInt(1e18)),
				)
				return deployment.ContractDeploy[*burn_mint_erc677.BurnMintERC677]{
					Address:  tokenAddress,
					Contract: token,
					Tv:       deployment.NewTypeAndVersion(BurnMintToken, deployment.Version1_0_0),
					Tx:       tx,
					Err:      err,
				}
			},
		)
		require.NoError(t, err)
		tokens[selector] = token
	}

	// Deploy MCMS setup & prerequisite contracts
	e, err := commonchangeset.ApplyChangesets(t, e, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(DeployPrerequisites),
			Config: DeployPrerequisiteConfig{
				Configs: prereqCfg,
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployMCMSWithTimelock),
			Config:    mcmsCfg,
		},
	})
	require.NoError(t, err)

	state, err := LoadOnchainState(e)
	require.NoError(t, err)

	// We only need the token admin registry to be owned by the timelock in these tests
	timelockOwnedContractsByChain := make(map[uint64][]common.Address)
	for _, selector := range selectors {
		timelockOwnedContractsByChain[selector] = []common.Address{state.Chains[selector].TokenAdminRegistry.Address()}
	}

	// Assemble map of addresses required for Timelock scheduling & execution
	timelockContracts := make(map[uint64]*proposalutils.TimelockExecutionContracts)
	for _, selector := range selectors {
		timelockContracts[selector] = &proposalutils.TimelockExecutionContracts{
			Timelock:  state.Chains[selector].Timelock,
			CallProxy: state.Chains[selector].CallProxy,
		}
	}

	// Transfer ownership of token admin registry to the Timelock
	e, err = commonchangeset.ApplyChangesets(t, e, timelockContracts, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.TransferToMCMSWithTimelock),
			Config: commonchangeset.TransferToMCMSWithTimelockConfig{
				ContractsByChain: timelockOwnedContractsByChain,
				MinDelay:         0,
			},
		},
	})
	require.NoError(t, err)

	return e, selectors[0], selectors[1], tokens, timelockContracts
}

// validateMemberOfBurnMintPair performs checks required to validate that a token pool is fully configured for cross-chain transfer.
// Assumes that the deployed token pools are burn-mint.
func validateMemberOfBurnMintPair(
	t *testing.T,
	state CCIPOnChainState,
	tokens map[uint64]*deployment.ContractDeploy[*burn_mint_erc677.BurnMintERC677],
	tokenSymbol TokenSymbol,
	chainSelector uint64,
	rate int64,
	capacity int64,
	externalAdmin *bind.TransactOpts,
) {
	timelockAddress := state.Chains[chainSelector].Timelock.Address()
	tokenPool, ok := state.Chains[chainSelector].BurnMintTokenPools[testTokenSymbol]
	require.True(t, ok)

	// Verify that the timelock is the owner
	owner, err := tokenPool.Owner(nil)
	require.NoError(t, err)
	require.Equal(t, timelockAddress, owner)

	// Fetch the supported remote chains
	supportedChains, err := tokenPool.GetSupportedChains(nil)
	require.NoError(t, err)

	// Verify that the rate limits and remote addresses are correct
	for _, supportedChain := range supportedChains {
		inboundConfig, err := tokenPool.GetCurrentInboundRateLimiterState(nil, supportedChain)
		require.NoError(t, err)
		require.Equal(t, true, inboundConfig.IsEnabled)
		require.Equal(t, big.NewInt(capacity), inboundConfig.Capacity)
		require.Equal(t, big.NewInt(rate), inboundConfig.Rate)

		outboundConfig, err := tokenPool.GetCurrentOutboundRateLimiterState(nil, supportedChain)
		require.NoError(t, err)
		require.Equal(t, true, outboundConfig.IsEnabled)
		require.Equal(t, big.NewInt(capacity), outboundConfig.Capacity)
		require.Equal(t, big.NewInt(rate), outboundConfig.Rate)

		remoteTokenAddress, err := tokenPool.GetRemoteToken(nil, supportedChain)
		require.NoError(t, err)
		require.Equal(t, tokens[supportedChain].Address.Bytes(), remoteTokenAddress)

		remoteBurnMintPool, ok := state.Chains[supportedChain].BurnMintTokenPools[testTokenSymbol]
		require.True(t, ok)
		remotePoolAddresses, err := tokenPool.GetRemotePools(nil, supportedChain)
		require.NoError(t, err)
		require.Equal(t, [][]byte{remoteBurnMintPool.Address().Bytes()}, remotePoolAddresses)
	}

	// Verify that the pool is set on the registry
	tokenConfigOnRegistry, err := state.Chains[chainSelector].TokenAdminRegistry.GetTokenConfig(nil, tokens[chainSelector].Address)
	require.NoError(t, err)
	require.Equal(t, timelockAddress, tokenConfigOnRegistry.Administrator)
	require.Equal(t, tokenPool.Address(), tokenConfigOnRegistry.TokenPool)

	if externalAdmin != nil {
		// Verify that the pending administrator is the external admin
		require.Equal(t, externalAdmin.From, tokenConfigOnRegistry.PendingAdministrator)

		// Accept ownership using deployer key
		_, err := tokenPool.AcceptOwnership(externalAdmin)
		require.NoError(t, err)
	}
}

func TestDeployTokenPoolContracts_DeployNew(t *testing.T) {
	t.Parallel()
	e, selectorA, selectorB, tokens, timelockContracts := setup2ChainEnvironment(t)

	e, err := commonchangeset.ApplyChangesets(t, e, timelockContracts, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(DeployTokenPoolContracts),
			Config: DeployTokenPoolContractsConfig{
				Symbol:        testTokenSymbol,
				TimelockDelay: 0 * time.Second,
				NewPools: map[uint64]NewTokenPoolInput{
					selectorA: {
						BaseTokenPoolInput: BaseTokenPoolInput{
							TokenAddress: tokens[selectorA].Address,
							RemoteChainsToAdd: RemoteChains{
								selectorB: createSymmetricRateLimits(100, 1000),
							},
						},
						Type:               BurnMintTokenPool,
						LocalTokenDecimals: localTokenDecimals,
					},
					selectorB: {
						BaseTokenPoolInput: BaseTokenPoolInput{
							TokenAddress: tokens[selectorB].Address,
							RemoteChainsToAdd: RemoteChains{
								selectorA: createSymmetricRateLimits(100, 1000),
							},
						},
						Type:               BurnMintTokenPool,
						LocalTokenDecimals: localTokenDecimals,
					},
				},
			},
		},
	})
	require.NoError(t, err)

	state, err := LoadOnchainState(e)
	require.NoError(t, err)

	for _, selector := range []uint64{selectorA, selectorB} {
		validateMemberOfBurnMintPair(t, state, tokens, testTokenSymbol, selector, 100, 1000, nil)
	}
}

func TestDeployTokenPoolContracts_DeployNewAndUpdateExisting(t *testing.T) {
	t.Parallel()
	e, selectorA, selectorB, tokens, timelockContracts := setup2ChainEnvironment(t)

	// Deploy the first token pool on chain A
	e, err := commonchangeset.ApplyChangesets(t, e, timelockContracts, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(DeployTokenPoolContracts),
			Config: DeployTokenPoolContractsConfig{
				Symbol:        testTokenSymbol,
				TimelockDelay: 0 * time.Second,
				NewPools: map[uint64]NewTokenPoolInput{
					selectorA: {
						BaseTokenPoolInput: BaseTokenPoolInput{
							TokenAddress: tokens[selectorA].Address,
						},
						Type:               BurnMintTokenPool,
						LocalTokenDecimals: localTokenDecimals,
					},
				},
			},
		},
	})
	require.NoError(t, err)

	state, err := LoadOnchainState(e)
	require.NoError(t, err)

	timelockAddress := state.Chains[selectorA].Timelock.Address()
	burnMintTokenPool := state.Chains[selectorA].BurnMintTokenPools[testTokenSymbol]

	// Verify that the timelock is the owner of the token pool on chain A
	owner, err := burnMintTokenPool.Owner(nil)
	require.NoError(t, err)
	require.Equal(t, timelockAddress, owner)

	// Verify that remote chains are not yet configured on the token pool on chain A
	supportedChains, err := burnMintTokenPool.GetSupportedChains(nil)
	require.NoError(t, err)
	require.Equal(t, []uint64{}, supportedChains)

	// Verify that the pool is set on the chain A registry
	tokenConfigOnRegistry, err := state.Chains[selectorA].TokenAdminRegistry.GetTokenConfig(nil, tokens[selectorA].Address)
	require.NoError(t, err)
	require.Equal(t, timelockAddress, tokenConfigOnRegistry.Administrator)
	require.Equal(t, burnMintTokenPool.Address(), tokenConfigOnRegistry.TokenPool)

	// Deploy the second token pool on chain B, configuring the token pool on chain A to point at it
	e, err = commonchangeset.ApplyChangesets(t, e, timelockContracts, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(DeployTokenPoolContracts),
			Config: DeployTokenPoolContractsConfig{
				Symbol:        testTokenSymbol,
				TimelockDelay: 0 * time.Second,
				ExistingPoolUpdates: map[uint64]BaseTokenPoolInput{
					selectorA: {
						TokenAddress: tokens[selectorA].Address,
						RemoteChainsToAdd: RemoteChains{
							selectorB: createSymmetricRateLimits(100, 1000),
						},
					},
				},
				NewPools: map[uint64]NewTokenPoolInput{
					selectorB: {
						BaseTokenPoolInput: BaseTokenPoolInput{
							TokenAddress: tokens[selectorB].Address,
							RemoteChainsToAdd: RemoteChains{
								selectorA: createSymmetricRateLimits(100, 1000),
							},
						},
						Type:               BurnMintTokenPool,
						LocalTokenDecimals: localTokenDecimals,
					},
				},
			},
		},
	})
	require.NoError(t, err)

	state, err = LoadOnchainState(e)
	require.NoError(t, err)

	for _, selector := range []uint64{selectorA, selectorB} {
		validateMemberOfBurnMintPair(t, state, tokens, testTokenSymbol, selector, 100, 1000, nil)
	}
}

func TestDeployTokenPoolContracts_DeployNewWithTransferToExternalAdmin(t *testing.T) {
	t.Parallel()
	e, selectorA, selectorB, tokens, timelockContracts := setup2ChainEnvironment(t)

	e, err := commonchangeset.ApplyChangesets(t, e, timelockContracts, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(DeployTokenPoolContracts),
			Config: DeployTokenPoolContractsConfig{
				Symbol:        testTokenSymbol,
				TimelockDelay: 0 * time.Second,
				NewPools: map[uint64]NewTokenPoolInput{
					selectorA: {
						BaseTokenPoolInput: BaseTokenPoolInput{
							TokenAddress: tokens[selectorA].Address,
							RemoteChainsToAdd: RemoteChains{
								selectorB: createSymmetricRateLimits(100, 1000),
							},
							ExternalAdmin: e.Chains[selectorA].DeployerKey.From,
						},
						Type:               BurnMintTokenPool,
						LocalTokenDecimals: localTokenDecimals,
					},
					selectorB: {
						BaseTokenPoolInput: BaseTokenPoolInput{
							TokenAddress: tokens[selectorB].Address,
							RemoteChainsToAdd: RemoteChains{
								selectorA: createSymmetricRateLimits(100, 1000),
							},
							ExternalAdmin: e.Chains[selectorB].DeployerKey.From,
						},
						Type:               BurnMintTokenPool,
						LocalTokenDecimals: localTokenDecimals,
					},
				},
			},
		},
	})
	require.NoError(t, err)

	state, err := LoadOnchainState(e)
	require.NoError(t, err)

	for _, selector := range []uint64{selectorA, selectorB} {
		validateMemberOfBurnMintPair(t, state, tokens, testTokenSymbol, selector, 100, 1000, e.Chains[selector].DeployerKey)
	}
}

func TestDeployTokenPoolContracts_KeepExistingWithTransferToExternalAdmin(t *testing.T) {
	t.Parallel()
	e, selectorA, selectorB, tokens, timelockContracts := setup2ChainEnvironment(t)

	// Initial deployment
	e, err := commonchangeset.ApplyChangesets(t, e, timelockContracts, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(DeployTokenPoolContracts),
			Config: DeployTokenPoolContractsConfig{
				Symbol:        testTokenSymbol,
				TimelockDelay: 0 * time.Second,
				NewPools: map[uint64]NewTokenPoolInput{
					selectorA: {
						BaseTokenPoolInput: BaseTokenPoolInput{
							TokenAddress: tokens[selectorA].Address,
							RemoteChainsToAdd: RemoteChains{
								selectorB: createSymmetricRateLimits(100, 1000),
							},
						},
						Type:               BurnMintTokenPool,
						LocalTokenDecimals: localTokenDecimals,
					},
					selectorB: {
						BaseTokenPoolInput: BaseTokenPoolInput{
							TokenAddress: tokens[selectorB].Address,
							RemoteChainsToAdd: RemoteChains{
								selectorA: createSymmetricRateLimits(100, 1000),
							},
						},
						Type:               BurnMintTokenPool,
						LocalTokenDecimals: localTokenDecimals,
					},
				},
			},
		},
	})
	require.NoError(t, err)

	state, err := LoadOnchainState(e)
	require.NoError(t, err)

	for _, selector := range []uint64{selectorA, selectorB} {
		validateMemberOfBurnMintPair(t, state, tokens, testTokenSymbol, selector, 100, 1000, nil)
	}

	// Transfer existing pools to an external administrator
	e, err = commonchangeset.ApplyChangesets(t, e, timelockContracts, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(DeployTokenPoolContracts),
			Config: DeployTokenPoolContractsConfig{
				Symbol:        testTokenSymbol,
				TimelockDelay: 0 * time.Second,
				ExistingPoolUpdates: map[uint64]BaseTokenPoolInput{
					selectorA: {
						TokenAddress:  tokens[selectorA].Address,
						ExternalAdmin: e.Chains[selectorA].DeployerKey.From,
					},
					selectorB: {
						TokenAddress:  tokens[selectorB].Address,
						ExternalAdmin: e.Chains[selectorB].DeployerKey.From,
					},
				},
			},
		},
	})
	require.NoError(t, err)

	state, err = LoadOnchainState(e)
	require.NoError(t, err)

	for _, selector := range []uint64{selectorA, selectorB} {
		validateMemberOfBurnMintPair(t, state, tokens, testTokenSymbol, selector, 100, 1000, e.Chains[selector].DeployerKey)
	}
}

func TestDeployTokenPoolContracts_RedeployNew(t *testing.T) {
	/*
		TODO: This use case is not yet supported.

		The LoadOnchainState function will potentially load the old token pool into the state mapping,
		as token pools are currently keyed by TokenSymbol. Pools of different versions will be tied to the same symbol.

		Two options to resolve...
		- Ability to remove a current address from the address book through a changeset.
		- Key pools by address instead of token symbol, same as other contracts
	*/
}

// RemoteChains.Validate: rate and capacity are non-zero when isEnabled is set to false
// RemoteChains.Validate: rate is greater than capacity
// RemoteChains.Validate: rate is 0

// NewTokenPoolInput.Validate: invalid pool type
// NewTokenPoolInput.Validate: accept liquidity must be defined for lock release type
// NewTokenPoolInput.Validate: accept liquidity must be nil for burn mint types

// DeployTokenPoolContractsConfig.Validate: chain selector is invalid
// DeployTokenPoolContractsConfig.Validate: chain selector maps have overlap

// deployTokenPool: deploy the 4 pool types successfully
