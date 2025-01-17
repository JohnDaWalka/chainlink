package changeset

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_pool"
)

// getAllTokenPoolsWithSymbol returns a list of all token pools tied to the given token symbol.
func GetAllTokenPoolsWithSymbol(
	chainState CCIPChainState,
	chainClient deployment.OnchainClient,
	symbol TokenSymbol,
) ([]*token_pool.TokenPool, error) {
	var tokenPools []*token_pool.TokenPool
	appendTokenPool := func(address common.Address) error {
		tp, err := token_pool.NewTokenPool(address, chainClient)
		if err != nil {
			return fmt.Errorf("failed to connect address %s with token pool bindings: %w", address, err)
		}
		tokenPools = append(tokenPools, tp)
		return nil
	}

	for _, pool := range chainState.BurnMintTokenPools[symbol] {
		err := appendTokenPool(pool.Address())
		if err != nil {
			return nil, err
		}
	}
	for _, pool := range chainState.BurnWithFromMintTokenPools[symbol] {
		err := appendTokenPool(pool.Address())
		if err != nil {
			return nil, err
		}
	}
	for _, pool := range chainState.BurnFromMintTokenPools[symbol] {
		err := appendTokenPool(pool.Address())
		if err != nil {
			return nil, err
		}
	}
	for _, pool := range chainState.LockReleaseTokenPools[symbol] {
		err := appendTokenPool(pool.Address())
		if err != nil {
			return nil, err
		}
	}

	return tokenPools, nil
}

// getTokenPoolWithSymbolAndAddress returns the token pool in the environment linked to both the given symbol and pool address.
func GetTokenPoolWithSymbolAndAddress(chainState CCIPChainState, chain deployment.Chain, symbol TokenSymbol, address common.Address) (*token_pool.TokenPool, error) {
	tokenPools, err := GetAllTokenPoolsWithSymbol(chainState, chain.Client, symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to get %s token pools on %s: %w", symbol, chain.Name(), err)
	}
	var desiredTokenPool *token_pool.TokenPool
	for _, tokenPool := range tokenPools {
		if tokenPool.Address() == address {
			desiredTokenPool = tokenPool
			break
		}
	}
	if desiredTokenPool == nil {
		return nil, fmt.Errorf("no token pool found with symbol %s and address %s on chain %s", symbol, address, chain.Name())
	}
	return desiredTokenPool, nil
}
