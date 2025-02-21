package capabilities_test

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"net"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/feeds_consumer"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	ctfconfig "github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"

	keystonecapabilities "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	libcontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	lidebug "github.com/smartcontractkit/chainlink/system-tests/lib/cre/debug"
	libdon "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	keystoneporconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/config/por"
	libjobs "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	keystonepor "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/por"
	libnode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	libenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	keystonetypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	libcrecli "github.com/smartcontractkit/chainlink/system-tests/lib/crecli"
	keystoneporcrecli "github.com/smartcontractkit/chainlink/system-tests/lib/crecli/por"
	libfunding "github.com/smartcontractkit/chainlink/system-tests/lib/funding"
	libtypes "github.com/smartcontractkit/chainlink/system-tests/lib/types"
)

const (
	cronCapabilityAssetFile            = "amd64_cron"
	ghReadTokenEnvVarName              = "GITHUB_READ_TOKEN"
	E2eJobDistributorImageEnvVarName   = "E2E_JD_IMAGE"
	E2eJobDistributorVersionEnvVarName = "E2E_JD_VERSION"
)

type TestConfig struct {
	BlockchainA                   *blockchain.Input                      `toml:"blockchain_a" validate:"required"`
	NodeSets                      []*ns.Input                            `toml:"nodesets" validate:"required"`
	WorkflowConfig                *WorkflowConfig                        `toml:"workflow_config" validate:"required"`
	JD                            *jd.Input                              `toml:"jd" validate:"required"`
	Fake                          *fake.Input                            `toml:"fake"`
	KeystoneContracts             *keystonetypes.KeystoneContractsInput  `toml:"keystone_contracts"`
	WorkflowRegistryConfiguration *keystonetypes.WorkflowRegistryInput   `toml:"workflow_registry_configuration"`
	FeedConsumer                  *keystonetypes.DeployFeedConsumerInput `toml:"feed_consumer"`
}

