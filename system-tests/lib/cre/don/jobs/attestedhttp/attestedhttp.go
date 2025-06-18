package attestedhttp

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	libnode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

var AttestedhttpJobSpecFactoryFn = func(attestedhttpBinaryPath string) types.JobSpecFactoryFn {
	return func(input *types.JobSpecFactoryInput) (types.DonsToJobSpecs, error) {
		return GenerateJobSpecs(
			input.DonTopology,
			attestedhttpBinaryPath,
		)
	}
}

func GenerateJobSpecs(donTopology *types.DonTopology, attestedhttpBinaryPath string) (types.DonsToJobSpecs, error) {
	if donTopology == nil {
		return nil, errors.New("topology is nil")
	}
	donToJobSpecs := make(types.DonsToJobSpecs)

	for _, donWithMetadata := range donTopology.DonsWithMetadata {
		workflowNodeSet, err := libnode.FindManyWithLabel(donWithMetadata.NodesMetadata, &types.Label{Key: libnode.NodeTypeKey, Value: types.WorkerNode}, libnode.EqualLabels)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find worker nodes")
		}

		for i, workerNode := range workflowNodeSet {
			nodeID, nodeIDErr := node.FindLabelValue(workerNode, node.NodeIDKey)
			if nodeIDErr != nil {
				return nil, errors.Wrap(nodeIDErr, "failed to get node id from labels")
			}

			url := "[YOUR_SECURE_URL]"
			pub := "[YOUR_SECURE_PUBLIC_KEY]"
			sig := []string{"YOUR_SIGNERS_1", "YOUR_SIGNERS_N"}
			id := "YOUR_INSTANCE_ID"
			c2 := fmt.Sprintf("'{\"SecureURL\":\"%s\",\"InstanceID\":\"%s\",\"PublicKey\":\"%s\",\"SignerPrivateKey\":\"%s\"}'", url, id, pub, sig[i])
			fmt.Println("Attested HTTP config:", c2)
			if flags.HasFlag(donWithMetadata.Flags, types.AttestedHTTPCapability) {
				donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobs.WorkerStandardCapability(nodeID, types.AttestedHTTPCapability, attestedhttpBinaryPath, c2))
			}
		}
	}

	return donToJobSpecs, nil
}
