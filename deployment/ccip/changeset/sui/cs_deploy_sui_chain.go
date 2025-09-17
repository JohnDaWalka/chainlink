package sui

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	offrampops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_offramp"
	onrampops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_onramp"
	routerops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_router"
	tokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_token_pool"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
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
		suiSigner := suiChain.Signer

		signerAddr, err := suiSigner.GetAddress()
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}

		deps := SuiDeps{
			AB: ab,
			SuiChain: sui_ops.OpTxDeps{
				Client: suiChain.Client,
				Signer: suiSigner,
				GetCallOpts: func() *bind.CallOpts {
					b := uint64(500_000_000)
					return &bind.CallOpts{
						WaitForExecution: true,
						GasBudget:        &b,
					}
				},
			},
			CCIPOnChainState: state,
		}

		// Deploy MCMS
		mcmsSeqReport, err := operations.ExecuteSequence(e.OperationsBundle, mcmsops.DeployMCMSSequence, deps.SuiChain, cld_ops.EmptyInput{})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy CCIP for Sui chain %d: %w", chainSel, err)
		}
		seqReports = append(seqReports, mcmsSeqReport.ExecutionReports...)

		// save MCMs address to the addressbook
		typeAndVersionMCMS := cldf.NewTypeAndVersion(shared.SuiMCMSType, deployment.Version1_0_0)
		err = deps.AB.Save(chainSel, mcmsSeqReport.Output.PackageId, typeAndVersionMCMS)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save MCMS address %s for Sui chain %d: %w", mcmsSeqReport.Output.PackageId, chainSel, err)
		}

		// Deploy Router
		// TODO: Maybe make this part of CCIP sequence
		routerReport, err := operations.ExecuteOperation(e.OperationsBundle, routerops.DeployCCIPRouterOp, deps.SuiChain, routerops.DeployCCIPRouterInput{
			McmsPackageId: mcmsSeqReport.Output.PackageId,
			McmsOwner:     signerAddr,
		})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy CCIP Router for Sui chain %d: %w", chainSel, err)
		}

		// save Router address to the addressbook
		typeAndVersionRouter := cldf.NewTypeAndVersion(shared.SuiCCIPRouterType, deployment.Version1_6_0)
		err = deps.AB.Save(chainSel, routerReport.Output.PackageId, typeAndVersionRouter)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save Router address %s for Sui chain %d: %w", routerReport.Output.PackageId, chainSel, err)
		}

		// Run DeployAndInitCCIpSequence
		ccipSeqInput := ccipops.DeployAndInitCCIPSeqInput{
			LinkTokenCoinMetadataObjectId: config.ContractParamsPerChain[chainSel].FeeQuoterParams.LinkTokenCoinMetadataObjectId,
			LocalChainSelector:            chainSel,
			DestChainSelector:             config.ContractParamsPerChain[chainSel].DestChainSelector,
			MaxFeeJuelsPerMsg:             config.ContractParamsPerChain[chainSel].FeeQuoterParams.MaxFeeJuelsPerMsg,
			TokenPriceStalenessThreshold:  config.ContractParamsPerChain[chainSel].FeeQuoterParams.TokenPriceStalenessThreshold,
			DeployCCIPInput: ccipops.DeployCCIPInput{
				McmsPackageId: mcmsSeqReport.Output.PackageId,
				McmsOwner:     signerAddr,
			},
			// Fee Quoter configuration
			AddMinFeeUsdCents:    []uint32{3000},
			AddMaxFeeUsdCents:    []uint32{30000},
			AddDeciBps:           []uint16{1000},
			AddDestGasOverhead:   []uint32{1000000},
			AddDestBytesOverhead: []uint32{1000},
			AddIsEnabled:         []bool{true},
			RemoveTokens:         []string{},

			// TODO: retrieve thesevalues from config
			// Fee Quoter destination chain configuration
			// values retried from here: https://github.com/smartcontractkit/chainlink-sui/pull/277/files#diff-5088e21cdbdb4efead9c1142c365a8717bcfe1a912cb0f2b54d0c0b7aad7e3c1R496
			IsEnabled:                         true,
			MaxNumberOfTokensPerMsg:           1,
			MaxDataBytes:                      30_000,
			MaxPerMsgGasLimit:                 3_000_000,
			DestGasOverhead:                   300_000,
			DestGasPerPayloadByteBase:         byte(16),
			DestGasPerPayloadByteHigh:         byte(40),
			DestGasPerPayloadByteThreshold:    uint16(3000),
			DestDataAvailabilityOverheadGas:   100,
			DestGasPerDataAvailabilityByte:    16,
			DestDataAvailabilityMultiplierBps: 1,
			ChainFamilySelector:               []byte{40, 18, 213, 44},
			EnforceOutOfOrder:                 false,
			DefaultTokenFeeUsdCents:           25,
			DefaultTokenDestGasOverhead:       90_000,
			DefaultTxGasLimit:                 200_000,
			GasMultiplierWeiPerEth:            1_000_000_000_000_000_000,
			GasPriceStalenessThreshold:        1_000_000,
			NetworkFeeUsdCents:                10,

			// apply_premium_multiplier_wei_per_eth_updates
			PremiumMultiplierWeiPerEth: []uint64{900_000_000_000_000_000},
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

		// save CCIP FeeQuoterCapObjectId address to the addressbook
		typeAndVersionCCIPFeeQuoterCapIdRef := cldf.NewTypeAndVersion(shared.SuiFeeQuoterCapType, deployment.Version1_6_0)
		err = deps.AB.Save(chainSel, ccipSeqReport.Output.Objects.FeeQuoterCapObjectId, typeAndVersionCCIPFeeQuoterCapIdRef)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save CCIP FeeQuoter CapId Id %s for Sui chain %d: %w", ccipSeqReport.Output.Objects.FeeQuoterCapObjectId, chainSel, err)
		}

		// No need to store rn
		// save CCIP TransferCapId address to the addressbook
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
				CCIPPackageId:      ccipSeqReport.Output.CCIPPackageId,
				MCMSPackageId:      mcmsSeqReport.Output.PackageId,
				MCMSOwnerPackageId: signerAddr,
			},
			OnRampInitializeInput: onrampops.OnRampInitializeInput{
				NonceManagerCapId:         ccipSeqReport.Output.Objects.NonceManagerCapObjectId,   // this is from NonceManager init Op
				SourceTransferCapId:       ccipSeqReport.Output.Objects.SourceTransferCapObjectId, // this is from CCIP package publish
				ChainSelector:             suiChain.Selector,
				FeeAggregator:             signerAddr,
				AllowListAdmin:            signerAddr,
				DestChainSelectors:        []uint64{config.ContractParamsPerChain[chainSel].DestChainSelector}, // TODOD add this in input instead of hardcoding
				DestChainEnabled:          []bool{true},
				DestChainAllowListEnabled: []bool{true},
			},
			ApplyDestChainConfigureOnRampInput: onrampops.ApplyDestChainConfigureOnRampInput{
				CCIPObjectRefId:           ccipSeqReport.Output.Objects.CCIPObjectRefObjectId,
				DestChainSelector:         []uint64{config.ContractParamsPerChain[chainSel].DestChainSelector},
				DestChainEnabled:          []bool{true},
				DestChainAllowListEnabled: []bool{false},
			},
			ApplyAllowListUpdatesInput: onrampops.ApplyAllowListUpdatesInput{
				CCIPObjectRefId:               ccipSeqReport.Output.Objects.CCIPObjectRefObjectId,
				DestChainSelector:             []uint64{config.ContractParamsPerChain[chainSel].DestChainSelector},
				DestChainAllowListEnabled:     []bool{false},
				DestChainAddAllowedSenders:    [][]string{{}},
				DestChainRemoveAllowedSenders: [][]string{{}},
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

		fmt.Println("ETH ONRAMP: ", deps.CCIPOnChainState.Chains[ccipSeqInput.DestChainSelector].OnRamp.Address())
		onRampBytes := [][]byte{deps.CCIPOnChainState.Chains[ccipSeqInput.DestChainSelector].OnRamp.Address().Bytes()} // ethereum chain onRamp bytes

		// Run DeployAndInitCCIPOffRampSequence
		ccipOffRampSeqInput := offrampops.DeployAndInitCCIPOffRampSeqInput{
			DeployCCIPOffRampInput: offrampops.DeployCCIPOffRampInput{
				CCIPPackageId: ccipSeqReport.Output.CCIPPackageId,
				MCMSPackageId: mcmsSeqReport.Output.PackageId,
			},
			InitializeOffRampInput: offrampops.InitializeOffRampInput{
				DestTransferCapId:                     ccipSeqReport.Output.Objects.DestTransferCapObjectId,
				FeeQuoterCapId:                        ccipSeqReport.Output.Objects.FeeQuoterCapObjectId,
				ChainSelector:                         suiChain.Selector,
				PremissionExecThresholdSeconds:        uint32(60 * 60 * 8),
				SourceChainSelectors:                  []uint64{config.ContractParamsPerChain[chainSel].DestChainSelector}, // this is ethereum
				SourceChainsIsEnabled:                 []bool{true},
				SourceChainsIsRMNVerificationDisabled: []bool{true},
				SourceChainsOnRamp:                    onRampBytes,
			},
		}
		ccipOffRampSeqReport, err := operations.ExecuteSequence(e.OperationsBundle, offrampops.DeployAndInitCCIPOffRampSequence, deps.SuiChain, ccipOffRampSeqInput)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy CCIP for Sui chain %d: %w", chainSel, err)
		}
		seqReports = append(seqReports, ccipOffRampSeqReport.ExecutionReports...)

		fmt.Println("SUI OFFRAMP: ", ccipOffRampSeqReport.Output.CCIPOffRampPackageId)

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

		// // NOT NEEDED FOR ARBITRARY MSG PASSING
		// // TODO abstract this into a different function
		// // Deploy CCIP TokenPool
		deployTp, err := operations.ExecuteOperation(e.OperationsBundle, tokenpoolops.DeployCCIPTokenPoolOp, deps.SuiChain,
			tokenpoolops.TokenPoolDeployInput{
				CCIPPackageId:    ccipSeqReport.Output.CCIPPackageId,
				MCMSAddress:      mcmsSeqReport.Output.PackageId,
				MCMSOwnerAddress: signerAddr,
			})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy TokenPool for Sui chain %d: %w", chainSel, err)
		}

		// save tokenPool address in addressbook
		typeAndVersionTokenPoolId := cldf.NewTypeAndVersion(shared.SuiTokenPoolType, deployment.Version1_6_0)
		err = deps.AB.Save(chainSel, deployTp.Output.PackageId, typeAndVersionTokenPoolId)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save offRamp StateObjectId %s for Sui chain %d: %w", deployTp.Output.PackageId, chainSel, err)
		}
	}

	fmt.Println("RAN CS_DEPLOY_SUI_CHAIN")
	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d DeploySuiChain) VerifyPreconditions(e cldf.Environment, config DeploySuiChainConfig) error {
	return nil
}
