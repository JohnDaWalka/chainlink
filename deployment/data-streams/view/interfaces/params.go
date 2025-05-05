package interfaces

import (
	"context"

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

// EthereumParams implements ChainParams
var _ ChainParams = (*EthereumParams)(nil)

// ContractView defines the base interface for any contract view
type ContractView interface {
	// SerializeView converts the view to a JSON string
	SerializeView() (string, error)
}

// ViewBuilder defines the interface for building contract views
type ViewBuilder interface {
	// BuildView constructs a view from blockchain data
	// Note the generic contract parameter - implement for specific contract types
	BuildView(ctx context.Context, contract interface{}, params ChainParams) (ContractView, error)
}
