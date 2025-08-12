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

func (c *Capability) JobSpecFn() cre.JobSpecFn {
	return computejobs.JobSpecFn
}

func (c *Capability) OptionalNodeConfigFn() cre.NodeConfigFn {
	return nil
}

func (c *Capability) OptionalGatewayHandlerConfigFn() cre.GatewayHandlerConfigFn {
	return nil
}

func (c *Capability) CapabilityRegistryV1ConfigFn() cre.CapabilityRegistryConfigFn {
	return computeregistry.CapabilityRegistryConfigFn
}
