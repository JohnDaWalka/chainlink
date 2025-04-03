package pricegetter

import (
	"context"
	"io"
	"math/big"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

type AllTokensPriceGetter interface {
	// GetJobSpecTokenPricesUSD returns all token prices defined in the jobspec.
	GetJobSpecTokenPricesUSD(ctx context.Context) (map[TokenID]*big.Int, error)
	// GetTokenPrices returns the prices for the provided tokens in USD.
	GetTokenPrices(ctx context.Context, tokens []TokenID) (map[TokenID]*big.Int, error)
	io.Closer
}

// TokenID is a struct that represents a token's address and chain ID.
type TokenID struct {
	Address cciptypes.Address
	ChainID uint64
}

func NewTokenID(address cciptypes.Address, chainID uint64) TokenID {
	return TokenID{
		Address: address,
		ChainID: chainID,
	}
}
