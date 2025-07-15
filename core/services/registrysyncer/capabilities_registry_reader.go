package registrysyncer

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	p2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
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

// CapabilityInfo represents capability information in a version-agnostic way
type CapabilityInfo struct {
	ID                    string
	LabelledName          string
	Version               string
	CapabilityType        uint8
	ResponseType          uint8
	ConfigurationContract common.Address
	IsDeprecated          bool
}

// DONInfo represents DON information in a version-agnostic way
type DONInfo struct {
	ID                       uint32
	ConfigCount              uint32
	F                        uint8
	IsPublic                 bool
	AcceptsWorkflows         bool
	NodeP2PIds               [][32]byte
	CapabilityConfigurations []VersionedCapabilityConfiguration
}

// VersionedCapabilityConfiguration represents capability configuration with version-specific fields
type VersionedCapabilityConfiguration struct {
	CapabilityId string
	Config       []byte
}

// NodeInfo represents node information in a version-agnostic way
type NodeInfo struct {
	NodeOperatorID      uint32
	P2PID               p2ptypes.PeerID
	Signer              [32]byte
	EncryptionPublicKey [32]byte
	HashedCapabilityIds [][32]byte
	CapabilityIds       []string // V2 specific field
	ConfigCount         uint32
	WorkflowDONId       uint32
	CapabilitiesDONIds  []uint32
}

// CapabilityConfiguration represents capability configuration in a version-agnostic way
// This type is already defined in local_registry.go, so we reuse it for consistency

// CapabilitiesRegistryReaderFactory creates version-specific readers
type CapabilitiesRegistryReaderFactory interface {
	NewCapabilitiesRegistryReader(
		ctx context.Context,
		contractReader types.ContractReader,
		contractAddress common.Address,
		version string,
	) (CapabilitiesRegistryReader, error)
}

// capabilitiesRegistryReaderFactory implements CapabilitiesRegistryReaderFactory
type capabilitiesRegistryReaderFactory struct{}

// NewCapabilitiesRegistryReaderFactory creates a new factory instance
func NewCapabilitiesRegistryReaderFactory() CapabilitiesRegistryReaderFactory {
	return &capabilitiesRegistryReaderFactory{}
}

// NewCapabilitiesRegistryReader creates a version-specific reader based on the version string
func (f *capabilitiesRegistryReaderFactory) NewCapabilitiesRegistryReader(
	ctx context.Context,
	contractReader types.ContractReader,
	contractAddress common.Address,
	version string,
) (CapabilitiesRegistryReader, error) {
	switch version {
	case "1":
		return NewCapabilitiesRegistryV1Reader(ctx, contractReader, contractAddress)
	case "2":
		return NewCapabilitiesRegistryV2Reader(ctx, contractReader, contractAddress)
	default:
		return nil, errors.New("unsupported capabilities registry version: " + version)
	}
}
