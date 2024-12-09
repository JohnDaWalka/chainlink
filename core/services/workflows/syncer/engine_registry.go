package syncer

import (
	"errors"
	"sync"
)

// IsReadyCloser is an abstraction for engines that can be checked for readiness and closed.
type IsReadyCloser interface {
	// Ready returns nil if the engine is ready to be used.
	Ready() error

	// Close closes the engine.
	Close() error
}
type engineRegistry struct {
	engines map[string]IsReadyCloser
	mu      sync.RWMutex
}

func newEngineRegistry() *engineRegistry {
	return &engineRegistry{
		engines: make(map[string]IsReadyCloser),
	}
}

// Add adds an engine to the registry.
func (r *engineRegistry) Add(id string, engine IsReadyCloser) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.engines[id] = engine
}

// Get retrieves an engine from the registry.
func (r *engineRegistry) Get(id string) (IsReadyCloser, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	engine, found := r.engines[id]
	if !found {
		return nil, errors.New("engine not found")
	}
	return engine, nil
}

// IsRunning is true if the engine exists and is ready.
func (r *engineRegistry) IsRunning(id string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	engine, found := r.engines[id]
	if !found {
		return false
	}

	return engine.Ready() == nil
}

// Pop removes an engine from the registry and returns the engine if found.
func (r *engineRegistry) Pop(id string) (IsReadyCloser, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	engine, ok := r.engines[id]
	if !ok {
		return nil, errors.New("remove failed: engine not found")
	}
	delete(r.engines, id)
	return engine, nil
}

// Close closes all engines in the registry.
func (r *engineRegistry) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	var err error
	for id, engine := range r.engines {
		closeErr := engine.Close()
		if closeErr != nil {
			err = errors.Join(err, closeErr)
		}
		delete(r.engines, id)
	}
	return err
}
