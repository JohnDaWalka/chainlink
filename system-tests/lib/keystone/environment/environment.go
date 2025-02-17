package environment

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"

	"github.com/smartcontractkit/chainlink/system-tests/lib/keystone/blockchain"
	libcontracts "github.com/smartcontractkit/chainlink/system-tests/lib/keystone/contracts"
	libdon "github.com/smartcontractkit/chainlink/system-tests/lib/keystone/don"
	libjobs "github.com/smartcontractkit/chainlink/system-tests/lib/keystone/don/jobs"
	libnode "github.com/smartcontractkit/chainlink/system-tests/lib/keystone/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/keystone/types"
)

// TODO think whether we should structure it in a way that enforces some order of execution,
// for example by making the outputs of one function inputs to another
func StartAndConfigure(
	t *testing.T,
	config types.KeystoneConfiguration,
	registerWorkflowFn types.KeystoneEnvironmentConsumerFn,
	prepareJobSpecsAndNodeConfigsFn types.JobAndConfigProducingFn,
) (*types.KeystoneEnvironment, error) {
	testLogger := framework.L
	keystoneEnv, err := Start(cldlogger.NewSingleFileLogger(t), testLogger, config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start keystone environment")
	}

	// Configure Workflow Registry and Feeds Consumer
	err = libcontracts.ConfigureWorkflowRegistry(testLogger, keystoneEnv)
	if err != nil {
		return nil, errors.Wrap(err, "failed to configure workflow registry")
	}

	// Register the workflow(s) with workflow registry and execute whaterver preparatory steps are needed
	err = registerWorkflowFn(keystoneEnv)
	if err != nil {
		return nil, errors.Wrap(err, "failed to register workflow")
	}

	donToConfigs, donToJobSpecs, err := prepareJobSpecsAndNodeConfigsFn(keystoneEnv)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare job specs and node configs")
	}

	err = libdon.Configure(t, testLogger, keystoneEnv, donToJobSpecs, donToConfigs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to configure nodes")
	}

	// CAUTION: It is crucial to configure OCR3 jobs on nodes before configuring the workflow contracts.
	// Wait for OCR listeners to be ready before setting the configuration.
	// If the ConfigSet event is missed, OCR protocol will not start.
	testLogger.Info().Msg("Waiting 30s for OCR listeners to be ready...")
	time.Sleep(30 * time.Second)
	testLogger.Info().Msg("Proceeding to set OCR3 configuration.")

	// Configure the Forwarder, OCR3 and Capabilities contracts
	err = libcontracts.ConfigureKeystone(keystoneEnv)
	if err != nil {
		return nil, errors.Wrap(err, "failed to configure keystone contracts")
	}

	return keystoneEnv, nil
}

func Start(cldLogger logger.Logger, testLogger zerolog.Logger, config types.KeystoneConfiguration) (*types.KeystoneEnvironment, error) {
	keystoneEnv := &types.KeystoneEnvironment{}
	keystoneEnv.GatewayConnectorData = &types.GatewayConnectorData{
		Path: "/node",
		Port: 5003,
	}

	bcInput, err := config.BlockchainInput()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get blockchain input")
	}

	// Create a new blockchain network and Seth client to interact with it
	err = blockchain.Start(bcInput, keystoneEnv)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start blockchain")
	}

	jdInput, err := config.JdInput()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get JD input")
	}

	// Start job distributor
	err = libjobs.StartJobDistributor(jdInput, keystoneEnv)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start job distributor")
	}

	nodeSetInput, err := config.NodeSetInput()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get node set input")
	}

	// Deploy the DONs
	err = libdon.Start(nodeSetInput, keystoneEnv)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start node sets")
	}

	// Prepare the CLD environment and figure out DON topology; configure chains for nodes and job distributor
	err = BuildTopologyAndCLDEnvironment(cldLogger, keystoneEnv)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build topology and CLD environment")
	}

	// Fund the nodes
	err = libdon.FundNodes(keystoneEnv)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fund nodes")
	}

	// Deploy keystone contracts (forwarder, capability registry, ocr3 capability, workflow registry)
	err = libcontracts.DeployKeystone(testLogger, keystoneEnv)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deploy keystone contracts")
	}

	return keystoneEnv, nil
}

