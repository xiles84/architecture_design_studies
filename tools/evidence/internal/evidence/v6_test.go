package evidence

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func loadV6Fixture(t *testing.T, root string) (*V2Document, any, *V2Confounds, *Document, *V2Document, *V2Document, *V2Document, *V2Document) {
	t.Helper()
	doc, raw, err := LoadV2(filepath.Join(root, "book", "evidence", "v6", "claims.json"))
	if err != nil {
		t.Fatal(err)
	}
	conf, err := LoadConfounds(filepath.Join(root, "book", "evidence", "v6", "confounds.json"))
	if err != nil {
		t.Fatal(err)
	}
	v1, _, err := Load(filepath.Join(root, "book", "evidence", "claims.json"))
	if err != nil {
		t.Fatal(err)
	}
	load := func(version string) *V2Document {
		d, _, err := LoadV2(filepath.Join(root, "book", "evidence", version, "claims.json"))
		if err != nil {
			t.Fatal(err)
		}
		return d
	}
	return doc, raw, conf, v1, load("v2"), load("v3"), load("v4"), load("v5")
}

func loadV6Schema(t *testing.T, root string) map[string]any {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, "book", "evidence", "v6", "schema.json"))
	if err != nil {
		t.Fatal(err)
	}
	var schema map[string]any
	if err := json.Unmarshal(data, &schema); err != nil {
		t.Fatal(err)
	}
	return schema
}

// TestV6RegistryValidates is the positive check: the committed v6 package must
// resolve every structured support key, every anchor, every supersession and every
// retirement against v1..v5, satisfy the v4 lifecycle and v5 repetition rules, and
// pass every method lint added in v6.
func TestV6RegistryValidates(t *testing.T) {
	root := repoRoot(t)
	doc, raw, conf, v1, v2, v3, v4, v5 := loadV6Fixture(t, root)
	schema := loadV6Schema(t, root)
	if errs := ValidateV6(root, schema, doc, raw, conf, v1, v2, v3, v4, v5); len(errs) > 0 {
		for _, e := range errs {
			t.Errorf("unexpected error: %v", e)
		}
	}
	// 29 carried from v5 plus v2-17 (the first bounded crossover) and v2-18 (the
	// first interaction).
	if len(doc.Claims) != 31 {
		t.Fatalf("want 31 claims, got %d", len(doc.Claims))
	}
}

// TestV6CarriesV5ForwardExceptTheReplicatedClaim pins that v6 adds method metadata
// without re-writing the corpus: every statement, limit, confound, trial count and
// anchor is identical to v5 except for the single claim the first run-level
// replication re-scoped (v2-13). A later pass that re-scopes another claim must add
// it here, in the supersession ledger, and in the analysis.
func TestV6CarriesV5ForwardExceptTheReplicatedClaim(t *testing.T) {
	root := repoRoot(t)
	v6, _, _, _, _, _, _, _ := loadV6Fixture(t, root)
	v5, _, err := LoadV2(filepath.Join(root, "book", "evidence", "v5", "claims.json"))
	if err != nil {
		t.Fatal(err)
	}
	// The claims added after v5, by id. Adding a claim is a deliberate act: it
	// must be listed here and in the supersession ledger, or this test fails.
	added := map[string]bool{
		"v2-17-optimistic-pessimistic-crossover": true,
		"v2-18-maintenance-size-interaction":     true,
	}
	if len(v6.Claims) != len(v5.Claims)+len(added) {
		t.Fatalf("claim count drifted: v6 %d v5 %d plus %d declared additions",
			len(v6.Claims), len(v5.Claims), len(added))
	}
	// The one claim the 2026-09-25 T1 replication corrected: re-scoped because new
	// evidence arrived. Its numbers and its limits both changed.
	// Claims re-scoped because new evidence arrived: their numbers, limits, trial
	// counts and anchors may all differ from v5. v2-06 was only reworded at the
	// v6 tag; the 2026-09-26 replication then re-scoped it, so it moved here.
	changed := map[string]bool{
		"v2-13-trigger-rollup-cost": true,
		"v2-06-hot-drop-designs":    true,
	}
	// Claims whose wording was made plainer without touching the evidence. Their
	// statement and limits may differ from v5; their confounds, trial counts and
	// anchors may not, which is what keeps a rewording from quietly carrying a
	// change of substance. Listed here and in the ledger, so neither a rewording
	// nor a re-scoping can pass as a silent edit.
	reworded := map[string]bool{
		"v2-01-normalized-index-first":              true,
		"v2-02-rollup-family":                       true,
		"v2-04-recency-maintained-index":            true,
		"v2-05-d3-read-advantage-repeated":          true,
		"v2-06-hot-drop-designs":                    true,
		"v2-07-refund-answerability-and-storage":    true,
		"v2-09-postgres-ledger-ratio-unsettled":     true,
		"v2-10-guarded-confirm-refusals":            true,
		"v2-11-section-seat-map":                    true,
		"v2-13-trigger-rollup-cost":                 true,
		"v2-14-cache-throughput-gain":               true,
		"v2-15-three-instance-staleness":            true,
		"v2-gap-01-ledger-load-multiplier-unstable": true,
		"v3-01-churn-crosses-real-ttl":              true,
		"v3-02-churn-600-through-vs-aside":          true,
		"v3-03-equal-total-framing-labelled":        true,
		"v3-04-yb-single-strict-relaxed":            true,
		"v3-05-yb-cluster3-strict-relaxed":          true,
		"v3-06-colocated-vs-noncolocated-locality":  true,
		"v3-gap-01-study05-remaining-regimes":       true,
	}
	for i := range v5.Claims {
		a, b := v6.Claims[i], v5.Claims[i]
		if a.ClaimID != b.ClaimID {
			t.Fatalf("claim order drifted at %d: %s vs %s", i, a.ClaimID, b.ClaimID)
		}
		structureSame := strings.Join(a.Confounds, "\x00") == strings.Join(b.Confounds, "\x00") &&
			a.Trials == b.Trials &&
			len(a.Anchors) == len(b.Anchors)
		textSame := a.Statement == b.Statement &&
			strings.Join(a.Limits, "\x00") == strings.Join(b.Limits, "\x00")
		same := textSame && structureSame
		// A re-scoped claim may change anything, because new evidence arrived.
		// A reworded claim may change its text and nothing else, which is the
		// rule that keeps a rewording from carrying a change of substance.
		if changed[a.ClaimID] {
			if same {
				t.Errorf("%s is listed as re-scoped but is identical to v5", a.ClaimID)
			}
			continue
		}
		if reworded[a.ClaimID] {
			if textSame {
				t.Errorf("%s is listed as reworded but its text is identical to v5", a.ClaimID)
			}
			if !structureSame {
				t.Errorf("%s is listed as reworded but its confounds, trial counts or anchors moved; a rewording may never move the evidence structure", a.ClaimID)
			}
			continue
		}
		if !same {
			t.Errorf("%s changed without being listed as re-scoped or reworded", a.ClaimID)
		}
	}
}

