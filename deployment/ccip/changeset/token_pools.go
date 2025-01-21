package changeset

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/erc20"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
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

// tokenPool defines behavior common to all token pools.
type tokenPool interface {
	GetToken(opts *bind.CallOpts) (common.Address, error)
	TypeAndVersion(*bind.CallOpts) (string, error)
}

// tokenPoolMetadata defines the token pool version version and symbol of the corresponding token.
type tokenPoolMetadata struct {
	Version semver.Version
	Symbol  TokenSymbol
}

// newTokenPoolWithMetadata returns a token pool along with its metadata.
func newTokenPoolWithMetadata[P tokenPool](
	newTokenPool func(address common.Address, backend bind.ContractBackend) (P, error),
	poolAddress common.Address,
	chainClient deployment.OnchainClient,
) (P, tokenPoolMetadata, error) {
	pool, err := newTokenPool(poolAddress, chainClient)
	if err != nil {
		return pool, tokenPoolMetadata{}, fmt.Errorf("failed to connect address %s with token pool bindings: %w", poolAddress, err)
	}
	tokenAddress, err := pool.GetToken(nil)
	if err != nil {
		return pool, tokenPoolMetadata{}, fmt.Errorf("failed to get token address from pool with address %s: %w", poolAddress, err)
	}
	typeAndVersionStr, err := pool.TypeAndVersion(nil)
	if err != nil {
		return pool, tokenPoolMetadata{}, fmt.Errorf("failed to get type and version from pool with address %s: %w", poolAddress, err)
	}
	_, versionStr, err := ccipconfig.ParseTypeAndVersion(typeAndVersionStr)
	if err != nil {
		return pool, tokenPoolMetadata{}, fmt.Errorf("failed to parse type and version of pool with address %s: %w", poolAddress, err)
	}
	version, err := semver.NewVersion(versionStr)
	if err != nil {
		return pool, tokenPoolMetadata{}, fmt.Errorf("failed parsing version %s of pool with address %s: %w", versionStr, poolAddress, err)
	}
	token, err := erc20.NewERC20(tokenAddress, chainClient)
	if err != nil {
		return pool, tokenPoolMetadata{}, fmt.Errorf("failed to connect address %s with ERC20 bindings: %w", tokenAddress, err)
	}
	symbol, err := token.Symbol(nil)
	if err != nil {
		return pool, tokenPoolMetadata{}, fmt.Errorf("failed to fetch symbol from token with address %s: %w", tokenAddress, err)
	}
	return pool, tokenPoolMetadata{
		Symbol:  TokenSymbol(symbol),
		Version: *version,
	}, nil
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
