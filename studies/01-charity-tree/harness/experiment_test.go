package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// This is a behavioral overload test: a service slower than demand must produce
// queueing or rejections, not silently reduce the scheduled arrival rate.
func TestScheduledLoadDoesNotFollowServiceCompletions(t *testing.T) {
	var invoked atomic.Int64
	r := scheduleArrivals(context.Background(), 100*time.Millisecond, 1000, 1, 2, func(context.Context, int64) error { invoked.Add(1); time.Sleep(5 * time.Millisecond); return nil })
	if r.Offered != 100 || r.Dropped == 0 || r.Accepted+r.Dropped != r.Offered || r.Completed+r.Errors != r.Accepted || r.Completed != invoked.Load() {
		t.Fatalf("lost offered work: %+v", r)
	}
	if r.Response.MeanMS < r.Service.MeanMS || r.QueueDelay.Count != r.Completed {
		t.Fatalf("scheduled latency excludes waiting: %+v", r)
	}
}

func TestMechanismControlsKeepUnchangedQueriesIdentical(t *testing.T) {
	// Catalogues now hold sixteen statements: the original twelve, plus
	// q13-q16 (RECENCY.md, v4). q15 joins the "changed" set exactly at the
	// D12 -> D13 transition, for the same reason q02/q05/q12 already do
	// there: that pair is where a direct charity_id filter first becomes
	// possible. q13/q14/q16 take no charity_id and are unaffected by it.
	for _, p := range []struct {
		a, b    string
		changed string
	}{
		{"d2_normalized_indexed", "d11_copied_key", ""},
		{"d11_copied_key", "d12_recency_index", ""},
		{"d12_recency_index", "d13_recency_sql", "q02,q05,q12,q15"},
		{"d13_recency_sql", "d17_sum_sql", "q08"},
		{"d17_sum_sql", "d14_sum_plain", ""},
		{"d14_sum_plain", "d15_sum_covering", ""},
		{"d15_sum_covering", "d3_flattened_fk", "q03,q04"},
	} {
		a, _ := readSQL(p.a, "queries.sql")
		b, _ := readSQL(p.b, "queries.sql")
		as, e := ParseCatalog(a)
		if e != nil {
			t.Fatal(e)
		}
		bs, e := ParseCatalog(b)
		if e != nil {
			t.Fatal(e)
		}
		if len(as) != 16 || len(bs) != 16 {
			t.Fatalf("%s/%s: catalogue lost a query (got %d/%d, want 16)", p.a, p.b, len(as), len(bs))
		}
		for i := range as {
			wantChanged := strings.Contains(p.changed, as[i].Name[:3])
			if (as[i].SQL != bs[i].SQL) != wantChanged {
				t.Errorf("%s -> %s: unexpected SQL change status for %s", p.a, p.b, as[i].Name)
			}
		}
	}
}

// TestLegacyReadScoreExcludesRecencyQueries is the geomean-unchanged test
// EH-02 requires: the study's original twelve-question read score must not
// move just because q13-q16 (RECENCY.md, v4) were added to a run's Reads.
func TestLegacyReadScoreExcludesRecencyQueries(t *testing.T) {
	base := &Run{Reads: []QueryResult{
		{Query: "q01_last_donation_global", OpsPerSec: 100},
		{Query: "q02_last_donation_charity", OpsPerSec: 200},
		{Query: "q03_top_donor_charity", OpsPerSec: 50},
		{Query: "q04_top_donors_leaderboard", OpsPerSec: 60},
		{Query: "q05_last_donor_charity", OpsPerSec: 300},
		{Query: "q06_person_first_last", OpsPerSec: 400},
		{Query: "q07_total_donated_global", OpsPerSec: 20},
		{Query: "q08_total_donated_charity", OpsPerSec: 150},
		{Query: "q09_person_recent_donations", OpsPerSec: 500},
		{Query: "q10_donation_by_id", OpsPerSec: 1000},
		{Query: "q11_person_donation_count", OpsPerSec: 250},
		{Query: "q12_charity_recent_feed", OpsPerSec: 90},
	}}
	want := geomeanReads(base)
	if want <= 0 {
		t.Fatal("fixture geomean is zero; test is not exercising anything")
	}

	withRecency := &Run{Reads: append(append([]QueryResult{}, base.Reads...),
		QueryResult{Query: "q13_donors_last_gift_window@trailing", OpsPerSec: 5},
		QueryResult{Query: "q13_donors_last_gift_window@historical", OpsPerSec: 5000},
		QueryResult{Query: "q14_donors_last_gift_window_count@trailing", OpsPerSec: 7},
		QueryResult{Query: "q14_donors_last_gift_window_count@historical", OpsPerSec: 7000},
		QueryResult{Query: "q15_charity_donors_last_gift_window@trailing", OpsPerSec: 3},
		QueryResult{Query: "q15_charity_donors_last_gift_window@historical", OpsPerSec: 3000},
		QueryResult{Query: "q16_lapsed_donors_count@trailing", OpsPerSec: 9},
		QueryResult{Query: "q16_lapsed_donors_count@historical", OpsPerSec: 9000},
	)}
	got := geomeanReads(withRecency)
	if got != want {
		t.Fatalf("adding recency queries changed the legacy read score: %.6f -> %.6f", want, got)
	}
}

