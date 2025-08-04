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

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/jobs"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"
	opjobs "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/jobs"
)

type DeployOCR3Capability struct {
	Env *cldf.Environment
}

type DeployOCR3CapabilityInput struct {
	Nodes []opjobs.DistributeOCRJobSpecSeqNode // Node to distribute the job specs to

	RegistryChainSel uint64
	MCMSConfig       *changeset.MCMSConfig
	RegistryRef      datastore.AddressRefKey
	OracleConfig     changeset.OracleConfig
	DONs             []contracts.ConfigureKeystoneDON

	Capabilities []capabilities_registry.CapabilitiesRegistryCapability

	// The following are needed for the OCR3 job spec distribution
	DomainKey            string
	EnvironmentLabel     string
	DONName              string
	ChainSelectorEVM     uint64
	ChainSelectorAptos   uint64
	BootstrapperOCR3Urls []string
	BootstrapCfgs        []opjobs.DistributeBootstrapJobSpecsSeqBootCfg
}

func (c DeployOCR3CapabilityInput) UseMCMS() bool {
	return c.MCMSConfig != nil
}

type DeployOCR3CapabilityOutput struct {
	JobSpecs              []jobs.OCR3JobConfigSpec
	BootstrapSpec         string
	Addresses             datastore.AddressRefStore
	MCMSTimelockProposals []mcms.TimelockProposal
	BatchOperation        *mcmstypes.BatchOperation
}

var DeployOCR3CapabilitySeq = operations.NewSequence[
	DeployOCR3CapabilityInput,
	DeployOCR3CapabilityOutput,
	DeployOCR3Capability,
](
	"deploy-ocr3-capability-seq",
	semver.MustParse("1.0.0"),
	"Deploy OCR3 Capability",
	func(b operations.Bundle, deps DeployOCR3Capability, input DeployOCR3CapabilityInput) (DeployOCR3CapabilityOutput, error) {
		ds := datastore.NewMemoryDataStore()

		chain, ok := deps.Env.BlockChains.EVMChains()[input.RegistryChainSel]
		if !ok {
			return DeployOCR3CapabilityOutput{}, fmt.Errorf("registry chain selector %d does not exist in environment", input.RegistryChainSel)
		}

		capabilitiesRegistry, err := changeset.LoadCapabilityRegistry(chain, *deps.Env, input.RegistryRef)
		if err != nil {
			return DeployOCR3CapabilityOutput{}, fmt.Errorf("failed to get capabilities registry contract: %w", err)
		}
		if input.UseMCMS() && capabilitiesRegistry.McmsContracts == nil {
			return DeployOCR3CapabilityOutput{}, fmt.Errorf("capabilities registry contract %s is not owned by MCMS", capabilitiesRegistry.Contract.Address())
		}

		capReport, err := operations.ExecuteOperation(b, contracts.AddCapabilitiesOp, contracts.AddCapabilitiesOpDeps{
			Chain:    chain,
			Contract: capabilitiesRegistry.Contract,
		}, contracts.AddCapabilitiesOpInput{
			ChainID:         input.RegistryChainSel,
			Capabilities:    input.Capabilities,
			ContractAddress: capabilitiesRegistry.Contract.Address(),
			UseMCMS:         input.UseMCMS(),
		})
		if err != nil {
			return DeployOCR3CapabilityOutput{}, fmt.Errorf("failed to add capabilities to capabilities registry: %w", err)
		}

		ocr3ContractReport, err := operations.ExecuteOperation(b, contracts.DeployOCR3Op, contracts.DeployOCR3OpDeps{
			Env: deps.Env,
		}, contracts.DeployOCR3OpInput{
			ChainSelector: input.RegistryChainSel,
		})
		if err != nil {
			return DeployOCR3CapabilityOutput{}, fmt.Errorf("failed to deploy OCR3 contract: %w", err)
		}
		ocr3Addresses, err := ocr3ContractReport.Output.Addresses.Fetch()
		if err != nil {
			return DeployOCR3CapabilityOutput{}, fmt.Errorf("failed to fetch OCR3 contract addresses: %w", err)
		}
		if len(ocr3Addresses) == 0 {
			return DeployOCR3CapabilityOutput{}, fmt.Errorf("no OCR3 capability address found for chain selector %d", input.RegistryChainSel)
		}
		ocr3Address := common.HexToAddress(ocr3Addresses[0].Address)

		for _, addr := range ocr3Addresses {
			addrErr := ds.Addresses().Add(addr)
			if addrErr != nil {
				return DeployOCR3CapabilityOutput{}, fmt.Errorf("failed to add OCR3 address %s to datastore: %w", addr.Address, addrErr)
			}
		}

		if err := ds.Merge(deps.Env.DataStore); err != nil {
			return DeployOCR3CapabilityOutput{}, fmt.Errorf("failed to merge datastore: %w", err)
		}
		deps.Env.DataStore = ds.Seal()

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
			return DeployOCR3CapabilityOutput{}, fmt.Errorf("failed to configure OCR3 contract: %w", err)
		}

		bootDistributionReport, err := operations.ExecuteSequence(b, opjobs.DistributeBootstrapJobSpecsSeq, opjobs.DistributeBootstrapJobSpecsSeqDeps{
			Offchain: deps.Env.Offchain,
		}, opjobs.DistributeBootstrapJobSpecsSeqInput{
			DONName:          input.DONName,
			DomainKey:        input.DomainKey,
			ContractID:       ocr3Address.Hex(),
			EnvironmentLabel: input.EnvironmentLabel,
			ChainSelectorEVM: input.ChainSelectorEVM,
			BootCfgs:         input.BootstrapCfgs,
		})
		if err != nil {
			return DeployOCR3CapabilityOutput{}, fmt.Errorf("failed to distribute bootstrap job specs: %w", err)
		}

		distributionReport, err := operations.ExecuteSequence(b, opjobs.DistributeOCRJobSpecSeq, opjobs.DistributeOCRJobSpecSeqDeps{
			Offchain: deps.Env.Offchain,
		}, opjobs.DistributeOCRJobSpecSeqInput{
			Nodes:                input.Nodes,
			DomainKey:            input.DomainKey,
			EnvironmentLabel:     input.EnvironmentLabel,
			DONName:              input.DONName,
			ContractID:           ocr3Address.Hex(),
			ChainSelectorEVM:     input.ChainSelectorEVM,
			ChainSelectorAptos:   input.ChainSelectorAptos,
			BootstrapperOCR3Urls: input.BootstrapperOCR3Urls,
		})
		if err != nil {
			return DeployOCR3CapabilityOutput{}, fmt.Errorf("failed to distribute OCR3 job specs: %w", err)
		}

		return DeployOCR3CapabilityOutput{
			JobSpecs:              distributionReport.Output.Specs,
			BootstrapSpec:         bootDistributionReport.Output.Spec,
			Addresses:             ocr3ContractReport.Output.Addresses,
			BatchOperation:        capReport.Output.BatchOperation,
			MCMSTimelockProposals: configOCR3ContractReport.Output.MCMSTimelockProposals,
		}, nil
	},
)
