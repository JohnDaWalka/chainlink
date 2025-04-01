package devenv

import (
	"context"
	"fmt"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/pkg/errors"
	"google.golang.org/grpc/credentials"
	"strings"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/types"
)

const (
	DevEnv = "devenv"
)

type EnvironmentConfig struct {
	Chains   []ChainConfig
	JDConfig JDConfig
}

func NewEnvironment(ctx func() context.Context, lggr logger.Logger, config EnvironmentConfig) (*deployment.Environment, *DON, error) {
	chains, err := NewChains(lggr, config.Chains)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create chains: %w", err)
	}
	offChain, err := NewJDClient(ctx(), config.JDConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create JD client: %w", err)
	}

	jd, ok := offChain.(*JobDistributor)
	if !ok {
		return nil, nil, errors.New("offchain client does not implement JobDistributor")
	}
	if jd == nil {
		return nil, nil, errors.New("offchain client is not set up")
	}
	var nodeIDs []string
	if jd.don != nil {
		// Gateway DON doesn't require any chain setup, and trying to create chains for it will fail,
		// because its nodes are missing chain-related configuration. Of course, we could add that configuration,
		// but its not how it is setup on production.
		if len(config.Chains) > 0 {
			err = jd.don.CreateSupportedChains(ctx(), config.Chains, *jd)
			if err != nil {
				return nil, nil, err
			}
		}
		nodeIDs = jd.don.NodeIds()
	}

	return deployment.NewEnvironment(
		DevEnv,
		lggr,
		deployment.NewMemoryAddressBook(),
		chains,
		nil, // sending nil for solana chains right now, we can build this when we need it
		nodeIDs,
		offChain,
		ctx,
		deployment.XXXGenerateTestOCRSecrets(),
	), jd.don, nil
}

type FullCLDEnvironmentOutput struct {
	Environment *deployment.Environment
	DonTopology *DonTopology
}

type FullCLDEnvironmentInput struct {
	JdOutput          *jd.Output
	BlockchainOutput  *blockchain.Output
	SethClient        *seth.Client
	NodeSetOutput     []*types.WrappedNodeOutput
	ExistingAddresses deployment.AddressBook
	Topology          *types.Topology
}

func (f *FullCLDEnvironmentInput) Validate() error {
	if f.JdOutput == nil {
		return errors.New("jd output not set")
	}
	if f.BlockchainOutput == nil {
		return errors.New("blockchain output not set")
	}
	if f.SethClient == nil {
		return errors.New("seth client not set")
	}
	if len(f.NodeSetOutput) == 0 {
		return errors.New("node set output not set")
	}
	if f.Topology == nil {
		return errors.New("topology not set")
	}
	if len(f.Topology.DonsMetadata) == 0 {
		return errors.New("metadata not set")
	}
	if f.Topology.WorkflowDONID == 0 {
		return errors.New("workflow don id not set")
	}
	return nil
}

func BuildFullCLDEnvironment(lgr logger.Logger, input *FullCLDEnvironmentInput, credentials credentials.TransportCredentials) (*FullCLDEnvironmentOutput, error) {
	if input == nil {
		return nil, errors.New("input is nil")
	}
	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "input validation failed")
	}

	envs := make([]*deployment.Environment, len(input.NodeSetOutput))
	dons := make([]*DON, len(input.NodeSetOutput))

	var allNodesInfo []NodeInfo
	chains := []ChainConfig{
		{
			ChainID:   input.SethClient.Cfg.Network.ChainID,
			ChainName: input.SethClient.Cfg.Network.Name,
			ChainType: strings.ToUpper(input.BlockchainOutput.Family),
			WSRPCs: []CribRPCs{{
				External: input.BlockchainOutput.Nodes[0].HostWSUrl,
				Internal: input.BlockchainOutput.Nodes[0].DockerInternalWSUrl,
			}},
			HTTPRPCs: []CribRPCs{{
				External: input.BlockchainOutput.Nodes[0].HostHTTPUrl,
				Internal: input.BlockchainOutput.Nodes[0].DockerInternalHTTPUrl,
			}},
			DeployerKey: input.SethClient.NewTXOpts(seth.WithNonce(nil)), // set nonce to nil, so that it will be fetched from the chain
		},
	}

	for i, nodeOutput := range input.NodeSetOutput {
		// assume that each nodeset has only one bootstrap node
		nodeInfo, err := GetNodeInfo(nodeOutput.Output, nodeOutput.NodeSetName, 1)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get node info")
		}
		allNodesInfo = append(allNodesInfo, nodeInfo...)

		// if DON has no capabilities we don't need to create chain configs (e.g. for gateway nodes)
		// we indicate to `NewEnvironment` that it should skip chain creation by passing an empty chain config
		if len(nodeOutput.Capabilities) == 0 {
			chains = []ChainConfig{}
		}

		jdConfig := JDConfig{
			GRPC:     input.JdOutput.HostGRPCUrl,
			WSRPC:    input.JdOutput.DockerWSRPCUrl,
			Creds:    credentials,
			NodeInfo: nodeInfo,
		}

		devenvConfig := EnvironmentConfig{
			JDConfig: jdConfig,
			Chains:   chains,
		}

		env, don, err := NewEnvironment(context.Background, lgr, devenvConfig)
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
		for j, node := range input.Topology.DonsMetadata[i].NodesMetadata {
			// both are required for job creation
			node.Labels = append(node.Labels, &types.Label{
				Key:   types.NodeIDKey,
				Value: don.NodeIds()[j],
			})

			node.Labels = append(node.Labels, &types.Label{
				Key:   types.NodeOCR2KeyBundleIDKey,
				Value: don.Nodes[j].Ocr2KeyBundleID,
			})

			node.Labels = append(node.Labels, &types.Label{
				Key:   types.NodeOCR2KeyBundleIDKey,
				Value: don.Nodes[j].Ocr2KeyBundleID,
			})
		}
	}

	var jd deployment.OffchainClient
	var err error

	if len(input.NodeSetOutput) > 0 {
		// We create a new instance of JD client using `allNodesInfo` instead of `nodeInfo` to ensure that it can interact with all nodes.
		// Otherwise, JD would fail to accept job proposals for unknown nodes, even though it would still propose jobs to them. And that
		// would be happening silently, without any error messages, and we wouldn't know about it until much later.
		jd, err = NewJDClient(context.Background(), JDConfig{
			GRPC:     input.JdOutput.HostGRPCUrl,
			WSRPC:    input.JdOutput.DockerWSRPCUrl,
			Creds:    credentials,
			NodeInfo: allNodesInfo,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to create JD client")
		}
	} else {
		jd = envs[0].Offchain
	}

	// we assume that all DONs run on the same chain and that there's only one chain
	output := &FullCLDEnvironmentOutput{
		Environment: &deployment.Environment{
			Name:              envs[0].Name,
			Logger:            envs[0].Logger,
			ExistingAddresses: input.ExistingAddresses,
			Chains:            envs[0].Chains,
			Offchain:          jd,
			OCRSecrets:        envs[0].OCRSecrets,
			GetContext:        envs[0].GetContext,
			NodeIDs:           nodeIDs,
		},
	}

	donTopology := &DonTopology{}
	donTopology.WorkflowDonID = input.Topology.WorkflowDONID

	for i, donMetadata := range input.Topology.DonsMetadata {
		donTopology.DonsWithMetadata = append(donTopology.DonsWithMetadata, &DonWithMetadata{
			DON:         dons[i],
			DonMetadata: donMetadata,
		})
	}

	output.DonTopology = donTopology

	return output, nil
}
