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

func (c *Capability) JobSpecFn() cre.JobSpecFn {
	return nil
}

func (c *Capability) OptionalNodeConfigFn() cre.NodeConfigFn {
	return nil
}

func (c *Capability) OptionalGatewayJobHandlerConfigFn() cre.GatewayHandlerConfigFn {
	return nil
}

func (c *Capability) CapabilityRegistryV1ConfigFn() cre.CapabilityRegistryConfigFn {
	return writeevmregistry.CapabilityRegistryConfigFn
}
