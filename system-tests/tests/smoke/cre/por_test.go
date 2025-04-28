package cre

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	df_changeset "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	df_changeset_types "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/rpc"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/data_feeds_cache"
	"github.com/smartcontractkit/chainlink/deployment"
	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"
	corevm "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"

	ctfconfig "github.com/smartcontractkit/chainlink-testing-framework/lib/config"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	libcontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	lidebug "github.com/smartcontractkit/chainlink/system-tests/lib/cre/debug"
	crechainreader "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/chainreader"
	crecompute "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/compute"
	creconsensus "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/consensus"
	crecron "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/cron"
	cregateway "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/gateway"
	creenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	keystonetypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	libcrecli "github.com/smartcontractkit/chainlink/system-tests/lib/crecli"
	libfunding "github.com/smartcontractkit/chainlink/system-tests/lib/funding"
	libtypes "github.com/smartcontractkit/chainlink/system-tests/lib/types"
)

var (
	SinglePoRDonCapabilitiesFlags = []string{"ocr3", "cron", "custom-compute", "write-evm"}
)

type CustomAnvilMiner struct {
	BlockSpeedSeconds int `toml:"block_speed_seconds"`
}

type TestConfig struct {
	Blockchains                   []*blockchain.Input                  `toml:"blockchains" validate:"required"`
	CustomAnvilMiner              *CustomAnvilMiner                    `toml:"custom_anvil_miner"`
	NodeSets                      []*ns.Input                          `toml:"nodesets" validate:"required"`
	WorkflowConfigs               []WorkflowConfig                     `toml:"workflow_configs" validate:"required"`
	JD                            *jd.Input                            `toml:"jd" validate:"required"`
	Fake                          *fake.Input                          `toml:"fake"`
	WorkflowRegistryConfiguration *keystonetypes.WorkflowRegistryInput `toml:"workflow_registry_configuration"`
	Infra                         *libtypes.InfraInput                 `toml:"infra" validate:"required"`
	DependenciesConfig            *DependenciesConfig                  `toml:"dependencies" validate:"required"`
}

// noCRENoCompilation is a custom validator for the tag "no_cre_no_compilation".
// It ensures that if UseCRECLI is false, then ShouldCompileNewWorkflow must also be false.
func noCRENoCompilation(fl validator.FieldLevel) bool {
	// Use Parent() to access the WorkflowConfig struct.
	wc, ok := fl.Parent().Interface().(WorkflowConfig)
	if !ok {
		return false
	}
	// If not using CRE CLI and ShouldCompileNewWorkflow is true, fail validation.
	if !wc.UseCRECLI && fl.Field().Bool() {
		return false
	}
	return true
}

func disabledInCI(fl validator.FieldLevel) bool {
	if os.Getenv("CI") == "true" {
		return !fl.Field().Bool()
	}

	return true
}

