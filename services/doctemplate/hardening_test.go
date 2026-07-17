package doctemplate

import (
	"strings"
	"testing"
)

// bodyDoc wraps a body fragment in a minimal w:document envelope.
func bodyDoc(inner string) string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
<w:body>
` + inner + `
</w:body>
</w:document>`
}

// renderBody builds a docx from the body fragment, runs it through the engine,
// and returns the processed word/document.xml as a string.
func renderBody(t *testing.T, inner string, data map[string]any) string {
	t.Helper()
	template := createTestDocx(t, bodyDoc(inner))
	result, err := ProcessTemplate(template, data)
	if err != nil {
		t.Fatalf("ProcessTemplate failed: %v", err)
	}
	archive, err := ReadDocxBytes(result)
	if err != nil {
		t.Fatalf("failed to read output docx: %v", err)
	}
	return archive.Content
}

func mustContain(t *testing.T, content string, wants ...string) {
	t.Helper()
	for _, w := range wants {
		if !strings.Contains(content, w) {
			t.Errorf("expected output to contain %q, but it did not", w)
		}
	}
}

func mustNotContain(t *testing.T, content string, unwanted ...string) {
	t.Helper()
	for _, u := range unwanted {
		if strings.Contains(content, u) {
			t.Errorf("expected output NOT to contain %q, but it did", u)
		}
	}
}

// TestContractShape_BodyLoopWithNestedMapsAndRowLoop exercises the exact
// converged placeholder contract: a dot-path BODY loop opener
// ({{#job_categories.academic.jobs}}) whose items carry nested maps read via
// dot path ({{job_template_phases.s1.task_outcome_numeric_value_max_derived}})
// AND a nested table-row loop ({{#outcome_criteria}}) whose row items also
// carry nested maps read via dot path.
func TestContractShape_BodyLoopWithNestedMapsAndRowLoop(t *testing.T) {
	inner := `<w:p><w:r><w:t>Report Header</w:t></w:r></w:p>
