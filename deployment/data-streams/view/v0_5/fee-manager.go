package v0_5

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/view/interfaces"

	feeManager "github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/fee_manager_v0_5_0"
)

type FeeManagerView struct {
	TypeAndVersion string         `json:"typeAndVersion,omitempty"`
	Address        common.Address `json:"address,omitempty"`
	Owner          common.Address `json:"owner,omitempty"`
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

// GenerateFeeManagerView generates a FeeManagerView from a FeeManager contract.
func GenerateFeeManagerView(rm *feeManager.FeeManager) (FeeManagerView, error) {
	if rm == nil {
		return FeeManagerView{}, errors.New("cannot generate view for nil FeeManager")
	}

	owner, err := rm.Owner(nil)
	if err != nil {
		return FeeManagerView{}, fmt.Errorf("failed to get owner for FeeManager: %w", err)
	}

	return FeeManagerView{
		Address:        rm.Address(),
		Owner:          owner,
		TypeAndVersion: "FeeManager 0.5.0",
	}, nil
}
