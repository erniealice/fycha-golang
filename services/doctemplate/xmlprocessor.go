package doctemplate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/beevik/etree"
)

// Regular expressions to find placeholders and loop markers in the XML text.
var (
	// placeholderRegex matches simple placeholders like {{key.path}} or {{ key.path }}.
	placeholderRegex = regexp.MustCompile(`{{\s*([^#{}][^{}]*?)\s*}}`)
	// loopStartRegex matches loop start markers like {{#key}}.
	loopStartRegex = regexp.MustCompile(`{{\s*#\s*([^{}]+)\s*}}`)
	// loopEndRegex matches loop end markers like {{/key}}.
	loopEndRegex = regexp.MustCompile(`{{\s*/\s*([^{}]+)\s*}}`)
)

// getPathValue retrieves a nested value from a map using a dot-separated path, returning the raw value.
func getPathValue(data map[string]any, path string) (any, bool) {
	parts := strings.Split(path, ".")
	var current any = data

	for _, part := range parts {
		currentMap, ok := current.(map[string]any)
		if !ok {
			return nil, false
		}
		value, exists := currentMap[part]
		if !exists {
			return nil, false
		}
		current = value
	}
	return current, true
}

// getReplaceValue retrieves a nested value and returns it as a string.
func getReplaceValue(data map[string]any, path string) (string, bool) {
	val, ok := getPathValue(data, path)
	if !ok {
		return "", false
	}
	if val == nil {
		return "", true // Treat nil as an empty string
	}
	return fmt.Sprintf("%v", val), true
}

