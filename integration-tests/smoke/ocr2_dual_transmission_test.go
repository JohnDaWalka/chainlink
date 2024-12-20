package smoke

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	gethwrappers "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/dual-transmission"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/testhelpers"
)

type SvrTestConfig struct {
	OverwriteCustomURL  bool   `toml:"overwrite_custom_url"`
	CustomURL           string `toml:"custom_url"`
	TestTimeoutMinutes  int16  `toml:"test_timeout_minutes" validate:"required"`
	ExpectedEventsCount int16  `toml:"expected_event_count" validate:"required"`
}

type Cfg struct {
	BlockchainA *blockchain.Input `toml:"blockchain_a" validate:"required"`
	Contracts   *actions.Input    `toml:"contracts" validate:"required"`
	NodeSet     *ns.Input         `toml:"nodeset" validate:"required"`
	SvrTestCfg  *SvrTestConfig    `toml:"test_config" validate:"required"`
}

var bootstrapJobSpec = `
type = "bootstrap"
schemaVersion = 1
name = "smoke OEV bootstrap"
externalJobID = "%s"
contractID = "%s"
relay = "%s"

[relayConfig]
chainID = %s
`
var oevJobSpec = `
type = "offchainreporting2"
schemaVersion = 1
name = "OEV job"
externalJobID = "%s"
forwardingAllowed = true
maxTaskDuration = "0s"
contractID = "%s"
relay = "%s"
ocrKeyBundleID = "%s"
pluginType = "median"
transmitterID = "%s"
p2pv2Bootstrappers = ["%s@%s"]

observationSource = """
 //randomness
    val1 [type="memo" value="10"]
    val2 [type="memo" value="20"]
    val3 [type="memo" value="30"]
    val4 [type="memo" value="40"]
    val5 [type="memo" value="50"]
    val6 [type="memo" value="60"]
    val7 [type="memo" value="70"]
    val8 [type="memo" value="80"]
    val9 [type="memo" value="90"]

    random1 [type="any"]
    random2 [type="any"]
    random3 [type="any"]

    val1 -> random1
    val2 -> random2
    val3 -> random3
    val4 -> random1
    val5 -> random2
    val6 -> random3
    val7 -> random1
    val8 -> random2
    val9 -> random3


    // data source 1
    ds1_multiply [type="multiply" times=100]

     // data source 2
    ds2_multiply [type="multiply" times=100]


    // data source 3
    ds3_multiply [type="multiply" times=100]


    random1 -> ds1_multiply -> answer
    random2 -> ds2_multiply -> answer
    random3 -> ds3_multiply -> answer

    answer [type=median]
"""

[relayConfig]
chainID = %s
enableDualTransmission = true

[relayConfig.dualTransmission]
contractAddress = "%s"
transmitterAddress = "%s"

[relayConfig.dualTransmission.meta]
hint = [ "calldata" ]
refund = [ "0xbc1Be4cC8790b0C99cff76100E0e6d01E32C6A2C:90" ]

[pluginConfig]
juelsPerFeeCoinSource = """
juels_per_fee_coin [type="sum" values=<[0]>];
"""
`

