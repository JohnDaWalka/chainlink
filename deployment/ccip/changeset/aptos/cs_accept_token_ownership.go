package aptos

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/mcms"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
	aptosmcms "github.com/smartcontractkit/mcms/sdk/aptos"
	mcmstypes "github.com/smartcontractkit/mcms/types"
)

var _ deployment.ChangeSetV2[config.AcceptTokenOwnershipConfig] = AcceptTokenOwnership{}

// AcceptTokenOwnership deploys token pools and sets up tokens on lanes
type AcceptTokenOwnership struct{}

func (cs AcceptTokenOwnership) VerifyPreconditions(env deployment.Environment, cfg config.AcceptTokenOwnershipConfig) error {
	return nil
}

func (cs AcceptTokenOwnership) Apply(env deployment.Environment, cfg config.AcceptTokenOwnershipConfig) (deployment.ChangesetOutput, error) {
	state, err := changeset.LoadOnchainState(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load Aptos onchain state: %w", err)
	}
	aptosChain := env.AptosChains[cfg.ChainSelector]
	deps := operation.AptosDeps{
		AptosChain:       aptosChain,
		CCIPOnChainState: state,
		OnChainState:     state.AptosChains[cfg.ChainSelector],
	}
	txs := []mcmstypes.Transaction{}

	// Deploy Aptos token and pool
	deployReport, err := operations.ExecuteOperation(env.OperationsBundle, operation.AcceptTokenOwnershipOp, deps, cfg.TokenAddress)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	txs = append(txs, deployReport.Output)

	mcmsOps := mcmstypes.BatchOperation{
		ChainSelector: mcmstypes.ChainSelector(deps.AptosChain.Selector),
		Transactions:  txs,
	}

	// Generate Aptos MCMS proposals
	proposal, err := utils.GenerateProposal(
		aptosChain.Client,
		state.AptosChains[cfg.ChainSelector].MCMSAddress,
		cfg.ChainSelector,
		[]mcmstypes.BatchOperation{mcmsOps},
		"Deploy Aptos MCMS and CCIP",
		aptosmcms.TimelockRoleProposer,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate MCMS proposal for Aptos chain %d: %w", cfg.ChainSelector, err)
	}

	return deployment.ChangesetOutput{
		MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
	}, nil
}