// extractPlaceholder extracts the key from a placeholder string, e.g., "{{propName}}" -> "propName".
func extractPlaceholder(text string) string {
	if matches := placeholderRegex.FindStringSubmatch(text); len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

// extractLoopMarker extracts the key from a loop marker, e.g., "{{#propName}}" -> "propName".
func extractLoopMarker(text string, prefix string) string {
	var re *regexp.Regexp
	if prefix == "#" {
		re = loopStartRegex
	} else if prefix == "/" {
		re = loopEndRegex
	} else {
		return ""
	}

	if matches := re.FindStringSubmatch(text); len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

// processParagraph handles placeholder replacement and loop detection within a <w:p> element.
// It performs cross-run text accumulation to correctly handle placeholders split by Word
// across multiple <w:r> elements (e.g., "{{" in one run, "client.name}}" in another).
//
// Returns:
//   - loopName: non-empty if a {{#key}} loop start marker was found
//   - endLoop: true if a {{/key}} loop end marker was found
func processParagraph(p *etree.Element, data map[string]any) (loopName string, endLoop bool) {
	allTextNodes := p.FindElements(".//t")
	if len(allTextNodes) == 0 {
		return "", false
	}

	var acc strings.Builder
	var accNodes []*etree.Element

	for _, t := range allTextNodes {
		text := t.Text()

		// If not accumulating and no placeholder opener, skip this node entirely
		if acc.Len() == 0 && !strings.Contains(text, "{{") {
			continue
		}

		// Add this node's text to the accumulator
		acc.WriteString(text)
		accNodes = append(accNodes, t)
		accumulated := acc.String()

		// Keep accumulating until we see a closing }}
		if !strings.Contains(accumulated, "}}") {
			continue
		}

		// We have at least one complete {{...}} — check what it is
		trimmed := strings.TrimSpace(accumulated)

		// Check for loop start marker: {{#key}}
		if m := loopStartRegex.FindStringSubmatch(trimmed); len(m) > 0 && m[0] == trimmed {
			clearNodes(accNodes)
			return strings.TrimSpace(m[1]), false
		}

		// Check for loop end marker: {{/key}}
		if m := loopEndRegex.FindStringSubmatch(trimmed); len(m) > 0 && m[0] == trimmed {
			clearNodes(accNodes)
			return "", true
		}

		// Replace all value placeholders {{key.path}} in the accumulated text
		replaced := replaceInText(accumulated, data)

		// Put the final result in the last text node, clear earlier ones
		accNodes[len(accNodes)-1].SetText(replaced)
		for i := 0; i < len(accNodes)-1; i++ {
			accNodes[i].SetText("")
		}

		// Reset for next potential placeholder group
		acc.Reset()
		accNodes = nil
	}

	return "", false
}

// replaceInText replaces all {{key.path}} placeholders in a string with values from data.
func replaceInText(text string, data map[string]any) string {
	return placeholderRegex.ReplaceAllStringFunc(text, func(match string) string {
		key := extractPlaceholder(match)
		if key == "" {
			return match
		}
		if val, ok := getReplaceValue(data, key); ok {
			return val
		}
		return match
	})
}

// clearNodes sets the text of all given elements to empty string.
func clearNodes(nodes []*etree.Element) {
	for _, n := range nodes {
		n.SetText("")
	}
}

// processElements recursively processes a slice of elements, applying templating logic.
func processElements(elements []*etree.Element, data map[string]any) {
	for _, el := range elements {
		switch el.Tag {
		case "p":
			processParagraph(el, data)
		case "tbl":
			processTable(el, data)
		}
	}
}

// processTable handles table row cloning for array loops and placeholder
// replacement in non-loop rows (e.g., header row, total row).
//
// Table structure expected:
//
//	<w:tbl>
//	  <w:tr>Header row (static)</w:tr>
//	  <w:tr>{{#items}}</w:tr>           ← start marker row
//	  <w:tr>{{description}} | {{amount}}</w:tr>  ← template row(s)
//	  <w:tr>{{/items}}</w:tr>           ← end marker row
//	  <w:tr>Total: {{total}}</w:tr>     ← static row with placeholders
//	</w:tbl>
func processTable(tbl *etree.Element, data map[string]any) {
	rows := tbl.FindElements("./tr")
	if len(rows) == 0 {
		return
	}

	// Scan rows to find loop markers
	loopKey := ""
	startIndex, endIndex := -1, -1

	for i, row := range rows {
		text := rowText(row)
		if key := extractLoopMarker(text, "#"); key != "" {
			loopKey = key
			startIndex = i
		} else if extractLoopMarker(text, "/") != "" && loopKey != "" {
			endIndex = i
			break
		}
	}

	// If no loop markers found, just process all rows for simple placeholders
	if startIndex == -1 || endIndex == -1 || loopKey == "" {
		for _, row := range rows {
			processRowCells(row, data)
		}
		return
	}

	// Process non-loop rows (header, total, etc.) for placeholder replacement
	for i, row := range rows {
		if i >= startIndex && i <= endIndex {
			continue // Skip loop marker and template rows
		}
		processRowCells(row, data)
	}

	// Collect template rows (between start and end markers, exclusive)
	var templateRows []*etree.Element
	for i := startIndex + 1; i < endIndex; i++ {
		templateRows = append(templateRows, rows[i])
	}

	if len(templateRows) == 0 {
		// No template rows between markers — just remove markers
		tbl.RemoveChild(rows[startIndex])
		tbl.RemoveChild(rows[endIndex])
		return
	}

	// Get the array data for looping
	rawLoopData, ok := getPathValue(data, loopKey)
	if !ok {
		// No data — remove all loop rows (markers + templates)
		for i := startIndex; i <= endIndex; i++ {
			tbl.RemoveChild(rows[i])
		}
		return
	}

	loopData, ok := rawLoopData.([]any)
	if !ok {
		// Try []map[string]any → []any conversion
		if mapSlice, isMapSlice := rawLoopData.([]map[string]any); isMapSlice {
			loopData = make([]any, len(mapSlice))
			for i, v := range mapSlice {
				loopData[i] = v
			}
		} else {
			return
		}
	}

	if len(loopData) == 0 {
		for i := startIndex; i <= endIndex; i++ {
			tbl.RemoveChild(rows[i])
		}
		return
	}

	// The anchor is the end marker row — we insert cloned rows before it
	anchor := rows[endIndex]

	// For each array item, clone template rows and process placeholders
	for _, itemData := range loopData {
		itemMap, ok := itemData.(map[string]any)
		if !ok {
			continue
		}

		for _, tmplRow := range templateRows {
			newRow := tmplRow.Copy()
			processRowCells(newRow, itemMap)
			tbl.InsertChild(anchor, newRow)
		}
	}

	// Remove original marker rows and template rows
	tbl.RemoveChild(rows[startIndex]) // {{#items}} row
	for _, tmplRow := range templateRows {
		tbl.RemoveChild(tmplRow)
	}
	tbl.RemoveChild(rows[endIndex]) // {{/items}} row
}

// processRowCells runs processParagraph on every cell's paragraphs in a table row.
func processRowCells(row *etree.Element, data map[string]any) {
	for _, cell := range row.FindElements(".//tc") {
		for _, p := range cell.FindElements("./p") {
			processParagraph(p, data)
		}
	}
}

// rowText concatenates all text content in a table row for marker detection.
func rowText(row *etree.Element) string {
	var sb strings.Builder
	for _, t := range row.FindElements(".//t") {
		sb.WriteString(t.Text())
	}
	return sb.String()
}

// ProcessBody is the main entry point for processing the document body.
// It uses clone-based iteration for body-level loops, supporting two-level
// nesting (body-level loop wrapping table-level loops).
//
// For body-level loops ({{#key}}...{{/key}} as standalone paragraphs):
//  1. Scans body children to find start/end marker paragraphs
//  2. Collects all elements between markers as the "template block"
//  3. Deep-copies the template block for each array item via element.Copy()
//  4. Processes each clone with the current item's data (resolving nested table loops)
//  5. Replaces the original markers and template block with processed clones
//  6. Processes remaining elements (before/after) with root data
func ProcessBody(body *etree.Element, data map[string]any) {
	bodyItems := body.ChildElements()

	// First pass: non-destructive scan for body-level loop markers
	startIdx, endIdx, loopKey := findBodyLoopMarkers(bodyItems)

	// No loop markers found — process all elements normally with root data
	if startIdx == -1 || endIdx == -1 {
		for _, el := range bodyItems {
			switch el.Tag {
			case "p":
				processParagraph(el, data)
			case "tbl":
				processTable(el, data)
			}
		}
		return
	}

	// Process elements before the loop with root data
	for i := 0; i < startIdx; i++ {
		el := bodyItems[i]
		switch el.Tag {
		case "p":
			processParagraph(el, data)
		case "tbl":
			processTable(el, data)
		}
	}

	// Collect the template block elements (between start and end markers, exclusive)
	var templateBlock []*etree.Element
	for i := startIdx + 1; i < endIdx; i++ {
		templateBlock = append(templateBlock, bodyItems[i])
	}

	// Get the array data for the loop key
	loopItems := resolveLoopData(data, loopKey)

	// The anchor for insertion is the end marker paragraph
	anchor := bodyItems[endIdx]

	// Clone and process template block for each array item.
	// Each clone gets processParagraph/processTable with the item's data,
	// so nested table loops (e.g., {{#items}} inside a table) resolve
	// against the current body-loop item's scope.
	for _, itemData := range loopItems {
		itemMap, ok := itemData.(map[string]any)
		if !ok {
			continue
		}
		for _, tmplEl := range templateBlock {
			clone := tmplEl.Copy()
			switch clone.Tag {
			case "p":
				processParagraph(clone, itemMap)
			case "tbl":
				processTable(clone, itemMap)
			}
			body.InsertChild(anchor, clone)
		}
	}

	// Remove original start marker, template block elements, and end marker
	body.RemoveChild(bodyItems[startIdx]) // {{#key}} paragraph
	for _, tmplEl := range templateBlock {
		body.RemoveChild(tmplEl)
	}
	body.RemoveChild(bodyItems[endIdx]) // {{/key}} paragraph

	// Process elements after the loop with root data.
	// Elements before the loop were already processed above.
	// Cloned elements were processed with item data during insertion.
	// Re-running processParagraph/processTable on those is safe because
	// already-replaced text won't match placeholder patterns.
	for i := endIdx + 1; i < len(bodyItems); i++ {
		el := bodyItems[i]
		switch el.Tag {
		case "p":
			processParagraph(el, data)
		case "tbl":
			processTable(el, data)
		}
	}
}

// findBodyLoopMarkers scans body child elements for the first {{#key}}/{{/key}} pair.
// It reads text without modifying the elements (non-destructive scan).
func findBodyLoopMarkers(elements []*etree.Element) (startIdx, endIdx int, key string) {
	startIdx = -1
	endIdx = -1

	for i, el := range elements {
		if el.Tag != "p" {
			continue
		}
		text := paragraphText(el)
		trimmed := strings.TrimSpace(text)

		if startIdx == -1 {
			if m := loopStartRegex.FindStringSubmatch(trimmed); len(m) > 0 && m[0] == trimmed {
				startIdx = i
				key = strings.TrimSpace(m[1])
			}
		} else {
			if m := loopEndRegex.FindStringSubmatch(trimmed); len(m) > 0 && m[0] == trimmed {
				endIdx = i
				return
			}
		}
	}
	return -1, -1, ""
}

// paragraphText concatenates all text content in a paragraph for marker detection,
// without modifying the element.
func paragraphText(p *etree.Element) string {
	var sb strings.Builder
	for _, t := range p.FindElements(".//t") {
		sb.WriteString(t.Text())
	}
	return sb.String()
}

// resolveLoopData extracts and normalizes the array data for a loop key.
func resolveLoopData(data map[string]any, key string) []any {
	raw, ok := getPathValue(data, key)
	if !ok {
		return nil
	}

	if ld, isSlice := raw.([]any); isSlice {
		return ld
	}
	if mapSlice, isMapSlice := raw.([]map[string]any); isMapSlice {
		result := make([]any, len(mapSlice))
		for i, v := range mapSlice {
			result[i] = v
		}
		return result
	}
	return nil
}
