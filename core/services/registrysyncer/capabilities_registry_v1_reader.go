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
		capConfigs := make([]VersionedCapabilityConfiguration, len(don.CapabilityConfigurations))
		for j, config := range don.CapabilityConfigurations {
			capConfigs[j] = VersionedCapabilityConfiguration{
				CapabilityId: fmt.Sprintf("%x", config.CapabilityId), // Convert bytes32 to hex string
				Config:       config.Config,
			}
		}

		result[i] = DONInfo{
			ID:                       don.Id,
			ConfigCount:              don.ConfigCount,
			F:                        don.F,
			IsPublic:                 don.IsPublic,
			AcceptsWorkflows:         don.AcceptsWorkflows,
			NodeP2PIds:               don.NodeP2PIds,
			CapabilityConfigurations: capConfigs,
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
		p2pId := p2ptypes.PeerID(node.P2pId)

		// Convert uint256 slice to uint32 slice
		capabilitiesDONIds := make([]uint32, len(node.CapabilitiesDONIds))
		for j, id := range node.CapabilitiesDONIds {
			capabilitiesDONIds[j] = uint32(id.Uint64())
		}

		result[i] = NodeInfo{
			NodeOperatorID:      node.NodeOperatorId,
			P2PID:               p2pId,
			Signer:              node.Signer,
			EncryptionPublicKey: node.EncryptionPublicKey,
			HashedCapabilityIds: node.HashedCapabilityIds,
			CapabilityIds:       nil, // V1 doesn't have string capability IDs
			ConfigCount:         node.ConfigCount,
			WorkflowDONId:       node.WorkflowDONId,
			CapabilitiesDONIds:  capabilitiesDONIds,
		}
	}

	return result, nil
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
