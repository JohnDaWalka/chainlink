package modsecrelayer

import (
	"context"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink-evm/pkg/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

// relayer is a service that monitors the source chain for CCIPMessageSent events,
// observes attestations for them on offchain storage, and executes them on the destination chain.
type relayer struct {
	lggr           logger.Logger
	lp             logpoller.LogPoller
	txm            txmgr.TxManager
	wg             sync.WaitGroup
	runCtx         context.Context
	runCtxCancel   context.CancelFunc
	eventSig       string
	onRampAddress  string
	offRampAddress string
}

var _ job.ServiceCtx = &relayer{}

func NewRelayer(lggr logger.Logger, lp logpoller.LogPoller, txm txmgr.TxManager, eventSig string, onRampAddress string, offRampAddress string) *relayer {
	return &relayer{
		lggr:           lggr,
		lp:             lp,
		txm:            txm,
		eventSig:       eventSig,
		onRampAddress:  onRampAddress,
		offRampAddress: offRampAddress,
	}
}

func (r *relayer) run(ctx context.Context) error {
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
func (r *relayer) Close() error {
	r.runCtxCancel()
	r.wg.Wait()
	return nil
}

// Start implements job.ServiceCtx.
func (r *relayer) Start(context.Context) error {
	r.wg.Add(1)
	r.runCtx, r.runCtxCancel = context.WithCancel(context.Background())
	go func() {
		defer r.wg.Done()
		r.run(r.runCtx)
	}()
	return nil
}
