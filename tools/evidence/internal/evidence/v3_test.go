package evidence

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func loadV3Fixture(t *testing.T, root string) (*V2Document, any, *V2Confounds, *Document, *V2Document) {
	t.Helper()
	doc, raw, err := LoadV2(filepath.Join(root, "book", "evidence", "v3", "claims.json"))
	if err != nil {
		t.Fatal(err)
	}
	conf, err := LoadConfounds(filepath.Join(root, "book", "evidence", "v3", "confounds.json"))
	if err != nil {
		t.Fatal(err)
	}
	v1, _, err := Load(filepath.Join(root, "book", "evidence", "claims.json"))
	if err != nil {
		t.Fatal(err)
	}
	v2, _, err := LoadV2(filepath.Join(root, "book", "evidence", "v2", "claims.json"))
	if err != nil {
		t.Fatal(err)
	}
	return doc, raw, conf, v1, v2
}

func loadV3Schema(t *testing.T, root string) map[string]any {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, "book", "evidence", "v3", "schema.json"))
	if err != nil {
		t.Fatal(err)
	}
	var schema map[string]any
	if err := json.Unmarshal(data, &schema); err != nil {
		t.Fatal(err)
	}
	return schema
}

// TestV3RegistryValidates is the positive check: the committed v3 package must
// resolve every structured support key, every anchor and every supersession
// against v1 or v2.
func TestV3RegistryValidates(t *testing.T) {
	root := repoRoot(t)
	doc, raw, conf, v1, v2 := loadV3Fixture(t, root)
	schema := loadV3Schema(t, root)
	if errs := ValidateV3(root, schema, doc, raw, conf, v1, v2); len(errs) > 0 {
		for _, e := range errs {
			t.Errorf("unexpected error: %v", e)
		}
	}
	if len(doc.Claims) < 29 {
		t.Fatalf("expected the v3 package to hold at least 29 claims, got %d", len(doc.Claims))
	}
	if doc.SchemaVersion != 3 {
		t.Fatalf("expected schema_version 3, got %d", doc.SchemaVersion)
	}
}

// TestV3SupersedesV2Resolves is the new closure rule: a v3 claim may supersede a
// v2 claim, and a v3 retirement may name a v2 claim. This is what v2 could not
// do and what makes the versioned update possible without editing v2.
func TestV3SupersedesV2Resolves(t *testing.T) {
	root := repoRoot(t)
	doc, raw, conf, v1, v2 := loadV3Fixture(t, root)
	schema := loadV3Schema(t, root)

	found := false
	for _, c := range doc.Claims {
		if c.ClaimID == "v3-gap-01-study05-remaining-regimes" {
			if c.Supersedes == nil || *c.Supersedes != "v2-gap-04-study05-unmeasured-regimes" {
				t.Fatalf("v3-gap-01 must supersede v2-gap-04, got %v", c.Supersedes)
			}
			found = true
		}
	}
	if !found {
		t.Fatal("v3-gap-01-study05-remaining-regimes is missing from v3")
	}

	unknown := "v2-does-not-exist"
	doc.Claims[0].Supersedes = &unknown
	errs := ValidateV3(root, schema, doc, raw, conf, v1, v2)
	if !hasSubstring(errs, "not in the predecessor registry") {
		t.Fatalf("expected an unknown-supersedes error, got %v", errs)
	}
}

// TestV3UnresolvedSupportKeyFails: the v3 resolver keeps the semantic gate that
// v1 lacked, so a claim cannot be added with a cell name that resolves nowhere.
func TestV3UnresolvedSupportKeyFails(t *testing.T) {
	root := repoRoot(t)
	doc, raw, conf, v1, v2 := loadV3Fixture(t, root)
	schema := loadV3Schema(t, root)
	for _, c := range doc.Claims {
		if c.ClaimID == "v3-01-churn-crosses-real-ttl" {
			c.Support[0].ResolvesTo = "this-token-does-not-exist-in-any-report"
		}
	}
	errs := ValidateV3(root, schema, doc, raw, conf, v1, v2)
	if !hasSubstring(errs, "does not appear in") {
		t.Fatalf("expected an unresolved-support-key error, got %v", errs)
	}
}
