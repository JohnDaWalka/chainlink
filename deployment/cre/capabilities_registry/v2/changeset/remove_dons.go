package changeset

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/pkg"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
)

var _ cldf.ChangeSetV2[RemoveDONsInput] = RemoveDONs{}

type RemoveDONsInput struct {
	Domain            string   `json:"domain" yaml:"domain"`
	Zone              string   `json:"zone" yaml:"zone"`
	DONNames          []string `json:"don_names" yaml:"don_names"`
	RegistryChainSel  uint64   `json:"registry_chain_selector" yaml:"registry_chain_selector"`
	RegistryQualifier string   `json:"registry_qualifier,omitempty" yaml:"registry_qualifier,omitempty"`

	MCMSConfig *ocr3.MCMSConfig `json:"mcms_config,omitempty" yaml:"mcms_config,omitempty"`
}

type RemoveDONs struct{}

func (u RemoveDONs) VerifyPreconditions(_ cldf.Environment, input RemoveDONsInput) error {
	if len(input.DONNames) == 0 {
		return fmt.Errorf("must specify at least one DON name")
	}
	return nil
}

func (u RemoveDONs) Apply(e cldf.Environment, input RemoveDONsInput) (cldf.ChangesetOutput, error) {
	// Get MCMS contracts if needed
	var mcmsContracts *commonchangeset.MCMSWithTimelockState
	if input.MCMSConfig != nil {
		var err error
		mcmsContracts, err = strategies.GetMCMSContracts(e, input.RegistryChainSel, input.RegistryQualifier)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get MCMS contracts: %w", err)
		}
	}

	registryRef := pkg.GetCapRegV2AddressRefKey(input.RegistryChainSel, input.RegistryQualifier)

	chainSel := registryRef.ChainSelector()
	registryAddressRef, err := e.DataStore.Addresses().Get(registryRef)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get registry address: %w", err)
	}
	chain, ok := e.BlockChains.EVMChains()[chainSel]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain not found for selector %d", chainSel)
	}

	capReg, err := capabilities_registry_v2.NewCapabilitiesRegistry(
		common.HexToAddress(registryAddressRef.Address), chain.Client,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create CapabilitiesRegistry: %w", err)
	}

	report, err := operations.ExecuteOperation(e.OperationsBundle, contracts.RemoveDONs, contracts.RemoveDONsDeps{
		Env:                  &e,
		MCMSContracts:        mcmsContracts,
		CapabilitiesRegistry: capReg,
	}, contracts.RemoveDONsInput{
		ChainSelector: input.RegistryChainSel,
		Domain:        input.Domain,
		Zone:          input.Zone,
		DONNames:      input.DONNames,
		MCMSConfig:    input.MCMSConfig,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to remove DONs: %w", err)
	}

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: report.Output.Proposals,
		Reports:               []operations.Report[any, any]{report.ToGenericReport()},
	}, nil
}