func TestSmoke(t *testing.T) {
	lggr := logging.GetTestLogger(t)
	pkey := os.Getenv("PRIVATE_KEY")
	require.NotEmpty(t, pkey, "private key is empty")
	in, err := framework.Load[Cfg](t)
	require.NoError(t, err)

	testTimeout, err := time.ParseDuration(fmt.Sprintf("%dm", in.SvrTestCfg.TestTimeoutMinutes)) //TODO: @george-dorin Fix Me!
	require.NoError(t, err)

	chainID := in.BlockchainA.ChainID
	require.NotEmpty(t, chainID, "blockchain_a.chain_id cannot be empty")
	chainFamily := in.BlockchainA.Out.Family
	require.NotEmpty(t, chainFamily, "need to specify a blockchain_a.out.family")

	// 1. Set up nodes & docker environment (5 nodes)
	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)

	// Replace customURL
	if in.SvrTestCfg.OverwriteCustomURL {
		require.NotEmpty(t, in.SvrTestCfg.CustomURL, "when overwrite_custom_url=true you need to provide a non empty custom_url")
		in.NodeSet.NodeSpecs[0].Node.UserConfigOverrides = fmt.Sprintf(in.NodeSet.NodeSpecs[0].Node.UserConfigOverrides, in.SvrTestCfg.CustomURL)
	} else {
		in.NodeSet.NodeSpecs[0].Node.UserConfigOverrides = fmt.Sprintf(in.NodeSet.NodeSpecs[0].Node.UserConfigOverrides, bc.Nodes[0].DockerInternalHTTPUrl)
	}

	out, err := ns.NewSharedDBNodeSet(in.NodeSet, bc)
	require.NoError(t, err)

	// connecting clients
	sethClient, err := seth.NewClientBuilder().
		WithGasPriceEstimations(true, 0, seth.Priority_Fast).
		WithRpcUrl(bc.Nodes[0].HostWSUrl).
		WithPrivateKeys([]string{pkey}).
		Build()
	require.NoError(t, err)
	nodeClients, err := clclient.New(out.CLNodes)
	require.NoError(t, err)
	bootstrapNode, workerNodes := nodeClients[0], nodeClients[1:]

	//2. Create secondary addresses
	primaryAddresses := make([]common.Address, 0)
	secondaryAddresses := make([]common.Address, 0)
	//Create secondary ETH key
	for i := range workerNodes {
		primary, err := workerNodes[i].PrimaryEthAddress()
		require.NoError(t, err, "Cannot get primary key")
		primaryAddresses = append(primaryAddresses, common.HexToAddress(primary))
		key, _, err := workerNodes[i].CreateTxKey(chainFamily, chainID)
		require.NoError(t, err, "Cannot create secondary key")
		secondaryAddresses = append(secondaryAddresses, common.HexToAddress(key.Data.Attributes.Address))
	}

	//3. Restart
	out, err = ns.UpgradeNodeSet(in.NodeSet, bc, time.Second*10)
	require.NoError(t, err, "Cannot restart nodes")

	// Reconnect to clients
	nodeClients, err = clclient.New(out.CLNodes)
	require.NoError(t, err, "Cannot create clients")
	bootstrapNode, workerNodes = nodeClients[0], nodeClients[1:]

	// 4. Fund addresses
	err = ns.FundNodes(sethClient.Client, nodeClients, pkey, 0.1)
	require.NoError(t, err)
	for i := range secondaryAddresses {
		err = ns.SendETH(sethClient.Client, pkey, secondaryAddresses[i].String(), big.NewFloat(0.1))
		require.NoError(t, err, "Cannot fund secondary address")
	}

	// 5. Deploy link contract
	linkContract, err := contracts.DeployLinkTokenContract(lggr, sethClient)
	require.NoError(t, err, "Error loading/deploying link token contract")

	// 6. Deploy forwarders
	var operators []common.Address
	operators, forwarders, _ := actions.DeployForwarderContracts(
		t, sethClient, common.HexToAddress(linkContract.Address()), len(workerNodes),
	)
	require.Equal(t, len(workerNodes), len(operators), "Number of operators should match number of nodes")
	require.Equal(t, len(workerNodes), len(forwarders), "Number of authorized forwarders should match number of nodes")

	// 7. Configure forwarders
	require.NoError(t, err, "Retrieving on-chain wallet addresses for chainlink nodes shouldn't fail")
	for i := range workerNodes {
		actions.AcceptAuthorizedReceiversOperator(
			t, lggr, sethClient, operators[i], forwarders[i], []common.Address{primaryAddresses[i], secondaryAddresses[i]})
		require.NoError(t, err, "Accepting Authorize Receivers on Operator shouldn't fail")

		chainIDBigInt, ok := new(big.Int).SetString(chainID, 10)
		require.True(t, ok, "ChainID cannot be converted to big.Int")
		_, _, err = workerNodes[i].TrackForwarder(chainIDBigInt, forwarders[i])
		require.NoError(t, err, "Cannot track forwarders")
	}

	// 8. Deploy dual agg
	//in.Contracts.URL = bc.Nodes[0].HostWSUrl
	ocrOffchainOptions := contracts.DefaultOffChainAggregatorOptions()
	oevContract, err := actions.NewDualAggregatorDeployment(sethClient, in.Contracts, linkContract.Address(), ocrOffchainOptions)
	require.NoError(t, err)
	dualAggContract, err := gethwrappers.NewDualAggregator(oevContract.Addresses[0], sethClient.Client)
	require.NoError(t, err)

	// 9. Configure dual agg
	config, err := BuildDualAggregatorOCR2ConfigLocal(workerNodes, ocrOffchainOptions)
	require.NoError(t, err, "error creating OEV on-chain config")
	_, err = dualAggContract.SetConfig(sethClient.NewTXOpts(), config.Signers, config.Transmitters, config.F, config.OnchainConfig, config.OffchainConfigVersion, config.OffchainConfig)
	require.NoError(t, err, "error configuring OEV contract")

	// 10. Create jobs
	// Bootstrap
	response, _, err2 := bootstrapNode.CreateJobRaw(fmt.Sprintf(bootstrapJobSpec, uuid.New().String(), oevContract.Addresses[0].String(), chainFamily, chainID))
	require.NoError(t, err2)
	require.Empty(t, response.Errors)

	// Feed job
	bootstrapPeerID, err := bootstrapNode.MustReadP2PKeys()
	require.NoError(t, err, "cannot get bootstrap peerID")
	require.Equal(t, 1, len(bootstrapPeerID.Data), "expected one bootstrap P2P key")

	for i := range workerNodes {
		ocr2Keys, err := workerNodes[i].MustReadOCR2Keys()
		require.NoError(t, err, "cannot fetch OCR2 keys")
		require.Equal(t, 1, len(ocr2Keys.Data), "expecting only one OCR2 key")
		response, _, err := workerNodes[i].CreateJobRaw(fmt.Sprintf(oevJobSpec, uuid.New().String(), oevContract.Addresses[0].String(), chainFamily, ocr2Keys.Data[0].ID, primaryAddresses[i], bootstrapPeerID.Data[0].Attributes.PeerID,
			strings.TrimPrefix(out.CLNodes[0].Node.DockerP2PUrl, "http://"), chainID, oevContract.Addresses[0].String(), secondaryAddresses[i]))
		require.NoError(t, err)
		require.Empty(t, response.Errors)
	}

	t.Run("test SVR transmissions and events", func(t *testing.T) {
		require.NoError(t, waitForDualAggregatorEvents(testcontext.Get(t), dualAggContract, in.SvrTestCfg.ExpectedEventsCount, testTimeout, lggr))
	})
}

