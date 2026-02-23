package fycha

import (
	"context"
	"errors"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

// ErrObjectNotFound is returned when a storage object does not exist.
var ErrObjectNotFound = errors.New("object not found")

// StorageReadResult holds the content and metadata of a downloaded file.
type StorageReadResult struct {
	Content     []byte
	ContentType string
}

// StorageReader reads objects from a storage backend.
// Implementations wrap provider-specific adapters (e.g., espyna StorageAdapter)
// to keep fycha provider-agnostic.
type StorageReader interface {
	ReadObject(ctx context.Context, containerName, objectKey string) (*StorageReadResult, error)
}

// StorageRouteRegistrar registers HTTP routes. Defined here to avoid
// circular imports with framework packages.
type StorageRouteRegistrar interface {
	HandleFunc(method, path string, handler http.HandlerFunc, middlewares ...string)
}

// StorageHandler serves files from a storage backend via HTTP.
// It is provider-agnostic — the actual storage implementation is injected
// via the StorageReader interface.
type StorageHandler struct {
	storage       StorageReader
	containerName string
	routePrefix   string
}

// NewStorageHandler creates a handler that serves files from storage.
//   - storage: the provider-agnostic reader (wraps espyna, GCS, S3, mock, etc.)
//   - containerName: the bucket/container to read from
//   - routePrefix: the URL prefix (e.g., "/storage/images")
func NewStorageHandler(storage StorageReader, containerName, routePrefix string) *StorageHandler {
	return &StorageHandler{
		storage:       storage,
		containerName: containerName,
		routePrefix:   strings.TrimRight(routePrefix, "/"),
	}
}

// RegisterRoutes registers the file serving route.
func (h *StorageHandler) RegisterRoutes(r StorageRouteRegistrar) {
	r.HandleFunc("GET", h.routePrefix+"/{path...}", h.serveFile)
}

// serveFile reads an object from storage and streams it to the HTTP response.
func (h *StorageHandler) serveFile(w http.ResponseWriter, r *http.Request) {
	objectKey := r.PathValue("path")
	if objectKey == "" {
		http.NotFound(w, r)
		return
	}

	// Reject path traversal attempts
	if strings.Contains(objectKey, "..") {
		http.NotFound(w, r)
		return
	}

	ctx := r.Context()
	result, err := h.storage.ReadObject(ctx, h.containerName, objectKey)
	if err != nil {
		if errors.Is(err, ErrObjectNotFound) {
			http.NotFound(w, r)
			return
		}
		// Client may have disconnected
		if ctx.Err() != nil {
			return
		}
		log.Printf("storage read error for %s: %v", objectKey, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Determine content type: prefer metadata, fall back to extension
	contentType := result.ContentType
	if contentType == "" || contentType == "application/octet-stream" {
		contentType = contentTypeFromExt(filepath.Ext(objectKey))
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Write(result.Content)
}

// contentTypeFromExt returns a MIME type for common file extensions.
func contentTypeFromExt(ext string) string {
	ext = strings.ToLower(ext)
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".webp":
		return "image/webp"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	case ".avif":
		return "image/avif"
	case ".pdf":
		return "application/pdf"
	default:
		ct := mime.TypeByExtension(ext)
		if ct != "" {
			return ct
		}
		return "application/octet-stream"
	}
}
