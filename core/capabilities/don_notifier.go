package capabilities

import (
	"context"
	"fmt"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

type DonNotifier interface {
	// Notify all subscribers of a new DON without blocking for a subscriber.
	NotifyDonSet(don capabilities.DON)
}

type DonNotifyWaitSubscriber interface {
	DonNotifier

	// Block until a new DON is received or the context is canceled.  The current
	// DON, if set, is returned immediately.
	WaitForDon(ctx context.Context) (capabilities.DON, error)

	// Subscribe returns a channel that will receive the latest DON.  Unsubscribe
	// by calling the returned function.
	Subscribe(ctx context.Context) (<-chan capabilities.DON, func(), error)
}

type donNotifier struct {
	lggr        logger.Logger
	mu          sync.Mutex
	don         *capabilities.DON
	subscribers map[chan capabilities.DON]struct{}
}

// TODO: wire in logger
func NewDonNotifier() *donNotifier {
	return &donNotifier{
		subscribers: make(map[chan capabilities.DON]struct{}),
	}
}

func (n *donNotifier) NotifyDonSet(don capabilities.DON) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.don = &don

	// Broadcast the new DON to all subscriber channels.
	for subCh := range n.subscribers {
		select {
		case subCh <- don:
		default:
		}
	}
}

func (n *donNotifier) Subscribe(ctx context.Context) (<-chan capabilities.DON, func(), error) {
	if ctx.Err() != nil {
		return nil, nil, ctx.Err()
	}

	subCh := make(chan capabilities.DON, 1)
	unsubscribe := func() {
		n.mu.Lock()
		defer n.mu.Unlock()
		if _, ok := n.subscribers[subCh]; ok {
			delete(n.subscribers, subCh)
			close(subCh)
		}
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	n.subscribers[subCh] = struct{}{}
	if n.don != nil {
		// Buffered so as not to block.
		subCh <- *n.don
	}

	return subCh, unsubscribe, nil
}

func (n *donNotifier) WaitForDon(ctx context.Context) (capabilities.DON, error) {
	n.mu.Lock()
	if n.don != nil {
		return *n.don, nil
	}
	n.mu.Unlock()

	subCh, unsubscribe, err := n.Subscribe(ctx)
	if err != nil {
		return capabilities.DON{}, fmt.Errorf("failed to subscribe to DON updates: %w", err)
	}
	defer unsubscribe()

	select {
	case <-ctx.Done():
		return capabilities.DON{}, fmt.Errorf("failed to wait for don: %w", ctx.Err())
	case don := <-subCh:
		return don, nil
	}
}
