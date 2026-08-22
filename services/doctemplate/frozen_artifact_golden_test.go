package doctemplate

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json/v2"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// ----------------------------------------------------------------------------
// Byte-stability goldens — the ACTUAL embedded frozen report-card artifacts
// ----------------------------------------------------------------------------
//
// Unlike TestByteStability_FrozenEmissions (synthetic invoice + handwritten
// report-card-shaped fragments, kept because they cost nothing), these goldens
// process the REAL template artifacts fayna embeds and renders live:
//
//	testdata/frozen-report-card-v1.docx    — vendored byte-identical copy of
//	  packages/fayna-golang/domain/operation/outcome_summary/document/
//	  report-card-template.docx (Template()/TemplateV1() embedded bytes)
//	  sha256 8d453d8f0905dbd6141c40abf2dd7439a6af82d4b83334e4d85a5d5b4403488c
//	testdata/frozen-report-card-v2.docx    — vendored copy of
//	  report-card-template-v2.docx (TemplateV2() embedded bytes)
//	  sha256 3f3f2cd7374ced58c76e4bf5bab7a693bc99fe40edad7ac6730f57e685a871cf
//	testdata/frozen-report-card-block.docx — vendored copy of
//	  outcome-summary-template-block.docx (TemplateBlock() embedded bytes)
//	  sha256 2ceb71e5cdc04122c2ea8e9a9ce3425bca235fdeee838c719c2915eb26340389
//	testdata/frozen-report-card-block.manifest.json — vendored copy of
//	  outcome-summary-template-block.manifest.json (ManifestBlock() bytes)
//	  sha256 6b809fc2791d1f30d78350680a10248e17de0905328bbe01f38af88d04fc72c9
//
// The artifacts are VENDORED (copied, with the source hashes above) because the
// dependency direction is fayna → fycha (fayna injects fycha's engine as a
// GenerateDoc closure; see fayna document/template.go package comment): fycha
// must never import fayna, so the test cannot read the embedded bytes via
// fayna's accessors. If fayna ever revs these artifacts, the frozen copies here
// intentionally do NOT change — they pin the ENGINE's behavior on these exact
// bytes, not fayna's current template inventory.
//
// Payloads are deterministic and fully-resolving so the processed parts carry
// zero residual {{...}} tokens (asserted below), which makes the goldens
// independent of the residual-token scrub (a no-op here) and directly
// comparable across the old and new engines:
//
//   - v1/v2: every value placeholder found in the artifact's parts resolves at
//     root scope to "V_<path>"; every loop key resolves to an EMPTY slice, so
//     loop markers and template rows are removed deterministically.
//   - block: the payload tree is built from the vendored manifest (the same
//     blank-guard contract fayna's builder uses), every scalar = "V_<path>",
//     every loop = exactly TWO items, recursively.
//
// PROVENANCE — the pinned hashes below were generated from the OLD (pre-Wave-6)
// engine, on 2026-07-19, without stashing (the Wave-6 diff stayed untouched in
// the working tree; the old engine was reconstructed from Git instead):
//
//  1. A throwaway module was created at
//     $SCRATCHPAD/head-harness/ with go.mod `module headharness`
//     (go 1.25.1, require github.com/beevik/etree v1.6.0 — same as fycha).
//  2. The four pre-Wave-6 engine sources were extracted from fycha HEAD
//     (commit 64081a525fa842dd3393e945da8b87a61eb5abc3, whose
//     services/doctemplate contains the OLD last-opener-break pairing, no
//     complete-stream scan, no residual scrub — verified by grepping the
//     extracted files for topLevelPairs/scrubResidualTemplateTokens: 0 hits):
//     git show HEAD:services/doctemplate/engine.go       > head-harness/doctemplate/engine.go
//     git show HEAD:services/doctemplate/xmlprocessor.go > head-harness/doctemplate/xmlprocessor.go
//     git show HEAD:services/doctemplate/docx.go         > head-harness/doctemplate/docx.go
//     git show HEAD:services/doctemplate/placeholder.go  > head-harness/doctemplate/placeholder.go
//  3. A main.go with payload builders and a part-hasher BYTE-EQUIVALENT to the
//     helpers in this file (frozenDeterministicValue/frozenSetDotPath/
//     frozenExtractTokens/frozenEmptyLoopPayload/frozenManifestPayload/
//     frozenHashProcessedParts) ran:
//     go run . <fycha>/services/doctemplate/testdata
//     Output (OLD engine):
//     goldenFrozenV1    = "01353c1301d4577580c4b3dbf355a29577a892f8e353062431ff23d403232444" // residual_tokens=0
//     goldenFrozenV2    = "e0356efed80095ecaaf080e5e29a3a22616c4df5730424c7b72111c11b104e7c" // residual_tokens=0
//     goldenFrozenBlock = "5a69272c77656615a17a41aa4a02f989c3c47fdf502492f4d91f719e955555fa" // residual_tokens=0
//  4. This test then ran against the NEW (Wave-6 + FIX-6) engine and produced
//     the identical three hashes — proving the refactor is byte-stable on the
//     real frozen artifacts. Any intentional emission change must update these
//     constants in the same commit, with a fresh provenance note.
const (
	goldenFrozenV1    = "01353c1301d4577580c4b3dbf355a29577a892f8e353062431ff23d403232444"
	goldenFrozenV2    = "e0356efed80095ecaaf080e5e29a3a22616c4df5730424c7b72111c11b104e7c"
	goldenFrozenBlock = "5a69272c77656615a17a41aa4a02f989c3c47fdf502492f4d91f719e955555fa"
)

