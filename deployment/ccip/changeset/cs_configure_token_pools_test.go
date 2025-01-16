package changeset

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

// createSymmetricRateLimits is a utility to quickly create a rate limiter config with equal inbound and outbound values.
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
) {
	timelockAddress := state.Chains[chainSelector].Timelock.Address()
	tokenPools, ok := state.Chains[chainSelector].BurnMintTokenPools[testTokenSymbol]
	require.True(t, ok)

	tokenPool := tokenPools[0]

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
		require.True(t, inboundConfig.IsEnabled)
		require.Equal(t, big.NewInt(capacity), inboundConfig.Capacity)
		require.Equal(t, big.NewInt(rate), inboundConfig.Rate)

		outboundConfig, err := tokenPool.GetCurrentOutboundRateLimiterState(nil, supportedChain)
		require.NoError(t, err)
		require.True(t, outboundConfig.IsEnabled)
		require.Equal(t, big.NewInt(capacity), outboundConfig.Capacity)
		require.Equal(t, big.NewInt(rate), outboundConfig.Rate)

		remoteTokenAddress, err := tokenPool.GetRemoteToken(nil, supportedChain)
		require.NoError(t, err)
		require.Equal(t, tokens[supportedChain].Address.Bytes(), remoteTokenAddress)

		remoteBurnMintPools, ok := state.Chains[supportedChain].BurnMintTokenPools[testTokenSymbol]
		require.True(t, ok)
		remoteBurnMintPool := remoteBurnMintPools[0]
		remotePoolAddresses, err := tokenPool.GetRemotePools(nil, supportedChain)
		require.NoError(t, err)
		require.Equal(t, [][]byte{remoteBurnMintPool.Address().Bytes()}, remotePoolAddresses)
	}
}

func TestValidateRemoteChains(t *testing.T) {
	t.Parallel()

	tests := []struct {
		IsEnabled bool
		Rate      *big.Int
		Capacity  *big.Int
		ErrStr    string
	}{
		{
			IsEnabled: false,
			Rate:      big.NewInt(1),
			Capacity:  big.NewInt(10),
			ErrStr:    "rate and capacity must be 0",
		},
		{
			IsEnabled: true,
			Rate:      big.NewInt(0),
			Capacity:  big.NewInt(10),
			ErrStr:    "rate must be greater than 0 and less than capacity",
		},
		{
			IsEnabled: true,
			Rate:      big.NewInt(11),
			Capacity:  big.NewInt(10),
			ErrStr:    "rate must be greater than 0 and less than capacity",
		},
	}

	for _, test := range tests {
		t.Run(test.ErrStr, func(t *testing.T) {
			remoteChains := RemoteChainsConfig{
				1: {
					Inbound: token_pool.RateLimiterConfig{
						IsEnabled: test.IsEnabled,
						Rate:      test.Rate,
						Capacity:  test.Capacity,
					},
					Outbound: token_pool.RateLimiterConfig{
						IsEnabled: test.IsEnabled,
						Rate:      test.Rate,
						Capacity:  test.Capacity,
					},
				},
			}

			err := remoteChains.Validate()
			require.Error(t, err)
			require.Contains(t, err.Error(), test.ErrStr)
		})
	}
}

