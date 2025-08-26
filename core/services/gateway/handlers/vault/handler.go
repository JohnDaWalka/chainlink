package vault

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/ratelimit"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	gw_handlers "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
)

const (
	HandlerType          = "vault"
	defaultCleanUpPeriod = 5 * time.Second
)

var (
	_                                 gw_handlers.Handler = (*handler)(nil)
	ErrInsufficientResponsesForQuorum                     = errors.New("insufficient valid responses to reach quorum")
)

type metrics struct {
	requestInternalError metric.Int64Counter
	requestUserError     metric.Int64Counter
	requestSuccess       metric.Int64Counter
}

func newMetrics() (*metrics, error) {
	requestInternalError, err := beholder.GetMeter().Int64Counter("gateway_vault_request_internal_error")
	if err != nil {
		return nil, fmt.Errorf("failed to register internal error counter: %w", err)
	}

	requestUserError, err := beholder.GetMeter().Int64Counter("gateway_vault_request_user_error")
	if err != nil {
		return nil, fmt.Errorf("failed to register user error counter: %w", err)
	}

	requestSuccess, err := beholder.GetMeter().Int64Counter("gateway_vault_request_success")
	if err != nil {
		return nil, fmt.Errorf("failed to register success counter: %w", err)
	}

	return &metrics{
		requestInternalError: requestInternalError,
		requestUserError:     requestUserError,
		requestSuccess:       requestSuccess,
	}, nil
}

type activeRequest struct {
	req       jsonrpc.Request[json.RawMessage]
	responses []*jsonrpc.Response[json.RawMessage]
	mu        sync.Mutex

	createdAt  time.Time
	callbackCh chan<- gw_handlers.UserCallbackPayload
}

type capabilitiesRegistry interface {
	DONForCapability(ctx context.Context, capabilityID string) (core.DON, []core.Node, error)
}

type handler struct {
	services.StateMachine
	methodConfig Config
	donConfig    *config.DONConfig
	don          gw_handlers.DON
	lggr         logger.Logger
	codec        api.JsonRPCCodec
	mu           sync.RWMutex
	stopCh       services.StopChan

	nodeRateLimiter *ratelimit.RateLimiter
	requestTimeout  time.Duration

	activeRequests map[string]*activeRequest
	metrics        *metrics

	capabilitiesRegistry capabilitiesRegistry
}

func (h *handler) HealthReport() map[string]error {
	return map[string]error{h.Name(): h.Healthy()}
}

func (h *handler) Name() string {
	return h.lggr.Name()
}

type SecretEntry struct {
	ID        string `json:"id"`
	Value     string `json:"value"`
	CreatedAt int64  `json:"created_at"`
}

type Config struct {
	NodeRateLimiter   ratelimit.RateLimiterConfig `json:"nodeRateLimiter"`
	RequestTimeoutSec int                         `json:"requestTimeoutSec"`
}

func NewHandler(methodConfig json.RawMessage, donConfig *config.DONConfig, don gw_handlers.DON, capabilitiesRegistry capabilitiesRegistry, lggr logger.Logger) (*handler, error) {
	var cfg Config
	if err := json.Unmarshal(methodConfig, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal method config: %w", err)
	}

	if cfg.RequestTimeoutSec == 0 {
		cfg.RequestTimeoutSec = 30
	}

	nodeRateLimiter, err := ratelimit.NewRateLimiter(cfg.NodeRateLimiter)
	if err != nil {
		return nil, fmt.Errorf("failed to create node rate limiter: %w", err)
	}

	metrics, err := newMetrics()
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics: %w", err)
	}

	return &handler{
		methodConfig:         cfg,
		donConfig:            donConfig,
		don:                  don,
		lggr:                 logger.Named(lggr, "VaultHandler:"+donConfig.DonId),
		requestTimeout:       time.Duration(cfg.RequestTimeoutSec) * time.Second,
		nodeRateLimiter:      nodeRateLimiter,
		activeRequests:       make(map[string]*activeRequest),
		mu:                   sync.RWMutex{},
		stopCh:               make(services.StopChan),
		metrics:              metrics,
		capabilitiesRegistry: capabilitiesRegistry,
	}, nil
}

