package registrysyncer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	p2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/capabilities/versioning"
	evmrelaytypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

// CapabilitiesRegistryReader defines the interface for reading from capabilities registry contracts
// across different versions. This interface abstracts the version-specific differences in the
// contract structure and provides a unified way to access capabilities, DONs, and nodes.
type CapabilitiesRegistryReader interface {
	// GetCapabilities returns all capabilities from the registry
	GetCapabilities(ctx context.Context) ([]CapabilityInfo, error)

	// GetDONs returns all DONs from the registry
	GetDONs(ctx context.Context) ([]DONInfo, error)

	// GetNodes returns all nodes from the registry
	GetNodes(ctx context.Context) ([]NodeInfo, error)

	// Address returns the contract address
	Address() common.Address

	// Close closes the reader and releases any resources
	Close() error
}

// CapabilityInfo represents capability information across all versions
// Version-specific fields are pointers and will be nil for versions that don't support them
type CapabilityInfo struct {
	// Common fields across all versions
	ID                    string
	LabelledName          string
	Version               string
	CapabilityType        uint8
	ResponseType          uint8
	ConfigurationContract common.Address
	IsDeprecated          bool

	// V1-specific fields
	HashedId *[32]byte `json:"hashedId,omitempty"` // Only populated for V1
}

// DONInfo represents DON information across all versions
// Version-specific fields are pointers and will be nil for versions that don't support them
type DONInfo struct {
	// Common fields across all versions
	ID                       uint32
	ConfigCount              uint32
	F                        uint8
	IsPublic                 bool
	AcceptsWorkflows         bool
	NodeP2PIds               []p2ptypes.PeerID
	CapabilityConfigurations []CapabilityConfiguration
}

// NodeInfo represents node information across all versions
// Version-specific fields are pointers and will be nil for versions that don't support them
type NodeInfo struct {
	// Common fields across all versions
	NodeOperatorID      uint32
	P2PID               p2ptypes.PeerID
	Signer              [32]byte
	EncryptionPublicKey [32]byte
	ConfigCount         uint32
	WorkflowDONId       uint32
	CapabilitiesDONIds  []uint32

	// V1-specific fields
	HashedCapabilityIds *[][32]byte `json:"hashedCapabilityIds,omitempty"` // V1 uses hashed IDs
}

// CapabilitiesRegistryReaderFactory creates version-specific readers
type CapabilitiesRegistryReaderFactory interface {
	NewCapabilitiesRegistryReader(
		ctx context.Context,
		relayer ContractReaderFactory,
		registryAddress string,
	) (CapabilitiesRegistryReader, error)
}

// capabilitiesRegistryReaderFactory implements CapabilitiesRegistryReaderFactory
type capabilitiesRegistryReaderFactory struct{}

// NewCapabilitiesRegistryReaderFactory creates a new factory instance
func NewCapabilitiesRegistryReaderFactory() CapabilitiesRegistryReaderFactory {
	return &capabilitiesRegistryReaderFactory{}
}

// NewCapabilitiesRegistryReader creates a version-specific reader by:
// 1. Detecting the contract version from the registry
// 2. Creating a version-specific contract reader configuration
// 3. Creating and binding the contract reader
// 4. Creating the appropriate version-specific registry reader
func (f *capabilitiesRegistryReaderFactory) NewCapabilitiesRegistryReader(
	ctx context.Context,
	relayer ContractReaderFactory,
	registryAddress string,
) (CapabilitiesRegistryReader, error) {
	contractAddress := common.HexToAddress(registryAddress)

	// Detect the contract version
	capabilitiesRegistryVersion, err := versioning.VerifyTypeAndVersion(
		ctx,
		registryAddress,
		relayer.NewContractReader,
		versioning.ContractType("CapabilitiesRegistry"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to detect contract version: %w", err)
	}

	// Create version-specific contract reader configuration
	var contractReaderConfig evmrelaytypes.ChainReaderConfig
	switch capabilitiesRegistryVersion.Major() {
	case 1:
		contractReaderConfig = buildV1ContractReaderConfig()
	case 2:
		contractReaderConfig = buildV2ContractReaderConfig()
	default:
		return nil, fmt.Errorf("unsupported capabilities registry version: %s", capabilitiesRegistryVersion.String())
	}

	// Create and configure the contract reader
	contractReaderConfigEncoded, err := json.Marshal(contractReaderConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal contract reader config: %w", err)
	}

	contractReader, err := relayer.NewContractReader(ctx, contractReaderConfigEncoded)
	if err != nil {
		return nil, fmt.Errorf("failed to create contract reader: %w", err)
	}

	// Create the bound contract
	capabilitiesContract := types.BoundContract{
		Address: registryAddress,
		Name:    "CapabilitiesRegistry",
	}

	err = contractReader.Bind(ctx, []types.BoundContract{capabilitiesContract})
	if err != nil {
		return nil, fmt.Errorf("failed to bind contract reader: %w", err)
	}

	// Create version-specific registry reader
	switch capabilitiesRegistryVersion.Major() {
	case 1:
		return NewCapabilitiesRegistryV1Reader(ctx, contractReader, contractAddress)
	case 2:
		return NewCapabilitiesRegistryV2Reader(ctx, contractReader, contractAddress)
	default:
		return nil, fmt.Errorf("unsupported capabilities registry version: %s", capabilitiesRegistryVersion.String())
	}
}

// buildV1ContractReaderConfig creates the contract reader configuration for V1 capabilities registry
func buildV1ContractReaderConfig() evmrelaytypes.ChainReaderConfig {
	return evmrelaytypes.ChainReaderConfig{
		Contracts: map[string]evmrelaytypes.ChainContractReader{
			"CapabilitiesRegistry": {
				ContractABI: kcr.CapabilitiesRegistryABI,
				Configs: map[string]*evmrelaytypes.ChainReaderDefinition{
					"getDONs": {
						ChainSpecificName: "getDONs",
					},
					"getCapabilities": {
						ChainSpecificName: "getCapabilities",
					},
					"getNodes": {
						ChainSpecificName: "getNodes",
					},
				},
			},
		},
	}
}

// buildV2ContractReaderConfig creates the contract reader configuration for V2 capabilities registry
func buildV2ContractReaderConfig() evmrelaytypes.ChainReaderConfig {
	// TODO: This will need to be updated with the actual V2 contract ABI
	// For now, we'll use the same structure as V1 but this will change
	// once the V2 contract bindings are available
	return evmrelaytypes.ChainReaderConfig{
		Contracts: map[string]evmrelaytypes.ChainContractReader{
			"CapabilitiesRegistry": {
				ContractABI: kcr.CapabilitiesRegistryABI, // TODO: Replace with V2 ABI
				Configs: map[string]*evmrelaytypes.ChainReaderDefinition{
					"getDONs": {
						ChainSpecificName: "getDONs",
					},
					"getCapabilities": {
						ChainSpecificName: "getCapabilities",
					},
					"getNodes": {
						ChainSpecificName: "getNodes",
					},
				},
			},
		},
	}
}
