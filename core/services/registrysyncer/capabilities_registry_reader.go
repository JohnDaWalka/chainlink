package registrysyncer

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	p2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	kcrv2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
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
	HashedID *[32]byte `json:"hashedId,omitempty"`

	// V2-specific fields
	Metadata *[]byte `json:"metadata,omitempty"`
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

	// V2-specific fields
	Name        *string   `json:"name,omitempty"`
	Config      *[]byte   `json:"config,omitempty"`
	DONFamilies *[]string `json:"donFamilies,omitempty"`

	// Version indicator
	Version string `json:"version"` // "v1" or "v2"
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
	CapabilitiesDONIds  []*big.Int

	// V1-specific fields
	HashedCapabilityIDs *[][32]byte `json:"hashedCapabilityIds,omitempty"` // V1 uses hashed IDs

	// V2-specific fields
	CapabilityIDs *[]string `json:"capabilityIds,omitempty"` // V2 uses string capability IDs

	// Version indicator
	Version string `json:"version"` // "v1" or "v2"
}

// Helper methods for NodeInfo
func (n *NodeInfo) IsV1() bool {
	return n.Version == "v1"
}

func (n *NodeInfo) IsV2() bool {
	return n.Version == "v2"
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
	return evmrelaytypes.ChainReaderConfig{
		Contracts: map[string]evmrelaytypes.ChainContractReader{
			"CapabilitiesRegistry": {
				ContractABI: kcrv2.CapabilitiesRegistryABI,
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
