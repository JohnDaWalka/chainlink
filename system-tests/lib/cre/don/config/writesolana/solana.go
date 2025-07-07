package writesolana

import (
	"strconv"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

// TODO: implement higher order function that would return function with this signature: func GenerateConfig(input types.GenerateConfigsInput) (types.NodeIndexToConfigOverride, error)

func GenerateConfig(input types.GenerateConfigsInput) (types.NodeIndexToConfigOverride, error) {
	configOverrides := make(types.NodeIndexToConfigOverride)

	// find worker nodes
	workflowNodeSet, err := node.FindManyWithLabel(input.DonMetadata.NodesMetadata, &types.Label{Key: node.NodeTypeKey, Value: types.WorkerNode}, node.EqualLabels)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find worker nodes")
	}

	for i := range workflowNodeSet {
		var nodeIndex int
		for _, label := range workflowNodeSet[i].Labels {
			if label.Key == node.IndexKey {
				nodeIndex, err = strconv.Atoi(label.Value)
				if err != nil {
					return nil, errors.Wrap(err, "failed to convert node index to int")
				}
			}
		}

		solanaAddress, solanaAddressErr := node.FindLabelValue(workflowNodeSet[i], node.NodeSolanaAddressKey)
		if solanaAddressErr != nil {
			return nil, errors.Wrap(solanaAddressErr, "failed to get solana address from labels")
		}

		if flags.HasFlag(input.Flags, types.WriteSolanaCapability) {
			// TODO: implement this correctly
			// TODO: pass solana address to the the config
			_ = solanaAddress
			configOverrides[nodeIndex] += config.WorkerSolana(
				[]any{},
			)
		}
	}

	return configOverrides, nil
}
