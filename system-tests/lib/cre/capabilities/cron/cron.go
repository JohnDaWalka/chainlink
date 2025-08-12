package cron

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	cronregistry "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry/v1/cron"
	cronjobs "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/cron"
)

type Capability struct{}

func New() *Capability {
	return &Capability{}
}

func (c *Capability) Validate() error {
	return nil
}

func (c *Capability) Flag() cre.CapabilityFlag {
	return cre.CronCapability
}
func (c *Capability) JobSpecFn() cre.JobSpecFn {
	return cronjobs.JobSpecFn
}

func (c *Capability) OptionalNodeConfigFn() cre.NodeConfigFn {
	return nil
}

func (c *Capability) OptionalGatewayHandlerConfigFn() cre.GatewayHandlerConfigFn {
	return nil
}

func (c *Capability) CapabilityRegistryV1ConfigFn() cre.CapabilityRegistryConfigFn {
	return cronregistry.CapabilityRegistryConfigFn
}
