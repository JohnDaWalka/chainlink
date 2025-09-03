package contracts

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	cap_reg_v2 "github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
	wf_reg_v2 "github.com/smartcontractkit/chainlink/deployment/cre/workflow_registry/v2/changeset/operations/contracts"
)

type DeployContractsSequenceDeps struct {
	Env *deployment.Environment
}

type (
	EVMChainID uint64
	Selector   uint64
)

// inputs and outputs have to be serializable, and must not contain sensitive data

type DeployContractsSequenceInput struct {
	RegistryChainSelector uint64
	ForwardersSelectors   []uint64
	DeployVaultOCR3       bool
	DeployEVMOCR3         bool
	EVMChainIDs           map[EVMChainID]Selector
	DeployConsensusOCR3   bool

	// WithV2Contracts if true will deploy Capability Registry and Workflow Registry V2
	WithV2Contracts bool
}

type DeployContractsSequenceOutput struct {
	// Not sure if we can serialize the address book without modifications, but whatever is returned needs to be serializable.
	// This could also be the address datastore instead.
	AddressBook deployment.AddressBook
	Datastore   datastore.DataStore // Keeping the address store for backward compatibility, as not everything has been migrated to address book
}

func updateAddresses(addr datastore.MutableAddressRefStore, as datastore.AddressRefStore, sourceAB deployment.AddressBook, ab deployment.AddressBook) error {
	addresses, err := as.Fetch()
	if err != nil {
		return err
	}
	for _, a := range addresses {
		if err := addr.Add(a); err != nil {
			return err
		}
	}

	return sourceAB.Merge(ab)
}

