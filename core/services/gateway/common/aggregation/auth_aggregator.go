package aggregation

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	gateway_common "github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
)

type AuthAggregator struct {
	services.StateMachine
	lggr      logger.Logger
	stopCh    services.StopChan
	threshold int
	mu        sync.RWMutex
	// observations is a map with observation digest as key
	observations map[string]*NodeObservations
	// observedAt is a map from node address to a map of workflow digest to last observed time
	observedAt      map[string]map[string]time.Time
	cleanupInterval time.Duration
}

func NewAuthAggregator(lggr logger.Logger, threshold int, cleanupInterval time.Duration) *AuthAggregator {
	return &AuthAggregator{
		lggr:            logger.Named(lggr, "AuthAggregator"),
		threshold:       threshold,
		observations:    make(map[string]*NodeObservations),
		observedAt:      make(map[string]map[string]time.Time),
		stopCh:          make(services.StopChan),
		cleanupInterval: cleanupInterval,
	}
}

func (agg *AuthAggregator) reapObservations() {
	agg.mu.Lock()
	defer agg.mu.Unlock()
	now := time.Now()
	var expiredCount int
	for node, digestObservedAt := range agg.observedAt {
		for digest, observedAt := range digestObservedAt {
			if now.Sub(observedAt) > agg.cleanupInterval {
				delete(digestObservedAt, digest)
				_, ok := agg.observations[digest]
				if !ok {
					agg.lggr.Warnw("Observation digest not found in observations", "digest", digest, "node", node)
					continue
				}
				agg.observations[digest].nodes.Remove(node)
				if len(agg.observations[digest].nodes) == 0 {
					delete(agg.observations, digest)
				}
				expiredCount++
			}
		}
	}
	if expiredCount > 0 {
		agg.lggr.Debugw("Removed expired callbacks", "count", expiredCount)
	}
}

func (agg *AuthAggregator) Start(ctx context.Context) error {
	return agg.StartOnce("AuthAggregator", func() error {
		agg.lggr.Info("Starting AuthAggregator")
		go func() {
			ticker := time.NewTicker(agg.cleanupInterval)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					agg.reapObservations()
				case <-agg.stopCh:
					return
				}
			}
		}()
		return nil
	})
}

func (agg *AuthAggregator) Close() error {
	return agg.StopOnce("AuthAggregator", func() error {
		agg.lggr.Info("Stopping AuthAggregator")
		close(agg.stopCh)
		return nil
	})
}

func (agg *AuthAggregator) Collect(o WorkflowAuthObservation, nodeAddress string) {
	agg.mu.Lock()
	defer agg.mu.Unlock()
	digest := o.Digest()
	_, ok := agg.observedAt[nodeAddress]
	if !ok {
		agg.observedAt[nodeAddress] = make(map[string]time.Time)
	}
	agg.observedAt[nodeAddress][digest] = time.Now()

	_, ok = agg.observations[digest]
	if !ok {
		agg.observations[digest] = &NodeObservations{
			observation: o,
			nodes:       make(StringSet),
		}
	}
	agg.observations[digest].nodes.Add(nodeAddress)
}

func (agg *AuthAggregator) Aggregate() []gateway_common.WorkflowAuthMetadata {
	agg.mu.RLock()
	defer agg.mu.RUnlock()

	var authKeys = make(map[string][]gateway_common.AuthorizedKey)
	for _, nodeObs := range agg.observations {
		if len(nodeObs.nodes) >= agg.threshold {
			authKeys[nodeObs.observation.WorkflowID] = append(authKeys[nodeObs.observation.WorkflowID], nodeObs.observation.AuthorizedKey)
		}
	}

	result := make([]gateway_common.WorkflowAuthMetadata, 0, len(authKeys))
	for workflowID, keys := range authKeys {
		result = append(result, gateway_common.WorkflowAuthMetadata{
			WorkflowID:     workflowID,
			AuthorizedKeys: keys,
		})
	}
	return result
}

type WorkflowAuthObservation struct {
	WorkflowID    string
	AuthorizedKey gateway_common.AuthorizedKey
}

type NodeObservations struct {
	observation WorkflowAuthObservation
	nodes       StringSet
}

func (w WorkflowAuthObservation) Digest() string {
	keyBytes := []byte(fmt.Sprintf("%s:%x", w.AuthorizedKey.KeyType, w.AuthorizedKey.PublicKey))
	data := []byte(w.WorkflowID)
	data = append(data, keyBytes...)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
