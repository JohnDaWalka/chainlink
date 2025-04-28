package memory

import (
	//"crypto/secp256k1"
	//"github.com/decred/dcrd/dcrec/secp256k1/v4"
	//"github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/xssnick/tonutils-go/ton/wallet"

	"github.com/hashicorp/consul/sdk/freeport"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tonaddress "github.com/xssnick/tonutils-go/address"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/xssnick/tonutils-go/ton"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

type TonChain struct {
	Client         *ton.APIClient
	DeployerWallet *wallet.Wallet
}

func getTestTonChainSelectors() []uint64 {
	// TODO: CTF to support different chain ids, need to investigate if it's possible (thru node config.yaml?)
	return []uint64{chainsel.TON_LOCALNET.Selector}
}

func createTonWallet(t *testing.T, useDefault bool) *wallet.Wallet {
	// TON wallet contract version
	ver := wallet.V5R1Final

	if useDefault {
		addressStr := blockchain.DefaultTonAccount
		defaultAddress, err := tonaddress.ParseAddr(addressStr)
		require.NoError(t, err)

		privateKeyStr := blockchain.DefaultTonPrivateKey
		privateKeyBytes, err := hex.DecodeString(privateKeyStr)
		require.NoError(t, err)
		privateKey := ed25519.NewKeyFromSeed(privateKeyBytes)

		t.Logf("Using default Ton account: %s %+v", addressStr, privateKeyBytes)

		wallet, err := wallet.FromPrivateKey(nil, privateKey, ver)
		//account, err := ton.NewAccountFromSigner(&crypto.Secp256k1PrivateKey{Inner: privateKey}, defaultAddress)
		require.NoError(t, err)
		return wallet
	} else {
		wallet, err := wallet.FromSeed(nil, "", ver, false)
		require.NoError(t, err)
		return wallet
	}
}

func GenerateChainsTon(t *testing.T, numChains int) map[uint64]deployment.TonChain {
	testTonChainSelectors := getTestTonChainSelectors()
	if len(testTonChainSelectors) < numChains {
		t.Fatalf("not enough test ton chain selectors available")
	}
	chains := make(map[uint64]deployment.TonChain)
	for i := 0; i < numChains; i++ {
		chainID := testTonChainSelectors[i]
		wallet := createTonWallet(t, true)

		nodeClient := tonChain(t, chainID, *wallet.Address())
		chains[chainID] = deployment.TonChain{
			Client: nodeClient,
			Wallet: wallet,
		}
	}
	t.Logf("Created %d Ton chains: %+v", len(chains), chains)
	return chains
}

func tonChain(t *testing.T, chainID uint64, adminAddress tonaddress.Address) *ton.APIClient {
	t.Helper()

	// TODO(ton): integrate Ton into CTF (https://smartcontract-it.atlassian.net/browse/NONEVM-1685)
	// initialize the docker network used by CTF
	err := framework.DefaultNetwork(once)
	require.NoError(t, err)

	maxRetries := 10
	var url string
	var port uint16
	var containerName string
	for i := 0; i < maxRetries; i++ {
		port = uint16(freeport.GetOne(t))

		bcInput := &blockchain.Input{
			Image:     "", // filled out by defaultTon function
			Type:      "ton",
			ChainID:   strconv.FormatUint(chainID, 10),
			PublicKey: adminAddress.String(),
			Port:      fmt.Sprintf("%d", port),
		}
		output, err := blockchain.NewBlockchainNetwork(bcInput)
		if err != nil {
			t.Logf("Error creating Ton network: %v", err)
			time.Sleep(time.Second)
			maxRetries -= 1
			continue
		}
		require.NoError(t, err)
		containerName = output.ContainerName
		testcontainers.CleanupContainer(t, output.Container)
		url = output.Nodes[0].ExternalHTTPUrl + "/v1"
		break
	}

	fmt.Printf("DEBUG: ton chain url: %s\n", url)

	client, err := ton.NewAPIClient().NewNodeClient(url, 0)
	require.NoError(t, err)

	var ready bool
	for i := 0; i < 30; i++ {
		time.Sleep(time.Second)
		_, err := client.GetChainId()
		if err != nil {
			t.Logf("API server not ready yet (attempt %d): %+v\n", i+1, err)
			continue
		}
		ready = true
		break
	}
	require.True(t, ready, "Ton network not ready")
	time.Sleep(15 * time.Second) // we have slot errors that force retries if the chain is not given enough time to boot

	// if we used the default account, it's already funded, but we need more to give to transmitters
	_, err = framework.ExecContainer(containerName, []string{"ton", "account", "fund-with-faucet", "--account", adminAddress.String(), "--amount", "100000000000"})
	require.NoError(t, err)

	return client
}

func createTonChainConfig(chainID string, chain deployment.TonChain) chainlink.RawConfig {
	chainConfig := chainlink.RawConfig{}

	chainConfig["Enabled"] = true
	chainConfig["ChainID"] = chainID
	chainConfig["NetworkName"] = "localnet"
	chainConfig["NetworkNameFull"] = "ton-localnet"
	chainConfig["Nodes"] = []any{
		map[string]any{
			"Name": "primary",
			// TODO(ton): fill out URL correctly
			"URL": "http://localhost:8080/v1",
		},
	}

	return chainConfig
}