// DeployContractsSequence is a sequence that deploys the Keystone contracts (OCR3, Capabilities Registry, Workflow Registry, Keystone Forwarder).
var DeployContractsSequence = operations.NewSequence(
	"deploy-keystone-contracts-seq",
	semver.MustParse("1.0.0"),
	"Deploy Keystone Contracts (BalanceReader, OCR3, DON Time, Vault-OCR3, EVM-OCR3, Capabilities Registry, Workflow Registry, Keystone Forwarder)",
	func(b operations.Bundle, deps DeployContractsSequenceDeps, input DeployContractsSequenceInput) (output DeployContractsSequenceOutput, err error) {
		ab := deployment.NewMemoryAddressBook()
		as := datastore.NewMemoryDataStore()

		// Capabilities Registry contract
		if input.WithV2Contracts {
			v2Report, err := operations.ExecuteOperation(b, cap_reg_v2.DeployCapabilitiesRegistry, cap_reg_v2.DeployCapabilitiesRegistryDeps(deps), cap_reg_v2.DeployCapabilitiesRegistryInput{
				ChainSelector: input.RegistryChainSelector,
			})
			if err != nil {
				return DeployContractsSequenceOutput{}, err
			}

			out, err := toV1Output(v2Report.Output)
			if err != nil {
				return DeployContractsSequenceOutput{}, err
			}

			if err = updateAddresses(as.Addresses(), out.Addresses, ab, out.AddressBook); err != nil {
				return DeployContractsSequenceOutput{}, err
			}
		} else {
			capabilitiesRegistryDeployReport, err := operations.ExecuteOperation(b, DeployCapabilityRegistryOp, DeployCapabilityRegistryOpDeps(deps), DeployCapabilityRegistryInput{ChainSelector: input.RegistryChainSelector})
			if err != nil {
				return DeployContractsSequenceOutput{}, err
			}
			err = updateAddresses(as.Addresses(), capabilitiesRegistryDeployReport.Output.Addresses, ab, capabilitiesRegistryDeployReport.Output.AddressBook)
			if err != nil {
				return DeployContractsSequenceOutput{}, err
			}
		}

		// OCR3 Contract
		ocr3DeployReport, err := operations.ExecuteOperation(b, DeployOCR3Op, DeployOCR3OpDeps(deps), DeployOCR3OpInput{ChainSelector: input.RegistryChainSelector, Qualifier: "capability_ocr3"})
		if err != nil {
			return DeployContractsSequenceOutput{}, err
		}
		err = updateAddresses(as.Addresses(), ocr3DeployReport.Output.Addresses, ab, ocr3DeployReport.Output.AddressBook)
		if err != nil {
			return DeployContractsSequenceOutput{}, err
		}

		// Workflow Registry contract
		if input.WithV2Contracts {
			v2Report, err := operations.ExecuteOperation(b, wf_reg_v2.DeployWorkflowRegistryOp, wf_reg_v2.DeployWorkflowRegistryOpDeps(deps), wf_reg_v2.DeployWorkflowRegistryOpInput{
				ChainSelector: input.RegistryChainSelector,
			})
			if err != nil {
				return DeployContractsSequenceOutput{}, err
			}

			out, err := toV1Output(v2Report.Output)
			if err != nil {
				return DeployContractsSequenceOutput{}, err
			}

			if err = updateAddresses(as.Addresses(), out.Addresses, ab, out.AddressBook); err != nil {
				return DeployContractsSequenceOutput{}, err
			}
		} else {
			workflowRegistryDeployReport, err := operations.ExecuteOperation(b, DeployWorkflowRegistryOp, DeployWorkflowRegistryOpDeps(deps), DeployWorkflowRegistryInput{ChainSelector: input.RegistryChainSelector})
			if err != nil {
				return DeployContractsSequenceOutput{}, err
			}
			err = updateAddresses(as.Addresses(), workflowRegistryDeployReport.Output.Addresses, ab, workflowRegistryDeployReport.Output.AddressBook)
			if err != nil {
				return DeployContractsSequenceOutput{}, err
			}
		}

		// Keystone Forwarder contract
		keystoneForwarderDeployReport, err := operations.ExecuteSequence(b, DeployKeystoneForwardersSequence, DeployKeystoneForwardersSequenceDeps(deps), DeployKeystoneForwardersInput{Targets: input.ForwardersSelectors})
		if err != nil {
			return DeployContractsSequenceOutput{}, err
		}
		err = updateAddresses(as.Addresses(), keystoneForwarderDeployReport.Output.Addresses, ab, keystoneForwarderDeployReport.Output.AddressBook)
		if err != nil {
			return DeployContractsSequenceOutput{}, err
		}

		// DON Time Contract - Copy of OCR3Capability
		donTimeDeployReport, err := operations.ExecuteOperation(b, DeployOCR3Op, DeployOCR3OpDeps(deps), DeployOCR3OpInput{ChainSelector: input.RegistryChainSelector, Qualifier: "DONTime"})
		if err != nil {
			return DeployContractsSequenceOutput{}, err
		}
		err = updateAddresses(as.Addresses(), donTimeDeployReport.Output.Addresses, ab, donTimeDeployReport.Output.AddressBook)
		if err != nil {
			return DeployContractsSequenceOutput{}, err
		}

		if input.DeployVaultOCR3 {
			// Vault OCR3 Contract
			vaultOCR3DeployReport, err := operations.ExecuteOperation(b, DeployOCR3Op, DeployOCR3OpDeps(deps), DeployOCR3OpInput{ChainSelector: input.RegistryChainSelector, Qualifier: "capability_vault"})
			if err != nil {
				return DeployContractsSequenceOutput{}, err
			}
			err = updateAddresses(as.Addresses(), vaultOCR3DeployReport.Output.Addresses, ab, vaultOCR3DeployReport.Output.AddressBook)
			if err != nil {
				return DeployContractsSequenceOutput{}, err
			}
		}

		if input.DeployEVMOCR3 {
			for chainID := range input.EVMChainIDs {
				// EVM cap OCR3 Contract
				qualifier := CapabilityContractIdentifier(uint64(chainID))
				// deploy OCR3 contract for each EVM chain on the registry chain to avoid a situation when more than 1 OCR contract (of any type) has the same address
				// because that violates a DB constraint for offchain reporting jobs
				// this can be removed once https://smartcontract-it.atlassian.net/browse/PRODCRE-804 is done and we can deploy OCR3 contract for each EVM chain on that chain
				evmOCR3DeployReport, err := operations.ExecuteOperation(b, DeployOCR3Op, DeployOCR3OpDeps(deps), DeployOCR3OpInput{ChainSelector: input.RegistryChainSelector, Qualifier: qualifier})
				if err != nil {
					return DeployContractsSequenceOutput{}, err
				}
				err = updateAddresses(as.Addresses(), evmOCR3DeployReport.Output.Addresses, ab, evmOCR3DeployReport.Output.AddressBook)
				if err != nil {
					return DeployContractsSequenceOutput{}, err
				}
			}
		}

		if input.DeployConsensusOCR3 {
			evmOCR3DeployReport, err := operations.ExecuteOperation(b, DeployOCR3Op, DeployOCR3OpDeps(deps), DeployOCR3OpInput{ChainSelector: input.RegistryChainSelector, Qualifier: "capability_consensus"})
			if err != nil {
				return DeployContractsSequenceOutput{}, err
			}
			err = updateAddresses(as.Addresses(), evmOCR3DeployReport.Output.Addresses, ab, evmOCR3DeployReport.Output.AddressBook)
			if err != nil {
				return DeployContractsSequenceOutput{}, err
			}
		}

		return DeployContractsSequenceOutput{
			AddressBook: ab,
			Datastore:   as.Seal(),
		}, nil
	},
)

