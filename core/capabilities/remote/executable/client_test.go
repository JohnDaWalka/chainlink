package executable_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-common/pkg/values"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/executable"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/registration"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/transmission"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

const (
	stepReferenceID1     = "step1"
	workflowID1          = "15c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0"
	workflowExecutionID1 = "95ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0abbadeed"
	workflowOwnerID      = "0xAA"
)

func Test_Client_DonTopologies(t *testing.T) {
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/CAPPL-363")
	ctx := testutils.Context(t)

	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_OneAtATime,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)

	responseTest := func(t *testing.T, response commoncap.CapabilityResponse, responseError error) {
		require.NoError(t, responseError)
		mp, err := response.Value.Unwrap()
		require.NoError(t, err)
		assert.Equal(t, "aValue1", mp.(map[string]any)["response"].(string))
	}

	capability := func() commoncap.ExecutableCapability { return &TestCapability{} }

	responseTimeOut := 10 * time.Minute

	method := func(caller commoncap.ExecutableCapability) {
		executeInputs, err := values.NewMap(map[string]any{"executeValue1": "aValue1"})
		require.NoError(t, err)
		executeMethod(ctx, caller, transmissionSchedule, executeInputs, responseTest, t)
	}

	remoteExecutableConfig := &commoncap.RemoteExecutableConfig{}

	testClient(t, remoteExecutableConfig, 1, responseTimeOut, 1, 0,
		capability, method)

	testClient(t, remoteExecutableConfig, 10, responseTimeOut, 1, 0,
		capability, method)

	testClient(t, remoteExecutableConfig, 1, responseTimeOut, 10, 3,
		capability, method)

	testClient(t, remoteExecutableConfig, 10, responseTimeOut, 10, 3,
		capability, method)

	testClient(t, remoteExecutableConfig, 10, responseTimeOut, 10, 9,
		capability, method)
}

func Test_Client_TransmissionSchedules(t *testing.T) {
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/CAPPL-363")
	ctx := testutils.Context(t)

	responseTest := func(t *testing.T, response commoncap.CapabilityResponse, responseError error) {
		require.NoError(t, responseError)
		mp, err := response.Value.Unwrap()
		require.NoError(t, err)
		assert.Equal(t, "aValue1", mp.(map[string]any)["response"].(string))
	}

	capability := &TestCapability{}

	responseTimeOut := 10 * time.Minute

	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_OneAtATime,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)

	testClient(t, &commoncap.RemoteExecutableConfig{}, 1, responseTimeOut, 1, 0,
		func() commoncap.ExecutableCapability { return capability }, func(caller commoncap.ExecutableCapability) {
			executeInputs, err2 := values.NewMap(map[string]any{"executeValue1": "aValue1"})
			require.NoError(t, err2)
			executeMethod(ctx, caller, transmissionSchedule, executeInputs, responseTest, t)
		})
	testClient(t, &commoncap.RemoteExecutableConfig{}, 10, responseTimeOut, 10, 3,
		func() commoncap.ExecutableCapability { return capability }, func(caller commoncap.ExecutableCapability) {
			executeInputs, err2 := values.NewMap(map[string]any{"executeValue1": "aValue1"})
			require.NoError(t, err2)
			executeMethod(ctx, caller, transmissionSchedule, executeInputs, responseTest, t)
		})

	transmissionSchedule, err = values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_AllAtOnce,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)

	testClient(t, &commoncap.RemoteExecutableConfig{}, 1, responseTimeOut, 1, 0,
		func() commoncap.ExecutableCapability { return capability }, func(caller commoncap.ExecutableCapability) {
			executeInputs, err := values.NewMap(map[string]any{"executeValue1": "aValue1"})
			require.NoError(t, err)
			executeMethod(ctx, caller, transmissionSchedule, executeInputs, responseTest, t)
		})
	testClient(t, &commoncap.RemoteExecutableConfig{}, 10, responseTimeOut, 10, 3,
		func() commoncap.ExecutableCapability { return capability }, func(caller commoncap.ExecutableCapability) {
			executeInputs, err := values.NewMap(map[string]any{"executeValue1": "aValue1"})
			require.NoError(t, err)
			executeMethod(ctx, caller, transmissionSchedule, executeInputs, responseTest, t)
		})
}

