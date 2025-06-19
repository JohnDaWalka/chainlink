package ccip

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pattonkan/sui-go/sui"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/messagingtest"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
	"github.com/stretchr/testify/require"
)

func Test_CCIPMessaging_Sui2EVM(t *testing.T) {
	// ctx := testhelpers.Context(t)
	e, _, _ := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithSuiChains(1),
	)

	evmChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyEVM))
	suiChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilySui))

	fmt.Println("EVM: ", evmChainSelectors[0])
	fmt.Println("Sui: ", suiChainSelectors[0])

	sourceChain := suiChainSelectors[0]
	destChain := evmChainSelectors[0]

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	t.Log("Source chain (Sui): ", sourceChain, "Dest chain (EVM): ", destChain)

	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)

	suiSenderAddr, err := rel.NewPrivateKeySigner(e.Env.BlockChains.SuiChains()[sourceChain].DeployerKey).GetAddress()
	require.NoError(t, err)

	suiSenderByte := sui.MustAddressFromHex(suiSenderAddr)
	var (
		nonce  uint64
		sender = common.LeftPadBytes(suiSenderByte[:], 32)
		out    messagingtest.TestCaseOutput
		setup  = messagingtest.NewTestSetupWithDeployedEnv(
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
		// _, err := testhelpers.LatestBlock(ctx, e.Env, destChain)
		// require.NoError(t, err)

		require.NoError(t, err)
		message := []byte("Hello EVM, from Sui!")
		out = messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:              setup,
				Nonce:                  &nonce,
				ValidationType:         messagingtest.ValidationTypeExec,
				FeeToken:               "0xa",
				Receiver:               state.Chains[destChain].Receiver.Address().Bytes(),
				MsgData:                message,
				ExtraArgs:              nil,
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
				ExtraAssertions: []func(t *testing.T){
					func(t *testing.T) {
						// iter, err := state.Chains[destChain].Receiver.FilterMessageReceived(&bind.FilterOpts{
						// 	Context: ctx,
						// 	Start:   latestHead,
						// })
						// require.NoError(t, err)
						// require.True(t, iter.Next())
						// MessageReceived doesn't emit the data unfortunately, so can't check that.
					},
				},
			},
		)
	})

	fmt.Printf("out: %v\n", out)
}
