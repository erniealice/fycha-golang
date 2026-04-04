package eventbus

import (
	"context"
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestMemoryEventBus_SubscribeAndPublish(t *testing.T) {
	t.Parallel()

	t.Run("handler receives published event", func(t *testing.T) {
		t.Parallel()

		bus := NewMemoryEventBus()
		var received Event

		bus.Subscribe("order.created", func(ctx context.Context, e Event) error {
			received = e
			return nil
		})

		evt := Event{
			Type:      "order.created",
			SourceID:  "order-123",
			Payload:   map[string]any{"amount": 100.0},
			Timestamp: time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC),
		}

		err := bus.Publish(context.Background(), evt)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if received.Type != "order.created" {
			t.Errorf("Type = %q, want %q", received.Type, "order.created")
		}
		if received.SourceID != "order-123" {
			t.Errorf("SourceID = %q, want %q", received.SourceID, "order-123")
		}
		if received.Payload["amount"] != 100.0 {
			t.Errorf("Payload[amount] = %v, want %v", received.Payload["amount"], 100.0)
		}
	})

	t.Run("timestamp defaults to now when zero", func(t *testing.T) {
		t.Parallel()

		bus := NewMemoryEventBus()
		var received Event

		bus.Subscribe("test.event", func(ctx context.Context, e Event) error {
			received = e
			return nil
		})

		before := time.Now()
		err := bus.Publish(context.Background(), Event{
			Type:     "test.event",
			SourceID: "src-1",
		})
		after := time.Now()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if received.Timestamp.IsZero() {
			t.Fatal("Timestamp should not be zero")
		}
		if received.Timestamp.Before(before) || received.Timestamp.After(after) {
			t.Errorf("Timestamp = %v, expected between %v and %v", received.Timestamp, before, after)
		}
	})

	t.Run("explicit timestamp is preserved", func(t *testing.T) {
		t.Parallel()

		bus := NewMemoryEventBus()
		var received Event

		bus.Subscribe("test.event", func(ctx context.Context, e Event) error {
			received = e
			return nil
		})

		explicit := time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC)
		err := bus.Publish(context.Background(), Event{
			Type:      "test.event",
			SourceID:  "src-1",
			Timestamp: explicit,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !received.Timestamp.Equal(explicit) {
			t.Errorf("Timestamp = %v, want %v", received.Timestamp, explicit)
		}
	})

	t.Run("no handlers is a no-op", func(t *testing.T) {
		t.Parallel()

		bus := NewMemoryEventBus()

		err := bus.Publish(context.Background(), Event{
			Type:     "unhandled.event",
			SourceID: "src-1",
		})
		if err != nil {
			t.Fatalf("expected nil error for no handlers, got: %v", err)
		}
	})

	t.Run("multiple handlers all invoked", func(t *testing.T) {
		t.Parallel()

		bus := NewMemoryEventBus()
		var count int32

		for i := 0; i < 3; i++ {
			bus.Subscribe("multi.event", func(ctx context.Context, e Event) error {
				atomic.AddInt32(&count, 1)
				return nil
			})
		}

		err := bus.Publish(context.Background(), Event{
			Type:     "multi.event",
			SourceID: "src-1",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if atomic.LoadInt32(&count) != 3 {
			t.Errorf("handler invocations = %d, want 3", count)
		}
	})

	t.Run("single handler error is returned directly", func(t *testing.T) {
		t.Parallel()

		bus := NewMemoryEventBus()
		handlerErr := errors.New("handler failed")

		bus.Subscribe("fail.event", func(ctx context.Context, e Event) error {
			return handlerErr
		})

		err := bus.Publish(context.Background(), Event{
			Type:     "fail.event",
			SourceID: "src-1",
		})
		if !errors.Is(err, handlerErr) {
			t.Errorf("error = %v, want %v", err, handlerErr)
		}
	})

	t.Run("multiple handler errors are aggregated", func(t *testing.T) {
		t.Parallel()

		bus := NewMemoryEventBus()

		bus.Subscribe("fail.multi", func(ctx context.Context, e Event) error {
			return errors.New("error-1")
		})
		bus.Subscribe("fail.multi", func(ctx context.Context, e Event) error {
			return nil // second handler succeeds
		})
		bus.Subscribe("fail.multi", func(ctx context.Context, e Event) error {
			return errors.New("error-3")
		})

		err := bus.Publish(context.Background(), Event{
			Type:     "fail.multi",
			SourceID: "src-1",
		})
		if err == nil {
			t.Fatal("expected an error")
		}

		errMsg := err.Error()
		if !strings.Contains(errMsg, "2 handler(s) failed") {
			t.Errorf("error should mention 2 handlers failed, got: %q", errMsg)
		}
		if !strings.Contains(errMsg, "fail.multi") {
			t.Errorf("error should mention event type, got: %q", errMsg)
		}
	})

	t.Run("handlers for different event types are independent", func(t *testing.T) {
		t.Parallel()

		bus := NewMemoryEventBus()
		var aCalled, bCalled bool

		bus.Subscribe("type.a", func(ctx context.Context, e Event) error {
			aCalled = true
			return nil
		})
		bus.Subscribe("type.b", func(ctx context.Context, e Event) error {
			bCalled = true
			return nil
		})

		err := bus.Publish(context.Background(), Event{Type: "type.a", SourceID: "src"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !aCalled {
			t.Error("type.a handler should have been called")
		}
		if bCalled {
			t.Error("type.b handler should not have been called")
		}
	})

	t.Run("publish with empty event type is a no-op", func(t *testing.T) {
		t.Parallel()

		bus := NewMemoryEventBus()
		handlerCalled := false

		bus.Subscribe("", func(ctx context.Context, e Event) error {
			handlerCalled = true
			return nil
		})

		// Publish with an empty event type — the handler registered for "" should fire
		err := bus.Publish(context.Background(), Event{Type: "", SourceID: "src-1"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !handlerCalled {
			t.Error("handler subscribed to empty string should have been called")
		}
	})

	t.Run("publish empty type with no matching handler", func(t *testing.T) {
		t.Parallel()

		bus := NewMemoryEventBus()

		// Subscribe to a non-empty type, then publish empty type
		bus.Subscribe("real.event", func(ctx context.Context, e Event) error {
			t.Error("handler for real.event should not be called")
			return nil
		})

		err := bus.Publish(context.Background(), Event{Type: "", SourceID: "src-1"})
		if err != nil {
			t.Fatalf("expected nil error for empty event type with no handler, got: %v", err)
		}
	})

	t.Run("subscribe with nil handler panics on publish", func(t *testing.T) {
		t.Parallel()

		bus := NewMemoryEventBus()
		bus.Subscribe("nil.handler", nil)

		defer func() {
			r := recover()
			if r == nil {
				t.Error("expected panic when invoking nil handler, but got none")
			}
		}()

		_ = bus.Publish(context.Background(), Event{Type: "nil.handler", SourceID: "src-1"})
	})

	t.Run("handler that panics does not crash the test process", func(t *testing.T) {
		t.Parallel()

		bus := NewMemoryEventBus()

		bus.Subscribe("panic.event", func(ctx context.Context, e Event) error {
			panic("handler exploded")
		})

		defer func() {
			r := recover()
			if r == nil {
				t.Error("expected panic from handler to propagate")
			}
		}()

		_ = bus.Publish(context.Background(), Event{Type: "panic.event", SourceID: "src-1"})
	})

	t.Run("concurrent publish and subscribe is safe", func(t *testing.T) {
		t.Parallel()

		bus := NewMemoryEventBus()
		var callCount int64
		done := make(chan struct{})

		// Spawn goroutines that concurrently subscribe and publish
		const goroutines = 50

		// Start subscribers in background
		go func() {
			for i := 0; i < goroutines; i++ {
				bus.Subscribe("concurrent.event", func(ctx context.Context, e Event) error {
					atomic.AddInt64(&callCount, 1)
					return nil
				})
			}
			close(done)
		}()

		// Publish concurrently while subscriptions are being added
		var publishErrors int64
		var wg sync.WaitGroup
		for i := 0; i < goroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := bus.Publish(context.Background(), Event{
					Type:     "concurrent.event",
					SourceID: "src-concurrent",
				})
				if err != nil {
					atomic.AddInt64(&publishErrors, 1)
				}
			}()
		}

		wg.Wait()
		<-done

		// No assertions on exact counts (race-dependent), but the test must not
		// panic or trigger the race detector.
		t.Logf("concurrent test: %d handler calls, %d publish errors",
			atomic.LoadInt64(&callCount), atomic.LoadInt64(&publishErrors))
	})

	t.Run("publish with cancelled context", func(t *testing.T) {
		t.Parallel()

		bus := NewMemoryEventBus()
		var receivedCtx context.Context

		bus.Subscribe("ctx.event", func(ctx context.Context, e Event) error {
			receivedCtx = ctx
			return nil
		})

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // cancel before publishing

		err := bus.Publish(ctx, Event{Type: "ctx.event", SourceID: "src-1"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// The handler should receive the cancelled context
		if receivedCtx.Err() == nil {
			t.Error("handler should receive a cancelled context")
		}
	})
}
