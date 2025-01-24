package changeset

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/usdc_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/erc20"
	"github.com/smartcontractkit/chainlink/v2/evm/utils"
)

var _ deployment.ChangeSet[DeployUSDCTokenPoolContractsConfig] = DeployUSDCTokenPoolContractsChangeset

// DeployUSDCTokenPoolInput defines all information required of the user to deploy a new token pool contract.
type DeployUSDCTokenPoolInput struct {
	// TokenMessenger is the address of the USDC token messenger contract.
	TokenMessenger common.Address
	// USDCTokenAddress is the address of the token for which we are deploying a USDC token pool.
	TokenAddress common.Address
	// AllowList is the optional list of addresses permitted to initiate a token transfer.
	// If omitted, all addresses will be permitted to transfer the token.
	AllowList []common.Address
	// ForceDeployment forces deployment of a new token pool, even if one already exists for the corresponding token in state.
	ForceDeployment bool
}

func (i DeployUSDCTokenPoolInput) Validate(ctx context.Context, chain deployment.Chain, state CCIPChainState) error {
	// Ensure that required fields are populated
	if i.TokenAddress == utils.ZeroAddress {
		return errors.New("token address must be defined")
	}
	if i.TokenMessenger == utils.ZeroAddress {
		return errors.New("token messenger must be defined")
	}

	// Validate the token exists and matches the USDC symbol
	token, err := erc20.NewERC20(i.TokenAddress, chain.Client)
	if err != nil {
		return fmt.Errorf("failed to connect address %s with erc20 bindings: %w", i.TokenAddress, err)
	}
	symbol, err := token.Symbol(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("failed to fetch symbol from token with address %s: %w", i.TokenAddress, err)
	}
	if symbol != string(USDCSymbol) {
		return fmt.Errorf("symbol of token with address %s (%s) is not USDC", i.TokenAddress, symbol)
	}

	// Check if a USDC token pool with the given type already exists
	if _, ok := state.USDCTokenPools[currentTokenPoolVersion]; ok {
		return fmt.Errorf("USDC token pool with version %s already exists on %s", currentTokenPoolVersion, chain.String())
	}

	return nil
}

// DeployUSDCTokenPoolContractsConfig defines the USDC token pool contracts that need to be deployed on each chain.
type DeployUSDCTokenPoolContractsConfig struct {
	// NewUSDCPools defines the per-chain configuration of each new USDC pool.
	NewUSDCPools map[uint64]DeployUSDCTokenPoolInput
}

func (c DeployUSDCTokenPoolContractsConfig) Validate(env deployment.Environment) error {
	state, err := LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	for chainSelector, poolConfig := range c.NewUSDCPools {
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
			return fmt.Errorf("chain with selector %d does not exist in state", chainSelector)
		}
		if router := chainState.Router; router == nil {
			return fmt.Errorf("missing router on %s", chain.String())
		}
		if rmnProxy := chainState.RMNProxy; rmnProxy == nil {
			return fmt.Errorf("missing rmnProxy on %s", chain.String())
		}
		err = poolConfig.Validate(env.GetContext(), chain, chainState)
		if err != nil {
			return fmt.Errorf("failed to validate USDC token pool config for chain selector %d: %w", chainSelector, err)
		}
	}
	return nil
}

// DeployUSDCTokenPoolContractsChangeset deploys new USDC pools across multiple chains.
func DeployUSDCTokenPoolContractsChangeset(env deployment.Environment, c DeployUSDCTokenPoolContractsConfig) (deployment.ChangesetOutput, error) {
	if err := c.Validate(env); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid DeployUSDCTokenPoolContractsConfig: %w", err)
	}
	newAddresses := deployment.NewMemoryAddressBook()

	state, err := LoadOnchainState(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	for chainSelector, poolConfig := range c.NewUSDCPools {
		chain := env.Chains[chainSelector]
		chainState := state.Chains[chainSelector]

		_, err := deployment.DeployContract(env.Logger, chain, newAddresses,
			func(chain deployment.Chain) deployment.ContractDeploy[*usdc_token_pool.USDCTokenPool] {
				poolAddress, tx, usdcTokenPool, err := usdc_token_pool.DeployUSDCTokenPool(
					chain.DeployerKey, chain.Client, poolConfig.TokenMessenger, poolConfig.TokenAddress,
					poolConfig.AllowList, chainState.RMNProxy.Address(), chainState.Router.Address(),
				)
				return deployment.ContractDeploy[*usdc_token_pool.USDCTokenPool]{
					Address:  poolAddress,
					Contract: usdcTokenPool,
					Tv:       deployment.NewTypeAndVersion(USDCTokenPool, currentTokenPoolVersion),
					Tx:       tx,
					Err:      err,
				}
			},
		)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy USDC token pool on %s: %w", chain, err)
		}
	}

	return deployment.ChangesetOutput{
		AddressBook: newAddresses,
	}, nil
}
