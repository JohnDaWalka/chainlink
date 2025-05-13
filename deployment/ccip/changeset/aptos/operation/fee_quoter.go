package operation

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip"
	aptos_fee_quoter "github.com/smartcontractkit/chainlink-aptos/bindings/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	aptosmcms "github.com/smartcontractkit/mcms/sdk/aptos"
	"github.com/smartcontractkit/mcms/types"
)

// UpdateFeeQuoterDestsInput contains configuration for updating FeeQuoter destination configs
type UpdateFeeQuoterDestsInput struct {
	MCMSAddress aptos.AccountAddress
	Updates     map[uint64]aptos_fee_quoter.DestChainConfig
}

// UpdateFeeQuoterDestsOp operation to update FeeQuoter destination configurations
var UpdateFeeQuoterDestsOp = operations.NewOperation(
	"update-fee-quoter-dests-op",
	Version1_0_0,
	"Updates FeeQuoter destination chain configurations",
	updateFeeQuoterDests,
)

func updateFeeQuoterDests(b operations.Bundle, deps AptosDeps, in UpdateFeeQuoterDestsInput) ([]types.Transaction, error) {
	// Bind CCIP Package
	ccipAddress := deps.OnChainState.CCIPAddress
	ccipBind := ccip.Bind(ccipAddress, deps.AptosChain.Client)

	// Process each destination chain config update
	var txs []types.Transaction

	for destChainSelector, destConfig := range in.Updates {
		// Encode the update operation
		moduleInfo, function, _, args, err := ccipBind.FeeQuoter().Encoder().ApplyDestChainConfigUpdates(
			destChainSelector,
			destConfig.IsEnabled,
			destConfig.MaxNumberOfTokensPerMsg,
			destConfig.MaxDataBytes,
			destConfig.MaxPerMsgGasLimit,
			destConfig.DestGasOverhead,
			destConfig.DestGasPerPayloadByteBase,
			destConfig.DestGasPerPayloadByteHigh,
			destConfig.DestGasPerPayloadByteThreshold,
			destConfig.DestDataAvailabilityOverheadGas,
			destConfig.DestGasPerDataAvailabilityByte,
			destConfig.DestDataAvailabilityMultiplierBps,
			destConfig.ChainFamilySelector,
			destConfig.EnforceOutOfOrder,
			destConfig.DefaultTokenFeeUsdCents,
			destConfig.DefaultTokenDestGasOverhead,
			destConfig.DefaultTxGasLimit,
			destConfig.GasMultiplierWeiPerEth,
			destConfig.GasPriceStalenessThreshold,
			destConfig.NetworkFeeUsdCents,
		)
		if err != nil {
			return []types.Transaction{}, fmt.Errorf("failed to encode ApplyDestChainConfigUpdates for chain %d: %w", destChainSelector, err)
		}

		additionalFields := aptosmcms.AdditionalFields{
			PackageName: moduleInfo.PackageName,
			ModuleName:  moduleInfo.ModuleName,
			Function:    function,
		}
		afBytes, err := json.Marshal(additionalFields)
		if err != nil {
			return []types.Transaction{}, fmt.Errorf("failed to marshal additional fields: %w", err)
		}

		txs = append(txs, types.Transaction{
			To:               ccipAddress.StringLong(),
			Data:             aptosmcms.ArgsToData(args),
			AdditionalFields: afBytes,
		})

		b.Logger.Infow("Adding FeeQuoter destination config update operation",
			"destChainSelector", destChainSelector,
			"isEnabled", destConfig.IsEnabled)
	}

	return txs, nil
}

// UpdateFeeQuoterPricesInput contains configuration for updating FeeQuoter price configs
type UpdateFeeQuoterPricesInput struct {
	MCMSAddress aptos.AccountAddress
	Prices      FeeQuoterPriceUpdatePerSource
}

type FeeQuoterPriceUpdatePerSource struct {
	TokenPrices map[string]*big.Int // token address (string) -> price
	GasPrices   map[uint64]*big.Int // dest chain -> gas price
}

// UpdateFeeQuoterPricesOp operation to update FeeQuoter prices
var UpdateFeeQuoterPricesOp = operations.NewOperation(
	"update-fee-quoter-prices-op",
	Version1_0_0,
	"Updates FeeQuoter token and gas prices",
	updateFeeQuoterPrices,
)