// TestV6PackageWithoutIsolationFails is the negative control for the no-package-
// attribution rule: a package result that neither withholds attribution nor names
// an isolation pair must not validate.
func TestV6PackageWithoutIsolationFails(t *testing.T) {
	root := repoRoot(t)
	doc, _, _, _, _, _, _, _ := loadV6Fixture(t, root)
	c := claimByID(t, doc, "v2-02-rollup-family")
	if c.EvidenceForm != "package" {
		t.Fatalf("fixture drift: v2-02 should be a package, is %q", c.EvidenceForm)
	}
	c.Attribution.Withheld = false
	c.Attribution.IsolatedBy = nil
	errs := validateEvidenceForm(doc)
	if len(errs) == 0 {
		t.Fatal("expected a package claim with no isolation pair and no withheld attribution to fail")
	}
	if !strings.Contains(errs[0].Error(), "package result may not be attributed") {
		t.Fatalf("error should explain the attribution rule, got: %v", errs[0])
	}
}

// TestV6ConfoundedClaimWithoutWithholdingFails is the other half: a claim that
// lists a confound must either withhold attribution or name the pair that isolates
// its headline.
func TestV6ConfoundedClaimWithoutWithholdingFails(t *testing.T) {
	root := repoRoot(t)
	doc, _, _, _, _, _, _, _ := loadV6Fixture(t, root)
	c := claimByID(t, doc, "v2-01-normalized-index-first")
	if len(c.Confounds) == 0 {
		t.Fatal("fixture drift: v2-01 should still list conf-01")
	}
	c.Attribution.Withheld = false
	c.Attribution.IsolatedBy = nil
	errs := validateEvidenceForm(doc)
	found := false
	for _, e := range errs {
		if strings.Contains(e.Error(), "confound") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected the confound rule to fire, got: %v", errs)
	}
}

// TestV6ControlCoveragePresentWithoutFiringFails covers the read-only-phase hole:
// controls that were present but never fired are not coverage, so a claim cannot
// record control_coverage present beside controls_fired 0.
func TestV6ControlCoveragePresentWithoutFiringFails(t *testing.T) {
	root := repoRoot(t)
	doc, _, _, _, _, _, _, _ := loadV6Fixture(t, root)
	c := claimByID(t, doc, "v2-07-refund-answerability-and-storage")
	if c.RunLevel.ControlsPresent == 0 || c.RunLevel.ControlsFired != 0 {
		t.Fatalf("fixture drift: v2-07 should have controls present and none fired, has %d/%d",
			c.RunLevel.ControlsPresent, c.RunLevel.ControlsFired)
	}
	if c.ControlCoverage != "absent" {
		t.Fatalf("fixture drift: v2-07 should record control_coverage absent, is %q", c.ControlCoverage)
	}
	c.ControlCoverage = "present"
	errs := validateControlCoverage(doc)
	if len(errs) == 0 {
		t.Fatal("expected control_coverage present beside controls_fired 0 to fail")
	}
	if !strings.Contains(errs[0].Error(), "never fired") {
		t.Fatalf("error should explain that an unfired control is not coverage, got: %v", errs[0])
	}
}

