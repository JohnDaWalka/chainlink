package lane

import "fmt"

type UpdateAptosLanes struct {
	UpdateFeeQuoterDestsConfig  []string
	UpdateFeeQuoterPricesConfig []string
	UpdateOnRampDestsConfig     []string
	UpdateOffRampSourcesConfig  []string
	UpdateRouterRampsConfig     []string
}

func (u *UpdateAptosLanes) SetUpdateFeeQuoterDestsConfig(laneCfg LaneConfig) {
	u.UpdateFeeQuoterDestsConfig = append(u.UpdateFeeQuoterDestsConfig, fmt.Sprintf("%d->%d", laneCfg.Source.Selector, laneCfg.Dest.Selector))
}

func (u *UpdateAptosLanes) SetUpdateFeeQuoterPricesConfig(laneCfg LaneConfig) {
	u.UpdateFeeQuoterPricesConfig = append(u.UpdateFeeQuoterPricesConfig, fmt.Sprintf("%d->%d", laneCfg.Source.Selector, laneCfg.Dest.Selector))
}

func (u *UpdateAptosLanes) SetUpdateOnRampDestsConfig(laneCfg LaneConfig) {
	u.UpdateOnRampDestsConfig = append(u.UpdateOnRampDestsConfig, fmt.Sprintf("%d->%d", laneCfg.Source.Selector, laneCfg.Dest.Selector))
}

func (u *UpdateAptosLanes) SetUpdateOffRampSourcesConfig(laneCfg LaneConfig) {
	u.UpdateOffRampSourcesConfig = append(u.UpdateOffRampSourcesConfig, fmt.Sprintf("%d->%d", laneCfg.Source.Selector, laneCfg.Dest.Selector))
}

func (u *UpdateAptosLanes) SetUpdateRouterRampsConfig(laneCfg LaneConfig) {
	u.UpdateRouterRampsConfig = append(u.UpdateRouterRampsConfig, fmt.Sprintf("%d->%d", laneCfg.Source.Selector, laneCfg.Dest.Selector))
}
