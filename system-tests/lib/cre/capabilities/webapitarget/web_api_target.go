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

func (c *Capability) JobSpecFn() cre.JobSpecFn {
	return webapijobs.TargetJobSpecFn
}

func (c *Capability) OptionalNodeConfigFn() cre.NodeConfigFn {
	return nil
}

func (c *Capability) OptionalGatewayHandlerConfigFn() cre.GatewayHandlerConfigFn {
	return nil
}

func (c *Capability) CapabilityRegistryV1ConfigFn() cre.CapabilityRegistryConfigFn {
	return webapiregistry.TargetCapabilityRegistryConfigFn
}
