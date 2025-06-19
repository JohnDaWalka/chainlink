package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-common/pkg/gateway/jsonrpc"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
)

// DummyHandler forwards each request/response without doing any checks.
type dummyHandler struct {
	donConfig      *config.DONConfig
	don            DON
	savedCallbacks map[string]*savedCallback
	mu             sync.Mutex
	lggr           logger.Logger
	codec          api.JsonRPCCodec
}

type savedCallback struct {
	id         string
	callbackCh chan<- UserCallbackPayload
}

var _ Handler = (*dummyHandler)(nil)

func NewDummyHandler(donConfig *config.DONConfig, don DON, lggr logger.Logger) (Handler, error) {
	return &dummyHandler{
		donConfig:      donConfig,
		don:            don,
		savedCallbacks: make(map[string]*savedCallback),
		lggr:           logger.Named(lggr, "DummyHandler."+donConfig.DonId),
		codec:          api.JsonRPCCodec{},
	}, nil
}

func (d *dummyHandler) HandleUserMessage(ctx context.Context, req *jsonrpc.Request, callbackCh chan<- UserCallbackPayload) error {
	if req.Params == nil {
		return errors.New("missing params attribute")
	}
	var msg api.Message
	err := json.Unmarshal(req.Params, &msg)
	if err != nil {
		return err
	}
	d.mu.Lock()
	d.savedCallbacks[msg.Body.MessageId] = &savedCallback{msg.Body.MessageId, callbackCh}
	don := d.don
	d.mu.Unlock()

	// Send to all nodes.
	for _, member := range d.donConfig.Members {
		err = multierr.Combine(err, don.SendToNode(ctx, member.Address, req))
	}
	return err
}

func (d *dummyHandler) HandleNodeMessage(ctx context.Context, resp *jsonrpc.Response, nodeAddr string) error {
	msg, err := d.codec.DecodeResponse(resp.Result)
	if err != nil {
		return err
	}
	d.mu.Lock()
	savedCb, found := d.savedCallbacks[msg.Body.MessageId]
	delete(d.savedCallbacks, msg.Body.MessageId)
	d.mu.Unlock()

	if found {
		// Send first response from a node back to the user, ignore any other ones.
		savedCb.callbackCh <- UserCallbackPayload{Resp: resp, ErrCode: api.NoError, ErrMsg: ""}
		close(savedCb.callbackCh)
	}
	return nil
}

func (d *dummyHandler) Start(context.Context) error {
	return nil
}

func (d *dummyHandler) Close() error {
	return nil
}
