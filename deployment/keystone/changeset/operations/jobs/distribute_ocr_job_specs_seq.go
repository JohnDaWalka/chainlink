package jobs

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/jobs"
)

type DistributeOCRJobSpecSeqDeps struct {
	NodeIDs  []string
	Offchain deployment.OffchainClient
}

type DistributeOCRJobSpecSeqInput struct {
	DomainKey            string
	EnvironmentLabel     string
	DONName              string
	ContractID           string
	ChainSelectorEVM     uint64
	ChainSelectorAptos   uint64
	BootstrapperOCR3Urls []string
}

type DistributeOCRJobSpecSeqOutput struct {
	Specs []jobs.OCR3JobConfigSpec
}

var DistributeOCRJobSpecSeq = operations.NewSequence[DistributeOCRJobSpecSeqInput, DistributeOCRJobSpecSeqOutput, DistributeOCRJobSpecSeqDeps](
	"distribute-ocr-job-spec-seq",
	semver.MustParse("1.0.0"),
	"Distribute OCR Job Specs",
	func(b operations.Bundle, deps DistributeOCRJobSpecSeqDeps, input DistributeOCRJobSpecSeqInput) (DistributeOCRJobSpecSeqOutput, error) {
		nodesByID := make(map[string]struct{})
		for _, nodeID := range deps.NodeIDs {
			nodesByID[nodeID] = struct{}{}
		}

		specs, err := jobs.BuildOCR3JobConfigSpecs(
			deps.Offchain, b.Logger, input.ContractID, input.ChainSelectorEVM, input.ChainSelectorAptos, deps.NodeIDs, input.BootstrapperOCR3Urls, input.DONName)
		if err != nil {
			return DistributeOCRJobSpecSeqOutput{}, fmt.Errorf("failed to build job specs: %w", err)
		}

		var mergedErrs error
		for _, spec := range specs {
			_, ok := nodesByID[spec.NodeID]
			if !ok {
				return DistributeOCRJobSpecSeqOutput{}, fmt.Errorf("node not found: %s", spec.NodeID)
			}

			_, opErr := operations.ExecuteOperation(b, DistributeOCRJobSpecOp, DistributeOCRJobSpecOpDeps{
				NodeID:   spec.NodeID,
				Offchain: deps.Offchain,
			}, DistributeOCRJobSpecOpInput{
				DomainKey:        input.DomainKey,
				EnvironmentLabel: input.EnvironmentLabel,
				Spec:             spec,
			})
			if opErr != nil {
				// Do not fail changeset if a single proposal fails, make it through all proposals.
				mergedErrs = fmt.Errorf("error proposing job to node %s spec %s: %w", spec.NodeID, spec.Spec, opErr)
				continue
			}
		}

		return DistributeOCRJobSpecSeqOutput{
			Specs: specs,
		}, mergedErrs
	},
)
