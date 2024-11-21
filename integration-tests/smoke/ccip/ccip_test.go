package smoke

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	commonutils "github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// TestAddLane covers the workflow of adding a lane between two chains and enabling it.
// It also covers the case where the onRamp is disabled on the OffRamp contract initially and then enabled.
func TestAddLane(t *testing.T) {
	t.Parallel()
	// We add more chains to the chainlink nodes than the number of chains where CCIP is deployed.
	e := changeset.NewMemoryEnvironmentWithJobsAndContracts(t, logger.TestLogger(t), 2, 4)
	// Here we have CR + nodes set up, but no CCIP contracts deployed.
	state, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)

	selectors := e.Env.AllChainSelectors()
	chain1, chain2 := selectors[0], selectors[1]

	// We expect no lanes available on any chain.
	for _, sel := range []uint64{chain1, chain2} {
		chain := state.Chains[sel]
		offRamps, err := chain.Router.GetOffRamps(nil)
		require.NoError(t, err)
		require.Len(t, offRamps, 0)
	}

	replayBlocks, err := changeset.LatestBlocksByChain(testcontext.Get(t), e.Env.Chains)
	require.NoError(t, err)

	// Add one lane from chain1 to chain 2 and send traffic.
	require.NoError(t, changeset.AddLaneWithDefaultPrices(e.Env, state, chain1, chain2))

	changeset.ReplayLogs(t, e.Env.Offchain, replayBlocks)
	time.Sleep(30 * time.Second)
	// disable the onRamp initially on OffRamp
	disableRampTx, err := state.Chains[chain2].OffRamp.ApplySourceChainConfigUpdates(e.Env.Chains[chain2].DeployerKey, []offramp.OffRampSourceChainConfigArgs{
		{
			Router:              state.Chains[chain2].Router.Address(),
			SourceChainSelector: chain1,
			IsEnabled:           false,
			OnRamp:              common.LeftPadBytes(state.Chains[chain1].OnRamp.Address().Bytes(), 32),
		},
	})
	_, err = deployment.ConfirmIfNoError(e.Env.Chains[chain2], disableRampTx, err)
	require.NoError(t, err)

	for _, sel := range []uint64{chain1, chain2} {
		chain := state.Chains[sel]
		offRamps, err := chain.Router.GetOffRamps(nil)
		require.NoError(t, err)
		if sel == chain2 {
			require.Len(t, offRamps, 1)
			srcCfg, err := chain.OffRamp.GetSourceChainConfig(nil, chain1)
			require.NoError(t, err)
			require.Equal(t, common.LeftPadBytes(state.Chains[chain1].OnRamp.Address().Bytes(), 32), srcCfg.OnRamp)
			require.False(t, srcCfg.IsEnabled)
		} else {
			require.Len(t, offRamps, 0)
		}
	}

	latesthdr, err := e.Env.Chains[chain2].Client.HeaderByNumber(testcontext.Get(t), nil)
	require.NoError(t, err)
	startBlock := latesthdr.Number.Uint64()
	// Send traffic on the first lane and it should not be processed by the plugin as onRamp is disabled
	// we will check this by confirming that the message is not executed by the end of the test
	msgSentEvent1 := changeset.TestSendRequest(t, e.Env, state, chain1, chain2, false, router.ClientEVM2AnyMessage{
		Receiver:     common.LeftPadBytes(state.Chains[chain2].Receiver.Address().Bytes(), 32),
		Data:         []byte("hello world"),
		TokenAmounts: nil,
		FeeToken:     common.HexToAddress("0x0"),
		ExtraArgs:    nil,
	})
	require.Equal(t, uint64(1), msgSentEvent1.SequenceNumber)

	// Add another lane
	require.NoError(t, changeset.AddLaneWithDefaultPrices(e.Env, state, chain2, chain1))

	// Send traffic on the second lane and it should succeed
	latesthdr, err = e.Env.Chains[chain1].Client.HeaderByNumber(testcontext.Get(t), nil)
	require.NoError(t, err)
	startBlock2 := latesthdr.Number.Uint64()
	msgSentEvent2 := changeset.TestSendRequest(t, e.Env, state, chain2, chain1, false, router.ClientEVM2AnyMessage{
		Receiver:     common.LeftPadBytes(state.Chains[chain2].Receiver.Address().Bytes(), 32),
		Data:         []byte("hello world"),
		TokenAmounts: nil,
		FeeToken:     common.HexToAddress("0x0"),
		ExtraArgs:    nil,
	})
	require.Equal(t, uint64(1), msgSentEvent2.SequenceNumber)
	require.NoError(t, commonutils.JustError(changeset.ConfirmExecWithSeqNr(t, e.Env.Chains[chain2], e.Env.Chains[chain1], state.Chains[chain1].OffRamp, &startBlock2, msgSentEvent2.SequenceNumber)))

	// now check for the previous message from chain 1 to chain 2 that it has not been executed till now as the onRamp was disabled
	changeset.ConfirmNoExecConsistentlyWithSeqNr(t, e.Env.Chains[chain1], e.Env.Chains[chain2], state.Chains[chain2].OffRamp, msgSentEvent1.SequenceNumber, 30*time.Second)

	// enable the onRamp on OffRamp
	enableRampTx, err := state.Chains[chain2].OffRamp.ApplySourceChainConfigUpdates(e.Env.Chains[chain2].DeployerKey, []offramp.OffRampSourceChainConfigArgs{
		{
			Router:              state.Chains[chain2].Router.Address(),
			SourceChainSelector: chain1,
			IsEnabled:           true,
			OnRamp:              common.LeftPadBytes(state.Chains[chain1].OnRamp.Address().Bytes(), 32),
		},
	})
	_, err = deployment.ConfirmIfNoError(e.Env.Chains[chain2], enableRampTx, err)
	require.NoError(t, err)

	srcCfg, err := state.Chains[chain2].OffRamp.GetSourceChainConfig(nil, chain1)
	require.NoError(t, err)
	require.Equal(t, common.LeftPadBytes(state.Chains[chain1].OnRamp.Address().Bytes(), 32), srcCfg.OnRamp)
	require.True(t, srcCfg.IsEnabled)

	// we need the replay here otherwise plugin is not able to locate the message
	changeset.ReplayLogs(t, e.Env.Offchain, replayBlocks)
	time.Sleep(30 * time.Second)
	// Now that the onRamp is enabled, the request should be processed
	require.NoError(t, commonutils.JustError(changeset.ConfirmExecWithSeqNr(t, e.Env.Chains[chain1], e.Env.Chains[chain2], state.Chains[chain2].OffRamp, &startBlock, msgSentEvent1.SequenceNumber)))
}

