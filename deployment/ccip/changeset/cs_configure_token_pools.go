package changeset

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	rbac_timelock "github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_pool"
)

var _ deployment.ChangeSet[ConfigureTokenPoolContractsConfig] = ConfigureTokenPoolContracts

// RateLimiterConfig defines the inbound and outbound rate limits for a remote chain.
type RateLimiterConfig struct {
	// Inbound is the rate limiter config for inbound transfers from a remote chain.
	Inbound token_pool.RateLimiterConfig
	// Outbound is the rate limiter config for outbound transfers to a remote chain.
	Outbound token_pool.RateLimiterConfig
}

// validateRateLimterConfig validates rate and capacity in accordance with on-chain code.
// see https://github.com/smartcontractkit/ccip/blob/ccip-develop/contracts/src/v0.8/ccip/libraries/RateLimiter.sol.
func validateRateLimiterConfig(rateLimiterConfig token_pool.RateLimiterConfig) error {
	zero := big.NewInt(0)
	if rateLimiterConfig.IsEnabled {
		if rateLimiterConfig.Rate.Cmp(rateLimiterConfig.Capacity) >= 0 || rateLimiterConfig.Rate.Cmp(zero) == 0 {
			return errors.New("rate must be greater than 0 and less than capacity")
		}
	} else {
		if rateLimiterConfig.Rate.Cmp(zero) != 0 || rateLimiterConfig.Capacity.Cmp(zero) != 0 {
			return errors.New("rate and capacity must be 0")
		}
	}
	return nil
}

// RemoteChainsConfig defines rate limits for remote chains.
type RemoteChainsConfig map[uint64]RateLimiterConfig

func (c RemoteChainsConfig) Validate() error {
	for chainSelector, chainConfig := range c {
		if err := validateRateLimiterConfig(chainConfig.Inbound); err != nil {
			return fmt.Errorf("validation of inbound rate limiter config for remote chain with selector %d failed: %w", chainSelector, err)
		}
		if err := validateRateLimiterConfig(chainConfig.Outbound); err != nil {
			return fmt.Errorf("validation of outbound rate limiter config for remote chain with selector %d failed: %w", chainSelector, err)
		}
	}
	return nil
}

// TokenPoolConfig defines all the information required of the user to configure a token pool.
type TokenPoolConfig struct {
	// ChainUpdates defines the chains and corresponding rate limits that should be defined on the token pool.
	ChainUpdates RemoteChainsConfig
	// PoolAddress is the address of the token pool that we plan to make active on the registry.
	PoolAddress common.Address
}

func (c TokenPoolConfig) Validate(ctx context.Context, chain deployment.Chain, state CCIPChainState, useMcms bool, tokenSymbol TokenSymbol) error {
	// Ensure that required fields are defined
	if c.PoolAddress == utils.ZeroAddress {
		return errors.New("pool address must be defined")
	}

	// Ensure that the given pool address and symbol are aligned and known to the environment
	_, err := getTokenPoolWithSymbolAndAddress(state, chain, tokenSymbol, c.PoolAddress)
	if err != nil {
		return fmt.Errorf("failed to find token pool on %s with symbol %s and address %s: %w", chain.Name(), tokenSymbol, c.PoolAddress, err)
	}

	// Validate that the token pool is owned by the address that will be actioning the transactions (i.e. Timelock or deployer key)
	if err := commoncs.ValidateOwnership(ctx, useMcms, chain.DeployerKey.From, state.Timelock.Address(), state.TokenAdminRegistry); err != nil {
		return fmt.Errorf("token pool with address %s on %s failed ownership validation: %w", c.PoolAddress, chain.Name(), err)
	}

	// Validate chain configurations, namely rate limits
	if err := c.ChainUpdates.Validate(); err != nil {
		return fmt.Errorf("failed to validate chain updates: %w", err)
	}

	return nil
}

// ConfigureTokenPoolContractsConfig is the configuration for the ConfigureTokenPoolContractsConfig changeset.
type ConfigureTokenPoolContractsConfig struct {
	// MCMS defines the delay to use for Timelock (if absent, the changeset will attempt to use the deployer key).
	MCMS *MCMSConfig
	// PoolUpdates defines the changes that we want to make to the token pool on a chain
	PoolUpdates map[uint64]TokenPoolConfig
	// Symbol is the symbol of the token of interest.
	TokenSymbol TokenSymbol
}

