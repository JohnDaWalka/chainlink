package logeventtrigger

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	logeventtriggerregistry "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry/v1/logevent"
	logeventtriggerjobs "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/logevent"
)

type Capability struct {
}

func (c *Capability) Flag() cre.CapabilityFlag {
	return cre.LogTriggerCapability
}

func (c *Capability) Validate() error {
	return nil
}

func (c *Capability) JobSpecFactoryFn() cre.JobSpecFactoryFn {
	return logeventtriggerjobs.JobSpecFn
}

func (c *Capability) OptionalNodeConfigFactoryFn() cre.NodeConfigFactoryFn {
	return nil
}

func (c *Capability) OptionalGatewayHandlerConfigFactoryFn() cre.GatewayHandlerConfigFactoryFn {
	return nil
}

func (c *Capability) CapabilityRegistryV1ConfigFactoryFn() cre.CapabilityRegistryConfigFactoryFn {
	return logeventtriggerregistry.CapabilityRegistryConfigFn
}
