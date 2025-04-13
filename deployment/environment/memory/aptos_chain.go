package memory

import (
	//"crypto/secp256k1"
	//"github.com/decred/dcrd/dcrec/secp256k1/v4"
	//"github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"

	//"strconv"
	"strings"
	"testing"

	//"time"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/aptos-labs/aptos-go-sdk/crypto"

	"github.com/stretchr/testify/require"
	//"github.com/testcontainers/testcontainers-go"

	//chainsel "github.com/smartcontractkit/chain-selectors"

	//"github.com/smartcontractkit/chainlink-testing-framework/framework"
	//"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

type AptosChain struct {
	Client      *aptos.NodeClient
	DeployerKey *aptos.Account
}

func getTestAptosChainSelectors() []uint64 {
	// TODO: CTF to support different chain ids, need to investigate if it's possible (thru node config.yaml?)
	//return []uint64{chainsel.APTOS_LOCALNET.Selector}
	//return []uint64{4457093679053095498}
	return []uint64{743186221051783445}
}

func createAptosAccount(t *testing.T, useDefault bool) *aptos.Account {
	if useDefault {
		addressStr := blockchain.DefaultAptosAccount
		var defaultAddress aptos.AccountAddress
		err := defaultAddress.ParseStringRelaxed(addressStr)
		require.NoError(t, err)

		privateKeyStr := blockchain.DefaultAptosPrivateKey
		privateKeyBytes, err := hex.DecodeString(strings.TrimPrefix(privateKeyStr, "0x"))
		require.NoError(t, err)
		privateKey := ed25519.NewKeyFromSeed(privateKeyBytes)

		t.Logf("Using default Aptos account: %s %+v", addressStr, privateKeyBytes)

		account, err := aptos.NewAccountFromSigner(&crypto.Ed25519PrivateKey{Inner: privateKey}, defaultAddress)
		require.NoError(t, err)
		return account
	} else {
		account, err := aptos.NewEd25519SingleSenderAccount()
		require.NoError(t, err)
		return account
	}
}

func GenerateChainsAptos(t *testing.T, numChains int) map[uint64]AptosChain {
	testAptosChainSelectors := getTestAptosChainSelectors()
	if len(testAptosChainSelectors) < numChains {
		t.Fatalf("not enough test aptos chain selectors available")
	}
	chains := make(map[uint64]AptosChain)
	for i := 0; i < numChains; i++ {
		chainID := testAptosChainSelectors[i]
		account := createAptosAccount(t, true)

		nodeClient := aptosChain(t, chainID, account.Address)
		chains[chainID] = AptosChain{
			Client:      nodeClient,
			DeployerKey: account,
		}
	}
	t.Logf("Created %d Aptos chains: %+v", len(chains), chains)
	return chains
}

func aptosChain(t *testing.T, chainSelector uint64, adminAddress aptos.AccountAddress) *aptos.NodeClient {
	t.Helper()

	// initialize the docker network used by CTF
	//err := framework.DefaultNetwork(once)
	//require.NoError(t, err)

	//maxRetries := 10
	//var url string
	//var port uint16
	//var containerName string
	//for i := 0; i < maxRetries; i++ {
	//// TODO(aptos): update CTF to be able to use the selected port
	////port = freeport.GetOne(t)
	//port = 8080

	//bcInput := &blockchain.Input{
	//Image: "", // filled out by defaultAptos function
	//Type:  "aptos",
	//// TODO(aptors): this should be chain id not chain selector?
	//ChainID:   strconv.FormatUint(chainSelector, 10),
	//PublicKey: adminAddress.String(),
	//Port:      fmt.Sprintf("%d", port),
	//}
	//output, err := blockchain.NewBlockchainNetwork(bcInput)
	//if err != nil {
	//t.Logf("Error creating Aptos network: %v", err)
	//time.Sleep(time.Second)
	//maxRetries -= 1
	//continue
	//}
	//require.NoError(t, err)
	//containerName = output.ContainerName
	//testcontainers.CleanupContainer(t, output.Container)
	//url = "https://api.devnet.aptoslabs.com/v1"
	////url = output.Nodes[0].HostHTTPUrl + "/v1"
	//break
	//}
	url := "https://api.testnet.aptoslabs.com/v1"

	fmt.Printf("DEBUG: aptos chain url: %s\n", url)
	client, err := aptos.NewNodeClient(url, 0)
	require.NoError(t, err)

	//var ready bool
	//for i := 0; i < 30; i++ {
	//time.Sleep(time.Second)
	//_, err := client.GetChainId()
	//if err != nil {
	//t.Logf("API server not ready yet (attempt %d): %+v\n", i+1, err)
	//continue
	//}
	//ready = true
	//break
	//}
	//require.True(t, ready, "Aptos network not ready")
	//time.Sleep(15 * time.Second) // we have slot errors that force retries if the chain is not given enough time to boot

	// if we used the default account, it's already funded, but we need more to give to transmitters
	//_, err = framework.ExecContainer(containerName, []string{"aptos", "account", "fund-with-faucet", "--account", adminAddress.String(), "--amount", "100000000000"})
	//require.NoError(t, err)

	return client
}

func createAptosChainConfig(chainID string, chain deployment.AptosChain) chainlink.RawConfig {
	chainConfig := chainlink.RawConfig{}

	chainConfig["Enabled"] = true
	chainConfig["ChainID"] = chainID
	chainConfig["NetworkName"] = "testnet"
	chainConfig["NetworkNameFull"] = "aptos-testnet"
	chainConfig["Nodes"] = []any{
		map[string]any{
			"Name": "primary",
			// TODO(aptos): fill out URL correctly
			//"URL": "http://localhost:8080/v1",
			"URL": "https://api.testnet.aptoslabs.com/v1",
		},
	}

	return chainConfig
}
