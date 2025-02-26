package mcmsnew

import (
	"fmt"
	"maps"
	"slices"

	"github.com/gagliardetto/solana-go"
	mcmBindings "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/mcm"
	routerBindings "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	"github.com/smartcontractkit/chainlink/deployment"
	mcms "github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcmssolanasdk "github.com/smartcontractkit/mcms/sdk/solana"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	state "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

var solanaAddress = mcmssolanasdk.ContractAddress

// TransferOwnership transfer the ownership of a Solana "contract"
func TransferOwnership(
	env deployment.Environment,
	cfg commontypes.TransferToMCMSWithTimelockConfigSolana,
) (deployment.ChangesetOutput, error) {
	mcmsState, err := state.MaybeLoadMCMSWithTimelockStateSolana(env, slices.Collect(maps.Keys(env.SolChains)))
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	batches := []mcmstypes.BatchOperation{}
	timelocks := map[uint64]string{}
	proposers := map[uint64]string{}
	inspectors := map[uint64]mcmssdk.Inspector{}
	for chainSelector, contractsToTransfer := range cfg.ContractsByChain {
		solChain, ok := env.SolChains[chainSelector]
		if !ok {
			return deployment.ChangesetOutput{}, fmt.Errorf("solana chain not found in environment (selector: %v)", chainSelector)
		}
		chainState, ok := mcmsState[chainSelector]
		if !ok {
			return deployment.ChangesetOutput{}, fmt.Errorf("chain state not found for selector: %v", chainSelector)
		}

		timelocks[chainSelector] = solanaAddress(chainState.TimelockProgram, mcmssolanasdk.PDASeed(chainState.TimelockSeed))
		proposers[chainSelector] = solanaAddress(chainState.McmProgram, mcmssolanasdk.PDASeed(chainState.ProposerMcmSeed))
		inspectors[chainSelector] = mcmssolanasdk.NewInspector(solChain.Client)

		var transactions []mcmstypes.Transaction
		for _, contract := range contractsToTransfer {
			contractTransactions, err := transferOwnershipMCMSTransaction(chainState, chainSelector, solChain, contract,
				chainState.TimelockProgram, chainState.TimelockSeed)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to create mcms transactions: %w", err)
			}

			transactions = append(transactions, contractTransactions...)
		}

		batches = append(batches, mcmstypes.BatchOperation{
			ChainSelector: mcmstypes.ChainSelector(chainSelector),
			Transactions:  transactions,
		})
	}

	proposal, err := proposalutils.BuildProposalFromBatchesV2(env, timelocks, proposers, inspectors, batches,
		"proposal to transfer ownership of contracts to timelock", cfg.MinDelay)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
	}

	return deployment.ChangesetOutput{MCMSTimelockProposals: []mcms.TimelockProposal{*proposal}}, nil
}

// transferOwnershipRouter transfers ownership of the router to the timelock.
func transferOwnershipMCMSTransaction(
	chainState *state.MCMSWithTimelockStateSolana,
	chainSelector uint64,
	solChain deployment.SolChain,
	contract commontypes.OwnableSolanaContract,
	timelockProgramID solana.PublicKey,
	timelockSeed state.PDASeed,
) ([]mcmstypes.Transaction, error) {
	timelockSignerPDA := state.GetTimelockSignerPDA(timelockProgramID, timelockSeed)
	auth := solChain.DeployerKey.PublicKey()

	instruction, err := transferOwnershipInstruction(contract.Seed, timelockSignerPDA, contract.OwnerPDA, auth)
	if err != nil {
		return nil, fmt.Errorf("failed to build transfer ownership instruction: %w", err)
	}
	transferMCMSTx, err := mcmssolanasdk.NewTransactionFromInstruction(instruction, string(contract.Type), []string{})
	if err != nil {
		return nil, fmt.Errorf("failed to build mcms transaction from transfer ownership instruction: %w", err)
	}

	instruction, err = acceptOwnershipInstruction(contract.Seed, contract.OwnerPDA, timelockSignerPDA)
	if err != nil {
		return nil, fmt.Errorf("failed to build accept ownership instruction: %w", err)
	}
	acceptMCMSTx, err := mcmssolanasdk.NewTransactionFromInstruction(instruction, string(contract.Type), []string{})
	if err != nil {
		return nil, fmt.Errorf("failed to build mcms transaction from accept ownership instruction: %w", err)
	}

	return []mcmstypes.Transaction{transferMCMSTx, acceptMCMSTx}, nil
}

func transferOwnershipInstruction(
	seed state.PDASeed, proposedOwner, ownerPDA, auth solana.PublicKey,
) (solana.Instruction, error) {
	if (seed == state.PDASeed{}) {
		return routerBindings.NewTransferOwnershipInstruction(proposedOwner, ownerPDA, auth).ValidateAndBuild()
	}
	return mcmBindings.NewTransferOwnershipInstruction(seed, proposedOwner, ownerPDA, auth).ValidateAndBuild()
}

func acceptOwnershipInstruction(seed state.PDASeed, ownerPDA, auth solana.PublicKey) (solana.Instruction, error) {
	if (seed == state.PDASeed{}) {
		return routerBindings.NewAcceptOwnershipInstruction(ownerPDA, auth).ValidateAndBuild()
	}
	return mcmBindings.NewAcceptOwnershipInstruction(seed, ownerPDA, auth).ValidateAndBuild()
}
