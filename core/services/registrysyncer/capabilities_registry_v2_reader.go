package registrysyncer

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	p2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	kcrv2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
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

// V2CapabilityMetadata represents the metadata structure for V2 capabilities
type V2CapabilityMetadata struct {
	CapabilityType uint8 `json:"capabilityType"`
	ResponseType   uint8 `json:"responseType"`
}

// parseCapabilityID parses a V2 capability ID string in format "labelledName@version"
// and extracts the components. For V2, the capability types come from metadata.
func parseCapabilityID(capabilityID string) (id, labelledName, version string, err error) {
	// V2 capability IDs are in format "labelledName@version"
	// Example: "data-streams-reports@1.0.0"
	parts := strings.Split(capabilityID, "@")
	if len(parts) != 2 {
		return "", "", "", fmt.Errorf("invalid capability ID format: %s", capabilityID)
	}

	id = capabilityID
	labelledName = parts[0]
	version = parts[1]

	return id, labelledName, version, nil
}

// parseCapabilityMetadata extracts capability type and response type from V2 metadata
func parseCapabilityMetadata(metadata []byte) (capabilityType, responseType uint8, err error) {
	if len(metadata) == 0 {
		// Default values if no metadata
		return 0, 0, nil
	}

	var meta V2CapabilityMetadata
	if err := json.Unmarshal(metadata, &meta); err != nil {
		// If we can't parse the metadata, use default values
		return 0, 0, nil
	}

	return meta.CapabilityType, meta.ResponseType, nil
}

// GetCapabilities returns all capabilities from the V2 registry
func (r *capabilitiesRegistryV2Reader) GetCapabilities(ctx context.Context) ([]CapabilityInfo, error) {
	var capabilities []kcrv2.CapabilitiesRegistryCapabilityInfo

	err := r.contractReader.GetLatestValue(
		ctx,
		r.boundContract.ReadIdentifier("getCapabilities"),
		primitives.Unconfirmed,
		nil,
		&capabilities,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get capabilities from V2 registry: %w", err)
	}

	result := make([]CapabilityInfo, len(capabilities))
	for i, cap := range capabilities {
		// Parse capability ID parts for V2
		id, labelledName, version, err := parseCapabilityID(cap.CapabilityId)
		if err != nil {
			return nil, fmt.Errorf("failed to parse capability ID %s: %w", cap.CapabilityId, err)
		}

		// Extract capability type and response type from metadata
		capabilityType, responseType, err := parseCapabilityMetadata(cap.Metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to parse capability metadata for %s: %w", cap.CapabilityId, err)
		}

		// Convert metadata to pointer for V2 field
		var metadata *[]byte
		if cap.Metadata != nil {
			metadata = &cap.Metadata
		}

		result[i] = CapabilityInfo{
			ID:                    id,
			LabelledName:          labelledName,
			Version:               version,
			CapabilityType:        capabilityType,
			ResponseType:          responseType,
			ConfigurationContract: cap.ConfigurationContract,
			IsDeprecated:          cap.IsDeprecated,
			Metadata:              metadata, // V2-specific field
		}
	}

	return result, nil
}

// GetDONs returns all DONs from the V2 registry
func (r *capabilitiesRegistryV2Reader) GetDONs(ctx context.Context) ([]DONInfo, error) {
	var dons []kcrv2.CapabilitiesRegistryDONInfo

	err := r.contractReader.GetLatestValue(
		ctx,
		r.boundContract.ReadIdentifier("getDONs"),
		primitives.Unconfirmed,
		nil,
		&dons,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get DONs from V2 registry: %w", err)
	}

	result := make([]DONInfo, len(dons))
	for i, don := range dons {
		capConfigs := make([]CapabilityConfiguration, len(don.CapabilityConfigurations))
		for j, config := range don.CapabilityConfigurations {
			capConfigID := config.CapabilityId // Store the string
			capConfigs[j] = CapabilityConfiguration{
				CapabilityIDString: &capConfigID, // Set the pointer to the V2 field
				Config:             config.Config,
			}
		}

		// Convert P2P IDs from [32]byte to PeerID
		nodeP2PIds := make([]p2ptypes.PeerID, len(don.NodeP2PIds))
		for j, p2pID := range don.NodeP2PIds {
			nodeP2PIds[j] = p2ptypes.PeerID(p2pID)
		}

		// Convert V2-specific fields to pointers
		var name *string
		if don.Name != "" {
			name = &don.Name
		}

		var config *[]byte
		if don.Config != nil {
			config = &don.Config
		}

		var donFamilies *[]string
		if len(don.DonFamilies) > 0 {
			donFamilies = &don.DonFamilies
		}

		result[i] = DONInfo{
			ID:                       don.Id,
			ConfigCount:              don.ConfigCount,
			F:                        don.F,
			IsPublic:                 don.IsPublic,
			AcceptsWorkflows:         don.AcceptsWorkflows,
			NodeP2PIds:               nodeP2PIds,
			CapabilityConfigurations: capConfigs,
			// V2-specific fields
			Name:        name,
			Config:      config,
			DONFamilies: donFamilies,
		}
	}

	return result, nil
}

// GetNodes returns all nodes from the V2 registry
func (r *capabilitiesRegistryV2Reader) GetNodes(ctx context.Context) ([]NodeInfo, error) {
	var nodes []kcrv2.INodeInfoProviderNodeInfo

	err := r.contractReader.GetLatestValue(
		ctx,
		r.boundContract.ReadIdentifier("getNodes"),
		primitives.Unconfirmed,
		nil,
		&nodes,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes from V2 registry: %w", err)
	}

	result := make([]NodeInfo, len(nodes))
	for i, node := range nodes {
		// Convert V2-specific fields to pointers
		var capabilityIDs *[]string
		if len(node.CapabilityIds) > 0 {
			capabilityIDs = &node.CapabilityIds
		}

		result[i] = NodeInfo{
			NodeOperatorID:      node.NodeOperatorId,
			P2PID:               p2ptypes.PeerID(node.P2pId),
			Signer:              node.Signer,
			EncryptionPublicKey: node.EncryptionPublicKey,
			ConfigCount:         node.ConfigCount,
			WorkflowDONId:       node.WorkflowDONId,
			CapabilitiesDONIds:  node.CapabilitiesDONIds,
			// V2-specific fields
			CapabilityIDs: capabilityIDs, // V2 uses string capability IDs
			Version:       "v2",          // Set V2 version
		}
	}

	return result, nil
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
