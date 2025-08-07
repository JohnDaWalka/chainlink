package cre

import (
	"context"
	"fmt"
	"math/big"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	writetarget "github.com/smartcontractkit/chainlink-solana/pkg/solana/write_target"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/rpc"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	df_cs "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/solana"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	ks_solana "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/solana"
	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	consensuscap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/consensus"
	mockcap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/mock"
	writesolcap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/writesolana"
	gatewayconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/config/gateway"
	solwriterconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/config/writesolana"
	creconsensus "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/consensus"
	cregateway "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/gateway"
	cremock "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/mock"
	creenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	mock_capability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/mock"
	mockcapability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/mock"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/mock/pb"
	"github.com/stretchr/testify/require"
)

var SinglePoRDonCapabilitiesFlagsSolana = []string{cre.OCR3Capability, cre.WriteSolanaCapability, cre.MockCapability}

func Test_WT_solana_with_mocked_capabilities(t *testing.T) {
	configErr := setConfigurationIfMissing("environment-one-don-multichain-solana-ci.toml")
	require.NoError(t, configErr, "failed to set CTF config")
	testLogger := framework.L

	// Load and validate test configuration
	in, err := framework.Load[TestConfig](t)
	require.NoError(t, err, "couldn't load test config")
	validateEnvVars(t)
	require.Len(t, in.NodeSets, 1, "expected 1 node set in the test config")
	// Assign all capabilities to the single node set
	var solChains []string
	for _, chain := range in.Blockchains {
		if chain.Type == "solana" {
			solChains = append(solChains, chain.ChainID)
		}
	}
	mustSetCapabilitiesFn := func(input []*ns.Input) []*cre.CapabilitiesAwareNodeSet {
		return []*cre.CapabilitiesAwareNodeSet{
			{
				Input:              input[0],
				SupportedSolChains: solChains,
				Capabilities:       SinglePoRDonCapabilitiesFlagsSolana,
				DONTypes:           []string{cre.WorkflowDON, cre.GatewayDON},
				BootstrapNodeIndex: 0, // not required, but set to make the configuration explicit
				GatewayNodeIndex:   0, // not required, but set to make the configuration explicit
			},
		}
	}

	feedIDs := make([]string, 0, len(in.WorkflowConfigs))
	for _, wc := range in.WorkflowConfigs {
		feedIDs = append(feedIDs, wc.FeedID)
	}

	capabilityFactoryFns := []cre.DONCapabilityWithConfigFactoryFn{
		consensuscap.OCR3CapabilityFactoryFn,
		mockcap.MockCapabilityFactoryFn,
	}

	for _, bc := range in.Blockchains {
		capabilityFactoryFns = append(capabilityFactoryFns, writesolcap.WriteSolanaCapabilityFactory(bc.ChainID))
	}

	setupOut := setupWTTestEnvironment(
		t,
		testLogger,
		in,
		mustSetCapabilitiesFn,
		capabilityFactoryFns,
	)
	// Log extra information that might help debugging
	// t.Cleanup(func() {
	//	debugTest(t, testLogger, setupOutput, in)
	// })

	mocksClient := mock_capability.NewMockCapabilityController(testLogger)
	mockClientsAddress := make([]string, 0)
	if in.Infra.Type == "docker" {
		for _, nodeSet := range in.NodeSets {
			if nodeSet.Name == "workflow" {
				for i, n := range nodeSet.NodeSpecs {
					if i == 0 {
						continue
					}
					if len(n.Node.CustomPorts) == 0 {
						panic("no custom port specified, mock capability running in kind must have a custom port in order to connect")
					}
					ports := strings.Split(n.Node.CustomPorts[0], ":")
					mockClientsAddress = append(mockClientsAddress, "127.0.0.1:"+ports[0])
				}
			}
		}
	}

	require.NoError(t, mocksClient.ConnectAll(mockClientsAddress, true, true), "could not connect to mock capabilities")
	fmt.Println("cap name", setupOut.WriteCap)
	fmt.Println("forwarder address", setupOut.ForwarderAddress, "forwarder state", setupOut.ForwarderState)
	fmt.Println("cache address", setupOut.CacheAddress, "cache state", setupOut.CacheState)

	//err = mocksClient.Execute(context.TODO(), &pb.ExecutableRequest{})
	val, err := values.WrapMap(struct {
		State    string // State pubkey of df cache
		Receiver string // df cache programID
		FeedIDs  []string
	}{
		setupOut.CacheState,
		setupOut.CacheAddress,
		[]string{in.WorkflowConfigs[0].FeedID},
	})
	require.NoError(t, err, "failed to wrap value")
	cfg, err := mockcapability.MapToBytes(val)
	require.NoError(t, err, "failed map to bytes config")
	ret, err := mocksClient.Nodes[1].API.Execute(context.TODO(), &pb.ExecutableRequest{
		ID:             setupOut.DeriveRemaining,
		CapabilityType: 4,
		Config:         cfg,
		Inputs:         []byte{},
		RequestMetadata: &pb.Metadata{
			WorkflowOwner: setupOut.WFOwner,
			WorkflowName:  setupOut.WFName,
		},
	})

	require.NoError(t, err, "execute capability failed")
	val, err = mockcapability.BytesToMap(ret.Value)
	require.NoError(t, err, "unmarshal response failed")

	var res solana.AccountMetaSlice

	err = val.Underlying["remaining_accounts"].UnwrapTo(&res)
	require.NoError(t, err, "unwrap response failed")

	require.Len(t, res, 3)
}

