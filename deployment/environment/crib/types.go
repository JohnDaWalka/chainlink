package crib

import (
	"crypto/tls"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/deployment/environment/types"
	"google.golang.org/grpc/credentials"

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

type DeployCCIPOutput struct {
	AddressBook deployment.AddressBookMap
	NodeIDs     []string
}

func NewDeployEnvironmentFromCribOutput(lggr logger.Logger, output DeployOutput, deployerKey string) (*deployment.Environment, error) {
	sethClients := make([]*seth.Client, 0)
	for _, chain := range output.BlockchainOutputs {
		if chain.Family == "evm" {
			sethClient, err := seth.NewClientBuilder().
				WithRpcUrl(chain.Nodes[0].HostWSUrl).
				WithPrivateKeys([]string{deployerKey}).
				Build()
			if err != nil {
				return nil, errors.Wrap(err, "failed to build sethClient")
			}
			sethClients = append(sethClients, sethClient)
		}
		// todo: add solana handling here
	}

	env, err := devenv.NewEnvironmentBuilder(lggr).
		WithNodeSetOutput([]*types.WrappedNodeOutput{output.NodesetOutput}).
		WithJobDistributor(output.JDOutput, credentials.NewTLS(&tls.Config{
			MinVersion: tls.VersionTLS12,
		})).
		WithBlockchains(output.BlockchainOutputs).
		WithSethClients(sethClients).
		WithExistingAddresses(output.AddressBook).
		Build()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create environment")
	}

	return env.Environment, nil
}
