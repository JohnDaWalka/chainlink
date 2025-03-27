package deployment

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"
	zkTypes "github.com/zksync-sdk/zksync2-go/types"
)

func PickXVMDeployFn[C any](
	chain Chain,
	deployEVM func(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, C, error),
	deployZKVM func(deployOpts *accounts.TransactOpts, client *clients.Client, wallet *accounts.Wallet, backend bind.ContractBackend, constructorArgs ...interface{}) (common.Address, *zkTypes.Receipt, C, error),
	constructorArgs ...interface{},
) (common.Address, *types.Transaction, C, error) {
	if !chain.IsZk {
		return deployEVM(chain.DeployerKey, chain.Client) // , constructorArgs...)
	}

	address, _, contract, err := deployZKVM(nil, chain.ClientZk, chain.DeployerKeyZk, chain.Client, constructorArgs...)
	return address, nil, contract, err
}
