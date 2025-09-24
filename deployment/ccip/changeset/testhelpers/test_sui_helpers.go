package testhelpers

import (
	"context"
	"fmt"
	"math/big"
	"strconv"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	suitx "github.com/block-vision/sui-go-sdk/transaction"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/message_hasher"
	suiBind "github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	suiofframp_helper "github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/ptb/offramp"
	suideps "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/sui"
	ccipclient "github.com/smartcontractkit/chainlink/deployment/ccip/shared/client"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

type SuiSendRequest struct {
	Receiver      []byte
	Data          []byte
	ExtraArgs     []byte
	FeeToken      string
	FeeTokenStore string
	TokenAmounts  []SuiTokenAmount
}

type SuiTokenAmount struct {
	Token  string
	Amount uint64
}

type RampMessageHeader struct {
	MessageId           []byte `json:"message_id"`
	SourceChainSelector string `json:"source_chain_selector"`
	DestChainSelector   string `json:"dest_chain_selector"`
	SequenceNumber      string `json:"sequence_number"`
	Nonce               string `json:"nonce"`
}

type Sui2AnyRampMessage struct {
	Header         RampMessageHeader `json:"header"`
	Sender         string            `json:"sender"`
	Data           []byte            `json:"data"`
	Receiver       []byte            `json:"receiver"`
	ExtraArgs      []byte            `json:"extra_args"`
	FeeToken       string            `json:"fee_token"`
	FeeTokenAmount string            `json:"fee_token_amount"`
	FeeValueJuels  string            `json:"fee_value_juels"`
}

type CCIPMessageSent struct {
	DestChainSelector string             `json:"dest_chain_selector"`
	SequenceNumber    string             `json:"sequence_number"`
	Message           Sui2AnyRampMessage `json:"message"`
}

func SendSuiRequestViaChainWriter(e cldf.Environment, cfg *ccipclient.CCIPSendReqConfig) (*ccipclient.AnyMsgSentEvent, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return &ccipclient.AnyMsgSentEvent{}, err
	}

	suiChains := e.BlockChains.SuiChains()
	suiChain := suiChains[cfg.SourceChain]

	deps := suideps.SuiDeps{
		SuiChain: sui_ops.OpTxDeps{
			Client: suiChain.Client,
			Signer: suiChain.Signer,
			GetCallOpts: func() *suiBind.CallOpts {
				b := uint64(400_000_000)
				return &suiBind.CallOpts{
					Signer:           suiChain.Signer,
					WaitForExecution: true,
					GasBudget:        &b,
				}
			},
		},
	}

	ccipObjectRefId := state.SuiChains[cfg.SourceChain].CCIPObjectRef
	ccipPackageId := state.SuiChains[cfg.SourceChain].CCIPAddress
	onRampPackageId := state.SuiChains[cfg.SourceChain].OnRampAddress
	onRampStateObjectId := state.SuiChains[cfg.SourceChain].OnRampStateObjectId
	linkTokenPkgId := state.SuiChains[cfg.SourceChain].LinkTokenAddress
	linkTokenObjectMetadataId := state.SuiChains[cfg.SourceChain].LinkTokenCoinMetadataId
	ccipOwnerCapId := state.SuiChains[cfg.SourceChain].CCIPOwnerCapObjectId

	bigIntSourceUsdPerToken, ok := new(big.Int).SetString("150000000000000000000000000000", 10)
	if !ok {
		return &ccipclient.AnyMsgSentEvent{}, fmt.Errorf("failed converting SourceUSDPerToken to bigInt")
	}

	bigIntGasUsdPerUnitGas, ok := new(big.Int).SetString("7500000000000", 10)
	if !ok {
		return &ccipclient.AnyMsgSentEvent{}, fmt.Errorf("failed converting GasUsdPerUnitGas to bigInt")
	}

	// Update Prices on FeeQuoter with minted LinkToken
	_, err = operations.ExecuteOperation(e.OperationsBundle, ccipops.FeeQuoterUpdatePricesWithOwnerCapOp, deps.SuiChain,
		ccipops.FeeQuoterUpdatePricesWithOwnerCapInput{
			CCIPPackageId:         ccipPackageId,
			CCIPObjectRef:         ccipObjectRefId,
			OwnerCapObjectId:      ccipOwnerCapId,
			SourceTokens:          []string{linkTokenObjectMetadataId},
			SourceUsdPerToken:     []*big.Int{bigIntSourceUsdPerToken},
			GasDestChainSelectors: []uint64{cfg.DestChain},
			GasUsdPerUnitGas:      []*big.Int{bigIntGasUsdPerUnitGas},
		})
	if err != nil {
		return &ccipclient.AnyMsgSentEvent{}, fmt.Errorf("failed to updatePrice for Sui chain %d: %w", cfg.SourceChain, err)
	}

	ctx := context.Background()
	client := sui.NewSuiClient(suiChain.URL)
	ptb := suitx.NewTransaction()
	ptb.SetSuiClient(client.(*sui.Client))

	ccipStateHelperContract, err := suiBind.NewBoundContract(
		ccipPackageId,
		"ccip",
		"onramp_state_helper",
		client,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create onramp state helper bound contract when appending PTB command: %w", err)
	}

	// Note: these will be different for token transfers
	typeArgsList := []string{}
	typeParamsList := []string{}
	paramTypes := []string{
		"address",
	}
	paramValues := []any{
		"0xf05ebbc239612bdcfc6eff5f6f4728e87bc56d25e6f9dfcce9cffd6cc3eeb3ca", // random sui address
	}

	onRampCreateTokenTransferParamsCall, err := ccipStateHelperContract.EncodeCallArgsWithGenerics(
		"create_token_transfer_params",
		typeArgsList,
		typeParamsList,
		paramTypes,
		paramValues,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to encode onRampCreateTokenTransferParamsCall call: %w", err)
	}

	extractedAny2SuiMessageResult, err := ccipStateHelperContract.AppendPTB(ctx, deps.SuiChain.GetCallOpts(), ptb, onRampCreateTokenTransferParamsCall)
	if err != nil {
		return nil, fmt.Errorf("failed to build PTB (get_token_param_data) using bindings: %w", err)
	}

	onRampContract, err := suiBind.NewBoundContract(
		onRampPackageId,
		"ccip_onramp",
		"onramp",
		client,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create onramp bound contract when appending PTB command: %w", err)
	}

	// normalize module
	normalizedModule, err := client.SuiGetNormalizedMoveModule(ctx, models.GetNormalizedMoveModuleRequest{
		Package:    onRampPackageId,
		ModuleName: "onramp",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get normalized module: %w", err)
	}

	functionSignature, ok := normalizedModule.ExposedFunctions["ccip_send"]
	if !ok {
		return nil, fmt.Errorf("missing function signature for receiver function not found in module (%s)", "ccip_send")
	}

	// Figure out the parameter types from the normalized module of the token pool
	paramTypes, err = suiofframp_helper.DecodeParameters(e.Logger, functionSignature.(map[string]any), "parameters")
	if err != nil {
		return nil, fmt.Errorf("failed to decode parameters for token pool function: %w", err)
	}

	msg := cfg.Message.(SuiSendRequest)

	typeArgsList = []string{linkTokenPkgId + "::link::LINK"}
	typeParamsList = []string{}
	paramValues = []any{
		suiBind.Object{Id: ccipObjectRefId},
		suiBind.Object{Id: onRampStateObjectId},
		suiBind.Object{Id: "0x6"},
		cfg.DestChain,
		[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0xdd, 0xbb, 0x6f, 0x35,
			0x8f, 0x29, 0x04, 0x08, 0xd7, 0x68, 0x47, 0xb4,
			0xf6, 0x02, 0xf0, 0xfd, 0x59, 0x92, 0x95, 0xfd,
		}, // receiver
		[]byte("hello evm from sui"),
		extractedAny2SuiMessageResult,                 // tokenParams
		suiBind.Object{Id: linkTokenObjectMetadataId}, // feeTokenMetadata
		suiBind.Object{Id: msg.FeeToken},
		[]byte{}, // extraArgs
	}

	encodedOnRampCCIPSendCall, err := onRampContract.EncodeCallArgsWithGenerics(
		"ccip_send",
		typeArgsList,
		typeParamsList,
		paramTypes,
		paramValues,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to encode receiver call: %w", err)
	}

	_, err = onRampContract.AppendPTB(ctx, deps.SuiChain.GetCallOpts(), ptb, encodedOnRampCCIPSendCall)
	if err != nil {
		return nil, fmt.Errorf("failed to build PTB (receiver call) using bindings: %w", err)
	}

	executeCCIPSend, err := suiBind.ExecutePTB(ctx, deps.SuiChain.GetCallOpts(), client, ptb)
	if err != nil {
		return nil, fmt.Errorf("failed to execute ccip_send with err: %w", err)
	}

	if len(executeCCIPSend.Events) == 0 {
		return nil, fmt.Errorf("no events returned from Sui CCIPSend")
	}

	suiEvent := executeCCIPSend.Events[0].ParsedJson

	seqStr, _ := suiEvent["sequence_number"].(string)
	seq, _ := strconv.ParseUint(seqStr, 10, 64)

	return &ccipclient.AnyMsgSentEvent{
		SequenceNumber: seq,
		RawEvent:       suiEvent,
	}, nil
}

func MakeSuiExtraArgs(gasLimit uint64, allowOOO bool) []byte {
	var clockObj [32]byte
	copy(clockObj[:], hexutil.MustDecode(
		"0x0000000000000000000000000000000000000000000000000000000000000006",
	))

	var stateObj [32]byte
	copy(stateObj[:], hexutil.MustDecode(
		"0xffa55df38c762e3c4ac661af441d19da5bd2a1bfbe1d6329c24cc10b4bb119be", // reciever CCIPReceiverStateObjectId
	))

	recieverObjectIds := [][32]byte{clockObj, stateObj}

	extraArgs, err := ccipevm.SerializeClientSUIExtraArgsV1(message_hasher.ClientSuiExtraArgsV1{
		GasLimit:                 new(big.Int).SetUint64(gasLimit),
		AllowOutOfOrderExecution: allowOOO,
		TokenReceiver:            [32]byte{},
		// TokenReceiver: [32]byte{255, 165, 93, 243, 140, 118, 46, 60,
		// 	74, 198, 97, 175, 68, 29, 25, 218,
		// 	91, 210, 161, 191, 190, 29, 99, 41,
		// 	194, 76, 193, 11, 75, 177, 25, 190}, // ObjectID i.e. CCIPReceiverStateObjectId
		ReceiverObjectIds: recieverObjectIds,
	})
	if err != nil {
		panic(err)
	}
	return extraArgs
}
