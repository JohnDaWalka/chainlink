package changesets_zksync

import (
	"context"
	"crypto/rand"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zksync-sdk/zksync2-go/contracts/contractdeployer"
	"github.com/zksync-sdk/zksync2-go/utils"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

type DeployFn[C any] func(chain deployment.Chain) deployment.ContractDeploy[*C]

func WrapDeployFn[C any](
	c deployment.Chain,
	deployEVM DeployFn[C],
	zkBytecode []byte,
	getAbi func() (*abi.ABI, error),
	args []interface{},
	newContract func(common.Address, bind.ContractBackend) (*C, error),
) DeployFn[C] {
	if c.IsZK {
		return wrapZKDeployFn(zkBytecode, args, getAbi, newContract)
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
) func(chain deployment.Chain) deployment.ContractDeploy[*C] {
	return func(chain deployment.Chain) deployment.ContractDeploy[*C] {
		contractDeployerAddress := utils.ContractDeployerAddress
		contractDeployer, err := contractdeployer.NewIContractDeployer(contractDeployerAddress, chain.Client)
		if err != nil {
			return wrapErr[C](err)
		}

		salt := make([]byte, 32)
		rand.Read(salt)

		input := make([]byte, len(bytecode))
		copy(input, bytecode)

		if len(args) > 0 {
			abi, err := getAbi()
			if err != nil {
				return wrapErr[C](err)
			}

			data, err := abi.Pack("", args...)
			if err != nil {
				return wrapErr[C](err)
			}

			input = append(input, data...)
		}

		bytecodeHash, err := utils.HashBytecode(input)
		if err != nil {
			return wrapErr[C](err)
		}

		chain.DeployerKey.GasPrice, err = chain.Client.SuggestGasPrice(context.Background())
		if err != nil {
			return wrapErr[C](err)
		}

		tx, err := contractDeployer.Create2(chain.DeployerKey, [32]byte(salt), [32]byte(bytecodeHash), input)
		if err != nil {
			return wrapErr[C](err)
		}

		receipt, err := bind.WaitMined(context.Background(), chain.Client, tx)
		if err != nil {
			return wrapErr[C](err)
		}

		address := common.HexToAddress("0x" + receipt.Logs[1].Topics[3].String()[2+32*2-40:])
		linkToken, err2 := newContract(address, chain.Client)

		return deployment.ContractDeploy[*C]{
			Address:  address,
			Contract: linkToken,
			Tx:       tx,
			Tv:       deployment.NewTypeAndVersion(types.LinkToken, deployment.Version1_0_0),
			Err:      err2,
		}
	}
}
