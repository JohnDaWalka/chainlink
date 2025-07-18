package registrysyncer

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	p2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
)

// capabilitiesRegistryV1Reader implements CapabilitiesRegistryReader for V1 contracts
type capabilitiesRegistryV1Reader struct {
	contractReader types.ContractReader
	address        common.Address
	boundContract  types.BoundContract
}

// NewCapabilitiesRegistryV1Reader creates a new V1 capabilities registry reader
func NewCapabilitiesRegistryV1Reader(
	ctx context.Context,
	contractReader types.ContractReader,
	contractAddress common.Address,
) (CapabilitiesRegistryReader, error) {
	boundContract := types.BoundContract{
		Address: contractAddress.Hex(),
		Name:    "CapabilitiesRegistry",
	}

	return &capabilitiesRegistryV1Reader{
		contractReader: contractReader,
		address:        contractAddress,
		boundContract:  boundContract,
	}, nil
}

// GetCapabilities returns all capabilities from the V1 registry
func (r *capabilitiesRegistryV1Reader) GetCapabilities(ctx context.Context) ([]CapabilityInfo, error) {
	var caps []kcr.CapabilitiesRegistryCapabilityInfo

	err := r.contractReader.GetLatestValue(
		ctx,
		r.boundContract.ReadIdentifier("getCapabilities"),
		primitives.Unconfirmed,
		nil,
		&caps,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get capabilities from V1 registry: %w", err)
	}

	result := make([]CapabilityInfo, len(caps))
	for i, cap := range caps {
		result[i] = CapabilityInfo{
			ID:                    fmt.Sprintf("%s@%s", cap.LabelledName, cap.Version),
			LabelledName:          cap.LabelledName,
			Version:               cap.Version,
			CapabilityType:        cap.CapabilityType,
			ResponseType:          cap.ResponseType,
			ConfigurationContract: cap.ConfigurationContract,
			IsDeprecated:          cap.IsDeprecated,
			// V1-specific fields
			HashedID: &cap.HashedId,
		}
	}

	return result, nil
}

// GetDONs returns all DONs from the V1 registry
func (r *capabilitiesRegistryV1Reader) GetDONs(ctx context.Context) ([]DONInfo, error) {
	var dons []kcr.CapabilitiesRegistryDONInfo

	err := r.contractReader.GetLatestValue(
		ctx,
		r.boundContract.ReadIdentifier("getDONs"),
		primitives.Unconfirmed,
		nil,
		&dons,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get DONs from V1 registry: %w", err)
	}

	result := make([]DONInfo, len(dons))
	for i, don := range dons {
		capConfigs := make([]CapabilityConfiguration, len(don.CapabilityConfigurations))
		for j, config := range don.CapabilityConfigurations {
			capConfigID := config.CapabilityId // Store the [32]byte
			capConfigs[j] = CapabilityConfiguration{
				CapabilityID:       &capConfigID, // Set the pointer to the V1 field
				Config:             config.Config,
				CapabilityIDString: nil, // V2 field is nil
				Version:            "v1",
			}
		}

		// Convert P2P IDs from [32]byte to PeerID
		nodeP2PIds := make([]p2ptypes.PeerID, len(don.NodeP2PIds))
		for k, id := range don.NodeP2PIds {
			nodeP2PIds[k] = p2ptypes.PeerID(id)
		}

		result[i] = DONInfo{
			ID:                       don.Id,
			ConfigCount:              don.ConfigCount,
			F:                        don.F,
			IsPublic:                 don.IsPublic,
			AcceptsWorkflows:         don.AcceptsWorkflows,
			NodeP2PIds:               nodeP2PIds,
			CapabilityConfigurations: capConfigs,
			// V2-specific fields are nil for V1
			Name:        nil,
			Config:      nil,
			DONFamilies: nil,
			Version:     "v1",
		}
	}

	return result, nil
}

