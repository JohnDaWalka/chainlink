package capabilities

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"

	"go.uber.org/multierr"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink-common/pkg/values/pb"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/webapi/webapicap"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
)

const (
	MethodComputeAction               = "compute_action"
	MethodWebAPITarget                = "web_api_target"
	MethodWebAPITrigger               = "web_api_trigger"
	MethodWebAPITriggerUpdateMetadata = "web_api_trigger_update_metadata"
)

type TriggersConfig struct {
	lastUpdatedAt time.Time

	// map of trigger ID to trigger config
	triggersConfig map[string]webapicap.TriggerConfig
}

type AllNodesTriggersConfig struct {
	// map of node address to node's triggerConfig
	triggersConfigMap map[string]TriggersConfig
}

type handler struct {
	capabilities.Validator[webapicap.TriggerConfig, struct{}, capabilities.TriggerResponse]
	config          HandlerConfig
	don             handlers.DON
	donConfig       *config.DONConfig
	savedCallbacks  map[string]*savedCallback
	mu              sync.Mutex
	lggr            logger.Logger
	httpClient      network.HTTPClient
	nodeRateLimiter *common.RateLimiter
	wg              sync.WaitGroup
	// each gateway node has a map of trigger IDs to trigger configs
	triggersConfig AllNodesTriggersConfig
	// One of the nodes' trigger configs that is part of consensus
	consensusConfig TriggersConfig
}

type HandlerConfig struct {
	// Is this part of standard capability so doesn't need to be defined here?
	F                       uint                     `json:"F"`
	NodeRateLimiter         common.RateLimiterConfig `json:"nodeRateLimiter"`
	MaxAllowedMessageAgeSec uint                     `json:"maxAllowedMessageAgeSec"`
}

type savedCallback struct {
	id         string
	callbackCh chan<- handlers.UserCallbackPayload
}

var _ handlers.Handler = (*handler)(nil)

func NewHandler(handlerConfig json.RawMessage, donConfig *config.DONConfig, don handlers.DON, httpClient network.HTTPClient, lggr logger.Logger) (*handler, error) {
	var cfg HandlerConfig
	err := json.Unmarshal(handlerConfig, &cfg)
	if err != nil {
		lggr.Errorf("error unmarshalling config: %s, err: %s", string(handlerConfig), err.Error())
		return nil, err
	}
	nodeRateLimiter, err := common.NewRateLimiter(cfg.NodeRateLimiter)
	if err != nil {
		return nil, err
	}

	return &handler{
		Validator:       capabilities.NewValidator[webapicap.TriggerConfig, struct{}, capabilities.TriggerResponse](capabilities.ValidatorArgs{}),
		config:          cfg,
		consensusConfig: TriggersConfig{triggersConfig: make(map[string]webapicap.TriggerConfig)},
		don:             don,
		donConfig:       donConfig,
		lggr:            lggr.Named("WebAPIHandler." + donConfig.DonId),
		httpClient:      httpClient,
		nodeRateLimiter: nodeRateLimiter,
		wg:              sync.WaitGroup{},
		savedCallbacks:  make(map[string]*savedCallback),
		triggersConfig: AllNodesTriggersConfig{
			triggersConfigMap: make(map[string]TriggersConfig),
		},
	}, nil
}