type WorkflowConfig struct {
	UseCRECLI                bool `toml:"use_cre_cli"`
	ShouldCompileNewWorkflow bool `toml:"should_compile_new_workflow" validate:"no_cre_no_compilation,disabled_in_ci"`
	// Tells the test where the workflow to compile is located
	WorkflowFolderLocation *string             `toml:"workflow_folder_location" validate:"required_if=ShouldCompileNewWorkflow true"`
	CompiledWorkflowConfig *CompiledConfig     `toml:"compiled_config" validate:"required_if=ShouldCompileNewWorkflow false"`
	DependenciesConfig     *DependenciesConfig `toml:"dependencies" validate:"required"`
	WorkflowName           string              `toml:"workflow_name" validate:"required" `
	FeedID                 string              `toml:"feed_id" validate:"required,startsnotwith=0x"`
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

// Defines relases/versions of test dependencies that will be downloaded from Github
type DependenciesConfig struct {
	CapabiltiesVersion string `toml:"capabilities_version" validate:"required"`
	CRECLIVersion      string `toml:"cre_cli_version" validate:"required"`
}

// Defines the location of already compiled workflow binary and config files
// They will be used if WorkflowConfig.ShouldCompileNewWorkflow is `false`
// Otherwise test will compile and upload a new workflow
type CompiledConfig struct {
	BinaryURL string `toml:"binary_url" validate:"required"`
	ConfigURL string `toml:"config_url" validate:"required"`
}

func validateEnvVars(t *testing.T, in *TestConfig) {
	require.NotEmpty(t, os.Getenv("PRIVATE_KEY"), "PRIVATE_KEY env var must be set")

	var ghReadToken string
	// this is a small hack to avoid changing the reusable workflow
	if os.Getenv("CI") == "true" {
		// This part should ideally happen outside of the test, but due to how our reusable e2e test workflow is structured now
		// we cannot execute this part in workflow steps (it doesn't support any pre-execution hooks)
		require.NotEmpty(t, os.Getenv(ctfconfig.E2E_TEST_CHAINLINK_IMAGE_ENV), "missing env var: "+ctfconfig.E2E_TEST_CHAINLINK_IMAGE_ENV)
		require.NotEmpty(t, os.Getenv(ctfconfig.E2E_TEST_CHAINLINK_VERSION_ENV), "missing env var: "+ctfconfig.E2E_TEST_CHAINLINK_VERSION_ENV)
		require.NotEmpty(t, os.Getenv(libjobs.E2eJobDistributorImageEnvVarName), "missing env var: "+libjobs.E2eJobDistributorImageEnvVarName)
		require.NotEmpty(t, os.Getenv(libjobs.E2eJobDistributorVersionEnvVarName), "missing env var: "+libjobs.E2eJobDistributorVersionEnvVarName)

		// disabled until we can figure out how to generate a gist read:write token in CI
		/*
		 This test can be run in two modes:
		 1. `existing` mode: it uses a workflow binary (and configuration) file that is already uploaded to Gist
		 2. `compile` mode: it compiles a new workflow binary and uploads it to Gist

		 For the `new` mode to work, the `GITHUB_API_TOKEN` env var must be set to a token that has `gist:read` and `gist:write` permissions, but this permissions
		 are tied to account not to repository. Currently, we have no service account in the CI at all. And using a token that's tied to personal account of a developer
		 is not a good idea. So, for now, we are only allowing the `existing` mode in CI.
		*/

		// we use this special function to subsitute a placeholder env variable with the actual environment variable name
		// it is defined in .github/e2e-tests.yml as '{{ env.GITHUB_API_TOKEN }}'
		ghReadToken = ctfconfig.MustReadEnvVar_String(ghReadTokenEnvVarName)
	} else {
		ghReadToken = os.Getenv(ghReadTokenEnvVarName)
	}

	require.NotEmpty(t, ghReadToken, ghReadTokenEnvVarName+" env var must be set")

	if in.WorkflowConfig.UseCRECLI {
		if in.WorkflowConfig.ShouldCompileNewWorkflow {
			gistWriteToken := os.Getenv("GIST_WRITE_TOKEN")
			require.NotEmpty(t, gistWriteToken, "GIST_WRITE_TOKEN must be set to use CRE CLI to compile workflows. It requires gist:read and gist:write permissions")
			err := os.Setenv("GITHUB_API_TOKEN", gistWriteToken)
			require.NoError(t, err, "failed to set GITHUB_API_TOKEN env var")
		}
	}
}

type binaryDownloadOutput struct {
	creCLIAbsPath string
}

// this is a small hack to avoid changing the reusable workflow, which doesn't allow to run any pre-execution hooks
func downloadBinaryFiles(in *TestConfig) (*binaryDownloadOutput, error) {
	var ghReadToken string
	if os.Getenv("CI") == "true" {
		ghReadToken = ctfconfig.MustReadEnvVar_String(ghReadTokenEnvVarName)
	} else {
		ghReadToken = os.Getenv(ghReadTokenEnvVarName)
	}

	_, err := keystonecapabilities.DownloadCapabilityFromRelease(ghReadToken, in.WorkflowConfig.DependenciesConfig.CapabiltiesVersion, cronCapabilityAssetFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to download cron capability. Make sure token has content:read permissions to the capabilities repo")
	}

	output := &binaryDownloadOutput{}

	if in.WorkflowConfig.UseCRECLI {
		output.creCLIAbsPath, err = libcrecli.DownloadAndInstallChainlinkCLI(ghReadToken, in.WorkflowConfig.DependenciesConfig.CRECLIVersion)
		if err != nil {
			return nil, errors.Wrap(err, "failed to download and install CRE CLI. Make sure token has content:read permissions to the dev-platform repo")
		}
	}

	return output, nil
}

type registerPoRWorkflowInput struct {
	*WorkflowConfig
	chainSelector               uint64
	workflowDonID               uint32
	feedID                      string
	workflowRegistryAddress     common.Address
	feedConsumerAddress         common.Address
	capabilitiesRegistryAddress common.Address
	priceProvider               PriceProvider
	sethClient                  *seth.Client
	deployerPrivateKey          string
	blockchain                  *blockchain.Output
	binaryDownloadOutput        binaryDownloadOutput
}

func registerPoRWorkflow(input registerPoRWorkflowInput) error {
	// Register workflow directly using the provided binary and config URLs
	// This is a legacy solution, probably we can remove it soon, but there's still quite a lot of people
	// who have no access to dev-platform repo, so they cannot use the CRE CLI
	if !input.WorkflowConfig.ShouldCompileNewWorkflow && !input.WorkflowConfig.UseCRECLI {
		err := libcontracts.RegisterWorkflow(input.sethClient, input.workflowRegistryAddress, input.workflowDonID, input.WorkflowConfig.WorkflowName, input.WorkflowConfig.CompiledWorkflowConfig.BinaryURL, input.WorkflowConfig.CompiledWorkflowConfig.ConfigURL)
		if err != nil {
			return errors.Wrap(err, "failed to register workflow")
		}

		return nil
	}

	// These two env vars are required by the CRE CLI
	err := os.Setenv("WORKFLOW_OWNER_ADDRESS", input.sethClient.MustGetRootKeyAddress().Hex())
	if err != nil {
		return errors.Wrap(err, "failed to set WORKFLOW_OWNER_ADDRESS env var")
	}

	err = os.Setenv("ETH_PRIVATE_KEY", input.deployerPrivateKey)
	if err != nil {
		return errors.Wrap(err, "failed to set ETH_PRIVATE_KEY")
	}

	// create CRE CLI settings file
	settingsFile, settingsErr := libcrecli.PrepareCRECLISettingsFile(input.sethClient.MustGetRootKeyAddress(), input.capabilitiesRegistryAddress, input.workflowRegistryAddress, input.workflowDonID, input.chainSelector, input.blockchain.Nodes[0].HostHTTPUrl)
	if settingsErr != nil {
		return errors.Wrap(settingsErr, "failed to create CRE CLI settings file")
	}

	var workflowURL string
	var workflowConfigURL string

	workflowConfigFile, configErr := keystoneporcrecli.CreateConfigFile(input.feedConsumerAddress, input.feedID, input.priceProvider.URL())
	if configErr != nil {
		return errors.Wrap(configErr, "failed to create workflow config file")
	}

	// compile and upload the workflow, if we are not using an existing one
	if input.WorkflowConfig.ShouldCompileNewWorkflow {
		compilationResult, err := libcrecli.CompileWorkflow(input.binaryDownloadOutput.creCLIAbsPath, *input.WorkflowConfig.WorkflowFolderLocation, workflowConfigFile, settingsFile)
		if err != nil {
			return errors.Wrap(err, "failed to compile workflow")
		}

		workflowURL = compilationResult.WorkflowURL
		workflowConfigURL = compilationResult.ConfigURL
	} else {
		workflowURL = input.WorkflowConfig.CompiledWorkflowConfig.BinaryURL
		workflowConfigURL = input.WorkflowConfig.CompiledWorkflowConfig.ConfigURL
	}

	registerErr := libcrecli.RegisterWorkflow(input.binaryDownloadOutput.creCLIAbsPath, input.WorkflowName, workflowURL, workflowConfigURL, settingsFile)
	if registerErr != nil {
		return errors.Wrap(registerErr, "failed to register workflow")
	}

	return nil
}

func logTestInfo(l zerolog.Logger, feedID, workflowName, feedConsumerAddr, forwarderAddr string) {
	l.Info().Msg("------ Test configuration:")
	l.Info().Msgf("Feed ID: %s", feedID)
	l.Info().Msgf("Workflow name: %s", workflowName)
	l.Info().Msgf("FeedConsumer address: %s", feedConsumerAddr)
	l.Info().Msgf("KeystoneForwarder address: %s", forwarderAddr)
}

func extraAllowedPortsAndIps(testLogger zerolog.Logger, fakePort int) ([]string, []int, error) {
	// we need to explicitly allow the port used by the fake data provider
	// and IP corresponding to host.docker.internal or the IP of the host machine, if we are running on Linux,
	// because that's where the fake data provider is running
	var hostIP string
	var err error

	system := runtime.GOOS
	switch system {
	case "darwin":
		hostIP = "192.168.65.1"
	case "linux":
		// for linux framework already returns an IP, so we don't need to resolve it,
		// but we need to remove the http:// prefix
		hostIP = strings.ReplaceAll(framework.HostDockerInternal(), "http://", "")
	default:
		err = fmt.Errorf("unsupported OS: %s", system)
	}
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to resolve host.docker.internal IP")
	}

	testLogger.Info().Msgf("Will allow IP %s and port %d for the fake data provider", hostIP, fakePort)

	ips, err := net.LookupIP("gist.githubusercontent.com")
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to resolve IP for gist.githubusercontent.com")
	}

	gistIPs := make([]string, len(ips))
	for i, ip := range ips {
		gistIPs[i] = ip.To4().String()
		testLogger.Debug().Msgf("Resolved IP for gist.githubusercontent.com: %s", gistIPs[i])
	}

	// we also need to explicitly allow Gist's IP
	return append(gistIPs, hostIP), []int{fakePort}, nil
}

