package modsecexecutor

import (
	"context"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/modsec/modsecstorage"
	"github.com/smartcontractkit/chainlink/v2/core/services/modsec/modsectypes"
)

// executor is a service that monitors the source chain for CCIPMessageSent events,
// observes attestations for them on offchain storage, and executes them on the destination chain.
type executor struct {
	lggr           logger.Logger
	lp             logpoller.LogPoller
	transmitter    modsectypes.Transmitter
	wg             sync.WaitGroup
	runCtx         context.Context
	runCtxCancel   context.CancelFunc
	eventSig       string
	onRampAddress  string
	offRampAddress string
	storage        modsecstorage.Storage
}

var _ job.ServiceCtx = &executor{}

func New(
	lggr logger.Logger,
	lp logpoller.LogPoller,
	transmitter modsectypes.Transmitter,
	eventSig string,
	onRampAddress string,
	offRampAddress string,
	storage modsecstorage.Storage,
) *executor {
	return &executor{
		lggr:           lggr,
		lp:             lp,
		transmitter:    transmitter,
		eventSig:       eventSig,
		onRampAddress:  onRampAddress,
		offRampAddress: offRampAddress,
		storage:        storage,
	}
}

func (r *executor) run(ctx context.Context) error {
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			r.lggr.Info("relaying modsec")
		}
	}
}

// Close implements job.ServiceCtx.
func (r *executor) Close() error {
	r.runCtxCancel()
	r.wg.Wait()
	return nil
}

// Start implements job.ServiceCtx.
func (r *executor) Start(context.Context) error {
	r.wg.Add(1)
	r.runCtx, r.runCtxCancel = context.WithCancel(context.Background())
	go func() {
		defer r.wg.Done()
		r.run(r.runCtx)
	}()
	return nil
}
