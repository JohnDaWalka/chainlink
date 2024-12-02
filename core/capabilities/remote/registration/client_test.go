package registration

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	types2 "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"

	libocrtypes "github.com/smartcontractkit/libocr/ragep2p/types"
)

func TestClient_RegisterWorkflow(t *testing.T) {
	lggr := logger.TestLogger(t)
	dispatcher := NewMockDispatcher()
	peer1 := libocrtypes.PeerID{'p', 'e', 'e', 'r', '1'}
	peer2 := libocrtypes.PeerID{'p', 'e', 'e', 'r', '2'}

	client := NewClient(lggr, types.MethodRegisterTrigger, 10*time.Millisecond, capabilities.CapabilityInfo{}, capabilities.DON{Members: []libocrtypes.PeerID{peer1, peer2}}, capabilities.DON{}, dispatcher, "test")
	servicetest.Run(t, client)

	err := client.RegisterWorkflow("workflow1", []byte("registerrequest"))
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		messages := dispatcher.GetMessages()

		// Check sent to both peers with the same number of requests
		if len(messages[peer1]) >= 3 {
			return len(messages[peer1]) == len(messages[peer2])
		}
		return false
	}, 60*time.Second, 10*time.Millisecond)

	assert.Equal(t, "registerrequest", string(dispatcher.GetMessages()[peer1][0].Payload))
}

func TestClient_UnregisterWorkflow(t *testing.T) {
	lggr := logger.TestLogger(t)
	dispatcher := NewMockDispatcher()

	peer1 := libocrtypes.PeerID{'p', 'e', 'e', 'r', '1'}
	client := NewClient(lggr, types.MethodRegisterTrigger, 10*time.Millisecond, capabilities.CapabilityInfo{}, capabilities.DON{Members: []libocrtypes.PeerID{peer1}}, capabilities.DON{}, dispatcher, "test")
	servicetest.Run(t, client)

	err := client.RegisterWorkflow("workflow1", []byte("request"))
	require.NoError(t, err)

	client.UnregisterWorkflow("workflow1")

	initialCount := len(dispatcher.GetMessages()[peer1])
	time.Sleep(100 * time.Microsecond)
	// If it has been unregistered then no new registration requests should be sent
	finalCount := len(dispatcher.GetMessages()[peer1])
	assert.Equal(t, initialCount, finalCount)
}

type MockDispatcher struct {
	mu       sync.Mutex
	messages map[types2.PeerID][]*types.MessageBody
}

func NewMockDispatcher() *MockDispatcher {
	return &MockDispatcher{
		messages: make(map[types2.PeerID][]*types.MessageBody),
	}
}

func (m *MockDispatcher) Send(peerID types2.PeerID, msgBody *types.MessageBody) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.messages[peerID] = append(m.messages[peerID], msgBody)
	return nil
}

func (m *MockDispatcher) GetMessages() map[types2.PeerID][]*types.MessageBody {
	m.mu.Lock()
	defer m.mu.Unlock()

	mapCopy := make(map[types2.PeerID][]*types.MessageBody)
	for k, v := range m.messages {
		mapCopy[k] = append([]*types.MessageBody(nil), v...)
	}
	return mapCopy
}
