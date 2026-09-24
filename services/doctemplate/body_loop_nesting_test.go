package doctemplate

import (
	"strings"
	"testing"
)

func para(text string) string {
	return `<w:p><w:r><w:t xml:space="preserve">` + text + `</w:t></w:r></w:p>`
}

// indexOrder asserts each want appears in content, in the given order.
func indexOrder(t *testing.T, content string, wants ...string) {
	t.Helper()
	from := 0
	for _, w := range wants {
		i := strings.Index(content[from:], w)
		if i < 0 {
			t.Fatalf("expected %q after offset %d (order %q)", w, from, wants)
		}
		from += i + len(w)
	}
}

// TestBodyLoopNestedParagraphLoops is the Progress Report shape: a subject loop
// whose block holds two paragraph loops (criteria, then rating descriptions).
// Before the recursive expander only the outer loop expanded; the inner markers
// were blanked and their item tokens scrubbed, leaving punctuation behind.
func TestBodyLoopNestedParagraphLoops(t *testing.T) {
	inner := para("Cover {{student_name}}") +
		para("{{#jobs}}") +
		para("SUBJECT {{name}}") +
		para("{{#assessments}}") +
		para("CRIT {{code}}: {{label}} /{{max}}") +
		para("{{/assessments}}") +
		para("{{#rating_descriptions}}") +
		para("BAND {{code}} ({{from}} → {{to}})") +
		para("{{/rating_descriptions}}") +
		para("TOTAL {{total_max}}") +
		para("{{/jobs}}") +
		para("Footer")

	content := renderBody(t, inner, map[string]any{
		"student_name": "Ana",
		"jobs": []any{
			map[string]any{
				"name": "Arts", "total_max": "32",
				"assessments": []any{
					map[string]any{"code": "A", "label": "Investigating", "max": "8"},
					map[string]any{"code": "B", "label": "Developing", "max": "8"},
				},
				"rating_descriptions": []any{
					map[string]any{"code": "A", "from": "0", "to": "2"},
				},
			},
			map[string]any{
				"name": "Math", "total_max": "32",
				"assessments": []any{
					map[string]any{"code": "A", "label": "Knowing", "max": "8"},
				},
				"rating_descriptions": []any{},
			},
		},
	})

	indexOrder(t, content,
		"Cover Ana",
		"SUBJECT Arts", "CRIT A: Investigating /8", "CRIT B: Developing /8", "BAND A (0 → 2)", "TOTAL 32",
		"SUBJECT Math", "CRIT A: Knowing /8", "TOTAL 32",
		"Footer")
	if got := strings.Count(content, "CRIT "); got != 3 {
		t.Errorf("criteria rows = %d, want 3", got)
	}
	if got := strings.Count(content, "BAND "); got != 1 {
		t.Errorf("band rows = %d, want 1 (Math has none)", got)
	}
	mustNotContain(t, content, "( → )")
	assertNoResidualTokens(t, content)
}

// TestBodyLoopSiblingLoops: two top-level loops in one body both expand (the
// old scanner only honoured the first; the second rendered once at root scope).
func TestBodyLoopSiblingLoops(t *testing.T) {
	inner := para("{{#jobs}}") + para("JOB {{name}}") + para("{{/jobs}}") +
		para("Between {{title}}") +
		para("{{#outcome_sections}}") + para("SECTION {{label}}") + para("{{/outcome_sections}}")

	content := renderBody(t, inner, map[string]any{
		"title": "T",
		"jobs":  []any{map[string]any{"name": "J1"}, map[string]any{"name": "J2"}},
		"outcome_sections": []any{
			map[string]any{"label": "S1"}, map[string]any{"label": "S2"}, map[string]any{"label": "S3"},
		},
	})
	indexOrder(t, content, "JOB J1", "JOB J2", "Between T", "SECTION S1", "SECTION S2", "SECTION S3")
	assertNoResidualTokens(t, content)
}

