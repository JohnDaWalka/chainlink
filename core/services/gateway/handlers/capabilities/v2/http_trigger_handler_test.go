package v2

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	gateway_common "github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	handlermocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/mocks"
)

func requireUserErrorSent(t *testing.T, callbackCh chan handlers.UserCallbackPayload, errorCode int) {
	select {
	case payload := <-callbackCh:
		require.NotEmpty(t, payload.RawResponse)
		fmt.Printf("Received error payload: %+v\n", payload.RawResponse)
		require.Equal(t, api.ErrorCode(errorCode), payload.ErrorCode)
	case <-t.Context().Done():
		t.Fatal("Expected error callback")
	}
}

func TestHttpTriggerHandler_HandleUserTriggerRequest(t *testing.T) {
	t.Run("successful trigger request", func(t *testing.T) {
		handler, mockDon := createTestTriggerHandler(t)
		callbackCh := make(chan<- handlers.UserCallbackPayload, 1)

		triggerReq := createTestTriggerRequest()
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		req := &jsonrpc.Request{
			Version: "2.0",
			ID:      "test-request-id",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  reqBytes,
		}

		// Mock DON to expect sends to all nodes
		mockDon.EXPECT().SendToNode(mock.Anything, "node1", req).Return(nil)
		mockDon.EXPECT().SendToNode(mock.Anything, "node2", req).Return(nil)
		mockDon.EXPECT().SendToNode(mock.Anything, "node3", req).Return(nil)

		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callbackCh)
		require.NoError(t, err)

		// Verify callback was saved
		executionID, err := workflows.EncodeExecutionID("test-workflow-id", "test-request-id")
		require.NoError(t, err)

		handler.callbacksMu.Lock()
		saved, exists := handler.callbacks[executionID]
		handler.callbacksMu.Unlock()

		require.True(t, exists)
		require.Equal(t, callbackCh, saved.callbackCh)
		require.NotNil(t, saved.responseAggregator)
	})

	t.Run("invalid JSON params", func(t *testing.T) {
		handler, _ := createTestTriggerHandler(t)
		callbackCh := make(chan handlers.UserCallbackPayload, 1)

		req := &jsonrpc.Request{
			Version: "2.0",
			ID:      "test-request-id",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  []byte(`{invalid json}`),
		}

		err := handler.HandleUserTriggerRequest(testutils.Context(t), req, callbackCh)
		require.Error(t, err)

		requireUserErrorSent(t, callbackCh, int(jsonrpc.ErrParse)) //nolint:gosec // safe to cast
	})

	t.Run("empty request ID", func(t *testing.T) {
		handler, _ := createTestTriggerHandler(t)
		callbackCh := make(chan handlers.UserCallbackPayload, 1)

		triggerReq := createTestTriggerRequest()
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		req := &jsonrpc.Request{
			Version: "2.0",
			ID:      "", // Empty ID
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  reqBytes,
		}

		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callbackCh)
		require.Error(t, err)
		require.Contains(t, err.Error(), "empty request ID")

		requireUserErrorSent(t, callbackCh, int(jsonrpc.ErrInvalidRequest)) //nolint:gosec // safe to cast
	})

	t.Run("request ID contains slash", func(t *testing.T) {
		handler, _ := createTestTriggerHandler(t)
		callbackCh := make(chan handlers.UserCallbackPayload, 1)

		triggerReq := gateway_common.HTTPTriggerRequest{
			Workflow: gateway_common.WorkflowSelector{
				WorkflowID: "test-workflow-id",
			},
			Input: []byte(`{"key": "value"}`),
		}
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		req := &jsonrpc.Request{
			Version: "2.0",
			ID:      "test/request/id", // Contains slashes
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  reqBytes,
		}

		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callbackCh)
		require.Error(t, err)
		require.Contains(t, err.Error(), "must not contain '/'")

		requireUserErrorSent(t, callbackCh, int(jsonrpc.ErrInvalidRequest)) //nolint:gosec // safe to cast
	})

	t.Run("invalid method", func(t *testing.T) {
		handler, _ := createTestTriggerHandler(t)
		callbackCh := make(chan handlers.UserCallbackPayload, 1)

		triggerReq := gateway_common.HTTPTriggerRequest{
			Workflow: gateway_common.WorkflowSelector{
				WorkflowID: "test-workflow-id",
			},
			Input: []byte(`{"key": "value"}`),
		}
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		req := &jsonrpc.Request{
			Version: "2.0",
			ID:      "test-request-id",
			Method:  "invalid-method",
			Params:  reqBytes,
		}

		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callbackCh)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid method")

		requireUserErrorSent(t, callbackCh, int(jsonrpc.ErrMethodNotFound)) //nolint:gosec // safe to cast
	})

	t.Run("duplicate request ID", func(t *testing.T) {
		handler, mockDon := createTestTriggerHandler(t)
		callbackCh1 := make(chan handlers.UserCallbackPayload, 1)
		callbackCh2 := make(chan handlers.UserCallbackPayload, 1)

		triggerReq := gateway_common.HTTPTriggerRequest{
			Workflow: gateway_common.WorkflowSelector{
				WorkflowID: "test-workflow-id",
			},
			Input: []byte(`{"key": "value"}`),
		}
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		req := &jsonrpc.Request{
			Version: "2.0",
			ID:      "test-request-id",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  reqBytes,
		}

		// First request should succeed
		mockDon.EXPECT().SendToNode(mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(3)
		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callbackCh1)
		require.NoError(t, err)

		// Second request with same ID should fail
		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callbackCh2)
		require.Error(t, err)
		require.Contains(t, err.Error(), "request ID already used")
		requireUserErrorSent(t, callbackCh2, int(jsonrpc.ErrInvalidRequest)) //nolint:gosec // safe to cast
	})

	t.Run("DON send failure", func(t *testing.T) {
		handler, mockDon := createTestTriggerHandler(t)
		callbackCh := make(chan handlers.UserCallbackPayload, 1)

		triggerReq := gateway_common.HTTPTriggerRequest{
			Workflow: gateway_common.WorkflowSelector{
				WorkflowID: "test-workflow-id",
			},
			Input: []byte(`{"key": "value"}`),
		}
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		req := &jsonrpc.Request{
			Version: "2.0",
			ID:      "test-request-id",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  reqBytes,
		}

		// Mock one node to fail
		mockDon.EXPECT().SendToNode(mock.Anything, "node1", req).Return(nil)
		mockDon.EXPECT().SendToNode(mock.Anything, "node2", req).Return(errors.New("send failed"))
		mockDon.EXPECT().SendToNode(mock.Anything, "node3", req).Return(nil)

		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callbackCh)
		require.Error(t, err)
		require.Contains(t, err.Error(), "send failed")
	})

	t.Run("invalid input JSON", func(t *testing.T) {
		handler, _ := createTestTriggerHandler(t)
		callbackCh := make(chan handlers.UserCallbackPayload, 1)

		reqBytes := []byte(`{"workflow":{"workflowID":"test-workflow-id"},"input":{"invalid json"}`)
		req := &jsonrpc.Request{
			Version: "2.0",
			ID:      "test-request-id",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  reqBytes,
		}

		err := handler.HandleUserTriggerRequest(testutils.Context(t), req, callbackCh)
		require.Error(t, err)
	})
}