func TestInitialDeployOnLocal(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	tenv, _, _ := testsetups.NewLocalDevEnvironmentWithDefaultPrice(t, lggr)
	e := tenv.Env
	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)

	// Add all lanes
	require.NoError(t, changeset.AddLanesForAll(e, state))
	// Need to keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)
	// Send a message from each chain to every other chain.
	expectedSeqNum := make(map[changeset.SourceDestPair]uint64)
	for src := range e.Chains {
		for dest, destChain := range e.Chains {
			if src == dest {
				continue
			}
			latesthdr, err := destChain.Client.HeaderByNumber(testcontext.Get(t), nil)
			require.NoError(t, err)
			block := latesthdr.Number.Uint64()
			startBlocks[dest] = &block
			msgSentEvent := changeset.TestSendRequest(t, e, state, src, dest, false, router.ClientEVM2AnyMessage{
				Receiver:     common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32),
				Data:         []byte("hello world"),
				TokenAmounts: nil,
				FeeToken:     common.HexToAddress("0x0"),
				ExtraArgs:    nil,
			})
			expectedSeqNum[changeset.SourceDestPair{
				SourceChainSelector: src,
				DestChainSelector:   dest,
			}] = msgSentEvent.SequenceNumber
		}
	}

	// Wait for all commit reports to land.
	changeset.ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)

	// After commit is reported on all chains, token prices should be updated in FeeQuoter.
	for dest := range e.Chains {
		linkAddress := state.Chains[dest].LinkToken.Address()
		feeQuoter := state.Chains[dest].FeeQuoter
		timestampedPrice, err := feeQuoter.GetTokenPrice(nil, linkAddress)
		require.NoError(t, err)
		require.Equal(t, changeset.MockLinkPrice, timestampedPrice.Value)
	}

	// Wait for all exec reports to land
	changeset.ConfirmExecWithSeqNrForAll(t, e, state, expectedSeqNum, startBlocks)

	// TODO: Apply the proposal.
}