type setupWTOutput struct {
	WriteCap        string
	DeriveRemaining string
	SolChainID      string

	ForwarderAddress string
	ForwarderState   string

	CacheAddress string
	CacheState   string

	WFName  string
	WFOwner string
	FeedID  string
}

func setupWTTestEnvironment(
	t *testing.T,
	testLogger zerolog.Logger,
	in *TestConfig,
	mustSetCapabilitiesFn func(input []*ns.Input) []*cre.CapabilitiesAwareNodeSet,
	capabilityFactoryFns []func([]string) []keystone_changeset.DONCapabilityWithConfig,
) *setupWTOutput {
	extraAllowedGatewayPorts := []int{}

	customBinariesPaths := map[string]string{}
	containerPath, pathErr := capabilities.DefaultContainerDirectory(in.Infra.Type)
	require.NoError(t, pathErr, "failed to get default container directory")
	var mockBinaryPathInTheContainer string
	if in.DependenciesConfig.MockCapapilityBinaryPath != "" {
		// where cron binary is located in the container
		mockBinaryPathInTheContainer = filepath.Join(containerPath, filepath.Base(in.DependenciesConfig.MockCapapilityBinaryPath))
		// where cron binary is located on the host
		customBinariesPaths[cre.MockCapability] = in.DependenciesConfig.MockCapapilityBinaryPath
	}

	t.Log("customBinariesPaths", customBinariesPaths)
	firstBlockchain := in.Blockchains[0]

	chainIDInt, err := strconv.Atoi(firstBlockchain.ChainID)
	require.NoError(t, err, "failed to convert chain ID to int")
	chainIDUint64 := libc.MustSafeUint64(int64(chainIDInt))
	cfg := solwriterconfig.GetGenerateConfig()

	universalSetupInput := creenv.SetupInput{
		CapabilitiesAwareNodeSets:            mustSetCapabilitiesFn(in.NodeSets),
		CapabilitiesContractFactoryFunctions: capabilityFactoryFns,
		BlockchainsInput:                     in.Blockchains,
		JdInput:                              *in.JD,
		InfraInput:                           *in.Infra,
		CustomBinariesPaths:                  customBinariesPaths,
		JobSpecFactoryFunctions: []cre.JobSpecFactoryFn{
			creconsensus.ConsensusJobSpecFactoryFn(chainIDUint64),
			cregateway.GatewayJobSpecFactoryFn(extraAllowedGatewayPorts, []string{}, []string{"0.0.0.0/0"}),
			cremock.MockJobSpecFactoryFn(mockBinaryPathInTheContainer),
		},
		ConfigFactoryFunctions: []cre.ConfigFactoryFn{
			gatewayconfig.GenerateConfig,
			cfg,
		},
	}

	universalSetupOutput, setupErr := creenv.SetupTestEnvironment(t.Context(), testLogger, cldlogger.NewSingleFileLogger(t), universalSetupInput)
	require.NoError(t, setupErr, "failed to setup test environment")

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
	out := &setupWTOutput{}
	wfName := [][10]uint8{{1, 2, 3}}
	wfDescription := [][32]uint8{{2, 3, 4}}
	wfOwner := [20]uint8{1}
	out.WFName = string(wfName[0][:])
	out.WFOwner = string(wfOwner[:])
	for _, bo := range universalSetupOutput.BlockchainOutput {
		if bo.ReadOnly {
			continue
		}

		if bo.SolChain != nil {
			chainID, err := bo.SolClient.GetGenesisHash(context.Background())
			require.NoError(t, err, "failed to get genesis hash")
			out.WriteCap = writetarget.GenerateWriteTargetName(chainID.String())
			out.DeriveRemaining = writetarget.GenerateDeriveRemainingName(chainID.String())
			forwarder, err := universalSetupOutput.CldEnvironment.DataStore.Addresses().Get(datastore.NewAddressRefKey(
				bo.SolChain.ChainSelector,
				ks_solana.ForwarderContract,
				semver.MustParse("1.0.0"),
				"test-forwarder",
			))
			require.NoError(t, err, "forwarder not found")
			forwarderState, err := universalSetupOutput.CldEnvironment.DataStore.Addresses().Get(datastore.NewAddressRefKey(
				bo.SolChain.ChainSelector,
				ks_solana.ForwarderState,
				semver.MustParse("1.0.0"),
				"test-forwarder",
			))
			out.ForwarderAddress = forwarder.Address
			out.ForwarderState = forwarderState.Address

			//df cache
			dfQualifier := "df-cache-qualifier"
			dfDeployOut, err := commonchangeset.RunChangeset(df_cs.DeployCache{}, *universalSetupOutput.CldEnvironment, &df_cs.DeployCacheRequest{
				ChainSel:   bo.SolChain.ChainSelector,
				Qualifier:  dfQualifier,
				Version:    "1.0.0",
				FeedAdmins: []solana.PublicKey{bo.SolChain.PrivateKey.PublicKey()},
			})
			require.NoError(t, err, "failed to deploy df cache")
			cacheID, err := dfDeployOut.DataStore.Addresses().Get(
				datastore.NewAddressRefKey(bo.SolChain.ChainSelector, df_cs.CacheContract, semver.MustParse("1.0.0"), dfQualifier))
			require.NoError(t, err, "df cache address not found")
			out.CacheAddress = cacheID.Address

			cacheState, err := dfDeployOut.DataStore.Addresses().Get(
				datastore.NewAddressRefKey(bo.SolChain.ChainSelector, df_cs.CacheState, semver.MustParse("1.0.0"), dfQualifier))
			require.NoError(t, err, "df cache state not found")

			out.CacheAddress = cacheID.Address
			out.CacheState = cacheState.Address
			ds := datastore.NewMemoryDataStore()
			ds.Merge(dfDeployOut.DataStore.Seal())
			ds.Merge(universalSetupOutput.CldEnvironment.DataStore)
			universalSetupOutput.CldEnvironment.DataStore = ds.Seal()
			feedIDin, ok := new(big.Int).SetString(in.WorkflowConfigs[0].FeedID, 0)
			require.True(t, ok, "invalid feedID")
			require.LessOrEqual(t, feedIDin.BitLen(), 128, "invalid feedID len")
			var feedID [16]uint8
			copy(feedID[:], feedIDin.Bytes())
			require.NoError(t, err, "failed to decode FeedID")
			_, err = commonchangeset.RunChangeset(df_cs.InitCacheDecimalReport{}, *universalSetupOutput.CldEnvironment,
				&df_cs.InitCacheDecimalReportRequest{
					ChainSel:  bo.SolChain.ChainSelector,
					Qualifier: dfQualifier,
					Version:   "1.0.0",
					FeedAdmin: bo.SolChain.PrivateKey.PublicKey(),
					DataIDs:   [][16]uint8{feedID},
				},
			)
			require.NoError(t, err, "failed to init decimal report")

			_, err = commonchangeset.RunChangeset(df_cs.ConfigureCacheDecimalReport{}, *universalSetupOutput.CldEnvironment,
				&df_cs.ConfigureCacheDecimalReportRequest{
					ChainSel:  bo.SolChain.ChainSelector,
					Qualifier: dfQualifier,
					Version:   "1.0.0",
					SenderList: []df_cs.Sender{
						{
							ProgramID: solana.MustPublicKeyFromBase58(forwarder.Address),
							StateID:   solana.MustPublicKeyFromBase58(forwarderState.Address),
						},
					},
					FeedAdmin:            bo.SolChain.PrivateKey.PublicKey(),
					DataIDs:              [][16]uint8{feedID},
					AllowedWorkflowOwner: [][20]uint8{wfOwner},
					AllowedWorkflowName:  wfName,
					Descriptions:         wfDescription,
				})
			require.NoError(t, err, "failed to configure decimal report")
		}
	}

	return out
}
