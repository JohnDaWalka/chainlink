package compute

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	computeregistry "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry/v1/compute"
	computejobs "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/compute"
)

type Capability struct {
}

func New() *Capability {
	return &Capability{}
}

func (c *Capability) Flag() cre.CapabilityFlag {
	return cre.CustomComputeCapability
}

func (c *Capability) Validate() error {
	return nil
}

func (c *Capability) JobSpecFactoryFn() cre.JobSpecFactoryFn {
	return computejobs.JobSpecFn
}

func (c *Capability) OptionalNodeConfigFactoryFn() cre.NodeConfigFactoryFn {
	return nil
}

func (c *Capability) OptionalGatewayHandlerConfigFactoryFn() cre.GatewayHandlerConfigFactoryFn {
	return nil
}

func (c *Capability) CapabilityRegistryV1ConfigFactoryFn() cre.CapabilityRegistryConfigFactoryFn {
	return computeregistry.CapabilityRegistryConfigFn
}
