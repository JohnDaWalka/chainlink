package consensus

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	consensusregistry "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry/v1/consensus"
	consensusjobs "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/consensus"
)

type CapabilityV1 struct {
	chainID uint64
}

func NewV1(chainID uint64) *CapabilityV1 {
	return &CapabilityV1{
		chainID: chainID,
	}
}

func (c *CapabilityV1) Validate() error {
	if c.chainID == 0 {
		return fmt.Errorf("chainID is required, got %d", c.chainID)
	}
	return nil
}

func (c *CapabilityV1) Flag() cre.CapabilityFlag {
	return cre.ConsensusCapability
}

func (c *CapabilityV1) JobSpecFactoryFn() cre.JobSpecFactoryFn {
	return consensusjobs.V1JobSpecFn(c.chainID)
}

func (c *CapabilityV1) OptionalNodeConfigFactoryFn() cre.NodeConfigFactoryFn {
	return nil
}

func (c *CapabilityV1) OptionalGatewayHandlerConfigFactoryFn() cre.GatewayHandlerConfigFactoryFn {
	return nil
}

func (c *CapabilityV1) CapabilityRegistryV1ConfigFactoryFn() cre.CapabilityRegistryConfigFactoryFn {
	return consensusregistry.CapabilityV1RegistryConfigFn
}

type CapabilityV2 struct {
	chainID uint64
}

func NewV2(chainID uint64) *CapabilityV2 {
	return &CapabilityV2{
		chainID: chainID,
	}
}

func (c *CapabilityV2) Validate() error {
	if c.chainID == 0 {
		return fmt.Errorf("chainID is required, got %d", c.chainID)
	}
	return nil
}

func (c *CapabilityV2) Flag() cre.CapabilityFlag {
	return cre.ConsensusCapabilityV2
}

func (c *CapabilityV2) JobSpecFactoryFn() cre.JobSpecFactoryFn {
	return consensusjobs.V2JobSpecFn
}

func (c *CapabilityV2) OptionalNodeConfigFactoryFn() cre.NodeConfigFactoryFn {
	return nil
}

func (c *CapabilityV2) OptionalGatewayHandlerConfigFactoryFn() cre.GatewayHandlerConfigFactoryFn {
	return nil
}

func (c *CapabilityV2) CapabilityRegistryV1ConfigFactoryFn() cre.CapabilityRegistryConfigFactoryFn {
	return consensusregistry.ConsensusV2CapabilityFactoryFn
}
