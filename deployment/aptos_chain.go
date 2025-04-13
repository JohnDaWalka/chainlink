package deployment

import (
	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_offramp"
	module_offramp "github.com/smartcontractkit/chainlink-aptos/bindings/ccip_offramp/offramp"
)

// AptosChain represents an Aptos chain.
type AptosChain struct {
	Selector       uint64
	Client         aptos.AptosRpcClient
	DeployerSigner aptos.TransactionSigner
}

func (c AptosChain) GetOfframpDynamicConfig(ccipAddress aptos.AccountAddress) (module_offramp.DynamicConfig, error) {
	offrampBindings := ccip_offramp.Bind(ccipAddress, c.Client)
	return offrampBindings.Offramp().GetDynamicConfig(&bind.CallOpts{})
}