func (h *handler) Start(ctx context.Context) error {
	return h.StartOnce("VaultHandler", func() error {
		h.lggr.Info("starting vault handler")
		go func() {
			ctx, cancel := h.stopCh.NewCtx()
			defer cancel()
			ticker := time.NewTicker(defaultCleanUpPeriod)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					h.removeExpiredRequests(ctx)
				case <-h.stopCh:
					return
				}
			}
		}()
		return nil
	})
}

func (h *handler) Close() error {
	return h.StopOnce("VaultHandler", func() error {
		h.lggr.Info("closing vault handler")
		close(h.stopCh)
		return nil
	})
}

// removeExpiredRequests removes expired requests from the pending requests map
func (h *handler) removeExpiredRequests(ctx context.Context) {
	h.mu.RLock()
	var expiredRequests []*activeRequest
	now := time.Now()
	for _, userRequest := range h.activeRequests {
		if now.Sub(userRequest.createdAt) > h.requestTimeout {
			expiredRequests = append(expiredRequests, userRequest)
		}
	}
	h.mu.RUnlock()

	for _, er := range expiredRequests {
		err := h.sendResponse(ctx, er, h.errorResponse(er.req, api.RequestTimeoutError, errors.New("request expired without getting any response")))
		if err != nil {
			h.lggr.Errorw("error sending response to user", "request_id", er.req.ID, "error", err)
		}
	}
}

func (h *handler) Methods() []string {
	return []string{
		MethodSecretsCreate,
		MethodSecretsGet,
		MethodSecretsUpdate,
		MethodSecretsDelete,
	}
}

func (h *handler) HandleLegacyUserMessage(ctx context.Context, msg *api.Message, callbackCh chan<- gw_handlers.UserCallbackPayload) error {
	return errors.New("vault handler does not support legacy messages")
}

func (h *handler) HandleJSONRPCUserMessage(ctx context.Context, req jsonrpc.Request[json.RawMessage], callbackCh chan<- gw_handlers.UserCallbackPayload) error {
	// Generate a unique ID for the request.
	// We do this ourselves to ensure the ID is unique and can't be tampered with by the user.
	if req.ID == "" {
		return errors.New("request ID cannot be empty")
	}

	h.lggr.Debugw("handling vault request", "method", req.Method, "requestID", req.ID)
	ar := &activeRequest{
		callbackCh: callbackCh,
		req:        req,
		createdAt:  time.Now(),
	}

	h.mu.Lock()
	h.activeRequests[req.ID] = ar
	h.mu.Unlock()
	switch req.Method {
	case MethodSecretsCreate:
		return h.handleSecretsCreate(ctx, ar)
	case MethodSecretsGet:
		return h.handleSecretsGet(ctx, ar)
	case MethodSecretsUpdate:
		return h.handleSecretsUpdate(ctx, ar)
	case MethodSecretsDelete:
		return h.handleSecretsDelete(ctx, ar)
	default:
		return h.sendResponse(ctx, ar, h.errorResponse(req, api.UnsupportedMethodError, errors.New("this method is unsupported: "+req.Method)))
	}
}

func (h *handler) HandleNodeMessage(ctx context.Context, resp *jsonrpc.Response[json.RawMessage], nodeAddr string) error {
	l := logger.With(h.lggr, "method", resp.Method, "requestID", resp.ID, "nodeAddr", nodeAddr)
	l.Debugw("handling node response")

	if !h.nodeRateLimiter.Allow(nodeAddr) {
		l.Debugw("node is rate limited", "nodeAddr", nodeAddr)
		return nil
	}

	h.mu.RLock()
	ar, ok := h.activeRequests[resp.ID]
	h.mu.RUnlock()
	if !ok {
		// Request is not found, so we don't need to send a response to the user
		// This might happen if the response is stale
		l.Errorw("no pending request found for ID")
		h.metrics.requestInternalError.Add(ctx, 1, metric.WithAttributes(
			attribute.String("don_id", h.donConfig.DonId),
			attribute.String("error", api.StaleNodeResponseError.String()),
		))
		return nil
	}

	err := h.validateUsingSignatures(ctx, *resp)
	if err == nil {
		l.Debugw("successfully validated signatures")
		h.sendSuccessResponse(ctx, l, ar, resp)
		return nil
	}

	l.Debugw("failed to validate signatures, falling back to quorum aggregation", "error", err)
	err = h.validateUsingQuorum(ctx, ar)
	switch {
	case errors.Is(err, ErrInsufficientResponsesForQuorum):
		ar.mu.Lock()
		l.Debugw("not enough valid responses to reach quorum", "error", err, "currentResponses", len(ar.responses))
		ar.responses = append(ar.responses, resp)
		ar.mu.Unlock()
	default:
		l.Errorw("error validating using quorum", "error", err)
		h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.HandlerError, fmt.Errorf("error validating using quorum: %w", err)))
	}

	l.Debugw("successfully reached quorum for response")
	h.sendSuccessResponse(ctx, l, ar, resp)
	return nil
}

