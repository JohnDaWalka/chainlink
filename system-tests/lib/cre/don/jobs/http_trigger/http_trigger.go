package http_trigger

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	crenode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

const triggerServiceConfigTemplate = `"""
{
	"incomingRateLimiter": {
		"globalBurst": 10,
		"globalRPS": 50,
		"perSenderBurst": 10,
		"perSenderRPS": 10
	},
	"outgoingRateLimiter": {
		"globalBurst": 10,
		"globalRPS": 50,
		"perSenderBurst": 10,
		"perSenderRPS": 10
	}
}
"""
`

var HttpTriggerJobSpecFactoryFn = func(httpTriggerBinaryPath string) types.JobSpecFactoryFn {
	return func(input *types.JobSpecFactoryInput) (types.DonsToJobSpecs, error) {
		return GenerateJobSpecs(
			input.DonTopology,
			httpTriggerBinaryPath,
		)
	}
}

func GenerateJobSpecs(donTopology *types.DonTopology, httpTriggerBinaryPath string) (types.DonsToJobSpecs, error) {
	if donTopology == nil {
		return nil, errors.New("topology is nil")
	}
	donToJobSpecs := make(types.DonsToJobSpecs)

	for _, donWithMetadata := range donTopology.DonsWithMetadata {
		if !flags.HasFlag(donWithMetadata.Flags, types.HttpTriggerCapability) {
			continue
		}
		workflowNodeSet, err := crenode.FindManyWithLabel(donWithMetadata.NodesMetadata, &types.Label{Key: crenode.NodeTypeKey, Value: types.WorkerNode}, crenode.EqualLabels)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find worker nodes")
		}

		for _, workerNode := range workflowNodeSet {
			nodeID, nodeIDErr := crenode.FindLabelValue(workerNode, crenode.NodeIDKey)
			if nodeIDErr != nil {
				return nil, errors.Wrap(nodeIDErr, "failed to get node id from labels")
			}

			donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobs.WorkerStandardCapability(nodeID, types.HttpTriggerCapability, httpTriggerBinaryPath, triggerServiceConfigTemplate))
		}
	}

	return donToJobSpecs, nil
}
