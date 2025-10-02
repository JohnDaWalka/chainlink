package common

import (
	"fmt"
	"sync"

	chainsel "github.com/smartcontractkit/chain-selectors"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

var _ cciptypes.AddressCodec = (*AddressCodecRegistry)(nil)

// AddressCodecRegistry is a singleton registry that manages ChainSpecificAddressCodec instances
// for different chain families. It implements the cciptypes.AddressCodec interface.
type AddressCodecRegistry struct {
	registeredAddressCodecMap map[string]ChainSpecificAddressCodec
	mu                        sync.RWMutex
}

var (
	addressCodecRegistryInstance *AddressCodecRegistry
	addressCodecRegistryOnce     sync.Once
)

// GetAddressCodecRegistry returns the singleton instance of AddressCodecRegistry.
func GetAddressCodecRegistry() *AddressCodecRegistry {
	addressCodecRegistryOnce.Do(func() {
		addressCodecRegistryInstance = &AddressCodecRegistry{
			registeredAddressCodecMap: make(map[string]ChainSpecificAddressCodec),
		}
	})
	return addressCodecRegistryInstance
}

// RegisterAddressCodecs updates multiple address codecs at once.
func (r *AddressCodecRegistry) RegisterAddressCodecs(addressCodecMap map[string]ChainSpecificAddressCodec) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Update existing codecs with new ones
	for family, codec := range addressCodecMap {
		r.registeredAddressCodecMap[family] = codec
	}
}

// GetRegisteredAddressCodecMap returns a copy of the registered address codec map for backward compatibility.
func (r *AddressCodecRegistry) GetRegisteredAddressCodecMap() map[string]ChainSpecificAddressCodec {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Return a copy to prevent external modification
	result := make(map[string]ChainSpecificAddressCodec)
	for family, codec := range r.registeredAddressCodecMap {
		result[family] = codec
	}
	return result
}

// ============ Implementation of cciptypes.AddressCodec interface ============

// AddressBytesToString converts an address from bytes to string
func (r *AddressCodecRegistry) AddressBytesToString(addr cciptypes.UnknownAddress, chainSelector cciptypes.ChainSelector) (string, error) {
	family, err := chainsel.GetSelectorFamily(uint64(chainSelector))
	if err != nil {
		return "", fmt.Errorf("failed to get chain family for selector %d: %w", chainSelector, err)
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	codec, exist := r.registeredAddressCodecMap[family]

	if !exist {
		return "", fmt.Errorf("unsupported family for address decode type %s", family)
	}

	return codec.AddressBytesToString(addr)
}

// TransmitterBytesToString converts a transmitter account from bytes to string
func (r *AddressCodecRegistry) TransmitterBytesToString(addr cciptypes.UnknownAddress, chainSelector cciptypes.ChainSelector) (string, error) {
	family, err := chainsel.GetSelectorFamily(uint64(chainSelector))
	if err != nil {
		return "", fmt.Errorf("failed to get chain family for selector %d: %w", chainSelector, err)
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	codec, exist := r.registeredAddressCodecMap[family]

	if !exist {
		return "", fmt.Errorf("unsupported family for transmitter decode type %s", family)
	}

	return codec.TransmitterBytesToString(addr)
}

// AddressStringToBytes converts an address from string to bytes
func (r *AddressCodecRegistry) AddressStringToBytes(addr string, chainSelector cciptypes.ChainSelector) (cciptypes.UnknownAddress, error) {
	family, err := chainsel.GetSelectorFamily(uint64(chainSelector))
	if err != nil {
		return nil, fmt.Errorf("failed to get chain family for selector %d: %w", chainSelector, err)
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	codec, exist := r.registeredAddressCodecMap[family]

	if !exist {
		return nil, fmt.Errorf("unsupported family for address decode type %s", family)
	}

	return codec.AddressStringToBytes(addr)
}

// OracleIDAsAddressBytes returns valid address bytes for a given chain selector and oracle ID.
// Used for making nil transmitters in the OCR config valid, it just means that this oracle does not support the destination chain.
func (r *AddressCodecRegistry) OracleIDAsAddressBytes(oracleID uint8, chainSelector cciptypes.ChainSelector) ([]byte, error) {
	family, err := chainsel.GetSelectorFamily(uint64(chainSelector))
	if err != nil {
		return nil, fmt.Errorf("failed to get chain family for selector %d: %w", chainSelector, err)
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	codec, exist := r.registeredAddressCodecMap[family]

	if !exist {
		return nil, fmt.Errorf("unsupported family for address decode type %s", family)
	}

	return codec.OracleIDAsAddressBytes(oracleID)
}
