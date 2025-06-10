package v2

import (
	"context"
	"errors"
	"time"

	sdkpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
)

type DisallowedExecutionHelper struct{}

var _ host.ExecutionHelper = DisallowedExecutionHelper{}

func (d DisallowedExecutionHelper) CallCapability(_ context.Context, _ *sdkpb.CapabilityRequest) (*sdkpb.CapabilityResponse, error) {
	return nil, errors.New("capability calls cannot be made during this execution")
}

func (d DisallowedExecutionHelper) GetId() string {
	return ""
}

func (d DisallowedExecutionHelper) GetNodeTime() time.Time {
}

func (d DisallowedExecutionHelper) GetDONTime() time.Time {
}
