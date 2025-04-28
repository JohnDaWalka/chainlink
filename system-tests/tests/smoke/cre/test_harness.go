package cre

import (
	"math/big"
	"os"

	"github.com/pkg/errors"

	df_changeset "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink/deployment"

	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	libcontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	keystonetypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	creworkflow "github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"
	libcrecli "github.com/smartcontractkit/chainlink/system-tests/lib/crecli"
	keystoneporcrecli "github.com/smartcontractkit/chainlink/system-tests/lib/crecli/por"
)

const (
	CronBinaryVersion   = "v1.0.2-alpha"
	CRECLIBinaryVersion = "v0.1.5"

	AuthorizationKeySecretName = "AUTH_KEY"
	// TODO: use once we can run these tests in CI (https://smartcontract-it.atlassian.net/browse/DX-589)
	// AuthorizationKey           = "12a-281j&@91.sj1:_}"
	AuthorizationKey = ""
)

type ReadContractTestInput struct {
	ExpectedFundingAmount *big.Int
}

type testHarness struct {
	readContractInput *ReadContractTestInput
}

// Defines the location of already compiled workflow binary and config files
// They will be used if WorkflowConfig.ShouldCompileNewWorkflow is `false`
// Otherwise test will compile and upload a new workflow
type CompiledConfig struct {
	BinaryURL  string `toml:"binary_url" validate:"required"`
	ConfigURL  string `toml:"config_url" validate:"required"`
	SecretsURL string `toml:"secrets_url"`
}

type WorkflowConfig struct {
	UseCRECLI bool `toml:"use_cre_cli"`
	/*
		These tests can be run in two modes:
		1. existing mode: it uses a workflow binary (and configuration) file that is already uploaded to Gist
		2. compile mode: it compiles a new workflow binary and uploads it to Gist

		For the "compile" mode to work, the `GIST_WRITE_TOKEN` env var must be set to a token that has `gist:read` and `gist:write` permissions, but this permissions
		are tied to account not to repository. Currently, we have no service account in the CI at all. And using a token that's tied to personal account of a developer
		is not a good idea. So, for now, we are only allowing the `existing` mode in CI.

		If you wish to use "compile" mode set `ShouldCompileNewWorkflow` to `true`, set `GIST_WRITE_TOKEN` env var and provide the path to the workflow folder.
	*/
	ShouldCompileNewWorkflow bool `toml:"should_compile_new_workflow" validate:"no_cre_no_compilation,disabled_in_ci"`
	// Tells the test where the workflow to compile is located
	WorkflowFolderLocation *string         `toml:"workflow_folder_location" validate:"required_if=ShouldCompileNewWorkflow true"`
	CompiledWorkflowConfig *CompiledConfig `toml:"compiled_config" validate:"required_if=ShouldCompileNewWorkflow false"`
	WorkflowName           string          `toml:"workflow_name" validate:"required" `
	FeedID                 string          `toml:"feed_id" validate:"required,startsnotwith=0x"`
}

type readContractInput struct {
	readTargetName       string
	contractReaderConfig string
	fundedAddress        string
	contractName         string
	contractMethod       string
}

type registerPoRWorkflowInput struct {
	WorkflowConfig
	chainSelector      uint64
	homeChainSelector  uint64
	writeTargetName    string
	workflowDonID      uint32
	feedID             string
	addressBook        deployment.AddressBook
	priceProvider      PriceProvider
	sethClient         *seth.Client
	deployerPrivateKey string
	creCLIAbsPath      string
	creCLIsettingsFile *os.File
	authKey            string
	readContractInput  *readContractInput
}

type workflowRegistrar struct {
	configBuilder ConfigBuilderFunc
}

type ConfigBuilderFunc func(registerPoRWorkflowInput) (*os.File, error)

