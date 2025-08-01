package jobs

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/jobs"
)

type DistributeOCRJobSpecSeqDeps struct {
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
	Nodes                []DistributeOCRJobSpecSeqNode
}

type DistributeOCRJobSpecSeqNode struct {
	ID       string
	P2PLabel string
}

type DistributeOCRJobSpecSeqOutput struct {
	Specs []jobs.OCR3JobConfigSpec
}

var DistributeOCRJobSpecSeq = operations.NewSequence[DistributeOCRJobSpecSeqInput, DistributeOCRJobSpecSeqOutput, DistributeOCRJobSpecSeqDeps](
	"distribute-ocr-job-spec-seq",
	semver.MustParse("1.0.0"),
	"Distribute OCR Job JobSpecs",
	func(b operations.Bundle, deps DistributeOCRJobSpecSeqDeps, input DistributeOCRJobSpecSeqInput) (DistributeOCRJobSpecSeqOutput, error) {
		nodeIdToP2PLabel := make(map[string]string)
		nodeIDs := make([]string, len(input.Nodes))
		for _, node := range input.Nodes {
			nodeIDs = append(nodeIDs, node.ID)
			nodeIdToP2PLabel[node.ID] = node.P2PLabel
		}

		specs, err := jobs.BuildOCR3JobConfigSpecs(
			deps.Offchain, b.Logger, input.ContractID, input.ChainSelectorEVM, input.ChainSelectorAptos, nodeIDs, input.BootstrapperOCR3Urls, input.DONName)
		if err != nil {
			return DistributeOCRJobSpecSeqOutput{}, fmt.Errorf("failed to build job specs: %w", err)
		}

		var mergedErrs error
		for _, spec := range specs {
			nodeLabel, ok := nodeIdToP2PLabel[spec.NodeID]
			if !ok {
				return DistributeOCRJobSpecSeqOutput{}, fmt.Errorf("node not found: %s", spec.NodeID)
			}

			_, opErr := operations.ExecuteOperation(b, DistributeJobSpecOp, DistributeJobSpecOpDeps{
				Offchain: deps.Offchain,
			}, DistributeJobSpecOpInput{
				NodeID:           spec.NodeID,
				NodeP2PLabel:     nodeLabel,
				DomainKey:        input.DomainKey,
				EnvironmentLabel: input.EnvironmentLabel,
				Spec:             spec.Spec,
			})
			if opErr != nil {
				// Do not fail the sequence if a single proposal fails, make it through all proposals.
				mergedErrs = fmt.Errorf("error proposing job to node %s spec %s: %w", spec.NodeID, spec.Spec, opErr)
				continue
			}
		}

		return DistributeOCRJobSpecSeqOutput{
			Specs: specs,
		}, mergedErrs
	},
)
