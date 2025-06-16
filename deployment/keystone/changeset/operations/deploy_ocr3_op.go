package operations

import (
	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

type DeployOCR3OpDeps struct {
}

type DeployOCR3OpInput struct {
}

type DeployOCR3OpOutput struct {
}

// DeployOCR3Op is an operation that deploys the OCR3 contract.
var DeployOCR3Op = operations.NewOperation[DeployOCR3OpInput, DeployOCR3OpOutput, DeployOCR3OpDeps](
	"deploy-ocr3-op",
	semver.MustParse("1.0.0"),
	"Deploy OCR3 Contract",
	func(b operations.Bundle, deps DeployOCR3OpDeps, input DeployOCR3OpInput) (DeployOCR3OpOutput, error) {
		// Here we would implement the logic to deploy the OCR3 contract.
		return DeployOCR3OpOutput{}, nil
	},
)
