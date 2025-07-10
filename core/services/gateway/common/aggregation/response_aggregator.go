package aggregation

import (
	"encoding/json"
	"errors"
	"fmt"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
)

// IdenticalNodeResponseAggregator collects node responses and aggregates identical responses.
// (Usually 2f+1, where f is the number of faulty nodes).
// NOT thread-safe.
type IdenticalNodeResponseAggregator struct {
	// responses is a map from response digest to a set of node addresses
	responses map[string]StringSet
	// nodeToResponse tracks which response key each node address is currently associated with
	nodeToResponse map[string]string
	threshold      int
}

func NewIdenticalNodeResponseAggregator(threshold int) (*IdenticalNodeResponseAggregator, error) {
	if threshold <= 0 {
		return nil, fmt.Errorf("threshold must be greater than 0, got %d", threshold)
	}
	return &IdenticalNodeResponseAggregator{
		responses:      make(map[string]StringSet),
		nodeToResponse: make(map[string]string),
		threshold:      threshold,
	}, nil
}

// CollectAndAggregate tracks responses by response content (hash) and node address.
// If the number of identical responses reaches the threshold, it returns the response.
// Otherwise, returns nil and no error.
// If a node provides a new response that differs from its previous one, the node is
// removed from its previous response group and added to the new one.
func (agg *IdenticalNodeResponseAggregator) CollectAndAggregate(
	resp *jsonrpc.Response[json.RawMessage],
	nodeAddress string) (*jsonrpc.Response[json.RawMessage], error) {
	if resp == nil {
		return nil, errors.New("response cannot be nil")
	}
	if nodeAddress == "" {
		return nil, errors.New("node address cannot be empty")
	}

	digest, err := resp.Digest()
	if err != nil {
		return nil, fmt.Errorf("error generating digest for response: %w", err)
	}

	// Check if the node already submitted a different response
	if oldKey, exists := agg.nodeToResponse[nodeAddress]; exists && oldKey != digest {
		if nodes, ok := agg.responses[oldKey]; ok {
			nodes.Remove(nodeAddress)
			// Clean up empty response groups
			if len(nodes) == 0 {
				delete(agg.responses, oldKey)
			}
		}
	}

	if _, ok := agg.responses[digest]; !ok {
		agg.responses[digest] = make(StringSet)
	}
	agg.responses[digest].Add(nodeAddress)
	agg.nodeToResponse[nodeAddress] = digest

	if len(agg.responses[digest]) >= agg.threshold {
		return resp, nil
	}

	return nil, nil
}
