package changeset

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"

	"github.com/smartcontractkit/chainlink/deployment"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_admin_registry"
)

var _ deployment.ChangeSet[ConfigureTokenAdminRegistryConfig] = ConfigureTokenAdminRegistryChangeset

// RegistryConfig defines a token and its state on the token admin registry
type RegistryConfig struct {
	// Type is the type of the token pool.
	Type deployment.ContractType
	// Version is the version of the token pool.
	Version semver.Version
	// Administrator is the address of the token administrator that should be set on the registry.
	Administrator common.Address
}

func (c RegistryConfig) Validate(ctx context.Context, chain deployment.Chain, state CCIPChainState, useMcms bool, tokenSymbol TokenSymbol) error {
	// Ensure that the inputted type is known
	if _, ok := tokenPoolTypes[c.Type]; !ok {
		return fmt.Errorf("%s is not a known token pool type", c.Type)
	}

	// Ensure that the inputted version is known
	if _, ok := tokenPoolVersions[c.Version]; !ok {
		return fmt.Errorf("%s is not a known token pool version", c.Version)
	}

	// Ensure that a pool with given symbol, type and version is known to the environment
	tokenPool, err := getTokenPoolFromSymbolTypeAndVersion(state, chain, tokenSymbol, c.Type, c.Version)
	if err != nil {
		return fmt.Errorf("failed to find token pool on %s with symbol %s, type %s, and version %s: %w", chain.String(), tokenSymbol, c.Type, c.Version, err)
	}

	// Validate that the token admin registry is owned by the address that will be actioning the transactions (i.e. Timelock or deployer key)
	if err := commoncs.ValidateOwnership(ctx, useMcms, chain.DeployerKey.From, state.Timelock.Address(), state.TokenAdminRegistry); err != nil {
		return fmt.Errorf("token admin registry failed ownership validation on %s: %w", chain.String(), err)
	}

	// Fetch information about the corresponding token and its state on the registry
	token, err := tokenPool.GetToken(nil)
	if err != nil {
		return fmt.Errorf("failed to get token from pool with address %s on chain %s: %w", tokenPool.Address(), chain.String(), err)
	}
	tokenConfig, err := state.TokenAdminRegistry.GetTokenConfig(nil, token)
	if err != nil {
		return fmt.Errorf("failed to get %s config from registry on chain %s: %w", tokenSymbol, chain.String(), err)
	}

	// To update the pool, one of the following must be true
	//   - We are the current admin of the token
	//   - The admin of the token doesn't exist yet and we are the owner of the registry
	registryOwner, err := state.TokenAdminRegistry.Owner(nil)
	if err != nil {
		return fmt.Errorf("failed to get owner of registry on chain %s: %w", chain.String(), err)
	}
	fromAddress := state.Timelock.Address() // again, "we" are either the Timelock or the deployer key
	if !useMcms {
		fromAddress = chain.DeployerKey.From
	}
	weCanBeAdmin := tokenConfig.Administrator == utils.ZeroAddress && fromAddress == registryOwner
	weAreAdmin := tokenConfig.Administrator == fromAddress
	if tokenConfig.TokenPool != tokenPool.Address() && !(weCanBeAdmin || weAreAdmin) {
		return fmt.Errorf("%s token pool with address %s is not set on the %s registry, but we can't set it because we do not control the admin address %s", tokenSymbol, tokenPool.Address(), chain, tokenConfig.Administrator)
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
			return fmt.Errorf("%s does not exist in state", chain.String())
		}
		if tokenAdminRegistry := chainState.TokenAdminRegistry; tokenAdminRegistry == nil {
			return fmt.Errorf("missing tokenAdminRegistry on %s", chain.String())
		}
		if c.MCMS != nil {
			if timelock := chainState.Timelock; timelock == nil {
				return fmt.Errorf("missing timelock on %s", chain.String())
			}
			if proposerMcm := chainState.ProposerMcm; proposerMcm == nil {
				return fmt.Errorf("missing proposerMcm on %s", chain.String())
			}
		}
		if err := registryUpdate.Validate(env.GetContext(), chain, chainState, c.MCMS != nil, c.TokenSymbol); err != nil {
			return fmt.Errorf("invalid pool update on %s: %w", chain.String(), err)
		}
	}

	return nil
}

