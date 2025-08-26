package registrysyncer

import (
	"testing"

	"github.com/smartcontractkit/libocr/ragep2p/types"
	"github.com/test-go/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

func TestLocalRegistry_DONForCapability(t *testing.T) {
	lggr := logger.Test(t)
	getPeerID := func() (types.PeerID, error) {
		return [32]byte{0: 1}, nil
	}
	idsToDons := map[DonID]DON{
		1: {
			DON: capabilities.DON{
				Name: "don1",
				ID:   1,
				F:    1,
				Members: []types.PeerID{
					{0: 1},
					{0: 2},
				},
			},
			CapabilityConfigurations: map[string]CapabilityConfiguration{
				"capabilityID@1.0.0": CapabilityConfiguration{},
			},
		},
		2: {
			DON: capabilities.DON{
				Name: "don2",
				ID:   2,
				F:    2,
				Members: []types.PeerID{
					{0: 3},
					{0: 4},
				},
			},
			CapabilityConfigurations: map[string]CapabilityConfiguration{
				"secondCapabilityID@1.0.0": CapabilityConfiguration{},
			},
		},
	}
	idsToNodes := map[types.PeerID]NodeInfo{
		{0: 1}: NodeInfo{
			NodeOperatorID: 0,
		},
		{0: 2}: NodeInfo{
			NodeOperatorID: 1,
		},
		{0: 3}: NodeInfo{
			NodeOperatorID: 2,
		},
		{0: 4}: NodeInfo{
			NodeOperatorID: 3,
		},
	}
	idsToCapabilities := map[string]Capability{
		"capabilityID@1.0.0": {
			ID:             "capabilityID@1.0.0",
			CapabilityType: capabilities.CapabilityTypeAction,
		},
		"secondCapabilityID@1.0.0": {
			ID:             "secondCapabilityID@1.0.0",
			CapabilityType: capabilities.CapabilityTypeAction,
		},
	}
	lr := NewLocalRegistry(lggr, getPeerID, idsToDons, idsToNodes, idsToCapabilities)

	_, _, err := lr.DONForCapability(t.Context(), "capabilityID@1.0.0")
	require.NoError(t, err)
}
