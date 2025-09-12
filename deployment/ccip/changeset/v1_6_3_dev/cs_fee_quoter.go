package v1_6_3_dev

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/fee_quoter"

	"github.com/smartcontractkit/chainlink/deployment"
	v1_6FeeQuoter "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	opsutil "github.com/smartcontractkit/chainlink/deployment/common/opsutils"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	ccipopsv1_6_3_dev "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_6_3_dev"
	ccipseqs "github.com/smartcontractkit/chainlink/deployment/ccip/sequence/evm/v1_6_3_dev"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commonopsutil "github.com/smartcontractkit/chainlink/deployment/common/opsutils"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

type ChainContractParams struct {
	FeeQuoterParams ccipopsv1_6_3_dev.FeeQuoterParamsSui
}

var (
	_ cldf.ChangeSet[UpdateFeeQuoterDestsConfig]  = UpdateFeeQuoterDestsChangeset
	_ cldf.ChangeSet[UpdateFeeQuoterPricesConfig] = UpdateFeeQuoterPricesChangeset

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

		weth9Addr := commonopsutil.GetAddressSafely(targetChainState.Weth9)
		timelockAddr := commonopsutil.GetAddressSafely(targetChainState.Timelock)
		offRampAddr := commonopsutil.GetAddressSafely(targetChainState.OffRamp)

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

		err = ab.Save(targetChain.Selector, feeQuoterAddress.String(), cldf.NewTypeAndVersion(shared.SuiSupportedFeeQuoter, deployment.Version1_6_3Dev))
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save address %s for chain %d: %w", feeQuoterAddress.String(), targetChain.Selector, err)
		}

		_, err = operations.ExecuteOperation(e.OperationsBundle, ccipopsv1_6_3_dev.SuiFeeQApplyAuthorizedCallerOp, targetChain, opsutil.EVMCallInput[fee_quoter.AuthorizedCallersAuthorizedCallerArgs]{
			ChainSelector: chainSelector,
			NoSend:        false,
			Address:       feeQuoterAddress,
			CallInput: fee_quoter.AuthorizedCallersAuthorizedCallerArgs{
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

type UpdateFeeQuoterDestsConfig struct {
	// UpdatesByChain is a mapping from source -> dest -> config update.
	UpdatesByChain map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig

	// Disallow mixing MCMS/non-MCMS per chain for simplicity.
	// (can still be achieved by calling this function multiple times)
	MCMS *proposalutils.TimelockConfig
}

func (cfg UpdateFeeQuoterDestsConfig) Validate(e cldf.Environment) error {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return err
	}
	supportedChains := state.SupportedChains()
	for chainSel, updates := range cfg.UpdatesByChain {
		chainState, ok := state.Chains[chainSel]
		if !ok {
			return fmt.Errorf("chain %d not found in onchain state", chainSel)
		}
		if chainState.TestRouter == nil {
			return fmt.Errorf("missing test router for chain %d", chainSel)
		}
		if chainState.Router == nil {
			return fmt.Errorf("missing router for chain %d", chainSel)
		}
		if chainState.OnRamp == nil {
			return fmt.Errorf("missing onramp onramp for chain %d", chainSel)
		}
		if err := commoncs.ValidateOwnership(e.GetContext(), cfg.MCMS != nil, e.BlockChains.EVMChains()[chainSel].DeployerKey.From, chainState.Timelock.Address(), chainState.FeeQuoter); err != nil {
			return err
		}

		for destination := range updates {
			// Destination cannot be an unknown destination.
			if _, ok := supportedChains[destination]; !ok {
				return fmt.Errorf("destination chain %d is not a supported %s", destination, chainState.OnRamp.Address())
			}
			sc, err := chainState.OnRamp.GetStaticConfig(&bind.CallOpts{Context: e.GetContext()})
			if err != nil {
				return fmt.Errorf("failed to get onramp static config %s: %w", chainState.OnRamp.Address(), err)
			}
			if destination == sc.ChainSelector {
				return errors.New("source and destination chain cannot be the same")
			}

			homeChainSelector, err := state.HomeChainSelector()
			if err != nil {
				return fmt.Errorf("failed to get home chain selector: %w", err)
			}
			execConfigs, err := v1_6FeeQuoter.GetAllActiveExecConfigs(e, homeChainSelector, destination)
			if err != nil {
				return fmt.Errorf("failed to get exec configs for destination chain %d: %w", destination, err)
			}

			for _, execOffchainCfg := range execConfigs {
				if execOffchainCfg.MultipleReportsEnabled {
					if !updates[destination].EnforceOutOfOrder {
						return errors.New("EnforceOutOfOrder must be true when MultipleReportsEnabled is true")
					}
				}
			}
		}
	}
	return nil
}

func (cfg UpdateFeeQuoterDestsConfig) ToSequenceInput(state stateview.CCIPOnChainState) ccipseqs.FeeQuoterApplyDestChainConfigUpdatesSequenceInput {
	updates := make(map[uint64]opsutil.EVMCallInput[[]fee_quoter.FeeQuoterDestChainConfigArgs], len(cfg.UpdatesByChain))
	for chainSel, destChainUpdates := range cfg.UpdatesByChain {
		args := make([]fee_quoter.FeeQuoterDestChainConfigArgs, len(destChainUpdates))
		i := 0
		for destChainSel, destChainUpdate := range destChainUpdates {
			args[i] = fee_quoter.FeeQuoterDestChainConfigArgs{
				DestChainSelector: destChainSel,
				DestChainConfig:   destChainUpdate,
			}
			i++
		}
		updates[chainSel] = opsutil.EVMCallInput[[]fee_quoter.FeeQuoterDestChainConfigArgs]{
			Address:       state.Chains[chainSel].FeeQuoter.Address(),
			ChainSelector: chainSel,
			CallInput:     args,
			NoSend:        cfg.MCMS != nil, // If MCMS exists, we do not want to send the transaction.
		}
	}

	return ccipseqs.FeeQuoterApplyDestChainConfigUpdatesSequenceInput{
		UpdatesByChain: updates,
	}
}

func UpdateFeeQuoterDestsChangeset(e cldf.Environment, cfg UpdateFeeQuoterDestsConfig) (cldf.ChangesetOutput, error) {
	output := cldf.ChangesetOutput{}

	if err := cfg.Validate(e); err != nil {
		return output, err
	}
	s, err := stateview.LoadOnchainState(e)
	if err != nil {
		return output, err
	}

	report, err := operations.ExecuteSequence(
		e.OperationsBundle,
		ccipseqs.SuiSupportedFeeQuoterApplyDestChainConfigUpdatesSequence,
		e.BlockChains.EVMChains(),
		cfg.ToSequenceInput(s),
	)

	return opsutil.AddEVMCallSequenceToCSOutput(e, output, report, err, s.EVMMCMSStateByChain(), cfg.MCMS, "Call ApplyDestChainConfigUpdates on FeeQuoters")
}

type UpdateFeeQuoterPricesConfig struct {
	PricesByChain map[uint64]FeeQuoterPriceUpdatePerSource // source -> PriceDetails
	MCMS          *proposalutils.TimelockConfig
}

type FeeQuoterPriceUpdatePerSource struct {
	TokenPrices map[common.Address]*big.Int // token address -> price
	GasPrices   map[uint64]*big.Int         // dest chain -> gas price
}

func (cfg UpdateFeeQuoterPricesConfig) Validate(e cldf.Environment) error {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return err
	}
	for chainSel, initialPrice := range cfg.PricesByChain {
		if err := cldf.IsValidChainSelector(chainSel); err != nil {
			return fmt.Errorf("invalid chain selector: %w", err)
		}
		chainState, ok := state.Chains[chainSel]
		if !ok {
			return fmt.Errorf("chain %d not found in onchain state", chainSel)
		}
		fq := chainState.FeeQuoter
		if fq == nil {
			return fmt.Errorf("missing fee quoter for chain %d", chainSel)
		}
		if err := commoncs.ValidateOwnership(e.GetContext(), cfg.MCMS != nil, e.BlockChains.EVMChains()[chainSel].DeployerKey.From, chainState.Timelock.Address(), chainState.FeeQuoter); err != nil {
			return err
		}
		// check that whether price updaters are set
		authCallers, err := fq.GetAllAuthorizedCallers(&bind.CallOpts{Context: e.GetContext()})
		if err != nil {
			return fmt.Errorf("failed to get authorized callers for chain %d: %w", chainSel, err)
		}
		if len(authCallers) == 0 {
			return fmt.Errorf("no authorized callers for chain %d", chainSel)
		}
		expectedAuthCaller := e.BlockChains.EVMChains()[chainSel].DeployerKey.From
		if cfg.MCMS != nil {
			expectedAuthCaller = chainState.Timelock.Address()
		}
		foundCaller := false
		for _, authCaller := range authCallers {
			if authCaller.Cmp(expectedAuthCaller) == 0 {
				foundCaller = true
			}
		}
		if !foundCaller {
			return fmt.Errorf("expected authorized caller %s not found for chain %d", expectedAuthCaller.String(), chainSel)
		}
		for token, price := range initialPrice.TokenPrices {
			if price == nil {
				return fmt.Errorf("token price for chain %d is nil", chainSel)
			}
			if token == (common.Address{}) {
				return fmt.Errorf("token address for chain %d is empty", chainSel)
			}
			contains, err := cldf.AddressBookContains(e.ExistingAddresses, chainSel, token.String())
			if err != nil {
				return fmt.Errorf("error checking address book for token %s: %w", token.String(), err)
			}
			if !contains {
				return fmt.Errorf("token %s not found in address book for chain %d", token.String(), chainSel)
			}
		}
		for dest, price := range initialPrice.GasPrices {
			if chainSel == dest {
				return errors.New("source and dest chain cannot be the same")
			}
			if err := cldf.IsValidChainSelector(dest); err != nil {
				return fmt.Errorf("invalid dest chain selector: %w", err)
			}
			if price == nil {
				return fmt.Errorf("gas price for chain %d is nil", chainSel)
			}
			if _, ok := state.SupportedChains()[dest]; !ok {
				return fmt.Errorf("dest chain %d not found in onchain state for chain %d", dest, chainSel)
			}
		}
	}

	return nil
}

func (cfg UpdateFeeQuoterPricesConfig) ToSequenceInput(state stateview.CCIPOnChainState) ccipseqs.FeeQuoterUpdatePricesSequenceInput {
	updates := make(map[uint64]opsutil.EVMCallInput[fee_quoter.InternalPriceUpdates], len(cfg.PricesByChain))
	for chainSel, prices := range cfg.PricesByChain {
		tokenPriceUpdates := make([]fee_quoter.InternalTokenPriceUpdate, len(prices.TokenPrices))
		i := 0
		for tokenAddress, price := range prices.TokenPrices {
			tokenPriceUpdates[i] = fee_quoter.InternalTokenPriceUpdate{
				SourceToken: tokenAddress,
				UsdPerToken: price,
			}
			i++
		}
		gasPriceUpdates := make([]fee_quoter.InternalGasPriceUpdate, len(prices.GasPrices))
		i = 0
		for destChainSelector, price := range prices.GasPrices {
			gasPriceUpdates[i] = fee_quoter.InternalGasPriceUpdate{
				DestChainSelector: destChainSelector,
				UsdPerUnitGas:     price,
			}
			i++
		}
		updates[chainSel] = opsutil.EVMCallInput[fee_quoter.InternalPriceUpdates]{
			ChainSelector: chainSel,
			Address:       state.Chains[chainSel].FeeQuoter.Address(),
			CallInput: fee_quoter.InternalPriceUpdates{
				TokenPriceUpdates: tokenPriceUpdates,
				GasPriceUpdates:   gasPriceUpdates,
			},
			NoSend: cfg.MCMS != nil, // If MCMS exists, we do not want to send the transaction.
		}
	}

	return ccipseqs.FeeQuoterUpdatePricesSequenceInput{
		UpdatesByChain: updates,
	}
}

func UpdateFeeQuoterPricesChangeset(e cldf.Environment, cfg UpdateFeeQuoterPricesConfig) (cldf.ChangesetOutput, error) {
	fmt.Println("UPDATING FEEQUOTER EVM")
	if err := cfg.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	s, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	report, err := operations.ExecuteSequence(
		e.OperationsBundle,
		ccipseqs.SuiSupportedFeeQuoterUpdatePricesSequence,
		e.BlockChains.EVMChains(),
		cfg.ToSequenceInput(s),
	)
	return opsutil.AddEVMCallSequenceToCSOutput(e, cldf.ChangesetOutput{}, report, err, s.EVMMCMSStateByChain(), cfg.MCMS, "Call UpdatePrices on FeeQuoters")
}
