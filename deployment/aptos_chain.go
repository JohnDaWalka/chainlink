package deployment

import (
	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	ccip_offramp "github.com/smartcontractkit/chainlink-aptos/bindings/ccip_offramp"
	module_offramp "github.com/smartcontractkit/chainlink-aptos/bindings/ccip_offramp/offramp"

	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

// ChainAdapterAptos is a Chainlink-specific extension of the generic AptosChain.
// It wraps the shared deployment.AptosChain and adds product-specific methods
// (e.g. for interacting with CCIP modules). This type should be the sole
// entry point for Aptos chain logic that depends on Chainlink bindings or behaviors.
type ChainAdapterAptos struct {
	deployment.AptosChain
}

func NewChainAdapterAptos(chain deployment.AptosChain) ChainAdapterAptos {
	return ChainAdapterAptos{AptosChain: chain}
}

func (a ChainAdapterAptos) GetOfframpDynamicConfig(ccipAddr aptos.AccountAddress) (module_offramp.DynamicConfig, error) {
	return ccip_offramp.
		Bind(ccipAddr, a.Client).
		Offramp().
		GetDynamicConfig(&bind.CallOpts{})
}
