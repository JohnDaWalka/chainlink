package ccipevm

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	commoncodec "github.com/smartcontractkit/chainlink-common/pkg/codec"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/codec"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

const execReportABI = `[{"name":"chainReports","type":"tuple[]","components":[{"name":"sourceChainSelector","type":"uint64","internalType":"uint64"},{"name":"messages","type":"tuple[]","components":[{"name":"header","type":"tuple","components":[{"name":"messageID","type":"bytes32","internalType":"bytes32"},{"name":"sourceChainSelector","type":"uint64","internalType":"uint64"},{"name":"destChainSelector","type":"uint64","internalType":"uint64"},{"name":"sequenceNumber","type":"uint64","internalType":"uint64"},{"name":"nonce","type":"uint64","internalType":"uint64"},{"name":"msgHash","type":"bytes32","internalType":"bytes32"},{"name":"onRamp","type":"bytes","internalType":"bytes"}]},{"name":"data","type":"bytes","internalType":"bytes"},{"name":"sender","type":"bytes","internalType":"bytes"},{"name":"receiver","type":"bytes","internalType":"bytes"},{"name":"extraArgs","type":"bytes","internalType":"bytes"},{"name":"feeToken","type":"bytes","internalType":"bytes"},{"name":"feeTokenAmount","type":"uint256","internalType":"uint256"},{"name":"feeValueJuels","type":"uint256","internalType":"uint256"},{"name":"tokenAmounts","type":"tuple[]","components":[{"name":"sourcePoolAddress","type":"bytes","internalType":"bytes"},{"name":"destTokenAddress","type":"bytes","internalType":"bytes"},{"name":"extraData","type":"bytes","internalType":"bytes"},{"name":"amount","type":"uint256","internalType":"uint256"},{"name":"destExecData","type":"bytes","internalType":"bytes"}]}]},{"name":"offchainTokenData","type":"bytes[][]","internalType":"bytes[][]"},{"name":"proofs","type":"bytes32[]","internalType":"bytes32[]"},{"name":"proofFlagBits","type":"uint256","internalType":"uint256"}]}]`

var execCodecConfig = types.CodecConfig{
	Configs: map[string]types.ChainCodecConfig{
		"ExecPluginReport": {
			TypeABI: execReportABI,
			ModifierConfigs: commoncodec.ModifiersConfig{
				&commoncodec.WrapperModifierConfig{
					Fields: map[string]string{
						"ChainReports.Messages.FeeTokenAmount":      "Int",
						"ChainReports.Messages.FeeValueJuels":       "Int",
						"ChainReports.Messages.TokenAmounts.Amount": "Int",
						"ChainReports.ProofFlagBits":                "Int",
					},
				},
			},
		},
	},
}

// ExecutePluginCodecV2 is a codec for encoding and decoding execute plugin reports with generic codec
type ExecutePluginCodecV2 struct{}

func NewExecutePluginCodecV2() *ExecutePluginCodecV2 {
	return &ExecutePluginCodecV2{}
}

func validate(report cciptypes.ExecutePluginReport) error {
	for i, chainReport := range report.ChainReports {
		if chainReport.ProofFlagBits.IsEmpty() {
			return errors.New("proof flag bits are empty")
		}

		for j, message := range chainReport.Messages {
			// optional fields
			if message.FeeToken == nil {
				report.ChainReports[i].Messages[j].FeeToken = []byte{}
			}

			if message.FeeValueJuels.IsEmpty() {
				report.ChainReports[i].Messages[j].FeeValueJuels = cciptypes.NewBigInt(big.NewInt(0))
			}

			if message.FeeTokenAmount.IsEmpty() {
				report.ChainReports[i].Messages[j].FeeTokenAmount = cciptypes.NewBigInt(big.NewInt(0))
			}

			// required fields
			if message.Sender == nil {
				return errors.New("message sender is nil")
			}

			for _, tokenAmount := range message.TokenAmounts {
				if tokenAmount.Amount.IsEmpty() {
					return fmt.Errorf("empty amount for token: %s", tokenAmount.DestTokenAddress)
				}

				_, err := abiDecodeUint32(tokenAmount.DestExecData)
				if err != nil {
					return fmt.Errorf("decode dest gas amount: %w", err)
				}
			}

			_, err := decodeExtraArgsV1V2(message.ExtraArgs)
			if err != nil {
				return fmt.Errorf("decode extra args to get gas limit: %w", err)
			}
		}
	}

	return nil
}

func (e *ExecutePluginCodecV2) Encode(ctx context.Context, report cciptypes.ExecutePluginReport) ([]byte, error) {
	if err := validate(report); err != nil {
		return nil, err
	}

	cd, err := codec.NewCodec(execCodecConfig)
	if err != nil {
		return nil, err
	}

	return cd.Encode(ctx, report, "ExecPluginReport")
}

func execPostProcess(report *cciptypes.ExecutePluginReport) error {
	if len(report.ChainReports) == 0 {
		return errors.New("chain reports is empty")
	}

	for i, evmChainReport := range report.ChainReports {
		for j := range evmChainReport.Messages {
			report.ChainReports[i].Messages[j].Header.MsgHash = cciptypes.Bytes32{}
			report.ChainReports[i].Messages[j].Header.OnRamp = cciptypes.UnknownAddress{}
			report.ChainReports[i].Messages[j].ExtraArgs = cciptypes.Bytes{}
			report.ChainReports[i].Messages[j].FeeToken = cciptypes.UnknownAddress{}
			report.ChainReports[i].Messages[j].FeeTokenAmount = cciptypes.BigInt{}
		}
	}

	return nil
}

func (e *ExecutePluginCodecV2) Decode(ctx context.Context, encodedReport []byte) (cciptypes.ExecutePluginReport, error) {
	report := cciptypes.ExecutePluginReport{}
	cd, err := codec.NewCodec(execCodecConfig)
	if err != nil {
		return report, err
	}

	err = cd.Decode(ctx, encodedReport, &report, "ExecPluginReport")
	if err != nil {
		return report, err
	}

	if err = execPostProcess(&report); err != nil {
		return report, err
	}

	return report, err
}

// Ensure ExecutePluginCodec implements the ExecutePluginCodec interface
var _ cciptypes.ExecutePluginCodec = (*ExecutePluginCodecV2)(nil)
