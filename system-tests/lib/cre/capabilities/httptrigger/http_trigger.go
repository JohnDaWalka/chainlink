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

func (c *Capability) JobSpecFn() cre.JobSpecFn {
	return httptriggerjobs.JobSpecFn
}

func (c *Capability) OptionalNodeConfigFn() cre.NodeConfigFn {
	return nil
}

func (c *Capability) OptionalGatewayHandlerConfigFn() cre.GatewayHandlerConfigFn {
	return nil
}

func (c *Capability) CapabilityRegistryV1ConfigFn() cre.CapabilityRegistryConfigFn {
	return httpregistry.CapabilityRegistryConfigFn
}
