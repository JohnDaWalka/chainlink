package txm

import (
	"errors"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txm/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txm/storage"
)

func TestLifecycle(t *testing.T) {
	t.Parallel()

	client := mocks.NewClient(t)
	ab := mocks.NewAttemptBuilder(t)
	config := Config{BlockTime: 10 * time.Millisecond}
	address1 := testutils.NewAddress()
	address2 := testutils.NewAddress()
	assert.NotEqual(t, address1, address2)
	addresses := []common.Address{address1, address2}
	keystore := mocks.NewKeystore(t)
	keystore.On("EnabledAddressesForChain", mock.Anything, mock.Anything).Return(addresses, nil)

	t.Run("fails to start if initial pending nonce call fails", func(t *testing.T) {
		txm := NewTxm(logger.Test(t), testutils.FixtureChainID, client, ab, nil, nil, config, keystore)
		client.On("PendingNonceAt", mock.Anything, address1).Return(uint64(0), errors.New("error")).Once()
		assert.Error(t, txm.Start(tests.Context(t)))
	})

	t.Run("tests lifecycle successfully without any transactions", func(t *testing.T) {
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		txStore := storage.NewInMemoryStoreManager(lggr, testutils.FixtureChainID)
		assert.NoError(t, txStore.Add(addresses...))
		txm := NewTxm(lggr, testutils.FixtureChainID, client, ab, txStore, nil, config, keystore)
		var nonce uint64 = 0
		// Start
		client.On("PendingNonceAt", mock.Anything, address1).Return(nonce, nil).Once()
		client.On("PendingNonceAt", mock.Anything, address2).Return(nonce, nil).Once()
		// backfill loop (may or may not be executed multiple times)
		client.On("NonceAt", mock.Anything, address1, mock.Anything).Return(nonce, nil)
		client.On("NonceAt", mock.Anything, address2, mock.Anything).Return(nonce, nil)

		servicetest.Run(t, txm)
		tests.AssertLogEventually(t, observedLogs, "Backfill time elapsed")
	})

}

func TestTrigger(t *testing.T) {
	t.Parallel()

	address := testutils.NewAddress()
	keystore := mocks.NewKeystore(t)
	keystore.On("EnabledAddressesForChain", mock.Anything, mock.Anything).Return([]common.Address{address}, nil)
	t.Run("Trigger fails if Txm is unstarted", func(t *testing.T) {
		txm := NewTxm(logger.Test(t), nil, nil, nil, nil, nil, Config{}, keystore)
		txm.Trigger(address)
		assert.Error(t, txm.Trigger(address), "Txm unstarted")
	})

	t.Run("executes Trigger", func(t *testing.T) {
		lggr := logger.Test(t)
		txStore := storage.NewInMemoryStoreManager(lggr, testutils.FixtureChainID)
		assert.NoError(t, txStore.Add(address))
		client := mocks.NewClient(t)
		ab := mocks.NewAttemptBuilder(t)
		config := Config{BlockTime: 1 * time.Minute, RetryBlockThreshold: 10}
		txm := NewTxm(lggr, testutils.FixtureChainID, client, ab, txStore, nil, config, keystore)
		var nonce uint64 = 0
		// Start
		client.On("PendingNonceAt", mock.Anything, address).Return(nonce, nil).Once()
		servicetest.Run(t, txm)
		assert.NoError(t, txm.Trigger(address))
	})
}