// TestV6RunReplicatedRequiresRunsOrderAndPrediction covers the replication rule: a
// run_replicated claim needs at least two whole runs, a non-fixed order, and a
// recorded prediction with a verdict.
func TestV6RunReplicatedRequiresRunsOrderAndPrediction(t *testing.T) {
	root := repoRoot(t)
	doc, _, _, _, _, _, _, _ := loadV6Fixture(t, root)
	c := claimByID(t, doc, "v2-05-d3-read-advantage-repeated")
	c.Strength = "run_replicated"
	c.Prediction = nil
	c.Replication.Runs = 1
	c.Replication.Order = "fixed"
	c.Replication.Tier = "single-run"
	joined := ""
	for _, e := range validateReplication(doc) {
		joined += e.Error() + "\n"
	}
	for _, want := range []string{"recorded prediction", "at least 2 whole runs", "not fixed"} {
		if !strings.Contains(joined, want) {
			t.Fatalf("expected an error containing %q, got:\n%s", want, joined)
		}
	}
}

// TestV6HardwareReplicatedRequiresHardwareTier prevents the two replication tiers
// from being conflated: a hardware_replicated claim must run on separate hosts.
func TestV6HardwareReplicatedRequiresHardwareTier(t *testing.T) {
	root := repoRoot(t)
	doc, _, _, _, _, _, _, _ := loadV6Fixture(t, root)
	c := claimByID(t, doc, "v2-05-d3-read-advantage-repeated")
	c.Strength = "hardware_replicated"
	c.Prediction = &V2Prediction{Text: "replicated across hosts", Source: "fixture", Verdict: "unresolved"}
	c.Replication.Runs = 2
	c.Replication.Order = "randomized"
	c.Replication.Tier = "run-level"
	errs := validateReplication(doc)
	found := false
	for _, e := range errs {
		if strings.Contains(e.Error(), "hardware-independent") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected the tier rule to fire, got: %v", errs)
	}
}

// TestV6CrossoverReversalNeedsLevelsAndInterval covers the crossover rule: a
// reversal must be bounded by at least five levels and state the crossing interval.
func TestV6CrossoverReversalNeedsLevelsAndInterval(t *testing.T) {
	root := repoRoot(t)
	doc, _, _, _, _, _, _, _ := loadV6Fixture(t, root)
	c := claimByID(t, doc, "v2-12-lock-vs-version-upper-bound")
	c.Crossover = &V2Crossover{Dimension: "writers", Levels: 3, ReversalObserved: true, CrossingInterval: nil}
	errs := validateCrossover(doc)
	if len(errs) < 2 {
		t.Fatalf("expected both the level and interval rules to fire, got: %v", errs)
	}
	// The same sweep is acceptable once it has the levels and an interval.
	interval := "between 1 and 4 writers"
	c.Crossover = &V2Crossover{Dimension: "writers", Levels: 5, ReversalObserved: true, CrossingInterval: &interval}
	if errs := validateCrossover(doc); len(errs) > 0 {
		t.Fatalf("a properly bounded crossover should pass, got: %v", errs)
	}
}

// TestV6EveryClaimCarriesMethodMetadata keeps the migration honest: every claim
// states its evidence form, an attribution record with a note, a control-coverage
// value, and a replication record whose tier matches its run count.
func TestV6EveryClaimCarriesMethodMetadata(t *testing.T) {
	root := repoRoot(t)
	doc, _, _, _, _, _, _, _ := loadV6Fixture(t, root)
	for _, c := range doc.Claims {
		if c.EvidenceForm == "" {
			t.Errorf("%s: missing evidence_form", c.ClaimID)
		}
		if c.Attribution.Note == "" {
			t.Errorf("%s: missing attribution note", c.ClaimID)
		}
		if c.ControlCoverage == "" {
			t.Errorf("%s: missing control_coverage", c.ClaimID)
		}
		if c.Replication == nil {
			t.Errorf("%s: missing replication record", c.ClaimID)
			continue
		}
		if c.Replication.Runs != c.Trials.WholeRunReplications {
			t.Errorf("%s: replication.runs %d disagrees with trials %d",
				c.ClaimID, c.Replication.Runs, c.Trials.WholeRunReplications)
		}
	}
}
