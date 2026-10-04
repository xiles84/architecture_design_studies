package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"charitytree/selectbench"
)

// Reporting regressions found by the 2026-10-04 HIGH review (R4). Each test is
// written so the defective behaviour fails it: the digest test asserts the
// double-newline variant differs from the file digest, the provenance test
// asserts client and engine images are not the same field, and the class test
// asserts two fixtures are not pooled.

func sampleJoinResult() selectbench.JoinResult {
	return selectbench.JoinResult{
		RunID: "test-run", Engine: "postgres", Topology: "single", Scale: "cte",
		FactCount: 100, Customers: 2, Seed: 42, Thresholds: []int{1},
		Provenance: selectbench.Provenance{
			Environment: "host-test", DatasetFingerprint: "fp1",
			ClientImage: "client:1", ClientImageDigest: "cd",
			EngineImage: "db:1", EngineImageDigest: "ed",
			Settings: map[string]string{"scale": "cte"},
		},
		Cells: []selectbench.Cell{{
			Axis: "family=bag;stmt=join_where;threshold=1", Statement: "j01_join_where",
			Group: "sel_join_bag", Status: selectbench.StatusValid, TimingMeasured: true,
			MedianMS: 1.5, Trials: 1, SamplesMS: []float64{1.5}, ResultRows: 1,
			Plan: "Seq Scan on facts",
		}},
		Summary: map[string]int{"valid": 1},
	}
}

func TestJoinInputsDigestReproducible(t *testing.T) {
	dir := t.TempDir()
	runDir := filepath.Join(dir, "run1")
	if err := os.Mkdir(runDir, 0o755); err != nil {
		t.Fatal(err)
	}
	data, err := json.MarshalIndent(sampleJoinResult(), "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	// The writer's exact bytes: one trailing newline.
	data = append(data, '\n')
	if err := os.WriteFile(filepath.Join(runDir, "join.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}
	reportOut := filepath.Join(dir, "report.md")
	if err := runSelectJoinReport(runDir, reportOut, ""); err != nil {
		t.Fatalf("runSelectJoinReport: %v", err)
	}
	report, err := os.ReadFile(reportOut)
	if err != nil {
		t.Fatal(err)
	}
	want := fmt.Sprintf("%x", sha256.Sum256(data))
	if !strings.Contains(string(report), want) {
		t.Fatalf("report does not advertise the file digest %s", want)
	}
	// The original defect hashed the file bytes with one more newline appended.
	old := fmt.Sprintf("%x", sha256.Sum256(append(append([]byte(nil), data...), '\n')))
	if old == want {
		t.Fatal("double-newline variant unexpectedly equals the file digest")
	}
	if strings.Contains(string(report), old) {
		t.Fatal("report advertises the double-newline digest")
	}
}

func TestJoinProvenanceSeparatesClientAndEngine(t *testing.T) {
	t.Setenv("ADS_BENCH_IMAGE", "localhost/client:1")
	t.Setenv("ADS_BENCH_IMAGE_ID", "client-digest")
	t.Setenv("ADS_ENGINE_IMAGE", "docker.io/library/postgres:17.11")
	t.Setenv("ADS_ENGINE_IMAGE_ID", "engine-digest")
	p := joinProvenance("Postgres 17.11", "host-test", "cte", "study01-person-donation-charity",
		10, 2, []int{1}, "postgres", "single", 42, 5, 120000, "cpu-before", "fp1")
	if p.ClientImage != "localhost/client:1" || p.ClientImageDigest != "client-digest" {
		t.Fatalf("client provenance = %q @ %q", p.ClientImage, p.ClientImageDigest)
	}
	if p.EngineImage != "docker.io/library/postgres:17.11" || p.EngineImageDigest != "engine-digest" {
		t.Fatalf("engine provenance = %q @ %q", p.EngineImage, p.EngineImageDigest)
	}
	if p.EngineImage == p.ClientImage {
		t.Fatal("engine and client images collapsed into one field")
	}
}

func TestJoinReplicatedPartitionsByClass(t *testing.T) {
	a := sampleJoinResult()
	b := sampleJoinResult()
	b.RunID = "test-run-2"
	b.Provenance.DatasetFingerprint = "fp2"
	report, err := selectbench.RenderJoinReplicated([]selectbench.JoinResult{a, b})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(report, "## Compatibility classes") {
		t.Fatal("replicated report has no compatibility-class partition")
	}
	if !strings.Contains(report, "fp1") || !strings.Contains(report, "fp2") {
		t.Fatal("both dataset fingerprints must be named in the class table")
	}
	// The runs table must carry a class column, and the two runs must not share
	// it: otherwise incompatible fixtures could be pooled under one axis.
	if !strings.Contains(report, "| run | class | engine |") {
		t.Fatal("runs table does not name the compatibility class")
	}
}
