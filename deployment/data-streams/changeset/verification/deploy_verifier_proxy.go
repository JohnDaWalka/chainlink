package verification

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/metadata"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/mcmsutil"
	ds "github.com/smartcontractkit/chainlink/deployment/datastore"
	"github.com/smartcontractkit/mcms"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/verifier_proxy_v0_5_0"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

// DeployVerifierProxyChangeset deploys VerifierProxy to the chains specified in the config.
var DeployVerifierProxyChangeset deployment.ChangeSetV2[DeployVerifierProxyConfig] = &verifierProxyDeploy{}

type verifierProxyDeploy struct{}
type DeployVerifierProxyConfig struct {
	// ChainsToDeploy is a list of chain selectors to deploy the contract to.
	ChainsToDeploy map[uint64]DeployVerifierProxy
	Ownership      types.OwnershipSettings
	Version        semver.Version
}

type DeployVerifierProxy struct {
	AccessControllerAddress common.Address
}

func (cfg DeployVerifierProxyConfig) Validate() error {
	switch cfg.Version {
	case deployment.Version0_5_0:
		// no-op
	default:
		return fmt.Errorf("unsupported contract version %s", cfg.Version)
	}
	if len(cfg.ChainsToDeploy) == 0 {
		return errors.New("ChainsToDeploy is empty")
	}
	for chain := range cfg.ChainsToDeploy {
		if err := deployment.IsValidChainSelector(chain); err != nil {
			return fmt.Errorf("invalid chain selector: %d - %w", chain, err)
		}
	}
	return nil
}

func (v *verifierProxyDeploy) Apply(e deployment.Environment, cc DeployVerifierProxyConfig) (deployment.ChangesetOutput, error) {
	// Create an in-memory data store to store the address references, contract metadata fow newly deployed contracts
	dataStore := ds.NewMemoryDataStore[
		metadata.SerializedContractMetadata,
		ds.DefaultMetadata,
	]()

	err := deploy(e, dataStore, cc)
	if err != nil {
		e.Logger.Errorw("Failed to deploy VerifierProxy", "err", err)
		return deployment.ChangesetOutput{}, deployment.MaybeDataErr(err)
	}

	var proposals []mcms.TimelockProposal
	if cc.Ownership.ShouldTransfer && cc.Ownership.MCMSProposalConfig != nil {
		filter := deployment.NewTypeAndVersion(types.VerifierProxy, deployment.Version0_5_0)
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

func (v *verifierProxyDeploy) VerifyPreconditions(_ deployment.Environment, cc DeployVerifierProxyConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid DeployVerifierProxyConfig: %w", err)
	}
	return nil
}

func deploy(e deployment.Environment,
	dataStore ds.MutableDataStore[metadata.SerializedContractMetadata, ds.DefaultMetadata],
	cfg DeployVerifierProxyConfig) error {
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid DeployVerifierProxyConfig: %w", err)
	}

	for chainSel, chainCfg := range cfg.ChainsToDeploy {
		chain, ok := e.Chains[chainSel]
		if !ok {
			return fmt.Errorf("chain not found for chain selector %d", chainSel)
		}
		verifierProxyMetadata := metadata.VerifierProxyMetadata{
			AccessControllerAddress: chainCfg.AccessControllerAddress.String(),
		}
		serialized, err := metadata.NewVerifierProxyMetadata(verifierProxyMetadata)
		if err != nil {
			return fmt.Errorf("failed to serialize verifier proxy metadata: %w", err)
		}
		_, err = changeset.DeployContractV2(e, dataStore, serialized, chain, verifyProxyDeployFn(chainCfg))
		if err != nil {
			return fmt.Errorf("failed to deploy verifier proxy: %w", err)
		}
	}

	return nil
}

// verifyProxyDeployFn returns a function that deploys a VerifyProxy contract.
func verifyProxyDeployFn(cfg DeployVerifierProxy) changeset.ContractDeployFn[*verifier_proxy_v0_5_0.VerifierProxy] {
	return func(chain deployment.Chain) *changeset.ContractDeployment[*verifier_proxy_v0_5_0.VerifierProxy] {
		addr, tx, contract, err := verifier_proxy_v0_5_0.DeployVerifierProxy(
			chain.DeployerKey,
			chain.Client,
			cfg.AccessControllerAddress,
		)
		if err != nil {
			return &changeset.ContractDeployment[*verifier_proxy_v0_5_0.VerifierProxy]{
				Err: err,
			}
		}
		return &changeset.ContractDeployment[*verifier_proxy_v0_5_0.VerifierProxy]{
			Address:  addr,
			Contract: contract,
			Tx:       tx,
			Tv:       deployment.NewTypeAndVersion(types.VerifierProxy, deployment.Version0_5_0),
			Err:      nil,
		}
	}
}
