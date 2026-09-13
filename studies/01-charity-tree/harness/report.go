package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ---------------------------------------------------------------------------
// Report generation
//
// Turns a results directory into the tables a reader actually needs. The raw
// JSON stays the source of truth; this only arranges it. Every number printed
// here can be traced back to exactly one file.
// ---------------------------------------------------------------------------

// designOrder fixes the column order everywhere in the report. Designs are laid
// out as a progression (plain tree -> indexed -> denormalised -> rolled up ->
// embedded), so reading a row left to right follows the argument.
var designOrder = []string{
	"d11_copied_key", "d12_recency_index", "d13_recency_sql", "d17_sum_sql", "d14_sum_plain", "d15_sum_covering", "d16_sum_rollup",
	"d1_normalized_minimal",
	"d2_normalized_indexed",
	"d3_flattened_fk",
	"d8_flattened_nofk",
	"d4_rollup_trigger",
	"d5_rollup_app",
	"d6_embedded_jsonb",
	"d9_embedded_hybrid",
	"d10_embedded_hybrid_locked",
	"d7_yb_child_colocated",
}

var designShort = map[string]string{
	"d11_copied_key": "D11 copied key", "d12_recency_index": "D12 date index", "d13_recency_sql": "D13 recency SQL", "d17_sum_sql": "D17 sum SQL", "d14_sum_plain": "D14 plain sum index", "d15_sum_covering": "D15 covering sum index", "d16_sum_rollup": "D16 sum rollup",
	"d1_normalized_minimal":      "D1 minimal",
	"d2_normalized_indexed":      "D2 indexed",
	"d3_flattened_fk":            "D3 flat+FK",
	"d8_flattened_nofk":          "D8 flat−FK",
	"d4_rollup_trigger":          "D4 rollup/trg",
	"d5_rollup_app":              "D5 rollup/app",
	"d6_embedded_jsonb":          "D6 embedded",
	"d9_embedded_hybrid":         "D9 hybrid",
	"d10_embedded_hybrid_locked": "D10 hybrid/locked",
	"d7_yb_child_colocated":      "D7 yb-coloc",
}

var topologyOrder = []string{"pg-single", "yb-single", "yb-cluster3"}

var topologyLabel = map[string]string{
	"pg-single":   "PostgreSQL, 1 node",
	"yb-single":   "YugabyteDB, 1 node (RF=1)",
	"yb-cluster3": "YugabyteDB, 3 nodes (RF=3)",
}

// queryQuestion maps each query back to the question it answers, so the tables
// read as answers to the study's questions rather than as opaque ids.
var queryQuestion = map[string]string{
	"q01_last_donation_global":    "Information of the last donation (global)",
	"q02_last_donation_charity":   "Information of the last donation (one charity)",
	"q03_top_donor_charity":       "Who donates the most",
	"q04_top_donors_leaderboard":  "Top-10 donor leaderboard",
	"q05_last_donor_charity":      "Who donated last",
	"q06_person_first_last":       "First and last donation of a person",
	"q07_total_donated_global":    "Total donated (global)",
	"q08_total_donated_charity":   "Total donated (one charity)",
	"q09_person_recent_donations": "A person's 20 most recent donations",
	"q10_donation_by_id":          "Donation by id (point lookup)",
	"q11_person_donation_count":   "How many donations a person made",
	"q12_charity_recent_feed":     "Charity activity feed (last 50, with names)",
}

type resultSet struct {
	runs map[string]map[string]*Run // topology -> design -> run
	// preserved orderings restricted to what the data actually contains
	topologies []string
	designs    []string
}

