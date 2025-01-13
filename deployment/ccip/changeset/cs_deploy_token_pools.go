package changeset

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/burn_from_mint_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/burn_with_from_mint_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_admin_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_pool"
)

/*
TODO: Other actions required to support token pools ->

- Removing support for chains (not as common, likely a separate changeset)
- Set rate limit admin (usually left as zero address, not urgent)
- Set rebalancer on lock-release pools (usually left as zero address, not urgent)
- Deploying a price feed for a token with BPS & reconfiguring a DON to support the feed

Deploying price feeds for tokens with BPS [REQUIRED]

Others, not as urgent

*/

var _ deployment.ChangeSet[DeployTokenPoolContractsConfig] = DeployTokenPoolContracts

// zero returns a zero-value big.Int.
func zero() *big.Int {
	return big.NewInt(0)
}

// zeroAddress returns a zero-value Ethereum address.
func zeroAddress() common.Address {
	return common.BigToAddress(zero())
}

// TokenChainConfig defines all information required to construct operations that fully configure the pool.
type TokenChainConfig struct {
	TokenAdminRegistry *token_admin_registry.TokenAdminRegistry
	TokenPool          *token_pool.TokenPool
	TokenAddress       common.Address
	TimelockAddress    common.Address
	ExternalAdmin      common.Address
	RegistryState      token_admin_registry.TokenAdminRegistryTokenConfig
	RemoteChainsToAdd  RemoteChains
	OwnedByTimelock    bool
}

// RateLimiterConfig defines the inbound and outbound rate limits for a remote chain.
type RateLimiterConfig struct {
	Inbound  token_pool.RateLimiterConfig
	Outbound token_pool.RateLimiterConfig
}

// validateRateLimterConfig validates rate and capacity in accordance with on-chain code.
// see https://github.com/smartcontractkit/ccip/blob/ccip-develop/contracts/src/v0.8/ccip/libraries/RateLimiter.sol.
func validateRateLimiterConfig(rateLimiterConfig token_pool.RateLimiterConfig) error {
	if rateLimiterConfig.IsEnabled {
		if rateLimiterConfig.Rate.Cmp(rateLimiterConfig.Capacity) >= 0 || rateLimiterConfig.Rate.Cmp(zero()) == 0 {
			return errors.New("rate must be greater than 0 and less than capacity")
		}
	} else {
		if rateLimiterConfig.Rate.Cmp(zero()) != 0 || rateLimiterConfig.Capacity.Cmp(zero()) != 0 {
			return errors.New("rate and capacity must be 0")
		}
	}
	return nil
}

// RemoteChainsConfig defines rate limits for remote chains.
type RemoteChains map[uint64]RateLimiterConfig

func (rc RemoteChains) Validate() error {
	for chainSelector, chainConfig := range rc {
		if err := validateRateLimiterConfig(chainConfig.Inbound); err != nil {
			return fmt.Errorf("validation of inbound rate limiter config for remote chain with selector %d failed: %w", chainSelector, err)
		}
		if err := validateRateLimiterConfig(chainConfig.Outbound); err != nil {
			return fmt.Errorf("validation of outbound rate limiter config for remote chain with selector %d failed: %w", chainSelector, err)
		}
	}
	return nil
}

// BaseTokenPoolInput defines all the information required of the user to configure new and existing pools.
type BaseTokenPoolInput struct {
	RemoteChainsToAdd RemoteChains
	TokenAddress      common.Address
	ExternalAdmin     common.Address
}

func (t BaseTokenPoolInput) Validate() error {
	if err := t.RemoteChainsToAdd.Validate(); err != nil {
		return fmt.Errorf("failed to validate remote chains config: %w", err)
	}
	return nil
}

// NewTokenPoolInput defines all information required of the user to deploy and configure a new token pool.
type NewTokenPoolInput struct {
	BaseTokenPoolInput
	Type               deployment.ContractType
	LocalTokenDecimals uint8
	AllowList          []common.Address
	AcceptLiquidity    *bool
}

