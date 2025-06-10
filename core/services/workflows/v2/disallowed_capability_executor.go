package v2

import (
	"context"
	"errors"
	"time"

	sdkpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
)

type DisallowedCapabilityExecutor struct{}

var _ host.ExecutionHelper = DisallowedCapabilityExecutor{}

func (d DisallowedCapabilityExecutor) CallCapability(_ context.Context, _ *sdkpb.CapabilityRequest) (*sdkpb.CapabilityResponse, error) {
	return nil, errors.New("capability calls cannot be made during this execution")
}

func (d DisallowedCapabilityExecutor) GetDONTime() time.Time {
	return time.Now()
}

func (d DisallowedCapabilityExecutor) GetId() string {
	return ""
}

func (d DisallowedCapabilityExecutor) GetNodeTime() time.Time {
	return time.Now()
}
