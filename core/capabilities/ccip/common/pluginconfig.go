package common

import (
	"fmt"
	"sync"

	sel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/config/env"
)

// PluginConfig holds the configuration for a plugin.
type PluginConfig struct {
	CommitPluginCodec          cciptypes.CommitPluginCodec
	ExecutePluginCodec         cciptypes.ExecutePluginCodec
	MessageHasher              cciptypes.MessageHasher
	TokenDataEncoder           cciptypes.TokenDataEncoder
	GasEstimateProvider        cciptypes.EstimateProvider
	RMNCrypto                  cciptypes.RMNCrypto
	ContractTransmitterFactory cctypes.ContractTransmitterFactory
	// PriceOnlyCommitFn optional method override for price only commit reports.
	PriceOnlyCommitFn string
	ChainRW           ChainRWProvider
	AddressCodec      ChainSpecificAddressCodec
	ExtraDataCodec    SourceChainExtraDataCodec
}

// PluginServices aggregates services for multiple chain families (singleton registry).
type PluginServices struct {
	PluginConfigs              map[string]PluginConfig // chainFamily -> PluginConfig
	AddrCodec                  *AddressCodecRegistry   // Pointer to singleton registry
	ChainRW                    MultiChainRW
	LOOPPCCIPProviderSupported map[string]bool
	mu                         sync.RWMutex
}

// InitFunction defines a function to initialize a PluginConfig.
type InitFunction func(logger.Logger, cciptypes.ExtraDataCodecBundle) PluginConfig

var registeredFactories = make(map[string]InitFunction)

// Singleton instance
var (
	pluginServicesInstance *PluginServices
	pluginServicesOnce     sync.Once
)

// RegisterPluginConfig registers a plugin config factory for a chain family.
func RegisterPluginConfig(chainFamily string, factory InitFunction) {
	registeredFactories[chainFamily] = factory
}

// GetPluginServicesRegistry returns the singleton PluginServices registry
func GetPluginServicesRegistry(lggr logger.Logger) *PluginServices {
	pluginServicesOnce.Do(func() {
		pluginServicesInstance = &PluginServices{
			PluginConfigs:              make(map[string]PluginConfig),
			LOOPPCCIPProviderSupported: make(map[string]bool),
		}

		// Initialize plugin services from registered factories. This does not include
		// services provided by CCIPProvider objects, which get set later.
		pluginServicesInstance.initializeFromFactories(lggr)
	})
	return pluginServicesInstance
}

// initializeFromFactories initializes the PluginServices from registered factories from the
// import injections in oraclecreator/plugin.go
func (ps *PluginServices) initializeFromFactories(lggr logger.Logger) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	extraDataCodecRegistry := GetExtraDataCodecRegistry()
	addressCodecRegistry := GetAddressCodecRegistry()
	addressCodecMap := make(map[string]ChainSpecificAddressCodec)
	chainRWProviderMap := make(map[string]ChainRWProvider)
	looppSupported := make(map[string]bool)

	for family, initFunc := range registeredFactories {
		config := initFunc(lggr, extraDataCodecRegistry)
		ps.PluginConfigs[family] = config
		looppSupported[family] = isLOOPPEnabledForFamily(family)

		// Register codecs and providers
		extraDataCodecRegistry.RegisterFamilyNoopCodec(family)
		if config.ExtraDataCodec != nil {
			extraDataCodecRegistry.RegisterCodec(family, config.ExtraDataCodec)
		}
		if config.AddressCodec != nil {
			addressCodecMap[family] = config.AddressCodec
		}
		if config.ChainRW != nil {
			chainRWProviderMap[family] = config.ChainRW
		}
	}

	// Update all at once to avoid multiple locks
	addressCodecRegistry.RegisterAddressCodecs(addressCodecMap)
	ps.AddrCodec = addressCodecRegistry

	ps.ChainRW = NewCRCW(chainRWProviderMap)
	ps.LOOPPCCIPProviderSupported = looppSupported
}

// UpdateCodecsFromCCIPProviders updates plugin configs with codecs from CCIPProvider objects
func (ps *PluginServices) UpdateCodecsFromCCIPProviders(
	ccipProviders map[cciptypes.ChainSelector]types.CCIPProvider,
) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	// Track which address codecs get updated
	updatedAddressCodecs := make(map[string]ChainSpecificAddressCodec)

	// Start with existing address codecs from the singleton registry
	if ps.AddrCodec != nil {
		existingCodecs := ps.AddrCodec.GetRegisteredAddressCodecMap()
		for family, codec := range existingCodecs {
			updatedAddressCodecs[family] = codec
		}
	}

	for chainSelector, provider := range ccipProviders {
		chainFamily, err := sel.GetSelectorFamily(uint64(chainSelector))
		if err != nil {
			return fmt.Errorf("failed to get chain family from chain selector %d: %w", chainSelector, err)
		}

		// Get the existing config for this chain family
		existingConfig, exists := ps.PluginConfigs[chainFamily]
		if !exists {
			return fmt.Errorf("no plugin config found for chain family %s", chainFamily)
		}

		// Update plugin config with services from CCIPProvider
		codec := provider.Codec()
		updatedConfig := existingConfig

		if commitCodec := codec.CommitPluginCodec; commitCodec != nil {
			updatedConfig.CommitPluginCodec = commitCodec
		}

		if execCodec := codec.ExecutePluginCodec; execCodec != nil {
			updatedConfig.ExecutePluginCodec = execCodec
		}

		// TODO: add back in once MessageHasher is supported in CCIPProvider
		//if hasher := codec.MessageHasher; hasher != nil {
		//	updatedConfig.MessageHasher = hasher
		//}

		if addrCodec := codec.ChainSpecificAddressCodec; addrCodec != nil {
			updatedConfig.AddressCodec = addrCodec
			updatedAddressCodecs[chainFamily] = addrCodec
		}

		// Update the config in the registry
		ps.PluginConfigs[chainFamily] = updatedConfig

		// Also update the extra data codec registry
		if extraDataCodec := codec.SourceChainExtraDataCodec; extraDataCodec != nil {
			edcr := GetExtraDataCodecRegistry()
			edcr.RegisterCodec(chainFamily, extraDataCodec)
		}
	}

	// Update the singleton address codec registry directly
	// This ensures all components that reference the singleton see the updates immediately
	ps.AddrCodec.RegisterAddressCodecs(updatedAddressCodecs)

	return nil
}

// GetPluginConfigForFamily returns the plugin config for a specific chain family
func (ps *PluginServices) GetPluginConfigForFamily(chainFamily string) (PluginConfig, error) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	config, exists := ps.PluginConfigs[chainFamily]
	if !exists {
		return PluginConfig{}, fmt.Errorf("no plugin config found for chain family %s", chainFamily)
	}
	return config, nil
}

func isLOOPPEnabledForFamily(chainFamily string) bool {
	switch chainFamily {
	case sel.FamilySolana:
		return env.SolanaPlugin.Cmd.Get() != ""
	case sel.FamilyTon:
		return env.TONPlugin.Cmd.Get() != ""
	default:
		return false
	}
}
