package operations

import (
	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

type DeployCapabilityRegistryOpDeps struct {
}

type DeployCapabilityRegistryInput struct {
}

type DeployCapabilityRegistryOutput struct {
}

// DeployCapabilityRegistryOp is an operation that deploys the Capability Registry contract.
var DeployCapabilityRegistryOp = operations.NewOperation[DeployCapabilityRegistryInput, DeployCapabilityRegistryOutput, DeployCapabilityRegistryOpDeps](
	"deploy-capability-registry-op",
	semver.MustParse("1.0.0"),
	"Deploy CapabilityRegistry Contract",
	func(b operations.Bundle, deps DeployCapabilityRegistryOpDeps, input DeployCapabilityRegistryInput) (DeployCapabilityRegistryOutput, error) {
		// Here we would implement the logic to deploy the OCR3 contract.
		return DeployCapabilityRegistryOutput{}, nil
	},
)
