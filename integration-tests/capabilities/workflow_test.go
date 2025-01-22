package capabilities_test

import (
	"bytes"
	"cmp"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-yaml/yaml"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	geth_types "github.com/ethereum/go-ethereum/core/types"
	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"

	cr_wrapper "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/feeds_consumer"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/forwarder"
	ocr3_capability "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/ocr3_capability"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/workflow/generated/workflow_registry_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	workflow_registry_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/workflowregistry"
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

type WorkflowConfig struct {
	UseChainlinkCLI bool                    `toml:"use_chainlink_cli"`
	ChainlinkCLI    *ChainlinkCLIConfig     `toml:"chainlink_cli"`
	UseExising      bool                    `toml:"use_existing"`
	Existing        *ExistingWorkflowConfig `toml:"existing"`
}

type ExistingWorkflowConfig struct {
	BinaryURL string `toml:"binary_url"`
	ConfigURL string `toml:"config_url"`
}

type ChainlinkCLIConfig struct {
	FolderLocation *string `toml:"folder_location"`
	Compile        bool    `toml:"compile"`
	ExistingWorkflowConfig
}

type WorkflowTestConfig struct {
	BlockchainA    *blockchain.Input `toml:"blockchain_a" validate:"required"`
	NodeSet        *ns.Input         `toml:"nodeset" validate:"required"`
	WorkflowConfig *WorkflowConfig   `toml:"workflow_config" validate:"required"`
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
		10*time.Second,             // DeltaResend: Time between resending unconfirmed transmissions
		1*time.Second,              // DeltaInitial: Initial delay before starting the first round
		5*time.Second,              // DeltaRound: Time between rounds within an epoch
		1*time.Second,              // DeltaGrace: Grace period for delayed transmissions
		5*time.Second,              // DeltaCertifiedCommitRequest: Time between certified commit requests
		10*time.Second,             // DeltaStage: Time between stages of the protocol
		uint64(10),                 // MaxRoundsPerEpoch: Maximum number of rounds per epoch
		transmissionSchedule,       // TransmissionSchedule: Transmission schedule
		oracleIdentities,           // Oracle identities with their public keys
		nil,                        // Plugin config (empty for now)
		&maxDurationInitialization, // MaxDurationInitialization: ???
		5*time.Second,              // MaxDurationQuery: Maximum duration for querying
		5*time.Second,              // MaxDurationObservation: Maximum duration for observation
		5*time.Second,              // MaxDurationAccept: Maximum duration for acceptance
		5*time.Second,              // MaxDurationTransmit: Maximum duration for transmission
		1,                          // F: Maximum number of faulty oracles
		nil,                        // OnChain config (empty for now)
	)
	require.NoError(t, err)

	// // values supplied by Alexandr Y
	// signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, err := ocr3confighelper.ContractSetConfigArgsForTests(
	// 	5*time.Second,              // DeltaProgress: Time between rounds
	// 	5*time.Second,              // DeltaResend: Time between resending unconfirmed transmissions
	// 	5*time.Second,              // DeltaInitial: Initial delay before starting the first round
	// 	2*time.Second,              // DeltaRound: Time between rounds within an epoch
	// 	500*time.Millisecond,       // DeltaGrace: Grace period for delayed transmissions
	// 	1*time.Second,              // DeltaCertifiedCommitRequest: Time between certified commit requests
	// 	30*time.Second,             // DeltaStage: Time between stages of the protocol
	// 	uint64(10),                 // MaxRoundsPerEpoch: Maximum number of rounds per epoch
	// 	transmissionSchedule,       // TransmissionSchedule: Transmission schedule
	// 	oracleIdentities,           // Oracle identities with their public keys
	// 	nil,                        // Plugin config (empty for now)
	// 	&maxDurationInitialization, // MaxDurationInitialization: ???
	// 	1*time.Second,              // MaxDurationQuery: Maximum duration for querying
	// 	1*time.Second,              // MaxDurationObservation: Maximum duration for observation
	// 	1*time.Second,              // MaxDurationAccept: Maximum duration for acceptance
	// 	1*time.Second,              // MaxDurationTransmit: Maximum duration for transmission
	// 	1,                          // F: Maximum number of faulty oracles
	// 	nil,                        // OnChain config (empty for now)
	// )
	// require.NoError(t, err)

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

func GenerateWorkflowIDFromStrings(owner string, name string, workflow []byte, config []byte, secretsURL string) (string, error) {
	ownerWithoutPrefix := owner
	if strings.HasPrefix(owner, "0x") {
		ownerWithoutPrefix = owner[2:]
	}

	ownerb, err := hex.DecodeString(ownerWithoutPrefix)
	if err != nil {
		return "", err
	}

	wid, err := GenerateWorkflowID(ownerb, name, workflow, config, secretsURL)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(wid[:]), nil
}

var (
	versionByte = byte(0)
)

func GenerateWorkflowID(owner []byte, name string, workflow []byte, config []byte, secretsURL string) ([32]byte, error) {
	s := sha256.New()
	_, err := s.Write(owner)
	if err != nil {
		return [32]byte{}, err
	}
	_, err = s.Write([]byte(name))
	if err != nil {
		return [32]byte{}, err
	}
	_, err = s.Write(workflow)
	if err != nil {
		return [32]byte{}, err
	}
	_, err = s.Write([]byte(config))
	if err != nil {
		return [32]byte{}, err
	}
	_, err = s.Write([]byte(secretsURL))
	if err != nil {
		return [32]byte{}, err
	}

	sha := [32]byte(s.Sum(nil))
	sha[0] = versionByte

	return sha, nil
}

func isInstalled(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func downloadAndDecode(url string) ([]byte, error) {
	// Step 1: Make an HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	// Step 2: Check the HTTP response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	// Step 3: Read the response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Step 4: Decode the base64 content
	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 content: %w", err)
	}

	return decoded, nil
}

type ChainlinkCliSettings struct {
	DevPlatform  DevPlatform  `yaml:"dev-platform"`
	UserWorkflow UserWorkflow `yaml:"user-workflow"`
	Logging      Logging      `yaml:"logging"`
	McmsConfig   McmsConfig   `yaml:"mcms-config"`
	Contracts    Contracts    `yaml:"contracts"`
	Rpcs         []Rpc        `yaml:"rpcs"`
}

type DevPlatform struct {
	CapabilitiesRegistryAddress string `yaml:"capabilities-registry-contract-address"`
	DonId                       uint32 `yaml:"don-id"`
	WorkflowRegistryAddress     string `yaml:"workflow-registry-contract-address"`
}

type UserWorkflow struct {
	WorkflowOwnerAddress string `yaml:"workflow-owner-address"`
}

type Logging struct {
	SethConfigPath string `yaml:"seth-config-path"`
}

type McmsConfig struct {
	ProposalsDirectory string `yaml:"proposals-directory"`
}

type Contracts struct {
	ContractRegistry []ContractRegistry `yaml:"registries"`
}

type ContractRegistry struct {
	Name          string `yaml:"name"`
	Address       string `yaml:"address"`
	ChainSelector uint64 `yaml:"chain-selector"`
}

type Rpc struct {
	ChainSelector uint64 `yaml:"chain-selector"`
	URL           string `yaml:"url"`
}

type PoRWorkflowConfig struct {
	FeedID          string `json:"feed_id"`
	URL             string `json:"url"`
	ConsumerAddress string `json:"consumer_address"`
}

func TestWorkflow(t *testing.T) {
	// without 0x prefix!
	feedID := "018bfe8840700040000000000000000000000000000000000000000000000000"
	feedBytes := common.HexToHash(feedID)

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

		if in.WorkflowConfig.UseChainlinkCLI {
			require.True(t, isInstalled("chainlink-cli"), "chainlink-cli is required for this test. Please install it, add to path and run again")
			require.NotEmpty(t, os.Getenv("GITHUB_API_TOKEN"), "GITHUB_API_TOKEN must be set to use chainlink-cli. It requires read/write Gist permissions")
			err := os.Setenv("WORKFLOW_OWNER_ADDRESS", sc.MustGetRootKeyAddress().Hex())
			require.NoError(t, err)

			err = os.Setenv("ETH_PRIVATE_KEY", pkey)
			require.NoError(t, err)
		}

		lgr := logger.TestLogger(t)
		require.NoError(t, err)
		addressBook := deployment.NewMemoryAddressBook()
		chainMap := make(map[uint64]deployment.Chain)
		ctx := context.Background()

		chainSelector, err := chainselectors.SelectorFromChainId(sc.Cfg.Network.ChainID)
		require.NoError(t, err)
		chainMap[chainSelector] = deployment.Chain{
			Selector:    chainSelector,
			Client:      sc.Client,
			DeployerKey: sc.NewTXOpts(seth.WithNonce(nil)), //set nonce to nil, so it will be fetched from the chain
			Confirm: func(tx *geth_types.Transaction) (uint64, error) {
				decoded, revertErr := sc.DecodeTx(tx)
				if revertErr != nil {
					return 0, revertErr
				}
				if decoded.Receipt == nil {
					return 0, fmt.Errorf("no receipt found for transaction %s even though it wasn't reverted. This should not happen", tx.Hash().String())
				}
				return decoded.Receipt.BlockNumber.Uint64(), nil
			},
		}

		ctfEnv := deployment.NewEnvironment("ctfV2", lgr, addressBook, chainMap, nil, nil, nil, func() context.Context { return ctx }, deployment.OCRSecrets{})

		// output, err := keystone_changeset.DeployCapabilityRegistry(*ctfEnv, chainSelector)
		// require.NoError(t, err)

		capRegAddr, tx, capabilitiesRegistryInstance, err := cr_wrapper.DeployCapabilitiesRegistry(sc.NewTXOpts(), sc.Client)
		_, decodeErr := sc.Decode(tx, err)
		require.NoError(t, decodeErr)

		allCaps := []cr_wrapper.CapabilitiesRegistryCapability{
			{
				LabelledName:   "offchain_reporting",
				Version:        "1.0.0",
				CapabilityType: 2, // CONSENSUS
				ResponseType:   0, // REPORT
			},
			{
				LabelledName:   "write_geth-testnet",
				Version:        "1.0.0",
				CapabilityType: 3, // TARGET
				ResponseType:   1, // OBSERVATION_IDENTICAL
			},
			{
				LabelledName:   "cron-trigger",
				Version:        "1.0.0",
				CapabilityType: uint8(0), // trigger
			},
			{
				LabelledName:   "custom-compute",
				Version:        "1.0.0",
				CapabilityType: uint8(1), // action
			},
		}

		tx, err = capabilitiesRegistryInstance.AddCapabilities(
			sc.NewTXOpts(),
			allCaps,
		)
		_, decodeErr = sc.Decode(tx, err)
		require.NoError(t, decodeErr)

		var hashedCapabilities [][32]byte

		for _, capability := range allCaps {
			hashed, err := capabilitiesRegistryInstance.GetHashedCapabilityId(
				sc.NewCallOpts(),
				capability.LabelledName,
				capability.Version,
			)
			require.NoError(t, err)
			hashedCapabilities = append(hashedCapabilities, hashed)
		}

		output, err := keystone_changeset.DeployForwarder(*ctfEnv, keystone_changeset.DeployForwarderRequest{
			ChainSelectors: []uint64{chainSelector},
		})
		require.NoError(t, err)

		addresses, err := output.AddressBook.AddressesForChain(chainSelector)
		require.NoError(t, err)

		var forwarderAddress common.Address
		for addrStr, tv := range addresses {
			fmt.Println("Address: ", addrStr)
			fmt.Println("Type and version: ", tv.String())
			if strings.Contains(tv.String(), "KeystoneForwarder") {
				forwarderAddress = common.HexToAddress(addrStr)
				break
			}
		}

		donID := uint32(1)
		workflowName := "abcdefgasd"

		output, err = workflow_registry_changeset.Deploy(*ctfEnv, chainSelector)
		require.NoError(t, err)

		addresses, err = output.AddressBook.AddressesForChain(chainSelector)
		require.NoError(t, err)

		var workflowRegistryAddr common.Address
		for addrStr, tv := range addresses {
			fmt.Println("Address: ", addrStr)
			fmt.Println("Type and version: ", tv.String())
			if strings.Contains(tv.String(), "WorkflowRegistry") {
				workflowRegistryAddr = common.HexToAddress(addrStr)
			}
		}

		// TODO really? why do I have to update that manually?
		ctfEnv.ExistingAddresses = output.AddressBook

		_, err = workflow_registry_changeset.UpdateAllowedDons(*ctfEnv, &workflow_registry_changeset.UpdateAllowedDonsRequest{
			RegistryChainSel: chainSelector,
			DonIDs:           []uint32{donID},
			Allowed:          true,
		})
		require.NoError(t, err)

		_, err = workflow_registry_changeset.UpdateAuthorizedAddresses(*ctfEnv, &workflow_registry_changeset.UpdateAuthorizedAddressesRequest{
			RegistryChainSel: chainSelector,
			Addresses:        []string{sc.MustGetRootKeyAddress().Hex()},
			Allowed:          true,
		})
		require.NoError(t, err)

		output, err = keystone_changeset.DeployFeedsConsumer(*ctfEnv, &keystone_changeset.DeployFeedsConsumerRequest{
			ChainSelector: chainSelector,
		})
		require.NoError(t, err)

		addresses, err = output.AddressBook.AddressesForChain(chainSelector)
		require.NoError(t, err)

		var feedsConsumerAddress common.Address
		for addrStr, tv := range addresses {
			fmt.Println("Address: ", addrStr)
			fmt.Println("Type and version: ", tv.String())
			if strings.Contains(tv.String(), "FeedConsumer") {
				feedsConsumerAddress = common.HexToAddress(addrStr)
			}
		}

		fmt.Println("Deployed feeds_consumer contract at", feedsConsumerAddress.Hex())

		var workflowNameBytes [10]byte
		var HashTruncateName = func(name string) string {
			// Compute SHA-256 hash of the input string
			hash := sha256.Sum256([]byte(name))

			// Encode as hex to ensure UTF8
			var hashBytes []byte = hash[:]
			resultHex := hex.EncodeToString(hashBytes)

			// Truncate to 10 bytes
			truncated := []byte(resultHex)[:10]
			return string(truncated)
		}

		truncated := HashTruncateName(workflowName)
		fmt.Println("Truncated name: ", truncated)
		fmt.Println("Workflow owner: ", sc.MustGetRootKeyAddress().Hex())
		fmt.Println("Workflow name: ", workflowName)
		fmt.Println("workflowNameBytes: ", string([]byte(truncated)))

		copy(workflowNameBytes[:], []byte(truncated))

		if in.WorkflowConfig.UseExising {
			require.NotEmpty(t, in.WorkflowConfig.Existing.BinaryURL)
			workFlowData, err := downloadAndDecode(in.WorkflowConfig.Existing.BinaryURL)
			require.NoError(t, err)

			var configData []byte
			if in.WorkflowConfig.Existing.ConfigURL != "" {
				configData, err = downloadAndDecode(in.WorkflowConfig.Existing.ConfigURL)
				require.NoError(t, err)
			}

			// use non-encoded workflow name
			workflowID, idErr := GenerateWorkflowIDFromStrings(sc.MustGetRootKeyAddress().Hex(), workflowName, workFlowData, configData, "")
			require.NoError(t, idErr)

			workflow_registryInstance, err := workflow_registry_wrapper.NewWorkflowRegistry(workflowRegistryAddr, sc.Client)
			require.NoError(t, err)

			// use non-encoded workflow name
			wrTx, wrErr := workflow_registryInstance.RegisterWorkflow(sc.NewTXOpts(), workflowName, [32]byte(common.Hex2Bytes(workflowID)), donID, uint8(0), in.WorkflowConfig.Existing.BinaryURL, in.WorkflowConfig.Existing.ConfigURL, "")
			_, decodeErr := sc.Decode(wrTx, wrErr)
			require.NoError(t, decodeErr)
		}

		if in.WorkflowConfig.UseChainlinkCLI {
			// create settings file
			settingsFile, err := os.CreateTemp("", ".chainlink-cli-settings.yaml")
			require.NoError(t, err)

			settings := ChainlinkCliSettings{
				DevPlatform: DevPlatform{
					CapabilitiesRegistryAddress: capRegAddr.Hex(),
					DonId:                       donID,
					WorkflowRegistryAddress:     workflowRegistryAddr.Hex(),
				},
				UserWorkflow: UserWorkflow{
					WorkflowOwnerAddress: sc.MustGetRootKeyAddress().Hex(),
				},
				Logging: Logging{},
				McmsConfig: McmsConfig{
					ProposalsDirectory: "./",
				},
				Contracts: Contracts{
					ContractRegistry: []ContractRegistry{
						{
							Name:          "CapabilitiesRegistry",
							Address:       capRegAddr.Hex(),
							ChainSelector: chainSelector,
						},
						{
							Name:          "WorkflowRegistry",
							Address:       workflowRegistryAddr.Hex(),
							ChainSelector: chainSelector,
						},
					},
				},
				Rpcs: []Rpc{
					{
						ChainSelector: chainSelector,
						URL:           bc.Nodes[0].HostHTTPUrl, // chainlink-cli doesn't work with WS
					},
				},
			}

			settingsMarshalled, err := yaml.Marshal(settings)
			require.NoError(t, err)

			_, err = settingsFile.Write(settingsMarshalled)
			require.NoError(t, err)

			var workflowGistUrl string
			var workflowConfigUrl string

			if in.WorkflowConfig.ChainlinkCLI.Compile {
				feedId := "0x018BFE88407000400000000000000000"

				configFile, err := os.CreateTemp("", "config.json")
				require.NoError(t, err)

				workflowConfig := PoRWorkflowConfig{
					FeedID:          feedId,
					URL:             "https://api.real-time-reserves.verinumus.io/v1/chainlink/proof-of-reserves/TrueUSD",
					ConsumerAddress: feedsConsumerAddress.Hex(),
				}

				configMarshalled, err := json.Marshal(workflowConfig)
				require.NoError(t, err)

				_, err = configFile.Write(configMarshalled)
				require.NoError(t, err)

				var outputBuffer bytes.Buffer

				compileCmd := exec.Command("chainlink-cli", "workflow", "compile", "-S", settingsFile.Name(), "-c", configFile.Name(), "main.go")
				compileCmd.Stdout = &outputBuffer
				compileCmd.Stderr = &outputBuffer
				compileCmd.Dir = *in.WorkflowConfig.ChainlinkCLI.FolderLocation
				err = compileCmd.Start()
				require.NoError(t, err)

				err = compileCmd.Wait()
				require.NoError(t, err)

				fmt.Println("Compile output:\n", outputBuffer.String())

				re := regexp.MustCompile(`Gist URL=([^\s]+)`)
				matches := re.FindAllStringSubmatch(outputBuffer.String(), -1)
				require.Len(t, matches, 2)

				ansiEscapePattern := `\x1b\[[0-9;]*m`
				re = regexp.MustCompile(ansiEscapePattern)

				workflowGistUrl = re.ReplaceAllString(matches[0][1], "")
				workflowConfigUrl = re.ReplaceAllString(matches[1][1], "")

				require.NotEmpty(t, workflowGistUrl)
				require.NotEmpty(t, workflowConfigUrl)
			} else {
				workflowGistUrl = in.WorkflowConfig.ChainlinkCLI.ExistingWorkflowConfig.BinaryURL
				workflowConfigUrl = in.WorkflowConfig.ChainlinkCLI.ExistingWorkflowConfig.ConfigURL
			}

			// deploy the workflow
			registerCmd := exec.Command("chainlink-cli", "workflow", "register", workflowName, "-b", workflowGistUrl, "-c", workflowConfigUrl, "-S", settingsFile.Name(), "-v")
			registerCmd.Stdout = os.Stdout
			registerCmd.Stderr = os.Stderr
			registerCmd.Dir = *in.WorkflowConfig.ChainlinkCLI.FolderLocation
			err = registerCmd.Run()
			require.NoError(t, err)
		}

		feedsConsumerInstance, err := feeds_consumer.NewKeystoneFeedsConsumer(feedsConsumerAddress, sc.Client)
		require.NoError(t, err)

		// here we need to use hex-encoded workflow name converted to []byte
		tx, err = feedsConsumerInstance.SetConfig(
			sc.NewTXOpts(),
			[]common.Address{forwarderAddress},           // allowed senders
			[]common.Address{sc.MustGetRootKeyAddress()}, // allowed workflow owners
			[][10]byte{workflowNameBytes},                // allowed workflow names
		)
		_, decodeErr = sc.Decode(tx, err)
		require.NoError(t, decodeErr)

		// TODO: When the capabilities registry address is provided:
		// - NOPs and nodes are added to the registry.
		// - Nodes are configured to listen to the registry for updates.
		nodeset, err := ns.NewSharedDBNodeSet(in.NodeSet, bc)
		require.NoError(t, err)

		nodeClients, err := clclient.New(nodeset.CLNodes)
		require.NoError(t, err)

		nodesInfo := getNodesInfo(t, nodeClients)

		for _, nodeInfo := range nodesInfo {
			_, err := actions.SendFunds(zerolog.Logger{}, sc, actions.FundsToSendPayload{
				ToAddress:  common.HexToAddress(nodeInfo.TransmitterAddress),
				Amount:     big.NewInt(5000000000000000000),
				PrivateKey: sc.MustGetRootPrivateKey(),
			})
			require.NoError(t, err)
		}

		bootstrapNodeInfo := nodesInfo[0]
		workflowNodesetInfo := nodesInfo[1:]

		// bootstrap node
		in.NodeSet.NodeSpecs[0].Node.TestConfigOverrides = fmt.Sprintf(`
				[Feature]
				LogPoller = true

				[OCR2]
				Enabled = true
				DatabaseTimeout = '1s'

				[P2P.V2]
				Enabled = true
				ListenAddresses = ['0.0.0.0:5001']
				DefaultBootstrappers = ['%s@localhost:5001']

				[Capabilities.Peering.V2]
				Enabled = true
				ListenAddresses = ['0.0.0.0:6690']
				DefaultBootstrappers = ['%s@localhost:6690']

				# This is needed for the target capability to be initialized
				[[EVM]]
				ChainID = '%s'

				[[EVM.Nodes]]
				Name = 'anvil'
				WSURL = '%s'
				HTTPURL = '%s'
			`,
			bootstrapNodeInfo.PeerID,
			bootstrapNodeInfo.PeerID,
			bc.ChainID,
			bc.Nodes[0].DockerInternalWSUrl,
			bc.Nodes[0].DockerInternalHTTPUrl,
		)

		for i := range workflowNodesetInfo {
			in.NodeSet.NodeSpecs[i+1].Node.TestConfigOverrides = fmt.Sprintf(`
				[Feature]
				LogPoller = true

				[OCR2]
				Enabled = true
				DatabaseTimeout = '1s'

				[P2P.V2]
				Enabled = true
				ListenAddresses = ['0.0.0.0:5001']
				# assuming that node0 is the bootstrap node
				DefaultBootstrappers = ['%s@node0:5001']

				[Capabilities.Peering.V2]
				Enabled = true
				ListenAddresses = ['0.0.0.0:6690']
				# assuming that node0 is the bootstrap node
				DefaultBootstrappers = ['%s@node0:6690']

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

				[Capabilities.ExternalRegistry]
				Address = '%s'
				NetworkID = 'evm'
				ChainID = '%s'

				[Capabilities.WorkflowRegistry]
				Address = "%s"
				NetworkID = "evm"
				ChainID = "%s"

				[Capabilities.GatewayConnector]
				DonID = "1"
				ChainIDForNodeKey = "%s"
				NodeAddress = '%s'

				[[Capabilities.GatewayConnector.Gateways]]
				Id = "por_gateway"
				URL = "%s"
			`,
				bootstrapNodeInfo.PeerID,
				bootstrapNodeInfo.PeerID,
				bc.ChainID,
				bc.Nodes[0].DockerInternalWSUrl,
				bc.Nodes[0].DockerInternalHTTPUrl,
				workflowNodesetInfo[i].TransmitterAddress,
				forwarderAddress.Hex(),
				capRegAddr,
				bc.ChainID,
				workflowRegistryAddr.Hex(),
				bc.ChainID,
				bc.ChainID,
				workflowNodesetInfo[i].TransmitterAddress,
				// assuming that node0 is the bootstrap node
				"ws://node0:5003/node",
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
		_, decodeErr = sc.Decode(tx, err)
		require.NoError(t, decodeErr)

		fmt.Println("Deployed ocr3_capability contract at", ocr3CapabilityAddress.Hex())

		// Add bootstrap spec to the last node
		bootstrapNode := nodeClients[0]

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			gatewayJobSpec := fmt.Sprintf(`
				type = "gateway"
				schemaVersion = 1
				name = "PoR Gateway"
				forwardingAllowed = false

				[gatewayConfig.ConnectionManagerConfig]
				AuthChallengeLen = 10
				AuthGatewayId = "por_gateway"
				AuthTimestampToleranceSec = 5
				HeartbeatIntervalSec = 20

				[[gatewayConfig.Dons]]
				DonId = "1"
				F = 1
				HandlerName = "web-api-capabilities"
					[gatewayConfig.Dons.HandlerConfig]
					MaxAllowedMessageAgeSec = 1_000

						[gatewayConfig.Dons.HandlerConfig.NodeRateLimiter]
						GlobalBurst = 10
						GlobalRPS = 50
						PerSenderBurst = 10
						PerSenderRPS = 10

					[[gatewayConfig.Dons.Members]]
					Address = "%s"
					Name = "Workflow Node 1"
					[[gatewayConfig.Dons.Members]]
					Address = "%s"
					Name = "Workflow Node 2"
					[[gatewayConfig.Dons.Members]]
					Address = "%s"
					Name = "Workflow Node 3"
					[[gatewayConfig.Dons.Members]]
					Address = "%s"
					Name = "Workflow Node 4"

				[gatewayConfig.NodeServerConfig]
				HandshakeTimeoutMillis = 1_000
				MaxRequestBytes = 100_000
				Path = "/node"
				Port = 5_003 #this is the port the other nodes will use to connect to the gateway
				ReadTimeoutMillis = 1_000
				RequestTimeoutMillis = 10_000
				WriteTimeoutMillis = 1_000

				[gatewayConfig.UserServerConfig]
				ContentTypeHeader = "application/jsonrpc"
				MaxRequestBytes = 100_000
				Path = "/"
				Port = 5_002
				ReadTimeoutMillis = 1_000
				RequestTimeoutMillis = 10_000
				WriteTimeoutMillis = 1_000

				[gatewayConfig.HTTPClientConfig]
				MaxResponseBytes = 100_000_000
			`,
				// ETH keys of the workflow nodes
				workflowNodesetInfo[0].TransmitterAddress,
				workflowNodesetInfo[1].TransmitterAddress,
				workflowNodesetInfo[2].TransmitterAddress,
				workflowNodesetInfo[3].TransmitterAddress,
			)

			r, _, err2 := bootstrapNode.CreateJobRaw(gatewayJobSpec)
			assert.NoError(t, err2)
			assert.Empty(t, r.Errors)

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
			r, _, err3 := bootstrapNode.CreateJobRaw(bootstrapJobSpec)
			assert.NoError(t, err3)
			assert.Empty(t, r.Errors)
		}()

		for i, nodeClient := range nodeClients {
			// Last node is a bootstrap node, so we skip it
			if i == 0 {
				continue
			}

			wg.Add(1)
			go func() {
				defer wg.Done()
				cronJobSpec := `
					type = "standardcapabilities"
					schemaVersion = 1
					name = "cron-capabilities"
					forwardingAllowed = false
					command = "/home/capabilities/cron-linux-amd64"
					config = ""
				`

				response, _, errCron := nodeClient.CreateJobRaw(cronJobSpec)
				assert.NoError(t, errCron)
				assert.Empty(t, response.Errors)

				computeJobSpec := `
					type = "standardcapabilities"
					schemaVersion = 1
					name = "compute-capabilities"
					forwardingAllowed = false
					command = "__builtin_custom-compute-action"
					config = """
					NumWorkers = 3
						[rateLimiter]
						globalRPS = 20.0
						globalBurst = 30
						perSenderRPS = 1.0
						perSenderBurst = 5
					"""
				`

				response, _, errCompute := nodeClient.CreateJobRaw(computeJobSpec)
				assert.NoError(t, errCompute)
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
					"node0:5001",
					nodesInfo[i].TransmitterAddress,
					bc.ChainID,
					nodesInfo[i].OcrKeyBundleID,
				)
				fmt.Println("consensusJobSpec", consensusJobSpec)
				response, _, errCons := nodeClient.CreateJobRaw(consensusJobSpec)
				assert.NoError(t, errCons)
				assert.Empty(t, response.Errors)
			}()
		}
		wg.Wait()

		// req := keystone_changeset.InitialContractsCfg{
		// 	RegistryChainSel: chainSelector,
		// }

		// _, err = keystone_changeset.ConfigureInitialContractsChangeset(*ctfEnv, req)
		// require.NoError(t, err)

		// allCapabilities := []kcr.CapabilitiesRegistryCapability{
		// 	{
		// 		LabelledName:   "offchain_reporting",
		// 		Version:        "1.0.0",
		// 		CapabilityType: 2, // CONSENSUS
		// 		ResponseType:   0, // REPORT
		// 	},
		// 	{
		// 		LabelledName:   "write_geth-testnet",
		// 		Version:        "1.0.0",
		// 		CapabilityType: 3, // TARGET
		// 		ResponseType:   1, // OBSERVATION_IDENTICAL
		// 	},
		// 	{
		// 		LabelledName:   "cron-trigger",
		// 		Version:        "1.0.0",
		// 		CapabilityType: uint8(0), // trigger
		// 	},
		// 	{
		// 		LabelledName:   "custom-compute",
		// 		Version:        "1.0.0",
		// 		CapabilityType: uint8(1), // action
		// 	},
		// }

		// p2pCapabilities := map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability{}

		// for i, node := range nodesInfo {
		// 	if i == 0 {
		// 		continue
		// 	}

		// 	peerId, err := p2pkey.MakePeerID(node.PeerID)
		// 	require.NoError(t, err)

		// 	p2pCapabilities[peerId] = allCapabilities
		// }

		// _, err = keystone_changeset.AppendNodeCapabilities(*ctfEnv, &keystone_changeset.AppendNodeCapabilitiesRequest{
		// 	RegistryChainSel:  chainSelector,
		// 	P2pToCapabilities: p2pCapabilities,
		// })
		// require.NoError(t, err)

		// keystone_changeset.AddCapabilities(*ctfEnv, capabilitiesRegistryInstance, []cr_wrapper.CapabilitiesRegistryCapability{
		// 	{
		// 		LabelledName:   "offchain_reporting",
		// 		Version:        "1.0.0",
		// 		CapabilityType: 2, // CONSENSUS
		// 		ResponseType:   0, // REPORT
		// 	},
		// 	{
		// 		LabelledName:   "write_geth-testnet",
		// 		Version:        "1.0.0",
		// 		CapabilityType: 3, // TARGET
		// 		ResponseType:   1, // OBSERVATION_IDENTICAL
		// 	},
		// 	{
		// 		LabelledName:   "cron-trigger",
		// 		Version:        "1.0.0",
		// 		CapabilityType: uint8(0), // trigger
		// 	},
		// 	{
		// 		LabelledName:   "custom-compute",
		// 		Version:        "1.0.0",
		// 		CapabilityType: uint8(1), // action
		// 	},
		// })

		// nops := []*kcr.CapabilitiesRegistryNodeOperator{
		// 	{
		// 		Admin: common.HexToAddress(sc.MustGetRootKeyAddress().Hex()),
		// 		Name:  "Admin",
		// 	},
		// }
		// nopToNodeId := map[kcr.CapabilitiesRegistryNodeOperator][]string{
		// 	{
		// 		Admin: common.HexToAddress(sc.MustGetRootKeyAddress().Hex()),
		// 		Name:  "Admin",
		// 	}: {nodesInfo[0].PeerID, nodesInfo[1].PeerID, nodesInfo[2].PeerID, nodesInfo[3].PeerID},
		// }

		// donToNodes := make(map[string][]deployment.Node)
		// donToNodes["1"] = []deployment.Node{}

		// donCaps := []keystone_changeset.RegisteredCapability{}

		// for _, cap := range allCapabilities {
		// 	c, err := keystone_changeset.FromCapabilitiesRegistryCapability(&cap, *ctfEnv, chainSelector)
		// 	require.NoError(t, err)

		// 	donCaps = append(donCaps, keystone_changeset.RegisteredCapability(*c))
		// }

		// donToCapability := map[string][]keystone_changeset.RegisteredCapability{
		// 	"1": donCaps,
		// }

		// _, err = keystone_changeset.RegisterNodes(lgr, &keystone_changeset.RegisterNodesRequest{
		// 	Env:                   ctfEnv,
		// 	RegistryChainSelector: chainSelector,
		// 	NopToNodeIDs:          nopToNodeId,
		// 	DonToNodes:            donToNodes,
		// 	DonToCapabilities:     donToCapability,
		// 	Nops:                  nops,
		// })
		// require.NoError(t, err)

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
				HashedCapabilityIds: hashedCapabilities,
			})

			donNodes = append(donNodes, peerID)
			signers = append(signers, node.Signer)
		}

		// Add NOPs to registry
		tx, err = capabilitiesRegistryInstance.AddNodeOperators(
			sc.NewTXOpts(),
			nopsToAdd,
		)
		_, decodeErr = sc.Decode(tx, err)
		require.NoError(t, decodeErr)

		// Add nodes to registry
		tx, err = capabilitiesRegistryInstance.AddNodes(
			sc.NewTXOpts(),
			nodesToAdd,
		)
		_, decodeErr = sc.Decode(tx, err)
		require.NoError(t, decodeErr)

		var capRegConfig []cr_wrapper.CapabilitiesRegistryCapabilityConfiguration
		for _, hashed := range hashedCapabilities {
			capRegConfig = append(capRegConfig, cr_wrapper.CapabilitiesRegistryCapabilityConfiguration{
				CapabilityId: hashed,
				Config:       []byte(""),
			})
		}

		// Add nodeset to registry
		tx, err = capabilitiesRegistryInstance.AddDON(
			sc.NewTXOpts(),
			donNodes,
			capRegConfig,
			true,     // is public
			true,     // accepts workflows
			uint8(1), // max number of malicious nodes
		)
		_, decodeErr = sc.Decode(tx, err)
		require.NoError(t, decodeErr)

		forwarderInstance, err := forwarder.NewKeystoneForwarder(forwarderAddress, sc.Client)
		require.NoError(t, err)

		_, err = sc.Decode(forwarderInstance.SetConfig(
			sc.NewTXOpts(),
			1,
			1,
			1,
			signers))
		require.NoError(t, err)

		// Wait for OCR listeners to be ready before setting the configuration.
		// If the ConfigSet event is missed, OCR protocol will not start.
		// TODO make it fluent!
		fmt.Println("Waiting 30s for OCR listeners to be ready...")
		time.Sleep(30 * time.Second)
		fmt.Println("Proceeding to set OCR3 configuration.")

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
		_, decodeErr = sc.Decode(tx, err)
		require.NoError(t, decodeErr)

		// It can take a while before the first report is produced, particularly on CI.
		timeout := 10 * time.Minute
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Node sent transaction

		startTime := time.Now()
		for {
			select {
			case <-ctx.Done():
				t.Fatalf("feed did not update, timeout after %s", timeout)
			case <-time.After(10 * time.Second):
				elapsed := time.Since(startTime).Round(time.Second)
				price, _, err := feedsConsumerInstance.GetPrice(
					sc.NewCallOpts(),
					feedBytes,
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
