package registrysyncer

import (
	"testing"

	"github.com/smartcontractkit/libocr/ragep2p/types"
	"github.com/test-go/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

func TestLocalRegistry_DONForCapability(t *testing.T) {
	lggr := logger.Test(t)
	getPeerID := func() (types.PeerID, error) {
		return [32]byte{0: 1}, nil
	}
	idsToDons := map[DonID]DON{}
	idsToNodes := map[types.PeerID]NodeInfo{}
	idsToCapabilities := map[string]Capability{}
	lr := NewLocalRegistry(lggr, getPeerID, idsToDons, idsToNodes, idsToCapabilities)

	_, _, err := lr.DONForCapability(t.Context(), "capabilityID@1.0.0")
	require.NoError(t, err)
}
