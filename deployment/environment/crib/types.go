package crib

import (
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
)

const (
	CRIB_ENV_NAME = "Crib Environment"
)

type DeployOutput struct {
	DON         devenv.DON
	Chains      []devenv.ChainConfig   // chain selector -> Chain Config
	AddressBook deployment.AddressBook // Addresses of all contracts
}

type DeployCCIPOutput struct {
	deployment.AddressBook `json:"addressbook"`
}

func NewDeployEnvironmentFromCribOutput(lggr logger.Logger, output DeployOutput) (*deployment.Environment, error) {
	var nodeIds = make([]string, 0)
	for _, n := range output.DON.Nodes {
		nodeIds = append(nodeIds, n.NodeId)
	}

	chains, err := devenv.NewChains(lggr, output.Chains)
	if err != nil {
		return nil, err
	}
	return deployment.NewEnvironment(
		CRIB_ENV_NAME,
		lggr,
		output.AddressBook,
		chains,
		nodeIds,
		nil, // todo: populate the offchain client using output.DON
	), nil
}
