package changeset

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	commonTypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

func ValidateForwarderForChain(env cldf.Environment, chainSelector uint64, forwarderAddress common.Address) error {
	evmChains := env.BlockChains.EVMChains()
	contractSetsResp, err := GetContractSets(env.Logger, &GetContractSetsRequest{
		Chains:      evmChains,
		AddressBook: env.ExistingAddresses,
	})
	if err != nil {
		return fmt.Errorf("failed to get contract sets: %w", err)
	}

	_, ok := env.BlockChains.EVMChains()[chainSelector]
	if !ok {
		return errors.New("chain not found in environment")
	}
	chainSets, ok := contractSetsResp.ContractSets[chainSelector]
	if !ok {
		return errors.New("chain not found in on chain sets")
	}
	if chainSets.Forwarder == nil {
		return errors.New("forwarders not found in on chain sets")
	}
	_, ok = chainSets.Forwarder[forwarderAddress]
	if !ok {
		return errors.New("contract not found in on chain sets")
	}
	return nil
}

func ValidateMCMSAddresses(ab cldf.AddressBook, chainSelector uint64) error {
	if _, err := cldf.SearchAddressBook(ab, chainSelector, commonTypes.RBACTimelock); err != nil {
		return fmt.Errorf("timelock not present on the chain %w", err)
	}
	if _, err := cldf.SearchAddressBook(ab, chainSelector, commonTypes.ProposerManyChainMultisig); err != nil {
		return fmt.Errorf("mcms proposer not present on the chain %w", err)
	}
	return nil
}
