package pdfconv

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// LibreOffice startup on the supported macOS development lane can take around
// 30 seconds before it starts laying out a generated document. Grade-sheet
// roster tables add enough work to exceed the former 60-second ceiling even
// though conversion is still progressing normally. Keep the native process
// bounded, but allow a complete heavy report to finish.
const libreOfficeConversionTimeout = 120 * time.Second

const (
	// maxAttempts allows exactly one retry, and only for instance-fault
	// signatures (see isInstanceFault). Timeouts and bad input never retry.
	maxAttempts = 2
	// defaultRetryBackoff gives a transient fork() EAGAIN time to clear.
	defaultRetryBackoff = 500 * time.Millisecond
	// waitDelay bounds how long Wait blocks draining output after the process
	// exits or is killed, in case a grandchild keeps the pipe open.
	waitDelay = 10 * time.Second
	// outputCaptureLimit caps captured soffice stdout+stderr; the diagnostic
	// lines ("Error: ...", "Cannot fork") appear early.
	outputCaptureLimit = 32 << 10
	// unhealthyFaultThreshold consecutive instance faults mark the converter
	// unhealthy (rc-docx-pdf-1 recycles at the same count).
	unhealthyFaultThreshold = 2
)

// Typed conversion failures. Every returned error wraps exactly one of these
// with %w (plus the captured soffice output), so callers classify with
// errors.Is rather than string matching.
var (
	// ErrInvalidInput: the bytes are not a zip/OOXML package. A file fault,
	// never retried.
	ErrInvalidInput = errors.New("input is not a DOCX (zip) package")
	// ErrBusy: no conversion slot freed up within the queue timeout, or the
	// caller gave up while queued.
	ErrBusy = errors.New("pdf conversion capacity busy")
	// ErrTimeout: LibreOffice exceeded the conversion deadline; its whole
	// process group was killed.
	ErrTimeout = errors.New("libreoffice conversion timed out")
	// ErrCannotFork: the libreoffice wrapper could not fork(2)
	// (EAGAIN/ENOMEM). Always an instance resource-exhaustion signal
	// (typically un-reaped zombies or too many concurrent soffice trees).
	ErrCannotFork = errors.New("libreoffice cannot fork (instance resource exhaustion)")
	// ErrSourceNotLoaded: LibreOffice exited 0 but printed "source file could
	// not be loaded". The input already passed the zip check, so this is the
	// downstream face of a wedged instance, not a bad document.
	ErrSourceNotLoaded = errors.New("libreoffice could not load a valid source document")
	// ErrConversionFailed: any other non-zero LibreOffice exit.
	ErrConversionFailed = errors.New("libreoffice conversion failed")
	// ErrNoOutput: LibreOffice reported success but produced no PDF, or the
	// output is not a PDF.
	ErrNoOutput = errors.New("libreoffice produced no valid PDF")
)

var (
	zipMagic = []byte("PK\x03\x04")
	pdfMagic = []byte("%PDF-")
)

// bootID is generated once per process and stamped into every conversion log
// line, so failures cluster by instance lifetime ("one wedged instance" vs
// "systemic") at a glance.
var bootID = func() string {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return "unknown"
	}
	return hex.EncodeToString(b)
}()

// commandSpec is one LibreOffice invocation.
type commandSpec struct {
	binary string
	args   []string
	env    []string
}

type conversionDeps struct {
	findBinary func() (string, error)
	// runCommand runs the spec until exit or ctx end and returns its captured
	// (bounded) stdout+stderr.
	runCommand   func(context.Context, commandSpec) (output string, err error)
	timeout      time.Duration
	queueTimeout time.Duration
	retryBackoff time.Duration
	limiter      *limiter
	faults       *faultTracker
	// tempRoot is the parent of each per-conversion workdir ("" = os.TempDir).
	tempRoot string
}

func productionConversionDeps() conversionDeps {
	return conversionDeps{
		findBinary:   findLibreOffice,
		runCommand:   runLibreOffice,
		timeout:      libreOfficeConversionTimeout,
		queueTimeout: defaultQueueTimeout,
		retryBackoff: defaultRetryBackoff,
		limiter:      defaultLimiter(),
		faults:       defaultFaults,
	}
}

