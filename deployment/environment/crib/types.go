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
	Chains      map[uint64]devenv.ChainConfig // chain selector -> Chain Config
	AddressBook deployment.AddressBook        // Addresses of all contracts
}

type DeployCCIPOutput struct {
	Ab deployment.AddressBook `json:"addressbook"`
}

func NewDeployEnvironmentFromCrib(lggr logger.Logger, output DeployOutput) (*deployment.Environment, error) {
	var nodeIds = make([]string, 0)
	for _, n := range output.DON.Nodes {
		nodeIds = append(nodeIds, n.NodeId)
	}

	var chains = make(map[uint64]deployment.Chain)
	for sel, chain := range output.Chains {
		multiClient, err := deployment.NewMultiClient(lggr, []deployment.RPC{{WSURL: chain.WSRPCs[0]}})
		if err != nil {
			return nil, err
		}
		chains[sel] = deployment.Chain{
			Selector:    sel,
			Client:      multiClient,
			DeployerKey: chain.DeployerKey,
		}
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
