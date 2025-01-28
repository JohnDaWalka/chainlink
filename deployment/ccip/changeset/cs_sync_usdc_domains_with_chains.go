package changeset

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/smartcontractkit/chainlink/deployment"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/usdc_token_pool"
)

var _ deployment.ChangeSet[SyncUSDCDomainsWithChainsConfig] = SyncUSDCDomainsWithChainsChangeset

type USDCChainConfig struct {
	// Version is the version of the USDC token pool.
	Version semver.Version
}

func (c USDCChainConfig) Validate(ctx context.Context, chain deployment.Chain, state CCIPChainState, useMcms bool, chainSelectorToDomainID map[uint64]uint32) error {
	usdcTokenPool, ok := state.USDCTokenPools[c.Version]
	if !ok {
		return fmt.Errorf("no USDC token pool found on %s with version %s", chain, c.Version)
	}

	if len(usdcTokenPool.Address().Bytes()) > 32 {
		// Will never be true for EVM
		return fmt.Errorf("expected USDC token pool address on %s (%s) to be less than 32 bytes", chain, usdcTokenPool.Address())
	}

	// Validate that the USDC token pool is owned by the address that will be actioning the transactions (i.e. Timelock or deployer key)
	if err := commoncs.ValidateOwnership(ctx, useMcms, chain.DeployerKey.From, state.Timelock.Address(), usdcTokenPool); err != nil {
		return fmt.Errorf("token pool with address %s on %s failed ownership validation: %w", usdcTokenPool.Address(), chain, err)
	}

	// Validate that each supported chain has a domain ID defined
	supportedChains, err := usdcTokenPool.GetSupportedChains(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("failed to get supported chains from USDC token pool on %s with address %s: %w", chain, usdcTokenPool.Address(), err)
	}
	for _, supportedChain := range supportedChains {
		if _, ok := chainSelectorToDomainID[supportedChain]; !ok {
			return fmt.Errorf("no USDC domain ID defined for chain with selector %d", supportedChain)
		}
	}

	return nil
}

// SyncUSDCDomainsWithChainsConfig defines the chain selector -> USDC domain mappings.
type SyncUSDCDomainsWithChainsConfig struct {
	// USDCConfigsByChain defines the USDC domain and pool version for each chain selector.
	USDCConfigsByChain map[uint64]USDCChainConfig
	// ChainSelectorToUSDCDomain maps chains selectors to their USDC domain identifiers.
	ChainSelectorToUSDCDomain map[uint64]uint32
	// MCMS defines the delay to use for Timelock (if absent, the changeset will attempt to use the deployer key).
	MCMS *MCMSConfig
}

func (c SyncUSDCDomainsWithChainsConfig) Validate(env deployment.Environment) error {
	state, err := LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	// Validate that all USDC configs inputted are for valid chains that define USDC pools.
	for chainSelector, config := range c.USDCConfigsByChain {
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
		if chainState.USDCTokenPools == nil {
			return fmt.Errorf("%s does not define any USDC token pools, config should be removed", chain)
		}
		if timelock := chainState.Timelock; timelock == nil {
			return fmt.Errorf("missing timelock on %s", chain.String())
		}
		if proposerMcm := chainState.ProposerMcm; proposerMcm == nil {
			return fmt.Errorf("missing proposerMcm on %s", chain.String())
		}
		if err = config.Validate(env.GetContext(), chain, chainState, c.MCMS != nil, c.ChainSelectorToUSDCDomain); err != nil {
			return fmt.Errorf("USDC config for %s is not valid: %w", chain, err)
		}
	}
	// Check that our input covers all chains that define USDC pools.
	for chainSelector, chainState := range state.Chains {
		if _, ok := c.USDCConfigsByChain[chainSelector]; !ok && chainState.USDCTokenPools != nil {
			return fmt.Errorf("no USDC chain config defined for %s, which does support USDC", env.Chains[chainSelector])
		}
	}
	return nil
}

// SyncUSDCDomainsWithChainsChangeset syncs domain support on specified USDC token pools with its chain support.
// As such, it is expected that ConfigureTokenPoolContractsChangeset is executed before running this changeset.
func SyncUSDCDomainsWithChainsChangeset(env deployment.Environment, c SyncUSDCDomainsWithChainsConfig) (deployment.ChangesetOutput, error) {
	if err := c.Validate(env); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid SyncUSDCDomainsWithChainsConfig: %w", err)
	}
	readOpts := &bind.CallOpts{Context: env.GetContext()}

	state, err := LoadOnchainState(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	deployerGroup := NewDeployerGroup(env, state, c.MCMS)

	for chainSelector, usdcChainConfig := range c.USDCConfigsByChain {
		chain := env.Chains[chainSelector]
		chainState := state.Chains[chainSelector]
		writeOpts, err := deployerGroup.GetDeployer(chainSelector)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to get transaction opts for %s", chain)
		}

		usdcTokenPool := chainState.USDCTokenPools[usdcChainConfig.Version]
		supportedChains, err := usdcTokenPool.GetSupportedChains(readOpts)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to fetch supported chains from USDC token pool with address %s on %s: %w", usdcTokenPool.Address(), chain, err)
		}

		domainUpdates := make([]usdc_token_pool.USDCTokenPoolDomainUpdate, 0)
		for _, remoteChainSelector := range supportedChains {
			remoteChainState := state.Chains[remoteChainSelector]
			remoteUSDCChainConfig := c.USDCConfigsByChain[remoteChainSelector]
			remoteUSDCTokenPool := remoteChainState.USDCTokenPools[remoteUSDCChainConfig.Version]

			var desiredAllowedCaller [32]byte
			remoteUSDCTokenPoolAddressBytes := remoteUSDCTokenPool.Address().Bytes()
			for i, j := len(desiredAllowedCaller)-len(remoteUSDCTokenPoolAddressBytes), 0; i < len(desiredAllowedCaller); i, j = i+1, j+1 {
				desiredAllowedCaller[i] = remoteUSDCTokenPoolAddressBytes[j]
			}
			desiredDomainIdentifier := c.ChainSelectorToUSDCDomain[remoteChainSelector]

			currentDomain, err := usdcTokenPool.GetDomain(readOpts, remoteChainSelector)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to fetch domain for %d from USDC token pool with address %s on %s: %w", remoteChainSelector, usdcTokenPool.Address(), chain, err)
			}
			// If any parameters are different, we need to add a setDomains call
			if currentDomain.AllowedCaller != desiredAllowedCaller ||
				currentDomain.DomainIdentifier != desiredDomainIdentifier {
				domainUpdates = append(domainUpdates, usdc_token_pool.USDCTokenPoolDomainUpdate{
					AllowedCaller:     desiredAllowedCaller,
					Enabled:           true,
					DomainIdentifier:  desiredDomainIdentifier,
					DestChainSelector: remoteChainSelector,
				})
			}
		}

		if len(domainUpdates) > 0 {
			_, err := usdcTokenPool.SetDomains(writeOpts, domainUpdates)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to create set domains operation on %s: %w", chain, err)
			}
		}
	}

	return deployerGroup.Enact("sync domain support with chain support on USDC token pools")
}
