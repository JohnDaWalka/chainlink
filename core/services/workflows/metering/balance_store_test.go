package metering

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

func TestBalanceStore(t *testing.T) {
	t.Parallel()

	zero := decimal.NewFromInt(0)
	one := decimal.NewFromInt(1)
	two := decimal.NewFromInt(2)
	five := decimal.NewFromInt(5)
	seven := decimal.NewFromInt(7)
	nine := decimal.NewFromInt(9)
	ten := decimal.NewFromInt(10)
	eleven := decimal.NewFromInt(11)

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		// 1 of resourceA is worth 2 credits
		// 2 credits is worth 1 of resourceA
		rate := decimal.NewFromInt(2)
		balanceStore := NewBalanceStore(ten, map[string]decimal.Decimal{"resourceA": rate}, logger.Test(t))

		assert.True(t, balanceStore.Get().Equal(ten), "initialization should set balance")
		assert.True(t, balanceStore.GetAs("resourceA").Equal(five), "rate should apply to balance")

		require.NoError(t, balanceStore.Add(one))
		assert.True(t, balanceStore.Get().Equal(eleven), "addition should update the balance")

		require.NoError(t, balanceStore.Minus(two))
		require.ErrorIs(t, balanceStore.Minus(decimal.NewFromInt(0)), ErrInvalidAmount)
		require.ErrorIs(t, balanceStore.Minus(decimal.NewFromInt(-1)), ErrInvalidAmount)
		assert.True(t, balanceStore.Get().Equal(nine), "subtraction should update the balance")

		require.NoError(t, balanceStore.AddAs("resourceA", one))
		require.ErrorIs(t, balanceStore.AddAs("resourceA", decimal.NewFromInt(0)), ErrInvalidAmount)
		require.ErrorIs(t, balanceStore.AddAs("resourceA", decimal.NewFromInt(-1)), ErrInvalidAmount)
		assert.True(t, balanceStore.Get().Equal(eleven), "addition by rate should update balance")

		require.NoError(t, balanceStore.MinusAs("resourceA", two))
		require.ErrorIs(t, balanceStore.MinusAs("resourceA", decimal.NewFromInt(0)), ErrInvalidAmount)
		require.ErrorIs(t, balanceStore.MinusAs("resourceA", decimal.NewFromInt(-1)), ErrInvalidAmount)
		assert.True(t, balanceStore.Get().Equal(seven), "subtraction by rate should update balance")
	})

	t.Run("handles unknown resources as 1:1", func(t *testing.T) {
		t.Parallel()

		lggr, logs := logger.TestObserved(t, zapcore.ErrorLevel)
		balanceStore := NewBalanceStore(ten, map[string]decimal.Decimal{}, lggr)

		assert.True(t, balanceStore.GetAs("").Equal(ten))
		require.NoError(t, balanceStore.MinusAs("", one))
		assert.True(t, balanceStore.Get().Equal(nine))
		assert.Len(t, logs.All(), 2)
	})

	t.Run("throws out negative conversion rates", func(t *testing.T) {
		t.Parallel()

		lggr, logs := logger.TestObserved(t, zapcore.ErrorLevel)
		balanceStore := NewBalanceStore(ten, map[string]decimal.Decimal{"resourceA": decimal.NewFromInt(-1)}, lggr)

		assert.True(t, balanceStore.GetAs("resourceA").Equal(ten))

		require.NoError(t, balanceStore.MinusAs("resourceA", one))
		assert.True(t, balanceStore.Get().Equal(nine))
		assert.Len(t, logs.All(), 3)
	})

	t.Run("cannot go negative by default", func(t *testing.T) {
		t.Parallel()

		balanceStore := NewBalanceStore(decimal.Zero, map[string]decimal.Decimal{"resourceA": decimal.NewFromInt(1)}, logger.Nop())

		require.ErrorIs(t, balanceStore.Minus(one), ErrInsufficientBalance)
		require.ErrorIs(t, balanceStore.MinusAs("resourceA", one), ErrInsufficientBalance)
	})

	t.Run("can go negative after allowNegative", func(t *testing.T) {
		t.Parallel()

		balanceStore := NewBalanceStore(decimal.Zero, map[string]decimal.Decimal{"resourceA": decimal.NewFromInt(1)}, logger.Nop())

		balanceStore.AllowNegative()
		require.NoError(t, balanceStore.Minus(one))
		assert.True(t, balanceStore.Get().Equal(decimal.NewFromInt(-1)))
	})

	t.Run("returns negative balances as 0 when converted to a resource", func(t *testing.T) {
		t.Parallel()

		balanceStore := NewBalanceStore(decimal.Zero, map[string]decimal.Decimal{"resourceA": ten}, logger.Nop())

		balanceStore.AllowNegative()
		require.NoError(t, balanceStore.Minus(one))
		assert.True(t, balanceStore.Get().Equal(decimal.NewFromInt(-1)))
		assert.True(t, balanceStore.GetAs("resourceA").Equal(zero))
	})

	t.Run("handles decimal rates", func(t *testing.T) {
		t.Parallel()

		// 1 of resource A is worth 0.1 credits
		rate := decimal.NewFromFloat(0.1)
		balanceStore := NewBalanceStore(ten, map[string]decimal.Decimal{"resourceA": rate}, logger.Nop())

		assert.True(t, balanceStore.Get().Equal(ten))
		assert.True(t, balanceStore.GetAs("resourceA").Equal(decimal.NewFromInt(100)))
	})

	t.Run("applies no rounding to result", func(t *testing.T) {
		t.Parallel()

		// 1 of resource A is worth 0.2 credits
		rate := decimal.NewFromFloat(0.2)
		balanceStore := NewBalanceStore(two, map[string]decimal.Decimal{"resourceA": rate}, logger.Nop())

		assert.True(t, balanceStore.Get().Equal(two))
		require.NoError(t, balanceStore.MinusAs("resourceA", one))
		assert.True(t, balanceStore.Get().Equal(decimal.NewFromFloat(1.8)), balanceStore.Get())
	})
}
