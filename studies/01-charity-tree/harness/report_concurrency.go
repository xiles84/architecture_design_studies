package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ---------------------------------------------------------------------------
// Experiment C report — concurrency control on the hot rollup row.
//
// Cells are identified by what they configured (design, strategy, isolation,
// connections), read back out of each result file, rather than by filename.
// Filenames are a convenience; the run record is the fact.
// ---------------------------------------------------------------------------

type ccell struct {
	run       *Run
	design    string
	strategy  string
	isolation string
	conns     int
}

var strategyBlurb = map[string]string{
	"blind":      "plain UPDATE; the engine's row lock decides the order",
	"forupdate":  "SELECT ... FOR UPDATE first — pessimistic, lock then act",
	"optimistic": "read version, UPDATE ... WHERE version = n, retry on loss",
	"—":          "maintained by an AFTER trigger inside the same statement",
}

var isolationName = map[string]string{
	"rc": "read committed", "rr": "repeatable read", "ser": "serializable",
}

func loadConcurrency(dir string) ([]ccell, error) {
	var out []ccell
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
		if len(r.Writes) == 0 {
			return nil
		}
		c := ccell{
			run:       &r,
			design:    r.DesignID,
			strategy:  optString(r.Options, "strategy"),
			isolation: optString(r.Options, "isolation"),
			conns:     int(optFloat(r.Options, "conns")),
		}
		// Only D5 lets the strategy vary; for every other design the label would
		// be misleading, since nothing in the cell honours it.
		if r.DesignID != "d5_rollup_app" {
			c.strategy = "—"
		}
		out = append(out, c)
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no write results found under %s", dir)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].conns != out[j].conns {
			return out[i].conns < out[j].conns
		}
		if out[i].design != out[j].design {
			return out[i].design < out[j].design
		}
		if out[i].strategy != out[j].strategy {
			return out[i].strategy < out[j].strategy
		}
		return out[i].isolation < out[j].isolation
	})
	return out, nil
}

func optString(m map[string]any, k string) string {
	if v, ok := m[k].(string); ok {
		return v
	}
	return ""
}

func optFloat(m map[string]any, k string) float64 {
	if v, ok := m[k].(float64); ok {
		return v
	}
	return 0
}

