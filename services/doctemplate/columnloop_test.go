package doctemplate

import (
	"errors"
	"strings"
	"testing"
)

// A matrix table: a static label column plus one column-loop column. The
// header loop iterates root "columns"; the body loop iterates each row's
// "cells" inside the "rows" row loop.
const columnLoopTable = `<w:tbl>
<w:tblGrid><w:gridCol w:w="3000"/><w:gridCol w:w="9000"/></w:tblGrid>
<w:tr><w:tc><w:tcPr><w:tcW w:w="3000" w:type="dxa"/></w:tcPr><w:p><w:r><w:t>{{label_heading}}</w:t></w:r></w:p></w:tc><w:tc><w:tcPr><w:tcW w:w="9000" w:type="dxa"/></w:tcPr><w:p><w:r><w:t>{{#columns}}</w:t></w:r><w:r><w:t>{{name}}</w:t></w:r><w:r><w:t>{{/columns}}</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>{{#rows}}</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:tcPr><w:tcW w:w="3000" w:type="dxa"/></w:tcPr><w:p><w:r><w:rPr><w:b w:val="{{bold}}"/></w:rPr><w:t>{{label}}</w:t></w:r></w:p></w:tc><w:tc><w:tcPr><w:tcW w:w="9000" w:type="dxa"/><w:shd w:val="clear" w:color="auto" w:fill="{{fill}}"/></w:tcPr><w:p><w:r><w:t>{{#cells}}</w:t></w:r><w:r><w:t>{{value}}</w:t></w:r><w:r><w:t>{{/cells}}</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>{{/rows}}</w:t></w:r></w:p></w:tc></w:tr>
</w:tbl>`

func matrixData(columns int, cellsPerRow ...int) map[string]any {
	cols := make([]any, columns)
	for i := range cols {
		cols[i] = map[string]any{"name": "Col" + string(rune('A'+i))}
	}
	rows := make([]any, len(cellsPerRow))
	for r, n := range cellsPerRow {
		cells := make([]any, n)
		for c := range cells {
			cells[c] = map[string]any{"value": "v" + string(rune('A'+r)) + string(rune('0'+c)), "fill": "00FF00"}
		}
		rows[r] = map[string]any{"label": "Row" + string(rune('A'+r)), "bold": "false", "cells": cells}
	}
	return map[string]any{"label_heading": "Name", "columns": cols, "rows": rows}
}

func TestColumnLoop_ExpandsHeaderAndRowsAndSplitsGrid(t *testing.T) {
	content := renderBody(t, columnLoopTable, matrixData(3, 3, 3))

	mustContain(t, content, "ColA", "ColB", "ColC", "vA0", "vA2", "vB1", "RowA", "RowB", `w:fill="00FF00"`)
	mustNotContain(t, content, "{{", "}}")
	if got := strings.Count(content, `<w:gridCol w:w="3000"/>`); got != 4 {
		t.Fatalf("grid columns of 3000 = %d, want 1 label + 3 split (9000/3)", got)
	}
	if got := strings.Count(content, `<w:tcW w:w="3000" w:type="dxa"/>`); got != 3+3*3 {
		t.Fatalf("cells of width 3000 = %d, want 3 label cells + 9 loop cells", got)
	}
}

func TestColumnLoop_SingleColumnKeepsFullWidth(t *testing.T) {
	content := renderBody(t, columnLoopTable, matrixData(1, 1))
	mustContain(t, content, "ColA", "vA0", `<w:gridCol w:w="9000"/>`)
	mustNotContain(t, content, "{{")
}

func TestColumnLoop_RemainderGoesToLastColumn(t *testing.T) {
	content := renderBody(t, columnLoopTable, matrixData(4, 4))
	// 9000 / 4 = 2250 each, no remainder; use 7 columns: 1285 x6 + 1290.
	if !strings.Contains(content, `<w:gridCol w:w="2250"/>`) {
		t.Fatalf("expected 2250-wide grid columns")
	}
	content = renderBody(t, columnLoopTable, matrixData(7, 7))
	mustContain(t, content, `<w:gridCol w:w="1285"/>`, `<w:gridCol w:w="1290"/>`)
}

