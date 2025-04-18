package changeset

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
)

// contractConstructor is a function type that takes an address and a client,
// returning the contract instance and an error.
type contractConstructor[T any] func(address common.Address, client bind.ContractBackend) (*T, error)

// getContractsFromAddrBook retrieves a list of contract instances of a specified type from the address book.
// It uses the provided constructor to initialize matching contracts for the given chain.
func getContractsFromAddrBook[T any](
	addrBook cldf.AddressBook,
	chain deployment.Chain,
	desiredType cldf.ContractType,
	constructor contractConstructor[T],
) ([]*T, error) {
	chainAddresses, err := addrBook.AddressesForChain(chain.Selector)
	if err != nil {
		return nil, fmt.Errorf("failed to get addresses for chain %d: %w", chain.Selector, err)
	}

	var contracts []*T
	for addr, typeAndVersion := range chainAddresses {
		if typeAndVersion.Type == desiredType {
			address := common.HexToAddress(addr)
			contractInstance, err := constructor(address, chain.Client)
			if err != nil {
				return nil, fmt.Errorf("failed to construct %s at %s: %w", desiredType, addr, err)
			}
			contracts = append(contracts, contractInstance)
		}
	}

	if len(contracts) == 0 {
		return nil, fmt.Errorf("no %s found for chain %d", desiredType, chain.Selector)
	}

	return contracts, nil
}
