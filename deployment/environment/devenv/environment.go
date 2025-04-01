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

type EnvironmentWithTopology struct {
	Environment *deployment.Environment
	DonTopology *DonTopology
}

type EnvironmentBuilder struct {
	JDOutput          *jd.Output
	BlockchainOutput  *blockchain.Output
	SethClient        *seth.Client
	NodeSetOutput     []*types.WrappedNodeOutput
	ExistingAddresses deployment.AddressBook
	Topology          *types.Topology
	errs              []string
	credentials       credentials.TransportCredentials
	logger            logger.Logger
}

func NewEnvironmentBuilder(lgr logger.Logger, credentials credentials.TransportCredentials) *EnvironmentBuilder {
	return &EnvironmentBuilder{
		logger:      lgr,
		credentials: credentials,
	}
}

func (b *EnvironmentBuilder) WithJDOutput(jdOutput *jd.Output) *EnvironmentBuilder {
	if jdOutput == nil {
		b.errs = append(b.errs, "jd output not set")
	}
	b.JDOutput = jdOutput
	return b
}

func (b *EnvironmentBuilder) WithBlockchainOutput(blockchainOutput *blockchain.Output) *EnvironmentBuilder {
	if blockchainOutput == nil {
		b.errs = append(b.errs, "blockchain output not set")
	}
	b.BlockchainOutput = blockchainOutput
	return b
}

func (b *EnvironmentBuilder) WithSethClient(sethClient *seth.Client) *EnvironmentBuilder {
	if sethClient == nil {
		b.errs = append(b.errs, "seth client not set")
	}
	b.SethClient = sethClient
	return b
}

func (b *EnvironmentBuilder) WithNodeSetOutput(nodeSetOutput []*types.WrappedNodeOutput) *EnvironmentBuilder {
	if nodeSetOutput == nil || len(b.NodeSetOutput) == 0 {
		b.errs = append(b.errs, "node set output not set")
	}
	b.NodeSetOutput = nodeSetOutput
	return b
}

func (b *EnvironmentBuilder) WithExistingAddresses(existingAddresses deployment.AddressBook) *EnvironmentBuilder {
	b.ExistingAddresses = existingAddresses
	return b
}

// WithTopology Topology is required for CRE DONs
func (b *EnvironmentBuilder) WithTopology(topology *types.Topology) *EnvironmentBuilder {
	if topology != nil {
		if len(topology.DonsMetadata) == 0 {
			b.errs = append(b.errs, "metadata not set")
		}
		if topology.WorkflowDONID == 0 {
			b.errs = append(b.errs, "workflow don id not set")
		}
	}

	b.Topology = topology
	return b
}

func (b *EnvironmentBuilder) Build() (*EnvironmentWithTopology, error) {
	if len(b.errs) > 0 {
		return nil, errors.New("validation errors: " + strings.Join(b.errs, ", "))
	}

	input := b

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
			GRPC:     input.JDOutput.HostGRPCUrl,
			WSRPC:    input.JDOutput.DockerWSRPCUrl,
			Creds:    b.credentials,
			NodeInfo: nodeInfo,
		}

		devenvConfig := EnvironmentConfig{
			JDConfig: jdConfig,
			Chains:   chains,
		}

		env, don, err := NewEnvironment(context.Background, b.logger, devenvConfig)
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
			GRPC:     input.JDOutput.HostGRPCUrl,
			WSRPC:    input.JDOutput.DockerWSRPCUrl,
			Creds:    b.credentials,
			NodeInfo: allNodesInfo,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to create JD client")
		}
	} else {
		jd = envs[0].Offchain
	}

	// we assume that all DONs run on the same chain and that there's only one chain
	output := &EnvironmentWithTopology{
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

	if b.Topology != nil {
		donTopology := &DonTopology{}
		donTopology.WorkflowDonID = input.Topology.WorkflowDONID

		for i, donMetadata := range input.Topology.DonsMetadata {
			donTopology.DonsWithMetadata = append(donTopology.DonsWithMetadata, &DonWithMetadata{
				DON:         dons[i],
				DonMetadata: donMetadata,
			})
		}

		output.DonTopology = donTopology
	}

	return output, nil
}
