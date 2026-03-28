package fycha

import (
	"context"
	"fmt"

	"github.com/erniealice/fycha-golang/services/doctemplate"
	"github.com/erniealice/fycha-golang/services/pdfconv"
)

// StorageReadWriter reads and writes objects from a storage backend.
// Implementations wrap provider-specific adapters (e.g., espyna StorageAdapter)
// to keep fycha provider-agnostic.
//
// This extends StorageReader (used by StorageHandler for HTTP file serving)
// with write capability needed for document generation output.
type StorageReadWriter interface {
	ReadObject(ctx context.Context, containerName, objectKey string) ([]byte, error)
	WriteObject(ctx context.Context, containerName, objectKey string, data []byte) error
}

// DocumentService orchestrates document template processing with storage I/O.
// It combines the pure doctemplate engine (bytes in/out) with storage
// for reading templates and writing results.
//
// The core processing is delegated to services/doctemplate.ProcessTemplate,
// which has zero I/O dependencies. This service adds the storage layer.
type DocumentService struct {
	storage StorageReadWriter
}

// NewDocumentService creates a DocumentService with the given storage backend.
// Pass nil for storage if you only need ProcessBytes (no storage I/O).
func NewDocumentService(storage StorageReadWriter) *DocumentService {
	return &DocumentService{storage: storage}
}

// ProcessBytes processes a DOCX template from raw bytes and returns the result as DOCX bytes.
// This is a convenience wrapper around doctemplate.ProcessTemplate — no storage needed.
func (s *DocumentService) ProcessBytes(templateData []byte, data map[string]any) ([]byte, error) {
	return doctemplate.ProcessTemplate(templateData, data)
}

// ProcessBytesToPDF processes a DOCX template and converts the result to PDF.
// Returns the PDF bytes, or an error if conversion fails.
// Requires LibreOffice to be installed (https://www.libreoffice.org/download/).
func (s *DocumentService) ProcessBytesToPDF(templateData []byte, data map[string]any) ([]byte, error) {
	docxBytes, err := doctemplate.ProcessTemplate(templateData, data)
	if err != nil {
		return nil, err
	}
	return convertToPDF(docxBytes)
}

// ProcessFromStorage reads a template from storage, processes it with the given data,
// and writes the result back to storage as DOCX.
//
// Parameters:
//   - templateContainer: storage container/bucket where the template lives
//   - templateKey: object key for the template (e.g., "templates/invoice.docx")
//   - outputContainer: storage container for the result (can be same as template container)
//   - outputKey: object key for the output (e.g., "output/invoice-123.docx")
//   - data: JSON-like map for placeholder replacement
func (s *DocumentService) ProcessFromStorage(
	ctx context.Context,
	templateContainer, templateKey string,
	outputContainer, outputKey string,
	data map[string]any,
) error {
	if s.storage == nil {
		return fmt.Errorf("storage not configured")
	}

	// Read template from storage
	templateData, err := s.storage.ReadObject(ctx, templateContainer, templateKey)
	if err != nil {
		return fmt.Errorf("reading template %s/%s: %w", templateContainer, templateKey, err)
	}

	// Process template
	result, err := doctemplate.ProcessTemplate(templateData, data)
	if err != nil {
		return fmt.Errorf("processing template: %w", err)
	}

	// Write result to storage
	if err := s.storage.WriteObject(ctx, outputContainer, outputKey, result); err != nil {
		return fmt.Errorf("writing output %s/%s: %w", outputContainer, outputKey, err)
	}

	return nil
}

// ProcessFromStorageToPDF reads a template from storage, processes it,
// converts to PDF, and writes the PDF back to storage.
func (s *DocumentService) ProcessFromStorageToPDF(
	ctx context.Context,
	templateContainer, templateKey string,
	outputContainer, outputKey string,
	data map[string]any,
) error {
	if s.storage == nil {
		return fmt.Errorf("storage not configured")
	}

	templateData, err := s.storage.ReadObject(ctx, templateContainer, templateKey)
	if err != nil {
		return fmt.Errorf("reading template %s/%s: %w", templateContainer, templateKey, err)
	}

	docxBytes, err := doctemplate.ProcessTemplate(templateData, data)
	if err != nil {
		return fmt.Errorf("processing template: %w", err)
	}

	pdfBytes, err := convertToPDF(docxBytes)
	if err != nil {
		return fmt.Errorf("converting to PDF: %w", err)
	}

	if err := s.storage.WriteObject(ctx, outputContainer, outputKey, pdfBytes); err != nil {
		return fmt.Errorf("writing output %s/%s: %w", outputContainer, outputKey, err)
	}

	return nil
}

// ProcessFromStorageToBytes reads a template from storage, processes it,
// and returns the result as DOCX bytes (for streaming to HTTP response, etc.).
func (s *DocumentService) ProcessFromStorageToBytes(
	ctx context.Context,
	templateContainer, templateKey string,
	data map[string]any,
) ([]byte, error) {
	if s.storage == nil {
		return nil, fmt.Errorf("storage not configured")
	}

	// Read template from storage
	templateData, err := s.storage.ReadObject(ctx, templateContainer, templateKey)
	if err != nil {
		return nil, fmt.Errorf("reading template %s/%s: %w", templateContainer, templateKey, err)
	}

	// Process and return
	result, err := doctemplate.ProcessTemplate(templateData, data)
	if err != nil {
		return nil, fmt.Errorf("processing template: %w", err)
	}

	return result, nil
}

// ProcessFromStorageToPDFBytes reads a template from storage, processes it,
// converts to PDF, and returns the PDF bytes for streaming to HTTP response.
func (s *DocumentService) ProcessFromStorageToPDFBytes(
	ctx context.Context,
	templateContainer, templateKey string,
	data map[string]any,
) ([]byte, error) {
	if s.storage == nil {
		return nil, fmt.Errorf("storage not configured")
	}

	templateData, err := s.storage.ReadObject(ctx, templateContainer, templateKey)
	if err != nil {
		return nil, fmt.Errorf("reading template %s/%s: %w", templateContainer, templateKey, err)
	}

	docxBytes, err := doctemplate.ProcessTemplate(templateData, data)
	if err != nil {
		return nil, fmt.Errorf("processing template: %w", err)
	}

	return convertToPDF(docxBytes)
}

// convertToPDF converts DOCX bytes to PDF using LibreOffice.
// Returns an error if LibreOffice is not installed (no silent fallback).
func convertToPDF(docxBytes []byte) ([]byte, error) {
	pdfBytes, ok, err := pdfconv.ConvertDocxToPDF(docxBytes)
	if err != nil {
		return nil, fmt.Errorf("PDF conversion failed: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf("PDF conversion unavailable: LibreOffice is not installed (see https://www.libreoffice.org/download/)")
	}
	return pdfBytes, nil
}
