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

// TODO: use logic from chainlink-common
func digest(r *jsonrpc.Response) (string, error) {
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
func (agg *identicalNodeResponseAggregator) CollectAndAggregate(resp *jsonrpc.Response, nodeAddress string) (*jsonrpc.Response, error) {
	key, err := digest(resp)
	if err != nil {
		return nil, fmt.Errorf("error generating digest for response: %w", err)
	}
	if _, ok := agg.responses[key]; !ok {
		agg.responses[key] = make(stringSet)
	}
	agg.responses[key].Add(nodeAddress)
	if len(agg.responses[key]) < agg.threshold {
		return nil, nil
	}
	return resp, nil
}
