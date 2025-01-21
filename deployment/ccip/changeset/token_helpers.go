package changeset

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_pool"
)

var currentTokenPoolVersion semver.Version = deployment.Version1_5_1

var tokenPoolTypes map[deployment.ContractType]struct{} = map[deployment.ContractType]struct{}{
	BurnMintTokenPool:         struct{}{},
	BurnWithFromMintTokenPool: struct{}{},
	BurnFromMintTokenPool:     struct{}{},
	LockReleaseTokenPool:      struct{}{},
}

var tokenPoolVersions map[semver.Version]struct{} = map[semver.Version]struct{}{
	deployment.Version1_5_1: struct{}{},
}

// getAllTokenPoolsWithSymbolAndVersion returns a list of all token pools tied to the given token symbol and semver version.
func getAllTokenPoolsWithSymbolAndVersion(
	chainState CCIPChainState,
	chainClient deployment.OnchainClient,
	symbol TokenSymbol,
	requestedVersion semver.Version,
) ([]*token_pool.TokenPool, error) {
	var tokenPools []*token_pool.TokenPool
	appendPoolIfVersionIsCorrect := func(address common.Address, version semver.Version) error {
		if version.Equal(&requestedVersion) {
			tp, err := token_pool.NewTokenPool(address, chainClient)
			if err != nil {
				return fmt.Errorf("failed to connect address %s with token pool bindings: %w", address, err)
			}
			tokenPools = append(tokenPools, tp)
			return nil
		}
		return nil
	}

	for version, pool := range chainState.BurnMintTokenPools[symbol] {
		err := appendPoolIfVersionIsCorrect(pool.Address(), version)
		if err != nil {
			return nil, err
		}
	}
	for version, pool := range chainState.BurnWithFromMintTokenPools[symbol] {
		err := appendPoolIfVersionIsCorrect(pool.Address(), version)
		if err != nil {
			return nil, err
		}
	}
	for version, pool := range chainState.BurnFromMintTokenPools[symbol] {
		err := appendPoolIfVersionIsCorrect(pool.Address(), version)
		if err != nil {
			return nil, err
		}
	}
	for version, pool := range chainState.LockReleaseTokenPools[symbol] {
		err := appendPoolIfVersionIsCorrect(pool.Address(), version)
		if err != nil {
			return nil, err
		}
	}

	return tokenPools, nil
}

// getTokenPoolFromSymbolTypeAndVersion returns the token pool in the environment linked to a particular symbol, type, and version
func getTokenPoolFromSymbolTypeAndVersion(
	chainState CCIPChainState,
	chain deployment.Chain,
	symbol TokenSymbol,
	poolType deployment.ContractType,
	version semver.Version,
) (*token_pool.TokenPool, error) {
	switch poolType {
	case BurnMintTokenPool:
		if tokenPools, ok := chainState.BurnMintTokenPools[symbol]; ok {
			if tokenPool, ok := tokenPools[version]; ok {
				return token_pool.NewTokenPool(tokenPool.Address(), chain.Client)
			}
		}
	case BurnFromMintTokenPool:
		if tokenPools, ok := chainState.BurnFromMintTokenPools[symbol]; ok {
			if tokenPool, ok := tokenPools[version]; ok {
				return token_pool.NewTokenPool(tokenPool.Address(), chain.Client)
			}
		}
	case BurnWithFromMintTokenPool:
		if tokenPools, ok := chainState.BurnWithFromMintTokenPools[symbol]; ok {
			if tokenPool, ok := tokenPools[version]; ok {
				return token_pool.NewTokenPool(tokenPool.Address(), chain.Client)
			}
		}
	case LockReleaseTokenPool:
		if tokenPools, ok := chainState.LockReleaseTokenPools[symbol]; ok {
			if tokenPool, ok := tokenPools[version]; ok {
				return token_pool.NewTokenPool(tokenPool.Address(), chain.Client)
			}
		}
	}

	return nil, fmt.Errorf("failed to find token pool with symbol %s, type %s, and version %s", symbol, poolType, version)
}
