package pdfconv

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// Remote conversion: when RemoteURLEnv is set, ConvertDocxToPDFContext POSTs
// the DOCX to a dedicated converter service (apps/pdf-converter) instead of
// running LibreOffice on this instance. The wire contract is shared with that
// server through the constants and the error-code table below, so both sides
// agree on one mapping and callers see the same sentinels as the local path.
const (
	// RemoteURLEnv is the converter's base URL (https, or http on loopback only).
	RemoteURLEnv = "FYCHA_PDF_CONVERTER_URL"
	// RemoteAPIKeyEnv is the shared bearer key: sent by the client, expected by the server.
	RemoteAPIKeyEnv = "FYCHA_PDF_CONVERTER_API_KEY"
	// RemoteConvertPath is appended to the base URL unless it already ends with it.
	RemoteConvertPath = "/v1/convert"
	// RequestIDHeader correlates the client and server log lines of one conversion.
	RequestIDHeader = "X-Request-Id"
	// DocxContentType is the request body type.
	DocxContentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"

	// remoteRequestTimeout bounds one HTTP attempt end to end. It sits just above
	// the converter's Cloud Run request timeout (300s) so the server's own 504
	// normally arrives first.
	remoteRequestTimeout = 310 * time.Second
	// maxRemoteResponseBytes caps how much of a response body is read.
	maxRemoteResponseBytes = 100 << 20
	// remoteMessageLimit caps the server message carried into a client error.
	remoteMessageLimit = 1 << 10
)

// Remote-only failures. Conversion failures reported by the server map back to
// the local sentinels (ErrBusy, ErrTimeout, ...) instead.
var (
	// ErrRemoteMisconfigured: RemoteURLEnv is set but unusable (no key, plain
	// http to a non-loopback host, embedded credentials, ...). Nothing is sent.
	ErrRemoteMisconfigured = errors.New("pdf converter remote is misconfigured")
	// ErrRemoteAuth: the converter rejected the API key (HTTP 401/403).
	ErrRemoteAuth = errors.New("pdf converter rejected the api key")
	// ErrRemoteUnavailable: the converter could not be reached, is recycling,
	// or answered with an unexpected status.
	ErrRemoteUnavailable = errors.New("pdf converter unavailable")
)

// Wire codes the server sends that have no sentinel of their own.
const (
	// CodeUnhealthy: the converter instance is recycling (retry lands elsewhere).
	CodeUnhealthy = "unhealthy"
	// CodeTooLarge: the request body exceeded the server limit.
	CodeTooLarge = "too_large"
	// CodeInternal: an error outside the table.
	CodeInternal = "internal"
)

// wireErrors is the one mapping between sentinel, wire code, HTTP status, and
// whether the client retries. ErrorCode takes the first entry whose sentinel
// matches, so each sentinel's primary code comes first.
var wireErrors = []struct {
	err       error
	code      string
	status    int
	retryable bool
}{
	{ErrInvalidInput, "invalid_input", http.StatusBadRequest, false},
	{ErrInvalidInput, CodeTooLarge, http.StatusRequestEntityTooLarge, false},
	{ErrRemoteAuth, "unauthorized", http.StatusUnauthorized, false},
	{ErrBusy, "busy", http.StatusServiceUnavailable, false},
	{ErrTimeout, "timeout", http.StatusGatewayTimeout, false},
	{ErrCannotFork, "cannot_fork", http.StatusInternalServerError, false},
	{ErrSourceNotLoaded, "source_not_loaded", http.StatusInternalServerError, false},
	{ErrConversionFailed, "conversion_failed", http.StatusInternalServerError, false},
	{ErrNoOutput, "no_output", http.StatusInternalServerError, false},
	{ErrRemoteUnavailable, "unavailable", http.StatusServiceUnavailable, true},
	{ErrRemoteUnavailable, CodeUnhealthy, http.StatusServiceUnavailable, true},
}

// ErrorCode returns the wire code and HTTP status the converter server reports
// for a conversion error.
func ErrorCode(err error) (code string, status int) {
	for _, w := range wireErrors {
		if errors.Is(err, w.err) {
			return w.code, w.status
		}
	}
	return CodeInternal, http.StatusInternalServerError
}

// StatusForCode returns the HTTP status for a wire code (500 when unknown).
func StatusForCode(code string) int {
	for _, w := range wireErrors {
		if w.code == code {
			return w.status
		}
	}
	return http.StatusInternalServerError
}

