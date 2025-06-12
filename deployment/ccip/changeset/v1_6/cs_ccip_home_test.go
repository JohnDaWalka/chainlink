package v1_6_test

import (
	"math/big"
	"regexp"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/maps"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	"github.com/smartcontractkit/chainlink-ccip/chainconfig"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	ccipops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_6"
	ccipseq "github.com/smartcontractkit/chainlink/deployment/ccip/sequence/evm/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

func TestInvalidOCR3Params(t *testing.T) {
	e, _ := testhelpers.NewMemoryEnvironment(t,
		testhelpers.WithPrerequisiteDeploymentOnly(nil))
	chain1 := e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]
	envNodes, err := deployment.NodeInfo(e.Env.NodeIDs, e.Env.Offchain)
	require.NoError(t, err)
	// Need to deploy prerequisites first so that we can form the USDC config
	// no proposals to be made, timelock can be passed as nil here
	e.Env, err = commonchangeset.Apply(t, e.Env, commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(v1_6.DeployHomeChainChangeset),
		v1_6.DeployHomeChainConfig{
			HomeChainSel:     e.HomeChainSel,
			RMNDynamicConfig: testhelpers.NewTestRMNDynamicConfig(),
			RMNStaticConfig:  testhelpers.NewTestRMNStaticConfig(),
			NodeOperators:    testhelpers.NewTestNodeOperator(e.Env.BlockChains.EVMChains()[e.HomeChainSel].DeployerKey.From),
			NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
				testhelpers.TestNodeOperator: envNodes.NonBootstraps().PeerIDs(),
			},
		},
	), commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(v1_6.DeployChainContractsChangeset),
		ccipseq.DeployChainContractsConfig{
			HomeChainSelector: e.HomeChainSel,
			ContractParamsPerChain: map[uint64]ccipseq.ChainContractParams{
				chain1: {
					FeeQuoterParams: ccipops.DefaultFeeQuoterParams(),
					OffRampParams:   ccipops.DefaultOffRampParams(),
				},
			},
		},
	))
	require.NoError(t, err)

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)
	nodes, err := deployment.NodeInfo(e.Env.NodeIDs, e.Env.Offchain)
	require.NoError(t, err)
	params := v1_6.DeriveOCRParamsForCommit(v1_6.SimulationTest, e.FeedChainSel, nil, nil)

	// tweak params to have invalid config
	// make DeltaRound greater than DeltaProgress
	params.OCRParameters.DeltaRound = params.OCRParameters.DeltaProgress + time.Duration(1)
	_, err = internal.BuildOCR3ConfigForCCIPHome(
		state.Chains[e.HomeChainSel].CCIPHome,
		e.Env.OCRSecrets,
		state.Chains[chain1].OffRamp.Address().Bytes(),
		chain1,
		nodes.NonBootstraps(),
		state.Chains[e.HomeChainSel].RMNHome.Address(),
		params.OCRParameters,
		params.CommitOffChainConfig,
		&globals.DefaultExecuteOffChainCfg,
		false,
	)
	require.Errorf(t, err, "expected error")
	pattern := `DeltaRound \(\d+\.\d+s\) must be less than DeltaProgress \(\d+s\)`
	matched, err1 := regexp.MatchString(pattern, err.Error())
	require.NoError(t, err1)
	require.True(t, matched)
}

