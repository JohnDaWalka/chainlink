package evm

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	evmregistry "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry/v1/evm"
	evmjobs "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/evm"
)

type Capability struct {
}

func New() *Capability {
	return &Capability{}
}

func (c *Capability) Validate() error {
	return nil
}

func (c *Capability) Flag() cre.CapabilityFlag {
	return cre.EVMCapability
}

func (c *Capability) JobSpecFactoryFn() cre.JobSpecFactoryFn {
	return evmjobs.JobSpecFn
}

func (c *Capability) OptionalNodeConfigFactoryFn() cre.NodeConfigFactoryFn {
	return nil
}

func (c *Capability) OptionalGatewayHandlerConfigFactoryFn() cre.GatewayHandlerConfigFactoryFn {
	return nil
}

func (c *Capability) CapabilityRegistryV1ConfigFactoryFn() cre.CapabilityRegistryConfigFactoryFn {
	return evmregistry.CapabilityRegistryConfigFn
}
