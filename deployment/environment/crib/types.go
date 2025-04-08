package crib

import (
	"crypto/tls"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc/credentials"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/deployment/environment/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
)

const (
	CRIB_ENV_NAME = "Crib Environment"
)

type DeployOutput struct {
	NodeIDs           []string
	Chains            []devenv.ChainConfig   // chain selector -> Chain Config
	AddressBook       deployment.AddressBook // Addresses of all contracts
	JDOutput          *jd.Output
	BlockchainOutputs types.ChainIDToBlockchainOutputs
	NodesetOutput     *types.WrappedNodeOutput
}

type CCIPOnChainDeployOutput struct {
	AddressBook deployment.AddressBookMap
	NodeIDs     []string
}

// BuildTopology In CCIP we don't need topology, but we use it here as a wrapper to provide Nodes with Metadata
func BuildTopology(nodeOutputs []*clnode.Output) *types.Topology {
	nodesMetadata := make([]*types.NodeMetadata, 0)

	// Add Node labels required to  build environment
	for i := range nodeOutputs {
		nodeWithLabels := types.NodeMetadata{}
		nodeType := types.WorkerNode
		if i == 0 {
			nodeType = types.BootstrapNode
		}

		nodeWithLabels.Labels = append(nodeWithLabels.Labels, &types.Label{
			Key:   types.NodeTypeKey,
			Value: nodeType,
		})

		nodesMetadata = append(nodesMetadata, &nodeWithLabels)
	}

	donsWithMetadata := []*types.DonMetadata{
		{
			NodesMetadata: nodesMetadata,
			Name:          "CCIP DON",
		},
	}

	return &types.Topology{
		// set some dummy ID as this is required field
		WorkflowDONID: uint32(1),
		DonsMetadata:  donsWithMetadata,
	}
}

func NewDeployEnvironmentFromCribOutput(lggr logger.Logger, output DeployOutput, deployerKey string) (*devenv.EnvironmentWithTopology, error) {
	sethClients := make([]*seth.Client, 0)
	for _, chain := range output.BlockchainOutputs {
		if chain.Family == "evm" {
			sethClient, err := seth.NewClientBuilder().
				WithRpcUrl(chain.Nodes[0].ExternalWSUrl).
				WithPrivateKeys([]string{deployerKey}).
				WithProtections(false, false, seth.MustMakeDuration(1*time.Minute)).
				Build()
			if err != nil {
				return nil, errors.Wrap(err, "failed to build sethClient")
			}
			sethClients = append(sethClients, sethClient)
		}
		// todo: add solana handling here
	}

	topology := BuildTopology(output.NodesetOutput.CLNodes)

	env, err := devenv.NewEnvironmentBuilder(lggr).
		WithNodeSetOutput([]*types.WrappedNodeOutput{
			output.NodesetOutput,
		}).
		WithTopology(topology).
		WithJobDistributor(output.JDOutput, credentials.NewTLS(&tls.Config{
			MinVersion: tls.VersionTLS12,
		})).
		WithBlockchains(output.BlockchainOutputs).
		WithSethClients(sethClients).
		WithExistingAddresses(output.AddressBook).
		Build()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build environment from CRIB deploy output")
	}

	return env, nil
}
