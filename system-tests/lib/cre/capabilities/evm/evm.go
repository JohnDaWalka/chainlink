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

func (c *Capability) JobSpecFn() cre.JobSpecFn {
	return evmjobs.JobSpecFn
}

func (c *Capability) OptionalNodeConfigFn() cre.NodeConfigFn {
	return nil
}

func (c *Capability) OptionalGatewayJobHandlerConfigFn() cre.GatewayHandlerConfigFn {
	return nil
}

func (c *Capability) CapabilityRegistryV1ConfigFn() cre.CapabilityRegistryConfigFn {
	return evmregistry.CapabilityRegistryConfigFn
}
