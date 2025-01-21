package txmgr_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

const maxQueued = 250

func TestTxmStore_BroadcastLifecycle(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	txStore := txmgr.NewTxStore(db, logger.Test(t))
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	ctx := tests.Context(t)
	_, from := cltest.MustInsertRandomKey(t, ethKeyStore)
	_, to := cltest.MustInsertRandomKey(t, ethKeyStore)
	var nonce uint64

	// Create
	r := NewEthTxRequest(from)
	_, err := txStore.CreateTx(ctx, r, maxQueued)
	require.NoError(t, err)

	// Count unconfirmed
	_, unconfirmedCount, err := txStore.FetchUnconfirmedTransactionAtNonceWithCount(ctx, 0, from, testutils.FixtureChainID)
	require.NoError(t, err)
	assert.Equal(t, 0, unconfirmedCount)

	// No matching transaction
	emptyTx, err := txStore.UpdateUnstartedTransactionWithNonce(ctx, to, nonce, testutils.FixtureChainID)
	require.NoError(t, err)
	require.Nil(t, emptyTx)
	emptyTx, err = txStore.UpdateUnstartedTransactionWithNonce(ctx, from, nonce, big.NewInt(99))
	require.NoError(t, err)
	require.Nil(t, emptyTx)
	// Start transaction
	tx, err := txStore.UpdateUnstartedTransactionWithNonce(ctx, from, nonce, testutils.FixtureChainID)
	require.NoError(t, err)
	require.Len(t, tx.Attempts, 0)
	require.NotNil(t, tx)

	// Count unconfirmed
	tx, unconfirmedCount, err = txStore.FetchUnconfirmedTransactionAtNonceWithCount(ctx, 0, from, testutils.FixtureChainID)
	require.NoError(t, err)
	assert.Equal(t, 1, unconfirmedCount)

	// Append attempt to transaction
	attempt := NewEthTxAttempt(tx)
	err = txStore.AppendAttemptToTransaction(ctx, tx, attempt)
	require.NoError(t, err)
	tx, unconfirmedCount, err = txStore.FetchUnconfirmedTransactionAtNonceWithCount(ctx, 0, from, testutils.FixtureChainID)
	require.NoError(t, err)
	require.Len(t, tx.Attempts, 1)
	assert.Equal(t, 1, unconfirmedCount)

	// Update transaction and attempt broadcast timestamps
	err = txStore.UpdateTransactionBroadcast(ctx, attempt.TxID, nonce, attempt.Hash, from)
	require.NoError(t, err)
}

func TestTxmStore_BackfillLifecycle(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	txStore := txmgr.NewTxStore(db, logger.Test(t))
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	ctx := tests.Context(t)
	_, from := cltest.MustInsertRandomKey(t, ethKeyStore)
	var nonce uint64

	// Create
	r := NewEthTxRequest(from)
	_, err := txStore.CreateTx(ctx, r, maxQueued)
	require.NoError(t, err)

	// Start transaction
	tx, err := txStore.UpdateUnstartedTransactionWithNonce(ctx, from, nonce, testutils.FixtureChainID)
	require.NoError(t, err)
	require.Len(t, tx.Attempts, 0)
	require.NotNil(t, tx)

	// Append attempt to transaction
	attempt := NewEthTxAttempt(tx)
	err = txStore.AppendAttemptToTransaction(ctx, tx, attempt)
	require.NoError(t, err)

	// Confirm transaction
	confirmedTransactions, unconfirmedTransactionIDs, err := txStore.MarkConfirmedAndReorgedTransactions(ctx, nonce+1, from, testutils.FixtureChainID)
	require.NoError(t, err)
	require.Len(t, confirmedTransactions, 1)
	require.Len(t, unconfirmedTransactionIDs, 0)

	// Unconfirm transaction
	confirmedTransactions, unconfirmedTransactionIDs, err = txStore.MarkConfirmedAndReorgedTransactions(ctx, nonce, from, testutils.FixtureChainID)
	require.NoError(t, err)
	require.Len(t, confirmedTransactions, 0)
	require.Len(t, unconfirmedTransactionIDs, 1)

	// Create empty
	tx, err = txStore.CreateEmptyUnconfirmedTransaction(ctx, from, 1, 1, testutils.FixtureChainID)
	require.NoError(t, err)
	require.NotNil(t, tx)

	// Count unconfirmed
	tx, unconfirmedCount, err := txStore.FetchUnconfirmedTransactionAtNonceWithCount(ctx, 0, from, testutils.FixtureChainID)
	require.NoError(t, err)
	require.NotNil(t, tx)
	assert.Equal(t, 2, unconfirmedCount)
}

func NewEthTxRequest(fromAddress common.Address) *types.TxRequest {
	return &types.TxRequest{
		FromAddress:       fromAddress,
		ToAddress:         testutils.NewAddress(),
		Data:              []byte{1, 2, 3},
		Value:             big.NewInt(10),
		SpecifiedGasLimit: uint64(1000000000),
		ChainID:           testutils.FixtureChainID,
	}
}

func NewEthTxAttempt(tx *types.Transaction) *types.Attempt {
	legacyTx := evmtypes.LegacyTx{
		Nonce:    *tx.Nonce,
		Value:    tx.Value,
		Gas:      tx.SpecifiedGasLimit,
		GasPrice: big.NewInt(1),
		Data:     tx.Data,
	}
	signedTransaction := evmtypes.NewTx(&legacyTx)
	return &types.Attempt{
		TxID:              tx.ID,
		Hash:              testutils.NewHash(),
		Fee:               gas.EvmFee{GasPrice: assets.NewWeiI(1)},
		GasLimit:          100,
		Type:              0x0,
		SignedTransaction: signedTransaction,
	}
}
