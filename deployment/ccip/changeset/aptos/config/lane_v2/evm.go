package lane_v2

import (
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
)

type UpdateEVMLanes struct {
	UpdateFeeQuoterDestsConfig  map[uint64]fee_quoter.FeeQuoterDestChainConfig
	UpdateFeeQuoterPricesConfig v1_6.FeeQuoterPriceUpdatePerSource
	UpdateOnRampDestsConfig     map[uint64]v1_6.OnRampDestinationUpdate
	UpdateOffRampSourcesConfig  map[uint64]v1_6.OffRampSourceUpdate
	UpdateRouterRampsConfig     map[uint64]v1_6.RouterUpdates
}

func (u *UpdateEVMLanes) SetUpdateFeeQuoterDestsConfig(laneCfg LaneConfig) {

}

func (u *UpdateEVMLanes) SetUpdateFeeQuoterPricesConfig(laneCfg LaneConfig) {

}

func (u *UpdateEVMLanes) SetUpdateOnRampDestsConfig(laneCfg LaneConfig) {
	// Setting the destination on the on ramp
	if u.UpdateOnRampDestsConfig == nil {
		u.UpdateOnRampDestsConfig = make(map[uint64]v1_6.OnRampDestinationUpdate)
	}
	u.UpdateOnRampDestsConfig[laneCfg.Dest.Selector] = v1_6.OnRampDestinationUpdate{
		IsEnabled:        !laneCfg.IsDisabled,
		TestRouter:       laneCfg.TestRouter,
		AllowListEnabled: laneCfg.Dest.AllowListEnabled,
	}
}

func (u *UpdateEVMLanes) SetUpdateOffRampSourcesConfig(laneCfg LaneConfig) {

}

func (u *UpdateEVMLanes) SetUpdateRouterRampsConfig(laneCfg LaneConfig) {

}
