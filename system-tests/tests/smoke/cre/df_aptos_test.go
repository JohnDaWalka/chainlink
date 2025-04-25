package cre

import (
	"bytes"
	crand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"text/template"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/rpc"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/deployment"
	clclient "github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	libcontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	lidebug "github.com/smartcontractkit/chainlink/system-tests/lib/cre/debug"
	creconsensus "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/consensus"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	creenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	keystonetypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	libtypes "github.com/smartcontractkit/chainlink/system-tests/lib/types"
)

type WorkflowConfigAptos struct {
	WorkflowName string `toml:"workflow_name" validate:"required" `
	FeedID       string `toml:"feed_id" validate:"required,startsnotwith=0x"`
}
type TestConfigAptos struct {
	BlockchainA                   *blockchain.Input                        `toml:"blockchain_a" validate:"required"`
	CustomAnvilMiner              *CustomAnvilMiner                        `toml:"custom_anvil_miner"`
	NodeSets                      []*ns.Input                              `toml:"nodesets" validate:"required"`
	JD                            *jd.Input                                `toml:"jd" validate:"required"`
	Fake                          *fake.Input                              `toml:"fake"`
	KeystoneContracts             *keystonetypes.KeystoneContractsInput    `toml:"keystone_contracts"`
	WorkflowRegistryConfiguration *keystonetypes.WorkflowRegistryInput     `toml:"workflow_registry_configuration"`
	DataFeedsCacheContract        *keystonetypes.DeployDataFeedsCacheInput `toml:"data_feeds_cache"`
	Infra                         *libtypes.InfraInput                     `toml:"infra" validate:"required"`
	WorkflowConfig                *WorkflowConfigAptos                     `toml:"workflow_config" validate:"required"`
	JobConfig                     *JobConfig                               `toml:"job_config"`
}

type JobConfig struct {
	NrOfStreams int `toml:"nr_of_streams"`
}

type registerDFAptosWorkflowInput struct {
	*WorkflowConfig
	chainSelector           uint64
	writeTargetName         string
	workflowDonID           uint32
	feedID                  string
	workflowRegistryAddress common.Address
	dataFeedsCacheAddress   common.Address
	sethClient              *seth.Client
	deployerPrivateKey      string
}

type configureDFAptosCacheInput struct {
	chainSelector         uint64
	fullCldEnvironment    *deployment.Environment
	forwarderAddress      common.Address
	dataFeedsCacheAddress common.Address
	sethClient            *seth.Client
	blockchain            *blockchain.Output
	settingsFile          *os.File
	deployerPrivateKey    string
	feedID                string
	workflowName          string
}

type dfAptosSetupOutput struct {
	dataFeedsCacheAddress common.Address
	forwarderAddress      common.Address
	sethClient            *seth.Client
	blockchainOutput      *blockchain.Output
	donTopology           *keystonetypes.DonTopology
	nodeOutput            []*keystonetypes.WrappedNodeOutput
	fakeDPURL             string
}

func configureDFAptosCacheContract(testLogger zerolog.Logger, input *configureDFAptosCacheInput) error {
	configInput := &keystonetypes.ConfigureDataFeedsCacheInput{
		CldEnv:                input.fullCldEnvironment,
		ChainSelector:         input.chainSelector,
		FeedIDs:               []string{input.feedID},
		Descriptions:          []string{"CRE aptos test feed"},
		DataFeedsCacheAddress: input.dataFeedsCacheAddress,
		AdminAddress:          input.sethClient.MustGetRootKeyAddress(),
		AllowedSenders:        []common.Address{input.forwarderAddress},
		AllowedWorkflowOwners: []common.Address{input.sethClient.MustGetRootKeyAddress()},
		AllowedWorkflowNames:  []string{input.workflowName},
	}

	_, configErr := libcontracts.ConfigureDataFeedsCache(testLogger, configInput)

	return configErr
}

