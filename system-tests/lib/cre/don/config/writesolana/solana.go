package writesolana

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

type Config struct {
}

func GetGenerateConfig(in Config) func(cre.GenerateConfigsInput) (cre.NodeIndexToConfigOverride, error) {
	return func(input cre.GenerateConfigsInput) (cre.NodeIndexToConfigOverride, error) {
		configOverrides := make(cre.NodeIndexToConfigOverride)
		if flags.HasFlag(input.Flags, cre.WriteSolanaCapability) {
			workerSolanaInputs := make([]*config.WorkerSolanaInput, 0)
			for chainSelector, bcOut := range input.BlockchainOutput {
				if bcOut.SolChain == nil {
					continue
				}

				chainID, err := bcOut.SolClient.GetGenesisHash(context.Background())
				if err != nil {
					return nil, errors.Wrap(err, "failed to get chainID from solana")
				}

				workerSolanaInputs = append(workerSolanaInputs, &config.WorkerSolanaInput{
					Name:    fmt.Sprintf("node-%d", chainSelector),
					ChainID: chainID.String(),
					NodeURL: bcOut.BlockchainOutput.Nodes[0].InternalHTTPUrl,
					// TODO PLEX-1622 add the rest solana inputs (forwarder, forwarder state from datastore) once changesets are integrated
				})
			}

			workflowNodeSet, err := node.FindManyWithLabel(input.DonMetadata.NodesMetadata, &cre.Label{Key: node.NodeTypeKey, Value: cre.WorkerNode}, node.EqualLabels)
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
						break
					}
				}

				configOverrides[nodeIndex] = config.WorkerSolana(workerSolanaInputs)
			}

			bootstrapNodes, err := node.FindManyWithLabel(input.DonMetadata.NodesMetadata, &cre.Label{Key: node.NodeTypeKey, Value: cre.BootstrapNode}, node.EqualLabels)
			if err != nil {
				return nil, errors.Wrap(err, "failed to find bootstrap nodes")
			}
			if len(bootstrapNodes) > 0 {
				bootstrapNode := bootstrapNodes[0]
				var nodeIndex int
				for _, label := range bootstrapNode.Labels {
					if label.Key == node.IndexKey {
						nodeIndex, err = strconv.Atoi(label.Value)
						if err != nil {
							return nil, errors.Wrap(err, "failed to convert node index to int")
						}
						break
					}
				}
				fmt.Println("bootstrap node idx", nodeIndex)
				configOverrides[nodeIndex] = config.WorkerSolana(workerSolanaInputs)
			}
		}
		return configOverrides, nil
	}
}
