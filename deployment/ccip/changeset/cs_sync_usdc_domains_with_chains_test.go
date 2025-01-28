package changeset_test

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func deployChainWithUSDCTokenPool(
	t *testing.T,
	lggr logger.Logger,
	e deployment.Environment,
	selector uint64,
	transferToTimelock bool,
) deployment.Environment {
	newAddresses := deployment.NewMemoryAddressBook()
	prereqCfg := []changeset.DeployPrerequisiteConfigPerChain{
		changeset.DeployPrerequisiteConfigPerChain{
			ChainSelector: selector,
		},
	}
	mcmsCfg := map[uint64]commontypes.MCMSWithTimelockConfig{
		selector: proposalutils.SingleGroupTimelockConfig(t),
	}

	e, err := commoncs.ApplyChangesets(t, e, nil, []commoncs.ChangesetApplication{
		{
			Changeset: commoncs.WrapChangeSet(changeset.DeployPrerequisitesChangeset),
			Config: changeset.DeployPrerequisiteConfig{
				Configs: prereqCfg,
			},
		},
		{
			Changeset: commoncs.WrapChangeSet(commoncs.DeployMCMSWithTimelock),
			Config:    mcmsCfg,
		},
	})
	require.NoError(t, err)

	usdcToken, tokenMessenger := testhelpers.DeployUSDCPrerequisites(t, lggr, e.Chains[selector], newAddresses)
	e, err = commonchangeset.ApplyChangesets(t, e, nil, []commonchangeset.ChangesetApplication{
		commonchangeset.ChangesetApplication{
			Changeset: commonchangeset.WrapChangeSet(changeset.DeployUSDCTokenPoolContractsChangeset),
			Config: changeset.DeployUSDCTokenPoolContractsConfig{
				NewUSDCPools: map[uint64]changeset.DeployUSDCTokenPoolInput{
					selector: changeset.DeployUSDCTokenPoolInput{
						TokenMessenger: tokenMessenger.Address,
						TokenAddress:   usdcToken.Address,
					},
				},
			},
		},
	})
	require.NoError(t, err)

	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)

	if transferToTimelock {
		// We only need the token pool owned by timelock in these tests
		timelockOwnedContractsByChain := make(map[uint64][]common.Address, 1)
		timelockOwnedContractsByChain[selector] = []common.Address{state.Chains[selector].USDCTokenPools[changeset.CurrentTokenPoolVersion].Address()}

		// Assemble map of addresses required for Timelock scheduling & execution
		timelockContracts := make(map[uint64]*proposalutils.TimelockExecutionContracts, 1)
		timelockContracts[selector] = &proposalutils.TimelockExecutionContracts{
			Timelock:  state.Chains[selector].Timelock,
			CallProxy: state.Chains[selector].CallProxy,
		}

		// Transfer ownership of token pools to timelock
		e, err = commoncs.ApplyChangesets(t, e, timelockContracts, []commoncs.ChangesetApplication{
			{
				Changeset: commoncs.WrapChangeSet(commoncs.TransferToMCMSWithTimelock),
				Config: commoncs.TransferToMCMSWithTimelockConfig{
					ContractsByChain: timelockOwnedContractsByChain,
					MinDelay:         0,
				},
			},
		})
		require.NoError(t, err)
	}

	return e
}

func TestValidateSyncUSDCDomainsWithChainsConfig(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 2,
	})
	selector := e.AllChainSelectors()[0]

	tests := []struct {
		Msg        string
		Input      changeset.SyncUSDCDomainsWithChainsConfig
		ErrStr     string
		DeployUSDC bool
	}{
		{
			Msg: "Chain selector is not valid",
			Input: changeset.SyncUSDCDomainsWithChainsConfig{
				USDCConfigsByChain: map[uint64]changeset.USDCChainConfig{
					0: changeset.USDCChainConfig{},
				},
			},
			ErrStr: "failed to validate chain selector 0",
		},
		{
			Msg: "Chain selector doesn't exist in environment",
			Input: changeset.SyncUSDCDomainsWithChainsConfig{
				USDCConfigsByChain: map[uint64]changeset.USDCChainConfig{
					5009297550715157269: changeset.USDCChainConfig{},
				},
			},
			ErrStr: "does not exist in environment",
		},
		{
			Msg: "Missing USDC in state",
			Input: changeset.SyncUSDCDomainsWithChainsConfig{
				USDCConfigsByChain: map[uint64]changeset.USDCChainConfig{
					selector: changeset.USDCChainConfig{},
				},
			},
			ErrStr: "does not define any USDC token pools, config should be removed",
		},
		{
			Msg: "Missing USDC in input",
			Input: changeset.SyncUSDCDomainsWithChainsConfig{
				USDCConfigsByChain: map[uint64]changeset.USDCChainConfig{},
			},
			DeployUSDC: true,
			ErrStr:     "which does support USDC",
		},
	}

	for _, test := range tests {
		t.Run(test.Msg, func(t *testing.T) {
			if test.DeployUSDC {
				e = deployChainWithUSDCTokenPool(t, lggr, e, selector, false)
			}

			err := test.Input.Validate(e)
			require.Contains(t, err.Error(), test.ErrStr)
		})
	}
}

