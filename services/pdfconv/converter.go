package pdfconv

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// ConvertDocxToPDF converts DOCX bytes to PDF bytes using LibreOffice headless.
// It auto-detects the OS to find the LibreOffice binary.
// If LibreOffice is not installed, it returns the original DOCX bytes with a false flag.
func ConvertDocxToPDF(docxBytes []byte) (pdfBytes []byte, ok bool, err error) {
	binary, err := findLibreOffice()
	if err != nil {
		log.Printf("pdfconv: LibreOffice not found, falling back to DOCX: %v", err)
		return docxBytes, false, nil
	}

	// Create temp directory for the conversion
	tmpDir, err := os.MkdirTemp("", "pdfconv-*")
	if err != nil {
		return nil, false, fmt.Errorf("pdfconv: creating temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Write DOCX to temp file
	docxPath := filepath.Join(tmpDir, "input.docx")
	if err := os.WriteFile(docxPath, docxBytes, 0644); err != nil {
		return nil, false, fmt.Errorf("pdfconv: writing temp docx: %w", err)
	}

	// Run LibreOffice headless conversion
	cmd := exec.Command(binary,
		"--headless",
		"--norestore",
		"--convert-to", "pdf",
		"--outdir", tmpDir,
		docxPath,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, false, fmt.Errorf("pdfconv: libreoffice conversion failed: %w", err)
	}

	// Read the resulting PDF
	pdfPath := filepath.Join(tmpDir, "input.pdf")
	pdfBytes, err = os.ReadFile(pdfPath)
	if err != nil {
		return nil, false, fmt.Errorf("pdfconv: reading converted PDF: %w", err)
	}

	return pdfBytes, true, nil
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

	// Check common Linux paths
	candidates := []string{
		"/usr/bin/soffice",
		"/usr/bin/libreoffice",
		"/usr/local/bin/soffice",
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
