package registrysyncer

import (
	"context"
	"errors"
	"fmt"

	"github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

type DonID uint32

type DON struct {
	capabilities.DON
	CapabilityConfigurations map[string]CapabilityConfiguration
}

type CapabilityConfiguration struct {
	// Common fields
	Config []byte

	// V1-specific fields
	CapabilityID *[32]byte `json:"capabilityId,omitempty"` // V1 uses [32]byte hash

	// V2-specific fields
	CapabilityIDString *string `json:"capabilityIdString,omitempty"` // V2 uses string capability ID

	// Version indicator
	Version string `json:"version"` // "v1" or "v2"
}

type Capability struct {
	ID             string
	CapabilityType capabilities.CapabilityType
}

type LocalRegistry struct {
	lggr              logger.Logger
	getPeerID         func() (types.PeerID, error)
	IDsToDONs         map[DonID]DON
	IDsToNodes        map[p2ptypes.PeerID]NodeInfo
	IDsToCapabilities map[string]Capability
}

func NewLocalRegistry(
	lggr logger.Logger,
	getPeerID func() (types.PeerID, error),
	IDsToDONs map[DonID]DON,
	IDsToNodes map[p2ptypes.PeerID]NodeInfo,
	IDsToCapabilities map[string]Capability,
) *LocalRegistry {
	return &LocalRegistry{
		lggr:              lggr,
		getPeerID:         getPeerID,
		IDsToDONs:         IDsToDONs,
		IDsToNodes:        IDsToNodes,
		IDsToCapabilities: IDsToCapabilities,
	}
}

func (l *LocalRegistry) LocalNode(ctx context.Context) (capabilities.Node, error) {
	pid, err := l.getPeerID()
	if err != nil {
		return capabilities.Node{}, fmt.Errorf("failed to get peer ID: %w", err)
	}

	return l.NodeByPeerID(ctx, p2ptypes.PeerID(pid))
}

func (l *LocalRegistry) NodeByPeerID(ctx context.Context, peerID types.PeerID) (capabilities.Node, error) {
	err := l.ensureNotEmpty()
	if err != nil {
		return capabilities.Node{}, err
	}
	nodeInfo, ok := l.IDsToNodes[peerID]
	if !ok {
		return capabilities.Node{}, errors.New("could not find peerID " + peerID.String())
	}

	var workflowDON capabilities.DON
	var capabilityDONs []capabilities.DON
	for _, d := range l.IDsToDONs {
		for _, p := range d.Members {
			if p == peerID {
				if d.AcceptsWorkflows {
					// The CapabilitiesRegistry enforces that the DON ID is strictly
					// greater than 0, so if the ID is 0, it means we've not set `workflowDON` initialized above yet.
					if workflowDON.ID == 0 {
						workflowDON = d.DON
						l.lggr.Debug("Workflow DON identified: %+v", workflowDON)
					} else {
						l.lggr.Errorf("Configuration error: node %s belongs to more than one workflowDON", peerID)
					}
				}

				capabilityDONs = append(capabilityDONs, d.DON)
			}
		}
	}

	return capabilities.Node{
		PeerID:              &peerID,
		NodeOperatorID:      nodeInfo.NodeOperatorID,
		Signer:              nodeInfo.Signer,
		EncryptionPublicKey: nodeInfo.EncryptionPublicKey,
		WorkflowDON:         workflowDON,
		CapabilityDONs:      capabilityDONs,
	}, nil
}

func (l *LocalRegistry) ConfigForCapability(ctx context.Context, capabilityID string, donID uint32) (CapabilityConfiguration, error) {
	err := l.ensureNotEmpty()
	if err != nil {
		return CapabilityConfiguration{}, err
	}
	d, ok := l.IDsToDONs[DonID(donID)]
	if !ok {
		return CapabilityConfiguration{}, fmt.Errorf("could not find don %d", donID)
	}

	cc, ok := d.CapabilityConfigurations[capabilityID]
	if !ok {
		return CapabilityConfiguration{}, fmt.Errorf("could not find capability configuration for capability %s and donID %d", capabilityID, donID)
	}

	return cc, nil
}

func (l *LocalRegistry) ensureNotEmpty() error {
	if len(l.IDsToDONs) == 0 {
		return errors.New("empty local registry. no DONs registered in the local registry")
	}
	if len(l.IDsToNodes) == 0 {
		return errors.New("empty local registry. no nodes registered in the local registry")
	}
	if len(l.IDsToCapabilities) == 0 {
		return errors.New("empty local registry. no capabilities registered in the local registry")
	}
	return nil
}

func (l *LocalRegistry) HasCapability(id string) (bool, error) {
	_, ok := l.IDsToCapabilities[id]
	return ok, nil
}

func (l *LocalRegistry) HasDON(id DonID) (bool, error) {
	_, ok := l.IDsToDONs[id]
	return ok, nil
}

func (l *LocalRegistry) CapabilityForID(id string) (capabilities.CapabilityInfo, error) {
	c, ok := l.IDsToCapabilities[id]
	if !ok {
		return capabilities.CapabilityInfo{}, fmt.Errorf("could not find capability with id %s", id)
	}

	return capabilities.CapabilityInfo{
		ID:             c.ID,
		CapabilityType: c.CapabilityType,
	}, nil
}

func (l *LocalRegistry) DONForID(id DonID) (capabilities.DON, error) {
	d, ok := l.IDsToDONs[id]
	if !ok {
		return capabilities.DON{}, fmt.Errorf("could not find DON with id %d", id)
	}

	return d.DON, nil
}

func (l *LocalRegistry) AllCapabilities() []capabilities.CapabilityInfo {
	capabilityInfos := make([]capabilities.CapabilityInfo, 0, len(l.IDsToCapabilities))
	for _, c := range l.IDsToCapabilities {
		capabilityInfos = append(capabilityInfos, capabilities.CapabilityInfo{
			ID:             c.ID,
			CapabilityType: c.CapabilityType,
		})
	}

	return capabilityInfos
}

func (l *LocalRegistry) AllDONs() []capabilities.DON {
	dons := make([]capabilities.DON, 0, len(l.IDsToDONs))
	for _, d := range l.IDsToDONs {
		dons = append(dons, d.DON)
	}

	return dons
}

func (l *LocalRegistry) AllNodes() []capabilities.Node {
	nodes := make([]capabilities.Node, 0, len(l.IDsToNodes))
	for peerID := range l.IDsToNodes {
		node, err := l.NodeByPeerID(context.Background(), peerID)
		if err != nil {
			l.lggr.Errorf("could not find node for peer %s: %s", peerID, err)
			continue
		}
		nodes = append(nodes, node)
	}

	return nodes
}
