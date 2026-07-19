package doctemplate

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"
)

// ----------------------------------------------------------------------------
// Byte-stability golden — frozen live report-card emissions
// ----------------------------------------------------------------------------

// goldenReportCardTemplate returns a report-card-shaped docx: a body-level loop
// ({{#job_categories.academic.jobs}}) wrapping a per-subject heading and a table
// with a nested row loop ({{#outcome_criteria}}) — the live v1/v2 emission shape
// (cf. TestContractShape_BodyLoopWithNestedMapsAndRowLoop).
func goldenReportCardTemplate(t *testing.T) []byte {
	t.Helper()
	inner := `<w:p><w:r><w:t>Report Header</w:t></w:r></w:p>
<w:p><w:r><w:t>{{#job_categories.academic.jobs}}</w:t></w:r></w:p>
<w:p><w:r><w:t>Subject: {{job_template_name_display}}</w:t></w:r></w:p>
<w:tbl>
<w:tblPr><w:tblStyle w:val="TableGrid"/></w:tblPr>
<w:tblGrid><w:gridCol w:w="5000"/><w:gridCol w:w="3000"/></w:tblGrid>
<w:tr><w:tc><w:p><w:r><w:t>Criterion</w:t></w:r></w:p></w:tc><w:tc><w:p><w:r><w:t>Mark</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>{{#outcome_criteria}}</w:t></w:r></w:p></w:tc><w:tc><w:p><w:r><w:t></w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>{{outcome_criteria_label_display}}</w:t></w:r></w:p></w:tc><w:tc><w:p><w:r><w:t>{{mark}}</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>{{/outcome_criteria}}</w:t></w:r></w:p></w:tc><w:tc><w:p><w:r><w:t></w:t></w:r></w:p></w:tc></w:tr>
</w:tbl>
<w:p><w:r><w:t>{{/job_categories.academic.jobs}}</w:t></w:r></w:p>
<w:p><w:r><w:t>Report Footer</w:t></w:r></w:p>`
	return createTestDocx(t, bodyDoc(inner))
}

func goldenReportCardData() map[string]any {
	return map[string]any{
		"job_categories": map[string]any{
			"academic": map[string]any{
				"jobs": []any{
					map[string]any{
						"job_template_name_display": "Mathematics",
						"outcome_criteria": []any{
							map[string]any{"outcome_criteria_label_display": "Knowing", "mark": "7"},
							map[string]any{"outcome_criteria_label_display": "Applying", "mark": "6"},
						},
					},
					map[string]any{
						"job_template_name_display": "Science",
						"outcome_criteria": []any{
							map[string]any{"outcome_criteria_label_display": "Inquiring", "mark": "5"},
						},
					},
				},
			},
		},
	}
}

func goldenInvoiceData() map[string]any {
	return map[string]any{
		"client": map[string]any{"name": "Acme Corporation", "address": "123 Business Ave, Manila"},
		"date":   "2026-03-08",
		"notes":  "Payment due within 30 days",
		"total":  "PHP 18,500.00",
		"items": []any{
			map[string]any{"description": "Frontend Development", "amount": "5000"},
			map[string]any{"description": "Backend Development", "amount": "8000"},
			map[string]any{"description": "Database Design", "amount": "3500"},
			map[string]any{"description": "Code Review", "amount": "2000"},
		},
	}
}

func hashDocumentXML(t *testing.T, docx []byte) string {
	t.Helper()
	arch, err := ReadDocxBytes(docx)
	if err != nil {
		t.Fatalf("read docx: %v", err)
	}
	sum := sha256.Sum256([]byte(arch.Content))
	return hex.EncodeToString(sum[:])
}