func (h *handler) sendSuccessResponse(ctx context.Context, l logger.Logger, ar *activeRequest, resp *jsonrpc.Response[json.RawMessage]) error {
	rawResponse, err := jsonrpc.EncodeResponse(resp)
	if err != nil {
		l.Errorw("failed to encode response", "error", err)
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.NodeReponseEncodingError, fmt.Errorf("failed to marshal response: %w", err)))
	}

	var errorCode api.ErrorCode
	if resp.Error != nil {
		errorCode = api.FromJSONRPCErrorCode(resp.Error.Code)
	} else {
		errorCode = api.NoError
	}

	l.Debugw("issued user callback", "errorCode", errorCode)
	successResp := gw_handlers.UserCallbackPayload{
		RawResponse: rawResponse,
		ErrorCode:   errorCode,
	}
	return h.sendResponse(ctx, ar, successResp)
}

func (h *handler) validateUsingQuorum(ctx context.Context, ar *activeRequest) error {
	don, _, err := h.capabilitiesRegistry.DONForCapability(ctx, vault.CapabilityID)
	if err != nil {
		return err
	}

	requiredQuorum := int(2*don.F + 1)

	ar.mu.Lock()
	defer ar.mu.Unlock()
	if len(ar.responses) < requiredQuorum {
		return ErrInsufficientResponsesForQuorum
	}

	shaToCount := map[string]int{}
	for _, r := range ar.responses {
		sha, err := r.Digest()
		if err != nil {
			return fmt.Errorf("failed to validate using quorum: failed to compute digest: %w", err)
		}
		if shaToCount[sha] >= requiredQuorum {
			return nil
		}
	}

	return ErrInsufficientResponsesForQuorum
}

// TODO: use the right response type
type Response struct {
	ID         string
	Error      string
	Payload    []byte
	Format     string
	Context    []byte
	Signatures [][]byte
}

// TODO: use the right version
func ValidateSignatures(resp *Response, allowedSigners []common.Address, minRequired int) error {
	if len(resp.Context) <= 64 {
		return fmt.Errorf("context too short: expected min 64 bytes, got %d bytes", len(resp.Context))
	}

	if len(resp.Signatures) < minRequired {
		return fmt.Errorf("not enough signatures: expected min %d, got %d", minRequired, len(resp.Signatures))
	}

	// The context contains:
	// 0:32 -> config digest
	// 32:64 -> epoch + round, namely:
	//   - 0:27 -> zero padding
	//   - 27:31 -> sequence number (big endian uint32)
	//   - 31:32 -> zero round value
	cd, epochRound := resp.Context[:32], resp.Context[32:64]
	configDigest, err := ocr2types.BytesToConfigDigest(cd)
	if err != nil {
		return fmt.Errorf("invalid config digest in signature: %w", err)
	}

	epoch := binary.BigEndian.Uint32(epochRound[27:31])
	round := uint8(epochRound[31])

	fullHash := ocr2key.ReportToSigData(ocr2types.ReportContext{
		ReportTimestamp: ocr2types.ReportTimestamp{
			ConfigDigest: configDigest,
			Epoch:        epoch,
			Round:        round,
		},
	}, resp.Payload)

	validSigners := map[common.Address]bool{}
	for _, s := range resp.Signatures {
		signerPubkey, err := crypto.SigToPub(fullHash, s)
		if err != nil {
			return fmt.Errorf("invalid signature: %w", err)
		}
		signerAddr := crypto.PubkeyToAddress(*signerPubkey)

		for _, as := range allowedSigners {
			if as.Hex() == signerAddr.Hex() {
				validSigners[signerAddr] = true
				break
			}
		}

		if len(validSigners) >= minRequired {
			return nil
		}
	}

	return fmt.Errorf("only %d valid signatures, need at least %d", len(validSigners), minRequired)
}

