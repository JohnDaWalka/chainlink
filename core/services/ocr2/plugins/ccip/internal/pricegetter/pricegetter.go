package pricegetter

import (
	"context"
	"math/big"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

type PriceGetter interface {
	cciptypes.PriceGetter
}

type AllTokensPriceGetter interface {
	PriceGetter
	// GetJobSpecTokenPricesUSD returns all token prices defined in the jobspec.
	GetJobSpecTokenPricesUSD(ctx context.Context) (map[TokenID]*big.Int, error)
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
