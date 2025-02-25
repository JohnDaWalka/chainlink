package environment

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	libnode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

func BuildFullCLDEnvironment(lgr logger.Logger, input *types.FullCLDEnvironmentInput) (*types.FullCLDEnvironmentOutput, error) {
	if input == nil {
		return nil, errors.New("input is nil")
	}
	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "input validation failed")
	}

	envs := make([]*deployment.Environment, len(input.NodeSetOutput))
	dons := make([]*devenv.DON, len(input.NodeSetOutput))

	for i, nodeOutput := range input.NodeSetOutput {
		// assume that each nodeset has only one bootstrap node
		nodeInfo, err := libnode.GetNodeInfo(nodeOutput.Output, nodeOutput.NodeSetName, 1)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get node info")
		}

		jdConfig := devenv.JDConfig{
			GRPC:     input.JdOutput.HostGRPCUrl,
			WSRPC:    input.JdOutput.DockerWSRPCUrl,
			Creds:    insecure.NewCredentials(),
			NodeInfo: nodeInfo,
		}

		devenvConfig := devenv.EnvironmentConfig{
			JDConfig: jdConfig,
			Chains: []devenv.ChainConfig{
				{
					ChainID:   input.SethClient.Cfg.Network.ChainID,
					ChainName: input.SethClient.Cfg.Network.Name,
					ChainType: strings.ToUpper(input.BlockchainOutput.Family),
					WSRPCs: []devenv.CribRPCs{{
						External: input.BlockchainOutput.Nodes[0].HostWSUrl,
						Internal: input.BlockchainOutput.Nodes[0].DockerInternalWSUrl,
					}},
					HTTPRPCs: []devenv.CribRPCs{{
						External: input.BlockchainOutput.Nodes[0].HostHTTPUrl,
						Internal: input.BlockchainOutput.Nodes[0].DockerInternalHTTPUrl,
					}},
					DeployerKey: input.SethClient.NewTXOpts(seth.WithNonce(nil)), // set nonce to nil, so that it will be fetched from the chain
				},
			},
		}

		env, don, err := devenv.NewEnvironment(context.Background, lgr, devenvConfig)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create environment")
		}

		envs[i] = env
		dons[i] = don
	}

	var nodeIDs []string
	for _, env := range envs {
		nodeIDs = append(nodeIDs, env.NodeIDs...)
	}

	for i, don := range dons {
		for j, node := range input.Topology.Metadata[i].NodesMetadata {
			// add only new labels, we need to avoid duplicates, because nodeInfo struct passed to libenv.BuildChainlinkDeploymentEnv
			// already contains some labels that we needed to add before to create config
			for _, donLabel := range don.Nodes[j].Labels() {
				if !node.HasLabel(donLabel) {
					node.Labels = append(node.Labels, donLabel)
				}
			}
		}
	}

	// we assume that all DONs run on the same chain and that there's only one chain
	// also, we don't care which instance of offchain client we use, because we have
	// only one instance of offchain client and we have just configured it to work
	// with nodes from all DONs
	output := &types.FullCLDEnvironmentOutput{
		Environment: &deployment.Environment{
			Name:              envs[0].Name,
			Logger:            envs[0].Logger,
			ExistingAddresses: input.ExistingAddresses,
			Chains:            envs[0].Chains,
			Offchain:          envs[0].Offchain,
			OCRSecrets:        envs[0].OCRSecrets,
			GetContext:        envs[0].GetContext,
			NodeIDs:           nodeIDs,
		},
	}

	donTopology := &types.DonTopology{}
	donTopology.WorkflowDonID = input.Topology.WorkflowDONID

	for i, donMetadata := range input.Topology.Metadata {
		donTopology.Dons = append(donTopology.Dons, &types.DonWithMetadata{
			DON:         dons[i],
			DonMetadata: donMetadata,
		})
	}

	output.DonTopology = donTopology

	return output, nil
}
