package environment

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	ctfconfig "github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"
	libcontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/crib"
	libdevenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/devenv"
	libdon "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	cretypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	libformat "github.com/smartcontractkit/chainlink/system-tests/lib/format"
	libfunding "github.com/smartcontractkit/chainlink/system-tests/lib/funding"
	libnix "github.com/smartcontractkit/chainlink/system-tests/lib/nix"
	libtypes "github.com/smartcontractkit/chainlink/system-tests/lib/types"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func SketchSetupTestEnvironment(
	ctx context.Context,
	testLogger zerolog.Logger,
	singeFileLogger *cldlogger.SingleFileLogger,
	input SetupInput,
) (*SetupOutput, error) {
	topologyErr := libdon.ValidateTopology(input.CapabilitiesAwareNodeSets, input.InfraInput)
	if topologyErr != nil {
		return nil, pkgerrors.Wrap(topologyErr, "failed to validate topology")
	}

	// Shell is only required, when using CRIB, because we want to run commands in the same "nix develop" context
	// We need to have this reference in the outer scope, because subsequent functions will need it
	var nixShell *libnix.Shell
	if input.InfraInput.InfraType == libtypes.CRIB {
		startNixShellInput := &cretypes.StartNixShellInput{
			InfraInput:     &input.InfraInput,
			CribConfigsDir: cribConfigsDir,
			PurgeNamespace: true,
		}

		var nixErr error
		nixShell, nixErr = crib.StartNixShell(startNixShellInput)
		if nixErr != nil {
			return nil, pkgerrors.Wrap(nixErr, "failed to start nix shell")
		}
	}

	defer func() {
		if nixShell != nil {
			_ = nixShell.Close()
		}
	}()

	allChainsCLDEnvironment := &cldf.Environment{
		Logger:            singeFileLogger,
		ExistingAddresses: cldf.NewMemoryAddressBook(),
		GetContext: func() context.Context {
			return ctx
		},
		// TODO: init operations bundle
	}

	bi := BlockchainsInput{
		infra:    &input.InfraInput,
		nixShell: nixShell,
	}
	bi.blockchainsInput = append(bi.blockchainsInput, input.BlockchainsInput...)

	startTime := time.Now()
	fmt.Print(libformat.PurpleText("\n[Stage 1/10] Starting %d blockchain(s)\n\n", len(bi.blockchainsInput)))

	// TODO: should this really be an operation?
	blkR, bcOutErr := operations.ExecuteOperation(allChainsCLDEnvironment.OperationsBundle, StartBlockchainsOp, StartBlockchainsDeps{
		logger:          zerolog.Logger{},
		singeFileLogger: singeFileLogger,
	}, bi)
	if bcOutErr != nil {
		return nil, pkgerrors.Wrap(bcOutErr, "failed to create blockchains")
	}
	blockchainsOutput := blkR.Output.Outputs
	homeChainOutput := blockchainsOutput[0]
	allChainsCLDEnvironment.BlockChains = chain.NewBlockChains(blkR.Output.Blockchains)

	fmt.Print(libformat.PurpleText("\n[Stage 1/10] Blockchains started in %.2f seconds\n", time.Since(startTime).Seconds()))
	startTime = time.Now()
	fmt.Print(libformat.PurpleText("[Stage 2/10] Deploying Keystone contracts\n\n"))

	seqR, err := operations.ExecuteSequence(
		allChainsCLDEnvironment.OperationsBundle,
		keystone_changeset.DeployKeystoneContractsSequence,
		keystone_changeset.DeployKeystoneContractsSequenceDeps{},
		keystone_changeset.DeployKeystoneContractsSequenceInput{HomeChainSelector: homeChainOutput.ChainSelector},
	)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to execute Keystone contracts deployment sequence")
	}
	// Merge the address book from the sequence output into the CLD environment
	if mergeErr := allChainsCLDEnvironment.ExistingAddresses.Merge(seqR.Output.AddressBook); mergeErr != nil { //nolint:staticcheck // won't migrate now
		return nil, pkgerrors.Wrap(mergeErr, "failed to merge address book from Keystone contracts deployment sequence")
	}

	var fwrChains []uint64

	// Deploy forwarders for all chains
	for _, bcOut := range blockchainsOutput {
		if bcOut.ChainSelector == homeChainOutput.ChainSelector {
			// Skip the home chain, because we already deployed the forwarder there
			continue
		}
		fwrChains = append(fwrChains, bcOut.ChainSelector)
	}
	frwR, err := operations.ExecuteSequence(
		allChainsCLDEnvironment.OperationsBundle,
		keystone_changeset.DeployKeystoneForwardersSequence,
		keystone_changeset.DeployKeystoneContractsSequenceDeps{},
		keystone_changeset.DeployKeystoneForwardersInput{Targets: fwrChains},
	)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to execute Keystone forwarders deployment sequence")
	}
	// Merge the address book from the sequence output into the CLD environment
	if mergeErr := allChainsCLDEnvironment.ExistingAddresses.Merge(frwR.Output.AddressBook); mergeErr != nil { //nolint:staticcheck // won't migrate now
		return nil, pkgerrors.Wrap(mergeErr, "failed to merge address book from Keystone forwarders deployment sequence")
	}

	fmt.Print(libformat.PurpleText("\n[Stage 2/10] Contracts deployed in %.2f seconds\n", time.Since(startTime).Seconds()))

	// Translate node input to structure required further down the road and put as much information
	// as we have at this point in labels. It will be used to generate node configs
	topologyReport, topologyErr := operations.ExecuteOperation(allChainsCLDEnvironment.OperationsBundle, BuildTopologyOp, BuildTopologyOpDeps{}, BuildTopologyOpInput{})
	if topologyErr != nil {
		return nil, pkgerrors.Wrap(topologyErr, "failed to build topology")
	}
	topology := topologyReport.Output.Topology

	startTime = time.Now()
	fmt.Print(libformat.PurpleText("[Stage 4/10] Preparing DON(s) configuration\n\n"))

	// Deploy the DONs
	// Hack for CI that allows us to dynamically set the chainlink image and version
	// CTFv2 currently doesn't support dynamic image and version setting
	if os.Getenv("CI") == "true" {
		// Due to how we pass custom env vars to reusable workflow we need to use placeholders, so first we need to resolve what's the name of the target environment variable
		// that stores chainlink version and then we can use it to resolve the image name
		for i := range input.CapabilitiesAwareNodeSets {
			image := fmt.Sprintf("%s:%s", os.Getenv(ctfconfig.E2E_TEST_CHAINLINK_IMAGE_ENV), ctfconfig.MustReadEnvVar_String(ctfconfig.E2E_TEST_CHAINLINK_VERSION_ENV))
			for j := range input.CapabilitiesAwareNodeSets[i].NodeSpecs {
				input.CapabilitiesAwareNodeSets[i].NodeSpecs[j].Node.Image = image
			}
		}
	}

	bcOuts := make(map[uint64]*blockchain.Output)
	sethClients := make(map[uint64]*seth.Client)
	for _, bcOut := range blockchainsOutput {
		bcOuts[bcOut.ChainSelector] = bcOut.BlockchainOutput
		sethClients[bcOut.ChainSelector] = bcOut.SethClient
	}

	fmt.Print(libformat.PurpleText("\n[Stage 4/10] DONs configuration prepared in %.2f seconds\n", time.Since(startTime).Seconds()))
	startTime = time.Now()
	fmt.Print(libformat.PurpleText("[Stage 5/10] Starting Job Distributor\n\n"))

	if input.InfraInput.InfraType == libtypes.CRIB {
		deployCribJdInput := &cretypes.DeployCribJdInput{
			JDInput:        &input.JdInput,
			NixShell:       nixShell,
			CribConfigsDir: cribConfigsDir,
		}

		var jdErr error
		input.JdInput.Out, jdErr = crib.DeployJd(deployCribJdInput)
		if jdErr != nil {
			return nil, pkgerrors.Wrap(jdErr, "failed to deploy JD with devspace")
		}
	}

	jdOutput, jdErr := CreateJobDistributor(&input.JdInput)
	if jdErr != nil {
		jdErr = fmt.Errorf("failed to start JD container for image %s: %w", input.JdInput.Image, jdErr)

		// useful end user messages
		if strings.Contains(jdErr.Error(), "pull access denied") || strings.Contains(jdErr.Error(), "may require 'docker login'") {
			jdErr = errors.Join(jdErr, errors.New("ensure that you either you have built the local image or you are logged into AWS with a profile that can read it (`aws sso login --profile <foo>)`"))
		}
		return nil, jdErr
	}

	fmt.Print(libformat.PurpleText("\n[Stage 5/10] Job Distributor started in %.2f seconds\n", time.Since(startTime).Seconds()))
	startTime = time.Now()
	fmt.Print(libformat.PurpleText("[Stage 6/10] Starting %d DON(s)\n\n", len(input.CapabilitiesAwareNodeSets)))

	if input.InfraInput.InfraType == libtypes.CRIB {
		testLogger.Info().Msg("Saving node configs and secret overrides")
		deployCribDonsInput := &cretypes.DeployCribDonsInput{
			Topology:       topology,
			NodeSetInputs:  input.CapabilitiesAwareNodeSets,
			NixShell:       nixShell,
			CribConfigsDir: cribConfigsDir,
		}

		var devspaceErr error
		input.CapabilitiesAwareNodeSets, devspaceErr = crib.DeployDons(deployCribDonsInput)
		if devspaceErr != nil {
			return nil, pkgerrors.Wrap(devspaceErr, "failed to deploy Dons with devspace")
		}
	}

	nodeSetOutput := make([]*cretypes.WrappedNodeOutput, 0, len(input.CapabilitiesAwareNodeSets))
	for _, nodeSetInput := range input.CapabilitiesAwareNodeSets {
		nodeset, nodesetErr := ns.NewSharedDBNodeSet(nodeSetInput.Input, homeChainOutput.BlockchainOutput)
		if nodesetErr != nil {
			return nil, pkgerrors.Wrapf(nodesetErr, "failed to create node set named %s", nodeSetInput.Name)
		}

		nodeSetOutput = append(nodeSetOutput, &cretypes.WrappedNodeOutput{
			Output:       nodeset,
			NodeSetName:  nodeSetInput.Name,
			Capabilities: nodeSetInput.Capabilities,
		})
	}

	// Prepare the CLD environment that's required by the keystone changeset
	// Ugly glue hack ¯\_(ツ)_/¯
	fullCldInput := &cretypes.FullCLDEnvironmentInput{
		JdOutput:          jdOutput,
		BlockchainOutputs: bcOuts,
		SethClients:       sethClients,
		NodeSetOutput:     nodeSetOutput,
		ExistingAddresses: allChainsCLDEnvironment.ExistingAddresses, //nolint:staticcheck // won't migrate now
		Topology:          topology,
	}

	// We need to use TLS for CRIB, because it exposes HTTPS endpoints
	var creds credentials.TransportCredentials
	if input.InfraInput.InfraType == libtypes.CRIB {
		creds = credentials.NewTLS(&tls.Config{
			MinVersion: tls.VersionTLS12,
		})
	} else {
		creds = insecure.NewCredentials()
	}

	fullCldOutput, cldErr := libdevenv.BuildFullCLDEnvironment(ctx, singeFileLogger, fullCldInput, creds)
	if cldErr != nil {
		return nil, pkgerrors.Wrap(cldErr, "failed to build full CLD environment")
	}

	fmt.Print(libformat.PurpleText("\n[Stage 6/10] DONs started in %.2f seconds\n", time.Since(startTime).Seconds()))
	startTime = time.Now()
	fmt.Print(libformat.PurpleText("[Stage 7/10] Funding Chainlink nodes\n\n"))

	// Fund the nodes
	concurrentNonceMap, concurrentNonceMapErr := NewConcurrentNonceMap(ctx, blockchainsOutput)
	if concurrentNonceMapErr != nil {
		return nil, pkgerrors.Wrap(concurrentNonceMapErr, "failed to create concurrent nonce map")
	}

	// Decrement the nonce for each chain, because we will increment it in the next loop
	for _, bcOut := range blockchainsOutput {
		concurrentNonceMap.Decrement(bcOut.ChainID)
	}

	errGroup := &errgroup.Group{}
	for _, metaDon := range fullCldOutput.DonTopology.DonsWithMetadata {
		for _, bcOut := range blockchainsOutput {
			for _, node := range metaDon.DON.Nodes {
				errGroup.Go(func() error {
					nodeAddress := node.AccountAddr[strconv.FormatUint(bcOut.ChainID, 10)]
					if nodeAddress == "" {
						return nil
					}

					nonce := concurrentNonceMap.Increment(bcOut.ChainID)

					_, fundingErr := libfunding.SendFunds(ctx, zerolog.Logger{}, bcOut.SethClient, libtypes.FundsToSend{
						ToAddress:  common.HexToAddress(nodeAddress),
						Amount:     big.NewInt(5000000000000000000),
						PrivateKey: bcOut.SethClient.MustGetRootPrivateKey(),
						Nonce:      ptr.Ptr(nonce),
					})
					if fundingErr != nil {
						return pkgerrors.Wrapf(fundingErr, "failed to fund node %s", nodeAddress)
					}
					return nil
				})
			}
		}
	}

	if err := errGroup.Wait(); err != nil {
		return nil, pkgerrors.Wrap(err, "failed to fund nodes")
	}

	fmt.Print(libformat.PurpleText("\n[Stage 7/10] Chainlink nodes funded in %.2f seconds\033[0m\n", time.Since(startTime).Seconds()))
	startTime = time.Now()
	fmt.Print(libformat.PurpleText("[Stage 8/10] Creating jobs with Job Distributor\n\n"))

	donToJobSpecs := make(cretypes.DonsToJobSpecs)

	for _, jobSpecGeneratingFn := range input.JobSpecFactoryFunctions {
		singleDonToJobSpecs, jobSpecsErr := jobSpecGeneratingFn(&cretypes.JobSpecFactoryInput{
			CldEnvironment:   fullCldOutput.Environment,
			BlockchainOutput: homeChainOutput.BlockchainOutput,
			DonTopology:      fullCldOutput.DonTopology,
			AddressBook:      allChainsCLDEnvironment.ExistingAddresses, //nolint:staticcheck // won't migrate now
		})
		if jobSpecsErr != nil {
			return nil, pkgerrors.Wrap(jobSpecsErr, "failed to generate job specs")
		}
		mergeJobSpecSlices(singleDonToJobSpecs, donToJobSpecs)
	}

	createJobsInput := cretypes.CreateJobsInput{
		CldEnv:        fullCldOutput.Environment,
		DonTopology:   fullCldOutput.DonTopology,
		DonToJobSpecs: donToJobSpecs,
	}

	jobsErr := libdon.CreateJobs(ctx, testLogger, createJobsInput)
	if jobsErr != nil {
		return nil, pkgerrors.Wrap(jobsErr, "failed to create jobs")
	}

	// CAUTION: It is crucial to configure OCR3 jobs on nodes before configuring the workflow contracts.
	// Wait for OCR listeners to be ready before setting the configuration.
	// If the ConfigSet event is missed, OCR protocol will not start.
	fmt.Print(libformat.PurpleText("\n[Stage 8/10] Jobs created in %.2f seconds\033[0m\n", time.Since(startTime).Seconds()))
	startTime = time.Now()
	fmt.Print(libformat.PurpleText("[Stage 9/10] Waiting for Log Poller to start tracking OCR3 contract\n\n"))

	for idx, nodeSetOut := range nodeSetOutput {
		if !flags.HasFlag(input.CapabilitiesAwareNodeSets[idx].Capabilities, cretypes.OCR3Capability) {
			continue
		}
		nsClients, cErr := clclient.New(nodeSetOut.CLNodes)
		if cErr != nil {
			return nil, pkgerrors.Wrap(cErr, "failed to create node set clients")
		}
		eg := &errgroup.Group{}
		for _, c := range nsClients {
			eg.Go(func() error {
				return c.WaitHealthy(".*ConfigWatcher", "passing", 100)
			})
		}
		if err := eg.Wait(); err != nil {
			return nil, pkgerrors.Wrap(err, "failed to wait for ConfigWatcher health check")
		}
	}

	fmt.Print(libformat.PurpleText("\n[Stage 9/10] Log Poller started in %.2f seconds\n", time.Since(startTime).Seconds()))
	startTime = time.Now()
	fmt.Print(libformat.PurpleText("[Stage 10/10] Configuring OCR3 and Keystone contracts\n\n"))

	// Configure Workflow Registry contract
	workflowRegistryInput := &cretypes.WorkflowRegistryInput{
		ChainSelector:  homeChainOutput.ChainSelector,
		CldEnv:         allChainsCLDEnvironment,
		AllowedDonIDs:  []uint32{topology.WorkflowDONID},
		WorkflowOwners: []common.Address{homeChainOutput.SethClient.MustGetRootKeyAddress()},
		Out: &cretypes.WorkflowRegistryOutput{
			ChainSelector:  homeChainOutput.ChainSelector,
			AllowedDonIDs:  []uint32{topology.WorkflowDONID},
			WorkflowOwners: []common.Address{homeChainOutput.SethClient.MustGetRootKeyAddress()},
		},
	}

	// TODO: properly setup the config reqs to config keystone contracts
	// TODO: make sure that the config seq fully replaces `libcontracts.ConfigureKeystone` here

	// Configure the Forwarder, OCR3 and Capabilities contracts
	configureKeystoneInput := cretypes.ConfigureKeystoneInput{
		ChainSelector: homeChainOutput.ChainSelector,
		CldEnv:        fullCldOutput.Environment,
		Topology:      topology,
	}

	if input.OCR3Config != nil {
		configureKeystoneInput.OCR3Config = *input.OCR3Config
	} else {
		ocr3Config, ocr3ConfigErr := libcontracts.DefaultOCR3Config(topology)
		if ocr3ConfigErr != nil {
			return nil, pkgerrors.Wrap(ocr3ConfigErr, "failed to generate default OCR3 config")
		}
		configureKeystoneInput.OCR3Config = *ocr3Config
	}

	_, keystoneErr := operations.ExecuteSequence(
		fullCldOutput.Environment.OperationsBundle,
		keystone_changeset.ConfigureKeystoneContractsSeq,
		keystone_changeset.ConfigureKeystoneContractsSequenceDeps{},
		keystone_changeset.ConfigureKeystoneContractsSequenceInput{},
	)
	if keystoneErr != nil {
		return nil, pkgerrors.Wrap(keystoneErr, "failed to configure keystone contracts")
	}

	fmt.Print(libformat.PurpleText("\n[Stage 10/10] OCR3 and Keystone contracts configured in %.2f seconds\n", time.Since(startTime).Seconds()))

	if input.InfraInput.InfraType != libtypes.CRIB {
		hasGateway := false
		for _, don := range fullCldOutput.DonTopology.DonsWithMetadata {
			if flags.HasFlag(don.Flags, cretypes.GatewayDON) {
				hasGateway = true
				break
			}
		}

		if hasGateway {
			startTime = time.Now()
			fmt.Print(libformat.PurpleText("[POST-SETUP] Waiting for all nodes to have expected Log Poller filters registered\n\n"))

			testLogger.Info().Msg("Waiting for all nodes to have expected log poller filters registered...")
			lpErr := waitForAllNodesToHaveExpectedFiltersRegistered(singeFileLogger, testLogger, homeChainOutput.ChainID, *fullCldOutput.DonTopology, input.CapabilitiesAwareNodeSets)
			if lpErr != nil {
				return nil, pkgerrors.Wrap(lpErr, "failed to wait for all nodes to have expected filters registered")
			}
			fmt.Print(libformat.PurpleText("\n[POST-SETUP] Wait finished in %.2f seconds\n\n", time.Since(startTime).Seconds()))
		}
	}

	return &SetupOutput{
		WorkflowRegistryConfigurationOutput: workflowRegistryInput.Out, // pass to caller, so that it can be optionally attached to TestConfig and saved to disk
		BlockchainOutput:                    blockchainsOutput,
		DonTopology:                         fullCldOutput.DonTopology,
		NodeOutput:                          nodeSetOutput,
		CldEnvironment:                      fullCldOutput.Environment,
	}, nil
}
