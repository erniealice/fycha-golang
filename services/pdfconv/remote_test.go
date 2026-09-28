package pdfconv

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

var (
	testDocx = []byte("PK\x03\x04fake-docx-body")
	testPDF  = []byte("%PDF-1.7 fake")
)

const testKey = "test-key-0123456789abcdef0123456789"

func testRemote(t *testing.T, h http.HandlerFunc) (remoteConfig, remoteDeps) {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	cfg, ok := parseRemoteConfig(srv.URL, testKey)
	if !ok || cfg.err != nil {
		t.Fatalf("parseRemoteConfig(%q) = %+v, %v", srv.URL, cfg, ok)
	}
	return cfg, remoteDeps{client: srv.Client(), retryBackoff: time.Millisecond}
}

func writeWireError(w http.ResponseWriter, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(StatusForCode(code))
	_ = json.NewEncoder(w).Encode(ErrorBody{Error: code, Message: message})
}

func TestRemoteConfigFromEnv(t *testing.T) {
	cases := []struct {
		name, url, key string
		configured     bool
		wantErr        bool
		wantEndpoint   string
	}{
		{name: "unset", url: "", key: testKey},
		{name: "blank", url: "   ", key: testKey},
		{name: "https base", url: "https://conv.example.run.app", key: testKey, configured: true, wantEndpoint: "https://conv.example.run.app/v1/convert"},
		{name: "https trailing slash", url: "https://conv.example.run.app/", key: testKey, configured: true, wantEndpoint: "https://conv.example.run.app/v1/convert"},
		{name: "full path kept", url: "https://conv.example.run.app/v1/convert", key: testKey, configured: true, wantEndpoint: "https://conv.example.run.app/v1/convert"},
		{name: "http loopback", url: "http://127.0.0.1:8095", key: testKey, configured: true, wantEndpoint: "http://127.0.0.1:8095/v1/convert"},
		{name: "http localhost", url: "http://localhost:8095", key: testKey, configured: true, wantEndpoint: "http://localhost:8095/v1/convert"},
		{name: "http remote refused", url: "http://conv.example.com", key: testKey, configured: true, wantErr: true},
		{name: "missing key", url: "https://conv.example.run.app", key: " ", configured: true, wantErr: true},
		{name: "userinfo refused", url: "https://u:p@conv.example.run.app", key: testKey, configured: true, wantErr: true},
		{name: "query refused", url: "https://conv.example.run.app?key=x", key: testKey, configured: true, wantErr: true},
		{name: "relative refused", url: "conv.example.run.app", key: testKey, configured: true, wantErr: true},
		{name: "other scheme refused", url: "ftp://conv.example.run.app", key: testKey, configured: true, wantErr: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv(RemoteURLEnv, tc.url)
			t.Setenv(RemoteAPIKeyEnv, tc.key)
			cfg, configured := remoteConfigFromEnv()
			if configured != tc.configured {
				t.Fatalf("configured = %v, want %v", configured, tc.configured)
			}
			if tc.wantErr {
				if !errors.Is(cfg.err, ErrRemoteMisconfigured) {
					t.Fatalf("err = %v, want ErrRemoteMisconfigured", cfg.err)
				}
				if strings.Contains(cfg.err.Error(), "u:p") {
					t.Fatalf("error leaks credentials: %v", cfg.err)
				}
				return
			}
			if cfg.err != nil {
				t.Fatalf("unexpected err: %v", cfg.err)
			}
			if cfg.endpoint != tc.wantEndpoint {
				t.Fatalf("endpoint = %q, want %q", cfg.endpoint, tc.wantEndpoint)
			}
		})
	}
}

func TestConvertRemoteMisconfiguredSendsNothing(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { calls.Add(1) }))
	defer srv.Close()
	t.Setenv(RemoteURLEnv, srv.URL)
	t.Setenv(RemoteAPIKeyEnv, "")
	_, ok, err := ConvertDocxToPDFContext(context.Background(), testDocx)
	if ok || !errors.Is(err, ErrRemoteMisconfigured) {
		t.Fatalf("ok=%v err=%v, want ErrRemoteMisconfigured", ok, err)
	}
	if calls.Load() != 0 {
		t.Fatalf("server called %d times", calls.Load())
	}
}

