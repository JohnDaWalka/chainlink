package jobs

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/jobs"
)

type DistributeOCRJobSpecSeqDeps struct {
	Nodes    []*nodev1.Node
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
		nodesByID := make(map[string]*nodev1.Node)
		for _, node := range deps.Nodes {
			nodesByID[node.Id] = node
		}

		specs, err := jobs.BuildOCR3JobConfigSpecs(
			deps.Offchain, b.Logger, input.ContractID, input.ChainSelectorEVM, input.ChainSelectorAptos, deps.Nodes, input.BootstrapperOCR3Urls, input.DONName)
		if err != nil {
			return DistributeOCRJobSpecSeqOutput{}, fmt.Errorf("failed to build job specs: %w", err)
		}

		var mergedErrs error
		for _, spec := range specs {
			node, ok := nodesByID[spec.NodeID]
			if !ok {
				return DistributeOCRJobSpecSeqOutput{}, fmt.Errorf("node not found: %s", spec.NodeID)
			}

			_, opErr := operations.ExecuteOperation(b, DistributeOCRJobSpecOp, DistributeOCRJobSpecOpDeps{
				Node:     node,
				Offchain: deps.Offchain,
			}, DistributeOCRJobSpecOpInput{
				DomainKey:        input.DomainKey,
				EnvironmentLabel: input.EnvironmentLabel,
				Spec:             spec,
			})
			if opErr != nil {
				// Do not fail changeset if a single proposal fails, make it through all proposals.
				mergedErrs = fmt.Errorf("error proposing job to node %s spec %s: %w", node.Id, spec.Spec, opErr)
				continue
			}
		}

		return DistributeOCRJobSpecSeqOutput{
			Specs: specs,
		}, mergedErrs
	},
)
