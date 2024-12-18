package executable_test

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mr-tron/base58"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/executable"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/transmission"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

func Test_RemoteExecutionCapability_DonTopologies(t *testing.T) {
	ctx := testutils.Context(t)

	responseTest := func(t *testing.T, response commoncap.CapabilityResponse, responseError error) {
		require.NoError(t, responseError)
		mp, err := response.Value.Unwrap()
		require.NoError(t, err)
		assert.Equal(t, "aValue1", mp.(map[string]any)["response"].(string))
	}

	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_OneAtATime,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)

	timeOut := 10 * time.Minute

	capability := &TestCapability{}

	var methods []func(ctx context.Context, caller commoncap.ExecutableCapability)

	methods = append(methods, func(ctx context.Context, caller commoncap.ExecutableCapability) {
		executeCapability(ctx, t, caller, transmissionSchedule, responseTest)
	})

	capabilityFactory := func() commoncap.ExecutableCapability { return capability }
	remoteExecutableConfig := &commoncap.RemoteExecutableConfig{}

	for _, method := range methods {
		// Test scenarios where the number of submissions is greater than or equal to F + 1
		wfDon, _ := setupDons(t, remoteExecutableConfig, capabilityFactory, 1, 0, timeOut, 1, 0, timeOut)
		wfDon.ExecuteMethodInParallelOnAllNodes(ctx, method)

		wfDon, _ = setupDons(t, remoteExecutableConfig, capabilityFactory, 4, 3, timeOut, 1, 0, timeOut)
		wfDon.ExecuteMethodInParallelOnAllNodes(ctx, method)
		wfDon, _ = setupDons(t, remoteExecutableConfig, capabilityFactory, 10, 3, timeOut, 1, 0, timeOut)
		wfDon.ExecuteMethodInParallelOnAllNodes(ctx, method)

		wfDon, _ = setupDons(t, remoteExecutableConfig, capabilityFactory, 1, 0, timeOut, 1, 0, timeOut)
		wfDon.ExecuteMethodInParallelOnAllNodes(ctx, method)
		wfDon, _ = setupDons(t, remoteExecutableConfig, capabilityFactory, 1, 0, timeOut, 4, 3, timeOut)
		wfDon.ExecuteMethodInParallelOnAllNodes(ctx, method)
		wfDon, _ = setupDons(t, remoteExecutableConfig, capabilityFactory, 1, 0, timeOut, 10, 3, timeOut)
		wfDon.ExecuteMethodInParallelOnAllNodes(ctx, method)

		wfDon, _ = setupDons(t, remoteExecutableConfig, capabilityFactory, 4, 3, timeOut, 4, 3, timeOut)
		wfDon.ExecuteMethodInParallelOnAllNodes(ctx, method)
		wfDon, _ = setupDons(t, remoteExecutableConfig, capabilityFactory, 10, 3, timeOut, 10, 3, timeOut)
		wfDon.ExecuteMethodInParallelOnAllNodes(ctx, method)
		wfDon, _ = setupDons(t, remoteExecutableConfig, capabilityFactory, 10, 9, timeOut, 10, 9, timeOut)
		wfDon.ExecuteMethodInParallelOnAllNodes(ctx, method)
	}
}

