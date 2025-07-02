package integrationtest

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	datastreamsllo "github.com/smartcontractkit/chainlink-data-streams/llo"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/data_feeds_cache"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/configurator"
	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	evmtestutils "github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr3/securemint"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/llo"
	"github.com/smartcontractkit/freeport"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/por_mock_ocr3plugin/por"
)

var (
	fNodes = uint8(1)
	nNodes = 4 // number of nodes (not including bootstrap)
)

// backendSimulator is an adapter that implements the chainlink-evm/pkg/types/types/Backend interface
// so that we can have minimal changes in the integration test which depends on an in-memory blockchain, while now we want to use anvil running in Docker (started by CRE dev env).
type backendSimulator struct {
	t           *testing.T
	chainClient clientSimulator
}

type clientSimulator struct {
	*ethclient.Client
}

var _ evmtypes.Backend = &backendSimulator{}

var _ simulated.Client = &clientSimulator{}

func (b *backendSimulator) Commit() common.Hash {
	b.t.Logf("Commit not implemented")
	return common.Hash{}
}

func (b *backendSimulator) Close() error {
	b.chainClient.Close()
	return nil
}

func (b *backendSimulator) Rollback() {
	b.t.Errorf("Rollback not implemented")
}

func (b *backendSimulator) Fork(parentHash common.Hash) error {
	b.t.Errorf("Fork not implemented")
	return nil
}

func (b *backendSimulator) AdjustTime(adjustment time.Duration) error {
	b.t.Errorf("AdjustTime not implemented")
	return nil
}

func (b *backendSimulator) Client() simulated.Client {
	return b.chainClient
}

func simulateBackend(t *testing.T, chainClient *ethclient.Client) evmtypes.Backend {
	return &backendSimulator{t: t, chainClient: clientSimulator{Client: chainClient}}
}

// TestIntegration_SecureMint_happy_path tests runs a small DON which runs the secure mint plugin
// and verifies that it can successfully create reports.
//
// Inspired by:
// * core/internal/features/ocr2/features_ocr2_test.go
// * core/services/ocr2/plugins/ocr2keeper/integration_21_test.go
func TestIntegration_SecureMint_happy_path(t *testing.T) {
	const salt = 100

	clientCSAKeys := make([]csakey.KeyV2, nNodes)
	clientPubKeys := make([]ed25519.PublicKey, nNodes)
	for i := range nNodes {
		k := big.NewInt(int64(salt + i))
		key := csakey.MustNewV2XXXTestingOnly(k)
		clientCSAKeys[i] = key
		clientPubKeys[i] = key.PublicKey
	}

	// steve, backend := setupBlockchain(t)
	steve, chainClient := connectToAnvil(t)
	backend := simulateBackend(t, chainClient)
	fromBlock, err := chainClient.BlockNumber(testutils.Context(t))
	require.NoError(t, err)
	t.Logf("Starting from block: %d", fromBlock)

	// Setup bootstrap
	// bootstrapCSAKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(salt - 1))
	// bootstrapNodePort := freeport.GetOne(t)
	bootstrapNodePort := 5002
	// appBootstrap, bootstrapPeerID, _, bootstrapKb, _ := setupNode(t, bootstrapNodePort, "bootstrap_securemint", backend, bootstrapCSAKey, nil)
	// bootstrapNode := node{app: appBootstrap, keyBundle: bootstrapKb}
	bootstrapPeerID := "12D3KooWQRiahhF1CrTex7P5gMuw4eg2sZ3Qn2XZehT6EuwPG6uL"

	p2pV2Bootstrappers := []commontypes.BootstrapperLocator{
		// Supply the bootstrap IP and port as a V2 peer address
		{PeerID: bootstrapPeerID, Addrs: []string{fmt.Sprintf("127.0.0.1:%d", bootstrapNodePort)}},
	}

	// Setup oracle nodes
	oracles, nodes := setupNodes(t, nNodes, backend, clientCSAKeys, func(c *chainlink.Config) {
		// inform node about bootstrap node
		c.P2P.V2.DefaultBootstrappers = &p2pV2Bootstrappers
	})

	// Setup capabilities registry and DON
	regAddress := common.HexToAddress("0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512")
	_, capRegAddress := setupSecureMintCapabilitiesRegistry(t, regAddress, steve, backend, nodes)
	t.Logf("Capabilities registry setup complete at: %s", capRegAddress.Hex())

	// Configure nodes to use the capabilities registry
	for i := range nodes {
		capRegAddressStr := capRegAddress.Hex()
		// Note: The capabilities registry configuration is already set in setupNode
		// The nodes are already configured to use the capabilities registry
		t.Logf("Node %d is configured to use capabilities registry at %s", i, capRegAddressStr)
	}

	allowedSenders := make([]common.Address, len(nodes))
	for i, node := range nodes {
		keys, err := node.app.GetKeyStore().Eth().EnabledKeysForChain(testutils.Context(t), testutils.SimulatedChainID)
		require.NoError(t, err)
		allowedSenders[i] = keys[0].Address // assuming the first key is the transmitter
	}

	_, configuratorAddress := setSecureMintOnchainConfigUsingOCR3Configurator(t, steve, backend, nodes, oracles)

	t.Logf("Creating bootstrap job with configurator address: %s", configuratorAddress.Hex())
	// bootstrapJob := createSecureMintBootstrapJob(t, bootstrapNode, configuratorAddress, testutils.SimulatedChainID.String(), fmt.Sprintf("%d", fromBlock))
	// t.Logf("Created bootstrap job: %s with id %d", bootstrapJob.Name.ValueOrZero(), bootstrapJob.ID)

	jobIDs := addSecureMintOCRJobs(t, nodes, configuratorAddress)

	t.Logf("jobIDs: %v", jobIDs)
	validateJobsRunningSuccessfully(t, nodes, jobIDs)

	t.Logf("Waiting for CRE Workflow to register itself as a subscriber to the secure mint trigger (in securemint/transmitter.go) and get triggered")
	time.Sleep(9 * time.Minute)
}