// GetNodes returns all nodes from the V1 registry
func (r *capabilitiesRegistryV1Reader) GetNodes(ctx context.Context) ([]NodeInfo, error) {
	var nodes []kcr.INodeInfoProviderNodeInfo

	err := r.contractReader.GetLatestValue(
		ctx,
		r.boundContract.ReadIdentifier("getNodes"),
		primitives.Unconfirmed,
		nil,
		&nodes,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes from V1 registry: %w", err)
	}

	result := make([]NodeInfo, len(nodes))
	for i, node := range nodes {
		// Convert p2pId from bytes32 to PeerID
		// In V1, the P2P ID is stored as [32]byte which can be directly converted to PeerID
		p2pID := p2ptypes.PeerID(node.P2pId)

		result[i] = NodeInfo{
			NodeOperatorID:      node.NodeOperatorId,
			P2PID:               p2pID,
			Signer:              node.Signer,
			EncryptionPublicKey: node.EncryptionPublicKey,
			ConfigCount:         node.ConfigCount,
			WorkflowDONId:       node.WorkflowDONId,
			CapabilitiesDONIds:  node.CapabilitiesDONIds,
			// V1-specific fields
			HashedCapabilityIDs: &node.HashedCapabilityIds,
			CapabilityIDs:       nil, // V2 data is nil
			Version:             "v1",
		}
	}

	return result, nil
}

// V2-only methods - return explicit errors for V1

// GetDONsInFamily returns an error as family operations are not supported in V1
func (r *capabilitiesRegistryV1Reader) GetDONsInFamily(ctx context.Context, donFamily string) ([]uint32, error) {
	return nil, fmt.Errorf("GetDONsInFamily: %w", ErrNotSupportedInV1)
}

// GetHistoricalDONInfo returns an error as historical data is not available in V1
func (r *capabilitiesRegistryV1Reader) GetHistoricalDONInfo(ctx context.Context, donID uint32, configCount uint32) (*DONInfo, error) {
	return nil, fmt.Errorf("GetHistoricalDONInfo: %w", ErrNotSupportedInV1)
}

// GetNode returns an error as single node lookup is not supported in V1
func (r *capabilitiesRegistryV1Reader) GetNode(ctx context.Context, p2pID [32]byte) (*NodeInfo, error) {
	return nil, fmt.Errorf("GetNode: %w", ErrNotSupportedInV1)
}

// GetNodeOperator returns an error as node operator information is not available in V1
func (r *capabilitiesRegistryV1Reader) GetNodeOperator(ctx context.Context, nodeOperatorID uint32) (*NodeOperator, error) {
	return nil, fmt.Errorf("GetNodeOperator: %w", ErrNotSupportedInV1)
}

// GetNodeOperators returns an error as node operator information is not available in V1
func (r *capabilitiesRegistryV1Reader) GetNodeOperators(ctx context.Context) ([]NodeOperator, error) {
	return nil, fmt.Errorf("GetNodeOperators: %w", ErrNotSupportedInV1)
}

// GetNodesByP2PIds returns an error as filtered node lookup is not supported in V1
func (r *capabilitiesRegistryV1Reader) GetNodesByP2PIds(ctx context.Context, p2pIDs [][32]byte) ([]NodeInfo, error) {
	return nil, fmt.Errorf("GetNodesByP2PIds: %w", ErrNotSupportedInV1)
}

// IsCapabilityDeprecated returns an error as deprecation checking is not supported in V1
func (r *capabilitiesRegistryV1Reader) IsCapabilityDeprecated(ctx context.Context, capabilityID string) (bool, error) {
	return false, fmt.Errorf("IsCapabilityDeprecated: %w", ErrNotSupportedInV1)
}

// Address returns the contract address
func (r *capabilitiesRegistryV1Reader) Address() common.Address {
	return r.address
}

// Close closes the reader and releases any resources
func (r *capabilitiesRegistryV1Reader) Close() error {
	// The contract reader is managed externally, so we don't close it here
	return nil
}
