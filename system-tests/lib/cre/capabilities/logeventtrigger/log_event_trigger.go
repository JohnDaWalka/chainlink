package logeventtrigger

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	logeventtriggerregistry "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry/v1/logevent"
	logeventtriggerjobs "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/logevent"
)

type Capability struct {
}

func New() *Capability {
	return &Capability{}
}

func (c *Capability) Flag() cre.CapabilityFlag {
	return cre.LogTriggerCapability
}

func (c *Capability) Validate() error {
	return nil
}

func (c *Capability) JobSpecFn() cre.JobSpecFn {
	return logeventtriggerjobs.JobSpecFn
}

func (c *Capability) OptionalNodeConfigFn() cre.NodeConfigFn {
	return nil
}

func (c *Capability) OptionalGatewayHandlerConfigFn() cre.GatewayHandlerConfigFn {
	return nil
}

func (c *Capability) CapabilityRegistryV1ConfigFn() cre.CapabilityRegistryConfigFn {
	return logeventtriggerregistry.CapabilityRegistryConfigFn
}
