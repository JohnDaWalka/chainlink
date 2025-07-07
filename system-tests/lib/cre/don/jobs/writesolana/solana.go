package solana

import (
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"

	libjobs "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	libnode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
)

// PROBABLY YOU DON'T NEED THE JOB SPEC, CONFIG SHOULD BE ENOUGH

// TODO potentially implement a higher order function that would return function with this signature: func GenerateWriteSolanaJobSpecs(input types.CreateJobsInput) (types.DonsToJobSpecs, error) like:
/*
var CronJobSpecFactoryFn = func(cronBinaryPath string) types.JobSpecFactoryFn {
	return func(input *types.JobSpecFactoryInput) (types.DonsToJobSpecs, error) {
		return GenerateJobSpecs(
			input.DonTopology,
			cronBinaryPath,
		)
	}
}
*/

var WriteSolanaJobSpecFactoryFn = func(solanaBinaryPathInTheContainer string) types.JobSpecFactoryFn {
	return func(input *types.JobSpecFactoryInput) (types.DonsToJobSpecs, error) {
		return GenerateWriteSolanaJobSpecs(input.DonTopology, solanaBinaryPathInTheContainer)
	}
}

func GenerateWriteSolanaJobSpecs(donTopology *types.DonTopology, solanaBinaryPathInTheContainer string) (types.DonsToJobSpecs, error) {
	if donTopology == nil {
		return nil, errors.New("topology is nil")
	}
	donToJobSpecs := make(types.DonsToJobSpecs)

	for _, donWithMetadata := range donTopology.DonsWithMetadata {
		workflowNodeSet, err := libnode.FindManyWithLabel(donWithMetadata.NodesMetadata, &types.Label{Key: libnode.NodeTypeKey, Value: types.WorkerNode}, libnode.EqualLabels)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find worker nodes")
		}

		for _, workerNode := range workflowNodeSet {
			nodeID, nodeIDErr := libnode.FindLabelValue(workerNode, libnode.NodeIDKey)
			if nodeIDErr != nil {
				return nil, errors.Wrap(nodeIDErr, "failed to get node id from labels")
			}

			if flags.HasFlag(donWithMetadata.Flags, types.WriteSolanaCapability) {
				// TODO: implement this correctly
				// TODO: pass solana binary path in the container to the job spec

				donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], libjobs.WorkerSolana(nodeID, ""))
			}
		}
	}

	return donToJobSpecs, nil
}
