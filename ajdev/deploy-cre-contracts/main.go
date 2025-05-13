// anvil --host 0.0.0.0 --port 8545 --chain-id 1337 --block-time 1
package main

import (
	"context"
	"log"
	"testing"
	"time"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"
	libcontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
)

func main() {
	log.Println("Starting deployment...")

	// create seth client
	sethc, err := seth.NewClientBuilder().
		WithRpcUrl("ws://0.0.0.0:8545").
		WithPrivateKeys([]string{"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}).
		// do not check if there's a pending nonce nor check node's health
		WithProtections(false, false, seth.MustMakeDuration(time.Second)).
		Build()
	if err != nil {
		log.Fatalf("Failed to create seth client: %v", err)
	}

	chainSelector, err := chainselectors.SelectorFromChainId(sethc.Cfg.Network.ChainID)
	if err != nil {
		log.Fatalf("failed to get chain selector for chain id %d", sethc.Cfg.Network.ChainID)
	}

	// create chains configs
	chainsConfigs := []devenv.ChainConfig{}

	chainsConfigs = append(chainsConfigs, devenv.ChainConfig{
		ChainID:   1337,
		ChainName: "anvil",
		ChainType: "EVM",
		WSRPCs: []devenv.CribRPCs{{
			External: "ws://0.0.0.0:8545",
			Internal: "ws://0.0.0.0:8545",
		}},
		HTTPRPCs: []devenv.CribRPCs{{
			External: "http://0.0.0.0:8545",
			Internal: "http://0.0.0.0:8545",
		}},
		// set nonce to nil, so that it will be fetched from the RPC node
		DeployerKey: sethc.NewTXOpts(seth.WithNonce(nil)),
	})

	log.Println("Chains configs:", chainsConfigs)

	// create chains
	sfLogger := cldlogger.NewSingleFileLogger(&testing.T{})
	chains, err := devenv.NewChains(sfLogger, chainsConfigs)
	if err != nil {
		log.Fatalf("Failed to create chains: %v", err)
	}

	// create deployment environment
	ctx := context.Background()
	depenv := &deployment.Environment{
		Logger:            sfLogger,
		Chains:            chains,
		ExistingAddresses: cldf.NewMemoryAddressBook(),
		GetContext: func() context.Context {
			return ctx
		},
	}

	// deploy contracts

	// deploy OCR3 contract
	ocr3Output, ocr3Err := keystone_changeset.DeployOCR3(*depenv, chainSelector) // //nolint:staticcheck // will migrate in DX-641
	if ocr3Err != nil {
		log.Fatalf("failed to deploy OCR3 contract: %v", ocr3Err)
	}

	mergeErr := depenv.ExistingAddresses.Merge(ocr3Output.AddressBook) //nolint:staticcheck // won't migrate now
	if mergeErr != nil {
		log.Fatalf("failed to merge address book: %v", mergeErr)
	}

	log.Printf("Deployed OCR3 contract on chain %d at %s", chainSelector, libcontracts.MustFindAddressesForChain(depenv.ExistingAddresses, chainSelector, keystone_changeset.OCR3Capability.String())) //nolint:staticcheck // won't migrate now
}
