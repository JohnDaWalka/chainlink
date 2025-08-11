package ccipaptos

import (
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipsui"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ocrimpls"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
)

// initializePluginConfig returns a PluginConfig for Aptos chains.
func initializePluginConfigFunc(chainselFamily string) ccipcommon.InitFunction {
	return func(lggr logger.Logger, extraDataCodec ccipcommon.ExtraDataCodec) ccipcommon.PluginConfig {
		var cwProvider ccipcommon.ChainRWProvider
		var transmitterFactory types.ContractTransmitterFactory

		if chainselFamily == chainsel.FamilyAptos {
			cwProvider = ChainCWProvider{}
			transmitterFactory = ocrimpls.NewAptosContractTransmitterFactory(extraDataCodec)
		} else {
			cwProvider = ccipsui.ChainCWProvider{}
			transmitterFactory = ocrimpls.NewSuiContractTransmitterFactory(extraDataCodec)
		}

		return ccipcommon.PluginConfig{
			CommitPluginCodec:          NewCommitPluginCodecV1(),
			ExecutePluginCodec:         NewExecutePluginCodecV1(extraDataCodec),
			MessageHasher:              NewMessageHasherV1(logger.Sugared(lggr).Named(chainselFamily).Named("MessageHasherV1"), extraDataCodec),
			TokenDataEncoder:           NewAptosTokenDataEncoder(),
			GasEstimateProvider:        NewGasEstimateProvider(),
			RMNCrypto:                  nil,
			ChainAccessorFactory:       AptosChainAccessorFactory{},
			ContractTransmitterFactory: transmitterFactory,
			ChainRW:                    cwProvider,
			ExtraDataCodec:             ExtraDataDecoder{},
			AddressCodec:               AddressCodec{},
		}
	}
}

func init() {
	// Register the Aptos and Sui plugin config factory
	ccipcommon.RegisterPluginConfig(chainsel.FamilyAptos, initializePluginConfigFunc(chainsel.FamilyAptos))
	ccipcommon.RegisterPluginConfig(chainsel.FamilySui, initializePluginConfigFunc(chainsel.FamilySui))
}

// public entry fun commit(
// 	caller: &signer,
// 	report_context: vector<vector<u8>>,
// 	report: vector<u8>,
// 	signatures: vector<vector<u8>>

// public fun commit(
// 	ref: &mut CCIPObjectRef,
// 	state: &mut OffRampState,
// 	clock: &clock::Clock,
// 	report_context: vector<vector<u8>>,
// 	report: vector<u8>,
// 	signatures: vector<vector<u8>>,
// 	ctx: &mut TxContext
// )
