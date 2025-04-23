package reward_manager

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/metadata"
	ds "github.com/smartcontractkit/chainlink/deployment/datastore"
	"github.com/smartcontractkit/mcms"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/mcmsutil"

	rewardManager "github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/reward_manager_v0_5_0"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

var DeployRewardManagerChangeset = deployment.CreateChangeSet(deployRewardManagerLogic, deployRewardManagerPrecondition)

type DeployRewardManager struct {
	LinkTokenAddress common.Address
}

type DeployRewardManagerConfig struct {
	ChainsToDeploy map[uint64]DeployRewardManager
	Ownership      types.OwnershipSettings
}

func (cc DeployRewardManagerConfig) Validate() error {
	if len(cc.ChainsToDeploy) == 0 {
		return errors.New("ChainsToDeploy is empty")
	}
	for chain := range cc.ChainsToDeploy {
		if err := deployment.IsValidChainSelector(chain); err != nil {
			return fmt.Errorf("invalid chain selector: %d - %w", chain, err)
		}
	}
	return nil
}

func deployRewardManagerLogic(e deployment.Environment, cc DeployRewardManagerConfig) (deployment.ChangesetOutput, error) {
	dataStore := ds.NewMemoryDataStore[
		metadata.SerializedContractMetadata,
		ds.DefaultMetadata,
	]()
	err := deployRewardManager(e, dataStore, cc)
	if err != nil {
		e.Logger.Errorw("Failed to deploy RewardManager", "err", err)
		return deployment.ChangesetOutput{}, deployment.MaybeDataErr(err)
	}

	var proposals []mcms.TimelockProposal
	if cc.Ownership.ShouldTransfer && cc.Ownership.MCMSProposalConfig != nil {
		filter := deployment.NewTypeAndVersion(types.RewardManager, deployment.Version0_5_0)
		res, err := mcmsutil.TransferToMCMSWithTimelockForTypeAndVersion(e, dataStore, filter, *cc.Ownership.MCMSProposalConfig)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to transfer ownership to MCMS: %w", err)
		}
		proposals = res.MCMSTimelockProposals
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

func deployRewardManagerPrecondition(_ deployment.Environment, cc DeployRewardManagerConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid DeployRewardManagerConfig: %w", err)
	}

	return nil
}

func deployRewardManager(e deployment.Environment,
	dataStore ds.MutableDataStore[metadata.SerializedContractMetadata, ds.DefaultMetadata],
	cc DeployRewardManagerConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid DeployRewardManagerConfig: %w", err)
	}

	for chainSel, chainCfg := range cc.ChainsToDeploy {
		chain, ok := e.Chains[chainSel]
		if !ok {
			return fmt.Errorf("chain not found for chain selector %d", chainSel)
		}
		rewardManagerMetadata := metadata.RewardManagerMetadata{}
		serialized, err := metadata.NewRewardManagerMetadata(rewardManagerMetadata)
		if err != nil {
			return fmt.Errorf("failed to serialize verifier proxy metadata: %w", err)
		}

		_, err = changeset.DeployContractV2(e, dataStore, serialized, chain, RewardManagerDeployFn(chainCfg.LinkTokenAddress))
		if err != nil {
			return err
		}
	}

	return nil
}

func RewardManagerDeployFn(linkAddress common.Address) changeset.ContractDeployFn[*rewardManager.RewardManager] {
	return func(chain deployment.Chain) *changeset.ContractDeployment[*rewardManager.RewardManager] {
		ccsAddr, ccsTx, ccs, err := rewardManager.DeployRewardManager(
			chain.DeployerKey,
			chain.Client,
			linkAddress,
		)
		if err != nil {
			return &changeset.ContractDeployment[*rewardManager.RewardManager]{
				Err: err,
			}
		}
		return &changeset.ContractDeployment[*rewardManager.RewardManager]{
			Address:  ccsAddr,
			Contract: ccs,
			Tx:       ccsTx,
			Tv:       deployment.NewTypeAndVersion(types.RewardManager, deployment.Version0_5_0),
			Err:      nil,
		}
	}
}
