package cre

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	writetarget "github.com/smartcontractkit/chainlink-solana/pkg/solana/write_target"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/rpc"
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
	mustSetCapabilitiesFn := func(input []*ns.Input) []*cre.CapabilitiesAwareNodeSet {
		return []*cre.CapabilitiesAwareNodeSet{
			{
				Input:              input[0],
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
	fmt.Println("forwarder address", setupOut.ForwarderAddress, "state", setupOut.ForwarderState)
	err = mocksClient.Execute(context.TODO(), &pb.ExecutableRequest{
		ID:              setupOut.DeriveRemaining,
		CapabilityType:  4,
		Config:          []byte{},
		Inputs:          []byte{},
		RequestMetadata: &pb.Metadata{},
	})
	require.NoError(t, err)
}

type setupWTOutput struct {
	WriteCap         string
	DeriveRemaining  string
	SolChainID       string
	ForwarderAddress string
	ForwarderState   string
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
			solwriterconfig.GetGenerateConfig(solwriterconfig.Config{}),
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
		}
	}

	// deploy df cache

	return out
}