func Test_Client_TimesOutIfInsufficientCapabilityPeerResponses(t *testing.T) {
	ctx := testutils.Context(t)

	responseTest := func(t *testing.T, response commoncap.CapabilityResponse, responseError error) {
		require.Error(t, responseError)
		require.ErrorIs(t, responseError, executable.ErrRequestExpired)
	}

	capability := &TestCapability{}

	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_AllAtOnce,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)

	// number of capability peers is less than F + 1

	testClient(t, &commoncap.RemoteExecutableConfig{}, 10, 1*time.Second, 10, 11,
		func() commoncap.ExecutableCapability { return capability },
		func(caller commoncap.ExecutableCapability) {
			executeInputs, err := values.NewMap(map[string]any{"executeValue1": "aValue1"})
			require.NoError(t, err)
			executeMethod(ctx, caller, transmissionSchedule, executeInputs, responseTest, t)
		})
}

func Test_Client_ContextCanceledBeforeQuorumReached(t *testing.T) {
	ctx, cancel := context.WithCancel(testutils.Context(t))

	responseTest := func(t *testing.T, response commoncap.CapabilityResponse, responseError error) {
		require.Error(t, responseError)
		require.ErrorIs(t, responseError, executable.ErrContextDoneBeforeResponseQuorum)
	}

	capability := &TestCapability{}
	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_AllAtOnce,
		"deltaStage": "20s",
	})
	require.NoError(t, err)

	cancel()
	testClient(t, &commoncap.RemoteExecutableConfig{}, 2, 20*time.Second, 2, 2,
		func() commoncap.ExecutableCapability { return capability },
		func(caller commoncap.ExecutableCapability) {
			executeInputs, err := values.NewMap(map[string]any{"executeValue1": "aValue1"})
			require.NoError(t, err)
			executeMethod(ctx, caller, transmissionSchedule, executeInputs, responseTest, t)
		})
}

func Test_Client_RegisterAndUnregisterWorkflows(t *testing.T) {
	ctx := testutils.Context(t)

	responseTest := func(t *testing.T, responseError error) {
		require.NoError(t, responseError)
	}

	responseTimeOut := 10 * time.Minute

	clientSideCapabilities := make([]commoncap.ExecutableCapability, 0)
	mux := sync.Mutex{}
	method := func(capability commoncap.ExecutableCapability) {
		mux.Lock()
		defer mux.Unlock()
		registerToWorkflowMethod(ctx, capability, responseTest, t)
		clientSideCapabilities = append(clientSideCapabilities, capability)
	}

	remoteExecutableConfig := &commoncap.RemoteExecutableConfig{
		RegistrationRefresh: 100 * time.Millisecond,
		RegistrationExpiry:  1 * time.Second,
	}

	var serverSideCapabilities []*TestCapability
	testClient(t, remoteExecutableConfig, 4, responseTimeOut, 4, 1,
		func() commoncap.ExecutableCapability {
			capability := &TestCapability{}
			serverSideCapabilities = append(serverSideCapabilities, capability)
			return capability
		}, method)

	require.Eventually(t, func() bool {
		for _, testCapability := range serverSideCapabilities {
			if len(testCapability.GetRegisterRequests()) != 1 {
				return false
			}
		}

		return true
	}, 10*time.Second, 100*time.Millisecond, "expected a registration request to be received by all server side capabilities")

	// Wait a few multiple of the refresh interval and confirm that the capabilities have only 1 registration request and 0 unregister requests
	time.Sleep(remoteExecutableConfig.RegistrationRefresh * 4)

	for _, testCapability := range serverSideCapabilities {
		assert.Len(t, testCapability.GetRegisterRequests(), 1)
		assert.Empty(t, testCapability.GetUnregisterRequests())
	}

	// Unregister from workflow
	for _, capability := range clientSideCapabilities {
		unregisterFromWorkflowMethod(ctx, capability, responseTest, t)
	}

	require.Eventually(t, func() bool {
		for _, testCapability := range serverSideCapabilities {
			if len(testCapability.GetUnregisterRequests()) != 1 {
				return false
			}
		}
		return true
	}, 10*time.Second, 100*time.Millisecond, "expected a registration request to be received by all server side capabilities")

	// Wait a few multiple of the refresh interval and confirm that the capabilities have only 1 registration request and 1 unregister requests
	time.Sleep(remoteExecutableConfig.RegistrationRefresh * 4)

	for _, testCapability := range serverSideCapabilities {
		assert.Len(t, testCapability.GetRegisterRequests(), 1)
		assert.Len(t, testCapability.GetUnregisterRequests(), 1)
	}
}

