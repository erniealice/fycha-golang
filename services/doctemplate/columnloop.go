package doctemplate

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/beevik/etree"
)

// maxColumnLoopItems bounds the TOTAL number of grid columns a table may have
// after a column loop expands: the static columns plus the expanded items.
// Word itself caps a table at 63 columns, so anything wider cannot be a
// legitimate document (rule 1).
const maxColumnLoopItems = 63

// ColumnLoopContractError is the structural marker for table column-loop
// failures (count mismatch, over-limit, unsupported cell shape) that must
// survive ProcessTemplate wrapping.
type ColumnLoopContractError interface {
	error
	IsColumnLoopContractError() bool
}

type columnLoopContractError struct{ reason string }

func (e *columnLoopContractError) Error() string {
	return "table column loop contract violation: " + e.reason
}

func (e *columnLoopContractError) IsColumnLoopContractError() bool { return true }

// ColumnLoopContractError is a small structural marker that callers (e.g.
// fayna) can detect via errors.As against a locally declared
// interface{ ColumnLoopContractError() bool }, without importing this
// package's types directly — the same pattern as
// styleTokenContractError.StyleContractError in xmlprocessor.go. Every
// column-loop rejection is constructed through columnLoopError below, so
// every one of them carries this marker.
func (e *columnLoopContractError) ColumnLoopContractError() bool {
	return true
}

func columnLoopError(format string, args ...any) error {
	return &columnLoopContractError{reason: fmt.Sprintf(format, args...)}
}

// columnLoopState is shared by every row of one table. All column loops in a
// table expand the same grid column and must produce the same number of cells,
// otherwise the rows would no longer line up with the table grid.
type columnLoopState struct {
	gridIndex        int // grid column of the template cell; -1 until known
	count            int // cells produced per loop; -1 until the first expansion
	totalGridColumns int // total columns in the table's original grid; -1 if unknown
}

func newColumnLoopState() *columnLoopState {
	return &columnLoopState{gridIndex: -1, count: -1, totalGridColumns: -1}
}

// loopToken is one {{#key}} or {{/key}} marker found in a cell's concatenated
// text, in document order.
type loopToken struct {
	start, end int
	key        string
	isOpen     bool
}

// scanLoopTokens finds every loop-start and loop-end marker in text, ordered
// by position, so a cell's marker stream can be validated as a whole (rule 5)
// instead of matching only the first opener and the last closer.
func scanLoopTokens(text string) []loopToken {
	var tokens []loopToken
	for _, m := range loopStartRegex.FindAllStringSubmatchIndex(text, -1) {
		tokens = append(tokens, loopToken{start: m[0], end: m[1], key: strings.TrimSpace(text[m[2]:m[3]]), isOpen: true})
	}
	for _, m := range loopEndRegex.FindAllStringSubmatchIndex(text, -1) {
		tokens = append(tokens, loopToken{start: m[0], end: m[1], key: strings.TrimSpace(text[m[2]:m[3]])})
	}
	sort.Slice(tokens, func(i, j int) bool { return tokens[i].start < tokens[j].start })
	return tokens
}

// classifyCell scans tc's whole text and reports whether it is a column-loop
// cell (its entire text is one {{#key}} ... {{/key}} pair), an ordinary
// row-loop marker cell (its entire text is a lone {{#key}} or {{/key}}), or
// neither. Any other marker shape in the cell — an opener with no closer
// (split), more than one pair (repeated), or a nested pair — is rejected with
// a typed error rather than silently misread, per rule 5. This runs BEFORE
// any row-loop scan sees the cell (rule 7), so a column-loop marker is never
// misread as a row-loop marker and a malformed cell is never mutated first.
func classifyCell(tc *etree.Element) (key string, isColumnLoop bool, err error) {
	trimmed := strings.TrimSpace(elementText(tc))
	tokens := scanLoopTokens(trimmed)
	if len(tokens) == 0 {
		return "", false, nil
	}

	first := tokens[0]
	if !first.isOpen || first.start != 0 {
		// Not a column-loop candidate at all: the cell's first marker is a
		// closer, or an opener is preceded by other text (e.g. a
		// mixed-content row-loop marker cell like "Note {{#items}} tail", or
		// a lone "{{/items}}" closer cell). Leave these exactly as before —
		// classification only engages once a cell visibly opens with
		// {{#key}} at position 0; anything else is the row-loop scanner's
		// concern, not a column-loop contract violation.
		return "", false, nil
	}

	if len(tokens) == 1 {
		if first.end == len(trimmed) {
			// A lone opener spanning the WHOLE cell is a row-loop marker
			// cell, not a column-loop attempt.
			return "", false, nil
		}
		return "", false, columnLoopError("cell %q has an incomplete loop marker: an opener with no matching closer in the same cell, or stray text around it", trimmed)
	}

	if len(tokens) > 2 {
		return "", false, columnLoopError("cell %q has nested or repeated loop markers", trimmed)
	}

	open, close := tokens[0], tokens[1]
	if close.isOpen {
		return "", false, columnLoopError("cell %q has loop markers out of order", trimmed)
	}
	if close.end != len(trimmed) {
		return "", false, columnLoopError("cell %q loop marker pair does not span the whole cell", trimmed)
	}
	if open.key != close.key {
		return "", false, columnLoopError("cell %q opens %q but closes %q", trimmed, open.key, close.key)
	}
	return open.key, true, nil
}