func Test_PromoteCandidate(t *testing.T) {
	for _, tc := range []struct {
		name        string
		mcmsEnabled bool
	}{
		{
			name:        "MCMS enabled",
			mcmsEnabled: true,
		},
		{
			name:        "MCMS disabled",
			mcmsEnabled: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := testcontext.Get(t)
			tenv, _ := testhelpers.NewMemoryEnvironment(t,
				testhelpers.WithNumOfChains(2),
				testhelpers.WithNumOfNodes(4))
			state, err := stateview.LoadOnchainState(tenv.Env)
			require.NoError(t, err)

			// Deploy to all chains.
			allChains := maps.Keys(tenv.Env.BlockChains.EVMChains())
			source := allChains[0]
			dest := allChains[1]

			if tc.mcmsEnabled {
				// Transfer ownership to timelock so that we can promote the zero digest later down the line.
				testhelpers.TransferToTimelock(t, tenv, state, []uint64{source, dest}, true)
			}

			var (
				capReg   = state.Chains[tenv.HomeChainSel].CapabilityRegistry
				ccipHome = state.Chains[tenv.HomeChainSel].CCIPHome
			)
			donID, err := internal.DonIDForChain(capReg, ccipHome, dest)
			require.NoError(t, err)
			require.NotEqual(t, uint32(0), donID)
			t.Logf("donID: %d", donID)
			candidateDigestCommitBefore, err := ccipHome.GetCandidateDigest(&bind.CallOpts{
				Context: ctx,
			}, donID, uint8(types.PluginTypeCCIPCommit))
			require.NoError(t, err)
			require.Equal(t, [32]byte{}, candidateDigestCommitBefore)
			ActiveDigestExecBefore, err := ccipHome.GetActiveDigest(&bind.CallOpts{
				Context: ctx,
			}, donID, uint8(types.PluginTypeCCIPExec))
			require.NoError(t, err)
			require.NotEqual(t, [32]byte{}, ActiveDigestExecBefore)

			var mcmsConfig *proposalutils.TimelockConfig
			if tc.mcmsEnabled {
				mcmsConfig = &proposalutils.TimelockConfig{
					MinDelay: 0,
				}
			}
			// promotes zero digest on commit and ensure exec is not affected
			_, err = commonchangeset.Apply(t, tenv.Env,
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(v1_6.PromoteCandidateChangeset),
					v1_6.PromoteCandidateChangesetConfig{
						HomeChainSelector: tenv.HomeChainSel,
						PluginInfo: []v1_6.PromoteCandidatePluginInfo{
							{
								RemoteChainSelectors:    []uint64{dest},
								PluginType:              types.PluginTypeCCIPCommit,
								AllowEmptyConfigPromote: true,
							},
						},
						MCMS: mcmsConfig,
					},
				),
			)
			require.NoError(t, err)

			// after promoting the zero digest, active digest should also be zero
			activeDigestCommit, err := ccipHome.GetActiveDigest(&bind.CallOpts{
				Context: ctx,
			}, donID, uint8(types.PluginTypeCCIPCommit))
			require.NoError(t, err)
			require.Equal(t, [32]byte{}, activeDigestCommit)

			activeDigestExec, err := ccipHome.GetActiveDigest(&bind.CallOpts{
				Context: ctx,
			}, donID, uint8(types.PluginTypeCCIPExec))
			require.NoError(t, err)
			require.Equal(t, ActiveDigestExecBefore, activeDigestExec)
		})
	}
}

