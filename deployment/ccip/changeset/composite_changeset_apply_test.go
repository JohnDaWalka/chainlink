package changeset

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment"
	ccdeploy "github.com/smartcontractkit/chainlink/deployment/ccip"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/stretchr/testify/require"
)

func TestInitialDeploy(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := ccdeploy.Context(t)
	tenv := ccdeploy.NewMemoryEnvironment(t, lggr, 3, 4, ccdeploy.MockLinkPrice, ccdeploy.MockWethPrice)
	e := tenv.Env

	state, err := ccdeploy.LoadOnchainState(tenv.Env)
	require.NoError(t, err)
	mcmsCfg := ccdeploy.NewTestMCMSConfig(t, e)
	changesets := []deployment.ChangesetApplication{
		{
			Changeset: deployment.WrapChangeSet(DeployPrerequisites),
			Config: DeployPrerequisiteConfig{
				ChainSelectors: tenv.Env.AllChainSelectors(),
			},
		},
		{
			Changeset: deployment.WrapChangeSet(DeployChainContracts),
			Config: DeployChainContractsConfig{
				ChainSelectors:    tenv.Env.AllChainSelectors(),
				HomeChainSelector: tenv.HomeChainSel,
				MCMSCfg:           mcmsCfg,
			},
		},
		{
			Changeset: deployment.WrapChangeSet(InitialAddChain),
			Config: ccdeploy.InitialAddChainConfig{
				HomeChainSel:   tenv.HomeChainSel,
				FeedChainSel:   tenv.FeedChainSel,
				ChainsToDeploy: tenv.Env.AllChainSelectors(),
				TokenConfig:    ccdeploy.NewTestTokenConfig(state.Chains[tenv.FeedChainSel].USDFeeds),
				MCMSConfig:     mcmsCfg,
				OCRSecrets:     deployment.XXXGenerateTestOCRSecrets(),
			},
		},
		{
			Changeset: deployment.WrapChangeSet(CCIPCapabilityJobspec),
		},
	}
	tenv.Env, err = ccdeploy.ApplyChangesets(ctx, tenv.Env, changesets)
	require.NoError(t, err)
	state, err = ccdeploy.LoadOnchainState(e)
	require.NoError(t, err)
	require.NotNil(t, state.Chains[tenv.HomeChainSel].LinkToken)
	// Ensure capreg logs are up to date.
	ccdeploy.ReplayLogs(t, e.Offchain, tenv.ReplayBlocks)
	// Add all lanes
	require.NoError(t, ccdeploy.AddLanesForAll(e, state))
	// Need to keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)
	// Send a message from each chain to every other chain.
	expectedSeqNum := make(map[uint64]uint64)

	for src := range e.Chains {
		for dest, destChain := range e.Chains {
			if src == dest {
				continue
			}
			latesthdr, err := destChain.Client.HeaderByNumber(testcontext.Get(t), nil)
			require.NoError(t, err)
			block := latesthdr.Number.Uint64()
			startBlocks[dest] = &block
			msgSentEvent := ccdeploy.TestSendRequest(t, e, state, src, dest, false, router.ClientEVM2AnyMessage{
				Receiver:     common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32),
				Data:         []byte("hello"),
				TokenAmounts: nil,
				FeeToken:     common.HexToAddress("0x0"),
				ExtraArgs:    nil,
			})
			expectedSeqNum[dest] = msgSentEvent.SequenceNumber
		}
	}

	// Wait for all commit reports to land.
	ccdeploy.ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)

	// Confirm token and gas prices are updated
	ccdeploy.ConfirmTokenPriceUpdatedForAll(t, e, state, startBlocks,
		ccdeploy.DefaultInitialPrices.LinkPrice, ccdeploy.DefaultInitialPrices.WethPrice)
	// TODO: Fix gas prices?
	//ccdeploy.ConfirmGasPriceUpdatedForAll(t, e, state, startBlocks)
	//
	//// Wait for all exec reports to land
	ccdeploy.ConfirmExecWithSeqNrForAll(t, e, state, expectedSeqNum, startBlocks)
}
