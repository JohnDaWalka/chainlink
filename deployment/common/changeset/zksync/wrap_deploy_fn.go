package changesets_zksync

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"

	"github.com/smartcontractkit/chainlink/deployment"
)

type DeployEVMFn[C any] func(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *C, error)
type DeployZKFn[C any] func(auth *bind.TransactOpts, ethClient *ethclient.Client, pk string, args ...interface{}) (common.Address, *types.Transaction, *C, error)

func PickDeployFn[C any](
	c deployment.Chain,
	deployEVM DeployEVMFn[C],
	deployZk DeployZKFn[C],
) (common.Address, *types.Transaction, *C, error) {
	if c.IsZK {
		return deployZk(c.DeployerKey, c.Client.(*ethclient.Client), c.DeployerPk)
	}

	return deployEVM(c.DeployerKey, c.Client)
}

type DeployFn[C any] func(chain deployment.Chain) deployment.ContractDeploy[*C]

func WrapDeployFn[C any](
	c deployment.Chain,
	deployEVM DeployFn[C],
	zkBytecode []byte,
	getAbi func() (*abi.ABI, error),
	args []interface{},
	newContract func(common.Address, bind.ContractBackend) (*C, error),
	tv deployment.TypeAndVersion,
) DeployFn[C] {
	if c.IsZK {
		return wrapZKDeployFn(zkBytecode, args, getAbi, newContract, tv)
	} else {
		return deployEVM
	}
}

func wrapErr[C any](err error) deployment.ContractDeploy[*C] {
	return deployment.ContractDeploy[*C]{
		Err: err,
	}
}

func wrapZKDeployFn[C any](
	bytecode []byte,
	args []interface{},
	getAbi func() (*abi.ABI, error),
	newContract func(common.Address, bind.ContractBackend) (*C, error),
	tv deployment.TypeAndVersion,
) func(chain deployment.Chain) deployment.ContractDeploy[*C] {
	return func(chain deployment.Chain) deployment.ContractDeploy[*C] {
		pk := common.Hex2Bytes(chain.DeployerPk)
		client := clients.NewClient(chain.Client.(*ethclient.Client).Client())
		wallet, err := accounts.NewWallet(pk, client, nil)
		if err != nil {
			return wrapErr[C](err)
		}

		var calldata []byte
		if len(args) > 0 {
			abi, err := getAbi()
			if err != nil {
				return wrapErr[C](err)
			}
			calldata, err = abi.Pack("", args...)
			if err != nil {
				return wrapErr[C](err)
			}
		}

		txHash, err := wallet.Deploy(nil, accounts.Create2Transaction{
			Bytecode: bytecode,
			Calldata: calldata})
		if err != nil {
			return wrapErr[C](err)
		}

		receipt, err := client.WaitMined(context.Background(), txHash)
		if err != nil {
			return wrapErr[C](err)
		}

		address := receipt.ContractAddress
		contract, err2 := newContract(address, chain.Client)

		return deployment.ContractDeploy[*C]{
			Address:  address,
			Contract: contract,
			Tx:       nil,
			Tv:       tv,
			Err:      err2,
		}
	}
}