<w:p><w:r><w:t>{{#job_categories.academic.jobs}}</w:t></w:r></w:p>
<w:p><w:r><w:t>Subject: {{job_template_name_display}}</w:t></w:r></w:p>
<w:p><w:r><w:t>PhaseMax: {{job_template_phases.s1.task_outcome_numeric_value_max_derived}}</w:t></w:r></w:p>
<w:tbl>
<w:tblPr><w:tblStyle w:val="TableGrid"/></w:tblPr>
<w:tblGrid><w:gridCol w:w="5000"/><w:gridCol w:w="3000"/></w:tblGrid>
<w:tr><w:tc><w:p><w:r><w:t>Criterion</w:t></w:r></w:p></w:tc><w:tc><w:p><w:r><w:t>Mark</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>{{#outcome_criteria}}</w:t></w:r></w:p></w:tc><w:tc><w:p><w:r><w:t></w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>{{outcome_criteria_label_display}}</w:t></w:r></w:p></w:tc><w:tc><w:p><w:r><w:t>{{job_template_phases.s1.task_outcome_numeric_value_max_derived}}</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>{{/outcome_criteria}}</w:t></w:r></w:p></w:tc><w:tc><w:p><w:r><w:t></w:t></w:r></w:p></w:tc></w:tr>
</w:tbl>
<w:p><w:r><w:t>{{/job_categories.academic.jobs}}</w:t></w:r></w:p>
<w:p><w:r><w:t>Report Footer</w:t></w:r></w:p>`

	data := map[string]any{
		"job_categories": map[string]any{
			"academic": map[string]any{
				"jobs": []any{
					map[string]any{
						"job_template_name_display": "Mathematics",
						"job_template_phases": map[string]any{
							"s1": map[string]any{"task_outcome_numeric_value_max_derived": "MATHS1MAX"},
						},
						"outcome_criteria": []any{
							map[string]any{
								"outcome_criteria_label_display": "CritKnowing",
								"job_template_phases": map[string]any{
									"s1": map[string]any{"task_outcome_numeric_value_max_derived": "KNOWMAX"},
								},
							},
							map[string]any{
								"outcome_criteria_label_display": "CritApplying",
								"job_template_phases": map[string]any{
									"s1": map[string]any{"task_outcome_numeric_value_max_derived": "APPLYMAX"},
								},
							},
						},
					},
					map[string]any{
						"job_template_name_display": "Science",
						"job_template_phases": map[string]any{
							"s1": map[string]any{"task_outcome_numeric_value_max_derived": "SCIS1MAX"},
						},
						"outcome_criteria": []any{
							map[string]any{
								"outcome_criteria_label_display": "CritInquiring",
								"job_template_phases": map[string]any{
									"s1": map[string]any{"task_outcome_numeric_value_max_derived": "INQMAX"},
								},
							},
						},
					},
				},
			},
		},
	}

	content := renderBody(t, inner, data)

	// Body loop iterated over both jobs (dot-path opener resolves).
	mustContain(t, content, "Mathematics", "Science")
	// Nested-map scalar read at BODY-ITEM scope.
	mustContain(t, content, "MATHS1MAX", "SCIS1MAX")
	// Nested table-row loop iterated over each job's criteria.
	mustContain(t, content, "CritKnowing", "CritApplying", "CritInquiring")
	// Nested-map scalar read at ROW-ITEM scope.
	mustContain(t, content, "KNOWMAX", "APPLYMAX", "INQMAX")
	// Static content preserved.
	mustContain(t, content, "Report Header", "Report Footer", "Criterion", "Mark")
	// No markers or raw leaf placeholders leak.
	mustNotContain(t, content,
		"{{#job_categories.academic.jobs}}", "{{/job_categories.academic.jobs}}",
		"{{#outcome_criteria}}", "{{/outcome_criteria}}",
		"{{job_template_name_display}}", "{{outcome_criteria_label_display}}",
		"{{job_template_phases.s1.task_outcome_numeric_value_max_derived}}",
	)

	// Science has exactly one criterion; Mathematics has two. The row loop
	// clones per item, so the single-criterion job must not spill extra rows.
	if got := strings.Count(content, "CritInquiring"); got != 1 {
		t.Errorf("expected 'CritInquiring' exactly once, got %d", got)
	}
}

// TestDeepScalarDotPath verifies a five-segment dot-path scalar resolves through
// arbitrarily nested map[string]any values with no depth cap.
func TestDeepScalarDotPath(t *testing.T) {
	inner := `<w:p><w:r><w:t>Deep: {{a.b.c.d.e}}</w:t></w:r></w:p>`
	data := map[string]any{
		"a": map[string]any{
			"b": map[string]any{
				"c": map[string]any{
					"d": map[string]any{
						"e": "DEEPLEAF",
					},
				},
			},
		},
	}
	content := renderBody(t, inner, data)
	mustContain(t, content, "DEEPLEAF")
	mustNotContain(t, content, "{{a.b.c.d.e}}")
}

// TestTableLoopValueNotASlice covers the fail-closed cleanup for a table-row
// loop whose key resolves to a value that is neither []any nor
// []map[string]any: a map, a scalar, and a wrong-typed slice. In every case the
// marker rows and the template row must be removed (no marker/placeholder leak),
// while the static header and tail rows survive.
func TestTableLoopValueNotASlice(t *testing.T) {
	inner := `<w:tbl>
<w:tblGrid><w:gridCol w:w="5000"/></w:tblGrid>
<w:tr><w:tc><w:p><w:r><w:t>HeaderCell</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>{{#items}}</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>ROWSTART {{field}} ROWEND</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>{{/items}}</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>StaticTail</w:t></w:r></w:p></w:tc></w:tr>
</w:tbl>`

	cases := map[string]any{
		"map":         map[string]any{"field": "x"},
		"scalar":      "some-scalar",
		"wrong-slice": []string{"a", "b"},
	}

	for name, itemsVal := range cases {
		t.Run(name, func(t *testing.T) {
			content := renderBody(t, inner, map[string]any{"items": itemsVal})
			// Markers and the template row must be gone entirely.
			mustNotContain(t, content,
				"{{#items}}", "{{/items}}", "{{field}}", "ROWSTART", "ROWEND")
			// Static rows survive.
			mustContain(t, content, "HeaderCell", "StaticTail")
		})
	}
}

// TestBodyLoopValueNotASlice is the body-loop analogue: a body-level loop key
// resolving to a map produces zero iterations and the markers are removed
// (resolveLoopData returns nil), never leaking markers or the item placeholder.
func TestBodyLoopValueNotASlice(t *testing.T) {
	inner := `<w:p><w:r><w:t>Before</w:t></w:r></w:p>
<w:p><w:r><w:t>{{#rows}}</w:t></w:r></w:p>
<w:p><w:r><w:t>ITEMROW {{field}}</w:t></w:r></w:p>
<w:p><w:r><w:t>{{/rows}}</w:t></w:r></w:p>
<w:p><w:r><w:t>After</w:t></w:r></w:p>`

	content := renderBody(t, inner, map[string]any{
		"rows": map[string]any{"field": "x"}, // a map, not a slice
	})
	mustNotContain(t, content, "{{#rows}}", "{{/rows}}", "{{field}}", "ITEMROW")
	mustContain(t, content, "Before", "After")
}

// TestMismatchedBodyCloseMarker documents the fail-closed pairing rule for body
// loops: an opener {{#items}} with only a {{/other}} present (no {{/items}}) is
// treated as NO loop — the region is rendered once at root scope and the stray
// marker paragraphs are blanked, rather than the wrong close marker truncating
// or wrongly closing the block. Pre-fix, {{/other}} would have closed the loop
// and expanded it len(items) times.
func TestMismatchedBodyCloseMarker(t *testing.T) {
	inner := `<w:p><w:r><w:t>Before</w:t></w:r></w:p>
<w:p><w:r><w:t>{{#items}}</w:t></w:r></w:p>
<w:p><w:r><w:t>ITEMROW {{title}}</w:t></w:r></w:p>
<w:p><w:r><w:t>{{/other}}</w:t></w:r></w:p>
<w:p><w:r><w:t>After</w:t></w:r></w:p>`

	content := renderBody(t, inner, map[string]any{
		"items": []any{
			map[string]any{"title": "A"},
			map[string]any{"title": "B"},
			map[string]any{"title": "C"},
		},
	})

	// Stray markers are blanked (exact-match branch), not left raw.
	mustNotContain(t, content, "{{#items}}", "{{/other}}")
	// No truncation: content before and after both survive.
	mustContain(t, content, "Before", "After")
	// Treated as NO loop → rendered exactly once, not expanded 3x.
	if got := strings.Count(content, "ITEMROW"); got != 1 {
		t.Errorf("expected the region rendered once (no loop), got ITEMROW x%d", got)
	}
	// Item-scope placeholder leaks verbatim under root scope (no fallback).
	mustContain(t, content, "{{title}}")
}

// TestMismatchedTableCloseMarker is the table-row analogue of the pairing rule:
// {{#items}} with only {{/other}} present is treated as a table with no loop —
// rows processed once for simple placeholders, marker cells blanked, static
// rows preserved. Pre-fix, {{/other}} would have closed the loop and expanded
// the template row per item.
func TestMismatchedTableCloseMarker(t *testing.T) {
	inner := `<w:tbl>
<w:tblGrid><w:gridCol w:w="5000"/></w:tblGrid>
<w:tr><w:tc><w:p><w:r><w:t>HeaderCell</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>{{#items}}</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>ITEMROW {{field}}</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>{{/other}}</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>StaticTail</w:t></w:r></w:p></w:tc></w:tr>
</w:tbl>`

	content := renderBody(t, inner, map[string]any{
		"items": []any{
			map[string]any{"field": "X"},
			map[string]any{"field": "Y"},
		},
	})

	// Marker cells blanked; no raw markers survive.
	mustNotContain(t, content, "{{#items}}", "{{/other}}")
	// Static rows preserved.
	mustContain(t, content, "HeaderCell", "StaticTail")
	// Treated as NO loop → template row rendered once, not expanded per item.
	if got := strings.Count(content, "ITEMROW"); got != 1 {
		t.Errorf("expected the template row rendered once (no loop), got ITEMROW x%d", got)
	}
	// Item-scope placeholder leaks verbatim under root scope (no fallback).
	mustContain(t, content, "{{field}}")
}

// TestMissingLeafInsideLoopLeaksVerbatim documents the no-fallback law: a
// placeholder inside a loop item that references a leaf absent from the item map
// is left verbatim (no parent/root scope chain), while a present leaf resolves.
func TestMissingLeafInsideLoopLeaksVerbatim(t *testing.T) {
	inner := `<w:p><w:r><w:t>{{#rows}}</w:t></w:r></w:p>
<w:p><w:r><w:t>PRESENT {{present}} MISSING {{absent.leaf}}</w:t></w:r></w:p>
<w:p><w:r><w:t>{{/rows}}</w:t></w:r></w:p>`

	content := renderBody(t, inner, map[string]any{
		"rows": []any{
			map[string]any{"present": "HERE"},
		},
	})

	// Present leaf resolves.
	mustContain(t, content, "HERE")
	// Missing leaf leaks verbatim — no fallback to any parent/root scope.
	mustContain(t, content, "{{absent.leaf}}")
	// Loop markers are still stripped normally.
	mustNotContain(t, content, "{{#rows}}", "{{/rows}}")
}