// ErrorBody is the JSON body of every non-200 converter response.
type ErrorBody struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

func errorForCode(code string) (sentinel error, retryable, known bool) {
	for _, w := range wireErrors {
		if w.code == code {
			return w.err, w.retryable, true
		}
	}
	return nil, false, false
}

type remoteConfig struct {
	endpoint string
	host     string
	apiKey   string
	err      error // set when configured but unusable (fail closed)
}

type remoteDeps struct {
	client       *http.Client
	retryBackoff time.Duration
}

// remoteHTTPClient never follows redirects: a converter POST has no reason to
// redirect, and following one could hand the bearer key to another path.
var remoteHTTPClient = &http.Client{
	Timeout: remoteRequestTimeout,
	CheckRedirect: func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

func productionRemoteDeps() remoteDeps {
	return remoteDeps{client: remoteHTTPClient, retryBackoff: defaultRetryBackoff}
}

// remoteConfigFromEnv reports whether remote conversion is configured. When it
// is but the values are unusable, cfg.err is set and nothing may be sent.
func remoteConfigFromEnv() (remoteConfig, bool) {
	return parseRemoteConfig(os.Getenv(RemoteURLEnv), os.Getenv(RemoteAPIKeyEnv))
}

func parseRemoteConfig(rawURL, apiKey string) (remoteConfig, bool) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return remoteConfig{}, false
	}
	fail := func(reason string) (remoteConfig, bool) {
		return remoteConfig{err: fmt.Errorf("pdfconv: %w: %s %s", ErrRemoteMisconfigured, RemoteURLEnv, reason)}, true
	}
	u, err := url.Parse(rawURL)
	if err != nil || u.Host == "" {
		// The raw value is not echoed: it may carry credentials.
		return fail("is not an absolute URL")
	}
	if u.User != nil {
		return fail("must not embed credentials (use " + RemoteAPIKeyEnv + ")")
	}
	if u.RawQuery != "" || u.Fragment != "" {
		return fail("must not carry a query or fragment")
	}
	switch u.Scheme {
	case "https":
	case "http":
		if !isLoopbackHost(u.Hostname()) {
			return fail("must use https (plain http is allowed only for a loopback host)")
		}
	default:
		return fail("must use https")
	}
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return fail("is set but " + RemoteAPIKeyEnv + " is empty")
	}
	path := strings.TrimRight(u.Path, "/")
	if !strings.HasSuffix(path, RemoteConvertPath) {
		path += RemoteConvertPath
	}
	u.Path, u.RawPath = path, ""
	return remoteConfig{endpoint: u.String(), host: u.Host, apiKey: apiKey}, true
}

