package txmgr

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"time"

	pkgerrors "github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-framework/chains/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink-framework/chains/txmgr/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txm/types"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
)

func (o *evmTxStore) AppendAttemptToTransaction(ctx context.Context, _ *types.Transaction, attempt *types.Attempt) error {
	var cancel context.CancelFunc
	ctx, cancel = o.stopCh.Ctx(ctx)
	defer cancel()

	var dbTxAttempt DbEthTxAttempt
	dbTxAttempt.FromAttempt(attempt)
	query, args, err := o.q.BindNamed(insertIntoEthTxAttemptsQuery, &dbTxAttempt)
	if err != nil {
		return pkgerrors.Wrap(err, "AppendAttemptToTransaction failed to bind named")
	}
	err = o.q.GetContext(ctx, &dbTxAttempt, query, args...)
	dbTxAttempt.ToAttempt(attempt)
	return pkgerrors.Wrap(err, "AppendAttemptToTransaction failed")
}

func (o *evmTxStore) CreateEmptyUnconfirmedTransaction(ctx context.Context, fromAddress common.Address, txNonce uint64, gasLimit uint64, chainID *big.Int) (tx *types.Transaction, err error) {
	var cancel context.CancelFunc
	ctx, cancel = o.stopCh.Ctx(ctx)
	defer cancel()

	now := time.Now()
	emptyTx := &types.Transaction{
		ChainID:           chainID,
		Nonce:             &txNonce,
		FromAddress:       fromAddress,
		ToAddress:         common.Address{},
		Value:             big.NewInt(0),
		SpecifiedGasLimit: gasLimit,
		CreatedAt:         now,
		State:             txmgr.TxUnconfirmed,

		// This is required for backwards compatibility
		Data:               []byte{0},
		InitialBroadcastAt: &now,
		LastBroadcastAt:    &now,
	}

	var dbTx DbEthTx
	dbTx.FromTransaction(emptyTx)

	query, args, err := o.q.BindNamed(`INSERT INTO evm.txes (nonce, from_address, to_address, encoded_payload, value, gas_limit, error, broadcast_at, initial_broadcast_at, created_at,
state, meta, subject, pipeline_task_run_id, min_confirmations, evm_chain_id, transmit_checker, idempotency_key, signal_callback, callback_completed) VALUES (
:nonce, :from_address, :to_address, :encoded_payload, :value, :gas_limit, :error, :broadcast_at, :initial_broadcast_at, :created_at, :state, :meta, :subject,
:pipeline_task_run_id, :min_confirmations, :evm_chain_id, :transmit_checker, :idempotency_key, :signal_callback, :callback_completed
) RETURNING *`, &dbTx)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "InsertTx failed to bind named")
	}
	err = o.q.GetContext(ctx, &dbTx, query, args...)
	tx = new(types.Transaction)
	dbTx.ToTransaction(tx)
	return tx, pkgerrors.Wrap(err, "InsertTx failed")
}

