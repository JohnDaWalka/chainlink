package v2

import (
	"encoding/json"
	"net/rpc"
	"sync"
)

type HTTPTriggerAuthorizer struct {
	mu        sync.RWMutex
	workflows map[string]*authorizedSenders
}

type authorizedSenders struct {
	keys map[string]struct{}
}

// NewHTTPTriggerAuthorizer creates a new HTTPTriggerAuthorizer.
func NewHTTPTriggerAuthorizer() *HTTPTriggerAuthorizer {
	return &HTTPTriggerAuthorizer{
		workflows: make(map[string]*authorizedSenders),
	}
}

// TODO: revisit function signature
func (a *HTTPTriggerAuthorizer) AddWorkflow(id string, keys []string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	authKeys := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		authKeys[k] = struct{}{}
	}
	a.workflows[id] = &authorizedSenders{keys: authKeys}
}

// RemoveWorkflow removes a workflow by ID.
func (a *HTTPTriggerAuthorizer) RemoveWorkflow(id string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.workflows, id)
}

func (a *HTTPTriggerAuthorizer) Authorize(workflowID, payload, signature string) bool {
	// TODO: Implement the authorization logic.
}

// HTTPTriggerAuthSyncer syncs workflow authorization metadata.
type HTTPTriggerAuthSyncer struct {
	authorizer *HTTPTriggerAuthorizer
}

// NewHTTPTriggerAuthSyncer creates a new HTTPTriggerAuthSyncer.
func NewHTTPTriggerAuthSyncer(authorizer *HTTPTriggerAuthorizer) *HTTPTriggerAuthSyncer {
	return &HTTPTriggerAuthSyncer{authorizer: authorizer}
}

// OnPushAuth handles pushed workflow auth metadata.
func (s *HTTPTriggerAuthSyncer) OnPushAuth(resp *rpc.Response) error {
	return s.processAuthResponse(resp)
}

// OnPullAuthResponse handles pulled workflow auth metadata.
func (s *HTTPTriggerAuthSyncer) OnPullAuthResponse(resp *rpc.Response) error {
	return s.processAuthResponse(resp)
}

// processAuthResponse processes the JSON-RPC response and updates the authorizer.
func (s *HTTPTriggerAuthSyncer) processAuthResponse(resp *rpc.Response) error {
	var meta map[string][]string
	if err := json.Unmarshal(resp.Result, &meta); err != nil {
		return err
	}
	for workflowID, keys := range meta {
		s.authorizer.AddWorkflow(workflowID, keys)
	}
	return nil
}

import (
	"context"
	"time"
)

// GatewayConnection defines the interface for making pull requests.
type GatewayConnection interface {
	PullAuthMetadata(ctx context.Context) (*rpc.Response, error)
}

// HTTPTriggerAuthorizerRunner periodically pulls auth metadata.
type HTTPTriggerAuthorizerRunner struct {
	authorizer *HTTPTriggerAuthorizer
	syncer     *HTTPTriggerAuthSyncer
	conn       GatewayConnection
	interval   time.Duration
	stopCh     chan struct{}
}

// NewHTTPTriggerAuthorizerRunner creates a new runner.
func NewHTTPTriggerAuthorizerRunner(authorizer *HTTPTriggerAuthorizer, syncer *HTTPTriggerAuthSyncer, conn GatewayConnection, interval time.Duration) *HTTPTriggerAuthorizerRunner {
	return &HTTPTriggerAuthorizerRunner{
		authorizer: authorizer,
		syncer:     syncer,
		conn:       conn,
		interval:   interval,
		stopCh:     make(chan struct{}),
	}
}

// Start begins the periodic pull loop.
func (r *HTTPTriggerAuthorizerRunner) Start() {
	go func() {
		ticker := time.NewTicker(r.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				ctx, cancel := context.WithTimeout(context.Background(), r.interval)
				resp, err := r.conn.PullAuthMetadata(ctx)
				cancel()
				if err == nil && resp != nil {
					_ = r.syncer.OnPullAuthResponse(resp)
				}
			case <-r.stopCh:
				return
			}
		}
	}()
}

// Stop halts the periodic pull loop.
func (r *HTTPTriggerAuthorizerRunner) Stop() {
	close(r.stopCh)
}