func setupDFAptosTestEnvironment(
	t *testing.T,
	testLogger zerolog.Logger,
	in *TestConfigAptos,
	mustSetCapabilitiesFn func(input []*ns.Input) []*keystonetypes.CapabilitiesAwareNodeSet,
	capabilityFactoryFns []func([]string) []keystone_changeset.DONCapabilityWithConfig,
	jobSpecFactoryFns []keystonetypes.JobSpecFactoryFn,
) *dfAptosSetupOutput {

	universalSetupInput := creenv.SetupInput{
		CapabilitiesAwareNodeSets:            mustSetCapabilitiesFn(in.NodeSets),
		CapabilitiesContractFactoryFunctions: capabilityFactoryFns,
		BlockchainsInput:                     *in.BlockchainA,
		JdInput:                              *in.JD,
		InfraInput:                           *in.Infra,
		JobSpecFactoryFunctions:              jobSpecFactoryFns,
	}

	universalSetupOutput, setupErr := creenv.SetupTestEnvironment(testcontext.Get(t), testLogger, cldlogger.NewSingleFileLogger(t), universalSetupInput)
	require.NoError(t, setupErr, "failed to setup test environment")

	if in.CustomAnvilMiner != nil {
		require.NotContains(t, in.BlockchainA.DockerCmdParamsOverrides, "-b", "custom_anvil_miner was specified but Anvil has '-b' key set, remove that parameter from 'docker_cmd_params' to run deployments instantly or remove custom_anvil_miner key from TOML config")
		require.Equal(t, "anvil", in.BlockchainA.Type, "custom_anvil_miner was specified but blockchain type is not Anvil")
		miner := rpc.NewRemoteAnvilMiner(universalSetupOutput.BlockchainOutput.BlockchainOutput.Nodes[0].ExternalHTTPUrl, nil)
		miner.MinePeriodically(time.Duration(in.CustomAnvilMiner.BlockSpeedSeconds) * time.Second)
	}

	deployDataFeedsInput := &keystonetypes.DeployDataFeedsCacheInput{
		ChainSelector: universalSetupOutput.BlockchainOutput.ChainSelector,
		CldEnv:        universalSetupOutput.CldEnvironment,
	}
	deployDataFeedsCacheOutput, dfErr := libcontracts.DeployDataFeedsCache(testLogger, deployDataFeedsInput)
	require.NoError(t, dfErr, "failed to deploy data feeds cache")

	dfConfigInput := &configureDFAptosCacheInput{
		chainSelector:         universalSetupOutput.BlockchainOutput.ChainSelector,
		fullCldEnvironment:    universalSetupOutput.CldEnvironment,
		forwarderAddress:      universalSetupOutput.KeystoneContractsOutput.ForwarderAddress,
		dataFeedsCacheAddress: deployDataFeedsCacheOutput.DataFeedsCacheAddress,
		sethClient:            universalSetupOutput.BlockchainOutput.SethClient,
		blockchain:            universalSetupOutput.BlockchainOutput.BlockchainOutput,
		deployerPrivateKey:    universalSetupOutput.BlockchainOutput.DeployerPrivateKey,
		feedID:                in.WorkflowConfig.FeedID,
	}
	dfConfigErr := configureDFAptosCacheContract(testLogger, dfConfigInput)
	require.NoError(t, dfConfigErr, "failed to configure data feeds cache")

	// Set inputs in the test config, so that they can be saved
	in.KeystoneContracts = &keystonetypes.KeystoneContractsInput{
		Out: universalSetupOutput.KeystoneContractsOutput,
	}
	in.DataFeedsCacheContract = &keystonetypes.DeployDataFeedsCacheInput{
		Out: &keystonetypes.DeployDataFeedsCacheOutput{
			DataFeedsCacheAddress: deployDataFeedsCacheOutput.DataFeedsCacheAddress,
		},
	}
	in.WorkflowRegistryConfiguration = &keystonetypes.WorkflowRegistryInput{
		Out: universalSetupOutput.WorkflowRegistryConfigurationOutput,
	}

	return &dfAptosSetupOutput{
		dataFeedsCacheAddress: deployDataFeedsCacheOutput.DataFeedsCacheAddress,
		forwarderAddress:      universalSetupOutput.KeystoneContractsOutput.ForwarderAddress,
		sethClient:            universalSetupOutput.BlockchainOutput.SethClient,
		blockchainOutput:      universalSetupOutput.BlockchainOutput.BlockchainOutput,
		donTopology:           universalSetupOutput.DonTopology,
		nodeOutput:            universalSetupOutput.NodeOutput,
	}
}