func setupBlockchain(t *testing.T) (
	*bind.TransactOpts,
	evmtypes.Backend,
) {
	steve := evmtestutils.MustNewSimTransactor(t) // config contract deployer and owner
	genesisData := gethtypes.GenesisAlloc{steve.From: {Balance: assets.Ether(1000).ToInt()}}
	backend := cltest.NewSimulatedBackend(t, genesisData, ethconfig.Defaults.Miner.GasCeil)
	backend.Commit()
	backend.Commit() // ensure starting block number at least 1

	return steve, backend
}

// connectToAnvil connects to anvil started by CRE
func connectToAnvil(t *testing.T) (*bind.TransactOpts, *ethclient.Client) {
	ctfKey, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	require.NoError(t, err)
	steve, err := bind.NewKeyedTransactorWithChainID(ctfKey, big.NewInt(1337))
	require.NoError(t, err)
	ethClient, err := ethclient.Dial("http://localhost:8545")
	require.NoError(t, err)
	return steve, ethClient
}

func setupNodes(t *testing.T, nNodes int, backend evmtypes.Backend, clientCSAKeys []csakey.KeyV2, f func(*chainlink.Config)) (oracles []confighelper.OracleIdentityExtra, nodes []node) {
	ports := freeport.GetN(t, nNodes)
	for i := range nNodes {
		app, peerID, transmitter, kb, observedLogs := setupNode(t, ports[i], fmt.Sprintf("oracle_securemint_%d", i), backend, clientCSAKeys[i], f)

		nodes = append(nodes, node{
			app:          app,
			clientPubKey: transmitter,
			keyBundle:    kb,
			observedLogs: observedLogs,
		})
		offchainPublicKey, err := hex.DecodeString(strings.TrimPrefix(kb.OnChainPublicKey(), "0x"))
		require.NoError(t, err)
		oracles = append(oracles, confighelper.OracleIdentityExtra{
			OracleIdentity: confighelper.OracleIdentity{
				OnchainPublicKey:  offchainPublicKey,
				TransmitAccount:   ocr2types.Account(hex.EncodeToString(transmitter[:])),
				OffchainPublicKey: kb.OffchainPublicKey(),
				PeerID:            peerID,
			},
			ConfigEncryptionPublicKey: kb.ConfigEncryptionPublicKey(),
		})
	}
	return
}