// TestBodyLoopNestedWithTableRowLoop: a table-row loop inside a nested body loop
// resolves against the innermost body item.
func TestBodyLoopNestedWithTableRowLoop(t *testing.T) {
	inner := para("{{#groups}}") + para("GROUP {{name}}") +
		para("{{#members}}") + para("MEMBER {{name}}") +
		`<w:tbl><w:tblGrid><w:gridCol w:w="5000"/></w:tblGrid>
<w:tr><w:tc>` + para("{{#scores}}") + `</w:tc></w:tr>
<w:tr><w:tc>` + para("SCORE {{value}}") + `</w:tc></w:tr>
<w:tr><w:tc>` + para("{{/scores}}") + `</w:tc></w:tr>
</w:tbl>` +
		para("{{/members}}") + para("{{/groups}}")

	content := renderBody(t, inner, map[string]any{
		"groups": []any{map[string]any{
			"name": "G1",
			"members": []any{
				map[string]any{"name": "M1", "scores": []any{map[string]any{"value": "1"}, map[string]any{"value": "2"}}},
				map[string]any{"name": "M2", "scores": []any{map[string]any{"value": "3"}}},
			},
		}},
	})
	indexOrder(t, content, "GROUP G1", "MEMBER M1", "SCORE 1", "SCORE 2", "MEMBER M2", "SCORE 3")
	if got := strings.Count(content, "<w:tbl>"); got != 2 {
		t.Errorf("tables = %d, want one per member (2)", got)
	}
	assertNoResidualTokens(t, content)
}

// TestBodyLoopNestedSameKey: a same-key opener inside a loop nests rather than
// closing early; the inner loop reads the item's own "items".
func TestBodyLoopNestedSameKey(t *testing.T) {
	inner := para("{{#items}}") + para("OUTER {{name}}") +
		para("{{#items}}") + para("INNER {{name}}") + para("{{/items}}") +
		para("END {{name}}") + para("{{/items}}")

	content := renderBody(t, inner, map[string]any{
		"items": []any{map[string]any{
			"name":  "O",
			"items": []any{map[string]any{"name": "I1"}, map[string]any{"name": "I2"}},
		}},
	})
	indexOrder(t, content, "OUTER O", "INNER I1", "INNER I2", "END O")
	assertNoResidualTokens(t, content)
}

// TestBodyLoopInnerMismatchedCloseFailsClosed: an inner opener without its own
// closer does not loop; the outer loop still expands and nothing leaks.
func TestBodyLoopInnerMismatchedCloseFailsClosed(t *testing.T) {
	inner := para("{{#jobs}}") + para("JOB {{name}}") +
		para("{{#assessments}}") + para("CRIT {{code}}") + para("{{/other}}") +
		para("{{/jobs}}") + para("After")

	content := renderBody(t, inner, map[string]any{
		"jobs": []any{
			map[string]any{"name": "J1", "assessments": []any{map[string]any{"code": "A"}, map[string]any{"code": "B"}}},
			map[string]any{"name": "J2"},
		},
	})
	indexOrder(t, content, "JOB J1", "JOB J2", "After")
	// Rendered once per job at job scope, never iterated over assessments.
	if got := strings.Count(content, "CRIT"); got != 2 {
		t.Errorf("CRIT rendered %d times, want 2 (once per job, no inner expansion)", got)
	}
	assertNoResidualTokens(t, content)
}

// TestBodyLoopNestedMissingOrEmptyData: an inner loop with absent, empty, or
// non-slice data renders nothing for that item and leaks nothing.
func TestBodyLoopNestedMissingOrEmptyData(t *testing.T) {
	inner := para("{{#jobs}}") + para("JOB {{name}}") +
		para("{{#assessments}}") + para("CRIT {{code}}") + para("{{/assessments}}") +
		para("{{/jobs}}")

	content := renderBody(t, inner, map[string]any{
		"jobs": []any{
			map[string]any{"name": "Absent"},
			map[string]any{"name": "Empty", "assessments": []any{}},
			map[string]any{"name": "Scalar", "assessments": "nope"},
		},
	})
	indexOrder(t, content, "JOB Absent", "JOB Empty", "JOB Scalar")
	mustNotContain(t, content, "CRIT")
	assertNoResidualTokens(t, content)
}

// TestBodyLoopDepthBounded: nesting past maxBodyLoopDepth is not expanded
// further (markers are blanked, tokens scrubbed) instead of recursing.
func TestBodyLoopDepthBounded(t *testing.T) {
	var open, close strings.Builder
	levels := maxBodyLoopDepth + 2
	for i := 0; i < levels; i++ {
		open.WriteString(para("{{#n}}"))
		close.WriteString(para("{{/n}}"))
	}
	inner := open.String() + para("LEAF {{v}}") + close.String()

	var build func(int) map[string]any
	build = func(level int) map[string]any {
		if level == 0 {
			return map[string]any{"v": "x"}
		}
		return map[string]any{"v": "x", "n": []any{build(level - 1)}}
	}
	content := renderBody(t, inner, build(levels))
	assertNoResidualTokens(t, content)
	if got := strings.Count(content, "LEAF"); got != 1 {
		t.Errorf("LEAF rendered %d times, want exactly 1", got)
	}
}
