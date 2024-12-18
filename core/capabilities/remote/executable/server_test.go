package executable_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/executable"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/registration"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

func Test_Server_ExcludesNonDeterministicInputAttributes(t *testing.T) {
	ctx := testutils.Context(t)

	numCapabilityPeers := 4

	callers, srvcs := testRemoteExecutableCapabilityServer(ctx, t, &commoncap.RemoteExecutableConfig{RequestHashExcludedAttributes: []string{"signed_report.Signatures"}},
		func() commoncap.ExecutableCapability { return &TestCapability{} }, 10, 9, numCapabilityPeers, 3, 10*time.Minute)

	for idx, caller := range callers {
		rawInputs := map[string]any{
			"signed_report": map[string]any{"Signatures": "sig" + strconv.Itoa(idx), "Price": 20},
		}

		inputs, err := values.NewMap(rawInputs)
		require.NoError(t, err)

		_, err = caller.Execute(context.Background(),
			commoncap.CapabilityRequest{
				Metadata: commoncap.RequestMetadata{
					WorkflowID:          workflowID1,
					WorkflowExecutionID: workflowExecutionID1,
				},
				Inputs: inputs,
			})
		require.NoError(t, err)
	}

	for _, caller := range callers {
		for i := 0; i < numCapabilityPeers; i++ {
			msg := <-caller.receivedMessages
			assert.Equal(t, remotetypes.Error_OK, msg.Error)
		}
	}
	closeServices(t, srvcs)
}

func Test_Server_Execute_RespondsAfterSufficientRequests(t *testing.T) {
	ctx := testutils.Context(t)

	numCapabilityPeers := 4

	callers, srvcs := testRemoteExecutableCapabilityServer(ctx, t, &commoncap.RemoteExecutableConfig{}, func() commoncap.ExecutableCapability { return &TestCapability{} }, 10, 9, numCapabilityPeers, 3, 10*time.Minute)

	for _, caller := range callers {
		_, err := caller.Execute(context.Background(),
			commoncap.CapabilityRequest{
				Metadata: commoncap.RequestMetadata{
					WorkflowID:          workflowID1,
					WorkflowExecutionID: workflowExecutionID1,
				},
			})
		require.NoError(t, err)
	}

	for _, caller := range callers {
		for i := 0; i < numCapabilityPeers; i++ {
			msg := <-caller.receivedMessages
			assert.Equal(t, remotetypes.Error_OK, msg.Error)
		}
	}
	closeServices(t, srvcs)
}

func Test_Server_RegisterToWorkflow(t *testing.T) {
	ctx := testutils.Context(t)

	numCapabilityPeers := 4

	remoteExecutableConfig := &commoncap.RemoteExecutableConfig{}
	remoteExecutableConfig.RegistrationRefresh = 100 * time.Millisecond

	var testCapabilities []*TestCapability

	clients, srvcs := testRemoteExecutableCapabilityServer(ctx, t, &commoncap.RemoteExecutableConfig{RegistrationRefresh: 100 * time.Millisecond},
		func() commoncap.ExecutableCapability {
			testCap := &TestCapability{}
			testCapabilities = append(testCapabilities, testCap)
			return testCap
		},
		10, 4, numCapabilityPeers, 3, 10*time.Minute)

	for _, caller := range clients {
		err := caller.RegisterToWorkflow(context.Background(), commoncap.RegisterToWorkflowRequest{
			Metadata: commoncap.RegistrationMetadata{
				WorkflowID:    workflowID1,
				ReferenceID:   stepReferenceID1,
				WorkflowOwner: workflowOwnerID,
			},
		})

		require.NoError(t, err)
	}

	require.Eventually(t, func() bool {
		for _, testCapability := range testCapabilities {
			if len(testCapability.GetRegisterRequests()) != 1 {
				return false
			}
		}

		return true
	}, 10*time.Second, 100*time.Millisecond, "expected one registration request to be received")

	// a short sleep to allow the registration refresh mechanism to run, then check that there is still one registration request
	time.Sleep(200 * time.Millisecond)

	for _, testCapability := range testCapabilities {
		assert.Len(t, testCapability.GetRegisterRequests(), 1)
	}

	closeServices(t, srvcs)
}