func (h *handler) validateUsingSignatures(ctx context.Context, resp jsonrpc.Response[json.RawMessage]) error {
	if resp.Result == nil {
		return errors.New("response result is nil: cannot validate signatures")
	}

	r := &Response{}
	err := json.Unmarshal(*resp.Result, r)
	if err != nil {
		return err
	}

	don, nodes, err := h.capabilitiesRegistry.DONForCapability(ctx, vault.CapabilityID)
	if err != nil {
		return err
	}

	signers := []common.Address{}
	for _, n := range nodes {
		signers = append(signers, common.BytesToAddress(n.Signer[0:20]))
	}

	err = ValidateSignatures(r, signers, int(don.F+1))
	if err != nil {
		return fmt.Errorf("failed to validate signatures: %w", err)
	}

	return nil
}

func (h *handler) handleSecretsCreate(ctx context.Context, ar *activeRequest) error {
	l := logger.With(h.lggr, "method", ar.req.Method, "requestID", ar.req.ID)
	var secretsCreateRequest SecretsCreateRequest
	if err := json.Unmarshal(*ar.req.Params, &secretsCreateRequest); err != nil {
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.UserMessageParseError, err))
	}

	if secretsCreateRequest.ID == "" || secretsCreateRequest.Value == "" || secretsCreateRequest.Owner == "" {
		l.Debugw("invalid request parameters: secret id, owner and value cannot be empty")
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.InvalidParamsError, errors.New("secret id and value cannot be empty")))
	}

	// At this point, we know that the request is valid and we can send it to the nodes
	return h.fanOutToVaultNodes(ctx, l, ar)
}

func (h *handler) handleSecretsUpdate(ctx context.Context, ar *activeRequest) error {
	l := logger.With(h.lggr, "method", ar.req.Method, "requestID", ar.req.ID)

	req := &vault.UpdateSecretsRequest{}
	if err := json.Unmarshal(*ar.req.Params, req); err != nil {
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.UserMessageParseError, err))
	}

	req.RequestId = ar.req.ID

	reqb, err := json.Marshal(req)
	if err != nil {
		l.Errorw("failed to marshal request", "error", err)
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.NodeReponseEncodingError, fmt.Errorf("failed to marshal request: %w", err)))
	}

	ar.req.Params = (*json.RawMessage)(&reqb)
	return h.fanOutToVaultNodes(ctx, l, ar)
}

func (h *handler) handleSecretsDelete(ctx context.Context, ar *activeRequest) error {
	l := logger.With(h.lggr, "method", ar.req.Method, "requestId", ar.req.ID)

	req := &vault.DeleteSecretsRequest{}
	if err := json.Unmarshal(*ar.req.Params, req); err != nil {
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.UserMessageParseError, err))
	}

	req.RequestId = ar.req.ID

	reqb, err := json.Marshal(req)
	if err != nil {
		l.Errorw("failed to marshal request", "error", err)
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.NodeReponseEncodingError, fmt.Errorf("failed to marshal request: %w", err)))
	}

	ar.req.Params = (*json.RawMessage)(&reqb)
	return h.fanOutToVaultNodes(ctx, l, ar)
}

func (h *handler) fanOutToVaultNodes(ctx context.Context, l logger.Logger, ar *activeRequest) error {
	var nodeErrors []error
	for _, node := range h.donConfig.Members {
		err := h.don.SendToNode(ctx, node.Address, &ar.req)
		if err != nil {
			nodeErrors = append(nodeErrors, err)
			l.Errorw("error sending request to node", "node", node.Address, "error", err)
		}
	}

	if len(nodeErrors) == len(h.donConfig.Members) && len(nodeErrors) > 0 {
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.FatalError, errors.New("failed to forward user request to nodes")))
	}

	l.Debugw("successfully forwarded request to Vault nodes")
	return nil
}

