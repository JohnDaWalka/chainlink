package syncer

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type backoffRecord struct {
	retryCount  int
	nextRetryAt time.Time
	signature   string
}

type backoffStrategy struct {
	clock              clockwork.Clock
	tickInterval       time.Duration
	maxBackoffInterval time.Duration
	backoffRecords     map[string]backoffRecord
	handler            func(ctx context.Context, event Event) error
	lggr               logger.Logger
}

func newBackoffStrategy(lggr logger.Logger, handler func(ctx context.Context, event Event) error, backoffInterval, maxBackoffInterval time.Duration) *backoffStrategy {
	return &backoffStrategy{
		clock:              clockwork.NewRealClock(),
		tickInterval:       backoffInterval,
		maxBackoffInterval: maxBackoffInterval,
		backoffRecords:     make(map[string]backoffRecord),
		handler:            handler,
		lggr:               lggr.Named("BackoffStrategy"),
	}
}

func (b *backoffStrategy) Apply(ctx context.Context, event Event) {
	key, sig, err := keysFor(event)
	if err != nil {
		b.lggr.Errorw("failed to apply backoff strategy: could not calculate keys for event", "err", err, "type", event.GetEventType())
		return
	}

	// We have an existing backoffRecord for the key, but it might correspond to a different
	// event. If it does, let's treat this as a reset and delete the record.
	// If not, it's the same record, so we'll inspect it when determining whether to execute or not.
	br, ok := b.backoffRecords[key]
	if ok && br.signature != sig {
		b.lggr.Debugw("found record matching workflowID with different signature, deleting...", "key", key, "oldSignature", br.signature, "newSignature", sig)
		delete(b.backoffRecords, key)
	}

	// If we don't have a backoffRecord, or the time listed in the backoffRecord
	// is in the past, we can execute now. Otherwise, log a message and continue.
	br, ok = b.backoffRecords[key]
	shouldExecuteNow := !ok || b.clock.Now().After(br.nextRetryAt)
	if !shouldExecuteNow {
		b.lggr.Debugw("skipping event as it is still in backoff", "key", key, "nextRetryAt", br.nextRetryAt, "retryCount", br.retryCount)
		return
	}

	err = b.handler(ctx, event)
	if err != nil {
		backoff := math.Pow(float64(2), float64(br.retryCount)) * float64(b.tickInterval)
		backoff = math.Min(backoff, float64(b.maxBackoffInterval))
		backoffDuration := time.Duration(backoff)

		retryCount := br.retryCount + 1
		nextRetryAt := b.clock.Now().Add(backoffDuration)
		b.backoffRecords[key] = backoffRecord{
			retryCount:  retryCount,
			nextRetryAt: nextRetryAt,
			signature:   sig,
		}
		b.lggr.Errorw("failed to handle event, backing off...", "err", err, "type", event.GetEventType(), "retryCount", retryCount, "nextRetryAt", nextRetryAt, "backoff", backoffDuration, "key", key)
	} else {
		delete(b.backoffRecords, key)
	}

	return
}

func idFor(owner []byte, name string) string {
	return fmt.Sprintf("%x-%s", owner, name)
}

// sigsFor generates two signatures for the given event
// - a `workflowKey`, which uniquely identifies a workflow by owner + name
// - a 'signature` which uniquely identifies the event
func keysFor(event Event) (string, string, error) {
	data := event.GetData()
	switch td := data.(type) {
	case WorkflowRegistryWorkflowRegisteredV1:
		return idFor(td.WorkflowOwner, td.WorkflowName), fmt.Sprintf("%s-%s-%s", WorkflowRegisteredEvent, td.WorkflowID.Hex(), toSpecStatus(td.Status)), nil
	case WorkflowRegistryWorkflowDeletedV1:
		return idFor(td.WorkflowOwner, td.WorkflowName), fmt.Sprintf("%s-%s", WorkflowDeletedEvent, td.WorkflowID.Hex()), nil
	case WorkflowRegistryWorkflowUpdatedV1:
		return idFor(td.WorkflowOwner, td.WorkflowName), fmt.Sprintf("%s-%s-%s-%s", WorkflowUpdatedEvent, td.OldWorkflowID.Hex(), td.NewWorkflowID.Hex(), toSpecStatus(td.Status)), nil
	case WorkflowRegistryWorkflowPausedV1:
		return idFor(td.WorkflowOwner, td.WorkflowName), fmt.Sprintf("%s-%s", WorkflowPausedEvent, td.WorkflowID.Hex()), nil
	case WorkflowRegistryWorkflowActivatedV1:
		return idFor(td.WorkflowOwner, td.WorkflowName), fmt.Sprintf("%s-%s", WorkflowPausedEvent, td.WorkflowID.Hex()), nil
	default:
		return "", "", fmt.Errorf("could not extract workflow ID from event type %T", event)
	}
}