func (t NewTokenPoolInput) Validate() error {
	if err := t.RemoteChainsToAdd.Validate(); err != nil {
		return fmt.Errorf("failed to validate remote chains config: %w", err)
	}
	if t.Type != BurnMintTokenPool && t.Type != BurnFromMintTokenPool && t.Type != BurnWithFromMintTokenPool && t.Type != LockReleaseTokenPool {
		return fmt.Errorf("%s is not a valid token pool type", t.Type)
	}
	if t.Type == LockReleaseTokenPool && t.AcceptLiquidity == nil {
		return errors.New("accept liquidity must be defined for lock release pools")
	}
	if t.Type != LockReleaseTokenPool && t.AcceptLiquidity != nil {
		return errors.New("accept liquidity must be nil for burn mint pools")
	}
	return nil
}

// DeployTokenPoolContractsConfig is the configuration for the DeployTokenPoolContracts changeset.
type DeployTokenPoolContractsConfig struct {
	Symbol              TokenSymbol
	TimelockDelay       time.Duration
	ExistingPoolUpdates map[uint64]BaseTokenPoolInput
	NewPools            map[uint64]NewTokenPoolInput
}

func (c DeployTokenPoolContractsConfig) Validate() error {
	seenChains := make(map[uint64]struct{})
	for chainSelector, chainConfig := range c.ExistingPoolUpdates {
		seenChains[chainSelector] = struct{}{}
		if err := chainConfig.Validate(); err != nil {
			return fmt.Errorf("chain with selector %d is invalid: %w", chainSelector, err)
		}
	}
	for chainSelector, chainConfig := range c.NewPools {
		if _, ok := seenChains[chainSelector]; ok {
			return fmt.Errorf("chain overlap exists between new pools and updates to existing pools")
		}
		if err := chainConfig.Validate(); err != nil {
			return fmt.Errorf("chain with selector %d is invalid: %w", chainSelector, err)
		}
	}

	return nil
}

// DeployTokenPoolContract deploys & configures new pools for a given token across multiple chains.
// The changeset will first deploy new token pools and transfer ownership of the pools to the Timelock.
// The outputted MCMS proposal will apply chain updates on each token pool and set new pools on the TokenAdminRegistry.
func DeployTokenPoolContracts(env deployment.Environment, c DeployTokenPoolContractsConfig) (deployment.ChangesetOutput, error) {
	if err := c.Validate(); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid DeployTokenPoolContractsConfig: %w", err)
	}
	newAddresses := deployment.NewMemoryAddressBook()

	state, err := LoadOnchainState(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	tokenChainConfigs := make(map[uint64]TokenChainConfig)

	// Deploy new token pools & transfer ownership of each new pool to the timelock
	for chainSelector, chainConfig := range c.NewPools {
		chainEnv, ok := env.Chains[chainSelector]
		if !ok {
			return deployment.ChangesetOutput{}, fmt.Errorf("no chain with selector %d found in environment", chainSelector)
		}
		chainState := state.Chains[chainSelector] // state is derived from env, no need to re-check
		tokenChainConfig, err := deployAndTransferTokenPoolToTimelock(env.Logger, chainEnv, chainState, newAddresses, chainConfig)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy and transfer token pool to timelock on %s: %w", chainEnv.Name(), err)
		}
		tokenChainConfigs[chainSelector] = tokenChainConfig
	}

	// Fetch addresses of the existing token pools via the token admin registry
	for chainSelector, chainConfig := range c.ExistingPoolUpdates {
		chainEnv, ok := env.Chains[chainSelector]
		if !ok {
			return deployment.ChangesetOutput{}, fmt.Errorf("no chain with selector %d found in environment", chainSelector)
		}
		chainState := state.Chains[chainSelector] // state is derived from env, no need to re-check
		tokenChainConfig, err := fetchAndValidateTimelockOwnedTokenPool(chainEnv, chainState, chainConfig)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to fetch token pool on %s for token with address %s: %w", chainEnv.Name(), chainConfig.TokenAddress, err)
		}
		tokenChainConfigs[chainSelector] = tokenChainConfig
	}

	var operations []timelock.BatchChainOperation
	timelocks := make(map[uint64]common.Address)
	proposers := make(map[uint64]*gethwrappers.ManyChainMultiSig)

	for chainSelector := range tokenChainConfigs {
		chainEnv := env.Chains[chainSelector] // chain selector has already been confirmed as a valid key
		batch, err := makeTokenPoolOperationsForChain(chainSelector, tokenChainConfigs)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to create operations for token pool on %s: %w", chainEnv.Name(), err)
		}
		proposers[chainSelector] = state.Chains[chainSelector].ProposerMcm
		timelocks[chainSelector] = state.Chains[chainSelector].Timelock.Address()

		operations = append(operations, timelock.BatchChainOperation{
			Batch:           batch,
			ChainIdentifier: mcms.ChainIdentifier(chainSelector),
		})
	}

	proposal, err := proposalutils.BuildProposalFromBatches(
		timelocks,
		proposers,
		operations,
		fmt.Sprintf("update token pool deployments for %s", c.Symbol),
		c.TimelockDelay,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
	}

	return deployment.ChangesetOutput{
		Proposals:   []timelock.MCMSWithTimelockProposal{*proposal},
		AddressBook: newAddresses,
		JobSpecs:    nil,
	}, nil
}

