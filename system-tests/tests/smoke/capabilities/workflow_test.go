package capabilities_test

import (
	"context"
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
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/feeds_consumer"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	ctfconfig "github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	libcrecli "github.com/smartcontractkit/chainlink/system-tests/lib/crecli"
	keystoneporcrecli "github.com/smartcontractkit/chainlink/system-tests/lib/crecli/por"
	keystonecapabilities "github.com/smartcontractkit/chainlink/system-tests/lib/keystone/capabilities"
	libcontracts "github.com/smartcontractkit/chainlink/system-tests/lib/keystone/contracts"
	lidebug "github.com/smartcontractkit/chainlink/system-tests/lib/keystone/debug"
	libdon "github.com/smartcontractkit/chainlink/system-tests/lib/keystone/don"
	keystoneporconfig "github.com/smartcontractkit/chainlink/system-tests/lib/keystone/don/config/por"
	libjobs "github.com/smartcontractkit/chainlink/system-tests/lib/keystone/don/jobs"
	keystonepor "github.com/smartcontractkit/chainlink/system-tests/lib/keystone/don/jobs/por"
	keystoneenv "github.com/smartcontractkit/chainlink/system-tests/lib/keystone/environment"
	keystonetypes "github.com/smartcontractkit/chainlink/system-tests/lib/keystone/types"
)

const (
	cronCapabilityAssetFile = "amd64_cron"
	ghReadTokenEnvVarName   = "GITHUB_READ_TOKEN"
)

type TestConfig struct {
	BlockchainA    *blockchain.Input                         `toml:"blockchain_a" validate:"required"`
	NodeSets       []*keystonetypes.CapabilitiesAwareNodeSet `toml:"nodesets" validate:"required"`
	WorkflowConfig *WorkflowConfig                           `toml:"workflow_config" validate:"required"`
	JD             *jd.Input                                 `toml:"jd" validate:"required"`
	PriceProvider  *PriceProviderConfig                      `toml:"price_provider"`
}

type WorkflowConfig struct {
	UseCRECLI                bool `toml:"use_cre_cli"`
	ShouldCompileNewWorkflow bool `toml:"should_compile_new_workflow"`
	// Tells the test where the workflow to compile is located
	WorkflowFolderLocation *string             `toml:"workflow_folder_location"`
	CompiledWorkflowConfig *CompiledConfig     `toml:"compiled_config"`
	DependenciesConfig     *DependenciesConfig `toml:"dependencies"`
	WorkflowName           string              `toml:"workflow_name" validate:"required" `
}

// Defines relases/versions of test dependencies that will be downloaded from Github
type DependenciesConfig struct {
	CapabiltiesVersion string `toml:"capabilities_version"`
	CRECLIVersion      string `toml:"cre_cli_version"`
}

// Defines the location of already compiled workflow binary and config files
// They will be used if WorkflowConfig.ShouldCompileNewWorkflow is `false`
// Otherwise test will compile and upload a new workflow
type CompiledConfig struct {
	BinaryURL string `toml:"binary_url"`
	ConfigURL string `toml:"config_url"`
}

type FakeConfig struct {
	*fake.Input
	Prices []float64 `toml:"prices"`
}

type PriceProviderConfig struct {
	Fake   *FakeConfig `toml:"fake"`
	FeedID string      `toml:"feed_id" validate:"required"`
	URL    string      `toml:"url"`
}