// func extraAllowedPortsAndIps(testLogger zerolog.Logger, fakePort int, nodeOutput *ns.Output) ([]string, []int, error) {
// 	// we need to explicitly allow the port used by the fake data provider
// 	// and IP corresponding to host.docker.internal or the IP of the host machine, if we are running on Linux,
// 	// because that's where the fake data provider is running
// 	var hostIP string
// 	var err error

// 	system := runtime.GOOS
// 	switch system {
// 	case "darwin":
// 		hostIP, err = libdon.ResolveHostDockerInternaIP(testLogger, nodeOutput)
// 	case "linux":
// 		// for linux framework already returns an IP, so we don't need to resolve it,
// 		// but we need to remove the http:// prefix
// 		hostIP = strings.ReplaceAll(framework.HostDockerInternal(), "http://", "")
// 	default:
// 		err = fmt.Errorf("unsupported OS: %s", system)
// 	}
// 	if err != nil {
// 		return nil, nil, errors.Wrap(err, "failed to resolve host.docker.internal IP")
// 	}

// 	testLogger.Info().Msgf("Will allow IP %s and port %d for the fake data provider", hostIP, fakePort)

// 	ips, err := net.LookupIP("gist.githubusercontent.com")
// 	if err != nil {
// 		return nil, nil, errors.Wrap(err, "failed to resolve IP for gist.githubusercontent.com")
// 	}

// 	gistIPs := make([]string, len(ips))
// 	for i, ip := range ips {
// 		gistIPs[i] = ip.To4().String()
// 		testLogger.Debug().Msgf("Resolved IP for gist.githubusercontent.com: %s", gistIPs[i])
// 	}

// 	// we also need to explicitly allow Gist's IP
// 	return append(gistIPs, hostIP), []int{fakePort}, nil
// }

type InfrastructureInput struct {
	jdInput         *jd.Input
	nodeSetInput    []*keystonetypes.CapabilitiesAwareNodeSet
	blockchainInput *blockchain.Input
}

type InfrastructureOutput struct {
	chainSelector      uint64
	blockchainOutput   *blockchain.Output
	jdOutput           *jd.Output
	sethClient         *seth.Client
	deployerPrivateKey string
	gatewayConnector   *keystonetypes.GatewayConnectorOutput
}

func CreateInfrastructure(
	cldLogger logger.Logger,
	testLogger zerolog.Logger,
	input InfrastructureInput,
) (*InfrastructureOutput, error) {
	if input.blockchainInput == nil {
		return nil, errors.New("blockchain input is nil")
	}

	if input.jdInput == nil {
		return nil, errors.New("JD input is nil")
	}

	if len(input.nodeSetInput) == 0 {
		return nil, errors.New("node set input is empty")
	}

	// Create a new blockchain network and Seth client to interact with it
	blockchainOutput, err := blockchain.NewBlockchainNetwork(input.blockchainInput)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create blockchain network")
	}

	pkey := os.Getenv("PRIVATE_KEY")
	if pkey == "" {
		return nil, errors.New("PRIVATE_KEY env var must be set")
	}

	sethClient, err := seth.NewClientBuilder().
		WithRpcUrl(blockchainOutput.Nodes[0].HostWSUrl).
		WithPrivateKeys([]string{pkey}).
		// do not check if there's a pending nonce nor check node health
		WithProtections(false, false, seth.MustMakeDuration(time.Second)).
		Build()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create seth client")
	}

	chainSelector, err := chainselectors.SelectorFromChainId(sethClient.Cfg.Network.ChainID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get chain selector for chain id %d", sethClient.Cfg.Network.ChainID)
	}

	// Start job distributor
	if os.Getenv("CI") == "true" {
		jdImage := ctfconfig.MustReadEnvVar_String(E2eJobDistributorImageEnvVarName)
		jdVersion := os.Getenv(E2eJobDistributorVersionEnvVarName)
		input.jdInput.Image = fmt.Sprintf("%s:%s", jdImage, jdVersion)
	}

	jdOutput, err := jd.NewJD(input.jdInput)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new job distributor")
	}

	// Deploy the DONs
	// Hack for CI that allows us to dynamically set the chainlink image and version
	// CTFv2 currently doesn't support dynamic image and version setting
	if os.Getenv("CI") == "true" {
		// Due to how we pass custom env vars to reusable workflow we need to use placeholders, so first we need to resolve what's the name of the target environment variable
		// that stores chainlink version and then we can use it to resolve the image name
		for i := range input.nodeSetInput {
			image := fmt.Sprintf("%s:%s", os.Getenv(ctfconfig.E2E_TEST_CHAINLINK_IMAGE_ENV), ctfconfig.MustReadEnvVar_String(ctfconfig.E2E_TEST_CHAINLINK_VERSION_ENV))
			for j := range input.nodeSetInput[i].NodeSpecs {
				input.nodeSetInput[i].NodeSpecs[j].Node.Image = image
			}
		}
	}

	return &InfrastructureOutput{
		chainSelector: chainSelector,
		// nodeOuput:          nodeOutput,
		blockchainOutput: blockchainOutput,
		jdOutput:         jdOutput,
		// cldEnv:             cldEnv,
		// donTopology:        donTopology,
		sethClient:         sethClient,
		deployerPrivateKey: pkey,
		gatewayConnector: &keystonetypes.GatewayConnectorOutput{
			Path: "/node",
			Port: 5003,
			// do not set the host, it will be resolved automatically
		},
	}, nil
}