func TestColumnLoop_ProcessTemplateFailsClosed(t *testing.T) {
	cases := map[string]map[string]any{
		"row count differs from header": matrixData(3, 3, 2),
		"too many columns":              matrixData(maxColumnLoopItems+1, maxColumnLoopItems+1),
		"loop value is not a list": func() map[string]any {
			d := matrixData(1, 1)
			d["columns"] = "not a list"
			return d
		}(),
	}
	for name, data := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := ProcessTemplate(createTestDocx(t, bodyDoc(columnLoopTable)), data)
			var contract ColumnLoopContractError
			if err == nil || !errors.As(err, &contract) || !contract.IsColumnLoopContractError() {
				t.Fatalf("err = %v, want a ColumnLoopContractError", err)
			}
		})
	}
}

func TestColumnLoop_RejectsSpannedTemplateCell(t *testing.T) {
	inner := `<w:tbl><w:tblGrid><w:gridCol w:w="1000"/><w:gridCol w:w="1000"/></w:tblGrid>
<w:tr><w:tc><w:tcPr><w:gridSpan w:val="2"/></w:tcPr><w:p><w:r><w:t>{{#columns}}</w:t></w:r><w:r><w:t>{{name}}</w:t></w:r><w:r><w:t>{{/columns}}</w:t></w:r></w:p></w:tc></w:tr>
</w:tbl>`
	_, err := ProcessTemplate(createTestDocx(t, bodyDoc(inner)), matrixData(2))
	var contract ColumnLoopContractError
	if !errors.As(err, &contract) {
		t.Fatalf("err = %v, want a ColumnLoopContractError", err)
	}
}

// A table without a column-loop cell must not be affected; a row-loop marker
// row (opener only) is never mistaken for a column loop.
func TestColumnLoop_RowMarkerCellIsNotAColumnLoop(t *testing.T) {
	inner := `<w:tbl><w:tblGrid><w:gridCol w:w="5000"/></w:tblGrid>
<w:tr><w:tc><w:p><w:r><w:t>{{#rows}}</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>{{label}}</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>{{/rows}}</w:t></w:r></w:p></w:tc></w:tr>
</w:tbl>`
	content := renderBody(t, inner, matrixData(0, 0, 0))
	mustContain(t, content, "RowA", "RowB", `<w:gridCol w:w="5000"/>`)
}

// ----------------------------------------------------------------------------
// Hardening rules 1-8 (fycha-golang.md, D12) — one focused test per rule.
// ----------------------------------------------------------------------------

// Rule 1: the limit applies to the TOTAL grid columns (static + expanded
// items), not the item count alone. columnLoopTable has one static (label)
// column plus the loop column, so 63 items — which the old items-only check
// would have allowed — produces 64 total columns and must be rejected; 62
// items (63 total) is the ceiling and must succeed.
func TestColumnLoop_TotalGridColumnLimit(t *testing.T) {
	_, err := ProcessTemplate(createTestDocx(t, bodyDoc(columnLoopTable)), matrixData(63, 63))
	var contract ColumnLoopContractError
	if err == nil || !errors.As(err, &contract) {
		t.Fatalf("err = %v, want a ColumnLoopContractError for 64 total columns", err)
	}

	content := renderBody(t, columnLoopTable, matrixData(62, 62))
	mustNotContain(t, content, "{{", "}}")
}