func Test_RemoteExecutionCapability_RegisterAndUnregisterWorkflow(t *testing.T) {
	ctx := testutils.Context(t)

	timeOut := 10 * time.Minute

	remoteExecutableConfig := &commoncap.RemoteExecutableConfig{
		RequestHashExcludedAttributes: []string{},
		RegistrationRefresh:           100 * time.Millisecond,
		RegistrationExpiry:            1 * time.Second,
	}

	var serverSideCapabilities []commoncap.ExecutableCapability

	wfDon, _ := setupDons(t, remoteExecutableConfig, func() commoncap.ExecutableCapability {
		testCapability := &TestCapability{}
		serverSideCapabilities = append(serverSideCapabilities, testCapability)
		return testCapability
	}, 4, 1, timeOut, 4, 1, timeOut)

	registerRequest := commoncap.RegisterToWorkflowRequest{
		Metadata: commoncap.RegistrationMetadata{
			WorkflowID:    workflowID1,
			ReferenceID:   stepReferenceID1,
			WorkflowOwner: workflowOwnerID,
		},
	}

	unregisterRequest := commoncap.UnregisterFromWorkflowRequest{
		Metadata: commoncap.RegistrationMetadata{
			WorkflowID:    workflowID1,
			ReferenceID:   stepReferenceID1,
			WorkflowOwner: workflowOwnerID,
		},
	}

	workflowNodes := wfDon.GetNodes()

	// Call RegisterToWorkflow on 2 clients
	err := workflowNodes[0].RegisterToWorkflow(ctx, registerRequest)
	require.NoError(t, err)
	err = workflowNodes[1].RegisterToWorkflow(ctx, registerRequest)
	require.NoError(t, err)

	// wait a couple of refresh intervals
	time.Sleep(2 * remoteExecutableConfig.RegistrationRefresh)

	// Should have no registrations on any server side capabilities
	for _, capability := range serverSideCapabilities {
		assert.Empty(t, capability.(*TestCapability).GetRegisterRequests())
	}

	// Subscribe the remaining 2 clients to the same workflow
	err = workflowNodes[2].RegisterToWorkflow(ctx, registerRequest)
	require.NoError(t, err)
	err = workflowNodes[3].RegisterToWorkflow(ctx, registerRequest)
	require.NoError(t, err)

	// Should eventually have 1 registration on each server side capability
	for _, capability := range serverSideCapabilities {
		require.Eventually(t, func() bool {
			return len(capability.(*TestCapability).GetRegisterRequests()) == 1
		}, 5*time.Second, 100*time.Millisecond)
	}

	// Unregister a client (leaving f+1 clients registered)
	err = workflowNodes[0].UnregisterFromWorkflow(ctx, unregisterRequest)
	require.NoError(t, err)

	// wait a couple of expiry intervals
	time.Sleep(2 * remoteExecutableConfig.RegistrationExpiry)

	// Should have no unregistration requests on any server side capabilities
	for _, capability := range serverSideCapabilities {
		assert.Empty(t, capability.(*TestCapability).GetUnregisterRequests())
	}

	// Unregister another client (leaving less than f+1 clients registered)
	err = workflowNodes[1].UnregisterFromWorkflow(ctx, unregisterRequest)
	require.NoError(t, err)

	// Should eventually have 1 unregistration on each server side capability
	for _, capability := range serverSideCapabilities {
		require.Eventually(t, func() bool {
			return len(capability.(*TestCapability).GetUnregisterRequests()) == 1
		}, 5*time.Second, 100*time.Millisecond)
	}

	// Unregister the remaining clients
	err = workflowNodes[2].UnregisterFromWorkflow(ctx, unregisterRequest)
	require.NoError(t, err)
	err = workflowNodes[3].UnregisterFromWorkflow(ctx, unregisterRequest)
	require.NoError(t, err)

	// wait a couple of expiry intervals
	time.Sleep(2 * remoteExecutableConfig.RegistrationExpiry)

	// confirm there is still only 1 unregister request on each server side capability
	for _, capability := range serverSideCapabilities {
		assert.Len(t, capability.(*TestCapability).GetUnregisterRequests(), 1)
	}

	// re-register all the clients
	for i := 0; i < len(workflowNodes); i++ {
		err = workflowNodes[i].RegisterToWorkflow(ctx, registerRequest)
		require.NoError(t, err)
	}

	// Should eventually have 2 registration requests on each server side capability
	for _, capability := range serverSideCapabilities {
		require.Eventually(t, func() bool {
			return len(capability.(*TestCapability).GetRegisterRequests()) == 2
		}, 5*time.Second, 100*time.Millisecond)
	}
}