// columnLoopCell is the error-swallowing form of classifyCell for callers that
// run strictly after every cell in the table has already been classified once
// without error (see rowHasColumnLoop / the table-level pre-pass).
func columnLoopCell(tc *etree.Element) (string, bool) {
	key, ok, _ := classifyCell(tc)
	return key, ok
}

func elementText(el *etree.Element) string {
	var sb strings.Builder
	for _, t := range el.FindElements(".//t") {
		sb.WriteString(t.Text())
	}
	return sb.String()
}

// rowHasColumnLoop reports whether any direct cell of row is a column loop,
// classifying every cell (rule 7) and surfacing the first marker-stream
// contract violation found (rule 5) rather than treating it as an ordinary
// cell.
func rowHasColumnLoop(row *etree.Element) (bool, error) {
	found := false
	for _, cell := range directChildren(row, "tc") {
		_, isLoop, err := classifyCell(cell)
		if err != nil {
			return false, err
		}
		if isLoop {
			found = true
		}
	}
	return found, nil
}

func directChildren(el *etree.Element, tag string) []*etree.Element {
	var out []*etree.Element
	for _, child := range el.ChildElements() {
		if child.Tag == tag {
			out = append(out, child)
		}
	}
	return out
}

func childByTag(el *etree.Element, tag string) *etree.Element {
	if el == nil {
		return nil
	}
	for _, child := range el.ChildElements() {
		if child.Tag == tag {
			return child
		}
	}
	return nil
}

// wordAttr finds el's attribute named key IN THE WORDPROCESSINGML NAMESPACE
// (rule 6) — not merely an attribute that happens to share the local name,
// which a foreign-namespaced decoy attribute could otherwise supply.
func wordAttr(el *etree.Element, key string, namespaces namespaceContext) *etree.Attr {
	if el == nil {
		return nil
	}
	for i := range el.Attr {
		attr := &el.Attr[i]
		if attr.Key != key {
			continue
		}
		if attr.Space == "" {
			continue
		}
		if namespaces[attr.Space] == wordprocessingMLNamespaceURI {
			return attr
		}
	}
	return nil
}

// intAttr reads a w:-namespaced integer attribute. A missing attribute
// returns fallback (an ordinary, documented default); a PRESENT but
// non-numeric value is a typed error, never a silent 0 (rule 6).
func intAttr(el *etree.Element, key string, namespaces namespaceContext, fallback int) (int, error) {
	attr := wordAttr(el, key, namespaces)
	if attr == nil {
		return fallback, nil
	}
	value, err := strconv.Atoi(strings.TrimSpace(attr.Value))
	if err != nil {
		return 0, columnLoopError("attribute %s has a malformed integer value %q", attr.FullKey(), attr.Value)
	}
	return value, nil
}

// requiredIntAttr is intAttr for attributes that MUST be present (a missing
// value is itself the contract violation, e.g. a loop cell's width).
func requiredIntAttr(el *etree.Element, key string, namespaces namespaceContext, context string) (int, error) {
	attr := wordAttr(el, key, namespaces)
	if attr == nil {
		return 0, columnLoopError("%s has no %s attribute", context, key)
	}
	value, err := strconv.Atoi(strings.TrimSpace(attr.Value))
	if err != nil {
		return 0, columnLoopError("%s attribute %s=%q is not a valid integer", context, key, attr.Value)
	}
	return value, nil
}

// cellGridSpan returns w:tcPr/w:gridSpan/@w:val (1 when absent).
func cellGridSpan(tc *etree.Element, namespaces namespaceContext) (int, error) {
	return intAttr(childByTag(childByTag(tc, "tcPr"), "gridSpan"), "val", namespaces, 1)
}

