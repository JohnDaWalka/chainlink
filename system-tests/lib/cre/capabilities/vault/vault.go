package vault

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	vaultregistry "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry/v1/vault"
	vaulthandler "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/gateway/handlers/vault"
	vaultjobs "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/vault"
)

func NewVaultCapability(chainID uint64) *Capability {
	return &Capability{
		chainID: chainID,
	}
}

type Capability struct {
	chainID uint64
}

func (c *Capability) Flag() cre.CapabilityFlag {
	return cre.VaultCapability
}

func (c *Capability) Validate() error {
	if c.chainID == 0 {
		return fmt.Errorf("chainID is required, got %d", c.chainID)
	}
	return nil
}

func (c *Capability) JobSpecFactoryFn() cre.JobSpecFactoryFn {
	return vaultjobs.JobSpecFn(c.chainID)
}

func (c *Capability) OptionalNodeConfigFactoryFn() cre.NodeConfigFactoryFn {
	return nil
}

func (c *Capability) OptionalGatewayHandlerConfigFactoryFn() cre.GatewayHandlerConfigFactoryFn {
	return vaulthandler.HandlerConfigFn
}

func (c *Capability) CapabilityRegistryV1ConfigFactoryFn() cre.CapabilityRegistryConfigFactoryFn {
	return vaultregistry.CapabilityRegistryConfigFn
}