func Test_RemoteExecutionCapability_RegisterAndUnregister_CapabilityNodeRestart(t *testing.T) {
	ctx := testutils.Context(t)

	timeOut := 10 * time.Minute

	remoteExecutableConfig := &commoncap.RemoteExecutableConfig{
		RequestHashExcludedAttributes: []string{},
		RegistrationRefresh:           100 * time.Millisecond,
		RegistrationExpiry:            1 * time.Second,
	}

	wfDon, capDon := setupDons(t, remoteExecutableConfig, func() commoncap.ExecutableCapability {
		testCapability := &TestCapability{}
		return testCapability
	}, 4, 1, timeOut, 4, 1, timeOut)

	registerRequest := commoncap.RegisterToWorkflowRequest{
		Metadata: commoncap.RegistrationMetadata{
			WorkflowID:    workflowID1,
			ReferenceID:   stepReferenceID1,
			WorkflowOwner: workflowOwnerID,
		},
	}

	workflowNodes := wfDon.GetNodes()

	// Call RegisterToWorkflow on all clients
	err := workflowNodes[0].RegisterToWorkflow(ctx, registerRequest)
	require.NoError(t, err)
	err = workflowNodes[1].RegisterToWorkflow(ctx, registerRequest)
	require.NoError(t, err)
	err = workflowNodes[2].RegisterToWorkflow(ctx, registerRequest)
	require.NoError(t, err)
	err = workflowNodes[3].RegisterToWorkflow(ctx, registerRequest)
	require.NoError(t, err)

	// Should eventually have 1 registration on each server side capability
	for _, node := range capDon.GetNodes() {
		require.Eventually(t, func() bool {
			return len(node.GetUnderlyingCapability().(*TestCapability).GetRegisterRequests()) == 1
		}, 5*time.Second, 100*time.Millisecond)
	}

	// Stop a single capability node
	capNodes := capDon.GetNodes()
	err = capNodes[0].Close()
	require.NoError(t, err)

	// Verify still have 1 registration on each server side capability
	for _, node := range capNodes {
		require.Eventually(t, func() bool {
			return len(node.GetUnderlyingCapability().(*TestCapability).GetRegisterRequests()) == 1
		}, 5*time.Second, 100*time.Millisecond)
	}

	// Restart the stopped capability node
	err = capNodes[0].Start(ctx)
	require.NoError(t, err)

	// The restarted nodes capability should eventually have 2 registrations, the latter one corresponding to the re-registration after restart
	require.Eventually(t, func() bool {
		return len(capNodes[0].GetUnderlyingCapability().(*TestCapability).GetRegisterRequests()) == 2
	}, 5*time.Second, 100*time.Millisecond)
}

func Test_RemoteExecutionCapability_RegisterAndUnregister_WorkflowNodeRestart(t *testing.T) {
	ctx := testutils.Context(t)

	timeOut := 10 * time.Minute

	remoteExecutableConfig := &commoncap.RemoteExecutableConfig{
		RequestHashExcludedAttributes: []string{},
		RegistrationRefresh:           100 * time.Millisecond,
		RegistrationExpiry:            1 * time.Second,
	}

	wfDon, capDon := setupDons(t, remoteExecutableConfig, func() commoncap.ExecutableCapability {
		testCapability := &TestCapability{}
		return testCapability
	}, 4, 1, timeOut, 4, 1, timeOut)

	registerRequest := commoncap.RegisterToWorkflowRequest{
		Metadata: commoncap.RegistrationMetadata{
			WorkflowID:    workflowID1,
			ReferenceID:   stepReferenceID1,
			WorkflowOwner: workflowOwnerID,
		},
	}

	workflowNodes := wfDon.GetNodes()

	// Call RegisterToWorkflow on all clients
	err := workflowNodes[0].RegisterToWorkflow(ctx, registerRequest)
	require.NoError(t, err)
	err = workflowNodes[1].RegisterToWorkflow(ctx, registerRequest)
	require.NoError(t, err)
	err = workflowNodes[2].RegisterToWorkflow(ctx, registerRequest)
	require.NoError(t, err)
	err = workflowNodes[3].RegisterToWorkflow(ctx, registerRequest)
	require.NoError(t, err)

	// Should eventually have 1 registration on each server side capability
	for _, node := range capDon.GetNodes() {
		require.Eventually(t, func() bool {
			return len(node.GetUnderlyingCapability().(*TestCapability).GetRegisterRequests()) == 1
		}, 5*time.Second, 100*time.Millisecond)
	}

	// Stop a single workflow node
	wfNodes := wfDon.GetNodes()
	err = wfNodes[0].Close()
	require.NoError(t, err)

	// sleep for a couple of registration expiry intervals
	time.Sleep(2 * remoteExecutableConfig.RegistrationExpiry)

	// Verify no unregister requests on any capability nodes
	for _, node := range capDon.GetNodes() {
		require.Eventually(t, func() bool {
			return len(node.GetUnderlyingCapability().(*TestCapability).GetUnregisterRequests()) == 0
		}, 5*time.Second, 100*time.Millisecond)
	}

	// Restart the stopped workflow node
	err = wfNodes[0].Start(ctx)
	require.NoError(t, err)

	// sleep for a couple of refresh intervals
	time.Sleep(2 * remoteExecutableConfig.RegistrationRefresh)

	// Verify still have 1 registration on each server side capability
	for _, node := range capDon.GetNodes() {
		require.Eventually(t, func() bool {
			return len(node.GetUnderlyingCapability().(*TestCapability).GetRegisterRequests()) == 1
		}, 5*time.Second, 100*time.Millisecond)
	}

	// Stop 2 of the workflow nodes
	err = wfNodes[1].Close()
	require.NoError(t, err)
	err = wfNodes[2].Close()
	require.NoError(t, err)

	// Eventually all capability nodes should have 1 unregister request
	for _, node := range capDon.GetNodes() {
		require.Eventually(t, func() bool {
			return len(node.GetUnderlyingCapability().(*TestCapability).GetUnregisterRequests()) == 1
		}, 5*time.Second, 100*time.Millisecond)
	}

	// Restart the stopped workflow nodes and register to workflow
	err = wfNodes[1].Start(ctx)
	require.NoError(t, err)
	err = wfNodes[2].Start(ctx)
	require.NoError(t, err)

	err = workflowNodes[1].RegisterToWorkflow(ctx, registerRequest)
	require.NoError(t, err)
	err = workflowNodes[2].RegisterToWorkflow(ctx, registerRequest)
	require.NoError(t, err)

	// Eventually all capability nodes show have 2 register requests, the latter one corresponding to the re-registration
	// after restart of the workflow nodes
	for _, node := range capDon.GetNodes() {
		require.Eventually(t, func() bool {
			return len(node.GetUnderlyingCapability().(*TestCapability).GetRegisterRequests()) == 2
		}, 5*time.Second, 100*time.Millisecond)
	}
}

