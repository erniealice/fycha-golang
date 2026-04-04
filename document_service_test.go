package fycha

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
)

// mockStorageReadWriter implements StorageReadWriter for testing.
type mockStorageReadWriter struct {
	readFunc  func(ctx context.Context, containerName, objectKey string) ([]byte, error)
	writeFunc func(ctx context.Context, containerName, objectKey string, data []byte) error
}

func (m *mockStorageReadWriter) ReadObject(ctx context.Context, containerName, objectKey string) ([]byte, error) {
	return m.readFunc(ctx, containerName, objectKey)
}

func (m *mockStorageReadWriter) WriteObject(ctx context.Context, containerName, objectKey string, data []byte) error {
	return m.writeFunc(ctx, containerName, objectKey, data)
}

func TestDocumentService_ProcessFromStorage_NilStorage(t *testing.T) {
	t.Parallel()

	svc := NewDocumentService(nil)

	err := svc.ProcessFromStorage(
		context.Background(),
		"templates", "invoice.docx",
		"output", "invoice-123.docx",
		map[string]any{"name": "Acme"},
	)
	if err == nil {
		t.Fatal("expected error for nil storage")
	}
	if err.Error() != "storage not configured" {
		t.Errorf("error = %q, want %q", err.Error(), "storage not configured")
	}
}

func TestDocumentService_ProcessFromStorageToBytes_NilStorage(t *testing.T) {
	t.Parallel()

	svc := NewDocumentService(nil)

	_, err := svc.ProcessFromStorageToBytes(
		context.Background(),
		"templates", "invoice.docx",
		map[string]any{"name": "Acme"},
	)
	if err == nil {
		t.Fatal("expected error for nil storage")
	}
	if err.Error() != "storage not configured" {
		t.Errorf("error = %q, want %q", err.Error(), "storage not configured")
	}
}

func TestDocumentService_ProcessFromStorage_ReadError(t *testing.T) {
	t.Parallel()

	storage := &mockStorageReadWriter{
		readFunc: func(ctx context.Context, containerName, objectKey string) ([]byte, error) {
			return nil, errors.New("blob not found")
		},
		writeFunc: func(ctx context.Context, containerName, objectKey string, data []byte) error {
			t.Fatal("write should not be called when read fails")
			return nil
		},
	}

	svc := NewDocumentService(storage)

	err := svc.ProcessFromStorage(
		context.Background(),
		"templates", "missing.docx",
		"output", "result.docx",
		map[string]any{},
	)
	if err == nil {
		t.Fatal("expected error when read fails")
	}
	if !strings.Contains(err.Error(), "reading template") {
		t.Errorf("error should mention reading template, got: %q", err.Error())
	}
}

func TestDocumentService_ProcessFromStorage_WriteError(t *testing.T) {
	t.Parallel()

	templateData := createMinimalDocx(t)

	storage := &mockStorageReadWriter{
		readFunc: func(ctx context.Context, containerName, objectKey string) ([]byte, error) {
			return templateData, nil
		},
		writeFunc: func(ctx context.Context, containerName, objectKey string, data []byte) error {
			return errors.New("write permission denied")
		},
	}

	svc := NewDocumentService(storage)

	err := svc.ProcessFromStorage(
		context.Background(),
		"templates", "invoice.docx",
		"output", "result.docx",
		map[string]any{"name": "Test"},
	)
	if err == nil {
		t.Fatal("expected error when write fails")
	}
	if !strings.Contains(err.Error(), "writing output") {
		t.Errorf("error should mention writing output, got: %q", err.Error())
	}
}