func TestValidateTokenPoolConfig(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 2,
	})

	e, selectorA, _, tokens, _ := setupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), true)

	e = deployTestTokenPools(t, e, map[uint64]DeployTokenPoolInput{
		selectorA: {
			Type:               BurnMintTokenPool,
			TokenAddress:       tokens[selectorA].Address,
			LocalTokenDecimals: 18,
		},
	}, true)

	state, err := LoadOnchainState(e)
	require.NoError(t, err)

	poolAddress := state.Chains[selectorA].BurnMintTokenPools[testTokenSymbol][0].Address()
	invalidPoolAddress := utils.RandomAddress()

	tests := []struct {
		UseMcms         bool
		TokenPoolConfig TokenPoolConfig
		ErrStr          string
		Msg             string
	}{
		{
			Msg:             "Pool address is missing",
			TokenPoolConfig: TokenPoolConfig{},
			ErrStr:          "pool address must be defined",
		},
		{
			Msg: "Pool address is unknown",
			TokenPoolConfig: TokenPoolConfig{
				PoolAddress: invalidPoolAddress,
			},
			ErrStr: fmt.Sprintf("failed to find token pool on %d with symbol %s and address %s", selectorA, testTokenSymbol, invalidPoolAddress),
		},
		{
			Msg: "Pool is not owned by required address",
			TokenPoolConfig: TokenPoolConfig{
				PoolAddress: poolAddress,
			},
			ErrStr: fmt.Sprintf("token pool with address %s on %d failed ownership validation", poolAddress, selectorA),
		},
	}

	for _, test := range tests {
		t.Run(test.Msg, func(t *testing.T) {
			err := test.TokenPoolConfig.Validate(e.GetContext(), e.Chains[selectorA], state.Chains[selectorA], test.UseMcms, testTokenSymbol)
			require.Error(t, err)
			require.ErrorContains(t, err, test.ErrStr)
		})
	}
}

func TestValidateConfigureTokenPoolContractsConfig(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 2,
	})

	tests := []struct {
		TokenSymbol TokenSymbol
		Input       ConfigureTokenPoolContractsConfig
		ErrStr      string
		Msg         string
	}{
		{
			Msg:    "Token symbol is missing",
			Input:  ConfigureTokenPoolContractsConfig{},
			ErrStr: "token symbol must be defined",
		},
		{
			Msg: "Chain selector is invalid",
			Input: ConfigureTokenPoolContractsConfig{
				TokenSymbol: testTokenSymbol,
				PoolUpdates: map[uint64]TokenPoolConfig{
					0: TokenPoolConfig{},
				},
			},
			ErrStr: "failed to validate chain selector 0",
		},
		{
			Msg: "Chain selector doesn't exist in environment",
			Input: ConfigureTokenPoolContractsConfig{
				TokenSymbol: testTokenSymbol,
				PoolUpdates: map[uint64]TokenPoolConfig{
					5009297550715157269: TokenPoolConfig{},
				},
			},
			ErrStr: "chain with selector 5009297550715157269 does not exist in environment",
		},
		{
			Msg: "Corresponding pool update missing",
			Input: ConfigureTokenPoolContractsConfig{
				TokenSymbol: testTokenSymbol,
				PoolUpdates: map[uint64]TokenPoolConfig{
					e.AllChainSelectors()[0]: TokenPoolConfig{
						ChainUpdates: RemoteChainsConfig{
							e.AllChainSelectors()[1]: RateLimiterConfig{},
						},
					},
				},
			},
			ErrStr: fmt.Sprintf("%d is expecting a pool update to be defined for chain with selector %d", e.AllChainSelectors()[0], e.AllChainSelectors()[1]),
		},
		{
			Msg: "Corresponding pool update missing a chain update",
			Input: ConfigureTokenPoolContractsConfig{
				TokenSymbol: testTokenSymbol,
				PoolUpdates: map[uint64]TokenPoolConfig{
					e.AllChainSelectors()[0]: TokenPoolConfig{
						ChainUpdates: RemoteChainsConfig{
							e.AllChainSelectors()[1]: RateLimiterConfig{},
						},
					},
					e.AllChainSelectors()[1]: TokenPoolConfig{},
				},
			},
			ErrStr: fmt.Sprintf("%d is expecting pool update on chain with selector %d to define a chain config pointing back to it", e.AllChainSelectors()[0], e.AllChainSelectors()[1]),
		},
		{
			Msg: "Token admin registry is missing",
			Input: ConfigureTokenPoolContractsConfig{
				TokenSymbol: testTokenSymbol,
				PoolUpdates: map[uint64]TokenPoolConfig{
					e.AllChainSelectors()[0]: TokenPoolConfig{
						ChainUpdates: RemoteChainsConfig{
							e.AllChainSelectors()[1]: RateLimiterConfig{},
						},
					},
					e.AllChainSelectors()[1]: TokenPoolConfig{
						ChainUpdates: RemoteChainsConfig{
							e.AllChainSelectors()[0]: RateLimiterConfig{},
						},
					},
				},
			},
			ErrStr: fmt.Sprintf("missing tokenAdminRegistry on %d", e.AllChainSelectors()[0]),
		},
	}

	for _, test := range tests {
		t.Run(test.Msg, func(t *testing.T) {
			err := test.Input.Validate(e)
			require.Contains(t, err.Error(), test.ErrStr)
		})
	}
}

