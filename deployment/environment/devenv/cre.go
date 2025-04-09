package devenv

import (
	"context"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"strings"

	"github.com/pkg/errors"
	"google.golang.org/grpc/credentials"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/types"
)

type EnvironmentWithTopology struct {
	Environment *deployment.Environment
	DonTopology *DonTopology
}

type CREEnvironmentBuilder struct {
	jdOutput          *jd.Output
	blockchainOutputs types.ChainIDToBlockchainOutputs
	sethClients       []*seth.Client
	nodeSetOutput     []*types.WrappedNodeOutput
	existingAddresses deployment.AddressBook
	topology          *types.Topology
	credentials       credentials.TransportCredentials
	logger            logger.Logger
	errs              []string
}

func NewCREEnvironmentBuilder(lgr logger.Logger) *CREEnvironmentBuilder {
	b := &CREEnvironmentBuilder{
		logger: lgr,
	}

	if lgr == nil {
		b.errs = append(b.errs, "logger not set")
	}
	return b
}

func (b *CREEnvironmentBuilder) WithJobDistributor(jdOutput *jd.Output, jdTransportCredentials credentials.TransportCredentials) *CREEnvironmentBuilder {
	if jdTransportCredentials == nil {
		b.errs = append(b.errs, "jd credentials not set")
	}
	if jdOutput == nil {
		b.errs = append(b.errs, "jd output not set")
		return b
	}
	if jdOutput.ExternalGRPCUrl == "" {
		b.errs = append(b.errs, "external gRPC url not set")
	}
	if jdOutput.InternalWSRPCUrl == "" {
		b.errs = append(b.errs, "internal wsRPC url not set")
	}

	b.jdOutput = jdOutput
	b.credentials = jdTransportCredentials
	return b
}

func (b *CREEnvironmentBuilder) WithBlockchains(blockchainOutputs types.ChainIDToBlockchainOutputs) *CREEnvironmentBuilder {
	if len(blockchainOutputs) == 0 {
		b.errs = append(b.errs, "blockchain outputs not set")
	}
	b.blockchainOutputs = blockchainOutputs
	return b
}

func (b *CREEnvironmentBuilder) WithSethClients(sethClients []*seth.Client) *CREEnvironmentBuilder {
	if len(sethClients) == 0 {
		b.errs = append(b.errs, "seth clients not set")
	}
	b.sethClients = sethClients
	return b
}

func (b *CREEnvironmentBuilder) WithNodeSets(nodeSetOutput []*types.WrappedNodeOutput) *CREEnvironmentBuilder {
	if nodeSetOutput == nil {
		b.errs = append(b.errs, "node set output not set")
	}
	if len(nodeSetOutput) == 0 {
		b.errs = append(b.errs, "node set outputs are empty")
	}
	b.nodeSetOutput = nodeSetOutput
	return b
}

func (b *CREEnvironmentBuilder) WithExistingAddresses(existingAddresses deployment.AddressBook) *CREEnvironmentBuilder {
	b.existingAddresses = existingAddresses
	return b
}

func (b *CREEnvironmentBuilder) WithTopology(topology *types.Topology) *CREEnvironmentBuilder {
	if topology != nil {
		if len(topology.DonsMetadata) == 0 {
			b.errs = append(b.errs, "metadata not set")
		}
		if topology.WorkflowDONID == 0 {
			b.errs = append(b.errs, "workflow don id not set")
		}
	}

	b.topology = topology
	return b
}

func (b *CREEnvironmentBuilder) Build() (*EnvironmentWithTopology, error) {
	if len(b.errs) > 0 {
		return nil, errors.New("validation errors: " + strings.Join(b.errs, ", "))
	}
	if b.topology == nil {
		return nil, errors.New("topology not set")
	}
	if b.blockchainOutputs == nil {
		return nil, errors.New("blockchain outputs not set")
	}
	if b.nodeSetOutput == nil {
		return nil, errors.New("nodeSetOutput not set")
	}
	if b.jdOutput == nil {
		return nil, errors.New("jd output not set")
	}

	envs := make([]*deployment.Environment, len(b.nodeSetOutput))
	dons := make([]*DON, len(b.nodeSetOutput))

	var allNodesInfo []NodeInfo

	sethClients := b.sethClients

	chains := chainsFromBlockchainOutputs(sethClients, b.blockchainOutputs)

	for idx, nodeOutput := range b.nodeSetOutput {
		// check how many bootstrap nodes we have in each DON
		bootstrapNodes, err := FindManyWithLabel(b.topology.DonsMetadata[idx].NodesMetadata, &types.Label{Key: types.NodeTypeKey, Value: types.BootstrapNode}, EqualLabels)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find bootstrap nodes")
		}
		nodeInfo, err := GetNodeInfo(nodeOutput.Output, nodeOutput.NodeSetName, len(bootstrapNodes))
		if err != nil {
			return nil, errors.Wrap(err, "failed to get node info")
		}
		allNodesInfo = append(allNodesInfo, nodeInfo...)

		// if DON has no capabilities, we don't need to create chain configs (e.g. for gateway nodes)
		// we indicate to `NewEnvironment` that it should skip chain creation by passing an empty chain config
		if len(nodeOutput.Capabilities) == 0 {
			chains = []ChainConfig{}
		}

		jdConfig := JDConfig{
			GRPC:     b.jdOutput.ExternalGRPCUrl,
			WSRPC:    b.jdOutput.InternalWSRPCUrl,
			Creds:    b.credentials,
			NodeInfo: nodeInfo,
		}

		devenvConfig := EnvironmentConfig{
			JDConfig: jdConfig,
			Chains:   chains,
		}

		b.logger.Infow("creating CLD environment")
		env, don, err := NewEnvironment(context.Background, b.logger, devenvConfig)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create a CLD environment")
		}

		envs[idx] = env
		dons[idx] = don
	}

	var nodeIDs []string
	for _, env := range envs {
		nodeIDs = append(nodeIDs, env.NodeIDs...)
	}

	for i, don := range dons {
		for j, node := range b.topology.DonsMetadata[i].NodesMetadata {
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

	if len(b.nodeSetOutput) > 0 {
		// We create a new instance of JD client using `allNodesInfo` instead of `nodeInfo` to ensure that it can interact with all nodes.
		// Otherwise, JD would fail to accept job proposals for unknown nodes, even though it would still propose jobs to them. And that
		// would be happening silently, without any error messages, and we wouldn't know about it until much later.
		jd, err = NewJDClient(context.Background(), JDConfig{
			GRPC:     b.jdOutput.ExternalGRPCUrl,
			WSRPC:    b.jdOutput.InternalWSRPCUrl,
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
			ExistingAddresses: b.existingAddresses,
			Chains:            envs[0].Chains,
			Offchain:          jd,
			OCRSecrets:        envs[0].OCRSecrets,
			GetContext:        envs[0].GetContext,
			NodeIDs:           nodeIDs,
		},
	}

	if b.topology != nil {
		donTopology := &DonTopology{
			GatewayConnectorOutput: b.topology.GatewayConnectorOutput,
			WorkflowDonID:          b.topology.WorkflowDONID,
		}

		for i, donMetadata := range b.topology.DonsMetadata {
			donTopology.DonsWithMetadata = append(donTopology.DonsWithMetadata, &DonWithMetadata{
				DON:         dons[i],
				DonMetadata: donMetadata,
			})
		}

		output.DonTopology = donTopology
	}

	return output, nil
}
