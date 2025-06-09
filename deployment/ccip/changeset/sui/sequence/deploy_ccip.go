package sequence

import (
	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/sui/operation"

	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

type DeployCCIPSeqObjects struct {
	CCIPObjectRefObjectId           string
	OwnerCapObjectId                string
	FeeQuoterCapObjectId            string
	FeeQuoterStateObjectId          string
	NonceManagerStateObjectId       string
	NonceManagerCapObjectId         string
	ReceiverRegistryStateObjectId   string
	RMNRemoteStateObjectId          string
	TokenAdminRegistryStateObjectId string
	SourceTransferCapObjectId       string
}

type DeployCCIPSeqOutput struct {
	CCIPPackageId string
	Objects       DeployCCIPSeqObjects
}

type DeployAndInitCCIPSeqInput struct {
	LinkTokenCoinMetadataObjectId string
	LocalChainSelector            uint64
	DestChainSelector             uint64
	operation.DeployCCIPInput
}

var DeployAndInitCCIPSequence = cld_ops.NewSequence(
	"sui-deploy-ccip-seq",
	semver.MustParse("0.1.0"),
	"Deploys and sets initial CCIP configuration",
	func(env cld_ops.Bundle, deps sui_ops.OpTxDeps, input DeployAndInitCCIPSeqInput) (DeployCCIPSeqOutput, error) {
		deployReport, err := cld_ops.ExecuteOperation(env, operation.DeployCCIPOp, deps, input.DeployCCIPInput)
		if err != nil {
			return DeployCCIPSeqOutput{}, err
		}
		return DeployCCIPSeqOutput{
			CCIPPackageId: deployReport.Output.PackageId,
			Objects: DeployCCIPSeqObjects{
				CCIPObjectRefObjectId: deployReport.Output.Objects.CCIPObjectRefObjectId,
				OwnerCapObjectId:      deployReport.Output.Objects.OwnerCapObjectId,
			},
		}, nil
	},
)
