package changeset

import (
	"errors"

	"github.com/Masterminds/semver/v3"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink/deployment"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	tonaddress "github.com/xssnick/tonutils-go/address"
)

const (
	TonCCIP     deployment.ContractType = "TonCCIP"
	TonReceiver deployment.ContractType = "TonReceiver"
)

// TonCCIPChainState holds a Go binding for all the currently deployed CCIP contracts
// on a chain. If a binding is nil, it means here is no such contract on the chain.
type or struct {
	LinkTokenAddress *tonaddress.Address
	CCIPAddress      *tonaddress.Address
	OffRamp          *tonaddress.Address
	Router           *tonaddress.Address

	// dummy receiver address
	ReceiverAddress *tonaddress.Address
}

func SaveOnchainStateTon(chainSelector uint64, tonState TonCCIPChainState, e deployment.Environment) error {
	ab := e.ExistingAddresses
	if !tonState.LinkTokenAddress.IsAddrNone() {
		ab.Save(chainSelector, tonState.LinkTokenAddress.String(), deployment.NewTypeAndVersion(commontypes.LinkToken, deployment.Version1_6_0))
	}
	if !tonState.CCIPAddress.IsAddrNone() {
		ab.Save(chainSelector, tonState.CCIPAddress.String(), deployment.NewTypeAndVersion(TonCCIP, deployment.Version1_6_0))
	}
	if !tonState.ReceiverAddress.IsAddrNone() {
		ab.Save(chainSelector, tonState.ReceiverAddress.String(), deployment.NewTypeAndVersion(TonReceiver, deployment.Version1_6_0))
	}
	return nil
}

func LoadOnchainStateTon(e deployment.Environment) (CCIPOnChainState, error) {
	state := CCIPOnChainState{
		TonChains: make(map[uint64]TonCCIPChainState),
	}
	for chainSelector, chain := range e.TonChains {
		addresses, err := e.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			// Chain not found in address book, initialize empty
			if !errors.Is(err, deployment.ErrChainNotFound) {
				return state, err
			}
			addresses = make(map[string]deployment.TypeAndVersion)
		}
		chainState, err := LoadChainStateTon(chain, addresses)
		if err != nil {
			return state, err
		}
		state.TonChains[chainSelector] = chainState
	}
	return state, nil
}

// LoadChainStateTon Loads all state for a TonChain into state
func LoadChainStateTon(chain deployment.TonChain, addresses map[string]deployment.TypeAndVersion) (TonCCIPChainState, error) {
	state := TonCCIPChainState{}

	// Most programs upgraded in place, but some are not so we always want to
	// load the latest version
	versions := make(map[deployment.ContractType]semver.Version)
	for addressStr, tvStr := range addresses {
		address, err := tonaddress.ParseAddr(addressStr)
		if err != nil {
			return state, err
		}

		switch tvStr.Type {
		case commontypes.LinkToken:
			state.LinkTokenAddress = address
		case TonCCIP:
			state.CCIPAddress = address
		case TonReceiver:
			state.ReceiverAddress = address
		case OffRamp:
			state.OffRamp = address
		case Router:
			state.Router = address
		default:
			log.Warn().Str("address", addressStr).Str("type", string(tvStr.Type)).Msg("Unknown Ton address type")
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
