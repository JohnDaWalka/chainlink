package modsecexecutor

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/verifier_events"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/modsec/modsecstorage"
	"github.com/smartcontractkit/chainlink/v2/core/services/modsec/modsectypes"
)

// executor is a service that monitors offhcain storage for new attestations,
// and executes them on the destination chain.
type executor struct {
	lggr           logger.Logger
	transmitter    modsectypes.Transmitter
	wg             sync.WaitGroup
	runCtx         context.Context
	runCtxCancel   context.CancelFunc
	offRampAddress string
	verifierEvents *verifier_events.VerifierEvents

	// TODO: in practice multiple storages and they are discovered.
	// Keeping it just a single storage for now to get an E2E test going.
	storage modsecstorage.Storage
}

var _ job.ServiceCtx = &executor{}

func New(
	lggr logger.Logger,
	transmitter modsectypes.Transmitter,
	offRampAddress string,
	verifierEvents *verifier_events.VerifierEvents,
	storage modsecstorage.Storage,
) *executor {
	return &executor{
		lggr:           lggr,
		transmitter:    transmitter,
		offRampAddress: offRampAddress,
		storage:        storage,
		verifierEvents: verifierEvents,
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
			if err := r.loop(ctx); err != nil {
				r.lggr.Error("failed to run executor loop", "error", err)
			}
		}
	}
}

func (r *executor) loop(ctx context.Context) error {
	all, err := r.storage.GetAll(ctx)
	if err != nil {
		r.lggr.Error("failed to get all attestations", "error", err)
		return err
	}

	r.lggr.Infow("got attestations", "count", len(all))
	for messageID, valuePayload := range all {
		var payload modsectypes.StorageValuePayload
		if err := json.Unmarshal(valuePayload, &payload); err != nil {
			r.lggr.Error("failed to unmarshal value payload", "error", err)
			continue
		}

		var messageHash [32]byte
		copy(messageHash[:], payload.MessageHash)
		// check if they're already executed onchain
		executed, err := r.verifierEvents.SMessageExecuted(&bind.CallOpts{Context: ctx}, messageHash)
		if err != nil {
			r.lggr.Error("failed to check if message was executed", "error", err)
			continue
		}
		if executed {
			r.lggr.Infow("message already executed, skipping", "messageID", messageID)
			continue
		}
		// execute the message via the transmitter
		err = r.transmitter.Transmit(ctx, payload)
		if err != nil {
			r.lggr.Error("failed to execute message", "error", err)
			continue
		}

		// TODO: inflight cache?
		r.lggr.Info("transmitted message", "messageID", messageID)
	}

	return nil
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