func TestBroadcastTransaction(t *testing.T) {
	t.Parallel()

	ctx := tests.Context(t)
	client := mocks.NewClient(t)
	ab := mocks.NewAttemptBuilder(t)
	config := Config{}
	address := testutils.NewAddress()
	keystore := mocks.NewKeystore(t)

	t.Run("fails if FetchUnconfirmedTransactionAtNonceWithCount for unconfirmed transactions fails", func(t *testing.T) {
		mTxStore := mocks.NewTxStore(t)
		mTxStore.On("FetchUnconfirmedTransactionAtNonceWithCount", mock.Anything, mock.Anything, mock.Anything).Return(nil, 0, errors.New("call failed")).Once()
		txm := NewTxm(logger.Test(t), testutils.FixtureChainID, client, ab, mTxStore, nil, config, keystore)
		bo, err := txm.broadcastTransaction(ctx, address)
		assert.Error(t, err)
		assert.False(t, bo)
		assert.Contains(t, err.Error(), "call failed")
	})

	t.Run("throws a warning and returns if unconfirmed transactions exceed maxInFlightTransactions", func(t *testing.T) {
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		mTxStore := mocks.NewTxStore(t)
		mTxStore.On("FetchUnconfirmedTransactionAtNonceWithCount", mock.Anything, mock.Anything, mock.Anything).Return(nil, int(maxInFlightTransactions+1), nil).Once()
		txm := NewTxm(lggr, testutils.FixtureChainID, client, ab, mTxStore, nil, config, keystore)
		txm.broadcastTransaction(ctx, address)
		tests.AssertLogEventually(t, observedLogs, "Reached transaction limit")
	})

	t.Run("checks pending nonce if unconfirmed transactions are more than 1/3 of maxInFlightTransactions", func(t *testing.T) {
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		mTxStore := mocks.NewTxStore(t)
		txm := NewTxm(lggr, testutils.FixtureChainID, client, ab, mTxStore, nil, config, keystore)
		txm.setNonce(address, 1)
		mTxStore.On("FetchUnconfirmedTransactionAtNonceWithCount", mock.Anything, mock.Anything, mock.Anything).Return(nil, int(maxInFlightTransactions/3), nil).Twice()

		client.On("PendingNonceAt", mock.Anything, address).Return(uint64(0), nil).Once() // LocalNonce: 1, PendingNonce: 0
		txm.broadcastTransaction(ctx, address)

		client.On("PendingNonceAt", mock.Anything, address).Return(uint64(1), nil).Once() // LocalNonce: 1, PendingNonce: 1
		mTxStore.On("UpdateUnstartedTransactionWithNonce", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Once()
		txm.broadcastTransaction(ctx, address)
		tests.AssertLogCountEventually(t, observedLogs, "Reached transaction limit.", 1)

	})

	t.Run("fails if UpdateUnstartedTransactionWithNonce fails", func(t *testing.T) {
		mTxStore := mocks.NewTxStore(t)
		mTxStore.On("FetchUnconfirmedTransactionAtNonceWithCount", mock.Anything, mock.Anything, mock.Anything).Return(nil, 0, nil).Once()
		txm := NewTxm(logger.Test(t), testutils.FixtureChainID, client, ab, mTxStore, nil, config, keystore)
		mTxStore.On("UpdateUnstartedTransactionWithNonce", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("call failed")).Once()
		bo, err := txm.broadcastTransaction(ctx, address)
		assert.False(t, bo)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "call failed")
	})

	t.Run("returns if there are no unstarted transactions", func(t *testing.T) {
		lggr := logger.Test(t)
		txStore := storage.NewInMemoryStoreManager(lggr, testutils.FixtureChainID)
		assert.NoError(t, txStore.Add(address))
		txm := NewTxm(lggr, testutils.FixtureChainID, client, ab, txStore, nil, config, keystore)
		bo, err := txm.broadcastTransaction(ctx, address)
		assert.NoError(t, err)
		assert.False(t, bo)
		assert.Equal(t, uint64(0), txm.getNonce(address))
	})
}

func TestBackfillTransactions(t *testing.T) {
	t.Parallel()

	ctx := tests.Context(t)
	client := mocks.NewClient(t)
	ab := mocks.NewAttemptBuilder(t)
	storage := mocks.NewTxStore(t)
	config := Config{}
	address := testutils.NewAddress()
	keystore := mocks.NewKeystore(t)

	t.Run("fails if latest nonce fetching fails", func(t *testing.T) {
		txm := NewTxm(logger.Test(t), testutils.FixtureChainID, client, ab, storage, nil, config, keystore)
		client.On("NonceAt", mock.Anything, address, mock.Anything).Return(uint64(0), errors.New("latest nonce fail")).Once()
		bo, err := txm.backfillTransactions(ctx, address)
		assert.Error(t, err)
		assert.False(t, bo)
		assert.Contains(t, err.Error(), "latest nonce fail")
	})

	t.Run("fails if MarkTransactionsConfirmed fails", func(t *testing.T) {
		txm := NewTxm(logger.Test(t), testutils.FixtureChainID, client, ab, storage, nil, config, keystore)
		client.On("NonceAt", mock.Anything, address, mock.Anything).Return(uint64(0), nil)
		storage.On("MarkTransactionsConfirmed", mock.Anything, mock.Anything, address).Return([]uint64{}, []uint64{}, errors.New("marking transactions confirmed failed"))
		bo, err := txm.backfillTransactions(ctx, address)
		assert.Error(t, err)
		assert.False(t, bo)
		assert.Contains(t, err.Error(), "marking transactions confirmed failed")
	})
}