// runLibreOffice executes soffice in its own process group so a timeout or
// cancel SIGKILLs the whole tree (wrapper, oosplash, soffice.bin, helpers)
// instead of orphaning it, and captures bounded output for classification.
func runLibreOffice(ctx context.Context, spec commandSpec) (string, error) {
	cmd := exec.CommandContext(ctx, spec.binary, spec.args...)
	cmd.Env = spec.env
	out := &cappedBuffer{limit: outputCaptureLimit}
	cmd.Stdout = out
	cmd.Stderr = out
	setProcessGroup(cmd)
	cmd.Cancel = func() error { return killProcessGroup(cmd) }
	cmd.WaitDelay = waitDelay
	err := cmd.Run()
	return out.String(), err
}

// ConvertDocxToPDF converts DOCX bytes to PDF bytes using LibreOffice headless.
// It auto-detects the OS to find the LibreOffice binary.
// If LibreOffice is not installed, it returns the original DOCX bytes with a false flag.
func ConvertDocxToPDF(docxBytes []byte) (pdfBytes []byte, ok bool, err error) {
	return ConvertDocxToPDFContext(context.Background(), docxBytes)
}

// ConvertDocxToPDFContext is ConvertDocxToPDF bounded by the caller's context:
// cancelling ctx (e.g. a client disconnect) stops queueing and kills an
// in-progress LibreOffice process group.
func ConvertDocxToPDFContext(ctx context.Context, docxBytes []byte) (pdfBytes []byte, ok bool, err error) {
	return convertDocxToPDFWithDeps(ctx, docxBytes, productionConversionDeps())
}

func convertDocxToPDFWithDeps(ctx context.Context, docxBytes []byte, deps conversionDeps) (pdfBytes []byte, ok bool, err error) {
	binary, err := deps.findBinary()
	if err != nil {
		log.Printf("pdfconv: LibreOffice not found, falling back to DOCX: %v", err)
		return docxBytes, false, nil
	}
	if !bytes.HasPrefix(docxBytes, zipMagic) {
		return nil, false, fmt.Errorf("pdfconv: %w (%d bytes, head %s)", ErrInvalidInput, len(docxBytes), headHex(docxBytes))
	}

	queuedAt := time.Now()
	release, err := deps.limiter.acquire(ctx, deps.queueTimeout)
	if err != nil {
		logResult(deps, "busy", 0, len(docxBytes), 0, time.Since(queuedAt), 0)
		return nil, false, err
	}
	defer release()
	queued := time.Since(queuedAt)

	started := time.Now()
	var attempt int
	for attempt = 1; ; attempt++ {
		pdfBytes, err = convertOnce(ctx, binary, docxBytes, deps)
		if err == nil || attempt >= maxAttempts || !isInstanceFault(err) {
			break
		}
		log.Printf("pdfconv: attempt %d instance fault, retrying with a fresh workdir: %v", attempt, err)
		if !sleepCtx(ctx, deps.retryBackoff) {
			err = fmt.Errorf("pdfconv: conversion canceled before retry: %w", ctx.Err())
			break
		}
	}

	deps.faults.record(err)
	logResult(deps, resultClass(err), attempt, len(docxBytes), len(pdfBytes), queued, time.Since(started))
	if err != nil {
		return nil, false, err
	}
	return pdfBytes, true, nil
}