func Test_RemoteExecutableCapability_TransmissionSchedules(t *testing.T) {
	ctx := testutils.Context(t)

	responseTest := func(t *testing.T, response commoncap.CapabilityResponse, responseError error) {
		require.NoError(t, responseError)
		mp, err := response.Value.Unwrap()
		require.NoError(t, err)
		assert.Equal(t, "aValue1", mp.(map[string]any)["response"].(string))
	}

	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_OneAtATime,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)

	timeOut := 10 * time.Minute

	capability := &TestCapability{}

	method := func(ctx context.Context, caller commoncap.ExecutableCapability) {
		executeCapability(ctx, t, caller, transmissionSchedule, responseTest)
	}
	wfDon, _ := setupDons(t, &commoncap.RemoteExecutableConfig{}, func() commoncap.ExecutableCapability { return capability }, 10, 9, timeOut, 10, 9, timeOut)
	wfDon.ExecuteMethodInParallelOnAllNodes(ctx, method)

	transmissionSchedule, err = values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_AllAtOnce,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)
	method = func(ctx context.Context, caller commoncap.ExecutableCapability) {
		executeCapability(ctx, t, caller, transmissionSchedule, responseTest)
	}

	wfDon, _ = setupDons(t, &commoncap.RemoteExecutableConfig{}, func() commoncap.ExecutableCapability { return capability }, 10, 9, timeOut, 10, 9, timeOut)
	wfDon.ExecuteMethodInParallelOnAllNodes(ctx, method)
}

func Test_RemoteExecutionCapability_CapabilityError(t *testing.T) {
	ctx := testutils.Context(t)

	capability := &TestErrorCapability{}

	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_AllAtOnce,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)

	var methods []func(ctx context.Context, caller commoncap.ExecutableCapability)

	methods = append(methods, func(ctx context.Context, caller commoncap.ExecutableCapability) {
		executeCapability(ctx, t, caller, transmissionSchedule, func(t *testing.T, responseCh commoncap.CapabilityResponse, responseError error) {
			assert.Equal(t, "error executing request: failed to execute capability", responseError.Error())
		})
	})

	for _, method := range methods {
		wfDon, _ := setupDons(t, &commoncap.RemoteExecutableConfig{}, func() commoncap.ExecutableCapability { return capability }, 10, 9, 10*time.Minute, 10, 9, 10*time.Minute)
		wfDon.ExecuteMethodInParallelOnAllNodes(ctx, method)
	}
}

func Test_RemoteExecutableCapability_RandomCapabilityError(t *testing.T) {
	ctx := testutils.Context(t)

	capability := &TestRandomErrorCapability{}

	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_AllAtOnce,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)

	var methods []func(ctx context.Context, caller commoncap.ExecutableCapability)

	methods = append(methods, func(ctx context.Context, caller commoncap.ExecutableCapability) {
		executeCapability(ctx, t, caller, transmissionSchedule, func(t *testing.T, responseCh commoncap.CapabilityResponse, responseError error) {
			assert.Equal(t, "error executing request: failed to execute capability", responseError.Error())
		})
	})

	for _, method := range methods {
		wfDon, _ := setupDons(t, &commoncap.RemoteExecutableConfig{}, func() commoncap.ExecutableCapability { return capability }, 10, 9, 1*time.Second, 10, 9, 10*time.Minute)
		wfDon.ExecuteMethodInParallelOnAllNodes(ctx, method)
	}
}

