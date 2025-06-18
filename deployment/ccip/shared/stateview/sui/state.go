package sui

import (
	"errors"
	"fmt"

	"github.com/pattonkan/sui-go/sui"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
)

type CCIPChainState struct {
	CCIPAddress             sui.Address
	CCIPObjectRef           sui.Address
	FeeQuoterCapId          sui.Address
	OnRampAddress           sui.Address
	OnRampStateObjectId     sui.Address
	OffRampAddress          sui.Address
	OffRampOwnerCapId       sui.Address
	OffRampStateObjectId    sui.Address
	LinkTokenAddress        sui.Address
	LinkTokenCoinMetadataId sui.Address
	LinkTokenTreasuryCapId  sui.Address
}

// LoadOnchainStatesui loads chain state for sui chains from env
func LoadOnchainStatesui(env cldf.Environment) (map[uint64]CCIPChainState, error) {
	rawChains := env.BlockChains.SuiChains()
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
	for addr, typeAndVersion := range addresses {
		// Parse address
		suiAddr := sui.MustAddressFromHex(addr)
		switch typeAndVersion.Type {

		case shared.SuiCCIPType:
			chainState.CCIPAddress = *suiAddr

		case shared.SuiCCIPObjectRefType:
			chainState.CCIPObjectRef = *suiAddr

		case shared.SuiFeeQuoterCapType:
			chainState.FeeQuoterCapId = *suiAddr

		case shared.SuiOnRampType:
			chainState.OnRampAddress = *suiAddr

		case shared.SuiOnRampStateObjectIdType:
			chainState.OnRampStateObjectId = *suiAddr

		case shared.SuiOffRampType:
			chainState.OffRampAddress = *suiAddr

		case shared.SuiOffRampStateObjectIdType:
			chainState.OffRampStateObjectId = *suiAddr

		case shared.SuiOffRampOwnerCapObjectIdType:
			chainState.OffRampOwnerCapId = *suiAddr

		case shared.SuiLinkTokenType:
			chainState.LinkTokenAddress = *suiAddr

		case shared.SuiLinkTokenObjectMetadataId:
			chainState.LinkTokenCoinMetadataId = *suiAddr

		case shared.SuiLinkTokenTreasuryCapId:
			chainState.LinkTokenTreasuryCapId = *suiAddr
		}
		// Set address based on type

	}
	return chainState, nil
}