// TestByteStability_FrozenEmissions pins the processed document.xml (the exact
// bytes rendered into the live report cards) to a SHA256 captured BEFORE the Q2
// pairing-stack + linear-cleanup refactor. The engine renders frozen v1/v2 report
// cards, so a change in these hashes is a byte-stability regression, not a
// cosmetic diff. Any intentional emission change must update these constants in
// the same commit.
func TestByteStability_FrozenEmissions(t *testing.T) {
	const (
		goldenInvoice    = "129fdd6745c04871e8fdba8ec5bfeed4d20732d5ddcb5686f21d96e9fd8231be"
		goldenReportCard = "be60e4c7b94cf732a646ce6ad1e041db52d51a9a407c567ce0247a82d1174e07"
	)

	inv, err := os.ReadFile("testdata/invoice-template.docx")
	if err != nil {
		t.Fatalf("read invoice template: %v", err)
	}
	invOut, err := ProcessTemplate(inv, goldenInvoiceData())
	if err != nil {
		t.Fatalf("process invoice: %v", err)
	}
	if got := hashDocumentXML(t, invOut); got != goldenInvoice {
		t.Errorf("invoice emission byte-stability regression:\n  want %s\n  got  %s", goldenInvoice, got)
	}

	rcOut, err := ProcessTemplate(goldenReportCardTemplate(t), goldenReportCardData())
	if err != nil {
		t.Fatalf("process report card: %v", err)
	}
	if got := hashDocumentXML(t, rcOut); got != goldenReportCard {
		t.Errorf("report-card emission byte-stability regression:\n  want %s\n  got  %s", goldenReportCard, got)
	}
}

// ----------------------------------------------------------------------------
// Q2(i) — table-loop pairing stack
// ----------------------------------------------------------------------------

// TestTableLoopWellNestedSelectsInner proves the pairing stack selects the
// INNERMOST complete pair as the table's row loop. A row sequence
// #outer, #inner, template, /inner, /outer expands the inner loop over its data
// while the enclosing #outer / /outer marker rows are blanked (handled as
// non-loop rows). Pre-fix, the last-opener-wins scan mispaired the markers.
func TestTableLoopWellNestedSelectsInner(t *testing.T) {
	inner := `<w:tbl>
<w:tblGrid><w:gridCol w:w="5000"/></w:tblGrid>
<w:tr><w:tc><w:p><w:r><w:t>HeaderCell</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>{{#outer}}</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>{{#inner}}</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>ROW {{field}}</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>{{/inner}}</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>{{/outer}}</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>StaticTail</w:t></w:r></w:p></w:tc></w:tr>
</w:tbl>`

	content := renderBody(t, inner, map[string]any{
		"inner": []any{
			map[string]any{"field": "A"},
			map[string]any{"field": "B"},
			map[string]any{"field": "C"},
		},
	})

	// Inner loop expanded exactly once per item.
	for _, v := range []string{"A", "B", "C"} {
		if !strings.Contains(content, "ROW "+v) {
			t.Errorf("expected inner-loop row for %q", v)
		}
	}
	if got := strings.Count(content, "ROW "); got != 3 {
		t.Errorf("expected inner loop expanded 3x, got ROW x%d", got)
	}
	// No markers (inner or outer) leak; outer markers are blanked.
	mustNotContain(t, content,
		"{{#outer}}", "{{/outer}}", "{{#inner}}", "{{/inner}}", "{{field}}")
	// Static rows survive.
	mustContain(t, content, "HeaderCell", "StaticTail")
}

// TestTableLoopCrossingCloseFailsClosed is the "crossing-close" regression: the
// exact row order #outer, #inner, /outer, /inner. A LIFO stack sees /outer while
// #inner is the nearest unclosed opener — an interleaved (malformed) close — so
// the table fails closed: no expansion, all four marker rows blanked, static
// rows preserved. Pre-fix, the last-opener-wins scan paired #inner with /inner,
// made the {{/outer}} row the template row, and left {{#outer}} raw in output.
func TestTableLoopCrossingCloseFailsClosed(t *testing.T) {
	inner := `<w:tbl>
<w:tblGrid><w:gridCol w:w="5000"/></w:tblGrid>
<w:tr><w:tc><w:p><w:r><w:t>HeaderCell</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>{{#outer}}</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>{{#inner}}</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>{{/outer}}</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>{{/inner}}</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>StaticTail</w:t></w:r></w:p></w:tc></w:tr>
</w:tbl>`

	content := renderBody(t, inner, map[string]any{
		"outer": []any{map[string]any{"x": "1"}},
		"inner": []any{map[string]any{"x": "2"}},
	})

	// Fail closed: no marker leaks (all blanked by the exact-match branch).
	mustNotContain(t, content, "{{#outer}}", "{{#inner}}", "{{/outer}}", "{{/inner}}")
	// Static rows preserved — no truncation, no mispaired removal.
	mustContain(t, content, "HeaderCell", "StaticTail")
}

