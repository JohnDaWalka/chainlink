package channel_config_store

import (
	"errors"
	"fmt"

	ds "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/metadata"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/mcmsutil"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/channel_config_store"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

// DeployChannelConfigStoreChangeset deploys ChannelConfigStore to the chains specified in the config.
var DeployChannelConfigStoreChangeset = deployment.CreateChangeSet(deployChannelConfigStoreLogic, deployChannelConfigStorePrecondition)

func (cc DeployChannelConfigStoreConfig) GetOwnershipConfig() types.OwnershipSettings {
	return cc.Ownership
}

type DeployChannelConfigStoreConfig struct {
	// ChainsToDeploy is a list of chain selectors to deploy the contract to.
	ChainsToDeploy []uint64
	Ownership      types.OwnershipSettings
}

func (cc DeployChannelConfigStoreConfig) Validate() error {
	if len(cc.ChainsToDeploy) == 0 {
		return errors.New("ChainsToDeploy is empty")
	}
	for _, chain := range cc.ChainsToDeploy {
		if err := deployment.IsValidChainSelector(chain); err != nil {
			return fmt.Errorf("invalid chain selector: %d - %w", chain, err)
		}
	}
	return nil
}

func deployChannelConfigStoreLogic(e deployment.Environment, cc DeployChannelConfigStoreConfig) (deployment.ChangesetOutput, error) {
	dataStore := ds.NewMemoryDataStore[metadata.SerializedContractMetadata, ds.DefaultMetadata]()
	err := performDeployment(e, dataStore, cc)
	if err != nil {
		e.Logger.Errorw("Failed to deploy ChannelConfigStore", "err", err)
		return deployment.ChangesetOutput{}, deployment.MaybeDataErr(err)
	}

	records, err := dataStore.Addresses().Fetch()
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to fetch addresses: %w", err)
	}
	proposals, err := mcmsutil.GetTransferOwnershipProposals(e, cc, records)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to transfer ownership to MCMS: %w", err)
	}

	sealedDs, err := ds.ToDefault(dataStore.Seal())
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to convert data store to default format: %w", err)
	}
	return deployment.ChangesetOutput{
		DataStore:             sealedDs,
		MCMSTimelockProposals: proposals,
	}, nil
}

func deployChannelConfigStorePrecondition(_ deployment.Environment, cc DeployChannelConfigStoreConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid DeployChannelConfigStoreConfig: %w", err)
	}

	return nil
}

func performDeployment(e deployment.Environment, dataStore ds.MutableDataStore[metadata.SerializedContractMetadata, ds.DefaultMetadata], cc DeployChannelConfigStoreConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid DeployChannelConfigStoreConfig: %w", err)
	}
	for _, chainSel := range cc.ChainsToDeploy {
		chain, ok := e.Chains[chainSel]
		if !ok {
			return fmt.Errorf("Chain not found for chain selector %d", chainSel)
		}
		res, err := changeset.DeployContract[*channel_config_store.ChannelConfigStore](e, dataStore, chain, channelConfigStoreDeployFn(), nil)
		if err != nil {
			return err
		}
		serialized, err := metadata.NewChannelConfigStoreMetadata(metadata.ChannelConfigStoreMetadata{Block: res.Block})
		if err != nil {
			return fmt.Errorf("failed to serialize verifier metadata: %w", err)
		}
		// Store ContractMetadata entry for the newly deployed contract
		if err = dataStore.ContractMetadata().Upsert(
			ds.ContractMetadata[metadata.SerializedContractMetadata]{
				ChainSelector: chain.Selector,
				Address:       res.Address.String(),
				Metadata:      *serialized,
			},
		); err != nil {
			return fmt.Errorf("failed to upser contract metadata: %w", err)
		}
	}

	return nil
}

// channelConfigStoreDeployFn returns a function that deploys a ChannelConfigStore contract.
func channelConfigStoreDeployFn() changeset.ContractDeployFn[*channel_config_store.ChannelConfigStore] {
	return func(chain deployment.Chain) *changeset.ContractDeployment[*channel_config_store.ChannelConfigStore] {
		ccsAddr, ccsTx, ccs, err := channel_config_store.DeployChannelConfigStore(
			chain.DeployerKey,
			chain.Client,
		)
		if err != nil {
			return &changeset.ContractDeployment[*channel_config_store.ChannelConfigStore]{
				Err: err,
			}
		}
		bn, err := chain.Confirm(ccsTx)
		return &changeset.ContractDeployment[*channel_config_store.ChannelConfigStore]{
			Address:  ccsAddr,
			Block:    bn,
			Contract: ccs,
			Tx:       ccsTx,
			Tv:       deployment.NewTypeAndVersion(types.ChannelConfigStore, deployment.Version1_0_0),
			Err:      nil,
		}
	}
}
