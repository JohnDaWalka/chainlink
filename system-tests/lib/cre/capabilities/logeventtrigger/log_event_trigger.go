package logeventtrigger

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	logeventtriggerregistry "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry/v1/logevent"
	factory "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
)

const logEventTriggerConfigTemplate = `"""
{
	"chainId": "{{.ChainID}}",
	"network": "{{.NetworkFamily}}",
	"lookbackBlocks": {{.LookbackBlocks}},
	"pollPeriod": {{.PollPeriod}}
}
"""`

type Capability struct {
}

func New() *Capability {
	return &Capability{}
}

func (c *Capability) Flag() cre.CapabilityFlag {
	return cre.LogTriggerCapability
}

func (c *Capability) Validate() error {
	return nil
}

func (c *Capability) JobSpecFn() cre.JobSpecFn {
	return factory.NewChainSpecificCapabilityJobSpecFactory(
		c.Flag(),
		logEventTriggerConfigTemplate,
		func(chainID uint64, _ *cre.NodeMetadata) map[string]any {
			return map[string]any{
				"ChainID":       chainID,
				"NetworkFamily": "evm",
			}
		},
		factory.BinaryPathBuilder,
	).GenerateJobSpecs
}

func (c *Capability) OptionalNodeConfigFn() cre.NodeConfigFn {
	return nil
}

func (c *Capability) OptionalGatewayJobHandlerConfigFn() cre.GatewayHandlerConfigFn {
	return nil
}

func (c *Capability) CapabilityRegistryV1ConfigFn() cre.CapabilityRegistryConfigFn {
	return logeventtriggerregistry.CapabilityRegistryConfigFn
}
