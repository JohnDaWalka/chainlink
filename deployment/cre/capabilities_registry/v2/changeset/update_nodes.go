package changeset

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

var _ cldf.ChangeSetV2[AddCapabilitiesInput] = AddCapabilities{}

type UpdateNodesInput struct {
	RegistryChainSel  uint64 `json:"registry_chain_sel" yaml:"registry_chain_sel"`
	RegistryQualifier string `json:"registry_qualifier" yaml:"registry_qualifier"`

	MCMSConfig *ocr3.MCMSConfig `json:"mcms_config" yaml:"mcms_config"`

	Nodes map[string]contracts.NodeConfig `json:"nodes" yaml:"nodes"`
}

type UpdateNodes struct{}

func (u UpdateNodes) VerifyPreconditions(_ cldf.Environment, config UpdateNodesInput) error {
	if len(config.Nodes) == 0 {
		return errors.New("nodes is required")
	}

	return nil
}

func (u UpdateNodes) Apply(e cldf.Environment, config UpdateNodesInput) (cldf.ChangesetOutput, error) {
	var mcmsContracts *commonchangeset.MCMSWithTimelockState
	if config.MCMSConfig != nil {
		var err error
		mcmsContracts, err = strategies.GetMCMSContracts(e, config.RegistryChainSel, emptyQualifier)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get MCMS contracts: %w", err)
		}
	}

	registryRef := pkg.GetCapRegV2AddressRefKey(config.RegistryChainSel, config.RegistryQualifier)

	registryAddressRef, err := e.DataStore.Addresses().Get(registryRef)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get registry address: %w", err)
	}

	chain, ok := e.BlockChains.EVMChains()[config.RegistryChainSel]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain not found for selector %d", config.RegistryChainSel)
	}

	capReg, err := capabilities_registry_v2.NewCapabilitiesRegistry(
		common.HexToAddress(registryAddressRef.Address), chain.Client,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create CapabilitiesRegistry: %w", err)
	}

	report, err := operations.ExecuteOperation(
		e.OperationsBundle,
		contracts.UpdateNodes,
		contracts.UpdateNodesDeps{
			Env:                  &e,
			MCMSContracts:        mcmsContracts,
			CapabilitiesRegistry: capReg,
		},
		contracts.UpdateNodesInput{
			ChainSelector: config.RegistryChainSel,
			NodesUpdates:  config.Nodes,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	return cldf.ChangesetOutput{
		Reports:               []operations.Report[any, any]{report.ToGenericReport()},
		MCMSTimelockProposals: report.Output.Proposals,
	}, nil
}
