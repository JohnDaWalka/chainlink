package don

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/keystone/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/keystone/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/keystone/types"
)

func globalBootstraperNodeData(donTopologies []*types.DonWithMetadata) (string, string, error) {
	if len(donTopologies) == 1 {
		// if there is only one DON, then the global bootstrapper is the bootstrap node of the DON
		peerID, err := node.ToP2PID(donTopologies[0].DON.Nodes[0], node.KeyExtractingTransformFn)
		if err != nil {
			return "", "", errors.Wrapf(err, "failed to get peer ID for node %s", donTopologies[0].DON.Nodes[0].Name)
		}

		return peerID, donTopologies[0].NodeOutput.CLNodes[0].Node.ContainerName, nil
	} else if len(donTopologies) > 1 {
		// if there's more than one DON, then peering capabilitity needs to point to the same bootstrap node
		// for all the DONs, and so we need to find it first. For us, it will always be the bootstrap node of the workflow DON.
		for _, donTopology := range donTopologies {
			if flags.HasFlag(donTopology.Flags, types.WorkflowDON) {
				peerID, err := node.ToP2PID(donTopology.DON.Nodes[0], node.KeyExtractingTransformFn)
				if err != nil {
					return "", "", errors.Wrapf(err, "failed to get peer ID for node %s", donTopology.DON.Nodes[0].Name)
				}

				return peerID, donTopology.NodeOutput.CLNodes[0].Node.ContainerName, nil
			}
		}

		return "", "", errors.New("expected at least one workflow DON")
	}

	return "", "", errors.New("expected at least one DON topology")
}

func FindPeeringData(donTopologies []*types.DonWithMetadata) (types.PeeringData, error) {
	globalBootstraperPeerID, globalBootstraperHost, err := globalBootstraperNodeData(donTopologies)
	if err != nil {
		return types.PeeringData{}, err
	}

	return types.PeeringData{
		GlobalBootstraperPeerID: globalBootstraperPeerID,
		GlobalBootstraperHost:   globalBootstraperHost,
	}, nil
}
