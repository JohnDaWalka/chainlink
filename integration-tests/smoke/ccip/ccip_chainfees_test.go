package smoke

import (
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func Test_CCIPChainFeeUpdates(t *testing.T) {
	// The outcome of the plugin is

	//type Outcome struct {
	//	// Each Gas Price is the combination of Execution and DataAvailability Fees using bitwise operations
	//	GasPrices []cciptypes.GasPriceChain `json:"gasPrices"`
	//}

	// And the observation algorithm is:

	// getFeeComponents from all the configured source chains.
	// getNativeTokenPrices from all the configured source chains
	// get current chain fee price updates

	// this price updates go to the fee quoter.
	// so if we have this two lanes
	// chainA -> chainB
	// chainB -> chainA
	// chainB stores chain fee of chainA and chainB of chain A

	// in this test setup we can have this two chains and two lanes.

	// we can start the commit plugin
	// then query the commitStore/feeQuoter to get chainFees

	// we can start/stop the plugin
	// and in the meantime we can update the chain fees by making contract txs

	lggr := logger.TestLogger(t)
	ctx := changeset.Context(t)

	e, _, _ := testsetups.NewLocalDevEnvironmentWithDefaultPrice(t, lggr, &changeset.TestConfigs{})
	state, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)

	allChainSelectors := maps.Keys(e.Env.Chains)
	assert.GreaterOrEqual(t, len(allChainSelectors), 2, "test requires at least 2 chains")

	sourceChain1 := allChainSelectors[0]
	sourceChain2 := allChainSelectors[1]

	require.NoError(t, changeset.AddLaneWithDefaultPricesAndFeeQuoterConfig(e.Env, state, sourceChain1, sourceChain2, false))
	require.NoError(t, changeset.AddLaneWithDefaultPricesAndFeeQuoterConfig(e.Env, state, sourceChain2, sourceChain1, false))

	// keep reading the chain fees on feeQuoter or offRamp

	for i := 0; i < 100; i++ {
		feeQuoter1 := state.Chains[sourceChain1].FeeQuoter
		chainFee2, err := feeQuoter1.GetDestinationChainGasPrice(&bind.CallOpts{Context: ctx}, sourceChain2)
		require.NoError(t, err)

		feeQuoter2 := state.Chains[sourceChain2].FeeQuoter
		chainFee1, err := feeQuoter2.GetDestinationChainGasPrice(&bind.CallOpts{Context: ctx}, sourceChain1)
		require.NoError(t, err)

		routerChain1 := state.Chains[sourceChain1].Router
		routerChain2 := state.Chains[sourceChain2].Router

		someMsg := router.ClientEVM2AnyMessage{
			Receiver:     utils.ZeroAddress.Bytes(),
			Data:         []byte(fmt.Sprintf("hello world %d", i)),
			TokenAmounts: nil,
			FeeToken:     common.HexToAddress("0x0"),
			ExtraArgs:    nil,
		}

		msgFeeToChain2, err := routerChain1.GetFee(&bind.CallOpts{Context: ctx}, sourceChain2, someMsg)
		require.NoError(t, err)

		msgFeeToChain1, err := routerChain2.GetFee(&bind.CallOpts{Context: ctx}, sourceChain1, someMsg)
		require.NoError(t, err)

		t.Logf("chainFee1 (stored in chain2): %v", chainFee1)
		t.Logf("chainFee2 (stored in chain1): %v", chainFee2)
		t.Logf("msgFeeToChain2: %v", msgFeeToChain2)
		t.Logf("msgFeeToChain1: %v", msgFeeToChain1)

		time.Sleep(100 * time.Millisecond)
	}
}
