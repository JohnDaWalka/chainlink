package sui

import (
	"errors"
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
)

type CCIPChainState struct {
	CCIPAddress      aptos.AccountAddress
	LinkTokenAddress aptos.AccountAddress
}

// LoadOnchainStatesui loads chain state for sui chains from env
func LoadOnchainStatesui(env cldf.Environment) (map[uint64]CCIPChainState, error) {
	rawChains, err := env.BlockChains.SuiChains()
	if err != nil {
		return nil, fmt.Errorf("failed to get SuiChains: %w", err)
	}

	suiChains := make(map[uint64]CCIPChainState)

	for chainSelector := range rawChains {
		addresses, err := env.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			// Chain not found in address book, initialize empty state
			if !errors.Is(err, cldf.ErrChainNotFound) {
				return nil, fmt.Errorf("failed to get addresses for chain %d: %w", chainSelector, err)
			}
			addresses = make(map[string]cldf.TypeAndVersion)
		}

		chainState, err := loadsuiChainStateFromAddresses(addresses)
		if err != nil {
			return nil, fmt.Errorf("failed to load chain state for chain %d: %w", chainSelector, err)
		}

		suiChains[chainSelector] = chainState
	}

	return suiChains, nil
}

func loadsuiChainStateFromAddresses(addresses map[string]cldf.TypeAndVersion) (CCIPChainState, error) {
	chainState := CCIPChainState{}
	for addrStr, typeAndVersion := range addresses {
		// Parse address
		address := &aptos.AccountAddress{}
		err := address.ParseStringRelaxed(addrStr)
		if err != nil {
			return chainState, fmt.Errorf("failed to parse address %s for %s: %w", addrStr, typeAndVersion.Type, err)
		}
		switch typeAndVersion.Type {
		case shared.AptosCCIPType:
			chainState.CCIPAddress = *address
		}
		// Set address based on type

	}
	return chainState, nil
}
