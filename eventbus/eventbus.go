// Package eventbus provides the event bus infrastructure for fycha.
// Phase 9: Auto-posting journal entries from operational domain events.
//
// Event flow:
//   Operational domain (revenue, expenditure, treasury, payroll, etc.)
//     → publishes Event to EventBus
//     → JournalPoster handler maps event type to debit/credit rules
//     → calls CreateJournalEntry use case
//     → journal entry is created as draft or posted automatically
package eventbus

import (
	"context"
	"time"
)

// Event represents a domain event published by an operational module.
// The Type field drives the journal posting logic in JournalPoster.
type Event struct {
	// Type is the event type string, e.g. "revenue.completed".
	// See the full list of recognized types in journal_poster.go.
	Type string

	// SourceID is the ID of the originating entity (revenue ID, loan ID, etc.).
	SourceID string

	// Payload carries event-specific data used to populate journal lines.
	// Keys are domain-specific; the journal poster documents expected keys
	// in each handler stub.
	Payload map[string]any

	// Timestamp is when the event occurred. Defaults to now if zero.
	Timestamp time.Time
}

// Handler is a function that processes a domain event.
// Handlers must be idempotent where possible — the event bus may retry
// failed deliveries in future implementations.
type Handler func(ctx context.Context, event Event) error

// EventBus is the interface for publishing and subscribing to domain events.
// Phase 9 uses MemoryEventBus; a persistent implementation can replace it.
type EventBus interface {
	// Publish sends an event to all registered handlers for the event type.
	// Returns the first handler error encountered; does not short-circuit.
	Publish(ctx context.Context, event Event) error

	// Subscribe registers a handler for a specific event type.
	// Multiple handlers may be registered for the same event type.
	Subscribe(eventType string, handler Handler)
}
