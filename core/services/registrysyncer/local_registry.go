package registrysyncer

import (
	"context"
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

type DonID uint32

type DON struct {
	capabilities.DON
	CapabilityConfigurations map[string]CapabilityConfiguration
}

type CapabilityConfiguration struct {
	Config []byte
}

type Capability struct {
	ID             string
	CapabilityType capabilities.CapabilityType
}

type LocalRegistry struct {
	lggr              logger.Logger
	getPeerID         func() (p2ptypes.PeerID, error)
	IDsToDONs         map[DonID]DON
	IDsToNodes        map[p2ptypes.PeerID]kcr.INodeInfoProviderNodeInfo
	IDsToCapabilities map[string]Capability
}

type RegistryDiff struct {
	Added   LocalRegistry `json:"added"`
	Removed LocalRegistry `json:"removed"`
}

func (l *LocalRegistry) Diff(other *LocalRegistry) (RegistryDiff, bool) {
	added := LocalRegistry{
		IDsToDONs:         make(map[DonID]DON),
		IDsToNodes:        make(map[p2ptypes.PeerID]kcr.INodeInfoProviderNodeInfo),
		IDsToCapabilities: make(map[string]Capability),
	}
	removed := LocalRegistry{
		IDsToDONs:         make(map[DonID]DON),
		IDsToNodes:        make(map[p2ptypes.PeerID]kcr.INodeInfoProviderNodeInfo),
		IDsToCapabilities: make(map[string]Capability),
	}

	// Compare IDsToDONs
	for key, value := range l.IDsToDONs {
		if _, exists := other.IDsToDONs[key]; !exists {
			removed.IDsToDONs[key] = value
		}
	}
	for key, value := range other.IDsToDONs {
		if _, exists := l.IDsToDONs[key]; !exists {
			added.IDsToDONs[key] = value
		}
	}

	// Compare IDsToNodes
	for key, value := range l.IDsToNodes {
		if _, exists := other.IDsToNodes[key]; !exists {
			removed.IDsToNodes[key] = value
		}
	}
	for key, value := range other.IDsToNodes {
		if _, exists := l.IDsToNodes[key]; !exists {
			added.IDsToNodes[key] = value
		}
	}

	// Compare IDsToCapabilities
	for key, value := range l.IDsToCapabilities {
		if _, exists := other.IDsToCapabilities[key]; !exists {
			removed.IDsToCapabilities[key] = value
		}
	}
	for key, value := range other.IDsToCapabilities {
		if _, exists := l.IDsToCapabilities[key]; !exists {
			added.IDsToCapabilities[key] = value
		}
	}

	notEmpty := func(r LocalRegistry) bool {
		return len(r.IDsToDONs) > 0 || len(r.IDsToNodes) > 0 || len(r.IDsToCapabilities) > 0
	}
	return RegistryDiff{
		Added:   added,
		Removed: removed,
	}, (notEmpty(added) || notEmpty(removed))
}

func NewLocalRegistry(
	lggr logger.Logger,
	getPeerID func() (p2ptypes.PeerID, error),
	IDsToDONs map[DonID]DON,
	IDsToNodes map[p2ptypes.PeerID]kcr.INodeInfoProviderNodeInfo,
	IDsToCapabilities map[string]Capability,
) LocalRegistry {
	return LocalRegistry{
		lggr:              lggr.Named("LocalRegistry"),
		getPeerID:         getPeerID,
		IDsToDONs:         IDsToDONs,
		IDsToNodes:        IDsToNodes,
		IDsToCapabilities: IDsToCapabilities,
	}
}

func (l *LocalRegistry) LocalNode(ctx context.Context) (capabilities.Node, error) {
	// Load the current nodes PeerWrapper, this gets us the current node's
	// PeerID, allowing us to contextualize registry information in terms of DON ownership
	// (eg. get my current DON configuration, etc).
	pid, err := l.getPeerID()
	if err != nil {
		return capabilities.Node{}, errors.New("unable to get local node: peerWrapper hasn't started yet")
	}

	var workflowDON capabilities.DON
	capabilityDONs := []capabilities.DON{}
	for _, d := range l.IDsToDONs {
		for _, p := range d.Members {
			if p == pid {
				if d.AcceptsWorkflows {
					// The CapabilitiesRegistry enforces that the DON ID is strictly
					// greater than 0, so if the ID is 0, it means we've not set `workflowDON` initialized above yet.
					if workflowDON.ID == 0 {
						workflowDON = d.DON
						l.lggr.Debug("Workflow DON identified: %+v", workflowDON)
					} else {
						l.lggr.Errorf("Configuration error: node %s belongs to more than one workflowDON", pid)
					}
				}

				capabilityDONs = append(capabilityDONs, d.DON)
			}
		}
	}

	return capabilities.Node{
		PeerID:         &pid,
		WorkflowDON:    workflowDON,
		CapabilityDONs: capabilityDONs,
	}, nil
}

func (l *LocalRegistry) ConfigForCapability(ctx context.Context, capabilityID string, donID uint32) (CapabilityConfiguration, error) {
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