func Test_SetCandidate(t *testing.T) {
	for _, tc := range []struct {
		name        string
		mcmsEnabled bool
	}{
		{
			name:        "MCMS enabled",
			mcmsEnabled: true,
		},
		{
			name:        "MCMS disabled",
			mcmsEnabled: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := testcontext.Get(t)
			tenv, _ := testhelpers.NewMemoryEnvironment(t,
				testhelpers.WithNumOfChains(2),
				testhelpers.WithNumOfNodes(4))
			state, err := stateview.LoadOnchainState(tenv.Env)
			require.NoError(t, err)

			// Deploy to all chains.
			allChains := maps.Keys(tenv.Env.BlockChains.EVMChains())
			source := allChains[0]
			dest := allChains[1]

			if tc.mcmsEnabled {
				// Transfer ownership to timelock so that we can promote the zero digest later down the line.
				testhelpers.TransferToTimelock(t, tenv, state, []uint64{source, dest}, true)
			}

			var (
				capReg   = state.Chains[tenv.HomeChainSel].CapabilityRegistry
				ccipHome = state.Chains[tenv.HomeChainSel].CCIPHome
			)
			donID, err := internal.DonIDForChain(capReg, ccipHome, dest)
			require.NoError(t, err)
			require.NotEqual(t, uint32(0), donID)
			candidateDigestCommitBefore, err := ccipHome.GetCandidateDigest(&bind.CallOpts{
				Context: ctx,
			}, donID, uint8(types.PluginTypeCCIPCommit))
			require.NoError(t, err)
			require.Equal(t, [32]byte{}, candidateDigestCommitBefore)
			candidateDigestExecBefore, err := ccipHome.GetCandidateDigest(&bind.CallOpts{
				Context: ctx,
			}, donID, uint8(types.PluginTypeCCIPExec))
			require.NoError(t, err)
			require.Equal(t, [32]byte{}, candidateDigestExecBefore)

			var mcmsConfig *proposalutils.TimelockConfig
			if tc.mcmsEnabled {
				mcmsConfig = &proposalutils.TimelockConfig{
					MinDelay: 0,
				}
			}

			tokenConfig := shared.NewTestTokenConfig(state.Chains[tenv.FeedChainSel].USDFeeds)

			_, err = commonchangeset.Apply(t, tenv.Env, commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_6.SetCandidateChangeset),
				v1_6.SetCandidateChangesetConfig{
					SetCandidateConfigBase: v1_6.SetCandidateConfigBase{
						HomeChainSelector: tenv.HomeChainSel,
						FeedChainSelector: tenv.FeedChainSel,
						MCMS:              mcmsConfig,
					},
					PluginInfo: []v1_6.SetCandidatePluginInfo{
						{
							OCRConfigPerRemoteChainSelector: map[uint64]v1_6.CCIPOCRParams{
								dest: v1_6.DeriveOCRParamsForCommit(v1_6.SimulationTest, tenv.FeedChainSel, tokenConfig.GetTokenInfo(logger.TestLogger(t),
									state.Chains[dest].LinkToken.Address(),
									state.Chains[dest].Weth9.Address()), nil),
							},
							PluginType: types.PluginTypeCCIPCommit,
						},
						{
							OCRConfigPerRemoteChainSelector: map[uint64]v1_6.CCIPOCRParams{
								dest: v1_6.DeriveOCRParamsForExec(v1_6.SimulationTest, nil, nil),
							},
							PluginType: types.PluginTypeCCIPExec,
						},
					},
				},
			), commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_6.SetCandidateChangeset),
				v1_6.SetCandidateChangesetConfig{
					SetCandidateConfigBase: v1_6.SetCandidateConfigBase{
						HomeChainSelector: tenv.HomeChainSel,
						FeedChainSelector: tenv.FeedChainSel,
						MCMS:              mcmsConfig,
					},
					PluginInfo: []v1_6.SetCandidatePluginInfo{
						{
							OCRConfigPerRemoteChainSelector: map[uint64]v1_6.CCIPOCRParams{
								dest: v1_6.DeriveOCRParamsForCommit(v1_6.SimulationTest, tenv.FeedChainSel, tokenConfig.GetTokenInfo(logger.TestLogger(t),
									state.Chains[dest].LinkToken.Address(),
									state.Chains[dest].Weth9.Address()), nil),
							},
							PluginType: types.PluginTypeCCIPCommit,
						},
						{
							OCRConfigPerRemoteChainSelector: map[uint64]v1_6.CCIPOCRParams{
								dest: v1_6.DeriveOCRParamsForExec(v1_6.SimulationTest, nil, nil),
							},
							PluginType: types.PluginTypeCCIPExec,
						},
					},
				},
			))
			require.NoError(t, err)
			// after setting a new candidate on both plugins, the candidate config digest
			// should be nonzero.
			candidateDigestCommitAfter, err := ccipHome.GetCandidateDigest(&bind.CallOpts{
				Context: ctx,
			}, donID, uint8(types.PluginTypeCCIPCommit))
			require.NoError(t, err)
			require.NotEqual(t, [32]byte{}, candidateDigestCommitAfter)
			require.NotEqual(t, candidateDigestCommitBefore, candidateDigestCommitAfter)

			candidateDigestExecAfter, err := ccipHome.GetCandidateDigest(&bind.CallOpts{
				Context: ctx,
			}, donID, uint8(types.PluginTypeCCIPExec))
			require.NoError(t, err)
			require.NotEqual(t, [32]byte{}, candidateDigestExecAfter)
			require.NotEqual(t, candidateDigestExecBefore, candidateDigestExecAfter)
		})
	}
}