func setupDons(t *testing.T,
	remoteExecutableConfig *commoncap.RemoteExecutableConfig,
	capabilityFactory func() commoncap.ExecutableCapability, numWorkflowPeers int, workflowDonF uint8, workflowNodeTimeout time.Duration,
	numCapabilityPeers int, capabilityDonF uint8, capabilityNodeResponseTimeout time.Duration) (*workflowDon, *capabilityDon) {
	lggr := logger.TestLogger(t)

	capabilityPeers := make([]p2ptypes.PeerID, numCapabilityPeers)
	for i := 0; i < numCapabilityPeers; i++ {
		capabilityPeerID := p2ptypes.PeerID{}
		require.NoError(t, capabilityPeerID.UnmarshalText([]byte(NewPeerID())))
		capabilityPeers[i] = capabilityPeerID
	}

	capabilityPeerID := p2ptypes.PeerID{}
	require.NoError(t, capabilityPeerID.UnmarshalText([]byte(NewPeerID())))

	capDonInfo := commoncap.DON{
		ID:      2,
		Members: capabilityPeers,
		F:       capabilityDonF,
	}

	capInfo := commoncap.CapabilityInfo{
		ID:             "cap_id@1.0.0",
		CapabilityType: commoncap.CapabilityTypeTarget,
		Description:    "Remote Target",
		DON:            &capDonInfo,
	}

	workflowPeers := make([]p2ptypes.PeerID, numWorkflowPeers)
	for i := 0; i < numWorkflowPeers; i++ {
		workflowPeerID := p2ptypes.PeerID{}
		require.NoError(t, workflowPeerID.UnmarshalText([]byte(NewPeerID())))
		workflowPeers[i] = workflowPeerID
	}

	workflowDonInfo := commoncap.DON{
		Members: workflowPeers,
		ID:      1,
		F:       workflowDonF,
	}

	broker := newTestAsyncMessageBroker(t, 1000)
	servicetest.Run(t, broker)

	workflowDONs := map[uint32]commoncap.DON{
		workflowDonInfo.ID: workflowDonInfo,
	}

	capDon := newCapabilityDon()
	for i := 0; i < numCapabilityPeers; i++ {
		node := newServerNode(capabilityPeers[i], broker, remoteExecutableConfig, capabilityFactory(), capInfo, capDonInfo, workflowDONs, capabilityNodeResponseTimeout, lggr)
		capDon.AddNode(node)
		servicetest.Run(t, node)
	}

	wfDon := newWorkflowDon(broker)
	for i := 0; i < numWorkflowPeers; i++ {
		node := newClientNode(workflowPeers[i], broker, remoteExecutableConfig, capInfo, capDonInfo, workflowDonInfo, workflowNodeTimeout, lggr)
		wfDon.AddNode(node)
		servicetest.Run(t, node)
	}

	return wfDon, capDon
}

type workflowNode interface {
	commoncap.ExecutableCapability
	Start(ctx context.Context) error
	Close() error
}

type client interface {
	remotetypes.Receiver
	commoncap.ExecutableCapability
	Start(ctx context.Context) error
	Close() error
}

type clientNode struct {
	client     client
	nodePeerID p2ptypes.PeerID
	broker     *testAsyncMessageBroker

	remoteExecutableConfig *commoncap.RemoteExecutableConfig
	remoteCapabilityInfo   commoncap.CapabilityInfo
	remoteDonInfo          commoncap.DON
	localDonInfo           commoncap.DON
	requestTimeout         time.Duration
	mux                    sync.Mutex
	running                bool
	lggr                   logger.Logger
}

func newClientNode(nodePeerID p2ptypes.PeerID, broker *testAsyncMessageBroker, remoteExecutableConfig *commoncap.RemoteExecutableConfig,
	remoteCapabilityInfo commoncap.CapabilityInfo,
	remoteDonInfo commoncap.DON,
	localDonInfo commoncap.DON,
	requestTimeout time.Duration,
	lggr logger.Logger) *clientNode {
	return &clientNode{
		nodePeerID:             nodePeerID,
		broker:                 broker,
		remoteExecutableConfig: remoteExecutableConfig,
		remoteCapabilityInfo:   remoteCapabilityInfo,
		remoteDonInfo:          remoteDonInfo,
		localDonInfo:           localDonInfo,
		requestTimeout:         requestTimeout,
		lggr:                   lggr,
	}
}

