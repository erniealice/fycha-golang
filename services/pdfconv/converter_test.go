package pdfconv

import (
	"bytes"
	"context"
	"errors"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// validDocx only needs the zip magic: the fakes never parse it.
var validDocx = []byte("PK\x03\x04fake-docx-body")

// testDeps returns isolated deps (own limiter + fault tracker) around run.
func testDeps(t *testing.T, run func(context.Context, commandSpec) (string, error)) conversionDeps {
	t.Helper()
	return conversionDeps{
		findBinary:   func() (string, error) { return "/test/soffice", nil },
		runCommand:   run,
		timeout:      time.Second,
		queueTimeout: time.Second,
		retryBackoff: time.Millisecond,
		limiter:      newLimiter(1),
		faults:       &faultTracker{},
		tempRoot:     t.TempDir(),
	}
}

// outDirOf returns the --outdir argument of a spec.
func outDirOf(spec commandSpec) string {
	i := slices.Index(spec.args, "--outdir")
	return spec.args[i+1]
}

func writePDF(spec commandSpec) error {
	return os.WriteFile(filepath.Join(outDirOf(spec), "input.pdf"), []byte("%PDF-test"), 0o600)
}

func envValue(env []string, key string) (string, int) {
	var val string
	count := 0
	for _, kv := range env {
		if k, v, _ := strings.Cut(kv, "="); k == key {
			val = v
			count++
		}
	}
	return val, count
}

func TestConvert_SuccessIsolatesProfileHomeAndTmp(t *testing.T) {
	t.Parallel()

	var got commandSpec
	deps := testDeps(t, func(_ context.Context, spec commandSpec) (string, error) {
		got = spec
		return "convert input.docx -> input.pdf using filter : writer_pdf_Export", writePDF(spec)
	})

	pdf, ok, err := convertDocxToPDFWithDeps(context.Background(), validDocx, deps)
	if err != nil || !ok {
		t.Fatalf("convert = (ok=%v, err=%v), want success", ok, err)
	}
	if !bytes.Equal(pdf, []byte("%PDF-test")) {
		t.Fatalf("pdf = %q, want %%PDF-test", pdf)
	}
	if got.binary != "/test/soffice" {
		t.Fatalf("binary = %q", got.binary)
	}

	workDir := outDirOf(got)
	if filepath.Dir(workDir) != deps.tempRoot {
		t.Fatalf("workdir %q not under temp root %q", workDir, deps.tempRoot)
	}
	profileURL := (&url.URL{Scheme: "file", Path: filepath.ToSlash(filepath.Join(workDir, "libreoffice-profile"))}).String()
	wantArgs := []string{
		"-env:UserInstallation=" + profileURL,
		"--headless", "--nologo", "--nodefault", "--nofirststartwizard", "--nolockcheck", "--norestore",
		"--convert-to", "pdf",
		"--outdir", workDir,
		filepath.Join(workDir, "input.docx"),
	}
	if !reflect.DeepEqual(got.args, wantArgs) {
		t.Fatalf("args = %#v\nwant %#v", got.args, wantArgs)
	}
	for key, want := range map[string]string{
		"HOME":   filepath.Join(workDir, "home"),
		"TMPDIR": filepath.Join(workDir, "tmp"),
	} {
		if val, n := envValue(got.env, key); val != want || n != 1 {
			t.Fatalf("env %s = %q (x%d), want exactly one %q", key, val, n, want)
		}
	}
	if _, err := os.Stat(workDir); !os.IsNotExist(err) {
		t.Fatalf("workdir %q not removed after conversion (stat err=%v)", workDir, err)
	}
}

func TestConvert_UnavailableFallsBack(t *testing.T) {
	t.Parallel()

	runCalled := false
	deps := testDeps(t, func(context.Context, commandSpec) (string, error) {
		runCalled = true
		return "", nil
	})
	deps.findBinary = func() (string, error) { return "", errors.New("not installed") }

	got, ok, err := convertDocxToPDFWithDeps(context.Background(), []byte("original-docx"), deps)
	if err != nil || ok {
		t.Fatalf("convert = (ok=%v, err=%v), want fallback", ok, err)
	}
	if !bytes.Equal(got, []byte("original-docx")) {
		t.Fatalf("fallback bytes = %q", got)
	}
	if runCalled {
		t.Fatal("runCommand called after unavailable binary")
	}
}

func TestConvert_InvalidInputRejectedWithoutRun(t *testing.T) {
	t.Parallel()

	for _, input := range [][]byte{nil, []byte("docx"), append([]byte{0xef, 0xbb, 0xbf}, validDocx...)} {
		runCalled := false
		deps := testDeps(t, func(context.Context, commandSpec) (string, error) {
			runCalled = true
			return "", nil
		})
		pdf, ok, err := convertDocxToPDFWithDeps(context.Background(), input, deps)
		if !errors.Is(err, ErrInvalidInput) || ok || pdf != nil {
			t.Fatalf("input %q: result = (%q, %v, %v), want ErrInvalidInput", input, pdf, ok, err)
		}
		if runCalled {
			t.Fatalf("input %q: soffice spawned for invalid input", input)
		}
	}
}

func TestConvert_CommandFailureNotRetried(t *testing.T) {
	t.Parallel()

	var calls atomic.Int32
	deps := testDeps(t, func(context.Context, commandSpec) (string, error) {
		calls.Add(1)
		return "Error: some filter blew up", errors.New("exit status 1")
	})

	pdf, ok, err := convertDocxToPDFWithDeps(context.Background(), validDocx, deps)
	if !errors.Is(err, ErrConversionFailed) || ok || pdf != nil {
		t.Fatalf("result = (%q, %v, %v), want ErrConversionFailed", pdf, ok, err)
	}
	if !strings.Contains(err.Error(), "some filter blew up") {
		t.Fatalf("error %q does not carry captured output", err)
	}
	if calls.Load() != 1 {
		t.Fatalf("attempts = %d, want 1 (no retry for a generic failure)", calls.Load())
	}
	if h := deps.faults.snapshot(deps.limiter); h.ConsecutiveInstanceFaults != 0 || h.Failures != 1 {
		t.Fatalf("health = %+v, want 1 failure and 0 instance faults", h)
	}
}

func TestConvert_TimeoutNotRetried(t *testing.T) {
	t.Parallel()

	var calls atomic.Int32
	deps := testDeps(t, func(ctx context.Context, _ commandSpec) (string, error) {
		calls.Add(1)
		<-ctx.Done()
		return "", errors.New("signal: killed")
	})
	deps.timeout = 10 * time.Millisecond

	_, ok, err := convertDocxToPDFWithDeps(context.Background(), validDocx, deps)
	if !errors.Is(err, ErrTimeout) || ok {
		t.Fatalf("err = %v, want ErrTimeout", err)
	}
	if calls.Load() != 1 {
		t.Fatalf("attempts = %d, want 1", calls.Load())
	}
}

func TestConvert_CallerCancelIsNotTimeout(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	deps := testDeps(t, func(runCtx context.Context, _ commandSpec) (string, error) {
		cancel()
		<-runCtx.Done()
		return "", errors.New("signal: killed")
	})

	_, _, err := convertDocxToPDFWithDeps(ctx, validDocx, deps)
	if !errors.Is(err, context.Canceled) || errors.Is(err, ErrTimeout) {
		t.Fatalf("err = %v, want context.Canceled and not ErrTimeout", err)
	}
}

func TestConvert_CannotForkRetriesOnceAndRecovers(t *testing.T) {
	t.Parallel()

	var workDirs []string
	deps := testDeps(t, func(_ context.Context, spec commandSpec) (string, error) {
		workDirs = append(workDirs, outDirOf(spec))
		if len(workDirs) == 1 {
			return "/usr/bin/libreoffice: 49: Cannot fork", errors.New("exit status 2")
		}
		return "", writePDF(spec)
	})
	deps.faults.consecutive.Store(1) // a prior fault; success must reset it

	pdf, ok, err := convertDocxToPDFWithDeps(context.Background(), validDocx, deps)
	if err != nil || !ok || !bytes.HasPrefix(pdf, pdfMagic) {
		t.Fatalf("result = (%q, %v, %v), want recovered PDF", pdf, ok, err)
	}
	if len(workDirs) != 2 || workDirs[0] == workDirs[1] {
		t.Fatalf("attempt workdirs = %q, want 2 distinct", workDirs)
	}
	if h := deps.faults.snapshot(deps.limiter); !h.Healthy || h.ConsecutiveInstanceFaults != 0 {
		t.Fatalf("health = %+v, want healthy with the fault run reset", h)
	}
}

func TestConvert_SourceNotLoadedExhaustsRetryAndMarksUnhealthy(t *testing.T) {
	t.Parallel()

	var calls atomic.Int32
	deps := testDeps(t, func(context.Context, commandSpec) (string, error) {
		calls.Add(1)
		return "Warning: failed to launch javaldx\nError: source file could not be loaded\n", nil
	})

	_, ok, err := convertDocxToPDFWithDeps(context.Background(), validDocx, deps)
	if !errors.Is(err, ErrSourceNotLoaded) || ok {
		t.Fatalf("err = %v, want ErrSourceNotLoaded", err)
	}
	if calls.Load() != maxAttempts {
		t.Fatalf("attempts = %d, want %d", calls.Load(), maxAttempts)
	}
	if h := deps.faults.snapshot(deps.limiter); !h.Healthy || h.ConsecutiveInstanceFaults != 1 {
		t.Fatalf("after 1 conversion health = %+v, want still healthy with 1 consecutive fault", h)
	}
	_, _, _ = convertDocxToPDFWithDeps(context.Background(), validDocx, deps)
	if h := deps.faults.snapshot(deps.limiter); h.Healthy || h.ConsecutiveInstanceFaults != 2 {
		t.Fatalf("after 2 conversions health = %+v, want unhealthy", h)
	}
}

func TestConvert_MissingOrInvalidOutput(t *testing.T) {
	t.Parallel()

	for name, run := range map[string]func(context.Context, commandSpec) (string, error){
		"missing": func(context.Context, commandSpec) (string, error) { return "", nil },
		"not-pdf": func(_ context.Context, spec commandSpec) (string, error) {
			return "", os.WriteFile(filepath.Join(outDirOf(spec), "input.pdf"), []byte("<html>"), 0o600)
		},
	} {
		_, ok, err := convertDocxToPDFWithDeps(context.Background(), validDocx, testDeps(t, run))
		if !errors.Is(err, ErrNoOutput) || ok {
			t.Fatalf("%s: err = %v, want ErrNoOutput", name, err)
		}
	}
}

func TestConvert_WaitDelayAfterCleanExitUsesOutput(t *testing.T) {
	t.Parallel()

	deps := testDeps(t, func(_ context.Context, spec commandSpec) (string, error) {
		return "", errors.Join(writePDF(spec), exec.ErrWaitDelay)
	})
	pdf, ok, err := convertDocxToPDFWithDeps(context.Background(), validDocx, deps)
	if err != nil || !ok || !bytes.HasPrefix(pdf, pdfMagic) {
		t.Fatalf("result = (%q, %v, %v), want the produced PDF", pdf, ok, err)
	}
}

func TestIsolatedEnvReplacesInheritedValues(t *testing.T) {
	t.Parallel()

	env := isolatedEnv([]string{"HOME=/root", "TMPDIR=/tmp", "PATH=/bin", "Temp=C:\\x"}, "/w/home", "/w/tmp")
	for key, want := range map[string]string{"HOME": "/w/home", "TMPDIR": "/w/tmp", "PATH": "/bin"} {
		if val, n := envValue(env, key); val != want || n != 1 {
			t.Fatalf("%s = %q (x%d), want %q", key, val, n, want)
		}
	}
	if _, n := envValue(env, "Temp"); n != 0 {
		t.Fatal("case-variant Temp not stripped")
	}
}

func TestCappedBufferTruncates(t *testing.T) {
	t.Parallel()

	c := &cappedBuffer{limit: 4}
	if n, _ := c.Write([]byte("abcdef")); n != 6 {
		t.Fatalf("Write reported %d, want full length", n)
	}
	if got := c.String(); got != "abcd…(truncated)" {
		t.Fatalf("String = %q", got)
	}
}
