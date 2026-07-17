package doctemplate

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"strings"
)

// Decompression limits for ReadDocxBytes. These are authoritative at the render
// boundary and independent of any upload-time validation: the archive bytes may
// have been supplied by an alternate writer or a replaced storage object, so the
// render path must never trust a declared (header) size and must cap the ACTUAL
// expanded bytes it reads.
//
//   - maxDocxEntries          bounds the number of ZIP entries (decompression-bomb
//     fan-out and per-entry allocation overhead).
//   - maxDocxEntryExpanded    bounds the ACTUAL expanded size of any single entry.
//   - maxDocxTotalExpanded    bounds the ACTUAL aggregate expanded size across all
//     entries, so many moderately-sized entries cannot sum to an OOM.
const (
	maxDocxEntries       = 2000
	maxDocxEntryExpanded = 64 << 20  // 64 MiB
	maxDocxTotalExpanded = 256 << 20 // 256 MiB
)

// DocxArchive holds the content and structure of a DOCX file.
type DocxArchive struct {
	Content string
	Headers map[string]string
	Footers map[string]string
	Images  map[string][]byte
	files   []*zip.File
}

// ReadDocxBytes reads a DOCX file from a byte slice and extracts its main
// components. It enforces authoritative decompression limits (entry count,
// per-entry and aggregate expanded size) and rejects unsafe entry paths so a
// crafted or replaced archive cannot exhaust memory or escape the archive root.
func ReadDocxBytes(data []byte) (*DocxArchive, error) {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}

	if len(reader.File) > maxDocxEntries {
		return nil, fmt.Errorf("docx: archive has %d entries, exceeds limit of %d", len(reader.File), maxDocxEntries)
	}

	archive := &DocxArchive{
		Headers: make(map[string]string),
		Footers: make(map[string]string),
		Images:  make(map[string][]byte),
		files:   reader.File,
	}

	var totalExpanded int64
	for _, file := range archive.files {
		if err := validateEntryPath(file.Name); err != nil {
			return nil, err
		}

		remaining := int64(maxDocxTotalExpanded) - totalExpanded
		contentBytes, err := readZipFileLimited(file, maxDocxEntryExpanded, remaining)
		if err != nil {
			return nil, err
		}
		totalExpanded += int64(len(contentBytes))

		switch {
		case file.Name == "word/document.xml":
			archive.Content = string(contentBytes)
		case strings.HasPrefix(file.Name, "word/header"):
			archive.Headers[file.Name] = string(contentBytes)
		case strings.HasPrefix(file.Name, "word/footer"):
			archive.Footers[file.Name] = string(contentBytes)
		case strings.HasPrefix(file.Name, "word/media/"):
			archive.Images[file.Name] = contentBytes
		}
	}

	return archive, nil
}

// validateEntryPath rejects absolute paths and parent-directory traversal in a
// ZIP entry name. Legitimate DOCX entries are workspace-relative (e.g.
// "word/document.xml"), so this never rejects a well-formed template.
func validateEntryPath(name string) error {
	if name == "" {
		return fmt.Errorf("docx: archive entry has an empty name")
	}
	norm := strings.ReplaceAll(name, "\\", "/")
	if strings.HasPrefix(norm, "/") {
		return fmt.Errorf("docx: archive entry %q uses an absolute path", name)
	}
	// Drive-letter absolute path (e.g. "C:/...").
	if len(norm) >= 2 && norm[1] == ':' {
		return fmt.Errorf("docx: archive entry %q uses an absolute path", name)
	}
	for _, seg := range strings.Split(norm, "/") {
		if seg == ".." {
			return fmt.Errorf("docx: archive entry %q escapes the archive root", name)
		}
	}
	return nil
}

// readZipFileLimited reads a ZIP entry, capping the ACTUAL expanded bytes at
// min(maxEntry, remainingTotal). It reads one byte past the effective cap so an
// over-limit expansion is detected by bytes actually read, never by the declared
// (and forgeable) UncompressedSize. A cap hit reports whether the per-entry or
// the aggregate limit was reached.
func readZipFileLimited(f *zip.File, maxEntry, remainingTotal int64) ([]byte, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	effectiveCap := maxEntry
	aggregate := false
	if remainingTotal < effectiveCap {
		effectiveCap = remainingTotal
		aggregate = true
	}

	content, err := io.ReadAll(io.LimitReader(rc, effectiveCap+1))
	if err != nil {
		return nil, err
	}
	if int64(len(content)) > effectiveCap {
		if aggregate {
			return nil, fmt.Errorf("docx: entry %q pushes total expanded size past limit of %d bytes", f.Name, int64(maxDocxTotalExpanded))
		}
		return nil, fmt.Errorf("docx: entry %q expands past per-entry limit of %d bytes", f.Name, maxEntry)
	}
	return content, nil
}

func readZipFile(f *zip.File) ([]byte, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return io.ReadAll(rc)
}

// WriteDocx creates a new DOCX file as a byte slice with modified content.
// It takes the original archive and the modified text content for the main document, headers, and footers.
func (archive *DocxArchive) WriteDocx(
	modifiedContent string,
	modifiedHeaders map[string]string,
	modifiedFooters map[string]string,
) ([]byte, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	for _, file := range archive.files {
		writer, err := zipWriter.Create(file.Name)
		if err != nil {
			return nil, err
		}

		var contentToWrite []byte
		var found bool

		// Check if the current file is one of the modified ones
		if file.Name == "word/document.xml" {
			contentToWrite = []byte(modifiedContent)
			found = true
		} else if content, ok := modifiedHeaders[file.Name]; ok {
			contentToWrite = []byte(content)
			found = true
		} else if content, ok := modifiedFooters[file.Name]; ok {
			contentToWrite = []byte(content)
			found = true
		}

		if found {
			if _, err := writer.Write(contentToWrite); err != nil {
				return nil, err
			}
		} else {
			original, err := readZipFile(file)
			if err != nil {
				return nil, err
			}
			if _, err := writer.Write(original); err != nil {
				return nil, err
			}
		}
	}

	if err := zipWriter.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
