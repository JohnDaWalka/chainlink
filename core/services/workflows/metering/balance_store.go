package metering

import (
	"errors"
	"sync"

	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrInvalidAmount       = errors.New("amount must be greater than 0")
)

type balanceStore struct {
	// meteringMode determines whether negative balances should return an error. meteringMode == true allows negative
	// balances.
	meteringMode bool
	// A balance of credits
	balance decimal.Decimal
	// Conversion rates of resource dimensions to number of units per credit
	conversions map[string]decimal.Decimal // TODO flip this
	lggr        logger.Logger
	mu          sync.RWMutex
}

type BalanceStore interface {
	Get() (balance int64)
	GetAs(unit string) (balance int64)
	Minus(amount int64) error
	MinusAs(unit string, amount int64) error
	Add(amount int64) error
	AddAs(unit string, amount int64) error
	AllowNegative()
}

var _ BalanceStore = (BalanceStore)(nil)

func NewBalanceStore(
	startingBalance decimal.Decimal,
	conversions map[string]decimal.Decimal,
	lggr logger.Logger,
) *balanceStore {
	// validations
	for resource, rate := range conversions {
		if rate.IsNegative() {
			// fail open
			lggr.Errorw("conversion rates must be a positive number, not using conversion", "resource", resource, "rate", rate)
			delete(conversions, resource)
		}
	}

	return &balanceStore{
		meteringMode: false,
		balance:      startingBalance,
		conversions:  conversions,
		lggr:         lggr,
	}
}

// convertToBalance converts a resource dimension amount to a credit amount.
// This method should only be used under a read lock.
func (bs *balanceStore) convertToBalance(fromUnit string, amount decimal.Decimal) decimal.Decimal {
	rate, ok := bs.conversions[fromUnit]
	if !ok {
		// Fail open, continue optimistically
		bs.lggr.Errorw("could not find conversion rate, continuing as 1:1; entering metering mode", "unit", fromUnit)
		rate = decimal.NewFromInt(1)
		bs.meteringMode = true
	}

	return amount.Mul(rate)
}

// ConvertToBalance converts a resource dimensions amount to a credit amount.
func (bs *balanceStore) ConvertToBalance(fromUnit string, amount decimal.Decimal) decimal.Decimal {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	return bs.convertToBalance(fromUnit, amount)
}

// convertFromBalance converts a credit amount to a resource dimensions amount.
// This method should only be used under a read lock.
func (bs *balanceStore) convertFromBalance(toUnit string, amount decimal.Decimal) decimal.Decimal {
	rate, ok := bs.conversions[toUnit]
	if !ok {
		// Fail open, continue optimistically
		bs.lggr.Errorw("could not find conversion rate, continuing as 1:1; entering metering mode", "unit", toUnit)
		rate = decimal.NewFromInt(1)
		bs.meteringMode = true
	}

	decimal.DivisionPrecision = defaultDecimalPrecision

	return amount.Div(rate)
}

// ConvertFromBalance converts a credit amount to a resource dimensions amount.
func (bs *balanceStore) ConvertFromBalance(toUnit string, amount decimal.Decimal) decimal.Decimal {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	return bs.convertFromBalance(toUnit, amount)
}

// Get returns the current credit balance
func (bs *balanceStore) Get() decimal.Decimal {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	return bs.balance
}

// GetAs returns the current universal credit balance expressed as a resource dimensions.
func (bs *balanceStore) GetAs(unit string) decimal.Decimal {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	if bs.balance.LessThanOrEqual(decimal.Zero) {
		return decimal.Zero
	}

	return bs.convertFromBalance(unit, bs.balance)
}

// Minus lowers the current credit balance.
func (bs *balanceStore) Minus(amount decimal.Decimal) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if amount.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidAmount
	}

	if amount.GreaterThan(bs.balance) && !bs.meteringMode {
		return ErrInsufficientBalance
	}

	bs.balance = bs.balance.Sub(amount)

	return nil
}

// MinusAs lowers the current credit balance based on an amount of resource dimensions.
func (bs *balanceStore) MinusAs(unit string, amount decimal.Decimal) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if amount.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidAmount
	}

	balToMinus := bs.convertToBalance(unit, amount)

	if balToMinus.GreaterThan(bs.balance) && !bs.meteringMode {
		return ErrInsufficientBalance
	}

	bs.balance = bs.balance.Sub(balToMinus)

	return nil
}

// Add increases the current credit balance.
func (bs *balanceStore) Add(amount decimal.Decimal) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if amount.LessThan(decimal.Zero) {
		return ErrInvalidAmount
	}

	bs.balance = bs.balance.Add(amount)

	return nil
}

// AddAs increases the current credit balance based on an amount of resource dimensions.
func (bs *balanceStore) AddAs(unit string, amount decimal.Decimal) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if amount.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidAmount
	}

	bs.balance = bs.balance.Add(bs.convertToBalance(unit, amount))

	return nil
}

// AllowNegative turns on the flag to allow negative balances.
func (bs *balanceStore) AllowNegative() {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	bs.meteringMode = true
}
