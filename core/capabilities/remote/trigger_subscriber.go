package remote

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/aggregation"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/messagecache"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/registration"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

const (
	defaultSendChannelBufferSize = 1000
)

// TriggerSubscriber is a shim for remote trigger capabilities.
// It translatesd between capability API calls and network messages.
// Its responsibilities are:
//  1. Periodically refresh all registrations for remote triggers.
//  2. Collect trigger events from remote nodes and aggregate responses via a customizable aggregator.
//
// TriggerSubscriber communicates with corresponding TriggerReceivers on remote nodes.
type triggerSubscriber struct {
	config             *commoncap.RemoteTriggerConfig
	capInfo            commoncap.CapabilityInfo
	capDonMembers      map[p2ptypes.PeerID]struct{}
	aggregator         types.Aggregator
	messageCache       *messagecache.MessageCache[triggerEventKey, p2ptypes.PeerID]
	mu                 sync.RWMutex // protects registeredWorkflows and messageCache
	stopCh             services.StopChan
	wg                 sync.WaitGroup
	lggr               logger.Logger
	registrationClient *registration.Client
	responseChannels   map[string]chan commoncap.TriggerResponse
}

type triggerEventKey struct {
	triggerEventID string
	workflowID     string
}

type TriggerSubscriber interface {
	commoncap.TriggerCapability
	Receive(ctx context.Context, msg *types.MessageBody)
}

var _ commoncap.TriggerCapability = &triggerSubscriber{}
var _ types.Receiver = &triggerSubscriber{}
var _ services.Service = &triggerSubscriber{}

// TODO makes this configurable with a default
const (
	maxBatchedWorkflowIDs = 1000
)

func NewTriggerSubscriber(config *commoncap.RemoteTriggerConfig, capInfo commoncap.CapabilityInfo, capDonInfo commoncap.DON, localDonInfo commoncap.DON, dispatcher types.Dispatcher, aggregator types.Aggregator, lggr logger.Logger) *triggerSubscriber {
	if aggregator == nil {
		lggr.Warnw("no aggregator provided, using default MODE aggregator", "capabilityId", capInfo.ID)
		aggregator = aggregation.NewDefaultModeAggregator(uint32(capDonInfo.F + 1))
	}
	if config == nil {
		lggr.Info("no config provided, using default values")
		config = &commoncap.RemoteTriggerConfig{}
	}
	config.ApplyDefaults()
	capDonMembers := make(map[p2ptypes.PeerID]struct{})
	for _, member := range capDonInfo.Members {
		capDonMembers[member] = struct{}{}
	}
	return &triggerSubscriber{
		config:             config,
		capInfo:            capInfo,
		capDonMembers:      capDonMembers,
		aggregator:         aggregator,
		messageCache:       messagecache.New[triggerEventKey, p2ptypes.PeerID](),
		stopCh:             make(services.StopChan),
		lggr:               lggr.Named("TriggerSubscriber"),
		registrationClient: registration.NewClient(lggr, types.MethodRegisterTrigger, config.RegistrationRefresh, capInfo, capDonInfo, localDonInfo, dispatcher, "TriggerSubscriber"),
		responseChannels:   make(map[string]chan commoncap.TriggerResponse),
	}
}

func (s *triggerSubscriber) Start(ctx context.Context) error {
	s.wg.Add(1)
	err := s.registrationClient.Start(ctx)
	if err != nil {
		return fmt.Errorf("failed to start capability register: %w", err)
	}

	go s.eventCleanupLoop()
	s.lggr.Info("TriggerSubscriber started")
	return nil
}

func (s *triggerSubscriber) Info(ctx context.Context) (commoncap.CapabilityInfo, error) {
	return s.capInfo, nil
}

func (s *triggerSubscriber) RegisterTrigger(ctx context.Context, request commoncap.TriggerRegistrationRequest) (<-chan commoncap.TriggerResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rawRequest, err := pb.MarshalTriggerRegistrationRequest(request)
	if err != nil {
		return nil, err
	}
	if request.Metadata.WorkflowID == "" {
		return nil, errors.New("empty workflowID")
	}
	if err := s.registrationClient.RegisterWorkflow(request.Metadata.WorkflowID, rawRequest); err != nil {
		return nil, fmt.Errorf("failed to register workflow: %w", err)
	}

	responseChannel := make(chan commoncap.TriggerResponse, defaultSendChannelBufferSize)
	s.responseChannels[request.Metadata.WorkflowID] = responseChannel
	return responseChannel, nil
}

func (s *triggerSubscriber) UnregisterTrigger(ctx context.Context, request commoncap.TriggerRegistrationRequest) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	responseChannel := s.responseChannels[request.Metadata.WorkflowID]
	if responseChannel != nil {
		close(responseChannel)
		delete(s.responseChannels, request.Metadata.WorkflowID)
	}

	s.registrationClient.UnregisterWorkflow(request.Metadata.WorkflowID)
	return nil
}