func BuildDualAggregatorOCR2ConfigLocal(workerNodes []*clclient.ChainlinkClient, ocrOffchainOptions contracts.OffchainOptions) (*contracts.OCRv2Config, error) {
	S, oracleIdentities, err := getOracleIdentitiesWithKeyIndexLocal(workerNodes, 0)
	if err != nil {
		return nil, err
	}
	signerKeys, _, f_, _, offchainConfigVersion, offchainConfig, err := confighelper.ContractSetConfigArgsForTests(
		30*time.Second,   // deltaProgress time.Duration,
		30*time.Second,   // deltaResend time.Duration,
		10*time.Second,   // deltaRound time.Duration,
		20*time.Second,   // deltaGrace time.Duration,
		20*time.Second,   // deltaStage time.Duration,
		3,                // rMax uint8,
		S,                // s []int,
		oracleIdentities, // oracles []OracleIdentityExtra,
		median.OffchainConfig{
			AlphaReportInfinite: false,
			AlphaReportPPB:      1,
			AlphaAcceptInfinite: false,
			AlphaAcceptPPB:      1,
			DeltaC:              time.Minute * 30,
		}.Encode(), // reportingPluginConfig []byte,
		nil,
		5*time.Second, // maxDurationQuery time.Duration,
		5*time.Second, // maxDurationObservation time.Duration,
		5*time.Second, // maxDurationReport time.Duration,
		5*time.Second, // maxDurationShouldAcceptFinalizedReport time.Duration,
		5*time.Second, // maxDurationShouldTransmitAcceptedReport time.Duration,
		1,             // f int,
		nil,           // The median reporting plugin has an empty onchain config
	)
	if err != nil {
		return nil, err
	}

	// Convert signers to addresses
	var signerAddresses []common.Address
	for _, signer := range signerKeys {
		signerAddresses = append(signerAddresses, common.BytesToAddress(signer))
	}

	// Replace transmitter with forwaders
	var transmitterAddresses []common.Address
	for i := range workerNodes {
		t, _, err := workerNodes[i].GetForwarders()
		if err != nil {
			return nil, errors.New("cannot get forwarder from node")
		}
		if len(t.Data) < 1 {
			return nil, errors.New("no forwarders found on node")
		}
		transmitterAddresses = append(transmitterAddresses, common.HexToAddress(t.Data[0].Attributes.Address))
	}

	onchainConfig, err := testhelpers.GenerateDefaultOCR2OnchainConfig(ocrOffchainOptions.MinimumAnswer, ocrOffchainOptions.MaximumAnswer)

	return &contracts.OCRv2Config{
		Signers:               signerAddresses,
		Transmitters:          transmitterAddresses,
		F:                     f_,
		OnchainConfig:         onchainConfig,
		OffchainConfigVersion: offchainConfigVersion,
		OffchainConfig:        []byte(fmt.Sprintf("0x%s", offchainConfig)),
	}, err
}

