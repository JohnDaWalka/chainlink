package v2

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	"go.uber.org/multierr"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	gateway_common "github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

var _ HTTPTriggerHandler = (*httpTriggerHandler)(nil)

type savedCallback struct {
	callbackCh         chan<- handlers.UserCallbackPayload
	createdAt          time.Time
	responseAggregator common.NodeResponseAggregator
}

type httpTriggerHandler struct {
	services.StateMachine
	config      ServiceConfig
	don         handlers.DON
	donConfig   *config.DONConfig
	lggr        logger.Logger
	callbacksMu sync.Mutex
	callbacks   map[string]savedCallback // requestID -> savedCallback
	stopCh      services.StopChan
}

type HTTPTriggerHandler interface {
	job.ServiceCtx
	HandleUserTriggerRequest(ctx context.Context, req *jsonrpc.Request, callbackCh chan<- handlers.UserCallbackPayload) error
	HandleNodeTriggerResponse(ctx context.Context, resp *jsonrpc.Response, nodeAddr string) error
}

func NewHTTPTriggerHandler(lggr logger.Logger, cfg ServiceConfig, donConfig *config.DONConfig, don handlers.DON) *httpTriggerHandler {
	return &httpTriggerHandler{
		lggr:      logger.Named(lggr, "RequestCallbacks"),
		callbacks: make(map[string]savedCallback),
		config:    cfg,
		don:       don,
		donConfig: donConfig,
	}
}

func (h *httpTriggerHandler) HandleUserTriggerRequest(ctx context.Context, req *jsonrpc.Request, callbackCh chan<- handlers.UserCallbackPayload) error {
	// TODO: PRODCRE-305 validate JWT against authorized keys
	triggerReq, err := h.validatedTriggerRequest(req, callbackCh)
	if err != nil {
		return err
	}
	// TODO: PRODCRE-475 support look-up of workflowID using workflowOwner/Label/Name. Rate-limiting using workflowOwner
	workflowID := triggerReq.Workflow.WorkflowID
	executionID, err := workflows.EncodeExecutionID(workflowID, req.ID)
	if err != nil {
		h.lggr.Errorw("error generating execution ID", "err", err)
		h.handleUserError(req.ID, jsonrpc.ErrInternal, internalErrorMessage, callbackCh)
		return nil
	}
	h.callbacksMu.Lock()
	_, found := h.callbacks[executionID]
	if found {
		h.callbacksMu.Unlock()
		h.handleUserError(req.ID, jsonrpc.ErrInvalidRequest, "request ID already used: "+req.ID, callbackCh)
		return errors.New("request ID already used: " + req.ID)
	}
	h.callbacks[executionID] = savedCallback{
		callbackCh:         callbackCh,
		createdAt:          time.Now(),
		responseAggregator: common.NewIdenticalNodeResponseAggregator(2*h.donConfig.F + 1),
	}
	h.callbacksMu.Unlock()
	// Send original request to all nodes
	for _, member := range h.donConfig.Members {
		err = multierr.Combine(err, h.don.SendToNode(ctx, member.Address, req))
	}
	return err
}

func (h *httpTriggerHandler) validatedTriggerRequest(req *jsonrpc.Request, callbackCh chan<- handlers.UserCallbackPayload) (*gateway_common.HTTPTriggerRequest, error) {
	var triggerReq gateway_common.HTTPTriggerRequest
	err := json.Unmarshal(req.Params, &triggerReq)
	if err != nil {
		h.handleUserError(req.ID, jsonrpc.ErrParse, "error decoding payload: "+err.Error(), callbackCh)
		return nil, err
	}
	if req.ID == "" {
		h.handleUserError(req.ID, jsonrpc.ErrInvalidParams, "empty request ID", callbackCh)
		return nil, errors.New("empty request ID")
	}
	// Request IDs from users must not contain "/", since this character is reserved
	// for internal node-to-node message routing (e.g., "http_action/{workflowID}/{uuid}").
	if strings.Contains(req.ID, "/") {
		h.handleUserError(req.ID, jsonrpc.ErrInvalidParams, "request ID must not contain '/'", callbackCh)
		return nil, errors.New("request ID must not contain '/'")
	}
	if req.Method != gateway_common.MethodWorkflowExecute {
		h.handleUserError(req.ID, jsonrpc.ErrMethodNotFound, "invalid method: "+req.Method, callbackCh)
		return nil, errors.New("invalid method: " + req.Method)
	}
	if !isValidJSON(triggerReq.Input) {
		h.lggr.Errorw("invalid params JSON", "params", triggerReq.Input)
		h.handleUserError(req.ID, jsonrpc.ErrInvalidRequest, "invalid params JSON", callbackCh)
		return nil, errors.New("invalid params JSON")
	}
	return &triggerReq, nil
}

