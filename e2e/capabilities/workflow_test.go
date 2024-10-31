package capabilities_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	capabilities_registry "github.com/smartcontractkit/chainlink/e2e/capabilities/components/evmcontracts/capabilities_registry_1_1_0"
	feeds_consumer "github.com/smartcontractkit/chainlink/e2e/capabilities/components/evmcontracts/feeds_consumer_1_0_0"
	forwarder "github.com/smartcontractkit/chainlink/e2e/capabilities/components/evmcontracts/forwarder_1_0_0"
	ocr3_capability "github.com/smartcontractkit/chainlink/e2e/capabilities/components/evmcontracts/ocr3_capability_1_0_0"
)

type WorkflowTestConfig struct {
	BlockchainA        *blockchain.Input `toml:"blockchain_a" validate:"required"`
	MockerDataProvider *fake.Input       `toml:"data_provider" validate:"required"`
	NodeSet            *ns.Input         `toml:"nodeset" validate:"required"`
}

type OCR3Config struct {
	Signers               [][]byte
	Transmitters          []common.Address
	F                     uint8
	OnchainConfig         []byte
	OffchainConfigVersion uint64
	OffchainConfig        []byte
}

func extractKey(value string) string {
	parts := strings.Split(value, "_")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return value
}

func generateOCR3Config(t *testing.T, nodes []*clclient.ChainlinkClient) (*OCR3Config, error) {
	oracleIdentities := []confighelper.OracleIdentityExtra{}
	transmissionSchedule := []int{}

	for i, node := range nodes {
		// TODO: Do not provide a bootstrap node to this func
		// We want to skip bootstrap node.
		if i == 0 {
			continue
		}
		transmissionSchedule = append(transmissionSchedule, 0)
		oracleIdentity := confighelper.OracleIdentityExtra{}
		// ocr2
		ocr2Keys, err := node.MustReadOCR2Keys()
		require.NoError(t, err)

		firstOCR2Key := ocr2Keys.Data[0].Attributes

		offchainPublicKeyBytes, err := hex.DecodeString(extractKey(firstOCR2Key.OffChainPublicKey))
		require.NoError(t, err)
		var offchainPublicKey [32]byte
		copy(offchainPublicKey[:], offchainPublicKeyBytes)
		oracleIdentity.OffchainPublicKey = offchainPublicKey

		onchainPubkey, err := hex.DecodeString(extractKey(firstOCR2Key.OnChainPublicKey))
		require.NoError(t, err)
		oracleIdentity.OnchainPublicKey = onchainPubkey

		sharedSecretEncryptionPublicKeyBytes, err := hex.DecodeString(extractKey(firstOCR2Key.ConfigPublicKey))
		require.NoError(t, err)
		var sharedSecretEncryptionPublicKey [32]byte
		copy(sharedSecretEncryptionPublicKey[:], sharedSecretEncryptionPublicKeyBytes)
		oracleIdentity.ConfigEncryptionPublicKey = sharedSecretEncryptionPublicKey

		// p2p
		p2pKeys, err := node.MustReadP2PKeys()
		require.NoError(t, err)
		oracleIdentity.PeerID = p2pKeys.Data[0].Attributes.PeerID

		// eth
		ethKeys, err := node.MustReadETHKeys()
		require.NoError(t, err)
		oracleIdentity.TransmitAccount = types.Account(ethKeys.Data[0].Attributes.Address)

		oracleIdentities = append(oracleIdentities, oracleIdentity)
	}

	maxDurationInitialization := 10 * time.Second

	// Generate OCR3 configuration arguments for testing
	signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, err := ocr3confighelper.ContractSetConfigArgsForTests(
		5*time.Second,              // DeltaProgress: Time between rounds
		5*time.Second,              // DeltaResend: Time between resending unconfirmed transmissions
		5*time.Second,              // DeltaInitial: Initial delay before starting the first round
		2*time.Second,              // DeltaRound: Time between rounds within an epoch
		500*time.Millisecond,       // DeltaGrace: Grace period for delayed transmissions
		1*time.Second,              // DeltaCertifiedCommitRequest: Time between certified commit requests
		30*time.Second,             // DeltaStage: Time between stages of the protocol
		uint64(10),                 // MaxRoundsPerEpoch: Maximum number of rounds per epoch
		transmissionSchedule,       // TransmissionSchedule: Transmission schedule
		oracleIdentities,           // Oracle identities with their public keys
		nil,                        // Plugin config (empty for now)
		&maxDurationInitialization, // MaxDurationInitialization: ???
		1*time.Second,              // MaxDurationQuery: Maximum duration for querying
		1*time.Second,              // MaxDurationObservation: Maximum duration for observation
		1*time.Second,              // MaxDurationAccept: Maximum duration for acceptance
		1*time.Second,              // MaxDurationTransmit: Maximum duration for transmission
		1,                          // F: Maximum number of faulty oracles
		nil,                        // OnChain config (empty for now)
	)
	require.NoError(t, err)

	// maxDurationInitialization *time.Duration,
	// maxDurationQuery time.Duration,
	// maxDurationObservation time.Duration,
	// maxDurationShouldAcceptAttestedReport time.Duration,
	// maxDurationShouldTransmitAcceptedReport time.Duration,

	signerAddresses := [][]byte{}
	for _, signer := range signers {
		signerAddresses = append(signerAddresses, signer)
	}

	transmitterAddresses := []common.Address{}
	for _, transmitter := range transmitters {
		transmitterAddresses = append(transmitterAddresses, common.HexToAddress(string(transmitter)))
	}

	return &OCR3Config{
		Signers:               signerAddresses,
		Transmitters:          transmitterAddresses,
		F:                     f,
		OnchainConfig:         onchainConfig,
		OffchainConfigVersion: offchainConfigVersion,
		OffchainConfig:        offchainConfig,
	}, nil
}

