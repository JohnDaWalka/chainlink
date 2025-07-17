package ccipsolana

import (
	"bytes"
	"context"
	"fmt"

	bin "github.com/gagliardetto/binary"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

type usdcAttestationPayload struct {
	Message     []byte
	Attestation []byte
}

type SolanaTokenDataEncoder struct{}

func NewSolanaTokenDataEncoder() SolanaTokenDataEncoder {
	return SolanaTokenDataEncoder{}
}

func (e SolanaTokenDataEncoder) EncodeUSDC(_ context.Context, message cciptypes.Bytes, attestation cciptypes.Bytes) (cciptypes.Bytes, error) {
	messageAndAttestation := usdcAttestationPayload{
		Message:     message,
		Attestation: attestation,
	}
	buf := new(bytes.Buffer)
	err := bin.NewBorshEncoder(buf).Encode(messageAndAttestation)
	if err != nil {
		return nil, fmt.Errorf("failed to borsh encode USDC message and attestation: %w", err)
	}
	return buf.Bytes(), nil
}
