package v2

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"go.uber.org/multierr"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/ratelimit"
	"github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	gateway_common "github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/webapi/webapicap"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

var _ jsonrpc.Service = (*service)(nil)
var _ NodeMessageHandler = (*service)(nil)

// NodeMessageHandler implements service-specific logic for managing messages from nodes.
type NodeMessageHandler interface {
	job.ServiceCtx

	// Handlers should not make any assumptions about goroutines calling HandleNodeMessage.
	// should be non-blocking
	HandleNodeMessage(ctx context.Context, req *jsonrpc.Response, nodeAddr string) error
}

type service struct {
	config          ServiceConfig
	don             handlers.DON
	donConfig       *config.DONConfig
	savedCallbacks  map[string]*savedCallback
	mu              sync.Mutex
	lggr            logger.Logger
	httpClient      network.HTTPClient
	nodeRateLimiter *ratelimit.RateLimiter
	wg              sync.WaitGroup
}

type ServiceConfig struct {
	NodeRateLimiter ratelimit.RateLimiterConfig `json:"nodeRateLimiter"`
}

type savedCallback struct {
	id         string
	callbackCh chan<- handlers.UserCallbackPayload
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
	return &service{
		config:          cfg,
		don:             don,
		donConfig:       donConfig,
		lggr:            logger.Named(lggr, "HTTPService"+donConfig.DonId),
		httpClient:      httpClient,
		nodeRateLimiter: nodeRateLimiter,
		wg:              sync.WaitGroup{},
		savedCallbacks:  make(map[string]*savedCallback),
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
				ExecutionError: true,
				ErrorMessage:   err.Error(),
			}
		} else {
			outboundResp = gateway.OutboundHTTPResponse{
				ExecutionError: false,
				StatusCode:     resp.StatusCode,
				Headers:        resp.Headers,
				Body:           resp.Body,
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

// TODO: can this return error?
func (h *service) HandleNodeMessage(ctx context.Context, resp *jsonrpc.Response, nodeAddr string) error {
	// TODO: what if error is received? log and do nothing?
	requestID := resp.ID
	if requestID == "" {
		return fmt.Errorf("received response with empty request ID from node %s", nodeAddr)
	}
	var outboundReq gateway.OutboundHTTPRequest
	err := json.Unmarshal(resp.Result, &outboundReq)
	if err != nil {
		return fmt.Errorf("failed to unmarshal HTTP request from node %s: %w", nodeAddr, err)
	}
	return h.handleOutgoingRequest(ctx, requestID, outboundReq, nodeAddr)
}

func (h *service) handleWebAPITriggerMessage(ctx context.Context, msg *api.Message, nodeAddr string) error {
	h.mu.Lock()
	savedCb, found := h.savedCallbacks[msg.Body.MessageId]
	delete(h.savedCallbacks, msg.Body.MessageId)
	h.mu.Unlock()

	if found {
		// Send first response from a node back to the user, ignore any other ones.
		// TODO: in practice, we should wait for at least 2F+1 nodes to respond and then return an aggregated response
		// back to the user.
		savedCb.callbackCh <- handlers.UserCallbackPayload{Msg: msg, ErrCode: api.NoError, ErrMsg: ""}
		close(savedCb.callbackCh)
	}
	return nil
}

func (h *service) HandleUserMessage(ctx context.Context, msg *api.Message, callbackCh chan<- handlers.UserCallbackPayload) error {
	h.mu.Lock()
	h.savedCallbacks[msg.Body.MessageId] = &savedCallback{msg.Body.MessageId, callbackCh}
	don := h.don
	h.mu.Unlock()
	body := msg.Body
	var payload webapicap.TriggerRequestPayload
	err := json.Unmarshal(body.Payload, &payload)
	if err != nil {
		h.lggr.Errorw("error decoding payload", "err", err)
		callbackCh <- handlers.UserCallbackPayload{Msg: msg, ErrCode: api.UserMessageParseError, ErrMsg: "error decoding payload " + err.Error()}
		close(callbackCh)
		return nil
	}

	if payload.Timestamp == 0 {
		h.lggr.Errorw("error decoding payload")
		callbackCh <- handlers.UserCallbackPayload{Msg: msg, ErrCode: api.UserMessageParseError, ErrMsg: "error decoding payload"}
		close(callbackCh)
		return nil
	}

	if uint(time.Now().Unix())-h.config.MaxAllowedMessageAgeSec > uint(payload.Timestamp) {
		callbackCh <- handlers.UserCallbackPayload{Msg: msg, ErrCode: api.HandlerError, ErrMsg: "stale message"}
		close(callbackCh)
		return nil
	}
	// TODO: apply allowlist and rate-limiting here
	if msg.Body.Method != MethodWebAPITrigger {
		h.lggr.Errorw("unsupported method", "method", body.Method)
		callbackCh <- handlers.UserCallbackPayload{Msg: msg, ErrCode: api.HandlerError, ErrMsg: "invalid method " + msg.Body.Method}
		close(callbackCh)
		return nil
	}
	req, err := common.ValidatedRequestFromMessage(msg)
	if err != nil {
		h.lggr.Errorw("error transforming message to request")
		callbackCh <- handlers.UserCallbackPayload{Msg: msg, ErrCode: api.UserMessageParseError, ErrMsg: "error transforming message to request"}
		close(callbackCh)
		return nil
	}
	// Send original request to all nodes
	for _, member := range h.donConfig.Members {
		err = multierr.Combine(err, don.SendToNode(ctx, member.Address, req))
	}
	return err
}

func (h *service) Start(context.Context) error {
	return nil
}

func (h *service) Close() error {
	h.wg.Wait()
	return nil
}