func updateFeeQuoterPrices(b operations.Bundle, deps AptosDeps, in UpdateFeeQuoterPricesInput) ([]types.Transaction, error) {
	var txs []types.Transaction

	// Bind CCIP Package
	ccipAddress := deps.OnChainState.CCIPAddress
	ccipBind := ccip.Bind(ccipAddress, deps.AptosChain.Client)

	//// Bind MCMS Package
	// mcmsAddress := deps.OnChainState.MCMSAddress
	// mcmsBind := mcms.Bind(mcmsAddress, deps.AptosChain.Client)

	// // Add CCIP Owner address to update token prices allow list
	// // TODO: add a check here if MCMS is already in the allow list
	// // TODO: don't call two contracts in the same OP
	// ccipOwnerAddress, err := mcmsBind.MCMSRegistry().GetRegisteredOwnerAddress(nil, ccipAddress)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get CCIP owner address: %w", err)
	// }
	// moduleInfo, function, _, args, err := ccipBind.Auth().Encoder().ApplyAllowedOfframpUpdates(nil, []aptos.AccountAddress{ccipOwnerAddress})
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to encode ApplyAllowedOfframpUpdates: %w", err)
	// }
	// additionalFields := aptosmcms.AdditionalFields{
	// 	PackageName: moduleInfo.PackageName,
	// 	ModuleName:  moduleInfo.ModuleName,
	// 	Function:    function,
	// }
	// afBytes, err := json.Marshal(additionalFields)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to marshal additional fields: %w", err)
	// }

	// txs = append(txs, types.Transaction{
	// 	To:               ccipAddress.StringLong(),
	// 	Data:             aptosmcms.ArgsToData(args),
	// 	AdditionalFields: afBytes,
	// })

	// // Convert token prices and gas prices to format expected by Aptos contract
	// var sourceTokens []aptos.AccountAddress
	// var sourceUsdPerToken []*big.Int
	// var gasDestChainSelectors []uint64
	// var gasUsdPerUnitGas []*big.Int

	// in.Prices.TokenPrices = map[string]*big.Int{"0x3b17dad1bdd88f337712cc2f6187bb741d56da467320373fd9198262cc93de76": big.NewInt(400000000)}
	// // Process token prices if any
	// for tokenAddr, price := range in.Prices.TokenPrices {
	// 	address := aptos.AccountAddress{}
	// 	err := address.ParseStringRelaxed(tokenAddr)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to parse Aptos token address %s: %w", tokenAddr, err)
	// 	}
	// 	sourceTokens = append(sourceTokens, address)
	// 	sourceUsdPerToken = append(sourceUsdPerToken, price)
	// }

	// // Process gas prices if any
	// for destChainSel, gasPrice := range in.Prices.GasPrices {
	// 	gasDestChainSelectors = append(gasDestChainSelectors, destChainSel)
	// 	gasUsdPerUnitGas = append(gasUsdPerUnitGas, gasPrice)
	// }

	// // Generate MCMS tx to update prices
	// if len(sourceTokens) == 0 && len(gasDestChainSelectors) == 0 {
	// 	b.Logger.Infow("No price updates to apply")
	// 	return txs, nil
	// }

	// // Encode the update tx
	// moduleInfo, function, _, args, err := ccipBind.FeeQuoter().Encoder().UpdatePrices(
	// 	sourceTokens,
	// 	sourceUsdPerToken,
	// 	gasDestChainSelectors,
	// 	gasUsdPerUnitGas,
	// )
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to encode UpdatePrices: %w", err)
	// }

	// additionalFields := aptosmcms.AdditionalFields{
	// 	PackageName: moduleInfo.PackageName,
	// 	ModuleName:  moduleInfo.ModuleName,
	// 	Function:    function,
	// }
	// afBytes, err := json.Marshal(additionalFields)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to marshal additional fields: %w", err)
	// }

	// txs = append(txs, types.Transaction{
	// 	To:               ccipAddress.StringLong(),
	// 	Data:             aptosmcms.ArgsToData(args),
	// 	AdditionalFields: afBytes,
	// })

	// // Encode the update tx
	// token := aptos.AccountAddress{}
	// _ = token.ParseStringRelaxed("0x3b17dad1bdd88f337712cc2f6187bb741d56da467320373fd9198262cc93de76")
	// moduleInfo, function, _, args, err = ccipBind.FeeQuoter().Encoder().ApplyPremiumMultiplierWeiPerEthUpdates([]aptos.AccountAddress{token}, []uint64{1})
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to encode UpdatePrices: %w", err)
	// }

	// additionalFields = aptosmcms.AdditionalFields{
	// 	PackageName: moduleInfo.PackageName,
	// 	ModuleName:  moduleInfo.ModuleName,
	// 	Function:    function,
	// }

	// afBytes, err = json.Marshal(additionalFields)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to marshal additional fields: %w", err)
	// }

	// txs = append(txs, types.Transaction{
	// 	To:               ccipAddress.StringLong(),
	// 	Data:             aptosmcms.ArgsToData(args),
	// 	AdditionalFields: afBytes,
	// })

	// Encode the update tx
	token := aptos.AccountAddress{}
	_ = token.ParseStringRelaxed("0x3b17dad1bdd88f337712cc2f6187bb741d56da467320373fd9198262cc93de76")
	moduleInfo, function, _, args, err := ccipBind.FeeQuoter().Encoder().ApplyTokenTransferFeeConfigUpdates(
		14767482510784806043,
		[]aptos.AccountAddress{token},
		[]uint32{1},
		[]uint32{100000},
		[]uint16{0},
		[]uint32{1000},
		[]uint32{1000},
		[]bool{true},
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to encode UpdatePrices: %w", err)
	}

	additionalFields := aptosmcms.AdditionalFields{
		PackageName: moduleInfo.PackageName,
		ModuleName:  moduleInfo.ModuleName,
		Function:    function,
	}

	afBytes, err := json.Marshal(additionalFields)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal additional fields: %w", err)
	}

	txs = append(txs, types.Transaction{
		To:               ccipAddress.StringLong(),
		Data:             aptosmcms.ArgsToData(args),
		AdditionalFields: afBytes,
	})

	return txs, nil
}