func isLoopbackHost(host string) bool {
	if strings.EqualFold(host, "localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

// convertRemote sends the DOCX to the remote converter. It returns ok=true on
// success and never ok=false: "LibreOffice absent" is a local-only condition.
func convertRemote(ctx context.Context, docxBytes []byte, cfg remoteConfig, deps remoteDeps) ([]byte, bool, error) {
	if cfg.err != nil {
		log.Printf("pdfconv: remote result=misconfigured: %v", cfg.err)
		return nil, false, cfg.err
	}
	if !bytes.HasPrefix(docxBytes, zipMagic) {
		return nil, false, fmt.Errorf("pdfconv: %w (%d bytes, head %s)", ErrInvalidInput, len(docxBytes), headHex(docxBytes))
	}

	reqID := newRequestID()
	started := time.Now()
	var (
		pdfBytes []byte
		status   int
		err      error
		attempt  int
	)
	for attempt = 1; ; attempt++ {
		var retryable bool
		pdfBytes, status, retryable, err = remoteAttempt(ctx, docxBytes, cfg, deps, reqID)
		if err == nil || !retryable || attempt >= maxAttempts || ctx.Err() != nil {
			break
		}
		log.Printf("pdfconv: remote attempt %d req=%s failed, retrying: %v", attempt, reqID, err)
		if !sleepCtx(ctx, deps.retryBackoff) {
			err = fmt.Errorf("pdfconv: remote conversion canceled before retry: %w", ctx.Err())
			break
		}
	}

	result := "ok"
	if err != nil {
		result, _ = ErrorCode(err)
		if ctx.Err() != nil {
			result = "canceled"
		}
	}
	log.Printf("pdfconv: remote result=%s attempts=%d status=%d in=%dB out=%dB took=%s host=%s req=%s",
		result, attempt, status, len(docxBytes), len(pdfBytes), time.Since(started).Round(time.Millisecond), cfg.host, reqID)
	if err != nil {
		return nil, false, err
	}
	return pdfBytes, true, nil
}

// remoteAttempt performs one HTTP round trip. retryable reports whether a
// second attempt could plausibly succeed (another instance, a transient blip).
func remoteAttempt(ctx context.Context, docxBytes []byte, cfg remoteConfig, deps remoteDeps, reqID string) (pdfBytes []byte, status int, retryable bool, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.endpoint, bytes.NewReader(docxBytes))
	if err != nil {
		return nil, 0, false, fmt.Errorf("pdfconv: %w: building request: %v", ErrRemoteMisconfigured, err)
	}
	req.Header.Set("Authorization", "Bearer "+cfg.apiKey)
	req.Header.Set("Content-Type", DocxContentType)
	req.Header.Set("Accept", "application/pdf")
	req.Header.Set(RequestIDHeader, reqID)

	resp, err := deps.client.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, 0, false, fmt.Errorf("pdfconv: remote conversion canceled: %w", ctx.Err())
		}
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return nil, 0, false, fmt.Errorf("pdfconv: remote %s: %w after %s: %v", cfg.host, ErrTimeout, deps.client.Timeout, err)
		}
		return nil, 0, true, fmt.Errorf("pdfconv: remote %s: %w: %v", cfg.host, ErrRemoteUnavailable, err)
	}
	defer resp.Body.Close()
	status = resp.StatusCode

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxRemoteResponseBytes+1))
	if err != nil {
		if ctx.Err() != nil {
			return nil, status, false, fmt.Errorf("pdfconv: remote conversion canceled: %w", ctx.Err())
		}
		return nil, status, true, fmt.Errorf("pdfconv: remote %s: %w: reading response: %v", cfg.host, ErrRemoteUnavailable, err)
	}

	if status == http.StatusOK {
		if len(body) > maxRemoteResponseBytes {
			return nil, status, false, fmt.Errorf("pdfconv: remote %s: %w: response exceeds %d bytes", cfg.host, ErrNoOutput, maxRemoteResponseBytes)
		}
		if !bytes.HasPrefix(body, pdfMagic) {
			return nil, status, false, fmt.Errorf("pdfconv: remote %s: %w: response head %s", cfg.host, ErrNoOutput, headHex(body))
		}
		return body, status, false, nil
	}

	sentinel, retryable, message := classifyRemoteFailure(status, body)
	detail := ""
	if message != "" {
		detail = "; message: " + message
	}
	return nil, status, retryable, fmt.Errorf("pdfconv: remote %s: %w (HTTP %d)%s", cfg.host, sentinel, status, detail)
}

// classifyRemoteFailure maps a non-200 response to a sentinel. The server's
// JSON code wins; otherwise (Cloud Run front-end pages, proxies) the status does.
func classifyRemoteFailure(status int, body []byte) (sentinel error, retryable bool, message string) {
	var eb ErrorBody
	if json.Unmarshal(body, &eb) == nil && eb.Error != "" {
		if s, r, ok := errorForCode(eb.Error); ok {
			return s, r, truncateMessage(eb.Message)
		}
		message = truncateMessage(eb.Error + ": " + eb.Message)
	}
	switch status {
	case http.StatusUnauthorized, http.StatusForbidden:
		return ErrRemoteAuth, false, message
	case http.StatusRequestEntityTooLarge:
		return ErrInvalidInput, false, message
	case http.StatusTooManyRequests:
		return ErrBusy, true, message
	case http.StatusGatewayTimeout:
		return ErrTimeout, false, message
	case http.StatusBadGateway, http.StatusServiceUnavailable:
		return ErrRemoteUnavailable, true, message
	}
	return ErrRemoteUnavailable, false, message
}

func truncateMessage(s string) string {
	s = strings.TrimSpace(s)
	if len(s) > remoteMessageLimit {
		return s[:remoteMessageLimit] + "…(truncated)"
	}
	return s
}

func newRequestID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "unknown"
	}
	return hex.EncodeToString(b)
}
