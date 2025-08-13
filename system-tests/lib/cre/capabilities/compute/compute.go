package compute

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	computeregistry "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry/v1/compute"
	factory "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
)

const customComputeConfigTemplate = `"""
NumWorkers = {{.NumWorkers}}
[rateLimiter]
globalRPS = {{.GlobalRPS}}
globalBurst = {{.GlobalBurst}}
perSenderRPS = {{.PerSenderRPS}}
perSenderBurst = {{.PerSenderBurst}}
"""`

type Capability struct {
}

func New() *Capability {
	return &Capability{}
}

func (c *Capability) Flag() cre.CapabilityFlag {
	return cre.CustomComputeCapability
}

func (c *Capability) Validate() error {
	return nil
}

func (c *Capability) JobSpecFn() cre.JobSpecFn {
	return factory.NewDonLevelCapabilityJobSpecFactory(
		c.Flag(),
		customComputeConfigTemplate,
		factory.NoOpExtractor, // No runtime values extraction needed
		func(_ *cre.JobSpecInput, _ cre.CapabilityConfig) (string, error) {
			return "__builtin_custom-compute-action", nil
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
	return computeregistry.CapabilityRegistryConfigFn
}