func (o *evmTxStore) CreateTx(ctx context.Context, txRequest *types.TxRequest, maxQueueSize uint64) (*types.Transaction, error) {
	var cancel context.CancelFunc
	ctx, cancel = o.stopCh.Ctx(ctx)
	defer cancel()
	var dbTx DbEthTx
	err := o.Transact(ctx, false, func(orm *evmTxStore) error {
		// Pipeline
		if txRequest.PipelineTaskRunID.Valid {
			err := orm.q.GetContext(ctx, &dbTx, `SELECT * FROM evm.txes WHERE pipeline_task_run_id = $1 AND evm_chain_id = $2`, txRequest.PipelineTaskRunID, txRequest.ChainID.String())
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					// if a previous transaction for this task run exists, immediately return it
					return nil
				}
				return pkgerrors.Wrap(err, "CreateTx")
			}
		}
		// Pruning
		pruned, qErr := pruneMaxQueue(ctx, maxQueueSize, orm)
		if qErr != nil {
			return pkgerrors.Wrap(qErr, "PruneUnstartedTxQueue failed")
		}
		if len(pruned) > 0 {
			o.logger.Warnf(fmt.Sprintf("Pruned %d old unstarted transactions", len(pruned)),
				"pruned-tx-ids", pruned,
			)
		}
		// Transaction creation
		err := orm.q.GetContext(ctx, &dbTx, `
INSERT INTO evm.txes (from_address, to_address, encoded_payload, value, gas_limit, state, created_at, meta, evm_chain_id, min_confirmations, pipeline_task_run_id, idempotency_key, signal_callback)
VALUES ($1,$2,$3,$4,$5,'unstarted',NOW(),$6,$7,$8,$9,$10,$11)
RETURNING "txes".*
`, txRequest.FromAddress, txRequest.ToAddress, txRequest.Data, assets.Eth(*txRequest.Value), txRequest.SpecifiedGasLimit, txRequest.Meta,
			txRequest.ChainID.String(), txRequest.MinConfirmations, txRequest.PipelineTaskRunID, txRequest.IdempotencyKey, txRequest.SignalCallback)
		if err != nil {
			return pkgerrors.Wrap(err, "CreateEthTransaction failed to insert evm tx")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	tx := new(types.Transaction)
	dbTx.ToTransaction(tx)
	return tx, err
}

func (o *evmTxStore) FetchUnconfirmedTransactionAtNonceWithCount(ctx context.Context, latestNonce uint64, fromAddress common.Address, chainID *big.Int) (txCopy *types.Transaction, unconfirmedCount int, err error) {
	var cancel context.CancelFunc
	ctx, cancel = o.stopCh.Ctx(ctx)
	defer cancel()

	var dbTx DbEthTx
	err = o.Transact(ctx, false, func(orm *evmTxStore) error {
		err = orm.q.GetContext(ctx, &dbTx, `SELECT * FROM evm.txes WHERE from_address=$1 AND evm_chain_id=$2 AND nonce=$3 AND state='unconfirmed'`, fromAddress, chainID.String(), latestNonce)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return pkgerrors.Wrap(err, "FetchUnconfirmedTransactionAtNonceWithCount failed")
			}
		} else {
			txCopy = new(types.Transaction)
			dbTx.ToTransaction(txCopy)
			if err = loadTxAttempts(ctx, orm, txCopy); err != nil {
				return pkgerrors.Wrap(err, "FetchUnconfirmedTransactionAtNonceWithCount failed to load evm.tx_attempts")
			}
		}
		if err = orm.q.GetContext(ctx, &unconfirmedCount, `SELECT count(*) FROM evm.txes WHERE from_address=$1 AND evm_chain_id=$2 AND state='unconfirmed'`, fromAddress, chainID.String()); err != nil {
			return pkgerrors.Wrap(err, "FetchUnconfirmedTransactionAtNonceWithCount failed to count unconfirmed txs")
		}
		return nil
	})
	return
}

func (o *evmTxStore) MarkConfirmedAndReorgedTransactions(ctx context.Context, latestNonce uint64, fromAddress common.Address, chainID *big.Int) (confirmedTxs []*types.Transaction, reorgTxIDs []uint64, err error) {
	var cancel context.CancelFunc
	ctx, cancel = o.stopCh.Ctx(ctx)
	defer cancel()

	err = o.Transact(ctx, true, func(orm *evmTxStore) error {
		// Re-orged transactions
		query := `UPDATE evm.txes SET state='unconfirmed' WHERE from_address=$1 AND state IN ('confirmed', 'finalized') AND nonce>=$2 AND evm_chain_id=$3 RETURNING id`
		err = o.q.SelectContext(ctx, &reorgTxIDs, query, fromAddress, latestNonce, chainID.String())
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return pkgerrors.Wrap(err, "MarkConfirmedAndReorgedTransactions failed to update reorged txes")
		}

		// Confirmed transactions
		var dbIncludedEtxs []DbEthTx
		query = `UPDATE evm.txes SET state='confirmed' WHERE state='unconfirmed' AND from_address=$1 AND nonce<$2 AND evm_chain_id=$3 RETURNING *`
		err = o.q.SelectContext(ctx, &dbIncludedEtxs, query, fromAddress, latestNonce, chainID.String())
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return pkgerrors.Wrap(err, "MarkConfirmedAndReorgedTransactions failed to update confirmed txes")
		}

		for _, tx := range dbIncludedEtxs {
			txCopy := new(types.Transaction)
			tx.ToTransaction(txCopy)
			confirmedTxs = append(confirmedTxs, txCopy)
		}
		return nil
	})
	return
}

func (o *evmTxStore) MarkUnconfirmedTransactionPurgeable(ctx context.Context, id uint64, _ common.Address) error {
	var cancel context.CancelFunc
	ctx, cancel = o.stopCh.Ctx(ctx)
	defer cancel()

	// TODO: add is_purgeable column
	if _, err := o.q.ExecContext(ctx, `UPDATE evm.txes SET is_purgeable WHERE id=$1`, id); err != nil {
		return err
	}
	return nil
}

