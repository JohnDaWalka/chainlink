package v0_5

//
//import (
//	"context"
//	"encoding/hex"
//	"math/big"
//	"testing"
//
//	"github.com/ethereum/go-ethereum/accounts/abi/bind"
//	"github.com/ethereum/go-ethereum/common"
//	fee_manager "github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/fee_manager_v0_5_0"
//	"github.com/smartcontractkit/chainlink/deployment/data-streams/view/v0_5/mocks"
//	"github.com/stretchr/testify/assert"
//)
//
//func TestFeeManagerViewGenerator_Generate(t *testing.T) {
//	// Create mock contract
//	mockContract := &mocks.MockFeeManager{
//		MockAddress: common.HexToAddress("0x1234567890123456789012345678901234567890"),
//
//		// Setup the mock function implementations
//		OwnerFunc: func(opts *bind.CallOpts) (common.Address, error) {
//			return common.HexToAddress("0xOwnerAddress"), nil
//		},
//
//		ILinkAddressFunc: func(opts *bind.CallOpts) (common.Address, error) {
//			return common.HexToAddress("0xLinkAddress"), nil
//		},
//
//		INativeAddressFunc: func(opts *bind.CallOpts) (common.Address, error) {
//			return common.HexToAddress("0xNativeAddress"), nil
//		},
//
//		TypeAndVersionFunc: func(opts *bind.CallOpts) (string, error) {
//			return "FeeManager 2.1.0", nil
//		},
//
//		SNativeSurchargeFunc: func(opts *bind.CallOpts) (*big.Int, error) {
//			return big.NewInt(100000000), nil
//		},
//
//		LinkAvailableForPaymentFunc: func(opts *bind.CallOpts) (*big.Int, error) {
//			return big.NewInt(5000000000), nil
//		},
//
//		// Mock discount functions
//		SGlobalDiscountsFunc: func(opts *bind.CallOpts, subscriber common.Address, token common.Address) (*big.Int, error) {
//			// Return a discount for a specific subscriber/token combination
//			if subscriber == common.HexToAddress("0xSubscriberAddress") &&
//				token == common.HexToAddress("0xLinkAddress") {
//				return big.NewInt(1000000000000000000), nil // 1 LINK (10^18)
//			}
//			return big.NewInt(0), nil
//		},
//
//		SSubscriberDiscountsFunc: func(opts *bind.CallOpts, subscriber common.Address, feedId [32]byte, token common.Address) (*big.Int, error) {
//			// Return a specific discount for certain combinations
//			if subscriber == common.HexToAddress("0xSubscriberAddress") {
//				// Example feedId for testing
//				testFeedId := [32]byte{}
//				copy(testFeedId[:], []byte("testFeedId"))
//
//				if string(feedId[:]) == string(testFeedId[:]) &&
//					token == common.HexToAddress("0xNativeAddress") {
//					return big.NewInt(500000000000000000), nil // 0.5 Native token
//				}
//			}
//			return big.NewInt(0), nil
//		},
//
//		// Mock event filtering
//		FilterSubscriberDiscountUpdatedFunc: func(opts *bind.FilterOpts, subscriber []common.Address, feedId [][32]byte) (*contracts.FeeManagerSubscriberDiscountUpdatedIterator, error) {
//			// Setup mock event data
//			zeroFeedId := [32]byte{}
//			testFeedId := [32]byte{}
//			copy(testFeedId[:], []byte("testFeedId"))
//
//			// Create mock events
//			events := []*fee_manager.FeeManagerSubscriberDiscountUpdated{
//				{
//					Subscriber: common.HexToAddress("0xSubscriberAddress"),
//					FeedId:     zeroFeedId, // Global discount
//					Token:      common.HexToAddress("0xLinkAddress"),
//					Discount:   uint64(1000000000000000000), // 1 LINK
//				},
//				{
//					Subscriber: common.HexToAddress("0xSubscriberAddress"),
//					FeedId:     testFeedId,
//					Token:      common.HexToAddress("0xNativeAddress"),
//					Discount:   uint64(500000000000000000), // 0.5 Native
//				},
//			}
//
//			return &mocks.MockSubscriberDiscountUpdatedIterator{
//				Events: events,
//				Index:  0,
//				Err:    nil,
//			}, nil
//		},
//	}
//
//	// Create the generator and parameters
//	generator := NewFeeManagerViewGenerator()
//	params := FeeManagerContext{FromBlock: 0, ToBlock: nil}
//
//	// Call the generate function
//	ctx := context.Background()
//	view, err := generator.Generate(ctx, mockContract, params)
//
//	// Assert results
//	assert.NoError(t, err)
//	assert.NotNil(t, view)
//
//	// Assert basic contract info
//	assert.Equal(t, "0x1234567890123456789012345678901234567890", view.Address)
//	assert.Equal(t, "0xOwnerAddress", view.Owner)
//	assert.Equal(t, "FeeManager 2.1.0", view.TypeAndVersion)
//
//	// Assert discount information
//	subscriberAddr := "0xSubscriberAddress"
//	assert.Contains(t, view.Discounts, subscriberAddr)
//
//	// Check global discount
//	assert.Contains(t, view.Discounts[subscriberAddr], "global")
//	assert.Equal(t, "1000000000000000000", view.Discounts[subscriberAddr]["global"].Link)
//	assert.True(t, view.Discounts[subscriberAddr]["global"].IsGlobal)
//
//	// Check feed-specific discount
//	feedIdHex := hex.EncodeToString([]byte("testFeedId"))
//	assert.Contains(t, view.Discounts[subscriberAddr], feedIdHex)
//	assert.Equal(t, "500000000000000000", view.Discounts[subscriberAddr][feedIdHex].Native)
//	assert.False(t, view.Discounts[subscriberAddr][feedIdHex].IsGlobal)
//}