type setupOutput struct {
	priceProvider        PriceProvider
	feedsConsumerAddress common.Address
	forwarderAddress     common.Address
	sethClient           *seth.Client
	blockchainOutput     *blockchain.Output
	donTopology          *keystonetypes.DonTopology
	nodeOutput           []*keystonetypes.WrappedNodeOutput
}

type NodeEthKeySelector struct {
	ChainSelector uint64 `toml:"ChainSelector"`
}
type NodeEthKey struct {
	JSON     string             `toml:"JSON"`
	Password string             `toml:"Password"`
	Selector NodeEthKeySelector `toml:"Selector"`
}

type NodeP2PKey struct {
	PrivateKey string `toml:"PrivateKey"`
	PublicKey  string `toml:"PublicKey"`
}

type p2pGenerationResult struct {
	encryptedP2PKeyJSONs [][]byte
	peerIDs              []string
	publicHexKeys        []string
	privateKeys          []string
}

func generateP2PKeys(pwd string, n int) (*p2pGenerationResult, error) {
	result := &p2pGenerationResult{}
	for i := 0; i < n; i++ {
		key, err := p2pkey.NewV2()
		if err != nil {
			return nil, err
		}
		d, err := key.ToEncryptedJSON(pwd, utils.DefaultScryptParams)
		if err != nil {
			return nil, err
		}
		result.encryptedP2PKeyJSONs = append(result.encryptedP2PKeyJSONs, d)
		result.peerIDs = append(result.peerIDs, key.PeerID().String())
		result.publicHexKeys = append(result.publicHexKeys, key.PublicKeyHex())
		result.privateKeys = append(result.privateKeys, hex.EncodeToString(key.Raw()))
	}
	return result, nil
}

// TODO update in the CTF, I need the public addresses
func NewETHKey(password string) ([]byte, common.Address, error) {
	privateKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	var address common.Address
	if err != nil {
		return nil, address, errors.Wrap(err, "failed to generate private key")
	}
	address = crypto.PubkeyToAddress(privateKey.PublicKey)
	jsonKey, err := keystore.EncryptKey(&keystore.Key{
		PrivateKey: privateKey,
		Address:    address,
	}, password, keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		return nil, address, errors.Wrap(err, "failed to encrypt the keystore")
	}
	return jsonKey, address, nil
}

func generateEVMKeys(pwd string, n int) ([][]byte, []common.Address, error) {
	encryptedEVMKeyJSONs := make([][]byte, 0)
	addresses := make([]common.Address, 0)
	for i := 0; i < n; i++ {
		key, addr, err := NewETHKey(pwd)
		if err != nil {
			return nil, addresses, err
		}
		encryptedEVMKeyJSONs = append(encryptedEVMKeyJSONs, key)
		addresses = append(addresses, addr)
	}
	return encryptedEVMKeyJSONs, addresses, nil
}