func (o *evmTxStore) UpdateTransactionBroadcast(ctx context.Context, txID uint64, _ uint64, attemptHash common.Hash, fromAddress common.Address) error {
	var cancel context.CancelFunc
	ctx, cancel = o.stopCh.Ctx(ctx)
	defer cancel()

	err := o.Transact(ctx, false, func(orm *evmTxStore) error {
		if _, qErr := o.q.ExecContext(ctx, `UPDATE evm.txes SET broadcast_at=NOW() WHERE id=$1`, txID); qErr != nil {
			return qErr
		}
		// TODO: Add broadcast_at column
		// if _, qErr := o.q.ExecContext(ctx, `UPDATE evm.tx_attempts SET broadcast_at=NOW() WHERE id=$1`, txID); qErr != nil {
		// 	return qErr
		// }
		return nil
	})
	return err
}

func (o *evmTxStore) UpdateUnstartedTransactionWithNonce(ctx context.Context, fromAddress common.Address, nonce uint64, chainID *big.Int) (*types.Transaction, error) {
	var cancel context.CancelFunc
	ctx, cancel = o.stopCh.Ctx(ctx)
	defer cancel()

	txCopy := new(types.Transaction)
	var dbEtx DbEthTx
	// We are forced to initialize broadcast_at and initial_broadcast_at to NOW because of a pre-existing constraint
	qErr := o.q.GetContext(ctx, &dbEtx, `UPDATE evm.txes SET state= 'unconfirmed', nonce=$1, broadcast_at=NOW(), initial_broadcast_at=NOW()
	WHERE id=(SELECT id FROM evm.txes WHERE from_address = $2 AND state = 'unstarted' AND evm_chain_id = $3 ORDER BY value ASC, created_at ASC, id ASC LIMIT 1) RETURNING *`, nonce, fromAddress, chainID.String())
	if qErr != nil {
		if errors.Is(qErr, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, pkgerrors.Wrap(qErr, "UpdateUnstartedTransactionWithNonce failed to update evm.txes")
	}
	dbEtx.ToTransaction(txCopy)
	return txCopy, nil
}

// FindTxWithIdempotency returns any broadcast ethtx with the given idempotencyKey and chainID
func (o *evmTxStore) FindTxWithIdempotency(ctx context.Context, idempotencyKey string, chainID *big.Int) (tx *types.Transaction, err error) {
	var cancel context.CancelFunc
	ctx, cancel = o.stopCh.Ctx(ctx)
	defer cancel()
	err = o.Transact(ctx, true, func(orm *evmTxStore) error {
		var dbEtx DbEthTx
		err = o.q.GetContext(ctx, &dbEtx, `SELECT * FROM evm.txes WHERE idempotency_key = $1 and evm_chain_id = $2`, idempotencyKey, chainID.String())
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil
			}
			return pkgerrors.Wrap(err, "FindTxWithIdempotency failed to load evm.txes")
		}
		tx = new(types.Transaction)
		dbEtx.ToTransaction(tx)
		if err = loadTxAttempts(ctx, orm, tx); err != nil {
			return pkgerrors.Wrap(err, "FindTxWithIdempotency failed to load evm.tx_attempts: %w")
		}
		return nil
	})
	return
}

func loadTxAttempts(ctx context.Context, o *evmTxStore, tx *types.Transaction) error {
	var dbTxAttempts []DbEthTxAttempt
	err := o.q.SelectContext(ctx, &dbTxAttempts, `SELECT * FROM evm.tx_attempts WHERE eth_tx_id=$1`, tx.ID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return pkgerrors.Wrap(err, "loadTxAttempts failed")
		}
	}
	for _, attempt := range dbTxAttempts {
		attemptCopy := new(types.Attempt)
		attempt.ToAttempt(attemptCopy)
		tx.Attempts = append(tx.Attempts, attemptCopy)
	}
	return nil
}

func pruneMaxQueue(ctx context.Context, queueSize uint64, o *evmTxStore) (ids []int64, err error) {
	err = o.q.SelectContext(ctx, &ids, `
DELETE FROM evm.txes
WHERE state = 'unstarted' AND
id < (
	SELECT min(id)
	FROM evm.txes
	WHERE state = 'unstarted'
	GROUP BY id
	ORDER BY id DESC
	LIMIT $1
) RETURNING id`, queueSize)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return
		}
		return ids, fmt.Errorf("pruneMaxQueue failed: %w", err)
	}
	return
}

