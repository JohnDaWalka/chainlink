package capabilities_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink/deployment"
	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/feeds_consumer"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	ctfconfig "github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	libcrecli "github.com/smartcontractkit/chainlink/system-tests/lib/crecli"
	keystoneporcrecli "github.com/smartcontractkit/chainlink/system-tests/lib/crecli/por"
	libfunding "github.com/smartcontractkit/chainlink/system-tests/lib/funding"
	keystonecapabilities "github.com/smartcontractkit/chainlink/system-tests/lib/keystone/capabilities"
	libcontracts "github.com/smartcontractkit/chainlink/system-tests/lib/keystone/contracts"
	lidebug "github.com/smartcontractkit/chainlink/system-tests/lib/keystone/debug"
	libdon "github.com/smartcontractkit/chainlink/system-tests/lib/keystone/don"
	keystoneporconfig "github.com/smartcontractkit/chainlink/system-tests/lib/keystone/don/config/por"
	libjobs "github.com/smartcontractkit/chainlink/system-tests/lib/keystone/don/jobs"
	keystonepor "github.com/smartcontractkit/chainlink/system-tests/lib/keystone/don/jobs/por"
	libenv "github.com/smartcontractkit/chainlink/system-tests/lib/keystone/environment"
	keystonetypes "github.com/smartcontractkit/chainlink/system-tests/lib/keystone/types"
	libtypes "github.com/smartcontractkit/chainlink/system-tests/lib/types"
)

const (
	cronCapabilityAssetFile            = "amd64_cron"
	ghReadTokenEnvVarName              = "GITHUB_READ_TOKEN"
	E2eJobDistributorImageEnvVarName   = "E2E_JD_IMAGE"
	E2eJobDistributorVersionEnvVarName = "E2E_JD_VERSION"
)

type TestConfig struct {
	BlockchainA    *blockchain.Input `toml:"blockchain_a" validate:"required"`
	NodeSets       []*ns.Input       `toml:"nodesets" validate:"required"`
	WorkflowConfig *WorkflowConfig   `toml:"workflow_config" validate:"required"`
	JD             *jd.Input         `toml:"jd" validate:"required"`
	Fake           *fake.Input       `toml:"fake"`
}