func TestWorkflow(t *testing.T) {
	t.Run("smoke test", func(t *testing.T) {
		in, err := framework.Load[WorkflowTestConfig](t)
		require.NoError(t, err)
		pkey := os.Getenv("PRIVATE_KEY")

		// deploy docker test environment
		bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
		require.NoError(t, err)

		// connect clients
		sc, err := seth.NewClientBuilder().
			WithRpcUrl(bc.Nodes[0].HostWSUrl).
			WithPrivateKeys([]string{pkey}).
			Build()
		require.NoError(t, err)

		capabilitiesRegistryAddress, tx, _, err := capabilities_registry.DeployCapabilitiesRegistry(
			sc.NewTXOpts(),
			sc.Client,
		)
		require.NoError(t, err)
		_, err = bind.WaitMined(context.Background(), sc.Client, tx)
		require.NoError(t, err)
		fmt.Println("Deployed capabilities_registry contract at", capabilitiesRegistryAddress.Hex())

		// TODO: When the capabilities registry address is provided:
		// - NOPs and nodes are added to the registry.
		// - Nodes are configured to listen to the registry for updates.
		nodeset, err := ns.NewSharedDBNodeSet(in.NodeSet, bc, "https://example.com") // TODO: Should not be a thing
		require.NoError(t, err)

		for i, n := range nodeset.CLNodes {
			fmt.Printf("Node %d --> %s\n", i, n.Node.HostURL)
			fmt.Printf("Node P2P %d --> %s\n", i, n.Node.HostP2PURL)
		}

		nodeClients, err := clclient.NewCLDefaultClients(nodeset.CLNodes, framework.L)
		require.NoError(t, err)

		ocr3CapabilityAddress, tx, ocr3CapabilityContract, err := ocr3_capability.DeployOCR3Capability(
			sc.NewTXOpts(),
			sc.Client,
		)
		require.NoError(t, err)
		_, err = bind.WaitMined(context.Background(), sc.Client, tx)
		require.NoError(t, err)
		fmt.Println("Deployed ocr3_capability contract at", ocr3CapabilityAddress.Hex())

		forwarderAddress, tx, _, err := forwarder.DeployKeystoneForwarder(
			sc.NewTXOpts(),
			sc.Client,
		)
		require.NoError(t, err)
		_, err = bind.WaitMined(context.Background(), sc.Client, tx)
		require.NoError(t, err)
		fmt.Println("Deployed forwarder contract at", forwarderAddress.Hex())

		feedsConsumerAddress, tx, _, err := feeds_consumer.DeployKeystoneFeedsConsumer(
			sc.NewTXOpts(),
			sc.Client,
		)
		require.NoError(t, err)
		_, err = bind.WaitMined(context.Background(), sc.Client, tx)
		require.NoError(t, err)
		fmt.Println("Deployed feeds_consumer contract at", feedsConsumerAddress.Hex())

		// Add bootstrap spec to the first node
		bootstrapNode := nodeClients[0]
		_, err = bootstrapNode.MustReadP2PKeys()
		require.NoError(t, err)
		fmt.Println("P2P keys fetched")
		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			bootstrapJobSpec := fmt.Sprintf(`
				type = "bootstrap"
				schemaVersion = 1
				name = "Botostrap"
				contractID = "%s"
				contractConfigTrackerPollInterval = "1s"
				contractConfigConfirmations = 1
				relay = "evm"
		
				[relayConfig]
				chainID = %s
			`, ocr3CapabilityAddress, bc.ChainID)
			fmt.Println("Creating bootstrap job spec", bootstrapJobSpec)
			r, _, err2 := bootstrapNode.CreateJobRaw(bootstrapJobSpec)
			require.NoError(t, err2)
			require.Equal(t, len(r.Errors), 0)
			fmt.Printf("Response from bootstrap node: %x\n", r)
		}()

		for i, nodeClient := range nodeClients {
			// First node is a bootstrap node, so we skip it
			if i == 0 {
				continue
			}

			wg.Add(1)
			go func() {
				defer wg.Done()
				scJobSpec := `
				type = "standardcapabilities"
				schemaVersion = 1
				name = "streams-capabilities"
				command="/streams"`
				fmt.Println("Creating standard capabilities job spec", scJobSpec)
				response, _, err2 := nodeClient.CreateJobRaw(scJobSpec)
				require.NoError(t, err2)
				require.Equal(t, len(response.Errors), 0)
				fmt.Printf("Response from node %d: %x\n", i+1, response)
			}()
		}
		wg.Wait()

		ocr3Config, err := generateOCR3Config(t, nodeClients)
		require.NoError(t, err)
		fmt.Println("ocr3Config", ocr3Config)

		// Configure KV store OCR contract
		tx, err = ocr3CapabilityContract.SetConfig(
			sc.NewTXOpts(),
			ocr3Config.Signers,
			ocr3Config.Transmitters,
			ocr3Config.F,
			ocr3Config.OnchainConfig,
			ocr3Config.OffchainConfigVersion,
			ocr3Config.OffchainConfig,
		)
		require.NoError(t, err)
		_, err = bind.WaitMined(context.Background(), sc.Client, tx)
		require.NoError(t, err)

		// Add bootstrap spec
		// ✅ 1. Deploy mock streams capability
		// ✅ 2. Add boostrap job spec
		// 3. Add OCR3 capability
		// 4. Add chain write capabilities
	})
}
