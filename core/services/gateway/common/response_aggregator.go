package common

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
)

var _ NodeResponseAggregator = (*identicalNodeResponseAggregator)(nil)

type NodeResponseAggregator interface {
	// CollectAndAggregate appends a node response to existing list of responses if exists
	// and tries to aggregate them into a single response.
	CollectAndAggregate(resp *jsonrpc.Response[json.RawMessage], nodeAddress string) (*jsonrpc.Response[json.RawMessage], error)
}

// identicalNodeResponseAggregator collects node responses and aggregates identical responses.
// (Usually 2f+1, where f is the number of faulty nodes).
// NOT thread-safe.
type identicalNodeResponseAggregator struct {
	responses map[string]stringSet
	// nodeToResponse tracks which response key each node address is currently associated with
	nodeToResponse map[string]string
	threshold      int
}

func NewIdenticalNodeResponseAggregator(threshold int) (NodeResponseAggregator, error) {
	if threshold <= 0 {
		return nil, fmt.Errorf("threshold must be greater than 0, got %d", threshold)
	}
	return &identicalNodeResponseAggregator{
		responses:      make(map[string]stringSet),
		nodeToResponse: make(map[string]string),
		threshold:      threshold,
	}, nil
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

// TODO: use logic from chainlink-common
func digest(r *jsonrpc.Response[json.RawMessage]) (string, error) {
	canonicalJSONBytes, err := json.Marshal(r)
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %w", err)
	}

	hasher := sha256.New()
	hasher.Write(canonicalJSONBytes)
	digestBytes := hasher.Sum(nil)

	return hex.EncodeToString(digestBytes), nil
}

// CollectAndAggregate tracks responses by response content (hash) and node address.
// If the number of identical responses reaches the threshold, it returns the response.
// Otherwise, returns nil and no error.
// If a node provides a new response that differs from its previous one, the node is
// removed from its previous response group and added to the new one.
func (agg *identicalNodeResponseAggregator) CollectAndAggregate(
	resp *jsonrpc.Response[json.RawMessage],
	nodeAddress string) (*jsonrpc.Response[json.RawMessage], error) {
	if resp == nil {
		return nil, fmt.Errorf("response cannot be nil")
	}
	if nodeAddress == "" {
		return nil, fmt.Errorf("node address cannot be empty")
	}

	key, err := digest(resp)
	if err != nil {
		return nil, fmt.Errorf("error generating digest for response: %w", err)
	}

	// Check if the node already submitted a different response
	if oldKey, exists := agg.nodeToResponse[nodeAddress]; exists && oldKey != key {
		if nodes, ok := agg.responses[oldKey]; ok {
			nodes.Remove(nodeAddress)
			// Clean up empty response groups
			if len(nodes) == 0 {
				delete(agg.responses, oldKey)
			}
		}
	}

	if _, ok := agg.responses[key]; !ok {
		agg.responses[key] = make(stringSet)
	}
	agg.responses[key].Add(nodeAddress)
	agg.nodeToResponse[nodeAddress] = key

	if len(agg.responses[key]) >= agg.threshold {
		return resp, nil
	}

	return nil, nil
}
