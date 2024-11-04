package test_env

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	tc "github.com/testcontainers/testcontainers-go"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env/job_distributor"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logstream"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/runid"

	"github.com/smartcontractkit/chainlink/integration-tests/testconfig/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"

	d "github.com/smartcontractkit/chainlink/integration-tests/docker"
)

var (
	ErrFundCLNode = "failed to fund CL node"
)

type ClusterTestEnv struct {
	Cfg           *TestEnvConfig
	DockerNetwork *tc.DockerNetwork
	LogStream     *logstream.LogStream
	TestConfig    ctf_config.GlobalTestConfig

	/* components */
	ClCluster              *ClCluster
	RmnCluster             *RmnCluster
	MockAdapter            *test_env.Killgrave
	PrivateEthereumConfigs []*ctf_config.EthereumNetworkConfig
	EVMNetworks            []*blockchain.EVMNetwork
	rpcProviders           map[int64]*test_env.RpcProvider
	JobDistributor         *job_distributor.Component
	l                      zerolog.Logger
	t                      *testing.T
	isSimulatedNetwork     bool
}

func NewTestEnv() (*ClusterTestEnv, error) {
	log.Logger = logging.GetLogger(nil, "CORE_DOCKER_ENV_LOG_LEVEL")
	network, err := docker.CreateNetwork(log.Logger)
	if err != nil {
		return nil, err
	}
	return &ClusterTestEnv{
		DockerNetwork: network,
		l:             log.Logger,
		ClCluster:     &ClCluster{},
		RmnCluster:    &RmnCluster{},
	}, nil
}

// WithTestEnvConfig sets the test environment cfg.
// Sets up private ethereum chain and MockAdapter containers with the provided cfg.
func (te *ClusterTestEnv) WithTestEnvConfig(cfg *TestEnvConfig) *ClusterTestEnv {
	te.Cfg = cfg
	if cfg.MockAdapter.ContainerName != "" {
		n := []string{te.DockerNetwork.Name}
		te.MockAdapter = test_env.NewKillgrave(n, te.Cfg.MockAdapter.ImpostersPath, test_env.WithContainerName(te.Cfg.MockAdapter.ContainerName), test_env.WithLogStream(te.LogStream))
	}
	return te
}

func (te *ClusterTestEnv) WithTestInstance(t *testing.T) *ClusterTestEnv {
	te.t = t
	te.l = logging.GetTestLogger(t)
	if te.MockAdapter != nil {
		te.MockAdapter.WithTestInstance(t)
	}
	return te
}

func (te *ClusterTestEnv) StartEthereumNetwork(cfg *ctf_config.EthereumNetworkConfig) (blockchain.EVMNetwork, test_env.RpcProvider, error) {
	// if environment is being restored from a previous state, use the existing config
	// this might fail terribly if temporary folders with chain data on the host machine were removed
	if te.Cfg != nil && te.Cfg.EthereumNetworkConfig != nil {
		cfg = te.Cfg.EthereumNetworkConfig
	}

	te.l.Info().
		Str("Execution Layer", string(*cfg.ExecutionLayer)).
		Str("Ethereum Version", string(*cfg.EthereumVersion)).
		Str("Custom Docker Images", fmt.Sprintf("%v", cfg.CustomDockerImages)).
		Msg("Starting Ethereum network")

	builder := test_env.NewEthereumNetworkBuilder()
	c, err := builder.WithExistingConfig(*cfg).
		WithTest(te.t).
		WithLogStream(te.LogStream).
		Build()
	if err != nil {
		return blockchain.EVMNetwork{}, test_env.RpcProvider{}, err
	}

	n, rpc, err := c.Start()

	if err != nil {
		return blockchain.EVMNetwork{}, test_env.RpcProvider{}, err
	}

	return n, rpc, nil
}

