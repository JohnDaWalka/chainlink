package registrysyncer

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
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

// V2-specific methods

// GetDONsInFamily returns DON IDs that belong to the specified family
func (r *capabilitiesRegistryV2Reader) GetDONsInFamily(ctx context.Context, donFamily string) ([]uint32, error) {
	var donIDsBig []big.Int

	err := r.contractReader.GetLatestValue(
		ctx,
		r.boundContract.ReadIdentifier("getDONsInFamily"),
		primitives.Unconfirmed,
		map[string]any{
			"donFamily": donFamily,
		},
		&donIDsBig,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get DONs in family %s: %w", donFamily, err)
	}

	// Convert []big.Int to []uint32
	donIDs := make([]uint32, len(donIDsBig))
	for i, bigID := range donIDsBig {
		if bigID.Cmp(big.NewInt(math.MaxUint32)) > 0 {
			return nil, fmt.Errorf("DON ID %s exceeds uint32 range", bigID.String())
		}
		donIDs[i] = uint32(bigID.Uint64()) // #nosec G115
	}

	return donIDs, nil
}

// GetHistoricalDONInfo returns historical DON information by DON ID and config count
func (r *capabilitiesRegistryV2Reader) GetHistoricalDONInfo(ctx context.Context, donID uint32, configCount uint32) (*DONInfo, error) {
	var don kcrv2.CapabilitiesRegistryDONInfo

	err := r.contractReader.GetLatestValue(
		ctx,
		r.boundContract.ReadIdentifier("getHistoricalDONInfo"),
		primitives.Unconfirmed,
		map[string]any{
			"donId":       donID,
			"configCount": configCount,
		},
		&don,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get historical DON info for DON %d, config %d: %w", donID, configCount, err)
	}

	// Convert node P2P IDs
	nodeP2PIDs := make([]p2ptypes.PeerID, len(don.NodeP2PIds))
	for i, nodeP2PID := range don.NodeP2PIds {
		nodeP2PIDs[i] = p2ptypes.PeerID(nodeP2PID)
	}

	// Convert capability configurations
	capabilities := make([]CapabilityConfiguration, len(don.CapabilityConfigurations))
	for i, cap := range don.CapabilityConfigurations {
		capabilities[i] = CapabilityConfiguration{
			CapabilityIDString: &cap.CapabilityId,
			Config:             cap.Config,
		}
	}

	// Convert DON families if present
	var donFamilies *[]string
	if len(don.DonFamilies) > 0 {
		families := make([]string, len(don.DonFamilies))
		copy(families, don.DonFamilies)
		donFamilies = &families
	}

	result := &DONInfo{
		ID:                       don.Id,
		ConfigCount:              don.ConfigCount,
		F:                        don.F,
		IsPublic:                 don.IsPublic,
		AcceptsWorkflows:         don.AcceptsWorkflows,
		NodeP2PIds:               nodeP2PIDs,
		CapabilityConfigurations: capabilities,
		// V2-specific fields
		Name:        &don.Name,
		Config:      &don.Config,
		DONFamilies: donFamilies,
		Version:     "v2",
	}

	return result, nil
}

// GetNode returns a single node by its P2P ID
func (r *capabilitiesRegistryV2Reader) GetNode(ctx context.Context, p2pID [32]byte) (*NodeInfo, error) {
	var node kcrv2.INodeInfoProviderNodeInfo

	err := r.contractReader.GetLatestValue(
		ctx,
		r.boundContract.ReadIdentifier("getNode"),
		primitives.Unconfirmed,
		map[string]any{
			"p2pId": p2pID,
		},
		&node,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get node with P2P ID %x: %w", p2pID, err)
	}

	// Convert capability IDs
	var capabilityIDs *[]string
	if len(node.CapabilityIds) > 0 {
		ids := make([]string, len(node.CapabilityIds))
		copy(ids, node.CapabilityIds)
		capabilityIDs = &ids
	}

	result := &NodeInfo{
		NodeOperatorID:      node.NodeOperatorId,
		P2PID:               p2ptypes.PeerID(node.P2pId),
		Signer:              node.Signer,
		EncryptionPublicKey: node.EncryptionPublicKey,
		ConfigCount:         node.ConfigCount,
		WorkflowDONId:       node.WorkflowDONId,
		CapabilitiesDONIds:  node.CapabilitiesDONIds,
		// V2-specific fields
		CapabilityIDs: capabilityIDs,
		Version:       "v2",
	}

	return result, nil
}

