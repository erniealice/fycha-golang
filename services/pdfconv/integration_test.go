package pdfconv

import (
	"archive/zip"
	"bytes"
	"context"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

// minimalDocx builds the smallest valid WordprocessingML package.
func minimalDocx(t *testing.T, text string) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for name, body := range map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/><Default Extension="xml" ContentType="application/xml"/><Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/></Types>`,
		"_rels/.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/></Relationships>`,
		"word/document.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:body><w:p><w:r><w:t>` + text + `</w:t></w:r></w:p></w:body></w:document>`,
	} {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write([]byte(body)); err != nil {
			t.Fatal(err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// TestIntegration_ConcurrentRealConversions runs several real LibreOffice
// instances at once (isolated profiles/HOME/TMPDIR) and checks every one yields
// a PDF and no soffice process referencing the test's temp root survives.
// Opt-in: PDFCONV_INTEGRATION=1 (LibreOffice cold start is ~30s on macOS).
func TestIntegration_ConcurrentRealConversions(t *testing.T) {
	if os.Getenv("PDFCONV_INTEGRATION") != "1" {
		t.Skip("set PDFCONV_INTEGRATION=1 to run real LibreOffice conversions")
	}
	binary, err := findLibreOffice()
	if err != nil {
		t.Skipf("LibreOffice unavailable: %v", err)
	}
	t.Logf("using %s", binary)

	workers := 3
	if n, err := strconv.Atoi(os.Getenv("PDFCONV_INTEGRATION_WORKERS")); err == nil && n > 0 {
		workers = n
	}
	deps := productionConversionDeps()
	deps.limiter = newLimiter(workers) // true concurrency: every worker runs its own soffice tree
	deps.faults = &faultTracker{}
	deps.tempRoot = t.TempDir()

	var wg sync.WaitGroup
	errs := make([]error, workers)
	pdfs := make([][]byte, workers)
	started := time.Now()
	for i := range workers {
		wg.Go(func() {
			var ok bool
			pdfs[i], ok, errs[i] = convertDocxToPDFWithDeps(context.Background(), minimalDocx(t, "pdfconv integration "+string(rune('A'+i))), deps)
			if errs[i] == nil && !ok {
				t.Errorf("worker %d: ok=false with nil error", i)
			}
		})
	}
	wg.Wait()
	t.Logf("%d concurrent conversions took %s", workers, time.Since(started))

	for i := range workers {
		if errs[i] != nil {
			t.Fatalf("worker %d: %v", i, errs[i])
		}
		if !bytes.HasPrefix(pdfs[i], pdfMagic) {
			t.Fatalf("worker %d: output is not a PDF (%d bytes)", i, len(pdfs[i]))
		}
	}
	if h := deps.faults.snapshot(deps.limiter); !h.Healthy || h.Conversions != int64(workers) || h.Failures != 0 {
		t.Fatalf("health = %+v", h)
	}

	// No LibreOffice process from this test may outlive its conversion.
	out, _ := exec.Command("ps", "-axo", "pid=,command=").Output()
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, deps.tempRoot) {
			t.Fatalf("leftover process references the test temp root: %s", strings.TrimSpace(line))
		}
	}
}