// Rule 2: width splitting must handle w:tcW units explicitly. dxa is
// supported; pct, auto, and a missing width are rejected with a typed error
// rather than producing a zero or stale width.
func TestColumnLoop_RejectsPctWidth(t *testing.T) {
	build := func(tcPr string) string {
		return `<w:tbl><w:tblGrid><w:gridCol w:w="4000"/></w:tblGrid>
<w:tr><w:tc>` + tcPr + `<w:p><w:r><w:t>{{#columns}}</w:t></w:r><w:r><w:t>{{name}}</w:t></w:r><w:r><w:t>{{/columns}}</w:t></w:r></w:p></w:tc></w:tr>
</w:tbl>`
	}
	cases := map[string]string{
		"pct":           `<w:tcPr><w:tcW w:w="4000" w:type="pct"/></w:tcPr>`,
		"auto":          `<w:tcPr><w:tcW w:w="4000" w:type="auto"/></w:tcPr>`,
		"missing value": `<w:tcPr><w:tcW w:type="dxa"/></w:tcPr>`,
		"no tcW":        ``,
	}
	for name, tcPr := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := ProcessTemplate(createTestDocx(t, bodyDoc(build(tcPr))), matrixData(2))
			var contract ColumnLoopContractError
			if err == nil || !errors.As(err, &contract) {
				t.Fatalf("err = %v, want a ColumnLoopContractError for width case %q", err, name)
			}
		})
	}
}

// Rule 3: a gridSpan, gridBefore, gridAfter, or vMerge in another (non-loop)
// row that crosses the loop's grid column is rejected — the engine cannot
// remap those to the expanded columns.
func TestColumnLoop_RejectsGridSpanCrossingLoop(t *testing.T) {
	cases := map[string]string{
		// The loop column is grid index 1 (label is index 0). A gridSpan of 2
		// starting at index 0 crosses it.
		"gridSpan": `<w:tr><w:tc><w:tcPr><w:gridSpan w:val="2"/></w:tcPr><w:p><w:r><w:t>Total</w:t></w:r></w:p></w:tc></w:tr>`,
		// A plain index-0 cell, then a vMerge cell landing exactly on the
		// loop column (index 1).
		"vMerge": `<w:tr><w:tc><w:p><w:r><w:t>Total</w:t></w:r></w:p></w:tc><w:tc><w:tcPr><w:vMerge/></w:tcPr><w:p><w:r><w:t>-</w:t></w:r></w:p></w:tc></w:tr>`,
		// gridBefore=2 skips columns [0,2) for this row, which swallows the
		// loop column (index 1) into its invisible prefix.
		"gridBefore": `<w:tr><w:trPr><w:gridBefore w:val="2"/></w:trPr><w:tc><w:p><w:r><w:t>Total</w:t></w:r></w:p></w:tc></w:tr>`,
	}
	for name, extraRow := range cases {
		t.Run(name, func(t *testing.T) {
			inner := `<w:tbl><w:tblGrid><w:gridCol w:w="3000"/><w:gridCol w:w="9000"/></w:tblGrid>
<w:tr><w:tc><w:tcPr><w:tcW w:w="3000" w:type="dxa"/></w:tcPr><w:p><w:r><w:t>Name</w:t></w:r></w:p></w:tc><w:tc><w:tcPr><w:tcW w:w="9000" w:type="dxa"/></w:tcPr><w:p><w:r><w:t>{{#columns}}</w:t></w:r><w:r><w:t>{{name}}</w:t></w:r><w:r><w:t>{{/columns}}</w:t></w:r></w:p></w:tc></w:tr>
` + extraRow + `
</w:tbl>`
			_, err := ProcessTemplate(createTestDocx(t, bodyDoc(inner)), matrixData(2))
			var contract ColumnLoopContractError
			if err == nil || !errors.As(err, &contract) {
				t.Fatalf("err = %v, want a ColumnLoopContractError for a %s crossing the loop column", err, name)
			}
		})
	}
}

// Rule 4: a nested table inside a column-loop cell is not supported.
func TestColumnLoop_RejectsNestedTable(t *testing.T) {
	inner := `<w:tbl><w:tblGrid><w:gridCol w:w="4000"/></w:tblGrid>
<w:tr><w:tc><w:tcPr><w:tcW w:w="4000" w:type="dxa"/></w:tcPr>
<w:p><w:r><w:t>{{#columns}}</w:t></w:r><w:r><w:t>{{name}}</w:t></w:r><w:r><w:t>{{/columns}}</w:t></w:r></w:p>
<w:tbl><w:tblGrid><w:gridCol w:w="1000"/></w:tblGrid><w:tr><w:tc><w:p/></w:tc></w:tr></w:tbl>
</w:tc></w:tr>
</w:tbl>`
	_, err := ProcessTemplate(createTestDocx(t, bodyDoc(inner)), matrixData(2))
	var contract ColumnLoopContractError
	if err == nil || !errors.As(err, &contract) {
		t.Fatalf("err = %v, want a ColumnLoopContractError for a nested table inside a column-loop cell", err)
	}
}

