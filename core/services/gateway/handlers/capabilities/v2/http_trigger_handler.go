package v2

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"go.uber.org/multierr"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	gateway_common "github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

var _ HTTPTriggerHandler = (*httpTriggerHandler)(nil)

type savedCallback struct {
	callbackCh chan<- *jsonrpc.Response
	createdAt  time.Time
}

type httpTriggerHandler struct {
	lggr        logger.Logger
	callbacksMu sync.Mutex
	callbacks   map[string]savedCallback // requestID -> savedCallback
}

type HTTPTriggerHandler interface {
	job.ServiceCtx
}

func NewHTTPTriggerHandler(lggr logger.Logger) *httpTriggerHandler {
	return &httpTriggerHandler{
		lggr:      logger.Named(lggr, "RequestCallbacks"),
		callbacks: make(map[string]savedCallback),
	}
}

func (h *httpTriggerHandler) HandleUserTriggerRequest(ctx context.Context, req *jsonrpc.Request, callbackCh chan<- *jsonrpc.Response) error {
	// TODO: is returning nil here correct?
	// TODO: PRODCRE-305 validate JWT against authorized keys
	triggerReq := h.validatedTriggerRequest(req, callbackCh)
	if triggerReq == nil {
		// Error already handled in validatedTriggerRequest
		return nil
	}
	// TODO: PRODCRE-475 support look-up of workflowID using workflowOwner/Label/Name. Rate-limiting using workflowOwner
	workflowID := triggerReq.Workflow.WorkflowID
	executionID, err := generateExecutionID(workflowID, req.ID)
	if err != nil {
		h.lggr.Errorw("error generating execution ID", "err", err)
		handleUserError(req.ID, InternalServerErrorCode, InternalServerErrorMessage, callbackCh)
		return nil
	}
	h.mu.Lock()
	_, found := h.savedCallbacks[executionID]
	if found {
		h.mu.Unlock()
		h.lggr.Debugw("callback already exists for execution ID", "executionID", executionID)
		handleUserError(req.ID, ValidationErrorCode, "request ID already used: "+req.ID, callbackCh)
	}
	h.savedCallbacks[executionID] = &savedCallback{
		callbackCh: callbackCh,
		createdAt:  time.Now(),
	}
	h.mu.Unlock()
	// Send original request to all nodes
	for _, member := range h.donConfig.Members {
		err = multierr.Combine(err, h.don.SendToNode(ctx, member.Address, req))
	}
	return err
}

func (h *httpTriggerHandler) validatedTriggerRequest(req *jsonrpc.Request, callbackCh chan<- *jsonrpc.Response) *gateway_common.HTTPTriggerRequest {
	var triggerReq gateway_common.HTTPTriggerRequest
	err := json.Unmarshal(req.Params, &triggerReq)
	if err != nil {
		h.lggr.Errorw("error decoding payload", "err", err)
		handleUserError(req.ID, ErrInvalidRequest, "error decoding payload: "+err.Error(), callbackCh)
		return nil
	}
	if req.ID == "" {
		h.lggr.Errorw("empty request ID", "method", req.Method)
		handleUserError(req.ID, ErrInvalidRequest, "empty request ID", callbackCh)
		return nil
	}
	if req.Method != gateway_common.MethodWorkflowExecute {
		h.lggr.Errorw("invalid method", "method", req.Method)
		handleUserError(req.ID, MethodNotFoundErrorCode, "invalid method: "+req.Method, callbackCh)
		return nil
	}
	if isValidJSON(triggerReq.Input) {
		h.lggr.Errorw("invalid params JSON", "params", triggerReq.Input)
		handleUserError(req.ID, ValidationErrorCode, "invalid params JSON", callbackCh)
		return nil
	}
	return &triggerReq
}

func (h *gatewayHandler) HandleNodeTriggerResponse(ctx context.Context, resp *jsonrpc.Response, nodeAddr string) error {
	h.lggr.Debugw("handling trigger response", "requestID", resp.ID, "nodeAddr", nodeAddr)
	resp, err := h.responseAggregator.CollectAndAggregate(resp.ID, resp, nodeAddr)
	if err != nil {
		h.lggr.Debugw("insufficient number of responses", "requestID", resp.ID, "nodeAddr", nodeAddr, "err", err)
		return nil
	}
	err = h.triggerCallbacks.SendResponse(ctx, resp.ID, resp)
	if err != nil {
		h.lggr.Errorw("error sending trigger response", "err", err, "requestID", resp.ID, "nodeAddr", nodeAddr)
		handleUserError(resp.ID, InternalServerErrorCode, InternalServerErrorMessage, resp.CallbackCh)
	}
	return nil
}

func isValidJSON(data []byte) bool {
	var val map[string]interface{}
	return json.Unmarshal(data, &val) == nil
}

func handleUserError(requestID string, code int64, message string, callbackCh chan<- *jsonrpc.Response) {
	callbackCh <- &jsonrpc.Response{
		Version: "2.0",
		ID:      requestID,
		Error: &jsonrpc.WireError{
			Code:    code,
			Message: message,
		},
	}
	close(callbackCh)
}