func validateJobsRunningSuccessfully(t *testing.T, nodes []node, jobIDs map[int]int32) {

	// 0. Add ourselves as a subscriber to the secure mint trigger capability
	transmissions := atomic.NewInt32(0)
	transmitter := securemint.SingletonTransmitter.Load().(capabilities.TriggerCapability)
	triggerConfig, err := values.NewMap(map[string]any{
		"workflowID":     "securemint-workflow",
		"maxFrequencyMs": 1000,
	})
	require.NoError(t, err)
	registerCh, err := transmitter.RegisterTrigger(testutils.Context(t), capabilities.TriggerRegistrationRequest{
		TriggerID: "securemint-trigger",
		Metadata: capabilities.RequestMetadata{
			WorkflowID: "securemint-workflow",
		},
		Config: triggerConfig,
	})
	require.NoError(t, err)
	go func() {
		for resp := range registerCh {
			t.Logf("Received trigger response: %+v", resp)
			outputs, err := resp.Event.Outputs.Unwrap()
			require.NoError(t, err)
			t.Logf("Received trigger response outputs: %+v", outputs)
			transmissions.Inc()
		}
	}()

	// 1. Assert no job spec errors
	for i, node := range nodes {
		jobs, _, err := node.app.JobORM().FindJobs(testutils.Context(t), 0, 1000)
		require.NoErrorf(t, err, "assert error finding jobs for node %d", i)
		t.Logf("%d jobs found for node %d", len(jobs), i)
		for _, j := range jobs {
			t.Logf("job %d on node %d oracle spec: %#v", j.ID, i, j.OCR2OracleSpec)
			t.Logf("job %d on node %d pipeline spec: %#v", j.ID, i, j.PipelineSpec)
		}
		// No spec errors
		for _, j := range jobs {
			ignore := 0
			for _, jse := range j.JobSpecErrors {
				// Non-fatal timing related error, ignore for testing.
				if strings.Contains(jse.Description, "leader's phase conflicts tGrace timeout") {
					ignore++
				} else {
					t.Errorf("assert error: job spec error on node %d: %v", i, jse)
				}
			}
			require.Lenf(t, j.JobSpecErrors, ignore, "assert error: job spec errors on node %d", i)
		}
	}

	t.Logf("No job spec errors identified for any node")

	runs, err := nodes[0].app.PipelineORM().GetAllRuns(testutils.Context(t))
	require.NoError(t, err, "assert error getting all runs")
	t.Logf("Found %d runs", len(runs))
	for _, run := range runs {
		t.Logf("Run ID: %d, Job ID: %d, Status: %s", run.ID, run.JobID, run.Status())
	}

	// 2. Assert that all the Secure Mint jobs get a run with valid values eventually
	var wg sync.WaitGroup
	for i, node := range nodes {
		wg.Add(1)
		go func() {
			defer wg.Done()

			pr := cltest.WaitForPipelineComplete(t, i, jobIDs[i], 1, 0, node.app.JobORM(), 30*time.Second, 1*time.Second)
			outputs, err := pr[0].Outputs.MarshalJSON()
			if !assert.NoError(t, err) {
				t.Logf("assert error marshalling outputs for job %d: %v", jobIDs[i], err)
				return
			}
			t.Logf("Pipeline itself is %+v", pr[0])
			t.Logf("Pipeline run outputs are %s", string(outputs))
		}()
	}
	t.Logf("waiting for pipeline runs to complete")
	wg.Wait()
	t.Logf("All pipeline runs completed successfully")

	// 3. Check that transmissions work
	expectedNumTransmissions := int32(4)
	gomega.NewWithT(t).Eventually(func() bool {
		numTransmissions := transmissions.Load()
		t.Logf("Number of (stub) report transmissions: %d", numTransmissions)
		return numTransmissions >= expectedNumTransmissions
	}, 30*time.Second, 1*time.Second).Should(
		gomega.BeTrue(),
		fmt.Sprintf("expected at least %d reports transmitted, but got less", expectedNumTransmissions),
	)
}

