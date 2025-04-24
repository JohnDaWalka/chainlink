package ton

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	tonaddress "github.com/xssnick/tonutils-go/address"
)

const (
	TonMCMSType     deployment.ContractType = "TonManyChainMultisig"
	TonCCIPType     deployment.ContractType = "TonCCIP"
	TonReceiverType deployment.ContractType = "TonReceiver"
)

type TonCCIPChainState struct {
	MCMSAddress      tonaddress.Address
	CCIPAddress      tonaddress.Address
	LinkTokenAddress tonaddress.Address

	// Test contracts
	TestRouterAddress tonaddress.Address
	ReceiverAddress   tonaddress.Address
}

// LoadOnchainStateTon loads chain state for Ton chains from env
func LoadOnchainStateTon(env deployment.Environment) (map[uint64]TonCCIPChainState, error) {
	tonChains := make(map[uint64]TonCCIPChainState)
	for chainSelector := range env.TonChains {
		addresses, err := env.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			// Chain not found in address book, initialize empty
			if !errors.Is(err, deployment.ErrChainNotFound) {
				return tonChains, err
			}
			addresses = make(map[string]deployment.TypeAndVersion)
		}
		chainState, err := loadTonChainStateFromAddresses(addresses)
		if err != nil {
			return tonChains, err
		}
		tonChains[chainSelector] = chainState
	}
	return tonChains, nil
}

func loadTonChainStateFromAddresses(addresses map[string]deployment.TypeAndVersion) (TonCCIPChainState, error) {
	chainState := TonCCIPChainState{}
	for addrStr, typeAndVersion := range addresses {
		// Parse address
		address, err := tonaddress.ParseAddr(addrStr)
		if err != nil {
			return chainState, fmt.Errorf("failed to parse address %s for %s: %w", addrStr, typeAndVersion.Type, err)
		}
		// Set address based on type
		switch typeAndVersion.Type {
		case TonMCMSType:
			chainState.MCMSAddress = *address
		case TonCCIPType:
			chainState.CCIPAddress = *address
		case commontypes.LinkToken:
			chainState.LinkTokenAddress = *address
		case TonReceiverType:
			chainState.ReceiverAddress = *address
		}
	}
	return chainState, nil
}
