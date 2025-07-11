package v2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/common/aggregation"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
)

type AuthHandler struct {
	services.StateMachine
	lggr             logger.Logger
	authorizedKeysMu sync.RWMutex
	// authorizedKeys maps workflow ID to a set of authorized keys
	authorizedKeys map[string]aggregation.StringSet
	agg            *aggregation.AuthAggregator
	config         ServiceConfig
	don            handlers.DON
	donConfig      *config.DONConfig
	stopCh         services.StopChan
}

// NewHTTPTriggerAuthorizer creates a new HTTPTriggerAuthorizer.
func NewAuthHandler(lggr logger.Logger, cfg ServiceConfig, don handlers.DON, donConfig *config.DONConfig) *AuthHandler {
	// f+1 identical responses from workflow are needed for a authorization key to be registered
	threshold := donConfig.F + 1
	return &AuthHandler{
		lggr:           logger.Named(lggr, "HTTPTriggerAuthHandler"),
		authorizedKeys: make(map[string]aggregation.StringSet),
		agg:            aggregation.NewAuthAggregator(lggr, threshold, time.Duration(cfg.CleanUpPeriodMs)*time.Millisecond),
		don:            don,
		donConfig:      donConfig,
		stopCh:         make(services.StopChan),
	}
}

func (h *AuthHandler) Authorize(workflowID, payload, signature string) bool {
	// TODO: Implement the authorization logic.
	return true
}

func (h *AuthHandler) syncAuthorizedKeys() {
	authData := h.agg.Aggregate()
	authorizedKeys := make(map[string]aggregation.StringSet)
	for _, data := range authData {
		h.authorizedKeys[data.WorkflowID] = make(aggregation.StringSet)
		for _, key := range data.AuthorizedKeys {
			h.authorizedKeys[data.WorkflowID].Add(key.PublicKey)
		}
	}
	h.authorizedKeysMu.Lock()
	defer h.authorizedKeysMu.Unlock()
	h.authorizedKeys = authorizedKeys
}

func (h *AuthHandler) sendAuthPullRequest(ctx context.Context) error {
	req := &jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      gateway.GetRequestID(gateway.MethodWorkflowPullAuthMetadata),
		Method:  gateway.MethodWorkflowPullAuthMetadata,
	}
	var combinedErr error
	for _, member := range h.donConfig.Members {
		err := h.don.SendToNode(ctx, member.Address, req)
		if err != nil {
			combinedErr = errors.Join(combinedErr, fmt.Errorf("failed to send auth pull request to node %s: %w", member.Address, err))
		}
	}
	return combinedErr
}

func (h *AuthHandler) OnAuthMetadataPush(ctx context.Context, resp *jsonrpc.Response[json.RawMessage], nodeAddr string) error {
	var authData gateway.WorkflowAuthMetadata
	if err := json.Unmarshal(*resp.Result, &authData); err != nil {
		return fmt.Errorf("failed to unmarshal auth metadata: %w", err)
	}
	h.lggr.Debugw("Received auth metadata push", "workflowID", authData.WorkflowID, "nodeAddr", nodeAddr)
	for _, key := range authData.AuthorizedKeys {
		h.agg.Collect(aggregation.WorkflowAuthObservation{
			WorkflowID:    authData.WorkflowID,
			AuthorizedKey: key,
		}, nodeAddr)
	}
	return nil
}

func (h *AuthHandler) OnAuthMetadataPullResponse(ctx context.Context, resp *jsonrpc.Response[json.RawMessage], nodeAddr string) error {
	var authData []gateway.WorkflowAuthMetadata
	if err := json.Unmarshal(*resp.Result, &authData); err != nil {
		return fmt.Errorf("failed to unmarshal auth metadata pull response: %w", err)
	}
	h.lggr.Debugw("Received auth metadata pull response", "nodeAddr", nodeAddr)
	for _, data := range authData {
		for _, key := range data.AuthorizedKeys {
			h.agg.Collect(aggregation.WorkflowAuthObservation{
				WorkflowID:    data.WorkflowID,
				AuthorizedKey: key,
			}, nodeAddr)
		}
	}
	return nil
}

// Start begins the periodic pull loop.
func (h *AuthHandler) Start(ctx context.Context) error {
	return h.StartOnce("AuthHandler", func() error {
		h.lggr.Info("Starting HTTP Trigger Authorizer")
		err := h.agg.Start(ctx)
		if err != nil {
			return err
		}
		h.runTicker(ctx, time.Duration(h.config.AuthPullIntervalMs)*time.Millisecond, func() {
			err2 := h.sendAuthPullRequest(ctx)
			if err2 != nil {
				h.lggr.Errorw("Failed to send auth pull request", "error", err2)
			}
		})
		h.runTicker(ctx, time.Duration(h.config.AuthAggregationIntervalMs)*time.Millisecond, h.syncAuthorizedKeys)
		return nil
	})
}

func (h *AuthHandler) runTicker(ctx context.Context, periodMs time.Duration, fn func()) {
	go func() {
		ticker := time.NewTicker(periodMs)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				fn()
			case <-h.stopCh:
				return
			}
		}
	}()
}

func (h *AuthHandler) Close() error {
	return h.StopOnce("AuthHandler", func() error {
		h.lggr.Info("Stopping HTTP Trigger Authorizer")
		if err := h.agg.Close(); err != nil {
			h.lggr.Errorw("Failed to close AuthAggregator", "error", err)
		}
		close(h.stopCh)
		return nil
	})
}
