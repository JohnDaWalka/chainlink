package por

import (
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/keystone/don/jobs"
	keystonenode "github.com/smartcontractkit/chainlink/system-tests/lib/keystone/don/node"
	keystoneflags "github.com/smartcontractkit/chainlink/system-tests/lib/keystone/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/keystone/types"
)

// If we wanted to by fancy we could also accept map[JobDescription]string that would get us the job spec
// if there's no job spec for the given JobDescription we would use the standard one, that could be easier
// than having to define the job spec for each JobDescription manually, in case someone wants to change one parameter
func Define(t *testing.T, ctfEnv *deployment.Environment, don *devenv.DON, nodeOutput *types.WrappedNodeOutput, bc *blockchain.Output, ocr3CapabilityAddress common.Address, donID uint32, flags []string, extraAllowedPorts []int, extraAllowedIps []string, cronCapBinName string, gatewayConnectorData types.GatewayConnectorData) types.DonJobs {
	donBootstrapNodePeerID, err := keystonenode.ToP2PID(don.Nodes[0], keystonenode.KeyExtractingTransformFn)
	require.NoError(t, err, "failed to get bootstrap node peer ID")

	donBootstrapNodeHost := nodeOutput.CLNodes[0].Node.ContainerName

	chainIDInt, err := strconv.Atoi(bc.ChainID)
	require.NoError(t, err, "failed to convert chain ID to int")
	chainIDUint64 := libc.MustSafeUint64(int64(chainIDInt))

	jobSpecs := make(types.DonJobs)

	// configuration of bootstrap node
	if keystoneflags.HasFlag(flags, types.OCR3Capability) {
		jobSpecs[types.JobDescription{Flag: types.OCR3Capability, NodeType: types.BootstrapNode}] = []*jobv1.ProposeJobRequest{jobs.BootstrapOCR3(don.Nodes[0].NodeID, ocr3CapabilityAddress, chainIDUint64)}
	}

	// if it's a workflow DON or it has custom compute capability, we need to create a gateway job
	if keystoneflags.HasFlag(flags, types.WorkflowDON) || keystoneflags.HasFlag(flags, types.CustomComputeCapability) {
		jobSpecs[types.JobDescription{Flag: types.WorkflowDON, NodeType: types.BootstrapNode}] = []*jobv1.ProposeJobRequest{jobs.BootstrapGateway(don, chainIDUint64, donID, extraAllowedPorts, extraAllowedIps, gatewayConnectorData)}
	}

	ocrPeeringData := types.OCRPeeringData{
		OCRBootstraperPeerID: donBootstrapNodePeerID,
		OCRBootstraperHost:   donBootstrapNodeHost,
		Port:                 5001,
	}

	// configuration of worker nodes
	for i, node := range don.Nodes {
		// First node is a bootstrap node, so we skip it
		if i == 0 {
			continue
		}

		if keystoneflags.HasFlag(flags, types.CronCapability) {
			jobSpec := jobs.WorkerStandardCapability(node.NodeID, "cron-capabilities", jobs.ExternalCapabilityPath(cronCapBinName), jobs.EmptyStdCapConfig)
			jobDesc := types.JobDescription{Flag: types.CronCapability, NodeType: types.WorkerNode}

			if _, ok := jobSpecs[jobDesc]; !ok {
				jobSpecs[jobDesc] = []*jobv1.ProposeJobRequest{jobSpec}
			} else {
				jobSpecs[jobDesc] = append(jobSpecs[jobDesc], jobSpec)
			}
		}

		if keystoneflags.HasFlag(flags, types.CustomComputeCapability) {
			config := `"""
				NumWorkers = 3
				[rateLimiter]
				globalRPS = 20.0
				globalBurst = 30
				perSenderRPS = 1.0
				perSenderBurst = 5
				"""`

			jobSpec := jobs.WorkerStandardCapability(node.NodeID, "custom-compute", "__builtin_custom-compute-action", config)
			jobDesc := types.JobDescription{Flag: types.CustomComputeCapability, NodeType: types.WorkerNode}

			if _, ok := jobSpecs[jobDesc]; !ok {
				jobSpecs[jobDesc] = []*jobv1.ProposeJobRequest{jobSpec}
			} else {
				jobSpecs[jobDesc] = append(jobSpecs[jobDesc], jobSpec)
			}
		}

		if keystoneflags.HasFlag(flags, types.OCR3Capability) {
			jobSpec := jobs.WorkerOCR3(node.NodeID, ocr3CapabilityAddress, common.HexToAddress(don.Nodes[i].AccountAddr[chainIDUint64]), node.Ocr2KeyBundleID, ocrPeeringData, chainIDUint64)
			jobDesc := types.JobDescription{Flag: types.OCR3Capability, NodeType: types.WorkerNode}

			if _, ok := jobSpecs[jobDesc]; !ok {
				jobSpecs[jobDesc] = []*jobv1.ProposeJobRequest{jobSpec}
			} else {
				jobSpecs[jobDesc] = append(jobSpecs[jobDesc], jobSpec)
			}
		}
	}

	return jobSpecs
}
