package changeset

import (
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	"github.com/smartcontractkit/chainlink/deployment"
)

// SolChainState holds a Go binding for all the currently deployed CCIP programs
// on a chain. If a binding is nil, it means here is no such contract on the chain.

var (
	SolRouter deployment.ContractType = "SolRouter"
)

type SolCCIPChainState struct {
	SolRouter *ccip_router.Instruction
}
