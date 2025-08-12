package gateway

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	gatewayconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/config/gateway"
	gatewayjobs "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/gateway"
)

func NewGatewayCapability(extraAllowedPorts []int, extraAllowedIPs []string, extraAllowedIPsCIDR []string) *Capability {
	return &Capability{
		extraAllowedPorts:   extraAllowedPorts,
		extraAllowedIPs:     extraAllowedIPs,
		extraAllowedIPsCIDR: extraAllowedIPsCIDR,
	}
}

type Capability struct {
	extraAllowedPorts   []int
	extraAllowedIPs     []string
	extraAllowedIPsCIDR []string
}

func (g *Capability) Flag() cre.CapabilityFlag {
	return cre.GatewayDON
}

func (g *Capability) Validate() error {
	return nil
}

func (g *Capability) JobSpecFactoryFn() cre.JobSpecFactoryFn {
	return gatewayjobs.JobSpecFn(g.extraAllowedPorts, g.extraAllowedIPs, g.extraAllowedIPsCIDR)
}

func (g *Capability) OptionalNodeConfigFactoryFn() cre.NodeConfigFactoryFn {
	return gatewayconfig.GenerateConfigFn
}

func (g *Capability) OptionalGatewayHandlerConfigFactoryFn() cre.GatewayHandlerConfigFactoryFn {
	return nil
}

func (g *Capability) CapabilityRegistryV1ConfigFactoryFn() cre.CapabilityRegistryConfigFactoryFn {
	return nil
}