func TestConvertRemoteSuccess(t *testing.T) {
	cfg, deps := testRemote(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != RemoteConvertPath {
			t.Errorf("got %s %s", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer "+testKey {
			t.Errorf("Authorization = %q", got)
		}
		if got := r.Header.Get("Content-Type"); got != DocxContentType {
			t.Errorf("Content-Type = %q", got)
		}
		if r.Header.Get(RequestIDHeader) == "" {
			t.Error("missing request id")
		}
		body, _ := io.ReadAll(r.Body)
		if string(body) != string(testDocx) {
			t.Errorf("body = %q", body)
		}
		w.Header().Set("Content-Type", "application/pdf")
		_, _ = w.Write(testPDF)
	})
	pdf, ok, err := convertRemote(context.Background(), testDocx, cfg, deps)
	if err != nil || !ok || string(pdf) != string(testPDF) {
		t.Fatalf("pdf=%q ok=%v err=%v", pdf, ok, err)
	}
}

func TestConvertRemoteInvalidInputNoRequest(t *testing.T) {
	var calls atomic.Int32
	cfg, deps := testRemote(t, func(http.ResponseWriter, *http.Request) { calls.Add(1) })
	_, ok, err := convertRemote(context.Background(), []byte("not a zip"), cfg, deps)
	if ok || !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("ok=%v err=%v", ok, err)
	}
	if calls.Load() != 0 {
		t.Fatalf("server called %d times", calls.Load())
	}
}

func TestConvertRemoteErrorCodes(t *testing.T) {
	for _, w := range wireErrors {
		t.Run(w.code, func(t *testing.T) {
			cfg, deps := testRemote(t, func(rw http.ResponseWriter, _ *http.Request) {
				writeWireError(rw, w.code, "detail for "+w.code)
			})
			_, ok, err := convertRemote(context.Background(), testDocx, cfg, deps)
			if ok || !errors.Is(err, w.err) {
				t.Fatalf("ok=%v err=%v, want %v", ok, err, w.err)
			}
			if !strings.Contains(err.Error(), "detail for "+w.code) {
				t.Fatalf("server message dropped: %v", err)
			}
		})
	}
	// Non-JSON bodies (e.g. a Cloud Run front-end page) classify by status.
	for status, want := range map[int]error{
		http.StatusUnauthorized:          ErrRemoteAuth,
		http.StatusForbidden:             ErrRemoteAuth,
		http.StatusRequestEntityTooLarge: ErrInvalidInput,
		http.StatusTooManyRequests:       ErrBusy,
		http.StatusGatewayTimeout:        ErrTimeout,
		http.StatusBadGateway:            ErrRemoteUnavailable,
		http.StatusNotFound:              ErrRemoteUnavailable,
	} {
		cfg, deps := testRemote(t, func(rw http.ResponseWriter, _ *http.Request) {
			http.Error(rw, "<html>front end</html>", status)
		})
		if _, _, err := convertRemote(context.Background(), testDocx, cfg, deps); !errors.Is(err, want) {
			t.Fatalf("status %d: err=%v, want %v", status, err, want)
		}
	}
}