func TestReporterPreservesIndependentLoadsAndFailures(t *testing.T) {
	dir := t.TempDir()
	results := filepath.Join(dir, "results")
	report := filepath.Join(dir, "reports", "test.md")
	for trial := 1; trial <= 3; trial++ {
		r := &Run{RunID: "test", Environment: "test-env", Topology: "pg-single", DesignID: "d3_flattened_fk", Verify: &VerifyReport{Passed: 14}, Dataset: map[string]any{"people": 5000, "donations": 108081}, Experiment: &ExperimentResult{Settings: ExperimentOptions{Trial: trial, Condition: "same", Mode: "reads"}, Keys: map[string]int64{}}}
		if trial == 3 {
			r.Error = "injected verification failure"
		}
		if err := writeJSON(filepath.Join(results, ratesString([]float64{float64(trial)})+".json"), r); err != nil {
			t.Fatal(err)
		}
	}
	fs, err := readExperiments(results)
	if err != nil || len(fs) != 3 {
		t.Fatalf("trial files collapsed: %d %v", len(fs), err)
	}
	if err := writeEnhancementReport(results, report); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(report)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "3 recorded cells; 1 cells") || !strings.Contains(string(b), "injected verification failure") {
		t.Fatalf("report hid failed trial: %s", b)
	}
}

func TestLongHistoryGenerationIsDeterministic(t *testing.T) {
	p := Profile{HistoryMultiplier: 3}
	a, e := GenerateProfile("tiny", 42, p)
	if e != nil {
		t.Fatal(e)
	}
	b, e := GenerateProfile("tiny", 42, p)
	if e != nil {
		t.Fatal(e)
	}
	if len(a.Donations) != len(b.Donations) || len(a.People) != 500 {
		t.Fatal("generation not reproducible")
	}
	for i := range a.Donations {
		if a.Donations[i].ID != b.Donations[i].ID || a.Donations[i].AmountCents != b.Donations[i].AmountCents || !a.Donations[i].DonatedAt.Equal(b.Donations[i].DonatedAt) {
			t.Fatalf("dataset differs at %d", i)
		}
	}
	for _, history := range a.DonationsByPerson {
		if len(history)%3 != 0 {
			t.Fatal("history multiplier not applied per donor")
		}
	}
}

// TestRecencyWindowRegimesDeriveFromTheDataset checks the two window regimes
// RECENCY.md section 2 defines: trailing ends at the dataset's own horizon;
// historical ends well before it. Both are derived from a fixed instant, not
// time.Now(), so two runs of this test (and two runs of the harness) agree.
func TestRecencyWindowRegimesDeriveFromTheDataset(t *testing.T) {
	horizon := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)

	tSince, tUntil := recencyWindowFor(horizon, "trailing")
	if !tUntil.Equal(horizon) {
		t.Fatalf("trailing regime should end at the dataset horizon, got %s", tUntil)
	}
	if got := tUntil.Sub(tSince); got != recencyWindow {
		t.Fatalf("trailing window width = %s, want %s", got, recencyWindow)
	}

	hSince, hUntil := recencyWindowFor(horizon, "historical")
	if !hUntil.Before(horizon) {
		t.Fatal("historical regime must end strictly before the dataset horizon")
	}
	if got := horizon.Sub(hUntil); got != recencyHistoricalLead {
		t.Fatalf("historical lead = %s, want %s", got, recencyHistoricalLead)
	}
	if got := hUntil.Sub(hSince); got != recencyWindow {
		t.Fatalf("historical window width = %s, want %s", got, recencyWindow)
	}

	// The trap RECENCY.md section 2 exists to catch: in the trailing regime,
	// "last gift is in the window" and "gave inside the window" are the same
	// set (nothing can be newer than the dataset's own end). That equivalence
	// must NOT hold for the historical regime -- a donor who gave in that
	// week and kept giving afterwards is in the "gave inside" set but not in
	// the "last gift is in the window" set.
	if !tUntil.Equal(horizon) {
		t.Fatal("trailing regime's accidental-correctness property depends on until == horizon")
	}
	stillGiving := hUntil.Add(30 * 24 * time.Hour) // well after the historical window, before the dataset horizon
	if !stillGiving.Before(horizon) {
		t.Fatal("test fixture error: stillGiving must be within the dataset's history")
	}
	// A donor whose last gift is AFTER the historical window closed does not
	// match recencyMatches for that window, even though they gave inside it
	// at some earlier point -- recencyMatches only ever sees the donor's
	// MAX, by construction, so this is really asserting the predicate takes
	// the donor's LAST gift and nothing earlier.
	if recencyMatches(stillGiving, 0, hSince, hUntil, 0) {
		t.Fatal("a donor whose true last gift is after the historical window must not match it")
	}
}

