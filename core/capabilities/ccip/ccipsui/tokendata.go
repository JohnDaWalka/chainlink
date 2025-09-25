package ccipsui

import (
	"context"

	"github.com/aptos-labs/aptos-go-sdk/bcs"
	ccipocr3common "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

type suiUSDCAttestationPayload struct {
	Message     []byte
	Attestation []byte
}

type SuiTokenDataEncoder struct{}

func NewSuiTokenDataEncoder() SuiTokenDataEncoder {
	return SuiTokenDataEncoder{}
}

func (e SuiTokenDataEncoder) EncodeUSDC(_ context.Context, message ccipocr3common.Bytes, attestation ccipocr3common.Bytes) (ccipocr3common.Bytes, error) {
	// Create the payload structure with message and attestation
	payload := suiUSDCAttestationPayload{
		Message:     message,
		Attestation: attestation,
	}

	// Create a BCS serializer
	s := &bcs.Serializer{}

	// Serialize the message bytes
	s.WriteBytes(payload.Message)

	// Serialize the attestation bytes
	s.WriteBytes(payload.Attestation)

	// Check for serialization errors
	if s.Error() != nil {
		return nil, s.Error()
	}

	// Return the BCS encoded bytes
	return s.ToBytes(), nil
}
