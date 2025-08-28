package sui

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui-internal/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui-internal/ops"
	burnminttokenpoolops "github.com/smartcontractkit/chainlink-sui-internal/ops/ccip_burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

var _ cldf.ChangeSetV2[DeploySuiBurnMintTpConfig] = DeploySuiBurnMintTp{}

// DeployAptosChain deploys Aptos chain packages and modules
type DeploySuiBurnMintTp struct{}

// Apply implements deployment.ChangeSetV2.
func (d DeploySuiBurnMintTp) Apply(e cldf.Environment, config DeploySuiBurnMintTpConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load Sui onchain state: %w", err)
	}

	ab := cldf.NewMemoryAddressBook()
	seqReports := make([]operations.Report[any, any], 0)

	suiChains := e.BlockChains.SuiChains()

	suiChain := suiChains[config.ChainSelector]
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
				b := uint64(400_000_000)
				return &bind.CallOpts{
					WaitForExecution: true,
					GasBudget:        &b,
				}
			},
		},
		CCIPOnChainState: state,
	}

	tokenPoolAddress := state.SuiChains[config.ChainSelector].TokenPoolAddress // BurnMintTokenPoolPackage
	ccipObjectRefId := state.SuiChains[config.ChainSelector].CCIPObjectRef

	linkTokenPkgId := "0xfd58da048fbf8d9c2749bc1fdaccf479a3e66065f5ab4a9e45c47a96921b882d"
	linkTokenObjectMetadataId := "0xabe4fe817da6fbd92f02ebcdf355822ce50f116b57dc5a043a119214c26019c9"
	linkTokenTreasuryCapId := "0x93a630e91e6d517cf3594cc8a224305485ff10e052bf11704dd0c2ce74556f0d"

	CCIPPackageId := state.SuiChains[config.ChainSelector].CCIPAddress
	MCMsPackageId := state.SuiChains[config.ChainSelector].MCMsAddress

	fmt.Println(
		"tokenPoolAddress:", tokenPoolAddress,
		"\nccipObjectRefId:", ccipObjectRefId,
		"\nlinkTokenPkgId:", linkTokenPkgId,
		"\nlinkTokenObjectMetadataId:", linkTokenObjectMetadataId,
		"\nlinkTokenTreasuryCapId:", linkTokenTreasuryCapId,
		"\nCCIPPackageId:", CCIPPackageId,
		"\nMCMsPackageId:", MCMsPackageId,
		config.RemoteChainSelector,
	)

	// // // Deploy BurnMint TP on SUI
	deployBurnMintTp, err := operations.ExecuteSequence(e.OperationsBundle, burnminttokenpoolops.DeployAndInitBurnMintTokenPoolSequence, deps.SuiChain,
		burnminttokenpoolops.DeployAndInitBurnMintTokenPoolInput{
			BurnMintTokenPoolDeployInput: burnminttokenpoolops.BurnMintTokenPoolDeployInput{
				CCIPPackageId:          CCIPPackageId,
				CCIPTokenPoolPackageId: tokenPoolAddress,
				MCMSAddress:            MCMsPackageId,
				MCMSOwnerAddress:       signerAddr,
			},

			CoinObjectTypeArg:      linkTokenPkgId + "::link_token::LINK_TOKEN",
			CCIPObjectRefObjectId:  ccipObjectRefId,
			CoinMetadataObjectId:   linkTokenObjectMetadataId,
			TreasuryCapObjectId:    linkTokenTreasuryCapId,
			TokenPoolAdministrator: signerAddr,

			// apply dest chain updates
			RemoteChainSelectorsToRemove: []uint64{},
			RemoteChainSelectorsToAdd:    []uint64{config.RemoteChainSelector},
			RemotePoolAddressesToAdd:     [][]string{{config.EVMTokenPool.String()}},
			RemoteTokenAddressesToAdd: []string{
				config.EVMToken.String(),
			},
			// set chain rate limiter configs
			RemoteChainSelectors: []uint64{config.RemoteChainSelector},
			OutboundIsEnableds:   []bool{false},
			OutboundCapacities:   []uint64{10000000000},
			OutboundRates:        []uint64{100},
			InboundIsEnableds:    []bool{false},
			InboundCapacities:    []uint64{10000000000},
			InboundRates:         []uint64{100},
		})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy BurnMintTP for Sui chain %d: %w", config.ChainSelector, err)
	}

	// save BnM TokenPool to addressbook
	typeAndVersionBurnMintTokenPool := cldf.NewTypeAndVersion(shared.SuiBnMTokenPoolType, deployment.Version1_5_1)
	err = deps.AB.Save(config.ChainSelector, deployBurnMintTp.Output.BurnMintTPPackageID, typeAndVersionBurnMintTokenPool)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save BurnMintTokenPool address %s for Sui chain %d: %w", deployBurnMintTp.Output.BurnMintTPPackageID, config.ChainSelector, err)
	}

	// save BnM TokenPool State to addressbook
	typeAndVersionBnMTpState := cldf.NewTypeAndVersion(shared.SuiBnMTokenPoolStateType, deployment.Version1_5_1)
	err = deps.AB.Save(config.ChainSelector, deployBurnMintTp.Output.Objects.StateObjectId, typeAndVersionBnMTpState)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save BurnMintTokenPoolState address %s for Sui chain %d: %w", deployBurnMintTp.Output.Objects.StateObjectId, config.ChainSelector, err)
	}
	seqReports = append(seqReports, deployBurnMintTp.ExecutionReports...)

	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d DeploySuiBurnMintTp) VerifyPreconditions(e cldf.Environment, config DeploySuiBurnMintTpConfig) error {
	return nil
}