func TestTokenTransfer(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	tenv, _, _ := testsetups.NewLocalDevEnvironmentWithDefaultPrice(t, lggr)
	e := tenv.Env
	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)

	srcToken, _, dstToken, _, err := changeset.DeployTransferableToken(
		lggr,
		tenv.Env.Chains,
		tenv.HomeChainSel,
		tenv.FeedChainSel,
		state,
		e.ExistingAddresses,
		"MY_TOKEN",
	)
	require.NoError(t, err)

	// Add all lanes
	require.NoError(t, changeset.AddLanesForAll(e, state))
	// Need to keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)
	// Send a message from each chain to every other chain.
	expectedSeqNum := make(map[changeset.SourceDestPair]uint64)

	twoCoins := new(big.Int).Mul(big.NewInt(1e18), big.NewInt(2))
	tx, err := srcToken.Mint(
		e.Chains[tenv.HomeChainSel].DeployerKey,
		e.Chains[tenv.HomeChainSel].DeployerKey.From,
		new(big.Int).Mul(twoCoins, big.NewInt(10)),
	)
	require.NoError(t, err)
	_, err = e.Chains[tenv.HomeChainSel].Confirm(tx)
	require.NoError(t, err)

	tx, err = dstToken.Mint(
		e.Chains[tenv.FeedChainSel].DeployerKey,
		e.Chains[tenv.FeedChainSel].DeployerKey.From,
		new(big.Int).Mul(twoCoins, big.NewInt(10)),
	)
	require.NoError(t, err)
	_, err = e.Chains[tenv.FeedChainSel].Confirm(tx)
	require.NoError(t, err)

	tx, err = srcToken.Approve(e.Chains[tenv.HomeChainSel].DeployerKey, state.Chains[tenv.HomeChainSel].Router.Address(), twoCoins)
	require.NoError(t, err)
	_, err = e.Chains[tenv.HomeChainSel].Confirm(tx)
	require.NoError(t, err)
	tx, err = dstToken.Approve(e.Chains[tenv.FeedChainSel].DeployerKey, state.Chains[tenv.FeedChainSel].Router.Address(), twoCoins)
	require.NoError(t, err)
	_, err = e.Chains[tenv.FeedChainSel].Confirm(tx)
	require.NoError(t, err)

	tokens := map[uint64][]router.ClientEVMTokenAmount{
		tenv.HomeChainSel: {{
			Token:  srcToken.Address(),
			Amount: twoCoins,
		}},
		tenv.FeedChainSel: {{
			Token:  dstToken.Address(),
			Amount: twoCoins,
		}},
	}

	for src := range e.Chains {
		for dest, destChain := range e.Chains {
			if src == dest {
				continue
			}
			latesthdr, err := destChain.Client.HeaderByNumber(testcontext.Get(t), nil)
			require.NoError(t, err)
			block := latesthdr.Number.Uint64()
			startBlocks[dest] = &block

			var (
				receiver = common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32)
				data     = []byte("hello world")
				feeToken = common.HexToAddress("0x0")
			)
			if src == tenv.HomeChainSel && dest == tenv.FeedChainSel {
				msgSentEvent := changeset.TestSendRequest(t, e, state, src, dest, false, router.ClientEVM2AnyMessage{
					Receiver:     receiver,
					Data:         data,
					TokenAmounts: tokens[src],
					FeeToken:     feeToken,
					ExtraArgs:    nil,
				})
				expectedSeqNum[changeset.SourceDestPair{
					SourceChainSelector: src,
					DestChainSelector:   dest,
				}] = msgSentEvent.SequenceNumber
			} else {
				msgSentEvent := changeset.TestSendRequest(t, e, state, src, dest, false, router.ClientEVM2AnyMessage{
					Receiver:     receiver,
					Data:         data,
					TokenAmounts: nil,
					FeeToken:     feeToken,
					ExtraArgs:    nil,
				})
				expectedSeqNum[changeset.SourceDestPair{
					SourceChainSelector: src,
					DestChainSelector:   dest,
				}] = msgSentEvent.SequenceNumber
			}
		}
	}

	// Wait for all commit reports to land.
	changeset.ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)

	// After commit is reported on all chains, token prices should be updated in FeeQuoter.
	for dest := range e.Chains {
		linkAddress := state.Chains[dest].LinkToken.Address()
		feeQuoter := state.Chains[dest].FeeQuoter
		timestampedPrice, err := feeQuoter.GetTokenPrice(nil, linkAddress)
		require.NoError(t, err)
		require.Equal(t, changeset.MockLinkPrice, timestampedPrice.Value)
	}

	// Wait for all exec reports to land
	changeset.ConfirmExecWithSeqNrForAll(t, e, state, expectedSeqNum, startBlocks)

	balance, err := dstToken.BalanceOf(nil, state.Chains[tenv.FeedChainSel].Receiver.Address())
	require.NoError(t, err)
	require.Equal(t, twoCoins, balance)
}
