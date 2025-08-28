package sui

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui-internal/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui-internal/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui-internal/ops/ccip"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

var _ cldf.ChangeSetV2[ApplyPremiumMultiplierWeiPerEthConfig] = ApplyPremiumMultiplierWeiPerEth{}

// DeployAptosChain deploys Aptos chain packages and modules
type ApplyPremiumMultiplierWeiPerEth struct{}

// Apply implements deployment.ChangeSetV2.
func (d ApplyPremiumMultiplierWeiPerEth) Apply(e cldf.Environment, config ApplyPremiumMultiplierWeiPerEthConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load Sui onchain state: %w", err)
	}

	ab := cldf.NewMemoryAddressBook()
	seqReports := make([]operations.Report[any, any], 0)

	suiChains := e.BlockChains.SuiChains()

	suiChain := suiChains[config.ChainSelector]
	suiSigner := suiChain.Signer

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

	// Run applyFeeTokenUpdate Operation
	applyFeeTokenUpdate, err := operations.ExecuteOperation(e.OperationsBundle, ccipops.FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesOp, deps.SuiChain,
		ccipops.FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesInput{
			CCIPPackageId:              config.CCIPPackageId,
			StateObjectId:              config.StateObjectId,
			OwnerCapObjectId:           config.OwnerCapObjectId,
			Tokens:                     config.Tokens,
			PremiumMultiplierWeiPerEth: config.PremiumMultiplierWeiPerEth,
		})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to applyFeeTokenUpdate for Sui chain %d: %w", config.ChainSelector, err)
	}

	seqReports = append(seqReports, applyFeeTokenUpdate.ToGenericReport())

	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d ApplyPremiumMultiplierWeiPerEth) VerifyPreconditions(e cldf.Environment, config ApplyPremiumMultiplierWeiPerEthConfig) error {
	return nil
}