// TestTableLoopStrayCloseFailsClosed covers a closer with no matching opener at
// all ({{/ghost}} first, then a well-formed {{#items}} loop). The stray closer
// makes the whole table fail closed rather than silently mispairing — every row
// is processed once and markers are blanked.
func TestTableLoopStrayCloseFailsClosed(t *testing.T) {
	inner := `<w:tbl>
<w:tblGrid><w:gridCol w:w="5000"/></w:tblGrid>
<w:tr><w:tc><w:p><w:r><w:t>HeaderCell</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>{{/ghost}}</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>{{#items}}</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>ROW {{field}}</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>{{/items}}</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>StaticTail</w:t></w:r></w:p></w:tc></w:tr>
</w:tbl>`

	content := renderBody(t, inner, map[string]any{
		"items": []any{map[string]any{"field": "A"}, map[string]any{"field": "B"}},
	})

	// Fail closed: the loop is NOT expanded (template row rendered once, verbatim).
	if got := strings.Count(content, "ROW "); got != 1 {
		t.Errorf("expected fail-closed (template row once), got ROW x%d", got)
	}
	mustNotContain(t, content, "{{/ghost}}", "{{#items}}", "{{/items}}")
	mustContain(t, content, "HeaderCell", "StaticTail")
}

// ----------------------------------------------------------------------------
// Q2(iii) — malformed marker structure fails closed AND zero-leaks
// ----------------------------------------------------------------------------
//
// The complete marker-stream scan classifies four malformed table shapes the old
// first-balanced-closer break missed or mis-handled — a trailing stray closer, a
// crossing close, an unclosed opener, and an unsupported second top-level pair —
// and fails them CLOSED (no expansion). The residual-token scrub then guarantees
// NO raw {{...}} reaches any part. Each shape is exercised with the marker in
// three physical layouts: exact (whole-cell), mixed-content (marker plus other
// cell text), and cross-run (marker split across <w:r> runs the way Word emits).

type markerForm int

const (
	formExact    markerForm = iota // {{#items}} is the entire cell text
	formMixed                      // Note {{#items}} tail — marker inside other text
	formCrossRun                   // {{ | #items | }} split across three runs
)

func (f markerForm) String() string {
	switch f {
	case formExact:
		return "exact"
	case formMixed:
		return "mixed-content"
	case formCrossRun:
		return "cross-run"
	default:
		return "unknown"
	}
}

// markerRow renders a loop marker (payload like "#items" or "/items") as a table
// row whose single cell carries the marker in the requested physical form.
func markerRow(payload string, form markerForm) string {
	var runs string
	switch form {
	case formExact:
		runs = "<w:r><w:t>{{" + payload + "}}</w:t></w:r>"
	case formMixed:
		runs = "<w:r><w:t>Note {{" + payload + "}} tail</w:t></w:r>"
	case formCrossRun:
		runs = "<w:r><w:t>{{</w:t></w:r><w:r><w:t>" + payload + "</w:t></w:r><w:r><w:t>}}</w:t></w:r>"
	}
	return "<w:tr><w:tc><w:p>" + runs + "</w:p></w:tc></w:tr>\n"
}

// staticCellRow is a plain (non-marker) row with the given text.
func staticCellRow(text string) string {
	return "<w:tr><w:tc><w:p><w:r><w:t>" + text + "</w:t></w:r></w:p></w:tc></w:tr>\n"
}