func TestCRE_DF_Aptos(t *testing.T) {
	testLogger := framework.L

	// Load and validate test configuration
	in, err := framework.Load[TestConfigAptos](t)
	require.NoError(t, err, "couldn't load test config")
	//validateEnvVars(t, in)
	require.Len(t, in.NodeSets, 3, "expected 1 node set in the test config")

	// Assign all capabilities to the single node set
	mustSetCapabilitiesFn := func(input []*ns.Input) []*keystonetypes.CapabilitiesAwareNodeSet {
		return []*keystonetypes.CapabilitiesAwareNodeSet{
			{
				Input:              input[0],
				Capabilities:       []string{keystonetypes.StreamTriggerV1},
				DONTypes:           []string{keystonetypes.CapabilitiesDON},
				BootstrapNodeIndex: 0,
			},
			{
				Input:              input[1],
				Capabilities:       []string{keystonetypes.OCR3Capability},
				DONTypes:           []string{keystonetypes.WorkflowDON},
				BootstrapNodeIndex: 0,
			},
			{
				Input:              input[2],
				Capabilities:       []string{keystonetypes.WriteEVMCapability},
				DONTypes:           []string{keystonetypes.CapabilitiesDON},
				BootstrapNodeIndex: 0,
			},
		}
	}

	fakeDataProvider := fake.Input{
		Port: 4431, //TODO @george-dorin: get from toml
		Out:  nil,
	}

	fakeDPURL, err := setupFakeDataProviderDF(testLogger, &fakeDataProvider)
	require.NoError(t, err, "cannot create fake data provider")

	chainIDInt, chainErr := strconv.Atoi(in.BlockchainA.ChainID)
	require.NoError(t, chainErr, "failed to convert chain ID to int")

	streams := make([]streamsJob, 0)
	for i := range in.JobConfig.NrOfStreams {
		_, feedID := NewFeedID(t)
		streams = append(streams, streamsJob{
			JobName:    fmt.Sprintf("stream job %d", i),
			JobID:      uuid.New().String(),
			StreamID:   i + 1,
			FeedID:     feedID,
			ContractID: "",
		})
	}

	DFJobSpecsFactoryFn := func(input *keystonetypes.JobSpecFactoryInput) (keystonetypes.DonsToJobSpecs, error) {
		donTojobSpecs := make(keystonetypes.DonsToJobSpecs, 0)
		for _, donWithMetadata := range input.DonTopology.DonsWithMetadata {
			jobSpecs := make(keystonetypes.DonJobs, 0)
			workflowNodeSet, err2 := node.FindManyWithLabel(donWithMetadata.NodesMetadata, &keystonetypes.Label{Key: node.NodeTypeKey, Value: keystonetypes.WorkerNode}, node.EqualLabels)
			if err2 != nil {
				// there should be no DON without worker nodes, even gateway DON is composed of a single worker node
				return nil, errors.Wrap(err2, "failed to find worker nodes")
			}
			for _, workerNode := range workflowNodeSet {
				nodeID, nodeIDErr := node.FindLabelValue(workerNode, node.NodeIDKey)
				if nodeIDErr != nil {
					return nil, errors.Wrap(nodeIDErr, "failed to get node id from labels")
				}
				if flags.HasFlag(donWithMetadata.Flags, keystonetypes.StreamTriggerV1) {
					for _, s := range streams {
						jobSpecs = append(jobSpecs, streamsTriggerV1Jobs(s, nodeID))
					}
				}
			}
			donTojobSpecs[donWithMetadata.ID] = jobSpecs

			//Create bridges
			if flags.HasFlag(donWithMetadata.Flags, keystonetypes.StreamTriggerV1) {
				for _, n := range donWithMetadata.DON.Nodes {
					err3 := n.CreateBridge(&clclient.BridgeTypeAttributes{
						Name: "bridge-1",
						URL:  fakeDPURL,
					})
					if err3 != nil {
						return nil, err3
					}
					err3 = n.CreateBridge(&clclient.BridgeTypeAttributes{
						Name: "bridge-2",
						URL:  fakeDPURL,
					})
					if err3 != nil {
						return nil, err3
					}
					err3 = n.CreateBridge(&clclient.BridgeTypeAttributes{
						Name: "bridge-3",
						URL:  fakeDPURL,
					})
					if err3 != nil {
						return nil, err3
					}
				}
			}
		}
		return donTojobSpecs, nil
	}

	require.NoError(t, err, "failed to convert chain ID to int")
	chainIDUint64 := libc.MustSafeUint64(int64(chainIDInt))

	setupOutput := setupDFAptosTestEnvironment(
		t,
		testLogger,
		in,
		mustSetCapabilitiesFn,
		[]keystonetypes.DONCapabilityWithConfigFactoryFn{libcontracts.StreamsTriggerV1CapabilityFactory, libcontracts.DefaultCapabilityFactoryFn, libcontracts.ChainWriterCapabilityFactory(libc.MustSafeUint64(int64(chainIDInt)))},
		[]keystonetypes.JobSpecFactoryFn{
			creconsensus.ConsensusJobSpecFactoryFn(chainIDUint64),
			DFJobSpecsFactoryFn,
		},
	)

	// Log extra information that might help debugging
	t.Cleanup(func() {
		if t.Failed() {
			logTestInfoDFAptos(testLogger, setupOutput.dataFeedsCacheAddress.Hex(), setupOutput.forwarderAddress.Hex())

			// log scanning is not supported for CRIB
			if in.Infra.InfraType == libtypes.CRIB {
				return
			}

			logDir := fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name())

			removeErr := os.RemoveAll(logDir)
			if removeErr != nil {
				testLogger.Error().Err(removeErr).Msg("failed to remove log directory")
				return
			}

			_, saveErr := framework.SaveContainerLogs(logDir)
			if saveErr != nil {
				testLogger.Error().Err(saveErr).Msg("failed to save container logs")
				return
			}

			debugDons := make([]*keystonetypes.DebugDon, 0, len(setupOutput.donTopology.DonsWithMetadata))
			for i, donWithMetadata := range setupOutput.donTopology.DonsWithMetadata {
				containerNames := make([]string, 0, len(donWithMetadata.NodesMetadata))
				for _, output := range setupOutput.nodeOutput[i].Output.CLNodes {
					containerNames = append(containerNames, output.Node.ContainerName)
				}
				debugDons = append(debugDons, &keystonetypes.DebugDon{
					NodesMetadata:  donWithMetadata.NodesMetadata,
					Flags:          donWithMetadata.Flags,
					ContainerNames: containerNames,
				})
			}

			debugInput := keystonetypes.DebugInput{
				DebugDons:        debugDons,
				BlockchainOutput: setupOutput.blockchainOutput,
				InfraInput:       in.Infra,
			}
			lidebug.PrintTestDebug(t.Name(), testLogger, debugInput)
		}
	})
	fmt.Println(setupOutput)
	fmt.Println("+_+_+_+_+_+_+_+_+")

	time.Sleep(time.Hour)
}
func logTestInfoDFAptos(l zerolog.Logger, dataFeedsCacheAddr, forwarderAddr string) {
	l.Info().Msg("------ Test configuration:")
	l.Info().Msgf("DataFeedsCache address: %s", dataFeedsCacheAddr)
	l.Info().Msgf("KeystoneForwarder address: %s", forwarderAddr)
}

