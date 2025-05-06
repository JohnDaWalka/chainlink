package mocks

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	fee_manager "github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/fee_manager_v0_5_0"
)

// MockFeeManager implements FeeManagerInterface for testing
type MockFeeManager struct {
	MockAddress common.Address

	OwnerFunc                   func(*bind.CallOpts) (common.Address, error)
	ILinkAddressFunc            func(*bind.CallOpts) (common.Address, error)
	INativeAddressFunc          func(*bind.CallOpts) (common.Address, error)
	TypeAndVersionFunc          func(*bind.CallOpts) (string, error)
	SNativeSurchargeFunc        func(*bind.CallOpts) (*big.Int, error)
	LinkAvailableForPaymentFunc func(*bind.CallOpts) (*big.Int, error)

	// Discount Fields
	SGlobalDiscountsFunc     func(*bind.CallOpts, common.Address, common.Address) (*big.Int, error)
	SSubscriberDiscountsFunc func(*bind.CallOpts, common.Address, [32]byte, common.Address) (*big.Int, error)

	// Events
	FilterSubscriberDiscountUpdatedFunc func(*bind.FilterOpts, []common.Address, [][32]byte) (*fee_manager.FeeManagerSubscriberDiscountUpdatedIterator, error)
}

// Implementation of key methods

func (m *MockFeeManager) Address() common.Address {
	return m.MockAddress
}

func (m *MockFeeManager) Owner(opts *bind.CallOpts) (common.Address, error) {
	return m.OwnerFunc(opts)
}

func (m *MockFeeManager) ILinkAddress(opts *bind.CallOpts) (common.Address, error) {
	return m.ILinkAddressFunc(opts)
}

func (m *MockFeeManager) INativeAddress(opts *bind.CallOpts) (common.Address, error) {
	return m.INativeAddressFunc(opts)
}

func (m *MockFeeManager) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	return m.TypeAndVersionFunc(opts)
}

func (m *MockFeeManager) SNativeSurcharge(opts *bind.CallOpts) (*big.Int, error) {
	return m.SNativeSurchargeFunc(opts)
}

func (m *MockFeeManager) LinkAvailableForPayment(opts *bind.CallOpts) (*big.Int, error) {
	return m.LinkAvailableForPaymentFunc(opts)
}

func (m *MockFeeManager) SGlobalDiscounts(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return m.SGlobalDiscountsFunc(opts, arg0, arg1)
}

func (m *MockFeeManager) SSubscriberDiscounts(opts *bind.CallOpts, arg0 common.Address, arg1 [32]byte, arg2 common.Address) (*big.Int, error) {
	return m.SSubscriberDiscountsFunc(opts, arg0, arg1, arg2)
}

func (m *MockFeeManager) FilterSubscriberDiscountUpdated(opts *bind.FilterOpts, subscriber []common.Address, feedId [][32]byte) (*fee_manager.FeeManagerSubscriberDiscountUpdatedIterator, error) {
	return m.FilterSubscriberDiscountUpdatedFunc(opts, subscriber, feedId)
}

// Add other method implementations as needed...

// Mock iterator
type MockSubscriberDiscountUpdatedIterator struct {
	Events []*fee_manager.FeeManagerSubscriberDiscountUpdated
	Index  int
	Err    error
}

func (it *MockSubscriberDiscountUpdatedIterator) Next() bool {
	if it.Index >= len(it.Events) {
		return false
	}
	it.Index++
	return true
}

func (it *MockSubscriberDiscountUpdatedIterator) Error() error {
	return it.Err
}

func (it *MockSubscriberDiscountUpdatedIterator) Close() error {
	return nil
}

func (it *MockSubscriberDiscountUpdatedIterator) Event() *fee_manager.FeeManagerSubscriberDiscountUpdated {
	if it.Index <= 0 || it.Index > len(it.Events) {
		return nil
	}
	return it.Events[it.Index-1]
}