func getOracleIdentitiesWithKeyIndexLocal(
	chainlinkNodes []*clclient.ChainlinkClient,
	keyIndex int,
) ([]int, []confighelper.OracleIdentityExtra, error) {
	S := make([]int, len(chainlinkNodes))
	oracleIdentities := make([]confighelper.OracleIdentityExtra, len(chainlinkNodes))
	sharedSecretEncryptionPublicKeys := make([]types.ConfigEncryptionPublicKey, len(chainlinkNodes))
	eg := &errgroup.Group{}
	for i, cl := range chainlinkNodes {
		index, chainlinkNode := i, cl
		eg.Go(func() error {
			addresses, err := chainlinkNode.EthAddresses()
			if err != nil {
				return err
			}
			ocr2Keys, err := chainlinkNode.MustReadOCR2Keys()
			if err != nil {
				return err
			}
			var ocr2Config nodeclient.OCR2KeyAttributes
			for _, key := range ocr2Keys.Data {
				if key.Attributes.ChainType == string(chaintype.EVM) {
					ocr2Config = nodeclient.OCR2KeyAttributes(key.Attributes)
					break
				}
			}

			keys, err := chainlinkNode.MustReadP2PKeys()
			if err != nil {
				return err
			}
			p2pKeyID := keys.Data[0].Attributes.PeerID

			offchainPkBytes, err := hex.DecodeString(strings.TrimPrefix(ocr2Config.OffChainPublicKey, "ocr2off_evm_"))
			if err != nil {
				return err
			}

			offchainPkBytesFixed := [ed25519.PublicKeySize]byte{}
			n := copy(offchainPkBytesFixed[:], offchainPkBytes)
			if n != ed25519.PublicKeySize {
				return fmt.Errorf("wrong number of elements copied")
			}

			configPkBytes, err := hex.DecodeString(strings.TrimPrefix(ocr2Config.ConfigPublicKey, "ocr2cfg_evm_"))
			if err != nil {
				return err
			}

			configPkBytesFixed := [ed25519.PublicKeySize]byte{}
			n = copy(configPkBytesFixed[:], configPkBytes)
			if n != ed25519.PublicKeySize {
				return fmt.Errorf("wrong number of elements copied")
			}

			onchainPkBytes, err := hex.DecodeString(strings.TrimPrefix(ocr2Config.OnChainPublicKey, "ocr2on_evm_"))
			if err != nil {
				return err
			}

			sharedSecretEncryptionPublicKeys[index] = configPkBytesFixed
			oracleIdentities[index] = confighelper.OracleIdentityExtra{
				OracleIdentity: confighelper.OracleIdentity{
					OnchainPublicKey:  onchainPkBytes,
					OffchainPublicKey: offchainPkBytesFixed,
					PeerID:            p2pKeyID,
					TransmitAccount:   types.Account(addresses[keyIndex]),
				},
				ConfigEncryptionPublicKey: configPkBytesFixed,
			}
			S[index] = 1
			return nil
		})
	}

	return S, oracleIdentities, eg.Wait()
}

func waitForDualAggregatorEvents(
	ctx context.Context,
	dualAggContract *gethwrappers.DualAggregator,
	nrOfEvents int16,
	timeout time.Duration,
	lggr zerolog.Logger,
) error {
	possibleRounds := make([]uint32, 0)
	for i := range 1000 {
		possibleRounds = append(possibleRounds, uint32(i))
	}

	startBlock := uint64(1) //TODO: @george-dorin Fix Me!
	lggr.Info().Msg("Waiting for Dual aggregator events")
	var confirmedPrimary, confirmedSecondary int16

	//Subscribe to events
	/*
		Primary first - NewTransmission
		Primary second - PrimaryFeedUnlocked
		Secondary first - SecondaryRoundIdUpdated NewTransmission
		Secondary second - SecondaryRoundIdUpdated
	*/

	ntSink := make(chan *gethwrappers.DualAggregatorNewTransmission, 100)
	newTransmissionsEvents, err := dualAggContract.WatchNewTransmission(&bind.WatchOpts{Context: ctx, Start: &startBlock}, ntSink, possibleRounds)
	if err != nil {
		return err
	}
	defer newTransmissionsEvents.Unsubscribe()

	sruSink := make(chan *gethwrappers.DualAggregatorSecondaryRoundIdUpdated, 100)
	sruEvents, err := dualAggContract.WatchSecondaryRoundIdUpdated(&bind.WatchOpts{Context: ctx, Start: &startBlock}, sruSink, possibleRounds)
	if err != nil {
		return err
	}
	defer sruEvents.Unsubscribe()

	for {
		select {
		case ret := <-ntSink:
			lggr.Info().Msg(fmt.Sprintf("Received NewTransmission event from %s for roundID %d", ret.Transmitter.String(), ret.AggregatorRoundId))
			confirmedPrimary++
			if confirmedPrimary > nrOfEvents && confirmedSecondary > nrOfEvents {
				return nil
			}
		case ret := <-sruSink:
			lggr.Info().Msg(fmt.Sprintf("Received SecondaryRoundIdUpdated  for roundID %d", ret.SecondaryRoundId))
			confirmedSecondary++
			if confirmedPrimary > nrOfEvents && confirmedSecondary > nrOfEvents {
				return nil
			}
		case <-time.After(timeout):
			return fmt.Errorf("timeout waiting for dual aggregator transmission events")
		}
	}

	return nil
}