func setupFakeDataProviderDF(testLogger zerolog.Logger, input *fake.Input) (string, error) {
	_, err := fake.NewFakeDataProvider(input)
	if err != nil {
		return "", errors.Wrap(err, "failed to set up fake data provider")
	}
	fakeAPIPath := "/fake/api/price"
	host := framework.HostDockerInternal()
	fakeFinalURL := fmt.Sprintf("%s:%d%s", host, input.Port, fakeAPIPath)

	getPriceResponseFn := func() map[string]interface{} {
		response := map[string]interface{}{
			"price": rand.Int(),
		}

		marshalled, mErr := json.Marshal(response)
		if mErr == nil {
			testLogger.Info().Msgf("Returning response: %s", string(marshalled))
		} else {
			testLogger.Info().Msgf("Returning response: %v", response)
		}

		return response
	}

	err = fake.Func("GET", fakeAPIPath, func(c *gin.Context) {
		c.JSON(200, getPriceResponseFn())
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to set up fake data provider")
	}

	return fakeFinalURL, nil
}

type streamsJob struct {
	JobName    string
	JobID      string
	StreamID   int
	FeedID     string
	ContractID string
}

func streamsTriggerV1Jobs(streamsJob streamsJob, nodeID string) *jobv1.ProposeJobRequest {
	job := `
type = "offchainreporting2"
schemaVersion = 1
name = "{{ .JobName }}"
externalJobID = "{{ .JobID }}"
forwardingAllowed = false
maxTaskDuration = "0s"
streamID = {{ .StreamID }}
feedID = "{{ .FeedID }}"
contractID = "{{ .ContractID }}"
relay = "evm"
pluginType = "mercury"
observationSource = """
// data source 1
ds1_payload [type=bridge name="bridge-1" timeout="50s" requestData=""];
ds1_benchmark [type=jsonparse path="price"];

// data source 2
ds2_payload [type=bridge name="bridge-2" timeout="50s" requestData=""];
ds2_benchmark [type=jsonparse path="price"];

// data source 3
ds3_payload [type=bridge name="bridge-3" timeout="50s" requestData=""];
ds3_benchmark [type=jsonparse path="price"];

ds1_payload -> ds1_benchmark -> benchmark_price;
ds2_payload -> ds2_benchmark -> benchmark_price;
ds3_payload -> ds3_benchmark -> benchmark_price;
benchmark_price [type=median allowedFaults=2 index=0];
"""

[relayConfig]
chainID = 1337
enableTriggerCapability = true
`
	tmpl, err := template.New("workflow").Parse(job)

	if err != nil {
		panic(err)
	}
	var renderedTemplate bytes.Buffer
	err = tmpl.Execute(&renderedTemplate, streamsJob)
	if err != nil {
		panic(err)
	}

	return &jobv1.ProposeJobRequest{
		NodeId: nodeID,
		Spec:   renderedTemplate.String()}
}
func NewFeedID(t *testing.T) ([32]byte, string) {
	buf := [32]byte{}
	_, err := crand.Read(buf[:])
	require.NoError(t, err, "cannot create feedID")
	return buf, "0x" + hex.EncodeToString(buf[:])
}
