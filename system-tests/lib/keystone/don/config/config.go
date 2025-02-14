package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"

	"github.com/smartcontractkit/chainlink/system-tests/lib/keystone/types"
)

func Set(t *testing.T, nodeInput *types.CapabilitiesAwareNodeSet, bc *blockchain.Output) *types.WrappedNodeOutput {
	nodeset, err := ns.UpgradeNodeSet(t, nodeInput.Input, bc, 5*time.Second)
	require.NoError(t, err, "failed to upgrade node set")

	return &types.WrappedNodeOutput{Output: nodeset, NodeSetName: nodeInput.Name, Capabilities: nodeInput.Capabilities}
}
