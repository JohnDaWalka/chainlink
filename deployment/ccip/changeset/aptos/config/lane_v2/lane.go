package lane

type LaneConfig struct {
	Source            ChainDefinition
	Dest              ChainDefinition
	IsDisabled        bool
	TestRouter        bool
	AptosExtraConfigs AptosExtraConfigs
}

type AptosExtraConfigs struct {
	OnRampVersion []byte
}

type UpdateLanesCfg interface {
	SetUpdateFeeQuoterDestsConfig(laneCfg LaneConfig)
	SetUpdateFeeQuoterPricesConfig(laneCfg LaneConfig)
	SetUpdateOnRampDestsConfig(laneCfg LaneConfig)
	SetUpdateOffRampSourcesConfig(laneCfg LaneConfig)
	SetUpdateRouterRampsConfig(laneCfg LaneConfig)
}

func ToUpdateLanesConfig(laneCfg LaneConfig, sourceChainUpdates, destChainUpdates UpdateLanesCfg) error {
	toSourceUpdates(laneCfg, sourceChainUpdates)
	toDestUpdates(laneCfg, destChainUpdates)

	return nil
}

func toSourceUpdates(laneCfg LaneConfig, sourceChainUpdates UpdateLanesCfg) error {
	// Setting the destination on the on ramp
	sourceChainUpdates.SetUpdateOnRampDestsConfig(laneCfg)
	// Setting gas prices updates
	sourceChainUpdates.SetUpdateFeeQuoterPricesConfig(laneCfg)
	// Setting the fee quoter destination on the source chain
	sourceChainUpdates.SetUpdateFeeQuoterDestsConfig(laneCfg)
	// Setting Router OnRamp updates
	sourceChainUpdates.SetUpdateRouterRampsConfig(laneCfg)

	return nil
}

func toDestUpdates(laneCfg LaneConfig, destChainUpdates UpdateLanesCfg) error {
	// Setting off ramp sources updates
	destChainUpdates.SetUpdateOffRampSourcesConfig(laneCfg)
	// Setting router off ramp updates
	destChainUpdates.SetUpdateRouterRampsConfig(laneCfg)
	return nil
}
