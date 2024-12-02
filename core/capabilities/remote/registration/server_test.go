package registration

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"

	libocrtypes "github.com/smartcontractkit/libocr/ragep2p/types"
)

func TestServer_Register(t *testing.T) {
	lggr := logger.TestLogger(t)
	target := &mockTarget{}

	peer1 := libocrtypes.PeerID{'p', 'e', 'e', 'r', '1'}
	peer2 := libocrtypes.PeerID{'p', 'e', 'e', 'r', '2'}
	peer3 := libocrtypes.PeerID{'p', 'e', 'e', 'r', '3'}
	peer4 := libocrtypes.PeerID{'p', 'e', 'e', 'r', '4'}
	peer5 := libocrtypes.PeerID{'p', 'e', 'e', 'r', '5'}

	workflowID1 := "15c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0"

	capInfo := capabilities.CapabilityInfo{ID: "test-capability"}
	registrationExpiry := 60 * time.Second
	workflowDONs := make(map[uint32]capabilities.DON)
	workflowDONs[1] = capabilities.DON{F: 1, Members: []p2ptypes.PeerID{peer1, peer2, peer3, peer4, peer5}}

	srv := NewServer(lggr, target, capInfo, registrationExpiry, workflowDONs, "test-server")
	servicetest.Run(t, srv)

	msg := &types.MessageBody{CallerDonId: 1, Payload: []byte("test-payload")}
	err := srv.Register(context.Background(), msg, peer1, workflowID1, "step1")
	require.NoError(t, err)
	assert.Empty(t, target.GetRegisterRequests())

	err = srv.Register(context.Background(), msg, peer2, workflowID1, "step1")
	require.NoError(t, err)
	assert.Empty(t, target.GetRegisterRequests())

	err = srv.Register(context.Background(), msg, peer3, workflowID1, "step1")
	require.NoError(t, err)
	// 2F+1 requests have been sent so register on the target should be called
	assert.Len(t, target.GetRegisterRequests(), 1)

	// Sending more requests should not result in the target receiving more register calls
	err = srv.Register(context.Background(), msg, peer4, workflowID1, "step1")
	require.NoError(t, err)
	err = srv.Register(context.Background(), msg, peer5, workflowID1, "step1")
	require.NoError(t, err)

	assert.Len(t, target.registerRequests, 1)
}

func TestServer_Unregister(t *testing.T) {
	lggr := logger.TestLogger(t)
	target := &mockTarget{}

	peer1 := libocrtypes.PeerID{'p', 'e', 'e', 'r', '1'}
	peer2 := libocrtypes.PeerID{'p', 'e', 'e', 'r', '2'}
	peer3 := libocrtypes.PeerID{'p', 'e', 'e', 'r', '3'}

	workflowID1 := "15c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0"

	capInfo := capabilities.CapabilityInfo{ID: "test-capability"}
	registrationExpiry := 10 * time.Millisecond
	workflowDONs := make(map[uint32]capabilities.DON)
	workflowDONs[1] = capabilities.DON{F: 1, Members: []p2ptypes.PeerID{peer1, peer2, peer3}}

	srv := NewServer(lggr, target, capInfo, registrationExpiry, workflowDONs, "test-server")
	servicetest.Run(t, srv)

	msg := &types.MessageBody{CallerDonId: 1, Payload: []byte("test-payload")}
	err := srv.Register(context.Background(), msg, peer1, workflowID1, "step1")
	require.NoError(t, err)
	err = srv.Register(context.Background(), msg, peer2, workflowID1, "step1")
	require.NoError(t, err)
	err = srv.Register(context.Background(), msg, peer3, workflowID1, "step1")
	require.NoError(t, err)

	assert.Eventually(t, func() bool { return len(target.GetUnregisterRequests()) == 1 }, 100*time.Millisecond, 10*time.Millisecond)
}

type mockTarget struct {
	registerRequests   [][]byte
	unregisterRequests [][]byte
	mux                sync.Mutex
}

func (m *mockTarget) Register(ctx context.Context, key Key, registerRequest []byte) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.registerRequests = append(m.registerRequests, registerRequest)
	return nil
}

func (m *mockTarget) Unregister(ctx context.Context, registerRequest []byte) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.unregisterRequests = append(m.unregisterRequests, registerRequest)
	return nil
}

func (m *mockTarget) GetRegisterRequests() [][]byte {
	m.mux.Lock()
	defer m.mux.Unlock()
	return m.registerRequests
}

func (m *mockTarget) GetUnregisterRequests() [][]byte {
	m.mux.Lock()
	defer m.mux.Unlock()
	return m.unregisterRequests
}
