package zksync

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	zkSyncTypes "github.com/zksync-sdk/zksync2-go/types"

	"github.com/zksync-sdk/zksync2-go/accounts"

	"github.com/smartcontractkit/chainlink/deployment"
)

func PickDeployFn[C any](
	chain deployment.Chain,
	evmDeployFn func(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, C, error),
	zkDeployFn func(auth *bind.TransactOpts, ethClient *ethclient.Client, wallet accounts.Wallet, args ...interface{}) (common.Address, *zkSyncTypes.Receipt, C, error),
	args ...interface{},
) (common.Address, *types.Transaction, C, error) {
	if chain.IsZk {
		address, _, c, err := zkDeployFn(chain.DeployerKey, chain.Client.(*ethclient.Client), *chain.DeployerKeyZk, args...)
		return address, nil, c, err
	}

	return evmDeployFn(chain.DeployerKey, chain.Client)
}