func TestValidateUSDCChainConfig(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 2,
	})
	selectors := e.AllChainSelectors()

	tests := []struct {
		Msg        string
		Input      changeset.USDCChainConfig
		ErrStr     string
		DeployUSDC bool
		UseMCMS    bool
	}{
		{
			Msg: "No USDC token pool found with version",
			Input: changeset.USDCChainConfig{
				Version: changeset.CurrentTokenPoolVersion,
			},
			ErrStr: "no USDC token pool found",
		},
		{
			Msg: "Not owned by expected owner",
			Input: changeset.USDCChainConfig{
				Version: changeset.CurrentTokenPoolVersion,
			},
			ErrStr:     "failed ownership validation",
			DeployUSDC: true,
			UseMCMS:    true,
		},
		{
			Msg: "No domain ID found for selector",
			Input: changeset.USDCChainConfig{
				Version: changeset.CurrentTokenPoolVersion,
			},
			ErrStr: "no USDC domain ID defined for chain with selector",
		},
	}

	for _, test := range tests {
		t.Run(test.Msg, func(t *testing.T) {
			if test.DeployUSDC {
				for _, selector := range selectors {
					e = deployChainWithUSDCTokenPool(t, lggr, e, selector, false)
				}
				var err error
				e, err = commoncs.ApplyChangesets(t, e, nil, []commoncs.ChangesetApplication{
					{
						Changeset: commoncs.WrapChangeSet(changeset.ConfigureTokenPoolContractsChangeset),
						Config: changeset.ConfigureTokenPoolContractsConfig{
							PoolUpdates: map[uint64]changeset.TokenPoolConfig{
								selectors[0]: changeset.TokenPoolConfig{
									ChainUpdates: changeset.RateLimiterPerChain{
										selectors[1]: testhelpers.CreateSymmetricRateLimits(0, 0),
									},
									Type:    changeset.USDCTokenPool,
									Version: changeset.CurrentTokenPoolVersion,
								},
								selectors[1]: changeset.TokenPoolConfig{
									ChainUpdates: changeset.RateLimiterPerChain{
										selectors[0]: testhelpers.CreateSymmetricRateLimits(0, 0),
									},
									Type:    changeset.USDCTokenPool,
									Version: changeset.CurrentTokenPoolVersion,
								},
							},
							TokenSymbol: "USDC",
						},
					},
				})
				require.NoError(t, err)
			}

			state, err := changeset.LoadOnchainState(e)
			require.NoError(t, err)

			err = test.Input.Validate(e.GetContext(), e.Chains[selectors[0]], state.Chains[selectors[0]], test.UseMCMS, map[uint64]uint32{})
			require.Contains(t, err.Error(), test.ErrStr)
		})
	}
}

