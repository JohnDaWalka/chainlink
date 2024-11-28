package capabilities_test

import (
	"bytes"
	"cmp"
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"os"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink/integration-tests/capabilities/components/evmcontracts/capabilities_registry"
	"github.com/smartcontractkit/chainlink/integration-tests/capabilities/components/evmcontracts/forwarder"
	"github.com/smartcontractkit/chainlink/integration-tests/capabilities/components/onchain"

	cr_wrapper "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	feeds_consumer "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/feeds_consumer"
	ocr3_capability "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/ocr3_capability"
)

// Copying this to avoid dependency on the core repo
func GetChainType(chainType string) (uint8, error) {
	switch chainType {
	case "evm":
		return 1, nil
	// case Solana:
	// 	return 2, nil
	// case Cosmos:
	// 	return 3, nil
	// case StarkNet:
	// 	return 4, nil
	// case Aptos:
	// 	return 5, nil
	default:
		return 0, fmt.Errorf("unexpected chaintype.ChainType: %#v", chainType)
	}
}

// Copying this to avoid dependency on the core repo
func MarshalMultichainPublicKey(ost map[string]types.OnchainPublicKey) (types.OnchainPublicKey, error) {
	pubKeys := make([][]byte, 0, len(ost))
	for k, pubKey := range ost {
		typ, err := GetChainType(k)
		if err != nil {
			// skipping unknown key type
			continue
		}
		buf := new(bytes.Buffer)
		if err = binary.Write(buf, binary.LittleEndian, typ); err != nil {
			return nil, err
		}
		length := len(pubKey)
		if length < 0 || length > math.MaxUint16 {
			return nil, errors.New("pubKey doesn't fit into uint16")
		}
		if err = binary.Write(buf, binary.LittleEndian, uint16(length)); err != nil {
			return nil, err
		}
		_, _ = buf.Write(pubKey)
		pubKeys = append(pubKeys, buf.Bytes())
	}
	// sort keys based on encoded type to make encoding deterministic
	slices.SortFunc(pubKeys, func(a, b []byte) int { return cmp.Compare(a[0], b[0]) })
	return bytes.Join(pubKeys, nil), nil
}

type WorkflowTestConfig struct {
	BlockchainA *blockchain.Input `toml:"blockchain_a" validate:"required"`
	NodeSet     *ns.Input         `toml:"nodeset" validate:"required"`
}

type OCR3Config struct {
	Signers               [][]byte
	Transmitters          []common.Address
	F                     uint8
	OnchainConfig         []byte
	OffchainConfigVersion uint64
	OffchainConfig        []byte
}

type NodeInfo struct {
	OcrKeyBundleID            string
	TransmitterAddress        string
	PeerID                    string
	Signer                    common.Address
	OffchainPublicKey         [32]byte
	OnchainPublicKey          types.OnchainPublicKey
	ConfigEncryptionPublicKey [32]byte
}

func extractKey(value string) string {
	parts := strings.Split(value, "_")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return value
}