func validateInputsAndEnvVars(t *testing.T, in *TestConfig) {
	require.NotEmpty(t, os.Getenv("PRIVATE_KEY"), "PRIVATE_KEY env var must be set")
	require.NotEmpty(t, in.WorkflowConfig.DependenciesConfig, "dependencies config must be set")

	if !in.WorkflowConfig.UseCRECLI {
		require.False(t, in.WorkflowConfig.ShouldCompileNewWorkflow, "if you are not using CRE CLI you cannot compile a new workflow")
	}

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
		require.False(t, in.WorkflowConfig.ShouldCompileNewWorkflow, "you cannot compile a new workflow in the CI as of now due to issues with generating a gist write token")

		// we use this special function to subsitute a placeholder env variable with the actual environment variable name
		// it is defined in .github/e2e-tests.yml as '{{ env.GITHUB_API_TOKEN }}'
		ghReadToken = ctfconfig.MustReadEnvVar_String(ghReadTokenEnvVarName)
	} else {
		ghReadToken = os.Getenv(ghReadTokenEnvVarName)
	}

	require.NotEmpty(t, ghReadToken, ghReadTokenEnvVarName+" env var must be set")
	require.NotEmpty(t, in.WorkflowConfig.DependenciesConfig.CapabiltiesVersion, "capabilities_version must be set in the dependencies config")

	_, err := keystonecapabilities.DownloadCapabilityFromRelease(ghReadToken, in.WorkflowConfig.DependenciesConfig.CapabiltiesVersion, cronCapabilityAssetFile)
	require.NoError(t, err, "failed to download cron capability. Make sure token has content:read permissions to the capabilities repo")

	if in.WorkflowConfig.UseCRECLI {
		require.NotEmpty(t, in.WorkflowConfig.DependenciesConfig.CRECLIVersion, "chainlink_cli_version must be set in the dependencies config")

		err = libcrecli.DownloadAndInstallChainlinkCLI(ghReadToken, in.WorkflowConfig.DependenciesConfig.CRECLIVersion)
		require.NoError(t, err, "failed to download and install CRE CLI. Make sure token has content:read permissions to the dev-platform repo")

		if in.WorkflowConfig.ShouldCompileNewWorkflow {
			gistWriteToken := os.Getenv("GIST_WRITE_TOKEN")
			require.NotEmpty(t, gistWriteToken, "GIST_WRITE_TOKEN must be set to use CRE CLI to compile workflows. It requires gist:read and gist:write permissions")
			err := os.Setenv("GITHUB_API_TOKEN", gistWriteToken)
			require.NoError(t, err, "failed to set GITHUB_API_TOKEN env var")
			require.NotEmpty(t, in.WorkflowConfig.WorkflowFolderLocation, "workflow_folder_location must be set, when compiling new workflow")
		}
	}

	if in.PriceProvider.Fake == nil {
		require.NotEmpty(t, in.PriceProvider.URL, "URL must be set in the price provider config, if fake provider is not used")
	}

	if len(in.NodeSets) == 1 {
		noneEmpty := in.NodeSets[0].DONType != "" && len(in.NodeSets[0].Capabilities) > 0
		bothEmpty := in.NodeSets[0].DONType == "" && len(in.NodeSets[0].Capabilities) == 0
		require.True(t, noneEmpty || bothEmpty, "either both DONType and Capabilities must be set or both must be empty, when using only one node set")
	} else {
		for _, nodeSet := range in.NodeSets {
			require.NotEmpty(t, nodeSet.Capabilities, "capabilities must be set for each node set")
			require.NotEmpty(t, nodeSet.DONType, "don_type must be set for each node set")
		}
	}

	// make sure the feed id is in the correct format
	in.PriceProvider.FeedID = strings.TrimPrefix(in.PriceProvider.FeedID, "0x")
}

