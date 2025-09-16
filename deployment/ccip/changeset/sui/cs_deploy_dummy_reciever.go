package sui

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

var _ cldf.ChangeSetV2[DeployDummyRecieverConfig] = DeployDummyReciever{}

// DeployAptosChain deploys Aptos chain packages and modules
type DeployDummyReciever struct{}

// Apply implements deployment.ChangeSetV2.
func (d DeployDummyReciever) Apply(e cldf.Environment, config DeployDummyRecieverConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load Sui onchain state: %w", err)
	}

	ab := cldf.NewMemoryAddressBook()
	seqReports := make([]operations.Report[any, any], 0)

	suiChain := e.BlockChains.SuiChains()[config.ChainSelector]

	deps := SuiDeps{
		AB: ab,
		SuiChain: sui_ops.OpTxDeps{
			Client: suiChain.Client,
			Signer: suiChain.Signer,
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

	// Run DummyReciever Operation
	_, err = operations.ExecuteOperation(e.OperationsBundle, ccipops.DeployCCIPDummyReceiverOp, deps.SuiChain, ccipops.DeployDummyReceiverInput{
		CCIPPackageId: config.CCIPPackageId,
		McmsPackageId: config.McmsPackageId,
		McmsOwner:     config.McmsOwner,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy LinkToken for Sui chain %d: %w", config.ChainSelector, err)
	}

	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d DeployDummyReciever) VerifyPreconditions(e cldf.Environment, config DeployDummyRecieverConfig) error {
	return nil
}
