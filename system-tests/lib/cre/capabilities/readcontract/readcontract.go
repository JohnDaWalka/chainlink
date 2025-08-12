package readcontract

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	readcontractregistry "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry/v1/readcontract"
	readcontractjobs "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/readcontract"
)

type Capability struct {
}

func New() *Capability {
	return &Capability{}
}

func (c *Capability) Flag() cre.CapabilityFlag {
	return cre.ReadContractCapability
}

func (c *Capability) Validate() error {
	return nil
}

func (c *Capability) JobSpecFn() cre.JobSpecFn {
	return readcontractjobs.JobSpecFn
}

func (c *Capability) OptionalNodeConfigFn() cre.NodeConfigFn {
	return nil
}

func (c *Capability) OptionalGatewayHandlerConfigFn() cre.GatewayHandlerConfigFn {
	return nil
}

func (c *Capability) CapabilityRegistryV1ConfigFn() cre.CapabilityRegistryConfigFn {
	return readcontractregistry.CapabilityRegistryConfigFn
}
