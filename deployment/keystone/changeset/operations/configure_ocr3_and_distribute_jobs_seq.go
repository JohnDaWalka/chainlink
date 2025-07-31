package operations

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/mcms"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/jobs"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"
	opjobs "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/jobs"
)

type ConfigureOCR3AndDistributeJobsSeqDeps struct {
	Env   *cldf.Environment
	Nodes []*nodev1.Node

	// DonCapabilities is used to add capabilities to capabilities registry.
	DonCapabilities []internal.DonCapabilities // externally sourced based on the environment
}

type ConfigureOCR3AndDistributeJobsSeqInput struct {
	RegistryChainSel        uint64
	MCMSConfig              *changeset.MCMSConfig
	RegistryContractAddress *common.Address
	OracleConfig            internal.OracleConfig
	DONs                    []contracts.ConfigureKeystoneDON

	// The following are needed for the OCR3 job spec distribution
	DomainKey            string
	EnvironmentLabel     string
	DONName              string
	ChainSelectorEVM     uint64
	ChainSelectorAptos   uint64
	BootstrapperOCR3Urls []string
}

func (c ConfigureOCR3AndDistributeJobsSeqInput) UseMCMS() bool {
	return c.MCMSConfig != nil
}

type ConfigureOCR3AndDistributeJobsSeqOutput struct {
	Specs                 []jobs.OCR3JobConfigSpec
	Addresses             datastore.AddressRefStore
	MCMSTimelockProposals []mcms.TimelockProposal
	BatchOperation        *mcmstypes.BatchOperation
}

var ConfigureOCR3AndDistributeJobsSeq = operations.NewSequence[
	ConfigureOCR3AndDistributeJobsSeqInput,
	ConfigureOCR3AndDistributeJobsSeqOutput,
	ConfigureOCR3AndDistributeJobsSeqDeps,
](
	"configure-ocr3-and-distribute-jobs-seq",
	semver.MustParse("1.0.0"),
	"Configure OCR3 and Distribute Jobs",
	func(b operations.Bundle, deps ConfigureOCR3AndDistributeJobsSeqDeps, input ConfigureOCR3AndDistributeJobsSeqInput) (ConfigureOCR3AndDistributeJobsSeqOutput, error) {
		chain, ok := deps.Env.BlockChains.EVMChains()[input.RegistryChainSel]
		if !ok {
			return ConfigureOCR3AndDistributeJobsSeqOutput{}, fmt.Errorf("registry chain selector %d does not exist in environment", input.RegistryChainSel)
		}

		capabilitiesRegistry, err := changeset.GetOwnedContractV2[*capabilities_registry.CapabilitiesRegistry](deps.Env.DataStore.Addresses(), chain, input.RegistryContractAddress.Hex())
		if err != nil {
			return ConfigureOCR3AndDistributeJobsSeqOutput{}, fmt.Errorf("failed to get capabilities registry contract: %w", err)
		}
		if input.UseMCMS() && capabilitiesRegistry.McmsContracts == nil {
			return ConfigureOCR3AndDistributeJobsSeqOutput{}, fmt.Errorf("capabilities registry contract %s is not owned by MCMS", capabilitiesRegistry.Contract.Address())
		}

		donInfos, err := internal.DonInfos(deps.DonCapabilities, deps.Env.Offchain)
		if err != nil {
			return ConfigureOCR3AndDistributeJobsSeqOutput{}, fmt.Errorf("failed to get don infos: %w", err)
		}

		donToCapabilities, err := internal.MapDonsToCaps(capabilitiesRegistry.Contract, donInfos)
		if err != nil {
			return ConfigureOCR3AndDistributeJobsSeqOutput{}, fmt.Errorf("failed to map dons to capabilities: %w", err)
		}

		capReport, err := operations.ExecuteOperation(b, contracts.AddCapabilitiesOp, contracts.AddCapabilitiesOpDeps{
			Chain:             chain,
			Contract:          capabilitiesRegistry.Contract,
			DonToCapabilities: donToCapabilities,
		}, contracts.AddCapabilitiesOpInput{
			ChainID:         input.RegistryChainSel,
			ContractAddress: capabilitiesRegistry.Contract.Address(),
			UseMCMS:         input.UseMCMS(),
		})
		if err != nil {
			return ConfigureOCR3AndDistributeJobsSeqOutput{}, fmt.Errorf("failed to add capabilities to capabilities registry: %w", err)
		}

		ocr3ContractReport, err := operations.ExecuteOperation(b, contracts.DeployOCR3Op, contracts.DeployOCR3OpDeps{
			Env: deps.Env,
		}, contracts.DeployOCR3OpInput{
			ChainSelector: input.RegistryChainSel,
		})
		if err != nil {
			return ConfigureOCR3AndDistributeJobsSeqOutput{}, fmt.Errorf("failed to deploy OCR3 contract: %w", err)
		}
		ocr3Addresses, err := ocr3ContractReport.Output.Addresses.Fetch()
		if err != nil {
			return ConfigureOCR3AndDistributeJobsSeqOutput{}, fmt.Errorf("failed to fetch OCR3 contract addresses: %w", err)
		}
		if len(ocr3Addresses) == 0 {
			return ConfigureOCR3AndDistributeJobsSeqOutput{}, fmt.Errorf("no OCR3 capability address found for chain selector %d", input.RegistryChainSel)
		}
		ocr3Address := common.HexToAddress(ocr3Addresses[0].Address)

		configOCR3ContractReport, err := operations.ExecuteOperation(
			b,
			contracts.ConfigureOCR3Op,
			contracts.ConfigureOCR3OpDeps{
				Env:      deps.Env,
				Registry: capabilitiesRegistry.Contract,
			},
			contracts.ConfigureOCR3OpInput{
				ContractAddress:  &ocr3Address,
				RegistryChainSel: input.RegistryChainSel,
				DONs:             input.DONs,
				Config:           &input.OracleConfig,
				DryRun:           false,
				MCMSConfig:       input.MCMSConfig,
			},
		)
		if err != nil {
			return ConfigureOCR3AndDistributeJobsSeqOutput{}, fmt.Errorf("failed to configure OCR3 contract: %w", err)
		}

		distributionReport, err := operations.ExecuteSequence(b, opjobs.DistributeOCRJobSpecSeq, opjobs.DistributeOCRJobSpecSeqDeps{
			Nodes:    deps.Nodes,
			Offchain: deps.Env.Offchain,
		}, opjobs.DistributeOCRJobSpecSeqInput{
			DomainKey:            input.DomainKey,
			EnvironmentLabel:     input.EnvironmentLabel,
			DONName:              input.DONName,
			ContractID:           ocr3Address.Hex(),
			ChainSelectorEVM:     input.ChainSelectorEVM,
			ChainSelectorAptos:   input.ChainSelectorAptos,
			BootstrapperOCR3Urls: input.BootstrapperOCR3Urls,
		})
		if err != nil {
			return ConfigureOCR3AndDistributeJobsSeqOutput{}, fmt.Errorf("failed to distribute OCR3 job specs: %w", err)
		}

		return ConfigureOCR3AndDistributeJobsSeqOutput{
			Specs:                 distributionReport.Output.Specs,
			Addresses:             ocr3ContractReport.Output.Addresses,
			BatchOperation:        capReport.Output.BatchOperation,
			MCMSTimelockProposals: configOCR3ContractReport.Output.MCMSTimelockProposals,
		}, nil
	},
)