func Test_Server_RegisterToWorkflow_Error(t *testing.T) {
	ctx := testutils.Context(t)

	numCapabilityPeers := 4

	remoteExecutableConfig := &commoncap.RemoteExecutableConfig{}
	remoteExecutableConfig.RegistrationRefresh = 100 * time.Millisecond

	var testCapabilities []*TestErrorCapability

	clients, srvcs := testRemoteExecutableCapabilityServer(ctx, t, &commoncap.RemoteExecutableConfig{RegistrationRefresh: 100 * time.Millisecond},
		func() commoncap.ExecutableCapability {
			testCap := &TestErrorCapability{}
			testCapabilities = append(testCapabilities, testCap)
			return testCap
		},
		10, 4, numCapabilityPeers, 3, 10*time.Minute)

	for _, caller := range clients {
		err := caller.RegisterToWorkflow(context.Background(), commoncap.RegisterToWorkflowRequest{
			Metadata: commoncap.RegistrationMetadata{
				WorkflowID:    workflowID1,
				ReferenceID:   stepReferenceID1,
				WorkflowOwner: workflowOwnerID,
			},
		})

		require.NoError(t, err)
	}

	// As the registration errors, the client should retry the registration request repeatedly
	require.Eventually(t, func() bool {
		for _, testCapability := range testCapabilities {
			if len(testCapability.GetRegisterRequests()) > 2 {
				return false
			}
		}

		return true
	}, 10*time.Second, 100*time.Millisecond, "expected more than 2 registration requests to be received")

	closeServices(t, srvcs)
}

func Test_Server_UnregisterFromWorkflowIsCalledWhenClientsAreShutdown(t *testing.T) {
	ctx := testutils.Context(t)

	numCapabilityPeers := 4

	remoteExecutableConfig := &commoncap.RemoteExecutableConfig{}
	remoteExecutableConfig.RegistrationRefresh = 100 * time.Millisecond

	var testCapabilities []*TestCapability

	clients, srvcs := testRemoteExecutableCapabilityServer(ctx, t,
		&commoncap.RemoteExecutableConfig{RegistrationRefresh: 50 * time.Millisecond, RegistrationExpiry: 500 * time.Millisecond},
		func() commoncap.ExecutableCapability {
			testCap := &TestCapability{}
			testCapabilities = append(testCapabilities, testCap)
			return testCap
		},
		10, 4, numCapabilityPeers, 3, 10*time.Minute)

	for _, caller := range clients {
		err := caller.RegisterToWorkflow(context.Background(), commoncap.RegisterToWorkflowRequest{
			Metadata: commoncap.RegistrationMetadata{
				WorkflowID:    workflowID1,
				ReferenceID:   stepReferenceID1,
				WorkflowOwner: workflowOwnerID,
			},
		})

		require.NoError(t, err)
	}

	require.Eventually(t, func() bool {
		for _, testCapability := range testCapabilities {
			if len(testCapability.GetRegisterRequests()) != 1 {
				return false
			}
		}

		return true
	}, 10*time.Second, 100*time.Millisecond, "expected one registration request to be received")

	for _, client := range clients {
		require.NoError(t, client.Close())
	}

	require.Eventually(t, func() bool {
		for _, testCapability := range testCapabilities {
			if len(testCapability.GetUnregisterRequests()) != 1 {
				return false
			}
		}

		return true
	}, 10*time.Second, 100*time.Millisecond, "expected one registration request to be received")

	// a short sleep greater than the expiry time then check that there is still only one unregistration request
	time.Sleep(1 * time.Second)

	for _, testCapability := range testCapabilities {
		assert.Len(t, testCapability.GetUnregisterRequests(), 1)
	}

	closeServices(t, srvcs)
}

