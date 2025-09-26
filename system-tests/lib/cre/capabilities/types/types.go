package types

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
)

type InstallableCapability interface {
	// Flag returns the unique identifier used in TOML configurations and internal references
	Flag() cre.CapabilityFlag

	// JobSpecFn returns a function that generates job specifications for this capability
	// based on the provided input configuration and topology. Most capabilities need this.
	// Exceptions include capabilities that are configured via the node config, like write-evm, aptos, tron or solana.
	JobSpecFn() jobs.JobSpecFn

	// NodeConfigTransformerFn returns a function to modify node-level configuration,
	// or nil if node config modification is not needed. Most capabilities don't need this.
	NodeConfigTransformerFn() cre.NodeConfigTransformerFn

	// GatewayJobHandlerConfigFn returns a function to configure gateway handlers in the gateway jobspec,
	// or nil if no gateway handler configuration is required for this capability. Only capabilities
	// that need to connect to external resources might need this.
	GatewayJobHandlerConfigFn() cre.GatewayHandlerConfigFn

	// CapabilityRegistryV1ConfigFn returns a function to generate capability registry
	// configuration for the v1 registry format
	CapabilityRegistryV1ConfigFn() cre.CapabilityRegistryConfigFn

	// CapabilityRegistryV2ConfigFn returns a function to generate capability registry
	// configuration for the v2 registry format
	CapabilityRegistryV2ConfigFn() cre.CapabilityRegistryConfigFn
}
