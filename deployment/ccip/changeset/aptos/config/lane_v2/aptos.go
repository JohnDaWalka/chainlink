package lane

import (
	"math/big"

	aptos_fee_quoter "github.com/smartcontractkit/chainlink-aptos/bindings/ccip/fee_quoter"
	aptos_router "github.com/smartcontractkit/chainlink-aptos/bindings/ccip_router/router"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
)

var defaultOnRampVersion = []byte{1, 6, 0}

type UpdateAptosLanes struct {
	UpdateFeeQuoterDestsConfig  operation.UpdateFeeQuoterDestsInput
	UpdateFeeQuoterPricesConfig operation.UpdateFeeQuoterPricesInput
	UpdateOnRampDestsConfig     operation.UpdateOnRampDestsInput
	UpdateOffRampSourcesConfig  operation.UpdateOffRampSourcesInput
	UpdateRouterRampsConfig     operation.UpdateRouterDestInput
}

func (u *UpdateAptosLanes) SetUpdateFeeQuoterDestsConfig(laneCfg LaneConfig) {
	// Setting the fee quoter destination on the source chain
	if u.UpdateFeeQuoterDestsConfig.Updates == nil {
		u.UpdateFeeQuoterDestsConfig.Updates = make(map[uint64]aptos_fee_quoter.DestChainConfig)
	}
	// Convert generic feequoter configs to Aptos format
	dfqc := laneCfg.Dest.FeeQuoterDestChainConfig
	afqc := aptos_fee_quoter.DestChainConfig{
		IsEnabled:                         dfqc.IsEnabled,
		MaxNumberOfTokensPerMsg:           dfqc.MaxNumberOfTokensPerMsg,
		MaxDataBytes:                      dfqc.MaxDataBytes,
		MaxPerMsgGasLimit:                 dfqc.MaxPerMsgGasLimit,
		DestGasOverhead:                   dfqc.DestGasOverhead,
		DestGasPerPayloadByteBase:         dfqc.DestGasPerPayloadByteBase,
		DestGasPerPayloadByteHigh:         dfqc.DestGasPerPayloadByteHigh,
		DestGasPerPayloadByteThreshold:    dfqc.DestGasPerPayloadByteThreshold,
		DestDataAvailabilityOverheadGas:   dfqc.DestDataAvailabilityOverheadGas,
		DestGasPerDataAvailabilityByte:    dfqc.DestGasPerDataAvailabilityByte,
		DestDataAvailabilityMultiplierBps: dfqc.DestDataAvailabilityMultiplierBps,
		ChainFamilySelector:               dfqc.ChainFamilySelector[:],
		EnforceOutOfOrder:                 dfqc.EnforceOutOfOrder,
		DefaultTokenFeeUsdCents:           dfqc.DefaultTokenFeeUSDCents,
		DefaultTokenDestGasOverhead:       dfqc.DefaultTokenDestGasOverhead,
		DefaultTxGasLimit:                 dfqc.DefaultTxGasLimit,
		GasMultiplierWeiPerEth:            dfqc.GasMultiplierWeiPerEth,
		GasPriceStalenessThreshold:        dfqc.GasPriceStalenessThreshold,
		NetworkFeeUsdCents:                dfqc.NetworkFeeUSDCents,
	}
	u.UpdateFeeQuoterDestsConfig.Updates[laneCfg.Dest.Selector] = afqc
}

func (u *UpdateAptosLanes) SetUpdateFeeQuoterPricesConfig(laneCfg LaneConfig) {
	// Setting gas prices updates
	if u.UpdateFeeQuoterPricesConfig.GasPrices == nil {
		u.UpdateFeeQuoterPricesConfig.GasPrices = make(map[uint64]*big.Int)
	}
	u.UpdateFeeQuoterPricesConfig.GasPrices[laneCfg.Dest.Selector] = laneCfg.Dest.GasPrice

	// Setting token prices updates
	if u.UpdateFeeQuoterPricesConfig.TokenPrices == nil {
		u.UpdateFeeQuoterPricesConfig.TokenPrices = make(map[string]*big.Int)
	}
	for tokenAddr, price := range laneCfg.Source.TokenPrices {
		u.UpdateFeeQuoterPricesConfig.TokenPrices[tokenAddr] = price
	}
}

func (u *UpdateAptosLanes) SetUpdateOnRampDestsConfig(laneCfg LaneConfig) {
	if u.UpdateOnRampDestsConfig.Updates == nil {
		u.UpdateOnRampDestsConfig.Updates = make(map[uint64]v1_6.OnRampDestinationUpdate)
	}
	u.UpdateOnRampDestsConfig.Updates[laneCfg.Dest.Selector] = v1_6.OnRampDestinationUpdate{
		IsEnabled:        !laneCfg.IsDisabled,
		TestRouter:       laneCfg.TestRouter,
		AllowListEnabled: laneCfg.Dest.AllowListEnabled,
	}
}

func (u *UpdateAptosLanes) SetUpdateOffRampSourcesConfig(laneCfg LaneConfig) {
	if u.UpdateOffRampSourcesConfig.Updates == nil {
		u.UpdateOffRampSourcesConfig.Updates = make(map[uint64]v1_6.OffRampSourceUpdate)
	}
	u.UpdateOffRampSourcesConfig.Updates[laneCfg.Source.Selector] = v1_6.OffRampSourceUpdate{
		IsEnabled:                 !laneCfg.IsDisabled,
		TestRouter:                laneCfg.TestRouter,
		IsRMNVerificationDisabled: laneCfg.Source.RMNVerificationDisabled,
	}
}

func (u *UpdateAptosLanes) SetUpdateRouterRampsConfig(laneCfg LaneConfig) {
	if u.UpdateRouterRampsConfig.Updates == nil {
		u.UpdateRouterRampsConfig.Updates = []aptos_router.OnRampSet{}
	}
	onRampVersion := laneCfg.AptosExtraConfigs.OnRampVersion
	if onRampVersion == nil {
		onRampVersion = defaultOnRampVersion
	}
	u.UpdateRouterRampsConfig.Updates = append(u.UpdateRouterRampsConfig.Updates, aptos_router.OnRampSet{
		DestChainSelector: laneCfg.Dest.Selector,
		OnRampVersion:     onRampVersion,
	})
}
