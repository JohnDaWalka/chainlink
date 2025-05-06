package v0_5

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	fee_manager "github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/fee_manager_v0_5_0"
	dsutil "github.com/smartcontractkit/chainlink/deployment/data-streams/utils"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/view/interfaces"
)

// FeeManagerView represents a view of the FeeManager contract state
type FeeManagerView struct {
	LinkAddress         string                               `json:"linkAddress"`
	NativeAddress       string                               `json:"nativeAddress"`
	ProxyAddress        string                               `json:"proxyAddress"`
	RewardManager       string                               `json:"rewardManager"`
	NativeSurcharge     string                               `json:"nativeSurcharge"`
	LinkAvailable       string                               `json:"linkAvailable"`
	TypeAndVersion      string                               `json:"typeAndVersion,omitempty"`
	Owner               common.Address                       `json:"owner,omitempty"`
	SubscriberDiscounts map[string]map[string]TokenDiscounts `json:"subscriberDiscounts"` // Map[subscriberAddress][feedId]TokenDiscounts
}

type TokenDiscounts struct {
	Link     string `json:"link"`
	Native   string `json:"native"`
	IsGlobal bool   `json:"isGlobal"`
}

// FeeManagerView implements the ContractView interface
var _ interfaces.ContractView = (*FeeManagerView)(nil)

// SerializeView serializes view to JSON
func (v FeeManagerView) SerializeView() (string, error) {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal contract view: %w", err)
	}
	return string(bytes), nil
}

// FeeManagerContext represents parameters for generating a FeeManager view
// In this simple case, we might not need parameters, but including as an example
type FeeManagerContext struct {
	FromBlock uint64
	ToBlock   *uint64
}

// FeeManagerViewGenerator implements ContractViewGenerator for FeeManager
type FeeManagerViewGenerator struct{}

// FeeManagerViewGenerator implements ContractViewGenerator
var _ interfaces.ContractViewGenerator[fee_manager.FeeManagerInterface, FeeManagerContext, FeeManagerView] = (*FeeManagerViewGenerator)(nil)

func NewFeeManagerViewGenerator() *FeeManagerViewGenerator {
	return &FeeManagerViewGenerator{}
}

func (f *FeeManagerViewGenerator) Generate(ctx context.Context, contract fee_manager.FeeManagerInterface, params FeeManagerContext) (FeeManagerView, error) {
	callOpts := &bind.CallOpts{Context: ctx}

	owner, err := contract.Owner(callOpts)
	if err != nil {
		return FeeManagerView{}, fmt.Errorf("failed to get owner: %w", err)
	}

	linkAddress, err := contract.ILinkAddress(callOpts)
	if err != nil {
		return FeeManagerView{}, fmt.Errorf("failed to get link address: %w", err)
	}

	nativeAddress, err := contract.INativeAddress(callOpts)
	if err != nil {
		return FeeManagerView{}, fmt.Errorf("failed to get native address: %w", err)
	}

	proxyAddress, err := contract.IProxyAddress(callOpts)
	if err != nil {
		return FeeManagerView{}, fmt.Errorf("failed to get proxy address: %w", err)
	}

	rewardManager, err := contract.IRewardManager(callOpts)
	if err != nil {
		return FeeManagerView{}, fmt.Errorf("failed to get reward manager: %w", err)
	}

	nativeSurcharge, err := contract.SNativeSurcharge(callOpts)
	if err != nil {
		return FeeManagerView{}, fmt.Errorf("failed to get native surcharge: %w", err)
	}

	linkAvailable, err := contract.LinkAvailableForPayment(callOpts)
	if err != nil {
		return FeeManagerView{}, fmt.Errorf("failed to get link available: %w", err)
	}

	typeAndVersion, err := contract.TypeAndVersion(callOpts)
	if err != nil {
		return FeeManagerView{}, fmt.Errorf("failed to get type and version: %w", err)
	}

	discounts, err := f.gatherOrganizedDiscounts(ctx, contract, params, linkAddress.Hex(), nativeAddress.Hex())
	if err != nil {
		return FeeManagerView{}, fmt.Errorf("failed to gather organized discounts: %w", err)
	}

	// Create and return the view
	view := &FeeManagerView{
		Owner:               owner,
		LinkAddress:         linkAddress.Hex(),
		NativeAddress:       nativeAddress.Hex(),
		ProxyAddress:        proxyAddress.Hex(),
		RewardManager:       rewardManager.Hex(),
		NativeSurcharge:     nativeSurcharge.String(),
		LinkAvailable:       linkAvailable.String(),
		TypeAndVersion:      typeAndVersion,
		SubscriberDiscounts: discounts,
	}

	return *view, nil
}