// cellGridIndex is the grid column where tc starts: w:trPr/w:gridBefore plus
// the spans of the cells before it.
func cellGridIndex(row, tc *etree.Element, namespaces namespaceContext) (int, error) {
	index, err := intAttr(childByTag(childByTag(row, "trPr"), "gridBefore"), "val", namespaces, 0)
	if err != nil {
		return -1, err
	}
	for _, cell := range directChildren(row, "tc") {
		if cell == tc {
			return index, nil
		}
		span, err := cellGridSpan(cell, namespaces)
		if err != nil {
			return -1, err
		}
		index += span
	}
	return -1, nil
}

// columnLoopCellWidthDXA returns a column-loop cell's tcW width in dxa units.
// Only dxa can be split and reassembled without changing the table's overall
// width, so pct, auto, and a missing tcW/width value are rejected until
// supported (rule 2).
func columnLoopCellWidthDXA(cell *etree.Element, namespaces namespaceContext) (int, error) {
	tcW := childByTag(childByTag(cell, "tcPr"), "tcW")
	if tcW == nil {
		return 0, columnLoopError("column loop cell has no w:tcW width; only dxa is supported")
	}
	if typeAttr := wordAttr(tcW, "type", namespaces); typeAttr != nil && typeAttr.Value != "dxa" {
		return 0, columnLoopError("column loop cell width type %q is not supported; only dxa", typeAttr.Value)
	}
	return requiredIntAttr(tcW, "w", namespaces, "column loop cell tcW")
}

// gridColWidthDXA returns a w:gridCol's width; gridCol widths are always dxa.
func gridColWidthDXA(gridCol *etree.Element, namespaces namespaceContext) (int, error) {
	return requiredIntAttr(gridCol, "w", namespaces, "column loop grid column")
}

// shares splits total into n integer parts; the last part takes the remainder.
func shares(total, n int) []int {
	out := make([]int, n)
	for i := range out {
		out[i] = total / n
	}
	if n > 0 {
		out[n-1] += total - (total/n)*n
	}
	return out
}

// splitSharesNonZero is shares with a rule-2 guard: a template cell too
// narrow to give every clone a positive width must be rejected, not silently
// produce zero- or stale-width columns.
func splitSharesNonZero(total, n int) ([]int, bool) {
	widths := shares(total, n)
	for _, w := range widths {
		if w <= 0 {
			return nil, false
		}
	}
	return widths, true
}

// planColumnLoopStructure validates the whole table's column-loop shape
// BEFORE any row is mutated: every column-loop cell spans exactly one grid
// column and shares the same grid column (rule 1's baseline, unchanged),
// contains no nested table (rule 4), and no OTHER row's cell crosses that
// grid column with a gridSpan/gridBefore/gridAfter/vMerge this engine cannot
// remap (rule 3). rowHasLoop is the rule-7 pre-pass classification: it is
// computed once, before the row-loop marker scan, and reused here.
func planColumnLoopStructure(tbl *etree.Element, rows []*etree.Element, rowHasLoop []bool, namespaces namespaceContext, state *columnLoopState) error {
	anyLoop := false
	for _, has := range rowHasLoop {
		if has {
			anyLoop = true
			break
		}
	}
	if !anyLoop {
		return nil
	}

	for i, row := range rows {
		if !rowHasLoop[i] {
			continue
		}
		for _, cell := range directChildren(row, "tc") {
			_, isLoop, err := classifyCell(cell)
			if err != nil {
				return err
			}
			if !isLoop {
				continue
			}
			if cell.FindElement(".//tbl") != nil {
				return columnLoopError("column loop cell contains a nested table, which is not supported")
			}
			span, err := cellGridSpan(cell, namespaces)
			if err != nil {
				return err
			}
			if span != 1 {
				return columnLoopError("column loop cell must span exactly one grid column")
			}
			index, err := cellGridIndex(row, cell, namespaces)
			if err != nil {
				return err
			}
			if state.gridIndex == -1 {
				state.gridIndex = index
			} else if state.gridIndex != index {
				return columnLoopError("column loops in one table must expand the same grid column")
			}
		}
	}

	if state.gridIndex == -1 {
		return nil
	}

	grid := childByTag(tbl, "tblGrid")
	if grid == nil {
		return columnLoopError("table with a column loop has no grid")
	}
	cols := directChildren(grid, "gridCol")
	state.totalGridColumns = len(cols)
	if state.gridIndex >= len(cols) {
		return columnLoopError("column loop grid column %d is outside the table grid", state.gridIndex+1)
	}

	for i, row := range rows {
		if rowHasLoop[i] {
			continue
		}
		if err := validateRowDoesNotCrossColumn(row, state.gridIndex, namespaces); err != nil {
			return err
		}
	}
	return nil
}