func Test_Server_InsufficientCallers(t *testing.T) {
	ctx := testutils.Context(t)

	numCapabilityPeers := 4

	callers, srvcs := testRemoteExecutableCapabilityServer(ctx, t, &commoncap.RemoteExecutableConfig{},
		func() commoncap.ExecutableCapability { return &TestCapability{} }, 10, 10, numCapabilityPeers, 3, 100*time.Millisecond)

	for _, caller := range callers {
		_, err := caller.Execute(context.Background(),
			commoncap.CapabilityRequest{
				Metadata: commoncap.RequestMetadata{
					WorkflowID:          workflowID1,
					WorkflowExecutionID: workflowExecutionID1,
				},
			})
		require.NoError(t, err)
	}

	for _, caller := range callers {
		for i := 0; i < numCapabilityPeers; i++ {
			msg := <-caller.receivedMessages
			assert.Equal(t, remotetypes.Error_TIMEOUT, msg.Error)
		}
	}
	closeServices(t, srvcs)
}

func Test_Server_CapabilityError(t *testing.T) {
	ctx := testutils.Context(t)

	numCapabilityPeers := 4

	callers, srvcs := testRemoteExecutableCapabilityServer(ctx, t, &commoncap.RemoteExecutableConfig{},
		func() commoncap.ExecutableCapability { return &TestErrorCapability{} }, 10, 9, numCapabilityPeers, 3, 100*time.Millisecond)

	for _, caller := range callers {
		_, err := caller.Execute(context.Background(),
			commoncap.CapabilityRequest{
				Metadata: commoncap.RequestMetadata{
					WorkflowID:          workflowID1,
					WorkflowExecutionID: workflowExecutionID1,
				},
			})
		require.NoError(t, err)
	}

	for _, caller := range callers {
		for i := 0; i < numCapabilityPeers; i++ {
			msg := <-caller.receivedMessages
			assert.Equal(t, remotetypes.Error_INTERNAL_ERROR, msg.Error)
		}
	}
	closeServices(t, srvcs)
}

func testRemoteExecutableCapabilityServer(ctx context.Context, t *testing.T,
	config *commoncap.RemoteExecutableConfig,
	capabilityFactory func() commoncap.ExecutableCapability,
	numWorkflowPeers int, workflowDonF uint8,
	numCapabilityPeers int, capabilityDonF uint8, capabilityNodeResponseTimeout time.Duration) ([]*serverTestClient, []services.Service) {
	lggr := logger.TestLogger(t)

	capabilityPeers := make([]p2ptypes.PeerID, numCapabilityPeers)
	for i := 0; i < numCapabilityPeers; i++ {
		capabilityPeerID := NewP2PPeerID(t)
		capabilityPeers[i] = capabilityPeerID
	}

	capDonInfo := commoncap.DON{
		ID:      1,
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
		workflowPeers[i] = NewP2PPeerID(t)
	}

	workflowDonInfo := commoncap.DON{
		Members: workflowPeers,
		ID:      2,
		F:       workflowDonF,
	}

	var srvcs []services.Service
	broker := newTestAsyncMessageBroker(t, 1000)
	err := broker.Start(context.Background())
	require.NoError(t, err)
	srvcs = append(srvcs, broker)

	workflowDONs := map[uint32]commoncap.DON{
		workflowDonInfo.ID: workflowDonInfo,
	}

	capabilityNodes := make([]remotetypes.Receiver, numCapabilityPeers)

	for i := 0; i < numCapabilityPeers; i++ {
		capabilityPeer := capabilityPeers[i]
		capabilityDispatcher := broker.NewDispatcherForNode(capabilityPeer)
		capabilityNode := executable.NewServer(config, capabilityPeer, capabilityFactory(), capInfo, capDonInfo, workflowDONs, capabilityDispatcher,
			capabilityNodeResponseTimeout, lggr)
		require.NoError(t, capabilityNode.Start(ctx))
		broker.RegisterReceiverNode(capabilityPeer, capabilityNode)
		capabilityNodes[i] = capabilityNode
		srvcs = append(srvcs, capabilityNode)
	}

	workflowNodes := make([]*serverTestClient, numWorkflowPeers)
	for i := 0; i < numWorkflowPeers; i++ {
		workflowPeerDispatcher := broker.NewDispatcherForNode(workflowPeers[i])
		workflowNode := newServerTestClient(lggr, workflowPeers[i], config.RegistrationRefresh, capInfo, capDonInfo, workflowDonInfo, workflowPeerDispatcher)
		broker.RegisterReceiverNode(workflowPeers[i], workflowNode)
		workflowNodes[i] = workflowNode
		servicetest.Run(t, workflowNode)
	}

	return workflowNodes, srvcs
}

