package ocrimpls

import (
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

// SuiCommitCallArgs defines the calldata structure for an Sui commit transaction.
type SuiCommitCallArgs struct {
	CCIPObjectRef string
	OffRampState  string
	Clock         string
	ReportContext [2][32]byte `mapstructure:"ReportContext"`
	Report        []byte      `mapstructure:"Report"`
	Signatures    [][96]byte  `mapstructure:"Signatures"`
}

// SuiExecCallArgs defines the calldata structure for an Sui execute transaction.
type SuiExecCallArgs struct {
	ReportContext [2][32]byte `mapstructure:"ReportContext"`
	Report        []byte      `mapstructure:"Report"`
}

// SuiContractTransmitterFactory implements the transmitter factory for Sui chains.
type SuiContractTransmitterFactory struct {
	extraDataCodec ccipcommon.ExtraDataCodec
}

func NewSuiContractTransmitterFactory(extraDataCodec ccipcommon.ExtraDataCodec) *SuiContractTransmitterFactory {
	return &SuiContractTransmitterFactory{
		extraDataCodec: extraDataCodec,
	}
}

// NewSuiCommitCalldataFunc returns a ToCalldataFunc for Sui commits that omits any Info object.
func NewSuiCommitCalldataFunc(commitMethod string) ToEd25519CalldataFunc {
	return func(
		rawReportCtx [2][32]byte,
		report ocr3types.ReportWithInfo[[]byte],
		signatures [][96]byte,
		_ ccipcommon.ExtraDataCodec,
	) (string, string, any, error) {
		return consts.ContractNameOffRamp,
			commitMethod,
			SuiCommitCallArgs{
				CCIPObjectRef: "0xe11769597fb38575477f2e76b95e42964a82b448b5b254f9cbba8def4d860673",
				OffRampState:  "0xf47cdce5f66fbd0dccc9d5fe483293f549e6a10a07de80410f22b35380b354aa",
				Clock:         "0x6",
				ReportContext: rawReportCtx,
				Report:        report.Report,
				Signatures:    signatures,
			},
			nil
	}
}

// NewCommitTransmitter constructs an Sui commit transmitter.
func (f *SuiContractTransmitterFactory) NewCommitTransmitter(
	lggr logger.Logger,
	cw types.ContractWriter,
	fromAccount ocrtypes.Account,
	offrampAddress string,
	commitMethod, _ string, // priceOnlyMethod is ignored for Sui
) ocr3types.ContractTransmitter[[]byte] {
	return &ccipTransmitter{
		lggr:                lggr,
		cw:                  cw,
		fromAccount:         fromAccount,
		offrampAddress:      offrampAddress,
		toEd25519CalldataFn: NewSuiCommitCalldataFunc(commitMethod),
		extraDataCodec:      f.extraDataCodec,
	}
}

// SuiExecCallDataFunc builds the execute call data for Sui
var SuiExecCallDataFunc = func(
	rawReportCtx [2][32]byte,
	report ocr3types.ReportWithInfo[[]byte],
	signatures [][96]byte,
	_ ccipcommon.ExtraDataCodec,
) (contract string, method string, args any, err error) {
	return consts.ContractNameOffRamp,
		consts.MethodExecute,
		SuiExecCallArgs{
			ReportContext: rawReportCtx,
			Report:        report.Report,
		}, nil
}

// NewExecTransmitter constructs an Sui execute transmitter.
func (f *SuiContractTransmitterFactory) NewExecTransmitter(
	lggr logger.Logger,
	cw types.ContractWriter,
	fromAccount ocrtypes.Account,
	offrampAddress string,
) ocr3types.ContractTransmitter[[]byte] {
	return &ccipTransmitter{
		lggr:                lggr,
		cw:                  cw,
		fromAccount:         fromAccount,
		offrampAddress:      offrampAddress,
		toEd25519CalldataFn: SuiExecCallDataFunc,
		extraDataCodec:      f.extraDataCodec,
	}
}
