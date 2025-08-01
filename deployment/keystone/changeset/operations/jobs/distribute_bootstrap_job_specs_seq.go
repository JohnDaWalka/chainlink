package jobs

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/jobs"
)

const (
	DefaultBootstrapJobName = "OCR3 MultiChain Capability Bootstrap"
)

type DistributeBootstrapJobSpecsSeqDeps struct {
	Offchain deployment.OffchainClient
}

type DistributeBootstrapJobSpecsSeqInput struct {
	DONName          string
	DomainKey        string
	ContractID       string
	EnvironmentLabel string
	ChainSelectorEVM uint64

	JobName  string // Optional job name, if not provided, the default will be used.
	BootCfgs []DistributeBootstrapJobSpecsSeqBootCfg
}

type DistributeBootstrapJobSpecsSeqBootCfg struct {
	P2PID    string
	NodeName string
}

type DistributeBootstrapJobSpecsSeqOutput struct {
	Spec string
}

var DistributeBootstrapJobSpecsSeq = operations.NewSequence[DistributeBootstrapJobSpecsSeqInput, DistributeBootstrapJobSpecsSeqOutput, DistributeBootstrapJobSpecsSeqDeps](
	"distribute-bootstrap-job-specs-seq",
	semver.MustParse("1.0.0"),
	"Distribute Bootstrap Job JobSpecs",
	func(b operations.Bundle, deps DistributeBootstrapJobSpecsSeqDeps, input DistributeBootstrapJobSpecsSeqInput) (DistributeBootstrapJobSpecsSeqOutput, error) {
		extJobID, err := jobs.BootstrapExternalJobID(input.DONName, input.ChainSelectorEVM)
		if err != nil {
			return DistributeBootstrapJobSpecsSeqOutput{}, fmt.Errorf("failed to get external job ID: %w", err)
		}

		chainID, err := chainsel.GetChainIDFromSelector(input.ChainSelectorEVM)
		if err != nil {
			return DistributeBootstrapJobSpecsSeqOutput{}, fmt.Errorf("failed to get chain ID from selector: %w", err)
		}

		jobName := DefaultBootstrapJobName
		if input.JobName != "" {
			jobName = input.JobName
		}
		spec, err := jobs.ResolveBootstrapJob(jobs.BootstrapCfg{
			JobName:       fmt.Sprintf(jobName+" (for DON %s)", input.DONName),
			ExternalJobID: extJobID,
			ChainID:       chainID,
			ContractID:    input.ContractID,
		})
		if err != nil {
			return DistributeBootstrapJobSpecsSeqOutput{}, fmt.Errorf("failed to resolve bootstrap job: %w", err)
		}

		if len(input.BootCfgs) == 0 {
			return DistributeBootstrapJobSpecsSeqOutput{}, fmt.Errorf("no bootstrap configurations provided")
		}

		var mergedErrs error
		for _, bootCfg := range input.BootCfgs {
			_, opErr := operations.ExecuteOperation(b, DistributeJobSpecOp, DistributeJobSpecOpDeps{
				Offchain: deps.Offchain,
			}, DistributeJobSpecOpInput{
				Spec:             spec,
				NodeP2PLabel:     bootCfg.P2PID,
				NodeName:         bootCfg.NodeName,
				DomainKey:        input.DomainKey,
				EnvironmentLabel: input.EnvironmentLabel,
			})
			if opErr != nil {
				// Do not fail the sequence if a single proposal fails, make it through all proposals.
				mergedErrs = fmt.Errorf("error proposing bootstrap job to node %s spec %s: %w", bootCfg.NodeName, spec, opErr)
				continue
			}
		}

		return DistributeBootstrapJobSpecsSeqOutput{
			Spec: spec,
		}, mergedErrs
	},
)
