package changeset

import (
	"errors"
	"fmt"

	ds "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/metadata"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/view"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/view/v0_5"
)

// SaveContractViews saves the contract views to the datastore.
var SaveContractViews = cldf.CreateChangeSet(saveViewsLogic, saveViewsPrecondition)

type SaveContractViewsConfig struct {
	Chains []uint64
}

func (cfg SaveContractViewsConfig) Validate() error {
	if len(cfg.Chains) == 0 {
		return errors.New("ConfigsByChain cannot be empty")
	}
	return nil
}

func saveViewsPrecondition(_ deployment.Environment, cc SaveContractViewsConfig) error {
	return cc.Validate()
}

func saveViewsLogic(e deployment.Environment, cfg SaveContractViewsConfig) (deployment.ChangesetOutput, error) {
	dataStore := ds.NewMemoryDataStore[metadata.SerializedContractMetadata, ds.DefaultMetadata]()
	state := DataStreamsOnChainState{Chains: make(map[uint64]DataStreamsChainState)}
	records, err := e.DataStore.Addresses().Fetch()
	if err != nil {
		return deployment.ChangesetOutput{}, errors.New("failed to fetch addresses")
	}

	addressesByChain := utils.AddressRefsToAddressByChain(records)

	envDatastore, err := ds.FromDefault[metadata.SerializedContractMetadata, ds.DefaultMetadata](e.DataStore)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to convert datastore: %w", err)
	}

	// Operation to Generate + Save Contract Views PER Chain.
	for _, chainSelector := range cfg.Chains {
		chainAddresses, ok := addressesByChain[chainSelector]
		if !ok {
			continue
		}
		chain := e.Chains[chainSelector]
		chainState, err := LoadChainState(e.Logger, chain, chainAddresses)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		state.Chains[chainSelector] = *chainState

		existingConfiguratorMetadata := make(map[view.Address]*metadata.GenericContractMetadata[v0_5.ConfiguratorView])
		configuratorContexts := make(map[view.Address]*ConfiguratorContext)
		for _, c := range chainState.Configurators {
			cm, err := envDatastore.ContractMetadata().Get(
				ds.NewContractMetadataKey(testutil.TestChain.Selector, c.Address().String()),
			)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to get contract metadata: %w", err)
			}
			contractMetadata, err := metadata.DeserializeMetadata[v0_5.ConfiguratorView](cm.Metadata)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to convert contract metadata: %w", err)
			}
			configuratorContexts[c.Address().String()] = &ConfiguratorContext{
				FromBlock: contractMetadata.Metadata.DeployBlock,
			}
			existingConfiguratorMetadata[c.Address().String()] = contractMetadata
		}

		configuratorViews, err := chainState.GenerateConfiguratorViews(e.GetContext(), configuratorContexts)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate configurator view: %w", err)
		}

		for address, contractView := range configuratorViews {
			existingMetadata := existingConfiguratorMetadata[address]
			contractMetadata := metadata.GenericContractMetadata[v0_5.ConfiguratorView]{
				Metadata: existingMetadata.Metadata,
				View:     contractView,
			}

			serialized, err := metadata.NewSerializedContractMetadata(contractMetadata)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to serialize contract metadata: %w", err)
			}

			if err = dataStore.ContractMetadata().Upsert(
				ds.ContractMetadata[metadata.SerializedContractMetadata]{
					ChainSelector: chain.Selector,
					Address:       address,
					Metadata:      *serialized,
				},
			); err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to upsert contract metadata: %w", err)
			}
		}
	}

	defaultDs, err := ds.ToDefault(dataStore.Seal())
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to convert data store to default format: %w", err)
	}

	return deployment.ChangesetOutput{DataStore: defaultDs}, nil

}
