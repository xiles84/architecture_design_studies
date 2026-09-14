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
	for _, p := range []struct {
		a, b    string
		changed string
	}{
		{"d2_normalized_indexed", "d11_copied_key", ""},
		{"d11_copied_key", "d12_recency_index", ""},
		{"d12_recency_index", "d13_recency_sql", "q02,q05,q12"},
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
		if len(as) != 12 || len(bs) != 12 {
			t.Fatal("catalogue lost a query")
		}
		for i := range as {
			wantChanged := strings.Contains(p.changed, as[i].Name[:3])
			if (as[i].SQL != bs[i].SQL) != wantChanged {
				t.Errorf("%s -> %s: unexpected SQL change status for %s", p.a, p.b, as[i].Name)
			}
		}
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
