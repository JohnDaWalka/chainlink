package por

import (
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/keystone/don/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/keystone/don/node"
	keystoneflags "github.com/smartcontractkit/chainlink/system-tests/lib/keystone/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/keystone/types"
)

func Define(t *testing.T,
	don *devenv.DON,
	nodeInput *types.CapabilitiesAwareNodeSet,
	nodeOutput *types.WrappedNodeOutput,
	bc *blockchain.Output,
	donID uint32,
	flags []string,
	peeringData types.PeeringData,
	capRegAddr,
	workflowRegistryAddr,
	forwarderAddress common.Address,
	gatewayConnectorData *types.GatewayConnectorData,
) types.NodeIndexToConfigOverrides {
	// prepare required variables
	donBootstrapNodePeerID, err := node.ToP2PID(don.Nodes[0], node.KeyExtractingTransformFn)
	require.NoError(t, err, "failed to get bootstrap node peer ID")

	donBootstrapNodeHost := nodeOutput.CLNodes[0].Node.ContainerName

	chainIDInt, err := strconv.Atoi(bc.ChainID)
	require.NoError(t, err, "failed to convert chain ID to int")
	chainIDUint64 := libc.MustSafeUint64(int64(chainIDInt))

	configOverrides := make(types.NodeIndexToConfigOverrides)

	// bootstrap node
	configOverrides[0] = config.BootstrapEVM(donBootstrapNodePeerID, chainIDUint64, capRegAddr, bc.Nodes[0].DockerInternalHTTPUrl, bc.Nodes[0].DockerInternalWSUrl)

	if keystoneflags.HasFlag(flags, types.WorkflowDON) {
		configOverrides[0] += config.BoostrapDon2DonPeering(peeringData)

		require.NotNil(t, gatewayConnectorData.Host, "gatewayConnectorData.Host is required for Workflow DON")
		gatewayConnectorData.Host = donBootstrapNodeHost
	}

	// worker nodes
	workflowNodeSet := don.Nodes[1:]
	for i := range workflowNodeSet {
		configOverrides[i+1] = config.WorkerEVM(donBootstrapNodePeerID, donBootstrapNodeHost, peeringData, chainIDUint64, capRegAddr, bc.Nodes[0].DockerInternalHTTPUrl, bc.Nodes[0].DockerInternalWSUrl)
		nodeEthAddr := common.HexToAddress(workflowNodeSet[i].AccountAddr[chainIDUint64])

		if keystoneflags.HasFlag(flags, types.WriteEVMCapability) {
			configOverrides[i+1] += config.WorkerWriteEMV(
				nodeEthAddr,
				forwarderAddress,
			)
		}

		// if it's workflow DON configure workflow registry
		if keystoneflags.HasFlag(flags, types.WorkflowDON) {
			configOverrides[i+1] += config.WorkerWorkflowRegistry(
				workflowRegistryAddr, chainIDUint64)
		}

		// workflow DON nodes always needs gateway connector, otherwise they won't be able to fetch the workflow
		// it's also required by custom compute, which can only run on workflow DON nodes
		if keystoneflags.HasFlag(flags, types.WorkflowDON) || keystoneflags.HasFlag(flags, types.CustomComputeCapability) {
			configOverrides[i+1] += config.WorkerGateway(
				nodeEthAddr,
				chainIDUint64,
				donID,
				*gatewayConnectorData,
			)
		}
	}

	return configOverrides
}
