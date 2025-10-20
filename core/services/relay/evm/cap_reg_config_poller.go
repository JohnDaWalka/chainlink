package evm

import (
	"context"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

var _ types.ConfigPoller = &CapRegConfigPoller{}

func NewCapRegConfigPoller(configTracker ocrtypes.ContractConfigTracker) *CapRegConfigPoller {
	return &CapRegConfigPoller{
		ContractConfigTracker: configTracker,
	}
}

type CapRegConfigPoller struct {
	ocrtypes.ContractConfigTracker
}

func (cp *CapRegConfigPoller) Start() {}

func (cp *CapRegConfigPoller) Close() error {
	return nil
}

func (cp *CapRegConfigPoller) Replay(ctx context.Context, fromBlock int64) error {
	return nil
}
