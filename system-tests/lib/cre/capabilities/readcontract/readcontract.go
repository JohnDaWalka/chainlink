package readcontract

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	readcontractregistry "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry/v1/readcontract"
	factory "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
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
	return factory.NewChainSpecificCapabilityJobSpecFactory(
		c.Flag(),
		`'{"chainId":{{.ChainID}},"network":"{{.NetworkFamily}}"}'`,
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

func (c *Capability) OptionalGatewayHandlerConfigFn() cre.GatewayHandlerConfigFn {
	return nil
}

func (c *Capability) CapabilityRegistryV1ConfigFn() cre.CapabilityRegistryConfigFn {
	return readcontractregistry.CapabilityRegistryConfigFn
}
