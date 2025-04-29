package utils

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment"
)

func MaybeFindEthAddress(ab deployment.AddressBook, chain uint64, typ deployment.ContractType) (common.Address, error) {
	addressHex, err := deployment.SearchAddressBook(ab, chain, typ)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to find contract %s address: %w", typ, err)
	}
	address := common.HexToAddress(addressHex)
	return address, nil
}

func EnvironmentAddresses(e deployment.Environment) (addresses map[string]deployment.TypeAndVersion, err error) {
	addresses = make(map[string]deployment.TypeAndVersion)
	records, err := e.DataStore.Addresses().Fetch()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch addresses from datastore: %w", err)
	}
	for _, record := range records {
		addresses[record.Address] = deployment.TypeAndVersion{
			Type:    deployment.ContractType(record.Type),
			Version: *record.Version,
		}
	}
	return addresses, nil
}
