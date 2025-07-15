package registrysyncer

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

// capabilitiesRegistryV2Reader implements CapabilitiesRegistryReader for V2 contracts
type capabilitiesRegistryV2Reader struct {
	contractReader types.ContractReader
	address        common.Address
	boundContract  types.BoundContract
}

// NewCapabilitiesRegistryV2Reader creates a new V2 capabilities registry reader
func NewCapabilitiesRegistryV2Reader(
	ctx context.Context,
	contractReader types.ContractReader,
	contractAddress common.Address,
) (CapabilitiesRegistryReader, error) {
	boundContract := types.BoundContract{
		Address: contractAddress.Hex(),
		Name:    "CapabilitiesRegistry",
	}

	return &capabilitiesRegistryV2Reader{
		contractReader: contractReader,
		address:        contractAddress,
		boundContract:  boundContract,
	}, nil
}

// GetCapabilities returns all capabilities from the V2 registry
func (r *capabilitiesRegistryV2Reader) GetCapabilities(ctx context.Context) ([]CapabilityInfo, error) {
	// TODO: Implement V2 capabilities reading
	// This will require V2 contract bindings and different field mappings
	return nil, errors.New("V2 capabilities reading not yet implemented")
}

// GetDONs returns all DONs from the V2 registry
func (r *capabilitiesRegistryV2Reader) GetDONs(ctx context.Context) ([]DONInfo, error) {
	// TODO: Implement V2 DON reading
	// This will require V2 contract bindings and different field mappings
	return nil, errors.New("V2 DON reading not yet implemented")
}

// GetNodes returns all nodes from the V2 registry
func (r *capabilitiesRegistryV2Reader) GetNodes(ctx context.Context) ([]NodeInfo, error) {
	// TODO: Implement V2 node reading
	// This will require V2 contract bindings and different field mappings
	// V2 uses string[] capabilityIds instead of bytes32[] hashedCapabilityIds
	return nil, errors.New("V2 node reading not yet implemented")
}

// Address returns the contract address
func (r *capabilitiesRegistryV2Reader) Address() common.Address {
	return r.address
}

// Close closes the reader and releases any resources
func (r *capabilitiesRegistryV2Reader) Close() error {
	// The contract reader is managed externally, so we don't close it here
	return nil
}
