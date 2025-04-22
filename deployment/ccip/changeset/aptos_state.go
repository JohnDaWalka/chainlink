package changeset

import (
	"errors"

	"github.com/Masterminds/semver/v3"
	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink/deployment"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

const (
	AptosCCIP     deployment.ContractType = "AptosCCIP"
	AptosReceiver deployment.ContractType = "AptosReceiver"
)

// AptosCCIPChainState holds a Go binding for all the currently deployed CCIP contracts
// on a chain. If a binding is nil, it means here is no such contract on the chain.
type AptosCCIPChainState struct {
	LinkTokenAddress aptos.AccountAddress
	CCIPAddress      aptos.AccountAddress

	// dummy receiver address
	ReceiverAddress aptos.AccountAddress
}

func SaveOnchainStateAptos(chainSelector uint64, aptosState AptosCCIPChainState, e deployment.Environment) error {
	ab := e.ExistingAddresses
	if aptosState.LinkTokenAddress != aptos.AccountZero {
		ab.Save(chainSelector, aptosState.LinkTokenAddress.String(), deployment.NewTypeAndVersion(commontypes.LinkToken, deployment.Version1_6_0))
	}
	if aptosState.CCIPAddress != aptos.AccountZero {
		ab.Save(chainSelector, aptosState.CCIPAddress.String(), deployment.NewTypeAndVersion(AptosCCIP, deployment.Version1_6_0))
	}
	if aptosState.ReceiverAddress != aptos.AccountZero {
		ab.Save(chainSelector, aptosState.ReceiverAddress.String(), deployment.NewTypeAndVersion(AptosReceiver, deployment.Version1_6_0))
	}
	return nil
}

func LoadOnchainStateAptos(e deployment.Environment) (CCIPOnChainState, error) {
	state := CCIPOnChainState{
		AptosChains: make(map[uint64]AptosCCIPChainState),
	}
	for chainSelector, chain := range e.AptosChains {
		addresses, err := e.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			// Chain not found in address book, initialize empty
			if !errors.Is(err, deployment.ErrChainNotFound) {
				return state, err
			}
			addresses = make(map[string]deployment.TypeAndVersion)
		}
		chainState, err := LoadChainStateAptos(chain, addresses)
		if err != nil {
			return state, err
		}
		state.AptosChains[chainSelector] = chainState
	}
	return state, nil
}

// LoadChainStateSolana Loads all state for a SolChain into state
func LoadChainStateAptos(chain deployment.AptosChain, addresses map[string]deployment.TypeAndVersion) (AptosCCIPChainState, error) {
	state := AptosCCIPChainState{}

	// Most programs upgraded in place, but some are not so we always want to
	// load the latest version
	versions := make(map[deployment.ContractType]semver.Version)
	for addressStr, tvStr := range addresses {
		var address aptos.AccountAddress
		if err := address.ParseStringRelaxed(addressStr); err != nil {
			return state, err
		}

		switch tvStr.Type {
		case commontypes.LinkToken:
			state.LinkTokenAddress = address
		case AptosCCIP:
			state.CCIPAddress = address
		case AptosReceiver:
			state.ReceiverAddress = address
		default:
			log.Warn().Str("address", addressStr).Str("type", string(tvStr.Type)).Msg("Unknown Aptos address type")
			continue
		}

		existingVersion, ok := versions[tvStr.Type]
		if ok {
			log.Warn().Str("existingVersion", existingVersion.String()).Str("type", string(tvStr.Type)).Msg("Duplicate address type found")
		}
		versions[tvStr.Type] = tvStr.Version
	}

	return state, nil
}
