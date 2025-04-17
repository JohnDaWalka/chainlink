package evm

import (
	"context"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commonservices "github.com/smartcontractkit/chainlink-common/pkg/services"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	evmtxmgr "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
)

type EVMChain struct {
	logger logger.Logger
	commonservices.StateMachine
	txm evmtxmgr.TxManager
}

func NewEVMChain(_ context.Context, lggr logger.Logger, txm evmtxmgr.TxManager) (*EVMChain, error) {
	return &EVMChain{
		logger: lggr,
		txm:    txm,
	}, nil
}

func (e *EVMChain) ReadContract(ctx context.Context, method string, encodedParams []byte) ([]byte, error) {
	return nil, errors.New("yet to implement")
}

// GetTransactionFee retrieves the fee of a transaction in the underlying chain's TXM
func (e *EVMChain) GetTransactionFee(ctx context.Context, transactionID string) (*commontypes.TransactionFee, error) {
	return e.txm.GetTransactionFee(ctx, transactionID)
}
func (e *EVMChain) Close() error {
	return e.StopOnce(e.Name(), func() error {
		return nil
	})
}

func (e *EVMChain) HealthReport() map[string]error {
	return map[string]error{e.Name(): e.Healthy()}
}

func (e *EVMChain) Name() string {
	return e.logger.Name()
}

func (e *EVMChain) Start(ctx context.Context) error {
	return e.StartOnce(e.Name(), func() error {
		return nil
	})
}
