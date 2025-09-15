package v1_6_3_dev

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	fq1_6_0 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	opsutil "github.com/smartcontractkit/chainlink/deployment/common/opsutils"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	ccipopsv1_6 "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_6"
	ccipopsv1_6_3_dev "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_6_3_dev"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

type ChainContractParams struct {
	FeeQuoterParams ccipopsv1_6_3_dev.FeeQuoterParamsSui
}

var (
	FeeQuoterWithSuiSupportChangeset = cldf.CreateChangeSet(deployFeeQuoterWithSuiSupportLogic, deployFeeQuoterWithSuiSupportPreCondition)
)

type FeeQuoterWithSuiSupportConfig struct {
	ContractParamsPerChain map[uint64]ChainContractParams
}

func deployFeeQuoterWithSuiSupportLogic(e cldf.Environment, config FeeQuoterWithSuiSupportConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		e.Logger.Errorw("Failed to load existing onchain state", "err", err)
		return cldf.ChangesetOutput{}, err
	}

	ab := cldf.NewMemoryAddressBook()

	for chainSelector, contractParams := range config.ContractParamsPerChain {
		targetChain := e.BlockChains.EVMChains()[chainSelector]

		targetChainState, chainExists := state.Chains[targetChain.Selector]
		if !chainExists {
			return cldf.ChangesetOutput{}, fmt.Errorf("chain %s not found in existing state, deploy the prerequisites first", targetChain.String())
		}

		// get the existing contract addresses
		linkAddr, err := targetChainState.LinkTokenAddress()
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}

		weth9Addr := opsutil.GetAddressSafely(targetChainState.Weth9)
		timelockAddr := opsutil.GetAddressSafely(targetChainState.Timelock)
		offRampAddr := opsutil.GetAddressSafely(targetChainState.OffRamp)

		report, err := operations.ExecuteOperation(e.OperationsBundle, ccipopsv1_6_3_dev.DeploySuiSupportedFeeQuoterOp, targetChain, opsutil.EVMDeployInput[ccipopsv1_6_3_dev.DeployFeeQInput]{
			ChainSelector: targetChain.Selector,
			DeployInput: ccipopsv1_6_3_dev.DeployFeeQInput{
				Chain:    targetChain.Selector,
				Params:   contractParams.FeeQuoterParams,
				LinkAddr: linkAddr,
				WethAddr: weth9Addr,
				// Allow timelock and deployer key to set prices.
				// Deployer key should be removed sometime after initial deployment
				PriceUpdaters: []common.Address{timelockAddr, targetChain.DeployerKey.From},
			},
		})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy fee quoter for %s: %w", targetChain, err)
		}

		feeQuoterAddress := report.Output.Address

		err = ab.Save(targetChain.Selector, feeQuoterAddress.String(), cldf.NewTypeAndVersion(shared.FeeQuoter, deployment.Version1_6_3Dev))
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save address %s for chain %d: %w", feeQuoterAddress.String(), targetChain.Selector, err)
		}

		_, err = operations.ExecuteOperation(e.OperationsBundle, ccipopsv1_6.FeeQApplyAuthorizedCallerOp, targetChain, opsutil.EVMCallInput[fq1_6_0.AuthorizedCallersAuthorizedCallerArgs]{
			ChainSelector: chainSelector,
			NoSend:        false,
			Address:       feeQuoterAddress,
			CallInput: fq1_6_0.AuthorizedCallersAuthorizedCallerArgs{
				AddedCallers: []common.Address{offRampAddr},
			},
		})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to set off ramp as authorized caller of FeeQuoter on chain %s: %w", targetChain, err)
		}
	}

	return cldf.ChangesetOutput{
		AddressBook: ab,
	}, nil
}

func deployFeeQuoterWithSuiSupportPreCondition(e cldf.Environment, _ FeeQuoterWithSuiSupportConfig) error {
	return nil
}
