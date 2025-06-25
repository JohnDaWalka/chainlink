package v2

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"go.uber.org/multierr"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/ratelimit"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	gateway_common "github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
)

var _ jsonrpc.Service = (*service)(nil)
var _ NodeMessageHandler = (*service)(nil)

const (
	ServiceName = "HTTPService"
)

const (
	ErrInvalidRequest          int64  = -32600 // InvalidRequest: invalid JSON, not retryable
	ValidationErrorCode        int64  = -32602 // ValidationError: invalid fields, not retryable
	MethodNotFoundErrorCode    int64  = -32601 // MethodNotFound: method does not exist, not retryable
	UnauthorizedErrorCode      int64  = -32001 // Unauthorized: invalid/missing JWT, not retryable
	ResourceNotFoundCode       int64  = -32004 // ResourceNotFound: workflowID does not exist, not retryable
	LimitExceededErrorCode     int64  = -32029 // LimitExceeded: rate limit exceeded, retryable
	InternalServerErrorCode    int64  = -32603 // InternalServerError: unexpected/unhandled error, retryable
	InternalServerErrorMessage string = "Internal Server Error"
)

// NodeMessageHandler implements service-specific logic for managing messages from nodes.
type NodeMessageHandler interface {
	// Handlers should not make any assumptions about goroutines calling HandleNodeMessage.
	// should be non-blocking
	HandleNodeMessage(ctx context.Context, req *jsonrpc.Response, nodeAddr string) error
}

type service struct {
	services.StateMachine
	config             ServiceConfig
	don                handlers.DON
	donConfig          *config.DONConfig
	savedCallbacks     map[string]*savedCallback
	mu                 sync.Mutex
	lggr               logger.Logger
	httpClient         network.HTTPClient
	nodeRateLimiter    *ratelimit.RateLimiter
	userRateLimiter    *ratelimit.RateLimiter
	wg                 sync.WaitGroup
	stopCh             services.StopChan
	responseAggregator common.NodeResponseAggregator
}

type ServiceConfig struct {
	NodeRateLimiter         ratelimit.RateLimiterConfig `json:"nodeRateLimiter"`
	UserRateLimiter         ratelimit.RateLimiterConfig `json:"userRateLimiter"`
	ResponseMaxAgeMs        int                         `json:"callbackMaxAgeMs"`
	ResponseCleanUpPeriodMs int                         `json:"callbackCleanUpPeriodMs"`
}
type savedCallback struct {
	callbackCh chan<- *jsonrpc.Response
	createdAt  time.Time
}

func NewHandler(handlerConfig json.RawMessage, donConfig *config.DONConfig, don handlers.DON, httpClient network.HTTPClient, lggr logger.Logger) (*service, error) {
	var cfg ServiceConfig
	err := json.Unmarshal(handlerConfig, &cfg)
	if err != nil {
		return nil, err
	}
	nodeRateLimiter, err := ratelimit.NewRateLimiter(cfg.NodeRateLimiter)
	if err != nil {
		return nil, err
	}
	userRateLimiter, err := ratelimit.NewRateLimiter(cfg.UserRateLimiter)
	if err != nil {
		return nil, err
	}
	responseAggregator := common.NewIdenticalNodeResponseAggregator(lggr, donConfig.F,
		cfg.ResponseMaxAgeMs, cfg.ResponseCleanUpPeriodMs)
	return &service{
		config:             cfg,
		don:                don,
		donConfig:          donConfig,
		lggr:               logger.Named(lggr, ServiceName+donConfig.DonId),
		httpClient:         httpClient,
		nodeRateLimiter:    nodeRateLimiter,
		userRateLimiter:    userRateLimiter,
		wg:                 sync.WaitGroup{},
		savedCallbacks:     make(map[string]*savedCallback),
		stopCh:             make(services.StopChan),
		responseAggregator: responseAggregator,
	}, nil
}

func (h *service) handleOutgoingRequest(ctx context.Context, requestID string, req gateway.OutboundHTTPRequest, nodeAddr string) error {
	h.lggr.Debugw("handling webAPI outgoing message", "requestID", requestID, "nodeAddr", nodeAddr)
	if !h.nodeRateLimiter.Allow(nodeAddr) {
		return fmt.Errorf("rate limit exceeded for node %s", nodeAddr)
	}
	timeout := time.Duration(req.TimeoutMs) * time.Millisecond
	httpReq := network.HTTPRequest{
		Method:           req.Method,
		URL:              req.URL,
		Headers:          req.Headers,
		Body:             req.Body,
		MaxResponseBytes: req.MaxResponseBytes,
		Timeout:          timeout,
	}

	// send response to node async
	h.wg.Add(1)
	go func() {
		defer h.wg.Done()
		// not cancelled when parent is cancelled to ensure the goroutine can finish
		newCtx := context.WithoutCancel(ctx)
		newCtx, cancel := context.WithTimeout(newCtx, timeout)
		defer cancel()
		l := logger.With(h.lggr, "url", req.URL, "requestID", requestID, "method", req.Method, "timeout", req.TimeoutMs)
		l.Debug("Sending request to client")
		var outboundResp gateway.OutboundHTTPResponse
		resp, err := h.httpClient.Send(ctx, httpReq)
		if err != nil {
			l.Errorw("error while sending HTTP request to external endpoint", "err", err)
			outboundResp = gateway.OutboundHTTPResponse{
				ErrorMessage: err.Error(),
			}
		} else {
			outboundResp = gateway.OutboundHTTPResponse{
				StatusCode: resp.StatusCode,
				Headers:    resp.Headers,
				Body:       resp.Body,
			}
		}
		params, err := json.Marshal(outboundResp)
		if err != nil {
			l.Errorw("failed to marshal HTTP response", "err", err)
			return
		}
		req := &jsonrpc.Request{
			Version: "2.0",
			ID:      requestID,
			Method:  gateway_common.MethodHTTPAction,
			Params:  params,
		}
		err = h.don.SendToNode(newCtx, nodeAddr, req)
		if err != nil {
			l.Errorw("failed to send to node", "err", err, "to", nodeAddr)
			return
		}
		l.Debugw("sent response to node", "to", nodeAddr)
	}()
	return nil
}

