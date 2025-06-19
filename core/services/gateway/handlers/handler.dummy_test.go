package handlers_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	hc "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
)

type testConnManager struct {
	handler     handlers.Handler
	sendCounter int
}

func (m *testConnManager) SetHandler(handler handlers.Handler) {
	m.handler = handler
}

func (m *testConnManager) SendToNode(ctx context.Context, nodeAddress string, resp *jsonrpc.Request) error {
	m.sendCounter++
	return nil
}

func TestDummyHandler_BasicFlow(t *testing.T) {
	t.Parallel()

	config := config.DONConfig{
		Members: []config.NodeConfig{
			{Name: "node one", Address: "addr_1"},
			{Name: "node two", Address: "addr_2"},
		},
	}

	connMgr := testConnManager{}
	handler, err := handlers.NewDummyHandler(&config, &connMgr, logger.Test(t))
	require.NoError(t, err)
	connMgr.SetHandler(handler)

	ctx := testutils.Context(t)

	// User request
	msg := api.Message{Body: api.MessageBody{MessageId: "1234"}}
	callbackCh := make(chan handlers.UserCallbackPayload, 1)
	require.NoError(t, handler.HandleUserMessage(ctx, &msg, callbackCh))
	require.Equal(t, 2, connMgr.sendCounter)

	// Responses from both nodes
	resp, err := hc.ValidatedResponseFromMessage(&msg)
	require.NoError(t, err)
	require.NoError(t, handler.HandleNodeMessage(ctx, resp, "addr_1"))
	require.NoError(t, handler.HandleNodeMessage(ctx, resp, "addr_2"))
	response := <-callbackCh
	require.Equal(t, "1234", response.Msg.Body.MessageId)
}
