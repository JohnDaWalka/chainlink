package writesolana

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

var WriteSolanaCapabilityFactory = func(chainID uint64) func(donFlags []string) []keystone_changeset.DONCapabilityWithConfig {
	return func(donFlags []string) []keystone_changeset.DONCapabilityWithConfig {
		var capabilities []keystone_changeset.DONCapabilityWithConfig

		if flags.HasFlag(donFlags, types.WriteSolanaCapability) {
			capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
				Capability: kcr.CapabilitiesRegistryCapability{
					LabelledName:   "write-solana", // TODO: change this to the correct name
					Version:        "1.0.0",        // TODO: change this to the correct name
					CapabilityType: 3,              // TARGET
					ResponseType:   1,              // OBSERVATION_IDENTICAL
				},
				Config: &capabilitiespb.CapabilityConfig{},
			})
		}

		return capabilities
	}
}
