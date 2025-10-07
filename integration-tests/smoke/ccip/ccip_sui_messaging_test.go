package ccip

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common/math"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	suiBind "github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_fee_quoter "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/fee_quoter"
	suiutil "github.com/smartcontractkit/chainlink-sui/bindings/utils"
	sui_deployment "github.com/smartcontractkit/chainlink-sui/deployment"
	sui_cs "github.com/smartcontractkit/chainlink-sui/deployment/changesets"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	linkops "github.com/smartcontractkit/chainlink-sui/deployment/ops/link"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	mlt "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/messagelimitationstest"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/messagingtest"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func Test_CCIP_Messaging_Sui2EVM(t *testing.T) {
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

	suiState, err := sui_deployment.LoadOnchainStatesui(e.Env)
	require.NoError(t, err)

	t.Log("Source chain (Sui): ", sourceChain, "Dest chain (EVM): ", destChain)

	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)

	suiSenderAddr, err := e.Env.BlockChains.SuiChains()[sourceChain].Signer.GetAddress()
	require.NoError(t, err)

	normalizedAddr, err := suiutil.ConvertStringToAddressBytes(suiSenderAddr)
	require.NoError(t, err)

	// SUI FeeToken
	// mint link token to use as feeToken
	_, output, err := commoncs.ApplyChangesets(t, e.Env, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(sui_cs.MintLinkToken{}, sui_cs.MintLinkTokenConfig{
			ChainSelector:  sourceChain,
			TokenPackageId: suiState[sourceChain].LinkTokenAddress,
			TreasuryCapId:  suiState[sourceChain].LinkTokenTreasuryCapId,
			Amount:         1000000000000, // 1000 Link with 1e9
		}),
	})
	require.NoError(t, err)

	rawOutput := output[0].Reports[0]
	outputMap, ok := rawOutput.Output.(sui_ops.OpTxResult[linkops.MintLinkTokenOutput])
	require.True(t, ok)

	var (
		nonce  uint64
		sender = common.LeftPadBytes(normalizedAddr[:], 32)
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

	// Get Sui FQ
	ctx := testcontext.Get(t)

	suiFeeQuoter, err := module_fee_quoter.NewFeeQuoter(
		state.SuiChains[sourceChain].CCIPAddress,
		e.Env.BlockChains.SuiChains()[sourceChain].Client)
	require.NoError(t, err)

	suiFeeQuoterDestChainConfig, err := suiFeeQuoter.DevInspect().GetDestChainConfig(ctx, &suiBind.CallOpts{
		Signer:           e.Env.BlockChains.SuiChains()[sourceChain].Signer,
		WaitForExecution: true,
	}, suiBind.Object{Id: state.SuiChains[sourceChain].CCIPObjectRef}, destChain)
	require.NoError(t, err, "Failed to get destination chain config")

	// For testing messages that revert on source
	mltTestSetup := mlt.NewTestSetup(
		t,
		state,
		sourceChain,
		destChain,
		common.HexToAddress(outputMap.Objects.MintedLinkTokenObjectId),
		suiFeeQuoterDestChainConfig,
		false, // testRouter
		true,  // validateResp
		mlt.WithDeployedEnv(e),
	)

	invalidDestChainSelectorTestSetup := mlt.NewTestSetup(
		t,
		state,
		sourceChain,
		destChain,
		common.HexToAddress("0x0"),
		suiFeeQuoterDestChainConfig,
		false, // testRouter
		true,  // validateResp
		mlt.WithDeployedEnv(e),
	)

	t.Run("Message to EVM - Should Succeed", func(t *testing.T) {
		out = messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:              setup,
				Nonce:                  &nonce,
				ValidationType:         messagingtest.ValidationTypeExec,
				Receiver:               state.Chains[destChain].Receiver.Address().Bytes(),
				ExtraArgs:              nil,
				Replayed:               true,
				FeeToken:               outputMap.Objects.MintedLinkTokenObjectId,
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
			},
		)
	})

	t.Run("Max Data Bytes - Should Succeed", func(t *testing.T) {
		latestHead, err := testhelpers.LatestBlock(ctx, e.Env, destChain)
		require.NoError(t, err)
		message := []byte(strings.Repeat("0", int(sui_deployment.DefaultCCIPSeqConfig.MaxDataBytes)))
		messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:      setup,
				ValidationType: messagingtest.ValidationTypeExec,
				FeeToken:       outputMap.Objects.MintedLinkTokenObjectId,
				Receiver:       state.Chains[destChain].Receiver.Address().Bytes(),
				MsgData:        message,
				// Just ensuring enough gas is provided to execute the message, doesn't matter if it's way too much
				ExtraArgs:              testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(300000), false),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
				ExtraAssertions: []func(t *testing.T){
					func(t *testing.T) { assertEvmMessageReceived(ctx, t, state, destChain, latestHead, message) },
				},
			},
		)
	})

	t.Run("Max Gas Limit - Should Succeed", func(t *testing.T) {
		latestHead, err := testhelpers.LatestBlock(ctx, e.Env, destChain)
		require.NoError(t, err)
		message := []byte("Hello EVM, from Sui!")
		messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:              setup,
				ValidationType:         messagingtest.ValidationTypeExec,
				FeeToken:               outputMap.Objects.MintedLinkTokenObjectId,
				Receiver:               state.Chains[destChain].Receiver.Address().Bytes(),
				MsgData:                message,
				ExtraArgs:              testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(int64(sui_deployment.DefaultCCIPSeqConfig.MaxPerMsgGasLimit)), false),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
				ExtraAssertions: []func(t *testing.T){
					func(t *testing.T) { assertEvmMessageReceived(ctx, t, state, destChain, latestHead, message) },
				},
			},
		)
	})

	t.Run("Max Data Bytes + 1 - Should Fail", func(t *testing.T) {
		message := []byte(strings.Repeat("0", int(sui_deployment.DefaultCCIPSeqConfig.MaxDataBytes)+1))
		mlt.Run(mlt.TestCase{
			TestSetup: mltTestSetup,
			Name:      "Max Data Bytes + 1 - Should Fail",
			Msg: testhelpers.SuiSendRequest{
				Receiver:  state.Chains[destChain].Receiver.Address().Bytes(),
				Data:      message,
				FeeToken:  outputMap.Objects.MintedLinkTokenObjectId,
				ExtraArgs: nil,
			},
			ExpRevert: true,
		})
	})

	t.Run("Max Data Bytes + 1 to EOA - Should Fail", func(t *testing.T) {
		message := []byte(strings.Repeat("0", int(sui_deployment.DefaultCCIPSeqConfig.MaxDataBytes)+1))
		mlt.Run(mlt.TestCase{
			TestSetup: mltTestSetup,
			Name:      "Max Data Bytes + 1 to EOA - Should Fail",
			Msg: testhelpers.SuiSendRequest{
				Receiver:  e.Env.BlockChains.EVMChains()[destChain].DeployerKey.From[:], // Sending to EOA
				Data:      message,
				FeeToken:  outputMap.Objects.MintedLinkTokenObjectId,
				ExtraArgs: nil,
			},
			ExpRevert: true,
		})
	})

	t.Run("Max Gas Limit + 1 - Should Fail", func(t *testing.T) {
		message := []byte("Hello EVM, from Sui!")
		mlt.Run(mlt.TestCase{
			TestSetup: mltTestSetup,
			Name:      "Max Gas Limit + 1 - Should Fail",
			Msg: testhelpers.SuiSendRequest{
				Receiver:  e.Env.BlockChains.EVMChains()[destChain].DeployerKey.From[:],
				Data:      message,
				FeeToken:  outputMap.Objects.MintedLinkTokenObjectId,
				ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(int64(sui_deployment.DefaultCCIPSeqConfig.MaxPerMsgGasLimit)+1), false),
			},
			ExpRevert: true,
		})
	})

	t.Run("Missing ExtraArgs - Should Fail", func(t *testing.T) {
		message := []byte("Hello EVM, from Sui!")
		mlt.Run(mlt.TestCase{
			TestSetup: mltTestSetup,
			Name:      "Missing ExtraArgs - Should Fail",
			Msg: testhelpers.SuiSendRequest{
				Receiver:  state.Chains[destChain].Receiver.Address().Bytes(),
				Data:      message,
				FeeToken:  outputMap.Objects.MintedLinkTokenObjectId,
				ExtraArgs: []byte{},
			},
			ExpRevert: true,
		})
	})

	t.Run("Send message to invalid receiver - Should Fail", func(t *testing.T) {
		message := []byte("Hello EVM, from Sui!")
		mlt.Run(mlt.TestCase{
			TestSetup: mltTestSetup,
			Name:      "Send message to invalid receiver - Should Fail",
			Msg: testhelpers.SuiSendRequest{
				Receiver:  []byte("0x0000"),
				Data:      message,
				FeeToken:  outputMap.Objects.MintedLinkTokenObjectId,
				ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(300000), false),
			},
			ExpRevert: true,
		})
	})

	t.Run("Send message to invalid chain selector - Should Fail", func(t *testing.T) {
		message := []byte("Hello EVM, from Sui!")
		mlt.Run(mlt.TestCase{
			TestSetup: invalidDestChainSelectorTestSetup,
			Name:      "Send message to invalid chain selector - Should Fail",
			Msg: testhelpers.SuiSendRequest{
				Receiver:  state.Chains[destChain].Receiver.Address().Bytes(),
				Data:      message,
				FeeToken:  outputMap.Objects.MintedLinkTokenObjectId,
				ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(300000), false),
			},
			ExpRevert: true,
		})
	})

	fmt.Printf("out: %v\n", out)
}

