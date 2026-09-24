package pdfconv

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestLimiter_SerializesConversions(t *testing.T) {
	t.Parallel()

	var active, peak atomic.Int32
	deps := testDeps(t, func(_ context.Context, spec commandSpec) (string, error) {
		n := active.Add(1)
		for {
			p := peak.Load()
			if n <= p || peak.CompareAndSwap(p, n) {
				break
			}
		}
		time.Sleep(20 * time.Millisecond)
		active.Add(-1)
		return "", writePDF(spec)
	})

	var wg sync.WaitGroup
	errs := make(chan error, 4)
	for range 4 {
		wg.Go(func() {
			if _, _, err := convertDocxToPDFWithDeps(context.Background(), validDocx, deps); err != nil {
				errs <- err
			}
		})
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Fatalf("conversion error: %v", err)
	}
	if peak.Load() != 1 {
		t.Fatalf("peak concurrent soffice = %d, want 1", peak.Load())
	}
	if deps.limiter.inFlight() != 0 {
		t.Fatalf("slots leaked: in flight = %d", deps.limiter.inFlight())
	}
}

func TestLimiter_QueueTimeoutAndCancel(t *testing.T) {
	t.Parallel()

	l := newLimiter(1)
	release, err := l.acquire(context.Background(), time.Second)
	if err != nil {
		t.Fatalf("first acquire: %v", err)
	}
	defer release()

	if _, err := l.acquire(context.Background(), 10*time.Millisecond); !errors.Is(err, ErrBusy) {
		t.Fatalf("timeout acquire err = %v, want ErrBusy", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := l.acquire(ctx, time.Second); !errors.Is(err, ErrBusy) || !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled acquire err = %v, want ErrBusy wrapping context.Canceled", err)
	}
}

func TestLimiterCapacityFromEnv(t *testing.T) {
	t.Parallel()

	for raw, want := range map[string]int{"": 1, " 3 ": 3, "0": 1, "-2": 1, "99": maxConcurrencyCeiling, "two": 1} {
		if got := capacityFromEnv(raw); got != want {
			t.Errorf("capacityFromEnv(%q) = %d, want %d", raw, got, want)
		}
	}
}
