package httpaction

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	httpactionregistry "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry/v1/httpaction"
	httpactionhandler "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/gateway/handlers/httpaction"
	factory "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
)

const httpActionConfigTemplate = `"""
{
	"proxyMode": "{{.ProxyMode}}",
	"incomingRateLimiter": {
		"globalBurst": {{.IncomingGlobalBurst}},
		"globalRPS": {{.IncomingGlobalRPS}},
		"perSenderBurst": {{.IncomingPerSenderBurst}},
		"perSenderRPS": {{.IncomingPerSenderRPS}}
	},
	"outgoingRateLimiter": {
		"globalBurst": {{.OutgoingGlobalBurst}},
		"globalRPS": {{.OutgoingGlobalRPS}},
		"perSenderBurst": {{.OutgoingPerSenderBurst}},
		"perSenderRPS": {{.OutgoingPerSenderRPS}}
	}
}
"""`

type Capability struct{}

func New() *Capability {
	return &Capability{}
}

func (c *Capability) Flag() cre.CapabilityFlag {
	return cre.HTTPActionCapability
}

func (c *Capability) Validate() error {
	return nil
}

func (c *Capability) JobSpecFn() cre.JobSpecFn {
	return factory.NewDonLevelCapabilityJobSpecFactory(
		c.Flag(),
		httpActionConfigTemplate,
		factory.NoOpExtractor, // No runtime values extraction needed
		factory.BinaryPathBuilder,
	).GenerateJobSpecs
}

func (c *Capability) OptionalNodeConfigFn() cre.NodeConfigFn {
	return nil
}

func (c *Capability) OptionalGatewayJobHandlerConfigFn() cre.GatewayHandlerConfigFn {
	return httpactionhandler.HandlerConfigFn
}

func (c *Capability) CapabilityRegistryV1ConfigFn() cre.CapabilityRegistryConfigFn {
	return httpactionregistry.CapabilityRegistryConfigFn
}
