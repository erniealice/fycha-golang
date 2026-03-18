package eventbus

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// MemoryEventBus is an in-process, synchronous event bus implementation.
// All handlers for an event type are invoked in registration order.
// This is the Phase 9 default; replace with a durable bus (e.g. outbox pattern)
// for production use cases that need delivery guarantees.
type MemoryEventBus struct {
	handlers map[string][]Handler
	mu       sync.RWMutex
}

// NewMemoryEventBus creates a new MemoryEventBus with no registered handlers.
func NewMemoryEventBus() *MemoryEventBus {
	return &MemoryEventBus{
		handlers: make(map[string][]Handler),
	}
}

// Subscribe registers a handler for a specific event type.
// Thread-safe; multiple handlers per event type are supported.
func (b *MemoryEventBus) Subscribe(eventType string, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

// Publish dispatches an event synchronously to all registered handlers.
// If Event.Timestamp is zero, it is set to the current time before dispatch.
// Returns a combined error if any handler fails; all handlers are invoked
// regardless of individual failures.
func (b *MemoryEventBus) Publish(ctx context.Context, event Event) error {
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	b.mu.RLock()
	handlers := b.handlers[event.Type]
	b.mu.RUnlock()

	var errs []error
	for _, h := range handlers {
		if err := h(ctx, event); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		return nil
	}
	if len(errs) == 1 {
		return errs[0]
	}
	return fmt.Errorf("eventbus: %d handler(s) failed for %q: first error: %w", len(errs), event.Type, errs[0])
}
