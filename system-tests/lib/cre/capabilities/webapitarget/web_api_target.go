package webapitarget

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	webapiregistry "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry/v1/webapi"
	factory "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
)

const webAPITargetConfigTemplate = `"""
[rateLimiter]
GlobalRPS = {{.GlobalRPS}}
GlobalBurst = {{.GlobalBurst}}
PerSenderRPS = {{.PerSenderRPS}}
PerSenderBurst = {{.PerSenderBurst}}
"""`

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
	return factory.NewDonLevelCapabilityJobSpecFactory(
		c.Flag(),
		webAPITargetConfigTemplate,
		factory.NoOpExtractor, // No runtime values extraction needed
		func(_ *cre.JobSpecInput, _ cre.CapabilityConfig) (string, error) {
			return "__builtin_web-api-target", nil
		},
	).GenerateJobSpecs
}

func (c *Capability) OptionalNodeConfigFn() cre.NodeConfigFn {
	return nil
}

func (c *Capability) OptionalGatewayJobHandlerConfigFn() cre.GatewayHandlerConfigFn {
	return nil
}

func (c *Capability) CapabilityRegistryV1ConfigFn() cre.CapabilityRegistryConfigFn {
	return webapiregistry.TargetCapabilityRegistryConfigFn
}
