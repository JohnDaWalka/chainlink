package v2

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
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
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
)

var _ handlers.Handler = (*gatewayHandler)(nil)

const (
	handlerName            = "HTTPCapabilityHandler"
	defaultCleanUpPeriodMs = 1000 * 60 * 10 // 10 minutes
)

type savedCallback struct {
	callbackCh chan<- *jsonrpc.Response
	createdAt  time.Time
}
type gatewayHandler struct {
	services.StateMachine
	config          ServiceConfig
	don             handlers.DON
	donConfig       *config.DONConfig
	lggr            logger.Logger
	httpClient      network.HTTPClient
	nodeRateLimiter *ratelimit.RateLimiter // Rate limiter for node requests (e.g. outgoing HTTP requests, HTTP trigger response, auth metadata exchange)
	userRateLimiter *ratelimit.RateLimiter // Rate limiter for user requests that trigger workflow executions
	wg              sync.WaitGroup
	stopCh          services.StopChan
	responseCache   ResponseCache
	responseAggregator common.NodeResponseAggregator
}

type ResponseCache interface {
	Set(req gateway_common.OutboundHTTPRequest, response gateway_common.OutboundHTTPResponse, ttl time.Duration)
	Get(req gateway_common.OutboundHTTPRequest) *gateway_common.OutboundHTTPResponse
	DeleteExpired() int
}

type ServiceConfig struct {
	NodeRateLimiter ratelimit.RateLimiterConfig `json:"nodeRateLimiter"`
	UserRateLimiter ratelimit.RateLimiterConfig `json:"userRateLimiter"`
	TriggerRequestMaxDurationMs        int                         `json:"triggerRequestMaxDurationMs"`
	CleanUpPeriodMs int                         `json:"cleanUpPeriodMs"`
}

func NewGatewayHandler(handlerConfig json.RawMessage, donConfig *config.DONConfig, don handlers.DON, httpClient network.HTTPClient, lggr logger.Logger) (*gatewayHandler, error) {
	var cfg ServiceConfig
	err := json.Unmarshal(handlerConfig, &cfg)
	if err != nil {
		return nil, err
	}
	cfg = WithDefaults(cfg)
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
	return &gatewayHandler{
		config:          cfg,
		don:             don,
		donConfig:       donConfig,
		lggr:            logger.With(logger.Named(lggr, handlerName), "donId", donConfig.DonId),
		httpClient:      httpClient,
		nodeRateLimiter: nodeRateLimiter,
		userRateLimiter: userRateLimiter,
		wg:              sync.WaitGroup{},
		stopCh:          make(services.StopChan),
		responseCache:   newResponseCache(lggr),
		callbacks: newCallbacks(lggr),
	}, nil
}

func WithDefaults(cfg ServiceConfig) ServiceConfig {
	if cfg.CleanUpPeriodMs == 0 {
		cfg.CleanUpPeriodMs = defaultCleanUpPeriodMs
	}
	// TODO: 
	return cfg
}

func (h *gatewayHandler) HandleNodeMessage(ctx context.Context, resp *jsonrpc.Response, nodeAddr string) error {
	if resp.ID == "" {
		return fmt.Errorf("received response with empty request ID from node %s", nodeAddr)
	}
	h.lggr.Debugw("handling incoming node message", "requestID", resp.ID, "nodeAddr", nodeAddr)
	var outboundReq gateway_common.OutboundHTTPRequest
	err := json.Unmarshal(resp.Result, &outboundReq)
	if err != nil {
		return fmt.Errorf("failed to unmarshal HTTP request from node %s: %w", nodeAddr, err)
	}
	return h.handleOutgoingRequest(ctx, resp.ID, outboundReq, nodeAddr)
}

func (h *gatewayHandler) HandleLegacyUserMessage(context.Context, *api.Message, chan<- handlers.UserCallbackPayload) error {
	return errors.New("HTTP capability gateway handler does not support legacy messages")
}

func (h *gatewayHandler) HandleJSONRPCUserMessage(context.Context, jsonrpc.Request, chan<- handlers.UserCallbackPayload) error {
	// TODO: Implement trigger request handling
	return nil
}

