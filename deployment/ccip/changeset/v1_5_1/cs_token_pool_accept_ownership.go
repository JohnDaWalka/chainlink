package v1_5_1

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/token_pool"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/deployergroup"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

// AcceptOwnershipInput defines the input for the AcceptOwnershipChangeset.
// The contract address is the address of the contract to accept ownership of (must implement OpenZeppelin Ownable2Step).
type AcceptOwnershipInputTokenPool struct {
	// MCMS is the Timelock config for the chain, which will be the new owner.
	MCMS *proposalutils.TimelockConfig
	// Contracts is a map of chainSelector to a slice of contract addresses to accept ownership for.
	Contracts map[uint64][]common.Address
}

// AcceptOwnershipChangeset accepts ownership of contracts (using OpenZeppelin 2-step Ownable) for the MCMS Timelock on each chain.
func AcceptOwnershipChangesetTokenPool(env cldf.Environment, input AcceptOwnershipInputTokenPool) (cldf.ChangesetOutput, error) {
	// Validate input parameters
	if input.Contracts == nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("contracts must be defined")
	}
	if len(input.Contracts) == 0 {
		return cldf.ChangesetOutput{}, fmt.Errorf("at least one contract must be defined")
	}

	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	deployerGroup := deployergroup.NewDeployerGroup(env, state, input.MCMS).WithDeploymentContext("accept ownership of contracts")

	for chainSelector, contractAddresses := range input.Contracts {
		chain := env.BlockChains.EVMChains()[chainSelector]
		opts, err := deployerGroup.GetDeployer(chainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get deployer for %s: %w", chain, err)
		}
		for _, contractAddr := range contractAddresses {
			// We assume the contract implements Ownable2Step (acceptOwnership)
			// Here, we use the token_pool ABI, but in general, this could be any Ownable2Step contract.
			client := chain.Client
			ownable, err := token_pool.NewTokenPool(contractAddr, client)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to bind contract at %s on chain %s: %w", contractAddr.Hex(), chain, err)
			}
			_, err = ownable.AcceptOwnership(opts)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to call acceptOwnership on %s (chain %s): %w", contractAddr.Hex(), chain, err)
			}
		}
	}

	return deployerGroup.Enact()
}