func setupTestEnvironment(t *testing.T, testLogger zerolog.Logger, in *TestConfig, priceProvider PriceProvider, binaryDownloadOutput binaryDownloadOutput, mustSetCapabilitiesFn func(input []*ns.Input) []*keystonetypes.CapabilitiesAwareNodeSet) *setupOutput {
	// Universal setup -- START
	envInput := InfrastructureInput{
		jdInput:         in.JD,
		nodeSetInput:    mustSetCapabilitiesFn(in.NodeSets),
		blockchainInput: in.BlockchainA,
	}
	singeFileLogger := cldlogger.NewSingleFileLogger(t)
	envOutput, err := CreateInfrastructure(singeFileLogger, testLogger, envInput)
	require.NoError(t, err, "failed to start environment")

	chainsConfig := []devenv.ChainConfig{
		{
			ChainID:   envOutput.sethClient.Cfg.Network.ChainID,
			ChainName: envOutput.sethClient.Cfg.Network.Name,
			ChainType: strings.ToUpper(envOutput.blockchainOutput.Family),
			WSRPCs: []devenv.CribRPCs{{
				External: envOutput.blockchainOutput.Nodes[0].HostWSUrl,
				Internal: envOutput.blockchainOutput.Nodes[0].DockerInternalWSUrl,
			}},
			HTTPRPCs: []devenv.CribRPCs{{
				External: envOutput.blockchainOutput.Nodes[0].HostHTTPUrl,
				Internal: envOutput.blockchainOutput.Nodes[0].DockerInternalHTTPUrl,
			}},
			DeployerKey: envOutput.sethClient.NewTXOpts(seth.WithNonce(nil)), // set nonce to nil, so that it will be fetched from the chain
		},
	}

	chains, err := devenv.NewChains(singeFileLogger, chainsConfig)
	require.NoError(t, err, "failed to create chains")

	chainsOnlyCld := &deployment.Environment{
		Logger:            singeFileLogger,
		Chains:            chains,
		ExistingAddresses: deployment.NewMemoryAddressBook(),
	}

	// Deploy keystone contracts (forwarder, capability registry, ocr3 capability, workflow registry)
	keystoneContractsInput := &keystonetypes.KeystoneContractsInput{
		ChainSelector: envOutput.chainSelector,
		CldEnv:        chainsOnlyCld,
	}

	keystoneContractsOutput, err := libcontracts.DeployKeystone(testLogger, keystoneContractsInput)
	require.NoError(t, err, "failed to deploy keystone contracts")

	// nodeInputs := mustSetCapabilitiesFn(in.NodeSets)
	topology, err := libdon.BuildTopology(envInput.nodeSetInput)
	require.NoError(t, err, "failed to build input DON topology")

	// Configure Workflow Registry
	workflowRegistryInput := &keystonetypes.WorkflowRegistryInput{
		ChainSelector:  envOutput.chainSelector,
		CldEnv:         chainsOnlyCld,
		AllowedDonIDs:  []uint32{topology.WorkflowDONID},
		WorkflowOwners: []common.Address{envOutput.sethClient.MustGetRootKeyAddress()},
	}

	_, err = libcontracts.ConfigureWorkflowRegistry(testLogger, workflowRegistryInput)
	require.NoError(t, err, "failed to configure workflow registry")

	// Allow extra IPs and ports for the fake data provider, which is running on host machine and requires explicit whitelisting
	var extraAllowedIPs []string
	var extraAllowedPorts []int
	if _, ok := priceProvider.(*FakePriceProvider); ok {
		extraAllowedIPs, extraAllowedPorts, err = extraAllowedPortsAndIps(testLogger, in.Fake.Port) // donTopology.MetaDons[0].NodeOutput.Output)
		require.NoError(t, err, "failed to get extra allowed ports and IPs")
	}

	// Generate keys
	donToP2PKeys := make(map[uint32]*p2pGenerationResult)
	donToEthKeys := make(map[uint32][][]byte)
	donToEthAddresses := make(map[uint32][]common.Address)
	for _, donMetadata := range topology.Metadata {
		p2pKeys, err := generateP2PKeys("", len(donMetadata.NodesMetadata))
		require.NoError(t, err, "failed to generate P2P keys")
		donToP2PKeys[donMetadata.ID] = p2pKeys

		ethKeys, addresses, err := generateEVMKeys("", len(donMetadata.NodesMetadata))
		require.NoError(t, err, "failed to generate EVM keys")
		donToEthKeys[donMetadata.ID] = ethKeys
		donToEthAddresses[donMetadata.ID] = addresses
	}

	for i, donMetadata := range topology.Metadata {
		for j := range donMetadata.NodesMetadata {
			nodeWithLabels := keystonetypes.NodeMetadata{}
			nodeType := keystonetypes.WorkerNode
			if j == 0 {
				nodeType = keystonetypes.BootstrapNode
			}
			nodeWithLabels.Labels = append(nodeWithLabels.Labels, &ptypes.Label{
				Key:   libnode.RoleLabelKey,
				Value: ptr.Ptr(nodeType),
			})
			nodeWithLabels.Labels = append(nodeWithLabels.Labels, &ptypes.Label{
				Key:   devenv.NodeLabelP2PIDType,
				Value: ptr.Ptr(donToP2PKeys[uint32(i+1)].peerIDs[j]),
			})
			nodeWithLabels.Labels = append(nodeWithLabels.Labels, &ptypes.Label{
				Key:   libnode.EthAddressKey,
				Value: ptr.Ptr(donToEthAddresses[uint32(i+1)][j].Hex()),
			})
			nodeWithLabels.Labels = append(nodeWithLabels.Labels, &ptypes.Label{
				Key:   libnode.NodeIndexKey,
				Value: ptr.Ptr(fmt.Sprint(j)),
			})
			nodeWithLabels.Labels = append(nodeWithLabels.Labels, &ptypes.Label{
				Key: libnode.HostLabelKey,
				// TODO this will only work with Docker, for CRIB we need a different approach
				Value: ptr.Ptr(fmt.Sprintf("%s-node%d", donMetadata.Name, j)),
			})

			topology.Metadata[i].NodesMetadata[j] = &nodeWithLabels
		}
	}

	peeringData, err := libdon.FindPeeringData(topology)
	require.NoError(t, err, "failed to get peering data")
	// prepare node configs
	donToConfigs := make(keystonetypes.DonsToConfigOverrides)

	var configErr error
	for _, donMetadata := range topology.Metadata {
		donToConfigs[donMetadata.ID], configErr = keystoneporconfig.GenerateConfigs(
			keystonetypes.GeneratePoRConfigsInput{
				DonMetadata:                 donMetadata,
				BlockchainOutput:            envOutput.blockchainOutput,
				DonID:                       donMetadata.ID,
				Flags:                       donMetadata.Flags,
				PeeringData:                 peeringData,
				CapabilitiesRegistryAddress: keystoneContractsOutput.CapabilitiesRegistryAddress,
				WorkflowRegistryAddress:     keystoneContractsOutput.WorkflowRegistryAddress,
				ForwarderAddress:            keystoneContractsOutput.ForwarderAddress,
				GatewayConnectorOutput:      envOutput.gatewayConnector,
			},
		)
		require.NoError(t, configErr, "failed to define config for DON %d", donMetadata.ID)
	}

	for i, donMetadata := range topology.Metadata {
		if configOverrides, ok := donToConfigs[donMetadata.ID]; ok {
			for j, configOverride := range configOverrides {
				if len(envInput.nodeSetInput[i].NodeSpecs)-1 < j {
					testLogger.Error().Msgf("Node %d has no config override", j)
					t.FailNow()
				}
				envInput.nodeSetInput[i].NodeSpecs[j].Node.TestConfigOverrides = configOverride

				ethKey := donToEthKeys[uint32(i+1)][j]
				p2pPublicKey := donToP2PKeys[uint32(i+1)].publicHexKeys[j]
				p2pPrivateKey := donToP2PKeys[uint32(i+1)].privateKeys[j]

				type NodeSecret struct {
					EthKey NodeEthKey `toml:"EthKey"`
					P2PKey NodeP2PKey `toml:"P2PKey"`
				}

				nodeSecret := NodeSecret{
					EthKey: NodeEthKey{
						JSON:     string(ethKey),
						Password: "",
						Selector: NodeEthKeySelector{
							ChainSelector: envOutput.chainSelector,
						},
					},
					P2PKey: NodeP2PKey{
						PrivateKey: p2pPrivateKey,
						PublicKey:  p2pPublicKey,
					},
				}

				nodeSecretString, err := toml.Marshal(nodeSecret)
				require.NoError(t, err, "failed to marshal node secrets")

				fmt.Println(string(nodeSecretString))

				envInput.nodeSetInput[i].NodeSpecs[j].Node.TestSecretsOverrides = string(nodeSecretString)
			}
		}
	}

	nodeOutput := make([]*keystonetypes.WrappedNodeOutput, 0, len(envInput.nodeSetInput))
	for _, nodeSetInput := range envInput.nodeSetInput {
		nodeset, nodesetErr := ns.NewSharedDBNodeSet(nodeSetInput.Input, envOutput.blockchainOutput)
		require.NoError(t, nodesetErr, "failed to deploy node set names %s", nodeSetInput.Name)

		nodeOutput = append(nodeOutput, &keystonetypes.WrappedNodeOutput{
			Output:       nodeset,
			NodeSetName:  nodeSetInput.Name,
			Capabilities: nodeSetInput.Capabilities,
		})
	}

	// for i, out := range nodeOutput {
	// 	c, err := clclient.New(out.Output.CLNodes)
	// 	require.NoError(t, err, "failed to create chainlink client")

	// 	err = clclient.ImportEVMKeys(c, donToEthKeys[uint32(i+1)], envOutput.blockchainOutput.ChainID)
	// 	require.NoError(t, err, "failed to import EVM keys")

	// 	err = clclient.ImportP2PKeys(c, donToP2PKeys[uint32(i+1)].encryptedP2PKeyJSONs)
	// 	require.NoError(t, err, "failed to import P2P keys")

	// 	// // delete the other p2p keys
	// 	// for j, client := range c {
	// 	// 	p2pKeys, err := client.MustReadP2PKeys()
	// 	// 	require.NoError(t, err, "failed to read P2P keys")

	// 	// 	fmt.Println("Genereated P2P keys for node", donToP2PKeys[uint32(i+1)].peerIDs[j])
	// 	// 	fmt.Println("P2P keys for node", j+1)
	// 	// 	for _, key := range p2pKeys.Data {

	// 	// 		// TODO remove old keys
	// 	// 	}
	// 	// }
	// }

	// Prepare the CLD environment and figure out DON topology; configure chains for nodes and job distributor
	// Ugly glue hack ¯\_(ツ)_/¯
	cldEnv, dons, err := libenv.BuildChainlinkDeploymentEnv(singeFileLogger, envOutput.jdOutput, nodeOutput, envOutput.blockchainOutput, envOutput.sethClient)
	require.NoError(t, err, "failed to build chainlink deployment environment")

	cldEnv.ExistingAddresses = chainsOnlyCld.ExistingAddresses

	donTopology := &keystonetypes.DonTopology{}
	donTopology.WorkflowDonID = topology.WorkflowDONID

	for i, donMetadata := range topology.Metadata {
		donTopology.Dons = append(donTopology.Dons, &keystonetypes.DonWithMetadata{
			DON:         dons[i],
			DonMetadata: donMetadata,
		})
	}

	// Fund the nodes
	for _, metaDon := range donTopology.Dons {
		for _, node := range metaDon.DON.Nodes {
			_, err := libfunding.SendFunds(zerolog.Logger{}, envOutput.sethClient, libtypes.FundsToSend{
				ToAddress:  common.HexToAddress(node.AccountAddr[envOutput.sethClient.Cfg.Network.ChainID]),
				Amount:     big.NewInt(5000000000000000000),
				PrivateKey: envOutput.sethClient.MustGetRootPrivateKey(),
			})
			require.NoError(t, err, "failed to send funds to node %s", node.AccountAddr[envOutput.sethClient.Cfg.Network.ChainID])
		}
	}

	// Workflow-specific configuration -- START
	deployFeedConsumerInput := &keystonetypes.DeployFeedConsumerInput{
		ChainSelector: envOutput.chainSelector,
		CldEnv:        chainsOnlyCld,
	}
	deployFeedsConsumerOutput, err := libcontracts.DeployFeedsConsumer(testLogger, deployFeedConsumerInput)
	require.NoError(t, err, "failed to deploy feeds consumer")

	configureFeedConsumerInput := &keystonetypes.ConfigureFeedConsumerInput{
		SethClient:            envOutput.sethClient,
		FeedConsumerAddress:   deployFeedsConsumerOutput.FeedConsumerAddress,
		AllowedSenders:        []common.Address{keystoneContractsOutput.ForwarderAddress},
		AllowedWorkflowOwners: []common.Address{envOutput.sethClient.MustGetRootKeyAddress()},
		AllowedWorkflowNames:  []string{in.WorkflowConfig.WorkflowName},
	}
	_, err = libcontracts.ConfigureFeedsConsumer(testLogger, configureFeedConsumerInput)
	require.NoError(t, err, "failed to configure feeds consumer")

	registerInput := registerPoRWorkflowInput{
		WorkflowConfig:              in.WorkflowConfig,
		chainSelector:               envOutput.chainSelector,
		workflowDonID:               donTopology.WorkflowDonID,
		feedID:                      in.WorkflowConfig.FeedID,
		workflowRegistryAddress:     keystoneContractsOutput.WorkflowRegistryAddress,
		feedConsumerAddress:         deployFeedsConsumerOutput.FeedConsumerAddress,
		capabilitiesRegistryAddress: keystoneContractsOutput.CapabilitiesRegistryAddress,
		priceProvider:               priceProvider,
		sethClient:                  envOutput.sethClient,
		deployerPrivateKey:          envOutput.deployerPrivateKey,
		blockchain:                  envOutput.blockchainOutput,
		binaryDownloadOutput:        binaryDownloadOutput,
	}

	err = registerPoRWorkflow(registerInput)
	require.NoError(t, err, "failed to register PoR workflow")
	// Workflow-specific configuration -- END

	donToJobSpecs := make(map[uint32]keystonetypes.DonJobs)
	var jobSpecsErr error
	for _, donWithMetadata := range donTopology.Dons {
		donToJobSpecs[donWithMetadata.ID], jobSpecsErr = keystonepor.GenerateJobSpecs(
			keystonetypes.GeneratePoRJobSpecsInput{
				CldEnv:                 cldEnv,
				DonWithMetadata:        *donWithMetadata,
				BlockchainOutput:       envOutput.blockchainOutput,
				DonID:                  donWithMetadata.ID,
				Flags:                  donWithMetadata.Flags,
				OCR3CapabilityAddress:  keystoneContractsOutput.OCR3CapabilityAddress,
				ExtraAllowedPorts:      extraAllowedPorts,
				ExtraAllowedIPs:        extraAllowedIPs,
				CronCapBinName:         cronCapabilityAssetFile,
				GatewayConnectorOutput: *envOutput.gatewayConnector,
			},
		)
		require.NoError(t, jobSpecsErr, "failed to define job specs for DON %d", donWithMetadata.ID)
	}

	// Configure nodes and create jobs
	createJobsInput := keystonetypes.CreateJobsInput{
		CldEnv:        cldEnv,
		DonTopology:   donTopology,
		DonToJobSpecs: donToJobSpecs,
	}

	err = libdon.CreateJobs(testLogger, createJobsInput)
	require.NoError(t, err, "failed to configure nodes and create jobs")

	// CAUTION: It is crucial to configure OCR3 jobs on nodes before configuring the workflow contracts.
	// Wait for OCR listeners to be ready before setting the configuration.
	// If the ConfigSet event is missed, OCR protocol will not start.
	testLogger.Info().Msg("Waiting 30s for OCR listeners to be ready...")
	time.Sleep(30 * time.Second)
	testLogger.Info().Msg("Proceeding to set OCR3 configuration.")

	// Configure the Forwarder, OCR3 and Capabilities contracts
	configureKeystoneInput := keystonetypes.ConfigureKeystoneInput{
		ChainSelector: envOutput.chainSelector,
		CldEnv:        cldEnv,
		Topology:      topology,
	}
	err = libcontracts.ConfigureKeystone(configureKeystoneInput)
	require.NoError(t, err, "failed to configure keystone contracts")

	// Set inputs in the test config, so that they can be saved
	in.KeystoneContracts = keystoneContractsInput
	in.FeedConsumer = deployFeedConsumerInput
	in.WorkflowRegistryConfiguration = workflowRegistryInput

	return &setupOutput{
		priceProvider:        priceProvider,
		feedsConsumerAddress: deployFeedsConsumerOutput.FeedConsumerAddress,
		forwarderAddress:     keystoneContractsOutput.ForwarderAddress,
		sethClient:           envOutput.sethClient,
		blockchainOutput:     envOutput.blockchainOutput,
		donTopology:          donTopology,
		nodeOutput:           nodeOutput,
	}
}