// GetNodeOperator returns node operator information by ID
func (r *capabilitiesRegistryV2Reader) GetNodeOperator(ctx context.Context, nodeOperatorID uint32) (*NodeOperator, error) {
	var nodeOp kcrv2.CapabilitiesRegistryNodeOperator

	err := r.contractReader.GetLatestValue(
		ctx,
		r.boundContract.ReadIdentifier("getNodeOperator"),
		primitives.Unconfirmed,
		map[string]any{
			"nodeOperatorId": nodeOperatorID,
		},
		&nodeOp,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get node operator %d: %w", nodeOperatorID, err)
	}

	result := &NodeOperator{
		Admin:   nodeOp.Admin,
		Name:    nodeOp.Name,
		Version: "v2",
	}

	return result, nil
}

// GetNodeOperators returns all node operators
func (r *capabilitiesRegistryV2Reader) GetNodeOperators(ctx context.Context) ([]NodeOperator, error) {
	var nodeOps []kcrv2.CapabilitiesRegistryNodeOperator

	err := r.contractReader.GetLatestValue(
		ctx,
		r.boundContract.ReadIdentifier("getNodeOperators"),
		primitives.Unconfirmed,
		nil,
		&nodeOps,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get node operators: %w", err)
	}

	result := make([]NodeOperator, len(nodeOps))
	for i, nodeOp := range nodeOps {
		result[i] = NodeOperator{
			Admin:   nodeOp.Admin,
			Name:    nodeOp.Name,
			Version: "v2",
		}
	}

	return result, nil
}

// GetNodesByP2PIds returns nodes filtered by the provided P2P IDs
func (r *capabilitiesRegistryV2Reader) GetNodesByP2PIds(ctx context.Context, p2pIDs [][32]byte) ([]NodeInfo, error) {
	var nodes []kcrv2.INodeInfoProviderNodeInfo

	err := r.contractReader.GetLatestValue(
		ctx,
		r.boundContract.ReadIdentifier("getNodesByP2PIds"),
		primitives.Unconfirmed,
		map[string]any{
			"p2pIds": p2pIDs,
		},
		&nodes,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes by P2P IDs: %w", err)
	}

	result := make([]NodeInfo, len(nodes))
	for i, node := range nodes {
		// Convert capability IDs
		var capabilityIDs *[]string
		if len(node.CapabilityIds) > 0 {
			ids := make([]string, len(node.CapabilityIds))
			copy(ids, node.CapabilityIds)
			capabilityIDs = &ids
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
			CapabilityIDs: capabilityIDs,
			Version:       "v2",
		}
	}

	return result, nil
}

// IsCapabilityDeprecated checks if a capability is marked as deprecated
func (r *capabilitiesRegistryV2Reader) IsCapabilityDeprecated(ctx context.Context, capabilityID string) (bool, error) {
	var isDeprecated bool

	err := r.contractReader.GetLatestValue(
		ctx,
		r.boundContract.ReadIdentifier("isCapabilityDeprecated"),
		primitives.Unconfirmed,
		map[string]any{
			"capabilityId": capabilityID,
		},
		&isDeprecated,
	)
	if err != nil {
		return false, fmt.Errorf("failed to check if capability %s is deprecated: %w", capabilityID, err)
	}

	return isDeprecated, nil
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
