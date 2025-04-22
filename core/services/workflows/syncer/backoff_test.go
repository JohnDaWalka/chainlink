package syncer

import (
	"context"
	"crypto/rand"
	"errors"
	"math"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
)

type handler struct {
	err       error
	callCount int
}

func (h *handler) Handle(ctx context.Context, event Event) error {
	h.callCount++
	return h.err
}

func mustRandomWID() types.WorkflowID {
	var randomBytes [32]byte
	_, err := rand.Read(randomBytes[:])
	if err != nil {
		panic(err)
	}
	return types.WorkflowID(randomBytes)
}

func Test_BackoffStrategy(t *testing.T) {
	h := &handler{err: errors.New("oops")}
	backoffInterval := 10 * time.Millisecond
	maxBackoffInterval := 100 * time.Millisecond
	bs := newBackoffStrategy(logger.TestLogger(t), h.Handle, backoffInterval, maxBackoffInterval)
	fakeClock := clockwork.NewFakeClock()
	bs.clock = fakeClock

	wid := mustRandomWID()
	owner := []byte("owner")
	name := "name"
	id := idFor(owner, name)
	bs.Apply(context.Background(), &workflowAsEvent{
		Data: WorkflowRegistryWorkflowRegisteredV1{
			WorkflowID:    wid,
			WorkflowOwner: owner,
			WorkflowName:  name,
			Status:        StatusActivated,
		},
		EventType: WorkflowRegisteredEvent,
	})

	// First try -- call the handler and create the record
	assert.Equal(t, 1, len(bs.backoffRecords))
	assert.Equal(t, 1, bs.backoffRecords[id].retryCount)
	assert.Equal(t, 1, h.callCount)

	// Second retry -- don't call the handler (we haven't hit the backoff threshold)
	bs.Apply(context.Background(), &workflowAsEvent{
		Data: WorkflowRegistryWorkflowRegisteredV1{
			WorkflowID:    wid,
			WorkflowOwner: owner,
			WorkflowName:  name,
			Status:        StatusActivated,
		},
		EventType: WorkflowRegisteredEvent,
	})
	assert.Equal(t, 1, bs.backoffRecords[id].retryCount)

	// Now advance the clock to the next retry time
	fakeClock.Advance(backoffInterval + 1*time.Millisecond)
	bs.Apply(context.Background(), &workflowAsEvent{
		Data: WorkflowRegistryWorkflowRegisteredV1{
			WorkflowID:    wid,
			WorkflowOwner: owner,
			WorkflowName:  name,
			Status:        StatusActivated,
		},
		EventType: WorkflowRegisteredEvent,
	})
	assert.Equal(t, 2, bs.backoffRecords[id].retryCount)
	assert.Equal(t, 2, h.callCount)

	// Skip ahead and confirm that the backoff is exponential
	br := bs.backoffRecords[id]
	br.retryCount = 5
	bs.backoffRecords[id] = br

	interval := math.Pow(float64(2), float64(br.retryCount)) * float64(backoffInterval)
	fakeClock.Advance(time.Duration(interval) + 1*time.Millisecond)
	bs.Apply(context.Background(), &workflowAsEvent{
		Data: WorkflowRegistryWorkflowRegisteredV1{
			WorkflowID:    wid,
			WorkflowOwner: owner,
			WorkflowName:  name,
			Status:        StatusActivated,
		},
		EventType: WorkflowRegisteredEvent,
	})
	assert.Equal(t, 6, bs.backoffRecords[id].retryCount)
	assert.Equal(t, 3, h.callCount)
}

func Test_BackoffStrategy_BackoffResets(t *testing.T) {
	h := &handler{err: errors.New("oops")}
	backoffInterval := 10 * time.Millisecond
	maxBackoffInterval := 100 * time.Millisecond
	bs := newBackoffStrategy(logger.TestLogger(t), h.Handle, backoffInterval, maxBackoffInterval)
	fakeClock := clockwork.NewFakeClock()
	bs.clock = fakeClock

	wid := mustRandomWID()
	owner := []byte("owner")
	name := "name"
	id := idFor(owner, name)
	bs.Apply(context.Background(), &workflowAsEvent{
		Data: WorkflowRegistryWorkflowRegisteredV1{
			WorkflowID:    wid,
			WorkflowOwner: owner,
			WorkflowName:  name,
			Status:        StatusActivated,
		},
		EventType: WorkflowRegisteredEvent,
	})

	assert.Equal(t, 1, len(bs.backoffRecords))
	assert.Equal(t, 1, bs.backoffRecords[id].retryCount)
	assert.Equal(t, 1, h.callCount)

	// The record should be removed if the handler succeeds
	fakeClock.Advance(backoffInterval + 1*time.Millisecond)
	h.err = nil
	bs.Apply(context.Background(), &workflowAsEvent{
		Data: WorkflowRegistryWorkflowRegisteredV1{
			WorkflowID:    wid,
			WorkflowOwner: owner,
			WorkflowName:  name,
			Status:        StatusActivated,
		},
		EventType: WorkflowRegisteredEvent,
	})
	assert.Len(t, bs.backoffRecords, 0)

	// Now check a reset if we get a different event
	// that fails
	h.err = errors.New("oops")
	bs.Apply(context.Background(), &workflowAsEvent{
		Data: WorkflowRegistryWorkflowRegisteredV1{
			WorkflowID:    wid,
			WorkflowOwner: owner,
			WorkflowName:  name,
			Status:        StatusActivated,
		},
		EventType: WorkflowRegisteredEvent,
	})

	assert.Equal(t, 1, len(bs.backoffRecords))
	assert.Equal(t, 1, bs.backoffRecords[id].retryCount)
	assert.Equal(t, 3, h.callCount)

	bs.Apply(context.Background(), &workflowAsEvent{
		Data: WorkflowRegistryWorkflowUpdatedV1{
			OldWorkflowID: wid,
			NewWorkflowID: mustRandomWID(),
			WorkflowOwner: owner,
			WorkflowName:  name,
			Status:        StatusActivated,
		},
		EventType: WorkflowUpdatedEvent,
	})

	assert.Equal(t, 1, len(bs.backoffRecords))
	assert.Equal(t, 1, bs.backoffRecords[id].retryCount)
	assert.Equal(t, 4, h.callCount)

	// Now check a reset if we get a different event
	// that succeeds -- there should be no record
	h.err = nil
	bs.Apply(context.Background(), &workflowAsEvent{
		Data: WorkflowRegistryWorkflowRegisteredV1{
			WorkflowID:    wid,
			WorkflowOwner: owner,
			WorkflowName:  name,
			Status:        StatusPaused,
		},
		EventType: WorkflowRegisteredEvent,
	})

	assert.Equal(t, 0, len(bs.backoffRecords))
	assert.Equal(t, 5, h.callCount)
}
