package registration

import (
	"context"
	"sync"
	"time"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

type clientRegistration struct {
	registrationRequest []byte
}

type registerDispatcher interface {
	Send(peerID p2ptypes.PeerID, msgBody *types.MessageBody) error
}

// Client is a shim for remote capabilities that support registration to a workflow.  It polls the Server to ensure
// the registration stays live.  In the current implementation the Server shim will unregister any workflow that has
// not been re-registered within the registrationExpiry interval.
type Client struct {
	services.StateMachine
	lggr                logger.Logger
	registrationMethod  string
	registrationRefresh time.Duration
	capInfo             commoncap.CapabilityInfo
	capDonInfo          commoncap.DON
	localDonInfo        commoncap.DON
	dispatcher          registerDispatcher
	registeredWorkflows map[string]*clientRegistration
	mu                  sync.RWMutex
	stopCh              services.StopChan
	wg                  sync.WaitGroup
}

func NewClient(lggr logger.Logger, registrationMethod string, registrationRefresh time.Duration, capInfo commoncap.CapabilityInfo, capDonInfo commoncap.DON,
	localDonInfo commoncap.DON, dispatcher registerDispatcher, registryType string) *Client {
	return &Client{
		lggr:                lggr.Named(registryType),
		registrationMethod:  registrationMethod,
		registrationRefresh: registrationRefresh,
		capInfo:             capInfo,
		capDonInfo:          capDonInfo,
		localDonInfo:        localDonInfo,
		dispatcher:          dispatcher,
		registeredWorkflows: make(map[string]*clientRegistration),
		stopCh:              make(services.StopChan),
	}
}

func (r *Client) Start(_ context.Context) error {
	return r.StartOnce(r.lggr.Name(), func() error {
		r.wg.Add(1)
		go func() {
			defer r.wg.Done()
			r.registrationLoop()
		}()
		r.lggr.Info("started")
		return nil
	})
}

func (r *Client) Close() error {
	return r.StopOnce(r.lggr.Name(), func() error {
		close(r.stopCh)
		r.wg.Wait()
		r.lggr.Info("closed")
		return nil
	})
}

func (r *Client) RegisterWorkflow(workflowID string, request []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.lggr.Infow("register workflow called", "capabilityId", r.capInfo.ID, "donId", r.capDonInfo.ID, "workflowID", workflowID)
	regState, ok := r.registeredWorkflows[workflowID]
	if !ok {
		regState = &clientRegistration{
			registrationRequest: request,
		}
		r.registeredWorkflows[workflowID] = regState
	} else {
		regState.registrationRequest = request
		r.lggr.Warnw("re-registering workflow", "capabilityId", r.capInfo.ID, "donId", r.capDonInfo.ID, "workflowID", workflowID)
	}

	return nil
}

func (r *Client) UnregisterWorkflow(workflowID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.lggr.Infow("unregister workflow called", "capabilityId", r.capInfo.ID, "donId", r.capDonInfo.ID, "workflowID", workflowID)
	delete(r.registeredWorkflows, workflowID)
	// Registrations will quickly expire on all remote nodes so it is currently considered unnecessary to send
	// unregister messages to the nodes
}

func (r *Client) registrationLoop() {
	ticker := time.NewTicker(r.registrationRefresh)
	defer ticker.Stop()
	for {
		select {
		case <-r.stopCh:
			return
		case <-ticker.C:
			r.mu.RLock()
			r.lggr.Infow("register for remote capability", "capabilityId", r.capInfo.ID, "donId", r.capDonInfo.ID, "nMembers", len(r.capDonInfo.Members), "nWorkflows", len(r.registeredWorkflows))
			if len(r.registeredWorkflows) == 0 {
				r.lggr.Infow("no workflows to register")
			}
			for _, registration := range r.registeredWorkflows {
				// NOTE: send to all by default, introduce different strategies later (KS-76)
				for _, peerID := range r.capDonInfo.Members {
					m := &types.MessageBody{
						CapabilityId:    r.capInfo.ID,
						CapabilityDonId: r.capDonInfo.ID,
						CallerDonId:     r.localDonInfo.ID,
						Method:          r.registrationMethod,
						Payload:         registration.registrationRequest,
					}
					err := r.dispatcher.Send(peerID, m)
					if err != nil {
						r.lggr.Errorw("failed to send message", "capabilityId", r.capInfo.ID, "donId", r.capDonInfo.ID, "peerId", peerID, "err", err)
					}
				}
			}
			r.mu.RUnlock()
		}
	}
}