// ConfigureTokenAdminRegistryChangeset configures updates administrators and token pools on the TokenAdminRegistry.
func ConfigureTokenAdminRegistryChangeset(env deployment.Environment, c ConfigureTokenAdminRegistryConfig) (deployment.ChangesetOutput, error) {
	if err := c.Validate(env); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid ConfigureTokenAdminRegistryConfig: %w", err)
	}
	state, err := LoadOnchainState(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	deployerGroup := NewDeployerGroup(env, state, c.MCMS)

	for chainSelector, registryUpdate := range c.RegistryUpdates {
		chain := env.Chains[chainSelector]
		chainState := state.Chains[chainSelector]

		tokenAdminRegistry := chainState.TokenAdminRegistry
		timelock := chainState.Timelock

		opts, err := deployerGroup.GetDeployer(chainSelector)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to get deployer for %s", chain)
		}

		err = createTokenAdminRegistryOps(opts, chainState, chain, tokenAdminRegistry, registryUpdate, timelock, c.MCMS, c.TokenSymbol)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to make operations to configure token admin registry on %s: %w", chain.String(), err)
		}
	}

	return deployerGroup.Enact(fmt.Sprintf("configure %s on token admin registries", c.TokenSymbol))
}

// createTokenAdminRegistryOps creates all transactions required to configure the tokenAdminRegistry on a chain,
// either applying the transactions with the deployer key or returning an MCMS proposal.
func createTokenAdminRegistryOps(
	opts *bind.TransactOpts,
	state CCIPChainState,
	chain deployment.Chain,
	tokenAdminRegistry *token_admin_registry.TokenAdminRegistry,
	registryUpdate RegistryConfig,
	timelock *gethwrappers.RBACTimelock,
	mcmsConfig *MCMSConfig,
	tokenSymbol TokenSymbol,
) error {
	tokenPool, err := getTokenPoolFromSymbolTypeAndVersion(state, chain, tokenSymbol, registryUpdate.Type, registryUpdate.Version)
	if err != nil {
		return fmt.Errorf("failed to find token pool on %s with symbol %s, type %s, and version %s: %w", chain.String(), tokenSymbol, registryUpdate.Type, registryUpdate.Version, err)
	}

	token, err := tokenPool.GetToken(nil)
	if err != nil {
		return fmt.Errorf("failed to get token from address %s on chain %s: %w", tokenPool.Address(), chain.String(), err)
	}

	tokenConfig, err := tokenAdminRegistry.GetTokenConfig(nil, token)
	if err != nil {
		return fmt.Errorf("failed to get %s config from registry on chain %s: %w", tokenSymbol, chain.String(), err)
	}

	fromAddress := timelock.Address()
	if mcmsConfig == nil {
		fromAddress = chain.DeployerKey.From
	}

	if tokenConfig.TokenPool != tokenPool.Address() {
		if tokenConfig.Administrator != fromAddress {
			if tokenConfig.PendingAdministrator != fromAddress {
				// Propose administrator
				_, err := tokenAdminRegistry.ProposeAdministrator(opts, token, fromAddress)
				if err != nil {
					return fmt.Errorf("failed to create proposeAdministrator transaction for %s on %s registry: %w", tokenSymbol, chain.String(), err)
				}
			}
			// Accept admin role
			_, err := tokenAdminRegistry.AcceptAdminRole(opts, token)
			if err != nil {
				return fmt.Errorf("failed to create acceptAdminRole transaction for %s on %s registry: %w", tokenSymbol, chain.String(), err)
			}
		}
		// Set pool
		_, err := tokenAdminRegistry.SetPool(opts, token, tokenPool.Address())
		if err != nil {
			return fmt.Errorf("failed to create setPool transaction for %s on %s registry: %w", tokenSymbol, chain.String(), err)
		}
	}

	if registryUpdate.Administrator != fromAddress {
		// Transfer admin role
		_, err := tokenAdminRegistry.TransferAdminRole(opts, token, registryUpdate.Administrator)
		if err != nil {
			return fmt.Errorf("failed to create transferAdminRole transaction for %s on %s registry: %w", tokenSymbol, chain.String(), err)
		}
	}

	return nil
}
