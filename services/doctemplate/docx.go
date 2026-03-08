package doctemplate

import (
	"archive/zip"
	"bytes"
	"io"
	"strings"
)

// DocxArchive holds the content and structure of a DOCX file.
type DocxArchive struct {
	Content string
	Headers map[string]string
	Footers map[string]string
	Images  map[string][]byte
	files   []*zip.File
}

// ReadDocxBytes reads a DOCX file from a byte slice and extracts its main components.
func ReadDocxBytes(data []byte) (*DocxArchive, error) {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}

	archive := &DocxArchive{
		Headers: make(map[string]string),
		Footers: make(map[string]string),
		Images:  make(map[string][]byte),
		files:   reader.File,
	}

	for _, file := range archive.files {
		contentBytes, err := readZipFile(file)
		if err != nil {
			return nil, err
		}

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
