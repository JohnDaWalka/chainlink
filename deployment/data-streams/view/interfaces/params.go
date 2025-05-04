package interfaces

import (
	"github.com/ethereum/go-ethereum/common"
	chainselectors "github.com/smartcontractkit/chain-selectors"
)

// ChainParams is a base interface for chain-specific parameters
type ChainParams interface {
	GetChainType() string
}

// EthereumParams contains Ethereum-specific parameters
type EthereumParams struct {
	FromBlock uint64
	ToBlock   *uint64
	Address   common.Address
}

func (p EthereumParams) GetChainType() string {
	return chainselectors.FamilyEVM
}

// Ensure EthereumParams implements ChainParams
var _ ChainParams = (*EthereumParams)(nil)