// AddTokenTransferFeeInput
type AddTokenTransferFeeInput struct {
	MCMSAddress aptos.AccountAddress
	FeeConfigs  map[uint64]FeeConfigs
}

type FeeConfigs struct {
	TokenAddress aptos.AccountAddress
	aptos_fee_quoter.TokenTransferFeeConfig
}

// AddTokenTransferFeeOp operation to update FeeQuoter prices
var AddTokenTransferFeeOp = operations.NewOperation(
	"add-token-transfer-fee-op",
	Version1_0_0,
	"Add token transfer fee",
	addTokenTransferFee,
)

func addTokenTransferFee(b operations.Bundle, deps AptosDeps, in AddTokenTransferFeeInput) ([]types.Transaction, error) {
	var txs []types.Transaction

	// Bind CCIP Package
	ccipAddress := deps.OnChainState.CCIPAddress
	ccipBind := ccip.Bind(ccipAddress, deps.AptosChain.Client)

	// Encode the update tx
	for destSel, configs := range in.FeeConfigs {
		moduleInfo, function, _, args, err := ccipBind.FeeQuoter().Encoder().ApplyTokenTransferFeeConfigUpdates(
			destSel,
			[]aptos.AccountAddress{configs.TokenAddress},
			[]uint32{configs.MinFeeUsdCents},    // addMinFeeUsdCents
			[]uint32{configs.MaxFeeUsdCents},    // addMaxFeeUsdCents
			[]uint16{configs.DeciBps},           // addDeciBps
			[]uint32{configs.DestGasOverhead},   // addDestGasOverhead
			[]uint32{configs.DestBytesOverhead}, // addDestBytesOverhead
			[]bool{configs.IsEnabled},           // addIsEnabled
			[]aptos.AccountAddress{},            // TODO: enable token removal
		)
		if err != nil {
			return nil, fmt.Errorf("failed to encode UpdatePrices: %w", err)
		}

		additionalFields := aptosmcms.AdditionalFields{
			PackageName: moduleInfo.PackageName,
			ModuleName:  moduleInfo.ModuleName,
			Function:    function,
		}

		afBytes, err := json.Marshal(additionalFields)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal additional fields: %w", err)
		}

		txs = append(txs, types.Transaction{
			To:               ccipAddress.StringLong(),
			Data:             aptosmcms.ArgsToData(args),
			AdditionalFields: afBytes,
		})
	}

	return txs, nil
}