func (s *triggerSubscriber) isWorkflowRegistered(workflowID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, found := s.responseChannels[workflowID]
	return found
}

var errResponseChannelNotFound = errors.New("response channel not found")

func (s *triggerSubscriber) sendResponse(workflowID string, response commoncap.TriggerResponse) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	responseChannel, found := s.responseChannels[workflowID]
	if !found {
		return errResponseChannelNotFound
	}

	responseChannel <- response
	return nil
}

func (s *triggerSubscriber) Receive(_ context.Context, msg *types.MessageBody) {
	sender, err := ToPeerID(msg.Sender)
	if err != nil {
		s.lggr.Errorw("failed to convert message sender to PeerID", "err", err)
		return
	}

	if _, found := s.capDonMembers[sender]; !found {
		s.lggr.Errorw("received message from unexpected node", "capabilityId", s.capInfo.ID, "sender", sender)
		return
	}
	if msg.Method == types.MethodTriggerEvent {
		meta := msg.GetTriggerEventMetadata()
		if meta == nil {
			s.lggr.Errorw("received message with invalid trigger metadata", "capabilityId", s.capInfo.ID, "sender", sender)
			return
		}
		if len(meta.WorkflowIds) > maxBatchedWorkflowIDs {
			s.lggr.Errorw("received message with too many workflow IDs - truncating", "capabilityId", s.capInfo.ID, "nWorkflows", len(meta.WorkflowIds), "sender", sender)
			meta.WorkflowIds = meta.WorkflowIds[:maxBatchedWorkflowIDs]
		}
		for _, workflowID := range meta.WorkflowIds {
			registered := s.isWorkflowRegistered(workflowID)
			if !registered {
				s.lggr.Errorw("received message for unregistered workflow", "capabilityId", s.capInfo.ID, "workflowID", SanitizeLogString(workflowID), "sender", sender)
				continue
			}
			key := triggerEventKey{
				triggerEventID: meta.TriggerEventId,
				workflowID:     workflowID,
			}
			nowMs := time.Now().UnixMilli()
			s.mu.Lock()
			creationTs := s.messageCache.Insert(key, sender, nowMs, msg.Payload)
			ready, payloads := s.messageCache.Ready(key, s.config.MinResponsesToAggregate, nowMs-s.config.MessageExpiry.Milliseconds(), true)
			s.mu.Unlock()
			if nowMs-creationTs > s.config.RegistrationExpiry.Milliseconds() {
				s.lggr.Warnw("received trigger event for an expired ID", "triggerEventID", meta.TriggerEventId, "capabilityId", s.capInfo.ID, "workflowId", workflowID, "sender", sender)
				continue
			}
			if ready {
				s.lggr.Debugw("trigger event ready to aggregate", "triggerEventID", meta.TriggerEventId, "capabilityId", s.capInfo.ID, "workflowId", workflowID)
				aggregatedResponse, err := s.aggregator.Aggregate(meta.TriggerEventId, payloads)
				if err != nil {
					s.lggr.Errorw("failed to aggregate responses", "triggerEventID", meta.TriggerEventId, "capabilityId", s.capInfo.ID, "workflowId", workflowID, "err", err)
					continue
				}
				s.lggr.Infow("remote trigger event aggregated", "triggerEventID", meta.TriggerEventId, "capabilityId", s.capInfo.ID, "workflowId", workflowID)
				err = s.sendResponse(workflowID, aggregatedResponse)
				// Possible that the response channel for the workflow was unregistered between the check that is was registered and here, so we ignore the error
				if err != nil && !errors.Is(err, errResponseChannelNotFound) {
					s.lggr.Errorw("failed to send response for workflow", "triggerEventID", meta.TriggerEventId, "capabilityId", s.capInfo.ID, "workflowID", workflowID, "err", err)
				}
			}
		}
	} else {
		s.lggr.Errorw("received trigger event with unknown method", "method", SanitizeLogString(msg.Method), "sender", sender)
	}
}

func (s *triggerSubscriber) eventCleanupLoop() {
	defer s.wg.Done()
	ticker := time.NewTicker(s.config.MessageExpiry)
	defer ticker.Stop()
	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.mu.Lock()
			s.messageCache.DeleteOlderThan(time.Now().UnixMilli() - s.config.MessageExpiry.Milliseconds())
			s.mu.Unlock()
		}
	}
}

func (s *triggerSubscriber) Close() error {
	close(s.stopCh)
	s.wg.Wait()
	err := s.registrationClient.Close()
	if err != nil {
		return fmt.Errorf("failed to close capability register: %w", err)
	}
	s.lggr.Info("TriggerSubscriber closed")
	return nil
}

func (s *triggerSubscriber) Ready() error {
	return nil
}

func (s *triggerSubscriber) HealthReport() map[string]error {
	return nil
}

func (s *triggerSubscriber) Name() string {
	return s.lggr.Name()
}
