package v2

import (
	"context"
	"errors"
	"sync"
	"time"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

type savedCallback struct {
	callbackCh chan<- *jsonrpc.Response
	createdAt  time.Time
}

type requestCallbacks struct {
	callbacks   map[string]savedCallback
	callbacksMu sync.Mutex
	lggr        logger.Logger
}

type RequestCallbacks interface {
	// SetCallback stores a callback channel for a given request ID.
	SetCallback(requestID string, callbackCh chan<- *jsonrpc.Response)
	// SendResponse sends a response to the callback channel for the given request ID.
	SendResponse(ctx context.Context, requestID string, resp *jsonrpc.Response) error
}

func newRequestCallbacks(lggr logger.Logger) *requestCallbacks {
	return &requestCallbacks{
		callbacks: make(map[string]savedCallback),
		lggr:      logger.Named(lggr, "RequestCallbacks"),
	}
}

func (rc *requestCallbacks) SetCallback(requestID string, callbackCh chan<- *jsonrpc.Response) {
	rc.callbacksMu.Lock()
	defer rc.callbacksMu.Unlock()

	rc.callbacks[requestID] = savedCallback{
		callbackCh: callbackCh,
		createdAt:  time.Now(),
	}
}

func (rc *requestCallbacks) SendResponse(ctx context.Context, requestID string, resp *jsonrpc.Response) error {
	rc.callbacksMu.Lock()
	defer rc.callbacksMu.Unlock()
	saved, exists := rc.callbacks[requestID]
	if !exists {
		return errors.New("callback not found for request ID: " + requestID)
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case saved.callbackCh <- resp:
		close(saved.callbackCh)
	}
	return nil
}
