package modsecverifier

import (
	"context"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/modsec/modsecstorage"
)

var _ job.ServiceCtx = &verifier{}

// verifier is a service that monitors the source chain for CCIPMessageSent events
// and pushes attestations to them to offchain storage.
type verifier struct {
	lggr          logger.Logger
	lp            logpoller.LogPoller
	wg            sync.WaitGroup
	runCtx        context.Context
	runCtxCancel  context.CancelFunc
	eventSig      string
	onRampAddress string
	storage       modsecstorage.Storage
}

func New(lggr logger.Logger, lp logpoller.LogPoller, eventSig string, onRampAddress string, storage modsecstorage.Storage) *verifier {
	return &verifier{
		lggr:          lggr,
		lp:            lp,
		eventSig:      eventSig,
		onRampAddress: onRampAddress,
		storage:       storage,
	}
}

func (v *verifier) run(ctx context.Context) error {
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			// TODO: implement logic
			v.lggr.Info("verifying modsec")
		}
	}
}

// Close implements job.ServiceCtx.
func (v *verifier) Close() error {
	v.runCtxCancel()
	v.wg.Wait()
	return nil
}

// Start implements job.ServiceCtx.
func (v *verifier) Start(context.Context) error {
	v.wg.Add(1)
	v.runCtx, v.runCtxCancel = context.WithCancel(context.Background())
	go func() {
		defer v.wg.Done()
		v.run(v.runCtx)
	}()
	return nil
}
