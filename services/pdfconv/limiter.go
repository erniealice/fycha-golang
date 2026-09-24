package pdfconv

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	// maxConcurrencyEnv bounds how many LibreOffice processes this process runs
	// at once. Each soffice tree costs roughly 150–300 MB RSS plus tens of
	// threads, and it shares the container with the whole app server, so the
	// default is 1 (rc-docx-pdf-1 "Track B": two concurrent trees exhausted the
	// fork() headroom). Values are clamped to [1, maxConcurrencyCeiling].
	maxConcurrencyEnv     = "PDFCONV_MAX_CONCURRENCY"
	defaultMaxConcurrency = 1
	maxConcurrencyCeiling = 8

	// defaultQueueTimeout caps how long a caller waits for a free slot before
	// getting ErrBusy, so a burst of report-card exports degrades into fast
	// errors instead of piling up goroutines behind a stuck conversion.
	defaultQueueTimeout = 90 * time.Second
)

// limiter is a counting semaphore over LibreOffice processes.
type limiter struct {
	slots chan struct{}
}

func newLimiter(capacity int) *limiter {
	if capacity < 1 {
		capacity = 1
	}
	return &limiter{slots: make(chan struct{}, capacity)}
}

// acquire blocks until a slot is free, the caller's ctx ends, or timeout
// elapses. Failures wrap ErrBusy (and the ctx error when the caller cancelled).
func (l *limiter) acquire(ctx context.Context, timeout time.Duration) (release func(), err error) {
	select {
	case l.slots <- struct{}{}:
		return l.release, nil
	default:
	}
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case l.slots <- struct{}{}:
		return l.release, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("pdfconv: %w while queued: %w", ErrBusy, ctx.Err())
	case <-timer.C:
		return nil, fmt.Errorf("pdfconv: %w: no conversion slot within %s", ErrBusy, timeout)
	}
}

func (l *limiter) release() { <-l.slots }

// inFlight reports how many slots are currently held.
func (l *limiter) inFlight() int { return len(l.slots) }

// capacityFromEnv parses PDFCONV_MAX_CONCURRENCY: empty → default, garbage →
// default (logged), otherwise clamped to [1, maxConcurrencyCeiling].
func capacityFromEnv(raw string) int {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return defaultMaxConcurrency
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		log.Printf("pdfconv: ignoring invalid %s=%q; using %d", maxConcurrencyEnv, raw, defaultMaxConcurrency)
		return defaultMaxConcurrency
	}
	return min(max(n, 1), maxConcurrencyCeiling)
}

// defaultLimiter is shared by every production conversion in this process.
var defaultLimiter = sync.OnceValue(func() *limiter {
	return newLimiter(capacityFromEnv(os.Getenv(maxConcurrencyEnv)))
})
