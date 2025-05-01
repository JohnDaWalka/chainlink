package v0_5_0

import (
	"errors"
	"fmt"

	ds "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/configurator"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/metadata"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/mcmsutil"
)

var DeployConfiguratorChangeset = deployment.CreateChangeSet(deployConfiguratorLogic, deployConfiguratorPrecondition)

type DeployConfiguratorConfig struct {
	ChainsToDeploy []uint64
	Ownership      types.OwnershipSettings
}

func (cc DeployConfiguratorConfig) GetOwnershipConfig() types.OwnershipSettings {
	return cc.Ownership
}

func (cc DeployConfiguratorConfig) Validate() error {
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

func deployConfiguratorLogic(e deployment.Environment, cc DeployConfiguratorConfig) (deployment.ChangesetOutput, error) {
	dataStore := ds.NewMemoryDataStore[
		metadata.SerializedContractMetadata,
		ds.DefaultMetadata,
	]()

	err := deploy(e, dataStore, cc)
	if err != nil {
		e.Logger.Errorw("Failed to deploy Configurator", "err", err)
		return deployment.ChangesetOutput{}, deployment.MaybeDataErr(err)
	}

	proposals, err := mcmsutil.GetTransferOwnershipProposals(e, cc, dataStore, deployment.NewTypeAndVersion(types.Configurator, deployment.Version0_5_0))
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

func deployConfiguratorPrecondition(_ deployment.Environment, cc DeployConfiguratorConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid DeployConfiguratorConfig: %w", err)
	}

	return nil
}

func deploy(e deployment.Environment, dataStore ds.MutableDataStore[metadata.SerializedContractMetadata, ds.DefaultMetadata], cc DeployConfiguratorConfig) error {
	for _, chainSel := range cc.ChainsToDeploy {
		chain, ok := e.Chains[chainSel]
		if !ok {
			return fmt.Errorf("chain not found for chain selector %d", chainSel)
		}
		md, err := metadata.NewConfiguratorMetadata(metadata.ConfiguratorMetadata{})
		if err != nil {
			return fmt.Errorf("failed to create metadata: %w", err)
		}
		options := &changeset.DeployOptions{ContractMetadata: md}
		_, err = changeset.DeployContract(e, dataStore, chain, DeployFn(), options)
		if err != nil {
			return fmt.Errorf("failed to deploy configurator: %w", err)
		}
	}
	return nil
}

func DeployFn() changeset.ContractDeployFn[*configurator.Configurator] {
	return func(chain deployment.Chain) *changeset.ContractDeployment[*configurator.Configurator] {
		ccsAddr, ccsTx, ccs, err := configurator.DeployConfigurator(
			chain.DeployerKey,
			chain.Client,
		)
		if err != nil {
			return &changeset.ContractDeployment[*configurator.Configurator]{
				Err: err,
			}
		}
		return &changeset.ContractDeployment[*configurator.Configurator]{
			Address:  ccsAddr,
			Contract: ccs,
			Tx:       ccsTx,
			Tv:       deployment.NewTypeAndVersion(types.Configurator, deployment.Version0_5_0),
			Err:      nil,
		}
	}
}
