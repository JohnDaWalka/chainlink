package vault

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jonboulle/clockwork"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/requests"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var _ capabilities.ExecutableCapability = (*Capability)(nil)

type Capability struct {
	lggr         logger.Logger
	clock        clockwork.Clock
	expiresAfter time.Duration
	handler      *requests.Handler[*Request, *Response]
}

func (s *Capability) Start(ctx context.Context) error {
	return s.handler.Start(ctx)
}

func (s *Capability) Close() error {
	return s.handler.Close()
}

func (s *Capability) Info(ctx context.Context) (capabilities.CapabilityInfo, error) {
	return capabilities.NewCapabilityInfo(vaultcommon.CapabilityID, capabilities.CapabilityTypeAction, "Vault Capability")
}

func (s *Capability) RegisterToWorkflow(ctx context.Context, request capabilities.RegisterToWorkflowRequest) error {
	// Left unimplemented as this method will never be called
	// for this capability
	return nil
}

func (s *Capability) UnregisterFromWorkflow(ctx context.Context, request capabilities.UnregisterFromWorkflowRequest) error {
	// Left unimplemented as this method will never be called
	// for this capability
	return nil
}

func (s *Capability) Execute(ctx context.Context, request capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
	if request.Payload == nil {
		return capabilities.CapabilityResponse{}, errors.New("capability does not support v1 requests")
	}

	if request.Method != MethodSecretsGet {
		return capabilities.CapabilityResponse{}, errors.New("unsupported method: can only call GetSecrets via capability interface")
	}

	r := &vaultcommon.GetSecretsRequest{}
	err := request.Payload.UnmarshalTo(r)
	if err != nil {
		return capabilities.CapabilityResponse{}, fmt.Errorf("could not unmarshal payload to GetSecretsRequest: %w", err)
	}

	// Validate the request: we only check that the request contains at least one secret request.
	// All other validations are done in the plugin and subject to consensus.
	if len(r.Requests) == 0 {
		return capabilities.CapabilityResponse{}, errors.New("no secret request specified in request")
	}

	// We need to generate sufficiently unique IDs accounting for two cases:
	// 1. called during the subscription phase, in which case the executionID will be blank
	// 2. called during execution, in which case it'll be present.
	// The reference ID is unique per phase, so we need to differentiate when generating
	// an ID.
	md := request.Metadata
	phaseOrExecution := md.WorkflowExecutionID
	if phaseOrExecution == "" {
		phaseOrExecution = "subscription"
	}
	id := fmt.Sprintf("%s::%s::%s", md.WorkflowID, phaseOrExecution, md.ReferenceID)

	resp, err := s.handleRequest(ctx, id, r)
	if err != nil {
		return capabilities.CapabilityResponse{}, err
	}

	// Note: we can drop the signatures from the response above here
	// since only a valid report will be successfully decryptable by the workflow DON.
	resppb := &vaultcommon.GetSecretsResponse{}
	err = proto.Unmarshal(resp.Payload, resppb)
	if err != nil {
		return capabilities.CapabilityResponse{}, fmt.Errorf("could not unmarshal response to GetSecretsResponse: %w", err)
	}

	anyproto, err := anypb.New(resppb)
	if err != nil {
		return capabilities.CapabilityResponse{}, fmt.Errorf("could not marshal response to anypb: %w", err)
	}

	return capabilities.CapabilityResponse{
		Payload: anyproto,
	}, nil
}

func (s *Capability) handleRequest(ctx context.Context, requestID string, request proto.Message) (*Response, error) {
	respCh := make(chan *Response, 1)
	s.handler.SendRequest(ctx, &Request{
		Payload:      request,
		ResponseChan: respCh,

		ExpiryTimeVal: s.clock.Now().Add(s.expiresAfter),
		IDVal:         requestID,
	})
	s.lggr.Debugw("sent request to OCR handler", "requestID", requestID)
	select {
	case <-ctx.Done():
		s.lggr.Debugw("request timed out", "requestID", requestID, "error", ctx.Err())
		return nil, ctx.Err()
	case resp := <-respCh:
		s.lggr.Debugw("received response for request", "requestID", requestID, "error", resp.Error)
		if resp.Error != "" {
			return nil, fmt.Errorf("error processing request %s: %w", requestID, errors.New(resp.Error))
		}

		return resp, nil
	}
}

func ValidateCreateSecretsRequest(request *vaultcommon.CreateSecretsRequest) error {
	if request.RequestId == "" {
		return errors.New("request ID must not be empty")
	}

	if len(request.EncryptedSecrets) == 0 {
		return errors.New("must have at least one request")
	}

	if len(request.EncryptedSecrets) >= MaxBatchSize {
		return fmt.Errorf("request batch size exceeds maximum of %d", MaxBatchSize)
	}

	uniqueIDs := map[string]bool{}
	for _, req := range request.EncryptedSecrets {
		if req.Id == nil {
			return errors.New("secret ID must not be nil")
		}

		if req.Id.Key == "" || req.Id.Owner == "" {
			return fmt.Errorf("secret ID must have both key and owner set: %v", req.Id)
		}

		if req.EncryptedValue == "" {
			return fmt.Errorf("encrypted value must be set for secret ID: %v", req.Id)
		}

		_, ok := uniqueIDs[KeyFor(req.Id)]
		if ok {
			return fmt.Errorf("duplicate secret ID found: %v", req.Id)
		}

		uniqueIDs[KeyFor(req.Id)] = true
	}

	return nil
}