func TestKeystoneWithOCR3Workflow_SingleDon_MockedPrice(t *testing.T) {
	testLogger := framework.L

	// Load and validate test configuration
	in, err := framework.Load[TestConfig](t)
	require.NoError(t, err, "couldn't load test config")
	validateEnvVars(t, in)
	require.Len(t, in.NodeSets, 1, "expected 1 node set in the test config")

	binaryDownloadOutput, err := downloadBinaryFiles(in)
	require.NoError(t, err, "failed to download binary files")

	// Assign all capabilities to the single node set
	mustSetCapabilitiesFn := func(input []*ns.Input) []*keystonetypes.CapabilitiesAwareNodeSet {
		return []*keystonetypes.CapabilitiesAwareNodeSet{
			{
				Input:        input[0],
				Capabilities: keystonetypes.SingleDonFlags,
				DONType:      keystonetypes.WorkflowDON,
			},
		}
	}

	priceProvider, priceErr := NewFakePriceProvider(testLogger, in.Fake)
	require.NoError(t, priceErr, "failed to create fake price provider")

	setupOutput := setupTestEnvironment(t, testLogger, in, priceProvider, *binaryDownloadOutput, mustSetCapabilitiesFn)

	// Log extra information that might help debugging
	t.Cleanup(func() {
		if t.Failed() {
			logTestInfo(testLogger, in.WorkflowConfig.FeedID, in.WorkflowConfig.WorkflowName, setupOutput.feedsConsumerAddress.Hex(), setupOutput.forwarderAddress.Hex())

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

			debugDons := make([]*keystonetypes.DebugDon, 0, len(setupOutput.donTopology.Dons))
			for i, donWithMetadata := range setupOutput.donTopology.Dons {
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
			}
			lidebug.PrintTestDebug(t.Name(), testLogger, debugInput)
		}
	})

	testLogger.Info().Msg("Waiting for feed to update...")
	timeout := 5 * time.Minute // It can take a while before the first report is produced, particularly on CI.

	feedsConsumerInstance, err := feeds_consumer.NewKeystoneFeedsConsumer(setupOutput.feedsConsumerAddress, setupOutput.sethClient.Client)
	require.NoError(t, err, "failed to create feeds consumer instance")

	startTime := time.Now()
	feedBytes := common.HexToHash(in.WorkflowConfig.FeedID)

	assert.Eventually(t, func() bool {
		elapsed := time.Since(startTime).Round(time.Second)
		price, _, err := feedsConsumerInstance.GetPrice(
			setupOutput.sethClient.NewCallOpts(),
			feedBytes,
		)
		require.NoError(t, err, "failed to get price from Keystone Consumer contract")

		hasNextPrice := setupOutput.priceProvider.NextPrice(price, elapsed)
		if !hasNextPrice {
			testLogger.Info().Msgf("Feed not updated yet, waiting for %s", elapsed)
		}

		return !hasNextPrice
	}, timeout, 10*time.Second, "feed did not update, timeout after: %s", timeout)

	require.EqualValues(t, priceProvider.ExpectedPrices(), priceProvider.ActualPrices(), "prices do not match")
	testLogger.Info().Msgf("All %d prices were found in the feed", len(priceProvider.ExpectedPrices()))
}

