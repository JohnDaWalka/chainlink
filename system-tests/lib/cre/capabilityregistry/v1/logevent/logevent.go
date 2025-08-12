package logevent

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

var CapabilityRegistryConfigFn = func(donFlags []string, nodeSetInput *cre.CapabilitiesAwareNodeSet) []keystone_changeset.DONCapabilityWithConfig {
	var capabilities []keystone_changeset.DONCapabilityWithConfig

	if nodeSetInput == nil || nodeSetInput.ChainCapabilities == nil {
		return capabilities
	}

	for _, chainID := range nodeSetInput.ChainCapabilities[cre.WriteEVMCapability].EnabledChains {
		if flags.HasFlag(donFlags, cre.LogTriggerCapability) {
			capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
				Capability: kcr.CapabilitiesRegistryCapability{
					LabelledName:   fmt.Sprintf("log-event-trigger-evm-%d", chainID),
					Version:        "1.0.0",
					CapabilityType: 0, // TRIGGER
					ResponseType:   0, // REPORT
				},
				Config: &capabilitiespb.CapabilityConfig{},
			})
		}
	}

	return capabilities
}
