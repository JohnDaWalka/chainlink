package ccip

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pattonkan/sui-go/sui"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/messagingtest"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"
)

func Test_CCIPMessaging_Sui2EVM(t *testing.T) {
	e, _, _ := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithSuiChains(1),
	)

	evmChainSelectors := maps.Keys(e.Env.Chains)
	suiChains := e.Env.BlockChains.SuiChains()
	suiChainSelectors := maps.Keys(suiChains)
	require.Equal(t, len(suiChainSelectors), 1)

	fmt.Println("EVM: ", evmChainSelectors)
	fmt.Println("Sui: ", suiChainSelectors)

	sourceChain := suiChainSelectors[0]
	destChain := evmChainSelectors[0]

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	t.Log("Source chain (Sui): ", sourceChain, "Dest chain (EVM): ", destChain)

	// testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)

	suiSenderAddr, err := rel.NewPrivateKeySigner(e.Env.BlockChains.SuiChains()[sourceChain].DeployerKey).GetAddress()
	require.NoError(t, err)

	suiSenderByte := sui.MustAddressFromHex(suiSenderAddr)
	var (
		replayed bool
		nonce    uint64
		sender   = common.LeftPadBytes(suiSenderByte[:], 32)
		_        messagingtest.TestCaseOutput
		setup    = messagingtest.NewTestSetupWithDeployedEnv(
			t,
			e,
			state,
			sourceChain,
			destChain,
			sender,
			false, // testRouter
		)
	)

	t.Run("Message to EVM", func(t *testing.T) {
		require.NoError(t, err)
		message := []byte("Hello EVM, from Sui!")
		_ = messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:              setup,
				Replayed:               replayed,
				Nonce:                  &nonce,
				ValidationType:         messagingtest.ValidationTypeExec,
				FeeToken:               "0xa",
				Receiver:               state.Chains[destChain].Receiver.Address().Bytes(),
				MsgData:                message,
				ExtraArgs:              nil,
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
			},
		)
	})
}
