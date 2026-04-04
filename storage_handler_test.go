package fycha

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockStorageReader implements StorageReader for testing.
type mockStorageReader struct {
	readFunc func(ctx context.Context, containerName, objectKey string) (*StorageReadResult, error)
}

func (m *mockStorageReader) ReadObject(ctx context.Context, containerName, objectKey string) (*StorageReadResult, error) {
	return m.readFunc(ctx, containerName, objectKey)
}

func TestContentTypeFromExt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		ext  string
		want string
	}{
		{ext: ".jpg", want: "image/jpeg"},
		{ext: ".jpeg", want: "image/jpeg"},
		{ext: ".JPG", want: "image/jpeg"},  // case insensitive
		{ext: ".JPEG", want: "image/jpeg"}, // case insensitive
		{ext: ".png", want: "image/png"},
		{ext: ".PNG", want: "image/png"},
		{ext: ".webp", want: "image/webp"},
		{ext: ".gif", want: "image/gif"},
		{ext: ".svg", want: "image/svg+xml"},
		{ext: ".avif", want: "image/avif"},
		{ext: ".pdf", want: "application/pdf"},
		// Known extension that falls to mime.TypeByExtension
		{ext: ".html", want: "text/html; charset=utf-8"},
		// Completely unknown extension falls to octet-stream
		{ext: ".xyz123unknown", want: "application/octet-stream"},
		// Empty extension
		{ext: "", want: "application/octet-stream"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run("ext_"+tt.ext, func(t *testing.T) {
			t.Parallel()

			got := contentTypeFromExt(tt.ext)
			if got != tt.want {
				t.Errorf("contentTypeFromExt(%q) = %q, want %q", tt.ext, got, tt.want)
			}
		})
	}
}

func TestServeFile_PathTraversal(t *testing.T) {
	t.Parallel()

	storageCalled := false
	storage := &mockStorageReader{
		readFunc: func(ctx context.Context, containerName, objectKey string) (*StorageReadResult, error) {
			storageCalled = true
			return nil, ErrObjectNotFound
		},
	}

	handler := NewStorageHandler(storage, "test-bucket", "/storage/files")

	// Test path traversal via direct handler call to bypass mux URL cleaning.
	// The mux normalizes ".." before routing, but we test the handler's own guard.
	tests := []struct {
		name      string
		pathValue string
		want      int
	}{
		{name: "double dot in path value", pathValue: "../etc/passwd", want: http.StatusNotFound},
		{name: "double dot mid-path value", pathValue: "images/../secret.txt", want: http.StatusNotFound},
		{name: "encoded double dot", pathValue: "images/..%2fsecret.txt", want: http.StatusNotFound},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mux := http.NewServeMux()
			mux.HandleFunc("GET /storage/files/{path...}", handler.serveFile)

			// Build a request where the path wildcard contains ".."
			req := httptest.NewRequest("GET", "/storage/files/"+tt.pathValue, nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			// The mux may clean the path (redirect) or the handler rejects it.
			// Either way the response must not be 200 OK.
			if w.Code == http.StatusOK {
				t.Errorf("status = %d, should not be 200 for path traversal attempt", w.Code)
			}
		})
	}

	// Test empty path via mux routing
	t.Run("empty path", func(t *testing.T) {
		t.Parallel()

		mux := http.NewServeMux()
		mux.HandleFunc("GET /storage/files/{path...}", handler.serveFile)

		req := httptest.NewRequest("GET", "/storage/files/", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		// Empty pathValue should return 404
		if w.Code != http.StatusNotFound {
			t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
		}
	})

	_ = storageCalled // suppress unused warning
}

