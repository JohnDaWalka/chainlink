package changeset

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipdeployment "github.com/smartcontractkit/chainlink/deployment/ccip"
)

// tokenChainConfig stores the token address and the desired admin of the token on a given chain.
type tokenChainConfig struct {
	TokenAddress common.Address
	DesiredAdmin common.Address
}

// createOperationForChain creates the operation required to set the proposed admin of the token to the correct address.
func createOperationForChain(chainSelector uint64, state ccipdeployment.CCIPOnChainState, tokenConfig tokenChainConfig) (*timelock.BatchChainOperation, error) {
	zeroAddress := common.Address{}
	chainState, ok := state.Chains[chainSelector]
	if !ok {
		return nil, fmt.Errorf("no state found for chain with selector %d", chainSelector)
	}
	if chainState.TokenAdminRegistry == nil {
		return nil, fmt.Errorf("chain with selector %d has no token admin registry", chainSelector)
	}
	if chainState.Timelock == nil {
		return nil, fmt.Errorf("chain with selector %d has no timelock", chainSelector)
	}
	tokenConfigOnRegistry, err := chainState.TokenAdminRegistry.GetTokenConfig(nil, tokenConfig.TokenAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch config from token admin registry with address %s for token with address %s: %w", chainState.TokenAdminRegistry.Address(), tokenConfig.TokenAddress, err)
	}
	if tokenConfigOnRegistry.PendingAdministrator == tokenConfig.DesiredAdmin {
		return nil, nil // If the pending administrator is already the desired admin, we have no action to take.
	}
	if tokenConfigOnRegistry.Administrator == zeroAddress {
		// If the administrator has never been set, we need to propose the desired admin.
		proposeAdministrator, err := chainState.TokenAdminRegistry.ProposeAdministrator(
			deployment.SimTransactOpts(),
			tokenConfig.TokenAddress,
			tokenConfig.DesiredAdmin,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create transaction to propose %s as admin for token with address %s: %w", tokenConfig.DesiredAdmin, tokenConfig.TokenAddress, err)
		}
		return &timelock.BatchChainOperation{
			ChainIdentifier: mcms.ChainIdentifier(chainSelector),
			Batch: []mcms.Operation{
				{
					To:    chainState.TokenAdminRegistry.Address(),
					Data:  proposeAdministrator.Data(),
					Value: big.NewInt(0),
				},
			},
		}, nil
	} else if tokenConfigOnRegistry.Administrator == chainState.Timelock.Address() {
		// If the current administrator is the Timelock, we need to transfer the admin role.
		transferAdminRole, err := chainState.TokenAdminRegistry.TransferAdminRole(
			deployment.SimTransactOpts(),
			tokenConfig.TokenAddress,
			tokenConfig.DesiredAdmin,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create transaction to transfer admin role of token with address %s to %s: %w", tokenConfig.TokenAddress, tokenConfig.DesiredAdmin, err)
		}
		return &timelock.BatchChainOperation{
			ChainIdentifier: mcms.ChainIdentifier(chainSelector),
			Batch: []mcms.Operation{
				{
					To:    chainState.TokenAdminRegistry.Address(),
					Data:  transferAdminRole.Data(),
					Value: big.NewInt(0),
				},
			},
		}, nil
	}

	// Otherwise, the current administrator is unknown and there is no action that we can take.
	return nil, fmt.Errorf("token with address %s has unknown admin (%s), unable to transfer admin role", tokenConfig.TokenAddress, tokenConfigOnRegistry.Administrator)
}

// ProposeOrTransferAdminRoleChangeSet creates an MCMS proposal to propose or transfer the admin role for a token.
func ProposeOrTransferAdminRoleChangeSet(env deployment.Environment, state ccipdeployment.CCIPOnChainState, tokenName string, tokenChainConfigs map[uint64]tokenChainConfig, delay time.Duration) (deployment.ChangesetOutput, error) {
	var batches []timelock.BatchChainOperation
	var multiErr error

	type result struct {
		err   error
		batch *timelock.BatchChainOperation
	}
	results := make(chan result, len(tokenChainConfigs))

	for chainSelector, config := range tokenChainConfigs {
		go func(chainSelector uint64, tokenChainConfig tokenChainConfig) {
			batchChainOperation, err := createOperationForChain(chainSelector, state, tokenChainConfig)
			results <- result{err: err, batch: batchChainOperation}
		}(chainSelector, config)
	}

	for i := 0; i < len(tokenChainConfigs); i++ {
		result := <-results
		if result.err != nil {
			multiErr = errors.Join(multiErr, result.err)
		}
		if result.batch != nil {
			batches = append(batches, *result.batch)
		}
	}

	if multiErr != nil {
		return deployment.ChangesetOutput{}, multiErr
	}

	proposal, err := ccipdeployment.BuildProposalFromBatches(state, batches, fmt.Sprintf("Proposal that updates the proposed administrator for %s across multiple chains", tokenName), delay)
	if err != nil || proposal == nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
	}
	return deployment.ChangesetOutput{
		Proposals: []timelock.MCMSWithTimelockProposal{*proposal},
	}, nil
}
