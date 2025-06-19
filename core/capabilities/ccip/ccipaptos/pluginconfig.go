package ccipaptos

import (
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipsui"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ocrimpls"
)

// initializePluginConfig returns a PluginConfig for Aptos chains.
func initializePluginConfigFunc(chainselFamily string) ccipcommon.InitFunction {
	return func(lggr logger.Logger, extraDataCodec ccipcommon.ExtraDataCodec) ccipcommon.PluginConfig {
		var cwProvider ccipcommon.ChainRWProvider
		if chainselFamily == chainsel.FamilyAptos {
			cwProvider = ChainCWProvider{}
		} else {
			cwProvider = ccipsui.ChainCWProvider{}
		}

		return ccipcommon.PluginConfig{
			CommitPluginCodec:          NewCommitPluginCodecV1(),
			ExecutePluginCodec:         NewExecutePluginCodecV1(extraDataCodec),
			MessageHasher:              NewMessageHasherV1(logger.Sugared(lggr).Named(chainsel.FamilyAptos).Named("MessageHasherV1"), extraDataCodec),
			TokenDataEncoder:           NewAptosTokenDataEncoder(),
			GasEstimateProvider:        NewGasEstimateProvider(),
			RMNCrypto:                  nil,
			ContractTransmitterFactory: ocrimpls.NewAptosContractTransmitterFactory(extraDataCodec),
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