// validateRowDoesNotCrossColumn rejects a non-column-loop row whose gridSpan,
// gridBefore, gridAfter, or vMerge would place a cell across the loop's grid
// column (rule 3). A plain, single-column cell that merely sits AT that
// column (span 1, no vMerge) is normal and left alone.
func validateRowDoesNotCrossColumn(row *etree.Element, gridIndex int, namespaces namespaceContext) error {
	trPr := childByTag(row, "trPr")
	gridBefore, err := intAttr(childByTag(trPr, "gridBefore"), "val", namespaces, 0)
	if err != nil {
		return err
	}
	gridAfter, err := intAttr(childByTag(trPr, "gridAfter"), "val", namespaces, 0)
	if err != nil {
		return err
	}
	if gridBefore > 0 && gridIndex < gridBefore {
		return columnLoopError("row gridBefore %d crosses the column-loop grid column %d", gridBefore, gridIndex+1)
	}

	index := gridBefore
	for _, cell := range directChildren(row, "tc") {
		span, err := cellGridSpan(cell, namespaces)
		if err != nil {
			return err
		}
		start, end := index, index+span
		if gridIndex >= start && gridIndex < end {
			if span > 1 {
				return columnLoopError("row cell spans the column-loop grid column %d (gridSpan %d)", gridIndex+1, span)
			}
			if vMerge := childByTag(childByTag(cell, "tcPr"), "vMerge"); vMerge != nil {
				return columnLoopError("row cell has vMerge across the column-loop grid column %d", gridIndex+1)
			}
		}
		index = end
	}

	if gridAfter > 0 && gridIndex >= index {
		return columnLoopError("row gridAfter %d crosses the column-loop grid column %d", gridAfter, gridIndex+1)
	}
	return nil
}

// expandColumnLoops replaces every column-loop cell of row with one resolved
// clone per item and returns the clones (which callers must not re-process).
// Structural validation (gridSpan, shared gridIndex, nested tables, crossing
// rows) has already run in planColumnLoopStructure; this only resolves data
// and widths, which are unknown until now.
func expandColumnLoops(row *etree.Element, data map[string]any, inherited namespaceContext, state *columnLoopState) (map[*etree.Element]bool, error) {
	namespaces := mergeNamespaceContexts(inherited, namespaceContextFor(row))
	clones := map[*etree.Element]bool{}
	for _, cell := range directChildren(row, "tc") {
		key, ok, err := classifyCell(cell)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}

		if state.gridIndex == -1 {
			// Defensive fallback for callers that process a row without first
			// running planColumnLoopStructure (e.g. a standalone row with no
			// enclosing table). The normal table path already has gridIndex
			// fixed before any row is expanded.
			span, err := cellGridSpan(cell, namespaces)
			if err != nil {
				return nil, err
			}
			if span != 1 {
				return nil, columnLoopError("column loop %q cell must span exactly one grid column", key)
			}
			index, err := cellGridIndex(row, cell, namespaces)
			if err != nil {
				return nil, err
			}
			state.gridIndex = index
		}

		var items []any
		if raw, present := getPathValue(data, key); present && raw != nil {
			items = resolveLoopData(data, key)
			if items == nil {
				return nil, columnLoopError("column loop %q is not a list", key)
			}
		}

		if state.count == -1 {
			state.count = len(items)
		} else if state.count != len(items) {
			return nil, columnLoopError("column loop %q has %d items; other rows of the table have %d", key, len(items), state.count)
		}

		total := len(items)
		if state.totalGridColumns > 0 {
			total += state.totalGridColumns - 1
		}
		if total > maxColumnLoopItems {
			return nil, columnLoopError("column loop %q would produce %d table columns; the limit is %d", key, total, maxColumnLoopItems)
		}

		template := cell.Copy()
		if err := stripColumnLoopMarkers(template, key); err != nil {
			return nil, err
		}

		if len(items) == 0 {
			// Rule 8: zero items removes the loop column entirely for this
			// row; no width to split, nothing to clone.
			row.RemoveChild(cell)
			continue
		}

		width, err := columnLoopCellWidthDXA(cell, namespaces)
		if err != nil {
			return nil, err
		}
		widths, ok := splitSharesNonZero(width, len(items))
		if !ok {
			return nil, columnLoopError("column loop %q cell width %d cannot be split into %d nonzero-width columns", key, width, len(items))
		}

		for i, item := range items {
			itemMap, isMap := item.(map[string]any)
			if !isMap {
				return nil, columnLoopError("column loop %q item %d is not an object", key, i+1)
			}
			clone := template.Copy()
			w := childByTag(childByTag(clone, "tcPr"), "tcW")
			if w == nil {
				return nil, columnLoopError("column loop %q cell lost its tcW while cloning", key)
			}
			attr := wordAttr(w, "w", namespaces)
			if attr == nil {
				return nil, columnLoopError("column loop %q cell lost its tcW width value while cloning", key)
			}
			attr.Value = strconv.Itoa(widths[i])
			// Insert first so the clone sees the row's namespace declarations.
			row.InsertChild(cell, clone)
			for _, p := range clone.FindElements(".//p") {
				processParagraph(p, itemMap)
			}
			if err := applyStyleAttributesWithNamespaceContext(clone, itemMap, inherited); err != nil {
				return nil, err
			}
			clones[clone] = true
		}
		row.RemoveChild(cell)
	}
	return clones, nil
}

