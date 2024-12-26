package changeset

import (
	"fmt"

	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_receiver"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/token_pool"
	"github.com/smartcontractkit/chainlink/deployment"
)

// SolChainState holds a Go binding for all the currently deployed CCIP programs
// on a chain. If a binding is nil, it means here is no such contract on the chain.

var (
	SolRouter    deployment.ContractType = "SolCcipRouter"
	SolReceiver  deployment.ContractType = "SolCcipReceiver"
	SolTokenPool deployment.ContractType = "SolTokenPool"
)

type SolCCIPRouter struct {
	SolanaChainSelector             uint64
	DefaultGasLimit                 ag_binary.Uint128
	DefaultAllowOutOfOrderExecution bool
	EnableExecutionAfter            int64
	// Accounts:
	Config                  ag_solanago.PublicKey
	State                   ag_solanago.PublicKey
	Authority               ag_solanago.PublicKey
	SystemProgram           ag_solanago.PublicKey
	Program                 ag_solanago.PublicKey
	ProgramData             ag_solanago.PublicKey
	ExternalExecutionConfig ag_solanago.PublicKey
	TokenPoolsSigner        ag_solanago.PublicKey
}

type SolCCIPChainState struct {
	CcipRouter   ag_solanago.PublicKey
	CcipReceiver ag_solanago.PublicKey
	TokenPool    ag_solanago.PublicKey
}

// TODO: Solana re-write
// we can add logic here but cleaner just to call LoadOnchainState_Sol for now ?
// the state will need to be defined separately for solana
// and LoadChainState() is completely different
func LoadOnchainStateSolana(e deployment.Environment) (CCIPOnChainState, error) {
	state := CCIPOnChainState{
		SolChains: make(map[uint64]SolCCIPChainState),
	}
	for chainSelector, chain := range e.SolChains {
		addresses, err := e.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			// Chain not found in address book, initialize empty
			if errors.Is(err, deployment.ErrChainNotFound) {
				addresses = make(map[string]deployment.TypeAndVersion)
			} else {
				return state, err
			}
		}
		chainState, err := LoadChainStateSolana(chain, addresses)
		if err != nil {
			return state, err
		}
		state.SolChains[chainSelector] = chainState

	}
	return state, nil
}

// LoadChainStateSolana Loads all state for a SolChain into state
func LoadChainStateSolana(chain deployment.SolChain, addresses map[string]deployment.TypeAndVersion) (SolCCIPChainState, error) {
	var state SolCCIPChainState
	for address, tvStr := range addresses {
		switch tvStr.String() {
		case deployment.NewTypeAndVersion(SolRouter, deployment.Version1_0_0).String():
			pub := ag_solanago.MustPublicKeyFromBase58(address)
			ccip_router.SetProgramID(pub)
			state.CcipRouter = pub
		case deployment.NewTypeAndVersion(SolReceiver, deployment.Version1_0_0).String():
			pub := ag_solanago.MustPublicKeyFromBase58(address)
			ccip_receiver.SetProgramID(pub)
			state.CcipReceiver = pub
		case deployment.NewTypeAndVersion(SolTokenPool, deployment.Version1_0_0).String():
			pub := ag_solanago.MustPublicKeyFromBase58(address)
			token_pool.SetProgramID(pub)
			state.TokenPool = pub
		default:
			return state, fmt.Errorf("unknown contract %s", tvStr)
		}
	}
	return state, nil
}
