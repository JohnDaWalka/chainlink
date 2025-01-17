package changeset_test

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

// createSymmetricRateLimits is a utility to quickly create a rate limiter config with equal inbound and outbound values.
func createSymmetricRateLimits(rate int64, capacity int64) changeset.RateLimiterConfig {
	return changeset.RateLimiterConfig{
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
func validateMemberOfBurnMintPair(
	t *testing.T,
	state changeset.CCIPOnChainState,
	tokenPool *burn_mint_token_pool.BurnMintTokenPool,
	expectedRemotePools []string,
	tokens map[uint64]*deployment.ContractDeploy[*burn_mint_erc677.BurnMintERC677],
	tokenSymbol changeset.TokenSymbol,
	chainSelector uint64,
	rate *big.Int,
	capacity *big.Int,
	expectedOwner common.Address,
) {
	// Verify that the owner is expected
	owner, err := tokenPool.Owner(nil)
	require.NoError(t, err)
	require.Equal(t, expectedOwner, owner)

	// Fetch the supported remote chains
	supportedChains, err := tokenPool.GetSupportedChains(nil)
	require.NoError(t, err)

	// Verify that the rate limits and remote addresses are correct
	for _, supportedChain := range supportedChains {
		inboundConfig, err := tokenPool.GetCurrentInboundRateLimiterState(nil, supportedChain)
		require.NoError(t, err)
		require.True(t, inboundConfig.IsEnabled)
		require.Equal(t, capacity, inboundConfig.Capacity)
		require.Equal(t, rate, inboundConfig.Rate)

		outboundConfig, err := tokenPool.GetCurrentOutboundRateLimiterState(nil, supportedChain)
		require.NoError(t, err)
		require.True(t, outboundConfig.IsEnabled)
		require.Equal(t, capacity, outboundConfig.Capacity)
		require.Equal(t, rate, outboundConfig.Rate)

		remoteTokenAddress, err := tokenPool.GetRemoteToken(nil, supportedChain)
		require.NoError(t, err)
		require.Equal(t, tokens[supportedChain].Address.Bytes(), remoteTokenAddress)

		remotePoolAddresses, err := tokenPool.GetRemotePools(nil, supportedChain)
		require.NoError(t, err)

		remotePoolsStr := make([]string, len(remotePoolAddresses))
		for i, remotePool := range remotePoolAddresses {
			remotePoolsStr[i] = common.HexToAddress(common.Bytes2Hex(remotePool)).String()
		}
		require.ElementsMatch(t, expectedRemotePools, remotePoolsStr)
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
			remoteChains := changeset.RemoteChainsConfig{
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

	e, selectorA, _, tokens, _ := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), true)

	e = testhelpers.DeployTestTokenPools(t, e, map[uint64]changeset.DeployTokenPoolInput{
		selectorA: {
			Type:               changeset.BurnMintTokenPool,
			TokenAddress:       tokens[selectorA].Address,
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
		},
	}, true)

	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)

	poolAddress := state.Chains[selectorA].BurnMintTokenPools[testhelpers.TestTokenSymbol][0].Address()
	invalidPoolAddress := utils.RandomAddress()

	tests := []struct {
		UseMcms         bool
		TokenPoolConfig changeset.TokenPoolConfig
		ErrStr          string
		Msg             string
	}{
		{
			Msg:             "Pool address is missing",
			TokenPoolConfig: changeset.TokenPoolConfig{},
			ErrStr:          "pool address must be defined",
		},
		{
			Msg: "Pool address is unknown",
			TokenPoolConfig: changeset.TokenPoolConfig{
				PoolAddress: invalidPoolAddress,
			},
			ErrStr: fmt.Sprintf("failed to find token pool on %d with symbol %s and address %s", selectorA, testhelpers.TestTokenSymbol, invalidPoolAddress),
		},
		{
			Msg: "Pool is not owned by required address",
			TokenPoolConfig: changeset.TokenPoolConfig{
				PoolAddress: poolAddress,
			},
			ErrStr: fmt.Sprintf("token pool with address %s on %d failed ownership validation", poolAddress, selectorA),
		},
	}

	for _, test := range tests {
		t.Run(test.Msg, func(t *testing.T) {
			err := test.TokenPoolConfig.Validate(e.GetContext(), e.Chains[selectorA], state.Chains[selectorA], test.UseMcms, testhelpers.TestTokenSymbol)
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
		TokenSymbol changeset.TokenSymbol
		Input       changeset.ConfigureTokenPoolContractsConfig
		ErrStr      string
		Msg         string
	}{
		{
			Msg:    "Token symbol is missing",
			Input:  changeset.ConfigureTokenPoolContractsConfig{},
			ErrStr: "token symbol must be defined",
		},
		{
			Msg: "Chain selector is invalid",
			Input: changeset.ConfigureTokenPoolContractsConfig{
				TokenSymbol: testhelpers.TestTokenSymbol,
				PoolUpdates: map[uint64]changeset.TokenPoolConfig{
					0: changeset.TokenPoolConfig{},
				},
			},
			ErrStr: "failed to validate chain selector 0",
		},
		{
			Msg: "Chain selector doesn't exist in environment",
			Input: changeset.ConfigureTokenPoolContractsConfig{
				TokenSymbol: testhelpers.TestTokenSymbol,
				PoolUpdates: map[uint64]changeset.TokenPoolConfig{
					5009297550715157269: changeset.TokenPoolConfig{},
				},
			},
			ErrStr: "chain with selector 5009297550715157269 does not exist in environment",
		},
		{
			Msg: "Corresponding pool update missing",
			Input: changeset.ConfigureTokenPoolContractsConfig{
				TokenSymbol: testhelpers.TestTokenSymbol,
				PoolUpdates: map[uint64]changeset.TokenPoolConfig{
					e.AllChainSelectors()[0]: changeset.TokenPoolConfig{
						ChainUpdates: changeset.RemoteChainsConfig{
							e.AllChainSelectors()[1]: changeset.RateLimiterConfig{},
						},
					},
				},
			},
			ErrStr: fmt.Sprintf("%d is expecting a pool update to be defined for chain with selector %d", e.AllChainSelectors()[0], e.AllChainSelectors()[1]),
		},
		{
			Msg: "Corresponding pool update missing a chain update",
			Input: changeset.ConfigureTokenPoolContractsConfig{
				TokenSymbol: testhelpers.TestTokenSymbol,
				PoolUpdates: map[uint64]changeset.TokenPoolConfig{
					e.AllChainSelectors()[0]: changeset.TokenPoolConfig{
						ChainUpdates: changeset.RemoteChainsConfig{
							e.AllChainSelectors()[1]: changeset.RateLimiterConfig{},
						},
					},
					e.AllChainSelectors()[1]: changeset.TokenPoolConfig{},
				},
			},
			ErrStr: fmt.Sprintf("%d is expecting pool update on chain with selector %d to define a chain config pointing back to it", e.AllChainSelectors()[0], e.AllChainSelectors()[1]),
		},
		{
			Msg: "Token admin registry is missing",
			Input: changeset.ConfigureTokenPoolContractsConfig{
				TokenSymbol: testhelpers.TestTokenSymbol,
				PoolUpdates: map[uint64]changeset.TokenPoolConfig{
					e.AllChainSelectors()[0]: changeset.TokenPoolConfig{
						ChainUpdates: changeset.RemoteChainsConfig{
							e.AllChainSelectors()[1]: changeset.RateLimiterConfig{},
						},
					},
					e.AllChainSelectors()[1]: changeset.TokenPoolConfig{
						ChainUpdates: changeset.RemoteChainsConfig{
							e.AllChainSelectors()[0]: changeset.RateLimiterConfig{},
						},
					},
				},
			},
			ErrStr: "missing tokenAdminRegistry",
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

	type regPass struct {
		SelectorA2B changeset.RateLimiterConfig
		SelectorB2A changeset.RateLimiterConfig
	}

	type updatePass struct {
		PoolIndexA  int
		PoolIndexB  int
		SelectorA2B changeset.RateLimiterConfig
		SelectorB2A changeset.RateLimiterConfig
	}

	tests := []struct {
		Msg              string
		RegistrationPass *regPass
		UpdatePass       *updatePass
	}{
		{
			Msg: "Configure new pools on registry",
			RegistrationPass: &regPass{
				SelectorA2B: createSymmetricRateLimits(100, 1000),
				SelectorB2A: createSymmetricRateLimits(100, 1000),
			},
		},
		{
			Msg: "Configure new pools on registry, update their rate limits",
			RegistrationPass: &regPass{
				SelectorA2B: createSymmetricRateLimits(100, 1000),
				SelectorB2A: createSymmetricRateLimits(100, 1000),
			},
			UpdatePass: &updatePass{
				PoolIndexA:  0,
				PoolIndexB:  0,
				SelectorA2B: createSymmetricRateLimits(200, 2000),
				SelectorB2A: createSymmetricRateLimits(200, 2000),
			},
		},
		{
			Msg: "Configure new pools on registry, update both pools",
			RegistrationPass: &regPass{
				SelectorA2B: createSymmetricRateLimits(100, 1000),
				SelectorB2A: createSymmetricRateLimits(100, 1000),
			},
			UpdatePass: &updatePass{
				PoolIndexA:  1,
				PoolIndexB:  1,
				SelectorA2B: createSymmetricRateLimits(100, 1000),
				SelectorB2A: createSymmetricRateLimits(100, 1000),
			},
		},
		{
			Msg: "Configure new pools on registry, update only one pool",
			RegistrationPass: &regPass{
				SelectorA2B: createSymmetricRateLimits(100, 1000),
				SelectorB2A: createSymmetricRateLimits(100, 1000),
			},
			UpdatePass: &updatePass{
				PoolIndexA:  0,
				PoolIndexB:  1,
				SelectorA2B: createSymmetricRateLimits(200, 2000),
				SelectorB2A: createSymmetricRateLimits(200, 2000),
			},
		},
	}

	for _, test := range tests {
		for _, mcmsConfig := range []*changeset.MCMSConfig{nil, &changeset.MCMSConfig{MinDelay: 0 * time.Second}} { // Run all tests with and without MCMS
			t.Run(test.Msg, func(t *testing.T) {
				e, selectorA, selectorB, tokens, timelockContracts := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), mcmsConfig != nil)

				e = testhelpers.DeployTestTokenPools(t, e, map[uint64]changeset.DeployTokenPoolInput{
					selectorA: {
						Type:               changeset.BurnMintTokenPool,
						TokenAddress:       tokens[selectorA].Address,
						LocalTokenDecimals: testhelpers.LocalTokenDecimals,
					},
					selectorB: {
						Type:               changeset.BurnMintTokenPool,
						TokenAddress:       tokens[selectorB].Address,
						LocalTokenDecimals: testhelpers.LocalTokenDecimals,
					},
				}, mcmsConfig != nil)

				e = testhelpers.DeployTestTokenPools(t, e, map[uint64]changeset.DeployTokenPoolInput{
					selectorA: {
						Type:               changeset.BurnMintTokenPool,
						TokenAddress:       tokens[selectorA].Address,
						LocalTokenDecimals: testhelpers.LocalTokenDecimals,
						ForceDeployment:    true,
					},
					selectorB: {
						Type:               changeset.BurnMintTokenPool,
						TokenAddress:       tokens[selectorB].Address,
						LocalTokenDecimals: testhelpers.LocalTokenDecimals,
						ForceDeployment:    true,
					},
				}, mcmsConfig != nil)

				state, err := changeset.LoadOnchainState(e)
				require.NoError(t, err)

				pools := map[uint64][]*burn_mint_token_pool.BurnMintTokenPool{
					selectorA: []*burn_mint_token_pool.BurnMintTokenPool{
						state.Chains[selectorA].BurnMintTokenPools[testhelpers.TestTokenSymbol][0],
						state.Chains[selectorA].BurnMintTokenPools[testhelpers.TestTokenSymbol][1],
					},
					selectorB: []*burn_mint_token_pool.BurnMintTokenPool{
						state.Chains[selectorB].BurnMintTokenPools[testhelpers.TestTokenSymbol][0],
						state.Chains[selectorB].BurnMintTokenPools[testhelpers.TestTokenSymbol][1],
					},
				}
				expectedOwners := make(map[uint64]common.Address, 2)
				if mcmsConfig != nil {
					expectedOwners[selectorA] = state.Chains[selectorA].Timelock.Address()
					expectedOwners[selectorB] = state.Chains[selectorB].Timelock.Address()
				} else {
					expectedOwners[selectorA] = e.Chains[selectorA].DeployerKey.From
					expectedOwners[selectorB] = e.Chains[selectorB].DeployerKey.From
				}

				if test.RegistrationPass != nil {
					// Configure & set the active pools on the registry
					e, err = commonchangeset.ApplyChangesets(t, e, timelockContracts, []commonchangeset.ChangesetApplication{
						{
							Changeset: commonchangeset.WrapChangeSet(changeset.ConfigureTokenPoolContracts),
							Config: changeset.ConfigureTokenPoolContractsConfig{
								TokenSymbol: testhelpers.TestTokenSymbol,
								MCMS:        mcmsConfig,
								PoolUpdates: map[uint64]changeset.TokenPoolConfig{
									selectorA: {
										PoolAddress: pools[selectorA][0].Address(),
										ChainUpdates: changeset.RemoteChainsConfig{
											selectorB: test.RegistrationPass.SelectorA2B,
										},
									},
									selectorB: {
										PoolAddress: pools[selectorB][0].Address(),
										ChainUpdates: changeset.RemoteChainsConfig{
											selectorA: test.RegistrationPass.SelectorB2A,
										},
									},
								},
							},
						},
						{
							Changeset: commonchangeset.WrapChangeSet(changeset.ConfigureTokenAdminRegistry),
							Config: changeset.ConfigureTokenAdminRegistryConfig{
								TokenSymbol: testhelpers.TestTokenSymbol,
								MCMS:        mcmsConfig,
								RegistryUpdates: map[uint64]changeset.RegistryConfig{
									selectorA: {
										PoolAddress: pools[selectorA][0].Address(),
									},
									selectorB: {
										PoolAddress: pools[selectorB][0].Address(),
									},
								},
							},
						},
					})
					require.NoError(t, err)

					for _, selector := range e.AllChainSelectors() {
						var remoteChainSelector uint64
						var rateLimiterConfig changeset.RateLimiterConfig
						switch selector {
						case selectorA:
							remoteChainSelector = selectorB
							rateLimiterConfig = test.RegistrationPass.SelectorA2B
						case selectorB:
							remoteChainSelector = selectorA
							rateLimiterConfig = test.RegistrationPass.SelectorB2A
						}
						validateMemberOfBurnMintPair(
							t,
							state,
							pools[selector][0],
							[]string{pools[remoteChainSelector][0].Address().Hex()},
							tokens,
							testhelpers.TestTokenSymbol,
							selector,
							rateLimiterConfig.Inbound.Rate, // inbound & outbound are the same in this test
							rateLimiterConfig.Inbound.Capacity,
							expectedOwners[selector],
						)
					}
				}

				if test.UpdatePass != nil {
					// Only configure, do not update registry
					e, err = commonchangeset.ApplyChangesets(t, e, timelockContracts, []commonchangeset.ChangesetApplication{
						{
							Changeset: commonchangeset.WrapChangeSet(changeset.ConfigureTokenPoolContracts),
							Config: changeset.ConfigureTokenPoolContractsConfig{
								TokenSymbol: testhelpers.TestTokenSymbol,
								MCMS:        mcmsConfig,
								PoolUpdates: map[uint64]changeset.TokenPoolConfig{
									selectorA: {
										PoolAddress: pools[selectorA][test.UpdatePass.PoolIndexA].Address(),
										ChainUpdates: changeset.RemoteChainsConfig{
											selectorB: test.UpdatePass.SelectorA2B,
										},
									},
									selectorB: {
										PoolAddress: pools[selectorB][test.UpdatePass.PoolIndexB].Address(),
										ChainUpdates: changeset.RemoteChainsConfig{
											selectorA: test.UpdatePass.SelectorB2A,
										},
									},
								},
							},
						},
					})
					require.NoError(t, err)

					for _, selector := range e.AllChainSelectors() {
						var poolIndex int
						var remotePoolIndex int
						var remoteChainSelector uint64
						var rateLimiterConfig changeset.RateLimiterConfig
						switch selector {
						case selectorA:
							remoteChainSelector = selectorB
							rateLimiterConfig = test.UpdatePass.SelectorA2B
							poolIndex = test.UpdatePass.PoolIndexA
							remotePoolIndex = test.UpdatePass.PoolIndexB
						case selectorB:
							remoteChainSelector = selectorA
							rateLimiterConfig = test.UpdatePass.SelectorB2A
							poolIndex = test.UpdatePass.PoolIndexB
							remotePoolIndex = test.UpdatePass.PoolIndexA
						}
						remotePoolAddresses := []string{pools[remoteChainSelector][0].Address().String()} // add registered pool by default
						if remotePoolIndex == 1 {                                                         // if remote pool address is being updated, we push the new address
							remotePoolAddresses = append(remotePoolAddresses, pools[remoteChainSelector][1].Address().String())
						}
						validateMemberOfBurnMintPair(
							t,
							state,
							pools[selector][poolIndex],
							remotePoolAddresses,
							tokens,
							testhelpers.TestTokenSymbol,
							selector,
							rateLimiterConfig.Inbound.Rate, // inbound & outbound are the same in this test
							rateLimiterConfig.Inbound.Capacity,
							expectedOwners[selector],
						)
					}
				}
			})
		}
	}
}