// Rule 5a: an opener with no matching closer in the same cell (a "split"
// marker) is rejected rather than left to leak or half-strip.
func TestColumnLoop_RejectsSplitMarkers(t *testing.T) {
	inner := `<w:tbl><w:tblGrid><w:gridCol w:w="4000"/></w:tblGrid>
<w:tr><w:tc><w:tcPr><w:tcW w:w="4000" w:type="dxa"/></w:tcPr><w:p><w:r><w:t>{{#columns}}</w:t></w:r><w:r><w:t>{{name}}</w:t></w:r></w:p></w:tc></w:tr>
</w:tbl>`
	_, err := ProcessTemplate(createTestDocx(t, bodyDoc(inner)), matrixData(2))
	var contract ColumnLoopContractError
	if err == nil || !errors.As(err, &contract) {
		t.Fatalf("err = %v, want a ColumnLoopContractError for an opener with no closer in the same cell", err)
	}
}

// Rule 5b: a repeated (or nested) same-key marker pair in one cell is
// rejected; the cell's full marker stream is parsed before anything mutates.
func TestColumnLoop_RejectsRepeatedMarker(t *testing.T) {
	inner := `<w:tbl><w:tblGrid><w:gridCol w:w="4000"/></w:tblGrid>
<w:tr><w:tc><w:tcPr><w:tcW w:w="4000" w:type="dxa"/></w:tcPr><w:p>` +
		`<w:r><w:t>{{#columns}}</w:t></w:r><w:r><w:t>{{name}}</w:t></w:r><w:r><w:t>{{/columns}}</w:t></w:r>` +
		`<w:r><w:t>{{#columns}}</w:t></w:r><w:r><w:t>{{name}}</w:t></w:r><w:r><w:t>{{/columns}}</w:t></w:r>` +
		`</w:p></w:tc></w:tr>
</w:tbl>`
	_, err := ProcessTemplate(createTestDocx(t, bodyDoc(inner)), matrixData(2))
	var contract ColumnLoopContractError
	if err == nil || !errors.As(err, &contract) {
		t.Fatalf("err = %v, want a ColumnLoopContractError for a repeated same-key marker pair in one cell", err)
	}
}

// Rule 6a: wordAttr must match the WordprocessingML namespace, not just the
// local attribute name. A foreign-namespaced decoy w:type/w:w pair ordered
// BEFORE the real w:-namespaced ones must be ignored.
func TestColumnLoop_NamespaceAwareAttr(t *testing.T) {
	inner := `<w:tbl xmlns:x="urn:example:decoy"><w:tblGrid><w:gridCol w:w="4000"/></w:tblGrid>
<w:tr><w:tc><w:tcPr><w:tcW x:type="pct" x:w="999999" w:type="dxa" w:w="4000"/></w:tcPr><w:p><w:r><w:t>{{#columns}}</w:t></w:r><w:r><w:t>{{name}}</w:t></w:r><w:r><w:t>{{/columns}}</w:t></w:r></w:p></w:tc></w:tr>
</w:tbl>`
	content := renderBody(t, inner, matrixData(2))
	mustNotContain(t, content, "{{", "}}")
	mustContain(t, content, "ColA", "ColB")
}

// Rule 6b: intAttr rejects a present-but-malformed integer with a typed
// error, never a silent 0/default.
func TestColumnLoop_MalformedIntAttr(t *testing.T) {
	inner := `<w:tbl><w:tblGrid><w:gridCol w:w="4000"/></w:tblGrid>
<w:tr><w:tc><w:tcPr><w:gridSpan w:val="two"/><w:tcW w:w="4000" w:type="dxa"/></w:tcPr><w:p><w:r><w:t>{{#columns}}</w:t></w:r><w:r><w:t>{{name}}</w:t></w:r><w:r><w:t>{{/columns}}</w:t></w:r></w:p></w:tc></w:tr>
</w:tbl>`
	_, err := ProcessTemplate(createTestDocx(t, bodyDoc(inner)), matrixData(2))
	var contract ColumnLoopContractError
	if err == nil || !errors.As(err, &contract) {
		t.Fatalf("err = %v, want a ColumnLoopContractError for a malformed integer attribute, not a silent 0", err)
	}
}