// Function to gather all discounts and organize them by subscriber and feedId
func (f *FeeManagerViewGenerator) gatherOrganizedDiscounts(ctx context.Context,
	contract fee_manager.FeeManagerInterface,
	params FeeManagerContext,
	linkAddress string,
	nativeAddress string) (map[string]map[string]TokenDiscounts, error) {

	// Create filter options
	filterOpts := &bind.FilterOpts{
		Start:   params.FromBlock,
		End:     params.ToBlock,
		Context: ctx,
	}

	// Get all subscriber discount events
	iterator, err := contract.FilterSubscriberDiscountUpdated(filterOpts, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to filter subscriber discount events: %w", err)
	}
	defer iterator.Close()

	// Get references to token addresses for comparison
	callOpts := &bind.CallOpts{Context: ctx}

	type discountKey struct {
		subscriber common.Address
		feedId     [32]byte
		token      common.Address
	}

	// Find all combinations of subscriber, feedId, and token
	discountMap := make(map[string]discountKey)
	for iterator.Next() {
		event := iterator.Event

		feedIdStr := dsutil.HexEncodeBytes32(event.FeedId)

		// Create a unique key for this combination
		key := fmt.Sprintf("%s-%s-%s", event.Subscriber.Hex(), feedIdStr, event.Token.Hex())

		// Store the combination
		discountMap[key] = discountKey{
			subscriber: event.Subscriber,
			feedId:     event.FeedId,
			token:      event.Token,
		}
	}

	if err := iterator.Error(); err != nil {
		return nil, fmt.Errorf("error iterating through events: %w", err)
	}

	// Map[subscriberAddress][feedId]TokenDiscounts
	result := make(map[string]map[string]TokenDiscounts)

	for _, combo := range discountMap {
		subscriberAddr := combo.subscriber.Hex()
		feedIdHex := dsutil.HexEncodeBytes32(combo.feedId)

		// global discount is set using feedId of all zeros
		isGlobalDiscount := false
		var zeroBytes [32]byte
		if combo.feedId == zeroBytes {
			isGlobalDiscount = true
			feedIdHex = "global"
		}

		if result[subscriberAddr] == nil {
			result[subscriberAddr] = make(map[string]TokenDiscounts)
		}

		tokenDiscounts := result[subscriberAddr][feedIdHex]
		tokenDiscounts.IsGlobal = isGlobalDiscount

		var discount *big.Int
		var err error

		if isGlobalDiscount {
			discount, err = contract.SGlobalDiscounts(callOpts, combo.subscriber, combo.token)
		} else {
			discount, err = contract.SSubscriberDiscounts(callOpts, combo.subscriber, combo.feedId, combo.token)
		}

		if err != nil {
			return nil, fmt.Errorf("failed to query discount: %w", err)
		}

		if combo.token.String() == linkAddress {
			tokenDiscounts.Link = discount.String()
		} else if combo.token.String() == nativeAddress {
			tokenDiscounts.Native = discount.String()
		}

		result[subscriberAddr][feedIdHex] = tokenDiscounts

	}

	return result, nil
}
