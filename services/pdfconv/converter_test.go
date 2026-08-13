package pdfconv

import (
	"bytes"
	"context"
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestConvertDocxToPDFWithDeps_SuccessUsesIsolatedProfile(t *testing.T) {
	t.Parallel()

	var gotBinary string
	var gotArgs []string
	deps := conversionDeps{
		findBinary: func() (string, error) {
			return "/test/soffice", nil
		},
		runCommand: func(_ context.Context, binary string, args ...string) error {
			gotBinary = binary
			gotArgs = append([]string(nil), args...)
			if len(args) < 2 {
				return errors.New("missing conversion arguments")
			}
			outDir := args[len(args)-2]
			return os.WriteFile(filepath.Join(outDir, "input.pdf"), []byte("%PDF-test"), 0o600)
		},
		timeout: time.Second,
	}

	pdf, ok, err := convertDocxToPDFWithDeps([]byte("docx"), deps)
	if err != nil {
		t.Fatalf("convertDocxToPDFWithDeps() error = %v", err)
	}
	if !ok {
		t.Fatal("convertDocxToPDFWithDeps() ok = false, want true")
	}
	if !bytes.Equal(pdf, []byte("%PDF-test")) {
		t.Fatalf("convertDocxToPDFWithDeps() bytes = %q, want %%PDF-test", pdf)
	}
	if gotBinary != "/test/soffice" {
		t.Fatalf("binary = %q, want /test/soffice", gotBinary)
	}
	if len(gotArgs) != 12 {
		t.Fatalf("argument count = %d, want 12: %q", len(gotArgs), gotArgs)
	}

	outDir := gotArgs[10]
	profileURL := (&url.URL{
		Scheme: "file",
		Path:   filepath.ToSlash(filepath.Join(outDir, "libreoffice-profile")),
	}).String()
	wantArgs := []string{
		"-env:UserInstallation=" + profileURL,
		"--headless",
		"--nologo",
		"--nodefault",
		"--nofirststartwizard",
		"--nolockcheck",
		"--norestore",
		"--convert-to", "pdf",
		"--outdir", outDir,
		filepath.Join(outDir, "input.docx"),
	}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Fatalf("arguments = %#v, want %#v", gotArgs, wantArgs)
	}
}

func TestConvertDocxToPDFWithDeps_UnavailableFallsBack(t *testing.T) {
	t.Parallel()

	docx := []byte("original-docx")
	runCalled := false
	deps := conversionDeps{
		findBinary: func() (string, error) {
			return "", errors.New("not installed")
		},
		runCommand: func(context.Context, string, ...string) error {
			runCalled = true
			return nil
		},
		timeout: time.Second,
	}

	got, ok, err := convertDocxToPDFWithDeps(docx, deps)
	if err != nil {
		t.Fatalf("convertDocxToPDFWithDeps() error = %v", err)
	}
	if ok {
		t.Fatal("convertDocxToPDFWithDeps() ok = true, want false")
	}
	if !bytes.Equal(got, docx) {
		t.Fatalf("fallback bytes = %q, want %q", got, docx)
	}
	if runCalled {
		t.Fatal("runCommand called after unavailable binary")
	}
}

func TestConvertDocxToPDFWithDeps_CommandFailure(t *testing.T) {
	t.Parallel()

	deps := conversionDeps{
		findBinary: func() (string, error) {
			return "/test/soffice", nil
		},
		runCommand: func(context.Context, string, ...string) error {
			return errors.New("exit status 1")
		},
		timeout: time.Second,
	}

	pdf, ok, err := convertDocxToPDFWithDeps([]byte("docx"), deps)
	if err == nil || !strings.Contains(err.Error(), "libreoffice conversion failed") {
		t.Fatalf("error = %v, want conversion failure", err)
	}
	if ok || pdf != nil {
		t.Fatalf("result = (%q, %v), want (nil, false)", pdf, ok)
	}
}

func TestConvertDocxToPDFWithDeps_Timeout(t *testing.T) {
	t.Parallel()

	deps := conversionDeps{
		findBinary: func() (string, error) {
			return "/test/soffice", nil
		},
		runCommand: func(ctx context.Context, _ string, _ ...string) error {
			<-ctx.Done()
			return ctx.Err()
		},
		timeout: 10 * time.Millisecond,
	}

	pdf, ok, err := convertDocxToPDFWithDeps([]byte("docx"), deps)
	if err == nil || !strings.Contains(err.Error(), "libreoffice conversion timed out") {
		t.Fatalf("error = %v, want timeout classification", err)
	}
	if ok || pdf != nil {
		t.Fatalf("result = (%q, %v), want (nil, false)", pdf, ok)
	}
}