func (w *clientNode) Start(ctx context.Context) error {
	w.mux.Lock()
	defer w.mux.Unlock()
	if !w.running {
		w.client = executable.NewClient(w.remoteExecutableConfig, w.remoteCapabilityInfo, w.remoteDonInfo, w.localDonInfo, w.broker.NewDispatcherForNode(w.nodePeerID), w.requestTimeout, w.lggr)
		w.broker.RegisterReceiverNode(w.nodePeerID, w.client)
		if err := w.client.Start(ctx); err != nil {
			return fmt.Errorf("failed to start client: %w", err)
		}
		w.running = true
	}

	return nil
}

func (w *clientNode) Close() error {
	w.mux.Lock()
	defer w.mux.Unlock()
	if w.running {
		w.broker.RemoveReceiverNode(w.nodePeerID)
		if err := w.client.Close(); err != nil {
			return fmt.Errorf("failed to close client: %w", err)
		}
		w.running = false
	}

	return nil
}

func (w *clientNode) Execute(ctx context.Context, request commoncap.CapabilityRequest) (commoncap.CapabilityResponse, error) {
	return w.client.Execute(ctx, request)
}

func (w *clientNode) RegisterToWorkflow(ctx context.Context, request commoncap.RegisterToWorkflowRequest) error {
	return w.client.RegisterToWorkflow(ctx, request)
}

func (w *clientNode) UnregisterFromWorkflow(ctx context.Context, request commoncap.UnregisterFromWorkflowRequest) error {
	return w.client.UnregisterFromWorkflow(ctx, request)
}

func (w *clientNode) Info(ctx context.Context) (commoncap.CapabilityInfo, error) {
	return w.client.Info(ctx)
}

type workflowDon struct {
	nodes []workflowNode
}

func newWorkflowDon(broker *testAsyncMessageBroker) *workflowDon {
	return &workflowDon{
		nodes: make([]workflowNode, 0),
	}
}

func (w *workflowDon) ExecuteMethodInParallelOnAllNodes(ctx context.Context, method func(ctx context.Context, caller commoncap.ExecutableCapability)) {
	wg := &sync.WaitGroup{}
	wg.Add(len(w.nodes))

	for _, node := range w.nodes {
		go func(caller commoncap.ExecutableCapability) {
			defer wg.Done()
			method(ctx, caller)
		}(node)
	}

	wg.Wait()
}

func (w *workflowDon) AddNode(wfNode workflowNode) {
	w.nodes = append(w.nodes, wfNode)
}

func (w *workflowDon) GetNodes() []workflowNode {
	return w.nodes
}

type server interface {
	remotetypes.Receiver
	Start(ctx context.Context) error
	Close() error
}

type serverNode struct {
	server     server
	nodePeerID p2ptypes.PeerID
	broker     *testAsyncMessageBroker

	remoteExecutableConfig *commoncap.RemoteExecutableConfig
	underlying             commoncap.ExecutableCapability
	capInfo                commoncap.CapabilityInfo
	localDonInfo           commoncap.DON
	workflowDONs           map[uint32]commoncap.DON
	requestTimeout         time.Duration
	mux                    sync.Mutex
	running                bool
	lggr                   logger.Logger
}

func newServerNode(nodePeerID p2ptypes.PeerID, broker *testAsyncMessageBroker, remoteExecutableConfig *commoncap.RemoteExecutableConfig,
	underlying commoncap.ExecutableCapability,
	capInfo commoncap.CapabilityInfo,
	localDonInfo commoncap.DON,
	workflowDONs map[uint32]commoncap.DON,
	requestTimeout time.Duration,
	lggr logger.Logger) *serverNode {
	return &serverNode{
		nodePeerID:             nodePeerID,
		broker:                 broker,
		remoteExecutableConfig: remoteExecutableConfig,
		underlying:             underlying,
		capInfo:                capInfo,
		localDonInfo:           localDonInfo,
		workflowDONs:           workflowDONs,
		requestTimeout:         requestTimeout,
		lggr:                   lggr,
	}
}

func (w *serverNode) GetUnderlyingCapability() commoncap.ExecutableCapability {
	return w.underlying
}

func (w *serverNode) Start(ctx context.Context) error {
	w.mux.Lock()
	defer w.mux.Unlock()
	if !w.running {
		w.server = executable.NewServer(w.remoteExecutableConfig, w.nodePeerID, w.underlying, w.capInfo, w.localDonInfo, w.workflowDONs, w.broker.NewDispatcherForNode(w.nodePeerID),
			w.requestTimeout, w.lggr)
		w.broker.RegisterReceiverNode(w.nodePeerID, w.server)
		if err := w.server.Start(ctx); err != nil {
			return fmt.Errorf("failed to start server: %w", err)
		}
		w.running = true
	}
	return nil
}

