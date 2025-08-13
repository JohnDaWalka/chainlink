package gateway

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	gatewayconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/config/gateway"
	gatewayjobs "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/gateway"
)

func New(extraAllowedPorts []int, extraAllowedIPs []string, extraAllowedIPsCIDR []string) *Capability {
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

func (g *Capability) JobSpecFn() cre.JobSpecFn {
	return gatewayjobs.JobSpecFn(g.extraAllowedPorts, g.extraAllowedIPs, g.extraAllowedIPsCIDR)
}

func (g *Capability) OptionalNodeConfigFn() cre.NodeConfigFn {
	return gatewayconfig.GenerateConfigFn
}

func (g *Capability) OptionalGatewayJobHandlerConfigFn() cre.GatewayHandlerConfigFn {
	return nil
}

func (g *Capability) CapabilityRegistryV1ConfigFn() cre.CapabilityRegistryConfigFn {
	return nil
}
