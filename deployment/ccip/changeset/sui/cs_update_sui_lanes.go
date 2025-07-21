package sui

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	seq "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/sequence"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

// TODO
var _ cldf.ChangeSetV2[config.UpdateAptosLanesConfig] = AddAptosLanes{}

// AddAptosLanes implements adding a new lane to an existing Aptos CCIP deployment
type AddAptosLanes struct{}

func (cs AddAptosLanes) VerifyPreconditions(env cldf.Environment, cfg config.UpdateAptosLanesConfig) error {
	return nil
}

// TODO
func (cs AddAptosLanes) Apply(env cldf.Environment, cfg config.UpdateAptosLanesConfig) (cldf.ChangesetOutput, error) {

	seqReports := make([]operations.Report[any, any], 0)

	// Add lane on EVM chains
	// TODO: applying a changeset within another changeset is an anti-pattern. Using it here until EVM is refactored into Operations
	evmUpdatesInput := config.ToEVMUpdateLanesConfig(cfg)
	_, err := v1_6.UpdateLanesLogic(env, cfg.EVMMCMSConfig, evmUpdatesInput)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	// Add lane on Aptos chains
	// Execute UpdateAptosLanesSequence for each aptos chain
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load Aptos onchain state: %w", err)
	}

	updateInputsByAptosChain := seq.ToAptosUpdateLanesConfig(state.AptosChains, cfg)
	for aptosChainSel, sequenceInput := range updateInputsByAptosChain {
		deps := operation.AptosDeps{
			AptosChain:       env.BlockChains.AptosChains()[aptosChainSel],
			CCIPOnChainState: state,
		}
		// Execute the sequence
		updateSeqReport, err := operations.ExecuteSequence(env.OperationsBundle, seq.UpdateAptosLanesSequence, deps, sequenceInput)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		seqReports = append(seqReports, updateSeqReport.ExecutionReports...)
	}

	return cldf.ChangesetOutput{
		Reports: seqReports,
	}, nil
}
