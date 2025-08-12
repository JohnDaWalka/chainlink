package httptrigger

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	httpregistry "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry/v1/httptrigger"
	httptriggerjobs "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/httptrigger"
)

type Capability struct {
}

func New() *Capability {
	return &Capability{}
}

func (c *Capability) Flag() cre.CapabilityFlag {
	return cre.HTTPTriggerCapability
}

func (c *Capability) Validate() error {
	return nil
}

func (c *Capability) JobSpecFactoryFn() cre.JobSpecFactoryFn {
	return httptriggerjobs.JobSpecFn
}

func (c *Capability) OptionalNodeConfigFactoryFn() cre.NodeConfigFactoryFn {
	return nil
}

func (c *Capability) OptionalGatewayHandlerConfigFactoryFn() cre.GatewayHandlerConfigFactoryFn {
	return nil
}

func (c *Capability) CapabilityRegistryV1ConfigFactoryFn() cre.CapabilityRegistryConfigFactoryFn {
	return httpregistry.CapabilityRegistryConfigFn
}