func TestValidateConfigureTokenPoolContracts(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 2,
	})

	e, selectorA, selectorB, tokens, timelockContracts := setupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), true)

	e = deployTestTokenPools(t, e, map[uint64]DeployTokenPoolInput{
		selectorA: {
			Type:               BurnMintTokenPool,
			TokenAddress:       tokens[selectorA].Address,
			LocalTokenDecimals: 18,
		},
		selectorB: {
			Type:               BurnMintTokenPool,
			TokenAddress:       tokens[selectorB].Address,
			LocalTokenDecimals: 18,
		},
	}, true)

	/*
		e = deployTestTokenPools(t, e, map[uint64]DeployTokenPoolInput{
			selectorA: {
				Type:               BurnMintTokenPool,
				TokenAddress:       tokens[selectorA].Address,
				LocalTokenDecimals: 18,
				ForceDeployment:    true,
			},
		}, true)
	*/

	state, err := LoadOnchainState(e)
	require.NoError(t, err)

	activePoolAddressA := state.Chains[selectorA].BurnMintTokenPools[testTokenSymbol][0].Address()
	activePoolAddressB := state.Chains[selectorB].BurnMintTokenPools[testTokenSymbol][0].Address()
	// upcomingPoolAddressA := state.Chains[selectorA].BurnMintTokenPools[testTokenSymbol][1].Address()
	// upcomingPoolAddressB := state.Chains[selectorB].BurnMintTokenPools[testTokenSymbol][1].Address()

	/*
		// Configure the active pools on the registry
		e, err = commonchangeset.ApplyChangesets(t, e, timelockContracts, []commonchangeset.ChangesetApplication{
			{
				Changeset: commonchangeset.WrapChangeSet(ConfigureTokenAdminRegistry),
				Config: ConfigureTokenAdminRegistryConfig{
					TokenSymbol: testTokenSymbol,
					MCMS: &MCMSConfig{
						MinDelay: 0 * time.Second,
					},
					RegistryUpdates: map[uint64]RegistryConfig{
						selectorA: {
							PoolAddress: activePoolAddressA,
						},
						selectorA: {
							PoolAddress: activePoolAddressB,
						},
					},
				},
			},
		})
	*/

	e, err = commonchangeset.ApplyChangesets(t, e, timelockContracts, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(ConfigureTokenPoolContracts),
			Config: ConfigureTokenPoolContractsConfig{
				TokenSymbol: testTokenSymbol,
				MCMS: &MCMSConfig{
					MinDelay: 0 * time.Second,
				},
				PoolUpdates: map[uint64]TokenPoolConfig{
					selectorA: {
						PoolAddress: activePoolAddressA,
						ChainUpdates: RemoteChainsConfig{
							selectorB: createSymmetricRateLimits(100, 1000),
						},
					},
					selectorA: {
						PoolAddress: activePoolAddressB,
						ChainUpdates: RemoteChainsConfig{
							selectorA: createSymmetricRateLimits(100, 1000),
						},
					},
				},
			},
		},
	})

	for _, selector := range e.AllChainSelectors() {
		validateMemberOfBurnMintPair(t, state, tokens, testTokenSymbol, selector, 100, 1000)
	}
}
