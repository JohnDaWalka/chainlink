package zksyncwrapper

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"
	"github.com/zksync-sdk/zksync2-go/types"
)

// this file is used as a template. see wrap_zk_bytecode.go before editing
func DeployPlaceholderContractNameZk(auth *bind.TransactOpts, ethClient *ethclient.Client, wallet accounts.Wallet, args ...interface{}) (common.Address, *types.Receipt, *PlaceholderContractName, error) {
	var calldata []byte
	if len(args) > 0 {
		abi, err := PlaceholderContractNameMetaData.GetAbi()
		if err != nil {
			return common.Address{}, nil, nil, err
		}
		calldata, err = abi.Pack("", args...)
		if err != nil {
			return common.Address{}, nil, nil, err
		}
	}

	txHash, err := wallet.Deploy(&accounts.TransactOpts{
		Nonce:     auth.Nonce,
		Value:     auth.Value,
		GasPrice:  auth.GasPrice,
		GasLimit:  auth.GasLimit,
		GasFeeCap: auth.GasFeeCap,
		GasTipCap: auth.GasTipCap,
		Context:   auth.Context,
	}, accounts.Create2Transaction{
		Bytecode: ZkBytecode,
		Calldata: calldata})
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	client := clients.NewClient(ethClient.Client())
	receipt, err := client.WaitMined(context.Background(), txHash)
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address := receipt.ContractAddress
	contract, err := NewPlaceholderContractName(address, ethClient)
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	return address, receipt, contract, nil
}
