package v2

import (
	"context"
	"errors"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	sdkpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
)

type disallowedExecutionHelper struct {
	*Engine
	TimeProvider
	SecretsFetcher
	lggr         logger.Logger
	userLogCount int
	userLogMu    sync.Mutex
}

func NewDisallowedExecutionHelper(lggr logger.Logger, engine *Engine, timeProvider TimeProvider, secretsFetcher SecretsFetcher) *disallowedExecutionHelper {
	return &disallowedExecutionHelper{
		Engine:         engine,
		TimeProvider:   timeProvider,
		SecretsFetcher: secretsFetcher,
		lggr:           lggr,
	}
}

var _ host.ExecutionHelper = &disallowedExecutionHelper{}

func (d *disallowedExecutionHelper) CallCapability(_ context.Context, _ *sdkpb.CapabilityRequest) (*sdkpb.CapabilityResponse, error) {
	return nil, errors.New("capability calls cannot be made during this execution")
}

func (d *disallowedExecutionHelper) GetWorkflowExecutionID() string {
	return ""
}

func (d *disallowedExecutionHelper) EmitUserLog(msg string) error {
	d.userLogMu.Lock()
	defer d.userLogMu.Unlock()
	d.userLogCount++
	return d.enqueueUserLog(msg, d.userLogCount, "") // empty execution ID for trigger subscription phase
}