func setSecureMintOnchainConfigUsingOCR3Configurator(t *testing.T, steve *bind.TransactOpts, backend evmtypes.Backend, nodes []node, oracles []confighelper.OracleIdentityExtra) (*configurator.Configurator, common.Address) {

	// 1. Deploy configurator contract
	configuratorAddress, _, configurator, err := configurator.DeployConfigurator(steve, backend.Client())
	require.NoError(t, err)
	backend.Commit()

	// Ensure we have finality depth worth of blocks to start.
	for range 5 {
		backend.Commit()
	}
	t.Logf("Deployed OCR3Configurator contract at: %s", configuratorAddress.Hex())

	// 2. Get the oracle config
	smPluginConfig := por.PorOffchainConfig{MaxChains: 5}
	smPluginConfigBytes, err := smPluginConfig.Serialize()
	require.NoError(t, err)

	// using the data streams llo codec for the validation about version and predecessor config digest in the Configurator contract: https://github.com/smartcontractkit/chainlink/blob/develop/contracts/src/v0.8/llo-feeds/v0.5.0/configuration/Configurator.sol#L116-L124
	onchainConfig, err := (&datastreamsllo.EVMOnchainConfigCodec{}).Encode(datastreamsllo.OnchainConfig{
		Version:                 1,
		PredecessorConfigDigest: nil,
	})
	require.NoError(t, err)

	signers, _, f, outOnchainConfig, offchainConfigVersion, offchainConfig, err := ocr3confighelper.ContractSetConfigArgsForTests(
		2*time.Second,        // deltaProgress,
		20*time.Second,       // deltaResend,
		400*time.Millisecond, // deltaInitial,
		500*time.Millisecond, // deltaRound,
		250*time.Millisecond, // deltaGrace,
		300*time.Millisecond, // deltaCertifiedCommitRequest,
		1*time.Minute,        // deltaStage,
		100,                  // rMax,
		[]int{len(oracles)},  // s,
		oracles,              // oracles,
		smPluginConfigBytes,  // reportingPluginConfig,
		nil,                  // maxDurationInitialization,
		250*time.Millisecond, // maxDurationQuery,
		1*time.Second,        // maxDurationObservation,
		1*time.Second,        // maxDurationShouldAcceptAttestedReport,
		1*time.Second,        // maxDurationShouldTransmitAcceptedReport,
		int(fNodes),          // f,
		onchainConfig,        // onchainConfig (binary blob containing configuration passed through to the ReportingPlugin and also available to the contract. Unlike ReportingPluginConfig which is only available offchain.)
	)
	require.NoError(t, err)

	// 3. Set config on the contract
	signerKeys := make([][]byte, len(signers))
	for i, signer := range signers {
		signerKeys[i] = signer
	}

	// use csa keys as transmitters, similar to LLO
	transmitters := make([][32]byte, nNodes)
	for i := range nNodes {
		transmitters[i] = nodes[i].clientPubKey
	}
	t.Logf("transmitters: %v", transmitters)

	configID := [32]byte{}
	copy(configID[:], common.FromHex("0x0000000000000000000000000000000000000000000000000000000000000001"))

	_, err = configurator.SetProductionConfig(steve, configID, signerKeys, transmitters, f, outOnchainConfig, offchainConfigVersion, offchainConfig)
	if err != nil {
		t.Logf("Error: %s", err)
		errString, err := rPCErrorFromError(err)
		require.NoError(t, err)
		t.Fatalf("Failed to configure contract: %s %s", errString, err)
	}

	// make sure config is finalized
	for range 5 {
		backend.Commit()
	}

	var topic common.Hash
	topic = llo.ProductionConfigSet

	var logs []gethtypes.Log
	gomega.NewWithT(t).Eventually(func() bool {
		logs, err = backend.Client().FilterLogs(testutils.Context(t), ethereum.FilterQuery{Addresses: []common.Address{configuratorAddress}, Topics: [][]common.Hash{[]common.Hash{topic, configID}}})
		return err == nil && len(logs) > 0
	}, 30*time.Second, 1*time.Second).Should(
		gomega.BeTrue(),
		fmt.Sprintf("expected at least 1 log, but got none, got error: %v", err),
	)

	// logs, err := backend.Client().FilterLogs(testutils.Context(t), ethereum.FilterQuery{Addresses: []common.Address{configuratorAddress}, Topics: [][]common.Hash{[]common.Hash{topic, configID}}})
	// require.NoError(t, err)
	// require.GreaterOrEqual(t, len(logs), 1)
	cfg, err := llo.DecodeProductionConfigSetLog(logs[len(logs)-1].Data)
	require.NoError(t, err)

	t.Logf("Configurator config digest: 0x%x", cfg.ConfigDigest)

	return configurator, configuratorAddress
}

func rPCErrorFromError(txError error) (string, error) {
	errBytes, err := json.Marshal(txError)
	if err != nil {
		return "", err
	}
	var callErr struct {
		Code    int
		Data    string `json:"data"`
		Message string `json:"message"`
	}
	err = json.Unmarshal(errBytes, &callErr)
	if err != nil {
		return "", err
	}
	// If the error data is blank
	if len(callErr.Data) == 0 {
		return callErr.Data, nil
	}
	// Some nodes prepend "Reverted " and we also remove the 0x
	trimmed := strings.TrimPrefix(callErr.Data, "Reverted ")[2:]
	data, err := hex.DecodeString(trimmed)
	if err != nil {
		return "", err
	}
	revert, err := abi.UnpackRevert(data)
	// If we can't decode the revert reason, return the raw data
	if err != nil {
		return callErr.Data, nil
	}
	return revert, nil
}

