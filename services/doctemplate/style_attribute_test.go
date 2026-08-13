package doctemplate

import (
	"errors"
	"strings"
	"testing"
)

type styleTokenContractErr interface {
	IsStyleTokenContractError() bool
}

func TestStyleAttributesResolveAtTableRowScope(t *testing.T) {
	inner := `<w:tbl>
<w:tblGrid><w:gridCol w:w="7000"/><w:gridCol w:w="1000"/></w:tblGrid>
<w:tr><w:tc><w:p><w:r><w:t>Heading</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>{{#rows}}</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:tcPr><w:shd w:val="clear" w:color="auto" w:fill="{{cell_fill}}"/></w:tcPr><w:p><w:r><w:rPr><w:color w:val="{{row_color}}"/><w:b w:val="{{row_bold}}"/></w:rPr><w:t>{{label}}</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>{{/rows}}</w:t></w:r></w:p></w:tc></w:tr>
</w:tbl>`

	content := renderBody(t, inner, map[string]any{
		"rows": []any{
			map[string]any{
				"cell_fill": "123456",
				"row_color": "ABCDEF",
				"row_bold":  "true",
				"label":     "Alpha",
			},
			map[string]any{
				"cell_fill": "654321",
				"row_color": "FEDCBA",
				"row_bold":  "false",
				"label":     "Beta",
			},
		},
	})

	mustContain(t, content, `w:fill="123456"`, `w:color w:val="ABCDEF"`, `w:b w:val="true"`, `Alpha`)
	mustContain(t, content, `w:fill="654321"`, `w:color w:val="FEDCBA"`, `w:b w:val="false"`, `Beta`)
	mustNotContain(t, content,
		"{{cell_fill}}", "{{row_color}}", "{{row_bold}}", "{{label}}")
}

func TestStyleAttributesResolveByNamespaceURIWithAlternatePrefix(t *testing.T) {
	t.Run("alternate prefix bound to WordprocessingML resolves", func(t *testing.T) {
		inner := `<w:tbl xmlns:x="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
<w:tblGrid><w:gridCol w:w="5000"/></w:tblGrid>
<w:tr><w:tc><w:p><w:r><w:t>{{#rows}}</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:tcPr><x:shd x:fill="{{cell_fill}}"/></w:tcPr><w:p><w:r><w:rPr><x:color x:val="{{row_color}}"/><x:b x:val="{{row_bold}}"/></w:rPr><w:t>{{label}}</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>{{/rows}}</w:t></w:r></w:p></w:tc></w:tr>
</w:tbl>`

		content := renderBody(t, inner, map[string]any{
			"rows": []any{map[string]any{
				"cell_fill": "123456",
				"row_color": "ABCDEF",
				"row_bold":  "true",
				"label":     "Alternate",
			}},
		})

		mustContain(t, content, `x:fill="123456"`, `x:color x:val="ABCDEF"`, `x:b x:val="true"`, "Alternate")
		mustNotContain(t, content, "{{cell_fill}}", "{{row_color}}", "{{row_bold}}", "{{label}}")
	})

	t.Run("same local names in another URI stay outside whitelist", func(t *testing.T) {
		inner := `<w:p xmlns:x="urn:ichizen:hostile-template"><w:r><w:rPr><x:color x:val="{{title_color}}"/></w:rPr><w:t>Title</w:t></w:r></w:p>`
		content := renderBody(t, inner, map[string]any{"title_color": "ABCDEF"})
		mustContain(t, content, `x:val="{{title_color}}"`)
	})
}

func TestStyleAttributeInterpolationIsWhitelistedAndWholeTokenOnly(t *testing.T) {
	t.Run("partial token on whitelisted attribute remains untouched", func(t *testing.T) {
		inner := `<w:tbl>
<w:tblGrid><w:gridCol w:w="5000"/></w:tblGrid>
<w:tr><w:tc><w:tcPr><w:shd w:val="clear" w:color="auto" w:fill="prefix-{{cell_fill}}"/></w:tcPr><w:p><w:r><w:rPr><w:color w:val="FFFFFF"/><w:b/></w:rPr><w:t>Cell</w:t></w:r></w:p></w:tc></w:tr>
</w:tbl>`

		content := renderBody(t, inner, map[string]any{
			"cell_fill": "123456",
		})
		mustContain(t, content, `w:fill="prefix-{{cell_fill}}"`)
	})

	t.Run("non-whitelisted token stays literal", func(t *testing.T) {
		inner := `<w:tbl>
<w:tblGrid><w:gridCol w:w="5000"/></w:tblGrid>
<w:tr><w:tc><w:tcPr><w:tcW w:w="{{cell_width}}"/></w:tcPr><w:p><w:r><w:t>Cell</w:t></w:r></w:p></w:tc></w:tr>
</w:tbl>`

		content := renderBody(t, inner, map[string]any{
			"cell_width": "2500",
		})

		mustContain(t, content, "{{cell_width}}")
		mustNotContain(t, content, "{{#cell_width}}", "{{/cell_width}}")
	})

	t.Run("root paragraph uses root data scope", func(t *testing.T) {
		inner := `<w:p><w:r><w:rPr><w:color w:val="{{title_color}}"/><w:b w:val="{{title_bold}}"/></w:rPr><w:t>Title</w:t></w:r></w:p>`
		content := renderBody(t, inner, map[string]any{
			"title_color": "ABCDEF",
			"title_bold":  false,
		})
		mustContain(t, content, `w:color w:val="ABCDEF"`, `w:b w:val="false"`)
	})
}

