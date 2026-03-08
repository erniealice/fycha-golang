package doctemplate

import (
	"archive/zip"
	"bytes"
	"os"
	"strings"
	"testing"
)

// testData returns the JSON data that matches the test template placeholders.
func testData() map[string]any {
	return map[string]any{
		"client": map[string]any{
			"name":    "Acme Corporation",
			"address": "123 Business Ave",
			"country": "Philippines",
			"skills":  "Web Development",
		},
		"developer": map[string]any{
			"address":        "456 Dev Street",
			"city_state_zip": "Manila, NCR 1000",
			"email":          "dev@example.com",
			"phone":          "+63 912 345 6789",
			"website":        "https://example.com",
		},
	}
}

func TestProcessTemplate_NestedJSON(t *testing.T) {
	// Read the test fixture
	templatePath := "../../../../references/zazzy-golang-v1/docs/testreplace-1.docx"
	templateData, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("failed to read test template: %v", err)
	}

	data := testData()

	// Process the template
	result, err := ProcessTemplate(templateData, data)
	if err != nil {
		t.Fatalf("ProcessTemplate failed: %v", err)
	}

	// Verify output is a valid DOCX (ZIP)
	if len(result) == 0 {
		t.Fatal("output is empty")
	}

	// Read back the output to inspect content
	archive, err := ReadDocxBytes(result)
	if err != nil {
		t.Fatalf("failed to read output docx: %v", err)
	}

	content := archive.Content

	// Verify nested JSON replacements happened
	checks := []struct {
		placeholder string
		expected    string
	}{
		{"{{client.name}}", "Acme Corporation"},
		{"{{client.address}}", "123 Business Ave"},
		{"{{client.country}}", "Philippines"},
		{"{{developer.address}}", "456 Dev Street"},
		{"{{developer.city_state_zip}}", "Manila, NCR 1000"},
		{"{{developer.email}}", "dev@example.com"},
		{"{{developer.phone}}", "+63 912 345 6789"},
		{"{{developer.website}}", "https://example.com"},
	}

	for _, check := range checks {
		// The placeholder should NOT appear in the output
		if strings.Contains(content, check.placeholder) {
			t.Errorf("placeholder %s was not replaced", check.placeholder)
		}
		// The expected value SHOULD appear in the output
		if !strings.Contains(content, check.expected) {
			t.Errorf("expected value %q not found in output for placeholder %s", check.expected, check.placeholder)
		}
	}

	// Write output for manual inspection
	outputPath := "../../../../references/zazzy-golang-v1/docs/testreplace-1-go-output.docx"
	if err := os.WriteFile(outputPath, result, 0644); err != nil {
		t.Logf("warning: could not write output file: %v", err)
	} else {
		t.Logf("output written to %s", outputPath)
	}
}

func TestGetNestedValue(t *testing.T) {
	data := map[string]any{
		"client": map[string]any{
			"name": "Acme",
			"address": map[string]any{
				"city": "Manila",
			},
		},
		"simple": "hello",
	}

	tests := []struct {
		keys     []string
		expected any
		found    bool
	}{
		{[]string{"simple"}, "hello", true},
		{[]string{"client", "name"}, "Acme", true},
		{[]string{"client", "address", "city"}, "Manila", true},
		{[]string{"missing"}, nil, false},
		{[]string{"client", "missing"}, nil, false},
	}

	for _, tt := range tests {
		val, found := getNestedValue(tt.keys, data)
		if found != tt.found {
			t.Errorf("getNestedValue(%v): found=%v, want %v", tt.keys, found, tt.found)
		}
		if found && val != tt.expected {
			t.Errorf("getNestedValue(%v): got %v, want %v", tt.keys, val, tt.expected)
		}
	}
}

func TestGetReplaceValue(t *testing.T) {
	data := map[string]any{
		"client": map[string]any{
			"name": "Acme",
		},
		"simple": "hello",
	}

	tests := []struct {
		path     string
		expected string
		found    bool
	}{
		{"simple", "hello", true},
		{"client.name", "Acme", true},
		{"missing", "", false},
	}

	for _, tt := range tests {
		val, found := getReplaceValue(data, tt.path)
		if found != tt.found {
			t.Errorf("getReplaceValue(%q): found=%v, want %v", tt.path, found, tt.found)
		}
		if val != tt.expected {
			t.Errorf("getReplaceValue(%q): got %q, want %q", tt.path, val, tt.expected)
		}
	}
}