func testClient(t *testing.T, remoteExecutableConfig *commoncap.RemoteExecutableConfig, numWorkflowPeers int, workflowNodeResponseTimeout time.Duration,
	numCapabilityPeers int, capabilityDonF uint8, capFactory func() commoncap.ExecutableCapability,
	method func(caller commoncap.ExecutableCapability)) []*clientTestServer {
	lggr := logger.TestLogger(t)
	remoteExecutableConfig.ApplyDefaults()

	capabilityPeers := make([]p2ptypes.PeerID, numCapabilityPeers)
	for i := 0; i < numCapabilityPeers; i++ {
		capabilityPeers[i] = NewP2PPeerID(t)
	}

	capDonInfo := commoncap.DON{
		ID:      1,
		Members: capabilityPeers,
		F:       capabilityDonF,
	}

	capInfo := commoncap.CapabilityInfo{
		ID:             "cap_id@1.0.0",
		CapabilityType: commoncap.CapabilityTypeAction,
		Description:    "Remote Executable Capability",
		DON:            &capDonInfo,
	}

	workflowPeers := make([]p2ptypes.PeerID, numWorkflowPeers)
	for i := 0; i < numWorkflowPeers; i++ {
		workflowPeers[i] = NewP2PPeerID(t)
	}

	workflowDonInfo := commoncap.DON{
		Members: workflowPeers,
		ID:      2,
	}

	broker := newTestAsyncMessageBroker(t, 100)

	testServers := make([]*clientTestServer, numCapabilityPeers)
	for i := 0; i < numCapabilityPeers; i++ {
		capabilityDispatcher := broker.NewDispatcherForNode(capabilityPeers[i])
		testServer := newTestServer(lggr, capabilityPeers[i], capInfo, remoteExecutableConfig.RegistrationExpiry, capabilityDispatcher, workflowDonInfo, capFactory())
		broker.RegisterReceiverNode(capabilityPeers[i], testServer)
		testServers[i] = testServer
		servicetest.Run(t, testServer)
	}

	clients := make([]commoncap.ExecutableCapability, numWorkflowPeers)

	for i := 0; i < numWorkflowPeers; i++ {
		workflowPeerDispatcher := broker.NewDispatcherForNode(workflowPeers[i])
		client := executable.NewClient(remoteExecutableConfig, capInfo, capDonInfo, workflowDonInfo, workflowPeerDispatcher,
			workflowNodeResponseTimeout, lggr)
		servicetest.Run(t, client)
		broker.RegisterReceiverNode(workflowPeers[i], client)
		clients[i] = client
	}

	servicetest.Run(t, broker)

	wg := &sync.WaitGroup{}
	wg.Add(len(clients))

	// Fire off all the requests
	for _, caller := range clients {
		go func(caller commoncap.ExecutableCapability) {
			defer wg.Done()
			method(caller)
		}(caller)
	}

	wg.Wait()

	return testServers
}

func registerToWorkflowMethod(ctx context.Context, caller commoncap.ExecutableCapability,
	responseTest func(t *testing.T, responseError error), t *testing.T) {
	err := caller.RegisterToWorkflow(ctx, commoncap.RegisterToWorkflowRequest{
		Metadata: commoncap.RegistrationMetadata{
			WorkflowID:    workflowID1,
			ReferenceID:   stepReferenceID1,
			WorkflowOwner: workflowOwnerID,
		},
	})

	responseTest(t, err)
}

func unregisterFromWorkflowMethod(ctx context.Context, caller commoncap.ExecutableCapability,
	responseTest func(t *testing.T, responseError error), t *testing.T) {
	err := caller.UnregisterFromWorkflow(ctx, commoncap.UnregisterFromWorkflowRequest{
		Metadata: commoncap.RegistrationMetadata{
			WorkflowID:    workflowID1,
			ReferenceID:   stepReferenceID1,
			WorkflowOwner: workflowOwnerID,
		},
	})

	responseTest(t, err)
}

func executeMethod(ctx context.Context, caller commoncap.ExecutableCapability, transmissionSchedule *values.Map,
	executeInputs *values.Map, responseTest func(t *testing.T, responseCh commoncap.CapabilityResponse, responseError error), t *testing.T) {
	responseCh, err := caller.Execute(ctx,
		commoncap.CapabilityRequest{
			Metadata: commoncap.RequestMetadata{
				WorkflowID:          workflowID1,
				WorkflowExecutionID: workflowExecutionID1,
				WorkflowOwner:       workflowOwnerID,
			},
			Config: transmissionSchedule,
			Inputs: executeInputs,
		})

	responseTest(t, responseCh, err)
}