func (w *serverNode) Close() error {
	w.mux.Lock()
	defer w.mux.Unlock()
	if w.running {
		w.broker.RemoveReceiverNode(w.nodePeerID)
		if err := w.server.Close(); err != nil {
			return fmt.Errorf("failed to close server: %w", err)
		}
		w.running = false
	}

	return nil
}

type capabilityNode interface {
	Start(ctx context.Context) error
	Close() error
	GetUnderlyingCapability() commoncap.ExecutableCapability
}

type capabilityDon struct {
	nodes []capabilityNode
}

func newCapabilityDon() *capabilityDon {
	return &capabilityDon{
		nodes: make([]capabilityNode, 0),
	}
}

func (c *capabilityDon) AddNode(node capabilityNode) {
	c.nodes = append(c.nodes, node)
}

func (c *capabilityDon) GetNodes() []capabilityNode {
	return c.nodes
}

type testAsyncMessageBroker struct {
	services.Service
	eng *services.Engine
	t   *testing.T

	mux    sync.Mutex
	nodes  map[p2ptypes.PeerID]remotetypes.Receiver
	sendCh chan *remotetypes.MessageBody
}

func newTestAsyncMessageBroker(t *testing.T, sendChBufferSize int) *testAsyncMessageBroker {
	b := &testAsyncMessageBroker{
		t:      t,
		nodes:  make(map[p2ptypes.PeerID]remotetypes.Receiver),
		sendCh: make(chan *remotetypes.MessageBody, sendChBufferSize),
	}
	b.Service, b.eng = services.Config{
		Name:  "testAsyncMessageBroker",
		Start: b.start,
	}.NewServiceEngine(logger.TestLogger(t))
	return b
}

func (a *testAsyncMessageBroker) start(ctx context.Context) error {
	a.eng.Go(func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-a.sendCh:
				receiverID := toPeerID(msg.Receiver)

				var receiver remotetypes.Receiver
				a.mux.Lock()
				receiver = a.nodes[receiverID]
				a.mux.Unlock()

				if receiver != nil {
					receiver.Receive(tests.Context(a.t), msg)
				}
			}
		}
	})
	return nil
}

func (a *testAsyncMessageBroker) NewDispatcherForNode(nodePeerID p2ptypes.PeerID) remotetypes.Dispatcher {
	return &nodeDispatcher{
		callerPeerID: nodePeerID,
		broker:       a,
	}
}

func (a *testAsyncMessageBroker) RegisterReceiverNode(nodePeerID p2ptypes.PeerID, node remotetypes.Receiver) {
	a.mux.Lock()
	defer a.mux.Unlock()
	if _, ok := a.nodes[nodePeerID]; ok {
		panic("node already registered")
	}

	a.nodes[nodePeerID] = node
}

func (a *testAsyncMessageBroker) RemoveReceiverNode(nodePeerID p2ptypes.PeerID) {
	a.mux.Lock()
	defer a.mux.Unlock()
	delete(a.nodes, nodePeerID)
}

func (a *testAsyncMessageBroker) Send(msg *remotetypes.MessageBody) {
	a.sendCh <- msg
}

func toPeerID(id []byte) p2ptypes.PeerID {
	return [32]byte(id)
}

type broker interface {
	Send(msg *remotetypes.MessageBody)
}

type nodeDispatcher struct {
	callerPeerID p2ptypes.PeerID
	broker       broker
}

func (t *nodeDispatcher) Name() string {
	return "nodeDispatcher"
}

func (t *nodeDispatcher) Start(ctx context.Context) error {
	return nil
}

func (t *nodeDispatcher) Close() error {
	return nil
}

func (t *nodeDispatcher) Ready() error {
	return nil
}

func (t *nodeDispatcher) HealthReport() map[string]error {
	return nil
}

func (t *nodeDispatcher) Send(peerID p2ptypes.PeerID, msgBody *remotetypes.MessageBody) error {
	msgBody.Version = 1
	msgBody.Sender = t.callerPeerID[:]
	msgBody.Receiver = peerID[:]
	msgBody.Timestamp = time.Now().UnixMilli()
	t.broker.Send(msgBody)
	return nil
}

func (t *nodeDispatcher) SetReceiver(capabilityID string, donID uint32, receiver remotetypes.Receiver) error {
	return nil
}
func (t *nodeDispatcher) RemoveReceiver(capabilityID string, donID uint32) {}

type abstractTestCapability struct {
}

func (t abstractTestCapability) Info(ctx context.Context) (commoncap.CapabilityInfo, error) {
	return commoncap.CapabilityInfo{}, nil
}

