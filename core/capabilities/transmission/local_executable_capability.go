package transmission

import (
	"context"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// LocalExecutableCapability handles the transmission protocol required for an executable capability that exists in the same node as
// the caller.
type LocalExecutableCapability struct {
	lggr logger.Logger
	capabilities.ExecutableCapability
	localNode    capabilities.Node
	capabilityID string
}

func NewLocalExecutableCapability(lggr logger.Logger, capabilityID string, localNode capabilities.Node, underlying capabilities.ExecutableCapability) *LocalExecutableCapability {
	return &LocalExecutableCapability{
		ExecutableCapability: underlying,
		capabilityID:         capabilityID,
		lggr:                 lggr,
		localNode:            localNode,
	}
}

func (l *LocalExecutableCapability) Execute(ctx context.Context, req capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
	if l.localNode.PeerID == nil || l.localNode.WorkflowDON.ID == 0 {
		l.lggr.Debugf("empty DON info, executing immediately")
		return l.ExecutableCapability.Execute(ctx, req)
	}

	if req.Config == nil || req.Config.Underlying["schedule"] == nil {
		l.lggr.Debug("no schedule found, executing immediately")
		return l.ExecutableCapability.Execute(ctx, req)
	}

	peerIDToTransmissionDelay, err := GetPeerIDToTransmissionDelay(l.localNode.WorkflowDON.Members, req)
	if err != nil {
		return capabilities.CapabilityResponse{}, fmt.Errorf("capability id: %s failed to get peer ID to transmission delay map: %w", l.capabilityID, err)
	}

	delay, existsForPeerID := peerIDToTransmissionDelay[*l.localNode.PeerID]
	if !existsForPeerID {
		return capabilities.CapabilityResponse{}, nil
	}

	select {
	case <-ctx.Done():
		return capabilities.CapabilityResponse{}, ctx.Err()
	case <-time.After(delay):
		response, err := l.ExecutableCapability.Execute(ctx, req)
		if err != nil {
			return response, err
		}

		// Set peer2peerID in the response metadata for local capabilities
		if len(response.Metadata.Metering) == 1 {
			response.Metadata.Metering[0].Peer2PeerID = l.localNode.PeerID.String()
		}

		return response, nil
	}
}
