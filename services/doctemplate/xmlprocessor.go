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
	// residualTokenRegex matches any complete {{...}} template construct (a loop
	// marker or an unresolved placeholder). It is used by the post-processing
	// zero-leak scrub. It deliberately forbids inner braces ([^{}]) so it can only
	// ever span a single, already-consolidated text node — see
	// scrubResidualTemplateTokens.
	residualTokenRegex = regexp.MustCompile(`{{[^{}]*}}`)
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

// scrubResidualTemplateTokens is the zero-leak backstop for the P7 all-parts
// {{/}} audit invariant. It removes every residual {{...}} construct left in the
// tree AFTER all resolvable placeholders were replaced and all well-formed loops
// expanded. Three leak vectors survive normal processing and are caught here:
//
//   - a malformed table (rejected by the complete marker-stream scan in
//     processTable and rendered as static rows) whose template placeholders
//     cannot resolve at root scope;
//   - a placeholder whose leaf is absent from the item/root data — the engine's
//     no-fallback law leaves it verbatim (replaceInText returns the match
//     unchanged), which the manifest-driven blank-seeding on the caller side is
//     built to pre-empt but which must never reach a rendered part;
//   - a loop/marker token embedded in mixed-content or split across runs that
//     processParagraph's exact-match branch does not blank.
//
// The scrub is deliberately node-scoped, not string-scoped over the serialized
// XML: processParagraph consolidates every completed {{...}} group into a single
// <w:t> node (its result goes into the last accumulated node, earlier ones are
// cleared), so residualTokenRegex — which forbids inner braces — always matches
// inside one text node's text. A string-level strip could match a token that
// visually spans element tags and delete the intervening markup; a node-level
// strip cannot. The engine never emits a lone unpaired "{{" (markers and
// placeholders always carry their "}}"), so removing complete tokens is
// sufficient to hold the zero-{{ invariant.
func scrubResidualTemplateTokens(root *etree.Element) {
	if root == nil {
		return
	}
	for _, t := range root.FindElements(".//t") {
		s := t.Text()
		if strings.Contains(s, "{{") {
			t.SetText(residualTokenRegex.ReplaceAllString(s, ""))
		}
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

	// Pair loop markers with a LIFO stack, scanning the COMPLETE marker stream
	// before selecting anything. A well-nested sequence selects the INNERMOST
	// complete pair (the first opener/closer pair to balance) as this table's row
	// loop; any enclosing markers are handled as ordinary non-loop rows and
	// blanked. The stream is only accepted after the WHOLE table has been proven
	// well-formed:
	//
	//   - every closer must match the nearest still-open opener of the SAME key
	//     (else a crossing/interleaved close, e.g. #outer,#inner,/outer,/inner);
	//   - no closer may appear with an empty stack (a stray/unmatched closer,
	//     whether leading like /ghost,#items,… or trailing like …,/items,/ghost);
	//   - no opener may remain unclosed at end of stream;
	//   - at most ONE top-level pair is supported (a second top-level opener/closer
	//     is rejected rather than silently ignored).
	//
	// Any of these makes the table malformed: it is treated as having no loop
	// (fail closed) — every row processed once, its marker cells blanked, its
	// non-resolving template placeholders removed by the residual-token scrub in
	// processXMLContent — rather than silently mispairing markers and rows or
	// expanding the wrong pair. This mirrors the fail-closed pairing the body-loop
	// scanner and the P4 hardening wave established for the render path, and the
	// pre-scan (nothing is mutated until the stream is proven valid) is what makes
	// the guarantee robust against the trailing-closer / unclosed-opener /
	// second-top-level shapes the old first-balanced-closer break missed.
	type openMarker struct {
		key   string
		index int
	}
	var stack []openMarker
	loopKey := ""
	startIndex, endIndex := -1, -1
	topLevelPairs := 0
	malformed := false

	for i, row := range rows {
		text := rowText(row)
		if key := extractLoopMarker(text, "#"); key != "" {
			stack = append(stack, openMarker{key: key, index: i})
			continue
		}
		if closeKey := extractLoopMarker(text, "/"); closeKey != "" {
			if len(stack) == 0 || stack[len(stack)-1].key != closeKey {
				// Closer with no matching open opener (stray closer), or a crossing
				// close whose key differs from the nearest unclosed opener: malformed.
				malformed = true
				break
			}
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			// First completed pair is the innermost; retain it as the candidate loop
			// but keep scanning so a later stray closer, unclosed opener, or a second
			// top-level pair is still classified malformed.
			if startIndex == -1 {
				startIndex, endIndex, loopKey = top.index, i, closeKey
			}
			if len(stack) == 0 {
				topLevelPairs++
			}
		}
	}

	// Full-stream validation: unclosed opener(s) still on the stack, or more than
	// one top-level pair, are malformed even though a prefix balanced.
	if !malformed {
		if len(stack) != 0 {
			malformed = true
		} else if topLevelPairs > 1 {
			malformed = true
		}
	}

	// No well-formed loop pair (no markers, an unclosed opener, a stray/crossing
	// closer, or a second top-level pair): process every row once for simple
	// placeholders. Marker paragraphs are blanked by processParagraph's exact-match
	// branch; any residual template placeholder that cannot resolve at this scope is
	// removed by the residual-token scrub after processing; static content survives.
	// Fail closed — never expand or truncate on ambiguous markers.
	if malformed || startIndex == -1 || endIndex == -1 || loopKey == "" {
		for _, row := range rows {
			processRowCells(row, data)
		}
		return
	}

	// Process non-loop rows (header, total, any enclosing markers, etc.) for
	// placeholder replacement. Rows inside [startIndex..endIndex] are this loop's
	// own markers and template rows, handled below.
	for i, row := range rows {
		if i >= startIndex && i <= endIndex {
			continue
		}
		processRowCells(row, data)
	}

	// Collect template rows (between start and end markers, exclusive)
	var templateRows []*etree.Element
	for i := startIndex + 1; i < endIndex; i++ {
		templateRows = append(templateRows, rows[i])
	}

	if len(templateRows) == 0 {
		// No template rows between markers — just remove the marker rows.
		spliceTableLoop(tbl, rows, startIndex, endIndex, nil, nil)
		return
	}

	// Get the array data for looping
	rawLoopData, ok := getPathValue(data, loopKey)
	if !ok {
		// No data — remove all loop rows (markers + templates).
		spliceTableLoop(tbl, rows, startIndex, endIndex, nil, nil)
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
			// Resolved value is neither []any nor []map[string]any (e.g. a
			// scalar, a plain map, or a typed slice such as []string). Fail
			// closed by removing the marker and template rows, identical to the
			// missing-key and empty-slice cleanup paths above — so no raw loop
			// markers or unresolved template placeholders leak into the output.
			spliceTableLoop(tbl, rows, startIndex, endIndex, nil, nil)
			return
		}
	}

	if len(loopData) == 0 {
		spliceTableLoop(tbl, rows, startIndex, endIndex, nil, nil)
		return
	}

	// Capture the child token immediately following the end-marker row, so the
	// cloned rows land in the exact slot the loop block occupied — preserving the
	// surrounding whitespace CharData and byte-for-byte output identical to the
	// previous insert-before-end-marker strategy.
	var afterAnchor etree.Token
	if idx := rows[endIndex].Index(); idx >= 0 && idx+1 < len(tbl.Child) {
		afterAnchor = tbl.Child[idx+1]
	}

	// Build every cloned row up front, without mutating the tree yet.
	var clones []*etree.Element
	for _, itemData := range loopData {
		itemMap, ok := itemData.(map[string]any)
		if !ok {
			continue
		}
		for _, tmplRow := range templateRows {
			newRow := tmplRow.Copy()
			processRowCells(newRow, itemMap)
			clones = append(clones, newRow)
		}
	}

	// Remove the marker + template rows and splice the clones into their slot in a
	// single linear rebuild (see spliceTableLoop).
	spliceTableLoop(tbl, rows, startIndex, endIndex, clones, afterAnchor)
}

