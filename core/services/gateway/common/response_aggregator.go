package common

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

var _ NodeResponseAggregator = (*identicalNodeResponseAggregator)(nil)

type NodeResponseAggregator interface {
	job.ServiceCtx
	// CollectAndAggregate appends a node response to existing list of responses if exists
	// and tries to aggregate them into a single response.
	CollectAndAggregate(requestID string, resp *jsonrpc.Response, nodeAddress string) (*jsonrpc.Response, error)
}

// identicalNodeResponseAggregator collects node responses and aggregates them
// if they are identical. It requires at least F+1 identical responses to return a response
// to the user.
// NOT thread-safe.
type identicalNodeResponseAggregator struct {
	services.StateMachine
	// responses is a map. requestID -> aggregatedResponses
	responses               map[string]aggregatedResponses
	f                       int
	stopCh                  services.StopChan
	lggr                    logger.Logger
	responseMaxAgeMs        int
	responseCleanUpPeriodMs int
}

func NewIdenticalNodeResponseAggregator(lggr logger.Logger, f int, responseMaxAgeMs int, responseCleanUpPeriodMs int) *identicalNodeResponseAggregator {
	return &identicalNodeResponseAggregator{
		responses:               make(map[string]aggregatedResponses),
		f:                       f,
		stopCh:                  make(services.StopChan),
		lggr:                    logger.Named(lggr, "IdenticalNodeResponseAggregator"),
		responseMaxAgeMs:        responseMaxAgeMs,
		responseCleanUpPeriodMs: responseCleanUpPeriodMs,
	}
}

func (a *identicalNodeResponseAggregator) Start(ctx context.Context) error {
	return a.StartOnce("IdenticalNodeResponseAggregator", func() error {
		a.lggr.Info("Starting IdenticalNodeResponseAggregator")
		go func() {
			ticker := time.NewTicker(time.Duration(a.responseCleanUpPeriodMs) * time.Millisecond)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					a.cleanUpExpiredResponses()
				case <-a.stopCh:
					a.lggr.Info("Stopping IdenticalNodeResponseAggregator")
					return
				}
			}
		}()
		return nil
	})
}

func (a *identicalNodeResponseAggregator) cleanUpExpiredResponses() {
	a.lggr.Debug("Cleaning up expired responses")
	now := time.Now()
	for requestID, aggregated := range a.responses {
		if now.Sub(aggregated.lastUpdated) > time.Duration(a.responseMaxAgeMs)*time.Millisecond {
			a.lggr.Debugw("Removing expired response", "requestID", requestID, "ageMs", now.Sub(aggregated.lastUpdated).Milliseconds())
			delete(a.responses, requestID)
		}
	}
}

func (a *identicalNodeResponseAggregator) Close() error {
	return a.StopOnce("IdenticalNodeResponseAggregator", func() error {
		a.lggr.Info("Closing IdenticalNodeResponseAggregator")
		close(a.stopCh)
		return nil
	})
}

type aggregatedResponses struct {
	// nodeAddressesByResp is a map. response hash -> node addresses
	nodeAddressesByResp map[string]stringSet
	lastUpdated         time.Time
}

// stringSet is a simple set implementation for strings.
type stringSet map[string]struct{}

func (s stringSet) Add(val string) {
	s[val] = struct{}{}
}

func (s stringSet) Contains(val string) bool {
	_, exists := s[val]
	return exists
}

func (s stringSet) Remove(val string) {
	delete(s, val)
}

func (s stringSet) Values() []string {
	values := make([]string, 0, len(s))
	for k := range s {
		values = append(values, k)
	}
	return values
}

func (a *identicalNodeResponseAggregator) CollectAndAggregate(requestID string, resp *jsonrpc.Response, nodeAddress string) (*jsonrpc.Response, error) {
	a.lggr.Debugw("Collecting node response", "requestID", requestID, "nodeAddress", nodeAddress)
	_, exists := a.responses[requestID]
	if !exists {
		a.responses[requestID] = aggregatedResponses{
			nodeAddressesByResp: make(map[string]stringSet),
		}
	}
	key, err := hashResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to hash response: %w", err)
	}
	aggregatedResp := a.responses[requestID]
	aggregatedResp.nodeAddressesByResp[key].Add(nodeAddress)
	aggregatedResp.lastUpdated = time.Now()
	a.responses[requestID] = aggregatedResp
	if len(aggregatedResp.nodeAddressesByResp[key]) < (2*a.f)+1 {
		return nil, fmt.Errorf("not enough responses to aggregate: got %d, need at least %d", len(a.responses), (2*a.f)+1)
	}
	return resp, nil
}

func hashResponse(resp *jsonrpc.Response) (string, error) {
	respJson, err := json.Marshal(resp)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %w", err)
	}
	s := sha256.New()
	_, err = s.Write([]byte(respJson))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(s.Sum(nil)), nil
}
