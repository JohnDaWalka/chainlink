package webapitrigger

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	webapiregistry "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry/v1/webapi"
	factory "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
)

const webAPITriggerConfigTemplate = `""`

type Capability struct {
}

func New() *Capability {
	return &Capability{}
}

func (c *Capability) Flag() cre.CapabilityFlag {
	return cre.WebAPITriggerCapability
}

func (c *Capability) Validate() error {
	return nil
}

func (c *Capability) JobSpecFn() cre.JobSpecFn {
	return factory.NewDonLevelCapabilityJobSpecFactory(
		c.Flag(),
		webAPITriggerConfigTemplate,
		factory.NoOpExtractor, // No runtime values extraction needed
		func(_ *cre.JobSpecInput, _ cre.CapabilityConfig) (string, error) {
			return "__builtin_web-api-trigger", nil
		},
	).GenerateJobSpecs
}

func (c *Capability) OptionalNodeConfigFn() cre.NodeConfigFn {
	return nil
}

func (c *Capability) OptionalGatewayHandlerConfigFn() cre.GatewayHandlerConfigFn {
	return nil
}

func (c *Capability) CapabilityRegistryV1ConfigFn() cre.CapabilityRegistryConfigFn {
	return webapiregistry.TriggerCapabilityRegistryConfigFn
}
