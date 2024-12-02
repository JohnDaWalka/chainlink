package registration

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/aggregation"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/messagecache"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/validation"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

type Key struct {
	CallerDonID   uint32
	WorkflowID    string
	StepReference string
}

type serverRegistration struct {
	registrationRequest []byte
}

type target interface {
	Register(ctx context.Context, key Key, registerRequest []byte) error
	Unregister(ctx context.Context, registerRequest []byte) error
}

// Server is a shim for remote capabilities that support registration to a workflow.  It aggregates registration requests
// and invokes the register method on the target capability when the minimum number of registrations are received (2f+1).
// The server will also periodically clean up expired registrations.  A registration is considered expired if it has not
// been aggregated within the registrationExpiry period, when a registration is expired unregister is called on the target
type Server struct {
	lggr               logger.Logger
	capInfo            commoncap.CapabilityInfo
	registrationExpiry time.Duration
	target             target
	registrations      map[Key]*serverRegistration
	messageCache       *messagecache.MessageCache[Key, p2ptypes.PeerID]
	membersCache       map[uint32]map[p2ptypes.PeerID]bool
	workflowDONs       map[uint32]commoncap.DON

	stopCh services.StopChan
	wg     sync.WaitGroup

	mu sync.RWMutex
}

func NewServer(lggr logger.Logger, target target, capInfo commoncap.CapabilityInfo, registrationExpiry time.Duration, workflowDONs map[uint32]commoncap.DON, serverType string) *Server {
	membersCache := make(map[uint32]map[p2ptypes.PeerID]bool)
	for id, don := range workflowDONs {
		cache := make(map[p2ptypes.PeerID]bool)
		for _, member := range don.Members {
			cache[member] = true
		}
		membersCache[id] = cache
	}

	return &Server{
		lggr:               lggr.Named(serverType),
		capInfo:            capInfo,
		target:             target,
		registrationExpiry: registrationExpiry,
		stopCh:             make(services.StopChan),
		registrations:      make(map[Key]*serverRegistration),
		messageCache:       messagecache.New[Key, p2ptypes.PeerID](),
		membersCache:       membersCache,
		workflowDONs:       workflowDONs,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.registrationCleanupLoop()
	}()

	s.lggr.Info("started")
	return nil
}

func (s *Server) Close() error {
	close(s.stopCh)
	s.wg.Wait()
	s.lggr.Info("closed")
	return nil
}

func (s *Server) Register(ctx context.Context, msg *types.MessageBody, sender p2ptypes.PeerID, workflowID string, stepReference string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	callerDon, ok := s.workflowDONs[msg.CallerDonId]
	if !ok {
		return errors.New("received a message from unsupported workflow DON")
	}
	if !s.membersCache[msg.CallerDonId][sender] {
		return errors.New("sender not a member of its workflow DON")
	}
	if err := validation.ValidateWorkflowOrExecutionID(workflowID); err != nil {
		return fmt.Errorf("received request with invalid workflow ID: %w", err)
	}

	s.lggr.Debugw("received registration", "capabilityId", s.capInfo.ID, "workflowId", workflowID, "sender", sender)
	key := Key{CallerDonID: msg.CallerDonId, WorkflowID: workflowID, StepReference: stepReference}
	nowMs := time.Now().UnixMilli()
	s.messageCache.Insert(key, sender, nowMs, msg.Payload)
	_, exists := s.registrations[key]
	if exists {
		s.lggr.Debugw("registration already exists", "capabilityId", s.capInfo.ID, "workflowId", workflowID)
		return nil
	}
	// NOTE: require 2F+1 by default, introduce different strategies later (KS-76)
	minRequired := uint32(2*callerDon.F + 1)
	ready, payloads := s.messageCache.Ready(key, minRequired, nowMs-s.registrationExpiry.Milliseconds(), false)
	if !ready {
		s.lggr.Debugw("not ready to aggregate yet", "capabilityId", s.capInfo.ID, "workflowId", workflowID, "minRequired", minRequired)
		return nil
	}
	aggregated, err := aggregation.AggregateModeRaw(payloads, uint32(callerDon.F+1))
	if err != nil {
		return fmt.Errorf("failed to aggregate registrations: %w", err)
	}
	err = s.target.Register(ctx, key, aggregated)
	if err != nil {
		return fmt.Errorf("failed to register request on target: %w", err)
	}

	s.registrations[key] = &serverRegistration{
		registrationRequest: aggregated,
	}
	s.lggr.Debugw("updated registration", "capabilityId", s.capInfo.ID, "workflowId", workflowID)
	return nil
}

func (s *Server) registrationCleanupLoop() {
	ticker := time.NewTicker(s.registrationExpiry)
	defer ticker.Stop()
	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			now := time.Now().UnixMilli()
			s.mu.Lock()
			for key, req := range s.registrations {
				callerDon := s.workflowDONs[key.CallerDonID]
				ready, _ := s.messageCache.Ready(key, uint32(2*callerDon.F+1), now-s.registrationExpiry.Milliseconds(), false)
				if !ready {
					s.lggr.Infow("registration expired", "capabilityId", s.capInfo.ID, "callerDonID", key.CallerDonID, "workflowId", key.WorkflowID)
					ctx, cancel := s.stopCh.NewCtx()
					err := s.target.Unregister(ctx, req.registrationRequest)
					cancel()
					s.lggr.Infow("unregistered", "capabilityId", s.capInfo.ID, "callerDonID", key.CallerDonID, "workflowId", key.WorkflowID, "err", err)
					delete(s.registrations, key)
					s.messageCache.Delete(key)
				}
			}
			s.mu.Unlock()
		}
	}
}