func malformedTable(rows string) string {
	return "<w:tbl>\n<w:tblGrid><w:gridCol w:w=\"5000\"/></w:tblGrid>\n" + rows + "</w:tbl>"
}

func TestTableLoopMalformedShapesZeroLeak(t *testing.T) {
	// A template row carrying an item-scope placeholder that CANNOT resolve at root
	// scope — the exact vector the Wave-6 review flagged as leaking {{field}}.
	tmpl := staticCellRow("ROW {{field}}")

	shapes := map[string]func(form markerForm) string{
		// Well-formed prefix (#items … /items) followed by a trailing stray closer.
		// Old code broke at /items and expanded; now the trailing /ghost is seen.
		"trailing-stray-closer": func(f markerForm) string {
			return malformedTable(staticCellRow("HeaderCell") +
				markerRow("#items", f) + tmpl + markerRow("/items", f) +
				markerRow("/ghost", f) + staticCellRow("StaticTail"))
		},
		// Crossing close: #outer,#inner,/outer,/inner — /outer arrives while #inner
		// is the nearest unclosed opener.
		"crossing-close": func(f markerForm) string {
			return malformedTable(staticCellRow("HeaderCell") +
				markerRow("#outer", f) + markerRow("#inner", f) + tmpl +
				markerRow("/outer", f) + markerRow("/inner", f) +
				staticCellRow("StaticTail"))
		},
		// Unclosed opener: #outer never closes (only the inner pair balances).
		"unclosed-opener": func(f markerForm) string {
			return malformedTable(staticCellRow("HeaderCell") +
				markerRow("#outer", f) + markerRow("#inner", f) + tmpl +
				markerRow("/inner", f) + staticCellRow("StaticTail"))
		},
		// Two separate top-level pairs in one table — explicitly unsupported.
		"second-top-level-loop": func(f markerForm) string {
			return malformedTable(staticCellRow("HeaderCell") +
				markerRow("#a", f) + tmpl + markerRow("/a", f) +
				markerRow("#b", f) + tmpl + markerRow("/b", f) +
				staticCellRow("StaticTail"))
		},
	}

	forms := []markerForm{formExact, formMixed, formCrossRun}

	for shapeName, build := range shapes {
		for _, form := range forms {
			t.Run(shapeName+"/"+form.String(), func(t *testing.T) {
				content := renderBody(t, build(form), map[string]any{
					// If any shape WRONGLY expanded, ROW would appear per item (>1).
					"items": []any{map[string]any{"field": "V1"}, map[string]any{"field": "V2"}},
					"a":     []any{map[string]any{"field": "A1"}, map[string]any{"field": "A2"}},
					"b":     []any{map[string]any{"field": "B1"}, map[string]any{"field": "B2"}},
				})

				// Zero leak against the BROAD sentinels — markers (any form),
				// unresolved {{field}}, and any mixed-content residue.
				assertNoResidualTokens(t, content)
				// Static content is never truncated.
				mustContain(t, content, "HeaderCell", "StaticTail")
				// Fail closed — the template row is rendered at most once per
				// occurrence, never expanded over the item slice.
				if got := strings.Count(content, "ROW "); got > strings.Count(build(form), "ROW ") {
					t.Errorf("template row was expanded (%d occurrences) — malformed table must fail closed", got)
				}
				// The would-be item values must never appear (no expansion happened).
				mustNotContain(t, content, "V1", "V2", "A1", "A2", "B1", "B2")
			})
		}
	}
}

// ----------------------------------------------------------------------------
// Q2(ii) — linear cleanup / expansion (no quadratic blow-up)
// ----------------------------------------------------------------------------