// Not used yet, in scope for chain writing
func setupDataFeedsCacheContract(t *testing.T, steve *bind.TransactOpts, backend evmtypes.Backend, allowedSenders []common.Address, workflowOwner, workflowName string) (
	common.Address, *data_feeds_cache.DataFeedsCache) {

	addr, _, dataFeedsCache, err := data_feeds_cache.DeployDataFeedsCache(steve, backend.Client())
	require.NoError(t, err)
	backend.Commit()

	var nameBytes [10]byte
	copy(nameBytes[:], workflowName)

	ownerAddr := common.HexToAddress(workflowOwner)

	_, err = dataFeedsCache.SetFeedAdmin(steve, ownerAddr, true)
	require.NoError(t, err)

	backend.Commit()

	metadatas := make([]data_feeds_cache.DataFeedsCacheWorkflowMetadata, len(allowedSenders))
	for i, sender := range allowedSenders {
		metadatas[i] =
			data_feeds_cache.DataFeedsCacheWorkflowMetadata{
				AllowedSender:        sender,
				AllowedWorkflowOwner: ownerAddr,
				AllowedWorkflowName:  nameBytes,
			}
	}

	feedIDBytes := [16]byte{}
	copy(feedIDBytes[:], common.FromHex("0xA1B2C3D4E5F600010203040506070809"))

	_, err = dataFeedsCache.SetDecimalFeedConfigs(steve, [][16]byte{feedIDBytes}, []string{"securemint"}, metadatas)
	if err != nil {
		errString, err := rPCErrorFromError(err)
		require.NoError(t, err)

		t.Fatalf("Failed to configure contract: %s", errString)
	}

	backend.Commit()

	return addr, dataFeedsCache
}

