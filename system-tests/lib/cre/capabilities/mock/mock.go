package mock

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	mockregistry "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry/v1/mock"
	mockjobs "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/mock"
)

type Capability struct {
}

func (c *Capability) Flag() cre.CapabilityFlag {
	return cre.MockCapability
}

func (c *Capability) Validate() error {
	return nil
}

func (c *Capability) JobSpecFactoryFn() cre.JobSpecFactoryFn {
	return mockjobs.JobSpecFn
}

func (c *Capability) OptionalNodeConfigFactoryFn() cre.NodeConfigFactoryFn {
	return nil
}

func (c *Capability) OptionalGatewayHandlerConfigFactoryFn() cre.GatewayHandlerConfigFactoryFn {
	return nil
}

func (c *Capability) CapabilityRegistryV1ConfigFactoryFn() cre.CapabilityRegistryConfigFactoryFn {
	return mockregistry.CapabilityRegistryConfigFn
}
