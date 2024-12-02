package remote

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"sync"
	"time"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/registration"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// TriggerPublisher manages all external users of a local trigger capability.
// Its responsibilities are:
//  1. Manage trigger registrations from external nodes (receive, store, aggregate, expire).
//  2. Send out events produced by an underlying, concrete trigger implementation.
//
// TriggerPublisher communicates with corresponding TriggerSubscribers on remote nodes.
type triggerPublisher struct {
	config       *commoncap.RemoteTriggerConfig
	underlying   commoncap.TriggerCapability
	capInfo      commoncap.CapabilityInfo
	capDonInfo   commoncap.DON
	workflowDONs map[uint32]commoncap.DON
	dispatcher   types.Dispatcher

	batchingQueue      map[[32]byte]*batchedResponse
	batchingEnabled    bool
	registrationServer *registration.Server
	bqMu               sync.Mutex // protects batchingQueue
	stopCh             services.StopChan
	wg                 sync.WaitGroup
	lggr               logger.Logger
}

type batchedResponse struct {
	rawResponse    []byte
	callerDonID    uint32
	triggerEventID string
	workflowIDs    []string
}

var _ types.ReceiverService = &triggerPublisher{}

const minAllowedBatchCollectionPeriod = 10 * time.Millisecond

func NewTriggerPublisher(config *commoncap.RemoteTriggerConfig, underlying commoncap.TriggerCapability, capInfo commoncap.CapabilityInfo, capDonInfo commoncap.DON, workflowDONs map[uint32]commoncap.DON, dispatcher types.Dispatcher, lggr logger.Logger) *triggerPublisher {
	if config == nil {
		lggr.Info("no config provided, using default values")
		config = &commoncap.RemoteTriggerConfig{}
	}
	config.ApplyDefaults()

	publisher := &triggerPublisher{
		config:          config,
		underlying:      underlying,
		capInfo:         capInfo,
		capDonInfo:      capDonInfo,
		workflowDONs:    workflowDONs,
		dispatcher:      dispatcher,
		batchingQueue:   make(map[[32]byte]*batchedResponse),
		batchingEnabled: config.MaxBatchSize > 1 && config.BatchCollectionPeriod >= minAllowedBatchCollectionPeriod,
		stopCh:          make(services.StopChan),
		lggr:            lggr.Named("TriggerPublisher"),
	}

	registrationServer := registration.NewServer(lggr, publisher, capInfo, config.RegistrationExpiry, workflowDONs, "TriggerPublisher")

	publisher.registrationServer = registrationServer

	return publisher
}

func (p *triggerPublisher) Start(ctx context.Context) error {
	err := p.registrationServer.Start(ctx)
	if err != nil {
		return fmt.Errorf("failed to start registration server: %w", err)
	}

	if p.batchingEnabled {
		p.wg.Add(1)
		go p.batchingLoop()
	}
	p.lggr.Info("TriggerPublisher started")
	return nil
}

func (p *triggerPublisher) Close() error {
	close(p.stopCh)
	p.wg.Wait()

	err := p.registrationServer.Close()
	if err != nil {
		p.lggr.Errorw("failed to close registration server", "err", err)
	}

	p.lggr.Info("TriggerPublisher closed")
	return nil
}

func (p *triggerPublisher) Receive(ctx context.Context, msg *types.MessageBody) {
	sender, err := ToPeerID(msg.Sender)
	if err != nil {
		p.lggr.Errorw("failed to convert message sender to PeerID", "err", err)
		return
	}

	if msg.Method == types.MethodRegisterTrigger {
		req, err := pb.UnmarshalTriggerRegistrationRequest(msg.Payload)
		if err != nil {
			p.lggr.Errorw("failed to unmarshal trigger registration request", "capabilityId", p.capInfo.ID, "err", err)
			return
		}

		workflowID := req.Metadata.WorkflowID
		err = p.registrationServer.Register(ctx, msg, sender, workflowID, "")
		if err != nil {
			p.lggr.Errorw("failed to register trigger", "capabilityId", p.capInfo.ID, "workflowID",
				SanitizeLogString(workflowID), "callerDonId", msg.CallerDonId, "sender", sender, "err", err)
		}
	} else {
		p.lggr.Errorw("received trigger request with unknown method", "method", SanitizeLogString(msg.Method), "sender", sender)
	}
}

func (p *triggerPublisher) Register(_ context.Context, key registration.Key, registerRequest []byte) error {
	unmarshalled, err := pb.UnmarshalTriggerRegistrationRequest(registerRequest)
	if err != nil {
		return fmt.Errorf("failed to unmarshal request: %w", err)
	}
	ctx, cancel := p.stopCh.NewCtx()
	callbackCh, err := p.underlying.RegisterTrigger(ctx, unmarshalled)
	cancel()
	if err != nil {
		return fmt.Errorf("failed to register trigger: %w", err)
	}

	p.wg.Add(1)
	go p.triggerEventLoop(callbackCh, key)
	return nil
}

