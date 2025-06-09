package sequence

import (
	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/sui/operation"

	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

type DeployAndInitCCIPOffRampSeqInput struct {
	operation.DeployCCIPOffRampInput
	operation.InitializeOffRampInput
	operation.SetOCR3ConfigInput
}

type DeployCCIPOffRampSeqObjects struct {
	ObjectCapId   string
	StateObjectId string
}

type DeployCCIPOffRampSeqOutput struct {
	CCIPOffRampPackageId string
	Objects              DeployCCIPOffRampSeqObjects
}

var DeployAndInitCCIPOffRampSequence = cld_ops.NewSequence(
	"sui-deploy-ccip-offramp-seq",
	semver.MustParse("0.1.0"),
	"Deploys and sets initial CCIP offRamp configuration",
	func(env cld_ops.Bundle, deps sui_ops.OpTxDeps, input DeployAndInitCCIPOffRampSeqInput) (DeployCCIPOffRampSeqOutput, error) {
		deployReport, err := cld_ops.ExecuteOperation(env, operation.DeployCCIPOffRampOp, deps, input.DeployCCIPOffRampInput)
		if err != nil {
			return DeployCCIPOffRampSeqOutput{}, err
		}

		// input.InitializeOffRampInput.OffRampPackageId = deployReport.Output.PackageId
		// input.InitializeOffRampInput.OwnerCapObjectId = deployReport.Output.Objects.OwnerCapObjectId
		// input.InitializeOffRampInput.OffRampStateId = deployReport.Output.Objects.CCIPOffRampStateObjectId

		// _, err = cld_ops.ExecuteOperation(env, operation.InitializeOffRampOp, deps, input.InitializeOffRampInput)
		// if err != nil {
		// 	return DeployCCIPOffRampSeqOutput{}, err
		// }

		// Note: Only need deploy + initialzie for now to run in a sequence.
		// We have a seperate changeset for set_ocr3_config

		// input.SetOCR3ConfigInput.OffRampPackageId = deployReport.Output.PackageId
		// input.SetOCR3ConfigInput.OwnerCapObjectId = deployReport.Output.Objects.OwnerCapObjectId
		// input.SetOCR3ConfigInput.OffRampStateId = deployReport.Output.Objects.CCIPOffRampStateObjectId
		// _, err = cld_ops.ExecuteOperation(env, operation.SetOCR3ConfigOp, deps, input.SetOCR3ConfigInput)
		// if err != nil {
		// 	return DeployCCIPOffRampSeqOutput{}, err
		// }

		return DeployCCIPOffRampSeqOutput{
			CCIPOffRampPackageId: deployReport.Output.PackageId,
			Objects: DeployCCIPOffRampSeqObjects{
				StateObjectId: deployReport.Output.Objects.CCIPOffRampStateObjectId,
				ObjectCapId:   deployReport.Output.Objects.OwnerCapObjectId,
			},
		}, nil
	},
)