func writeConcurrencyReport(dir, outPath string) error {
	cells, err := loadConcurrency(dir)
	if err != nil {
		return err
	}
	var b strings.Builder
	ref := cells[0].run

	fmt.Fprintf(&b, "# Experiment C — concurrency control on the hot rollup row\n\n")
	fmt.Fprintf(&b, "| | |\n|---|---|\n")
	fmt.Fprintf(&b, "| Run id | `%s` |\n", ref.RunID)
	fmt.Fprintf(&b, "| Environment | [`%s`](../../../docs/environments/%s.md) |\n", ref.Environment, ref.Environment)
	fmt.Fprintf(&b, "| Topology | %s |\n", topologyLabel[ref.Topology])
	fmt.Fprintf(&b, "| Engine | `%s` |\n", firstLine(ref.EngineVer))
	fmt.Fprintf(&b, "| Scale | `%s` — %v people, %v donations |\n", ref.Scale, ref.Dataset["people"], ref.Dataset["donations"])
	fmt.Fprintf(&b, "| Operation | `insert` — append one donation and update both parent rollups |\n")

	fmt.Fprintf(&b, "\n## The setup\n\n")
	fmt.Fprintf(&b, "Ten charities; every donation in the system updates one of their rows. The largest\n")
	fmt.Fprintf(&b, "charity holds 30%% of all donors, so under load its row is a point every writer\n")
	fmt.Fprintf(&b, "must pass through. The contention is structural, not contrived — it is what\n")
	fmt.Fprintf(&b, "consolidating an aggregate onto a parent *means*.\n\n")
	fmt.Fprintf(&b, "| Row | What maintains the aggregates |\n|---|---|\n")
	fmt.Fprintf(&b, "| `D3 (no rollup)` | nothing — the ceiling. What the insert costs if no aggregate exists. |\n")
	fmt.Fprintf(&b, "| `D4 trigger` | an AFTER trigger, inside the same statement |\n")
	for _, s := range []string{"blind", "forupdate", "optimistic"} {
		fmt.Fprintf(&b, "| `D5 %s` | %s |\n", s, strategyBlurb[s])
	}

	// -- contention curve -----------------------------------------------------
	fmt.Fprintf(&b, "\n## Throughput as writers are added (inserts/s)\n\n")
	fmt.Fprintf(&b, "A strategy that wins at 4 writers and collapses at 32 has not solved anything.\n\n")

	connSet := map[int]bool{}
	variantSet := map[string]bool{}
	byKey := map[string]*ccell{}
	for i := range cells {
		c := &cells[i]
		if c.isolation != "rc" {
			continue
		}
		v := variantLabel(c)
		connSet[c.conns] = true
		variantSet[v] = true
		byKey[fmt.Sprintf("%s|%d", v, c.conns)] = c
	}
	conns := sortedInts(connSet)
	variants := orderVariants(variantSet)

	fmt.Fprintf(&b, "| Writers | ")
	for _, v := range variants {
		fmt.Fprintf(&b, "%s | ", v)
	}
	fmt.Fprintf(&b, "\n|---:|")
	for range variants {
		fmt.Fprintf(&b, "---:|")
	}
	fmt.Fprintf(&b, "\n")
	for _, n := range conns {
		fmt.Fprintf(&b, "| %d | ", n)
		for _, v := range variants {
			c := byKey[fmt.Sprintf("%s|%d", v, n)]
			if c == nil {
				fmt.Fprintf(&b, "— | ")
				continue
			}
			fmt.Fprintf(&b, "%s | ", fmtOps(c.run.Writes[0].OpsPerSec))
		}
		fmt.Fprintf(&b, "\n")
	}

	// Contention shows up in the tail long before it shows up in the median: a
	// strategy whose throughput is holding while its p99.9 climbs is already
	// failing for some fraction of users.
	for _, tail := range []struct {
		name  string
		title string
		pick  func(LatencyStats) float64
	}{
		{"p50", "median latency (ms)", func(l LatencyStats) float64 { return l.P50MS }},
		{"p99", "p99 latency (ms)", func(l LatencyStats) float64 { return l.P99MS }},
		{"p999", "p99.9 latency (ms)", func(l LatencyStats) float64 { return l.P999MS }},
		{"max", "worst observed (ms)", func(l LatencyStats) float64 { return l.MaxMS }},
	} {
		fmt.Fprintf(&b, "\n### %s\n\n| Writers | ", tail.title)
		for _, v := range variants {
			fmt.Fprintf(&b, "%s | ", v)
		}
		fmt.Fprintf(&b, "\n|---:|")
		for range variants {
			fmt.Fprintf(&b, "---:|")
		}
		fmt.Fprintf(&b, "\n")
		for _, n := range conns {
			fmt.Fprintf(&b, "| %d | ", n)
			for _, v := range variants {
				c := byKey[fmt.Sprintf("%s|%d", v, n)]
				if c == nil {
					fmt.Fprintf(&b, "— | ")
					continue
				}
				fmt.Fprintf(&b, "%s | ", fmtMS(tail.pick(c.run.Writes[0].Latency)))
			}
			fmt.Fprintf(&b, "\n")
		}
		if tail.name == "p999" {
			fmt.Fprintf(&b, "\nA dash means fewer than 10 000 samples in that cell — too few for p99.9 to be\n")
			fmt.Fprintf(&b, "anything but the worst request relabelled. Use the `worst observed` row instead.\n")
		}
	}
	fmt.Fprintf(&b, "\n> Tails here are closed-loop and therefore optimistic: a stalled server stalls\n")
	fmt.Fprintf(&b, "> the load generator too, so the stall is counted once rather than charged to\n")
	fmt.Fprintf(&b, "> every request it would have delayed. Compare across strategies, which all pay\n")
	fmt.Fprintf(&b, "> the same bias; do not read these as SLO numbers.\n")

	// -- isolation ------------------------------------------------------------
	fmt.Fprintf(&b, "\n## Isolation level at fixed concurrency\n\n")
	fmt.Fprintf(&b, "Held at 8 writers so the isolation level is the only variable. Stricter isolation\n")
	fmt.Fprintf(&b, "is expected to cost *retries* rather than latency, which is why retries are counted\n")
	fmt.Fprintf(&b, "here rather than hidden inside the throughput number.\n\n")
	fmt.Fprintf(&b, "| Strategy | Isolation | inserts/s | p99 ms | Retries | Errors |\n")
	fmt.Fprintf(&b, "|---|---|---:|---:|---:|---:|\n")
	for i := range cells {
		c := &cells[i]
		if c.conns != 8 || c.design != "d5_rollup_app" {
			continue
		}
		w := c.run.Writes[0]
		fmt.Fprintf(&b, "| %s | %s | %s | %s | %d | %d |\n",
			c.strategy, isolationName[c.isolation], fmtOps(w.OpsPerSec),
			fmtMS(w.Latency.P99MS), w.Retries, w.Errors)
	}

	// -- correctness ----------------------------------------------------------
	fmt.Fprintf(&b, "\n## Did the aggregates survive? (the result that matters)\n\n")
	fmt.Fprintf(&b, "After each contended run, every rollup is recomputed from the donation table and\n")
	fmt.Fprintf(&b, "compared. **A strategy that is faster and silently loses increments has not won\n")
	fmt.Fprintf(&b, "anything** — it has published a total that is wrong in a way nothing in the system\n")
	fmt.Fprintf(&b, "would ever notice. `drift` is what the rollups claim was donated minus what the\n")
	fmt.Fprintf(&b, "donation rows actually add up to.\n\n")
	fmt.Fprintf(&b, "| Writers | Variant | Isolation | Person rows wrong | Charity rows wrong | Drift (cents) |\n")
	fmt.Fprintf(&b, "|---:|---|---|---:|---:|---:|\n")
	anyDrift := false
	for i := range cells {
		c := &cells[i]
		a := c.run.Audit
		if a == nil || !a.Ran {
			continue
		}
		mark := ""
		if a.PersonMismatches+a.CharityMismatches > 0 {
			mark, anyDrift = " ❌", true
		}
		fmt.Fprintf(&b, "| %d | %s | %s | %d | %d | %d%s |\n",
			c.conns, variantLabel(c), isolationName[c.isolation],
			a.PersonMismatches, a.CharityMismatches, a.DriftCents, mark)
	}
	if !anyDrift {
		fmt.Fprintf(&b, "\nNo strategy lost an update in this run. That is a result about *these* strategies\n")
		fmt.Fprintf(&b, "on *this* engine, not a general guarantee: every one of them performs its\n")
		fmt.Fprintf(&b, "increment inside a single UPDATE statement, which the engine executes atomically\n")
		fmt.Fprintf(&b, "under its own row lock. A read-modify-write split across two statements — the\n")
		fmt.Fprintf(&b, "shape most application code reaches for first — has no such protection.\n")
	}

	fmt.Fprintf(&b, "\n---\n\nRaw results: `%s/`.\n", filepath.Base(dir))
	writeProvenance(&b, dir, filepath.Join(filepath.Dir(outPath), "analyses"), ref.RunID)
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(outPath, []byte(b.String()), 0o644)
}

func variantLabel(c *ccell) string {
	switch c.design {
	case "d3_flattened_fk":
		return "D3 no rollup"
	case "d4_rollup_trigger":
		return "D4 trigger"
	case "d5_rollup_app":
		return "D5 " + c.strategy
	}
	return designShort[c.design]
}

func orderVariants(set map[string]bool) []string {
	pref := []string{"D3 no rollup", "D4 trigger", "D5 blind", "D5 forupdate", "D5 optimistic"}
	var out []string
	for _, p := range pref {
		if set[p] {
			out = append(out, p)
			delete(set, p)
		}
	}
	var rest []string
	for k := range set {
		rest = append(rest, k)
	}
	sort.Strings(rest)
	return append(out, rest...)
}

func sortedInts(set map[int]bool) []int {
	var out []int
	for k := range set {
		out = append(out, k)
	}
	sort.Ints(out)
	return out
}