func registerPoRWorkflow(t *testing.T, in *TestConfig, workflowName string, keystoneEnv *keystonetypes.KeystoneEnvironment, priceProvider PriceProvider) {
	require.NotNil(t, keystoneEnv, "keystone environment must be set")
	require.NotNil(t, keystoneEnv.Environment, "environment must be set")
	require.NotNil(t, keystoneEnv.SethClient, "seth client must be set")
	require.NotNil(t, keystoneEnv.Blockchain, "blockchain must be set")
	require.NotEmpty(t, keystoneEnv.ChainSelector, "chain selector must be set")
	require.NotNil(t, keystoneEnv.KeystoneContractAddresses, "keystone contract addresses must be set")
	require.NotEmpty(t, keystoneEnv.KeystoneContractAddresses.CapabilitiesRegistryAddress, "capabilities registry address must be set")
	require.NotEmpty(t, keystoneEnv.KeystoneContractAddresses.WorkflowRegistryAddress, "workflow registry address must be set")
	require.NotEmpty(t, keystoneEnv.KeystoneContractAddresses.FeedsConsumerAddress, "feed consumer address must be set")
	require.NotEmpty(t, keystoneEnv.DeployerPrivateKey, "deployer private key must be set")
	require.NotEmpty(t, keystoneEnv.WorkflowDONID, "workflow DON ID must be set")

	// Register workflow directly using the provided binary and config URLs
	// This is a legacy solution, probably we can remove it soon, but there's still quite a lot of people
	// who have no access to dev-platform repo, so they cannot use the CRE CLI
	if !in.WorkflowConfig.ShouldCompileNewWorkflow && !in.WorkflowConfig.UseCRECLI {
		libcontracts.RegisterWorkflow(t, keystoneEnv.SethClient, keystoneEnv.KeystoneContractAddresses.WorkflowRegistryAddress, keystoneEnv.WorkflowDONID, workflowName, in.WorkflowConfig.CompiledWorkflowConfig.BinaryURL, in.WorkflowConfig.CompiledWorkflowConfig.ConfigURL)

		return
	}

	// These two env vars are required by the CRE CLI
	err := os.Setenv("WORKFLOW_OWNER_ADDRESS", keystoneEnv.SethClient.MustGetRootKeyAddress().Hex())
	require.NoError(t, err, "failed to set WORKFLOW_OWNER_ADDRESS env var")

	err = os.Setenv("ETH_PRIVATE_KEY", keystoneEnv.DeployerPrivateKey)
	require.NoError(t, err, "failed to set ETH_PRIVATE_KEY env var")

	// create CRE CLI settings file
	settingsFile := libcrecli.PrepareCRECLISettingsFile(t, keystoneEnv.SethClient.MustGetRootKeyAddress(), keystoneEnv.KeystoneContractAddresses.CapabilitiesRegistryAddress, keystoneEnv.KeystoneContractAddresses.WorkflowRegistryAddress, keystoneEnv.WorkflowDONID, keystoneEnv.ChainSelector, keystoneEnv.Blockchain.Nodes[0].HostHTTPUrl)

	var workflowURL string
	var workflowConfigURL string

	workflowConfigFile := keystoneporcrecli.CreateConfigFile(t, keystoneEnv.KeystoneContractAddresses.FeedsConsumerAddress, in.PriceProvider.FeedID, priceProvider.URL())

	// compile and upload the workflow, if we are not using an existing one
	if in.WorkflowConfig.ShouldCompileNewWorkflow {
		workflowURL, workflowConfigURL = libcrecli.CompileWorkflow(t, *in.WorkflowConfig.WorkflowFolderLocation, workflowConfigFile, settingsFile)
	} else {
		workflowURL = in.WorkflowConfig.CompiledWorkflowConfig.BinaryURL
		workflowConfigURL = in.WorkflowConfig.CompiledWorkflowConfig.ConfigURL
	}

	libcrecli.RegisterWorkflow(t, workflowName, workflowURL, workflowConfigURL, settingsFile)
}