func getNodesInfo(
	t *testing.T,
	nodes []*clclient.ChainlinkClient,
) (nodesInfo []NodeInfo) {
	nodesInfo = make([]NodeInfo, len(nodes))

	for i, node := range nodes {
		// OCR Keys
		ocr2Keys, err := node.MustReadOCR2Keys()
		require.NoError(t, err)
		nodesInfo[i].OcrKeyBundleID = ocr2Keys.Data[0].ID

		firstOCR2Key := ocr2Keys.Data[0].Attributes
		nodesInfo[i].Signer = common.HexToAddress(extractKey(firstOCR2Key.OnChainPublicKey))

		pubKeys := make(map[string]types.OnchainPublicKey)
		ethOnchainPubKey, err := hex.DecodeString(extractKey(firstOCR2Key.OnChainPublicKey))
		require.NoError(t, err)
		pubKeys["evm"] = ethOnchainPubKey

		multichainPubKey, err := MarshalMultichainPublicKey(pubKeys)
		require.NoError(t, err)
		nodesInfo[i].OnchainPublicKey = multichainPubKey

		offchainPublicKeyBytes, err := hex.DecodeString(extractKey(firstOCR2Key.OffChainPublicKey))
		require.NoError(t, err)
		var offchainPublicKey [32]byte
		copy(offchainPublicKey[:], offchainPublicKeyBytes)
		nodesInfo[i].OffchainPublicKey = offchainPublicKey

		sharedSecretEncryptionPublicKeyBytes, err := hex.DecodeString(extractKey(firstOCR2Key.ConfigPublicKey))
		require.NoError(t, err)
		var sharedSecretEncryptionPublicKey [32]byte
		copy(sharedSecretEncryptionPublicKey[:], sharedSecretEncryptionPublicKeyBytes)
		nodesInfo[i].ConfigEncryptionPublicKey = sharedSecretEncryptionPublicKey

		// ETH Keys
		ethKeys, err := node.MustReadETHKeys()
		require.NoError(t, err)
		nodesInfo[i].TransmitterAddress = ethKeys.Data[0].Attributes.Address

		// P2P Keys
		p2pKeys, err := node.MustReadP2PKeys()
		require.NoError(t, err)
		nodesInfo[i].PeerID = p2pKeys.Data[0].Attributes.PeerID
	}

	return nodesInfo
}

func generateOCR3Config(
	t *testing.T,
	nodesInfo []NodeInfo,
) (config *OCR3Config) {
	oracleIdentities := []confighelper.OracleIdentityExtra{}
	transmissionSchedule := []int{}

	for _, nodeInfo := range nodesInfo {
		transmissionSchedule = append(transmissionSchedule, 1)
		oracleIdentity := confighelper.OracleIdentityExtra{}
		oracleIdentity.OffchainPublicKey = nodeInfo.OffchainPublicKey
		oracleIdentity.OnchainPublicKey = nodeInfo.OnchainPublicKey
		oracleIdentity.ConfigEncryptionPublicKey = nodeInfo.ConfigEncryptionPublicKey
		oracleIdentity.PeerID = nodeInfo.PeerID
		oracleIdentity.TransmitAccount = types.Account(nodeInfo.TransmitterAddress)
		oracleIdentities = append(oracleIdentities, oracleIdentity)
	}

	maxDurationInitialization := 10 * time.Second

	// Generate OCR3 configuration arguments for testing
	signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, err := ocr3confighelper.ContractSetConfigArgsForTests(
		20*time.Second,             // DeltaProgress: Time between rounds
		20*time.Second,             // DeltaResend: Time between resending unconfirmed transmissions
		20*time.Second,             // DeltaInitial: Initial delay before starting the first round
		6*time.Second,              // DeltaRound: Time between rounds within an epoch
		1*time.Second,              // DeltaGrace: Grace period for delayed transmissions
		4*time.Second,              // DeltaCertifiedCommitRequest: Time between certified commit requests
		60*time.Second,             // DeltaStage: Time between stages of the protocol
		uint64(10),                 // MaxRoundsPerEpoch: Maximum number of rounds per epoch
		transmissionSchedule,       // TransmissionSchedule: Transmission schedule
		oracleIdentities,           // Oracle identities with their public keys
		nil,                        // Plugin config (empty for now)
		&maxDurationInitialization, // MaxDurationInitialization: ???
		4*time.Second,              // MaxDurationQuery: Maximum duration for querying
		4*time.Second,              // MaxDurationObservation: Maximum duration for observation
		4*time.Second,              // MaxDurationAccept: Maximum duration for acceptance
		4*time.Second,              // MaxDurationTransmit: Maximum duration for transmission
		1,                          // F: Maximum number of faulty oracles
		nil,                        // OnChain config (empty for now)
	)
	require.NoError(t, err)

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
	}
}

