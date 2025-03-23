package deployment

import (
	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip"
	module_offramp "github.com/smartcontractkit/chainlink-aptos/bindings/ccip/offramp"
)

// AptosChain represents an Aptos chain.
type AptosChain struct {
	Selector       uint64
	Client         aptos.AptosRpcClient
	DeployerSigner aptos.TransactionSigner
}

func (c AptosChain) GetOfframpDynamicConfig(ccipAddress aptos.AccountAddress) (module_offramp.DynamicConfig, error) {
	ccipBindings := ccip.Bind(ccipAddress, c.Client)
	return ccipBindings.Offramp.GetDynamicConfig(&bind.CallOpts{})
}
