package doctemplate

import (
	"fmt"

	"github.com/beevik/etree"
)

// ProcessTemplate takes a DOCX template as bytes and a data map,
// performs placeholder replacement, and returns the processed DOCX as bytes.
func ProcessTemplate(templateData []byte, data map[string]any) ([]byte, error) {
	// Step 1: Read the DOCX archive
	archive, err := ReadDocxBytes(templateData)
	if err != nil {
		return nil, fmt.Errorf("reading docx: %w", err)
	}

	// Step 2: Process the main document content
	processedContent, err := processXMLContent(archive.Content, data)
	if err != nil {
		return nil, fmt.Errorf("processing document.xml: %w", err)
	}

	// Step 3: Process headers
	processedHeaders := make(map[string]string, len(archive.Headers))
	for name, content := range archive.Headers {
		processed, err := processXMLContent(content, data)
		if err != nil {
			return nil, fmt.Errorf("processing header %s: %w", name, err)
		}
		processedHeaders[name] = processed
	}

	// Step 4: Process footers
	processedFooters := make(map[string]string, len(archive.Footers))
	for name, content := range archive.Footers {
		processed, err := processXMLContent(content, data)
		if err != nil {
			return nil, fmt.Errorf("processing footer %s: %w", name, err)
		}
		processedFooters[name] = processed
	}

	// Step 5: Write the modified DOCX
	return archive.WriteDocx(processedContent, processedHeaders, processedFooters)
}

// processXMLContent parses an OOXML string, processes placeholders using
// the etree-based XML processor, and serializes back to string.
func processXMLContent(xmlContent string, data map[string]any) (string, error) {
	doc := etree.NewDocument()
	doc.ReadSettings.PreserveCData = true

	if err := doc.ReadFromString(xmlContent); err != nil {
		return "", fmt.Errorf("parsing XML: %w", err)
	}

	// Find the body element — it may be w:body inside w:document,
	// or for headers/footers the root structure differs.
	body := doc.FindElement("//body")
	if body == nil {
		// Try alternative: process all paragraphs/tables at document level
		root := doc.Root()
		if root != nil {
			processAllElements(root, data)
		}
	} else {
		ProcessBody(body, data)
	}

	doc.WriteSettings.CanonicalEndTags = false
	doc.WriteSettings.CanonicalText = false
	doc.WriteSettings.CanonicalAttrVal = false

	result, err := doc.WriteToString()
	if err != nil {
		return "", fmt.Errorf("serializing XML: %w", err)
	}

	return result, nil
}

// processAllElements walks all child elements recursively, processing
// paragraphs and tables it finds. Used for headers/footers where the
// structure differs from the main document body.
func processAllElements(el *etree.Element, data map[string]any) {
	for _, child := range el.ChildElements() {
		switch child.Tag {
		case "p":
			processParagraph(child, data)
		case "tbl":
			processTable(child, data)
		default:
			processAllElements(child, data)
		}
	}
}
