package crib

import (
	"crypto/tls"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc/credentials"

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

type CCIPInfraAndOnChainDeployOutput struct {
	// todo: Replace usages of NodeIDs, Chains to rely on new CTF based output types like BlockchainOutputs and NodesetOutput
	NodeIDs           []string
	Chains            []devenv.ChainConfig   // chain selector -> Chain Config
	AddressBook       deployment.AddressBook // Addresses of all contracts
	JDOutput          *jd.Output
	BlockchainOutputs types.ChainIDToBlockchainOutputs
	NodesetOutput     *simple_node_set.Output
}

type CCIPOnChainDeployOutput struct {
	AddressBook deployment.AddressBookMap
	NodeIDs     []string
}

func NewDeployEnvironmentFromCribOutput(lggr logger.Logger, output CCIPInfraAndOnChainDeployOutput, deployerKey string) (*deployment.Environment, *devenv.DON, error) {
	sethClients := make([]*seth.Client, 0)
	for _, chain := range output.BlockchainOutputs {
		if chain.Family == "evm" {
			sethClient, err := seth.NewClientBuilder().
				WithRpcUrl(chain.Nodes[0].ExternalWSUrl).
				WithPrivateKeys([]string{deployerKey}).
				WithProtections(false, false, seth.MustMakeDuration(1*time.Minute)).
				Build()
			if err != nil {
				return nil, nil, errors.Wrap(err, "failed to build sethClient")
			}
			sethClients = append(sethClients, sethClient)
		}
		// todo: add solana handling here
	}

	env, don, err := devenv.NewCCIPEnvironmentBuilder(lggr).
		WithNodeSet(output.NodesetOutput).
		WithJobDistributor(output.JDOutput, credentials.NewTLS(&tls.Config{
			MinVersion: tls.VersionTLS12,
		})).
		WithBlockchains(output.BlockchainOutputs).
		WithSethClients(sethClients).
		WithExistingAddresses(output.AddressBook).
		Build()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to build environment from CRIB deploy output")
	}

	return env, don, nil
}