func (te *ClusterTestEnv) StartJobDistributor(cfg *ccip.JDConfig) error {
	jdDB, err := test_env.NewPostgresDb(
		[]string{te.DockerNetwork.Name},
		test_env.WithPostgresDbName(cfg.GetJDDBName()),
		test_env.WithPostgresImageVersion(cfg.GetJDDBVersion()),
	)
	if err != nil {
		return fmt.Errorf("failed to create postgres db for job-distributor: %w", err)
	}
	err = jdDB.StartContainer()
	if err != nil {
		return fmt.Errorf("failed to start postgres db for job-distributor: %w", err)
	}
	jd := job_distributor.New([]string{te.DockerNetwork.Name},
		job_distributor.WithImage(cfg.GetJDImage()),
		job_distributor.WithVersion(cfg.GetJDVersion()),
		job_distributor.WithDBURL(jdDB.InternalURL.String()),
	)
	jd.LogStream = te.LogStream
	err = jd.StartContainer()
	if err != nil {
		return fmt.Errorf("failed to start job-distributor: %w", err)
	}
	te.JobDistributor = jd
	return nil
}

func (te *ClusterTestEnv) StartMockAdapter() error {
	return te.MockAdapter.StartContainer()
}

func (te *ClusterTestEnv) StartClCluster(nodeConfig *chainlink.Config, count int, secretsConfig string, testconfig ctf_config.GlobalTestConfig, opts ...ClNodeOption) error {
	if te.Cfg != nil && te.Cfg.ClCluster != nil {
		te.ClCluster = te.Cfg.ClCluster
	} else {
		// prepend the postgres version option from the toml config
		if testconfig.GetChainlinkImageConfig().PostgresVersion != nil && *testconfig.GetChainlinkImageConfig().PostgresVersion != "" {
			opts = append([]func(c *ClNode){
				func(c *ClNode) {
					c.PostgresDb.EnvComponent.ContainerVersion = *testconfig.GetChainlinkImageConfig().PostgresVersion
				},
			}, opts...)
		}
		opts = append(opts, WithSecrets(secretsConfig))
		te.ClCluster = &ClCluster{}
		for i := 0; i < count; i++ {
			ocrNode, err := NewClNode([]string{te.DockerNetwork.Name}, *testconfig.GetChainlinkImageConfig().Image, *testconfig.GetChainlinkImageConfig().Version, nodeConfig, te.LogStream, opts...)
			if err != nil {
				return err
			}
			te.ClCluster.Nodes = append(te.ClCluster.Nodes, ocrNode)
		}
	}

	// Set test logger
	if te.t != nil {
		for _, n := range te.ClCluster.Nodes {
			n.SetTestLogger(te.t)
		}
	}

	// Start/attach node containers
	te.l.Info().Msgf("Starting CL cluster with %d nodes", count)
	return te.ClCluster.Start()
}