// Simple client that only responds once it has received a message from each workflow peer
type clientTestServer struct {
	services.StateMachine
	lggr               logger.Logger
	peerID             p2ptypes.PeerID
	dispatcher         remotetypes.Dispatcher
	workflowDonInfo    commoncap.DON
	messageIDToSenders map[string]map[p2ptypes.PeerID]bool

	executableCapability commoncap.ExecutableCapability

	registrationServer *registration.Server

	mux sync.Mutex
}

func newTestServer(lggr logger.Logger, peerID p2ptypes.PeerID, capInfo commoncap.CapabilityInfo,
	registrationExpiry time.Duration, dispatcher remotetypes.Dispatcher, workflowDonInfo commoncap.DON,
	executableCapability commoncap.ExecutableCapability) *clientTestServer {
	target := &executable.TargetAdapter{Capability: executableCapability}

	workflowDONs := map[uint32]commoncap.DON{
		workflowDonInfo.ID: workflowDonInfo,
	}

	return &clientTestServer{
		lggr:                 lggr,
		dispatcher:           dispatcher,
		workflowDonInfo:      workflowDonInfo,
		peerID:               peerID,
		messageIDToSenders:   make(map[string]map[p2ptypes.PeerID]bool),
		executableCapability: executableCapability,
		registrationServer:   registration.NewServer(lggr, target, capInfo, registrationExpiry, workflowDONs, "testExecutableServer"),
	}
}

func (r *clientTestServer) Start(ctx context.Context) error {
	return r.StartOnce(r.peerID.String(), func() error {
		if err := r.registrationServer.Start(ctx); err != nil {
			return fmt.Errorf("failed to start capability register: %w", err)
		}
		return nil
	})
}

func (r *clientTestServer) Close() error {
	r.IfNotStopped(func() {
		if err := r.registrationServer.Close(); err != nil {
			r.lggr.Errorf("failed to close capability register: %v", err)
		}
	})

	return nil
}

func (r *clientTestServer) Receive(ctx context.Context, msg *remotetypes.MessageBody) {
	r.mux.Lock()
	defer r.mux.Unlock()

	sender := toPeerID(msg.Sender)

	switch msg.Method {
	case remotetypes.MethodExecute:
		messageID, err := executable.GetMessageID(msg)
		if err != nil {
			panic(err)
		}

		if r.messageIDToSenders[messageID] == nil {
			r.messageIDToSenders[messageID] = make(map[p2ptypes.PeerID]bool)
		}

		sendersOfMessageID := r.messageIDToSenders[messageID]
		if sendersOfMessageID[sender] {
			panic("received duplicate message")
		}

		sendersOfMessageID[sender] = true
		if len(r.messageIDToSenders[messageID]) == len(r.workflowDonInfo.Members) {
			capabilityRequest, err := pb.UnmarshalCapabilityRequest(msg.Payload)
			if err != nil {
				panic(err)
			}
			resp, responseErr := r.executableCapability.Execute(context.Background(), capabilityRequest)
			payload, marshalErr := pb.MarshalCapabilityResponse(resp)
			r.sendResponse(messageID, responseErr, payload, marshalErr)
		}

	case remotetypes.MethodRegisterToWorkflow:
		registerRequest, err := pb.UnmarshalRegisterToWorkflowRequest(msg.Payload)
		if err != nil {
			panic(err)
		}

		err = r.registrationServer.Register(ctx, msg, sender, registerRequest.Metadata.WorkflowID, registerRequest.Metadata.ReferenceID)
		if err != nil {
			panic(err)
		}
	case remotetypes.MethodUnregisterFromWorkflow:
		panic("unexpected call, client should explicitly unregister from workflow, expiration of registration is expected to take care of this")
	default:
		panic("unknown method")
	}
}

func (r *clientTestServer) sendResponse(messageID string, responseErr error,
	payload []byte, marshalErr error) {
	for receiver := range r.messageIDToSenders[messageID] {
		var responseMsg = &remotetypes.MessageBody{
			CapabilityId:    "cap_id@1.0.0",
			CapabilityDonId: 1,
			CallerDonId:     r.workflowDonInfo.ID,
			Method:          remotetypes.MethodExecute,
			MessageId:       []byte(messageID),
			Sender:          r.peerID[:],
			Receiver:        receiver[:],
		}

		if responseErr != nil {
			responseMsg.Error = remotetypes.Error_INTERNAL_ERROR
		} else {
			if marshalErr != nil {
				panic(marshalErr)
			}
			responseMsg.Payload = payload
		}

		err := r.dispatcher.Send(receiver, responseMsg)
		if err != nil {
			panic(err)
		}
	}
}
