package mock

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	mockregistry "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry/v1/mock"
	factory "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
)

const mockConfigTemplate = `"""port={{.Port}}"""`

type Capability struct {
}

func New() *Capability {
	return &Capability{}
}

func (c *Capability) Flag() cre.CapabilityFlag {
	return cre.MockCapability
}

func (c *Capability) Validate() error {
	return nil
}

func (c *Capability) JobSpecFn() cre.JobSpecFn {
	return factory.NewDonLevelFactory(
		c.Flag(),
		mockConfigTemplate,
		factory.NoOpExtractor, // No runtime values extraction needed
		factory.BinaryPathBuilder,
	).GenerateJobSpecs
}

func (c *Capability) OptionalNodeConfigFn() cre.NodeConfigFn {
	return nil
}

func (c *Capability) OptionalGatewayHandlerConfigFn() cre.GatewayHandlerConfigFn {
	return nil
}

func (c *Capability) CapabilityRegistryV1ConfigFn() cre.CapabilityRegistryConfigFn {
	return mockregistry.CapabilityRegistryConfigFn
}