func TestWorkflow(t *testing.T) {
	workflowOwner := "0x00000000000000000000000000000000000000aa"
	workflowName := "ccipethsep"
	feedID := "0x0003fbba4fce42f65d6032b18aee53efdf526cc734ad296cb57565979d883bdd"

	t.Run("Keystoen workflow test", func(t *testing.T) {
		in, err := framework.Load[WorkflowTestConfig](t)
		require.NoError(t, err)
		pkey := os.Getenv("PRIVATE_KEY")

		bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
		require.NoError(t, err)

		sc, err := seth.NewClientBuilder().
			WithRpcUrl(bc.Nodes[0].HostWSUrl).
			WithPrivateKeys([]string{pkey}).
			Build()
		require.NoError(t, err)

		capabilitiesRegistryInstance, err := capabilities_registry.Deploy(sc)
		require.NoError(t, err)

		require.NoError(t, capabilitiesRegistryInstance.AddCapabilities(
			[]cr_wrapper.CapabilitiesRegistryCapability{
				{
					LabelledName:   "mock-streams-trigger",
					Version:        "1.0.0",
					CapabilityType: 0, // TRIGGER
					ResponseType:   0, // REPORT
				},
				{
					LabelledName:   "offchain_reporting",
					Version:        "1.0.0",
					CapabilityType: 2, // CONSENSUS
					ResponseType:   0, // REPORT
				},
				{
					LabelledName:   "write_31337",
					Version:        "1.0.0",
					CapabilityType: 3, // TARGET
					ResponseType:   1, // OBSERVATION_IDENTICAL
				},
			},
		))

		forwarderInstance, err := forwarder.Deploy(sc)
		require.NoError(t, err)

		// TODO: When the capabilities registry address is provided:
		// - NOPs and nodes are added to the registry.
		// - Nodes are configured to listen to the registry for updates.
		nodeset, err := ns.NewSharedDBNodeSet(in.NodeSet, bc)
		require.NoError(t, err)

		nodeClients, err := clclient.New(nodeset.CLNodes)
		require.NoError(t, err)

		err = onchain.FundNodes(sc, nodeClients, pkey, 5)
		require.NoError(t, err)

		nodesInfo := getNodesInfo(t, nodeClients)

		bootstrapNodeInfo := nodesInfo[0]
		workflowNodesetInfo := nodesInfo[1:]

		for i := range workflowNodesetInfo {
			in.NodeSet.NodeSpecs[i].Node.TestConfigOverrides = fmt.Sprintf(`
				[Feature]
				LogPoller = true

				[OCR2]
				Enabled = true
				DatabaseTimeout = '1s'

				[P2P.V2]
				Enabled = true
				ListenAddresses = ['0.0.0.0:6690']

				# This is needed for the target capability to be initialized
				[[EVM]]
				ChainID = '%s'

				[[EVM.Nodes]]
				Name = 'anvil'
				WSURL = '%s'
				HTTPURL = '%s'

				[EVM.Workflow]
				FromAddress = '%s'
				ForwarderAddress = '%s'
				GasLimitDefault = 400_000

				# This is needed for external registry
				[Capabilities.ExternalRegistry]
				Address = '%s'
				NetworkID = 'evm'
				ChainID = '%s'
			`,
				bc.ChainID,
				bc.Nodes[0].DockerInternalWSUrl,
				bc.Nodes[0].DockerInternalHTTPUrl,
				nodesInfo[i].TransmitterAddress,
				forwarderInstance.Address,
				capabilitiesRegistryInstance.Address,
				bc.ChainID,
			)
		}

		nodeset, err = ns.UpgradeNodeSet(in.NodeSet, bc, 5*time.Second)
		require.NoError(t, err)
		nodeClients, err = clclient.New(nodeset.CLNodes)
		require.NoError(t, err)

		ocr3CapabilityAddress, tx, ocr3CapabilityContract, err := ocr3_capability.DeployOCR3Capability(
			sc.NewTXOpts(),
			sc.Client,
		)
		require.NoError(t, err)
		_, err = bind.WaitMined(context.Background(), sc.Client, tx)
		require.NoError(t, err)
		fmt.Println("Deployed ocr3_capability contract at", ocr3CapabilityAddress.Hex())

		feedsConsumerAddress, tx, feedsConsumerContract, err := feeds_consumer.DeployKeystoneFeedsConsumer(
			sc.NewTXOpts(),
			sc.Client,
		)
		require.NoError(t, err)
		_, err = bind.WaitMined(context.Background(), sc.Client, tx)
		require.NoError(t, err)
		fmt.Println("Deployed feeds_consumer contract at", feedsConsumerAddress.Hex())

		var workflowNameBytes [10]byte
		copy(workflowNameBytes[:], []byte(workflowName))

		tx, err = feedsConsumerContract.SetConfig(
			sc.NewTXOpts(),
			[]common.Address{forwarderInstance.Address},
			[]common.Address{common.HexToAddress(workflowOwner)},
			[][10]byte{workflowNameBytes},
		)
		require.NoError(t, err)
		_, err = bind.WaitMined(context.Background(), sc.Client, tx)
		require.NoError(t, err)

		// Add bootstrap spec to the first node
		bootstrapNode := nodeClients[0]

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
				providerType = "ocr3-capability"
			`, ocr3CapabilityAddress, bc.ChainID)
			r, _, err2 := bootstrapNode.CreateJobRaw(bootstrapJobSpec)
			assert.NoError(t, err2)
			assert.Empty(t, r.Errors)
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
					command="/home/capabilities/streams-linux-amd64"
				`
				response, _, err2 := nodeClient.CreateJobRaw(scJobSpec)
				assert.NoError(t, err2)
				assert.Empty(t, response.Errors)

				consensusJobSpec := fmt.Sprintf(`
					type = "offchainreporting2"
					schemaVersion = 1
					name = "Keystone OCR3 Consensus Capability"
					contractID = "%s"
					ocrKeyBundleID = "%s"
					p2pv2Bootstrappers = [
						"%s@%s",
					]
					relay = "evm"
					pluginType = "plugin"
					transmitterID = "%s"

					[relayConfig]
					chainID = "%s"

					[pluginConfig]
					command = "/usr/local/bin/chainlink-ocr3-capability"
					ocrVersion = 3
					pluginName = "ocr-capability"
					providerType = "ocr3-capability"
					telemetryType = "plugin"

					[onchainSigningStrategy]
					strategyName = 'multi-chain'
					[onchainSigningStrategy.config]
					evm = "%s"
					`,
					ocr3CapabilityAddress,
					nodesInfo[i].OcrKeyBundleID,
					bootstrapNodeInfo.PeerID,
					strings.TrimPrefix(nodeset.CLNodes[0].Node.DockerP2PUrl, "http://"),
					nodesInfo[i].TransmitterAddress,
					bc.ChainID,
					nodesInfo[i].OcrKeyBundleID,
				)
				fmt.Println("consensusJobSpec", consensusJobSpec)
				response, _, err2 = nodeClient.CreateJobRaw(consensusJobSpec)
				assert.NoError(t, err2)
				assert.Empty(t, response.Errors)

				workflowSpec := fmt.Sprintf(`
type = "workflow"
schemaVersion = 1
name = "Keystone CCIP Feeds Workflow"
forwardingAllowed = false
workflow = """
name: %s
owner: '%s'
triggers:
  - id: mock-streams-trigger@1.0.0
    config:
      maxFrequencyMs: 15000
      feedIds:
        - '%s'
consensus:
  - id: offchain_reporting@1.0.0
    ref: ccip_feeds
    inputs:
      observations:
        - $(trigger.outputs)
    config:
      report_id: '0001'
      key_id: evm
      aggregation_method: data_feeds
      aggregation_config:
        allowedPartialStaleness: '0.5'
        feeds:
          '%s':
            deviation: '0.05'
            heartbeat: 3600
            remappedID: '0x666666666666'
      encoder: EVM
      encoder_config:
        abi: '(bytes32 FeedID, uint224 Price, uint32 Timestamp)[] Reports'
targets:
  - id: write_%s@1.0.0
    inputs:
      signed_report: $(ccip_feeds.outputs)
    config:
      address: '%s'
      deltaStage: 45s
      schedule: oneAtATime
"""`,
					workflowName,
					workflowOwner,
					feedID,
					feedID,
					bc.ChainID,
					feedsConsumerAddress,
				)
				response, _, err2 = nodeClient.CreateJobRaw(workflowSpec)
				assert.NoError(t, err2)
				assert.Empty(t, response.Errors)
			}()
		}
		wg.Wait()

		var nopsToAdd []cr_wrapper.CapabilitiesRegistryNodeOperator
		var nodesToAdd []cr_wrapper.CapabilitiesRegistryNodeParams
		var donNodes [][32]byte
		var signers []common.Address

		for i, node := range nodesInfo {
			if i == 0 {
				continue
			}
			nopsToAdd = append(nopsToAdd, cr_wrapper.CapabilitiesRegistryNodeOperator{
				Admin: common.HexToAddress(node.TransmitterAddress),
				Name:  fmt.Sprintf("NOP %d", i),
			})

			var peerID ragetypes.PeerID
			err = peerID.UnmarshalText([]byte(node.PeerID))
			require.NoError(t, err)

			nodesToAdd = append(nodesToAdd, cr_wrapper.CapabilitiesRegistryNodeParams{
				NodeOperatorId:      uint32(i), //nolint:gosec // disable G115
				Signer:              common.BytesToHash(node.Signer.Bytes()),
				P2pId:               peerID,
				EncryptionPublicKey: [32]byte{1, 2, 3, 4, 5},
				HashedCapabilityIds: capabilitiesRegistryInstance.ExistingHashedCapabilitiesIDs,
			})

			donNodes = append(donNodes, peerID)
			signers = append(signers, node.Signer)
		}

		// Add NOPs to registry
		tx, err = capabilitiesRegistryInstance.Contract.AddNodeOperators(
			sc.NewTXOpts(),
			nopsToAdd,
		)
		require.NoError(t, err)
		_, err = bind.WaitMined(context.Background(), sc.Client, tx)
		require.NoError(t, err)

		// Add nodes to registry
		tx, err = capabilitiesRegistryInstance.Contract.AddNodes(
			sc.NewTXOpts(),
			nodesToAdd,
		)
		require.NoError(t, err)
		_, err = bind.WaitMined(context.Background(), sc.Client, tx)
		require.NoError(t, err)

		// Add nodeset to registry
		tx, err = capabilitiesRegistryInstance.Contract.AddDON(
			sc.NewTXOpts(),
			donNodes,
			[]cr_wrapper.CapabilitiesRegistryCapabilityConfiguration{
				{
					CapabilityId: capabilitiesRegistryInstance.ExistingHashedCapabilitiesIDs[0],
					Config:       []byte(""),
				},
				{
					CapabilityId: capabilitiesRegistryInstance.ExistingHashedCapabilitiesIDs[1],
					Config:       []byte(""),
				},
				{
					CapabilityId: capabilitiesRegistryInstance.ExistingHashedCapabilitiesIDs[2],
					Config:       []byte(""),
				},
			},
			true,
			true,
			uint8(1),
		)
		require.NoError(t, err)
		_, err = bind.WaitMined(context.Background(), sc.Client, tx)
		require.NoError(t, err)

		require.NoError(t, forwarderInstance.SetConfig(
			1,
			1,
			1,
			signers,
		))

		// Configure OCR capability contract
		ocr3Config := generateOCR3Config(t, workflowNodesetInfo)
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

		// It can take a while before the first report is produced, particularly on CI.
		timeout := 5 * time.Minute
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		startTime := time.Now()
		for {
			select {
			case <-ctx.Done():
				t.Fatalf("feed did not update, timeout after %s", timeout)
			case <-time.After(5 * time.Second):
				elapsed := time.Since(startTime).Round(time.Second)
				price, _, err := feedsConsumerContract.GetPrice(
					sc.NewCallOpts(),
					common.HexToHash(feedID),
				)
				require.NoError(t, err)

				if price.String() != "0" {
					fmt.Printf("Feed updated after %s - price set, price=%s\n", elapsed, price)
					return
				}
				fmt.Printf("Feed not updated yet, waiting for %s\n", elapsed)
			}
		}
	})
}
