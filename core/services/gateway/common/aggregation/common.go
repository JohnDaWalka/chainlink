package aggregation

type NodeResponseAggregator interface {
	// CollectAndAggregate appends a node response (as a string) to existing list of responses if exists
	// and tries to aggregate them into a single response string.
	CollectAndAggregate(resp string, nodeAddress string) (string, error)
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