func Test_RevokeCandidate(t *testing.T) {
	for _, tc := range []struct {
		name        string
		mcmsEnabled bool
	}{
		{
			name:        "MCMS enabled",
			mcmsEnabled: true,
		},
		{
			name:        "MCMS disabled",
			mcmsEnabled: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := testcontext.Get(t)
			tenv, _ := testhelpers.NewMemoryEnvironment(t,
				testhelpers.WithNumOfChains(2),
				testhelpers.WithNumOfNodes(4))
			state, err := stateview.LoadOnchainState(tenv.Env)
			require.NoError(t, err)

			// Deploy to all chains.
			allChains := maps.Keys(tenv.Env.BlockChains.EVMChains())
			source := allChains[0]
			dest := allChains[1]

			if tc.mcmsEnabled {
				// Transfer ownership to timelock so that we can promote the zero digest later down the line.
				testhelpers.TransferToTimelock(t, tenv, state, []uint64{source, dest}, true)
			}

			var (
				capReg   = state.Chains[tenv.HomeChainSel].CapabilityRegistry
				ccipHome = state.Chains[tenv.HomeChainSel].CCIPHome
			)
			donID, err := internal.DonIDForChain(capReg, ccipHome, dest)
			require.NoError(t, err)
			require.NotEqual(t, uint32(0), donID)
			candidateDigestCommitBefore, err := ccipHome.GetCandidateDigest(&bind.CallOpts{
				Context: ctx,
			}, donID, uint8(types.PluginTypeCCIPCommit))
			require.NoError(t, err)
			require.Equal(t, [32]byte{}, candidateDigestCommitBefore)
			candidateDigestExecBefore, err := ccipHome.GetCandidateDigest(&bind.CallOpts{
				Context: ctx,
			}, donID, uint8(types.PluginTypeCCIPExec))
			require.NoError(t, err)
			require.Equal(t, [32]byte{}, candidateDigestExecBefore)

			var mcmsConfig *proposalutils.TimelockConfig
			if tc.mcmsEnabled {
				mcmsConfig = &proposalutils.TimelockConfig{
					MinDelay: 0,
				}
			}
			tokenConfig := shared.NewTestTokenConfig(state.Chains[tenv.FeedChainSel].USDFeeds)
			_, err = commonchangeset.Apply(t, tenv.Env,
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(v1_6.SetCandidateChangeset),
					v1_6.SetCandidateChangesetConfig{
						SetCandidateConfigBase: v1_6.SetCandidateConfigBase{
							HomeChainSelector: tenv.HomeChainSel,
							FeedChainSelector: tenv.FeedChainSel,
							MCMS:              mcmsConfig,
						},
						PluginInfo: []v1_6.SetCandidatePluginInfo{
							{
								OCRConfigPerRemoteChainSelector: map[uint64]v1_6.CCIPOCRParams{
									dest: v1_6.DeriveOCRParamsForCommit(v1_6.SimulationTest, tenv.FeedChainSel, tokenConfig.GetTokenInfo(logger.TestLogger(t),
										state.Chains[dest].LinkToken.Address(),
										state.Chains[dest].Weth9.Address()), nil),
								},
								PluginType: types.PluginTypeCCIPCommit,
							},
							{
								OCRConfigPerRemoteChainSelector: map[uint64]v1_6.CCIPOCRParams{
									dest: v1_6.DeriveOCRParamsForExec(v1_6.SimulationTest, nil, nil),
								},
								PluginType: types.PluginTypeCCIPExec,
							},
						},
					},
				),
			)
			require.NoError(t, err)

			// after setting a new candidate on both plugins, the candidate config digest
			// should be nonzero.
			candidateDigestCommitAfter, err := ccipHome.GetCandidateDigest(&bind.CallOpts{
				Context: ctx,
			}, donID, uint8(types.PluginTypeCCIPCommit))
			require.NoError(t, err)
			require.NotEqual(t, [32]byte{}, candidateDigestCommitAfter)
			require.NotEqual(t, candidateDigestCommitBefore, candidateDigestCommitAfter)

			candidateDigestExecAfter, err := ccipHome.GetCandidateDigest(&bind.CallOpts{
				Context: ctx,
			}, donID, uint8(types.PluginTypeCCIPExec))
			require.NoError(t, err)
			require.NotEqual(t, [32]byte{}, candidateDigestExecAfter)
			require.NotEqual(t, candidateDigestExecBefore, candidateDigestExecAfter)

			// next we can revoke candidate - this should set the candidate digest back to zero
			_, err = commonchangeset.Apply(t, tenv.Env, commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_6.RevokeCandidateChangeset),
				v1_6.RevokeCandidateChangesetConfig{
					HomeChainSelector:   tenv.HomeChainSel,
					RemoteChainSelector: dest,
					PluginType:          types.PluginTypeCCIPCommit,
					MCMS:                mcmsConfig,
				},
			), commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_6.RevokeCandidateChangeset),
				v1_6.RevokeCandidateChangesetConfig{
					HomeChainSelector:   tenv.HomeChainSel,
					RemoteChainSelector: dest,
					PluginType:          types.PluginTypeCCIPExec,
					MCMS:                mcmsConfig,
				},
			))
			require.NoError(t, err)

			// after revoking the candidate, the candidate digest should be zero
			candidateDigestCommitAfterRevoke, err := ccipHome.GetCandidateDigest(&bind.CallOpts{
				Context: ctx,
			}, donID, uint8(types.PluginTypeCCIPCommit))
			require.NoError(t, err)
			require.Equal(t, [32]byte{}, candidateDigestCommitAfterRevoke)

			candidateDigestExecAfterRevoke, err := ccipHome.GetCandidateDigest(&bind.CallOpts{
				Context: ctx,
			}, donID, uint8(types.PluginTypeCCIPExec))
			require.NoError(t, err)
			require.Equal(t, [32]byte{}, candidateDigestExecAfterRevoke)
		})
	}
}