// setupSecureMintCapabilitiesRegistry connects to an existing capabilities registry at the given address
// and registers the secure mint capability, node operator, nodes, and DON.
func setupSecureMintCapabilitiesRegistry(t *testing.T, regAddress common.Address, steve *bind.TransactOpts, backend evmtypes.Backend, nodes []node) (*kcr.CapabilitiesRegistry, common.Address) {
	capReg, err := kcr.NewCapabilitiesRegistry(regAddress, backend.Client())
	require.NoError(t, err)
	t.Logf("Connected to CapabilitiesRegistry at: %s", regAddress.Hex())

	// Add secure mint capability (if not already present)
	secureMintCapability := kcr.CapabilitiesRegistryCapability{
		LabelledName:   "securemint-trigger",
		Version:        "1.0.0",
		CapabilityType: 0, // TRIGGER
		ResponseType:   0, // REPORT
	}
	_, err = capReg.AddCapabilities(steve, []kcr.CapabilitiesRegistryCapability{secureMintCapability})
	// Ignore error if already exists (check for both string message and custom error code)
	if err != nil {
		errString := err.Error()
		if strings.Contains(errString, "already exists") || strings.Contains(errString, "0xebf52551") {
			t.Logf("Secure mint capability already exists, skipping...")
		} else {
			t.Logf("Error adding secure mint capability: %v", err)
			errString, err := rPCErrorFromError(err)
			require.NoError(t, err)
			t.Fatalf("Failed to add secure mint capability: %s", errString)
		}
	}
	backend.Commit()

	// Get the hashed capability ID
	hashedCapabilityID, err := capReg.GetHashedCapabilityId(nil, secureMintCapability.LabelledName, secureMintCapability.Version)
	require.NoError(t, err)
	t.Logf("Secure mint capability ID: %x", hashedCapabilityID)

	// Add node operator (if not already present)
	_, err = capReg.AddNodeOperators(steve, []kcr.CapabilitiesRegistryNodeOperator{{
		Admin: steve.From,
		Name:  "securemint-nop",
	}})
	if err != nil {
		errString := err.Error()
		if strings.Contains(errString, "already exists") || strings.Contains(errString, "0x") {
			t.Logf("Node operator already exists, skipping...")
		} else {
			require.NoError(t, err)
		}
	}
	backend.Commit()

	time.Sleep(1 * time.Second) // wait for one block to be committed

	// Get the node operator ID from the event
	it, err := capReg.FilterNodeOperatorAdded(nil, nil, nil)
	require.NoError(t, err)
	var nodeOperatorID uint32
	for it.Next() {
		if it.Event.Name == "securemint-nop" {
			nodeOperatorID = it.Event.NodeOperatorId
			break
		}
	}
	require.NotZero(t, nodeOperatorID)
	t.Logf("Node operator ID: %d", nodeOperatorID)

	// Add nodes to the registry (if not already present)
	var nodeParams []kcr.CapabilitiesRegistryNodeParams
	var peerIDs [][32]byte
	for i, node := range nodes {
		p2pKeys, err := node.app.GetKeyStore().P2P().GetAll()
		require.NoError(t, err)
		require.Len(t, p2pKeys, 1, "Expected exactly one P2P key per node")
		peerID := p2pKeys[0].PeerID()
		nodeParam := kcr.CapabilitiesRegistryNodeParams{
			NodeOperatorId:      nodeOperatorID,
			Signer:              testutils.Random32Byte(),
			P2pId:               peerID,
			EncryptionPublicKey: testutils.Random32Byte(),
			HashedCapabilityIds: [][32]byte{hashedCapabilityID},
		}
		nodeParams = append(nodeParams, nodeParam)
		peerIDs = append(peerIDs, peerID)
		t.Logf("Added node %d with peer ID: %x", i, peerID)
	}
	_, err = capReg.AddNodes(steve, nodeParams)
	if err != nil {
		errString := err.Error()
		if strings.Contains(errString, "already exists") || strings.Contains(errString, "0x") {
			t.Logf("Nodes already exist, skipping...")
		} else {
			require.NoError(t, err)
		}
	}
	time.Sleep(1 * time.Second) // wait for one block to be committed
	backend.Commit()

	onChainNodes, err := capReg.GetNodes(nil)
	require.NoError(t, err)
	// t.Logf("Nodes: %+v", onChainNodes)
	for _, node := range onChainNodes {
		for _, capabilityId := range node.HashedCapabilityIds {
			t.Logf("Node with supported capability ids %x: %x", node.P2pId, capabilityId)
		}
	}

	// Create capability configuration
	capabilityConfig := &capabilitiespb.CapabilityConfig{
		DefaultConfig: values.Proto(values.EmptyMap()).GetMapValue(),
		RemoteConfig: &capabilitiespb.CapabilityConfig_RemoteTriggerConfig{
			RemoteTriggerConfig: &capabilitiespb.RemoteTriggerConfig{
				RegistrationRefresh:     nil, // Will use default
				RegistrationExpiry:      nil, // Will use default
				MinResponsesToAggregate: 0,   // F + 1
				MessageExpiry:           nil, // Will use default
			},
		},
	}
	configBytes, err := proto.Marshal(capabilityConfig)
	require.NoError(t, err)

	// Add DON as a capability DON (isPublic=true, acceptsWorkflows=false)
	_, err = capReg.AddDON(steve, peerIDs, []kcr.CapabilitiesRegistryCapabilityConfiguration{{
		CapabilityId: hashedCapabilityID,
		Config:       configBytes,
	}}, true, false, fNodes)
	if err != nil {
		errString := err.Error()
		if strings.Contains(errString, "already exists") || strings.Contains(errString, "0x") {
			t.Logf("DON already exists, skipping...")
		} else {
			require.NoError(t, err)
		}
	}
	time.Sleep(1 * time.Second) // wait for one block to be committed
	backend.Commit()

	onChainDONs, err := capReg.GetDONs(nil)
	require.NoError(t, err)
	// t.Logf("DONs: %v", onChainDONs)
	for _, don := range onChainDONs {
		strP2pIds := make([]string, len(don.NodeP2PIds))
		for i, nodeP2pId := range don.NodeP2PIds {
			strP2pIds[i] = hex.EncodeToString(nodeP2pId[:])
		}
		t.Logf("DON %d has %d nodes with p2p IDs %v", don.Id, len(don.NodeP2PIds), strP2pIds)
		for _, config := range don.CapabilityConfigurations {
			t.Logf("DON %d has capability config %x with config %v", don.Id, config.CapabilityId, config.Config)
		}
	}

	t.Logf("Created capability DON with %d nodes", len(peerIDs))

	return capReg, regAddress
}