func TestConvertRemoteRetryPolicy(t *testing.T) {
	cases := []struct {
		name      string
		first     func(http.ResponseWriter)
		wantCalls int32
		wantOK    bool
	}{
		{"unhealthy retried", func(w http.ResponseWriter) { writeWireError(w, CodeUnhealthy, "") }, 2, true},
		{"503 page retried", func(w http.ResponseWriter) { http.Error(w, "x", http.StatusServiceUnavailable) }, 2, true},
		{"busy not retried", func(w http.ResponseWriter) { writeWireError(w, "busy", "") }, 1, false},
		{"timeout not retried", func(w http.ResponseWriter) { writeWireError(w, "timeout", "") }, 1, false},
		{"auth not retried", func(w http.ResponseWriter) { writeWireError(w, "unauthorized", "") }, 1, false},
		{"conversion failure not retried", func(w http.ResponseWriter) { writeWireError(w, "source_not_loaded", "") }, 1, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32
			cfg, deps := testRemote(t, func(w http.ResponseWriter, _ *http.Request) {
				if calls.Add(1) == 1 {
					tc.first(w)
					return
				}
				_, _ = w.Write(testPDF)
			})
			_, ok, err := convertRemote(context.Background(), testDocx, cfg, deps)
			if calls.Load() != tc.wantCalls || ok != tc.wantOK {
				t.Fatalf("calls=%d ok=%v err=%v, want calls=%d ok=%v", calls.Load(), ok, err, tc.wantCalls, tc.wantOK)
			}
		})
	}

	// A refused connection is a transport failure: retried once, then reported.
	srv := httptest.NewServer(http.NotFoundHandler())
	cfg, _ := parseRemoteConfig(srv.URL, testKey)
	srv.Close()
	_, _, err := convertRemote(context.Background(), testDocx, cfg, remoteDeps{client: &http.Client{Timeout: time.Second}, retryBackoff: time.Millisecond})
	if !errors.Is(err, ErrRemoteUnavailable) {
		t.Fatalf("closed server: err=%v, want ErrRemoteUnavailable", err)
	}
}

func TestConvertRemoteNonPDFBody(t *testing.T) {
	cfg, deps := testRemote(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("<html>login</html>"))
	})
	if _, ok, err := convertRemote(context.Background(), testDocx, cfg, deps); ok || !errors.Is(err, ErrNoOutput) {
		t.Fatalf("ok=%v err=%v, want ErrNoOutput", ok, err)
	}
}

func TestConvertRemoteNoRedirect(t *testing.T) {
	cfg, _ := testRemote(t, func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/elsewhere", http.StatusTemporaryRedirect)
	})
	deps := productionRemoteDeps()
	deps.retryBackoff = time.Millisecond
	if _, ok, err := convertRemote(context.Background(), testDocx, cfg, deps); ok || !errors.Is(err, ErrRemoteUnavailable) {
		t.Fatalf("ok=%v err=%v, want ErrRemoteUnavailable (redirect not followed)", ok, err)
	}
}

func TestConvertRemoteContextCancel(t *testing.T) {
	release := make(chan struct{})
	cfg, deps := testRemote(t, func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
		case <-release:
		}
	})
	defer close(release)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	started := time.Now()
	_, ok, err := convertRemote(ctx, testDocx, cfg, deps)
	if ok || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("ok=%v err=%v, want context deadline", ok, err)
	}
	if time.Since(started) > 5*time.Second {
		t.Fatalf("cancel took %s", time.Since(started))
	}
}

func TestErrorCodeRoundTrip(t *testing.T) {
	for _, w := range wireErrors {
		sentinel, retryable, known := errorForCode(w.code)
		if !known || sentinel != w.err || retryable != w.retryable {
			t.Fatalf("errorForCode(%q) = %v,%v,%v", w.code, sentinel, retryable, known)
		}
		if got := StatusForCode(w.code); got != w.status {
			t.Fatalf("StatusForCode(%q) = %d, want %d", w.code, got, w.status)
		}
	}
	for _, sentinel := range []error{ErrInvalidInput, ErrBusy, ErrTimeout, ErrCannotFork, ErrSourceNotLoaded, ErrConversionFailed, ErrNoOutput, ErrRemoteUnavailable, ErrRemoteAuth} {
		code, _ := ErrorCode(errors.Join(errors.New("wrapped"), sentinel))
		back, _, _ := errorForCode(code)
		if back != sentinel {
			t.Fatalf("%v -> %q -> %v", sentinel, code, back)
		}
	}
	if code, status := ErrorCode(errors.New("other")); code != CodeInternal || status != http.StatusInternalServerError {
		t.Fatalf("unknown error -> %q %d", code, status)
	}
}