func Test_UpdateChainConfigs(t *testing.T) {
	for _, tc := range []struct {
		name        string
		mcmsEnabled bool
	}{
		{
			name:        "MCMS enabled",
			mcmsEnabled: true,
		},
		{
			name:        "MCMS disabled",
			mcmsEnabled: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(3))
			state, err := stateview.LoadOnchainState(tenv.Env)
			require.NoError(t, err)

			allChains := maps.Keys(tenv.Env.BlockChains.EVMChains())
			source := allChains[0]
			dest := allChains[1]
			otherChain := allChains[2]

			if tc.mcmsEnabled {
				// Transfer ownership to timelock so that we can promote the zero digest later down the line.
				testhelpers.TransferToTimelock(t, tenv, state, []uint64{source, dest}, true)
			}

			ccipHome := state.Chains[tenv.HomeChainSel].CCIPHome
			otherChainConfig, err := ccipHome.GetChainConfig(nil, otherChain)
			require.NoError(t, err)
			assert.NotZero(t, otherChainConfig.FChain)

			var mcmsConfig *proposalutils.TimelockConfig
			if tc.mcmsEnabled {
				mcmsConfig = &proposalutils.TimelockConfig{
					MinDelay: 0,
				}
			}
			_, err = commonchangeset.Apply(t, tenv.Env,
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(v1_6.UpdateChainConfigChangeset),
					v1_6.UpdateChainConfigConfig{
						HomeChainSelector:  tenv.HomeChainSel,
						RemoteChainRemoves: []uint64{otherChain},
						RemoteChainAdds:    make(map[uint64]v1_6.ChainConfig),
						MCMS:               mcmsConfig,
					},
				),
			)
			require.NoError(t, err)

			// other chain should be gone
			chainConfigAfter, err := ccipHome.GetChainConfig(nil, otherChain)
			require.NoError(t, err)
			assert.Zero(t, chainConfigAfter.FChain)

			// Lets add it back now.
			_, err = commonchangeset.Apply(t, tenv.Env,
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(v1_6.UpdateChainConfigChangeset),
					v1_6.UpdateChainConfigConfig{
						HomeChainSelector:  tenv.HomeChainSel,
						RemoteChainRemoves: []uint64{},
						RemoteChainAdds: map[uint64]v1_6.ChainConfig{
							otherChain: {
								EncodableChainConfig: chainconfig.ChainConfig{
									GasPriceDeviationPPB:    cciptypes.BigInt{Int: big.NewInt(testhelpers.DefaultGasPriceDeviationPPB)},
									DAGasPriceDeviationPPB:  cciptypes.BigInt{Int: big.NewInt(testhelpers.DefaultDAGasPriceDeviationPPB)},
									OptimisticConfirmations: globals.OptimisticConfirmations,
								},
								FChain:  otherChainConfig.FChain,
								Readers: otherChainConfig.Readers,
							},
						},
						MCMS: mcmsConfig,
					},
				),
			)
			require.NoError(t, err)

			chainConfigAfter2, err := ccipHome.GetChainConfig(nil, otherChain)
			require.NoError(t, err)
			assert.Equal(t, chainConfigAfter2.FChain, otherChainConfig.FChain)
			assert.Equal(t, chainConfigAfter2.Readers, otherChainConfig.Readers)
			assert.Equal(t, chainConfigAfter2.Config, otherChainConfig.Config)
		})
	}
}

