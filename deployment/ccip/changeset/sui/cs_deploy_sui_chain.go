package sui

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/ops/ccip"
	lockreleasetokenpoolops "github.com/smartcontractkit/chainlink-sui/ops/ccip_lock_release_token_pool"
	offrampops "github.com/smartcontractkit/chainlink-sui/ops/ccip_offramp"
	onrampops "github.com/smartcontractkit/chainlink-sui/ops/ccip_onramp"
	routerops "github.com/smartcontractkit/chainlink-sui/ops/ccip_router"
	tokenpoolops "github.com/smartcontractkit/chainlink-sui/ops/ccip_token_pool"
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
					b := uint64(400_000_000)
					return bind.TxOpts{
						GasBudget: &b,
					}
				},
			},
			CCIPOnChainState: state,
		}

		// Deploy Router
		// TODO: Maybe make this part of CCIP sequence
		routerReport, err := operations.ExecuteOperation(e.OperationsBundle, routerops.DeployCCIPRouterOp, deps.SuiChain, cld_ops.EmptyInput{})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy CCIP Router for Sui chain %d: %w", chainSel, err)
		}

		// save Router address to the addressbook
		typeAndVersionRouter := cldf.NewTypeAndVersion(shared.SuiCCIPRouterType, deployment.Version1_6_0)
		err = deps.AB.Save(chainSel, routerReport.Output.PackageId, typeAndVersionRouter)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save Router address %s for Sui chain %d: %w", routerReport.Output.PackageId, chainSel, err)
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

		// Run DeployAndInitCCIpSequence
		ccipSeqInput := ccipops.DeployAndInitCCIPSeqInput{
			LinkTokenCoinMetadataObjectId: config.ContractParamsPerChain[chainSel].FeeQuoterParams.LinkTokenCoinMetadataObjectId,
			LocalChainSelector:            config.ContractParamsPerChain[chainSel].OnRampParams.ChainSelector,
			DestChainSelector:             uint64(909606746561742123),
			MaxFeeJuelsPerMsg:             config.ContractParamsPerChain[chainSel].FeeQuoterParams.MaxFeeJuelsPerMsg,
			TokenPriceStalenessThreshold:  config.ContractParamsPerChain[chainSel].FeeQuoterParams.TokenPriceStalenessThreshold,
			DeployCCIPInput: ccipops.DeployCCIPInput{
				McmsPackageId: mcmsSeqReport.Output.PackageId,
				McmsOwner:     "0x2",
			},
			// Fee Quoter destination chain configuration
			IsEnabled:                         true,
			MaxNumberOfTokensPerMsg:           2,
			MaxDataBytes:                      2000,
			MaxPerMsgGasLimit:                 5000000,
			DestGasOverhead:                   1000000,
			DestGasPerPayloadByteBase:         byte(2),
			DestGasPerPayloadByteHigh:         byte(5),
			DestGasPerPayloadByteThreshold:    uint16(10),
			DestDataAvailabilityOverheadGas:   300000,
			DestGasPerDataAvailabilityByte:    1,
			DestDataAvailabilityMultiplierBps: 1,
			ChainFamilySelector:               []byte{40, 18, 213, 44},
			EnforceOutOfOrder:                 false,
			DefaultTokenFeeUsdCents:           3,
			DefaultTokenDestGasOverhead:       100000,
			DefaultTxGasLimit:                 500000,
			GasMultiplierWeiPerEth:            100,
			GasPriceStalenessThreshold:        1000000000,
			NetworkFeeUsdCents:                10,

			// apply_premium_multiplier_wei_per_eth_updates
			PremiumMultiplierWeiPerEth: []uint64{10},
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
				CCIPPackageId:      ccipSeqReport.Output.CCIPPackageId,
				MCMSPackageId:      mcmsSeqReport.Output.PackageId,
				MCMSOwnerPackageId: "0x2",
			},
			OnRampInitializeInput: onrampops.OnRampInitializeInput{
				NonceManagerCapId:         ccipSeqReport.Output.Objects.NonceManagerCapObjectId,   // this is from NonceManager init Op
				SourceTransferCapId:       ccipSeqReport.Output.Objects.SourceTransferCapObjectId, // this is from CCIP package publish
				ChainSelector:             suiChain.Selector,
				FeeAggregator:             suiSignerAddr,
				AllowListAdmin:            suiSignerAddr,
				DestChainSelectors:        []uint64{909606746561742123}, // TODOD add this in input instead of hardcoding
				DestChainEnabled:          []bool{true},
				DestChainAllowListEnabled: []bool{true},
			},
			ApplyDestChainConfigureOnRampInput: onrampops.ApplyDestChainConfigureOnRampInput{
				DestChainSelector:         []uint64{909606746561742123},
				DestChainEnabled:          []bool{true},
				DestChainAllowListEnabled: []bool{false},
			},
			ApplyAllowListUpdatesInput: onrampops.ApplyAllowListUpdatesInput{
				DestChainSelector:             []uint64{909606746561742123},
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

		// Run DeployAndInitCCIPOffRampSequence
		ccipOffRampSeqInput := offrampops.DeployAndInitCCIPOffRampSeqInput{
			DeployCCIPOffRampInput: offrampops.DeployCCIPOffRampInput{
				CCIPPackageId: ccipSeqReport.Output.CCIPPackageId,
				MCMSPackageId: mcmsSeqReport.Output.PackageId,
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

		// NOT NEEDED FOR ARBITRARY MSG PASSING
		// TODO abstract this into a different function
		// Deploy CCIP TokenPool
		deployTp, err := operations.ExecuteOperation(e.OperationsBundle, tokenpoolops.DeployCCIPTokenPoolOp, deps.SuiChain,
			tokenpoolops.TokenPoolDeployInput{
				CCIPPackageId:     ccipSeqReport.Output.CCIPPackageId,
				CCIPRouterAddress: routerReport.Output.PackageId,
				MCMSAddress:       mcmsSeqReport.Output.PackageId,
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

		// linkTokenTreasuryCapId := state.SuiChains[suiChain.Selector].LinkTokenTreasuryCapId.String()
		linkTokenObjectMetadataId := state.SuiChains[suiChain.Selector].LinkTokenCoinMetadataId.String()
		// linkTokenPkgId := state.SuiChains[suiChain.Selector].LinkTokenAddress.String()

		// // Deploy LockRelease TP

		deployLockReleaseTp, err := operations.ExecuteOperation(e.OperationsBundle, lockreleasetokenpoolops.DeployCCIPLockReleaseTokenPoolOp, deps.SuiChain,
			lockreleasetokenpoolops.LockReleaseTokenPoolDeployInput{
				CCIPPackageId:          ccipSeqReport.Output.CCIPPackageId,
				CCIPRouterAddress:      routerReport.Output.PackageId,
				CCIPTokenPoolPackageId: deployTp.Output.PackageId,
				LockReleaseLocalToken:  linkTokenObjectMetadataId,
				MCMSAddress:            mcmsSeqReport.Output.PackageId,
			})
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}

		// deployLockReleaseTp, err := operations.ExecuteSequence(e.OperationsBundle, lockreleasetokenpoolops.DeployAndInitLockReleaseTokenPoolSequence, deps.SuiChain,
		// 	lockreleasetokenpoolops.DeployAndInitLockReleaseTokenPoolInput{
		// 		LockReleaseTokenPoolDeployInput: lockreleasetokenpoolops.LockReleaseTokenPoolDeployInput{
		// 			CCIPPackageId:          ccipSeqReport.Output.CCIPPackageId,
		// 			CCIPRouterAddress:      routerReport.Output.PackageId,
		// 			CCIPTokenPoolPackageId: deployTp.Output.PackageId,
		// 			LockReleaseLocalToken:  linkTokenObjectMetadataId,
		// 			MCMSAddress:            mcmsSeqReport.Output.PackageId,
		// 		},

		// deployLockReleaseTp, err := operations.ExecuteSequence(e.OperationsBundle, lockreleasetokenpoolops.DeployAndInitLockReleaseTokenPoolSequence, deps.SuiChain,
		// 	lockreleasetokenpoolops.DeployAndInitLockReleaseTokenPoolInput{
		// 		LockReleaseTokenPoolDeployInput: lockreleasetokenpoolops.LockReleaseTokenPoolDeployInput{
		// 			CCIPPackageId:          ccipSeqReport.Output.CCIPPackageId,
		// 			CCIPRouterAddress:      routerReport.Output.PackageId,
		// 			CCIPTokenPoolPackageId: deployTp.Output.PackageId,
		// 			LockReleaseLocalToken:  linkTokenObjectMetadataId,
		// 			MCMSAddress:            mcmsSeqReport.Output.PackageId,
		// 		},

		// 		CoinObjectTypeArg:     linkTokenPkgId + "::link_token::LINK_TOKEN",
		// 		CCIPObjectRefObjectId: ccipSeqReport.Output.Objects.CCIPObjectRefObjectId,
		// 		CoinMetadataObjectId:  linkTokenObjectMetadataId,
		// 		TreasuryCapObjectId:   linkTokenTreasuryCapId,
		// 		TokenPoolPackageId:    deployTp.Output.PackageId,
		// 		Rebalancer:            "",

		// 		// apply dest chain updates
		// 		RemoteChainSelectorsToRemove: []uint64{},
		// 		RemoteChainSelectorsToAdd:    []uint64{909606746561742123},
		// 		RemotePoolAddressesToAdd: [][][]byte{
		// 			{
		// 				[]byte{
		// 					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// 					0x00, 0x00, 0x00, 0x00, 0xaf, 0x46, 0xbf, 0x6d,
		// 					0x19, 0x92, 0x1e, 0x30, 0xbc, 0x5c, 0xc0, 0x04,
		// 					0x3d, 0xc6, 0xde, 0x91, 0xed, 0xf0, 0x0c, 0x98,
		// 				},
		// 			},
		// 		},
		// 		RemoteTokenAddressesToAdd: [][]byte{
		// 			[]byte{
		// 				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// 				0x00, 0x00, 0x00, 0x00, 0x77, 0x98, 0x77, 0xa7,
		// 				0xb0, 0xd9, 0xe8, 0x60, 0x31, 0x69, 0xdd, 0xbd,
		// 				0x78, 0x36, 0xe4, 0x78, 0xb4, 0x62, 0x47, 0x89,
		// 			},
		// 		},

		// 		// set chain rate limiter configs
		// 		RemoteChainSelectors: []uint64{909606746561742123},
		// 		OutboundIsEnableds:   []bool{true},
		// 		OutboundCapacities:   []uint64{1000},
		// 		OutboundRates:        []uint64{1000},
		// 		InboundIsEnableds:    []bool{true},
		// 		InboundCapacities:    []uint64{100},
		// 		InboundRates:         []uint64{1000},
		// 	})
		// if err != nil {
		// 	return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy TokenPool for Sui chain %d: %w", suiChain.Selector, err)
		// }

		// save LockRelease PackageId to the addressbook
		typeAndVersionLockReleasePackageId := cldf.NewTypeAndVersion(shared.SuiLockReleaseTPType, deployment.Version1_6_0)
		err = deps.AB.Save(chainSel, deployLockReleaseTp.Output.PackageId, typeAndVersionLockReleasePackageId)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save LockRelease PackageId %s for Sui chain %d: %w", deployLockReleaseTp.Output.PackageId, chainSel, err)
		}

		// // save LockRelease stateObjectId to the addressbook
		// typeAndVersionLockReleaseStateId := cldf.NewTypeAndVersion(shared.SuiLockReleaseTPStateType, deployment.Version1_6_0)
		// err = deps.AB.Save(chainSel, deployLockReleaseTp.Output.Objects.StateObjectId, typeAndVersionLockReleaseStateId)
		// if err != nil {
		// 	return cldf.ChangesetOutput{}, fmt.Errorf("failed to save LockRelase StateObjectId %s for Sui chain %d: %w", deployLockReleaseTp.Output.Objects.StateObjectId, chainSel, err)
		// }

		// seqReports = append(seqReports, deployLockReleaseTp.ExecutionReports...)

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
