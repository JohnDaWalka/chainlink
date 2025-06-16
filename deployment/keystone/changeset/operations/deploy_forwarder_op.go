package operations

import (
	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

type DeployForwarderOpDeps struct{}

type DeployForwarderOpInput struct {
}

type DeployForwarderOpOutput struct {
}

// DeployKeystoneForwarderOp is an operation that deploys the Keystone Forwarder contract.
var DeployKeystoneForwarderOp = operations.NewOperation[DeployForwarderOpInput, DeployForwarderOpOutput, DeployForwarderOpDeps](
	"deploy-keystone-forwarder-op",
	semver.MustParse("1.0.0"),
	"Deploy KeystoneForwarder Contract",
	func(b operations.Bundle, deps DeployForwarderOpDeps, input DeployForwarderOpInput) (DeployForwarderOpOutput, error) {
		// Here we would implement the logic to deploy the Keystone Forwarder contract.
		return DeployForwarderOpOutput{}, nil
	},
)
