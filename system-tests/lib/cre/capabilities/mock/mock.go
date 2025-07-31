package mock

import (
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

var MockCapabilityFactoryFn = func(donFlags []string) []keystone_changeset.DONCapabilityWithConfig {
	var capabilities []keystone_changeset.DONCapabilityWithConfig

	capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
		Capability: kcr.CapabilitiesRegistryCapability{
			LabelledName:   "mock",
			Version:        "1.0.0",
			CapabilityType: 0, // TRIGGER
		},
		Config: &capabilitiespb.CapabilityConfig{},
	})

	return capabilities
}
