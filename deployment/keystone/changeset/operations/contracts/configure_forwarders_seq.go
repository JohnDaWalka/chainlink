package contracts

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"
	"github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
)

type ConfigureForwardersSeqDeps struct {
	Env      *cldf.Environment
	Registry *capabilities_registry.CapabilitiesRegistry
}

type ConfigureForwardersSeqInput struct {
	WFDonName        string
	WFNodeIDs        []string
	RegistryChainSel uint64

	// MCMSConfig is optional. If non-nil, the changes will be proposed using MCMS.
	MCMSConfig *changeset.MCMSConfig
	// Chains is optional. Defines chains for which request will be executed. If empty, runs for all available chains.
	Chains map[uint64]struct{}
}

func (i ConfigureForwardersSeqInput) Validate() error {
	if len(i.WFNodeIDs) == 0 {
		return fmt.Errorf("WFNodeIDs must not be empty")
	}
	return nil
}

func (i ConfigureForwardersSeqInput) UseMCMS() bool {
	return i.MCMSConfig != nil
}

type ConfigureForwardersSeqOutput struct {
	MCMSTimelockProposals []mcms.TimelockProposal
}

var ConfigureForwardersSeq = operations.NewSequence[ConfigureForwardersSeqInput, ConfigureForwardersSeqOutput, ConfigureForwardersSeqDeps](
	"configure-forwarders-seq",
	semver.MustParse("1.0.0"),
	"Configure Keystone Forwarders",
	func(b operations.Bundle, deps ConfigureForwardersSeqDeps, input ConfigureForwardersSeqInput) (ConfigureForwardersSeqOutput, error) {
		wfDON, err := internal.NewRegisteredDon(*deps.Env, internal.RegisteredDonConfig{
			NodeIDs:          input.WFNodeIDs,
			Name:             input.WFDonName,
			RegistryChainSel: input.RegistryChainSel,
			Registry:         deps.Registry,
		})
		if err != nil {
			return ConfigureForwardersSeqOutput{}, fmt.Errorf("configure-forwarders-seq failed: %w", err)
		}

		// configure forwarders for all EVM chains
		evmChain := deps.Env.BlockChains.EVMChains()
		opPerChain := make(map[uint64]mcmstypes.BatchOperation)
		forwarderContracts := make(map[uint64]*changeset.OwnedContract[*forwarder.KeystoneForwarder])

		for _, chain := range evmChain {
			if _, shouldInclude := input.Chains[chain.Selector]; len(input.Chains) > 0 && !shouldInclude {
				continue
			}

			addressesRefs := deps.Env.DataStore.Addresses().Filter(
				datastore.AddressRefByChainSelector(chain.Selector),
				datastore.AddressRefByType(datastore.ContractType(changeset.KeystoneForwarder)),
			)
			if len(addressesRefs) == 0 {
				return ConfigureForwardersSeqOutput{}, fmt.Errorf("configure-forwarders-seq failed: no KeystoneForwarder contract found for chain selector %d", chain.Selector)
			}

			for _, addrRef := range addressesRefs {
				contract, err := changeset.GetOwnedContractV2[*forwarder.KeystoneForwarder](deps.Env.DataStore.Addresses(), chain, addrRef.Address)
				if err != nil {
					return ConfigureForwardersSeqOutput{}, fmt.Errorf("configure-forwarders-seq failed: failed to get KeystoneForwarder contract for chain selector %d: %w", chain.Selector, err)
				}

				fwrReport, err := operations.ExecuteOperation(b, ConfigureForwarderOp, ConfigureForwarderOpDeps{
					Env:      deps.Env,
					Chain:    &chain,
					Contract: contract.Contract,
				}, ConfigureForwarderOpInput{
					Dons:    []internal.RegisteredDon{*wfDON},
					UseMCMS: input.UseMCMS(),
				})
				if err != nil {
					return ConfigureForwardersSeqOutput{}, fmt.Errorf("configure-forwarders-seq failed for chain selector %d: %w", chain.Selector, err)
				}

				opPerChain[chain.Selector] = fwrReport.Output.BatchOperation
				forwarderContracts[chain.Selector] = contract
			}
		}

		var out ConfigureForwardersSeqOutput
		if input.UseMCMS() {
			if len(opPerChain) == 0 {
				return out, fmt.Errorf("configure-forwarders-seq failed: no operations generated for MCMS")
			}

			for chainSelector, op := range opPerChain {
				fwr, ok := forwarderContracts[chainSelector]
				if !ok {
					return out, fmt.Errorf("configure-forwarders-seq failed: expected configured forwarder address for chain selector %d", chainSelector)
				}
				if fwr.McmsContracts == nil {
					return out, fmt.Errorf("configure-forwarders-seq failed: expected forwarder contract %s to be owned by MCMS for chain selector %d", fwr.Contract.Address(), chainSelector)
				}
				timelocksPerChain := map[uint64]string{
					chainSelector: fwr.McmsContracts.Timelock.Address().Hex(),
				}
				proposerMCMSes := map[uint64]string{
					chainSelector: fwr.McmsContracts.ProposerMcm.Address().Hex(),
				}
				inspector, err := proposalutils.McmsInspectorForChain(*deps.Env, chainSelector)
				if err != nil {
					return out, err
				}
				inspectorPerChain := map[uint64]mcmssdk.Inspector{
					chainSelector: inspector,
				}

				proposal, err := proposalutils.BuildProposalFromBatchesV2(
					*deps.Env,
					timelocksPerChain,
					proposerMCMSes,
					inspectorPerChain,
					[]mcmstypes.BatchOperation{op},
					"proposal to set forwarder config",
					proposalutils.TimelockConfig{
						MinDelay: input.MCMSConfig.MinDuration,
					},
				)
				if err != nil {
					return out, fmt.Errorf("configure-forwarders-seq failed: failed to build proposal: %w", err)
				}
				out.MCMSTimelockProposals = append(out.MCMSTimelockProposals, *proposal)
			}
		}

		return out, nil
	},
)

type ConfigureForwarderOpDeps struct {
	Env      *cldf.Environment
	Chain    *evm.Chain
	Contract *forwarder.KeystoneForwarder
}

type ConfigureForwarderOpInput struct {
	Dons    []internal.RegisteredDon
	UseMCMS bool
}

type ConfigureForwarderOpOutput struct {
	BatchOperation mcmstypes.BatchOperation
}

var ConfigureForwarderOp = operations.NewOperation[ConfigureForwarderOpInput, ConfigureForwarderOpOutput, ConfigureForwarderOpDeps](
	"configure-forwarder-op",
	semver.MustParse("1.0.0"),
	"Configure Keystone Forwarder",
	func(b operations.Bundle, deps ConfigureForwarderOpDeps, input ConfigureForwarderOpInput) (ConfigureForwarderOpOutput, error) {
		ops, err := internal.ConfigureForwarder(b.Logger, *deps.Chain, deps.Contract, input.Dons, input.UseMCMS)
		if err != nil {
			return ConfigureForwarderOpOutput{}, fmt.Errorf("configure-forwarder-op failed: failed to configure forwarder for chain selector %d: %w", deps.Chain.Selector, err)
		}
		return ConfigureForwarderOpOutput{BatchOperation: ops[deps.Chain.Selector]}, nil
	},
)
