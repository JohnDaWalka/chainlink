package ccipaptos

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/ccip_aptos_utils"
)

var aptosUtilsABI = abihelpers.MustParseABI(ccip_aptos_utils.AptosUtilsABI)

// ExecutePluginCodecV1 is a codec for encoding and decoding execute plugin reports.
// Compatible with ccip::offramp version 1.6.0
type ExecutePluginCodecV1 struct {
	extraDataCodec ccipcommon.ExtraDataCodec
}

func NewExecutePluginCodecV1(extraDataCodec ccipcommon.ExtraDataCodec) *ExecutePluginCodecV1 {
	return &ExecutePluginCodecV1{
		extraDataCodec: extraDataCodec,
	}
}

func (e *ExecutePluginCodecV1) Encode(ctx context.Context, report cciptypes.ExecutePluginReport) ([]byte, error) {
	aptosReport := make([]ccip_aptos_utils.AptosUtilsExecutionReport, 0, len(report.ChainReports))

	for _, chainReport := range report.ChainReports {
		if chainReport.ProofFlagBits.IsEmpty() {
			return nil, errors.New("proof flag bits are empty")
		}

		aptosProofs := make([][32]byte, 0, len(chainReport.Proofs))
		for _, proof := range chainReport.Proofs {
			aptosProofs = append(aptosProofs, proof)
		}

		if len(chainReport.Messages) != 1 {
			return nil, fmt.Errorf("only single report message expected, got %d", len(chainReport.Messages))
		}

		if len(chainReport.OffchainTokenData) != 1 {
			return nil, fmt.Errorf("only single group of offchain token data expected, got %d", len(chainReport.OffchainTokenData))
		}

		message := chainReport.Messages[0]
		offchainTokenData := chainReport.OffchainTokenData[0]

		receiver, err := addressBytesToBytes32(message.Receiver)
		if err != nil {
			return nil, fmt.Errorf("failed to convert receiver address to bytes32: %w", err)
		}

		tokenAmounts := make([]ccip_aptos_utils.AptosUtilsAny2AptosTokenTransfer, 0, len(message.TokenAmounts))
		for _, tokenAmount := range message.TokenAmounts {
			if tokenAmount.Amount.IsEmpty() {
				return nil, fmt.Errorf("empty amount for token: %s", tokenAmount.DestTokenAddress)
			}

			destExecDataDecodedMap, err2 := e.extraDataCodec.DecodeTokenAmountDestExecData(tokenAmount.DestExecData, chainReport.SourceChainSelector)
			if err2 != nil {
				return nil, fmt.Errorf("failed to decode dest exec data: %w", err)
			}

			destGasAmount, err2 := extractDestGasAmountFromMap(destExecDataDecodedMap)
			if err2 != nil {
				return nil, fmt.Errorf("decode dest gas amount: %w", err)
			}

			destTokenAddress, err2 := addressBytesToBytes32(tokenAmount.DestTokenAddress)
			if err2 != nil {
				return nil, fmt.Errorf("failed to convert token address to bytes32: %w", err)
			}

			tokenAmounts = append(tokenAmounts, ccip_aptos_utils.AptosUtilsAny2AptosTokenTransfer{
				SourcePoolAddress: tokenAmount.SourcePoolAddress,
				DestTokenAddress:  destTokenAddress,
				ExtraData:         tokenAmount.ExtraData,
				Amount:            tokenAmount.Amount.Int,
				DestGasAmount:     destGasAmount,
			})
		}

		decodedExtraArgsMap, err := e.extraDataCodec.DecodeExtraArgs(message.ExtraArgs, chainReport.SourceChainSelector)
		if err != nil {
			return nil, err
		}

		gasLimit, err := parseExtraDataMap(decodedExtraArgsMap)
		if err != nil {
			return nil, fmt.Errorf("decode extra args to get gas limit: %w", err)
		}

		aptosMessage := ccip_aptos_utils.AptosUtilsAny2AptosRampMessage{
			Header: ccip_aptos_utils.AptosUtilsRampMessageHeader{
				MessageId:           message.Header.MessageID,
				SourceChainSelector: uint64(message.Header.SourceChainSelector),
				DestChainSelector:   uint64(message.Header.DestChainSelector),
				SequenceNumber:      uint64(message.Header.SequenceNumber),
				Nonce:               message.Header.Nonce,
			},
			Sender:       common.LeftPadBytes(message.Sender, 32), // todo: make it chain-agnostic
			Data:         message.Data,
			Receiver:     receiver,
			GasLimit:     gasLimit,
			TokenAmounts: tokenAmounts,
		}

		aptosChainReport := ccip_aptos_utils.AptosUtilsExecutionReport{
			SourceChainSelector: uint64(chainReport.SourceChainSelector),
			Message:             aptosMessage,
			OffchainTokenData:   offchainTokenData,
			Proofs:              aptosProofs,
		}
		aptosReport = append(aptosReport, aptosChainReport)
	}

	packed, err := aptosUtilsABI.Pack("exposeExecutionReport", aptosReport)
	if err != nil {
		return nil, fmt.Errorf("failed to pack execution report: %w", err)
	}

	return packed[4:], nil
}