// fetchAndValidateTimelockOwnedTokenPool constructs a token's chain configuration based on the current pool set on the registry,
// asserting that the pool is current owned by the Timelock.
func fetchAndValidateTimelockOwnedTokenPool(
	chainEnv deployment.Chain,
	chainState CCIPChainState,
	chainConfig BaseTokenPoolInput,
) (TokenChainConfig, error) {
	// Assumes that timelock & tokenAdminRegistry are not nil
	// TODO: Should we assert this?
	tokenAdminRegistry := chainState.TokenAdminRegistry
	timelock := chainState.Timelock

	tokenConfigOnRegistry, err := tokenAdminRegistry.GetTokenConfig(nil, chainConfig.TokenAddress)
	if err != nil {
		return TokenChainConfig{}, fmt.Errorf("failed to get config on %s registry for token with address %s: %w", chainEnv.Name(), chainConfig.TokenAddress, err)
	}
	if tokenConfigOnRegistry.TokenPool.Cmp(zeroAddress()) == 0 {
		return TokenChainConfig{}, fmt.Errorf("token with address %s is not set on %s registry: %w", chainConfig.TokenAddress, chainEnv.Name(), err)
	}
	tokenPool, err := token_pool.NewTokenPool(tokenConfigOnRegistry.TokenPool, chainEnv.Client)
	if err != nil {
		return TokenChainConfig{}, fmt.Errorf("failed to connect address %s on %s with token pool bindings: %w", tokenConfigOnRegistry.TokenPool, chainEnv.Name(), err)
	}
	owner, err := tokenPool.Owner(nil)
	if err != nil {
		return TokenChainConfig{}, fmt.Errorf("failed to fetch owner from token pool with address %s on %s: %w", tokenConfigOnRegistry.TokenPool, chainEnv.Name(), err)
	}
	if owner.Cmp(timelock.Address()) != 0 {
		return TokenChainConfig{}, fmt.Errorf("token pool with address %s on %s is not owned by the timelock", tokenConfigOnRegistry.TokenPool, chainEnv.Name())
	}
	return TokenChainConfig{
		TokenAdminRegistry: tokenAdminRegistry,
		TokenAddress:       chainConfig.TokenAddress,
		TokenPool:          tokenPool,
		RegistryState:      tokenConfigOnRegistry,
		OwnedByTimelock:    true,
		RemoteChainsToAdd:  chainConfig.RemoteChainsToAdd,
		TimelockAddress:    timelock.Address(),
		ExternalAdmin:      chainConfig.ExternalAdmin,
	}, nil
}

