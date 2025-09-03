package worker

import (
	"fmt"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	libformat "github.com/smartcontractkit/chainlink/system-tests/lib/format"
)

// Future represents a typed result from a background task
type Future[T any] struct {
	value   T
	err     error
	done    chan struct{}
	isReady bool
	mutex   sync.RWMutex
}

// Get blocks until the task completes and returns the typed result
func (f *Future[T]) Get() (T, error) {
	<-f.done
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	return f.value, f.err
}

// IsReady returns true if the task has completed (non-blocking)
func (f *Future[T]) IsReady() bool {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	return f.isReady
}

// TypedBackgroundRunner manages execution of background tasks with full type safety
type TypedBackgroundRunner struct {
	wg             sync.WaitGroup
	panicCh        chan panicInfo
	errorCh        chan errorInfo
	taskNames      []string
	taskNamesMutex sync.Mutex
}

type panicInfo struct {
	taskName   string
	panicValue any
	panicStack []byte
}

type errorInfo struct {
	taskName   string
	errorValue error
}

// NewTypedBackgroundRunner creates a new typed background runner
func NewTypedBackgroundRunner() *TypedBackgroundRunner {
	return &TypedBackgroundRunner{
		panicCh: make(chan panicInfo, 10),
		errorCh: make(chan errorInfo, 10),
	}
}

// Add starts a background task and returns a typed Future.
func Add[T any](runner *TypedBackgroundRunner, name string, fn func() (T, error)) *Future[T] {
	future := &Future[T]{
		done: make(chan struct{}),
	}

	runner.taskNamesMutex.Lock()
	runner.taskNames = append(runner.taskNames, name)
	runner.taskNamesMutex.Unlock()

	runner.wg.Add(1)
	go func() {
		defer runner.wg.Done()
		defer close(future.done)
		defer func() {
			if p := recover(); p != nil {
				runner.panicCh <- panicInfo{
					taskName:   name,
					panicValue: p,
					panicStack: debug.Stack(),
				}
			}
		}()

		fmt.Print(libformat.PurpleText("\n---> [BACKGROUND] Starting %s\n\n", name))
		start := time.Now()

		value, err := fn()

		future.mutex.Lock()
		future.value = value
		future.err = err
		future.isReady = true
		future.mutex.Unlock()

		if err == nil {
			fmt.Print(libformat.PurpleText("\n<--- [BACKGROUND] %s completed in %.2f seconds\n\n", name, time.Since(start).Seconds()))
		} else {
			runner.errorCh <- errorInfo{
				taskName:   name,
				errorValue: err,
			}
			fmt.Print(libformat.RedText("\n<--- [BACKGROUND] %s failed in %.2f seconds: %v\n\n", name, time.Since(start).Seconds(), err))
		}
	}()

	return future
}

// AddVoid starts a background task that doesn't return a value.
// Returns a Future[struct{}] for consistency.
func AddVoid(runner *TypedBackgroundRunner, name string, fn func() error) *Future[struct{}] {
	return Add(runner, name, func() (struct{}, error) {
		return struct{}{}, fn()
	})
}

// Wait blocks until all background tasks complete.
// Returns an error if any task panicked or errored. Individual task errors can also be checked via Future.Get()
func (tbr *TypedBackgroundRunner) Wait() error {
	// Close panic channel when all goroutines finish
	go func() {
		tbr.wg.Wait()
		close(tbr.panicCh)
		close(tbr.errorCh)
	}()

	// Check for panics
	for panicInfo := range tbr.panicCh {
		// Print the original stack trace from the background goroutine
		if panicInfo.panicStack != nil {
			fmt.Fprintf(os.Stderr, "Original panic stack trace from background task '%s':\n%s\n", panicInfo.taskName, panicInfo.panicStack)
		}
		// we want to surface panic as soon as possible to avoid executing operations, which are doomed to fail anyway
		// since we assume all tasks must execute successfully
		panic(fmt.Sprintf("Background task '%s' panicked: %v", panicInfo.taskName, panicInfo.panicValue))
	}

	// Check for errors
	var err error
	for errorInfo := range tbr.errorCh {
		if errorInfo.errorValue != nil {
			err = multierror.Append(err, errors.Wrapf(errorInfo.errorValue, "Background task '%s' failed", errorInfo.taskName))
		}
	}

	return err
}

// GetTaskNames returns the names of all tasks that have been added
func (tbr *TypedBackgroundRunner) GetTaskNames() []string {
	tbr.taskNamesMutex.Lock()
	defer tbr.taskNamesMutex.Unlock()
	names := make([]string, len(tbr.taskNames))
	copy(names, tbr.taskNames)
	return names
}
