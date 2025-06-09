package operation

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_onramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_onramp/onramp"
	"github.com/smartcontractkit/chainlink-sui/bindings/packages/onramp"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

type DeployCCIPOnRampObjects struct {
	// State Object
	OwnerCapObjectId        string
	CCIPOnrampStateObjectId string
}

type DeployCCIPOnRampInput struct {
	CCIPPackageId string
	MCMSPackageId string
}

var deployHandlerOnRamp = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input DeployCCIPOnRampInput) (output sui_ops.OpTxResult[DeployCCIPOnRampObjects], err error) {
	onRampPackage, tx, err := onramp.PublishOnramp(
		b.GetContext(),
		deps.GetTxOpts(),
		deps.Signer,
		deps.Client,
		input.CCIPPackageId,
		input.MCMSPackageId,
	)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{}, err
	}

	// TODO: We should move the object ID finding logic into the binding package
	obj1, err1 := bind.FindObjectIdFromPublishTx(*tx, "onramp", "OwnerCap")
	obj2, err2 := bind.FindObjectIdFromPublishTx(*tx, "onramp", "OnRampState")

	if err1 != nil || err2 != nil {
		return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{}, fmt.Errorf("failed to find object IDs in publish tx: %w", err)
	}

	return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{
		Digest:    tx.Digest.String(),
		PackageId: onRampPackage.Address().String(),
		Objects: DeployCCIPOnRampObjects{
			OwnerCapObjectId:        obj1,
			CCIPOnrampStateObjectId: obj2,
		},
	}, err
}

type OnRampInitializeInput struct {
	OnRampPackageId           string
	OnRampStateId             string
	OwnerCapObjectId          string
	NonceManagerCapId         string
	SourceTransferCapId       string
	ChainSelector             uint64
	FeeAggregator             string
	AllowListAdmin            string
	DestChainSelectors        []uint64
	DestChainEnabled          []bool
	DestChainAllowListEnabled []bool
}

var onRampInitializeHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input OnRampInitializeInput) (output sui_ops.OpTxResult[DeployCCIPOnRampObjects], err error) {
	onRampPackage, err := module_onramp.NewOnramp(input.OnRampPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{}, err
	}

	call := onRampPackage.Initialize(
		bind.Object{Id: input.OnRampStateId},
		bind.Object{Id: input.OwnerCapObjectId},
		bind.Object{Id: input.NonceManagerCapId},
		bind.Object{Id: input.SourceTransferCapId},
		input.ChainSelector,
		input.FeeAggregator,
		input.AllowListAdmin,
		input.DestChainSelectors,
		input.DestChainEnabled,
		input.DestChainAllowListEnabled,
	)

	tx, err := call.Execute(b.GetContext(), deps.GetTxOpts(), deps.Signer, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{}, fmt.Errorf("failed to execute onRamp initialization: %w", err)
	}

	return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{
		Digest:    tx.Digest.String(),
		PackageId: input.OnRampPackageId,
		Objects:   DeployCCIPOnRampObjects{},
	}, err
}

var DeployCCIPOnRampOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-on-ramp", "package", "deploy"),
	semver.MustParse("0.1.0"),
	"Deploys the CCIP onRamp package",
	deployHandlerOnRamp,
)

var OnRampInitializeOP = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-on-ramp", "package", "initialize"),
	semver.MustParse("0.1.0"),
	"Initialize the CCIP onRamp package",
	onRampInitializeHandler,
)