func BuildTopologyAndCLDEnvironment(lgr logger.Logger, keystoneEnv *types.KeystoneEnvironment) error {
	err := buildChainlinkDeploymentEnv(lgr, keystoneEnv)
	if err != nil {
		return errors.Wrap(err, "failed to build chainlink deployment environment")
	}
	err = libdon.BuildDONTopology(keystoneEnv)
	if err != nil {
		return errors.Wrap(err, "failed to build DON topology")
	}

	return nil
}

func buildChainlinkDeploymentEnv(lgr logger.Logger, keystoneEnv *types.KeystoneEnvironment) error {
	if keystoneEnv == nil {
		return errors.New("keystone environment must be set")
	}

	envs := make([]*deployment.Environment, len(keystoneEnv.MustWrappedNodeOutput()))
	keystoneEnv.Dons = make([]*devenv.DON, len(keystoneEnv.MustWrappedNodeOutput()))

	for i, nodeOutput := range keystoneEnv.MustWrappedNodeOutput() {
		// assume that each nodeset has only one bootstrap node
		nodeInfo, err := libnode.GetNodeInfo(nodeOutput.Output, nodeOutput.NodeSetName, 1)
		if err != nil {
			return errors.Wrap(err, "failed to get node info")
		}

		jdConfig := devenv.JDConfig{
			GRPC:     keystoneEnv.MustJD().HostGRPCUrl,
			WSRPC:    keystoneEnv.MustJD().DockerWSRPCUrl,
			Creds:    insecure.NewCredentials(),
			NodeInfo: nodeInfo,
		}

		devenvConfig := devenv.EnvironmentConfig{
			JDConfig: jdConfig,
			Chains: []devenv.ChainConfig{
				{
					ChainID:   keystoneEnv.MustSethClient().Cfg.Network.ChainID,
					ChainName: keystoneEnv.MustSethClient().Cfg.Network.Name,
					ChainType: strings.ToUpper(keystoneEnv.MustBlockchain().Family),
					WSRPCs: []devenv.CribRPCs{{
						External: keystoneEnv.MustBlockchain().Nodes[0].HostWSUrl,
						Internal: keystoneEnv.MustBlockchain().Nodes[0].DockerInternalWSUrl,
					}},
					HTTPRPCs: []devenv.CribRPCs{{
						External: keystoneEnv.MustBlockchain().Nodes[0].HostHTTPUrl,
						Internal: keystoneEnv.MustBlockchain().Nodes[0].DockerInternalHTTPUrl,
					}},
					DeployerKey: keystoneEnv.MustSethClient().NewTXOpts(seth.WithNonce(nil)), // set nonce to nil, so that it will be fetched from the chain
				},
			},
		}

		env, don, err := devenv.NewEnvironment(context.Background, lgr, devenvConfig)
		if err != nil {
			return errors.Wrap(err, "failed to create environment")
		}

		envs[i] = env
		keystoneEnv.Dons[i] = don
	}

	var nodeIDs []string
	for _, env := range envs {
		nodeIDs = append(nodeIDs, env.NodeIDs...)
	}

	// we assume that all DONs run on the same chain and that there's only one chain
	// also, we don't care which instance of offchain client we use, because we have
	// only one instance of offchain client and we have just configured it to work
	// with nodes from all DONs
	keystoneEnv.Environment = &deployment.Environment{
		Name:              envs[0].Name,
		Logger:            envs[0].Logger,
		ExistingAddresses: envs[0].ExistingAddresses,
		Chains:            envs[0].Chains,
		Offchain:          envs[0].Offchain,
		OCRSecrets:        envs[0].OCRSecrets,
		GetContext:        envs[0].GetContext,
		NodeIDs:           nodeIDs,
	}

	return nil
}
