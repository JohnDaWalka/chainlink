package httpaction

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	httpactionregistry "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry/v1/httpaction"
	httpactionhandler "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/gateway/handlers/httpaction"
	httpactionjobs "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/httpaction"
)

type Capability struct{}

func New() *Capability {
	return &Capability{}
}

func (c *Capability) Flag() cre.CapabilityFlag {
	return cre.HTTPActionCapability
}

func (c *Capability) Validate() error {
	return nil
}

func (c *Capability) JobSpecFn() cre.JobSpecFn {
	return httpactionjobs.JobSpecFn
}

func (c *Capability) OptionalNodeConfigFn() cre.NodeConfigFn {
	return nil
}

func (c *Capability) OptionalGatewayHandlerConfigFn() cre.GatewayHandlerConfigFn {
	return httpactionhandler.HandlerConfigFn
}

func (c *Capability) CapabilityRegistryV1ConfigFn() cre.CapabilityRegistryConfigFn {
	return httpactionregistry.CapabilityRegistryConfigFn
}