// Rule 7: column-loop cells are classified structurally BEFORE the row-loop
// marker scan. If the header's column-loop cell text ("{{#columns}}...
// {{/columns}}") leaked into row-loop marker scanning, it would push an
// unclosed "columns" opener onto the row-loop stack and falsely mark the
// WHOLE table malformed — failing the genuine, unrelated "items" row loop
// closed instead of expanding it.
func TestColumnLoop_ClassifiedBeforeRowLoop(t *testing.T) {
	inner := `<w:tbl>
<w:tblGrid><w:gridCol w:w="4000"/><w:gridCol w:w="4000"/></w:tblGrid>
<w:tr><w:tc><w:tcPr><w:tcW w:w="4000" w:type="dxa"/></w:tcPr><w:p><w:r><w:t>{{#columns}}</w:t></w:r><w:r><w:t>{{name}}</w:t></w:r><w:r><w:t>{{/columns}}</w:t></w:r></w:p></w:tc><w:tc><w:p><w:r><w:t>Header</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>{{#items}}</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>{{field}}</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>{{/items}}</w:t></w:r></w:p></w:tc></w:tr>
</w:tbl>`
	content := renderBody(t, inner, map[string]any{
		"columns": []any{map[string]any{"name": "ColA"}, map[string]any{"name": "ColB"}},
		"items": []any{
			map[string]any{"field": "V1"},
			map[string]any{"field": "V2"},
			map[string]any{"field": "V3"},
		},
	})
	mustNotContain(t, content, "{{", "}}")
	mustContain(t, content, "ColA", "ColB", "V1", "V2", "V3")
	if got := strings.Count(content, "V1"); got != 1 {
		t.Fatalf("expected the items row loop to expand exactly once per item, got V1 x%d (content=%s)", got, content)
	}
}

// Rule 8: zero items leaves a valid table — the loop grid column removed and
// every affected row's cell count consistent with the shrunken grid.
func TestColumnLoop_ZeroItemsValidTable(t *testing.T) {
	content := renderBody(t, columnLoopTable, matrixData(0, 0))
	mustNotContain(t, content, "{{", "}}")
	mustContain(t, content, "Name", "RowA")

	if got := strings.Count(content, "<w:gridCol"); got != 1 {
		t.Fatalf("grid columns = %d, want 1 (label only, loop column removed)", got)
	}
	if !strings.Contains(content, `<w:gridCol w:w="3000"/>`) {
		t.Fatalf("expected the surviving grid column to keep the label width; content=%s", content)
	}

	for _, tr := range strings.Split(content, "<w:tr>")[1:] {
		row := tr
		if end := strings.Index(row, "</w:tr>"); end != -1 {
			row = row[:end]
		}
		if got := strings.Count(row, "<w:tc>"); got != 1 {
			t.Fatalf("row cell count = %d, want 1 (loop cell removed), row=%q", got, row)
		}
	}
}

// The exported ColumnLoopContractError interface requires importing this
// package. fayna classifies the error without that import via a
// structurally-matching marker method instead (the same pattern as
// StyleContractError for style-token errors): errors.As against a locally
// declared interface{ ColumnLoopContractError() bool } must still find it.
func TestColumnLoop_MarkerMethodDetectableWithoutImport(t *testing.T) {
	_, err := ProcessTemplate(createTestDocx(t, bodyDoc(columnLoopTable)), matrixData(64, 64))
	if err == nil {
		t.Fatal("expected an error for over-limit column loop")
	}
	var marker interface{ ColumnLoopContractError() bool }
	if !errors.As(err, &marker) || !marker.ColumnLoopContractError() {
		t.Fatalf("err = %v, want errors.As to find a ColumnLoopContractError() bool marker", err)
	}
}
