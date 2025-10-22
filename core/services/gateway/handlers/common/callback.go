package common

import (
	"context"
	"errors"
	"sync/atomic"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
)

type Callback struct {
	ch chan handlers.UserCallbackPayload

	sent       atomic.Bool
	waitCalled atomic.Bool
	stopCh     services.StopChan // Signal when Wait() has finished
}

func (c *Callback) SendResponse(payload handlers.UserCallbackPayload) error {
	if !c.sent.CompareAndSwap(false, true) {
		return errors.New("response already sent: each callback can only be used once")
	}
	// Check if Wait() has already been called due to timeout, then return an error
	select {
	case <-c.stopCh:
		return errors.New("receiver is no longer waiting: Wait() has already finished")
	case c.ch <- payload:
		return nil
	}
}

func (c *Callback) Wait(ctx context.Context) (handlers.UserCallbackPayload, error) {
	if !c.waitCalled.CompareAndSwap(false, true) {
		return handlers.UserCallbackPayload{}, errors.New("Wait can only be called once per Callback instance")
	}
	defer close(c.stopCh)

	select {
	case <-ctx.Done():
		return handlers.UserCallbackPayload{}, ctx.Err()
	case r := <-c.ch:
		return r, nil
	}
}

func NewCallback() *Callback {
	ch := make(chan handlers.UserCallbackPayload, 1)
	stopCh := make(services.StopChan)
	return &Callback{ch: ch, stopCh: stopCh}
}
