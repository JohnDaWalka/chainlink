package v0_5

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/reward_manager_v0_5_0"
)

// RecipientInfo represents the data for each recipient
type RecipientInfo struct {
	PayeeAddress common.Address `json:"payeeAddress"`
	Weight       string         `json:"weight"`
}

// RewardManagerView represents a processed view of reward manager data
type RewardManagerView struct {
	RecipientWeights map[string][]RecipientInfo `json:"recipientWeights"` // poolId -> recipient info
}

// SerializeView serializes view to JSON
func (r RewardManagerView) SerializeView() (string, error) {
	bytes, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal contract view: %w", err)
	}
	return string(bytes), nil
}

// RewardManagerViewParams defines parameters for generating the view
type RewardManagerViewParams struct {
	FromBlock uint64
	ToBlock   *uint64
}

// RewardManagerContract defines a minimal interface
type RewardManagerContract interface {
	// Methods to get pool IDs
	SRegisteredPoolIds(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error)

	// Methods to get pool and recipient data
	STotalRewardRecipientFees(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error)
	SRewardRecipientWeights(opts *bind.CallOpts, arg0 [32]byte, arg1 common.Address) (*big.Int, error)
	SRewardRecipientWeightsSet(opts *bind.CallOpts, arg0 [32]byte) (bool, error)

	// Event filtering methods
	FilterRewardRecipientsUpdated(opts *bind.FilterOpts, poolId [][32]byte) (*reward_manager_v0_5_0.RewardManagerRewardRecipientsUpdatedIterator, error)
	FilterRewardsClaimed(opts *bind.FilterOpts, poolId [][32]byte, recipient []common.Address) (*reward_manager_v0_5_0.RewardManagerRewardsClaimedIterator, error)
}

// RewardManagerViewGenerator generates views for reward manager contracts
type RewardManagerViewGenerator struct {
	contract RewardManagerContract
}

// NewRewardManagerViewGenerator creates a new view generator
func NewRewardManagerViewGenerator(contract RewardManagerContract) *RewardManagerViewGenerator {
	return &RewardManagerViewGenerator{
		contract: contract,
	}
}

// Generate creates a view based on the provided parameters
func (g *RewardManagerViewGenerator) Generate(ctx context.Context, params RewardManagerViewParams) (RewardManagerView, error) {
	view := RewardManagerView{}

	// First, collect registered pool IDs
	poolIDs, err := g.getPoolIDs(ctx)
	if err != nil {
		return view, fmt.Errorf("failed to get registered pool IDs: %w", err)
	}

	recipientWeights, err := g.getRecipientWeights(ctx, poolIDs, params)
	if err != nil {
		return view, fmt.Errorf("failed to get recipient weights: %w", err)
	}

	view.RecipientWeights = recipientWeights

	return view, nil
}

func (g *RewardManagerViewGenerator) getRecipientWeights(ctx context.Context, poolIds [][32]byte, params RewardManagerViewParams) (map[string][]RecipientInfo, error) {
	filterOpts := &bind.FilterOpts{
		Context: ctx,
		Start:   params.FromBlock,
		End:     params.ToBlock,
	}

	recipientWeights := make(map[string][]RecipientInfo)

	updateIterator, err := g.contract.FilterRewardRecipientsUpdated(filterOpts, poolIds)
	if err != nil {
		return nil, fmt.Errorf("failed to filter reward recipients updated events: %w", err)
	}
	defer updateIterator.Close()

	// Process update events to collect all recipients and their weights
	poolRecipients := make(map[[32]byte][]common.Address)

	for updateIterator.Next() {
		event := updateIterator.Event
		poolID := event.PoolId

		// Convert poolID to hex string for the view
		poolIDHex := fmt.Sprintf("%#x", poolID)

		// Initialize the recipient list for this pool
		if _, exists := recipientWeights[poolIDHex]; !exists {
			recipientWeights[poolIDHex] = []RecipientInfo{}
		}

		// Track recipients for each pool for later use
		if _, exists := poolRecipients[poolID]; !exists {
			poolRecipients[poolID] = []common.Address{}
		}

		// Process each recipient
		for _, r := range event.NewRewardRecipients {
			recipient := r.Addr
			weight := r.Weight

			// Add to the pool's recipients if not already there
			found := false
			for _, addr := range poolRecipients[poolID] {
				if addr == recipient {
					found = true
					break
				}
			}
			if !found {
				poolRecipients[poolID] = append(poolRecipients[poolID], recipient)
			}

			recipientInfo := RecipientInfo{
				PayeeAddress: recipient,
				Weight:       strconv.FormatUint(weight, 10),
			}

			// Check if this recipient is already in the list (from a previous event)
			found = false
			for i, info := range recipientWeights[poolIDHex] {
				if info.PayeeAddress == recipient {
					// Update the existing entry
					recipientWeights[poolIDHex][i] = recipientInfo
					found = true
					break
				}
			}

			// Add if not found
			if !found {
				recipientWeights[poolIDHex] = append(recipientWeights[poolIDHex], recipientInfo)
			}
		}
	}

	// Check for errors in the iterator
	if err := updateIterator.Error(); err != nil {
		return nil, fmt.Errorf("error iterating through reward recipients updated events: %w", err)
	}

	return recipientWeights, nil

}

func (g *RewardManagerViewGenerator) getPoolIDs(ctx context.Context) ([][32]byte, error) {
	var poolIDs [][32]byte
	for i := int64(0); ; i++ {
		poolID, err := g.contract.SRegisteredPoolIds(&bind.CallOpts{Context: ctx}, big.NewInt(i))
		if err != nil {
			// We'll see a revert when we reach the end of registered pool IDs.
			// Note: this is a simplification. We may want to check for a specific error type instead.
			break
		}

		// Check if the weights have been set for this pool
		weightsSet, err := g.contract.SRewardRecipientWeightsSet(&bind.CallOpts{Context: ctx}, poolID)
		if err != nil {
			return nil, fmt.Errorf("failed to check if weights are set for pool %x: %w", poolID, err)
		}

		// Only include pools that have weights set
		if weightsSet {
			poolIDs = append(poolIDs, poolID)
		}
	}
	return poolIDs, nil
}
