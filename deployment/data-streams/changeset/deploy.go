package changeset

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/deployment"
	ds "github.com/smartcontractkit/chainlink/deployment/datastore"
)

type (
	// Contract covers contracts such as channel_config_store.ChannelConfigStore and fee_manager.FeeManager.
	Contract interface {
		// Caller:
		Owner(opts *bind.CallOpts) (common.Address, error)
		TypeAndVersion(opts *bind.CallOpts) (string, error)

		// Transactor:
		AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)
		TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)
	}

	ContractDeployFn[C Contract] func(chain deployment.Chain) *ContractDeployment[C]

	ContractDeployment[C Contract] struct {
		Address  common.Address
		Contract C
		Tx       *types.Transaction
		Tv       deployment.TypeAndVersion
		Err      error
	}
)

var _ deployment.ChangeSetV2[DeployChannelConfigStoreConfig] = DeployChannelConfigStore{}

// DeployContractV2 deploys a contract and saves the address to datastore.
func DeployContractV2[C Contract](
	e deployment.Environment,
	dataStore ds.MutableDataStore[SerializedContractMetadata, ds.DefaultMetadata],
	metadata SerializedContractMetadata,
	ab deployment.AddressBook,
	chain deployment.Chain,
	deployFn ContractDeployFn[C],
) (*ContractDeployment[C], error) {
	contractDeployment := deployFn(chain)
	if contractDeployment.Err != nil {
		e.Logger.Errorw("Failed to deploy contract", "err", contractDeployment.Err, "chain", chain.Selector)
		return nil, contractDeployment.Err
	}
	_, err := chain.Confirm(contractDeployment.Tx)
	if err != nil {
		e.Logger.Errorw("Failed to confirm deployment", "err", err)
		return nil, err
	}
	e.Logger.Infow("Deployed contract", "Contract", contractDeployment.Tv.String(), "addr", contractDeployment.Address.String(), "chain", chain.String())

	// Store Address
	if err = dataStore.Addresses().Add(
		ds.AddressRef{
			ChainSelector: chain.Selector,
			Address:       contractDeployment.Address.String(),
			Type:          ds.ContractType(contractDeployment.Tv.Type),
			Version:       &contractDeployment.Tv.Version,
		},
	); err != nil {
		e.Logger.Errorw("Failed to save contract address", "err", err)
		return nil, err
	}

	// Add a new ContractMetadata entry for the newly deployed contract
	if err = dataStore.ContractMetadata().Add(
		ds.ContractMetadata[SerializedContractMetadata]{
			ChainSelector: chain.Selector,
			Address:       contractDeployment.Address.String(),
			Metadata:      metadata,
		},
	); err != nil {
		return nil, fmt.Errorf("failed to save contract metadata: %w", err)
	}

	// Maintained for some existing backwards compatibility. Remove after fully migrated to datastore
	err = ab.Save(chain.Selector, contractDeployment.Address.String(), contractDeployment.Tv)
	if err != nil {
		e.Logger.Errorw("Failed to save contract address", "err", err)
		return nil, err
	}

	return contractDeployment, nil
}

// DeployContract deploys a contract and saves the address to the address book.
//
// Note that this function modifies the given address book variable, so it should be passed by reference.
func DeployContract[C Contract](
	e deployment.Environment,
	ab deployment.AddressBook,
	chain deployment.Chain,
	deployFn ContractDeployFn[C],
) (*ContractDeployment[C], error) {
	contractDeployment := deployFn(chain)
	if contractDeployment.Err != nil {
		e.Logger.Errorw("Failed to deploy contract", "err", contractDeployment.Err, "chain", chain.Selector)
		return nil, contractDeployment.Err
	}
	_, err := chain.Confirm(contractDeployment.Tx)
	if err != nil {
		e.Logger.Errorw("Failed to confirm deployment", "err", err)
		return nil, err
	}
	e.Logger.Infow("Deployed contract", "Contract", contractDeployment.Tv.String(), "addr", contractDeployment.Address.String(), "chain", chain.String())
	err = ab.Save(chain.Selector, contractDeployment.Address.String(), contractDeployment.Tv)
	if err != nil {
		e.Logger.Errorw("Failed to save contract address", "err", err)
		return nil, err
	}
	return contractDeployment, nil
}
