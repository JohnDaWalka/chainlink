package smoke

import (
	"context"
	"fmt"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"google.golang.org/grpc/credentials/insecure"
	"math/big"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/subosito/gotenv"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	ctfv2_blockchain "github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	ctfconfig "github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	ctf_testenv "github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env/job_distributor"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logstream"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	ccipactions "github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/testconfig/ccip"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type CTFV2Config struct {
	BlockchainNetworks []*ctfv2_blockchain.Input `toml:"Networks" validate:"required"`
	//ChainA  *ctfv2_blockchain.Input `toml:"SIMULATED_1" validate:"required"`
	//ChainB  *ctfv2_blockchain.Input `toml:"SIMULATED_2" validate:"required"`
	//ChainC  *ctfv2_blockchain.Input `toml:"SIMULATED_3" validate:"required"`
	NodeSet *ns.Input `toml:"nodeset" validate:"required"`
}

func TestSmoke(t *testing.T) {
	lggr := logger.TestLogger(t)
	in, err := framework.Load[CTFV2Config](t)
	require.NoError(t, err)
	tCfg := &changeset.TestConfigs{}

	if _, err := os.Stat(".env"); err == nil || !os.IsNotExist(err) {
		require.NoError(t, gotenv.Load(".env"), "Error loading .env file")
	}
	testEnv, err := test_env.NewTestEnv()
	require.NoError(t, err, "Error creating test env")
	cfg, err := tc.GetChainAndTestTypeSpecificConfig("Smoke", tc.CCIP)
	loggingCfg := cfg.GetLoggingConfig()
	testEnv.LogStream, err = logstream.NewLogStream(t, loggingCfg)
	require.NoError(t, err, "failed to create logstream ")

	require.NoError(t, err, "Error getting config")

	testEnv.MockAdapter = ctf_testenv.NewKillgrave([]string{testEnv.DockerNetwork.Name}, "", ctf_testenv.WithLogStream(testEnv.LogStream))

	err = testEnv.StartMockAdapter()
	require.NoError(t, err, "failed to create killgrave instance ")

	jdDB, err := ctf_testenv.NewPostgresDb(
		[]string{testEnv.DockerNetwork.Name},
		ctf_testenv.WithPostgresDbName(ccip.DEFAULT_DB_NAME),
		ctf_testenv.WithPostgresImageVersion(ccip.DEFAULT_DB_VERSION),
	)
	require.NoError(t, err, "failed to create postgres db for job-distributor")
	err = jdDB.StartContainer()
	require.NoError(t, err, "failed to start postgres db for job-distributor")
	//time.Sleep(4*time.Second)
	jd := job_distributor.New([]string{testEnv.DockerNetwork.Name},
		job_distributor.WithImage(ctfconfig.MustReadEnvVar_String(ccip.E2E_JD_IMAGE)),
		job_distributor.WithVersion(ctfconfig.MustReadEnvVar_String(ccip.E2E_JD_VERSION)),
		job_distributor.WithDBURL(jdDB.InternalURL.String()),
	)
	jd.LogStream = testEnv.LogStream
	fmt.Println("Starting job-distributor container with DB URL:", jdDB.InternalURL.String())
	err = jd.StartContainer()
	require.NoError(t, err, "failed to start job-distributor")
	testEnv.JobDistributor = jd

	//evmNetworks := networks.MustGetSelectedNetworkConfig(cfg.GetNetworkConfig())

	var blockchains []*ctfv2_blockchain.Output
	for _, bcInput := range in.BlockchainNetworks {
		bc, err := ctfv2_blockchain.NewBlockchainNetwork(bcInput)
		require.NoError(t, err)
		blockchains = append(blockchains, bc)
	}

	jdConfig := devenv.JDConfig{
		GRPC:  cfg.CCIP.JobDistributorConfig.GetJDGRPC(),
		WSRPC: cfg.CCIP.JobDistributorConfig.GetJDWSRPC(),
	}
	// TODO : move this as a part of test_env setup with an input in testconfig
	// if JD is not provided, we will spin up a new JD
	if jdConfig.GRPC == "" || jdConfig.WSRPC == "" {
		jdConfig = devenv.JDConfig{
			GRPC: jd.Grpc,
			// we will use internal wsrpc for nodes on same docker network to connect to JD
			WSRPC: jd.InternalWSRPC,
			Creds: insecure.NewCredentials(),
		}
	}
	require.NotEmpty(t, jdConfig, "JD config is empty")
	evmNetworks := networks.MustGetSelectedNetworkConfig(cfg.GetNetworkConfig())
	networkPvtKeys := make(map[int64]string)
	networkName := make(map[int64]string)
	for _, net := range evmNetworks {
		require.Greater(t, len(net.PrivateKeys), 0, "No private keys found for network")
		networkPvtKeys[net.ChainID] = net.PrivateKeys[0]
		networkName[net.ChainID] = net.Name
	}
	var chains []devenv.ChainConfig
	for _, chain := range blockchains {
		chainID, err := strconv.ParseInt(chain.ChainID, 10, 64)
		require.NoError(t, err, "invalid chain id")
		pvtKeyStr, exists := networkPvtKeys[chainID]
		require.Truef(t, exists, "Private key not found for chain id %d", chainID)
		pvtKey, err := crypto.HexToECDSA(pvtKeyStr)
		require.NoError(t, err)
		deployer, err := bind.NewKeyedTransactorWithChainID(pvtKey, big.NewInt(chainID))
		require.NoError(t, err)
		chains = append(chains, devenv.ChainConfig{
			ChainID:     uint64(chainID),
			ChainName:   networkName[chainID],
			ChainType:   devenv.EVMChainType,
			WSRPCs:      []string{chain.Nodes[0].HostWSUrl},
			HTTPRPCs:    []string{chain.Nodes[0].HostHTTPUrl},
			DeployerKey: deployer,
		})
	}
	homeChainSelector, err := cfg.CCIP.GetHomeChainSelector(evmNetworks)
	require.NoError(t, err, "Error getting home chain selector")
	feedChainSelector, err := cfg.CCIP.GetFeedChainSelector(evmNetworks)
	require.NoError(t, err, "Error getting feed chain selector")
	envConfig := devenv.EnvironmentConfig{
		Chains:            chains,
		JDConfig:          jdConfig,
		HomeChainSelector: homeChainSelector,
		FeedChainSelector: feedChainSelector,
	}

	deploymentChains, err := devenv.NewChains(lggr, envConfig.Chains)
	require.NoError(t, err)
	// locate the home chain
	homeChainSel := envConfig.HomeChainSelector
	require.NotEmpty(t, homeChainSel, "homeChainSel should not be empty")
	feedSel := envConfig.FeedChainSelector
	require.NotEmpty(t, feedSel, "feedSel should not be empty")
	replayBlocks, err := changeset.LatestBlocksByChain(testcontext.Get(t), deploymentChains)
	require.NoError(t, err)
	_ = replayBlocks
	ab := deployment.NewMemoryAddressBook()
	crConfig := changeset.DeployTestContracts(t, lggr, ab, homeChainSel, feedSel, deploymentChains, changeset.MockLinkPrice, changeset.MockWethPrice)

	// start the chainlink nodes with the CR address
	err = testsetups.StartChainlinkNodes(t, &envConfig, crConfig, testEnv, cfg)
	require.NoError(t, err)

	e, don, err := devenv.NewEnvironment(func() context.Context { return testcontext.Get(t) }, lggr, envConfig)
	require.NoError(t, err)
	require.NotNil(t, e)
	e.ExistingAddresses = ab

	// fund the nodes
	zeroLogLggr := logging.GetTestLogger(t)
	testsetups.FundNodes(t, zeroLogLggr, testEnv, cfg, don.PluginNodes())

	env := *e
	envNodes, err := deployment.NodeInfo(env.NodeIDs, env.Offchain)
	require.NoError(t, err)
	allChains := env.AllChainSelectors()
	var usdcChains []uint64
	if tCfg.IsUSDC {
		usdcChains = allChains
	}
	mcmsCfgPerChain := commontypes.MCMSWithTimelockConfig{
		Canceller:         commonchangeset.SingleGroupMCMS(t),
		Bypasser:          commonchangeset.SingleGroupMCMS(t),
		Proposer:          commonchangeset.SingleGroupMCMS(t),
		TimelockExecutors: env.AllDeployerKeys(),
		TimelockMinDelay:  big.NewInt(0),
	}
	mcmsCfg := make(map[uint64]commontypes.MCMSWithTimelockConfig)
	for _, c := range env.AllChainSelectors() {
		mcmsCfg[c] = mcmsCfgPerChain
	}
	// Need to deploy prerequisites first so that we can form the USDC config
	// no proposals to be made, timelock can be passed as nil here
	env, err = commonchangeset.ApplyChangesets(t, env, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.DeployHomeChain),
			Config: changeset.DeployHomeChainConfig{
				HomeChainSel:     homeChainSel,
				RMNStaticConfig:  changeset.NewTestRMNStaticConfig(),
				RMNDynamicConfig: changeset.NewTestRMNDynamicConfig(),
				NodeOperators:    changeset.NewTestNodeOperator(chains[homeChainSel].DeployerKey.From),
				NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
					"NodeOperator": envNodes.NonBootstraps().PeerIDs(),
				},
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.DeployPrerequisites),
			Config: changeset.DeployPrerequisiteConfig{
				ChainSelectors: allChains,
				Opts: []changeset.PrerequisiteOpt{
					changeset.WithUSDCChains(usdcChains),
					changeset.WithMulticall3(tCfg.IsMultiCall3),
				},
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployMCMSWithTimelock),
			Config:    mcmsCfg,
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.DeployChainContracts),
			Config: changeset.DeployChainContractsConfig{
				ChainSelectors:    allChains,
				HomeChainSelector: homeChainSel,
			},
		},
	})
	require.NoError(t, err)

	state, err := changeset.LoadOnchainState(env)
	require.NoError(t, err)
	tokenConfig := changeset.NewTestTokenConfig(state.Chains[feedSel].USDFeeds)
	usdcCCTPConfig := make(map[cciptypes.ChainSelector]pluginconfig.USDCCCTPTokenConfig)
	timelocksPerChain := make(map[uint64]*gethwrappers.RBACTimelock)
	ocrParams := make(map[uint64]changeset.CCIPOCRParams)
	for _, chain := range usdcChains {
		require.NotNil(t, state.Chains[chain].MockUSDCTokenMessenger)
		require.NotNil(t, state.Chains[chain].MockUSDCTransmitter)
		require.NotNil(t, state.Chains[chain].USDCTokenPool)
		usdcCCTPConfig[cciptypes.ChainSelector(chain)] = pluginconfig.USDCCCTPTokenConfig{
			SourcePoolAddress:            state.Chains[chain].USDCTokenPool.Address().String(),
			SourceMessageTransmitterAddr: state.Chains[chain].MockUSDCTransmitter.Address().String(),
		}
	}
	var usdcAttestationCfg changeset.USDCAttestationConfig
	if len(usdcChains) > 0 {
		var endpoint string
		err = ccipactions.SetMockServerWithUSDCAttestation(testEnv.MockAdapter, nil)
		require.NoError(t, err)
		endpoint = testEnv.MockAdapter.InternalEndpoint
		usdcAttestationCfg = changeset.USDCAttestationConfig{
			API:         endpoint,
			APITimeout:  commonconfig.MustNewDuration(time.Second),
			APIInterval: commonconfig.MustNewDuration(500 * time.Millisecond),
		}
	}
	require.NotNil(t, state.Chains[feedSel].LinkToken)
	require.NotNil(t, state.Chains[feedSel].Weth9)

	for _, chain := range allChains {
		timelocksPerChain[chain] = state.Chains[chain].Timelock
		tokenInfo := tokenConfig.GetTokenInfo(env.Logger, state.Chains[chain].LinkToken, state.Chains[chain].Weth9)
		ocrParams[chain] = changeset.DefaultOCRParams(feedSel, tokenInfo)
	}
	// Deploy second set of changesets to deploy and configure the CCIP contracts.
	env, err = commonchangeset.ApplyChangesets(t, env, timelocksPerChain, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.ConfigureNewChains),
			Config: changeset.NewChainsConfig{
				HomeChainSel:   homeChainSel,
				FeedChainSel:   feedSel,
				ChainsToDeploy: allChains,
				TokenConfig:    tokenConfig,
				OCRSecrets:     deployment.XXXGenerateTestOCRSecrets(),
				OCRParams:      ocrParams,
				USDCConfig: changeset.USDCConfig{
					EnabledChains:         usdcChains,
					USDCAttestationConfig: usdcAttestationCfg,
					CCTPTokenConfig:       usdcCCTPConfig,
				},
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.CCIPCapabilityJobspec),
		},
	})
	require.NoError(t, err)

	// Ensure capreg logs are up to date.
	changeset.ReplayLogs(t, e.Offchain, replayBlocks)

	tEnv := changeset.DeployedEnv{
		Env:          env,
		HomeChainSel: homeChainSel,
		FeedChainSel: feedSel,
		ReplayBlocks: replayBlocks,
	}
	_ = tEnv
	////_, err = fake.NewFakeDataProvider(in.MockerDataProvider)
	////require.NoError(t, err)
	//out, err := ns.NewSharedDBNodeSet(in.NodeSet, blockchains[0])
	//require.NoError(t, err)
	//
	//t.Run("test something", func(t *testing.T) {
	//	for _, n := range out.CLNodes {
	//		require.NotEmpty(t, n.Node.HostURL)
	//	}
	//})
}