// TestRecencyExpectedMatchesAgainstHandBuiltDataset is the expected-set
// computation, checked against donors whose window membership is chosen by
// hand rather than generated, so the test data and the code under test
// cannot share a bug.
func TestRecencyExpectedMatchesAgainstHandBuiltDataset(t *testing.T) {
	day := func(n int) time.Time { return time.Date(2026, 1, n, 0, 0, 0, 0, time.UTC) }
	since, until := day(10), day(17) // [10, 17)

	lastByPerson := map[int64]time.Time{
		1: day(15), // inside the window
		2: day(9),  // before the window (lapsed candidate)
		3: day(17), // exactly at "until" -- excluded, half-open interval
		4: day(10), // exactly at "since" -- included, half-open interval
		5: day(20), // well after the window
		6: day(15), // TIE with donor 1
	}
	personCharityID := map[int64]int64{1: 100, 2: 100, 3: 100, 4: 200, 5: 100, 6: 200}

	global := recencyExpectedMatches(lastByPerson, personCharityID, since, until, 0)
	gotIDs := map[int64]bool{}
	for _, r := range global {
		gotIDs[r.PersonID] = true
	}
	wantIDs := map[int64]bool{1: true, 4: true, 6: true}
	if len(gotIDs) != len(wantIDs) {
		t.Fatalf("global window match set = %v, want %v", gotIDs, wantIDs)
	}
	for id := range wantIDs {
		if !gotIDs[id] {
			t.Fatalf("donor %d should match the global window, matches: %v", id, gotIDs)
		}
	}
	// Newest-first order, with the tie (donors 1 and 6, both day 15) sorted
	// ahead of donor 4 (day 10).
	if len(global) != 3 || global[2].PersonID != 4 {
		t.Fatalf("expected order newest-first with donor 4 last, got %+v", global)
	}

	// Charity-scoped: only donor 4 belongs to charity 200 and matches the
	// window (donor 6 also belongs to 200 but is a day-15 tie, present too).
	scoped := recencyExpectedMatches(lastByPerson, personCharityID, since, until, 200)
	if len(scoped) != 2 {
		t.Fatalf("charity-scoped match set = %+v, want donors 4 and 6", scoped)
	}

	lapsed := recencyExpectedLapsedCount(lastByPerson, since)
	if lapsed != 1 { // only donor 2, day 9, is strictly before day 10
		t.Fatalf("lapsed count = %d, want 1 (donor 2 only)", lapsed)
	}
}

// TestDiagnoseFlagMismatchClassifiesBySymptom is the recency audit's
// classification logic (RECENCY.md section 5), unit tested directly since it
// needs no database.
func TestDiagnoseFlagMismatchClassifiesBySymptom(t *testing.T) {
	cases := []struct {
		count int
		want  string
	}{
		{0, "no flagged donation: the flag was cleared but never moved forward"},
		{1, "flagged donation is not the donor's newest"},
		{2, "more than one flagged donation for this donor -- the unguarded race (RECENCY.md D21)"},
		{5, "more than one flagged donation for this donor -- the unguarded race (RECENCY.md D21)"},
	}
	for _, c := range cases {
		if got := diagnoseFlagMismatch(c.count); got != c.want {
			t.Errorf("diagnoseFlagMismatch(%d) = %q, want %q", c.count, got, c.want)
		}
	}
}

func TestExperimentReportDoesNotPresentDocDBSizeAsZeroBytes(t *testing.T) {
	dir := t.TempDir()
	results := filepath.Join(dir, "results")
	report := filepath.Join(dir, "reports", "test.md")
	r := &Run{RunID: "test", Environment: "test-env", Engine: "yugabyte", Topology: "yb-single", DesignID: "d3_flattened_fk", Verify: &VerifyReport{Passed: 14}, Stats: &DBStats{}, Dataset: map[string]any{"people": 5000, "donations": 108081}, Experiment: &ExperimentResult{Settings: ExperimentOptions{Trial: 1, Condition: "same", Mode: "verify"}, Keys: map[string]int64{}}}
	if err := writeJSON(filepath.Join(results, "trial.json"), r); err != nil {
		t.Fatal(err)
	}
	if err := writeEnhancementReport(results, report); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(report)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "DocDB physical storage bytes were not collected") || !strings.Contains(string(b), "| 5000 | 108081 | unavailable |") {
		t.Fatalf("missing storage became a measured zero: %s", b)
	}
}