func (s *Capability) CreateSecrets(ctx context.Context, request *vaultcommon.CreateSecretsRequest) (*Response, error) {
	s.lggr.Debugw("executing vault capability call", "method", "CreateSecrets", "requestID", request.RequestId)
	// TODO: validate that secrets are encrypted with the correct key
	err := ValidateCreateSecretsRequest(request)
	if err != nil {
		return nil, fmt.Errorf("failed to validate create secrets request: %w", err)
	}

	return s.handleRequest(ctx, request.RequestId, request)
}

func ValidateUpdateSecretsRequest(request *vaultcommon.UpdateSecretsRequest) error {
	if request.RequestId == "" {
		return errors.New("request ID must not be empty")
	}

	if len(request.EncryptedSecrets) == 0 {
		return errors.New("must have at least one request")
	}

	if len(request.EncryptedSecrets) >= MaxBatchSize {
		return fmt.Errorf("request batch size exceeds maximum of %d", MaxBatchSize)
	}

	uniqueIDs := map[string]bool{}
	for _, req := range request.EncryptedSecrets {
		if req.Id == nil {
			return errors.New("secret ID must not be nil")
		}

		if req.Id.Key == "" || req.Id.Owner == "" {
			return fmt.Errorf("secret ID must have both key and owner set: %v", req.Id)
		}

		if req.EncryptedValue == "" {
			return fmt.Errorf("encrypted value must be set for secret ID: %v", req.Id)
		}

		_, ok := uniqueIDs[KeyFor(req.Id)]
		if ok {
			return fmt.Errorf("duplicate secret ID found: %v", req.Id)
		}

		uniqueIDs[KeyFor(req.Id)] = true
	}

	return nil
}

func (s *Capability) UpdateSecrets(ctx context.Context, request *vaultcommon.UpdateSecretsRequest) (*Response, error) {
	s.lggr.Debugw("executing vault capability call", "method", "UpdateSecrets", "requestID", request.RequestId)
	// TODO: validate that secrets are encrypted with the correct key
	err := ValidateUpdateSecretsRequest(request)
	if err != nil {
		return nil, fmt.Errorf("failed to validate update secrets request: %w", err)
	}

	return s.handleRequest(ctx, request.RequestId, request)
}

func ValidateDeleteSecretsRequest(request *vaultcommon.DeleteSecretsRequest) error {
	if request.RequestId == "" {
		return errors.New("request ID must not be empty")
	}

	if len(request.Ids) == 0 {
		return errors.New("must have at least one request")
	}

	if len(request.Ids) >= MaxBatchSize {
		return fmt.Errorf("request batch size exceeds maximum of %d", MaxBatchSize)
	}

	uniqueIDs := map[string]bool{}
	for _, id := range request.Ids {
		if id.Key == "" || id.Owner == "" {
			return fmt.Errorf("secret ID must have both key and owner set: %v", id)
		}

		_, ok := uniqueIDs[KeyFor(id)]
		if ok {
			return fmt.Errorf("duplicate secret ID found: %v", id)
		}

		uniqueIDs[KeyFor(id)] = true
	}

	return nil
}

func (s *Capability) DeleteSecrets(ctx context.Context, request *vaultcommon.DeleteSecretsRequest) (*Response, error) {
	s.lggr.Debugw("executing vault capability call", "method", "DeleteSecrets", "requestID", request.RequestId)
	err := ValidateDeleteSecretsRequest(request)
	if err != nil {
		return nil, fmt.Errorf("failed to validate delete secrets request: %w", err)
	}

	return s.handleRequest(ctx, request.RequestId, request)
}

func ValidateListSecretIdentifiersRequest(request *vaultcommon.ListSecretIdentifiersRequest) error {
	if request.RequestId == "" {
		return errors.New("request ID must not be empty")
	}

	if request.Owner == "" {
		return errors.New("owner must not be empty")
	}

	return nil
}

func (s *Capability) ListSecretIdentifiers(ctx context.Context, request *vaultcommon.ListSecretIdentifiersRequest) (*Response, error) {
	s.lggr.Debugw("executing vault capability call", "method", "ListSecretIdentifiers", "requestID", request.RequestId)
	err := ValidateListSecretIdentifiersRequest(request)
	if err != nil {
		return nil, fmt.Errorf("failed to validate list secret identifiers request: %w", err)
	}

	return s.handleRequest(ctx, request.RequestId, request)
}

func ValidateGetSecretsRequest(request *vaultcommon.GetSecretsRequest) error {
	if len(request.Requests) == 0 {
		return errors.New("no GetSecret request specified in request")
	}
	if len(request.Requests) >= MaxBatchSize {
		return fmt.Errorf("request batch size exceeds maximum of %d", MaxBatchSize)
	}

	for _, req := range request.Requests {
		id := req.Id
		if req == nil {
			return errors.New("secret ID must not be nil")
		}
		if id.Key == "" || id.Owner == "" {
			return fmt.Errorf("secret ID must have both key and owner set: %v", id)
		}
	}

	return nil
}

func (s *Capability) GetSecrets(ctx context.Context, requestID string, request *vaultcommon.GetSecretsRequest) (*Response, error) {
	s.lggr.Debugw("executing vault capability call", "method", "GetSecrets")
	err := ValidateGetSecretsRequest(request)
	if err != nil {
		return nil, fmt.Errorf("failed to validate get secrets request: %w", err)
	}

	return s.handleRequest(ctx, requestID, request)
}

func NewCapability(
	lggr logger.Logger,
	clock clockwork.Clock,
	expiresAfter time.Duration,
	handler *requests.Handler[*Request, *Response],
) *Capability {
	return &Capability{
		lggr:         lggr.Named("VaultCapability"),
		clock:        clock,
		expiresAfter: expiresAfter,
		handler:      handler,
	}
}