func TestSyncUSDCDomainsWithChainsChangeset(t *testing.T) {
	t.Parallel()

	for _, mcmsConfig := range []*changeset.MCMSConfig{nil, &changeset.MCMSConfig{MinDelay: 0 * time.Second}} {
		msg := "Sync domains without MCMS"
		if mcmsConfig != nil {
			msg = "Sync domains with MCMS"
		}

		t.Run(msg, func(t *testing.T) {
			lggr := logger.TestLogger(t)
			e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
				Chains: 2,
			})
			selectors := e.AllChainSelectors()

			timelockContracts := make(map[uint64]*proposalutils.TimelockExecutionContracts, len(selectors))
			for _, selector := range selectors {
				e = deployChainWithUSDCTokenPool(t, lggr, e, selector, mcmsConfig != nil)
				state, err := changeset.LoadOnchainState(e)
				require.NoError(t, err)
				timelockContracts[selector] = &proposalutils.TimelockExecutionContracts{
					Timelock:  state.Chains[selector].Timelock,
					CallProxy: state.Chains[selector].CallProxy,
				}
			}

			e, err := commoncs.ApplyChangesets(t, e, timelockContracts, []commoncs.ChangesetApplication{
				{
					Changeset: commoncs.WrapChangeSet(changeset.ConfigureTokenPoolContractsChangeset),
					Config: changeset.ConfigureTokenPoolContractsConfig{
						MCMS: mcmsConfig,
						PoolUpdates: map[uint64]changeset.TokenPoolConfig{
							selectors[0]: changeset.TokenPoolConfig{
								ChainUpdates: changeset.RateLimiterPerChain{
									selectors[1]: testhelpers.CreateSymmetricRateLimits(0, 0),
								},
								Type:    changeset.USDCTokenPool,
								Version: changeset.CurrentTokenPoolVersion,
							},
							selectors[1]: changeset.TokenPoolConfig{
								ChainUpdates: changeset.RateLimiterPerChain{
									selectors[0]: testhelpers.CreateSymmetricRateLimits(0, 0),
								},
								Type:    changeset.USDCTokenPool,
								Version: changeset.CurrentTokenPoolVersion,
							},
						},
						TokenSymbol: "USDC",
					},
				},
			})
			require.NoError(t, err)

			e, err = commoncs.ApplyChangesets(t, e, timelockContracts, []commoncs.ChangesetApplication{
				{
					Changeset: commoncs.WrapChangeSet(changeset.SyncUSDCDomainsWithChainsChangeset),
					Config: changeset.SyncUSDCDomainsWithChainsConfig{
						MCMS: mcmsConfig,
						USDCConfigsByChain: map[uint64]changeset.USDCChainConfig{
							selectors[0]: {
								Version: changeset.CurrentTokenPoolVersion,
							},
							selectors[1]: {
								Version: changeset.CurrentTokenPoolVersion,
							},
						},
						ChainSelectorToUSDCDomain: map[uint64]uint32{
							selectors[0]: 1,
							selectors[1]: 2,
						},
					},
				},
			})
			require.NoError(t, err)

			state, err := changeset.LoadOnchainState(e)
			require.NoError(t, err)

			for i, selector := range selectors {
				remoteSelector := selectors[0]
				if i == 0 {
					remoteSelector = selectors[1]
				}
				remoteDomain := uint32(1)
				if i == 0 {
					remoteDomain = 2
				}
				usdcTokenPool := state.Chains[selector].USDCTokenPools[changeset.CurrentTokenPoolVersion]
				remoteUsdcTokenPool := state.Chains[remoteSelector].USDCTokenPools[changeset.CurrentTokenPoolVersion]
				domain, err := usdcTokenPool.GetDomain(nil, remoteSelector)
				allowedCaller := make([]byte, 32)
				bytesCopied := copy(allowedCaller, domain.AllowedCaller[:])
				require.Equal(t, 32, bytesCopied)
				require.NoError(t, err)
				require.True(t, domain.Enabled)
				require.Equal(t, remoteDomain, domain.DomainIdentifier)
				require.Equal(t, remoteUsdcTokenPool.Address(), common.BytesToAddress(allowedCaller))
			}

			// Idempotency check
			output, err := changeset.SyncUSDCDomainsWithChainsChangeset(e, changeset.SyncUSDCDomainsWithChainsConfig{
				MCMS: mcmsConfig,
				USDCConfigsByChain: map[uint64]changeset.USDCChainConfig{
					selectors[0]: {
						Version: changeset.CurrentTokenPoolVersion,
					},
					selectors[1]: {
						Version: changeset.CurrentTokenPoolVersion,
					},
				},
				ChainSelectorToUSDCDomain: map[uint64]uint32{
					selectors[0]: 1,
					selectors[1]: 2,
				},
			})
			require.NoError(t, err)
			require.Empty(t, output.Proposals)
		})
	}
}
