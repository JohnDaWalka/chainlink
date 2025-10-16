package ccipton

import (
	"context"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/pkg/contractreader"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipnoop"

	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"

	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

var _ ccipcommon.CCIPProviderWrapper = (*tonCCIPProviderWrapper)(nil)

// tonCCIPProviderWrapper wraps TON chain services.
type tonCCIPProviderWrapper struct {
	chainFamily    string
	ccipProvider   types.CCIPProvider
	extraDataCodec cciptypes.ExtraDataCodecBundle
}

// NewTonCCIPProviderWrapper creates a new TON CCIPProviderWrapper
func NewTonCCIPProviderWrapper(
	ctx context.Context,
	_ logger.Logger,
	_ cciptypes.ChainSelector,
	relayer loop.Relayer,
	cargs types.CCIPProviderArgs,
	_ ccipcommon.LegacyPluginServicesArgs,
	extraDataCodec cciptypes.ExtraDataCodecBundle,
) (ccipcommon.CCIPProviderWrapper, error) {
	w := &tonCCIPProviderWrapper{
		chainFamily:    chainsel.FamilyTon,
		extraDataCodec: extraDataCodec,
	}

	var ccipProvider types.CCIPProvider
	var err error
	if w.CCIPProviderSupported() {
		// LOOPP is configured for TON, create CCIPProvider from relayer
		ccipProvider, err = relayer.NewCCIPProvider(ctx, cargs)
		if err != nil {
			return nil, err
		}
		w.ccipProvider = ccipProvider
	}

	return w, nil
}

func (w *tonCCIPProviderWrapper) CCIPProvider() types.CCIPProvider {
	return w.ccipProvider
}

func (w *tonCCIPProviderWrapper) ChainFamily() string {
	return w.chainFamily
}

func (w *tonCCIPProviderWrapper) CCIPProviderSupported() bool {
	return true
}

func (w *tonCCIPProviderWrapper) GasEstimateProvider() cciptypes.EstimateProvider {
	return ccipnoop.NewGasEstimateProvider(w.extraDataCodec)
}

func (w *tonCCIPProviderWrapper) RMNCrypto() cciptypes.RMNCrypto {
	return nil
}

func (w *tonCCIPProviderWrapper) PriceOnlyCommitFn() string {
	return ""
}

func (w *tonCCIPProviderWrapper) ContractTransmitterFactory() cctypes.ContractTransmitterFactory {
	return nil
}

func (w *tonCCIPProviderWrapper) ContractReader() contractreader.Extended {
	return nil
}

func (w *tonCCIPProviderWrapper) ContractWriter() types.ContractWriter {
	return nil
}

func init() {
	// Register the TON wrapper factory
	ccipcommon.RegisterCCIPProviderWrapperFactory(chainsel.FamilyTon, NewTonCCIPProviderWrapper)
}