func TestDocumentService_ProcessFromStorage_Success(t *testing.T) {
	t.Parallel()

	templateData := createMinimalDocx(t)
	var writtenData []byte
	var writtenContainer, writtenKey string

	storage := &mockStorageReadWriter{
		readFunc: func(ctx context.Context, containerName, objectKey string) ([]byte, error) {
			if containerName != "templates" {
				t.Errorf("read containerName = %q, want %q", containerName, "templates")
			}
			if objectKey != "invoice.docx" {
				t.Errorf("read objectKey = %q, want %q", objectKey, "invoice.docx")
			}
			return templateData, nil
		},
		writeFunc: func(ctx context.Context, containerName, objectKey string, data []byte) error {
			writtenContainer = containerName
			writtenKey = objectKey
			writtenData = data
			return nil
		},
	}

	svc := NewDocumentService(storage)

	err := svc.ProcessFromStorage(
		context.Background(),
		"templates", "invoice.docx",
		"output", "result.docx",
		map[string]any{"name": "Test Corp"},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if writtenContainer != "output" {
		t.Errorf("write containerName = %q, want %q", writtenContainer, "output")
	}
	if writtenKey != "result.docx" {
		t.Errorf("write objectKey = %q, want %q", writtenKey, "result.docx")
	}
	if len(writtenData) == 0 {
		t.Error("written data should not be empty")
	}
}

func TestDocumentService_ProcessFromStorageToBytes_Success(t *testing.T) {
	t.Parallel()

	templateData := createMinimalDocx(t)

	storage := &mockStorageReadWriter{
		readFunc: func(ctx context.Context, containerName, objectKey string) ([]byte, error) {
			return templateData, nil
		},
		writeFunc: func(ctx context.Context, containerName, objectKey string, data []byte) error {
			t.Fatal("write should not be called for ToBytes")
			return nil
		},
	}

	svc := NewDocumentService(storage)

	result, err := svc.ProcessFromStorageToBytes(
		context.Background(),
		"templates", "invoice.docx",
		map[string]any{"name": "Test"},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) == 0 {
		t.Error("result should not be empty")
	}
}

func TestDocumentService_ProcessBytes_NilStorage(t *testing.T) {
	t.Parallel()

	// ProcessBytes should work even with nil storage
	svc := NewDocumentService(nil)

	templateData := createMinimalDocx(t)

	result, err := svc.ProcessBytes(templateData, map[string]any{"name": "Test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) == 0 {
		t.Error("result should not be empty")
	}
}

func TestDocumentService_ProcessBytes_EmptySlice(t *testing.T) {
	t.Parallel()

	svc := NewDocumentService(nil)

	_, err := svc.ProcessBytes([]byte{}, map[string]any{"name": "Test"})
	if err == nil {
		t.Fatal("expected error for empty byte slice")
	}
}

func TestDocumentService_ProcessBytes_NilSlice(t *testing.T) {
	t.Parallel()

	svc := NewDocumentService(nil)

	_, err := svc.ProcessBytes(nil, map[string]any{"name": "Test"})
	if err == nil {
		t.Fatal("expected error for nil byte slice")
	}
}

func TestDocumentService_ProcessBytes_InvalidDocx(t *testing.T) {
	t.Parallel()

	svc := NewDocumentService(nil)

	// Not a ZIP file at all — just random text bytes
	notAZip := []byte("This is definitely not a valid DOCX or ZIP file.")

	_, err := svc.ProcessBytes(notAZip, map[string]any{"name": "Test"})
	if err == nil {
		t.Fatal("expected error for invalid DOCX (not a ZIP)")
	}
}

func TestDocumentService_ProcessBytes_CorruptZip(t *testing.T) {
	t.Parallel()

	svc := NewDocumentService(nil)

	// Start with ZIP magic bytes but then corrupt the rest
	corruptZip := []byte("PK\x03\x04corrupted-data-here-not-a-real-zip")

	_, err := svc.ProcessBytes(corruptZip, map[string]any{"name": "Test"})
	if err == nil {
		t.Fatal("expected error for corrupt ZIP data")
	}
}

func TestDocumentService_ProcessFromStorage_EmptyContainerName(t *testing.T) {
	t.Parallel()

	readCalled := false
	storage := &mockStorageReadWriter{
		readFunc: func(ctx context.Context, containerName, objectKey string) ([]byte, error) {
			readCalled = true
			if containerName == "" {
				return nil, errors.New("container name is empty")
			}
			return createMinimalDocx(t), nil
		},
		writeFunc: func(ctx context.Context, containerName, objectKey string, data []byte) error {
			return nil
		},
	}

	svc := NewDocumentService(storage)

	err := svc.ProcessFromStorage(
		context.Background(),
		"", "template.docx",
		"output", "result.docx",
		map[string]any{"name": "Test"},
	)
	if err == nil {
		t.Fatal("expected error for empty container name")
	}
	if !readCalled {
		t.Error("read should have been called (validation is in storage layer)")
	}
}

func TestDocumentService_ProcessFromStorage_EmptyKey(t *testing.T) {
	t.Parallel()

	storage := &mockStorageReadWriter{
		readFunc: func(ctx context.Context, containerName, objectKey string) ([]byte, error) {
			if objectKey == "" {
				return nil, errors.New("object key is empty")
			}
			return createMinimalDocx(t), nil
		},
		writeFunc: func(ctx context.Context, containerName, objectKey string, data []byte) error {
			return nil
		},
	}

	svc := NewDocumentService(storage)

	err := svc.ProcessFromStorage(
		context.Background(),
		"templates", "",
		"output", "result.docx",
		map[string]any{"name": "Test"},
	)
	if err == nil {
		t.Fatal("expected error for empty template key")
	}
	if !strings.Contains(err.Error(), "reading template") {
		t.Errorf("error should mention reading template, got: %q", err.Error())
	}
}

func TestDocumentService_ProcessFromStorage_EmptyOutputKey(t *testing.T) {
	t.Parallel()

	templateData := createMinimalDocx(t)

	storage := &mockStorageReadWriter{
		readFunc: func(ctx context.Context, containerName, objectKey string) ([]byte, error) {
			return templateData, nil
		},
		writeFunc: func(ctx context.Context, containerName, objectKey string, data []byte) error {
			if objectKey == "" {
				return errors.New("output key is empty")
			}
			return nil
		},
	}

	svc := NewDocumentService(storage)

	err := svc.ProcessFromStorage(
		context.Background(),
		"templates", "template.docx",
		"output", "",
		map[string]any{"name": "Test"},
	)
	if err == nil {
		t.Fatal("expected error for empty output key")
	}
	if !strings.Contains(err.Error(), "writing output") {
		t.Errorf("error should mention writing output, got: %q", err.Error())
	}
}

func TestDocumentService_ProcessBytes_MalformedPlaceholders(t *testing.T) {
	t.Parallel()

	// Create a DOCX with malformed/unclosed placeholders
	contentTypesXML := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/><Default Extension="xml" ContentType="application/xml"/><Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/></Types>`

	relsXML := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/></Relationships>`

	documentRelsXML := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"></Relationships>`

	tests := []struct {
		name       string
		docContent string
	}{
		{
			name: "unclosed placeholder",
			docContent: `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
<w:body><w:p><w:r><w:t>Hello {{name</w:t></w:r></w:p></w:body></w:document>`,
		},
		{
			name: "nested braces",
			docContent: `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
<w:body><w:p><w:r><w:t>Hello {{{{name}}}}</w:t></w:r></w:p></w:body></w:document>`,
		},
		{
			name: "placeholder with special chars",
			docContent: `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
<w:body><w:p><w:r><w:t>Hello {{na&lt;me}}</w:t></w:r></w:p></w:body></w:document>`,
		},
	}

	svc := NewDocumentService(nil)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			buf := new(bytes.Buffer)
			w := zip.NewWriter(buf)

			files := map[string]string{
				"[Content_Types].xml":          contentTypesXML,
				"_rels/.rels":                  relsXML,
				"word/document.xml":            tt.docContent,
				"word/_rels/document.xml.rels": documentRelsXML,
			}

			for name, content := range files {
				writer, err := w.Create(name)
				if err != nil {
					t.Fatalf("failed to create zip entry %s: %v", name, err)
				}
				if _, err := writer.Write([]byte(content)); err != nil {
					t.Fatalf("failed to write zip entry %s: %v", name, err)
				}
			}

			if err := w.Close(); err != nil {
				t.Fatalf("failed to close zip: %v", err)
			}

			// ProcessBytes should not panic on malformed placeholders.
			// It may succeed (leaving the placeholder as-is) or return an error.
			_, _ = svc.ProcessBytes(buf.Bytes(), map[string]any{"name": "Test"})
		})
	}
}