func (p *triggerPublisher) Unregister(ctx context.Context, registerRequest []byte) error {
	unmarshalled, err := pb.UnmarshalTriggerRegistrationRequest(registerRequest)
	if err != nil {
		return fmt.Errorf("failed to unmarshal registration request: %w", err)
	}

	return p.underlying.UnregisterTrigger(ctx, unmarshalled)
}

func (p *triggerPublisher) triggerEventLoop(callbackCh <-chan commoncap.TriggerResponse, key registration.Key) {
	defer p.wg.Done()
	for {
		select {
		case <-p.stopCh:
			return
		case response, ok := <-callbackCh:
			if !ok {
				p.lggr.Infow("triggerEventLoop channel closed", "capabilityId", p.capInfo.ID, "workflowID", key.WorkflowID)
				return
			}
			triggerEvent := response.Event
			p.lggr.Debugw("received trigger event", "capabilityId", p.capInfo.ID, "workflowID", key.WorkflowID, "triggerEventID", triggerEvent.ID)
			marshaledResponse, err := pb.MarshalTriggerResponse(response)
			if err != nil {
				p.lggr.Debugw("can't marshal trigger event", "err", err)
				break
			}

			if p.batchingEnabled {
				p.enqueueForBatching(marshaledResponse, key, triggerEvent.ID)
			} else {
				// a single-element "batch"
				p.sendBatch(&batchedResponse{
					rawResponse:    marshaledResponse,
					callerDonID:    key.CallerDonID,
					triggerEventID: triggerEvent.ID,
					workflowIDs:    []string{key.WorkflowID},
				})
			}
		}
	}
}

func (p *triggerPublisher) enqueueForBatching(rawResponse []byte, key registration.Key, triggerEventID string) {
	// put in batching queue, group by hash(callerDonId, triggerEventID, response)
	combined := make([]byte, 4)
	binary.LittleEndian.PutUint32(combined, key.CallerDonID)
	combined = append(combined, []byte(triggerEventID)...)
	combined = append(combined, rawResponse...)
	sha := sha256.Sum256(combined)
	p.bqMu.Lock()
	elem, exists := p.batchingQueue[sha]
	if !exists {
		elem = &batchedResponse{
			rawResponse:    rawResponse,
			callerDonID:    key.CallerDonID,
			triggerEventID: triggerEventID,
			workflowIDs:    []string{key.WorkflowID},
		}
		p.batchingQueue[sha] = elem
	} else {
		elem.workflowIDs = append(elem.workflowIDs, key.WorkflowID)
	}
	p.bqMu.Unlock()
}

func (p *triggerPublisher) sendBatch(resp *batchedResponse) {
	for len(resp.workflowIDs) > 0 {
		idBatch := resp.workflowIDs
		if p.batchingEnabled && int64(len(idBatch)) > int64(p.config.MaxBatchSize) {
			idBatch = idBatch[:p.config.MaxBatchSize]
			resp.workflowIDs = resp.workflowIDs[p.config.MaxBatchSize:]
		} else {
			resp.workflowIDs = nil
		}
		msg := &types.MessageBody{
			CapabilityId:    p.capInfo.ID,
			CapabilityDonId: p.capDonInfo.ID,
			CallerDonId:     resp.callerDonID,
			Method:          types.MethodTriggerEvent,
			Payload:         resp.rawResponse,
			Metadata: &types.MessageBody_TriggerEventMetadata{
				TriggerEventMetadata: &types.TriggerEventMetadata{
					WorkflowIds:    idBatch,
					TriggerEventId: resp.triggerEventID,
				},
			},
		}
		// NOTE: send to all nodes by default, introduce different strategies later (KS-76)
		for _, peerID := range p.workflowDONs[resp.callerDonID].Members {
			err := p.dispatcher.Send(peerID, msg)
			if err != nil {
				p.lggr.Errorw("failed to send trigger event", "capabilityId", p.capInfo.ID, "peerID", peerID, "err", err)
			}
		}
	}
}

func (p *triggerPublisher) batchingLoop() {
	defer p.wg.Done()
	ticker := time.NewTicker(p.config.BatchCollectionPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-p.stopCh:
			return
		case <-ticker.C:
			p.bqMu.Lock()
			queue := p.batchingQueue
			p.batchingQueue = make(map[[32]byte]*batchedResponse)
			p.bqMu.Unlock()

			for _, elem := range queue {
				p.sendBatch(elem)
			}
		}
	}
}

func (p *triggerPublisher) Ready() error {
	return nil
}

func (p *triggerPublisher) HealthReport() map[string]error {
	return nil
}

func (p *triggerPublisher) Name() string {
	return p.lggr.Name()
}