type WorkflowConfig struct {
	UseCRECLI                bool `toml:"use_cre_cli"`
	ShouldCompileNewWorkflow bool `toml:"should_compile_new_workflow" validate:"no_cre_no_compilation,disabled_in_ci"`
	// Tells the test where the workflow to compile is located
	WorkflowFolderLocation *string             `toml:"workflow_folder_location" validate:"required_if=ShouldCompileNewWorkflow true"`
	CompiledWorkflowConfig *CompiledConfig     `toml:"compiled_config" validate:"required_if=ShouldCompileNewWorkflow false"`
	DependenciesConfig     *DependenciesConfig `toml:"dependencies" validate:"required"`
	WorkflowName           string              `toml:"workflow_name" validate:"required" `
	FeedID                 string              `toml:"feed_id" validate:"required,no0xPrefix"`
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

func no0xPrefix(fl validator.FieldLevel) bool {
	return !strings.HasPrefix(fl.Field().String(), "0x")
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
	err = framework.Validator.RegisterValidation("no0xPrefix", no0xPrefix)
	if err != nil {
		panic(errors.Wrap(err, "failed to register no0xPrefix validator"))
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

// this is a small hack to avoid changing the reusable workflow, which doesn't allow to run any pre-execution hooks
func downloadBinaryFiles(in *TestConfig) error {
	var ghReadToken string
	if os.Getenv("CI") == "true" {
		ghReadToken = ctfconfig.MustReadEnvVar_String(ghReadTokenEnvVarName)
	} else {
		ghReadToken = os.Getenv(ghReadTokenEnvVarName)
	}

	_, err := keystonecapabilities.DownloadCapabilityFromRelease(ghReadToken, in.WorkflowConfig.DependenciesConfig.CapabiltiesVersion, cronCapabilityAssetFile)
	if err != nil {
		return errors.Wrap(err, "failed to download cron capability. Make sure token has content:read permissions to the capabilities repo")
	}

	if in.WorkflowConfig.UseCRECLI {
		err = libcrecli.DownloadAndInstallChainlinkCLI(ghReadToken, in.WorkflowConfig.DependenciesConfig.CRECLIVersion)
		if err != nil {
			return errors.Wrap(err, "failed to download and install CRE CLI. Make sure token has content:read permissions to the dev-platform repo")
		}
	}

	return nil
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
		compilationResult, err := libcrecli.CompileWorkflow(*input.WorkflowConfig.WorkflowFolderLocation, workflowConfigFile, settingsFile)
		if err != nil {
			return errors.Wrap(err, "failed to compile workflow")
		}

		workflowURL = compilationResult.WorkflowURL
		workflowConfigURL = compilationResult.ConfigURL
	} else {
		workflowURL = input.WorkflowConfig.CompiledWorkflowConfig.BinaryURL
		workflowConfigURL = input.WorkflowConfig.CompiledWorkflowConfig.ConfigURL
	}

	registerErr := libcrecli.RegisterWorkflow(input.WorkflowName, workflowURL, workflowConfigURL, settingsFile)
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

func setupFakeDataProvider(testLogger zerolog.Logger, input *fake.Input, expectedPrices []float64, priceIndex *int) (string, error) {
	_, err := fake.NewFakeDataProvider(input)
	if err != nil {
		return "", errors.Wrap(err, "failed to set up fake data provider")
	}
	fakeAPIPath := "/fake/api/price"
	host := framework.HostDockerInternal()
	fakeFinalURL := fmt.Sprintf("%s:%d%s", host, input.Port, fakeAPIPath)

	getPriceResponseFn := func() map[string]interface{} {
		response := map[string]interface{}{
			"accountName": "TrueUSD",
			"totalTrust":  expectedPrices[*priceIndex],
			"ripcord":     false,
			"updatedAt":   time.Now().Format(time.RFC3339),
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

// PriceProvider abstracts away the logic of checking whether the feed has been correctly updated
// and it also returns port and URL of the price provider. This is so, because when using a mocked
// price provider we need start a separate service and whitelist its port and IP with the gateway job.
// Also, since it's a mocked price provider we can now check whether the feed has been correctly updated
// instead of only checking whether it has some price that's != 0.
type PriceProvider interface {
	URL() string
	NextPrice(price *big.Int, elapsed time.Duration) bool
	ExpectedPrices() []*big.Int
	ActualPrices() []*big.Int
}

// TrueUSDPriceProvider is a PriceProvider implementation that uses a live feed to get the price
type TrueUSDPriceProvider struct {
	testLogger   zerolog.Logger
	url          string
	actualPrices []*big.Int
}

func NewTrueUSDPriceProvider(testLogger zerolog.Logger) PriceProvider {
	return &TrueUSDPriceProvider{
		testLogger: testLogger,
		url:        "https://api.real-time-reserves.verinumus.io/v1/chainlink/proof-of-reserves/TrueUSD",
	}
}

func (l *TrueUSDPriceProvider) NextPrice(price *big.Int, elapsed time.Duration) bool {
	// if price is nil or 0 it means that the feed hasn't been updated yet
	if price == nil || price.Cmp(big.NewInt(0)) == 0 {
		return true
	}

	l.testLogger.Info().Msgf("Feed updated after %s - price set, price=%s", elapsed, price)
	l.actualPrices = append(l.actualPrices, price)

	// no other price to return, we are done
	return false
}

func (l *TrueUSDPriceProvider) URL() string {
	return l.url
}

func (l *TrueUSDPriceProvider) ExpectedPrices() []*big.Int {
	// we don't have a way to check the price in the live feed, so we always assume it's correct
	// as long as it's != 0. And we only wait for the first price to be set.
	return l.actualPrices
}

func (l *TrueUSDPriceProvider) ActualPrices() []*big.Int {
	// we don't have a way to check the price in the live feed, so we always assume it's correct
	// as long as it's != 0. And we only wait for the first price to be set.
	return l.actualPrices
}

// FakePriceProvider is a PriceProvider implementation that uses a mocked feed to get the price
// It returns a configured price sequence and makes sure that the feed has been correctly updated
type FakePriceProvider struct {
	t              *testing.T
	testLogger     zerolog.Logger
	priceIndex     *int
	url            string
	expectedPrices []*big.Int
	actualPrices   []*big.Int
}

func NewFakePriceProvider(testLogger zerolog.Logger, in *TestConfig) (PriceProvider, error) {
	priceIndex := ptr.Ptr(0)
	expectedPricesFloat64 := []float64{182.9, 122.01}
	expectedPrices := make([]*big.Int, len(expectedPricesFloat64))
	for i, p := range expectedPricesFloat64 {
		// convert float64 to big.Int by multiplying by 100
		// just like the PoR workflow does
		expectedPrices[i] = libc.Float64ToBigInt(p)
	}

	url, err := setupFakeDataProvider(testLogger, in.Fake, expectedPricesFloat64, priceIndex)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set up fake data provider")
	}

	return &FakePriceProvider{
		testLogger:     testLogger,
		expectedPrices: expectedPrices,
		priceIndex:     priceIndex,
		url:            url,
	}, nil
}

func (f *FakePriceProvider) priceAlreadyFound(price *big.Int) bool {
	for _, p := range f.actualPrices {
		if p.Cmp(price) == 0 {
			return true
		}
	}

	return false
}

func (f *FakePriceProvider) NextPrice(price *big.Int, elapsed time.Duration) bool {
	// if price is nil or 0 it means that the feed hasn't been updated yet
	if price == nil || price.Cmp(big.NewInt(0)) == 0 {
		return true
	}

	if !f.priceAlreadyFound(price) {
		f.testLogger.Info().Msgf("Feed updated after %s - price set, price=%s", elapsed, price)
		f.actualPrices = append(f.actualPrices, price)

		if len(f.actualPrices) == len(f.expectedPrices) {
			// all prices found, nothing more to check
			return false
		}

		require.Less(f.t, len(f.actualPrices), len(f.expectedPrices), "more prices found than expected")
		f.testLogger.Info().Msgf("Changing price provider price to %s", f.expectedPrices[len(f.actualPrices)].String())
		*f.priceIndex = len(f.actualPrices)

		// set new price and continue checking
		return true
	}

	// continue checking, price not updated yet
	return true
}

func (f *FakePriceProvider) ActualPrices() []*big.Int {
	return f.actualPrices
}

func (f *FakePriceProvider) ExpectedPrices() []*big.Int {
	return f.expectedPrices
}

func (f *FakePriceProvider) URL() string {
	return f.url
}

func extraAllowedPortsAndIps(testLogger zerolog.Logger, fakePort int, nodeOutput *ns.Output) ([]string, []int, error) {
	// we need to explicitly allow the port used by the fake data provider
	// and IP corresponding to host.docker.internal or the IP of the host machine, if we are running on Linux,
	// because that's where the fake data provider is running
	var hostIP string
	var err error

	system := runtime.GOOS
	switch system {
	case "darwin":
		hostIP, err = libdon.ResolveHostDockerInternaIP(testLogger, nodeOutput)
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

type InfrastructureInput struct {
	jdInput         *jd.Input
	nodeSetInput    []*keystonetypes.CapabilitiesAwareNodeSet
	blockchainInput *blockchain.Input
}

type InfrastructureOutput struct {
	chainSelector      uint64
	nodeOuput          []*keystonetypes.WrappedNodeOutput
	blockchainOutput   *blockchain.Output
	jdOutput           *jd.Output
	cldEnv             *deployment.Environment
	donTopology        *keystonetypes.DonTopology
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

	nodeOutput := make([]*keystonetypes.WrappedNodeOutput, 0, len(input.nodeSetInput))
	for _, nsInput := range input.nodeSetInput {
		nodeset, nodesetErr := ns.NewSharedDBNodeSet(nsInput.Input, blockchainOutput)
		if nodesetErr != nil {
			return nil, errors.Wrapf(nodesetErr, "failed to deploy node set names %s", nsInput.Name)
		}

		nodeOutput = append(nodeOutput, &keystonetypes.WrappedNodeOutput{
			Output:       nodeset,
			NodeSetName:  nsInput.Name,
			Capabilities: nsInput.Capabilities,
		})
	}

	// Prepare the CLD environment and figure out DON topology; configure chains for nodes and job distributor
	// Ugly glue hack ¯\_(ツ)_/¯
	cldEnv, donTopology, err := libenv.BuildTopologyAndCLDEnvironment(cldLogger, input.nodeSetInput, jdOutput, nodeOutput, blockchainOutput, sethClient)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build topology and CLD environment")
	}

	// Fund the nodes
	for _, metaDon := range donTopology.MetaDons {
		for _, node := range metaDon.DON.Nodes {
			_, err := libfunding.SendFunds(zerolog.Logger{}, sethClient, libtypes.FundsToSend{
				ToAddress:  common.HexToAddress(node.AccountAddr[sethClient.Cfg.Network.ChainID]),
				Amount:     big.NewInt(5000000000000000000),
				PrivateKey: sethClient.MustGetRootPrivateKey(),
			})
			if err != nil {
				return nil, errors.Wrapf(err, "failed to send funds to node %s", node.AccountAddr[sethClient.Cfg.Network.ChainID])
			}
		}
	}

	return &InfrastructureOutput{
		chainSelector:      chainSelector,
		nodeOuput:          nodeOutput,
		blockchainOutput:   blockchainOutput,
		jdOutput:           jdOutput,
		cldEnv:             cldEnv,
		donTopology:        donTopology,
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
}

// TODO to each input add output and cache, same way sergey did in ctfv2
func setupTestEnvironment(t *testing.T, testLogger zerolog.Logger, in *TestConfig, priceProvider PriceProvider, mustSetCapabilitiesFn func(input []*ns.Input) []*keystonetypes.CapabilitiesAwareNodeSet) *setupOutput {
	// Universal setup -- START
	envInput := InfrastructureInput{
		jdInput:         in.JD,
		nodeSetInput:    mustSetCapabilitiesFn(in.NodeSets),
		blockchainInput: in.BlockchainA,
	}
	envOutput, err := CreateInfrastructure(cldlogger.NewSingleFileLogger(t), testLogger, envInput)
	require.NoError(t, err, "failed to start environment")

	// Deploy keystone contracts (forwarder, capability registry, ocr3 capability, workflow registry)
	keystoneContractsInput := keystonetypes.KeystoneContractsInput{
		ChainSelector: envOutput.chainSelector,
		CldEnv:        envOutput.cldEnv,
	}

	keystoneContractsOutput, err := libcontracts.DeployKeystone(testLogger, keystoneContractsInput)
	require.NoError(t, err, "failed to deploy keystone contracts")

	// Configure Workflow Registry
	workflowRegistryInput := keystonetypes.WorkflowRegistryInput{
		ChainSelector:  envOutput.chainSelector,
		CldEnv:         envOutput.cldEnv,
		AllowedDonIDs:  []uint32{envOutput.donTopology.WorkflowDONID},
		WorkflowOwners: []common.Address{envOutput.sethClient.MustGetRootKeyAddress()},
	}

	err = libcontracts.ConfigureWorkflowRegistry(testLogger, workflowRegistryInput)
	require.NoError(t, err, "failed to configure workflow registry")
	// Universal setup -- END

	// Workflow-specific configuration -- START
	deployFeedConsumerInput := keystonetypes.DeployFeedConsumerInput{
		ChainSelector: envOutput.chainSelector,
		CldEnv:        envOutput.cldEnv,
	}
	deployFeedsConsumerOutput, err := libcontracts.DeployFeedsConsumer(testLogger, deployFeedConsumerInput)
	require.NoError(t, err, "failed to deploy feeds consumer")

	configureFeedConsumerInput := keystonetypes.ConfigureFeedConsumerInput{
		SethClient:            envOutput.sethClient,
		FeedConsumerAddress:   deployFeedsConsumerOutput.Address,
		AllowedSenders:        []common.Address{keystoneContractsOutput.ForwarderAddress},
		AllowedWorkflowOwners: []common.Address{envOutput.sethClient.MustGetRootKeyAddress()},
		AllowedWorkflowNames:  []string{in.WorkflowConfig.WorkflowName},
	}
	err = libcontracts.ConfigureFeedsConsumer(testLogger, configureFeedConsumerInput)
	require.NoError(t, err, "failed to configure feeds consumer")

	registerInput := registerPoRWorkflowInput{
		WorkflowConfig:              in.WorkflowConfig,
		chainSelector:               envOutput.chainSelector,
		workflowDonID:               envOutput.donTopology.WorkflowDONID,
		feedID:                      in.WorkflowConfig.FeedID,
		workflowRegistryAddress:     keystoneContractsOutput.WorkflowRegistryAddress,
		feedConsumerAddress:         deployFeedsConsumerOutput.Address,
		capabilitiesRegistryAddress: keystoneContractsOutput.CapabilitiesRegistryAddress,
		priceProvider:               priceProvider,
		sethClient:                  envOutput.sethClient,
		deployerPrivateKey:          envOutput.deployerPrivateKey,
		blockchain:                  envOutput.blockchainOutput,
	}

	err = registerPoRWorkflow(registerInput)
	require.NoError(t, err, "failed to register PoR workflow")
	// Workflow-specific configuration -- END

	// Universal setup -- CONTINUED
	// Allow extra IPs and ports for the fake data provider, which is running on host machine and requires explicit whitelisting
	var extraAllowedIPs []string
	var extraAllowedPorts []int
	if _, ok := priceProvider.(*FakePriceProvider); ok {
		extraAllowedIPs, extraAllowedPorts, err = extraAllowedPortsAndIps(testLogger, in.Fake.Port, envOutput.donTopology.MetaDons[0].NodeOutput.Output)
		require.NoError(t, err, "failed to get extra allowed ports and IPs")
	}

	// Prepare job specs and node configs
	configsAndJobsInput := jobsAndConfigsInput{
		donTopology:                 envOutput.donTopology,
		blockchainOutput:            envOutput.blockchainOutput,
		gatewayConnectorOutput:      envOutput.gatewayConnector,
		workflowRegistryAddress:     keystoneContractsOutput.WorkflowRegistryAddress,
		forwarderAddress:            keystoneContractsOutput.ForwarderAddress,
		capabilitiesRegistryAddress: keystoneContractsOutput.CapabilitiesRegistryAddress,
		ocr3capabilityAddress:       keystoneContractsOutput.OCR3CapabilityAddress,
		cldEnv:                      envOutput.cldEnv,
		extraAllowedIPs:             extraAllowedIPs,
		extraAllowedPorts:           extraAllowedPorts,
	}

	configsAndJobsOutput, err := prepareJobSpecsAndNodeConfigs(configsAndJobsInput)
	require.NoError(t, err, "failed to prepare job specs and node configs")

	// Configure nodes and create jobs
	configureDonInput := keystonetypes.ConfigureDonInput{
		CldEnv:               envOutput.cldEnv,
		BlockchainOutput:     envOutput.blockchainOutput,
		JdOutput:             envOutput.jdOutput,
		DonTopology:          envOutput.donTopology,
		DonToJobSpecs:        configsAndJobsOutput.donToJobSpecs,
		DonToConfigOverrides: configsAndJobsOutput.nodeToConfigOverrides,
	}
	configureDonOutput, err := libdon.Configure(t, testLogger, configureDonInput)
	require.NoError(t, err, "failed to configure nodes and create jobs")

	_ = configureDonOutput

	// CAUTION: It is crucial to configure OCR3 jobs on nodes before configuring the workflow contracts.
	// Wait for OCR listeners to be ready before setting the configuration.
	// If the ConfigSet event is missed, OCR protocol will not start.
	testLogger.Info().Msg("Waiting 30s for OCR listeners to be ready...")
	time.Sleep(30 * time.Second)
	testLogger.Info().Msg("Proceeding to set OCR3 configuration.")

	// Configure the Forwarder, OCR3 and Capabilities contracts
	configureKeystoneInput := keystonetypes.ConfigureKeystoneInput{
		ChainSelector: envOutput.chainSelector,
		CldEnv:        envOutput.cldEnv,
		DonTopology:   envOutput.donTopology,
	}
	err = libcontracts.ConfigureKeystone(configureKeystoneInput)
	require.NoError(t, err, "failed to configure keystone contracts")

	return &setupOutput{
		priceProvider:        priceProvider,
		feedsConsumerAddress: deployFeedsConsumerOutput.Address,
		forwarderAddress:     keystoneContractsOutput.ForwarderAddress,
		sethClient:           envOutput.sethClient,
		blockchainOutput:     envOutput.blockchainOutput,
		donTopology:          envOutput.donTopology,
	}
}

type jobsAndConfigsInput struct {
	donTopology                 *keystonetypes.DonTopology
	blockchainOutput            *blockchain.Output
	gatewayConnectorOutput      *keystonetypes.GatewayConnectorOutput
	workflowRegistryAddress     common.Address
	forwarderAddress            common.Address
	capabilitiesRegistryAddress common.Address
	ocr3capabilityAddress       common.Address
	cldEnv                      *deployment.Environment
	extraAllowedIPs             []string
	extraAllowedPorts           []int
}

type jobsAndConfigsOutput struct {
	donToJobSpecs         keystonetypes.DonsToJobSpecs
	nodeToConfigOverrides keystonetypes.DonsToConfigOverrides
}

func prepareJobSpecsAndNodeConfigs(input jobsAndConfigsInput) (*jobsAndConfigsOutput, error) {
	peeringData, err := libdon.FindPeeringData(input.donTopology.MetaDons)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get peering data")
	}

	// prepare node configs
	donToConfigs := make(keystonetypes.DonsToConfigOverrides)
	var configErr error
	for _, donTopology := range input.donTopology.MetaDons {
		donToConfigs[donTopology.ID], configErr = keystoneporconfig.Define(
			donTopology.DON,
			donTopology.NodeInput,
			donTopology.NodeOutput,
			input.blockchainOutput,
			donTopology.ID,
			donTopology.Flags,
			peeringData,
			input.capabilitiesRegistryAddress,
			input.workflowRegistryAddress,
			input.forwarderAddress,
			input.gatewayConnectorOutput,
		)
		if configErr != nil {
			return nil, errors.Wrapf(configErr, "failed to define config for DON %d", donTopology.ID)
		}
	}

	// define jobs
	donToJobSpecs := make(map[uint32]keystonetypes.DonJobs)
	var jobSpecsErr error
	for _, donTopology := range input.donTopology.MetaDons {
		donToJobSpecs[donTopology.ID], jobSpecsErr = keystonepor.Define(
			input.cldEnv,
			donTopology.DON,
			donTopology.NodeOutput,
			input.blockchainOutput,
			input.ocr3capabilityAddress,
			donTopology.ID,
			donTopology.Flags,
			input.extraAllowedPorts,
			input.extraAllowedIPs,
			cronCapabilityAssetFile,
			*input.gatewayConnectorOutput,
		)
		if jobSpecsErr != nil {
			return nil, errors.Wrapf(jobSpecsErr, "failed to define job specs for DON %d", donTopology.ID)
		}
	}

	return &jobsAndConfigsOutput{
		nodeToConfigOverrides: donToConfigs,
		donToJobSpecs:         donToJobSpecs,
	}, nil
}
func TestKeystoneWithOCR3Workflow_SingleDon_MockedPrice(t *testing.T) {
	testLogger := framework.L

	// Load and validate test configuration
	in, err := framework.Load[TestConfig](t)
	require.NoError(t, err, "couldn't load test config")
	validateEnvVars(t, in)
	require.Len(t, in.NodeSets, 1, "expected 1 node set in the test config")

	err = downloadBinaryFiles(in)
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

	priceProvider, priceErr := NewFakePriceProvider(testLogger, in)
	require.NoError(t, priceErr, "failed to create fake price provider")

	setupOutput := setupTestEnvironment(t, testLogger, in, priceProvider, mustSetCapabilitiesFn)

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

			debugInput := keystonetypes.DebugInput{
				DonTopology:      setupOutput.donTopology,
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

	err = downloadBinaryFiles(in)
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
	setupOutput := setupTestEnvironment(t, testLogger, in, priceProvider, mustSetCapabilitiesFn)

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

			debugInput := keystonetypes.DebugInput{
				DonTopology:      setupOutput.donTopology,
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
