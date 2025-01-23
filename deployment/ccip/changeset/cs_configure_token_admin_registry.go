package changeset

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_admin_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_pool"
)

var _ deployment.ChangeSet[ConfigureTokenAdminRegistryConfig] = ConfigureTokenAdminRegistryChangeset

// RegistryConfig defines a token and its state on the token admin registry
type RegistryConfig struct {
	// Type is the type of the token pool.
	Type deployment.ContractType
	// Version is the version of the token pool.
	Version semver.Version
	// ExternalAdministrator is the address of a 3rd party token administrator that should be set on the registry.
	ExternalAdministrator common.Address
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

	fromAddress := state.Timelock.Address() // "We" are either the Timelock or the deployer key
	if !useMcms {
		fromAddress = chain.DeployerKey.From
	}

	// Running this changeset has possible motivations: we want to update the pool for a token or transfer the admin rights of a token.
	// It doesn't really matter if we are doing one or both, so long as we are able to perform the action(s).

	// To perform these actions, we have to be the admin of the token. There are three ways this can happen:
	//   1. We are already the admin of the token (no action)
	//   2. We are the proposed admin of the token (just have to accept)
	//   3. We can become the admin of the token (have to propose and accept), which requires us to be the owner of the registry and for the token to be admin-less.
	// The following code checks these conditions.

	if tokenConfig.Administrator == fromAddress || tokenConfig.PendingAdministrator == fromAddress {
		// We are already the administrator / pending administrator & will be able to perform any actions required
		return nil
	}

	// If we are not admin / pending admin, we must set ourselves as admin of the token, which requires two things to be true.
	//   1. We own the token admin registry
	//   2. An admin musn't exist yet
	// We've already validated that we own the registry during ValidateOwnership, so we only need to check the 2nd condition
	if tokenConfig.Administrator != utils.ZeroAddress {
		return fmt.Errorf("unable to set %s as admin of %s token on %s: token already has an administrator (%s)", fromAddress, tokenSymbol, chain, tokenConfig.Administrator)
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

	chainConfigs, err := getConfigurationsByChain(env, state, c, deployerGroup)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to fetch configurations for each chain: %w", err)
	}

	// Propose admin pass
	for chainSelector := range c.RegistryUpdates {
		cc := chainConfigs[chainSelector]

		if cc.TokenConfigOnRegistry.Administrator != cc.Sender && cc.TokenConfigOnRegistry.PendingAdministrator != cc.Sender {
			_, err := cc.State.TokenAdminRegistry.ProposeAdministrator(cc.Opts, cc.TokenAddress, cc.Sender)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to create proposeAdministrator transaction for %s on %s registry: %w", c.TokenSymbol, cc.Chain, err)
			}
		}
	}

	proposeAdminOutput, err := deployerGroup.Enact(fmt.Sprintf("propose admin for %s on token admin registries", c.TokenSymbol))
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("propose admin for %s on token admin registries: %w", c.TokenSymbol, err)
	}

	// Accept admin pass
	for chainSelector := range c.RegistryUpdates {
		cc := chainConfigs[chainSelector]

		if cc.TokenConfigOnRegistry.Administrator != cc.Sender {
			_, err := cc.State.TokenAdminRegistry.AcceptAdminRole(cc.Opts, cc.TokenAddress)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to create acceptAdminRole transaction for %s on %s registry: %w", c.TokenSymbol, cc.Chain, err)
			}
		}
	}

	acceptAdminOutput, err := deployerGroup.Enact(fmt.Sprintf("accept admin rights for %s on token admin registries", c.TokenSymbol))
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to accept admin rights for %s on token admin registries: %w", c.TokenSymbol, err)
	}

	// Configuration pass (set pool, transfer admin role to 3rd party)
	for chainSelector, registryUpdate := range c.RegistryUpdates {
		cc := chainConfigs[chainSelector]

		// Only set the pool if we need to
		if cc.TokenConfigOnRegistry.TokenPool != cc.TokenPool.Address() {
			_, err := cc.State.TokenAdminRegistry.SetPool(cc.Opts, cc.TokenAddress, cc.TokenPool.Address())
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to create setPool transaction for %s on %s registry: %w", c.TokenSymbol, cc.Chain, err)
			}
		}

		// Only set the administrator to an external address if we need to
		if registryUpdate.ExternalAdministrator != cc.Sender {
			_, err := cc.State.TokenAdminRegistry.TransferAdminRole(cc.Opts, cc.TokenAddress, registryUpdate.ExternalAdministrator)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to create transferAdminRole transaction for %s on %s registry: %w", c.TokenSymbol, cc.Chain, err)
			}
		}
	}

	configurationOutput, err := deployerGroup.Enact(fmt.Sprintf("configure %s on token admin registries", c.TokenSymbol))
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to configure %s on token admin registries: %w", c.TokenSymbol, err)
	}

	if c.MCMS != nil {
		// Pre-allocate the proposal slice with the correct capacity
		totalProposals := len(proposeAdminOutput.Proposals) +
			len(acceptAdminOutput.Proposals) +
			len(configurationOutput.Proposals)
		proposals := make([]timelock.MCMSWithTimelockProposal, 0, totalProposals)
		proposals = append(proposals, proposeAdminOutput.Proposals...)
		proposals = append(proposals, acceptAdminOutput.Proposals...)
		proposals = append(proposals, configurationOutput.Proposals...)

		return deployment.ChangesetOutput{
			Proposals: proposals,
		}, nil
	}

	return deployment.ChangesetOutput{}, nil
}

