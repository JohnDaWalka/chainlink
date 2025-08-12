package cron

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	cronregistry "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry/v1/cron"
	cronjobs "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/cron"
)

type Capability struct{}

func (c *Capability) Validate() error {
	return nil
}

func (c *Capability) Flag() cre.CapabilityFlag {
	return cre.CronCapability
}
func (c *Capability) JobSpecFactoryFn() cre.JobSpecFactoryFn {
	return cronjobs.JobSpecFn
}

func (c *Capability) OptionalNodeConfigFactoryFn() cre.NodeConfigFactoryFn {
	return nil
}

func (c *Capability) OptionalGatewayHandlerConfigFactoryFn() cre.GatewayHandlerConfigFactoryFn {
	return nil
}

func (c *Capability) CapabilityRegistryV1ConfigFactoryFn() cre.CapabilityRegistryConfigFactoryFn {
	return cronregistry.CapabilityRegistryConfigFn
}