func TestHttpTriggerHandler_HandleNodeTriggerResponse(t *testing.T) {
	t.Run("successful aggregation", func(t *testing.T) {
		handler, mockDon := createTestTriggerHandler(t)
		callbackCh := make(chan handlers.UserCallbackPayload, 1)

		// First, create a trigger request to set up the callback
		triggerReq := gateway_common.HTTPTriggerRequest{
			Workflow: gateway_common.WorkflowSelector{
				WorkflowID: "test-workflow-id",
			},
			Input: []byte(`{"key": "value"}`),
		}
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		req := &jsonrpc.Request{
			Version: "2.0",
			ID:      "test-request-id",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  reqBytes,
		}

		mockDon.EXPECT().SendToNode(mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(3)
		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callbackCh)
		require.NoError(t, err)

		executionID, err := workflows.EncodeExecutionID("test-workflow-id", "test-request-id")
		require.NoError(t, err)

		// Create node responses
		nodeResp := &jsonrpc.Response{
			Version: "2.0",
			ID:      executionID,
			Result:  json.RawMessage(`{"result":"success"}`),
		}

		// Send responses from multiple nodes (need 2f+1 = 3 for f=1)
		err = handler.HandleNodeTriggerResponse(testutils.Context(t), nodeResp, "node1")
		require.NoError(t, err)

		err = handler.HandleNodeTriggerResponse(testutils.Context(t), nodeResp, "node2")
		require.NoError(t, err)

		// Third response should trigger aggregation
		err = handler.HandleNodeTriggerResponse(testutils.Context(t), nodeResp, "node3")
		require.NoError(t, err)

		// Check that callback was called
		select {
		case payload := <-callbackCh:
			require.NotEmpty(t, payload.RawResponse)
			require.Equal(t, api.NoError, payload.ErrorCode)

			var resp jsonrpc.Response
			err := json.Unmarshal(payload.RawResponse, &resp)
			require.NoError(t, err)
			require.Equal(t, nodeResp.Result, resp.Result)
		case <-t.Context().Done():
			t.Fatal("Expected callback")
		}
	})

	t.Run("callback not found", func(t *testing.T) {
		handler, _ := createTestTriggerHandler(t)

		nodeResp := &jsonrpc.Response{
			Version: "2.0",
			ID:      "nonexistent-execution-id",
			Result:  json.RawMessage(`{"result": "success"}`),
		}

		err := handler.HandleNodeTriggerResponse(testutils.Context(t), nodeResp, "node1")
		require.Error(t, err)
		require.Contains(t, err.Error(), "callback not found")
	})
}