// buildRowLoopTable builds a one-table docx with a header row, a {{#items}}
// loop over `tmplRows` template rows, and a static tail row.
func buildRowLoopTable(tmplRows int) string {
	var sb strings.Builder
	sb.WriteString("<w:tbl>\n<w:tblGrid><w:gridCol w:w=\"5000\"/></w:tblGrid>\n")
	sb.WriteString("<w:tr><w:tc><w:p><w:r><w:t>HeaderCell</w:t></w:r></w:p></w:tc></w:tr>\n")
	sb.WriteString("<w:tr><w:tc><w:p><w:r><w:t>{{#items}}</w:t></w:r></w:p></w:tc></w:tr>\n")
	for i := 0; i < tmplRows; i++ {
		fmt.Fprintf(&sb, "<w:tr><w:tc><w:p><w:r><w:t>R%d {{field}}</w:t></w:r></w:p></w:tc></w:tr>\n", i)
	}
	sb.WriteString("<w:tr><w:tc><w:p><w:r><w:t>{{/items}}</w:t></w:r></w:p></w:tc></w:tr>\n")
	sb.WriteString("<w:tr><w:tc><w:p><w:r><w:t>StaticTail</w:t></w:r></w:p></w:tc></w:tr>\n")
	sb.WriteString("</w:tbl>")
	return sb.String()
}

// minRenderTime renders the workload `runs` times and returns the MINIMUM
// wall-clock duration. The minimum (not mean) is the standard noise-robust
// estimator for a deterministic workload: scheduling stalls, GC pauses, and
// other-process interference only ever ADD time, so the fastest of k runs is
// the closest observation of the true cost. Each iteration starts from a
// freshly-collected heap (runtime.GC()) so garbage carried over from earlier
// tests or the previous iteration cannot bill its collection time to this run —
// without that, the LARGER workload systematically absorbs more carried-over GC
// debt and the ratio drifts upward under full-suite memory pressure.
// Correctness is asserted by the caller on the returned content of the LAST run.
func minRenderTime(t *testing.T, runs int, inner string, data map[string]any) (time.Duration, string) {
	t.Helper()
	best := time.Duration(1<<63 - 1)
	var content string
	for i := 0; i < runs; i++ {
		runtime.GC()
		start := time.Now()
		content = renderBody(t, inner, data)
		if d := time.Since(start); d < best {
			best = d
		}
	}
	return best, content
}

// complexityRatioCeiling is the pass/fail line for the size-doubling checks
// below. Doubling the input size multiplies a linear stage's cost by ~2 (often
// less, since fixed per-render overhead — ZIP envelope, XML parse setup — does
// not double), while a re-introduced quadratic stage multiplies it by ~4. The
// review measured the OLD quadratic cleanup at 0.93–1.17 s for 8,000 rows, so
// at these sizes the timed region dwarfs timer resolution and the two growth
// classes are separated by a full factor of ~2 around this threshold: generous
// enough not to flake on a loaded/slow CI box (linear needs a >1.5x relative
// slowdown of ONLY the doubled run, surviving the best-of-k minimum AND the
// per-run runtime.GC() heap leveling, to false-positive; observed linear ratios
// are 1.7–2.4), tight enough that quadratic growth fails decisively (the OLD
// engine measured ratio 3.71 on the cleanup shape via the provenance harness) —
// the design the Wave-6 review required in place of the old 20-second absolute
// ceilings, which the quadratic implementation passed with enormous headroom.
const complexityRatioCeiling = 3.0

// perfRuns is the best-of-k repetition count for each size.
const perfRuns = 5

