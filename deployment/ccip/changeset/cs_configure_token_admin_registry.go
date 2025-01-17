package changeset

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	rbac_timelock "github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_admin_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_pool"
)

var _ deployment.ChangeSet[ConfigureTokenAdminRegistryConfig] = ConfigureTokenAdminRegistry

// RegistryConfig defines a token and its state on the token admin registry
type RegistryConfig struct {
	// PoolAddress is the address of the token pool that should be set on the registry.
	PoolAddress common.Address
	// Administrator is the address of the token administrator that should be set on the registry.
	Administrator common.Address
}

func (c RegistryConfig) Validate(ctx context.Context, chain deployment.Chain, state CCIPChainState, useMcms bool, tokenSymbol TokenSymbol) error {
	// Ensure that required fields are populated
	if c.PoolAddress == utils.ZeroAddress {
		return errors.New("pool address must be defined")
	}

	// Ensure that the given pool address and symbol are aligned and known to the environment
	tokenPool, err := GetTokenPoolWithSymbolAndAddress(state, chain, tokenSymbol, c.PoolAddress)
	if err != nil {
		return fmt.Errorf("failed to find token pool on %s with symbol %s and address %s: %w", chain.Name(), tokenSymbol, c.PoolAddress, err)
	}

	// Validate that the token admin registry is owned by the address that will be actioning the transactions (i.e. Timelock or deployer key)
	if err := commoncs.ValidateOwnership(ctx, useMcms, chain.DeployerKey.From, state.Timelock.Address(), state.TokenAdminRegistry); err != nil {
		return fmt.Errorf("token admin registry failed ownership validation on %s: %w", chain.Name(), err)
	}

	// Fetch information about the corresponding token and its state on the registry
	token, err := tokenPool.GetToken(nil)
	if err != nil {
		return fmt.Errorf("failed to get token from pool with address %s on chain %s: %w", c.PoolAddress, chain.Name(), err)
	}
	tokenConfig, err := state.TokenAdminRegistry.GetTokenConfig(nil, token)
	if err != nil {
		return fmt.Errorf("failed to get %s config from registry on chain %s: %w", tokenSymbol, chain.Name(), err)
	}

	// To update the pool, one of the following must be true
	//   - We are the current admin of the token
	//   - The admin of the token doesn't exist yet and we are the owner of the registry
	registryOwner, err := state.TokenAdminRegistry.Owner(nil)
	if err != nil {
		return fmt.Errorf("failed to get owner of registry on chain %s: %w", chain.Name(), err)
	}
	fromAddress := state.Timelock.Address() // again, "we" are either the Timelock or the deployer key
	if !useMcms {
		fromAddress = chain.DeployerKey.From
	}
	weCanBeAdmin := tokenConfig.Administrator == utils.ZeroAddress && fromAddress == registryOwner
	weAreAdmin := tokenConfig.Administrator == fromAddress
	if tokenConfig.TokenPool != c.PoolAddress && !(weCanBeAdmin || weAreAdmin) {
		return fmt.Errorf("address %s is unable to be the admin of %s on %s", fromAddress, tokenSymbol, chain.Name())
	}
	return nil
}

// ConfigureTokenAdminRegistryConfig is the configuration for the ConfigureTokenAdminRegistry changeset.
type ConfigureTokenAdminRegistryConfig struct {
	// MCMS defines the delay to use for Timelock (if absent, the changeset will attempt to use the deployer key).
	MCMS *MCMSConfig
	// RegistryUpdates defines the desired state of the registry on each given chain.
	RegistryUpdates map[uint64]RegistryConfig
	// TokenSymbol is the symbol of the token of interest.
	TokenSymbol TokenSymbol
}

func (c ConfigureTokenAdminRegistryConfig) Validate(env deployment.Environment) error {
	if c.TokenSymbol == "" {
		return errors.New("token symbol must be defined")
	}
	state, err := LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	for chainSelector, registryUpdate := range c.RegistryUpdates {
		err := deployment.IsValidChainSelector(chainSelector)
		if err != nil {
			return fmt.Errorf("failed to validate chain selector %d: %w", chainSelector, err)
		}
		chain, ok := env.Chains[chainSelector]
		if !ok {
			return fmt.Errorf("chain with selector %d does not exist in environment", chainSelector)
		}
		chainState, ok := state.Chains[chainSelector]
		if !ok {
			return fmt.Errorf("%s does not exist in state", chain.Name())
		}
		if tokenAdminRegistry := chainState.TokenAdminRegistry; tokenAdminRegistry == nil {
			return fmt.Errorf("missing tokenAdminRegistry on %s", chain.Name())
		}
		if timelock := chainState.Timelock; timelock == nil {
			return fmt.Errorf("missing timelock on %s", chain.Name())
		}
		if proposerMcm := chainState.ProposerMcm; proposerMcm == nil {
			return fmt.Errorf("missing proposerMcm on %s", chain.Name())
		}
		if err := registryUpdate.Validate(env.GetContext(), chain, chainState, c.MCMS != nil, c.TokenSymbol); err != nil {
			return fmt.Errorf("invalid pool update on %s: %w", chain.Name(), err)
		}
	}

	return nil
}

