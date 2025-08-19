package readcontract

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

var CapabilityRegistryConfigFn = func(_ []string, nodeSetInput *cre.CapabilitiesAwareNodeSet) ([]keystone_changeset.DONCapabilityWithConfig, error) {
	var capabilities []keystone_changeset.DONCapabilityWithConfig

	if nodeSetInput == nil || nodeSetInput.ChainCapabilities == nil {
		return nil, errors.New("node set input is nil or chain capabilities is nil")
	}

	if _, ok := nodeSetInput.ChainCapabilities[cre.ReadContractCapability]; !ok {
		return nil, nil
	}

	for _, chainID := range nodeSetInput.ChainCapabilities[cre.ReadContractCapability].EnabledChains {
		capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   fmt.Sprintf("read-contract-evm-%d", chainID),
				Version:        "1.0.0",
				CapabilityType: 1, // ACTION
			},
			Config: &capabilitiespb.CapabilityConfig{},
		})
	}

	return capabilities, nil
}