func (h *handler) handleSecretsGet(ctx context.Context, ar *activeRequest) error {
	l := logger.With(h.lggr, "method", ar.req.Method, "requestID", ar.req.ID)
	var secretsGetRequest SecretsGetRequest
	if err := json.Unmarshal(*ar.req.Params, &secretsGetRequest); err != nil {
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.UserMessageParseError, err))
	}

	if secretsGetRequest.ID == "" || secretsGetRequest.Owner == "" {
		l.Debugw("invalid request parameters: secret id and owner cannot be empty")
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.InvalidParamsError, errors.New("secret id and value cannot be empty")))
	}

	// At this point, we know that the request is valid and we can send it to the nodes
	var nodeErrors []error
	for _, node := range h.donConfig.Members {
		err := h.don.SendToNode(ctx, node.Address, &ar.req)
		if err != nil {
			nodeErrors = append(nodeErrors, err)
			l.Errorw("error sending request to node", "node", node.Address, "error", err)
		}
	}

	if len(nodeErrors) == len(h.donConfig.Members) && len(nodeErrors) > 0 {
		return h.sendResponse(ctx, ar, h.errorResponse(ar.req, api.FatalError, errors.New("failed to forward user request to nodes")))
	}

	l.Debugw("successfully forwarded request to Vault nodes")
	return nil
}

func (h *handler) errorResponse(
	req jsonrpc.Request[json.RawMessage],
	errorCode api.ErrorCode,
	errs ...error,
) gw_handlers.UserCallbackPayload {
	err := errors.New("unknown error")
	if len(errs) > 0 && errs[0] != nil {
		err = errs[0]
	}

	switch errorCode {
	case api.FatalError:
	case api.NodeReponseEncodingError:
		h.lggr.Errorw(err.Error(), "request_id", req.ID)
		// Intentionally hide the error from the user
		err = errors.New(errorCode.String())
	case api.InvalidParamsError:
		h.lggr.Errorw("invalid params", "request_id", req.ID, "params", string(*req.Params))
		err = errors.New("invalid params error: " + err.Error())
	case api.UnsupportedMethodError:
		h.lggr.Errorw("unsupported method", "request_id", req.ID, "method", req.Method)
		err = errors.New("unsupported method: " + req.Method)
	case api.UserMessageParseError:
		h.lggr.Errorw("user message parse error", "request_id", req.ID, "error", err.Error())
		err = errors.New("user message parse error: " + err.Error())
	case api.NoError:
	case api.UnsupportedDONIdError:
	case api.HandlerError:
	case api.RequestTimeoutError:
	case api.StaleNodeResponseError:
		// Unused in this handler
	}

	return gw_handlers.UserCallbackPayload{
		RawResponse: h.codec.EncodeNewErrorResponse(
			req.ID,
			api.ToJSONRPCErrorCode(errorCode),
			err.Error(),
			nil,
		),
		ErrorCode: errorCode,
	}
}

func (h *handler) sendResponse(ctx context.Context, userRequest *activeRequest, resp gw_handlers.UserCallbackPayload) error {
	switch resp.ErrorCode {
	case api.StaleNodeResponseError:
	case api.FatalError:
	case api.NodeReponseEncodingError:
	case api.RequestTimeoutError:
	case api.HandlerError:
		h.metrics.requestInternalError.Add(ctx, 1, metric.WithAttributes(
			attribute.String("don_id", h.donConfig.DonId),
			attribute.String("error", resp.ErrorCode.String()),
		))
	case api.InvalidParamsError:
	case api.UnsupportedMethodError:
	case api.UserMessageParseError:
	case api.UnsupportedDONIdError:
		h.metrics.requestUserError.Add(ctx, 1, metric.WithAttributes(
			attribute.String("don_id", h.donConfig.DonId),
		))
	case api.NoError:
		h.metrics.requestSuccess.Add(ctx, 1, metric.WithAttributes(
			attribute.String("don_id", h.donConfig.DonId),
		))
	}

	select {
	case userRequest.callbackCh <- resp:
		h.lggr.Debugw("sent response", "request_id", userRequest.req.ID, "error_code", resp.ErrorCode, "raw_response", string(resp.RawResponse))
		h.mu.Lock()
		delete(h.activeRequests, userRequest.req.ID)
		h.mu.Unlock()
		return nil
	case <-ctx.Done():
		h.mu.Lock()
		delete(h.activeRequests, userRequest.req.ID)
		h.mu.Unlock()
		return ctx.Err()
	}
}
