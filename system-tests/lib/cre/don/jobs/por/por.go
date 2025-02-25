package por

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	creflags "github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

func GenerateAllJobSpecs(input *types.GeneratePoRJobSpecsInputs) (types.DonsToJobSpecs, error) {
	if input == nil {
		return nil, errors.New("input is nil")
	}
	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "input validation failed")
	}
	donToJobSpecs := make(types.DonsToJobSpecs)

	gatewayConnectorData := input.GatewayConnectorOutput
	// DON_LOOP:
	// 	for _, donWithMetadata := range input.DonsWithMetadata {
	// 		if creflags.HasFlag(donWithMetadata.DonMetadata.Flags, types.GatewayDON) {
	// 			for _, nodeMetadata := range donWithMetadata.DonMetadata.NodesMetadata {
	// 				isGatewayNode := false
	// 			GATEWAY_LOOP:
	// 				for _, label := range nodeMetadata.Labels {
	// 					// check if the node is a gateway node
	// 					if label.Key == node.ExtraRolesKey && node.LabelContains(label.Value, ptr.Ptr(types.GatewayNode)) {
	// 						isGatewayNode = true
	// 						break GATEWAY_LOOP
	// 					}
	// 				}

	// 				// we support only ONE gateway node, even though multiple are supported by the CRE
	// 				if isGatewayNode {
	// 					for _, label := range nodeMetadata.Labels {
	// 						if label.Key == node.HostLabelKey {
	// 							gatewayConnectorData.Host = *label.Value
	// 							break DON_LOOP
	// 						}
	// 					}

	// 					return nil, errors.New("failed to get gateway node host from labels")
	// 				}
	// 			}
	// 		}
	// 	}

	// 	if gatewayConnectorData.Host == "" {
	// 		return nil, errors.New("failed to find gateway node")
	// 	}

	for _, donWithMetadata := range input.DonsWithMetadata {
		// if it's a workflow DON or it has custom compute capability, it needs access to gateway connector
		if creflags.HasFlag(donWithMetadata.Flags, types.WorkflowDON) || creflags.HasFlag(donWithMetadata.Flags, types.CustomComputeCapability) {
			workflowNodeSet, err := node.FindManyWithLabel(donWithMetadata.NodesMetadata, &ptypes.Label{Key: devenv.NodeLabelKeyType, Value: ptr.Ptr(string(devenv.NodeLabelValuePlugin))}, node.EqualLabels)
			if err != nil {
				return nil, errors.Wrap(err, "failed to find worker nodes")
			}

			ethAddresses := make([]string, len(workflowNodeSet))
			for i, n := range workflowNodeSet {
				for _, label := range n.Labels {
					if label.Key == node.EthAddressKey {
						if label.Value == nil {
							return nil, errors.New("eth address label value is nil")
						}
						if *label.Value == "" {
							return nil, errors.New("eth address label value is empty")
						}
						ethAddresses[i] = *label.Value
						break
					}
				}
			}
			gatewayConnectorData.Dons = append(gatewayConnectorData.Dons, types.GatewayConnectorDons{
				MembersEthAddresses: ethAddresses,
				ID:                  donWithMetadata.DonMetadata.ID,
			})
		}
	}

	for _, donWithMetadata := range input.DonsWithMetadata {
		jobSpecs, err := generateDonJobSpecs(types.GeneratePoRJobSpecsInput{
			CldEnv:                 input.CldEnv,
			DonWithMetadata:        *donWithMetadata,
			BlockchainOutput:       input.BlockchainOutput,
			OCR3CapabilityAddress:  input.OCR3CapabilityAddress,
			CronCapBinName:         input.CronCapBinName,
			ExtraAllowedPorts:      input.ExtraAllowedPorts,
			ExtraAllowedIPs:        input.ExtraAllowedIPs,
			GatewayConnectorOutput: gatewayConnectorData,
		})
		if err != nil {
			return nil, errors.Wrapf(err, "failed to generate job specs for don %d", donWithMetadata.DonMetadata.ID)
		}

		donToJobSpecs[donWithMetadata.DonMetadata.ID] = jobSpecs
	}

	return donToJobSpecs, nil
}