// Starts a RMN cluster in the context of the coreqa tool. There's limited value
// to this, as the RMN nodes will likely enter a reset loop as they expect a very particular
// set of conditions, such as an already deployed RMN contract on chain, and preexisting, correct
// shared and local chain config. Still, it may be useful to maintain them in that reset loop
// until the deployment is done manually (in that case, you should ensure to use a docker image
// that keeps the RMN binary restarting until these conditions are established.)
func (te *ClusterTestEnv) StartRmnCluster(count int) error {
	sharedConfig := RmnNodeSharedConfig{}
	localConfig := RmnNodeLocalConfig{}

	// TODO: This configuration is not too meaningful; it's
	// just cobbled together to get the RMN cluster as far as possible without
	// actual contracts deployed on-chain.
	for i := 0; i < len(te.EVMNetworks); i++ {
		net := te.EVMNetworks[i]

		// Kind of hacky, but we override the net name to match Balthazar's expectation based on the
		// chain ID for the specific test IDs used. TODO: properly match Balthazar's expectations.
		name := net.Name
		switch net.ChainID {
		case 1337:
			name = "DevnetAlpha"
		case 2337:
			name = "DevnetBeta"
		case 3337:
			name = "DevnetGamma"
		}

		sharedChain := SharedChain{
			Name:                         name,
			MaxTaggedRootsPerVoteToBless: 10,
			AfnType:                      "V1_0",
			// Must be replaced with the AFN contract address on this chain once deployed. This is a sepolia stub
			AfnContract:      "0xb4d360459f32dd641ef5a6985ffbac5c4e5521aa",
			InflightTime:     Duration{Minutes: 5},
			MaxFreshBlockAge: Duration{Seconds: 24},
			UponFinalityViolationVoteToCurseOnOtherChainsWithLegacyContracts: true,
			Stability: StabilityConfig{
				Type:              "FinalityTag",
				SoftConfirmations: 0,
			},
			BlessFeeConfig: FeeConfig{
				Type: "Eip1559",
				MaxFeePerGas: &Gwei{
					Gwei: 1000,
				},
				MaxPriorityFeePerGas: &Gwei{
					Gwei: 50,
				},
			},
			CurseFeeConfig: FeeConfig{
				Type: "Eip1559",
				MaxFeePerGas: &Gwei{
					Gwei: 2000,
				},
				MaxPriorityFeePerGas: &Gwei{
					Gwei: 200,
				},
			},
		}

		sharedConfig.Chains = append(sharedConfig.Chains, sharedChain)
		localChain := Chain{
			Name: name,
			RPCS: te.rpcProviders[net.ChainID].PrivateHttpUrls(),
		}
		localConfig.Chains = append(localConfig.Chains, localChain)
		// No lanes since this is just a barebones test environment. Lane config must
		// be added manually.
	}

	for i := 0; i < count; i++ {
		nodeDefaultCName := fmt.Sprintf("%s-%s", "rmn-node", uuid.NewString()[0:8])
		// TODO: make the image name/version configurable
		rmnNode, err := NewRmnNode([]string{te.DockerNetwork.Name}, nodeDefaultCName, "rmn", "latest", sharedConfig, localConfig)
		if err != nil {
			return err
		}
		te.RmnCluster.Nodes = append(te.RmnCluster.Nodes, rmnNode)
	}
	te.l.Info().Msgf("Starting RMN cluster with %d nodes", count)
	return te.RmnCluster.Start()
}

func (te *ClusterTestEnv) Terminate() error {
	// TESTCONTAINERS_RYUK_DISABLED=false by default so ryuk will remove all
	// the containers and the Network
	return nil
}

type CleanupOpts struct {
	TestName string
}

// Cleanup cleans the environment up after it's done being used, mainly for returning funds when on live networks and logs.
func (te *ClusterTestEnv) Cleanup(opts CleanupOpts) error {
	te.l.Info().Msg("Cleaning up test environment")

	runIdErr := runid.RemoveLocalRunId(te.TestConfig.GetLoggingConfig().RunId)
	if runIdErr != nil {
		te.l.Warn().Msgf("Failed to remove .run.id file due to: %s (not a big deal, you can still remove it manually)", runIdErr.Error())
	}

	if te.t == nil {
		return fmt.Errorf("cannot cleanup test environment without a testing.T")
	}

	if te.ClCluster == nil || len(te.ClCluster.Nodes) == 0 {
		return fmt.Errorf("chainlink nodes are nil, unable cleanup chainlink nodes")
	}

	te.logWhetherAllContainersAreRunning()

	err := te.handleNodeCoverageReports(opts.TestName)
	if err != nil {
		te.l.Error().Err(err).Msg("Error handling node coverage reports")
	}

	return nil
}

