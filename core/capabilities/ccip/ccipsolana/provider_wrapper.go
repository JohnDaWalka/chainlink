package ccipsolana

import (
	"context"
	"fmt"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-ccip/pkg/chainaccessor"
	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	"github.com/smartcontractkit/chainlink-ccip/pkg/contractreader"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ocrimpls"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/config/env"

	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

var _ ccipcommon.CCIPProviderWrapper = (*solanaCCIPProviderWrapper)(nil)

// solanaCCIPProviderWrapper wraps Solana chain services.
// Solana may or may not use LOOPP CCIPProvider depending on configuration.
type solanaCCIPProviderWrapper struct {
	chainFamily    string
	ccipProvider   types.CCIPProvider
	legacyServices *solanaLegacyServices
	extraDataCodec cciptypes.ExtraDataCodecBundle
}

// solanaBackwardsCompatibleCCIPProvider provides a backwards compatible CCIPProvider for use when LOOPP CCIPProvider
// is not yet supported by directly implemented ChainAccessor, ContractTransmitter, Codec methods.
type solanaBackwardsCompatibleCCIPProvider struct {
	codec               ccipocr3.Codec
	chainAccessor       cciptypes.ChainAccessor
	contractTransmitter ocr3types.ContractTransmitter[[]byte]
}

type solanaLegacyServices struct {
	gasEstimateProvider        cciptypes.EstimateProvider
	contractTransmitterFactory cctypes.ContractTransmitterFactory
	contractReader             contractreader.Extended
	contractWriter             types.ContractWriter
}

// NewSolanaCCIPProviderWrapper creates a new Solana CCIPProviderWrapper
func NewSolanaCCIPProviderWrapper(
	ctx context.Context,
	lggr logger.Logger,
	chainSelector cciptypes.ChainSelector,
	relayer loop.Relayer,
	cargs types.CCIPProviderArgs,
	largs ccipcommon.LegacyPluginServicesArgs,
	extraDataCodec cciptypes.ExtraDataCodecBundle,
) (ccipcommon.CCIPProviderWrapper, error) {
	lggr.Debugw("Creating new Solana CCIPProviderWrapper", "chainSelector", chainSelector)
	w := &solanaCCIPProviderWrapper{
		chainFamily:    chainsel.FamilySolana,
		extraDataCodec: extraDataCodec,
	}

	var ccipProvider types.CCIPProvider
	var err error
	if w.CCIPProviderSupported() {
		// LOOPP is configured for Solana, create CCIPProvider from relayer
		ccipProvider, err = relayer.NewCCIPProvider(ctx, cargs)
		if err != nil {
			return nil, err
		}
		w.ccipProvider = ccipProvider
	} else {
		// LOOPP not configured, use legacy services only
		chainRWProvider := ChainRWProvider{}
		cr, err := chainRWProvider.GetChainReader(ctx, largs.ChainReaderOpts)
		if err != nil {
			return nil, fmt.Errorf("error getting chain reader: %w", err)
		}
		ce, err := ccipcommon.WrapContractReaderInObservedExtended(lggr, cr, chainSelector)
		if err != nil {
			return nil, fmt.Errorf("wrapping contract reader in observed extended: %w", err)
		}
		cw, err := chainRWProvider.GetChainWriter(ctx, largs.ChainWriterOpts)
		if err != nil {
			return nil, err
		}
		da, err := chainaccessor.NewDefaultAccessor(
			logger.Sugared(lggr).Named(chainsel.FamilySolana).Named("DefaultAccessor"),
			chainSelector,
			ce,
			cw,
			nil, // TODO: use AddressCodecRegistry
		)

		transmitter, err := w.getTransmitterForPluginType(
			lggr,
			cargs.PluginType,
			chainSelector,
			largs,
			extraDataCodec,
			cw,
		)
		if err != nil {
			return nil, fmt.Errorf("getting transmitter for plugin type %d: %w", cargs.PluginType, err)
		}
		w.ccipProvider = &solanaBackwardsCompatibleCCIPProvider{
			codec: ccipocr3.Codec{
				ChainSpecificAddressCodec: AddressCodec{}, // TODO: maybe implement AddressCodecRegistry
				CommitPluginCodec:         NewCommitPluginCodecV1(),
				ExecutePluginCodec:        NewExecutePluginCodecV1(w.extraDataCodec),
				TokenDataEncoder:          NewSolanaTokenDataEncoder(),
				SourceChainExtraDataCodec: ExtraDataDecoder{},
				MessageHasher:             NewMessageHasherV1(logger.Sugared(lggr).Named(chainsel.FamilySolana).Named("MessageHasherV1"), extraDataCodec),
			},
			chainAccessor:       da,
			contractTransmitter: transmitter,
		}
		w.legacyServices = &solanaLegacyServices{
			gasEstimateProvider:        NewGasEstimateProvider(extraDataCodec),
			contractReader:             ce,
			contractWriter:             cw,
			contractTransmitterFactory: ocrimpls.NewSVMContractTransmitterFactory(extraDataCodec),
		}
	}

	return w, nil
}

func (w *solanaCCIPProviderWrapper) getTransmitterForPluginType(
	lggr logger.Logger,
	pluginType cciptypes.PluginType,
	chainSelector cciptypes.ChainSelector,
	largs ccipcommon.LegacyPluginServicesArgs,
	extraDataCodec cciptypes.ExtraDataCodecBundle,
	contractWriter types.ContractWriter,
) (ocr3types.ContractTransmitter[[]byte], error) {
	if chainSelector != largs.DestChainSelector {
		// Transmitter is only needed for dest chain
		return nil, nil
	}
	if len(largs.DestFromAccounts) == 0 {
		return nil, fmt.Errorf("transmitter array is empty for dest chain selector %d", largs.DestChainSelector)
	}
	transmitterFactory := ocrimpls.NewSVMContractTransmitterFactory(extraDataCodec)
	switch pluginType {
	case cciptypes.PluginTypeCCIPCommit:
		return transmitterFactory.NewCommitTransmitter(
			logger.Sugared(lggr).Named("CCIPCommitTransmitter").
				Named(largs.RelayID.String()).
				Named(fmt.Sprintf("%d", chainSelector)),
			contractWriter,
			ocrtypes.Account(largs.DestFromAccounts[0]),
			largs.OffRampAddressString,
			consts.MethodCommit,
			w.PriceOnlyCommitFn(),
		), nil
	case cciptypes.PluginTypeCCIPExec:
		return nil, nil // TODO: implement
	default:
		return nil, fmt.Errorf("unsupported plugin type: %d", pluginType)
	}
}

func (w *solanaCCIPProviderWrapper) CCIPProvider() types.CCIPProvider {
	return w.ccipProvider
}

func (p *solanaBackwardsCompatibleCCIPProvider) ChainAccessor() cciptypes.ChainAccessor {
	return p.chainAccessor
}

func (p *solanaBackwardsCompatibleCCIPProvider) ContractTransmitter() ocr3types.ContractTransmitter[[]byte] {
	return p.contractTransmitter
}

func (p *solanaBackwardsCompatibleCCIPProvider) Codec() cciptypes.Codec {
	return p.codec
}

func (w *solanaCCIPProviderWrapper) ChainFamily() string {
	return w.chainFamily
}

func (w *solanaCCIPProviderWrapper) CCIPProviderSupported() bool {
	return env.SolanaPlugin.Cmd.Get() != ""
}

func (w *solanaCCIPProviderWrapper) GasEstimateProvider() cciptypes.EstimateProvider {
	if w.legacyServices == nil {
		return nil
	}
	return w.legacyServices.gasEstimateProvider
}

func (w *solanaCCIPProviderWrapper) RMNCrypto() cciptypes.RMNCrypto {
	return nil // Solana doesn't support RMN
}

func (w *solanaCCIPProviderWrapper) ContractReader() contractreader.Extended {
	if w.legacyServices == nil {
		return nil
	}
	return w.legacyServices.contractReader
}

func (w *solanaCCIPProviderWrapper) ContractWriter() types.ContractWriter {
	if w.legacyServices == nil {
		return nil
	}
	return w.legacyServices.contractWriter
}

func (w *solanaCCIPProviderWrapper) PriceOnlyCommitFn() string {
	return consts.MethodCommitPriceOnly
}

func (w *solanaCCIPProviderWrapper) ContractTransmitterFactory() cctypes.ContractTransmitterFactory {
	if w.legacyServices == nil {
		return nil
	}
	return w.legacyServices.contractTransmitterFactory
}

func init() {
	// Register the Solana wrapper factory
	ccipcommon.RegisterCCIPProviderWrapperFactory(chainsel.FamilySolana, NewSolanaCCIPProviderWrapper)
}

// TODO: the following backwards compat methods are not used here, consider creating separate CCIPProviderServices
// interface that only contains ChainAccessor, ContractTransmitter, Codec methods and leaves the lifecycle methods out.

func (p *solanaBackwardsCompatibleCCIPProvider) Start(ctx context.Context) error {
	// TODO: start CR/CW here?
	return nil
}

func (p *solanaBackwardsCompatibleCCIPProvider) Close() error {
	// TODO: maybe close CR/CW here?
	return nil
}

func (p *solanaBackwardsCompatibleCCIPProvider) Ready() error {
	return nil
}

func (p *solanaBackwardsCompatibleCCIPProvider) HealthReport() map[string]error {
	return nil
}

func (p *solanaBackwardsCompatibleCCIPProvider) Name() string {
	return "SolanaBackwardsCompatibleCCIPProvider"
}
