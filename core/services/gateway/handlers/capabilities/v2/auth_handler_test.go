package v2

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	gateway_common "github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/common/aggregation"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/mocks"
)

func createTestAuthHandler(t *testing.T) (*AuthHandler, *mocks.DON, *config.DONConfig) {
	lggr := logger.Test(t)
	mockDon := mocks.NewDON(t)

	donConfig := &config.DONConfig{
		F: 1,
		Members: []config.NodeConfig{
			{Address: "node1"},
			{Address: "node2"},
			{Address: "node3"},
		},
	}

	cfg := WithDefaults(ServiceConfig{})
	handler := NewAuthHandler(lggr, cfg, mockDon, donConfig)
	return handler, mockDon, donConfig
}

func TestSyncAuthorizedKeys(t *testing.T) {
	handler, _, _ := createTestAuthHandler(t)

	// Test when aggregator has no data
	handler.syncAuthorizedKeys()
	require.Empty(t, handler.authorizedKeys)

	// Start the aggregator to enable data collection
	ctx := testutils.Context(t)
	err := handler.agg.Start(ctx)
	require.NoError(t, err)
	defer handler.agg.Close()

	// Add some test data to aggregator
	workflowID := "test-workflow"
	observation := aggregation.WorkflowAuthObservation{
		WorkflowID: workflowID,
		AuthorizedKey: gateway_common.AuthorizedKey{
			PublicKey: "key1",
		},
	}

	// Collect enough observations to meet threshold (F+1 = 2)
	err = handler.agg.Collect(observation, "node1")
	require.NoError(t, err)
	err = handler.agg.Collect(observation, "node2")
	require.NoError(t, err)
	handler.syncAuthorizedKeys()

	workflowKeys, exists := handler.authorizedKeys[workflowID]

	require.True(t, exists)
	require.True(t, workflowKeys.Contains("key1"))
}

func TestSyncAuthorizedKeysMultipleWorkflows(t *testing.T) {
	handler, _, _ := createTestAuthHandler(t)

	ctx := testutils.Context(t)
	err := handler.agg.Start(ctx)
	require.NoError(t, err)
	defer handler.agg.Close()

	// Add observations for multiple workflows
	workflows := []string{"workflow1", "workflow2"}
	keys := []string{"key1", "key2", "key3"}

	for _, workflowID := range workflows {
		for _, key := range keys {
			observation := aggregation.WorkflowAuthObservation{
				WorkflowID: workflowID,
				AuthorizedKey: gateway_common.AuthorizedKey{
					PublicKey: key,
				},
			}
			err = handler.agg.Collect(observation, "node1")
			require.NoError(t, err)
			err = handler.agg.Collect(observation, "node2")
			require.NoError(t, err)
		}
	}
	handler.syncAuthorizedKeys()

	for _, workflowID := range workflows {
		workflowKeys, exists := handler.authorizedKeys[workflowID]
		require.True(t, exists, "Workflow %s should exist", workflowID)

		for _, key := range keys {
			require.True(t, workflowKeys.Contains(key), "Key %s should exist for workflow %s", key, workflowID)
		}
	}
}

func TestSendAuthPullRequest(t *testing.T) {
	handler, mockDon, donConfig := createTestAuthHandler(t)
	ctx := testutils.Context(t)
	for _, member := range donConfig.Members {
		mockDon.EXPECT().SendToNode(ctx, member.Address, mock.Anything).Return(nil).Once()
	}

	err := handler.sendAuthPullRequest(ctx)
	require.NoError(t, err)
	mockDon.AssertExpectations(t)
}

func TestSendAuthPullRequestWithErrors(t *testing.T) {
	handler, mockDon, donConfig := createTestAuthHandler(t)
	ctx := testutils.Context(t)

	// Mock errors for some nodes
	expectedErrors := []error{
		errors.New("connection failed"),
		nil,
		errors.New("timeout"),
	}

	for i, member := range donConfig.Members {
		mockDon.EXPECT().SendToNode(ctx, member.Address, mock.Anything).Return(expectedErrors[i]).Once()
	}

	err := handler.sendAuthPullRequest(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "connection failed")
	require.Contains(t, err.Error(), "timeout")
	require.NotContains(t, err.Error(), "node2")
	mockDon.AssertExpectations(t)
}

func TestSendAuthPullRequestVerifyPayload(t *testing.T) {
	handler, mockDon, donConfig := createTestAuthHandler(t)
	ctx := testutils.Context(t)
	// Capture the request payload
	var capturedReq *jsonrpc.Request[json.RawMessage]
	mockDon.On("SendToNode", ctx, mock.AnythingOfType("string"), mock.Anything).
		Run(func(args mock.Arguments) {
			capturedReq = args.Get(2).(*jsonrpc.Request[json.RawMessage])
		}).Return(nil)

	err := handler.sendAuthPullRequest(ctx)
	require.NoError(t, err)

	require.Equal(t, jsonrpc.JsonRpcVersion, capturedReq.Version)
	require.Equal(t, gateway.MethodWorkflowPullAuthMetadata, capturedReq.Method)
	require.NotEmpty(t, capturedReq.ID)

	mockDon.AssertNumberOfCalls(t, "SendToNode", len(donConfig.Members))
}

