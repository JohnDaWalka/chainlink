package aptos

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	seq "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/sequence"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
	"github.com/smartcontractkit/mcms"
	aptosmcms "github.com/smartcontractkit/mcms/sdk/aptos"
	mcmstypes "github.com/smartcontractkit/mcms/types"
)

var _ deployment.ChangeSetV2[config.AddTokenPoolConfig] = AddTokenPool{}

// AddTokenPool deploys token pools and sets up tokens on lanes
type AddTokenPool struct{}

func (cs AddTokenPool) VerifyPreconditions(env deployment.Environment, cfg config.AddTokenPoolConfig) error {
	return nil
}

func (cs AddTokenPool) Apply(env deployment.Environment, cfg config.AddTokenPoolConfig) (deployment.ChangesetOutput, error) {
	state, err := changeset.LoadOnchainState(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load Aptos onchain state: %w", err)
	}

	aptosChain := env.AptosChains[cfg.ChainSelector]
	ab := deployment.NewMemoryAddressBook()
	seqReports := make([]operations.Report[any, any], 0)
	proposals := make([]mcms.TimelockProposal, 0)

	mcmsOperations := []mcmstypes.BatchOperation{}

	deps := operation.AptosDeps{
		AB:               ab,
		AptosChain:       aptosChain,
		CCIPOnChainState: state,
	}

	// Deploy Aptos token and pool
	depInput := seq.DeployTokenPoolSeqInput{}
	deploySeq, err := operations.ExecuteSequence(env.OperationsBundle, seq.DeployAptosTokenPoolSequence, deps, depInput)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	seqReports = append(seqReports, deploySeq.ExecutionReports...)
	mcmsOperations = append(mcmsOperations, deploySeq.Output.MCMSOperations...)

	// Connect token pools EVM -> Aptos
	connInput := seq.ConnectTokenPoolSeqInput{
		TokenPoolAddress: deploySeq.Output.TokenPoolAddress,
		PoolType:         cfg.PoolType,
		RemotePools:      toRemotePools(cfg.RemoteChainTokenPoolConfig),
	}
	connectSeq, err := operations.ExecuteSequence(env.OperationsBundle, seq.ConnectTokenPoolSequence, deps, connInput)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	seqReports = append(seqReports, connectSeq.ExecutionReports...)
	mcmsOperations = append(mcmsOperations, deploySeq.Output.MCMSOperations...)

	// Generate Aptos MCMS proposals
	proposal, err := utils.GenerateProposal(
		aptosChain.Client,
		state.AptosChains[cfg.ChainSelector].MCMSAddress,
		cfg.ChainSelector,
		mcmsOperations,
		"Deploy Aptos MCMS and CCIP",
		aptosmcms.TimelockRoleProposer,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate MCMS proposal for Aptos chain %d: %w", cfg.ChainSelector, err)
	}
	proposals = append(proposals, *proposal)

	return deployment.ChangesetOutput{
		AddressBook:           ab,
		MCMSTimelockProposals: proposals,
		Reports:               seqReports,
	}, nil
}

func toRemotePools(cfg config.RemoteChainTokenPoolConfig) map[uint64]operation.RemotePool {
	remotePools := make(map[uint64]operation.RemotePool)
	for chainSelector, remoteConfig := range cfg.EVMRemoteConfigs {
		remotePools[chainSelector] = operation.RemotePool{
			RemotePoolAddress:  remoteConfig.TokenPoolAddress.Bytes(),
			RemoteTokenAddress: remoteConfig.TokenAddress.Bytes(),
			RateLimiterConfig:  remoteConfig.RateLimiterConfig,
		}
	}
	return remotePools
}