// TestTableLoopLargeExpansionBounded guards the expansion path's growth class
// with a size-doubling complexity-ratio check: render N items and 2N items
// (40 template rows each → 8,000 vs 16,000 expanded rows) and require
// time(2N)/time(N) < complexityRatioCeiling. The old ascending per-row
// RemoveChild + insert-before-anchor shape was the O(n²) risk here; the linear
// splice keeps the ratio ~2. Correctness of the doubled expansion is asserted
// on the same run that is timed.
func TestTableLoopLargeExpansionBounded(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping expansion complexity-ratio test in -short mode")
	}
	const (
		tmplRows = 40
		items    = 200 // N: 40 * 200 = 8,000 expanded rows; 2N doubles the items
	)
	inner := buildRowLoopTable(tmplRows)

	makeItems := func(n int) map[string]any {
		itemData := make([]any, n)
		for i := range itemData {
			itemData[i] = map[string]any{"field": fmt.Sprintf("V%d", i)}
		}
		return map[string]any{"items": itemData}
	}

	timeN, _ := minRenderTime(t, perfRuns, inner, makeItems(items))
	time2N, content := minRenderTime(t, perfRuns, inner, makeItems(2*items))

	// Correctness on the doubled (timed) workload: header + tail + every
	// template row expanded per item, zero residual tokens.
	mustContain(t, content, "HeaderCell", "StaticTail")
	assertNoResidualTokens(t, content)
	// Spot-check first and last item, first and last template row.
	mustContain(t, content, "R0 V0", "R39 V0", "R0 V399", "R39 V399")
	// Every expanded data cell present: tmplRows * 2N "R" rows.
	if got := strings.Count(content, "</w:tr>"); got < tmplRows*2*items {
		t.Errorf("expected at least %d expanded rows, got %d", tmplRows*2*items, got)
	}

	ratio := float64(time2N) / float64(timeN)
	if ratio >= complexityRatioCeiling {
		t.Errorf("expansion time grew %.2fx when input doubled (N=%s, 2N=%s); ratio >= %.1f indicates super-linear (quadratic ~4x) growth",
			ratio, timeN, time2N, complexityRatioCeiling)
	}
	t.Logf("expansion: N=%d rows in %s, 2N=%d rows in %s, ratio %.2f (ceiling %.1f)",
		tmplRows*items, timeN, tmplRows*2*items, time2N, ratio, complexityRatioCeiling)
}

// TestTableLoopLargeFailClosedCleanupBounded is the fail-closed cleanup
// analogue, on the shape where the review MEASURED the old code quadratic
// (0.93–1.17 s at 8,000 rows vs 73–93 ms linear): a large loop block whose key
// resolves to a non-slice must be removed whole. Size-doubling ratio check:
// removing a 2N-row block must cost < complexityRatioCeiling times an N-row
// block. The old ascending removal was O(block²) — ratio ~4, decisively
// failing; the descending-detach rebuild is linear — ratio ~2.
func TestTableLoopLargeFailClosedCleanupBounded(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping cleanup complexity-ratio test in -short mode")
	}
	// N = the review's measured quadratic size (old engine: 0.93–1.17 s here);
	// 2N doubles it, where the old engine would spend ~4x N's time (measured
	// ratio 3.71 via the provenance harness) while the linear rebuild spends ~2x.
	// The large absolute timed regions (tens of ms at N even on this machine)
	// keep the ratio steady against scheduler/GC noise.
	const blockRows = 8000

	// items resolves to a scalar → not a slice → fail-closed whole-block removal.
	data := map[string]any{"items": "not-a-slice"}

	timeN, _ := minRenderTime(t, perfRuns, buildRowLoopTable(blockRows), data)
	time2N, content := minRenderTime(t, perfRuns, buildRowLoopTable(2*blockRows), data)

	// Correctness on the doubled (timed) workload: the entire loop block
	// (markers + all template rows) is gone; static rows survive; zero residual.
	mustNotContain(t, content, "R0 ", "R15999 ")
	assertNoResidualTokens(t, content)
	mustContain(t, content, "HeaderCell", "StaticTail")

	ratio := float64(time2N) / float64(timeN)
	if ratio >= complexityRatioCeiling {
		t.Errorf("fail-closed cleanup time grew %.2fx when block doubled (N=%s, 2N=%s); ratio >= %.1f indicates super-linear (quadratic ~4x) growth",
			ratio, timeN, time2N, complexityRatioCeiling)
	}
	t.Logf("cleanup: N=%d-row block in %s, 2N=%d-row block in %s, ratio %.2f (ceiling %.1f)",
		blockRows, timeN, 2*blockRows, time2N, ratio, complexityRatioCeiling)
}