func (e *ExecutePluginCodecV1) Decode(ctx context.Context, encodedReport []byte) (cciptypes.ExecutePluginReport, error) {
	method, ok := aptosUtilsABI.Methods["exposeExecutionReport"]
	if !ok {
		return cciptypes.ExecutePluginReport{}, errors.New("missing method exposeExecutionReport")
	}

	unpacked, err := method.Inputs.Unpack(encodedReport)
	if err != nil {
		return cciptypes.ExecutePluginReport{}, fmt.Errorf("failed to unpack execution report: %w", err)
	}
	if len(unpacked) != 1 {
		return cciptypes.ExecutePluginReport{}, fmt.Errorf("expected 1 argument, got %d", len(unpacked))
	}

	aptosReportRaw := abi.ConvertType(unpacked[0], new([]ccip_aptos_utils.AptosUtilsExecutionReport))
	aptosReportPtr, is := aptosReportRaw.(*[]ccip_aptos_utils.AptosUtilsExecutionReport)
	if !is {
		return cciptypes.ExecutePluginReport{}, fmt.Errorf("got an unexpected report type %T", unpacked[0])
	}
	if aptosReportPtr == nil {
		return cciptypes.ExecutePluginReport{}, errors.New("aptos report is nil")
	}

	aptosReport := *aptosReportPtr
	executeReport := cciptypes.ExecutePluginReport{
		ChainReports: make([]cciptypes.ExecutePluginReportSingleChain, 0, len(aptosReport)),
	}

	for _, aptosChainReport := range aptosReport {
		proofs := make([]cciptypes.Bytes32, 0, len(aptosChainReport.Proofs))
		for _, proof := range aptosChainReport.Proofs {
			proofs = append(proofs, proof)
		}

		aptosMessage := aptosChainReport.Message

		tokenAmounts := make([]cciptypes.RampTokenAmount, 0, len(aptosMessage.TokenAmounts))
		for _, tokenAmount := range aptosMessage.TokenAmounts {
			destData, err := abiEncodeUint32(tokenAmount.DestGasAmount)
			if err != nil {
				return cciptypes.ExecutePluginReport{}, fmt.Errorf("abi encode dest gas amount: %w", err)
			}
			tokenAmounts = append(tokenAmounts, cciptypes.RampTokenAmount{
				SourcePoolAddress: tokenAmount.SourcePoolAddress,
				DestTokenAddress:  tokenAmount.DestTokenAddress[:],
				ExtraData:         tokenAmount.ExtraData,
				Amount:            cciptypes.NewBigInt(tokenAmount.Amount),
				DestExecData:      destData,
			})
		}

		message := cciptypes.Message{
			Header: cciptypes.RampMessageHeader{
				MessageID:           aptosMessage.Header.MessageId,
				SourceChainSelector: cciptypes.ChainSelector(aptosMessage.Header.SourceChainSelector),
				DestChainSelector:   cciptypes.ChainSelector(aptosMessage.Header.DestChainSelector),
				SequenceNumber:      cciptypes.SeqNum(aptosMessage.Header.SequenceNumber),
				Nonce:               aptosMessage.Header.Nonce,
				MsgHash:             cciptypes.Bytes32{},        // todo: info not available, but not required atm
				OnRamp:              cciptypes.UnknownAddress{}, // todo: info not available, but not required atm
			},
			Sender:         aptosMessage.Sender,
			Data:           aptosMessage.Data,
			Receiver:       aptosMessage.Receiver[:],
			ExtraArgs:      cciptypes.Bytes{},          // <-- todo: info not available, but not required atm
			FeeToken:       cciptypes.UnknownAddress{}, // <-- todo: info not available, but not required atm
			FeeTokenAmount: cciptypes.BigInt{},         // <-- todo: info not available, but not required atm
			TokenAmounts:   tokenAmounts,
		}

		chainReport := cciptypes.ExecutePluginReportSingleChain{
			SourceChainSelector: cciptypes.ChainSelector(aptosChainReport.SourceChainSelector),
			Messages:            []cciptypes.Message{message},
			OffchainTokenData:   [][][]byte{aptosChainReport.OffchainTokenData},
			Proofs:              proofs,
			ProofFlagBits:       cciptypes.NewBigInt(big.NewInt(0)),
		}

		executeReport.ChainReports = append(executeReport.ChainReports, chainReport)
	}

	return executeReport, nil
}

// Ensure ExecutePluginCodec implements the ExecutePluginCodec interface
var _ cciptypes.ExecutePluginCodec = (*ExecutePluginCodecV1)(nil)
