package cre

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/rpc"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/deployment"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	libcontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	lidebug "github.com/smartcontractkit/chainlink/system-tests/lib/cre/debug"
	creconsensus "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/consensus"
	creenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
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
) *dfAptosSetupOutput {
	chainIDInt, err := strconv.Atoi(in.BlockchainA.ChainID)
	require.NoError(t, err, "failed to convert chain ID to int")
	chainIDUint64 := libc.MustSafeUint64(int64(chainIDInt))

	universalSetupInput := creenv.SetupInput{
		CapabilitiesAwareNodeSets:            mustSetCapabilitiesFn(in.NodeSets),
		CapabilitiesContractFactoryFunctions: capabilityFactoryFns,
		BlockchainsInput:                     *in.BlockchainA,
		JdInput:                              *in.JD,
		InfraInput:                           *in.Infra,
		JobSpecFactoryFunctions: []keystonetypes.JobSpecFactoryFn{
			creconsensus.ConsensusJobSpecFactoryFn(chainIDUint64),
		},
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

	chainIDInt, chainErr := strconv.Atoi(in.BlockchainA.ChainID)
	require.NoError(t, chainErr, "failed to convert chain ID to int")

	setupOutput := setupDFAptosTestEnvironment(
		t,
		testLogger,
		in,
		mustSetCapabilitiesFn,
		[]keystonetypes.DONCapabilityWithConfigFactoryFn{libcontracts.StreamsTriggerV1CapabilityFactory, libcontracts.DefaultCapabilityFactoryFn, libcontracts.ChainWriterCapabilityFactory(libc.MustSafeUint64(int64(chainIDInt)))},
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
}
func logTestInfoDFAptos(l zerolog.Logger, dataFeedsCacheAddr, forwarderAddr string) {
	l.Info().Msg("------ Test configuration:")
	l.Info().Msgf("DataFeedsCache address: %s", dataFeedsCacheAddr)
	l.Info().Msgf("KeystoneForwarder address: %s", forwarderAddr)
}