// deployAndTransferTokenPoolToTimelock deploys a token pool and transfers ownership to the timelock using CCIP chain state.
func deployAndTransferTokenPoolToTimelock(
	logger logger.Logger,
	chainEnv deployment.Chain,
	chainState CCIPChainState,
	addressBook deployment.AddressBook,
	chainConfig NewTokenPoolInput,
) (TokenChainConfig, error) {
	// Assumes that router, rmnProxy, tokenAdminRegistry, & timelock will not be nil
	// TODO: Should we assert this?
	router := chainState.Router
	rmnProxy := chainState.RMNProxy
	timelock := chainState.Timelock
	tokenAdminRegistry := chainState.TokenAdminRegistry

	tokenConfigOnRegistry, err := tokenAdminRegistry.GetTokenConfig(nil, chainConfig.TokenAddress)
	if err != nil {
		return TokenChainConfig{}, fmt.Errorf("failed to get config on %s registry for token with address %s: %w", chainEnv.Name(), chainConfig.TokenAddress, err)
	}
	tokenPoolDeployment, err := deployTokenPool(
		logger,
		chainEnv,
		addressBook,
		chainConfig,
		router.Address(),
		rmnProxy.Address(),
	)
	if err != nil {
		return TokenChainConfig{}, fmt.Errorf("failed to deploy token pool on %s: %w", chainEnv.Name(), err)
	}
	tx, err := tokenPoolDeployment.Contract.TransferOwnership(chainEnv.DeployerKey, timelock.Address())
	if err != nil {
		return TokenChainConfig{}, fmt.Errorf("failed to transfer ownership of token pool to timelock on %s: %w", chainEnv.Name(), err)
	}
	_, err = chainEnv.Confirm(tx)
	if err != nil {
		return TokenChainConfig{}, fmt.Errorf("failed to confirm ownership transfer of token pool to timelock on %s: %w", chainEnv.Name(), err)
	}

	return TokenChainConfig{
		TokenAdminRegistry: tokenAdminRegistry,
		TokenAddress:       chainConfig.TokenAddress,
		TokenPool:          tokenPoolDeployment.Contract,
		RegistryState:      tokenConfigOnRegistry,
		OwnedByTimelock:    false,
		RemoteChainsToAdd:  chainConfig.RemoteChainsToAdd,
		TimelockAddress:    timelock.Address(),
		ExternalAdmin:      chainConfig.ExternalAdmin,
	}, nil
}

// deployTokenPool deploys a token pool contract based on a given type & configuration.
func deployTokenPool(
	logger logger.Logger,
	chainEnv deployment.Chain,
	addressBook deployment.AddressBook,
	chainConfig NewTokenPoolInput,
	routerAddress common.Address,
	rmnProxyAddress common.Address,
) (*deployment.ContractDeploy[*token_pool.TokenPool], error) {
	return deployment.DeployContract(logger, chainEnv, addressBook,
		func(chain deployment.Chain) deployment.ContractDeploy[*token_pool.TokenPool] {
			var tpAddr common.Address
			var tx *types.Transaction
			var err error
			switch chainConfig.Type {
			case BurnMintTokenPool:
				tpAddr, tx, _, err = burn_mint_token_pool.DeployBurnMintTokenPool(
					chain.DeployerKey,
					chain.Client,
					chainConfig.TokenAddress,
					chainConfig.LocalTokenDecimals,
					chainConfig.AllowList,
					rmnProxyAddress,
					routerAddress,
				)
			case BurnWithFromMintTokenPool:
				tpAddr, tx, _, err = burn_with_from_mint_token_pool.DeployBurnWithFromMintTokenPool(
					chain.DeployerKey,
					chain.Client,
					chainConfig.TokenAddress,
					chainConfig.LocalTokenDecimals,
					chainConfig.AllowList,
					rmnProxyAddress,
					routerAddress,
				)
			case BurnFromMintTokenPool:
				tpAddr, tx, _, err = burn_from_mint_token_pool.DeployBurnFromMintTokenPool(
					chain.DeployerKey,
					chain.Client,
					chainConfig.TokenAddress,
					chainConfig.LocalTokenDecimals,
					chainConfig.AllowList,
					rmnProxyAddress,
					routerAddress,
				)
			case LockReleaseTokenPool:
				tpAddr, tx, _, err = lock_release_token_pool.DeployLockReleaseTokenPool(
					chain.DeployerKey,
					chain.Client,
					chainConfig.TokenAddress,
					chainConfig.LocalTokenDecimals,
					chainConfig.AllowList,
					rmnProxyAddress,
					*chainConfig.AcceptLiquidity,
					routerAddress,
				)
			}
			var tp *token_pool.TokenPool
			if err == nil { // prevents overwriting the error (also, if there were an error with deployment, converting to an abstract token pool wouldn't be useful)
				tp, err = token_pool.NewTokenPool(tpAddr, chain.Client)
			}
			return deployment.ContractDeploy[*token_pool.TokenPool]{
				Address: tpAddr, Contract: tp, Tv: deployment.NewTypeAndVersion(chainConfig.Type, deployment.Version1_5_1), Tx: tx, Err: err,
			}
		},
	)
}

