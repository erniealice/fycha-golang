# doctemplate — Pure Go DOCX Template Engine

A pure Go package for processing Microsoft Word (.docx) templates with JSON data. Replaces `{{placeholders}}` with values and expands `{{#loops}}` to duplicate table rows — producing ready-to-use documents like invoices, contracts, and reports.

Zero external service dependencies. No Node.js, no LibreOffice, no HTTP calls. Just Go + [etree](https://github.com/beevik/etree).

## Quick Start

```go
import "github.com/erniealice/fycha-golang/services/doctemplate"

// Read your Word template
template, _ := os.ReadFile("invoice-template.docx")

// Prepare your data
data := map[string]any{
    "client": map[string]any{
        "name":    "Acme Corporation",
        "address": "123 Business Ave, Manila",
    },
    "date":  "2026-03-08",
    "total": "₱18,500.00",
    "items": []any{
        map[string]any{"description": "Frontend Development", "amount": "₱5,000.00"},
        map[string]any{"description": "Backend Development",  "amount": "₱8,000.00"},
        map[string]any{"description": "Database Design",      "amount": "₱3,500.00"},
        map[string]any{"description": "Code Review",          "amount": "₱2,000.00"},
    },
}

// Process and save
result, err := doctemplate.ProcessTemplate(template, data)
if err != nil {
    log.Fatal(err)
}
os.WriteFile("invoice-output.docx", result, 0644)
```

## Features

| Feature | Status | Description |
|---------|--------|-------------|
| Simple placeholders | Done | `{{name}}` → "Alice" |
| Nested paths | Done | `{{client.name}}` → "Acme Corporation" |
| Cross-run accumulation | Done | Handles Word splitting `{{` across multiple XML runs |
| Table row loops | Done | `{{#items}}...{{/items}}` duplicates rows per array item |
| Non-loop row replacement | Done | Static rows (headers, totals) with `{{total}}` are processed |
| Header/footer processing | Done | Placeholders in document headers/footers are replaced |
| OOXML preservation | Done | All namespaces (w:, w14:, mc:, etc.) preserved on roundtrip |
| Body-level loops | Done | `{{#section}}...{{/section}}` for paragraph-level looping |
| Image replacement | Planned | — |

## Template Syntax

### Simple Placeholders

In your Word document, type `{{key}}` anywhere. The engine replaces it with the corresponding value from your data map.

```
Invoice for {{client.name}}
Address: {{client.address}}
Date: {{date}}
```

With data:
```json
{
  "client": { "name": "Acme Corp", "address": "123 Main St" },
  "date": "2026-03-08"
}
```

Produces:
```
Invoice for Acme Corp
Address: 123 Main St
Date: 2026-03-08
```

**Nested paths** use dot notation: `{{client.name}}` traverses `data["client"]["name"]`.

### Table Row Loops

For repeating table rows (like invoice line items), use three special rows:

| Row | Content | Purpose |
|-----|---------|---------|
| Marker row | `{{#items}}` in first cell | Marks start of loop — this row is removed from output |
| Template row | `{{description}}` \| `{{amount}}` | Cloned once per array item — placeholders resolved per item |
| Marker row | `{{/items}}` in first cell | Marks end of loop — this row is removed from output |

#### Example: Invoice Table

**Template in Word:**

| Description | Amount |
|-------------|--------|
| `{{#items}}` | |
| `{{description}}` | `{{amount}}` |
| `{{/items}}` | |
| **Total** | **`{{total}}`** |

**Data:**
```json
{
  "items": [
    { "description": "Frontend Development", "amount": "₱5,000.00" },
    { "description": "Backend Development",  "amount": "₱8,000.00" }
  ],
  "total": "₱13,000.00"
}
```

**Output:**

| Description | Amount |
|-------------|--------|
| Frontend Development | ₱5,000.00 |
| Backend Development | ₱8,000.00 |
| **Total** | **₱13,000.00** |

Key behaviors:
- **Marker rows are removed** — `{{#items}}` and `{{/items}}` rows don't appear in output
- **Template rows are cloned** — one copy per array item, with placeholders resolved from that item's data
- **Static rows are preserved** — the header row ("Description" / "Amount") and total row stay as-is
- **Non-loop placeholders in static rows work** — `{{total}}` in the total row is replaced from root data
- **No explicit mapping needed** — `{{description}}` inside the loop auto-resolves from the current array item

### Cross-Run Handling

Microsoft Word frequently splits text across multiple XML `<w:r>` (run) elements — especially when spell-check or formatting changes are involved. For example, `{{client.name}}` might be stored as:

```xml
<w:r><w:t>{{</w:t></w:r>
<w:proofErr w:type="spellStart"/>
<w:r><w:t>client.name</w:t></w:r>
<w:proofErr w:type="spellEnd"/>
<w:r><w:t>}}</w:t></w:r>
```

The engine handles this transparently by accumulating text across runs until a complete `{{...}}` placeholder is found. Template authors don't need to worry about this — just type the placeholder normally in Word.

## API Reference

### Core Function

```go
func ProcessTemplate(templateData []byte, data map[string]any) ([]byte, error)
```

The main entry point. Takes a DOCX file as bytes and a data map, returns the processed DOCX as bytes.

- **templateData**: Raw bytes of the .docx template file
- **data**: Key-value map for placeholder replacement. Supports nested maps and arrays.
- **Returns**: Processed .docx file as bytes, or an error

This function:
1. Reads the DOCX ZIP archive
2. Extracts `word/document.xml`, headers, and footers
3. Parses each XML part with etree (preserves all OOXML attributes/namespaces)
4. Processes paragraphs: cross-run text accumulation + placeholder replacement
5. Processes tables: detects `{{#key}}`/`{{/key}}` markers, clones rows per array item
6. Serializes XML back to string
7. Writes a new DOCX ZIP with the modified content

### DOCX Archive (Low-Level)

For advanced use cases where you need to inspect or manipulate the archive directly:

```go
// Read a DOCX file into its components
archive, err := doctemplate.ReadDocxBytes(data)

// Access the raw XML content
fmt.Println(archive.Content)  // word/document.xml as string
fmt.Println(archive.Headers)  // map[filename]xmlContent
fmt.Println(archive.Footers)  // map[filename]xmlContent

// Write back with modified content
result, err := archive.WriteDocx(
    modifiedContent,   // new document.xml content
    modifiedHeaders,   // map of modified headers
    modifiedFooters,   // map of modified footers
)
```

## Storage Integration (DocumentService)

For applications that read templates from and write results to cloud storage (GCS, S3, Azure, local filesystem), the parent `fycha` package provides `DocumentService`:

```go
import fycha "github.com/erniealice/fycha-golang"

// 1. Implement the StorageReadWriter interface (or use an adapter)
type StorageReadWriter interface {
    ReadObject(ctx context.Context, containerName, objectKey string) ([]byte, error)
    WriteObject(ctx context.Context, containerName, objectKey string, data []byte) error
}

// 2. Create the service
docService := fycha.NewDocumentService(myStorageAdapter)

// 3a. Process from storage → storage
err := docService.ProcessFromStorage(ctx,
    "templates", "invoice-template.docx",  // source
    "output",    "invoice-123.docx",       // destination
    data,
)

// 3b. Process from storage → bytes (for HTTP download)
result, err := docService.ProcessFromStorageToBytes(ctx,
    "templates", "invoice-template.docx",
    data,
)

// 3c. Process from bytes directly (no storage needed)
result, err := docService.ProcessBytes(templateBytes, data)
```

### Wiring with Espyna StorageAdapter

In your app's composition root (e.g., `container.go`), bridge the espyna `StorageAdapter` to fycha's `StorageReadWriter`:

```go
// storageReadWriter adapts espyna's StorageAdapter to fycha's interface
type storageReadWriter struct {
    adapter *consumer.StorageAdapter
}

func (s *storageReadWriter) ReadObject(ctx context.Context, container, key string) ([]byte, error) {
    return s.adapter.DownloadObject(ctx, container, key)
}

func (s *storageReadWriter) WriteObject(ctx context.Context, container, key string, data []byte) error {
    _, err := s.adapter.UploadObject(ctx, container, key, data)
    return err
}

// In container setup:
docService := fycha.NewDocumentService(&storageReadWriter{adapter: storageAdapter})
```

## Creating Templates

### In Microsoft Word

1. Open Word and create your document layout (invoice, contract, report, etc.)
2. Type placeholders using `{{double.curly.braces}}` syntax
3. For nested data, use dot notation: `{{client.name}}`, `{{client.address}}`
4. For repeating table rows:
   - Add a row with `{{#arrayName}}` in the first cell
   - Add your data row(s) with `{{field}}` placeholders
   - Add a row with `{{/arrayName}}` in the first cell
5. Save as .docx

### Tips for Template Authors

- **Don't worry about text splitting** — the engine handles Word's internal XML fragmentation
- **Keep placeholders simple** — use `{{name}}` not `{{data.items[0].name}}`
- **Loop markers go in their own rows** — don't mix `{{#items}}` with data in the same cell
- **Static rows work normally** — header rows, total rows, any row outside loop markers keeps its formatting and gets placeholder replacement
- **Formatting is preserved** — bold, italic, colors, cell shading all carry through
- **Test with 2+ items** — to verify row duplication works correctly

### Data Format

The data map uses Go's `map[string]any` type. When consuming from JSON:

```go
var data map[string]any
json.Unmarshal(jsonBytes, &data)

// json.Unmarshal produces map[string]any for objects
// and []any for arrays — exactly what ProcessTemplate expects
```

Supported value types:
- `string` — used directly
- `int`, `float64`, etc. — converted via `fmt.Sprintf("%v", val)`
- `map[string]any` — accessed via dot notation (`{{key.subkey}}`)
- `[]any` of `map[string]any` — used for table row loops (`{{#key}}...{{/key}}`)

## Architecture

```
┌─────────────────────────────────────────────┐
│  Consumer App (retail-admin, service-admin)  │
│  ┌───────────────────────────────────────┐  │
│  │ container.go (composition root)       │  │
│  │ - wires StorageAdapter → fycha iface  │  │
│  └───────────────────┬───────────────────┘  │
│                      │                       │
├──────────────────────┼───────────────────────┤
│  fycha-golang-ryta   │                       │
│  ┌───────────────────┴───────────────────┐  │
│  │ DocumentService                        │  │
│  │ - ProcessFromStorage (storage I/O)     │  │
│  │ - ProcessFromStorageToBytes            │  │
│  │ - ProcessBytes (no I/O)               │  │
│  │                                        │  │
│  │ StorageReadWriter interface            │  │
│  └───────────────────┬───────────────────┘  │
│                      │                       │
│  ┌───────────────────┴───────────────────┐  │
│  │ services/doctemplate (pure library)    │  │
│  │ - ProcessTemplate([]byte, data)       │  │
│  │ - ReadDocxBytes / WriteDocx           │  │
│  │ - XML processor (etree-based)         │  │
│  │ - Zero I/O, zero storage deps        │  │
│  └───────────────────────────────────────┘  │
├─────────────────────────────────────────────┤
│  espyna-golang-ryta                          │
│  - StorageAdapter (GCS, S3, Azure, Mock)    │
└─────────────────────────────────────────────┘
```

**Layer 1 — `services/doctemplate/`**: Pure library. Accepts `[]byte`, returns `[]byte`. No I/O, no context, no storage. Fully reusable and independently testable.

**Layer 2 — `fycha/DocumentService`**: Orchestrates storage I/O + template processing. Defines `StorageReadWriter` interface to stay provider-agnostic.

**Layer 3 — App composition root**: Bridges espyna's StorageAdapter to fycha's interface.

## File Structure

```
packages/fycha-golang-ryta/
├── document_service.go          # DocumentService + StorageReadWriter interface
├── services/
│   └── doctemplate/
│       ├── engine.go            # ProcessTemplate — public API entry point
│       ├── docx.go              # DOCX ZIP read/write (ReadDocxBytes, WriteDocx)
│       ├── placeholder.go       # Regex-based {{key.path}} replacement
│       ├── xmlprocessor.go      # Cross-run accumulation, table loops, body loops
│       ├── engine_test.go       # Tests (5 passing)
│       ├── README.md            # This file
│       └── testdata/
│           ├── gen_template.go  # Generates test DOCX templates
│           ├── invoice-template.docx
│           └── invoice-output.docx
```

## Running Tests

```bash
cd packages/fycha-golang-ryta
go test ./services/doctemplate/ -v
```

Expected output:
```
=== RUN   TestProcessTemplate_NestedJSON
--- PASS: TestProcessTemplate_NestedJSON
=== RUN   TestGetNestedValue
--- PASS: TestGetNestedValue
=== RUN   TestGetReplaceValue
--- PASS: TestGetReplaceValue
=== RUN   TestProcessTemplate_InvoiceTableLoop
--- PASS: TestProcessTemplate_InvoiceTableLoop
=== RUN   TestReplacePlaceholders
--- PASS: TestReplacePlaceholders
PASS
```

## Dependencies

- [github.com/beevik/etree](https://github.com/beevik/etree) — DOM-style XML manipulation that preserves OOXML namespaces, attributes, and element order during roundtrip. Single package, no transitive dependencies.
