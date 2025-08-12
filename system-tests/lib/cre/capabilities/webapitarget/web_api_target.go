package webapitarget

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	webapiregistry "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry/v1/webapi"
	webapijobs "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/webapi"
)

type Capability struct {
}

func New() *Capability {
	return &Capability{}
}

func (c *Capability) Flag() cre.CapabilityFlag {
	return cre.WebAPITargetCapability
}

func (c *Capability) Validate() error {
	return nil
}

func (c *Capability) JobSpecFactoryFn() cre.JobSpecFactoryFn {
	return webapijobs.TargetJobSpecFn
}

func (c *Capability) OptionalNodeConfigFactoryFn() cre.NodeConfigFactoryFn {
	return nil
}

func (c *Capability) OptionalGatewayHandlerConfigFactoryFn() cre.GatewayHandlerConfigFactoryFn {
	return nil
}

func (c *Capability) CapabilityRegistryV1ConfigFactoryFn() cre.CapabilityRegistryConfigFactoryFn {
	return webapiregistry.TargetCapabilityRegistryConfigFn
}