// makeTokenPoolOperationsForChain constructs a batch of MCMS operations to configure a token pool on a chain.
func makeTokenPoolOperationsForChain(
	chainSelector uint64,
	tokenChainConfigs map[uint64]TokenChainConfig,
) ([]mcms.Operation, error) {
	var batch []mcms.Operation
	tokenChainConfig, ok := tokenChainConfigs[chainSelector]
	if !ok {
		return []mcms.Operation{}, fmt.Errorf("no token found on chain with selector %d", chainSelector)
	}

	// Accept ownership if the timelock does not currently own the pool
	if !tokenChainConfig.OwnedByTimelock {
		acceptOwnershipTx, err := tokenChainConfig.TokenPool.AcceptOwnership(deployment.SimTransactOpts())
		if err != nil {
			return []mcms.Operation{}, fmt.Errorf("failed to create acceptOwnership transaction: %w", err)
		}
		batch = append(batch, mcms.Operation{
			To:    tokenChainConfig.TokenPool.Address(),
			Data:  acceptOwnershipTx.Data(),
			Value: zero(),
		})
	}

	// Apply chain updates on the token pool
	var chainUpdates []token_pool.TokenPoolChainUpdate
	for remoteChainSelector, remoteChainConfig := range tokenChainConfig.RemoteChainsToAdd {
		remoteTokenConfig, ok := tokenChainConfigs[remoteChainSelector]
		if !ok {
			return []mcms.Operation{}, fmt.Errorf("no token found on remote chain with selector %d", remoteChainSelector)
		}
		remotePoolAddresses := [][]byte{remoteTokenConfig.TokenPool.Address().Bytes()}
		// If the token pool on the remote chain's registry is not the current pool nor the zero address, we should add it as a supported remote pool to avoid downtime
		if remoteTokenConfig.RegistryState.TokenPool.Cmp(remoteTokenConfig.TokenPool.Address()) != 0 && remoteTokenConfig.RegistryState.TokenPool.Cmp(zeroAddress()) != 0 {
			remotePoolAddresses = append(remotePoolAddresses, remoteTokenConfig.RegistryState.TokenPool.Bytes())
		}
		chainUpdates = append(chainUpdates, token_pool.TokenPoolChainUpdate{
			RemoteChainSelector:       remoteChainSelector,
			InboundRateLimiterConfig:  remoteChainConfig.Inbound,
			OutboundRateLimiterConfig: remoteChainConfig.Outbound,
			RemoteTokenAddress:        remoteTokenConfig.TokenAddress.Bytes(),
			RemotePoolAddresses:       remotePoolAddresses,
		})
	}
	if len(chainUpdates) > 0 {
		applyChainUpdatesTx, err := tokenChainConfig.TokenPool.ApplyChainUpdates(deployment.SimTransactOpts(), []uint64{}, chainUpdates)
		if err != nil {
			return []mcms.Operation{}, fmt.Errorf("failed to create applyChainUpdates transaction: %w", err)
		}
		batch = append(batch, mcms.Operation{
			To:    tokenChainConfig.TokenPool.Address(),
			Data:  applyChainUpdatesTx.Data(),
			Value: zero(),
		})
	}

	// Set the administrator of the token to the timelock (if it hasn't been set before)
	noExistingAdmin := tokenChainConfig.RegistryState.Administrator.Cmp(zeroAddress()) == 0
	if noExistingAdmin {
		proposeAdministratorTx, err := tokenChainConfig.TokenAdminRegistry.ProposeAdministrator(deployment.SimTransactOpts(), tokenChainConfig.TokenAddress, tokenChainConfig.TimelockAddress)
		if err != nil {
			return []mcms.Operation{}, fmt.Errorf("failed to create proposeAdministrator transaction: %w", err)
		}
		batch = append(batch, mcms.Operation{
			To:    tokenChainConfig.TokenAdminRegistry.Address(),
			Data:  proposeAdministratorTx.Data(),
			Value: zero(),
		})
		acceptAdminRoleTx, err := tokenChainConfig.TokenAdminRegistry.AcceptAdminRole(deployment.SimTransactOpts(), tokenChainConfig.TokenAddress)
		if err != nil {
			return []mcms.Operation{}, fmt.Errorf("failed to create acceptAdminRole transaction: %w", err)
		}
		batch = append(batch, mcms.Operation{
			To:    tokenChainConfig.TokenAdminRegistry.Address(),
			Data:  acceptAdminRoleTx.Data(),
			Value: zero(),
		})
	}
	isTimelockAdmin := noExistingAdmin || tokenChainConfig.RegistryState.Administrator.Cmp(tokenChainConfig.TimelockAddress) == 0

	// Set the pool if the timelock is admin at this point & pool isn't already set
	if isTimelockAdmin && tokenChainConfig.RegistryState.TokenPool.Cmp(tokenChainConfig.TokenPool.Address()) != 0 {
		setPoolTx, err := tokenChainConfig.TokenAdminRegistry.SetPool(deployment.SimTransactOpts(), tokenChainConfig.TokenAddress, tokenChainConfig.TokenPool.Address())
		if err != nil {
			return []mcms.Operation{}, fmt.Errorf("failed to create setPool transaction: %w", err)
		}
		batch = append(batch, mcms.Operation{
			To:    tokenChainConfig.TokenAdminRegistry.Address(),
			Data:  setPoolTx.Data(),
			Value: zero(),
		})
	}

	// If an external admin is specified & timelock is currently the admin, transfer ownership of the pool and admin rights on the registry.
	// The timelock would be the owner of the pool at this point, so we don't need to check ownership there.
	if isTimelockAdmin && tokenChainConfig.ExternalAdmin.Cmp(zeroAddress()) != 0 {
		transferAdminRoleTx, err := tokenChainConfig.TokenAdminRegistry.TransferAdminRole(deployment.SimTransactOpts(), tokenChainConfig.TokenAddress, tokenChainConfig.ExternalAdmin)
		if err != nil {
			return []mcms.Operation{}, fmt.Errorf("failed to create transferAdminRole transaction: %w", err)
		}
		batch = append(batch, mcms.Operation{
			To:    tokenChainConfig.TokenAdminRegistry.Address(),
			Data:  transferAdminRoleTx.Data(),
			Value: zero(),
		})
		transferOwnershipTx, err := tokenChainConfig.TokenPool.TransferOwnership(deployment.SimTransactOpts(), tokenChainConfig.ExternalAdmin)
		if err != nil {
			return []mcms.Operation{}, fmt.Errorf("failed to create transferOwnership transaction: %w", err)
		}
		batch = append(batch, mcms.Operation{
			To:    tokenChainConfig.TokenPool.Address(),
			Data:  transferOwnershipTx.Data(),
			Value: zero(),
		})
	}

	return batch, nil
}