func TestProcessTemplate_InvoiceTableLoop(t *testing.T) {
	templateData, err := os.ReadFile("testdata/invoice-template.docx")
	if err != nil {
		t.Fatalf("failed to read invoice template: %v", err)
	}

	data := map[string]any{
		"client": map[string]any{
			"name":    "Acme Corporation",
			"address": "123 Business Ave, Manila",
		},
		"date":  "2026-03-08",
		"notes": "Payment due within 30 days",
		"total": "₱18,500.00",
		"items": []any{
			map[string]any{"description": "Frontend Development", "amount": "₱5,000.00"},
			map[string]any{"description": "Backend Development", "amount": "₱8,000.00"},
			map[string]any{"description": "Database Design", "amount": "₱3,500.00"},
			map[string]any{"description": "Code Review", "amount": "₱2,000.00"},
		},
	}

	result, err := ProcessTemplate(templateData, data)
	if err != nil {
		t.Fatalf("ProcessTemplate failed: %v", err)
	}

	archive, err := ReadDocxBytes(result)
	if err != nil {
		t.Fatalf("failed to read output docx: %v", err)
	}

	content := archive.Content

	// 1. Verify simple placeholder replacements
	simpleChecks := []struct {
		placeholder string
		expected    string
	}{
		{"{{client.name}}", "Acme Corporation"},
		{"{{client.address}}", "123 Business Ave, Manila"},
		{"{{date}}", "2026-03-08"},
		{"{{notes}}", "Payment due within 30 days"},
		{"{{total}}", "₱18,500.00"},
	}

	for _, check := range simpleChecks {
		if !strings.Contains(content, check.expected) {
			t.Errorf("expected value %q not found for %s", check.expected, check.placeholder)
		}
	}

	// 2. Verify loop markers are stripped
	loopMarkers := []string{"{{#items}}", "{{/items}}"}
	for _, marker := range loopMarkers {
		if strings.Contains(content, marker) {
			t.Errorf("loop marker %s should be stripped from output", marker)
		}
	}

	// 3. Verify all 4 item descriptions appear (table row duplication)
	itemChecks := []string{
		"Frontend Development",
		"Backend Development",
		"Database Design",
		"Code Review",
	}
	for _, item := range itemChecks {
		if !strings.Contains(content, item) {
			t.Errorf("item %q not found in output — table row duplication may have failed", item)
		}
	}

	// 4. Verify all 4 amounts appear
	amountChecks := []string{
		"₱5,000.00",
		"₱8,000.00",
		"₱3,500.00",
		"₱2,000.00",
	}
	for _, amount := range amountChecks {
		if !strings.Contains(content, amount) {
			t.Errorf("amount %q not found in output", amount)
		}
	}

	// 5. Verify static header row preserved
	if !strings.Contains(content, "Description") {
		t.Error("table header 'Description' missing from output")
	}

	// 6. Count <w:tr> elements — expect: 1 header + 4 data rows + 1 total = 6
	trCount := strings.Count(content, "<w:tr")
	t.Logf("table rows (<w:tr>) in output: %d", trCount)
	if trCount < 6 {
		t.Errorf("expected at least 6 <w:tr> elements (1 header + 4 data + 1 total), got %d", trCount)
	}

	// Write output for manual inspection
	outputPath := "testdata/invoice-output.docx"
	if err := os.WriteFile(outputPath, result, 0644); err != nil {
		t.Logf("warning: could not write output: %v", err)
	} else {
		t.Logf("output written to %s — open in Word/LibreOffice to verify", outputPath)
	}
}