// create PoR workflow-specific config file
func buildPoRConfig(input registerPoRWorkflowInput) (*os.File, error) {
	var secretNameToUse *string
	if input.authKey != "" {
		secretNameToUse = ptr.Ptr(AuthorizationKeySecretName)
	}

	dataFeedsCacheAddress, dataFeedsCacheErr := crecontracts.FindAddressesForChain(input.addressBook, input.chainSelector, df_changeset.DataFeedsCache.String())
	if dataFeedsCacheErr != nil {
		return nil, errors.Wrapf(dataFeedsCacheErr, "failed to find data feeds cache address for chain %d", input.chainSelector)
	}

	workflowConfigFile, configErr := keystoneporcrecli.CreateConfigFile(
		&libcrecli.PoRWorkflowConfig{
			FeedID:            input.FeedID,
			URL:               input.priceProvider.URL(),
			WriteTargetName:   input.writeTargetName,
			ConsumerAddress:   dataFeedsCacheAddress.Hex(),
			AuthKeySecretName: secretNameToUse,
		},
	)
	if configErr != nil {
		return nil, errors.Wrap(configErr, "failed to create workflow config file")
	}
	return workflowConfigFile, nil
}

func buildReadContractConfig(input registerPoRWorkflowInput) (*os.File, error) {
	if input.readContractInput == nil {
		return nil, errors.New("cannot build read contract config from nil input")
	}

	var secretNameToUse *string
	if input.authKey != "" {
		secretNameToUse = ptr.Ptr(AuthorizationKeySecretName)
	}

	dataFeedsCacheAddress, dataFeedsCacheErr := crecontracts.FindAddressesForChain(input.addressBook, input.chainSelector, df_changeset.DataFeedsCache.String())
	if dataFeedsCacheErr != nil {
		return nil, errors.Wrapf(dataFeedsCacheErr, "failed to find data feeds cache address for chain %d", input.chainSelector)
	}

	balanceReaderAddress, err := crecontracts.FindAddressesForChain(input.addressBook, input.chainSelector, "BalanceReader")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find balance reader address for chain %d", input.chainSelector)
	}

	workflowConfigFile, configErr := keystoneporcrecli.CreateConfigFile(
		&libcrecli.PoRWorkflowConfig{
			FeedID:            input.FeedID,
			WriteTargetName:   input.writeTargetName,
			ConsumerAddress:   dataFeedsCacheAddress.Hex(),
			AuthKeySecretName: secretNameToUse,
			ReadBalanceReaderConfig: libcrecli.ReadBalanceReaderConfig{
				ReadTargetName:       input.readContractInput.readTargetName,
				ContractReaderConfig: input.readContractInput.contractReaderConfig,
				FundedAddress:        input.readContractInput.fundedAddress,
				ContractAddress:      balanceReaderAddress.Hex(),
				ContractName:         input.readContractInput.contractName,
				ContractMethod:       input.readContractInput.contractMethod,
			},
		},
	)
	if configErr != nil {
		return nil, errors.Wrap(configErr, "failed to create workflow config file")
	}
	return workflowConfigFile, nil
}

