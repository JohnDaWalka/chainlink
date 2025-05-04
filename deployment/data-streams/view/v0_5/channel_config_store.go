package v0_5

import "github.com/ethereum/go-ethereum/common"

type ChannelConfigStoreView struct {
	TypeAndVersion string         `json:"typeAndVersion,omitempty"`
	Owner          common.Address `json:"owner,omitempty"`
}