func TestByteStability_FrozenArtifacts(t *testing.T) {
	cases := []struct {
		name     string
		file     string
		payload  func(t *testing.T, artifact []byte) map[string]any
		expected string
	}{
		{
			name: "v1",
			file: "testdata/frozen-report-card-v1.docx",
			payload: func(t *testing.T, artifact []byte) map[string]any {
				return frozenEmptyLoopPayload(t, artifact)
			},
			expected: goldenFrozenV1,
		},
		{
			name: "v2",
			file: "testdata/frozen-report-card-v2.docx",
			payload: func(t *testing.T, artifact []byte) map[string]any {
				return frozenEmptyLoopPayload(t, artifact)
			},
			expected: goldenFrozenV2,
		},
		{
			name: "block",
			file: "testdata/frozen-report-card-block.docx",
			payload: func(t *testing.T, _ []byte) map[string]any {
				manifest, err := os.ReadFile("testdata/frozen-report-card-block.manifest.json")
				if err != nil {
					t.Fatalf("read block manifest: %v", err)
				}
				return frozenManifestPayload(t, manifest)
			},
			expected: goldenFrozenBlock,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			artifact, err := os.ReadFile(tc.file)
			if err != nil {
				t.Fatalf("read frozen artifact: %v", err)
			}
			hash, residual := frozenHashProcessedParts(t, artifact, tc.payload(t, artifact))
			// The payloads fully resolve by construction; zero residual makes the
			// golden independent of the residual-token scrub AND doubles as a
			// zero-leak check on the real artifact.
			if residual != 0 {
				t.Errorf("processed parts carry %d residual {{ / }} sentinels; payload no longer fully resolves the artifact", residual)
			}
			if hash != tc.expected {
				t.Errorf("frozen %s artifact emission byte-stability regression:\n  want %s\n  got  %s",
					tc.name, tc.expected, hash)
			}
		})
	}
}

// ---- deterministic payload construction (byte-equivalent to the provenance
// ---- harness documented above; do not change without regenerating goldens) ----

var (
	// frozenTokenRe finds complete {{...}} constructs in tag-stripped part text.
	frozenTokenRe = regexp.MustCompile(`\{\{.*?\}\}`)
	// frozenTagRe strips XML tags so placeholders split across <w:r> runs re-join.
	frozenTagRe = regexp.MustCompile(`<[^>]+>`)
)

func frozenDeterministicValue(path string) string { return "V_" + path }

// frozenSetDotPath writes val at a dot-separated path, materializing nested maps.
func frozenSetDotPath(root map[string]any, path string, val any) {
	parts := strings.Split(path, ".")
	m := root
	for i := 0; i < len(parts)-1; i++ {
		next, ok := m[parts[i]].(map[string]any)
		if !ok {
			next = map[string]any{}
			m[parts[i]] = next
		}
		m = next
	}
	m[parts[len(parts)-1]] = val
}

