package generated_zks

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	// zkSyncAccounts "github.com/zksync-sdk/zksync2-go/accounts"
	// zkSyncClient "github.com/zksync-sdk/zksync2-go/clients"
	// zktypes "github.com/zksync-sdk/zksync2-go/types"
)

type CustomTransaction struct {
	*types.Transaction
	CustomHash common.Hash
}

func (tx *CustomTransaction) Hash() common.Hash {
	return tx.CustomHash
}

// func ConvertToTransaction(resp zktypes.TransactionResponse) *CustomTransaction {
// 	dtx := &types.DynamicFeeTx{
// 		ChainID:   resp.ChainID.ToInt(),
// 		Nonce:     uint64(resp.Nonce),
// 		GasTipCap: resp.MaxPriorityFeePerGas.ToInt(),
// 		GasFeeCap: resp.MaxFeePerGas.ToInt(),
// 		To:        &resp.To,
// 		Value:     resp.Value.ToInt(),
// 		Data:      resp.Data,
// 		Gas:       uint64(resp.Gas),
// 	}

// 	tx := types.NewTx(dtx)
// 	customTransaction := CustomTransaction{Transaction: tx, CustomHash: resp.Hash}
// 	return &customTransaction
// }

// func IsZkSync(backend bind.ContractBackend) bool {
// 	client, ok := backend.(*ethclient.Client)
// 	if !ok {
// 		return false
// 	}
// 	chainId, err := client.ChainID(context.Background())
// 	if err != nil {
// 		return false
// 	}
// 	switch chainId.Uint64() {

// 	case 324, 280, 300:
// 		return true
// 	}
// 	return false
// }

// type ZkSyncContract struct {
// 	zkclient *zkSyncClient.Client
// 	wallet   *zkSyncAccounts.Wallet
// }

// // NewZkSyncContract creates and returns a new ZkSyncContract instance.
// func NewZkSyncContract(auth *bind.TransactOpts, backend bind.ContractBackend, contractBytes []byte, contractAbi *abi.ABI, params ...interface{}) (common.Address, *CustomTransaction, *bind.BoundContract, error) {
// 	client, ok := backend.(*ethclient.Client)
// 	if !ok {
// 		return common.Address{}, nil, nil, fmt.Errorf("backend is not an *ethclient.Client")
// 	}

// 	// Retrieve wallet from context safely
// 	walletValue := auth.Context.Value("wallet")
// 	wallet, ok := walletValue.(*zkSyncAccounts.Wallet)
// 	if !ok || wallet == nil {
// 		return common.Address{}, nil, nil, fmt.Errorf("wallet not found in context or invalid type")
// 	}

// 	zkclient := zkSyncClient.NewClient(client.Client())

// 	constructor, _ := contractAbi.Pack("", params...)

// 	hash, _ := wallet.DeployWithCreate(nil, zkSyncAccounts.CreateTransaction{
// 		Bytecode: contractBytes,
// 		Calldata: constructor,
// 	})
// 	receipt, _ := zkclient.WaitMined(context.Background(), hash)
// 	tx, _, _ := zkclient.TransactionByHash(context.Background(), hash)
// 	ethTx := ConvertToTransaction(*tx)
// 	address := receipt.ContractAddress
// 	contractBind := bind.NewBoundContract(address, *contractAbi, backend, backend, backend)
// 	return address, ethTx, contractBind, nil
// }
