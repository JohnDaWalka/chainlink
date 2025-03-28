package deployment

import (
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"
	zkTypes "github.com/zksync-sdk/zksync2-go/types"
)

// type DeployEVMFn[C any] func(auth *bind.TransactOpts, backend bind.ContractBackend, args ...any) (common.Address, *types.Transaction, *C, error)
type DeployZKVMFn[C any] func(deployOpts *accounts.TransactOpts, client *clients.Client, wallet *accounts.Wallet, backend bind.ContractBackend, args ...any) (common.Address, *zkTypes.Receipt, *C, error)

func PickXVMDeployFn[C any](
	chain Chain,
	deployEVM any,
	deployZKVM DeployZKVMFn[C],
	args ...any,
) (common.Address, *types.Transaction, *C, error) {
	if !chain.IsZk {
		deployEVMReflection := reflect.ValueOf(deployEVM)
		values := make([]reflect.Value, len(args)+2)
		values[0] = reflect.ValueOf(chain.DeployerKey)
		values[1] = reflect.ValueOf(chain.Client)
		for i, arg := range args {
			values[i+2] = reflect.ValueOf(arg)
		}
		result := deployEVMReflection.Call(values)
		if result[3].Interface() != nil {
			return common.Address{}, nil, nil, result[3].Interface().(error)
		}
		return result[0].Interface().(common.Address), result[1].Interface().(*types.Transaction), result[2].Interface().(*C), nil
	}

	address, _, contract, err := deployZKVM(nil, chain.ClientZk, chain.DeployerKeyZk, chain.Client, args...)
	return address, nil, contract, err
}
