package devenv

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/devenv/dctl"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/devenv/griddle"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
	"gopkg.in/yaml.v3"
)

// relative path to the directory where generated griddle configs will be stored
const griddleConfigsDir = "./configs/griddle-devenv/"

// Bootstrap bootstraps the development environment based on the provided infra input.
func Bootstrap(infraInput infra.Provider) error {
	// connect to telepresence
	dctlClient := dctl.NewDctlClient()
	err := dctlClient.NetworkConnect()
	if err != nil {
		return errors.Wrapf(err, "failed to connect to telepresence")
	}

	return nil
}

func DeployBlockchain(infraIn infra.Provider, input *cre.DeployGriddleDevenvBlockchainInput) (*blockchain.Output, error) {
	griddleBaseConfigDir := "base"

	// generate griddle yaml config for the blockchain
	template := griddle.DeployTemplate{
		Instances: []griddle.Instance{
			{
				Name:                "anvil-1337",
				Chart:               "app-template",
				Version:             "4.3.0",
				Repository:          "https://bjw-s-labs.github.io/helm-charts",
				LocalRepositoryName: "bjw-s",
				Config: []string{
					filepath.Join(griddleBaseConfigDir, "/values/bjw-s/blockchain.anvil.base.yaml"),
					filepath.Join(griddleBaseConfigDir, "/values/bjw-s/blockchain.anvil.1337.yaml"),
				},
			},
		},
	}

	griddleConfig := griddle.ConfigFromInfraInputWithTemplate(infraIn, "blockchain", template)
	yamlData, err := yaml.Marshal(&griddleConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal griddle config")
	}

	// save yamlData to a generated config dir
	configFilePath := filepath.Join(griddleConfigsDir, "griddle-blockchain.yaml")

	err = os.WriteFile(configFilePath, yamlData, 0644)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to write griddle config to %s", configFilePath)
	}

	// deploy the blockchain using dctl wrapper
	dctlClient := dctl.NewDctlClient()
	err = dctlClient.DeployApply(configFilePath, "blockchain", infraIn.GriddleDevenvInput.Namespace, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to deploy blockchain using dctl")
	}

	// return the blockchain output
	return &blockchain.Output{
		Type:    input.BlockchainInput.Type,
		Family:  "evm",
		ChainID: input.BlockchainInput.ChainID,
		Nodes: []*blockchain.Node{
			{
				InternalWSUrl:   "ws://anvil-1337:8545",
				ExternalWSUrl:   "ws://anvil-1337:8545",
				InternalHTTPUrl: "http://anvil-1337:8545",
				ExternalHTTPUrl: "http://anvil-1337:8545",
			},
		},
	}, nil

}
