package jobs

import (
	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/jobs/offchain"

	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/jobs"
)

type DistributeJobSpecOpDeps struct {
	Offchain deployment.OffchainClient
}

type DistributeJobSpecOpInput struct {
	DomainKey        string
	EnvironmentLabel string
	// NodeID is optional. On CLD, the loaded bootstrap config has no IDs, only names.
	// See: https://github.com/smartcontractkit/chainlink-deployments/blob/main/domains/keystone/staging/inputs/bootstrappers.json
	NodeID       string
	NodeName     string
	NodeP2PLabel string
	Spec         string
}

type DistributeJobSpecOpOutput struct {
	Spec string
}

var DistributeJobSpecOp = operations.NewOperation[DistributeJobSpecOpInput, DistributeJobSpecOpOutput, DistributeJobSpecOpDeps](
	"distribute-job-spec-op",
	semver.MustParse("1.0.0"),
	"Distribute Job Spec",
	func(b operations.Bundle, deps DistributeJobSpecOpDeps, input DistributeJobSpecOpInput) (DistributeJobSpecOpOutput, error) {
		p2pID := input.NodeP2PLabel
		b.Logger.Debugw("Proposing job", "nodeName", input.NodeName, "nodeID", input.NodeID, "domain", input.DomainKey, "environment", input.EnvironmentLabel)
		req := jobs.ProposeJobRequest{
			Job:            input.Spec,
			DomainKey:      input.DomainKey,
			Environment:    input.EnvironmentLabel,
			NodeLabels:     map[string]string{offchain.P2pIDLabel: p2pID},
			OffchainClient: deps.Offchain,
			Lggr:           b.Logger,
		}
		if input.NodeID != "" {
			req.NodeIDs = []string{input.NodeID}
		}
		err := jobs.ProposeJob(b.GetContext(), req)
		if err != nil {
			return DistributeJobSpecOpOutput{}, err
		}
		return DistributeJobSpecOpOutput{
			Spec: input.Spec,
		}, nil
	},
)