// convertOnce runs one LibreOffice conversion in a fresh workdir that holds the
// input, the output, and an isolated user profile, HOME, and TMPDIR, so no
// LibreOffice or fontconfig state is shared between runs (a half-written shared
// artifact from an interrupted run was a wedge candidate in rc-docx-pdf-1).
func convertOnce(ctx context.Context, binary string, docxBytes []byte, deps conversionDeps) ([]byte, error) {
	workDir, err := os.MkdirTemp(deps.tempRoot, "pdfconv-*")
	if err != nil {
		return nil, fmt.Errorf("pdfconv: creating temp dir: %w", err)
	}
	defer os.RemoveAll(workDir)

	docxPath := filepath.Join(workDir, "input.docx")
	if err := os.WriteFile(docxPath, docxBytes, 0o600); err != nil {
		return nil, fmt.Errorf("pdfconv: writing temp docx: %w", err)
	}
	homeDir := filepath.Join(workDir, "home")
	tmpDir := filepath.Join(workDir, "tmp")
	for _, dir := range []string{homeDir, tmpDir} {
		if err := os.Mkdir(dir, 0o700); err != nil {
			return nil, fmt.Errorf("pdfconv: creating %s: %w", dir, err)
		}
	}

	// Reusing the interactive/default profile can block indefinitely on a stale
	// singleton lock, and concurrent runs would otherwise attach to one office
	// process. A unique profile makes each run an independent instance.
	profileDir := filepath.Join(workDir, "libreoffice-profile")
	profileURL := (&url.URL{Scheme: "file", Path: filepath.ToSlash(profileDir)}).String()
	spec := commandSpec{
		binary: binary,
		args: []string{
			"-env:UserInstallation=" + profileURL,
			"--headless",
			"--nologo",
			"--nodefault",
			"--nofirststartwizard",
			"--nolockcheck",
			"--norestore",
			"--convert-to", "pdf",
			"--outdir", workDir,
			docxPath,
		},
		env: isolatedEnv(os.Environ(), homeDir, tmpDir),
	}

	runCtx, cancel := context.WithTimeout(ctx, deps.timeout)
	defer cancel()
	output, runErr := deps.runCommand(runCtx, spec)
	detail := outputDetail(output)
	if errors.Is(runErr, exec.ErrWaitDelay) {
		// soffice exited 0 but a lingering helper held the output pipe past
		// waitDelay; judge the run by its output like any clean exit.
		log.Printf("pdfconv: soffice exited but its output pipe stayed open past %s (lingering helper)", waitDelay)
		runErr = nil
	}

	switch {
	case ctx.Err() != nil:
		return nil, fmt.Errorf("pdfconv: conversion canceled: %w", ctx.Err())
	case errors.Is(runCtx.Err(), context.DeadlineExceeded):
		return nil, fmt.Errorf("pdfconv: %w after %s%s", ErrTimeout, deps.timeout, detail)
	case strings.Contains(output, "Cannot fork"):
		return nil, fmt.Errorf("pdfconv: %w (exit: %v)%s", ErrCannotFork, runErr, detail)
	case strings.Contains(output, "source file could not be loaded"):
		return nil, fmt.Errorf("pdfconv: %w (exit: %v)%s", ErrSourceNotLoaded, runErr, detail)
	case runErr != nil:
		return nil, fmt.Errorf("pdfconv: %w: %v%s", ErrConversionFailed, runErr, detail)
	}

	pdfBytes, err := os.ReadFile(filepath.Join(workDir, "input.pdf"))
	if err != nil {
		return nil, fmt.Errorf("pdfconv: %w: %v%s", ErrNoOutput, err, detail)
	}
	if !bytes.HasPrefix(pdfBytes, pdfMagic) {
		return nil, fmt.Errorf("pdfconv: %w: output head %s%s", ErrNoOutput, headHex(pdfBytes), detail)
	}
	return pdfBytes, nil
}

// isInstanceFault reports failures that indicate the host (not the document)
// is unhealthy. They are the only failures retried and counted toward Health.
func isInstanceFault(err error) bool {
	return errors.Is(err, ErrCannotFork) || errors.Is(err, ErrSourceNotLoaded)
}

// isolatedEnv returns base with HOME/TMPDIR (and the Windows temp vars)
// pointed into the per-run workdir.
func isolatedEnv(base []string, homeDir, tmpDir string) []string {
	env := make([]string, 0, len(base)+4)
	for _, kv := range base {
		key, _, _ := strings.Cut(kv, "=")
		switch strings.ToUpper(key) {
		case "HOME", "TMPDIR", "TMP", "TEMP":
			continue
		}
		env = append(env, kv)
	}
	return append(env, "HOME="+homeDir, "TMPDIR="+tmpDir, "TMP="+tmpDir, "TEMP="+tmpDir)
}

// Health is a point-in-time view of the converter for liveness wiring and logs.
type Health struct {
	Healthy                   bool
	ConsecutiveInstanceFaults int64
	Conversions               int64
	Failures                  int64
	InFlight                  int
}

// CurrentHealth reports whether recent conversions indicate a wedged instance
// (unhealthyFaultThreshold consecutive instance faults). An app can surface
// this on a liveness endpoint so the platform recycles the instance; pdfconv
// itself never exits the process, because it shares it with the app server.
func CurrentHealth() Health {
	return defaultFaults.snapshot(defaultLimiter())
}

// faultTracker counts outcomes; a success resets the consecutive-fault run.
type faultTracker struct {
	consecutive atomic.Int64
	conversions atomic.Int64
	failures    atomic.Int64
}

var defaultFaults = &faultTracker{}

func (f *faultTracker) record(err error) {
	f.conversions.Add(1)
	switch {
	case err == nil:
		f.consecutive.Store(0)
	case isInstanceFault(err):
		f.failures.Add(1)
		if n := f.consecutive.Add(1); n >= unhealthyFaultThreshold {
			log.Printf("pdfconv: WEDGED %d consecutive instance faults on boot=%s; valid documents are failing, so this host needs recycling: %v", n, bootID, err)
		}
	default:
		f.failures.Add(1)
	}
}

