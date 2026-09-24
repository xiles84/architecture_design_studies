package evidence

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func loadV5Fixture(t *testing.T, root string) (*V2Document, any, *V2Confounds, *Document, *V2Document, *V2Document, *V2Document) {
	t.Helper()
	doc, raw, err := LoadV2(filepath.Join(root, "book", "evidence", "v5", "claims.json"))
	if err != nil {
		t.Fatal(err)
	}
	conf, err := LoadConfounds(filepath.Join(root, "book", "evidence", "v5", "confounds.json"))
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
	v3, _, err := LoadV2(filepath.Join(root, "book", "evidence", "v3", "claims.json"))
	if err != nil {
		t.Fatal(err)
	}
	v4, _, err := LoadV2(filepath.Join(root, "book", "evidence", "v4", "claims.json"))
	if err != nil {
		t.Fatal(err)
	}
	return doc, raw, conf, v1, v2, v3, v4
}

func loadV5Schema(t *testing.T, root string) map[string]any {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, "book", "evidence", "v5", "schema.json"))
	if err != nil {
		t.Fatal(err)
	}
	var schema map[string]any
	if err := json.Unmarshal(data, &schema); err != nil {
		t.Fatal(err)
	}
	return schema
}

func claimByID(t *testing.T, doc *V2Document, id string) *V2Claim {
	t.Helper()
	for i := range doc.Claims {
		if doc.Claims[i].ClaimID == id {
			return doc.Claims[i]
		}
	}
	t.Fatalf("claim %s not found", id)
	return nil
}

// TestV5RegistryValidates is the positive check: the committed v5 package must
// resolve every structured support key, every anchor, every supersession and every
// retirement against v1, v2, v3 or v4, satisfy the v4 lifecycle rules, and pass the
// new repetition-versus-trials lint.
func TestV5RegistryValidates(t *testing.T) {
	root := repoRoot(t)
	doc, raw, conf, v1, v2, v3, v4 := loadV5Fixture(t, root)
	schema := loadV5Schema(t, root)
	if errs := ValidateV5(root, schema, doc, raw, conf, v1, v2, v3, v4); len(errs) > 0 {
		for _, e := range errs {
			t.Errorf("unexpected error: %v", e)
		}
	}
	if len(doc.Claims) != 29 {
		t.Fatalf("want 29 claims, got %d", len(doc.Claims))
	}
}

// TestV5RepetitionLintCatchesTheV207Shape is the negative control for the new
// lint: restoring v2-07's old limit, which cited three runs beside Trials: 1 and a
// single anchor, must fail. This is the defect the third review round found.
func TestV5RepetitionLintCatchesTheV207Shape(t *testing.T) {
	root := repoRoot(t)
	doc, _, _, _, _, _, _ := loadV5Fixture(t, root)
	c := claimByID(t, doc, "v2-07-refund-answerability-and-storage")
	if c.Trials.WholeRunReplications != 1 {
		t.Fatalf("fixture drift: v2-07 should still declare one whole-run replication")
	}
	for i, l := range c.Limits {
		if strings.Contains(l, "relation-size difference measured in the primary run") {
			c.Limits[i] = "the +33-35% storage delta is deterministic and stable across three runs"
		}
	}
	errs := validateRepetition(doc)
	if len(errs) == 0 {
		t.Fatal("expected the lint to fail a limit citing three runs beside Trials: 1")
	}
	if !strings.Contains(errs[0].Error(), "three runs") {
		t.Fatalf("error should name the offending phrase, got: %v", errs[0])
	}
}

// TestV5RepetitionLintAllowsACorroboratingAnchor is the other half of the rule:
// repetition is legitimate when the claim can show it. The same sentence passes
// once a corroborating anchor is present.
func TestV5RepetitionLintAllowsACorroboratingAnchor(t *testing.T) {
	root := repoRoot(t)
	doc, _, _, _, _, _, _ := loadV5Fixture(t, root)
	c := claimByID(t, doc, "v2-07-refund-answerability-and-storage")
	c.Limits = append(c.Limits, "corroborated by two sibling runs")
	c.Anchors = append(c.Anchors, V2Anchor{RunID: "20260913T021010Z", Role: "corroborating", InputsDigest: "fixture"})
	if errs := validateRepetition(doc); len(errs) > 0 {
		// Only the appended sentence is in play here; the claim's own limits carry
		// no other repetition. Any error means the escape hatch is not working.
		for _, e := range errs {
			if strings.Contains(e.Error(), "two sibling runs") {
				t.Fatalf("a corroborating anchor should permit the repetition: %v", e)
			}
		}
	}
}

// TestV5CarriesV4ClaimsForwardExceptTheFive pins the sweep's result: exactly five
// claims changed, and every other claim is identical to v4. A future pass that
// re-scopes a sixth claim without touching this test and the ledger should fail.
func TestV5CarriesV4ClaimsForwardExceptTheFive(t *testing.T) {
	root := repoRoot(t)
	v5, _, _, _, _, _, _ := loadV5Fixture(t, root)
	v4, _, err := LoadV2(filepath.Join(root, "book", "evidence", "v4", "claims.json"))
	if err != nil {
		t.Fatal(err)
	}
	changed := map[string]bool{
		"v2-07-refund-answerability-and-storage":   true,
		"v2-10-guarded-confirm-refusals":           true,
		"v2-14-cache-throughput-gain":              true,
		"v2-16-strict-freshness-read-cost":         true,
		"v3-06-colocated-vs-noncolocated-locality": true,
	}
	var differing []string
	for i := range v5.Claims {
		a, b := v5.Claims[i], v4.Claims[i]
		if a.ClaimID != b.ClaimID {
			t.Fatalf("claim order drifted at %d: %s vs %s", i, a.ClaimID, b.ClaimID)
		}
		same := a.Statement == b.Statement && strings.Join(a.Limits, "\x00") == strings.Join(b.Limits, "\x00")
		if !same && !changed[a.ClaimID] {
			differing = append(differing, a.ClaimID)
		}
		if same && changed[a.ClaimID] {
			t.Errorf("%s is listed as re-scoped but is identical to v4", a.ClaimID)
		}
	}
	if len(differing) > 0 {
		t.Fatalf("claims changed without being listed in the ledger: %v", differing)
	}
}

// TestV5KeepsTheV306Statement guards the disposition that differed from the
// review: the statement stands on the schema's hash primary key, so it must not be
// softened to "is designed to" and the limits must state the basis.
func TestV5KeepsTheV306Statement(t *testing.T) {
	root := repoRoot(t)
	doc, _, _, _, _, _, _ := loadV5Fixture(t, root)
	c := claimByID(t, doc, "v3-06-colocated-vs-noncolocated-locality")
	if !strings.Contains(c.Statement, "keeps one donor's donations in one tablet") {
		t.Errorf("the statement should be retained, not softened: %s", c.Statement)
	}
	joined := strings.Join(c.Limits, " ")
	if !strings.Contains(joined, "PRIMARY KEY ((person_id) HASH") {
		t.Error("a limit must state that the mapping follows from the primary key's hash column")
	}
	if strings.Contains(joined, "the placement labels are intent") {
		t.Error("the old limit, which left the statement's basis unexplained, should be gone")
	}
}
