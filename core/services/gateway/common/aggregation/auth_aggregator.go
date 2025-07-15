package aggregation

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
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
	// observations is a map that tracks auth data from workflow nodes.
	// keyed by workflow digest
	observations map[string]*NodeObservations
	// observedAt is a map from node address to a map of workflow digest to last observed time
	// This is used to clean up old observations that are no longer relevant.
	observedAt      map[string]map[string]time.Time
	cleanupInterval time.Duration
}

func NewAuthAggregator(lggr logger.Logger, threshold int, cleanupInterval time.Duration) *AuthAggregator {
	if threshold <= 0 {
		panic(fmt.Sprintf("threshold must be greater than 0, got %d", threshold))
	}
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
				delete(agg.observedAt[node], digest)
				if len(agg.observedAt[node]) == 0 {
					delete(agg.observedAt, node)
				}
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

func (agg *AuthAggregator) Start(context.Context) error {
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

// Collect adds an observation from a workflow node to the aggregator.
func (agg *AuthAggregator) Collect(o WorkflowAuthObservation, nodeAddress string) error {
	if o.WorkflowID == "" {
		return errors.New("workflow ID cannot be empty")
	}
	if o.AuthorizedKey.PublicKey == "" {
		return errors.New("authorized key public key cannot be empty")
	}
	if nodeAddress == "" {
		return errors.New("node address cannot be empty")
	}
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
	return nil
}

// Aggregate returns the aggregated auth metadata for workflows that have reached the threshold.
func (agg *AuthAggregator) Aggregate() ([]gateway_common.WorkflowAuthMetadata, error) {
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
	return result, nil
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
