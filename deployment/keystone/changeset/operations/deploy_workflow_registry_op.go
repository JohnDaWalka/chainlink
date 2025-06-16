package operations

import (
	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

type DeployWorkflowRegistryOpDeps struct {
}

type DeployWorkflowRegistryInput struct {
}

type DeployWorkflowRegistryOutput struct {
}

// DeployWorkflowRegistryOp is an operation that deploys the Workflow Registry contract.
var DeployWorkflowRegistryOp = operations.NewOperation[DeployWorkflowRegistryInput, DeployWorkflowRegistryOutput, DeployWorkflowRegistryOpDeps](
	"deploy-workflow-registry-op",
	semver.MustParse("1.0.0"),
	"Deploy WorkflowRegistry Contract",
	func(b operations.Bundle, deps DeployWorkflowRegistryOpDeps, input DeployWorkflowRegistryInput) (DeployWorkflowRegistryOutput, error) {
		// Here we would implement the logic to deploy the Workflow Registry contract.
		return DeployWorkflowRegistryOutput{}, nil
	},
)
