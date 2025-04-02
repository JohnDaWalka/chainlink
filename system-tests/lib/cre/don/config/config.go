package config

import (
	"testing"
	"time"

	cldtypes "github.com/smartcontractkit/chainlink/deployment/environment/types"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

func Set(t *testing.T, nodeInput *types.CapabilitiesAwareNodeSet, bc *blockchain.Output) (*cldtypes.WrappedNodeOutput, error) {
	nodeset, err := ns.UpgradeNodeSet(t, nodeInput.Input, bc, 5*time.Second)
	if err != nil {
		return nil, errors.Wrap(err, "failed to upgrade node set")
	}

	return &cldtypes.WrappedNodeOutput{Output: nodeset, NodeSetName: nodeInput.Name, Capabilities: nodeInput.Capabilities}, nil
}
