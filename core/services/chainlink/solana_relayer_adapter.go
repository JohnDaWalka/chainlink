package chainlink

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana"
)

// SolanaRelayerAdapter wraps the Solana relayer to implement missing methods required by types.Relayer interface
type SolanaRelayerAdapter struct {
	*solana.Relayer
}

// NewSolanaRelayerAdapter creates a new adapter for the Solana relayer
func NewSolanaRelayerAdapter(relayer *solana.Relayer) *SolanaRelayerAdapter {
	return &SolanaRelayerAdapter{Relayer: relayer}
}

// NewCCIPProvider implements the missing method required by types.Relayer interface
func (s *SolanaRelayerAdapter) NewCCIPProvider(ctx context.Context, rargs types.RelayArgs) (types.CCIPProvider, error) {
	// For now, return an error as Solana CCIP provider is not implemented
	// This can be implemented when Solana CCIP support is added
	return nil, fmt.Errorf("NewCCIPProvider not implemented for Solana relayer")
}

// TON implements the missing method required by types.Relayer interface
func (s *SolanaRelayerAdapter) TON() (types.TONService, error) {
	// Solana relayer doesn't support TON
	return nil, fmt.Errorf("TON not supported by Solana relayer")
}