func Test_ValidateMultipleReportsVsEnforceOutOfOrder(t *testing.T) {
	// Set up environment with 3 chains (home + 2 others)
	tenv, _ := testhelpers.NewMemoryEnvironment(t,
		testhelpers.WithNumOfChains(3),
		testhelpers.WithNumOfNodes(4))

	state, err := stateview.LoadOnchainState(tenv.Env)
	require.NoError(t, err)

	allChains := maps.Keys(tenv.Env.BlockChains.EVMChains())
	source := allChains[0]
	dest1 := allChains[1]
	dest2 := allChains[2]

	// First, set up lanes between source and both destinations
	// By default, EnforceOutOfOrder is false
	_, err = commonchangeset.Apply(t, tenv.Env,
		commonchangeset.Configure(
			v1_6.UpdateBidirectionalLanesChangeset,
			v1_6.UpdateBidirectionalLanesConfig{
				Lanes: []v1_6.BidirectionalLaneDefinition{
					{
						Chains: [2]v1_6.ChainDefinition{
							{
								Selector:                 source,
								ConnectionConfig:         v1_6.ConnectionConfig{},
								GasPrice:                 testhelpers.DefaultGasPrice,
								FeeQuoterDestChainConfig: v1_6.DefaultFeeQuoterDestChainConfig(true),
							},
							{
								Selector:                 dest1,
								ConnectionConfig:         v1_6.ConnectionConfig{},
								GasPrice:                 testhelpers.DefaultGasPrice,
								FeeQuoterDestChainConfig: v1_6.DefaultFeeQuoterDestChainConfig(true),
							},
						},
					},
					{
						Chains: [2]v1_6.ChainDefinition{
							{
								Selector:                 source,
								ConnectionConfig:         v1_6.ConnectionConfig{},
								GasPrice:                 testhelpers.DefaultGasPrice,
								FeeQuoterDestChainConfig: v1_6.DefaultFeeQuoterDestChainConfig(true),
							},
							{
								Selector:                 dest2,
								ConnectionConfig:         v1_6.ConnectionConfig{},
								GasPrice:                 testhelpers.DefaultGasPrice,
								FeeQuoterDestChainConfig: v1_6.DefaultFeeQuoterDestChainConfig(true),
							},
						},
					},
				},
			},
		),
	)
	require.NoError(t, err)

	// Verify that EnforceOutOfOrder is false by default
	feeQuoter := state.Chains[source].FeeQuoter
	destConfig1, err := feeQuoter.GetDestChainConfig(nil, dest1)
	require.NoError(t, err)
	require.False(t, destConfig1.EnforceOutOfOrder, "EnforceOutOfOrder should be false by default")

	destConfig2, err := feeQuoter.GetDestChainConfig(nil, dest2)
	require.NoError(t, err)
	require.False(t, destConfig2.EnforceOutOfOrder, "EnforceOutOfOrder should be false by default")

	// Test Case 1: MultipleReportsEnabled=true with EnforceOutOfOrder=false should fail
	t.Run("MultipleReportsEnabled with EnforceOutOfOrder false should fail", func(t *testing.T) {
		execParams := v1_6.DeriveOCRParamsForExec(v1_6.SimulationTest, nil, func(params v1_6.CCIPOCRParams) v1_6.CCIPOCRParams {
			params.ExecuteOffChainConfig.MultipleReportsEnabled = true
			return params
		})

		_, err = commonchangeset.Apply(t, tenv.Env,
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_6.SetCandidateChangeset),
				v1_6.SetCandidateChangesetConfig{
					SetCandidateConfigBase: v1_6.SetCandidateConfigBase{
						HomeChainSelector: tenv.HomeChainSel,
						FeedChainSelector: tenv.FeedChainSel,
					},
					PluginInfo: []v1_6.SetCandidatePluginInfo{
						{
							OCRConfigPerRemoteChainSelector: map[uint64]v1_6.CCIPOCRParams{
								source: execParams,
							},
							PluginType: types.PluginTypeCCIPExec,
						},
					},
				},
			),
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "MultipleReportsEnabled is set to true")
		require.Contains(t, err.Error(), "EnforceOutOfOrder=false")
	})

	// Test Case 2: Update one destination to have EnforceOutOfOrder=true, should still fail
	t.Run("MultipleReportsEnabled with mixed EnforceOutOfOrder should fail", func(t *testing.T) {
		// Update only dest1 to have EnforceOutOfOrder=true
		feeQuoterDestConfig := v1_6.DefaultFeeQuoterDestChainConfig(true)
		feeQuoterDestConfig.EnforceOutOfOrder = true

		_, err = commonchangeset.Apply(t, tenv.Env,
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_6.UpdateFeeQuoterDestsChangeset),
				v1_6.UpdateFeeQuoterDestsConfig{
					UpdatesByChain: map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig{
						source: {
							dest1: feeQuoterDestConfig,
						},
					},
				},
			),
		)
		require.NoError(t, err)

		// Try to set MultipleReportsEnabled=true again, should still fail because dest2 has EnforceOutOfOrder=false
		execParams := v1_6.DeriveOCRParamsForExec(v1_6.SimulationTest, nil, func(params v1_6.CCIPOCRParams) v1_6.CCIPOCRParams {
			params.ExecuteOffChainConfig.MultipleReportsEnabled = true
			return params
		})

		_, err = commonchangeset.Apply(t, tenv.Env,
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_6.SetCandidateChangeset),
				v1_6.SetCandidateChangesetConfig{
					SetCandidateConfigBase: v1_6.SetCandidateConfigBase{
						HomeChainSelector: tenv.HomeChainSel,
						FeedChainSelector: tenv.FeedChainSel,
					},
					PluginInfo: []v1_6.SetCandidatePluginInfo{
						{
							OCRConfigPerRemoteChainSelector: map[uint64]v1_6.CCIPOCRParams{
								source: execParams,
							},
							PluginType: types.PluginTypeCCIPExec,
						},
					},
				},
			),
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "destination chain")
		require.Contains(t, err.Error(), "EnforceOutOfOrder=false")
	})

	// Test Case 3: Update all destinations to have EnforceOutOfOrder=true, should succeed
	t.Run("MultipleReportsEnabled with all EnforceOutOfOrder true should succeed", func(t *testing.T) {
		// Update dest2 to also have EnforceOutOfOrder=true
		feeQuoterDestConfig := v1_6.DefaultFeeQuoterDestChainConfig(true)
		feeQuoterDestConfig.EnforceOutOfOrder = true

		_, err = commonchangeset.Apply(t, tenv.Env,
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_6.UpdateFeeQuoterDestsChangeset),
				v1_6.UpdateFeeQuoterDestsConfig{
					UpdatesByChain: map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig{
						source: {
							dest2: feeQuoterDestConfig,
						},
					},
				},
			),
		)
		require.NoError(t, err)

		// Now setting MultipleReportsEnabled=true should succeed
		execParams := v1_6.DeriveOCRParamsForExec(v1_6.SimulationTest, nil, func(params v1_6.CCIPOCRParams) v1_6.CCIPOCRParams {
			params.ExecuteOffChainConfig.MultipleReportsEnabled = true
			return params
		})

		_, err = commonchangeset.Apply(t, tenv.Env,
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_6.SetCandidateChangeset),
				v1_6.SetCandidateChangesetConfig{
					SetCandidateConfigBase: v1_6.SetCandidateConfigBase{
						HomeChainSelector: tenv.HomeChainSel,
						FeedChainSelector: tenv.FeedChainSel,
					},
					PluginInfo: []v1_6.SetCandidatePluginInfo{
						{
							OCRConfigPerRemoteChainSelector: map[uint64]v1_6.CCIPOCRParams{
								source: execParams,
							},
							PluginType: types.PluginTypeCCIPExec,
						},
					},
				},
			),
		)
		require.NoError(t, err, "Should succeed when all destinations have EnforceOutOfOrder=true")
	})

	// Test Case 4: MultipleReportsEnabled=false should always work regardless of EnforceOutOfOrder
	t.Run("MultipleReportsEnabled false should always succeed", func(t *testing.T) {
		// Reset dest2 back to EnforceOutOfOrder=false
		feeQuoterDestConfig := v1_6.DefaultFeeQuoterDestChainConfig(true)
		feeQuoterDestConfig.EnforceOutOfOrder = false

		_, err = commonchangeset.Apply(t, tenv.Env,
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_6.UpdateFeeQuoterDestsChangeset),
				v1_6.UpdateFeeQuoterDestsConfig{
					UpdatesByChain: map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig{
						source: {
							dest2: feeQuoterDestConfig,
						},
					},
				},
			),
		)
		require.NoError(t, err)

		// Setting MultipleReportsEnabled=false should work
		execParams := v1_6.DeriveOCRParamsForExec(v1_6.SimulationTest, nil, func(params v1_6.CCIPOCRParams) v1_6.CCIPOCRParams {
			params.ExecuteOffChainConfig.MultipleReportsEnabled = false
			return params
		})

		_, err = commonchangeset.Apply(t, tenv.Env,
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_6.SetCandidateChangeset),
				v1_6.SetCandidateChangesetConfig{
					SetCandidateConfigBase: v1_6.SetCandidateConfigBase{
						HomeChainSelector: tenv.HomeChainSel,
						FeedChainSelector: tenv.FeedChainSel,
					},
					PluginInfo: []v1_6.SetCandidatePluginInfo{
						{
							OCRConfigPerRemoteChainSelector: map[uint64]v1_6.CCIPOCRParams{
								source: execParams,
							},
							PluginType: types.PluginTypeCCIPExec,
						},
					},
				},
			),
		)
		require.NoError(t, err, "Should succeed when MultipleReportsEnabled=false regardless of EnforceOutOfOrder")
	})
}