func registerNoCRENoCompilationTranslation(v *validator.Validate, trans ut.Translator) {
	_ = v.RegisterTranslation("no_cre_no_compilation", trans, func(ut ut.Translator) error {
		return ut.Add("no_cre_no_compilation", "{0} must be false when UseCRECLI is false, it is not possible to compile a workflow without it", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("no_cre_no_compilation", fe.Field())
		return t
	})
}

func registerNoFolderLocationTranslation(v *validator.Validate, trans ut.Translator) {
	_ = v.RegisterTranslation("folder_required_if_compiling", trans, func(ut ut.Translator) error {
		return ut.Add("folder_required_if_compiling", "{0} must set, when compiling the workflow", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("folder_required_if_compiling", fe.Field())
		return t
	})
}

func init() {
	err := framework.Validator.RegisterValidation("no_cre_no_compilation", noCRENoCompilation)
	if err != nil {
		panic(errors.Wrap(err, "failed to register no_cre_no_compilation validator"))
	}
	err = framework.Validator.RegisterValidation("disabled_in_ci", disabledInCI)
	if err != nil {
		panic(errors.Wrap(err, "failed to register disabled_in_ci validator"))
	}

	if framework.ValidatorTranslator != nil {
		registerNoCRENoCompilationTranslation(framework.Validator, framework.ValidatorTranslator)
		registerNoFolderLocationTranslation(framework.Validator, framework.ValidatorTranslator)
	}
}

// Defines the location of the binary files that are required to run the test
// When test runs in CI hardcoded versions will be downloaded before the test starts
// Command that downloads them is part of "test_cmd" in .github/e2e-tests.yml file
type DependenciesConfig struct {
	ReadContractCapabilityBinaryPath string `toml:"read_contract_capability_binary_path"`
	CronCapabilityBinaryPath         string `toml:"cron_capability_binary_path"`
	CRECLIBinaryPath                 string `toml:"cre_cli_binary_path" validate:"required"`
}

func validateEnvVars(t *testing.T, in *TestConfig) {
	require.NotEmpty(t, os.Getenv("PRIVATE_KEY"), "PRIVATE_KEY env var must be set")

	// this is a small hack to avoid changing the reusable workflow
	if os.Getenv("CI") == "true" {
		// This part should ideally happen outside of the test, but due to how our reusable e2e test workflow is structured now
		// we cannot execute this part in workflow steps (it doesn't support any pre-execution hooks)
		require.NotEmpty(t, os.Getenv(ctfconfig.E2E_TEST_CHAINLINK_IMAGE_ENV), "missing env var: "+ctfconfig.E2E_TEST_CHAINLINK_IMAGE_ENV)
		require.NotEmpty(t, os.Getenv(ctfconfig.E2E_TEST_CHAINLINK_VERSION_ENV), "missing env var: "+ctfconfig.E2E_TEST_CHAINLINK_VERSION_ENV)
		require.NotEmpty(t, os.Getenv(creenv.E2eJobDistributorImageEnvVarName), "missing env var: "+creenv.E2eJobDistributorImageEnvVarName)
		require.NotEmpty(t, os.Getenv(creenv.E2eJobDistributorVersionEnvVarName), "missing env var: "+creenv.E2eJobDistributorVersionEnvVarName)
	}

	for _, workflowConfig := range in.WorkflowConfigs {
		if workflowConfig.UseCRECLI {
			if workflowConfig.ShouldCompileNewWorkflow {
				gistWriteToken := os.Getenv("GIST_WRITE_TOKEN")
				require.NotEmpty(t, gistWriteToken, "GIST_WRITE_TOKEN must be set to use CRE CLI to compile workflows. It requires gist:read and gist:write permissions")
				err := os.Setenv("CRE_GITHUB_API_TOKEN", gistWriteToken)
				require.NoError(t, err, "failed to set CRE_GITHUB_API_TOKEN env var")

				// set it only for the first workflow config, since it will be used for all workflows
				break
			}
		}
	}
}

type configureDataFeedsCacheInput struct {
	useCRECLI          bool
	chainSelector      uint64
	fullCldEnvironment *deployment.Environment
	workflowName       string
	feedID             string
	sethClient         *seth.Client
	blockchain         *blockchain.Output
	creCLIAbsPath      string
	settingsFile       *os.File
	deployerPrivateKey string
}

func configureDataFeedsCacheContract(testLogger zerolog.Logger, input *configureDataFeedsCacheInput) error {
	chainIDInt, intErr := strconv.Atoi(input.blockchain.ChainID)
	if intErr != nil {
		return errors.Wrap(intErr, "failed to convert chain ID to int")
	}

	forwarderAddress, forwarderErr := crecontracts.FindAddressesForChain(input.fullCldEnvironment.ExistingAddresses, input.chainSelector, keystone_changeset.KeystoneForwarder.String()) //nolint:staticcheck // won't migrate now
	if forwarderErr != nil {
		return errors.Wrapf(forwarderErr, "failed to find forwarder address for chain %d", input.chainSelector)
	}

	dataFeedsCacheAddress, dataFeedsCacheErr := crecontracts.FindAddressesForChain(input.fullCldEnvironment.ExistingAddresses, input.chainSelector, df_changeset.DataFeedsCache.String()) //nolint:staticcheck // won't migrate now
	if dataFeedsCacheErr != nil {
		return errors.Wrapf(dataFeedsCacheErr, "failed to find data feeds cache address for chain %d", input.chainSelector)
	}

	if input.useCRECLI {
		// These two env vars are required by the CRE CLI
		err := os.Setenv("CRE_ETH_PRIVATE_KEY", input.deployerPrivateKey)
		if err != nil {
			return errors.Wrap(err, "failed to set CRE_ETH_PRIVATE_KEY")
		}

		dfAdminErr := libcrecli.SetFeedAdmin(input.creCLIAbsPath, chainIDInt, input.sethClient.MustGetRootKeyAddress(), input.settingsFile)
		if dfAdminErr != nil {
			return errors.Wrap(dfAdminErr, "failed to set feed admin")
		}

		cleanFeedID := strings.TrimPrefix(input.feedID, "0x")

		// Ensure the feed ID is long enough
		if len(cleanFeedID) < 14 { // Need at least 7 bytes (14 hex chars)
			return fmt.Errorf("feed ID too short: %s", input.feedID)
		} else if len(cleanFeedID) > 32 {
			cleanFeedID = cleanFeedID[:32]
		}

		// Extract decimals from feed ID
		decimals, decimalsErr := df_changeset.GetDecimalsFromFeedID(cleanFeedID)
		if decimalsErr != nil {
			return errors.Wrapf(decimalsErr, "failed to get decimals from feed ID %s", input.feedID)
		}

		dfConfigErr := libcrecli.SetFeedConfig(
			input.creCLIAbsPath,
			input.feedID,
			strconv.Itoa(int(decimals)),
			"PoR test feed",
			chainIDInt,
			[]common.Address{forwarderAddress},
			[]common.Address{input.sethClient.MustGetRootKeyAddress()},
			[]string{input.workflowName},
			input.settingsFile,
		)

		if dfConfigErr != nil {
			return errors.Wrap(dfConfigErr, "failed to set feed config")
		}

		return nil
	}

	configInput := &keystonetypes.ConfigureDataFeedsCacheInput{
		CldEnv:                input.fullCldEnvironment,
		ChainSelector:         input.chainSelector,
		FeedIDs:               []string{input.feedID},
		Descriptions:          []string{"PoR test feed"},
		DataFeedsCacheAddress: dataFeedsCacheAddress,
		AdminAddress:          input.sethClient.MustGetRootKeyAddress(),
		AllowedSenders:        []common.Address{forwarderAddress},
		AllowedWorkflowNames:  []string{input.workflowName},
		AllowedWorkflowOwners: []common.Address{input.sethClient.MustGetRootKeyAddress()},
	}

	_, configErr := libcontracts.ConfigureDataFeedsCache(testLogger, configInput)

	return configErr
}

func logTestInfo(l zerolog.Logger, feedID, workflowName, dataFeedsCacheAddr, forwarderAddr string) {
	l.Info().Msg("------ Test configuration:")
	l.Info().Msgf("Feed ID: %s", feedID)
	l.Info().Msgf("Workflow name: %s", workflowName)
	l.Info().Msgf("DataFeedsCache address: %s", dataFeedsCacheAddr)
	l.Info().Msgf("KeystoneForwarder address: %s", forwarderAddr)
}

type porSetupOutput struct {
	priceProvider                   PriceProvider
	addressBook                     deployment.AddressBook
	chainSelectorToSethClient       map[uint64]*seth.Client
	chainSelectorToBlockchainOutput map[uint64]*blockchain.Output
	donTopology                     *keystonetypes.DonTopology
	nodeOutput                      []*keystonetypes.WrappedNodeOutput
	chainSelectorToWorkflowConfig   map[uint64]WorkflowConfig
}

func (th *testHarness) setupPoRTestEnvironment(
	t *testing.T,
	testLogger zerolog.Logger,
	in *TestConfig,
	priceProvider PriceProvider,
	mustSetCapabilitiesFn func(input []*ns.Input) []*keystonetypes.CapabilitiesAwareNodeSet,
	capabilityFactoryFns []func([]string) []keystone_changeset.DONCapabilityWithConfig,
) *porSetupOutput {
	extraAllowedPorts := make([]int, 0)
	if in.Fake != nil {
		if _, ok := priceProvider.(*FakePriceProvider); ok {
			extraAllowedPorts = append(extraAllowedPorts, in.Fake.Port)
		}
	}

	customBinariesPaths := map[string]string{}
	containerPath, pathErr := capabilities.DefaultContainerDirectory(in.Infra.InfraType)
	require.NoError(t, pathErr, "failed to get default container directory")

	var cronBinaryPathInTheContainer string
	if in.DependenciesConfig.CronCapabilityBinaryPath != "" {
		// where cron binary is located in the container
		cronBinaryPathInTheContainer = filepath.Join(containerPath, filepath.Base(in.DependenciesConfig.CronCapabilityBinaryPath))
		// where cron binary is located on the host
		customBinariesPaths[keystonetypes.CronCapability] = in.DependenciesConfig.CronCapabilityBinaryPath
	} else {
		// assume that if cron binary is already in the image it is in the default location and has default name
		cronBinaryPathInTheContainer = filepath.Join(containerPath, "cron")
	}

	var readContractBinaryPathInTheContainer string
	if in.DependenciesConfig.CronCapabilityBinaryPath != "" {
		// where cron binary is located in the container
		readContractBinaryPathInTheContainer = filepath.Join(containerPath, filepath.Base(in.DependenciesConfig.ReadContractCapabilityBinaryPath))

		// where cron binary is located on the host
		customBinariesPaths[keystonetypes.ReadContractCapability] = in.DependenciesConfig.ReadContractCapabilityBinaryPath
	} else {
		// assume that if cron binary is already in the image it is in the default location and has default name
		readContractBinaryPathInTheContainer = filepath.Join(containerPath, "cron")
	}

	firstBlockchain := in.Blockchains[0]

	chainIDInt, err := strconv.Atoi(firstBlockchain.ChainID)
	require.NoError(t, err, "failed to convert chain ID to int")
	chainIDUint64 := libc.MustSafeUint64(int64(chainIDInt))

	jobSpecFactoryFuncs := []keystonetypes.JobSpecFactoryFn{
		creconsensus.ConsensusJobSpecFactoryFn(chainIDUint64),
		crecron.CronJobSpecFactoryFn(cronBinaryPathInTheContainer),
		cregateway.GatewayJobSpecFactoryFn(extraAllowedPorts, []string{}, []string{"0.0.0.0/0"}),
		crecompute.ComputeJobSpecFactoryFn,
	}

	if th.readContractInput != nil {
		jobSpecFactoryFuncs = append(
			jobSpecFactoryFuncs,
			crechainreader.ChainReaderJobSpecFactoryFn(chainIDUint64, "evm", "", readContractBinaryPathInTheContainer),
		)
	}

	universalSetupInput := creenv.SetupInput{
		CapabilitiesAwareNodeSets:            mustSetCapabilitiesFn(in.NodeSets),
		CapabilitiesContractFactoryFunctions: capabilityFactoryFns,
		BlockchainsInput:                     in.Blockchains,
		JdInput:                              *in.JD,
		InfraInput:                           *in.Infra,
		CustomBinariesPaths:                  customBinariesPaths,
		ExtraAllowedPorts:                    extraAllowedPorts,
		JobSpecFactoryFunctions:              jobSpecFactoryFuncs,
	}

	universalSetupOutput, setupErr := creenv.SetupTestEnvironment(testcontext.Get(t), testLogger, cldlogger.NewSingleFileLogger(t), universalSetupInput)
	require.NoError(t, setupErr, "failed to setup test environment")
	homeChainOutput := universalSetupOutput.BlockchainOutput[0]

	if in.CustomAnvilMiner != nil {
		for _, bi := range universalSetupInput.BlockchainsInput {
			if bi.Type == blockchain.TypeAnvil {
				require.NotContains(t, bi.DockerCmdParamsOverrides, "-b", "custom_anvil_miner was specified but Anvil has '-b' key set, remove that parameter from 'docker_cmd_params' to run deployments instantly or remove custom_anvil_miner key from TOML config")
			}
		}
		for _, bo := range universalSetupOutput.BlockchainOutput {
			if bo.BlockchainOutput.Type == blockchain.TypeAnvil {
				miner := rpc.NewRemoteAnvilMiner(bo.BlockchainOutput.Nodes[0].ExternalHTTPUrl, nil)
				miner.MinePeriodically(time.Duration(in.CustomAnvilMiner.BlockSpeedSeconds) * time.Second)
			}
		}
	}

	chainSelectorToWorkflowConfig := make(map[uint64]WorkflowConfig)
	chainSelectorToSethClient := make(map[uint64]*seth.Client)
	chainSelectorToBlockchainOutput := make(map[uint64]*blockchain.Output)

	for idx, bo := range universalSetupOutput.BlockchainOutput {
		chainSelectorToWorkflowConfig[bo.ChainSelector] = in.WorkflowConfigs[idx]
		chainSelectorToSethClient[bo.ChainSelector] = bo.SethClient
		chainSelectorToBlockchainOutput[bo.ChainSelector] = bo.BlockchainOutput

		deployConfig := df_changeset_types.DeployConfig{
			ChainsToDeploy: []uint64{bo.ChainSelector},
			Labels:         []string{"data-feeds"}, // label required by the changeset
		}

		dfOutput, dfErr := df_changeset.RunChangeset(df_changeset.DeployCacheChangeset, *universalSetupOutput.CldEnvironment, deployConfig)
		require.NoError(t, dfErr, "failed to deploy data feed cache contract")

		mergeErr := universalSetupOutput.CldEnvironment.ExistingAddresses.Merge(dfOutput.AddressBook) //nolint:staticcheck // won't migrate now
		require.NoError(t, mergeErr, "failed to merge address book")

		wfRegistrar := (&workflowRegistrar{
			configBuilder: buildPoRConfig,
		})

		var rci *readContractInput
		if th.readContractInput != nil {
			// override the config builder
			wfRegistrar.configBuilder = buildReadContractConfig

			// Deploy a balance reader and merge to address book
			br, err := keystone_changeset.DeployBalanceReader(*universalSetupOutput.CldEnvironment, keystone_changeset.DeployBalanceReaderRequest{
				ChainSelectors: []uint64{bo.ChainSelector},
			})
			require.NoError(t, err, "failed to deploy balance reader contract")

			require.NoError(t,
				universalSetupOutput.CldEnvironment.ExistingAddresses.Merge(br.AddressBook), //nolint:staticcheck // won't migrate yet
				"failed to merge address book with balance reader",
			)

			// create a new address and fund it
			pub, _, err := seth.NewAddress()
			require.NoError(t, err, "failed to generate new address")
			fundedAddress := common.HexToAddress(pub)

			_, fundingErr := libfunding.SendFunds(zerolog.Logger{}, bo.SethClient, libtypes.FundsToSend{
				ToAddress:  fundedAddress,
				Amount:     th.readContractInput.ExpectedFundingAmount,
				PrivateKey: bo.SethClient.MustGetRootPrivateKey(),
			})

			require.NoError(t, fundingErr, "failed to fund address %s", fundedAddress)
			rci = &readContractInput{
				fundedAddress:        fundedAddress.Hex(),
				readTargetName:       fmt.Sprintf("read-contract-%s-%d@1.0.0", "evm", bo.ChainID),
				contractReaderConfig: `{"contracts":{"BalanceReader":{"contractABI":"[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"addresses\",\"type\":\"address[]\"}],\"name\":\"getNativeBalances\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]","contractPollingFilter":{"genericEventNames":null,"pollingFilter":{"topic2":null,"topic3":null,"topic4":null,"retention":"0s","maxLogsKept":0,"logsPerBlock":0}},"configs":{"getNativeBalances":"{  \"chainSpecificName\": \"getNativeBalances\"}"}}}}`,
				contractName:         "BalanceReader",
				contractMethod:       "getNativeBalances",
			}
		}

		var creCLIAbsPath string
		var creCLISettingsFile *os.File
		if in.WorkflowConfigs[idx].UseCRECLI {
			// make sure that path is indeed absolute
			var pathErr error
			creCLIAbsPath, pathErr = filepath.Abs(in.DependenciesConfig.CRECLIBinaryPath)
			require.NoError(t, pathErr, "failed to get absolute path for CRE CLI")

			// create CRE CLI settings file
			var settingsErr error
			creCLISettingsFile, settingsErr = libcrecli.PrepareCRECLISettingsFile(
				bo.SethClient.MustGetRootKeyAddress(),
				universalSetupOutput.CldEnvironment.ExistingAddresses, //nolint:staticcheck // won't migrate now
				universalSetupOutput.DonTopology.WorkflowDonID,
				homeChainOutput.ChainSelector,
				map[uint64]string{
					homeChainOutput.ChainSelector: homeChainOutput.BlockchainOutput.Nodes[0].ExternalHTTPUrl,
					bo.ChainSelector:              bo.BlockchainOutput.Nodes[0].ExternalHTTPUrl,
				},
			)
			require.NoError(t, settingsErr, "failed to create CRE CLI settings file")
		}

		dfConfigInput := &configureDataFeedsCacheInput{
			useCRECLI:          in.WorkflowConfigs[idx].UseCRECLI,
			chainSelector:      bo.ChainSelector,
			fullCldEnvironment: universalSetupOutput.CldEnvironment,
			workflowName:       in.WorkflowConfigs[idx].WorkflowName,
			feedID:             in.WorkflowConfigs[idx].FeedID,
			sethClient:         bo.SethClient,
			blockchain:         bo.BlockchainOutput,
			creCLIAbsPath:      creCLIAbsPath,
			settingsFile:       creCLISettingsFile,
			deployerPrivateKey: bo.DeployerPrivateKey,
		}
		dfConfigErr := configureDataFeedsCacheContract(testLogger, dfConfigInput)
		require.NoError(t, dfConfigErr, "failed to configure data feeds cache")

		registerInput := registerPoRWorkflowInput{
			WorkflowConfig:     in.WorkflowConfigs[idx],
			homeChainSelector:  homeChainOutput.ChainSelector,
			chainSelector:      bo.ChainSelector,
			workflowDonID:      universalSetupOutput.DonTopology.WorkflowDonID,
			feedID:             in.WorkflowConfigs[idx].FeedID,
			addressBook:        universalSetupOutput.CldEnvironment.ExistingAddresses, //nolint:staticcheck // won't migrate now
			priceProvider:      priceProvider,
			sethClient:         bo.SethClient,
			deployerPrivateKey: bo.DeployerPrivateKey,
			creCLIAbsPath:      creCLIAbsPath,
			creCLIsettingsFile: creCLISettingsFile,
			writeTargetName:    corevm.GenerateWriteTargetName(bo.ChainID),
			readContractInput:  rci,
		}

		workflowErr := wfRegistrar.registerWorkflow(registerInput)
		require.NoError(t, workflowErr, "failed to register PoR workflow")
	}
	// Workflow-specific configuration -- END

	// TODO use address book to save the contract addresses

	return &porSetupOutput{
		priceProvider:                   priceProvider,
		chainSelectorToSethClient:       chainSelectorToSethClient,
		chainSelectorToBlockchainOutput: chainSelectorToBlockchainOutput,
		donTopology:                     universalSetupOutput.DonTopology,
		nodeOutput:                      universalSetupOutput.NodeOutput,
		addressBook:                     universalSetupOutput.CldEnvironment.ExistingAddresses, //nolint:staticcheck // won't migrate now
		chainSelectorToWorkflowConfig:   chainSelectorToWorkflowConfig,
	}
}

// config file to use: environment-one-don-multichain.toml
func TestCRE_OCR3_PoR_Workflow_SingleDon_MultipleWriters_MockedPrice(t *testing.T) {
	testLogger := framework.L

	// Load and validate test configuration
	in, err := framework.Load[TestConfig](t)
	require.NoError(t, err, "couldn't load test config")
	validateEnvVars(t, in)
	require.Len(t, in.NodeSets, 1, "expected 1 node set in the test config")

	// Assign all capabilities to the single node set
	mustSetCapabilitiesFn := func(input []*ns.Input) []*keystonetypes.CapabilitiesAwareNodeSet {
		return []*keystonetypes.CapabilitiesAwareNodeSet{
			{
				Input:              input[0],
				Capabilities:       SinglePoRDonCapabilitiesFlags,
				DONTypes:           []string{keystonetypes.WorkflowDON, keystonetypes.GatewayDON},
				BootstrapNodeIndex: 0, // not required, but set to make the configuration explicit
				GatewayNodeIndex:   0, // not required, but set to make the configuration explicit
			},
		}
	}

	feedIDs := make([]string, 0, len(in.WorkflowConfigs))
	for _, wc := range in.WorkflowConfigs {
		feedIDs = append(feedIDs, wc.FeedID)
	}

	priceProvider, priceErr := NewFakePriceProvider(testLogger, in.Fake, AuthorizationKey, feedIDs)
	require.NoError(t, priceErr, "failed to create fake price provider")

	homeChain := in.Blockchains[0]
	targetChain := in.Blockchains[1]
	homeChainID, chainErr := strconv.Atoi(homeChain.ChainID)
	require.NoError(t, chainErr, "failed to convert home chain ID to int")
	targetChainID, chainErr := strconv.Atoi(targetChain.ChainID)
	require.NoError(t, chainErr, "failed to convert target chain ID to int")

	setupOutput := new(testHarness).setupPoRTestEnvironment(
		t,
		testLogger,
		in,
		priceProvider,
		mustSetCapabilitiesFn,
		[]keystonetypes.DONCapabilityWithConfigFactoryFn{
			libcontracts.DefaultCapabilityFactoryFn,
			libcontracts.ChainWriterCapabilityFactory(libc.MustSafeUint64(int64(homeChainID))),
			libcontracts.ChainWriterCapabilityFactory(libc.MustSafeUint64(int64(targetChainID))),
		},
	)

	// Log extra information that might help debugging
	t.Cleanup(func() {
		debugTest(t, testLogger, setupOutput, in)
	})

	waitForFeedUpdate(t, testLogger, priceProvider, setupOutput, 5*time.Minute)
}

// Verifies that workflow can read a given contract and write it to a feed.  Feed is read by fetching the latest bundle.
// Config file to use: environment-one-don-read-contract.toml
func TestCRE_OCR3_ReadBalance_Workflow_SingleDon_MockedPrice(t *testing.T) {
	testLogger := framework.L
	expectedReadAmount := big.NewInt(99)

	// Load and validate test configuration
	in, err := framework.Load[TestConfig](t)
	require.NoError(t, err, "couldn't load test config")
	validateEnvVars(t, in)
	require.Len(t, in.NodeSets, 1, "expected 1 node set in the test config")

	// Assign all capabilities to the single node set
	mustSetCapabilitiesFn := func(input []*ns.Input) []*keystonetypes.CapabilitiesAwareNodeSet {
		return []*keystonetypes.CapabilitiesAwareNodeSet{
			{
				Input:              input[0],
				Capabilities:       append(SinglePoRDonCapabilitiesFlags, "read-contract"),
				DONTypes:           []string{keystonetypes.WorkflowDON, keystonetypes.GatewayDON},
				BootstrapNodeIndex: 0, // not required, but set to make the configuration explicit
				GatewayNodeIndex:   0, // not required, but set to make the configuration explicit
			},
		}
	}

	// fake price provider without a data provider, price will be read on chain
	priceProvider := &FakePriceProvider{
		testLogger: testLogger,
		expectedPrices: map[string][]*big.Int{
			cleanFeedID(in.WorkflowConfigs[0].FeedID): {
				expectedReadAmount,
			},
		},
	}

	chainIDInt, chainErr := strconv.Atoi(in.Blockchains[0].ChainID)
	require.NoError(t, chainErr, "failed to convert chain ID to int")

	setupOutput := (&testHarness{
		readContractInput: &ReadContractTestInput{
			ExpectedFundingAmount: expectedReadAmount,
		},
	}).setupPoRTestEnvironment(
		t,
		testLogger,
		in,
		priceProvider,
		mustSetCapabilitiesFn,
		[]keystonetypes.DONCapabilityWithConfigFactoryFn{
			libcontracts.DefaultCapabilityFactoryFn,
			libcontracts.ChainWriterCapabilityFactory(libc.MustSafeUint64(int64(chainIDInt))),
			libcontracts.ChainReaderCapabilityFactory(libc.MustSafeUint64(int64(chainIDInt)), "evm"),
		},
	)

	// Log extra information that might help debugging
	t.Cleanup(func() {
		debugTest(t, testLogger, setupOutput, in)
	})

	waitForFeedBundleUpdate(t, testLogger, priceProvider, setupOutput, 5*time.Minute)
}

// config file to use: environment-gateway-don.toml
func TestCRE_OCR3_PoR_Workflow_GatewayDon_MockedPrice(t *testing.T) {
	testLogger := framework.L

	// Load and validate test configuration
	in, err := framework.Load[TestConfig](t)
	require.NoError(t, err, "couldn't load test config")
	validateEnvVars(t, in)
	require.Len(t, in.NodeSets, 2, "expected 2 node sets in the test config")

	// Assign all capabilities to the single node set
	mustSetCapabilitiesFn := func(input []*ns.Input) []*keystonetypes.CapabilitiesAwareNodeSet {
		return []*keystonetypes.CapabilitiesAwareNodeSet{
			{
				Input:              input[0],
				Capabilities:       SinglePoRDonCapabilitiesFlags,
				DONTypes:           []string{keystonetypes.WorkflowDON},
				BootstrapNodeIndex: 0,
			},
			{
				Input:              input[1],
				Capabilities:       []string{},
				DONTypes:           []string{keystonetypes.GatewayDON}, // <----- it's crucial to set the correct DON type
				BootstrapNodeIndex: -1,                                 // <----- it's crucial to indicate there's no bootstrap node
				GatewayNodeIndex:   0,
			},
		}
	}

	priceProvider, priceErr := NewFakePriceProvider(testLogger, in.Fake, AuthorizationKey, []string{in.WorkflowConfigs[0].FeedID})
	require.NoError(t, priceErr, "failed to create fake price provider")

	firstBlockchain := in.Blockchains[0]
	chainIDInt, chainErr := strconv.Atoi(firstBlockchain.ChainID)
	require.NoError(t, chainErr, "failed to convert chain ID to int")

	setupOutput := new(testHarness).setupPoRTestEnvironment(t, testLogger, in, priceProvider, mustSetCapabilitiesFn, []keystonetypes.DONCapabilityWithConfigFactoryFn{libcontracts.DefaultCapabilityFactoryFn, libcontracts.ChainWriterCapabilityFactory(libc.MustSafeUint64(int64(chainIDInt)))})

	// Log extra information that might help debugging
	t.Cleanup(func() {
		debugTest(t, testLogger, setupOutput, in)
	})

	waitForFeedUpdate(t, testLogger, priceProvider, setupOutput, 5*time.Minute)
}

// config file to use: environment-capabilities-don.toml
func TestCRE_OCR3_PoR_Workflow_CapabilitiesDons_LivePrice(t *testing.T) {
	testLogger := framework.L

	// Load and validate test configuration
	in, err := framework.Load[TestConfig](t)
	require.NoError(t, err, "couldn't load test config")
	validateEnvVars(t, in)
	require.Len(t, in.NodeSets, 3, "expected 3 node sets in the test config")

	mustSetCapabilitiesFn := func(input []*ns.Input) []*keystonetypes.CapabilitiesAwareNodeSet {
		return []*keystonetypes.CapabilitiesAwareNodeSet{
			{
				Input:              input[0],
				Capabilities:       []string{keystonetypes.OCR3Capability, keystonetypes.CustomComputeCapability, keystonetypes.CronCapability},
				DONTypes:           []string{keystonetypes.WorkflowDON},
				BootstrapNodeIndex: 0,
			},
			{
				Input:              input[1],
				Capabilities:       []string{keystonetypes.WriteEVMCapability},
				DONTypes:           []string{keystonetypes.CapabilitiesDON}, // <----- it's crucial to set the correct DON type
				BootstrapNodeIndex: -1,                                      // <----- indicate that capabilities DON doesn't have a bootstrap node and will use the global bootstrap node
			},
			{
				Input:              input[2],
				Capabilities:       []string{},
				DONTypes:           []string{keystonetypes.GatewayDON}, // <----- it's crucial to set the correct DON type
				BootstrapNodeIndex: -1,                                 // <----- it's crucial to indicate there's no bootstrap node for the gateway DON
				GatewayNodeIndex:   0,
			},
		}
	}

	firstBlockchain := in.Blockchains[0]
	chainIDInt, chainErr := strconv.Atoi(firstBlockchain.ChainID)
	require.NoError(t, chainErr, "failed to convert chain ID to int")

	priceProvider := NewTrueUSDPriceProvider(testLogger, []string{in.WorkflowConfigs[0].FeedID})
	setupOutput := new(testHarness).setupPoRTestEnvironment(t, testLogger, in, priceProvider, mustSetCapabilitiesFn, []keystonetypes.DONCapabilityWithConfigFactoryFn{libcontracts.DefaultCapabilityFactoryFn, libcontracts.ChainWriterCapabilityFactory(libc.MustSafeUint64(int64(chainIDInt)))})

	// Log extra information that might help debugging
	t.Cleanup(func() {
		debugTest(t, testLogger, setupOutput, in)
	})

	waitForFeedUpdate(t, testLogger, priceProvider, setupOutput, 5*time.Minute)
}

func waitForFeedUpdate(t *testing.T, testLogger zerolog.Logger, priceProvider PriceProvider, setupOutput *porSetupOutput, timeout time.Duration) {
	for chainSelector, workflowConfig := range setupOutput.chainSelectorToWorkflowConfig {
		testLogger.Info().Msgf("Waiting for feed %s to update...", workflowConfig.FeedID)
		timeout := 5 * time.Minute // It can take a while before the first report is produced, particularly on CI.

		dataFeedsCacheAddresses, dataFeedsCacheErr := crecontracts.FindAddressesForChain(setupOutput.addressBook, chainSelector, df_changeset.DataFeedsCache.String())
		require.NoError(t, dataFeedsCacheErr, "failed to find data feeds cache address for chain %d", chainSelector)

		dataFeedsCacheInstance, instanceErr := data_feeds_cache.NewDataFeedsCache(dataFeedsCacheAddresses, setupOutput.chainSelectorToSethClient[chainSelector].Client)
		require.NoError(t, instanceErr, "failed to create data feeds cache instance")

		startTime := time.Now()
		assert.Eventually(t, func() bool {
			elapsed := time.Since(startTime).Round(time.Second)
			price, err := dataFeedsCacheInstance.GetLatestAnswer(setupOutput.chainSelectorToSethClient[chainSelector].NewCallOpts(), [16]byte(common.Hex2Bytes(workflowConfig.FeedID)))
			require.NoError(t, err, "failed to get price from Data Feeds Cache contract")

			// if there are no more prices to be found, we can stop waiting
			return !setupOutput.priceProvider.NextPrice(workflowConfig.FeedID, price, elapsed)
		}, timeout, 10*time.Second, "feed %s did not update, timeout after: %s", workflowConfig.FeedID, timeout)

		require.EqualValues(t, priceProvider.ExpectedPrices(workflowConfig.FeedID), priceProvider.ActualPrices(workflowConfig.FeedID), "prices do not match")
		testLogger.Info().Msgf("All %d prices were found in the feed %s", len(priceProvider.ExpectedPrices(workflowConfig.FeedID)), workflowConfig.FeedID)
	}
}

func waitForFeedBundleUpdate(t *testing.T, testLogger zerolog.Logger, priceProvider PriceProvider, setupOutput *porSetupOutput, timeout time.Duration) {
	for chainSelector, wfConfig := range setupOutput.chainSelectorToWorkflowConfig {
		testLogger.Info().Msgf("Waiting for feed %s to update...", wfConfig.FeedID)

		dataFeedsCacheAddresses, dataFeedsCacheErr := crecontracts.FindAddressesForChain(setupOutput.addressBook, chainSelector, df_changeset.DataFeedsCache.String())
		require.NoError(t, dataFeedsCacheErr, "failed to find data feeds cache address for chain %d", chainSelector)

		dataFeedsCacheInstance, instanceErr := data_feeds_cache.NewDataFeedsCache(dataFeedsCacheAddresses, setupOutput.chainSelectorToSethClient[chainSelector].Client)
		require.NoError(t, instanceErr, "failed to create data feeds cache instance")

		startTime := time.Now()
		assert.Eventually(t, func() bool {
			elapsed := time.Since(startTime).Round(time.Second)
			bundle, err := dataFeedsCacheInstance.GetLatestBundle(setupOutput.chainSelectorToSethClient[chainSelector].NewCallOpts(), [16]byte(common.Hex2Bytes(wfConfig.FeedID)))
			require.NoError(t, err, "failed to get price from Data Feeds Cache contract")
			price := new(big.Int).SetBytes(bundle)
			// if there are no more prices to be found, we can stop waiting
			return !setupOutput.priceProvider.NextPrice(wfConfig.FeedID, price, elapsed)
		}, timeout, 10*time.Second, "feed %s did not update, timeout after: %s", wfConfig.FeedID, timeout)

		require.EqualValues(t, priceProvider.ExpectedPrices(wfConfig.FeedID), priceProvider.ActualPrices(wfConfig.FeedID), "prices do not match")
		testLogger.Info().Msgf("All %d prices were found in the feed %s", len(priceProvider.ExpectedPrices(wfConfig.FeedID)), wfConfig.FeedID)
	}
}

func debugTest(t *testing.T, testLogger zerolog.Logger, setupOutput *porSetupOutput, in *TestConfig) {
	if t.Failed() {
		counter := 0
		for chainSelector, workflowConfig := range setupOutput.chainSelectorToWorkflowConfig {
			dataFeedsCacheAddresses, dataFeedsCacheErr := crecontracts.FindAddressesForChain(setupOutput.addressBook, chainSelector, df_changeset.DataFeedsCache.String())
			require.NoError(t, dataFeedsCacheErr, "failed to find data feeds cache address for chain %d", chainSelector)

			forwarderAddresses, forwarderErr := crecontracts.FindAddressesForChain(setupOutput.addressBook, chainSelector, keystone_changeset.KeystoneForwarder.String())
			require.NoError(t, forwarderErr, "failed to find forwarder address for chain %d", chainSelector)

			logTestInfo(testLogger, workflowConfig.FeedID, workflowConfig.WorkflowName, dataFeedsCacheAddresses.Hex(), forwarderAddresses.Hex())
			counter++
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
				BlockchainOutput: setupOutput.chainSelectorToBlockchainOutput[chainSelector],
				InfraInput:       in.Infra,
			}
			lidebug.PrintTestDebug(t.Name(), testLogger, debugInput)
		}
	}
}