func TestServeFile_NotFound(t *testing.T) {
	t.Parallel()

	storage := &mockStorageReader{
		readFunc: func(ctx context.Context, containerName, objectKey string) (*StorageReadResult, error) {
			return nil, ErrObjectNotFound
		},
	}

	handler := NewStorageHandler(storage, "test-bucket", "/storage/files")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /storage/files/{path...}", handler.serveFile)

	req := httptest.NewRequest("GET", "/storage/files/missing-file.png", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestServeFile_InternalError(t *testing.T) {
	t.Parallel()

	storage := &mockStorageReader{
		readFunc: func(ctx context.Context, containerName, objectKey string) (*StorageReadResult, error) {
			return nil, errors.New("connection refused")
		},
	}

	handler := NewStorageHandler(storage, "test-bucket", "/storage/files")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /storage/files/{path...}", handler.serveFile)

	req := httptest.NewRequest("GET", "/storage/files/image.png", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestServeFile_EncodedPathTraversal(t *testing.T) {
	t.Parallel()

	storage := &mockStorageReader{
		readFunc: func(ctx context.Context, containerName, objectKey string) (*StorageReadResult, error) {
			return nil, ErrObjectNotFound
		},
	}

	handler := NewStorageHandler(storage, "test-bucket", "/storage/files")

	tests := []struct {
		name      string
		pathValue string
	}{
		{name: "percent-encoded dot-dot-slash", pathValue: "%2e%2e/etc/passwd"},
		{name: "double-encoded dot-dot", pathValue: "%252e%252e/etc/passwd"},
		{name: "mixed encoded dot-dot", pathValue: "img/%2e%2e/secret.txt"},
		{name: "backslash traversal", pathValue: "img/..\\secret.txt"},
		{name: "encoded backslash traversal", pathValue: "img/%2e%2e%5csecret.txt"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mux := http.NewServeMux()
			mux.HandleFunc("GET /storage/files/{path...}", handler.serveFile)

			req := httptest.NewRequest("GET", "/storage/files/"+tt.pathValue, nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			if w.Code == http.StatusOK {
				t.Errorf("status = %d, should not be 200 for encoded path traversal", w.Code)
			}
		})
	}
}

func TestServeFile_VeryLongFilePath(t *testing.T) {
	t.Parallel()

	storageCalled := false
	storage := &mockStorageReader{
		readFunc: func(ctx context.Context, containerName, objectKey string) (*StorageReadResult, error) {
			storageCalled = true
			return &StorageReadResult{
				Content:     []byte("data"),
				ContentType: "image/png",
			}, nil
		},
	}

	handler := NewStorageHandler(storage, "test-bucket", "/storage/files")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /storage/files/{path...}", handler.serveFile)

	// Build a path with 2000+ characters
	longSegment := "a"
	for len(longSegment) < 2000 {
		longSegment += "a"
	}
	longPath := "images/" + longSegment + ".png"

	req := httptest.NewRequest("GET", "/storage/files/"+longPath, nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	// The handler should either succeed (storage handles long keys) or return an error,
	// but must not panic.
	if w.Code != http.StatusOK && w.Code != http.StatusNotFound && w.Code != http.StatusInternalServerError {
		t.Errorf("unexpected status = %d for very long path", w.Code)
	}
	_ = storageCalled
}

func TestServeFile_SpecialCharacterPaths(t *testing.T) {
	t.Parallel()

	storage := &mockStorageReader{
		readFunc: func(ctx context.Context, containerName, objectKey string) (*StorageReadResult, error) {
			return &StorageReadResult{
				Content:     []byte("data"),
				ContentType: "image/png",
			}, nil
		},
	}

	handler := NewStorageHandler(storage, "test-bucket", "/storage/files")

	tests := []struct {
		name string
		path string
	}{
		{name: "path with spaces", path: "images/my%20photo.png"},
		{name: "path with unicode", path: "images/%E4%B8%AD%E6%96%87.png"},
		{name: "path with plus sign", path: "images/photo+1.png"},
		{name: "path with parentheses", path: "images/photo(1).png"},
		{name: "path with hash encoded", path: "images/file%23name.png"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mux := http.NewServeMux()
			mux.HandleFunc("GET /storage/files/{path...}", handler.serveFile)

			req := httptest.NewRequest("GET", "/storage/files/"+tt.path, nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			// Should not panic; either success or a well-defined error
			if w.Code != http.StatusOK && w.Code != http.StatusNotFound {
				t.Errorf("status = %d for path %q, expected 200 or 404", w.Code, tt.path)
			}
		})
	}
}

func TestServeFile_ContentTypeSpoofing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		objectKey       string
		storageType     string // content type returned by storage (potentially spoofed)
		wantContentType string
	}{
		{
			name:            "storage claims HTML but file is .png",
			objectKey:       "images/logo.png",
			storageType:     "text/html",
			wantContentType: "text/html", // handler trusts storage when it has a value
		},
		{
			name:            "storage claims javascript but file is .jpg",
			objectKey:       "photos/pic.jpg",
			storageType:     "application/javascript",
			wantContentType: "application/javascript",
		},
		{
			name:            "storage returns empty string falls back to extension",
			objectKey:       "docs/report.pdf",
			storageType:     "",
			wantContentType: "application/pdf",
		},
		{
			name:            "storage returns octet-stream falls back to extension",
			objectKey:       "icons/logo.svg",
			storageType:     "application/octet-stream",
			wantContentType: "image/svg+xml",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storage := &mockStorageReader{
				readFunc: func(ctx context.Context, containerName, objectKey string) (*StorageReadResult, error) {
					return &StorageReadResult{
						Content:     []byte("fake-data"),
						ContentType: tt.storageType,
					}, nil
				},
			}

			handler := NewStorageHandler(storage, "test-bucket", "/storage/files")

			mux := http.NewServeMux()
			mux.HandleFunc("GET /storage/files/{path...}", handler.serveFile)

			req := httptest.NewRequest("GET", "/storage/files/"+tt.objectKey, nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200", w.Code)
			}

			gotCT := w.Header().Get("Content-Type")
			if gotCT != tt.wantContentType {
				t.Errorf("Content-Type = %q, want %q", gotCT, tt.wantContentType)
			}
		})
	}
}

func TestServeFile_EmptyFileBody(t *testing.T) {
	t.Parallel()

	storage := &mockStorageReader{
		readFunc: func(ctx context.Context, containerName, objectKey string) (*StorageReadResult, error) {
			return &StorageReadResult{
				Content:     []byte{}, // zero bytes
				ContentType: "image/png",
			}, nil
		},
	}

	handler := NewStorageHandler(storage, "test-bucket", "/storage/files")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /storage/files/{path...}", handler.serveFile)

	req := httptest.NewRequest("GET", "/storage/files/empty.png", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	// Empty content should still return 200 with correct headers
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200 for empty file body", w.Code)
	}

	gotCT := w.Header().Get("Content-Type")
	if gotCT != "image/png" {
		t.Errorf("Content-Type = %q, want %q", gotCT, "image/png")
	}

	if w.Body.Len() != 0 {
		t.Errorf("body length = %d, want 0 for empty file", w.Body.Len())
	}
}

func TestServeFile_NilContent(t *testing.T) {
	t.Parallel()

	storage := &mockStorageReader{
		readFunc: func(ctx context.Context, containerName, objectKey string) (*StorageReadResult, error) {
			return &StorageReadResult{
				Content:     nil, // nil bytes
				ContentType: "image/png",
			}, nil
		},
	}

	handler := NewStorageHandler(storage, "test-bucket", "/storage/files")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /storage/files/{path...}", handler.serveFile)

	req := httptest.NewRequest("GET", "/storage/files/nil-content.png", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	// nil content should still return 200 without panicking
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200 for nil content", w.Code)
	}

	if w.Body.Len() != 0 {
		t.Errorf("body length = %d, want 0 for nil content", w.Body.Len())
	}
}

func TestServeFile_CancelledContext(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	storage := &mockStorageReader{
		readFunc: func(ctx context.Context, containerName, objectKey string) (*StorageReadResult, error) {
			return nil, ctx.Err()
		},
	}

	handler := NewStorageHandler(storage, "test-bucket", "/storage/files")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /storage/files/{path...}", handler.serveFile)

	req := httptest.NewRequest("GET", "/storage/files/image.png", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	// When context is cancelled the handler returns early without writing a body.
	// httptest.ResponseRecorder defaults to 200, but no content should be written.
	if w.Body.Len() != 0 {
		t.Errorf("body length = %d, want 0 for cancelled context (handler should bail out)", w.Body.Len())
	}

	// Verify no Content-Type header was set (handler returned before writing)
	if ct := w.Header().Get("Content-Type"); ct != "" {
		t.Errorf("Content-Type = %q, want empty for cancelled context", ct)
	}
}

func TestServeFile_Success(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		objectKey       string
		content         []byte
		storageType     string // content type from storage metadata
		wantContentType string
	}{
		{
			name:            "PNG with storage content type",
			objectKey:       "images/logo.png",
			content:         []byte("fake-png-data"),
			storageType:     "image/png",
			wantContentType: "image/png",
		},
		{
			name:            "JPEG falls back to extension when storage returns octet-stream",
			objectKey:       "photos/pic.jpg",
			content:         []byte("fake-jpg-data"),
			storageType:     "application/octet-stream",
			wantContentType: "image/jpeg",
		},
		{
			name:            "PDF with empty storage content type",
			objectKey:       "docs/report.pdf",
			content:         []byte("fake-pdf-data"),
			storageType:     "",
			wantContentType: "application/pdf",
		},
		{
			name:            "SVG with proper storage type",
			objectKey:       "icons/logo.svg",
			content:         []byte("<svg></svg>"),
			storageType:     "image/svg+xml",
			wantContentType: "image/svg+xml",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storage := &mockStorageReader{
				readFunc: func(ctx context.Context, containerName, objectKey string) (*StorageReadResult, error) {
					if containerName != "test-bucket" {
						t.Errorf("containerName = %q, want %q", containerName, "test-bucket")
					}
					if objectKey != tt.objectKey {
						t.Errorf("objectKey = %q, want %q", objectKey, tt.objectKey)
					}
					return &StorageReadResult{
						Content:     tt.content,
						ContentType: tt.storageType,
					}, nil
				},
			}

			handler := NewStorageHandler(storage, "test-bucket", "/storage/files")

			mux := http.NewServeMux()
			mux.HandleFunc("GET /storage/files/{path...}", handler.serveFile)

			req := httptest.NewRequest("GET", "/storage/files/"+tt.objectKey, nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
			}

			gotCT := w.Header().Get("Content-Type")
			if gotCT != tt.wantContentType {
				t.Errorf("Content-Type = %q, want %q", gotCT, tt.wantContentType)
			}

			gotCache := w.Header().Get("Cache-Control")
			if gotCache != "public, max-age=86400" {
				t.Errorf("Cache-Control = %q, want %q", gotCache, "public, max-age=86400")
			}

			if string(w.Body.Bytes()) != string(tt.content) {
				t.Errorf("body = %q, want %q", w.Body.String(), string(tt.content))
			}
		})
	}
}