func Test_CCIP_Messaging_EVM2Sui(t *testing.T) {
	lggr := logger.TestLogger(t)
	e, _, _ := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithSuiChains(1),
	)

	evmChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyEVM))
	suiChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilySui))

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	sourceChain := evmChainSelectors[0]
	destChain := suiChainSelectors[0]

	lggr.Debug("Source chain (EVM): ", sourceChain, "Dest chain (Sui): ", destChain)

	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)

	var (
		nonce  uint64
		sender = common.LeftPadBytes(e.Env.BlockChains.EVMChains()[sourceChain].DeployerKey.From.Bytes(), 32)
		setup  = messagingtest.NewTestSetupWithDeployedEnv(
			t,
			e,
			state,
			sourceChain,
			destChain,
			sender,
			false, // test router
		)

		// Tokens
		nativeFeeToken = "0x0"
		evmLinkToken   = state.Chains[sourceChain].LinkToken
		wethToken      = state.Chains[sourceChain].Weth9
	)

	// Deploy SUI Reciever
	_, output, err := commoncs.ApplyChangesets(t, e.Env, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(sui_cs.DeployDummyReciever{}, sui_cs.DeployDummyRecieverConfig{
			SuiChainSelector: destChain,
			McmsOwner:        "0x1",
		}),
	})
	require.NoError(t, err)

	rawOutput := output[0].Reports[0]

	outputMap, ok := rawOutput.Output.(sui_ops.OpTxResult[ccipops.DeployDummyReceiverObjects])
	require.True(t, ok)

	id := strings.TrimPrefix(outputMap.PackageId, "0x")
	receiverByteDecoded, err := hex.DecodeString(id)
	require.NoError(t, err)

	// register the reciever
	_, _, err = commoncs.ApplyChangesets(t, e.Env, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(sui_cs.RegisterDummyReciever{}, sui_cs.RegisterDummyReceiverConfig{
			SuiChainSelector:       destChain,
			CCIPObjectRefObjectId:  state.SuiChains[destChain].CCIPObjectRef,
			DummyReceiverPackageId: outputMap.PackageId,
		}),
	})
	require.NoError(t, err)

	receiverByte := receiverByteDecoded

	var clockObj [32]byte
	copy(clockObj[:], hexutil.MustDecode(
		"0x0000000000000000000000000000000000000000000000000000000000000006",
	))

	var stateObj [32]byte
	copy(stateObj[:], hexutil.MustDecode(
		outputMap.Objects.CCIPReceiverStateObjectId,
	))

	recieverObjectIds := [][32]byte{clockObj, stateObj}

	t.Run("Message to Sui - Should Succeed", func(t *testing.T) {
		// ccipChainState := state.SuiChains[destChain]
		message := []byte("Hello Sui, from EVM!")
		messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:              setup,
				Nonce:                  &nonce,
				ValidationType:         messagingtest.ValidationTypeExec,
				Receiver:               receiverByte,
				MsgData:                message,
				ExtraArgs:              testhelpers.MakeSuiExtraArgs(1000000, true, recieverObjectIds, [32]byte{}),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
			},
		)
	})

	ctx := testcontext.Get(t)
	srcFeeQuoterDestChainConfig, err := state.Chains[sourceChain].FeeQuoter.GetDestChainConfig(&bind.CallOpts{Context: ctx}, destChain)
	require.NoError(t, err, "Failed to get destination chain config")

	// grant mint role
	tx, err := evmLinkToken.GrantMintRole(e.Env.BlockChains.EVMChains()[sourceChain].DeployerKey, common.BytesToAddress(sender))
	_, err = cldf.ConfirmIfNoError(e.Env.BlockChains.EVMChains()[sourceChain], tx, err)
	require.NoError(t, err)

	// mint token and approve to router
	tx, err = evmLinkToken.Mint(e.Env.BlockChains.EVMChains()[sourceChain].DeployerKey, common.BytesToAddress(sender), deployment.E18Mult(10_000))
	_, err = cldf.ConfirmIfNoError(e.Env.BlockChains.EVMChains()[sourceChain], tx, err)
	require.NoError(t, err)

	tx, err = evmLinkToken.Approve(e.Env.BlockChains.EVMChains()[sourceChain].DeployerKey, state.Chains[sourceChain].Router.Address(), math.MaxBig256)
	_, err = cldf.ConfirmIfNoError(e.Env.BlockChains.EVMChains()[sourceChain], tx, err)
	require.NoError(t, err)

	// Deposit 1 ETH to get WETH
	wethTransactOpts := *e.Env.BlockChains.EVMChains()[sourceChain].DeployerKey
	wethTransactOpts.Value = deployment.E18Mult(1)
	tx, err = wethToken.Deposit(&wethTransactOpts)
	_, err = cldf.ConfirmIfNoError(e.Env.BlockChains.EVMChains()[sourceChain], tx, err)
	require.NoError(t, err)

	tx, err = wethToken.Approve(e.Env.BlockChains.EVMChains()[sourceChain].DeployerKey, state.Chains[sourceChain].Router.Address(), math.MaxBig256)
	_, err = cldf.ConfirmIfNoError(e.Env.BlockChains.EVMChains()[sourceChain], tx, err)
	require.NoError(t, err)

	// For testing messages that revert on source
	mltTestSetup := mlt.NewTestSetup(
		t,
		state,
		sourceChain,
		destChain,
		common.HexToAddress("0x0"),
		srcFeeQuoterDestChainConfig,
		false, // testRouter
		true,  // validateResp
		mlt.WithDeployedEnv(e),
	)

	invalidDestChainSelectorTestSetup := mlt.NewTestSetup(
		t,
		state,
		sourceChain,
		destChain,
		common.HexToAddress("0x0"),
		srcFeeQuoterDestChainConfig,
		false, // testRouter
		true,  // validateResp
		mlt.WithDeployedEnv(e),
	)

	t.Run("Max Data Bytes - Should Succeed", func(t *testing.T) {
		message := []byte(strings.Repeat("0", int(srcFeeQuoterDestChainConfig.MaxDataBytes)))
		messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:      setup,
				Nonce:          &nonce,
				ValidationType: messagingtest.ValidationTypeExec,
				Receiver:       receiverByte,
				MsgData:        message,
				// true for out of order execution, which is necessary and enforced for Sui
				ExtraArgs:              testhelpers.MakeSuiExtraArgs(1000000, true, recieverObjectIds, [32]byte{}),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
			},
		)
	})

	t.Run("Fee Token (LINK) - Should Succeed", func(t *testing.T) {
		message := []byte("Hello Sui, from EVM!")
		messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:      setup,
				Nonce:          &nonce,
				ValidationType: messagingtest.ValidationTypeExec,
				Receiver:       receiverByte,
				MsgData:        message,
				// true for out of order execution, which is necessary and enforced for Sui
				ExtraArgs:              testhelpers.MakeSuiExtraArgs(1000000, true, recieverObjectIds, [32]byte{}),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
				FeeToken:               state.Chains[sourceChain].LinkToken.Address().String(),
			},
		)
	})

	t.Run("Max Data Bytes + 1 - Should Fail", func(t *testing.T) {
		message := []byte(strings.Repeat("0", int(srcFeeQuoterDestChainConfig.MaxDataBytes)+1))
		mlt.Run(mlt.TestCase{
			TestSetup: mltTestSetup,
			Name:      "Max Data Bytes + 1 - Should Fail",
			Msg: router.ClientEVM2AnyMessage{
				Receiver:  receiverByte,
				Data:      message,
				FeeToken:  common.HexToAddress(nativeFeeToken),
				ExtraArgs: testhelpers.MakeSuiExtraArgs(uint64(srcFeeQuoterDestChainConfig.MaxPerMsgGasLimit)+1, true, recieverObjectIds, [32]byte{}),
			},
			ExpRevert: true,
		})
	})

	t.Run("Max Data Bytes + 1 to EOA - Should Fail", func(t *testing.T) {
		message := []byte(strings.Repeat("0", int(srcFeeQuoterDestChainConfig.MaxDataBytes)+1))
		mlt.Run(mlt.TestCase{
			TestSetup: mltTestSetup,
			Name:      "Max Data Bytes + 1 to EOA - Should Fail",
			Msg: router.ClientEVM2AnyMessage{
				Receiver:  receiverByte, // Sending to EOA
				Data:      message,
				FeeToken:  common.HexToAddress(nativeFeeToken),
				ExtraArgs: testhelpers.MakeSuiExtraArgs(uint64(srcFeeQuoterDestChainConfig.MaxPerMsgGasLimit)+1, true, recieverObjectIds, [32]byte{}),
			},
			ExpRevert: true,
		})
	})

	t.Run("Missing ExtraArgs - Should Fail", func(t *testing.T) {
		message := []byte("Hello Sui, from EVM!")
		mlt.Run(mlt.TestCase{
			TestSetup: mltTestSetup,
			Name:      "Missing ExtraArgs - Should Fail",
			Msg: router.ClientEVM2AnyMessage{
				Receiver:  receiverByte,
				Data:      message,
				FeeToken:  common.HexToAddress(nativeFeeToken),
				ExtraArgs: []byte{},
			},
			ExpRevert: true,
		})
	})

	t.Run("OutOfOrder Execution False - Should Fail", func(t *testing.T) {
		message := []byte("Hello Sui, from EVM!")
		mlt.Run(mlt.TestCase{
			TestSetup: mltTestSetup,
			Name:      "OutOfOrder Execution False - Should Fail",
			Msg: router.ClientEVM2AnyMessage{
				Receiver:  receiverByte,
				Data:      message,
				FeeToken:  common.HexToAddress(nativeFeeToken),
				ExtraArgs: testhelpers.MakeSuiExtraArgs(100000, false, recieverObjectIds, [32]byte{}),
			},
			ExpRevert: true,
		})
	})

	t.Run("Send message to invalid receiver - Should Fail", func(t *testing.T) {
		message := []byte("Hello Sui, from EVM!")
		mlt.Run(mlt.TestCase{
			TestSetup: mltTestSetup,
			Name:      "Send message to invalid receiver - Should Fail",
			Msg: router.ClientEVM2AnyMessage{
				Receiver:  []byte("0x000"),
				Data:      message,
				FeeToken:  common.HexToAddress(nativeFeeToken),
				ExtraArgs: testhelpers.MakeSuiExtraArgs(100000, false, recieverObjectIds, [32]byte{}),
			},
			ExpRevert: true,
		})
	})

	t.Run("Send message to invalid chain selector - Should Fail", func(t *testing.T) {
		message := []byte("Hello Sui, from EVM!")
		mlt.Run(mlt.TestCase{
			TestSetup: invalidDestChainSelectorTestSetup,
			Name:      "Send message to invalid chain selector - Should Fail",
			Msg: router.ClientEVM2AnyMessage{
				Receiver:  receiverByte,
				Data:      message,
				FeeToken:  common.HexToAddress(nativeFeeToken),
				ExtraArgs: testhelpers.MakeSuiExtraArgs(100000, false, recieverObjectIds, [32]byte{}),
			},
			ExpRevert: true,
		})
	})
}