// frozenExtractTokens inventories the artifact's value-placeholder paths and
// loop keys across document.xml, headers, and footers.
func frozenExtractTokens(t *testing.T, docx []byte) (valuePaths, loopKeys []string) {
	t.Helper()
	arch, err := ReadDocxBytes(docx)
	if err != nil {
		t.Fatalf("read artifact for token inventory: %v", err)
	}
	parts := []string{arch.Content}
	for _, h := range arch.Headers {
		parts = append(parts, h)
	}
	for _, f := range arch.Footers {
		parts = append(parts, f)
	}
	vset := map[string]bool{}
	lset := map[string]bool{}
	for _, p := range parts {
		text := frozenTagRe.ReplaceAllString(p, "")
		for _, m := range frozenTokenRe.FindAllString(text, -1) {
			body := strings.TrimSpace(m[2 : len(m)-2])
			switch {
			case strings.HasPrefix(body, "#"):
				lset[strings.TrimSpace(body[1:])] = true
			case strings.HasPrefix(body, "/"):
				// close marker: ignore
			default:
				vset[body] = true
			}
		}
	}
	for k := range vset {
		valuePaths = append(valuePaths, k)
	}
	for k := range lset {
		loopKeys = append(loopKeys, k)
	}
	sort.Strings(valuePaths)
	sort.Strings(loopKeys)
	return
}

// frozenEmptyLoopPayload resolves every value placeholder at root scope to
// "V_<path>" and every loop key to an EMPTY slice (markers + template rows are
// removed deterministically) — fully-resolving with zero residual by design.
func frozenEmptyLoopPayload(t *testing.T, docx []byte) map[string]any {
	t.Helper()
	valuePaths, loopKeys := frozenExtractTokens(t, docx)
	root := map[string]any{}
	for _, p := range valuePaths {
		frozenSetDotPath(root, p, frozenDeterministicValue(p))
	}
	for _, k := range loopKeys {
		frozenSetDotPath(root, k, []any{})
	}
	return root
}

// frozenManifestNode mirrors the blank-guard manifest schema fayna ships next to
// the block artifact (scalars + nested loops).
type frozenManifestNode struct {
	Scalars []string                      `json:"scalars"`
	Loops   map[string]frozenManifestNode `json:"loops"`
}

func frozenBuildNode(n frozenManifestNode) map[string]any {
	m := map[string]any{}
	for _, s := range n.Scalars {
		frozenSetDotPath(m, s, frozenDeterministicValue(s))
	}
	keys := make([]string, 0, len(n.Loops))
	for k := range n.Loops {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		sub := n.Loops[k]
		frozenSetDotPath(m, k, []any{frozenBuildNode(sub), frozenBuildNode(sub)})
	}
	return m
}

// frozenManifestPayload builds the block payload tree from the vendored
// manifest: every scalar = "V_<path>", every loop = exactly two items.
func frozenManifestPayload(t *testing.T, manifestJSON []byte) map[string]any {
	t.Helper()
	var root frozenManifestNode
	if err := json.Unmarshal(manifestJSON, &root); err != nil {
		t.Fatalf("parse block manifest: %v", err)
	}
	return frozenBuildNode(root)
}

// frozenHashProcessedParts processes the template and hashes ALL processed XML
// parts (document.xml + headers + footers, name-prefixed, headers/footers in
// sorted-name order) with SHA256. It also counts residual "{{"/"}}" sentinels
// across those parts.
func frozenHashProcessedParts(t *testing.T, templateBytes []byte, data map[string]any) (string, int) {
	t.Helper()
	out, err := ProcessTemplate(templateBytes, data)
	if err != nil {
		t.Fatalf("process frozen artifact: %v", err)
	}
	arch, err := ReadDocxBytes(out)
	if err != nil {
		t.Fatalf("read processed output: %v", err)
	}
	h := sha256.New()
	io.WriteString(h, "part:word/document.xml\n")
	io.WriteString(h, arch.Content)
	hnames := make([]string, 0, len(arch.Headers))
	for n := range arch.Headers {
		hnames = append(hnames, n)
	}
	sort.Strings(hnames)
	for _, n := range hnames {
		io.WriteString(h, "\npart:"+n+"\n")
		io.WriteString(h, arch.Headers[n])
	}
	fnames := make([]string, 0, len(arch.Footers))
	for n := range arch.Footers {
		fnames = append(fnames, n)
	}
	sort.Strings(fnames)
	for _, n := range fnames {
		io.WriteString(h, "\npart:"+n+"\n")
		io.WriteString(h, arch.Footers[n])
	}
	all := arch.Content
	for _, n := range hnames {
		all += arch.Headers[n]
	}
	for _, n := range fnames {
		all += arch.Footers[n]
	}
	residual := strings.Count(all, "{{") + strings.Count(all, "}}")
	return hex.EncodeToString(h.Sum(nil)), residual
}