// sendHTTPMessageToClient is an outgoing message from the gateway to external endpoints
// returns message to be sent back to the capability node
func (h *handler) sendHTTPMessageToClient(ctx context.Context, req network.HTTPRequest, msg *api.Message) (*api.Message, error) {
	var payload Response
	resp, err := h.httpClient.Send(ctx, req)
	if err != nil {
		return nil, err
	}
	payload = Response{
		ExecutionError: false,
		StatusCode:     resp.StatusCode,
		Headers:        resp.Headers,
		Body:           resp.Body,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &api.Message{
		Body: api.MessageBody{
			MessageId: msg.Body.MessageId,
			Method:    msg.Body.Method,
			DonId:     msg.Body.DonId,
			Payload:   payloadBytes,
		},
	}, nil
}

func (h *handler) handleWebAPITriggerMessage(ctx context.Context, msg *api.Message, nodeAddr string) error {
	h.mu.Lock()
	savedCb, found := h.savedCallbacks[msg.Body.MessageId]
	delete(h.savedCallbacks, msg.Body.MessageId)
	h.mu.Unlock()

	if found {
		// Send first response from a node back to the user, ignore any other ones.
		// TODO: in practice, we should wait for at least 2F+1 nodes to respond and then return an aggregated response
		// back to the user.  https://smartcontract-it.atlassian.net/browse/CAPPL-237
		savedCb.callbackCh <- handlers.UserCallbackPayload{Msg: msg, ErrCode: api.NoError, ErrMsg: ""}
		close(savedCb.callbackCh)
	}
	return nil
}

func (h *handler) handleWebAPIOutgoingMessage(ctx context.Context, msg *api.Message, nodeAddr string) error {
	h.lggr.Debugw("handling webAPI outgoing message", "messageId", msg.Body.MessageId, "nodeAddr", nodeAddr)
	if !h.nodeRateLimiter.Allow(nodeAddr) {
		return fmt.Errorf("rate limit exceeded for node %s", nodeAddr)
	}
	var payload Request
	err := json.Unmarshal(msg.Body.Payload, &payload)
	if err != nil {
		return err
	}

	timeout := time.Duration(payload.TimeoutMs) * time.Millisecond
	req := network.HTTPRequest{
		Method:  payload.Method,
		URL:     payload.URL,
		Headers: payload.Headers,
		Body:    payload.Body,
		Timeout: timeout,
	}

	// send response to node async
	h.wg.Add(1)
	go func() {
		defer h.wg.Done()
		// not cancelled when parent is cancelled to ensure the goroutine can finish
		newCtx := context.WithoutCancel(ctx)
		newCtx, cancel := context.WithTimeout(newCtx, timeout)
		defer cancel()
		l := h.lggr.With("url", payload.URL, "messageId", msg.Body.MessageId, "method", payload.Method)
		respMsg, err := h.sendHTTPMessageToClient(newCtx, req, msg)
		if err != nil {
			l.Errorw("error while sending HTTP request to external endpoint", "err", err)
			payload := Response{
				ExecutionError: true,
				ErrorMessage:   err.Error(),
			}
			payloadBytes, err2 := json.Marshal(payload)
			if err2 != nil {
				// should not happen
				l.Errorw("error while marshalling payload", "err", err2)
				return
			}
			respMsg = &api.Message{
				Body: api.MessageBody{
					MessageId: msg.Body.MessageId,
					Method:    msg.Body.Method,
					DonId:     msg.Body.DonId,
					Payload:   payloadBytes,
				},
			}
		}
		// this signature is not verified by the node because
		// WS connection between gateway and node are already verified
		respMsg.Signature = msg.Signature

		err = h.don.SendToNode(newCtx, nodeAddr, respMsg)
		if err != nil {
			l.Errorw("failed to send to node", "err", err, "to", nodeAddr)
			return
		}
		l.Debugw("sent response to node", "to", nodeAddr)
	}()
	return nil
}

//	body := api.MessageBody{
//		MessageId: types.RandomID().String(),
//		DonId:     h.connector.DonID(),
//		Method:    webapicapabilities.MethodWebAPITriggerUpdateMetadata,
//		Receiver:  gatewayIDStr,
//		Payload:   payloadJSON,
//	}
func (h *handler) handleWebAPITriggerUpdateMetadata(ctx context.Context, msg *api.Message, nodeAddr string) error {
	body := msg.Body
	h.lggr.Debugw("handleWebAPITriggerUpdateMetadata", "body", body, "payload", string(body.Payload))

	var payload map[string]string
	err := json.Unmarshal(body.Payload, &payload)
	donID := body.DonId
	if err != nil {
		h.lggr.Errorw("error decoding payload", "err", err, "payload", string(body.Payload))
		return err
	}
	triggersConfig := make(map[string]webapicap.TriggerConfig)
	for triggerID, configPbString := range payload {
		pbBytes, err := base64.StdEncoding.DecodeString(configPbString)
		if err != nil {
			h.lggr.Errorw("error decoding pb bytes", "err", err)
			return err
		}
		pb := &pb.Map{}
		err = proto.Unmarshal(pbBytes, pb)
		if err != nil {
			h.lggr.Errorw("error unmarshalling", "err", err)
			return err
		}
		vmap, err := values.FromMapValueProto(pb)
		if err != nil {
			h.lggr.Errorw("error FromMapValueProto", "err", err)
			return err
		}
		if vmap == nil {
			h.lggr.Errorw("error FromMapValueProto nil vmap")
			return err
		}
		reqConfig, err := h.ValidateConfig(vmap)
		if err != nil {
			h.lggr.Errorw("error validating config", "err", err)
			return err
		}
		triggersConfig[triggerID] = *reqConfig
	}
	h.triggersConfig.triggersConfigMap[donID] = TriggersConfig{lastUpdatedAt: time.Now(), triggersConfig: triggersConfig}
	return nil
}

type triggerStruct struct {
	config webapicap.TriggerConfig
	count  uint
}

func (h *handler) updateTriggerConsensus() {
	// calculate the mode of the trigger configs to determine consensus

	// 1. Make list of node configs that haven't expired
	triggers := h.triggersConfig.triggersConfigMap

	// 2. Make list of hashes of each node's config.
	// Update, make a list of hashes of each node's triggerId and config.
	// 3. Count the number of times each hash appears.

	// triggerId => triggerConfig hash => triggerStruct
	triggersHashmap := make(map[string]map[string]triggerStruct)

	// We need to know which triggerConfig has the most votes for each triggerId.
	// We need to know the count of each triggerConfig for each triggerId.
	// Need a structure that is for each triggerId, there is a set of triggerConfigs, their hashes and the count of each hash.
	// convert map to list of triggerConfigs to triggerConfig

	// Performance is a concern as this is O(nodes) * O(triggers)

	for _, triggerConfig := range triggers {
		for triggerID, triggerConfig := range triggerConfig.triggersConfig {
			_, isTriggerInHashMap := triggersHashmap[triggerID]
			if !isTriggerInHashMap {
				triggersHashmap[triggerID] = make(map[string]triggerStruct)
			}
			s := fmt.Sprintf("%v", triggerConfig)
			h := sha256.New()
			h.Write([]byte(s))
			hash := hex.EncodeToString(h.Sum(nil))
			t, isHashInHashMap := triggersHashmap[triggerID][hash]
			if !isHashInHashMap {
				triggersHashmap[triggerID][hash] = triggerStruct{config: triggerConfig, count: 1}
			} else {
				t.count++
				triggersHashmap[triggerID][hash] = t
			}
		}
	}

	type kv struct {
		Hash  string
		Value uint
	}
	// 4. Find the hash that appears the most.
	// https://stackoverflow.com/questions/18695346/how-can-i-sort-a-mapstringint-by-its-values

	for triggerID, triggersHashmapByHash := range triggersHashmap {
		var ss []kv
		for triggerConfigHash, v := range triggersHashmapByHash {
			ss = append(ss, kv{triggerConfigHash, v.count})
		}

		sort.Slice(ss, func(i, j int) bool {
			return ss[i].Value > ss[j].Value
		})
		// 5. Check if the hash that appears the most appears more than F times.
		if ss[0].Value < h.config.F+1 {
			// 6. If it doesn't, do nothing
			return
		}
		// 7. If it does, update the consensus config to that config.
		// TODO: how do stale triggerIds get cleaned up? https://smartcontract-it.atlassian.net/browse/CAPPL-236
		h.consensusConfig.triggersConfig[triggerID] = triggersHashmapByHash[ss[0].Hash].config
		h.consensusConfig.lastUpdatedAt = time.Now()
	}
}
func (h *handler) HandleNodeMessage(ctx context.Context, msg *api.Message, nodeAddr string) error {
	switch msg.Body.Method {
	case MethodWebAPITrigger:
		return h.handleWebAPITriggerMessage(ctx, msg, nodeAddr)
	case MethodWebAPITarget, MethodComputeAction:
		return h.handleWebAPIOutgoingMessage(ctx, msg, nodeAddr)
	case MethodWebAPITriggerUpdateMetadata:
		return h.handleWebAPITriggerUpdateMetadata(ctx, msg, nodeAddr)
	default:
		return fmt.Errorf("unsupported method: %s", msg.Body.Method)
	}
}

func (h *handler) Start(context.Context) error {
	return nil
}

func (h *handler) Close() error {
	h.wg.Wait()
	return nil
}

func (h *handler) HandleUserMessage(ctx context.Context, msg *api.Message, callbackCh chan<- handlers.UserCallbackPayload) error {
	h.mu.Lock()
	h.savedCallbacks[msg.Body.MessageId] = &savedCallback{msg.Body.MessageId, callbackCh}
	don := h.don
	h.mu.Unlock()
	body := msg.Body
	var payload webapicap.TriggerRequestPayload
	err := json.Unmarshal(body.Payload, &payload)
	if err != nil {
		h.lggr.Errorw("error decoding payload", "err", err)
		callbackCh <- handlers.UserCallbackPayload{Msg: msg, ErrCode: api.UserMessageParseError, ErrMsg: fmt.Sprintf("error decoding payload %s", err.Error())}
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
	// TODO: apply allowlist and rate-limiting here  https://smartcontract-it.atlassian.net/browse/CAPPL-238
	if msg.Body.Method != MethodWebAPITrigger {
		h.lggr.Errorw("unsupported method", "method", body.Method)
		callbackCh <- handlers.UserCallbackPayload{Msg: msg, ErrCode: api.HandlerError, ErrMsg: fmt.Sprintf("invalid method %s", msg.Body.Method)}
		close(callbackCh)
		return nil
	}

	// Send to all nodes.
	for _, member := range h.donConfig.Members {
		err = multierr.Combine(err, don.SendToNode(ctx, member.Address, msg))
	}
	return err
}
