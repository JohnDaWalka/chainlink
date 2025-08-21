package config

import (
	"math/big"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
)

type Config struct {
	ChainSelector   uint64
	ContractAddress []byte
	AccountAddress  []byte
	ExpectedBalance *big.Int
	TxHash          []byte
	ExpectedReceipt *evm.Receipt
	ExpectedTx      *evm.Transaction
}
