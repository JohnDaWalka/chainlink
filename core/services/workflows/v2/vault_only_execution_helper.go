package v2

import (
	"context"
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/secrets"
	sdkpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/workflowkey"
	"google.golang.org/protobuf/types/known/anypb"
)

type VaultOnlyExecutionHelper struct {
	*baseExecutionHelper
	WorkflowOwner string
	Key           workflowkey.Key
}

func NewVaultOnlyExecutionHelper(
	engine *Engine,
	workflowOwner string,
	key workflowkey.Key,
) *VaultOnlyExecutionHelper {
	return &VaultOnlyExecutionHelper{
		baseExecutionHelper: &baseExecutionHelper{
			Engine: engine,
			// No execution ID when calling the Vault during the registration step. How to handle this?
			// WorkflowExecutionID: WorkflowExecutionID,
		},
		WorkflowOwner: workflowOwner,
		Key:           key,
	}
}

var _ host.ExecutionHelper = &VaultOnlyExecutionHelper{}

func (v *VaultOnlyExecutionHelper) CallCapability(ctx context.Context, request *sdkpb.CapabilityRequest) (*sdkpb.CapabilityResponse, error) {
	// TODO: needs a semaphore to limit concurrent calls but only when called during the registration step.
	if request.Id != "vault@1.0.0" || request.Method != "GetSecrets" {
		return nil, errors.New("capability calls cannot be made during this execution")
	}

	p := &sdkpb.GetSecretRequest{}
	err := request.Payload.UnmarshalTo(p)
	if err != nil {
		return nil, fmt.Errorf("unexpected secrets payload for vault.GetSecrets: %w", err)
	}

	sreq := &secrets.GetSecretRequest{
		Id:             p.Id,
		Namespace:      p.Namespace,
		Owner:          v.WorkflowOwner,
		EncryptionKeys: []string{v.Key.PublicKeyString()},
	}

	request.Payload, err = anypb.New(sreq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal GetSecretRequest to request.Payload: %w", err)
	}

	r, err := v.baseExecutionHelper.CallCapability(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to execute vault.GetSecrets: %w", err)
	}

	if r.GetError() != "" {
		return r, nil
	}

	sr := &secrets.GetSecretResponse{}
	err = r.GetPayload().UnmarshalTo(sr)
	if err != nil {
		return nil, fmt.Errorf("unexpected secrets response for vault.GetSecrets: %w", err)
	}

	secret, err := decryptShares(v.Key, sr.EncryptedDecryptionKeyShares)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt secret shares: %w", err)
	}

	sdkr := &sdkpb.GetSecretResponse{
		Id:        sr.Id,
		Namespace: sr.Namespace,
		Owner:     sr.Owner,
		Value:     secret,
	}

	anysdkr, err := anypb.New(sdkr)
	if err != nil {
		return nil, fmt.Errorf("failed to wrap in any: %w", err)
	}

	r.Response = &sdkpb.CapabilityResponse_Payload{
		Payload: anysdkr,
	}

	return r, nil
}

func decryptShares(key workflowkey.Key, shares []*secrets.EncryptedDecryptionShare) (string, error) {
	// Implement the decryption logic here using the provided key and shares.
	// This is a placeholder for the actual decryption logic.
	return "decrypted_secret_value", nil
}

func (d VaultOnlyExecutionHelper) GetID() string {
	return ""
}
