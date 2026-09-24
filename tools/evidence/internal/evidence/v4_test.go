package evidence

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func loadV4Fixture(t *testing.T, root string) (*V2Document, any, *V2Confounds, *Document, *V2Document, *V2Document) {
	t.Helper()
	doc, raw, err := LoadV2(filepath.Join(root, "book", "evidence", "v4", "claims.json"))
	if err != nil {
		t.Fatal(err)
	}
	conf, err := LoadConfounds(filepath.Join(root, "book", "evidence", "v4", "confounds.json"))
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
	return doc, raw, conf, v1, v2, v3
}

func loadV4Schema(t *testing.T, root string) map[string]any {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, "book", "evidence", "v4", "schema.json"))
	if err != nil {
		t.Fatal(err)
	}
	var schema map[string]any
	if err := json.Unmarshal(data, &schema); err != nil {
		t.Fatal(err)
	}
	return schema
}

// TestV4RegistryValidates is the positive check: the committed v4 package must
// resolve every structured support key, every anchor, every supersession and
// every retirement against v1, v2 or v3, and satisfy the v4 lifecycle rules.
func TestV4RegistryValidates(t *testing.T) {
	root := repoRoot(t)
	doc, raw, conf, v1, v2, v3 := loadV4Fixture(t, root)
	schema := loadV4Schema(t, root)
	if errs := ValidateV4(root, schema, doc, raw, conf, v1, v2, v3); len(errs) > 0 {
		for _, e := range errs {
			t.Errorf("unexpected error: %v", e)
		}
	}
	if len(doc.Claims) != 29 {
		t.Fatalf("expected the v4 package to hold 29 claims (two retired, two added), got %d", len(doc.Claims))
	}
	if doc.SchemaVersion != 4 {
		t.Fatalf("expected schema_version 4, got %d", doc.SchemaVersion)
	}
}

// TestV4RetiresTheStalePlacementGaps is the regression test for the defect v4
// exists to repair. v3 printed "Study 05's placement pair was never executed"
// and "balanced client access across cluster endpoints is not proven" beside the
// v3-06 and v3-05 claims that contradict them. Those two claims must be retired
// here, and their successors must exist.
func TestV4RetiresTheStalePlacementGaps(t *testing.T) {
	root := repoRoot(t)
	doc, _, _, _, _, _ := loadV4Fixture(t, root)

	retired := map[string]V2Retired{}
	for _, r := range doc.RetiredPredecessorClaims {
		retired[r.ClaimID] = r
	}
	for _, id := range []string{"v2-gap-05-colocation-unverified", "v2-gap-07-real-network-balanced-endpoints"} {
		r, ok := retired[id]
		if !ok {
			t.Fatalf("%s must be retired in v4", id)
		}
		if len(r.SupersededBy) == 0 {
			t.Errorf("%s: a retirement must name its successor", id)
		}
		if len(r.ClosedDimensions) == 0 {
			t.Errorf("%s: a retirement must name the clause it closes", id)
		}
	}
	for _, c := range doc.Claims {
		if c.ClaimID == "v2-gap-05-colocation-unverified" || c.ClaimID == "v2-gap-07-real-network-balanced-endpoints" {
			t.Fatalf("%s is retired but still active in claims[]", c.ClaimID)
		}
	}
	// The successors must be active, and must supersede the retired claim.
	want := map[string]string{
		"v4-gap-01-physical-placement-unverified": "v2-gap-05-colocation-unverified",
		"v4-gap-02-endpoint-distribution-partial": "v2-gap-07-real-network-balanced-endpoints",
	}
	for _, c := range doc.Claims {
		supersedes, ok := want[c.ClaimID]
		if !ok {
			continue
		}
		delete(want, c.ClaimID)
		if c.Supersedes == nil || *c.Supersedes != supersedes {
			t.Errorf("%s must supersede %s, got %v", c.ClaimID, supersedes, c.Supersedes)
		}
		if c.Status != "active" {
			t.Errorf("%s should be active, got %q", c.ClaimID, c.Status)
		}
		if c.Strength != "gap" {
			t.Errorf("%s should be a gap, got %q", c.ClaimID, c.Strength)
		}
	}
	if len(want) > 0 {
		t.Fatalf("successor gap claims missing from v4: %v", want)
	}
}