// spliceTableLoop rebuilds tbl's direct-child list in a single O(n) pass: it drops
// the marker + template rows rows[startIndex..endIndex] and inserts `clones` into
// the exact slot the block occupied — immediately before `afterAnchor` (the child
// that followed the end-marker row), or at the end if afterAnchor is nil.
//
// etree's RemoveChildAt and InsertChild reindex EVERY following sibling on each
// call, so removing or expanding a large loop block one node at a time is O(n²)
// (the interspersed whitespace CharData that must be preserved for byte-stability
// are re-indexed on every removal). Detaching all children from the END is O(1)
// per node (nothing follows the last slot to reindex), and re-adding the final
// ordered list via AddChild is O(1) per node, so the whole rebuild is linear —
// an accepted 64 MiB document.xml with many loop rows no longer burns quadratic
// CPU during cleanup.
//
// Non-removed tokens keep their original relative order (including the block's
// interspersed whitespace CharData), and clones are spliced exactly where the old
// insert-before-end-marker + per-row-remove strategy left them, so the serialized
// output is byte-identical — pinned by TestByteStability_FrozenEmissions.
func spliceTableLoop(tbl *etree.Element, rows []*etree.Element, startIndex, endIndex int, clones []*etree.Element, afterAnchor etree.Token) {
	remove := make(map[etree.Token]bool, endIndex-startIndex+1)
	for i := startIndex; i <= endIndex; i++ {
		remove[rows[i]] = true
	}

	final := make([]etree.Token, 0, len(tbl.Child)+len(clones))
	spliced := false
	for _, tok := range tbl.Child {
		if afterAnchor != nil && tok == afterAnchor {
			for _, c := range clones {
				final = append(final, c)
			}
			spliced = true
		}
		if !remove[tok] {
			final = append(final, tok)
		}
	}
	// afterAnchor nil (block ended at the last child) or not found: clones go last.
	if !spliced {
		for _, c := range clones {
			final = append(final, c)
		}
	}

	// Detach every current child from the end (O(1) each), then re-add the final
	// ordered list (O(1) each). AddChild re-parents and re-indexes each token.
	for len(tbl.Child) > 0 {
		tbl.RemoveChildAt(len(tbl.Child) - 1)
	}
	for _, tok := range final {
		tbl.AddChild(tok)
	}
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
			// The close marker must carry the SAME key as the open marker. A
			// mismatched {{/otherkey}} does not close this loop; scanning
			// continues for the matching {{/key}}. If no matching close is ever
			// found we fall through to the final return and report no loop
			// (-1, -1, ""), so ProcessBody treats the body as having no loop:
			// every element is processed once at root scope and stray marker
			// paragraphs are blanked by processParagraph's exact-match branch.
			// This is fail closed (matching the no-marker path), never
			// truncating the template block at the wrong close marker.
			if m := loopEndRegex.FindStringSubmatch(trimmed); len(m) > 0 && m[0] == trimmed && strings.TrimSpace(m[1]) == key {
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
