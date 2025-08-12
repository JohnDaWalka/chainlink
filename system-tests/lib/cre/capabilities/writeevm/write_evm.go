package writeevm

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	writeevmregistry "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry/v1/writeevm"
)

type Capability struct {
}

func New() *Capability {
	return &Capability{}
}

func (c *Capability) Flag() cre.CapabilityFlag {
	return cre.WriteEVMCapability
}

func (c *Capability) Validate() error {
	return nil
}

func (c *Capability) JobSpecFactoryFn() cre.JobSpecFactoryFn {
	return nil
}

func (c *Capability) OptionalNodeConfigFactoryFn() cre.NodeConfigFactoryFn {
	return nil
}

func (c *Capability) OptionalGatewayHandlerConfigFactoryFn() cre.GatewayHandlerConfigFactoryFn {
	return nil
}

func (c *Capability) CapabilityRegistryV1ConfigFactoryFn() cre.CapabilityRegistryConfigFactoryFn {
	return writeevmregistry.CapabilityRegistryConfigFn
}