func TestDocumentService_ProcessBytes_NilDataMap(t *testing.T) {
	t.Parallel()

	svc := NewDocumentService(nil)
	templateData := createMinimalDocx(t)

	// Should not panic with nil data map
	result, err := svc.ProcessBytes(templateData, nil)
	if err != nil {
		t.Fatalf("unexpected error with nil data map: %v", err)
	}
	if len(result) == 0 {
		t.Error("result should not be empty")
	}
}

func TestDocumentService_ProcessBytes_EmptyDataMap(t *testing.T) {
	t.Parallel()

	svc := NewDocumentService(nil)
	templateData := createMinimalDocx(t)

	// Empty data map should leave placeholders as-is but not error
	result, err := svc.ProcessBytes(templateData, map[string]any{})
	if err != nil {
		t.Fatalf("unexpected error with empty data map: %v", err)
	}
	if len(result) == 0 {
		t.Error("result should not be empty")
	}
}

func TestDocumentService_ProcessFromStorageToBytes_EmptyKey(t *testing.T) {
	t.Parallel()

	storage := &mockStorageReadWriter{
		readFunc: func(ctx context.Context, containerName, objectKey string) ([]byte, error) {
			if objectKey == "" {
				return nil, errors.New("object key is empty")
			}
			return createMinimalDocx(t), nil
		},
		writeFunc: func(ctx context.Context, containerName, objectKey string, data []byte) error {
			t.Fatal("write should not be called for ToBytes")
			return nil
		},
	}

	svc := NewDocumentService(storage)

	_, err := svc.ProcessFromStorageToBytes(
		context.Background(),
		"templates", "",
		map[string]any{"name": "Test"},
	)
	if err == nil {
		t.Fatal("expected error for empty template key")
	}
}

// createMinimalDocx builds a minimal valid DOCX archive in memory.
// It contains just enough structure for doctemplate.ProcessTemplate to succeed.
func createMinimalDocx(t *testing.T) []byte {
	t.Helper()

	contentTypesXML := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/><Default Extension="xml" ContentType="application/xml"/><Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/></Types>`

	relsXML := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/></Relationships>`

	documentXML := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
<w:body>
<w:p><w:r><w:t>Hello {{name}}</w:t></w:r></w:p>
</w:body>
</w:document>`

	documentRelsXML := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"></Relationships>`

	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	files := map[string]string{
		"[Content_Types].xml":          contentTypesXML,
		"_rels/.rels":                  relsXML,
		"word/document.xml":            documentXML,
		"word/_rels/document.xml.rels": documentRelsXML,
	}

	for name, content := range files {
		writer, err := w.Create(name)
		if err != nil {
			t.Fatalf("failed to create zip entry %s: %v", name, err)
		}
		if _, err := writer.Write([]byte(content)); err != nil {
			t.Fatalf("failed to write zip entry %s: %v", name, err)
		}
	}

	if err := w.Close(); err != nil {
		t.Fatalf("failed to close zip: %v", err)
	}

	return buf.Bytes()
}