func TestHttpTriggerHandler_ServiceLifecycle(t *testing.T) {
	t.Run("start and stop", func(t *testing.T) {
		handler, _ := createTestTriggerHandler(t)

		ctx := testutils.Context(t)
		err := handler.Start(ctx)
		require.NoError(t, err)

		err = handler.Close()
		require.NoError(t, err)
	})

	t.Run("double start and close should errors", func(t *testing.T) {
		handler, _ := createTestTriggerHandler(t)

		ctx := testutils.Context(t)
		err := handler.Start(ctx)
		require.NoError(t, err)

		err = handler.Start(ctx)
		require.Error(t, err)

		err = handler.Close()
		require.NoError(t, err)

		err = handler.Close()
		require.Error(t, err)
	})
}

func TestHttpTriggerHandler_ReapExpiredCallbacks(t *testing.T) {
	t.Run("reap expired callbacks", func(t *testing.T) {
		cfg := ServiceConfig{
			CleanUpPeriodMs:             100,
			MaxTriggerRequestDurationMs: 50,
		}
		handler, mockDon := createTestTriggerHandlerWithConfig(t, cfg)
		callbackCh := make(chan handlers.UserCallbackPayload, 1)

		// Add a callback
		triggerReq := gateway_common.HTTPTriggerRequest{
			Workflow: gateway_common.WorkflowSelector{
				WorkflowID: "test-workflow-id",
			},
			Input: []byte(`{"key": "value"}`),
		}
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		req := &jsonrpc.Request{
			Version: "2.0",
			ID:      "test-request-id",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  reqBytes,
		}

		mockDon.EXPECT().SendToNode(mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(3)
		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callbackCh)
		require.NoError(t, err)

		executionID, err := workflows.EncodeExecutionID("test-workflow-id", "test-request-id")
		require.NoError(t, err)

		// Manually set the callback's createdAt to the past to simulate expiration
		handler.callbacksMu.Lock()
		if cb, exists := handler.callbacks[executionID]; exists {
			cb.createdAt = time.Now().Add(-time.Duration(cfg.MaxTriggerRequestDurationMs+1) * time.Millisecond)
			handler.callbacks[executionID] = cb
		}
		handler.callbacksMu.Unlock()

		// Manually trigger reaping
		handler.reapExpiredCallbacks()

		// Verify callback was removed
		handler.callbacksMu.Lock()
		_, exists := handler.callbacks[executionID]
		handler.callbacksMu.Unlock()
		require.False(t, exists)
	})

	t.Run("keep non-expired callbacks", func(t *testing.T) {
		cfg := ServiceConfig{
			CleanUpPeriodMs:             100,
			MaxTriggerRequestDurationMs: 300000,
		}
		handler, mockDon := createTestTriggerHandlerWithConfig(t, cfg)
		callbackCh := make(chan handlers.UserCallbackPayload, 1)

		// Add a callback
		triggerReq := gateway_common.HTTPTriggerRequest{
			Workflow: gateway_common.WorkflowSelector{
				WorkflowID: "test-workflow-id",
			},
			Input: []byte(`{"key": "value"}`),
		}
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		req := &jsonrpc.Request{
			Version: "2.0",
			ID:      "test-request-id",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  reqBytes,
		}

		mockDon.EXPECT().SendToNode(mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(3)
		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callbackCh)
		require.NoError(t, err)

		executionID, err := workflows.EncodeExecutionID("test-workflow-id", "test-request-id")
		require.NoError(t, err)

		// Optionally, set createdAt to now (should not be expired)
		handler.callbacksMu.Lock()
		if cb, exists := handler.callbacks[executionID]; exists {
			cb.createdAt = time.Now()
			handler.callbacks[executionID] = cb
		}
		handler.callbacksMu.Unlock()

		// Manually trigger reaping
		handler.reapExpiredCallbacks()

		// Verify callback still exists
		handler.callbacksMu.Lock()
		_, exists := handler.callbacks[executionID]
		handler.callbacksMu.Unlock()
		require.True(t, exists)
	})
}

func TestIsValidJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected bool
	}{
		{
			name:     "valid JSON object",
			input:    []byte(`{"key": "value"}`),
			expected: true,
		},
		{
			name:     "valid JSON array",
			input:    []byte(`[1, 2, 3]`),
			expected: true,
		},
		{
			name:     "invalid JSON",
			input:    []byte(`{invalid}`),
			expected: false,
		},
		{
			name:     "empty object",
			input:    []byte(`{}`),
			expected: true,
		},
		{
			name:     "null",
			input:    []byte(`null`),
			expected: false,
		},
		{
			name:     "empty string",
			input:    []byte(``),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidJSON(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

// Helper functions

func createTestTriggerRequest() gateway_common.HTTPTriggerRequest {
	return gateway_common.HTTPTriggerRequest{
		Workflow: gateway_common.WorkflowSelector{
			WorkflowID: "test-workflow-id",
		},
		Input: []byte(`{"key": "value"}`),
	}
}

func createTestTriggerRequestWithWorkflowID(workflowID string) gateway_common.HTTPTriggerRequest {
	return gateway_common.HTTPTriggerRequest{
		Workflow: gateway_common.WorkflowSelector{
			WorkflowID: workflowID,
		},
		Input: []byte(`{"key": "value"}`),
	}
}

func createTestTriggerRequestWithInput(input []byte) gateway_common.HTTPTriggerRequest {
	return gateway_common.HTTPTriggerRequest{
		Workflow: gateway_common.WorkflowSelector{
			WorkflowID: "test-workflow-id",
		},
		Input: input,
	}
}

func createTestTriggerHandler(t *testing.T) (*httpTriggerHandler, *handlermocks.DON) {
	cfg := ServiceConfig{
		CleanUpPeriodMs:             60000,
		MaxTriggerRequestDurationMs: 300000,
	}
	return createTestTriggerHandlerWithConfig(t, cfg)
}

func createTestTriggerHandlerWithConfig(t *testing.T, cfg ServiceConfig) (*httpTriggerHandler, *handlermocks.DON) {
	donConfig := &config.DONConfig{
		DonId: "test-don",
		F:     1, // This means we need 2f+1 = 3 responses for consensus
		Members: []config.NodeConfig{
			{Address: "node1"},
			{Address: "node2"},
			{Address: "node3"},
		},
	}
	mockDon := handlermocks.NewDON(t)
	lggr := logger.Test(t)

	handler := NewHTTPTriggerHandler(lggr, cfg, donConfig, mockDon)
	return handler, mockDon
}