// Add is required by the inMemoryStore
func (o *evmTxStore) Add(_ ...common.Address) error {
	return nil
}

func (db *DbEthTx) FromTransaction(tx *types.Transaction) {
	//nolint:gosec // IDs are generated in the DB and fit the corresponding DB struct
	db.ID = int64(tx.ID)
	db.IdempotencyKey = tx.IdempotencyKey
	db.FromAddress = tx.FromAddress
	db.ToAddress = tx.ToAddress
	db.EncodedPayload = tx.Data
	db.GasLimit = tx.SpecifiedGasLimit
	db.BroadcastAt = tx.LastBroadcastAt
	db.CreatedAt = tx.CreatedAt
	db.Meta = tx.Meta
	db.Subject = tx.Subject
	db.PipelineTaskRunID = tx.PipelineTaskRunID
	db.MinConfirmations = tx.MinConfirmations
	db.InitialBroadcastAt = tx.InitialBroadcastAt
	db.SignalCallback = tx.SignalCallback
	db.CallbackCompleted = tx.CallbackCompleted

	db.Value = assets.Eth(*tx.Value)
	db.State = tx.State

	if tx.ChainID != nil {
		db.EVMChainID = *ubig.New(tx.ChainID)
	}
	if tx.Nonce != nil {
		c := *tx.Nonce
		//nolint:gosec // suppress this until migration
		n := int64(c)
		db.Nonce = &n
	}
}

func (db DbEthTx) ToTransaction(tx *types.Transaction) {
	tx.ID = uint64(db.ID)
	if db.Nonce != nil {
		c := *db.Nonce
		n := uint64(c)
		tx.Nonce = &n
	}
	tx.IdempotencyKey = db.IdempotencyKey
	tx.FromAddress = db.FromAddress
	tx.ToAddress = db.ToAddress
	tx.Data = db.EncodedPayload
	tx.Value = db.Value.ToInt()
	tx.SpecifiedGasLimit = db.GasLimit

	tx.CreatedAt = db.CreatedAt
	tx.InitialBroadcastAt = db.BroadcastAt
	tx.LastBroadcastAt = db.BroadcastAt

	tx.State = db.State
	tx.Meta = db.Meta
	tx.Subject = db.Subject
	tx.PipelineTaskRunID = db.PipelineTaskRunID
	tx.MinConfirmations = db.MinConfirmations
	tx.ChainID = db.EVMChainID.ToInt()
	tx.SignalCallback = db.SignalCallback
	tx.CallbackCompleted = db.CallbackCompleted
}

func (db *DbEthTxAttempt) FromAttempt(attempt *types.Attempt) {
	//nolint:gosec // IDs are generated in the DB and fit the corresponding DB struct
	db.ID = int64(attempt.ID)
	//nolint:gosec // IDs are generated in the DB and fit the corresponding DB struct
	db.EthTxID = int64(attempt.TxID)
	db.GasPrice = attempt.Fee.GasPrice
	db.SignedRawTx = attempt.SignedTransaction.Data()
	db.Hash = attempt.Hash
	db.CreatedAt = attempt.CreatedAt
	db.ChainSpecificGasLimit = attempt.GasLimit
	db.GasTipCap = attempt.Fee.GasTipCap
	db.GasFeeCap = attempt.Fee.GasFeeCap
	db.TxType = int(attempt.Type)
	// This is necessary to keep attempts backwards compatible
	db.State = txmgrtypes.TxAttemptBroadcast.String()
}

func (db *DbEthTxAttempt) ToAttempt(attempt *types.Attempt) {
	attempt.ID = uint64(db.ID)
	attempt.TxID = uint64(db.EthTxID)
	attempt.Hash = db.Hash
	attempt.Fee = gas.EvmFee{
		GasPrice:   db.GasPrice,
		DynamicFee: gas.DynamicFee{GasTipCap: db.GasTipCap, GasFeeCap: db.GasFeeCap},
	}
	attempt.GasLimit = db.ChainSpecificGasLimit
	attempt.Type = byte(db.TxType)
	// attempt.SignedTransaction = db.SignedRawTx
	attempt.CreatedAt = db.CreatedAt
}
