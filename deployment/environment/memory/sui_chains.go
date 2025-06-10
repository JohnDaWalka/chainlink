package memory

import (
	"context"
	"crypto/ed25519"
	"errors"
	"testing"
	"time"

	"github.com/pattonkan/sui-go/suiclient"
	"github.com/pattonkan/sui-go/suisigner"

	"github.com/smartcontractkit/freeport"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"

	suichain "github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

func getTestSuiChainSelectors() []uint64 {
	// TODO: CTF to support different chain ids, need to investigate if it's possible (thru node config.yaml?)
	return []uint64{chainsel.SUI_LOCALNET.Selector}
}

func GenerateChainsSui(t *testing.T, numChains int) map[uint64]suichain.Chain {
	testSuiChainSelectors := getTestSuiChainSelectors()
	if len(testSuiChainSelectors) < numChains {
		t.Fatalf("not enough test sui chain selectors available")
	}
	chains := make(map[uint64]suichain.Chain)
	for i := 0; i < numChains; i++ {
		selector := testSuiChainSelectors[i]
		chainID, err := chainsel.GetChainIDFromSelector(selector)
		require.NoError(t, err)

		url, _, privateKey, client := suiChain(t, chainID)
		chains[selector] = suichain.Chain{
			ChainMetadata: suichain.ChainMetadata{
				Selector: selector,
			},
			Client:      client,
			DeployerKey: privateKey,
			URL:         url,
			Confirm: func(txHash string, opts ...any) error {
				return errors.New("TODO: sui Confirm")
			},
		}
	}
	t.Logf("Created %d Sui chains: %+v", len(chains), chains)
	return chains
}

func suiChain(t *testing.T, chainID string) (string, string, ed25519.PrivateKey, *suiclient.ClientImpl) {
	t.Helper()

	// initialize the docker network used by CTF
	err := framework.DefaultNetwork(once)
	require.NoError(t, err)

	maxRetries := 10
	var url string
	var suiAddress string
	var mnemonic string
	for i := 0; i < maxRetries; i++ {
		// reserve all the ports we need explicitly to avoid port conflicts in other tests
		ports := freeport.GetN(t, 2)

		bcInput := &blockchain.Input{
			Image: "", // filled out by defaultSui function
			Type:  "sui",
			// TODO: this is unused, can it be applied?
			ChainID: chainID,
		}
		output, err := blockchain.NewBlockchainNetwork(bcInput)
		if err != nil {
			t.Logf("Error creating Sui network: %v", err)
			freeport.Return(ports)
			time.Sleep(time.Second)
			maxRetries -= 1
			continue
		}
		require.NoError(t, err)
		testcontainers.CleanupContainer(t, output.Container)
		url = output.Nodes[0].ExternalHTTPUrl

		suiWalletInfo := output.NetworkSpecificData.SuiAccount
		mnemonic = suiWalletInfo.Mnemonic
		suiAddress = suiWalletInfo.SuiAddress
		break
	}

	suiSigner, err := suisigner.NewSignerWithMnemonic(mnemonic, suisigner.KeySchemeFlagEd25519)
	require.NoError(t, err)
	suiPrivateKey := suiSigner.PrivateKey()

	client := suiclient.NewClient(url)

	var ready bool
	for i := 0; i < 30; i++ {
		time.Sleep(time.Second)
		receivedChainID, err := client.GetChainIdentifier(context.Background())
		if err != nil {
			t.Logf("API server not ready yet (attempt %d): %+v\n", i+1, err)
			continue
		}
		// we can't compare receivedChainID to chainID because it's generated from a new genesis block
		// checkpoint each time
		// TODO: could we keep the same genesis config each time when starting the container?
		t.Logf("Successfully fetched chain id: %s", receivedChainID)
		ready = true
		break
	}
	require.True(t, ready, "Sui network not ready")
	time.Sleep(15 * time.Second) // we have slot errors that force retries if the chain is not given enough time to boot
	return url, suiAddress, suiPrivateKey, client
}

func createSuiChainConfig(chainID string, chain suichain.Chain) chainlink.RawConfig {
	chainConfig := chainlink.RawConfig{}

	chainConfig["Enabled"] = true
	chainConfig["ChainID"] = chainID
	chainConfig["NetworkName"] = "sui-localnet"
	chainConfig["NetworkNameFull"] = "sui-localnet"
	chainConfig["Nodes"] = []any{
		map[string]any{
			"Name": "primary",
			"URL":  chain.URL,
		},
	}

	return chainConfig
}