// handleNodeCoverageReports handles the coverage reports for the chainlink nodes
func (te *ClusterTestEnv) handleNodeCoverageReports(testName string) error {
	testName = strings.ReplaceAll(testName, "/", "_")
	showHTMLCoverageReport := te.TestConfig.GetLoggingConfig().ShowHTMLCoverageReport != nil && *te.TestConfig.GetLoggingConfig().ShowHTMLCoverageReport
	isCI := os.Getenv("CI") != ""

	te.l.Info().
		Bool("showCoverageReportFlag", showHTMLCoverageReport).
		Bool("isCI", isCI).
		Bool("show", showHTMLCoverageReport || isCI).
		Msg("Checking if coverage report should be shown")

	var covHelper *d.NodeCoverageHelper

	if showHTMLCoverageReport || isCI {
		// Stop all nodes in the chainlink cluster.
		// This is needed to get go coverage profile from the node containers https://go.dev/doc/build-cover#FAQ
		// TODO: fix this as it results in: ERR LOG AFTER TEST ENDED ... INF 🐳 Stopping container
		err := te.ClCluster.Stop()
		if err != nil {
			return err
		}

		clDir, err := getChainlinkDir()
		if err != nil {
			return err
		}

		var coverageRootDir string
		if os.Getenv("GO_COVERAGE_DEST_DIR") != "" {
			coverageRootDir = filepath.Join(os.Getenv("GO_COVERAGE_DEST_DIR"), testName)
		} else {
			coverageRootDir = filepath.Join(clDir, ".covdata", testName)
		}

		var containers []tc.Container
		for _, node := range te.ClCluster.Nodes {
			containers = append(containers, node.Container)
		}

		covHelper, err = d.NewNodeCoverageHelper(context.Background(), containers, clDir, coverageRootDir)
		if err != nil {
			return err
		}
	}

	// Show html coverage report when flag is set (local runs)
	if showHTMLCoverageReport {
		path, err := covHelper.SaveMergedHTMLReport()
		if err != nil {
			return err
		}
		te.l.Info().Str("testName", testName).Str("filePath", path).Msg("Chainlink node coverage html report saved")
	}

	// Save percentage coverage report when running in CI
	if isCI {
		// Save coverage percentage to a file to show in the CI
		path, err := covHelper.SaveMergedCoveragePercentage()
		if err != nil {
			te.l.Error().Err(err).Str("testName", testName).Msg("Failed to save coverage percentage for test")
		} else {
			te.l.Info().Str("testName", testName).Str("filePath", path).Msg("Chainlink node coverage percentage report saved")
		}
	}

	return nil
}

// getChainlinkDir returns the path to the chainlink directory
func getChainlinkDir() (string, error) {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return "", fmt.Errorf("cannot determine the path of the calling file")
	}
	dir := filepath.Dir(filename)
	chainlinkDir := filepath.Clean(filepath.Join(dir, "../../.."))
	return chainlinkDir, nil
}

func (te *ClusterTestEnv) logWhetherAllContainersAreRunning() {
	for _, node := range te.ClCluster.Nodes {
		if node.Container == nil {
			continue
		}

		isCLRunning := node.Container.IsRunning()
		isDBRunning := node.PostgresDb.Container.IsRunning()

		if !isCLRunning {
			te.l.Warn().Str("Node", node.ContainerName).Msg("Chainlink node was not running, when test ended")
		}

		if !isDBRunning {
			te.l.Warn().Str("Node", node.ContainerName).Msg("Postgres DB is not running, when test ended")
		}
	}
}

func (te *ClusterTestEnv) GetRpcProvider(chainId int64) (*test_env.RpcProvider, error) {
	if rpc, ok := te.rpcProviders[chainId]; ok {
		return rpc, nil
	}

	return nil, fmt.Errorf("no RPC provider available for chain ID %d", chainId)
}

func (te *ClusterTestEnv) GetFirstEvmNetwork() (*blockchain.EVMNetwork, error) {
	if len(te.EVMNetworks) == 0 {
		return nil, fmt.Errorf("no EVM networks available")
	}

	return te.EVMNetworks[0], nil
}

func (te *ClusterTestEnv) GetEVMNetworkForChainId(chainId int64) (*blockchain.EVMNetwork, error) {
	for _, network := range te.EVMNetworks {
		if network.ChainID == chainId {
			return network, nil
		}
	}

	return nil, fmt.Errorf("no EVM network available for chain ID %d", chainId)
}