// stripColumnLoopMarkers removes the {{#key}} and {{/key}} markers from the
// template cell. Each marker must sit in a single text node (a marker split
// across runs is rejected rather than half-removed).
func stripColumnLoopMarkers(cell *etree.Element, key string) error {
	for _, t := range cell.FindElements(".//t") {
		text := t.Text()
		text = loopStartRegex.ReplaceAllStringFunc(text, func(m string) string {
			if extractLoopMarker(m, "#") == key {
				return ""
			}
			return m
		})
		text = loopEndRegex.ReplaceAllStringFunc(text, func(m string) string {
			if extractLoopMarker(m, "/") == key {
				return ""
			}
			return m
		})
		t.SetText(text)
	}
	remaining := elementText(cell)
	if strings.Contains(remaining, "{{#") || strings.Contains(remaining, "{{/") {
		return columnLoopError("column loop %q marker is split across runs or nested", key)
	}
	return nil
}

// applyColumnLoopGrid rewrites w:tblGrid after the table's column loops have
// expanded: the template grid column becomes `count` columns sharing its
// width, so the table keeps its overall width. Zero items removes the loop
// grid column entirely, leaving the grid consistent with every row's now
// loop-cell-free cell count (rule 8).
func applyColumnLoopGrid(tbl *etree.Element, state *columnLoopState, namespaces namespaceContext) error {
	if state == nil || state.gridIndex < 0 || state.count < 0 {
		return nil
	}
	grid := childByTag(tbl, "tblGrid")
	if grid == nil {
		return columnLoopError("table with a column loop has no grid")
	}
	cols := directChildren(grid, "gridCol")
	if state.gridIndex >= len(cols) {
		return columnLoopError("column loop grid column %d is outside the table grid", state.gridIndex+1)
	}
	template := cols[state.gridIndex]

	if state.count == 0 {
		grid.RemoveChild(template)
		return nil
	}

	width, err := gridColWidthDXA(template, namespaces)
	if err != nil {
		return err
	}
	widths, ok := splitSharesNonZero(width, state.count)
	if !ok {
		return columnLoopError("column loop grid column width %d cannot be split into %d nonzero-width columns", width, state.count)
	}
	for i := 0; i < state.count; i++ {
		clone := template.Copy()
		attr := wordAttr(clone, "w", namespaces)
		if attr == nil {
			return columnLoopError("column loop grid column lost its width while cloning")
		}
		attr.Value = strconv.Itoa(widths[i])
		grid.InsertChild(template, clone)
	}
	grid.RemoveChild(template)
	return nil
}

// processColumnLoopRow expands the row's column loops, then processes the
// remaining (non-clone) cells and row properties against data, as a row
// without column loops is processed.
func processColumnLoopRow(row *etree.Element, data map[string]any, inherited namespaceContext, state *columnLoopState) error {
	clones, err := expandColumnLoops(row, data, inherited, state)
	if err != nil {
		return err
	}
	for _, cell := range row.FindElements(".//tc") {
		if insideAny(cell, row, clones) {
			continue
		}
		for _, p := range cell.FindElements("./p") {
			processParagraph(p, data)
		}
	}
	for _, child := range row.ChildElements() {
		if clones[child] {
			continue
		}
		if err := applyStyleAttributesWithNamespaceContext(child, data, inherited); err != nil {
			return err
		}
	}
	return nil
}

func insideAny(el, stop *etree.Element, set map[*etree.Element]bool) bool {
	for current := el; current != nil && current != stop; current = current.Parent() {
		if set[current] {
			return true
		}
	}
	return false
}
