package v1_5_1

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/token_pool"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

var (
	// SetRouterOnTokenPoolChangeset sets the router address on token pool contracts to the test router or the main router.
	SetRouterOnTokenPoolChangeset = cldf.CreateChangeSet(setRouterLogic, setRouterPrecondition)
)

// RouterUpdateOnPool is a struct that holds the information needed to update the router on a token pool.
type RouterUpdateOnPool struct {
	// TokenSymbol is symbol of the pool's underlying token.
	TokenSymbol changeset.TokenSymbol
	// PoolType is the type of the token pool contract.
	PoolType cldf.ContractType
	// Version is the version of the token pool contract.
	PoolVersion semver.Version
	// TestRouter indicates whether to set the test router or the main router.
	TestRouter bool
}

// SetRouterOnTokenPool is a struct that holds the information needed to update the router on token pools on multiple chains.
type SetRouterOnTokenPoolConfig struct {
	// PoolUpdates is a map of chain selector to a list of router updates on token pools.
	PoolUpdates map[uint64][]RouterUpdateOnPool
	// MCMS is the MCMS config (nil if not using MCMS).
	MCMS *proposalutils.TimelockConfig
}

func setRouterLogic(e deployment.Environment, c SetRouterOnTokenPoolConfig) (deployment.ChangesetOutput, error) {
	state, err := changeset.LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	deployerGroup := changeset.NewDeployerGroup(e, state, c.MCMS).WithDeploymentContext("propose admin role for tokens on token admin registries")

	for chainSel, updates := range c.PoolUpdates {
		chain := e.Chains[chainSel]
		chainState := state.Chains[chainSel]

		opts, err := deployerGroup.GetDeployer(chainSel)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to get deployer for %s", chain)
		}

		for _, update := range updates {
			poolAddress, ok := changeset.GetTokenPoolAddressFromSymbolTypeAndVersion(chainState, chain, update.TokenSymbol, update.PoolType, update.PoolVersion)
			if !ok {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to get token pool address for %s on %s: %w", update.TokenSymbol, chain, err)
			}
			pool, err := token_pool.NewTokenPool(poolAddress, chain.Client)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to create token pool bindings for %s on %s: %w", update.TokenSymbol, chain, err)
			}
			if update.TestRouter {
				_, err = pool.SetRouter(opts, chainState.TestRouter.Address())
				if err != nil {
					return deployment.ChangesetOutput{}, fmt.Errorf("failed to set test router on %s on %s: %w", update.TokenSymbol, chain, err)
				}
			} else {
				_, err = pool.SetRouter(opts, chainState.Router.Address())
				if err != nil {
					return deployment.ChangesetOutput{}, fmt.Errorf("failed to set main router on %s on %s: %w", update.TokenSymbol, chain, err)
				}
			}
		}
	}

	return deployerGroup.Enact()
}

func setRouterPrecondition(e deployment.Environment, c SetRouterOnTokenPoolConfig) error {
	state, err := changeset.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	for chainSel, updates := range c.PoolUpdates {
		err := changeset.ValidateChain(e, state, chainSel, c.MCMS)
		if err != nil {
			return fmt.Errorf("failed to validate chain %s: %w", chainSel, err)
		}
		chain := e.Chains[chainSel]

		chainState, ok := state.Chains[chainSel]
		if !ok {
			return fmt.Errorf("%s not found in onchain state", chain)
		}

		for _, update := range updates {
			if update.TestRouter {
				if chainState.TestRouter == nil {
					return fmt.Errorf("test router not found on %s", chain)
				}
			} else {
				if chainState.Router == nil {
					return fmt.Errorf("main router not found on %s", chain)
				}
			}
			poolAddress, ok := changeset.GetTokenPoolAddressFromSymbolTypeAndVersion(chainState, chain, update.TokenSymbol, update.PoolType, update.PoolVersion)
			if !ok {
				return fmt.Errorf("failed to get token pool address for %s on %s: %w", update.TokenSymbol, chain, err)
			}
			pool, err := token_pool.NewTokenPool(poolAddress, chain.Client)
			if err != nil {
				return fmt.Errorf("failed to create token pool bindings for %s on %s: %w", update.TokenSymbol, chain, err)
			}
			err = commoncs.ValidateOwnership(e.GetContext(), c.MCMS != nil, chain.DeployerKey.From, chainState.Timelock.Address(), pool)
			if err != nil {
				return fmt.Errorf("failed to validate ownership of %s on %s: %w", update.TokenSymbol, chain, err)
			}
		}
	}

	return nil
}
