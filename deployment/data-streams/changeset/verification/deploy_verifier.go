package verification

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/metadata"
	ds "github.com/smartcontractkit/chainlink/deployment/datastore"
	"github.com/smartcontractkit/mcms"

	verifier "github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/verifier_v0_5_0"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/mcmsutil"
)

var DeployVerifierChangeset = deployment.CreateChangeSet(deployVerifierLogic, deployVerifierPrecondition)

type DeployVerifier struct {
	VerifierProxyAddress common.Address
}

type DeployVerifierConfig struct {
	ChainsToDeploy map[uint64]DeployVerifier
	Ownership      types.OwnershipSettings
}

func (cc DeployVerifierConfig) Validate() error {
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

func deployVerifierLogic(e deployment.Environment, cc DeployVerifierConfig) (deployment.ChangesetOutput, error) {
	dataStore := ds.NewMemoryDataStore[
		metadata.SerializedContractMetadata,
		ds.DefaultMetadata,
	]()
	err := deployVerifier(e, dataStore, cc)
	if err != nil {
		e.Logger.Errorw("Failed to deploy Verifier", "err", err)
		return deployment.ChangesetOutput{}, deployment.MaybeDataErr(err)
	}

	var proposals []mcms.TimelockProposal
	if cc.Ownership.ShouldTransfer && cc.Ownership.MCMSProposalConfig != nil {
		filter := deployment.NewTypeAndVersion(types.Verifier, deployment.Version0_5_0)
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

func deployVerifierPrecondition(_ deployment.Environment, cc DeployVerifierConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid DeployVerifierConfig: %w", err)
	}

	return nil
}

func deployVerifier(e deployment.Environment, dataStore ds.MutableDataStore[metadata.SerializedContractMetadata, ds.DefaultMetadata], cc DeployVerifierConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid DeployVerifierConfig: %w", err)
	}

	for chainSel, chainCfg := range cc.ChainsToDeploy {
		chain, ok := e.Chains[chainSel]
		if !ok {
			return fmt.Errorf("chain not found for chain selector %d", chainSel)
		}
		verifierMetadata := metadata.VerifierMetadata{
			VerifierProxyAddress: chainCfg.VerifierProxyAddress.String(),
		}
		serialized, err := metadata.NewVerifierMetadata(verifierMetadata)
		if err != nil {
			return fmt.Errorf("failed to serialize verifier metadata: %w", err)
		}
		_, err = changeset.DeployContractV2(e, dataStore, serialized, chain, VerifierDeployFn(chainCfg.VerifierProxyAddress))
		if err != nil {
			return fmt.Errorf("failed to deploy verifier: %w", err)
		}
	}

	return nil
}

func VerifierDeployFn(verifierProxyAddress common.Address) changeset.ContractDeployFn[*verifier.Verifier] {
	return func(chain deployment.Chain) *changeset.ContractDeployment[*verifier.Verifier] {
		ccsAddr, ccsTx, ccs, err := verifier.DeployVerifier(
			chain.DeployerKey,
			chain.Client,
			verifierProxyAddress,
		)
		if err != nil {
			return &changeset.ContractDeployment[*verifier.Verifier]{
				Err: err,
			}
		}
		return &changeset.ContractDeployment[*verifier.Verifier]{
			Address:  ccsAddr,
			Contract: ccs,
			Tx:       ccsTx,
			Tv:       deployment.NewTypeAndVersion(types.Verifier, deployment.Version0_5_0),
			Err:      nil,
		}
	}
}