func TestOnAuthMetadataPush(t *testing.T) {
	handler, _, _ := createTestAuthHandler(t)
	ctx := testutils.Context(t)

	err := handler.agg.Start(ctx)
	require.NoError(t, err)
	defer handler.agg.Close()

	authData := gateway.WorkflowAuthMetadata{
		WorkflowID: "test-workflow",
		AuthorizedKeys: []gateway_common.AuthorizedKey{
			{PublicKey: "key1"},
			{PublicKey: "key2"},
		},
	}

	result, err := json.Marshal(authData)
	require.NoError(t, err)

	rawResult := json.RawMessage(result)
	resp := &jsonrpc.Response[json.RawMessage]{
		Result: &rawResult,
	}

	err = handler.OnAuthMetadataPush(ctx, resp, "node1")
	require.NoError(t, err)
}

func TestOnAuthMetadataPushInvalidJSON(t *testing.T) {
	handler, _, _ := createTestAuthHandler(t)
	ctx := testutils.Context(t)

	invalidJSON := json.RawMessage(`{"invalid": json}`)
	resp := &jsonrpc.Response[json.RawMessage]{
		Result: &invalidJSON,
	}

	err := handler.OnAuthMetadataPush(ctx, resp, "node1")
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to unmarshal auth metadata")
}

func TestOnAuthMetadataPullResponse(t *testing.T) {
	handler, _, _ := createTestAuthHandler(t)
	ctx := testutils.Context(t)

	err := handler.agg.Start(ctx)
	require.NoError(t, err)
	defer handler.agg.Close()

	authData := []gateway.WorkflowAuthMetadata{
		{
			WorkflowID: "workflow1",
			AuthorizedKeys: []gateway_common.AuthorizedKey{
				{PublicKey: "key1"},
			},
		},
		{
			WorkflowID: "workflow2",
			AuthorizedKeys: []gateway_common.AuthorizedKey{
				{PublicKey: "key2"},
				{PublicKey: "key3"},
			},
		},
	}

	result, err := json.Marshal(authData)
	require.NoError(t, err)

	rawResult := json.RawMessage(result)
	resp := &jsonrpc.Response[json.RawMessage]{
		Result: &rawResult,
	}

	err = handler.OnAuthMetadataPullResponse(ctx, resp, "node1")
	require.NoError(t, err)
}

func TestOnAuthMetadataPullResponseInvalidJSON(t *testing.T) {
	handler, _, _ := createTestAuthHandler(t)
	ctx := testutils.Context(t)

	invalidJSON := json.RawMessage(`[{"invalid": json}]`)
	resp := &jsonrpc.Response[json.RawMessage]{
		Result: &invalidJSON,
	}

	err := handler.OnAuthMetadataPullResponse(ctx, resp, "node1")
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to unmarshal auth metadata pull response")
}

func TestStartAndClose(t *testing.T) {
	handler, _, _ := createTestAuthHandler(t)
	ctx := testutils.Context(t)

	err := handler.Start(ctx)
	require.NoError(t, err)
	require.NoError(t, handler.Ready())
	err = handler.Start(ctx) // Should error on second start
	require.Error(t, err)

	err = handler.Close()
	require.NoError(t, err)
	require.Error(t, handler.Ready())
	err = handler.Close() // Should error on second close
	require.Error(t, err)
}

func TestAuthHandlerIntegration(t *testing.T) {
	handler, mockDon, _ := createTestAuthHandler(t)
	ctx := testutils.Context(t)

	// Mock successful DON interactions
	mockDon.On("SendToNode", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*jsonrpc.Request[json.RawMessage]")).
		Return(nil).Maybe()

	err := handler.Start(ctx)
	require.NoError(t, err)
	defer handler.Close()

	// Simulate receiving auth metadata push
	authData := gateway.WorkflowAuthMetadata{
		WorkflowID: "integration-workflow",
		AuthorizedKeys: []gateway_common.AuthorizedKey{
			{PublicKey: "integration-key1"},
			{PublicKey: "integration-key2"},
		},
	}

	result, err := json.Marshal(authData)
	require.NoError(t, err)

	rawResult := json.RawMessage(result)
	resp := &jsonrpc.Response[json.RawMessage]{
		Result: &rawResult,
	}

	// Push from multiple nodes to meet threshold
	err = handler.OnAuthMetadataPush(ctx, resp, "node1")
	require.NoError(t, err)
	err = handler.OnAuthMetadataPush(ctx, resp, "node2")
	require.NoError(t, err)
	time.Sleep(100 * time.Millisecond)
	handler.syncAuthorizedKeys()

	// Verify keys are available for authorization
	workflowKeys, exists := handler.authorizedKeys["integration-workflow"]

	require.True(t, exists)
	require.True(t, workflowKeys.Contains("integration-key1"))
	require.True(t, workflowKeys.Contains("integration-key2"))
}