// TestV4GapMustDeclareKind: every gap declares what sort of gap it is, so the
// book can badge a schema limitation differently from a regime nobody ran.
func TestV4GapMustDeclareKind(t *testing.T) {
	root := repoRoot(t)
	doc, raw, conf, v1, v2, v3 := loadV4Fixture(t, root)
	schema := loadV4Schema(t, root)

	for _, c := range doc.Claims {
		if c.Strength == "gap" && len(c.GapKind) == 0 {
			t.Fatalf("fixture precondition: %s has no gap_kind", c.ClaimID)
		}
	}
	for _, c := range doc.Claims {
		if c.ClaimID == "v2-gap-02-hold-funnel-unanswerable" {
			if len(c.GapKind) != 1 || c.GapKind[0] != "schema_limitation" {
				t.Fatalf("v2-gap-02 is a schema limitation, got %v", c.GapKind)
			}
			c.GapKind = nil
		}
	}
	errs := ValidateV4(root, schema, doc, raw, conf, v1, v2, v3)
	if !hasSubstring(errs, "must declare gap_kind") {
		t.Fatalf("expected a missing-gap_kind error, got %v", errs)
	}
}

// TestV4RetiredStatusCannotStayActive: the lifecycle states that retire a claim
// may not be used as a label inside claims[], because claims[] is exactly what
// the book renders.
func TestV4RetiredStatusCannotStayActive(t *testing.T) {
	root := repoRoot(t)
	doc, raw, conf, v1, v2, v3 := loadV4Fixture(t, root)
	schema := loadV4Schema(t, root)

	doc.Claims[0].Status = "retired"
	errs := ValidateV4(root, schema, doc, raw, conf, v1, v2, v3)
	if !hasSubstring(errs, "must not appear in claims[]") {
		t.Fatalf("expected a retired-in-claims error, got %v", errs)
	}
}

// TestV4PartialSupersessionNeedsSuccessors: a claim marked partially
// superseded must say by what and what remains open, or the label carries no
// information.
func TestV4PartialSupersessionNeedsSuccessors(t *testing.T) {
	root := repoRoot(t)
	doc, raw, conf, v1, v2, v3 := loadV4Fixture(t, root)
	schema := loadV4Schema(t, root)

	patched := false
	for _, c := range doc.Claims {
		if c.ClaimID == "v3-gap-01-study05-remaining-regimes" {
			if c.Status != "partially_superseded" || len(c.SupersededBy) != 2 || len(c.RemainingDimensions) == 0 {
				t.Fatalf("v3-gap-01 must be partially superseded by its two successors, got status=%q superseded_by=%v remaining=%v",
					c.Status, c.SupersededBy, c.RemainingDimensions)
			}
			c.SupersededBy = nil
			patched = true
		}
	}
	if !patched {
		t.Fatal("v3-gap-01-study05-remaining-regimes is missing from v4")
	}
	errs := ValidateV4(root, schema, doc, raw, conf, v1, v2, v3)
	if !hasSubstring(errs, "requires superseded_by") {
		t.Fatalf("expected a missing-superseded_by error, got %v", errs)
	}
}

// TestV4SuccessorMustResolve: a retirement or a partial supersession may not
// point at a claim that does not exist in this registry.
func TestV4SuccessorMustResolve(t *testing.T) {
	root := repoRoot(t)
	doc, raw, conf, v1, v2, v3 := loadV4Fixture(t, root)
	schema := loadV4Schema(t, root)

	doc.RetiredPredecessorClaims[len(doc.RetiredPredecessorClaims)-1].SupersededBy = []string{"v4-does-not-exist"}
	errs := ValidateV4(root, schema, doc, raw, conf, v1, v2, v3)
	if !hasSubstring(errs, "is not an active claim in this registry") {
		t.Fatalf("expected an unresolved-successor error, got %v", errs)
	}
}

// TestV4CarriesV3ClaimsForwardUnchanged is the guarantee that a lifecycle
// update is not a licence to rewrite evidence: every claim that still exists in
// v4 and in v3 must have the same statement.
func TestV4CarriesV3ClaimsForwardUnchanged(t *testing.T) {
	root := repoRoot(t)
	_, _, _, _, _, v3 := loadV4Fixture(t, root)
	v4, _, err := LoadV2(filepath.Join(root, "book", "evidence", "v4", "claims.json"))
	if err != nil {
		t.Fatal(err)
	}
	byID := map[string]string{}
	for _, c := range v3.Claims {
		byID[c.ClaimID] = c.Statement
	}
	carried := 0
	for _, c := range v4.Claims {
		old, ok := byID[c.ClaimID]
		if !ok {
			continue
		}
		carried++
		if old != c.Statement {
			t.Errorf("%s: statement changed from v3 to v4:\n  v3: %s\n  v4: %s", c.ClaimID, old, c.Statement)
		}
	}
	// 27 of v3's 29 claims survive; two gaps are retired, and the two v4 gaps
	// are new.
	if carried != 27 {
		t.Fatalf("expected 27 v3 claims carried forward, got %d", carried)
	}
}