func TestProcessTemplateRejectsMissingOrInvalidStyleAttributeValue(t *testing.T) {
	t.Run("missing style token value", func(t *testing.T) {
		inner := `<w:tbl>
<w:tblGrid><w:gridCol w:w="5000"/></w:tblGrid>
<w:tr><w:tc><w:tcPr><w:shd w:val="clear" w:color="auto" w:fill="{{cell_fill}}"/></w:tcPr><w:p><w:r><w:rPr><w:color w:val="{{row_color}}"/><w:b w:val="{{row_bold}}"/></w:rPr><w:t>Cell</w:t></w:r></w:p></w:tc></w:tr>
</w:tbl>`

		template := createTestDocx(t, bodyDoc(inner))
		_, err := ProcessTemplate(template, map[string]any{
			"row_color": "ABCDEF",
			"row_bold":  "true",
		})
		if err == nil {
			t.Fatal("expected missing style-token error")
		}

		var marker styleTokenContractErr
		if !errors.As(err, &marker) {
			t.Fatalf("expected StyleTokenContractError from ProcessTemplate, got: %v", err)
		}
	})

	t.Run("invalid style token value", func(t *testing.T) {
		inner := `<w:tbl>
<w:tblGrid><w:gridCol w:w="5000"/></w:tblGrid>
<w:tr><w:tc><w:tcPr><w:shd w:val="clear" w:color="auto" w:fill="{{cell_fill}}"/></w:tcPr><w:p><w:r><w:rPr><w:color w:val="{{row_color}}"/><w:b w:val="{{row_bold}}"/></w:rPr><w:t>Cell</w:t></w:r></w:p></w:tc></w:tr>
</w:tbl>`

		_, err := ProcessTemplate(createTestDocx(t, bodyDoc(inner)), map[string]any{
			"cell_fill": "bad",
			"row_color": "zzzzzz",
			"row_bold":  "true",
		})
		if err == nil {
			t.Fatal("expected invalid style-token error")
		}

		var marker styleTokenContractErr
		if !errors.As(err, &marker) {
			t.Fatalf("expected StyleTokenContractError from ProcessTemplate, got: %v", err)
		}
	})

	t.Run("invalid bold token value", func(t *testing.T) {
		inner := `<w:tbl>
<w:tblGrid><w:gridCol w:w="5000"/></w:tblGrid>
<w:tr><w:tc><w:tcPr><w:shd w:val="clear" w:color="auto" w:fill="{{cell_fill}}"/></w:tcPr><w:p><w:r><w:rPr><w:color w:val="{{row_color}}"/><w:b w:val="{{row_bold}}"/></w:rPr><w:t>Cell</w:t></w:r></w:p></w:tc></w:tr>
</w:tbl>`

		_, err := ProcessTemplate(createTestDocx(t, bodyDoc(inner)), map[string]any{
			"cell_fill": "123456",
			"row_color": "ABCDEF",
			"row_bold":  "TRUE",
		})
		if err == nil {
			t.Fatal("expected invalid bold style-token error")
		}

		var marker styleTokenContractErr
		if !errors.As(err, &marker) {
			t.Fatalf("expected StyleTokenContractError from ProcessTemplate, got: %v", err)
		}
	})

	t.Run("style values are not trimmed into validity", func(t *testing.T) {
		inner := `<w:p><w:r><w:rPr><w:color w:val="{{title_color}}"/></w:rPr><w:t>Title</w:t></w:r></w:p>`
		_, err := ProcessTemplate(createTestDocx(t, bodyDoc(inner)), map[string]any{
			"title_color": " ABCDEF ",
		})
		if err == nil {
			t.Fatal("expected exact-grammar style-token error")
		}

		var marker styleTokenContractErr
		if !errors.As(err, &marker) || !marker.IsStyleTokenContractError() {
			t.Fatalf("expected wrapped style-contract marker, got: %v", err)
		}
	})
}

func TestStyleAttributeEscapesInXMLLibrary(t *testing.T) {
	inner := `<w:tbl>
<w:tblGrid><w:gridCol w:w="5000"/></w:tblGrid>
<w:tr><w:tc><w:tcPr><w:shd w:val="clear" w:color="auto" w:fill="123456"/></w:tcPr><w:p><w:r><w:rPr><w:color w:val="ABCDEF"/><w:b w:val="true"/></w:rPr><w:t>Plain cell</w:t></w:r></w:p></w:tc></w:tr>
</w:tbl>`

	content := renderBody(t, inner, map[string]any{})
	mustContain(t, content, `w:shd w:val="clear" w:color="auto" w:fill="123456"`)
	mustContain(t, content, `w:color w:val="ABCDEF"`)
	if strings.Contains(content, "prefix{{") {
		t.Fatal("unexpected token text in escaped output")
	}
}
