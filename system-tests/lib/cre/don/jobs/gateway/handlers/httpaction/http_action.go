package httpaction

import (
	coregateway "github.com/smartcontractkit/chainlink/v2/core/services/gateway"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
)

var HandlerConfigFn = func(donMetadata []*cre.DonMetadata) (cre.HandlerTypeToConfig, error) {
	// if any of the DONs have http action or http trigger capability, we need to add a http handler to the jobspec for the gateway node
	if !don.AnyDonHasCapability(donMetadata, cre.HTTPActionCapability) && !don.AnyDonHasCapability(donMetadata, cre.HTTPTriggerCapability) {
		return nil, nil
	}

	return map[string]string{coregateway.HTTPCapabilityType: `
ServiceName = "workflows"
[gatewayConfig.Dons.Handlers.Config]
maxTriggerRequestDurationMs = 5_000
[gatewayConfig.Dons.Handlers.Config.NodeRateLimiter]
globalBurst = 10
globalRPS = 50
perSenderBurst = 10
perSenderRPS = 10
[gatewayConfig.Dons.Handlers.Config.UserRateLimiter]
globalBurst = 10
globalRPS = 50
perSenderBurst = 10
perSenderRPS = 10`}, nil
}