// ConfigureTokenAdminRegistry configures updates administrators and token pools on the TokenAdminRegistry.
func ConfigureTokenAdminRegistry(env deployment.Environment, c ConfigureTokenAdminRegistryConfig) (deployment.ChangesetOutput, error) {
	if err := c.Validate(env); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid ConfigureTokenAdminRegistryConfig: %w", err)
	}
	state, err := LoadOnchainState(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	var batches []timelock.BatchChainOperation
	timelocks := make(map[uint64]common.Address)
	proposers := make(map[uint64]*gethwrappers.ManyChainMultiSig)
	for chainSelector, registryUpdate := range c.RegistryUpdates {
		chain := env.Chains[chainSelector]
		chainState := state.Chains[chainSelector]

		tokenAdminRegistry := chainState.TokenAdminRegistry
		timelock := chainState.Timelock
		proposerMcm := chainState.ProposerMcm

		operations, err := createTokenAdminRegistryOps(chain, tokenAdminRegistry, registryUpdate, timelock, c.MCMS, c.TokenSymbol)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to make operations to configure token admin registry on %s: %w", chain.Name(), err)
		}
		if len(operations) > 0 {
			proposers[chainSelector] = proposerMcm
			timelocks[chainSelector] = timelock.Address()
			batches = append(batches, rbac_timelock.BatchChainOperation{
				Batch:           operations,
				ChainIdentifier: mcms.ChainIdentifier(chainSelector),
			})
		}
	}

	if len(batches) > 0 {
		proposal, err := proposalutils.BuildProposalFromBatches(
			timelocks,
			proposers,
			batches,
			fmt.Sprintf("configure %s on token admin registries", c.TokenSymbol),
			c.MCMS.MinDelay,
		)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return deployment.ChangesetOutput{
			Proposals:   []rbac_timelock.MCMSWithTimelockProposal{*proposal},
			AddressBook: nil,
			JobSpecs:    nil,
		}, nil
	}

	return deployment.ChangesetOutput{
		Proposals:   []timelock.MCMSWithTimelockProposal{},
		AddressBook: nil,
		JobSpecs:    nil,
	}, nil
}

// createTokenAdminRegistryOps creates all transactions required to configure the tokenAdminRegistry on a chain,
// either applying the transactions with the deployer key or returning an MCMS proposal.
func createTokenAdminRegistryOps(
	chain deployment.Chain,
	tokenAdminRegistry *token_admin_registry.TokenAdminRegistry,
	registryUpdate RegistryConfig,
	timelock *gethwrappers.RBACTimelock,
	mcmsConfig *MCMSConfig,
	tokenSymbol TokenSymbol,
) ([]mcms.Operation, error) {
	// Create opts and handler
	opts, handle, err := MakeTxOptsAndHandlerForContract(
		tokenAdminRegistry.Address(),
		chain,
		mcmsConfig,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to make tx opts and handler for registry on %s: %w", chain.Name(), err)
	}

	tokenPool, err := token_pool.NewTokenPool(registryUpdate.PoolAddress, chain.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to connect address %s on chain %s with token pool bindings: %w", registryUpdate.PoolAddress, chain.Name(), err)
	}
	token, err := tokenPool.GetToken(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get token from address %s on chain %s: %w", registryUpdate.PoolAddress, chain.Name(), err)
	}

	tokenConfig, err := tokenAdminRegistry.GetTokenConfig(nil, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get %s config from registry on chain %s: %w", tokenSymbol, chain.Name(), err)
	}

	fromAddress := timelock.Address()
	if opts.From == chain.DeployerKey.From {
		fromAddress = chain.DeployerKey.From
	}

	var operations []mcms.Operation

	if tokenConfig.TokenPool != registryUpdate.PoolAddress {
		if tokenConfig.Administrator != fromAddress {
			if tokenConfig.PendingAdministrator != fromAddress {
				// Propose administrator
				tx, err := tokenAdminRegistry.ProposeAdministrator(opts, token, fromAddress)
				if err != nil {
					return nil, fmt.Errorf("failed to create proposeAdministrator transaction for %s on %s registry: %w", tokenSymbol, chain.Name(), err)
				}
				mcmsOp, err := handle(tx)
				if err != nil {
					return nil, fmt.Errorf("failed to handle proposeAdministrator transaction for %s on %s registry: %w", tokenSymbol, chain.Name(), err)
				}
				if mcmsOp != nil {
					operations = append(operations, *mcmsOp)
				}
			}
			// Accept admin role
			tx, err := tokenAdminRegistry.AcceptAdminRole(opts, token)
			if err != nil {
				return nil, fmt.Errorf("failed to create acceptAdminRole transaction for %s on %s registry: %w", tokenSymbol, chain.Name(), err)
			}
			mcmsOp, err := handle(tx)
			if err != nil {
				return nil, fmt.Errorf("failed to handle acceptAdminRole transaction for %s on %s registry: %w", tokenSymbol, chain.Name(), err)
			}
			if mcmsOp != nil {
				operations = append(operations, *mcmsOp)
			}
		}
		// Set pool
		tx, err := tokenAdminRegistry.SetPool(opts, token, registryUpdate.PoolAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to create setPool transaction for %s on %s registry: %w", tokenSymbol, chain.Name(), err)
		}
		mcmsOp, err := handle(tx)
		if err != nil {
			return nil, fmt.Errorf("failed to handle setPool transaction for %s on %s registry: %w", tokenSymbol, chain.Name(), err)
		}
		if mcmsOp != nil {
			operations = append(operations, *mcmsOp)
		}
	}

	if registryUpdate.Administrator != fromAddress {
		// Transfer admin role
		tx, err := tokenAdminRegistry.TransferAdminRole(opts, token, registryUpdate.Administrator)
		if err != nil {
			return nil, fmt.Errorf("failed to create transferAdminRole transaction for %s on %s registry: %w", tokenSymbol, chain.Name(), err)
		}
		mcmsOp, err := handle(tx)
		if err != nil {
			return nil, fmt.Errorf("failed to handle transferAdminRole transaction for %s on %s registry: %w", tokenSymbol, chain.Name(), err)
		}
		if mcmsOp != nil {
			operations = append(operations, *mcmsOp)
		}
	}

	return operations, nil
}