func (f *faultTracker) snapshot(l *limiter) Health {
	n := f.consecutive.Load()
	return Health{
		Healthy:                   n < unhealthyFaultThreshold,
		ConsecutiveInstanceFaults: n,
		Conversions:               f.conversions.Load(),
		Failures:                  f.failures.Load(),
		InFlight:                  l.inFlight(),
	}
}

func resultClass(err error) string {
	for _, c := range []struct {
		target error
		name   string
	}{
		{ErrTimeout, "timeout"},
		{ErrCannotFork, "cannot_fork"},
		{ErrSourceNotLoaded, "source_not_loaded"},
		{ErrConversionFailed, "failed"},
		{ErrNoOutput, "no_output"},
		{context.Canceled, "canceled"},
		{context.DeadlineExceeded, "canceled"},
	} {
		if errors.Is(err, c.target) {
			return c.name
		}
	}
	if err != nil {
		return "error"
	}
	return "ok"
}

// logResult writes one structured line per conversion, including the resource
// snapshot that exposes fork() pressure (zombies, pids, threads, commit).
func logResult(deps conversionDeps, result string, attempts, inBytes, outBytes int, queued, took time.Duration) {
	diag, zombies := diagSnapshot()
	warn := ""
	if zombies >= zombieWarnThreshold {
		warn = " WARN=zombies_accumulating(no_init_reaper?_run_under_tini)"
	}
	log.Printf("pdfconv: result=%s attempts=%d in=%dB out=%dB queued=%s took=%s in_flight=%d consecutive_faults=%d %s%s",
		result, attempts, inBytes, outBytes, queued.Round(time.Millisecond), took.Round(time.Millisecond),
		deps.limiter.inFlight(), deps.faults.consecutive.Load(), diag, warn)
}

func sleepCtx(ctx context.Context, d time.Duration) bool {
	if d <= 0 {
		return ctx.Err() == nil
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		return true
	case <-ctx.Done():
		return false
	}
}

func outputDetail(output string) string {
	output = strings.TrimSpace(output)
	if output == "" {
		return ""
	}
	return "; output: " + output
}

func headHex(b []byte) string {
	return hex.EncodeToString(b[:min(len(b), 8)])
}

// cappedBuffer keeps the first limit bytes written and discards the rest while
// still reporting full writes, so soffice never blocks on a full pipe.
type cappedBuffer struct {
	mu        sync.Mutex
	buf       bytes.Buffer
	limit     int
	truncated bool
}

func (c *cappedBuffer) Write(p []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if room := c.limit - c.buf.Len(); room > 0 {
		c.buf.Write(p[:min(len(p), room)])
		c.truncated = c.truncated || len(p) > room
	} else if len(p) > 0 {
		c.truncated = true
	}
	return len(p), nil
}

func (c *cappedBuffer) String() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.truncated {
		return c.buf.String() + "…(truncated)"
	}
	return c.buf.String()
}

// findLibreOffice locates the LibreOffice binary based on the OS.
func findLibreOffice() (string, error) {
	if runtime.GOOS == "windows" {
		// Common Windows install paths
		candidates := []string{
			`C:\Program Files\LibreOffice\program\soffice.exe`,
			`C:\Program Files (x86)\LibreOffice\program\soffice.exe`,
		}
		for _, p := range candidates {
			if _, err := os.Stat(p); err == nil {
				return p, nil
			}
		}
		return "", fmt.Errorf("LibreOffice not found at standard Windows paths")
	}

	// Linux / macOS — look for soffice or libreoffice in PATH
	for _, name := range []string{"soffice", "libreoffice"} {
		if p, err := exec.LookPath(name); err == nil {
			return p, nil
		}
	}

	// Check common install paths (the macOS app bundle is not on PATH by default)
	candidates := []string{
		"/usr/bin/soffice",
		"/usr/bin/libreoffice",
		"/usr/local/bin/soffice",
		"/usr/lib/libreoffice/program/soffice",
		"/Applications/LibreOffice.app/Contents/MacOS/soffice",
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}

	return "", fmt.Errorf("LibreOffice not found in PATH or standard locations")
}

// ReplaceExtension replaces the file extension with .pdf.
// Exported for use in callers that need to adjust filenames.
func ReplaceExtension(filename string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		return filename + ".pdf"
	}
	return strings.TrimSuffix(filename, ext) + ".pdf"
}
