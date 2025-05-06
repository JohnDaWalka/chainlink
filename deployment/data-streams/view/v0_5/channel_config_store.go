package v0_5

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/view/interfaces"
)

type ChannelConfigStoreView struct {
	TypeAndVersion string         `json:"typeAndVersion,omitempty"`
	Owner          common.Address `json:"owner,omitempty"`
}

// ChannelConfigStoreView implements the ContractView interface
var _ interfaces.ContractView = (*ChannelConfigStoreView)(nil)

// SerializeView serializes view to JSON
func (v ChannelConfigStoreView) SerializeView() (string, error) {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal contract view: %w", err)
	}
	return string(bytes), nil
}
