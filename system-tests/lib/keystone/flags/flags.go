package flags

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/system-tests/lib/keystone/types"
)

func DONTopologyWithFlag(donTopologies []*types.DONTopology, flag string) []*types.DONTopology {
	var result []*types.DONTopology

	for _, donTopology := range donTopologies {
		if HasFlag(donTopology.Flags, flag) {
			result = append(result, donTopology)
		}
	}

	return result
}

func HasFlag(values []string, flag string) bool {
	return slices.Contains(values, flag)
}

func MustOneDONTopologyWithFlag(t *testing.T, donTopologies []*types.DONTopology, flag string) *types.DONTopology {
	donTopologies = DONTopologyWithFlag(donTopologies, flag)
	require.Len(t, donTopologies, 1, "expected exactly one DON topology with flag %d", flag)

	return donTopologies[0]
}

func NodeSetFlags(nodeSet *types.CapabilitiesAwareNodeSet) ([]string, error) {
	var stringCaps []string
	if len(nodeSet.Capabilities) == 0 && nodeSet.DONType == "" {
		// if no flags are set, we assign all known capabilities to the DON
		return types.SingleDonFlags, nil
	}

	stringCaps = append(stringCaps, append(nodeSet.Capabilities, nodeSet.DONType)...)
	return stringCaps, nil
}
