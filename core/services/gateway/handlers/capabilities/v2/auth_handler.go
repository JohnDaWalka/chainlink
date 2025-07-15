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

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(lggr logger.Logger, cfg ServiceConfig, don handlers.DON, donConfig *config.DONConfig) *AuthHandler {
	// f+1 identical responses from workflow are needed for a authorization key to be registered
	threshold := donConfig.F + 1
	return &AuthHandler{
		lggr:           logger.Named(lggr, "HTTPTriggerAuthHandler"),
		authorizedKeys: make(map[string]aggregation.StringSet),
		agg:            aggregation.NewAuthAggregator(lggr, threshold, time.Duration(cfg.CleanUpPeriodMs)*time.Millisecond),
		don:            don,
		donConfig:      donConfig,
		config:         cfg,
		stopCh:         make(services.StopChan),
	}
}

func (h *AuthHandler) Authorize(workflowID, payload, signature string) bool {
	// TODO: PRODCRE-305 Implement authorization logic
	return true
}

// syncAuthorizedKeys aggregates the authorized keys from the AuthAggregator and updates the local cache.
// Should be called periodically to keep the authorized keys up to date.
func (h *AuthHandler) syncAuthorizedKeys() {
	authData, err := h.agg.Aggregate()
	if err != nil {
		h.lggr.Errorw("Failed to aggregate auth data", "error", err)
		return
	}
	authorizedKeys := make(map[string]aggregation.StringSet)
	for _, data := range authData {
		authorizedKeys[data.WorkflowID] = make(aggregation.StringSet)
		for _, key := range data.AuthorizedKeys {
			authorizedKeys[data.WorkflowID].Add(key.PublicKey)
		}
	}
	h.authorizedKeysMu.Lock()
	defer h.authorizedKeysMu.Unlock()
	h.authorizedKeys = authorizedKeys
}

// sendAuthPullRequest sends a request to all nodes in the DON to pull the latest auth metadata.
// no retries are performed, as the caller is expected to poll periodically.
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

// OnAuthMetadataPush handles the push of auth metadata from a node when a new workflow is registered
func (h *AuthHandler) OnAuthMetadataPush(ctx context.Context, resp *jsonrpc.Response[json.RawMessage], nodeAddr string) error {
	var authData gateway.WorkflowAuthMetadata
	if err := json.Unmarshal(*resp.Result, &authData); err != nil {
		return fmt.Errorf("failed to unmarshal auth metadata: %w", err)
	}
	h.lggr.Debugw("Received auth metadata push", "workflowID", authData.WorkflowID, "nodeAddr", nodeAddr)
	var combinedErr error
	for _, key := range authData.AuthorizedKeys {
		err := h.agg.Collect(aggregation.WorkflowAuthObservation{
			WorkflowID:    authData.WorkflowID,
			AuthorizedKey: key,
		}, nodeAddr)
		if err != nil {
			combinedErr = errors.Join(combinedErr, fmt.Errorf("failed to collect auth observation: %w", err))
		}
	}
	return combinedErr
}

// OnAuthMetadataPullResponse handles the response to the auth metadata pull request.
func (h *AuthHandler) OnAuthMetadataPullResponse(ctx context.Context, resp *jsonrpc.Response[json.RawMessage], nodeAddr string) error {
	var authData []gateway.WorkflowAuthMetadata
	if err := json.Unmarshal(*resp.Result, &authData); err != nil {
		return fmt.Errorf("failed to unmarshal auth metadata pull response: %w", err)
	}
	h.lggr.Debugw("Received auth metadata pull response", "nodeAddr", nodeAddr)
	var combinedErr error
	for _, data := range authData {
		for _, key := range data.AuthorizedKeys {
			err := h.agg.Collect(aggregation.WorkflowAuthObservation{
				WorkflowID:    data.WorkflowID,
				AuthorizedKey: key,
			}, nodeAddr)
			combinedErr = errors.Join(combinedErr, err)
		}
	}
	return combinedErr
}

// Start begins the periodic pull loop.
func (h *AuthHandler) Start(ctx context.Context) error {
	return h.StartOnce("AuthHandler", func() error {
		h.lggr.Info("Starting HTTP Trigger Authorizer")
		err := h.agg.Start(ctx)
		if err != nil {
			return err
		}
		h.runTicker(time.Duration(h.config.AuthPullIntervalMs)*time.Millisecond, func() {
			err2 := h.sendAuthPullRequest(ctx)
			if err2 != nil {
				h.lggr.Errorw("Failed to send auth pull request", "error", err2)
			}
		})
		h.runTicker(time.Duration(h.config.AuthAggregationIntervalMs)*time.Millisecond, h.syncAuthorizedKeys)
		return nil
	})
}

func (h *AuthHandler) runTicker(period time.Duration, fn func()) {
	go func() {
		ticker := time.NewTicker(period)
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