func (h *gatewayHandler) handleOutgoingRequest(ctx context.Context, requestID string, req gateway_common.OutboundHTTPRequest, nodeAddr string) error {
	h.lggr.Debugw("handling webAPI outgoing message", "requestID", requestID, "nodeAddr", nodeAddr)
	if !h.nodeRateLimiter.Allow(nodeAddr) {
		return fmt.Errorf("rate limit exceeded for node %s", nodeAddr)
	}
	if req.CacheSettings.ReadFromCache {
		cached := h.responseCache.Get(req)
		if cached != nil {
			h.lggr.Debugw("Using cached HTTP response", "requestID", requestID, "nodeAddr", nodeAddr)
			return h.sendResponseToNode(ctx, requestID, *cached, nodeAddr)
		}
	}

	timeout := time.Duration(req.TimeoutMs) * time.Millisecond
	httpReq := network.HTTPRequest{
		Method:           req.Method,
		URL:              req.URL,
		Headers:          req.Headers,
		Body:             req.Body,
		MaxResponseBytes: req.MaxResponseBytes,
		Timeout:          timeout,

	// send response to node async
	h.wg.Add(1)
	go func() {
		defer h.wg.Done()
		// not cancelled when parent is cancelled to ensure the goroutine can finish
		newCtx := context.WithoutCancel(ctx)
		newCtx, cancel := context.WithTimeout(newCtx, timeout)
		defer cancel()
		l := logger.With(h.lggr, "requestID", requestID, "method", req.Method, "timeout", req.TimeoutMs)
		l.Debug("Sending request to client")
		var outboundResp gateway_common.OutboundHTTPResponse
		resp, err := h.httpClient.Send(ctx, httpReq)
		if err != nil {
			l.Errorw("error while sending HTTP request to external endpoint", "err", err)
			outboundResp = gateway_common.OutboundHTTPResponse{
				ErrorMessage: err.Error(),
			}
		} else {
			outboundResp = gateway_common.OutboundHTTPResponse{
				StatusCode: resp.StatusCode,
				Headers:    resp.Headers,
				Body:       resp.Body,
			}
			if req.CacheSettings.StoreInCache && isCacheableStatusCode(resp.StatusCode) {
				cacheTTLMs := req.CacheSettings.TTLMs
				if cacheTTLMs > 0 {
					h.responseCache.Set(req, outboundResp, time.Duration(cacheTTLMs)*time.Millisecond)
					l.Debugw("Cached HTTP response", "ttlMs", cacheTTLMs)
				}
			}
		}
		err = h.sendResponseToNode(newCtx, requestID, outboundResp, nodeAddr)
		if err != nil {
			l.Errorw("error sending response to node", "err", err, "nodeAddr", nodeAddr)
		}
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

func (h *gatewayHandler) HealthReport() map[string]error {
	return map[string]error{handlerName: h.Healthy()}
}

func (h *gatewayHandler) Name() string {
	return handlerName
}

func (h *gatewayHandler) Start(context.Context) error {
	return h.StartOnce(handlerName, func() error {
		h.lggr.Info("Starting " + handlerName)
		err := h.responseAggregator.Start(ctx)
		if err != nil {
			return err
		}
		go func() {
			ticker := time.NewTicker(time.Duration(h.config.CleanUpPeriodMs) * time.Millisecond)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					TODO: 
					// h.reapExpiredCallbacks()
					h.responseCache.DeleteExpired()
				case <-h.stopCh:
					return
				}
			}
		}()
		return nil
	})
}

func (h *gatewayHandler) Close() error {
	return h.StopOnce(handlerName, func() error {
		h.lggr.Info("Closing " + handlerName)
		err := h.responseAggregator.Close()
		close(h.stopCh)
		h.wg.Wait()
		return nil
	})
}

func (h *gatewayHandler) sendResponseToNode(ctx context.Context, requestID string, resp gateway_common.OutboundHTTPResponse, nodeAddr string) error {
	params, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	req := &jsonrpc.Request{
		Version: "2.0",
		ID:      requestID,
		Method:  gateway_common.MethodHTTPAction,
		Params:  params,
	}

	err = h.don.SendToNode(ctx, nodeAddr, req)
	if err != nil {
		return err
	}

	h.lggr.Debugw("sent response to node", "to", nodeAddr)
	return nil
}

// isCacheableStatusCode returns true if the HTTP status code indicates a cacheable response.
// This includes successful responses (2xx) and client errors (4xx)
func isCacheableStatusCode(statusCode int) bool {
	return (statusCode >= 200 && statusCode < 300) || (statusCode >= 400 && statusCode < 500)
}
