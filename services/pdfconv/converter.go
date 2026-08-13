package pdfconv

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// LibreOffice startup on the supported macOS development lane can take around
// 30 seconds before it starts laying out a generated document. Grade-sheet
// roster tables add enough work to exceed the former 60-second ceiling even
// though conversion is still progressing normally. Keep the native process
// bounded, but allow a complete heavy report to finish.
const libreOfficeConversionTimeout = 120 * time.Second

type conversionDeps struct {
	findBinary func() (string, error)
	runCommand func(context.Context, string, ...string) error
	timeout    time.Duration
}

func productionConversionDeps() conversionDeps {
	return conversionDeps{
		findBinary: findLibreOffice,
		runCommand: func(ctx context.Context, name string, args ...string) error {
			return exec.CommandContext(ctx, name, args...).Run()
		},
		timeout: libreOfficeConversionTimeout,
	}
}

// ConvertDocxToPDF converts DOCX bytes to PDF bytes using LibreOffice headless.
// It auto-detects the OS to find the LibreOffice binary.
// If LibreOffice is not installed, it returns the original DOCX bytes with a false flag.
func ConvertDocxToPDF(docxBytes []byte) (pdfBytes []byte, ok bool, err error) {
	return convertDocxToPDFWithDeps(docxBytes, productionConversionDeps())
}

func convertDocxToPDFWithDeps(docxBytes []byte, deps conversionDeps) (pdfBytes []byte, ok bool, err error) {
	binary, err := deps.findBinary()
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

	// Give every conversion an isolated LibreOffice user profile. Reusing the
	// interactive/default profile can block indefinitely on a stale singleton
	// lock, and concurrent requests would otherwise contend for the same office
	// process. The deadline bounds native-process failure on every provider lane.
	profileDir := filepath.Join(tmpDir, "libreoffice-profile")
	profileURL := (&url.URL{Scheme: "file", Path: filepath.ToSlash(profileDir)}).String()
	ctx, cancel := context.WithTimeout(context.Background(), deps.timeout)
	defer cancel()
	args := []string{
		"-env:UserInstallation=" + profileURL,
		"--headless",
		"--nologo",
		"--nodefault",
		"--nofirststartwizard",
		"--nolockcheck",
		"--norestore",
		"--convert-to", "pdf",
		"--outdir", tmpDir,
		docxPath,
	}

	if err := deps.runCommand(ctx, binary, args...); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, false, fmt.Errorf("pdfconv: libreoffice conversion timed out")
		}
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