func startNodeSets(t *testing.T, nsInputs []*keystonetypes.CapabilitiesAwareNodeSet, keystoneEnv *keystonetypes.KeystoneEnvironment) {
	require.NotNil(t, keystoneEnv, "keystone environment must be set")
	require.NotNil(t, keystoneEnv.Blockchain, "blockchain environment must be set")

	// Hack for CI that allows us to dynamically set the chainlink image and version
	// CTFv2 currently doesn't support dynamic image and version setting
	if os.Getenv("CI") == "true" {
		// Due to how we pass custom env vars to reusable workflow we need to use placeholders, so first we need to resolve what's the name of the target environment variable
		// that stores chainlink version and then we can use it to resolve the image name
		for i := range nsInputs {
			image := fmt.Sprintf("%s:%s", os.Getenv(ctfconfig.E2E_TEST_CHAINLINK_IMAGE_ENV), ctfconfig.MustReadEnvVar_String(ctfconfig.E2E_TEST_CHAINLINK_VERSION_ENV))
			for j := range nsInputs[i].NodeSpecs {
				nsInputs[i].NodeSpecs[j].Node.Image = image
			}
		}
	}

	for _, nsInput := range nsInputs {
		nodeset, err := ns.NewSharedDBNodeSet(nsInput.Input, keystoneEnv.Blockchain)
		require.NoError(t, err, "failed to deploy node set")

		keystoneEnv.NodeInput = append(keystoneEnv.NodeInput, nsInput)
		keystoneEnv.WrappedNodeOutput = append(keystoneEnv.WrappedNodeOutput, &keystonetypes.WrappedNodeOutput{
			Output:       nodeset,
			NodeSetName:  nsInput.Name,
			Capabilities: nsInput.Capabilities,
		})
	}
}

func logTestInfo(l zerolog.Logger, feedID, workflowName, feedConsumerAddr, forwarderAddr string) {
	l.Info().Msg("------ Test configuration:")
	l.Info().Msgf("Feed ID: %s", feedID)
	l.Info().Msgf("Workflow name: %s", workflowName)
	l.Info().Msgf("FeedConsumer address: %s", feedConsumerAddr)
	l.Info().Msgf("KeystoneForwarder address: %s", forwarderAddr)
}