func (wr *workflowRegistrar) registerWorkflow(input registerPoRWorkflowInput) error {
	// Register workflow directly using the provided binary URL and optionally config and secrets URLs
	// This is a legacy solution, probably we can remove it soon, but there's still quite a lot of people
	// who have no access to dev-platform repo, so they cannot use the CRE CLI
	if !input.WorkflowConfig.ShouldCompileNewWorkflow && !input.WorkflowConfig.UseCRECLI {
		workflowRegistryAddress, workflowRegistryErr := crecontracts.FindAddressesForChain(input.addressBook, input.chainSelector, keystone_changeset.WorkflowRegistry.String())
		if workflowRegistryErr != nil {
			return errors.Wrapf(workflowRegistryErr, "failed to find workflow registry address for chain %d", input.chainSelector)
		}

		err := libcontracts.RegisterWorkflow(
			input.sethClient,
			workflowRegistryAddress,
			input.workflowDonID,
			input.WorkflowConfig.WorkflowName,
			input.WorkflowConfig.CompiledWorkflowConfig.BinaryURL,
			&input.WorkflowConfig.CompiledWorkflowConfig.ConfigURL,
			&input.WorkflowConfig.CompiledWorkflowConfig.SecretsURL,
		)
		if err != nil {
			return errors.Wrap(err, "failed to register workflow")
		}

		return nil
	}

	// This env var is required by the CRE CLI
	err := os.Setenv("CRE_ETH_PRIVATE_KEY", input.deployerPrivateKey)
	if err != nil {
		return errors.Wrap(err, "failed to set CRE_ETH_PRIVATE_KEY")
	}

	workflowConfigFile, configErr := wr.configBuilder(input)
	if configErr != nil {
		return errors.Wrap(configErr, "failed to create workflow config file")
	}
	workflowConfigFilePath := workflowConfigFile.Name()

	// indicate to the CRE CLI that the secret will be shared between all nodes in the workflow by using specific suffix
	authKeyEnvVarName := AuthorizationKeySecretName + libcrecli.SharedSecretEnvVarSuffix

	var secretsFilePath *string
	if input.authKey != "" {
		// create workflow-specific secrets file using the CRE CLI, which contains a mapping of secret names to environment variables that hold them
		// secrets will be read from the environment variables by the CRE CLI and encoded using nodes' public keys and when workflow executes it will
		// be able to read all secrets, which after decoding will be set as environment variables with names specified in the secrets file
		secrets := map[string][]string{
			AuthorizationKeySecretName: {authKeyEnvVarName},
		}

		secretsFile, secretsErr := libcrecli.CreateSecretsFile(secrets)
		if secretsErr != nil {
			return errors.Wrap(secretsErr, "failed to create secrets file")
		}
		secretsFilePath = ptr.Ptr(secretsFile.Name())
	}

	workflowRegistryAddress, workflowRegistryErr := crecontracts.FindAddressesForChain(input.addressBook, input.homeChainSelector, keystone_changeset.WorkflowRegistry.String())
	if workflowRegistryErr != nil {
		return errors.Wrapf(workflowRegistryErr, "failed to find workflow registry address for chain %d", input.homeChainSelector)
	}

	registerWorkflowInput := keystonetypes.RegisterWorkflowWithCRECLIInput{
		ChainSelector:            input.chainSelector,
		WorkflowDonID:            input.workflowDonID,
		WorkflowRegistryAddress:  workflowRegistryAddress,
		WorkflowOwnerAddress:     input.sethClient.MustGetRootKeyAddress(),
		CRECLIPrivateKey:         input.deployerPrivateKey,
		CRECLIAbsPath:            input.creCLIAbsPath,
		CRESettingsFile:          input.creCLIsettingsFile,
		WorkflowName:             input.WorkflowConfig.WorkflowName,
		ShouldCompileNewWorkflow: input.WorkflowConfig.ShouldCompileNewWorkflow,
	}

	if input.WorkflowConfig.ShouldCompileNewWorkflow {
		registerWorkflowInput.NewWorkflow = &keystonetypes.NewWorkflow{
			FolderLocation:   *input.WorkflowConfig.WorkflowFolderLocation,
			WorkflowFileName: "main.go",
			ConfigFilePath:   &workflowConfigFilePath,
			SecretsFilePath:  secretsFilePath,
			Secrets: map[string]string{
				authKeyEnvVarName: input.authKey,
			},
		}
	} else {
		registerWorkflowInput.ExistingWorkflow = &keystonetypes.ExistingWorkflow{
			BinaryURL:  input.WorkflowConfig.CompiledWorkflowConfig.BinaryURL,
			ConfigURL:  &input.WorkflowConfig.CompiledWorkflowConfig.ConfigURL,
			SecretsURL: &input.WorkflowConfig.CompiledWorkflowConfig.SecretsURL,
		}
	}

	registerErr := creworkflow.RegisterWithCRECLI(registerWorkflowInput)
	if registerErr != nil {
		return errors.Wrap(registerErr, "failed to register workflow with CRE CLI")
	}

	return nil
}