func TestKeystoneWithOCR3Workflow_TwoDons_LivePrice(t *testing.T) {
	testLogger := framework.L

	// Load and validate test configuration
	in, err := framework.Load[TestConfig](t)
	require.NoError(t, err, "couldn't load test config")
	validateEnvVars(t, in)
	require.Len(t, in.NodeSets, 2, "expected 2 node sets in the test config")

	binaryDownloadOutput, err := downloadBinaryFiles(in)
	require.NoError(t, err, "failed to download binary files")

	mustSetCapabilitiesFn := func(input []*ns.Input) []*keystonetypes.CapabilitiesAwareNodeSet {
		return []*keystonetypes.CapabilitiesAwareNodeSet{
			{
				Input:        input[0],
				Capabilities: []string{keystonetypes.OCR3Capability, keystonetypes.CustomComputeCapability, keystonetypes.CronCapability},
				DONType:      keystonetypes.WorkflowDON,
			},
			{
				Input:        input[1],
				Capabilities: []string{keystonetypes.WriteEVMCapability},
				DONType:      keystonetypes.CapabilitiesDON, // <----- it's crucial to set the correct DON type
			},
		}
	}

	priceProvider := NewTrueUSDPriceProvider(testLogger)
	setupOutput := setupTestEnvironment(t, testLogger, in, priceProvider, *binaryDownloadOutput, mustSetCapabilitiesFn)

	// Log extra information that might help debugging
	t.Cleanup(func() {
		if t.Failed() {
			logTestInfo(testLogger, in.WorkflowConfig.FeedID, in.WorkflowConfig.WorkflowName, setupOutput.feedsConsumerAddress.Hex(), setupOutput.forwarderAddress.Hex())

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

			debugDons := make([]*keystonetypes.DebugDon, 0, len(setupOutput.donTopology.Dons))
			for i, donWithMetadata := range setupOutput.donTopology.Dons {
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
			}
			lidebug.PrintTestDebug(t.Name(), testLogger, debugInput)
		}
	})

	testLogger.Info().Msg("Waiting for feed to update...")
	timeout := 5 * time.Minute // It can take a while before the first report is produced, particularly on CI.

	feedsConsumerInstance, err := feeds_consumer.NewKeystoneFeedsConsumer(setupOutput.feedsConsumerAddress, setupOutput.sethClient.Client)
	require.NoError(t, err, "failed to create feeds consumer instance")

	startTime := time.Now()
	feedBytes := common.HexToHash(in.WorkflowConfig.FeedID)

	assert.Eventually(t, func() bool {
		elapsed := time.Since(startTime).Round(time.Second)
		price, _, err := feedsConsumerInstance.GetPrice(
			setupOutput.sethClient.NewCallOpts(),
			feedBytes,
		)
		require.NoError(t, err, "failed to get price from Keystone Consumer contract")

		hasNextPrice := setupOutput.priceProvider.NextPrice(price, elapsed)
		if !hasNextPrice {
			testLogger.Info().Msgf("Feed not updated yet, waiting for %s", elapsed)
		}

		return !hasNextPrice
	}, timeout, 10*time.Second, "feed did not update, timeout after: %s", timeout)

	require.EqualValues(t, priceProvider.ExpectedPrices(), priceProvider.ActualPrices(), "prices do not match")
	testLogger.Info().Msgf("All %d prices were found in the feed", len(priceProvider.ExpectedPrices()))
}
