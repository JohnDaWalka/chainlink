package v2

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/dontime"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestDonTimeProvider_GetDONTime(t *testing.T) {
	lggr := logger.TestLogger(t)
	workflowExecutionID := "test-workflow-123"

	t.Run("successful case - RequestDonTime returns valid response", func(t *testing.T) {
		store := dontime.NewStore(20 * time.Minute)
		provider := NewDonTimeProvider(store, workflowExecutionID, lggr)
		expectedTimestamp := int64(1609459200000) // 2021-01-01 00:00:00 UTC in milliseconds

		go func() {
			for {
				time.Sleep(10 * time.Millisecond)
				req := store.GetRequest(workflowExecutionID)
				if req != nil {
					resp := dontime.Response{
						WorkflowExecutionID: workflowExecutionID,
						SeqNum:              0,
						Timestamp:           expectedTimestamp,
						Err:                 nil,
					}
					req.SendResponse(t.Context(), resp)
					return
				}
			}
		}()

		result, err := provider.GetDONTime()

		require.NoError(t, err)
		expectedTime := time.Unix(0, expectedTimestamp*int64(time.Millisecond)).UTC()
		assert.Equal(t, expectedTime, result)
		assert.Equal(t, 1, provider.timeSeqNum) // Should increment after call
	})

	t.Run("error case with fallback to node time", func(t *testing.T) {
		store := dontime.NewStore(20 * time.Minute)
		provider := NewDonTimeProvider(store, workflowExecutionID, lggr)

		go func() {
			for {
				time.Sleep(10 * time.Millisecond)
				req := store.GetRequest(workflowExecutionID)
				if req != nil {
					resp := dontime.Response{
						WorkflowExecutionID: workflowExecutionID,
						SeqNum:              0,
						Timestamp:           0,
						Err:                 errors.New("request failed"),
					}
					req.SendResponse(t.Context(), resp)
					return
				}
			}
		}()

		result, err := provider.GetDONTime()

		require.NoError(t, err)
		// Should fallback to node time (which defaults to zero in new store)
		expectedTime := time.Unix(0, 0).UTC()
		assert.Equal(t, expectedTime, result)
	})

	t.Run("timestamp exceeds maxTimestampMillis", func(t *testing.T) {
		store := dontime.NewStore(20 * time.Minute)
		provider := NewDonTimeProvider(store, workflowExecutionID, lggr)

		// Use a timestamp that exceeds maxTimestampMillis (4_200_000_000_000)
		excessiveTimestamp := int64(5_000_000_000_000) // Year 2128

		go func() {
			for {
				time.Sleep(10 * time.Millisecond)
				req := store.GetRequest(workflowExecutionID)
				if req != nil {
					resp := dontime.Response{
						WorkflowExecutionID: workflowExecutionID,
						SeqNum:              0,
						Timestamp:           excessiveTimestamp,
						Err:                 nil,
					}
					req.SendResponse(t.Context(), resp)
					return
				}
			}
		}()

		result, err := provider.GetDONTime()

		require.NoError(t, err)
		// Should be truncated to maxTimestampMillis
		expectedTime := time.Unix(0, maxTimestampMillis*int64(time.Millisecond)).UTC()
		assert.Equal(t, expectedTime, result)
	})
}