func (h *service) handleTriggerResponse(ctx context.Context, resp *jsonrpc.Response, nodeAddr string) error {
	h.lggr.Debugw("handling trigger response", "requestID", resp.ID, "nodeAddr", nodeAddr)
	h.mu.Lock()
	defer h.mu.Unlock()
	callback, exists := h.savedCallbacks[resp.ID]
	if !exists {
		h.lggr.Errorw("no callback found for request ID", "requestID", resp.ID, "nodeAddr", nodeAddr)
		return nil
	}
	resp, err := h.responseAggregator.CollectAndAggregate(resp.ID, resp, nodeAddr)
	if err != nil {
		h.lggr.Debugw("insufficient number of responses", "requestID", resp.ID, "nodeAddr", nodeAddr, "err", err)
	}
	select {
	case <-ctx.Done():
		h.lggr.Warnw("context cancelled while sending response to user", "requestID", resp.ID, "nodeAddr", nodeAddr)
		return ctx.Err()
	case callback.callbackCh <- resp:
		close(callback.callbackCh)
	}
	return nil
}

// TODO: can this return error?
func (h *service) HandleNodeMessage(ctx context.Context, resp *jsonrpc.Response, nodeAddr string) error {
	requestID := resp.ID
	if requestID == "" {
		return fmt.Errorf("received response with empty request ID from node %s", nodeAddr)
	}
	h.lggr.Debugw("handling incoming node message", "requestID", requestID, "nodeAddr", nodeAddr)
	h.mu.Lock()
	_, exists := h.savedCallbacks[requestID]
	h.mu.Unlock()
	if exists {
		err := h.handleTriggerResponse(ctx, resp, nodeAddr)
		if err != nil {
			return err
		}
	}
	var outboundReq gateway.OutboundHTTPRequest
	err := json.Unmarshal(resp.Result, &outboundReq)
	if err != nil {
		return fmt.Errorf("failed to unmarshal HTTP request from node %s: %w", nodeAddr, err)
	}
	return h.handleOutgoingRequest(ctx, requestID, outboundReq, nodeAddr)
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

func (h *service) validatedTriggerRequest(req *jsonrpc.Request, callbackCh chan<- *jsonrpc.Response) *gateway_common.HTTPTriggerRequest {
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

func (h *service) HandleUserRequest(ctx context.Context, req *jsonrpc.Request, callbackCh chan<- *jsonrpc.Response) error {
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

// reapExpiredCallbacks removes callbacks that are older than the maximum age
func (h *service) reapExpiredCallbacks() {
	h.mu.Lock()
	defer h.mu.Unlock()
	now := time.Now()
	var expiredCount int
	for executionID, callback := range h.savedCallbacks {
		if now.Sub(callback.createdAt) > time.Duration(h.config.ResponseMaxAgeMs)*time.Millisecond {
			delete(h.savedCallbacks, executionID)
			expiredCount++
		}
	}
	if expiredCount > 0 {
		h.lggr.Infow("Removed expired callbacks", "count", expiredCount, "remaining", len(h.savedCallbacks))
	}
}

func isValidJSON(data []byte) bool {
	var val map[string]interface{}
	return json.Unmarshal(data, &val) == nil
}

func generateExecutionID(workflowID, triggerEventID string) (string, error) {
	s := sha256.New()
	_, err := s.Write([]byte(workflowID))
	if err != nil {
		return "", err
	}

	_, err = s.Write([]byte(triggerEventID))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(s.Sum(nil)), nil
}

func (h *service) HealthReport() map[string]error {
	return map[string]error{ServiceName: h.Healthy()}
}

func (h *service) Name() string {
	return ServiceName
}

func (h *service) Start(ctx context.Context) error {
	return h.StartOnce(ServiceName, func() error {
		h.lggr.Info("Starting " + ServiceName)
		err := h.responseAggregator.Start(ctx)
		if err != nil {
			return err
		}
		// Start a goroutine to periodically clean up expired callbacks
		go func() {
			ticker := time.NewTicker(time.Duration(h.config.ResponseCleanUpPeriodMs) * time.Millisecond)
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

func (h *service) Close() error {
	return h.StopOnce(ServiceName, func() error {
		h.lggr.Info("Closing " + ServiceName)
		err := h.responseAggregator.Close()
		if err != nil {
			h.lggr.Errorw("error closing response aggregator", "err", err)
		}
		close(h.stopCh)
		h.wg.Wait()
		return nil
	})
}
