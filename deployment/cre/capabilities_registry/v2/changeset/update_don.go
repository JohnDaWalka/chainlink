package changeset

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

var _ cldf.ChangeSetV2[UpdateDONInput] = UpdateDON{}

type UpdateDONInput struct {
	RegistryQualifier string `json:"registry_qualifier" yaml:"registry_qualifier"`
	RegistryChainSel  uint64 `json:"registry_chain_sel" yaml:"registry_chain_sel"`

	// P2PIDs are the peer ids that compose the don. Optional, only provided if the DON composition is changing.
	P2PIDs            []p2pkey.PeerID              `json:"p2p_ids" yaml:"p2p_ids"`
	CapabilityConfigs []contracts.CapabilityConfig `json:"capability_configs" yaml:"capability_configs"`

	// DonName to update, this is required
	DonName string `json:"don_name" yaml:"don_name"`

	NewDonName string `json:"new_don_name" yaml:"new_don_name"`

	// F is the fault tolerance level
	// if omitted, the existing value fetched from the registry is used
	F uint8 `json:"f" yaml:"f"`

	// IsPublic indicates whether the DON is public or private
	// If omitted, the existing value fetched from the registry is used
	IsPublic bool `json:"is_public" yaml:"is_public"`

	// Force indicates whether to force the update even if we cannot validate that all forwarder contracts are ready to accept the new configure version.
	// This is very dangerous, and could break the whole platform if the forwarders are not ready. Be very careful with this option.
	Force bool `json:"force" yaml:"force"`

	MCMSConfig *ocr3.MCMSConfig `json:"mcms_config,omitempty" yaml:"mcms_config,omitempty"`
}

type UpdateDON struct{}

func (u UpdateDON) VerifyPreconditions(_ cldf.Environment, input UpdateDONInput) error {
	if input.DonName == "" {
		return errors.New("missing DON name")
	}
	return nil
}

func (u UpdateDON) Apply(e cldf.Environment, input UpdateDONInput) (cldf.ChangesetOutput, error) {
	registryRef := pkg.GetCapRegV2AddressRefKey(input.RegistryChainSel, input.RegistryQualifier)

	registryAddressRef, err := e.DataStore.Addresses().Get(registryRef)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get registry address: %w", err)
	}

	chainSel := registryRef.ChainSelector()
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

	updateDonReport, err := operations.ExecuteOperation(
		e.OperationsBundle,
		contracts.UpdateDON,
		contracts.UpdateDONDeps{
			Env:                  &e,
			CapabilitiesRegistry: capReg,
		},
		contracts.UpdateDONInput{
			ChainSelector:     chainSel,
			P2PIDs:            input.P2PIDs,
			CapabilityConfigs: input.CapabilityConfigs,
			DonName:           input.DonName,
			NewDonName:        input.NewDonName,
			F:                 input.F,
			IsPublic:          input.IsPublic,
			Force:             input.Force,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to update don: %w", err)
	}

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: updateDonReport.Output.Proposals,
		Reports:               []operations.Report[any, any]{updateDonReport.ToGenericReport()},
	}, nil
}
