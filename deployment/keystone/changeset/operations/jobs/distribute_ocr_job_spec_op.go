package jobs

import (
	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/jobs"
)

type DistributeOCRJobSpecOpDeps struct {
	NodeID   string // Node ID to distribute the job spec to
	Offchain deployment.OffchainClient
}

type DistributeOCRJobSpecOpInput struct {
	DomainKey        string
	EnvironmentLabel string
	Spec             jobs.OCR3JobConfigSpec
}

type DistributeOCRJobSpecOpOutput struct {
	Spec jobs.OCR3JobConfigSpec
}

var DistributeOCRJobSpecOp = operations.NewOperation[DistributeOCRJobSpecOpInput, DistributeOCRJobSpecOpOutput, DistributeOCRJobSpecOpDeps](
	"distribute-ocr-job-spec-op",
	semver.MustParse("1.0.0"),
	"Distribute OCR Job Spec",
	func(b operations.Bundle, deps DistributeOCRJobSpecOpDeps, input DistributeOCRJobSpecOpInput) (DistributeOCRJobSpecOpOutput, error) {
		b.Logger.Debugw("Proposing job", "nodeID", deps.NodeID, "domain", input.DomainKey, "environment", input.EnvironmentLabel)
		req := jobs.ProposeJobRequest{
			Job:            input.Spec.Spec,
			DomainKey:      input.DomainKey,
			Environment:    input.EnvironmentLabel,
			NodeIDs:        []string{deps.NodeID},
			OffchainClient: deps.Offchain,
			Lggr:           b.Logger,
		}
		err := jobs.ProposeJob(b.GetContext(), req)
		if err != nil {
			return DistributeOCRJobSpecOpOutput{}, err
		}
		return DistributeOCRJobSpecOpOutput{
			Spec: input.Spec,
		}, nil
	},
)
