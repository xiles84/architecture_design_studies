package evidence

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func loadV2Fixture(t *testing.T, root string) (*V2Document, any, *V2Confounds, *Document) {
	t.Helper()
	doc, raw, err := LoadV2(filepath.Join(root, "book", "evidence", "v2", "claims.json"))
	if err != nil {
		t.Fatal(err)
	}
	conf, err := LoadConfounds(filepath.Join(root, "book", "evidence", "v2", "confounds.json"))
	if err != nil {
		t.Fatal(err)
	}
	v1, _, err := Load(filepath.Join(root, "book", "evidence", "claims.json"))
	if err != nil {
		t.Fatal(err)
	}
	return doc, raw, conf, v1
}

// TestV2RegistryValidates is the positive check: the committed correction
// package must resolve every structured support key and anchor.
func TestV2RegistryValidates(t *testing.T) {
	root := repoRoot(t)
	doc, raw, conf, v1 := loadV2Fixture(t, root)
	var schema map[string]any
	data, err := os.ReadFile(filepath.Join(root, "book", "evidence", "v2", "schema.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &schema); err != nil {
		t.Fatal(err)
	}
	if errs := ValidateV2(root, schema, doc, raw, conf, v1); len(errs) > 0 {
		for _, e := range errs {
			t.Errorf("unexpected error: %v", e)
		}
	}
	if len(doc.Claims) < 20 {
		t.Fatalf("expected the correction package to hold at least 20 claims, got %d", len(doc.Claims))
	}
}

// TestV2UnresolvedSupportKeyFails is the semantic resolver's negative control:
// a support key whose token is not in the cited report must be rejected. This
// is exactly the v1 failure mode (a curated cell name that resolves nowhere).
func TestV2UnresolvedSupportKeyFails(t *testing.T) {
	root := repoRoot(t)
	doc, raw, conf, v1 := loadV2Fixture(t, root)
	doc.Claims[0].Support[0].ResolvesTo = "this-token-does-not-exist-in-any-report"
	errs := ValidateV2(root, map[string]any{}, doc, raw, conf, v1)
	if !hasSubstring(errs, "does not appear in") {
		t.Fatalf("expected an unresolved-support-key error, got %v", errs)
	}
}

// TestV2LegacyFlagMismatchFails: a claim that claims complete provenance while
// its run has no tag must be rejected.
func TestV2LegacyFlagMismatchFails(t *testing.T) {
	root := repoRoot(t)
	doc, raw, conf, v1 := loadV2Fixture(t, root)
	for _, c := range doc.Claims {
		if c.LegacyProvenanceIncomplete {
			c.LegacyProvenanceIncomplete = false
			break
		}
	}
	errs := ValidateV2(root, map[string]any{}, doc, raw, conf, v1)
	if !hasSubstring(errs, "legacy_provenance_incomplete") {
		t.Fatalf("expected a legacy-provenance error, got %v", errs)
	}
}

// TestV2UnreferencedConfoundFails: a registered confound nobody cites is dead.
func TestV2UnreferencedConfoundFails(t *testing.T) {
	root := repoRoot(t)
	doc, raw, conf, v1 := loadV2Fixture(t, root)
	for _, c := range doc.Claims {
		c.Confounds = nil
	}
	errs := ValidateV2(root, map[string]any{}, doc, raw, conf, v1)
	if !hasSubstring(errs, "not referenced by any claim") {
		t.Fatalf("expected an unreferenced-confound error, got %v", errs)
	}
}

// TestV2UnknownSupersedesFails: a v2 claim may only supersede a v1 claim that
// exists, so the correction cannot invent provenance.
func TestV2UnknownSupersedesFails(t *testing.T) {
	root := repoRoot(t)
	doc, raw, conf, v1 := loadV2Fixture(t, root)
	unknown := "claim-does-not-exist"
	doc.Claims[0].Supersedes = &unknown
	errs := ValidateV2(root, map[string]any{}, doc, raw, conf, v1)
	if !hasSubstring(errs, "not in the v1 registry") {
		t.Fatalf("expected an unknown-supersedes error, got %v", errs)
	}
}

// TestV1DiagnosticRegression keeps the mechanical v1 audit the brainstorm
// conclusion computed: 11 of 12 numeric claims have a name that resolves
// nowhere, and 36 of 43 names are unresolved. A green validator on v1 would
// have hidden this; the fixture makes the failure mode reproducible.
func TestV1DiagnosticRegression(t *testing.T) {
	root := repoRoot(t)
	v1, _, err := Load(filepath.Join(root, "book", "evidence", "claims.json"))
	if err != nil {
		t.Fatal(err)
	}
	rep, err := V1Diagnostic(root, v1)
	if err != nil {
		t.Fatal(err)
	}
	if rep.NumericClaims != 12 || rep.ClaimsWithMissing != 11 || rep.Names != 43 || rep.NamesMissing != 36 {
		t.Fatalf("v1 diagnostic drifted: claims=%d missing=%d names=%d unresolved=%d (want 12/11/43/36)",
			rep.NumericClaims, rep.ClaimsWithMissing, rep.Names, rep.NamesMissing)
	}
}

func hasSubstring(errs []error, sub string) bool {
	for _, e := range errs {
		if strings.Contains(e.Error(), sub) {
			return true
		}
	}
	return false
}
