package evm

import (
	"fmt"
	"strconv"

	chainselectors "github.com/smartcontractkit/chain-selectors"

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
		if flags.HasFlag(donFlags, cre.EVMCapability) {
			selector, err := chainselectors.SelectorFromChainId(chainID)
			if err != nil {
				fmt.Printf("Error getting selector from chainID: %d, err: %s\n", chainID, err.Error())
				selector = 0
			}

			capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
				Capability: kcr.CapabilitiesRegistryCapability{
					LabelledName:   "evm" + ":ChainSelector:" + strconv.FormatUint(selector, 10),
					Version:        "1.0.0",
					CapabilityType: 3, // TARGET
					ResponseType:   1, // OBSERVATION_IDENTICAL
				},
				Config: &capabilitiespb.CapabilityConfig{},
			})
		}
	}

	return capabilities
}
