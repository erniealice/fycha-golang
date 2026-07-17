package doctemplate

import (
	"archive/zip"
	"bytes"
	"io"
	"strings"
	"testing"
)

// zeroReader yields an unbounded stream of zero bytes. Used to write large,
// highly-compressible ZIP entries without allocating the expanded payload in the
// test process.
type zeroReader struct{}

func (zeroReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = 0
	}
	return len(p), nil
}

// buildZip assembles an in-memory ZIP from name→size entries (bytes of zeros).
func buildZip(t *testing.T, entries []struct {
	name string
	size int64
}) []byte {
	t.Helper()
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)
	for _, e := range entries {
		fw, err := w.Create(e.name)
		if err != nil {
			t.Fatalf("create zip entry %q: %v", e.name, err)
		}
		if e.size > 0 {
			if _, err := io.CopyN(fw, zeroReader{}, e.size); err != nil {
				t.Fatalf("write zip entry %q: %v", e.name, err)
			}
		}
	}
	if err := w.Close(); err != nil {
		t.Fatalf("close zip writer: %v", err)
	}
	return buf.Bytes()
}

// TestReadDocxBytes_NormalTemplatePasses proves the limits do not reject a
// well-formed template.
func TestReadDocxBytes_NormalTemplatePasses(t *testing.T) {
	doc := createTestDocx(t, bodyDoc(`<w:p><w:r><w:t>hello</w:t></w:r></w:p>`))
	archive, err := ReadDocxBytes(doc)
	if err != nil {
		t.Fatalf("normal template rejected: %v", err)
	}
	if !strings.Contains(archive.Content, "hello") {
		t.Errorf("expected document content to survive, got %q", archive.Content)
	}
}

// TestReadDocxBytes_OversizedEntryRejected proves a single entry whose ACTUAL
// expanded size exceeds the per-entry cap is rejected (declared size is never
// trusted; the read is bounded by bytes actually produced).
func TestReadDocxBytes_OversizedEntryRejected(t *testing.T) {
	data := buildZip(t, []struct {
		name string
		size int64
	}{
		{"word/document.xml", maxDocxEntryExpanded + 1},
	})
	_, err := ReadDocxBytes(data)
	if err == nil {
		t.Fatalf("expected per-entry expansion limit error, got nil")
	}
	if !strings.Contains(err.Error(), "per-entry limit") {
		t.Errorf("expected per-entry limit error, got: %v", err)
	}
}

// TestReadDocxBytes_TooManyEntriesRejected proves the entry-count cap.
func TestReadDocxBytes_TooManyEntriesRejected(t *testing.T) {
	entries := make([]struct {
		name string
		size int64
	}, maxDocxEntries+1)
	for i := range entries {
		entries[i].name = "part/" + itoa(i) + ".xml"
		entries[i].size = 0
	}
	_, err := ReadDocxBytes(buildZip(t, entries))
	if err == nil {
		t.Fatalf("expected too-many-entries error, got nil")
	}
	if !strings.Contains(err.Error(), "exceeds limit") {
		t.Errorf("expected entry-count limit error, got: %v", err)
	}
}

// TestReadDocxBytes_AggregateOverflowRejected proves that many entries under the
// per-entry cap still cannot sum past the aggregate cap. Entries are not stored
// (non word/* names) so peak test memory stays near one entry.
func TestReadDocxBytes_AggregateOverflowRejected(t *testing.T) {
	const each = int64(maxDocxEntryExpanded) // 64 MiB
	// 5 * 64 MiB = 320 MiB > 256 MiB aggregate cap; the entry that crosses the
	// remaining budget is bounded to remaining+1 bytes.
	entries := make([]struct {
		name string
		size int64
	}, 5)
	for i := range entries {
		entries[i].name = "filler/" + itoa(i) + ".bin"
		entries[i].size = each
	}
	_, err := ReadDocxBytes(buildZip(t, entries))
	if err == nil {
		t.Fatalf("expected aggregate expansion limit error, got nil")
	}
	if !strings.Contains(err.Error(), "total expanded size") {
		t.Errorf("expected aggregate limit error, got: %v", err)
	}
}

// TestReadDocxBytes_TraversalPathRejected proves a parent-directory entry name is
// rejected before its bytes are read.
func TestReadDocxBytes_TraversalPathRejected(t *testing.T) {
	data := buildZip(t, []struct {
		name string
		size int64
	}{
		{"../evil.xml", 4},
	})
	_, err := ReadDocxBytes(data)
	if err == nil {
		t.Fatalf("expected traversal-path error, got nil")
	}
	if !strings.Contains(err.Error(), "escapes the archive root") {
		t.Errorf("expected traversal error, got: %v", err)
	}
}

// TestReadDocxBytes_AbsolutePathRejected proves an absolute entry name is
// rejected.
func TestReadDocxBytes_AbsolutePathRejected(t *testing.T) {
	data := buildZip(t, []struct {
		name string
		size int64
	}{
		{"/etc/passwd", 4},
	})
	_, err := ReadDocxBytes(data)
	if err == nil {
		t.Fatalf("expected absolute-path error, got nil")
	}
	if !strings.Contains(err.Error(), "absolute path") {
		t.Errorf("expected absolute-path error, got: %v", err)
	}
}

// itoa is a tiny local int→string to avoid importing strconv for one call.
func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	var b [20]byte
	pos := len(b)
	for i > 0 {
		pos--
		b[pos] = byte('0' + i%10)
		i /= 10
	}
	return string(b[pos:])
}