type TestCapability struct {
	abstractTestCapability
	registerRequests   []commoncap.RegisterToWorkflowRequest
	unregisterRequests []commoncap.UnregisterFromWorkflowRequest
	mu                 sync.Mutex
}

func (t *TestCapability) Execute(ctx context.Context, request commoncap.CapabilityRequest) (commoncap.CapabilityResponse, error) {
	value := request.Inputs.Underlying["executeValue1"]
	response, err := values.NewMap(map[string]any{"response": value})
	if err != nil {
		return commoncap.CapabilityResponse{}, err
	}
	return commoncap.CapabilityResponse{
		Value: response,
	}, nil
}

func (t *TestCapability) RegisterToWorkflow(ctx context.Context, request commoncap.RegisterToWorkflowRequest) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.registerRequests = append(t.registerRequests, request)
	return nil
}

func (t *TestCapability) UnregisterFromWorkflow(ctx context.Context, request commoncap.UnregisterFromWorkflowRequest) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.unregisterRequests = append(t.unregisterRequests, request)
	return nil
}

func (t *TestCapability) GetRegisterRequests() []commoncap.RegisterToWorkflowRequest {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.registerRequests
}

func (t *TestCapability) GetUnregisterRequests() []commoncap.UnregisterFromWorkflowRequest {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.unregisterRequests
}

type TestErrorCapability struct {
	abstractTestCapability
	registerRequests   []commoncap.RegisterToWorkflowRequest
	unregisterRequests []commoncap.UnregisterFromWorkflowRequest
	mu                 sync.Mutex
}

func (t *TestErrorCapability) Execute(ctx context.Context, request commoncap.CapabilityRequest) (commoncap.CapabilityResponse, error) {
	return commoncap.CapabilityResponse{}, errors.New("an error")
}

func (t *TestErrorCapability) RegisterToWorkflow(ctx context.Context, request commoncap.RegisterToWorkflowRequest) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.registerRequests = append(t.registerRequests, request)
	return nil
}

func (t *TestErrorCapability) UnregisterFromWorkflow(ctx context.Context, request commoncap.UnregisterFromWorkflowRequest) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.unregisterRequests = append(t.unregisterRequests, request)
	return nil
}

func (t *TestErrorCapability) GetRegisterRequests() []commoncap.RegisterToWorkflowRequest {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.registerRequests
}

func (t *TestErrorCapability) GetUnregisterRequests() []commoncap.UnregisterFromWorkflowRequest {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.unregisterRequests
}

type TestRandomErrorCapability struct {
	abstractTestCapability
}

func (t TestRandomErrorCapability) Execute(ctx context.Context, request commoncap.CapabilityRequest) (commoncap.CapabilityResponse, error) {
	return commoncap.CapabilityResponse{}, errors.New(uuid.New().String())
}

func (t TestRandomErrorCapability) RegisterToWorkflow(ctx context.Context, request commoncap.RegisterToWorkflowRequest) error {
	return errors.New(uuid.New().String())
}

func (t TestRandomErrorCapability) UnregisterFromWorkflow(ctx context.Context, request commoncap.UnregisterFromWorkflowRequest) error {
	return errors.New(uuid.New().String())
}

func NewP2PPeerID(t *testing.T) p2ptypes.PeerID {
	id := p2ptypes.PeerID{}
	require.NoError(t, id.UnmarshalText([]byte(NewPeerID())))
	return id
}

func NewPeerID() string {
	var privKey [32]byte
	_, err := rand.Read(privKey[:])
	if err != nil {
		panic(err)
	}

	peerID := append(libp2pMagic(), privKey[:]...)

	return base58.Encode(peerID)
}

func libp2pMagic() []byte {
	return []byte{0x00, 0x24, 0x08, 0x01, 0x12, 0x20}
}

func executeCapability(ctx context.Context, t *testing.T, caller commoncap.ExecutableCapability, transmissionSchedule *values.Map, responseTest func(t *testing.T, response commoncap.CapabilityResponse, responseError error)) {
	executeInputs, err := values.NewMap(
		map[string]any{
			"executeValue1": "aValue1",
		},
	)
	require.NoError(t, err)
	response, err := caller.Execute(ctx,
		commoncap.CapabilityRequest{
			Metadata: commoncap.RequestMetadata{
				WorkflowID:          workflowID1,
				WorkflowExecutionID: workflowExecutionID1,
			},
			Config: transmissionSchedule,
			Inputs: executeInputs,
		})

	responseTest(t, response, err)
}
