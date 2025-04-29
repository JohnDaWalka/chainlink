package utils

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink/deployment"
)

// ValidateContract validates a contract's existence and type
func ValidateContract(e deployment.Environment, chainSel uint64, contractAddr string, expectedType deployment.ContractType, expectedVersion semver.Version) error {
	records, err := e.DataStore.Addresses().Fetch()
	if err != nil {
		return fmt.Errorf("failed to fetch addresses from datastore: %w", err)
	}

	var tv *deployment.TypeAndVersion
	for _, record := range records {
		if record.Address == contractAddr {
			tv = &deployment.TypeAndVersion{
				Type:    deployment.ContractType(record.Type),
				Version: *record.Version,
			}
			break
		}
	}
	if tv == nil {
		return fmt.Errorf("unable to find contract %s in datastore", contractAddr)
	}

	if tv.Type != expectedType || tv.Version != expectedVersion {
		return fmt.Errorf(
			"unexpected contract type %s for %s on chain selector %d)",
			tv,
			expectedType,
			chainSel,
		)
	}

	return nil
}
