package aggregation

import (
	"sync"
)

var _ NodeResponseAggregator = (*AuthAggregator)(nil)

type AuthReport struct {
	WorkflowID     string
	AuthorizedKeys []string
	NodeAddress    string
}

type AggregatedAuthResult struct {
	WorkflowID    string
	AuthorizedKey string
	NodeAddresses []string
	ThresholdMet  bool
}

type AuthAggregator struct {
	mu        sync.Mutex
	threshold int
	// map[workflowID][authorizedKey]set[nodeAddress]
	reports map[string]map[string]map[string]struct{}
}

func NewAuthAggregator(threshold int, cleanupInterval time.Duration) *AuthAggregator {
	return &AuthAggregator{
		threshold:       threshold,
		reports:         make(map[string]map[string]map[string]struct{}),
		cleanupInterval: cleanupInterval,
		stopCh:          make(chan struct{}),
	}
}

var _ services.Service = (*AuthAggregator)(nil)

func (a *AuthAggregator) Start(ctx context.Context) error {
	a.mu.Lock()
	if a.running {
		a.mu.Unlock()
		return nil
	}
	a.running = true
	a.mu.Unlock()

	go a.cleanupLoop()
	return nil
}

func (a *AuthAggregator) Close() error {
	a.mu.Lock()
	if !a.running {
		a.mu.Unlock()
		return nil
	}
	a.running = false
	close(a.stopCh)
	a.mu.Unlock()
	return nil
}

func (a *AuthAggregator) Name() string { return "AuthAggregator" }
func (a *AuthAggregator) HealthReport() map[string]error { return map[string]error{} }

type cleanupable interface {
	cleanup()
}

func (a *AuthAggregator) cleanupLoop() {
	ticker := time.NewTicker(a.cleanupInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			a.cleanup()
		case <-a.stopCh:
			return
		}
	}
}

// cleanup removes old entries from the reports map.
// You can customize this logic as needed.
func (a *AuthAggregator) cleanup() {
	a.mu.Lock()
	defer a.mu.Unlock()
	// Example: clear all reports (customize as needed)
	a.reports = make(map[string]map[string]map[string]struct{})
}

// Add these fields to AuthAggregator struct:
// cleanupInterval time.Duration
// stopCh chan struct{}
// running bool

func (a *AuthAggregator) CollectAndAggregate(resp string, nodeAddress string) (string, error)
	a.mu.Lock()
	defer a.mu.Unlock()

	var results []AggregatedAuthResult

	for _, report := range reports {
		for _, key := range report.AuthorizedKeys {
			if _, ok := a.reports[report.WorkflowID]; !ok {
				a.reports[report.WorkflowID] = make(map[string]map[string]struct{})
			}
			if _, ok := a.reports[report.WorkflowID][key]; !ok {
				a.reports[report.WorkflowID][key] = make(map[string]struct{})
			}
			a.reports[report.WorkflowID][key][report.NodeAddress] = struct{}{}

			// Check if threshold met
			if len(a.reports[report.WorkflowID][key]) == a.threshold {
				nodeAddresses := make([]string, 0, len(a.reports[report.WorkflowID][key]))
				for addr := range a.reports[report.WorkflowID][key] {
					nodeAddresses = append(nodeAddresses, addr)
				}
				results = append(results, AggregatedAuthResult{
					WorkflowID:    report.WorkflowID,
					AuthorizedKey: key,
					NodeAddresses: nodeAddresses,
					ThresholdMet:  true,
				})
			}
		}
	}

	return results
}
