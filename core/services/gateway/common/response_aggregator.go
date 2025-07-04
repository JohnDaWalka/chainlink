package common

import (
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
)

var _ NodeResponseAggregator = (*identicalNodeResponseAggregator)(nil)

type NodeResponseAggregator interface {
	// CollectAndAggregate appends a node response to existing list of responses if exists
	// and tries to aggregate them into a single response.
	CollectAndAggregate(resp *jsonrpc.Response, nodeAddress string) (*jsonrpc.Response, error)
}

// identicalNodeResponseAggregator collects node responses and aggregates identical responses.
// (Usually 2f+1, where f is the number of faulty nodes).
// NOT thread-safe.
type identicalNodeResponseAggregator struct {
	responses map[string]stringSet
	threshold int
}

func NewIdenticalNodeResponseAggregator(threshold int) NodeResponseAggregator {
	return &identicalNodeResponseAggregator{
		responses: make(map[string]stringSet),
		threshold: threshold,
	}
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

func (a *identicalNodeResponseAggregator) CollectAndAggregate(resp *jsonrpc.Response, nodeAddress string) (*jsonrpc.Response, error) {
	key := "TODO: USE DIGEST() FROM CHAINLINK-COMMON"
	a.responses[key].Add(nodeAddress)
	if len(a.responses[key]) < a.threshold {
		return nil, nil
	}
	return resp, nil
}

// TODO: Remove:
// func (a *identicalNodeResponseAggregator) Start(ctx context.Context) error {
// 	return a.StartOnce("IdenticalNodeResponseAggregator", func() error {
// 		a.lggr.Info("Starting IdenticalNodeResponseAggregator")
// 		go func() {
// 			ticker := time.NewTicker(time.Duration(a.responseCleanUpPeriodMs) * time.Millisecond)
// 			defer ticker.Stop()
// 			for {
// 				select {
// 				case <-ticker.C:
// 					a.cleanUpExpiredResponses()
// 				case <-a.stopCh:
// 					a.lggr.Info("Stopping IdenticalNodeResponseAggregator")
// 					return
// 				}
// 			}
// 		}()
// 		return nil
// 	})
// }

// func (a *identicalNodeResponseAggregator) cleanUpExpiredResponses() {
// 	now := time.Now()
// 	for requestID, aggregated := range a.responses {
// 		if now.Sub(aggregated.lastUpdated) > time.Duration(a.responseMaxAgeMs)*time.Millisecond {
// 			a.lggr.Debugw("Removing expired response", "requestID", requestID, "ageMs", now.Sub(aggregated.lastUpdated).Milliseconds())
// 			delete(a.responses, requestID)
// 		}
// 	}
// }

// func (a *identicalNodeResponseAggregator) Close() error {
// 	return a.StopOnce("IdenticalNodeResponseAggregator", func() error {
// 		a.lggr.Info("Closing IdenticalNodeResponseAggregator")
// 		close(a.stopCh)
// 		return nil
// 	})
// }