func CapabilityContractIdentifier(chainID uint64) string {
	return fmt.Sprintf("capability_evm_%d", chainID)
}

type DeprecatedOutput struct {
	Addresses   datastore.AddressRefStore
	AddressBook deployment.AddressBook
}

// toV1Output transforms a v2 output to a common output format that uses the deprecated
// address book.
func toV1Output(in any) (DeprecatedOutput, error) {
	ab := deployment.NewMemoryAddressBook()
	ds := datastore.NewMemoryDataStore()
	labels := deployment.NewLabelSet()
	var r datastore.AddressRef

	switch v := in.(type) {
	case cap_reg_v2.DeployCapabilitiesRegistryOutput:
		r = datastore.AddressRef{
			ChainSelector: v.ChainSelector,
			Address:       v.Address,
			Type:          datastore.ContractType(v.Type),
			Version:       semver.MustParse(v.Version),
			Qualifier:     v.Qualifier,
			Labels:        datastore.NewLabelSet(v.Labels...),
		}
		for _, l := range v.Labels {
			labels.Add(l)
		}
	case wf_reg_v2.DeployWorkflowRegistryOpOutput:
		r = datastore.AddressRef{
			ChainSelector: v.ChainSelector,
			Address:       v.Address,
			Type:          datastore.ContractType(v.Type),
			Version:       semver.MustParse(v.Version),
			Qualifier:     v.Qualifier,
			Labels:        datastore.NewLabelSet(v.Labels...),
		}
		for _, l := range v.Labels {
			labels.Add(l)
		}
	default:
		return DeprecatedOutput{}, fmt.Errorf("unsupported input type for transform: %T", in)
	}

	if err := ds.Addresses().Add(r); err != nil {
		return DeprecatedOutput{}, fmt.Errorf("failed to add address ref: %w", err)
	}

	if err := ab.Save(r.ChainSelector, r.Address, deployment.TypeAndVersion{
		Type:    deployment.ContractType(r.Type),
		Version: *r.Version,
		Labels:  labels,
	}); err != nil {
		return DeprecatedOutput{}, fmt.Errorf("failed to save address to address book: %w", err)
	}

	return DeprecatedOutput{
		Addresses:   ds.Addresses(),
		AddressBook: ab,
	}, nil
}
