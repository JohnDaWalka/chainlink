package v0_5

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/view/interfaces"

	rewardManager "github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/reward_manager_v0_5_0"
)

type RewardManagerView struct {
	TypeAndVersion string         `json:"typeAndVersion,omitempty"`
	Address        common.Address `json:"address,omitempty"`
	Owner          common.Address `json:"owner,omitempty"`
}

// RewardManagerView implements the ContractView interface
var _ interfaces.ContractView = (*RewardManagerView)(nil)

// SerializeView serializes view to JSON
func (v RewardManagerView) SerializeView() (string, error) {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal contract view: %w", err)
	}
	return string(bytes), nil
}

// GenerateRewardManagerView generates a RewardManagerView from a RewardManager contract.
func GenerateRewardManagerView(rm *rewardManager.RewardManager) (RewardManagerView, error) {
	if rm == nil {
		return RewardManagerView{}, errors.New("cannot generate view for nil RewardManager")
	}

	owner, err := rm.Owner(nil)
	if err != nil {
		return RewardManagerView{}, fmt.Errorf("failed to get owner for RewardManager: %w", err)
	}

	return RewardManagerView{
		Address:        rm.Address(),
		Owner:          owner,
		TypeAndVersion: "RewardManager 0.5.0",
	}, nil
}