// chainConfig defines the configuration needed to create operations on a chain
type chainConfig struct {
	TokenPool             *token_pool.TokenPool
	TokenAddress          common.Address
	TokenConfigOnRegistry token_admin_registry.TokenAdminRegistryTokenConfig
	Sender                common.Address
	State                 CCIPChainState
	Chain                 deployment.Chain
	Opts                  *bind.TransactOpts
}

// getConfigurationsByChain fetches the configuration required to create operations for each chain
func getConfigurationsByChain(
	env deployment.Environment,
	state CCIPOnChainState,
	c ConfigureTokenAdminRegistryConfig,
	deployerGroup *DeployerGroup,
) (map[uint64]chainConfig, error) {
	chainConfigs := make(map[uint64]chainConfig, len(c.RegistryUpdates))

	for chainSelector, registryUpdate := range c.RegistryUpdates {
		chain := env.Chains[chainSelector]
		chainState := state.Chains[chainSelector]

		opts, err := deployerGroup.GetDeployer(chainSelector)
		if err != nil {
			return map[uint64]chainConfig{}, fmt.Errorf("failed to get deployer for %s", chain)
		}

		tokenPool, err := getTokenPoolFromSymbolTypeAndVersion(chainState, chain, c.TokenSymbol, registryUpdate.Type, registryUpdate.Version)
		if err != nil {
			return map[uint64]chainConfig{}, fmt.Errorf("failed to find token pool on %s with symbol %s, type %s, and version %s: %w", chain.String(), c.TokenSymbol, registryUpdate.Type, registryUpdate.Version, err)
		}

		token, err := tokenPool.GetToken(nil)
		if err != nil {
			return map[uint64]chainConfig{}, fmt.Errorf("failed to get token from address %s on chain %s: %w", tokenPool.Address(), chain.String(), err)
		}

		tokenConfig, err := chainState.TokenAdminRegistry.GetTokenConfig(nil, token)
		if err != nil {
			return map[uint64]chainConfig{}, fmt.Errorf("failed to get %s config from registry on chain %s: %w", c.TokenSymbol, chain.String(), err)
		}

		sender := chainState.Timelock.Address()
		if c.MCMS == nil {
			sender = chain.DeployerKey.From
		}

		chainConfigs[chainSelector] = chainConfig{
			TokenPool:             tokenPool,
			TokenAddress:          token,
			TokenConfigOnRegistry: tokenConfig,
			Sender:                sender,
			State:                 chainState,
			Chain:                 chain,
			Opts:                  opts,
		}
	}

	return chainConfigs, nil
}
