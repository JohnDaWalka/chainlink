package v1_6

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/siloed_lock_release_token_pool"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/deployergroup"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

var (
	_ cldf.ChangeSet[SiloedLockReleaseTokenPoolUpdateDesignationsChangesetConfig] = SiloedLockReleaseTokenPoolUpdateDesignations
	_ cldf.ChangeSet[SiloedLockReleaseTokenPoolSetRebalancerChangesetConfig]      = SiloedLockReleaseTokenPoolSetRebalancer
)

type SiloedLockReleaseTokenPoolUpdateDesignationsConfig struct {
	Removes []uint64
	Adds    []siloed_lock_release_token_pool.SiloedLockReleaseTokenPoolSiloConfigUpdate
}

type SiloedLockReleaseTokenPoolUpdateDesignationsChangesetConfig struct {
	Tokens map[uint64]map[shared.TokenSymbol]SiloedLockReleaseTokenPoolUpdateDesignationsConfig
	MCMS   *proposalutils.TimelockConfig
}

type SiloedLockReleaseTokenPoolSetRebalancerChangesetConfig struct {
	Tokens map[uint64]map[shared.TokenSymbol]common.Address
	MCMS   *proposalutils.TimelockConfig
}

func validateSiloedLockReleaseTokenPool(e cldf.Environment, state stateview.CCIPOnChainState, chainSelector uint64, mcmsCfg *proposalutils.TimelockConfig, token shared.TokenSymbol) error {
	if token == "" {
		return errors.New("token symbol cannot be empty")
	}

	if err := stateview.ValidateChain(e, state, chainSelector, mcmsCfg); err != nil {
		return fmt.Errorf("failed to validate chain with selector %d: %w", chainSelector, err)
	}

	chain, ok := e.BlockChains.EVMChains()[chainSelector]
	if !ok {
		return fmt.Errorf("chain with selector %d does not exist in environment", chainSelector)
	}

	chainState, ok := state.Chains[chainSelector]
	if !ok {
		return fmt.Errorf("%s does not exist in state", chain)
	}

	pool, ok := chainState.SiloedLockReleaseTokenPool[token][deployment.Version1_6_0]
	if !ok {
		return fmt.Errorf("siloed lock release token pool for token %s does not exist", token)
	}

	tokenAddress, err := pool.GetToken(&bind.CallOpts{Context: e.GetContext()})
	if err != nil {
		return fmt.Errorf("failed to get token address from pool with address %s: %w", pool.Address(), err)
	}

	return validateTokenSymbol(e.GetContext(), chain, tokenAddress, token)
}

func (c SiloedLockReleaseTokenPoolUpdateDesignationsChangesetConfig) Validate(e cldf.Environment) error {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	for chainSelector, tokens := range c.Tokens {
		for token, designationConfig := range tokens {
			if err := validateSiloedLockReleaseTokenPool(e, state, chainSelector, c.MCMS, token); err != nil {
				return err
			}

			chainState, _ := state.EVMChainState(chainSelector)
			pool := chainState.SiloedLockReleaseTokenPool[token][deployment.Version1_6_0]

			for _, remoteChainSelector := range designationConfig.Removes {
				remoteChain, ok := e.BlockChains.EVMChains()[remoteChainSelector]
				if !ok {
					return fmt.Errorf("chain with selector %d does not exist in environment", remoteChainSelector)
				}

				isSiloed, err := pool.IsSiloed(&bind.CallOpts{Context: e.GetContext()}, remoteChainSelector)
				if err != nil {
					return fmt.Errorf("failed to check if chain %s is soled: %w", remoteChain, err)
				}

				if !isSiloed {
					return fmt.Errorf("chain %s is not siloed", remoteChain)
				}
			}

			for _, add := range designationConfig.Adds {
				remoteChain, ok := e.BlockChains.EVMChains()[add.RemoteChainSelector]
				if !ok {
					return fmt.Errorf("chain with selector %d does not exist in environment", add.RemoteChainSelector)
				}

				isSiloed, err := pool.IsSiloed(&bind.CallOpts{Context: e.GetContext()}, add.RemoteChainSelector)
				if err != nil {
					return fmt.Errorf("failed to check if chain %s is soled: %w", remoteChain, err)
				}

				if isSiloed {
					return fmt.Errorf("chain %s is siloed", remoteChain)
				}

				isSupportedChain, err := pool.IsSupportedChain(&bind.CallOpts{Context: e.GetContext()}, add.RemoteChainSelector)
				if err != nil {
					return fmt.Errorf("failed to check if chain %s is supported: %w", remoteChain, err)
				}

				if !isSupportedChain {
					return fmt.Errorf("chain %s is not supported", remoteChain)
				}

				if add.Rebalancer == (common.Address{}) {
					return errors.New("rebalancer address cannot be empty")
				}
			}
		}
	}

	return nil
}

