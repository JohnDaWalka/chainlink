package flags

import (
	"slices"
	"strings"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

func HasFlag(values []string, flag string) bool {
	return slices.Contains(values, flag)
}

func HasOnlyOneFlag(values []string, flag string) bool {
	return slices.Contains(values, flag) && len(values) == 1
}

// HasFlagForChain checks if a capability is enabled for a specific chain on a nodeset.
// Returns true if the capability is listed in global Capabilities, or if it's enabled
// in ChainCapabilities for the given chain ID.
func HasFlagForChain(nodeSet *cre.CapabilitiesAwareNodeSet, capability string, chainID uint64) bool {
	if nodeSet == nil {
		return false
	}
	if HasFlag(nodeSet.ComputedCapabilities, capability) {
		return true
	}
	if nodeSet.ChainCapabilities == nil {
		return false
	}
	cfg, ok := nodeSet.ChainCapabilities[capability]
	if !ok || cfg == nil {
		return false
	}
	return slices.Contains(cfg.EnabledChains, chainID)
}

func HasFlagForAnyChain(values []string, capability string) bool {
	if HasFlag(values, capability) {
		return true
	}

	for _, value := range values {
		if strings.HasPrefix(value, capability+"-") {
			return true
		}
	}

	return false
}

func DonMetadataWithFlag(donTopologies []*cre.DonMetadata, flag string) []*cre.DonMetadata {
	var result []*cre.DonMetadata

	for _, donTopology := range donTopologies {
		if HasFlagForAnyChain(donTopology.Flags, flag) {
			result = append(result, donTopology)
		}
	}

	return result
}

func OneDonMetadataWithFlag(donTopologies []*cre.DonMetadata, flag string) (*cre.DonMetadata, error) {
	donTopologies = DonMetadataWithFlag(donTopologies, flag)
	if len(donTopologies) != 1 {
		return nil, errors.Errorf("expected exactly one DON topology with flag %s, got %d", flag, len(donTopologies))
	}

	return donTopologies[0], nil
}

func NodeSetFlags(nodeSet *cre.CapabilitiesAwareNodeSet) ([]string, error) {
	var stringCaps []string

	return append(stringCaps, append(nodeSet.ComputedCapabilities, nodeSet.DONTypes...)...), nil
}