func (h *httpTriggerHandler) HandleNodeTriggerResponse(ctx context.Context, resp *jsonrpc.Response, nodeAddr string) error {
	h.lggr.Debugw("handling trigger response", "requestID", resp.ID, "nodeAddr", nodeAddr)
	h.callbacksMu.Lock()
	defer h.callbacksMu.Unlock()
	saved, exists := h.callbacks[resp.ID]
	if !exists {
		return errors.New("callback not found for request ID: " + resp.ID)
	}
	resp, err := saved.responseAggregator.CollectAndAggregate(resp, nodeAddr)
	if err != nil {
		return err
	}
	if resp == nil {
		h.lggr.Debugw("Not enough responses to aggregate", "requestID", resp.ID, "nodeAddress", nodeAddr)
		return nil
	}
	rawResp, err := json.Marshal(resp)
	if err != nil {
		return errors.New("failed to marshal response: " + err.Error())
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case saved.callbackCh <- handlers.UserCallbackPayload{
		RawResponse: rawResp,
		ErrorCode:   api.NoError,
	}:
		close(saved.callbackCh)
	}
	return nil
}

func (h *httpTriggerHandler) Start(ctx context.Context) error {
	return h.StartOnce("HTTPTriggerHandler", func() error {
		h.lggr.Info("Starting HTTPTriggerHandler")
		go func() {
			ticker := time.NewTicker(time.Duration(h.config.CleanUpPeriodMs) * time.Millisecond)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					h.reapExpiredCallbacks()
				case <-h.stopCh:
					return
				}
			}
		}()
		return nil
	})
}

func (h *httpTriggerHandler) Close() error {
	return h.StopOnce("HTTPTriggerHandler", func() error {
		h.lggr.Info("Closing HTTPTriggerHandler")
		close(h.stopCh)
		return nil
	})
}

// reapExpiredCallbacks removes callbacks that are older than the maximum age
func (h *httpTriggerHandler) reapExpiredCallbacks() {
	h.callbacksMu.Lock()
	defer h.callbacksMu.Unlock()
	now := time.Now()
	var expiredCount int
	for executionID, callback := range h.callbacks {
		if now.Sub(callback.createdAt) > time.Duration(h.config.MaxTriggerRequestDurationMs)*time.Millisecond {
			delete(h.callbacks, executionID)
			expiredCount++
		}
	}
	if expiredCount > 0 {
		h.lggr.Infow("Removed expired callbacks", "count", expiredCount, "remaining", len(h.callbacks))
	}
}

func isValidJSON(data []byte) bool {
	var val map[string]interface{}
	return json.Unmarshal(data, &val) != nil
}

func (h *httpTriggerHandler) handleUserError(requestID string, code int64, message string, callbackCh chan<- handlers.UserCallbackPayload) {
	resp := &jsonrpc.Response{
		Version: "2.0",
		ID:      requestID,
		Error: &jsonrpc.WireError{
			Code:    code,
			Message: message,
		},
	}
	rawResp, err := json.Marshal(resp)
	if err != nil {
		h.lggr.Errorw("failed to marshal error response", "err", err, "requestID", requestID)
		return
	}
	callbackCh <- handlers.UserCallbackPayload{
		RawResponse: rawResp,
		ErrorCode:   api.ErrorCode(code),
	}
}