func loadResults(dir string) (*resultSet, error) {
	rs := &resultSet{runs: map[string]map[string]*Run{}}
	err := filepath.WalkDir(dir, func(path string, e os.DirEntry, err error) error {
		if err != nil || e.IsDir() || !strings.HasSuffix(path, ".json") {
			return err
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		var r Run
		if err := json.Unmarshal(b, &r); err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		if r.Topology == "" || r.DesignID == "" {
			return nil
		}
		if rs.runs[r.Topology] == nil {
			rs.runs[r.Topology] = map[string]*Run{}
		}
		rs.runs[r.Topology][r.DesignID] = &r
		return nil
	})
	if err != nil {
		return nil, err
	}
	for _, t := range topologyOrder {
		if _, ok := rs.runs[t]; ok {
			rs.topologies = append(rs.topologies, t)
		}
	}
	seen := map[string]bool{}
	for _, d := range designOrder {
		for _, t := range rs.topologies {
			if _, ok := rs.runs[t][d]; ok && !seen[d] {
				seen[d] = true
				rs.designs = append(rs.designs, d)
			}
		}
	}
	if len(rs.topologies) == 0 {
		return nil, fmt.Errorf("no result files found under %s", dir)
	}
	return rs, nil
}

func (rs *resultSet) designsIn(topo string) []string {
	var out []string
	for _, d := range rs.designs {
		if _, ok := rs.runs[topo][d]; ok {
			out = append(out, d)
		}
	}
	return out
}

func (rs *resultSet) any() *Run {
	for _, t := range rs.topologies {
		for _, d := range rs.designsIn(t) {
			return rs.runs[t][d]
		}
	}
	return nil
}

func readOf(r *Run, q string) *QueryResult {
	if r == nil {
		return nil
	}
	for i := range r.Reads {
		if r.Reads[i].Query == q {
			return &r.Reads[i]
		}
	}
	return nil
}

func writeOf(r *Run, op string) *WriteResult {
	if r == nil {
		return nil
	}
	for i := range r.Writes {
		if r.Writes[i].Op == op {
			return &r.Writes[i]
		}
	}
	return nil
}

// fmtOps renders throughput at a precision that matches how much it varies:
// four significant figures on a number like 48 090 is false precision.
func fmtOps(v float64) string {
	switch {
	case v <= 0:
		return "—"
	case v >= 10000:
		return fmt.Sprintf("%.0fk", v/1000)
	case v >= 1000:
		return fmt.Sprintf("%.1fk", v/1000)
	case v >= 100:
		return fmt.Sprintf("%.0f", v)
	default:
		return fmt.Sprintf("%.1f", v)
	}
}

func fmtMS(v float64) string {
	if v <= 0 {
		return "—"
	}
	if v < 10 {
		return fmt.Sprintf("%.2f", v)
	}
	return fmt.Sprintf("%.0f", v)
}

func fmtBytes(b int64) string {
	if b <= 0 {
		return "—"
	}
	const u = 1024
	if b < u {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(u), 0
	for n := b / u; n >= u; n /= u {
		div *= u
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGT"[exp])
}

// ratio expresses b relative to a as a human-facing multiplier. Speedups and
// slowdowns are both rendered as "Nx faster/slower" rather than as a percentage,
// because the effects in this study span four orders of magnitude.
func ratio(a, b float64) string {
	if a <= 0 || b <= 0 {
		return "—"
	}
	if b >= a {
		return fmt.Sprintf("%.1fx faster", b/a)
	}
	return fmt.Sprintf("%.1fx slower", a/b)
}

func writeReport(resultsDir, outPath string) error {
	rs, err := loadResults(resultsDir)
	if err != nil {
		return err
	}
	var b strings.Builder
	ref := rs.any()

	fmt.Fprintf(&b, "# Study 01 — results: tree structures under a charity → person → donation schema\n\n")
	fmt.Fprintf(&b, "Generated from `%s`. Every number here comes from exactly one JSON file in that\n", filepath.Base(resultsDir))
	fmt.Fprintf(&b, "directory; the JSON is the source of truth and this file only arranges it.\n\n")

	// -- provenance -----------------------------------------------------------
	fmt.Fprintf(&b, "## What was measured\n\n")
	fmt.Fprintf(&b, "| | |\n|---|---|\n")
	fmt.Fprintf(&b, "| Run id | `%s` |\n", ref.RunID)
	fmt.Fprintf(&b, "| Environment | [`%s`](../../../docs/environments/%s.md) |\n", ref.Environment, ref.Environment)
	fmt.Fprintf(&b, "| Dataset scale | `%s` — %v charities, %v people, %v donations |\n",
		ref.Scale, ref.Dataset["charities"], ref.Dataset["people"], ref.Dataset["donations"])
	fmt.Fprintf(&b, "| Donations per person | mean %.1f, max %v |\n",
		toF(ref.Dataset["avg_donations_person"]), ref.Dataset["max_donations_person"])
	fmt.Fprintf(&b, "| Dataset seed | %v (identical data in every cell) |\n", ref.Dataset["seed"])
	fmt.Fprintf(&b, "| Client connections | %v |\n", ref.Options["conns"])
	fmt.Fprintf(&b, "| Duration per query | %v measured, %v warmup discarded |\n", ref.Options["duration"], ref.Options["warmup"])
	fmt.Fprintf(&b, "\n**Engines**\n\n| Topology | Version |\n|---|---|\n")
	for _, t := range rs.topologies {
		for _, d := range rs.designsIn(t) {
			fmt.Fprintf(&b, "| %s | `%s` |\n", topologyLabel[t], firstLine(rs.runs[t][d].EngineVer))
			break
		}
	}

	// -- correctness ----------------------------------------------------------
	fmt.Fprintf(&b, "\n## Correctness gate\n\n")
	fmt.Fprintf(&b, "Every design must return the same answers, checked against values computed\n")
	fmt.Fprintf(&b, "independently in Go from the generated dataset. Timings are only reported for\n")
	fmt.Fprintf(&b, "cells that passed; a failing cell aborts before it is measured.\n\n")
	fmt.Fprintf(&b, "| Topology | Design | Checks passed |\n|---|---|---|\n")
	allPass := true
	for _, t := range rs.topologies {
		for _, d := range rs.designsIn(t) {
			r := rs.runs[t][d]
			p, f := 0, 0
			if r.Verify != nil {
				p, f = r.Verify.Passed, r.Verify.Failed
			}
			mark := "✅"
			if f > 0 {
				mark, allPass = "❌", false
			}
			fmt.Fprintf(&b, "| %s | %s | %s %d/%d |\n", topologyLabel[t], designShort[d], mark, p, p+f)
		}
	}
	if allPass {
		fmt.Fprintf(&b, "\nAll cells passed.\n")
	}

	// -- the trade-off at a glance --------------------------------------------
	writeTradeoff(&b, rs)

	// -- read throughput ------------------------------------------------------
	fmt.Fprintf(&b, "\n## Read throughput by question (ops/s, higher is better)\n\n")
	for _, t := range rs.topologies {
		ds := rs.designsIn(t)
		fmt.Fprintf(&b, "### %s\n\n", topologyLabel[t])
		fmt.Fprintf(&b, "| Question | ")
		for _, d := range ds {
			fmt.Fprintf(&b, "%s | ", designShort[d])
		}
		fmt.Fprintf(&b, "\n|---|")
		for range ds {
			fmt.Fprintf(&b, "---:|")
		}
		fmt.Fprintf(&b, "\n")
		for _, q := range sortedQueries(rs, t) {
			fmt.Fprintf(&b, "| %s | ", queryQuestion[q])
			for _, d := range ds {
				qr := readOf(rs.runs[t][d], q)
				if qr == nil {
					fmt.Fprintf(&b, "— | ")
					continue
				}
				fmt.Fprintf(&b, "%s | ", fmtOps(qr.OpsPerSec))
			}
			fmt.Fprintf(&b, "\n")
		}
		fmt.Fprintf(&b, "\n")
	}

	// -- tail latency ---------------------------------------------------------
	fmt.Fprintf(&b, "## Read tail latency (p99, ms, lower is better)\n\n")
	fmt.Fprintf(&b, "Means hide the cases a design handles badly; this is where design decisions\n")
	fmt.Fprintf(&b, "usually show up first.\n\n")
	for _, t := range rs.topologies {
		ds := rs.designsIn(t)
		fmt.Fprintf(&b, "### %s\n\n| Question | ", topologyLabel[t])
		for _, d := range ds {
			fmt.Fprintf(&b, "%s | ", designShort[d])
		}
		fmt.Fprintf(&b, "\n|---|")
		for range ds {
			fmt.Fprintf(&b, "---:|")
		}
		fmt.Fprintf(&b, "\n")
		for _, q := range sortedQueries(rs, t) {
			fmt.Fprintf(&b, "| %s | ", queryQuestion[q])
			for _, d := range ds {
				qr := readOf(rs.runs[t][d], q)
				if qr == nil {
					fmt.Fprintf(&b, "— | ")
					continue
				}
				fmt.Fprintf(&b, "%s | ", fmtMS(qr.Latency.P99MS))
			}
			fmt.Fprintf(&b, "\n")
		}
		fmt.Fprintf(&b, "\n")
	}

	// -- writes ---------------------------------------------------------------
	fmt.Fprintf(&b, "## Write throughput (ops/s, higher is better)\n\n")
	fmt.Fprintf(&b, "`insert` appends a donation, `update` corrects one amount, `delete` removes one\n")
	fmt.Fprintf(&b, "donation. Reads are only half the story: every design that made a read cheap\n")
	fmt.Fprintf(&b, "paid for it somewhere in this table.\n\n")
	fmt.Fprintf(&b, "A dash in the p99.9 or p99.99 row means **too few samples to support that\n")
	fmt.Fprintf(&b, "percentile**, not a missing measurement: a quantile is only reported when at\n")
	fmt.Fprintf(&b, "least ten samples sit beyond it (10 000 samples for p99.9, 100 000 for p99.99).\n")
	fmt.Fprintf(&b, "Below that a \"p99.9\" would be the single worst request relabelled, which reads\n")
	fmt.Fprintf(&b, "as far more authoritative than it is. `max` is always shown, and for a\n")
	fmt.Fprintf(&b, "low-throughput cell it is the only honest tail figure available.\n\n")
	fmt.Fprintf(&b, "> **These tails are optimistic, by construction.** The harness is closed-loop:\n")
	fmt.Fprintf(&b, "> each worker waits for its own response before issuing the next request. When\n")
	fmt.Fprintf(&b, "> the server stalls, the load generator stalls with it, so the stall is recorded\n")
	fmt.Fprintf(&b, "> once instead of being charged to every request that would have arrived during\n")
	fmt.Fprintf(&b, "> it — the coordinated-omission effect. It barely touches p50, biases p99 a\n")
	fmt.Fprintf(&b, "> little, and biases p99.9 and beyond substantially. Treat deep tails here as a\n")
	fmt.Fprintf(&b, "> floor on what a real open-loop client would see, and compare them **between\n")
	fmt.Fprintf(&b, "> designs** (all of which pay the same bias) rather than as absolute SLO figures.\n\n")
	for _, t := range rs.topologies {
		ds := rs.designsIn(t)
		fmt.Fprintf(&b, "### %s\n\n| Operation | ", topologyLabel[t])
		for _, d := range ds {
			fmt.Fprintf(&b, "%s | ", designShort[d])
		}
		fmt.Fprintf(&b, "\n|---|")
		for range ds {
			fmt.Fprintf(&b, "---:|")
		}
		fmt.Fprintf(&b, "\n")
		for _, op := range writeOpOrder {
			fmt.Fprintf(&b, "| %s | ", writeOpLabel[op])
			for _, d := range ds {
				wr := writeOf(rs.runs[t][d], op)
				if wr == nil {
					fmt.Fprintf(&b, "— | ")
					continue
				}
				fmt.Fprintf(&b, "%s | ", fmtOps(wr.OpsPerSec))
			}
			fmt.Fprintf(&b, "\n")
		}
		// The insert tail, progressively deeper. Writes are where contention
		// lives, so this is the path whose tail is worth following furthest.
		for _, tail := range []struct {
			label string
			pick  func(LatencyStats) float64
		}{
			{"insert p99 (ms)", func(l LatencyStats) float64 { return l.P99MS }},
			{"insert p99.9 (ms)", func(l LatencyStats) float64 { return l.P999MS }},
			{"insert p99.99 (ms)", func(l LatencyStats) float64 { return l.P9999MS }},
			{"insert max (ms)", func(l LatencyStats) float64 { return l.MaxMS }},
		} {
			fmt.Fprintf(&b, "| %s | ", tail.label)
			for _, d := range ds {
				wr := writeOf(rs.runs[t][d], "insert")
				if wr == nil {
					fmt.Fprintf(&b, "— | ")
					continue
				}
				fmt.Fprintf(&b, "%s | ", fmtMS(tail.pick(wr.Latency)))
			}
			fmt.Fprintf(&b, "\n")
		}
		fmt.Fprintf(&b, "\n")
	}

	// -- storage and load -----------------------------------------------------
	fmt.Fprintf(&b, "## Storage footprint and load time\n\n")
	fmt.Fprintf(&b, "Read speed bought with redundancy is paid for in bytes and in how long the\n")
	fmt.Fprintf(&b, "data takes to land. Sizes are PostgreSQL only — YugabyteDB stores data in\n")
	fmt.Fprintf(&b, "DocDB rather than PostgreSQL heap files, so `pg_total_relation_size()` does not\n")
	fmt.Fprintf(&b, "describe it.\n\n")
	for _, t := range rs.topologies {
		ds := rs.designsIn(t)
		fmt.Fprintf(&b, "### %s\n\n", topologyLabel[t])
		fmt.Fprintf(&b, "| Design | Tables+indexes | Index bytes | Indexes on `donation` | Load total | COPY | Index build | Rollup |\n")
		fmt.Fprintf(&b, "|---|---:|---:|---:|---:|---:|---:|---:|\n")
		for _, d := range ds {
			r := rs.runs[t][d]
			size, idxBytes, nIdx := "—", "—", "—"
			if r.Stats != nil {
				if r.Stats.TotalBytes > 0 {
					size = fmtBytes(r.Stats.TotalBytes)
				}
				var ib int64
				n := 0
				for _, rel := range r.Stats.Relations {
					if rel.Kind != "index" {
						continue
					}
					ib += rel.TotalBytes
					if rel.Table == "donation" {
						n++
					}
				}
				if ib > 0 {
					idxBytes = fmtBytes(ib)
				}
				if len(r.Stats.Relations) > 0 {
					nIdx = fmt.Sprint(n)
				}
			}
			if r.Load == nil {
				fmt.Fprintf(&b, "| %s | %s | %s | %s | — | — | — | — |\n", designShort[d], size, idxBytes, nIdx)
				continue
			}
			fmt.Fprintf(&b, "| %s | %s | %s | %s | %.1fs | %.1fs | %.1fs | %.1fs |\n",
				designShort[d], size, idxBytes, nIdx, r.Load.TotalMS/1000, r.Load.CopyMS/1000,
				r.Load.IndexMS/1000, r.Load.RollupMS/1000)
		}
		fmt.Fprintf(&b, "\n")
	}
	fmt.Fprintf(&b, "**Indexes on `donation`** counts the index entries that must be maintained on every\n")
	fmt.Fprintf(&b, "donation INSERT — including the primary key. It is the structural driver of write\n")
	fmt.Fprintf(&b, "cost, and the number to look at when a write result surprises you.\n\n")

	// -- controlled pairs -----------------------------------------------------
	writePairs(&b, rs)

	fmt.Fprintf(&b, "\n---\n\n")
	fmt.Fprintf(&b, "Query plans for every cell are in `%s/<topology>/plans/<design>.txt`,\n", filepath.Base(resultsDir))
	fmt.Fprintf(&b, "captured with `EXPLAIN (ANALYZE, BUFFERS)` on PostgreSQL and\n")
	fmt.Fprintf(&b, "`EXPLAIN (ANALYZE, DIST)` on YugabyteDB.\n")

	writeProvenance(&b, resultsDir, filepath.Join(filepath.Dir(outPath), "analyses"), ref.RunID)

	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(outPath, []byte(b.String()), 0o644)
}

// geomeanReads summarises a design's read performance in one number.
//
// Geometric, not arithmetic: these throughputs span four orders of magnitude
// (35 ops/s to 40 000 ops/s in the same column), and an arithmetic mean would be
// entirely determined by whichever query happens to be fastest, telling you
// nothing about the eleven others. The geometric mean treats a 10x gain on a
// slow query and a 10x gain on a fast one as equally valuable, which is the right
// weighting when you do not yet know the workload mix.
func geomeanReads(r *Run) float64 {
	if r == nil || len(r.Reads) == 0 {
		return 0
	}
	sum, n := 0.0, 0
	for i := range r.Reads {
		if v := r.Reads[i].OpsPerSec; v > 0 {
			sum += math.Log(v)
			n++
		}
	}
	if n == 0 {
		return 0
	}
	return math.Exp(sum / float64(n))
}

// writeTradeoff is the one table to read if you read only one. Every design in
// this study buys read speed with write cost, storage, or both; putting all
// three side by side is the only way to see what a design actually costs.
func writeTradeoff(b *strings.Builder, rs *resultSet) {
	fmt.Fprintf(b, "\n## The trade-off at a glance\n\n")
	fmt.Fprintf(b, "Every design here buys read speed with something — write throughput, bytes on\n")
	fmt.Fprintf(b, "disk, or index maintenance on the hot path. This table puts the purchase and the\n")
	fmt.Fprintf(b, "price in the same row.\n\n")
	fmt.Fprintf(b, "*Read score* is the **geometric** mean across all twelve questions. These\n")
	fmt.Fprintf(b, "throughputs span four orders of magnitude, so an arithmetic mean would be decided\n")
	fmt.Fprintf(b, "entirely by the fastest query and say nothing about the other eleven.\n\n")

	for _, t := range rs.topologies {
		ds := rs.designsIn(t)
		if len(ds) == 0 {
			continue
		}
		base := geomeanReads(rs.runs[t][ds[0]])
		baseIns := 0.0
		if w := writeOf(rs.runs[t][ds[0]], "insert"); w != nil {
			baseIns = w.OpsPerSec
		}

		fmt.Fprintf(b, "### %s\n\n", topologyLabel[t])
		fmt.Fprintf(b, "| Design | Read score | vs %s | Insert/s | vs %s | Indexes on `donation` | Storage |\n",
			designShort[ds[0]], designShort[ds[0]])
		fmt.Fprintf(b, "|---|---:|---:|---:|---:|---:|---:|\n")
		for _, d := range ds {
			r := rs.runs[t][d]
			g := geomeanReads(r)
			ins := 0.0
			if w := writeOf(r, "insert"); w != nil {
				ins = w.OpsPerSec
			}
			relRead, relWrite := "—", "—"
			if base > 0 && g > 0 {
				relRead = fmt.Sprintf("%.1fx", g/base)
			}
			if baseIns > 0 && ins > 0 {
				relWrite = fmt.Sprintf("%.2fx", ins/baseIns)
			}
			nIdx, size := "—", "—"
			if r.Stats != nil {
				n := 0
				for _, rel := range r.Stats.Relations {
					if rel.Kind == "index" && rel.Table == "donation" {
						n++
					}
				}
				if len(r.Stats.Relations) > 0 {
					nIdx = fmt.Sprint(n)
				}
				if r.Stats.TotalBytes > 0 {
					size = fmtBytes(r.Stats.TotalBytes)
				}
			}
			fmt.Fprintf(b, "| %s | %s | %s | %s | %s | %s | %s |\n",
				designShort[d], fmtOps(g), relRead, fmtOps(ins), relWrite, nIdx, size)
		}
		fmt.Fprintf(b, "\n")
	}
}

// pair is a comparison between two designs that differ in exactly one decision.
type pair struct{ a, bb, question string }

var pairs = []pair{
	{"d2_normalized_indexed", "d11_copied_key", "Copied key and its foreign key, before charity access paths"},
	{"d11_copied_key", "d12_recency_index", "Charity/date index with identical read SQL"},
	{"d12_recency_index", "d13_recency_sql", "Recency SQL rewrites with identical indexes"},
	{"d13_recency_sql", "d17_sum_sql", "Sum SQL rewrite with identical indexes"},
	{"d17_sum_sql", "d14_sum_plain", "Plain charity index with identical SQL"},
	{"d14_sum_plain", "d15_sum_covering", "INCLUDE amount payload with identical SQL"},
	{"d15_sum_covering", "d3_flattened_fk", "Ranking SQL rewrites with identical indexes"},
	{"d3_flattened_fk", "d16_sum_rollup", "Synchronous sum maintenance on both parents"},
	{"d16_sum_rollup", "d4_rollup_trigger", "Full rollup package versus sums only"},
	{"d1_normalized_minimal", "d2_normalized_indexed", "What do secondary indexes alone buy? (identical SQL)"},
	{"d2_normalized_indexed", "d3_flattened_fk", "What does denormalising the grandparent key buy?"},
	{"d3_flattened_fk", "d8_flattened_nofk", "What does enforcing foreign keys cost?"},
	{"d3_flattened_fk", "d4_rollup_trigger", "What do consolidated aggregates buy and cost?"},
	{"d4_rollup_trigger", "d5_rollup_app", "Does it matter whether a trigger or the app maintains them?"},
	{"d3_flattened_fk", "d6_embedded_jsonb", "What does embedding ALL children in the parent buy and cost?"},
	{"d6_embedded_jsonb", "d9_embedded_hybrid", "Does bounding the embedded slice fix what embedding broke?"},
	{"d3_flattened_fk", "d9_embedded_hybrid", "What does a bounded embedded cache buy and cost?"},
	{"d9_embedded_hybrid", "d10_embedded_hybrid_locked", "What does making the cache trigger concurrency-correct cost?"},
	{"d3_flattened_fk", "d7_yb_child_colocated", "What does sharding children next to their parent buy?"},
}

// noiseFloor measures this run's own error bar, empirically, from the D4/D5
// control pair.
//
// D4 and D5 store the same aggregates in the same columns and read them with
// byte-identical SQL. Every read difference between them is therefore
// measurement noise by construction -- cell-to-cell drift, thermal throttling,
// an autovacuum pass, the scheduler moving a container onto an efficiency core.
//
// Using that to set the significance threshold is far better than picking a
// round number: the threshold is derived from THIS run on THIS machine, and it
// rises automatically when the machine was noisy. It returns the worst observed
// ratio, which is deliberately conservative -- a pair difference smaller than
// the largest error the control exhibited cannot be argued from.
func noiseFloor(rs *resultSet) (worst, typical float64, n int) {
	var ratios []float64
	for _, t := range rs.topologies {
		a, bb := rs.runs[t]["d4_rollup_trigger"], rs.runs[t]["d5_rollup_app"]
		if a == nil || bb == nil {
			continue
		}
		for _, q := range sortedQueries(rs, t) {
			qa, qb := readOf(a, q), readOf(bb, q)
			if qa == nil || qb == nil || qa.OpsPerSec <= 0 || qb.OpsPerSec <= 0 {
				continue
			}
			r := qa.OpsPerSec / qb.OpsPerSec
			if r < 1 {
				r = 1 / r
			}
			ratios = append(ratios, r)
		}
	}
	if len(ratios) == 0 {
		return 1.3, 1.3, 0 // fall back to a conventional threshold
	}
	sort.Float64s(ratios)
	return ratios[len(ratios)-1], ratios[len(ratios)/2], len(ratios)
}

func writePairs(b *strings.Builder, rs *resultSet) {
	worst, typical, n := noiseFloor(rs)

	fmt.Fprintf(b, "## Controlled pairs — one decision at a time\n\n")
	fmt.Fprintf(b, "Each pair below differs by exactly one design decision, so a measured\n")
	fmt.Fprintf(b, "difference has exactly one explanation. This is the part of the report worth\n")
	fmt.Fprintf(b, "reading if you read nothing else.\n\n")

	fmt.Fprintf(b, "### This run's error bar\n\n")
	if n > 0 {
		fmt.Fprintf(b, "**D4 and D5 read identically by construction** — same columns, byte-identical\n")
		fmt.Fprintf(b, "query SQL. Every read difference measured between them is therefore noise, which\n")
		fmt.Fprintf(b, "makes that pair a free calibration of this run's own error bar:\n\n")
		fmt.Fprintf(b, "| | |\n|---|---|\n")
		fmt.Fprintf(b, "| Control comparisons | %d |\n", n)
		fmt.Fprintf(b, "| Typical (median) disagreement | **%.2fx** |\n", typical)
		fmt.Fprintf(b, "| Worst disagreement | **%.2fx** |\n", worst)
		fmt.Fprintf(b, "\nDifferences below the worst control disagreement (**%.2fx**) are suppressed in the\n", worst)
		fmt.Fprintf(b, "tables that follow. A gap smaller than the largest error two *identical* designs\n")
		fmt.Fprintf(b, "showed cannot be attributed to a design decision.\n\n")
		if worst > 1.4 {
			fmt.Fprintf(b, "> ⚠️ A worst-case control disagreement of %.2fx is high. It means this run can\n", worst)
			fmt.Fprintf(b, "> only support order-of-magnitude claims, not tens-of-percent ones. Comparisons\n")
			fmt.Fprintf(b, "> that matter should be re-measured with `-trials N`, which reports a median\n")
			fmt.Fprintf(b, "> across repeats and a per-cell spread.\n\n")
		}
	} else {
		fmt.Fprintf(b, "The D4/D5 control pair is not present in this run, so the threshold below is a\n")
		fmt.Fprintf(b, "conventional %.2fx rather than one measured from the data.\n\n", worst)
	}

	for _, p := range pairs {
		var rows []string
		for _, t := range rs.topologies {
			ra, rb := rs.runs[t][p.a], rs.runs[t][p.bb]
			if ra == nil || rb == nil {
				continue
			}
			for _, q := range sortedQueries(rs, t) {
				qa, qb := readOf(ra, q), readOf(rb, q)
				if qa == nil || qb == nil || qa.OpsPerSec <= 0 || qb.OpsPerSec <= 0 {
					continue
				}
				// Only surface differences big enough to be a signal rather than
				// scheduling noise on a shared laptop.
				r := qb.OpsPerSec / qa.OpsPerSec
				if r > 1/worst && r < worst {
					continue
				}
				rows = append(rows, fmt.Sprintf("| %s | %s | %s | %s | **%s** |",
					topologyLabel[t], queryQuestion[q], fmtOps(qa.OpsPerSec), fmtOps(qb.OpsPerSec),
					ratio(qa.OpsPerSec, qb.OpsPerSec)))
			}
			// Writes are ALWAYS listed, even when the difference is below the
			// error bar. Every design in this study buys read speed with write
			// cost, so the write row is the price tag. Suppressing it as "not
			// significant" would turn the most important half of a trade-off into
			// silence, and a reader would be left to assume it was never measured.
			for _, op := range writeOpOrder {
				wa, wb := writeOf(ra, op), writeOf(rb, op)
				if wa == nil || wb == nil {
					continue
				}
				if wa.Skipped != "" || wb.Skipped != "" {
					rows = append(rows, fmt.Sprintf("| %s | %s | — | — | *not measured — %s* |",
						topologyLabel[t], writeOpLabel[op], firstSentence(wa.Skipped+wb.Skipped)))
					continue
				}
				if wa.OpsPerSec <= 0 || wb.OpsPerSec <= 0 {
					continue
				}
				r := wb.OpsPerSec / wa.OpsPerSec
				change := fmt.Sprintf("**%s**", ratio(wa.OpsPerSec, wb.OpsPerSec))
				if r > 1/worst && r < worst {
					change = fmt.Sprintf("no measurable change (%.2fx, inside the %.2fx error bar)", r, worst)
				}
				rows = append(rows, fmt.Sprintf("| %s | %s | %s | %s | %s |",
					topologyLabel[t], writeOpLabel[op], fmtOps(wa.OpsPerSec), fmtOps(wb.OpsPerSec), change))
			}
		}
		if len(rows) == 0 {
			continue
		}
		fmt.Fprintf(b, "### %s → %s\n\n*%s*\n\n", designShort[p.a], designShort[p.bb], p.question)
		fmt.Fprintf(b, "| Topology | Operation | %s | %s | Change |\n|---|---|---:|---:|---|\n",
			designShort[p.a], designShort[p.bb])
		for _, r := range rows {
			fmt.Fprintln(b, r)
		}
		fmt.Fprintf(b, "\nRead differences below this run's measured %.2fx error bar are omitted. **Write\n", worst)
		fmt.Fprintf(b, "rows are always shown**, even when the change is inside the error bar — writes are\n")
		fmt.Fprintf(b, "the price every one of these designs pays for its read speed, and a hidden price\n")
		fmt.Fprintf(b, "tag reads as no price at all.\n\n")
	}
}

func sortedQueries(rs *resultSet, topo string) []string {
	seen := map[string]bool{}
	var out []string
	for _, d := range rs.designsIn(topo) {
		for _, q := range rs.runs[topo][d].Reads {
			if !seen[q.Query] {
				seen[q.Query] = true
				out = append(out, q.Query)
			}
		}
	}
	sort.Strings(out)
	return out
}

// firstSentence trims a long explanation down to something that fits in a table
// cell without losing the reason.
func firstSentence(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.IndexByte(s, ';'); i > 0 {
		return s[:i]
	}
	if len(s) > 120 {
		return s[:117] + "..."
	}
	return s
}

func toF(v any) float64 {
	if f, ok := v.(float64); ok {
		return f
	}
	return 0
}

// writeOpOrder fixes how write operations are listed, grouped by where the write
// originates. The child-originating ops come first because they are the hot path;
// the parent-originating ops follow because they are rarer but are where the
// embedded and rollup designs hide their worst cases.
var writeOpOrder = []string{"insert", "update", "delete", "update_person", "delete_person"}

// writeOpLabel names each operation in terms of what a user did, not what SQL ran.
var writeOpLabel = map[string]string{
	"insert":        "insert donation",
	"update":        "correct a donation amount",
	"delete":        "remove one donation",
	"update_person": "donor edits their profile",
	"delete_person": "erase a donor and all their donations",
}