func closeServices(t *testing.T, srvcs []services.Service) {
	for _, srv := range srvcs {
		require.NoError(t, srv.Close())
	}
}

type serverTestClient struct {
	services.StateMachine
	lggr               logger.Logger
	peerID             p2ptypes.PeerID
	dispatcher         remotetypes.Dispatcher
	capabilityDonInfo  commoncap.DON
	receivedMessages   chan *remotetypes.MessageBody
	callerDonID        string
	registrationClient *registration.Client
}

func (r *serverTestClient) Receive(_ context.Context, msg *remotetypes.MessageBody) {
	r.receivedMessages <- msg
}

func newServerTestClient(lggr logger.Logger, peerID p2ptypes.PeerID, registrationRefresh time.Duration, capInfo commoncap.CapabilityInfo,
	capabilityDonInfo commoncap.DON,
	workflowDonInfo commoncap.DON,
	dispatcher remotetypes.Dispatcher) *serverTestClient {
	registrationClient := registration.NewClient(lggr, remotetypes.MethodRegisterToWorkflow, registrationRefresh, capInfo, capabilityDonInfo, workflowDonInfo, dispatcher, "serverTestClient")

	return &serverTestClient{lggr: lggr, peerID: peerID, dispatcher: dispatcher, capabilityDonInfo: capabilityDonInfo,
		receivedMessages: make(chan *remotetypes.MessageBody, 100), callerDonID: "workflow-don",
		registrationClient: registrationClient}
}

func (r *serverTestClient) Info(ctx context.Context) (commoncap.CapabilityInfo, error) {
	panic("not implemented")
}

func (r *serverTestClient) Start(ctx context.Context) error {
	return r.StartOnce(r.peerID.String(), func() error {
		if err := r.registrationClient.Start(ctx); err != nil {
			return fmt.Errorf("failed to start capability register: %w", err)
		}
		return nil
	})
}

func (r *serverTestClient) Close() error {
	r.IfNotStopped(func() {
		if err := r.registrationClient.Close(); err != nil {
			r.lggr.Errorf("failed to close capability register: %v", err)
		}
	})

	return nil
}

func (r *serverTestClient) RegisterToWorkflow(ctx context.Context, req commoncap.RegisterToWorkflowRequest) error {
	rawRequest, err := pb.MarshalRegisterToWorkflowRequest(req)
	if err != nil {
		return err
	}

	err = r.registrationClient.RegisterWorkflow(req.Metadata.WorkflowID, rawRequest)
	if err != nil {
		return err
	}

	return nil
}

func (r *serverTestClient) UnregisterFromWorkflow(ctx context.Context, req commoncap.UnregisterFromWorkflowRequest) error {
	r.registrationClient.UnregisterWorkflow(req.Metadata.WorkflowID)

	return nil
}

func (r *serverTestClient) Execute(ctx context.Context, req commoncap.CapabilityRequest) (<-chan commoncap.CapabilityResponse, error) {
	rawRequest, err := pb.MarshalCapabilityRequest(req)
	if err != nil {
		return nil, err
	}

	messageID := remotetypes.MethodExecute + ":" + req.Metadata.WorkflowExecutionID

	for _, node := range r.capabilityDonInfo.Members {
		message := &remotetypes.MessageBody{
			CapabilityId:    "capability-id",
			CapabilityDonId: 1,
			CallerDonId:     2,
			Method:          remotetypes.MethodExecute,
			Payload:         rawRequest,
			MessageId:       []byte(messageID),
			Sender:          r.peerID[:],
			Receiver:        node[:],
		}

		if err = r.dispatcher.Send(node, message); err != nil {
			return nil, err
		}
	}

	return nil, nil
}