// If we wanted to by fancy we could also accept map[JobDescription]string that would get us the job spec
// if there's no job spec for the given JobDescription we would use the standard one, that could be easier
// than having to define the job spec for each JobDescription manually, in case someone wants to change one parameter
func generateDonJobSpecs(input types.GeneratePoRJobSpecsInput) (types.DonJobs, error) {
	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "input validation failed")
	}
	jobSpecs := make(types.DonJobs)

	chainIDInt, err := strconv.Atoi(input.BlockchainOutput.ChainID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert chain ID to int")
	}
	chainIDUint64 := libc.MustSafeUint64(int64(chainIDInt))

	bootstrapNode, err := node.FindOneWithLabel(input.DonWithMetadata.NodesMetadata, &ptypes.Label{Key: devenv.NodeLabelKeyType, Value: ptr.Ptr(string(devenv.NodeLabelValueBootstrap))}, node.EqualLabels)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find bootstrap node")
	}

	donBootstrapNodePeerID, err := node.ToP2PID(bootstrapNode, node.KeyExtractingTransformFn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get bootstrap node peer ID")
	}

	var donBootstrapNodeHost string
	for _, label := range bootstrapNode.Labels {
		if label.Key == node.HostLabelKey {
			donBootstrapNodeHost = *label.Value
			break
		}
	}

	if donBootstrapNodeHost == "" {
		return nil, errors.New("failed to get bootstrap node host from labels")
	}

	var bootstrapNodeID string
	for _, label := range bootstrapNode.Labels {
		if label.Key == devenv.NodeIDKeyType {
			bootstrapNodeID = *label.Value
			break
		}
	}

	if bootstrapNodeID == "" {
		return nil, errors.New("failed to get bootstrap node id from labels")
	}

	workflowNodeSet, err := node.FindManyWithLabel(input.DonWithMetadata.NodesMetadata, &ptypes.Label{Key: devenv.NodeLabelKeyType, Value: ptr.Ptr(string(devenv.NodeLabelValuePlugin))}, node.EqualLabels)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find worker nodes")
	}

	// configuration of bootstrap node
	if creflags.HasFlag(input.DonWithMetadata.Flags, types.OCR3Capability) {
		jobSpecs[types.JobDescription{Flag: types.OCR3Capability, NodeType: types.BootstrapNode}] = []*jobv1.ProposeJobRequest{jobs.BootstrapOCR3(bootstrapNodeID, input.OCR3CapabilityAddress, chainIDUint64)}
	}

	// TODO this is an oversimplification, we should create this job on each node that's also a gateway node, and eth addresses should include all nodes that might need this gateway,
	// and these nodes would be for sure all worker nodes, but potentially from more than 1 DON

	// if it's a workflow DON or it has custom compute capability, we need to create a gateway job
	// if creflags.HasFlag(input.DonWithMetadata.Flags, types.WorkflowDON) || creflags.HasFlag(input.DonWithMetadata.Flags, types.CustomComputeCapability) {
	// ethAddresses := make([]string, len(workflowNodeSet))
	// for i, n := range workflowNodeSet {
	// 	for _, label := range n.Labels {
	// 		if label.Key == node.EthAddressKey {
	// 			if label.Value == nil {
	// 				return nil, errors.New("eth address label value is nil")
	// 			}
	// 			if *label.Value == "" {
	// 				return nil, errors.New("eth address label value is empty")
	// 			}
	// 			ethAddresses[i] = *label.Value
	// 			break
	// 		}
	// 	}
	// }

	if creflags.HasFlag(input.DonWithMetadata.Flags, types.GatewayDON) {
		gatewayNode, err := node.FindOneWithLabel(input.DonWithMetadata.NodesMetadata, &ptypes.Label{Key: node.ExtraRolesKey, Value: ptr.Ptr(string(types.GatewayNode))}, node.LabelContains)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find bootstrap node")
		}

		var gatewayNodeID string
		for _, label := range gatewayNode.Labels {
			if label.Key == devenv.NodeIDKeyType {
				gatewayNodeID = *label.Value
				break
			}
		}

		if gatewayNodeID == "" {
			return nil, errors.New("failed to get gateway node id from labels")
		}

		jobSpecs[types.JobDescription{Flag: types.WorkflowDON, NodeType: types.BootstrapNode}] = []*jobv1.ProposeJobRequest{jobs.Gateway(gatewayNodeID, chainIDUint64, input.DonWithMetadata.ID, input.ExtraAllowedPorts, input.ExtraAllowedIPs, input.GatewayConnectorOutput)}
	}

	ocrPeeringData := types.OCRPeeringData{
		OCRBootstraperPeerID: donBootstrapNodePeerID,
		OCRBootstraperHost:   donBootstrapNodeHost,
		Port:                 5001,
	}

	// configuration of worker nodes
	for _, workerNode := range workflowNodeSet {
		var nodeID string
		for _, label := range workerNode.Labels {
			if label.Key == devenv.NodeIDKeyType {
				nodeID = *label.Value
				break
			}
		}

		if nodeID == "" {
			return nil, errors.New("failed to get node id from labels")
		}

		if creflags.HasFlag(input.DonWithMetadata.Flags, types.CronCapability) {
			jobSpec := jobs.WorkerStandardCapability(nodeID, "cron-capability", jobs.ExternalCapabilityPath(input.CronCapBinName), jobs.EmptyStdCapConfig)
			jobDesc := types.JobDescription{Flag: types.CronCapability, NodeType: types.WorkerNode}

			if _, ok := jobSpecs[jobDesc]; !ok {
				jobSpecs[jobDesc] = []*jobv1.ProposeJobRequest{jobSpec}
			} else {
				jobSpecs[jobDesc] = append(jobSpecs[jobDesc], jobSpec)
			}
		}

		if creflags.HasFlag(input.DonWithMetadata.Flags, types.CustomComputeCapability) {
			config := `"""
				NumWorkers = 3
				[rateLimiter]
				globalRPS = 20.0
				globalBurst = 30
				perSenderRPS = 1.0
				perSenderBurst = 5
				"""`

			jobSpec := jobs.WorkerStandardCapability(nodeID, "custom-compute", "__builtin_custom-compute-action", config)
			jobDesc := types.JobDescription{Flag: types.CustomComputeCapability, NodeType: types.WorkerNode}

			if _, ok := jobSpecs[jobDesc]; !ok {
				jobSpecs[jobDesc] = []*jobv1.ProposeJobRequest{jobSpec}
			} else {
				jobSpecs[jobDesc] = append(jobSpecs[jobDesc], jobSpec)
			}
		}

		var nodeEthAddr common.Address
		for _, label := range workerNode.Labels {
			if label.Key == node.EthAddressKey {
				if label.Value == nil {
					return nil, errors.New("eth address label value is nil")
				}
				if *label.Value == "" {
					return nil, errors.New("eth address label value is empty")
				}
				nodeEthAddr = common.HexToAddress(*label.Value)
				break
			}
		}

		var ocr2KeyBundleID string
		for _, label := range workerNode.Labels {
			if label.Key == devenv.NodeOCR2KeyBundleIDType {
				if label.Value == nil {
					return nil, errors.New("ocr2 key bundle id label value is nil")
				}
				if *label.Value == "" {
					return nil, errors.New("ocr2 key bundle id label value is empty")
				}
				ocr2KeyBundleID = *label.Value
				break
			}
		}

		if creflags.HasFlag(input.DonWithMetadata.Flags, types.OCR3Capability) {
			jobSpec := jobs.WorkerOCR3(nodeID, input.OCR3CapabilityAddress, nodeEthAddr, ocr2KeyBundleID, ocrPeeringData, chainIDUint64)
			jobDesc := types.JobDescription{Flag: types.OCR3Capability, NodeType: types.WorkerNode}

			if _, ok := jobSpecs[jobDesc]; !ok {
				jobSpecs[jobDesc] = []*jobv1.ProposeJobRequest{jobSpec}
			} else {
				jobSpecs[jobDesc] = append(jobSpecs[jobDesc], jobSpec)
			}
		}
	}

	return jobSpecs, nil
}
