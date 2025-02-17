package blockchain

import (
	"os"

	"github.com/pkg/errors"
	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink/system-tests/lib/keystone/types"
)

func Start(blockchainInput *blockchain.Input, keystoneEnv *types.KeystoneEnvironment) error {
	bc, err := blockchain.NewBlockchainNetwork(blockchainInput)
	if err != nil {
		return errors.Wrap(err, "failed to create blockchain network")
	}

	pkey := os.Getenv("PRIVATE_KEY")
	if pkey == "" {
		return errors.New("PRIVATE_KEY env var must be set")
	}

	sc, err := seth.NewClientBuilder().
		WithRpcUrl(bc.Nodes[0].HostWSUrl).
		WithPrivateKeys([]string{pkey}).
		Build()
	if err != nil {
		return errors.Wrap(err, "failed to create seth client")
	}

	chainSelector, err := chainselectors.SelectorFromChainId(sc.Cfg.Network.ChainID)
	if err != nil {
		return errors.Wrapf(err, "failed to get chain selector for chain id %d", sc.Cfg.Network.ChainID)
	}

	keystoneEnv.Blockchain = bc
	keystoneEnv.SethClient = sc
	keystoneEnv.DeployerPrivateKey = pkey
	keystoneEnv.ChainSelector = chainSelector

	return nil
}
