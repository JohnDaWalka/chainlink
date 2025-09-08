package environment

import (
	"context"
	"fmt"
	"os"

	"github.com/alitto/pond/v2"
	"github.com/ethereum/go-ethereum/common"
	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/credentials/insecure"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ctfconfig "github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	ks_contracts_op "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/crib"
	libdevenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/devenv"
	libdon "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	libformat "github.com/smartcontractkit/chainlink/system-tests/lib/format"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

const (
	GithubReadTokenEnvVarName          = "GITHUB_READ_TOKEN"
	E2eJobDistributorImageEnvVarName   = "E2E_JD_IMAGE"
	E2eJobDistributorVersionEnvVarName = "E2E_JD_VERSION"
	cribConfigsDir                     = "crib-configs"
)

type SetupOutput struct {
	WorkflowRegistryConfigurationOutput *cre.WorkflowRegistryOutput
	CldEnvironment                      *cldf.Environment
	BlockchainOutput                    []*cre.WrappedBlockchainOutput
	DonTopology                         *cre.DonTopology
	NodeOutput                          []*cre.WrappedNodeOutput
	InfraInput                          infra.Input
}

type SetupInput struct {
	CapabilitiesAwareNodeSets []*cre.CapabilitiesAwareNodeSet
	BlockchainsInput          []blockchain.Input
	JdInput                   jd.Input
	InfraInput                infra.Input
	ContractVersions          map[string]string
	WithV2Registries          bool
	OCR3Config                *keystone_changeset.OracleConfig
	DONTimeConfig             *keystone_changeset.OracleConfig
	VaultOCR3Config           *keystone_changeset.OracleConfig
	CapabilityConfigs         cre.CapabilityConfigs
	CopyCapabilityBinaries    bool // if true, copy capability binaries to the containers (if false, we assume that the plugins image already has them)
	Capabilities              []cre.InstallableCapability

	// Deprecated: use Capabilities []cre.InstallableCapability instead
	ConfigFactoryFunctions []cre.NodeConfigTransformerFn
	// Deprecated: use Capabilities []cre.InstallableCapability instead
	JobSpecFactoryFunctions []cre.JobSpecFn
	// Deprecated: use Capabilities []cre.InstallableCapability instead
	CapabilitiesContractFactoryFunctions []cre.CapabilityRegistryConfigFn
}

// func mustGetAddress(dataStore datastore.MutableDataStore, chainSel uint64, contractType string, version string, qualifier string) string {
// 	key := datastore.NewAddressRefKey(
// 		chainSel,
// 		datastore.ContractType(contractType),
// 		semver.MustParse(version),
// 		qualifier,
// 	)
// 	addrRef, err := dataStore.Addresses().Get(key)
// 	if err != nil {
// 		panic(fmt.Sprintf("Failed to get %s %s (qualifier=%s) address for chain %d: %s", contractType, version, qualifier, chainSel, err.Error()))
// 	}
// 	return addrRef.Address
// }

var stageCount = 7

func SetupTestEnvironment(
	ctx context.Context,
	testLogger zerolog.Logger,
	singleFileLogger logger.Logger,
	input SetupInput,
) (*SetupOutput, error) {
	topologyErr := libdon.ValidateTopology(input.CapabilitiesAwareNodeSets, input.InfraInput)
	if topologyErr != nil {
		return nil, pkgerrors.Wrap(topologyErr, "failed to validate topology")
	}

	if input.InfraInput.Type == infra.CRIB {
		// In CRIB v2 we no longer rely on devspace to create a namespace so we need to do it before deploying
		err := crib.Bootstrap(&input.InfraInput)
		if err != nil {
			return nil, pkgerrors.Wrap(err, "failed to create namespace")
		}
	}

	stageGen := NewStageGen(stageCount, "STAGE")
	fmt.Print(libformat.PurpleText("%s", stageGen.Wrap("Starting %d blockchain(s)", len(input.BlockchainsInput))))

	startBlockchainsOutput, bcOutErr := StartBlockchains(BlockchainLoggers{
		lggr:       testLogger,
		singleFile: singleFileLogger,
	}, BlockchainsInput{
		infra:            &input.InfraInput,
		blockchainsInput: input.BlockchainsInput,
	})
	if bcOutErr != nil {
		return nil, pkgerrors.Wrap(bcOutErr, "failed to start blockchains")
	}

	homeChainOutput := startBlockchainsOutput.BlockChainOutputs[0]
	homeChainSelector := homeChainOutput.ChainSelector

	// blockchainOutputs := startBlockchainsOutput.BlockChainOutputs
	//homeChainOutput := blockchainOutputs[0]
	// blockChains := startBlockchainsOutput.BlockChains

	// memoryDatastore := datastore.NewMemoryDataStore()
	// allChainsCLDEnvironment := &cldf.Environment{
	// 	Logger:            singleFileLogger,
	// 	ExistingAddresses: cldf.NewMemoryAddressBook(),
	// 	DataStore:         memoryDatastore.Seal(),
	// 	GetContext: func() context.Context {
	// 		return ctx
	// 	},
	// 	BlockChains: cldf_chain.NewBlockChains(blockChains),
	// }
	// allChainsCLDEnvironment.OperationsBundle = operations.NewBundle(allChainsCLDEnvironment.GetContext, singleFileLogger, operations.NewMemoryReporter())

	fmt.Print(libformat.PurpleText("%s", stageGen.WrapAndNext("Blockchains started in %.2f seconds", stageGen.Elapsed().Seconds())))

	// DEPLOY CONTRACTS
	fmt.Print(libformat.PurpleText("%s", stageGen.Wrap("Deploying Keystone contracts")))

	// evmForwardersSelectors := make([]uint64, 0)
	// solForwardersSelectors := make([]uint64, 0)
	// for _, bcOut := range blockchainOutputs {
	// 	for _, donMetadata := range input.CapabilitiesAwareNodeSets {
	// 		if slices.Contains(evmForwardersSelectors, bcOut.ChainSelector) {
	// 			continue
	// 		}
	// 		if flags.RequiresForwarderContract(donMetadata.ComputedCapabilities, bcOut.ChainID) {
	// 			evmForwardersSelectors = append(evmForwardersSelectors, bcOut.ChainSelector)
	// 		}
	// 	}
	// }

	// var allNodeFlags []string
	// for i := range input.CapabilitiesAwareNodeSets {
	// 	nodeFlags, err := flags.NodeSetFlags(input.CapabilitiesAwareNodeSets[i])
	// 	if err != nil {
	// 		continue
	// 	}
	// 	allNodeFlags = append(allNodeFlags, nodeFlags...)
	// }
	// vaultOCR3AddrFlag := flags.HasFlag(allNodeFlags, cre.VaultCapability)
	// evmOCR3AddrFlag := flags.HasFlagForAnyChain(allNodeFlags, cre.EVMCapability)
	// consensusV2AddrFlag := flags.HasFlag(allNodeFlags, cre.ConsensusCapabilityV2)

	// chainsWithEVMCapability := make(map[ks_contracts_op.EVMChainID]ks_contracts_op.Selector)
	// for _, chain := range blockchainOutputs {
	// 	for _, donMetadata := range input.CapabilitiesAwareNodeSets {
	// 		if flags.HasFlagForChain(donMetadata.ComputedCapabilities, cre.EVMCapability, chain.ChainID) {
	// 			if chainsWithEVMCapability[ks_contracts_op.EVMChainID(chain.ChainID)] != 0 {
	// 				continue
	// 			}
	// 			chainsWithEVMCapability[ks_contracts_op.EVMChainID(chain.ChainID)] = ks_contracts_op.Selector(chain.ChainSelector)
	// 		}
	// 	}
	// }

	// // use CLD to deploy the registry contracts, which are required before constructing the node TOML configs
	// homeChainSelector := homeChainOutput.ChainSelector
	// deployRegistrySeq := ks_contracts_op.DeployRegistryContractsSequence
	// if input.WithV2Registries {
	// 	deployRegistrySeq = ks_contracts_op.DeployV2RegistryContractsSequence
	// }

	// registryContractsReport, seqErr := operations.ExecuteSequence(
	// 	allChainsCLDEnvironment.OperationsBundle,
	// 	deployRegistrySeq,
	// 	ks_contracts_op.DeployContractsSequenceDeps{
	// 		Env: allChainsCLDEnvironment,
	// 	},
	// 	ks_contracts_op.DeployRegistryContractsSequenceInput{
	// 		RegistryChainSelector: homeChainSelector,
	// 	},
	// )
	// if seqErr != nil {
	// 	return nil, pkgerrors.Wrap(seqErr, "failed to deploy Keystone contracts")
	// }

	// if err := allChainsCLDEnvironment.ExistingAddresses.Merge(registryContractsReport.Output.AddressBook); err != nil { //nolint:staticcheck // won't migrate now
	// 	return nil, pkgerrors.Wrap(err, "failed to merge address book with Keystone contracts addresses")
	// }

	// if err := memoryDatastore.Merge(registryContractsReport.Output.Datastore); err != nil {
	// 	return nil, pkgerrors.Wrap(err, "failed to merge datastore with Keystone contracts addresses")
	// }

	// // deploy evm forwarders
	// evmForwardersReport, seqErr := operations.ExecuteSequence(
	// 	allChainsCLDEnvironment.OperationsBundle,
	// 	ks_contracts_op.DeployKeystoneForwardersSequence,
	// 	ks_contracts_op.DeployKeystoneForwardersSequenceDeps{
	// 		Env: allChainsCLDEnvironment,
	// 	},
	// 	ks_contracts_op.DeployKeystoneForwardersInput{
	// 		Targets: evmForwardersSelectors,
	// 	},
	// )
	// if seqErr != nil {
	// 	return nil, pkgerrors.Wrap(seqErr, "failed to deploy evm forwarder")
	// }

	// if seqErr = allChainsCLDEnvironment.ExistingAddresses.Merge(evmForwardersReport.Output.AddressBook); seqErr != nil { //nolint:staticcheck // won't migrate now
	// 	return nil, pkgerrors.Wrap(seqErr, "failed to merge address book with Keystone contracts addresses")
	// }

	// if seqErr = memoryDatastore.Merge(evmForwardersReport.Output.Datastore); seqErr != nil {
	// 	return nil, pkgerrors.Wrap(seqErr, "failed to merge datastore with Keystone contracts addresses")
	// }

	// // deploy solana forwarders
	// for _, sel := range solForwardersSelectors {
	// 	out, err := operations.ExecuteSequence(
	// 		allChainsCLDEnvironment.OperationsBundle,
	// 		ks_sol_seq.DeployForwarderSeq,
	// 		ks_sol_op.Deps{
	// 			Env:       *allChainsCLDEnvironment,
	// 			Chain:     allChainsCLDEnvironment.BlockChains.SolanaChains()[sel],
	// 			Datastore: memoryDatastore.Seal(),
	// 		},
	// 		ks_sol_seq.DeployForwarderSeqInput{
	// 			ChainSel:    sel,
	// 			ProgramName: deployment.KeystoneForwarderProgramName,
	// 		},
	// 	)
	// 	if err != nil {
	// 		return nil, pkgerrors.Wrap(err, "failed to deploy sol forwarder")
	// 	}

	// 	err = memoryDatastore.AddressRefStore.Add(datastore.AddressRef{
	// 		Address:       out.Output.ProgramID.String(),
	// 		ChainSelector: sel,
	// 		Version:       semver.MustParse(input.ContractVersions[ks_sol.ForwarderContract.String()]),
	// 		Qualifier:     ks_sol.DefaultForwarderQualifier,
	// 		Type:          ks_sol.ForwarderContract,
	// 	})
	// 	if err != nil {
	// 		return nil, pkgerrors.Wrap(err, "failed to add address to the datastore for Solana Forwarder contract")
	// 	}

	// 	err = memoryDatastore.AddressRefStore.Add(datastore.AddressRef{
	// 		Address:       out.Output.State.String(),
	// 		ChainSelector: sel,
	// 		Version:       semver.MustParse(input.ContractVersions[ks_sol.ForwarderState.String()]),
	// 		Qualifier:     ks_sol.DefaultForwarderQualifier,
	// 		Type:          ks_sol.ForwarderState,
	// 	})
	// 	if err != nil {
	// 		return nil, pkgerrors.Wrap(err, "failed to add address to the datastore for Solana Forwarder state")
	// 	}

	// 	testLogger.Info().Msgf("Deployed Forwarder %s contract on Solana chain chain %d programID: %s state: %s", input.ContractVersions[ks_sol.ForwarderContract.String()], sel, out.Output.ProgramID.String(), out.Output.State.String())
	// }

	// deploy the various ocr contracts
	// TODO move this deeper into the stack when we have all the p2p ids and can deploy and configure in one sequence
	// deploy OCR3 contract
	// we deploy OCR3 contract with a qualifier, so that we can distinguish it from other OCR3 contracts (Vault, EVM, ConsensusV2)
	// TODO track the qualifiers in vars/consts rather than raw strings
	// _, seqErr = deployOCR3Contract("capability_ocr3", homeChainSelector, allChainsCLDEnvironment, memoryDatastore)
	// if seqErr != nil {
	// 	return nil, fmt.Errorf("failed to deploy OCR3 contract %w", seqErr)
	// }
	// // deploy DONTime contract
	// _, seqErr = deployOCR3Contract("DONTime", homeChainSelector, allChainsCLDEnvironment, memoryDatastore)
	// if seqErr != nil {
	// 	return nil, fmt.Errorf("failed to deploy DONTime contract %w", seqErr)
	// }
	// if vaultOCR3AddrFlag {
	// 	_, seqErr = deployOCR3Contract("capability_vault", homeChainSelector, allChainsCLDEnvironment, memoryDatastore)
	// 	if seqErr != nil {
	// 		return nil, fmt.Errorf("failed to deploy Vault OCR3 contract %w", seqErr)
	// 	}
	// }
	// if evmOCR3AddrFlag {
	// 	for chainID, selector := range chainsWithEVMCapability {
	// 		qualifier := ks_contracts_op.CapabilityContractIdentifier(uint64(chainID))
	// 		_, seqErr = deployOCR3Contract(qualifier, uint64(selector), allChainsCLDEnvironment, memoryDatastore)
	// 		if seqErr != nil {
	// 			return nil, fmt.Errorf("failed to deploy EVM OCR3 contract for chainID %d, selector %d: %w", chainID, selector, seqErr)
	// 		}
	// 	}
	// }
	// if consensusV2AddrFlag {
	// 	_, seqErr = deployOCR3Contract("capability_consensus", homeChainSelector, allChainsCLDEnvironment, memoryDatastore)
	// 	if seqErr != nil {
	// 		return nil, fmt.Errorf("failed to deploy Consensus V2 OCR3 contract %w", seqErr)
	// 	}
	// }
	// allChainsCLDEnvironment.DataStore = memoryDatastore.Seal()

	// ocr3Addr := mustGetAddress(memoryDatastore, homeChainSelector, keystone_changeset.OCR3Capability.String(), input.ContractVersions[keystone_changeset.OCR3Capability.String()], "capability_ocr3")
	// testLogger.Info().Msgf("Deployed OCR3 %s contract on chain %d at %s", input.ContractVersions[keystone_changeset.OCR3Capability.String()], homeChainSelector, ocr3Addr)

	// donTimeAddr := mustGetAddress(memoryDatastore, homeChainSelector, keystone_changeset.OCR3Capability.String(), input.ContractVersions[keystone_changeset.OCR3Capability.String()], "DONTime")
	// testLogger.Info().Msgf("Deployed OCR3 %s (DON Time) contract on chain %d at %s", input.ContractVersions[keystone_changeset.OCR3Capability.String()], homeChainSelector, donTimeAddr)

	// wfRegAddr := mustGetAddress(memoryDatastore, homeChainSelector, keystone_changeset.WorkflowRegistry.String(), input.ContractVersions[keystone_changeset.WorkflowRegistry.String()], "")
	// testLogger.Info().Msgf("Deployed Workflow Registry %s contract on chain %d at %s", input.ContractVersions[keystone_changeset.WorkflowRegistry.String()], homeChainSelector, wfRegAddr)

	// capRegAddr := mustGetAddress(memoryDatastore, homeChainSelector, keystone_changeset.CapabilitiesRegistry.String(), input.ContractVersions[keystone_changeset.CapabilitiesRegistry.String()], "")
	// testLogger.Info().Msgf("Deployed Capabilities Registry %s contract on chain %d at %s", input.ContractVersions[keystone_changeset.CapabilitiesRegistry.String()], homeChainSelector, capRegAddr)

	// var vaultOCR3CommonAddr common.Address
	// if vaultOCR3AddrFlag {
	// 	vaultOCR3Addr := mustGetAddress(memoryDatastore, homeChainSelector, keystone_changeset.OCR3Capability.String(), input.ContractVersions[keystone_changeset.OCR3Capability.String()], "capability_vault")
	// 	testLogger.Info().Msgf("Deployed OCR3 %s (Vault) contract on chain %d at %s", input.ContractVersions[keystone_changeset.OCR3Capability.String()], homeChainSelector, vaultOCR3Addr)
	// 	vaultOCR3CommonAddr = common.HexToAddress(vaultOCR3Addr)
	// }

	// evmOCR3CommonAddresses := make(map[uint64]common.Address)
	// if evmOCR3AddrFlag {
	// 	for chainID := range chainsWithEVMCapability {
	// 		qualifier := ks_contracts_op.CapabilityContractIdentifier(uint64(chainID))
	// 		// we have deployed OCR3 contract for each EVM chain on the registry chain to avoid a situation when more than 1 OCR contract (of any type) has the same address
	// 		// because that violates a DB constraint for offchain reporting jobs
	// 		evmOCR3Addr := mustGetAddress(memoryDatastore, homeChainSelector, keystone_changeset.OCR3Capability.String(), "1.0.0", qualifier)
	// 		testLogger.Info().Msgf("Deployed EVM OCR3 contract on chainID: %d, selector: %d, at: %s", chainID, homeChainSelector, evmOCR3Addr)
	// 		evmOCR3CommonAddresses[homeChainSelector] = common.HexToAddress(evmOCR3Addr)
	// 	}
	// }
	// var consensusV2OCR3CommonAddr common.Address
	// if consensusV2AddrFlag {
	// 	consensusV2OCR3Addr := mustGetAddress(memoryDatastore, homeChainSelector, keystone_changeset.OCR3Capability.String(), input.ContractVersions[keystone_changeset.OCR3Capability.String()], "capability_consensus")
	// 	testLogger.Info().Msgf("Deployed Consensus V2 OCR3 %s contract on chain %d at %s", input.ContractVersions[keystone_changeset.OCR3Capability.String()], homeChainSelector, consensusV2OCR3Addr)
	// 	consensusV2OCR3CommonAddr = common.HexToAddress(consensusV2OCR3Addr)
	// }

	// for _, forwarderSelector := range evmForwardersSelectors {
	// 	forwarderAddr := mustGetAddress(memoryDatastore, forwarderSelector, keystone_changeset.KeystoneForwarder.String(), input.ContractVersions[keystone_changeset.KeystoneForwarder.String()], "")
	// 	testLogger.Info().Msgf("Deployed Forwarder %s contract on chain %d at %s", input.ContractVersions[keystone_changeset.KeystoneForwarder.String()], forwarderSelector, forwarderAddr)
	// }

	// for _, forwarderSelector := range solForwardersSelectors {
	// 	forwarderAddr := mustGetAddress(memoryDatastore, forwarderSelector, ks_sol.ForwarderContract.String(), input.ContractVersions[ks_sol.ForwarderContract.String()], ks_sol.DefaultForwarderQualifier)
	// 	forwarderStateAddr := mustGetAddress(memoryDatastore, forwarderSelector, ks_sol.ForwarderState.String(), input.ContractVersions[ks_sol.ForwarderState.String()], ks_sol.DefaultForwarderQualifier)
	// 	testLogger.Info().Msgf("Deployed Forwarder %s contract on Solana chain %d at %s state %s", input.ContractVersions[ks_sol.ForwarderContract.String()], forwarderSelector, forwarderAddr, forwarderStateAddr)
	// }

	deployKeystoneContractsOutput, deployErr := crecontracts.DeployKeystoneContracts(
		ctx,
		testLogger,
		singleFileLogger,
		crecontracts.DeployKeystoneContractsInput{
			CldfBlockchains:           startBlockchainsOutput.BlockChains,
			CtfBlockchains:            startBlockchainsOutput.BlockChainOutputs,
			ContractVersions:          input.ContractVersions,
			WithV2Registries:          input.WithV2Registries,
			CapabilitiesAwareNodeSets: input.CapabilitiesAwareNodeSets,
		},
	)
	if deployErr != nil {
		return nil, pkgerrors.Wrap(deployErr, "failed to deploy Keystone contracts")
	}

	// // get chainIDs, they'll be used for identifying ETH keys and Forwarder addresses
	// // and also for creating the CLD environment
	// evmChainIDs := make([]int, 0)
	// bcOuts := make(map[uint64]*cre.WrappedBlockchainOutput)
	// sethClients := make(map[uint64]*seth.Client)
	// solClients := make(map[uint64]*solrpc.Client)
	// solChainIDs := make([]string, 0)
	// for _, bcOut := range startBlockchainsOutput.BlockChainOutputs {
	// 	if bcOut.SolChain != nil {
	// 		sel := bcOut.SolChain.ChainSelector
	// 		bcOuts[sel] = bcOut
	// 		solClients[sel] = bcOut.SolClient
	// 		bcOuts[sel].ChainSelector = sel
	// 		bcOuts[sel].SolChain = bcOut.SolChain
	// 		bcOuts[sel].SolChain.ArtifactsDir = bcOut.SolChain.ArtifactsDir
	// 		solChainIDs = append(solChainIDs, bcOut.SolChain.ChainID)
	// 		continue
	// 	}
	// 	bcOuts[bcOut.ChainSelector] = bcOut
	// 	evmChainIDs = append(evmChainIDs, libc.MustSafeInt(bcOut.ChainID))
	// 	sethClients[bcOut.ChainSelector] = bcOut.SethClient
	// }

	// Translate node input to structure required further down the road and put as much information
	// as we have at this point in labels. It will be used to generate node configs
	topology, updatedNodeSets, topoErr := BuildTopology(
		homeChainSelector,
		input.CapabilitiesAwareNodeSets,
		input.InfraInput,
		startBlockchainsOutput.BlockChainOutputs,
		deployKeystoneContractsOutput.Env.ExistingAddresses, //nolint:staticcheck // won't migrate now
		deployKeystoneContractsOutput.Env.DataStore,
		input.Capabilities,
		input.CapabilityConfigs,
		input.CopyCapabilityBinaries,
	)
	if topoErr != nil {
		return nil, pkgerrors.Wrap(topoErr, "failed to build topology")
	}

	fmt.Print(libformat.PurpleText("%s", stageGen.WrapAndNext("DONs configuration prepared in %.2f seconds", stageGen.Elapsed().Seconds())))

	wfPool := pond.NewResultPool[*cre.WorkflowRegistryOutput](1)
	wfTask := wfPool.SubmitErr(func() (*cre.WorkflowRegistryOutput, error) {
		fmt.Print(libformat.PurpleText("\n---> [BACKGROUND] Starting Workflow Registry Contract Configuration\n\n"))
		defer fmt.Print(libformat.PurpleText("\n---> [BACKGROUND] Finished Workflow Registry Contract Configuration\n\n"))

		return crecontracts.ConfigureWorkflowRegistry(
			ctx,
			testLogger,
			singleFileLogger,
			&cre.WorkflowRegistryInput{
				ContractAddress: common.HexToAddress(crecontracts.MustGetAddressFromDataStore(deployKeystoneContractsOutput.Env.DataStore, homeChainSelector, keystone_changeset.WorkflowRegistry.String(), input.ContractVersions[keystone_changeset.WorkflowRegistry.String()], "")),
				ChainSelector:   homeChainOutput.ChainSelector,
				CldEnv:          deployKeystoneContractsOutput.Env,
				AllowedDonIDs:   []uint64{topology.WorkflowDONID},
				WorkflowOwners:  []common.Address{homeChainOutput.SethClient.MustGetRootKeyAddress()},
			},
		)
	})

	// JOB DISTRIBUTOR + JOBS (creation and distribution to nodes)
	fmt.Print(libformat.PurpleText("%s", stageGen.Wrap("Starting Job Distributor, DONs and creating Jobs with Job Distributor")))

	jdOutput, nodeSetOutput, jobsSeqErr := StartDONsAndJD(
		testLogger,
		input.JdInput,
		homeChainOutput.BlockchainOutput,
		topology,
		input.InfraInput,
		updatedNodeSets,
	)
	if jobsSeqErr != nil {
		return nil, pkgerrors.Wrap(jobsSeqErr, "failed to setup jobs")
	}

	// append the nodeset output, so that later it can be stored in the cached output, so that we can use the environment again without running setup
	for idx, nsOut := range nodeSetOutput {
		input.CapabilitiesAwareNodeSets[idx].Out = nsOut.Output
	}

	for idx, bcOut := range startBlockchainsOutput.BlockChainOutputs {
		input.BlockchainsInput[idx].Out = bcOut.BlockchainOutput
	}

	// append the jd output, so that later it can be stored in the cached output, so that we can use the environment again without running setup
	input.JdInput.Out = jdOutput

	// Prepare the CLD environment that's required by the keystone changeset
	// Ugly glue hack ¯\_(ツ)_/¯
	fullCldInput := &cre.FullCLDEnvironmentInput{
		JdOutput:          jdOutput,
		BlockchainOutputs: startBlockchainsOutput.BlockChainOutputs,
		// BlockchainOutputs: startBlockchainsOutput.BlockChainOutputs,
		// SethClients:       sethClients,
		// SolClients:        solClients,
		NodeSetOutput:     nodeSetOutput,
		ExistingAddresses: deployKeystoneContractsOutput.Env.ExistingAddresses, //nolint:staticcheck // won't migrate now
		Datastore:         deployKeystoneContractsOutput.Env.DataStore,
		Topology:          topology,
		OperationsBundle:  deployKeystoneContractsOutput.Env.OperationsBundle,
	}

	fullCldOutput, cldErr := libdevenv.BuildFullCLDEnvironment(ctx, singleFileLogger, fullCldInput, insecure.NewCredentials())
	if cldErr != nil {
		return nil, pkgerrors.Wrap(cldErr, "failed to build full CLD environment")
	}

	createJobsInput := CreateJobsWithJdOpInput{}

	jobSpecFactoryFunctions := make([]cre.JobSpecFn, 0)
	for _, capability := range input.Capabilities {
		jobSpecFactoryFunctions = append(jobSpecFactoryFunctions, capability.JobSpecFn())
	}

	// Deprecated, use Capabilities instead
	jobSpecFactoryFunctions = append(jobSpecFactoryFunctions, input.JobSpecFactoryFunctions...)

	createJobsDeps := CreateJobsWithJdOpDeps{
		Logger:                    testLogger,
		SingleFileLogger:          singleFileLogger,
		HomeChainBlockchainOutput: homeChainOutput.BlockchainOutput,
		AddressBook:               deployKeystoneContractsOutput.Env.ExistingAddresses, //nolint:staticcheck // won't migrate now
		JobSpecFactoryFunctions:   jobSpecFactoryFunctions,
		FullCLDEnvOutput:          fullCldOutput,
		CapabilitiesAwareNodeSets: input.CapabilitiesAwareNodeSets,
		InfraInput:                &input.InfraInput,
		CapabilitiesConfigs:       input.CapabilityConfigs,
		Capabilities:              input.Capabilities,
	}
	_, createJobsErr := operations.ExecuteOperation(deployKeystoneContractsOutput.Env.OperationsBundle, CreateJobsWithJdOp, createJobsDeps, createJobsInput)
	if createJobsErr != nil {
		return nil, pkgerrors.Wrap(createJobsErr, "failed to create jobs with Job Distributor")
	}

	// CAUTION: It is crucial to configure OCR3 jobs on nodes before configuring the workflow contracts.
	// Wait for OCR listeners to be ready before setting the configuration.
	// If the ConfigSet event is missed, OCR protocol will not start.
	fmt.Print(libformat.PurpleText("%s", stageGen.WrapAndNext("Jobs created in %.2f seconds", stageGen.Elapsed().Seconds())))

	// This operation cannot execute in the background, because it uses master private keys and we want to avoid nonce issues
	fmt.Print(libformat.PurpleText("%s", stageGen.Wrap("Preparing Chainlink Node funding")))
	preFundingOutput, prefundErr := operations.ExecuteOperation(fullCldOutput.Environment.OperationsBundle, PrepareCLNodesFundingOp, PrepareFundCLNodesOpDeps{
		TestLogger:        testLogger,
		Env:               fullCldOutput.Environment,
		BlockchainOutputs: startBlockchainsOutput.BlockChainOutputs,
		DonTopology:       fullCldOutput.DonTopology,
	}, PrepareFundCLNodesOpInput{FundingPerChainFamilyForEachNode: map[string]uint64{
		"evm":    10000000000000000, // 0.01 ETH
		"solana": 50_000_000,        // 0.05 SOL
	}})
	if prefundErr != nil {
		return nil, pkgerrors.Wrap(prefundErr, "failed to prepare funding of CL nodes")
	}
	fmt.Print(libformat.PurpleText("%s", stageGen.WrapAndNext("Chainlink Node funding prepared in %.2f seconds", stageGen.Elapsed().Seconds())))

	bkgErrPool := pond.NewPool(10)
	fundNodesTaskErr := bkgErrPool.SubmitErr(func() error {
		fmt.Print(libformat.PurpleText("\n---> [BACKGROUND] Funding Chainlink nodes\n\n"))
		defer fmt.Print(libformat.PurpleText("\n---> [BACKGROUND] Finished Funding Chainlink nodes\n\n"))

		_, fundErr := operations.ExecuteOperation(fullCldOutput.Environment.OperationsBundle, FundCLNodesOp, FundCLNodesOpDeps{
			TestLogger:        testLogger,
			Env:               fullCldOutput.Environment,
			BlockchainOutputs: startBlockchainsOutput.BlockChainOutputs,
			DonTopology:       fullCldOutput.DonTopology,
		}, FundCLNodesOpInput{
			FundingAmountPerChainFamily: preFundingOutput.Output.FundingPerChainFamilyForEachNode,
			PrivateKeyPerChainFamily:    preFundingOutput.Output.PrivateKeysPerChainFamily,
		})

		return fundErr
	})

	fmt.Print(libformat.PurpleText("%s", stageGen.Wrap("Waiting for Log Poller to start tracking OCR3 contract")))

	for idx, nodeSetOut := range nodeSetOutput {
		if !flags.HasFlag(updatedNodeSets[idx].ComputedCapabilities, cre.ConsensusCapability) || !flags.HasFlag(updatedNodeSets[idx].ComputedCapabilities, cre.VaultCapability) {
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
		if waitErr := eg.Wait(); waitErr != nil {
			return nil, pkgerrors.Wrap(waitErr, "failed to wait for ConfigWatcher health check")
		}
	}

	fmt.Print(libformat.PurpleText("%s", stageGen.WrapAndNext("Log Poller started in %.2f seconds", stageGen.Elapsed().Seconds())))

	// wait for log poller filters to be registered in the background, because we don't need it them at this stage yet

	filterRegErr := bkgErrPool.SubmitErr(func() error {
		fmt.Print(libformat.PurpleText("\n---> [BACKGROUND] Waiting for all nodes to have expected LogPoller filters registered\n\n"))
		defer fmt.Print(libformat.PurpleText("\n---> [BACKGROUND] All nodes have expected LogPoller filters registered\n\n"))

		return crecontracts.WaitForWorkflowRegistryFiltersRegistration(testLogger, singleFileLogger, input.InfraInput.Type, homeChainOutput.ChainID, fullCldOutput.DonTopology, updatedNodeSets)
	})

	fmt.Print(libformat.PurpleText("%s", stageGen.Wrap("Configuring OCR3 and Keystone contracts")))

	// Configure the Forwarder, OCR3 and Capabilities contracts
	configureKeystoneInput := cre.ConfigureKeystoneInput{
		ChainSelector:               homeChainSelector,
		CldEnv:                      fullCldOutput.Environment,
		Topology:                    topology,
		CapabilitiesRegistryAddress: ptr.Ptr(crecontracts.MustGetAddressFromMemoryDataStore(deployKeystoneContractsOutput.MemoryDataStore, homeChainSelector, keystone_changeset.CapabilitiesRegistry.String(), input.ContractVersions[keystone_changeset.CapabilitiesRegistry.String()], "")),
		OCR3Address:                 ptr.Ptr(crecontracts.MustGetAddressFromMemoryDataStore(deployKeystoneContractsOutput.MemoryDataStore, homeChainSelector, keystone_changeset.OCR3Capability.String(), input.ContractVersions[keystone_changeset.OCR3Capability.String()], crecontracts.OCR3ContractQualifier)),
		DONTimeAddress:              ptr.Ptr(crecontracts.MustGetAddressFromMemoryDataStore(deployKeystoneContractsOutput.MemoryDataStore, homeChainSelector, keystone_changeset.OCR3Capability.String(), input.ContractVersions[keystone_changeset.OCR3Capability.String()], crecontracts.DONTimeContractQualifier)),
		VaultOCR3Address:            crecontracts.MightGetAddressFromMemoryDataStore(deployKeystoneContractsOutput.MemoryDataStore, homeChainSelector, keystone_changeset.OCR3Capability.String(), input.ContractVersions[keystone_changeset.OCR3Capability.String()], crecontracts.VaultOCR3ContractQualifier),
		EVMOCR3Addresses:            evmOCR3AddressesFromDataStore(startBlockchainsOutput.BlockChainOutputs, updatedNodeSets, deployKeystoneContractsOutput.MemoryDataStore, homeChainSelector),
		ConsensusV2OCR3Address:      crecontracts.MightGetAddressFromMemoryDataStore(deployKeystoneContractsOutput.MemoryDataStore, homeChainSelector, keystone_changeset.OCR3Capability.String(), input.ContractVersions[keystone_changeset.OCR3Capability.String()], crecontracts.ConsensusV2ContractQualifier),
		NodeSets:                    input.CapabilitiesAwareNodeSets,
	}

	if input.OCR3Config != nil {
		configureKeystoneInput.OCR3Config = *input.OCR3Config
	} else {
		ocr3Config, ocr3ConfigErr := crecontracts.DefaultOCR3Config(topology)
		if ocr3ConfigErr != nil {
			return nil, pkgerrors.Wrap(ocr3ConfigErr, "failed to generate default OCR3 config")
		}
		configureKeystoneInput.OCR3Config = *ocr3Config
	}

	if input.DONTimeConfig != nil {
		configureKeystoneInput.DONTimeConfig = *input.DONTimeConfig
	} else {
		donTimeConfig, donTimeConfigErr := crecontracts.DefaultOCR3Config(topology)
		donTimeConfig.DeltaRoundMillis = 0 // Fastest rounds possible
		if donTimeConfigErr != nil {
			return nil, pkgerrors.Wrap(donTimeConfigErr, "failed to generate default DON Time config")
		}
		configureKeystoneInput.DONTimeConfig = *donTimeConfig
	}

	ocr3Config, ocr3ConfigErr := crecontracts.DefaultOCR3Config(topology)
	if ocr3ConfigErr != nil {
		return nil, pkgerrors.Wrap(ocr3ConfigErr, "failed to generate default OCR3 config")
	}
	configureKeystoneInput.VaultOCR3Config = *ocr3Config

	defaultOcr3Config, defaultOcr3ConfigErr := crecontracts.DefaultOCR3Config(topology)
	if defaultOcr3ConfigErr != nil {
		return nil, pkgerrors.Wrap(defaultOcr3ConfigErr, "failed to generate default OCR3 config for EVM")
	}
	configureKeystoneInput.EVMOCR3Config = *defaultOcr3Config
	configureKeystoneInput.ConsensusV2OCR3Config = *defaultOcr3Config

	capabilitiesContractFactoryFunctions := make([]cre.CapabilityRegistryConfigFn, 0)
	for _, capability := range input.Capabilities {
		configFn := capability.CapabilityRegistryV1ConfigFn()
		if input.WithV2Registries {
			configFn = capability.CapabilityRegistryV2ConfigFn()
		}

		capabilitiesContractFactoryFunctions = append(capabilitiesContractFactoryFunctions, configFn)
	}

	// Deprecated, use Capabilities instead
	capabilitiesContractFactoryFunctions = append(capabilitiesContractFactoryFunctions, input.CapabilitiesContractFactoryFunctions...)

	keystoneErr := crecontracts.ConfigureKeystone(configureKeystoneInput, capabilitiesContractFactoryFunctions)
	if keystoneErr != nil {
		return nil, pkgerrors.Wrap(keystoneErr, "failed to configure keystone contracts")
	}

	fmt.Print(libformat.PurpleText("%s", stageGen.WrapAndNext("OCR3 and Keystone contracts configured in %.2f seconds", stageGen.Elapsed().Seconds())))

	fmt.Print(libformat.PurpleText("%s", stageGen.Wrap("Writing bootstrapping data into disk (address book, data store, etc...)")))

	artifactPath, artifactErr := DumpArtifact(
		deployKeystoneContractsOutput.MemoryDataStore.AddressRefStore,
		deployKeystoneContractsOutput.Env.ExistingAddresses, //nolint:staticcheck // won't migrate now
		*jdOutput,
		*fullCldOutput.DonTopology,
		fullCldOutput.Environment.Offchain,
		capabilitiesContractFactoryFunctions,
		input.CapabilitiesAwareNodeSets,
	)
	if artifactErr != nil {
		testLogger.Error().Err(artifactErr).Msg("failed to generate artifact")
		fmt.Print(libformat.PurpleText("%s", stageGen.WrapAndNext("Failed to write bootstrapping data into disk in %.2f seconds", stageGen.Elapsed().Seconds())))
	} else {
		testLogger.Info().Msgf("Environment artifact saved to %s", artifactPath)
		fmt.Print(libformat.PurpleText("%s", stageGen.WrapAndNext("Wrote bootstrapping data into disk in %.2f seconds", stageGen.Elapsed().Seconds())))
	}

	wfPool.StopAndWait()
	workflowRegistryConfigurationOutput, wfRegistrationErr := wfTask.Wait()
	if wfRegistrationErr != nil {
		return nil, pkgerrors.Wrap(wfRegistrationErr, "failed to configure workflow registry")
	}

	bkgErrPool.StopAndWait()
	if err := fundNodesTaskErr.Wait(); err != nil {
		return nil, pkgerrors.Wrap(err, "failed to fund chainlink nodes")
	}
	if err := filterRegErr.Wait(); err != nil {
		return nil, pkgerrors.Wrap(err, "failed while waiting for log poller filters to be registered")
	}

	return &SetupOutput{
		WorkflowRegistryConfigurationOutput: workflowRegistryConfigurationOutput, // pass to caller, so that it can be optionally attached to TestConfig and saved to disk
		BlockchainOutput:                    startBlockchainsOutput.BlockChainOutputs,
		DonTopology:                         fullCldOutput.DonTopology,
		NodeOutput:                          nodeSetOutput,
		CldEnvironment:                      fullCldOutput.Environment,
	}, nil
}

func evmOCR3AddressesFromDataStore(blockchains []*cre.WrappedBlockchainOutput, nodeSets []*cre.CapabilitiesAwareNodeSet, ds *datastore.MemoryDataStore, homeChainSelector uint64) map[uint64]common.Address {
	chainsWithEVMCapability := crecontracts.ChainsWithEVMCapability(blockchains, nodeSets)
	evmOCR3CommonAddresses := make(map[uint64]common.Address)
	for chainID := range chainsWithEVMCapability {
		qualifier := ks_contracts_op.CapabilityContractIdentifier(uint64(chainID))
		// we have deployed OCR3 contract for each EVM chain on the registry chain to avoid a situation when more than 1 OCR contract (of any type) has the same address
		// because that violates a DB constraint for offchain reporting jobs
		evmOCR3Addr := crecontracts.MustGetAddressFromMemoryDataStore(ds, homeChainSelector, keystone_changeset.OCR3Capability.String(), "1.0.0", qualifier)
		evmOCR3CommonAddresses[homeChainSelector] = evmOCR3Addr
	}

	return evmOCR3CommonAddresses
}

func CreateJobDistributor(input *jd.Input) (*jd.Output, error) {
	if os.Getenv("CI") == "true" {
		jdImage := ctfconfig.MustReadEnvVar_String(E2eJobDistributorImageEnvVarName)
		jdVersion := os.Getenv(E2eJobDistributorVersionEnvVarName)
		input.Image = fmt.Sprintf("%s:%s", jdImage, jdVersion)
	}

	jdOutput, err := jd.NewJD(input)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to create new job distributor")
	}

	return jdOutput, nil
}

func mergeJobSpecSlices(from, to cre.DonsToJobSpecs) {
	for fromDonID, fromJobSpecs := range from {
		if _, ok := to[fromDonID]; !ok {
			to[fromDonID] = make([]*jobv1.ProposeJobRequest, 0)
		}
		to[fromDonID] = append(to[fromDonID], fromJobSpecs...)
	}
}

// func deployOCR3Contract(qualifier string, selector uint64, env *cldf.Environment, ds datastore.MutableDataStore) (*ks_contracts_op.DeployOCR3ContractSequenceOutput, error) {
// 	ocr3DeployReport, err := operations.ExecuteSequence(
// 		env.OperationsBundle,
// 		ks_contracts_op.DeployOCR3ContractsSequence,
// 		ks_contracts_op.DeployOCR3ContractSequenceDeps{
// 			Env: env,
// 		},
// 		ks_contracts_op.DeployOCR3ContractSequenceInput{
// 			ChainSelector: selector,
// 			Qualifier:     qualifier,
// 		},
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to deploy OCR3 contract '%s' on chain %d: %w", qualifier, selector, err)
// 	}
// 	// TODO: CRE-742 remove address book
// 	if err = env.ExistingAddresses.Merge(ocr3DeployReport.Output.AddressBook); err != nil { //nolint:staticcheck // won't migrate now
// 		return nil, fmt.Errorf("failed to merge address book with OCR3 contract address for '%s' on chain %d: %w", qualifier, selector, err)
// 	}
// 	if err = ds.Merge(ocr3DeployReport.Output.Datastore); err != nil {
// 		return nil, fmt.Errorf("failed to merge datastore with OCR3 contract address for '%s' on chain %d: %w", qualifier, selector, err)
// 	}
// 	return &ocr3DeployReport.Output, nil
// }
