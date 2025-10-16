package common

import (
	"context"

	"github.com/smartcontractkit/chainlink-ccip/pkg/contractreader"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
)

// CCIPProviderWrapper wraps both the new CCIPProvider (if available via LOOPP)
// and the legacy PluginServices components. This enables gradual migration from
// PluginServices to CCIPProvider while maintaining backward compatibility.
//
// Each chain family implements this interface to provide:
// - A CCIPProvider if the chain supports LOOPP (e.g., TON, Solana)
// - Legacy plugin services for components not yet in CCIPProvider
type CCIPProviderWrapper interface {
	// CCIPProvider returns the CCIPProvider instance for this chain. It will return the LOOPP
	// CCIPProvider if supported and configured, otherwise it will return a backwards-compatible
	// CCIPProvider that uses legacy services.
	CCIPProvider() types.CCIPProvider

	// ChainFamily returns the chain family identifier (e.g., "evm", "solana", "ton")
	ChainFamily() string

	// CCIPProviderSupported returns true if this chain family supports LOOPPed CCIPProvider.
	CCIPProviderSupported() bool

	LegacyPluginServices
}

// LegacyPluginServices provides access to plugin services that have not yet
// migrated to CCIPProvider. As components migrate to CCIPProvider, they should
// be removed from this interface.
type LegacyPluginServices interface {
	// GasEstimateProvider provides gas estimates for transaction execution
	GasEstimateProvider() cciptypes.EstimateProvider

	// RMNCrypto provides RMN (Risk Management Network) cryptographic operations
	// May return nil for chains that don't support RMN
	RMNCrypto() cciptypes.RMNCrypto

	// PriceOnlyCommitFn returns the method name for price-only commit functions
	// This is chain-specific (e.g., Solana uses a different method name)
	// Returns empty string if not applicable
	PriceOnlyCommitFn() string

	// ContractTransmitterFactory returns the contract transmitter factory
	// This is used to create commit and exec transmitters when CCIPProvider doesn't provide one
	ContractTransmitterFactory() cctypes.ContractTransmitterFactory

	// ContractReader returns an extended ContractReader. Not all chain families will implement this.
	ContractReader() contractreader.Extended

	// ContractWriter returns a ContractWriter. Not all chain families will implement this.
	ContractWriter() types.ContractWriter
}

type LegacyPluginServicesArgs struct {
	ChainReaderOpts      ChainReaderProviderOpts
	ChainWriterOpts      ChainWriterProviderOpts
	DestChainSelector    cciptypes.ChainSelector
	DestFromAccounts     []string
	OffRampAddressString string
	RelayID              types.RelayID
}

// CCIPProviderWrapperFactory is a function type that creates a CCIPProviderWrapper
// for a specific chain. Each chain family registers its factory.
type CCIPProviderWrapperFactory func(
	ctx context.Context,
	lggr logger.Logger,
	chainSelector cciptypes.ChainSelector,
	relayer loop.Relayer,
	cargs types.CCIPProviderArgs,
	largs LegacyPluginServicesArgs,
	extraDataCodec cciptypes.ExtraDataCodecBundle,
) (CCIPProviderWrapper, error)

var registeredWrapperFactories = make(map[string]CCIPProviderWrapperFactory)

// RegisterCCIPProviderWrapperFactory registers a wrapper factory for a chain family
func RegisterCCIPProviderWrapperFactory(chainFamily string, factory CCIPProviderWrapperFactory) {
	registeredWrapperFactories[chainFamily] = factory
}

// GetCCIPProviderWrapperFactory returns the registered factory for a chain family
func GetCCIPProviderWrapperFactory(chainFamily string) (CCIPProviderWrapperFactory, bool) {
	factory, exists := registeredWrapperFactories[chainFamily]
	return factory, exists
}