func (c ConfigureTokenPoolContractsConfig) Validate(env deployment.Environment) error {
	if c.TokenSymbol == "" {
		return errors.New("token symbol must be defined")
	}
	state, err := LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	for chainSelector, poolUpdate := range c.PoolUpdates {
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
		for remoteChainSelector := range poolUpdate.ChainUpdates {
			remotePoolUpdate, ok := c.PoolUpdates[remoteChainSelector]
			if !ok {
				return fmt.Errorf("%s is expecting a pool update to be defined for chain with selector %d", chain.Name(), remoteChainSelector)
			}
			missingErr := fmt.Errorf("%s is expecting pool update on chain with selector %d to define a chain config pointing back to it", chain.Name(), remoteChainSelector)
			if remotePoolUpdate.ChainUpdates == nil {
				return missingErr
			}
			if _, ok := remotePoolUpdate.ChainUpdates[chainSelector]; !ok {
				return missingErr
			}
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
		if err := poolUpdate.Validate(env.GetContext(), chain, chainState, c.MCMS != nil, c.TokenSymbol); err != nil {
			return fmt.Errorf("invalid pool update on %s: %w", chain.Name(), err)
		}
	}

	return nil
}

// ConfigureTokenPoolContracts configures pools for a given token across multiple chains.
// The outputted MCMS proposal will update chain configurations on each pool, encompassing new chain additions and rate limit changes.
// Removing chain support is not in scope for this changeset.
func ConfigureTokenPoolContracts(env deployment.Environment, c ConfigureTokenPoolContractsConfig) (deployment.ChangesetOutput, error) {
	if err := c.Validate(env); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid ConfigureTokenPoolContractsConfig: %w", err)
	}
	state, err := LoadOnchainState(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	var batches []timelock.BatchChainOperation
	timelocks := make(map[uint64]common.Address)
	proposers := make(map[uint64]*gethwrappers.ManyChainMultiSig)
	for chainSelector := range c.PoolUpdates {
		chain := env.Chains[chainSelector]
		chainState := state.Chains[chainSelector]

		operations, err := configureTokenPool(env.Chains, state, c, chainSelector)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to make operations to configure token admin registry on %s: %w", chain.Name(), err)
		}
		if len(operations) > 0 {
			proposers[chainSelector] = chainState.ProposerMcm
			timelocks[chainSelector] = chainState.Timelock.Address()
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
			fmt.Sprintf("configure %s token pools", c.TokenSymbol),
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

// configureTokenPool creates all transactions required to configure the desired token pool on a chain,
// either applying the transactions with the deployer key or returning an MCMS proposal.
func configureTokenPool(
	chains map[uint64]deployment.Chain,
	state CCIPOnChainState,
	config ConfigureTokenPoolContractsConfig,
	chainSelector uint64,
) ([]mcms.Operation, error) {
	poolUpdate := config.PoolUpdates[chainSelector]
	poolAddress := poolUpdate.PoolAddress
	chain := chains[chainSelector]
	tokenPool, err := token_pool.NewTokenPool(poolAddress, chain.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to connect pool with address %s on %s with token pool bindings: %w", poolAddress, chain.Name(), err)
	}
	remoteTokenAddress, err := tokenPool.GetToken(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get token from pool with address %s on %s: %w", poolAddress, chain.Name(), err)
	}
	tokenAdminRegistry := state.Chains[chainSelector].TokenAdminRegistry
	tokenConfig, err := tokenAdminRegistry.GetTokenConfig(nil, remoteTokenAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get %s config from registry on %s: %w", config.TokenSymbol, chain.Name(), err)
	}

	// Create opts and handler
	opts, handle, err := makeTxOptsAndHandlerForContract(
		poolAddress,
		chain,
		chain.DeployerKey,
		config.MCMS,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to make transaction opts and handler for pool with address %s on %s: %w", poolAddress, chain.Name(), err)
	}

	// For adding chain support
	var chainAdditions []token_pool.TokenPoolChainUpdate
	// For updating rate limits
	var remoteChainSelectorsToUpdate []uint64
	var updatedOutboundConfigs []token_pool.RateLimiterConfig
	var updatedInboundConfigs []token_pool.RateLimiterConfig
	// For adding remote pools
	remotePoolAddressAdditions := make(map[uint64]common.Address)

	for remoteChainSelector, chainUpdate := range poolUpdate.ChainUpdates {
		isSupportedChain, err := tokenPool.IsSupportedChain(nil, remoteChainSelector)
		if err != nil {
			return nil, fmt.Errorf("failed to check if %d is supported on pool with address %s on %s: %w", remoteChainSelector, poolAddress, chain.Name(), err)
		}
		remoteChain := chains[remoteChainSelector]
		remotePoolAddress := config.PoolUpdates[remoteChainSelector].PoolAddress
		remoteTokenPool, err := token_pool.NewTokenPool(remotePoolAddress, remoteChain.Client)
		if err != nil {
			return nil, fmt.Errorf("failed to connect pool with address %s on %s with token pool bindings: %w", poolAddress, remoteChain.Name(), err)
		}
		remoteTokenAddress, err := remoteTokenPool.GetToken(nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get token from pool with address %s on %s: %w", remotePoolAddress, remoteChain.Name(), err)
		}
		remoteTokenAdminRegistry := state.Chains[remoteChainSelector].TokenAdminRegistry
		remoteTokenConfig, err := remoteTokenAdminRegistry.GetTokenConfig(nil, remoteTokenAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to get %s config from registry on %s: %w", config.TokenSymbol, remoteChain.Name(), err)
		}
		if isSupportedChain {
			// Just update the rate limits if the chain is already supported
			remoteChainSelectorsToUpdate = append(remoteChainSelectorsToUpdate, remoteChainSelector)
			updatedOutboundConfigs = append(updatedOutboundConfigs, chainUpdate.Outbound)
			updatedInboundConfigs = append(updatedInboundConfigs, chainUpdate.Inbound)
			// Also, add a new remote pool if the token pool on the remote chain is being updated
			if remoteTokenConfig.TokenPool != utils.ZeroAddress && remoteTokenConfig.TokenPool != remotePoolAddress {
				remotePoolAddressAdditions[remoteChainSelector] = remotePoolAddress
			}
		} else {
			// Add chain support if it doesn't yet exist
			// First, we need to assemble a list of valid remote pools
			// The desired token pool on the remote chain is added by default
			var remotePoolAddresses [][]byte
			remotePoolAddresses = append(remotePoolAddresses, remotePoolAddress.Bytes())
			// If the desired token pool is updating an old one, we still need to support the remote pool addresses that the old pool supported to ensure 0 downtime
			if tokenConfig.TokenPool != utils.ZeroAddress && tokenConfig.TokenPool != poolAddress {
				activeTokenPool, err := token_pool.NewTokenPool(tokenConfig.TokenPool, chain.Client)
				if err != nil {
					return nil, fmt.Errorf("failed to connect pool with address %s on %s with token pool bindings: %w", tokenConfig.TokenPool, chain.Name(), err)
				}
				remotePoolAddressesOnChain, err := activeTokenPool.GetRemotePools(nil, remoteChainSelector)
				if err != nil {
					return nil, fmt.Errorf("failed to fetch remote pools from token pool with address %s on chain %s: %w", tokenConfig.TokenPool, chain.Name(), err)
				}
				remotePoolAddresses = append(remotePoolAddresses, remotePoolAddressesOnChain...)
			}
			chainAdditions = append(chainAdditions, token_pool.TokenPoolChainUpdate{
				RemoteChainSelector:       remoteChainSelector,
				InboundRateLimiterConfig:  chainUpdate.Inbound,
				OutboundRateLimiterConfig: chainUpdate.Outbound,
				RemoteTokenAddress:        remoteTokenAddress.Bytes(),
				RemotePoolAddresses:       remotePoolAddresses,
			})
		}
	}

	var operations []mcms.Operation

	// Handle new chain support
	if len(chainAdditions) > 0 {
		tx, err := tokenPool.ApplyChainUpdates(opts, []uint64{}, chainAdditions)
		if err != nil {
			return nil, fmt.Errorf("failed to create applyChainUpdates transaction for token pool with address %s: %w", poolAddress, err)
		}
		mcmsOp, err := handle(tx)
		if err != nil {
			return nil, fmt.Errorf("failed to handle applyChainUpdates transaction for token pool with address %s: %w", poolAddress, err)
		}
		if mcmsOp != nil {
			operations = append(operations, *mcmsOp)
		}
	}

	// Handle updates to existing chain support
	if len(remoteChainSelectorsToUpdate) > 0 {
		tx, err := tokenPool.SetChainRateLimiterConfigs(opts, remoteChainSelectorsToUpdate, updatedOutboundConfigs, updatedInboundConfigs)
		if err != nil {
			return nil, fmt.Errorf("failed to create setChainRateLimiterConfigs transaction for token pool with address %s: %w", poolAddress, err)
		}
		mcmsOp, err := handle(tx)
		if err != nil {
			return nil, fmt.Errorf("failed to handle setChainRateLimiterConfigs transaction for token pool with address %s: %w", poolAddress, err)
		}
		if mcmsOp != nil {
			operations = append(operations, *mcmsOp)
		}
	}

	// Handle remote pool additions
	for remoteChainSelector, remotePoolAddress := range remotePoolAddressAdditions {
		tx, err := tokenPool.AddRemotePool(opts, remoteChainSelector, remotePoolAddress.Bytes())
		if err != nil {
			return nil, fmt.Errorf("failed to create addRemotePool transaction for token pool with address %s: %w", poolAddress, err)
		}
		mcmsOp, err := handle(tx)
		if err != nil {
			return nil, fmt.Errorf("failed to handle addRemotePool transaction for token pool with address %s: %w", poolAddress, err)
		}
		if mcmsOp != nil {
			operations = append(operations, *mcmsOp)
		}
	}

	return operations, nil
}
