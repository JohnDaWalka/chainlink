package vault

import (
	coregateway "github.com/smartcontractkit/chainlink/v2/core/services/gateway"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
)

var HandlerConfigFn = func(donMetadata []*cre.DonMetadata) (cre.HandlerTypeToConfig, error) {
	if !don.AnyDonHasCapability(donMetadata, cre.VaultCapability) {
		return nil, nil
	}

	return map[string]string{coregateway.VaultHandlerType: `
[gatewayConfig.Dons.Handlers.Config]
requestTimeoutSec = 30
[gatewayConfig.Dons.Handlers.Config.NodeRateLimiter]
globalBurst = 10
globalRPS = 50
perSenderBurst = 10
perSenderRPS = 10
`}, nil
}
