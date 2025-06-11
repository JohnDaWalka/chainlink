package sui

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/ops/ccip"
	offrampops "github.com/smartcontractkit/chainlink-sui/ops/ccip_offramp"
	onrampops "github.com/smartcontractkit/chainlink-sui/ops/ccip_onramp"
	mcmsops "github.com/smartcontractkit/chainlink-sui/ops/mcms"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

var _ cldf.ChangeSetV2[DeploySuiChainConfig] = DeploySuiChain{}

// DeployAptosChain deploys Aptos chain packages and modules
type DeploySuiChain struct{}

// Apply implements deployment.ChangeSetV2.
func (d DeploySuiChain) Apply(e cldf.Environment, config DeploySuiChainConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load Sui onchain state: %w", err)
	}

	ab := cldf.NewMemoryAddressBook()
	seqReports := make([]operations.Report[any, any], 0)

	for chainSel := range config.ContractParamsPerChain {
		suiChains := e.BlockChains.SuiChains()

		suiChain := suiChains[chainSel]
		suiSigner := rel.NewPrivateKeySigner(suiChain.DeployerKey)

		suiSignerAddr, err := suiSigner.GetAddress()
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}

		deps := SuiDeps{
			AB: ab,
			SuiChain: sui_ops.OpTxDeps{
				Client: *suiChain.Client,
				Signer: suiSigner,
				GetTxOpts: func() bind.TxOpts {
					b := uint64(300_000_000)
					return bind.TxOpts{
						GasBudget: &b,
					}
				},
			},
			CCIPOnChainState: state,
		}

		// Run DeployAndInitCCIpSequence
		ccipSeqInput := ccipops.DeployAndInitCCIPSeqInput{
			LinkTokenCoinMetadataObjectId: config.ContractParamsPerChain[chainSel].FeeQuoterParams.LinkTokenCoinMetadataObjectId,
			LocalChainSelector:            config.ContractParamsPerChain[chainSel].OnRampParams.ChainSelector,
			DestChainSelector:             suiChain.Selector,
			MaxFeeJuelsPerMsg:             config.ContractParamsPerChain[chainSel].FeeQuoterParams.MaxFeeJuelsPerMsg,
			TokenPriceStalenessThreshold:  config.ContractParamsPerChain[chainSel].FeeQuoterParams.TokenPriceStalenessThreshold,
			DeployCCIPInput: ccipops.DeployCCIPInput{
				McmsPackageId: "0x2",
			},
		}

		ccipSeqReport, err := operations.ExecuteSequence(e.OperationsBundle, ccipops.DeployAndInitCCIPSequence, deps.SuiChain, ccipSeqInput)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy CCIP for Sui chain %d: %w", chainSel, err)
		}
		seqReports = append(seqReports, ccipSeqReport.ExecutionReports...)

		// save CCIP address to the addressbook
		typeAndVersionCCIP := cldf.NewTypeAndVersion(shared.SuiCCIPType, deployment.Version1_6_0)
		err = deps.AB.Save(chainSel, ccipSeqReport.Output.CCIPPackageId, typeAndVersionCCIP)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save CCIP address %s for Sui chain %d: %w", ccipSeqReport.Output.CCIPPackageId, chainSel, err)
		}

		// save CCIP ObjectRef address to the addressbook
		typeAndVersionCCIPObjectRef := cldf.NewTypeAndVersion(shared.SuiCCIPObjectRefType, deployment.Version1_6_0)
		err = deps.AB.Save(chainSel, ccipSeqReport.Output.Objects.CCIPObjectRefObjectId, typeAndVersionCCIPObjectRef)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save CCIP objectRef Id %s for Sui chain %d: %w", ccipSeqReport.Output.Objects.CCIPObjectRefObjectId, chainSel, err)
		}

		// No need to store rn
		// // save CCIP TransferCapId address to the addressbook
		// typeAndVersionTransferCapId := cldf.NewTypeAndVersion(shared.SuiCCIPTransferCapIdType, deployment.Version1_6_0)
		// err = deps.AB.Save(chainSel, ccipSeqReport.Output.Objects.SourceTransferCapObjectId, typeAndVersionTransferCapId)
		// if err != nil {
		// 	return cldf.ChangesetOutput{}, fmt.Errorf("failed to save CCIP TransferCapId Id %s for Sui chain %d: %w", ccipSeqReport.Output.Objects.SourceTransferCapObjectId, chainSel, err)
		// }

		// // save CCIP NonceManagerCapObjectId address to the addressbook
		// typeAndVersionNonceManagerCapObjectId := cldf.NewTypeAndVersion(shared.SuiCCIPObjectRefType, deployment.Version1_6_0)
		// err = deps.AB.Save(chainSel, ccipSeqReport.Output.Objects.NonceManagerCapObjectId, typeAndVersionNonceManagerCapObjectId)
		// if err != nil {
		// 	return cldf.ChangesetOutput{}, fmt.Errorf("failed to save CCIP objectRef Id %s for Sui chain %d: %w", ccipSeqReport.Output.Objects.CCIPObjectRefObjectId, chainSel, err)
		// }

		// Run DeployAndInitCCIPOnRampSequence
		ccipOnRampSeqInput := onrampops.DeployAndInitCCIPOnRampSeqInput{
			DeployCCIPOnRampInput: onrampops.DeployCCIPOnRampInput{
				CCIPPackageId: ccipSeqReport.Output.CCIPPackageId,
			},

			OnRampInitializeInput: onrampops.OnRampInitializeInput{
				NonceManagerCapId:         ccipSeqReport.Output.Objects.NonceManagerCapObjectId,   // this is from NonceManager init Op
				SourceTransferCapId:       ccipSeqReport.Output.Objects.SourceTransferCapObjectId, // this is from CCIP package publish
				ChainSelector:             suiChain.Selector,
				FeeAggregator:             suiSignerAddr,
				AllowListAdmin:            suiSignerAddr,
				DestChainSelectors:        []uint64{909606746561742123},
				DestChainEnabled:          []bool{true},
				DestChainAllowListEnabled: []bool{true},
			},
		}

		ccipOnRampSeqReport, err := operations.ExecuteSequence(e.OperationsBundle, onrampops.DeployAndInitCCIPOnRampSequence, deps.SuiChain, ccipOnRampSeqInput)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy CCIP for Sui chain %d: %w", chainSel, err)
		}
		seqReports = append(seqReports, ccipOnRampSeqReport.ExecutionReports...)

		// save onRamp address to the addressbook
		typeAndVersionOnRamp := cldf.NewTypeAndVersion(shared.SuiOnRampType, deployment.Version1_6_0)
		err = deps.AB.Save(chainSel, ccipOnRampSeqReport.Output.CCIPOnRampPackageId, typeAndVersionOnRamp)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save onRamp address %s for Sui chain %d: %w", ccipOnRampSeqReport.Output.CCIPOnRampPackageId, chainSel, err)
		}

		// save onRampStateId address to the addressbook
		typeAndVersionOnRampStateId := cldf.NewTypeAndVersion(shared.SuiOnRampStateObjectIdType, deployment.Version1_6_0)
		err = deps.AB.Save(chainSel, ccipOnRampSeqReport.Output.Objects.StateObjectId, typeAndVersionOnRampStateId)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save onRamp state object Id  %s for Sui chain %d: %w", ccipOnRampSeqReport.Output.Objects.StateObjectId, chainSel, err)
		}

		// Run DeployMCMSSequence
		mcmsReport, err := operations.ExecuteSequence(e.OperationsBundle, mcmsops.DeployMCMSSequence, deps.SuiChain, cld_ops.EmptyInput{})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy CCIP for Sui chain %d: %w", chainSel, err)
		}
		seqReports = append(seqReports, mcmsReport.ExecutionReports...)

		// save onRamp address to the addressbook
		typeAndVersionMCMs := cldf.NewTypeAndVersion(shared.SuiMCMSType, deployment.Version1_6_0)
		err = deps.AB.Save(chainSel, mcmsReport.Output.PackageId, typeAndVersionMCMs)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save MCMs address %s for Sui chain %d: %w", mcmsReport.Output.PackageId, chainSel, err)
		}

		// Run DeployAndInitCCIPOffRampSequence
		ccipOffRampSeqInput := offrampops.DeployAndInitCCIPOffRampSeqInput{
			DeployCCIPOffRampInput: offrampops.DeployCCIPOffRampInput{
				CCIPPackageId: ccipSeqReport.Output.CCIPPackageId,
				MCMSPackageId: mcmsReport.Output.PackageId,
			},
		}
		ccipOffRampSeqReport, err := operations.ExecuteSequence(e.OperationsBundle, offrampops.DeployAndInitCCIPOffRampSequence, deps.SuiChain, ccipOffRampSeqInput)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy CCIP for Sui chain %d: %w", chainSel, err)
		}
		seqReports = append(seqReports, ccipOffRampSeqReport.ExecutionReports...)

		// save offRamp address to the addressbook
		typeAndVersionOffRamp := cldf.NewTypeAndVersion(shared.SuiOffRampType, deployment.Version1_6_0)
		err = deps.AB.Save(chainSel, ccipOffRampSeqReport.Output.CCIPOffRampPackageId, typeAndVersionOffRamp)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save offRamp address %s for Sui chain %d: %w", ccipOffRampSeqReport.Output.CCIPOffRampPackageId, chainSel, err)
		}

		// save offRamp ownerCapId to the addressbook
		typeAndVersionOffRampOwnerCapId := cldf.NewTypeAndVersion(shared.SuiOffRampOwnerCapObjectIdType, deployment.Version1_6_0)
		err = deps.AB.Save(chainSel, ccipOffRampSeqReport.Output.Objects.OwnerCapId, typeAndVersionOffRampOwnerCapId)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save offRamp ObjectCapId address %s for Sui chain %d: %w", ccipOffRampSeqReport.Output.CCIPOffRampPackageId, chainSel, err)
		}

		// save offRamp stateObjectId to the addressbook
		typeAndVersionOffRampObjectStateId := cldf.NewTypeAndVersion(shared.SuiOffRampStateObjectIdType, deployment.Version1_6_0)
		err = deps.AB.Save(chainSel, ccipOffRampSeqReport.Output.Objects.StateObjectId, typeAndVersionOffRampObjectStateId)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save offRamp StateObjectId %s for Sui chain %d: %w", ccipOffRampSeqReport.Output.Objects.StateObjectId, chainSel, err)
		}

	}
	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d DeploySuiChain) VerifyPreconditions(e cldf.Environment, config DeploySuiChainConfig) error {
	return nil
}