// createTestDocx builds a minimal DOCX in memory from raw document.xml content.
func createTestDocx(t *testing.T, documentXML string) []byte {
	t.Helper()

	contentTypesXML := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/><Default Extension="xml" ContentType="application/xml"/><Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/></Types>`

	relsXML := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/></Relationships>`

	documentRelsXML := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"></Relationships>`

	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	files := map[string]string{
		"[Content_Types].xml":          contentTypesXML,
		"_rels/.rels":                  relsXML,
		"word/document.xml":            documentXML,
		"word/_rels/document.xml.rels": documentRelsXML,
	}

	for name, content := range files {
		writer, err := w.Create(name)
		if err != nil {
			t.Fatalf("failed to create zip entry %s: %v", name, err)
		}
		if _, err := writer.Write([]byte(content)); err != nil {
			t.Fatalf("failed to write zip entry %s: %v", name, err)
		}
	}

	if err := w.Close(); err != nil {
		t.Fatalf("failed to close zip: %v", err)
	}

	return buf.Bytes()
}

func TestProcessTemplate_BodyLevelLoop(t *testing.T) {
	// Template: {{#sections}} paragraph, {{title}} paragraph, {{/sections}} paragraph
	documentXML := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
<w:body>
<w:p><w:r><w:t>Header Text</w:t></w:r></w:p>
<w:p><w:r><w:t>{{#sections}}</w:t></w:r></w:p>
<w:p><w:r><w:t>Title: {{title}}</w:t></w:r></w:p>
<w:p><w:r><w:t>{{/sections}}</w:t></w:r></w:p>
<w:p><w:r><w:t>Footer Text</w:t></w:r></w:p>
</w:body>
</w:document>`

	templateData := createTestDocx(t, documentXML)

	data := map[string]any{
		"sections": []any{
			map[string]any{"title": "Section 1"},
			map[string]any{"title": "Section 2"},
			map[string]any{"title": "Section 3"},
		},
	}

	result, err := ProcessTemplate(templateData, data)
	if err != nil {
		t.Fatalf("ProcessTemplate failed: %v", err)
	}

	archive, err := ReadDocxBytes(result)
	if err != nil {
		t.Fatalf("failed to read output docx: %v", err)
	}

	content := archive.Content

	// 1. All 3 section titles must appear
	for _, title := range []string{"Section 1", "Section 2", "Section 3"} {
		if !strings.Contains(content, title) {
			t.Errorf("expected %q in output, not found", title)
		}
	}

	// 2. Loop markers must be stripped
	for _, marker := range []string{"{{#sections}}", "{{/sections}}"} {
		if strings.Contains(content, marker) {
			t.Errorf("loop marker %s should be stripped from output", marker)
		}
	}

	// 3. Static content before and after the loop must be preserved
	if !strings.Contains(content, "Header Text") {
		t.Error("static header text before loop is missing")
	}
	if !strings.Contains(content, "Footer Text") {
		t.Error("static footer text after loop is missing")
	}

	// 4. Count occurrences of "Title:" — should be exactly 3
	titleCount := strings.Count(content, "Title:")
	if titleCount != 3 {
		t.Errorf("expected 3 'Title:' occurrences, got %d", titleCount)
	}

	t.Logf("output content length: %d bytes", len(content))
}

func TestProcessTemplate_TwoLevelNesting(t *testing.T) {
	// Template: body-level {{#invoices}} wrapping a {{client_name}} paragraph
	// and a table with {{#items}}/{{/items}} row loop
	documentXML := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
<w:body>
<w:p><w:r><w:t>Multi-Invoice Report</w:t></w:r></w:p>
<w:p><w:r><w:t>{{#invoices}}</w:t></w:r></w:p>
<w:p><w:r><w:t>Client: {{client_name}}</w:t></w:r></w:p>
<w:tbl>
<w:tblPr><w:tblStyle w:val="TableGrid"/></w:tblPr>
<w:tblGrid><w:gridCol w:w="5000"/><w:gridCol w:w="3000"/></w:tblGrid>
<w:tr><w:tc><w:p><w:r><w:t>Description</w:t></w:r></w:p></w:tc><w:tc><w:p><w:r><w:t>Amount</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>{{#items}}</w:t></w:r></w:p></w:tc><w:tc><w:p><w:r><w:t></w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>{{description}}</w:t></w:r></w:p></w:tc><w:tc><w:p><w:r><w:t>{{amount}}</w:t></w:r></w:p></w:tc></w:tr>
<w:tr><w:tc><w:p><w:r><w:t>{{/items}}</w:t></w:r></w:p></w:tc><w:tc><w:p><w:r><w:t></w:t></w:r></w:p></w:tc></w:tr>
</w:tbl>
<w:p><w:r><w:t>{{/invoices}}</w:t></w:r></w:p>
<w:p><w:r><w:t>End of Report</w:t></w:r></w:p>
</w:body>
</w:document>`

	templateData := createTestDocx(t, documentXML)

	data := map[string]any{
		"invoices": []any{
			map[string]any{
				"client_name": "Acme Corp",
				"items": []any{
					map[string]any{"description": "Dev", "amount": "5000"},
				},
			},
			map[string]any{
				"client_name": "Beta Corp",
				"items": []any{
					map[string]any{"description": "Design", "amount": "3000"},
					map[string]any{"description": "QA", "amount": "2000"},
				},
			},
		},
	}

	result, err := ProcessTemplate(templateData, data)
	if err != nil {
		t.Fatalf("ProcessTemplate failed: %v", err)
	}

	archive, err := ReadDocxBytes(result)
	if err != nil {
		t.Fatalf("failed to read output docx: %v", err)
	}

	content := archive.Content

	// 1. Both client names must appear
	if !strings.Contains(content, "Acme Corp") {
		t.Error("expected 'Acme Corp' in output")
	}
	if !strings.Contains(content, "Beta Corp") {
		t.Error("expected 'Beta Corp' in output")
	}

	// 2. All item descriptions must appear
	for _, desc := range []string{"Dev", "Design", "QA"} {
		if !strings.Contains(content, desc) {
			t.Errorf("expected item description %q in output", desc)
		}
	}

	// 3. All amounts must appear
	for _, amount := range []string{"5000", "3000", "2000"} {
		if !strings.Contains(content, amount) {
			t.Errorf("expected amount %q in output", amount)
		}
	}

	// 4. Loop markers must be stripped
	for _, marker := range []string{"{{#invoices}}", "{{/invoices}}", "{{#items}}", "{{/items}}"} {
		if strings.Contains(content, marker) {
			t.Errorf("loop marker %s should be stripped from output", marker)
		}
	}

	// 5. Static content must be preserved
	if !strings.Contains(content, "Multi-Invoice Report") {
		t.Error("report header missing")
	}
	if !strings.Contains(content, "End of Report") {
		t.Error("report footer missing")
	}

	// 6. Count table rows: 2 tables expected
	//    Table 1 (Acme): 1 header + 1 data row = 2 rows
	//    Table 2 (Beta): 1 header + 2 data rows = 3 rows
	//    Total: 5 <w:tr> elements
	trCount := strings.Count(content, "<w:tr")
	t.Logf("table rows (<w:tr>) in output: %d", trCount)
	if trCount < 5 {
		t.Errorf("expected at least 5 <w:tr> elements (2 tables: 1+1 header + 1+2 data), got %d", trCount)
	}

	t.Logf("output content length: %d bytes", len(content))
}

func TestReplacePlaceholders(t *testing.T) {
	data := map[string]any{
		"name":  "Alice",
		"age":   "30",
		"items": []any{}, // should be skipped by placeholder replacer
	}

	input := "Hello {{name}}, age {{age}}. Loop: {{#items}}{{/items}}"
	result := ReplacePlaceholders(input, data)

	if !strings.Contains(result, "Hello Alice") {
		t.Errorf("expected 'Hello Alice', got: %s", result)
	}
	if !strings.Contains(result, "age 30") {
		t.Errorf("expected 'age 30', got: %s", result)
	}
	// Loop markers should be preserved
	if !strings.Contains(result, "{{#items}}") {
		t.Errorf("loop start marker should be preserved, got: %s", result)
	}
	if !strings.Contains(result, "{{/items}}") {
		t.Errorf("loop end marker should be preserved, got: %s", result)
	}
}