func setupFakeDataProvider(t *testing.T, testLogger zerolog.Logger, in *TestConfig, priceIndex *int) string {
	_, err := fake.NewFakeDataProvider(in.PriceProvider.Fake.Input)
	require.NoError(t, err, "failed to set up fake data provider")
	fakeAPIPath := "/fake/api/price"
	host := framework.HostDockerInternal()
	fakeFinalURL := fmt.Sprintf("%s:%d%s", host, in.PriceProvider.Fake.Port, fakeAPIPath)

	getPriceResponseFn := func() map[string]interface{} {
		response := map[string]interface{}{
			"accountName": "TrueUSD",
			"totalTrust":  in.PriceProvider.Fake.Prices[*priceIndex],
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
	require.NoError(t, err, "failed to set up fake data provider")

	return fakeFinalURL
}

func setupPriceProvider(t *testing.T, testLogger zerolog.Logger, in *TestConfig) PriceProvider {
	if in.PriceProvider.Fake != nil {
		return NewFakePriceProvider(t, testLogger, in)
	}

	return NewLivePriceProvider(t, testLogger, in)
}

// PriceProvider abstracts away the logic of checking whether the feed has been correctly updated
// and it also returns port and URL of the price provider. This is so, because when using a mocked
// price provider we need start a separate service and whitelist its port and IP with the gateway job.
// Also, since it's a mocked price provider we can now check whether the feed has been correctly updated
// instead of only checking whether it has some price that's != 0.
type PriceProvider interface {
	URL() string
	NextPrice(price *big.Int, elapsed time.Duration) bool
	CheckPrices()
}

// LivePriceProvider is a PriceProvider implementation that uses a live feed to get the price, typically http://api.real-time-reserves.verinumus.io
type LivePriceProvider struct {
	t            *testing.T
	testLogger   zerolog.Logger
	url          string
	actualPrices []*big.Int
}

func NewLivePriceProvider(t *testing.T, testLogger zerolog.Logger, in *TestConfig) PriceProvider {
	return &LivePriceProvider{
		testLogger: testLogger,
		url:        in.PriceProvider.URL,
		t:          t,
	}
}

func (l *LivePriceProvider) NextPrice(price *big.Int, elapsed time.Duration) bool {
	// if price is nil or 0 it means that the feed hasn't been updated yet
	if price == nil || price.Cmp(big.NewInt(0)) == 0 {
		return true
	}

	l.testLogger.Info().Msgf("Feed updated after %s - price set, price=%s", elapsed, price)
	l.actualPrices = append(l.actualPrices, price)

	// no other price to return, we are done
	return false
}

func (l *LivePriceProvider) URL() string {
	return l.url
}

func (l *LivePriceProvider) CheckPrices() {
	// we don't have a way to check the price in the live feed, so we always assume it's correct
	// as long as it's != 0. And we only wait for the first price to be set.
	require.NotEmpty(l.t, l.actualPrices, "no prices found in the feed")
	require.NotEqual(l.t, l.actualPrices[0], big.NewInt(0), "price found in the feed is 0")
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

func NewFakePriceProvider(t *testing.T, testLogger zerolog.Logger, in *TestConfig) PriceProvider {
	priceIndex := ptr.Ptr(0)
	expectedPrices := make([]*big.Int, len(in.PriceProvider.Fake.Prices))
	for i, p := range in.PriceProvider.Fake.Prices {
		// convert float64 to big.Int by multiplying by 100
		// just like the PoR workflow does
		expectedPrices[i] = libc.Float64ToBigInt(p)
	}

	return &FakePriceProvider{
		t:              t,
		testLogger:     testLogger,
		expectedPrices: expectedPrices,
		priceIndex:     priceIndex,
		url:            setupFakeDataProvider(t, testLogger, in, priceIndex),
	}
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

func (f *FakePriceProvider) CheckPrices() {
	require.EqualValues(f.t, f.expectedPrices, f.actualPrices, "prices found in the feed do not match prices set in the mock")
	f.testLogger.Info().Msgf("All %d mocked prices were found in the feed", len(f.expectedPrices))
}

func (f *FakePriceProvider) URL() string {
	return f.url
}

func startBlockchain(t *testing.T, in *TestConfig, keystoneEnv *keystonetypes.KeystoneEnvironment) {
	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err, "failed to create blockchain network")

	pkey := os.Getenv("PRIVATE_KEY")
	require.NotEmpty(t, pkey, "private key must not be empty")

	sc, err := seth.NewClientBuilder().
		WithRpcUrl(bc.Nodes[0].HostWSUrl).
		WithPrivateKeys([]string{pkey}).
		Build()
	require.NoError(t, err, "failed to create seth client")

	chainSelector, err := chainselectors.SelectorFromChainId(sc.Cfg.Network.ChainID)
	require.NoError(t, err, "failed to get chain selector for chain id %d", sc.Cfg.Network.ChainID)

	keystoneEnv.Blockchain = bc
	keystoneEnv.SethClient = sc
	keystoneEnv.DeployerPrivateKey = pkey
	keystoneEnv.ChainSelector = chainSelector
}

func extraAllowedPortsAndIps(t *testing.T, testLogger zerolog.Logger, in *TestConfig, nodeOutput *ns.Output) ([]string, []int) {
	// no need to allow anything, if we are using live feed
	if in.PriceProvider.Fake == nil {
		return nil, nil
	}

	// we need to explicitly allow the port used by the fake data provider
	// and IP corresponding to host.docker.internal or the IP of the host machine, if we are running on Linux,
	// because that's where the fake data provider is running
	var hostIP string
	var err error

	system := runtime.GOOS
	switch system {
	case "darwin":
		hostIP, err = libdon.ResolveHostDockerInternaIP(testLogger, nodeOutput)
		require.NoError(t, err, "failed to resolve host.docker.internal IP")
	case "linux":
		// for linux framework already returns an IP, so we don't need to resolve it,
		// but we need to remove the http:// prefix
		hostIP = strings.ReplaceAll(framework.HostDockerInternal(), "http://", "")
	default:
		err = fmt.Errorf("unsupported OS: %s", system)
	}
	require.NoError(t, err, "failed to resolve host.docker.internal IP")

	testLogger.Info().Msgf("Will allow IP %s and port %d for the fake data provider", hostIP, in.PriceProvider.Fake.Port)

	ips, err := net.LookupIP("gist.githubusercontent.com")
	require.NoError(t, err, "failed to resolve IP for gist.githubusercontent.com")

	gistIPs := make([]string, len(ips))
	for i, ip := range ips {
		gistIPs[i] = ip.To4().String()
		testLogger.Debug().Msgf("Resolved IP for gist.githubusercontent.com: %s", gistIPs[i])
	}

	// we also need to explicitly allow Gist's IP
	return append(gistIPs, hostIP), []int{in.PriceProvider.Fake.Port}
}

// TODO think whether we should structure it in a way that envforces some order of execution,
// for example by making the outputs of one function inputs to another
func prepareTestEnvironment(t *testing.T, testLogger zerolog.Logger, in *TestConfig) (*keystonetypes.KeystoneEnvironment, PriceProvider) {
	keystoneEnv := &keystonetypes.KeystoneEnvironment{}
	keystoneEnv.GatewayConnectorData = &keystonetypes.GatewayConnectorData{
		Path: "/node",
		Port: 5003,
	}

	// Create a new blockchain network and Seth client to interact with it
	startBlockchain(t, in, keystoneEnv)

	// Get either a no-op price provider (for live endpoint)
	// or a fake price provider (for mock endpoint)
	priceProvider := setupPriceProvider(t, testLogger, in)

	// Start job distributor
	libjobs.StartJobDistributor(t, in.JD, keystoneEnv)

	// Deploy the DONs
	startNodeSets(t, in.NodeSets, keystoneEnv)

	// Prepare the CLD environment and figure out DON topology; configure chains for nodes and job distributor
	keystoneenv.BuildTopologyAndCLDEnvironment(t, keystoneEnv)

	// Fund the nodes
	libdon.FundNodes(t, keystoneEnv)

	// Deploy keystone contracts (forwarder, capability registry, ocr3 capability, workflow registry)
	libcontracts.DeployKeystone(t, testLogger, keystoneEnv)

	// Separated from Keystone deployment because it will soon be replaced with DF Cache
	libcontracts.DeployFeedsConsumer(t, testLogger, keystoneEnv)

	// Configure Workflow Registry and Feeds Consumer
	libcontracts.ConfigureWorkflowRegistry(t, testLogger, keystoneEnv)
	libcontracts.ConfigureFeedsConsumer(t, testLogger, in.WorkflowConfig.WorkflowName, keystoneEnv)

	// Register the workflow (either via CRE CLI or by calling the workflow registry directly; using only workflow DON id)
	registerPoRWorkflow(t, in, in.WorkflowConfig.WorkflowName, keystoneEnv, priceProvider)

	donToConfigs, donToJobSpecs := prepareJobSpecsAndNodeConfigs(t, testLogger, in, keystoneEnv)
	libdon.Configure(t, testLogger, keystoneEnv, donToJobSpecs, donToConfigs)

	// CAUTION: It is crucial to configure OCR3 jobs on nodes before configuring the workflow contracts.
	// Wait for OCR listeners to be ready before setting the configuration.
	// If the ConfigSet event is missed, OCR protocol will not start.
	// TODO make it fluent!
	testLogger.Info().Msg("Waiting 30s for OCR listeners to be ready...")
	time.Sleep(30 * time.Second)
	testLogger.Info().Msg("Proceeding to set OCR3 configuration.")

	// Configure the Forwarder, OCR3 and Capabilities contracts
	libcontracts.ConfigureKeystone(t, keystoneEnv)

	return keystoneEnv, priceProvider
}

func prepareJobSpecsAndNodeConfigs(t *testing.T, testLogger zerolog.Logger, in *TestConfig, keystoneEnv *keystonetypes.KeystoneEnvironment) (keystonetypes.DonsToConfigOverrides, map[uint32]keystonetypes.DonJobs) {
	ips, ports := extraAllowedPortsAndIps(t, testLogger, in, keystoneEnv.WrappedNodeOutput[0].Output)

	peeringData, err := libdon.FindPeeringData(keystoneEnv.DONTopology)
	require.NoError(t, err, "failed to get peering data")

	donToConfigs := make(keystonetypes.DonsToConfigOverrides)
	for _, donTopology := range keystoneEnv.DONTopology {
		donToConfigs[donTopology.ID] = keystoneporconfig.Define(t,
			donTopology.DON,
			donTopology.NodeInput,
			donTopology.NodeOutput,
			keystoneEnv.Blockchain,
			donTopology.ID,
			donTopology.Flags,
			peeringData,
			keystoneEnv.KeystoneContractAddresses.CapabilitiesRegistryAddress,
			keystoneEnv.KeystoneContractAddresses.WorkflowRegistryAddress,
			keystoneEnv.KeystoneContractAddresses.ForwarderAddress,
			keystoneEnv.GatewayConnectorData,
		)
	}

	// define jobs
	donToJobSpecs := make(map[uint32]keystonetypes.DonJobs)
	for _, donTopology := range keystoneEnv.DONTopology {
		jobSpecs := keystonepor.Define(t,
			keystoneEnv.Environment,
			donTopology.DON,
			donTopology.NodeOutput,
			keystoneEnv.Blockchain,
			keystoneEnv.KeystoneContractAddresses.OCR3CapabilityAddress,
			donTopology.ID,
			donTopology.Flags,
			ports,
			ips,
			cronCapabilityAssetFile,
			*keystoneEnv.GatewayConnectorData,
		)
		donToJobSpecs[donTopology.ID] = jobSpecs
	}

	return donToConfigs, donToJobSpecs
}
func TestKeystoneWithOCR3Workflow(t *testing.T) {
	testLogger := framework.L

	// Load test configuration
	in, err := framework.Load[TestConfig](t)
	require.NoError(t, err, "couldn't load test config")
	validateInputsAndEnvVars(t, in)

	keystoneEnv, priceProvider := prepareTestEnvironment(t, testLogger, in)

	// Log extra information that might help debugging
	t.Cleanup(func() {
		if t.Failed() {
			logTestInfo(testLogger, in.PriceProvider.FeedID, in.WorkflowConfig.WorkflowName, keystoneEnv.KeystoneContractAddresses.FeedsConsumerAddress.Hex(), keystoneEnv.KeystoneContractAddresses.ForwarderAddress.Hex())

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

			lidebug.PrintTestDebug(t, testLogger, keystoneEnv)
		}
	})

	// It can take a while before the first report is produced, particularly on CI.
	timeout := 10 * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	feedsConsumerInstance, err := feeds_consumer.NewKeystoneFeedsConsumer(keystoneEnv.KeystoneContractAddresses.FeedsConsumerAddress, keystoneEnv.SethClient.Client)
	require.NoError(t, err, "failed to create feeds consumer instance")

	testLogger.Info().Msg("Waiting for feed to update...")
	startTime := time.Now()
	feedBytes := common.HexToHash(in.PriceProvider.FeedID)

	for {
		select {
		case <-ctx.Done():
			testLogger.Error().Msgf("feed did not update, timeout after %s", timeout)
			t.FailNow()
		case <-time.After(10 * time.Second):
			elapsed := time.Since(startTime).Round(time.Second)
			price, _, err := feedsConsumerInstance.GetPrice(
				keystoneEnv.SethClient.NewCallOpts(),
				feedBytes,
			)
			require.NoError(t, err, "failed to get price from Keystone Consumer contract")

			if !priceProvider.NextPrice(price, elapsed) {
				// check if all expected prices were found and finish the test
				priceProvider.CheckPrices()
				return
			}
			testLogger.Info().Msgf("Feed not updated yet, waiting for %s", elapsed)
		}
	}
}