func (c SiloedLockReleaseTokenPoolSetRebalancerChangesetConfig) Validate(e cldf.Environment) error {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	for chainSelector, tokens := range c.Tokens {
		for token, newRebalancer := range tokens {
			if err := validateSiloedLockReleaseTokenPool(e, state, chainSelector, c.MCMS, token); err != nil {
				return err
			}

			if newRebalancer == (common.Address{}) {
				return errors.New("rebalancer address cannot be empty")
			}

			chainState, _ := state.EVMChainState(chainSelector)
			pool := chainState.SiloedLockReleaseTokenPool[token][deployment.Version1_6_0]

			rebalancer, err := pool.GetRebalancer(&bind.CallOpts{Context: e.GetContext()})
			if err != nil {
				return fmt.Errorf("failed to get rebalancer address from pool with address %s: %w", pool.Address(), err)
			}

			if rebalancer == newRebalancer {
				return errors.New("rebalancer address is unchanged")
			}
		}
	}

	return nil
}

// SiloedLockReleaseTokenPoolUpdateDesignations updates designations for chains on whether to mark funds as Siloed or not.
func SiloedLockReleaseTokenPoolUpdateDesignations(e cldf.Environment, c SiloedLockReleaseTokenPoolUpdateDesignationsChangesetConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid SiloedLockReleaseTokenPoolUpdateDesignationsChangesetConfig: %w", err)
	}

	state, _ := stateview.LoadOnchainState(e)
	deployerGroup := deployergroup.NewDeployerGroup(e, state, c.MCMS).WithDeploymentContext("update designations for siloed lock release token pool")

	for chainSelector, tokens := range c.Tokens {
		opts, err := deployerGroup.GetDeployer(chainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get deployer for chain with selector %d: %w", chainSelector, err)
		}

		chainState, _ := state.EVMChainState(chainSelector)
		for token, designationConfig := range tokens {
			pool := chainState.SiloedLockReleaseTokenPool[token][deployment.Version1_6_0]

			if _, err := pool.UpdateSiloDesignations(opts, designationConfig.Removes, designationConfig.Adds); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to update designations for token %s on chain with selector %d: %w", token, chainSelector, err)
			}
		}
	}

	return deployerGroup.Enact()
}

// SiloedLockReleaseTokenPoolSetRebalancer sets the rebalancer for a token pool.
func SiloedLockReleaseTokenPoolSetRebalancer(e cldf.Environment, c SiloedLockReleaseTokenPoolSetRebalancerChangesetConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid SiloedLockReleaseTokenPoolSetRebalancerChangesetConfig: %w", err)
	}

	state, _ := stateview.LoadOnchainState(e)
	deployerGroup := deployergroup.NewDeployerGroup(e, state, c.MCMS).WithDeploymentContext("set rebalancer for siloed lock release token pool")

	for chainSelector, tokens := range c.Tokens {
		opts, err := deployerGroup.GetDeployer(chainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get deployer for chain with selector %d: %w", chainSelector, err)
		}

		chainState, _ := state.EVMChainState(chainSelector)
		for token, newRebalancer := range tokens {
			pool := chainState.SiloedLockReleaseTokenPool[token][deployment.Version1_6_0]

			if _, err := pool.SetRebalancer(opts, newRebalancer); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to set rebalancer for token %s on chain with selector %d: %w", token, chainSelector, err)
			}
		}
	}

	return deployerGroup.Enact()
}
