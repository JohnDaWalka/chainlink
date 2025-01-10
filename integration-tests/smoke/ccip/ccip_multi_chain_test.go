package smoke

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
)

func TestWithFewerLanesThanAllPossibleLanesBetweenChains(t *testing.T) {
	e, _ := testsetups.NewIntegrationEnvironment(t, changeset.WithChains(3))
	state, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)
	require.Len(t, e.Env.AllChainSelectors(), 3)
	src1, chain2 := e.Env.AllChainSelectors()[0], e.Env.AllChainSelectors()[1]
	// Add only one lane between the first two chains
	changeset.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, src1, chain2, true)
	startBlocks := make(map[uint64]*uint64)
	latesthdr, err := e.Env.Chains[chain2].Client.HeaderByNumber(testcontext.Get(t), nil)
	require.NoError(t, err)
	block := latesthdr.Number.Uint64()
	startBlocks[chain2] = &block
	expectedSeqNumExec := make(map[changeset.SourceDestPair][]uint64)
	// Send a message from the first chain to the second chain
	msgSentEvent, err := changeset.DoSendRequest(
		t, e.Env, state,
		changeset.WithSourceChain(src1),
		changeset.WithDestChain(chain2),
		changeset.WithTestRouter(true),
		changeset.WithEvm2AnyMessage(router.ClientEVM2AnyMessage{
			Receiver:     common.LeftPadBytes(state.Chains[chain2].Receiver.Address().Bytes(), 32),
			Data:         []byte("hello"),
			TokenAmounts: nil,
			FeeToken:     common.HexToAddress("0x0"),
			ExtraArgs:    nil,
		}))
	require.NoError(t, err)

	expectedSeqNumExec[changeset.SourceDestPair{
		SourceChainSelector: src1,
		DestChainSelector:   chain2,
	}] = []uint64{msgSentEvent.SequenceNumber}

	// Wait for all exec reports to land
	changeset.ConfirmExecWithSeqNrsForAll(t, e.Env, state, expectedSeqNumExec, startBlocks)
}
