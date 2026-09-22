package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"adsplatform/adapters/filestore"
	"adsplatform/adapters/markdown"
	"adsplatform/core/provenance"
)

// A generated report holds measurements and no interpretation (methodology 11).
// Its TL;DR lists measured facts selected by fixed rules -- which cells failed,
// whether the negative controls fired, which designs broke an invariant, the
// highest and lowest valid result per condition. It ranks by the number and says
// nothing about which design is better; what the numbers mean is written
// separately, in a signed analysis under reports/analyses/.

func writeReport(resultsDir, out string) error {
	if resultsDir == "" {
		return fmt.Errorf("report: -results is required")
	}
	if out == "" {
		return fmt.Errorf("report: -report-out is required")
	}
	cells, err := filestore.ReadJSONTree[Result](resultsDir)
	if err != nil {
		return fmt.Errorf("read results: %w", err)
	}
	if len(cells) == 0 {
		return fmt.Errorf("no results under %s", resultsDir)
	}
	digest, fileCount, err := provenance.Digest(os.DirFS(resultsDir))
	if err != nil {
		return fmt.Errorf("inputs digest: %w", err)
	}

	sort.Slice(cells, func(i, j int) bool {
		if cells[i].Topology != cells[j].Topology {
			return cells[i].Topology < cells[j].Topology
		}
		return cells[i].DesignShort < cells[j].DesignShort
	})

	analysesDir := filepath.Join(resultsDir, "..", "..", "reports", "analyses")
	var analyses []provenance.Analysis
	if st, err := os.Stat(analysesDir); err == nil && st.IsDir() {
		if a, err := provenance.LoadAnalyses(os.DirFS(analysesDir), cells[0].RunID); err == nil {
			analyses = a
		}
	}

	var b strings.Builder
	runID := cells[0].RunID

	fmt.Fprintf(&b, "# Study 04 — configuration portal — run `%s`\n\n", runID)
	fmt.Fprintf(&b, "Generated %s from %d result file(s). This report contains measurements only.\n\n",
		time.Now().UTC().Format(time.RFC3339), fileCount)

	b.WriteString("| Field | Value |\n|---|---|\n")
	fmt.Fprintf(&b, "| Study | %s |\n", cells[0].Study)
	fmt.Fprintf(&b, "| Run id | `%s` |\n", runID)
	fmt.Fprintf(&b, "| Environment | `%s` |\n", cells[0].Environment)
	fmt.Fprintf(&b, "| Repository commit | `%s` |\n", cells[0].RepoCommit)
	fmt.Fprintf(&b, "| `git describe` | `%s` |\n", cells[0].RepoDesc)
	fmt.Fprintf(&b, "| Working tree dirty at start | %t |\n", cells[0].RepoDirty)
	if cells[0].RunTag != "" {
		fmt.Fprintf(&b, "| Run tag | `%s` |\n", cells[0].RunTag)
	}
	fmt.Fprintf(&b, "| **Inputs digest** | `%s` |\n", digest)
	fmt.Fprintf(&b, "| Cells | %d |\n\n", len(cells))

	// ---------------------------------------------------------------- TL;DR
	b.WriteString("## TL;DR — measured facts, selected by fixed rules\n\n")
	failed := []string{}
	controlsFired, controlsMissed := []string{}, []string{}
	invariantBreakers := []string{}
	for _, c := range cells {
		if c.Error != "" {
			failed = append(failed, fmt.Sprintf("`%s/%s` (%s)", c.Topology, c.Design, firstLine(c.Error)))
		}
		// A negative control breaking an invariant is the control doing its job,
		// reported in the control line above. Listing it here as well would read
		// as a design having broken the contract -- and an audit that runs in more
		// than one phase would list the same failure once per phase.
		if c.Gate.Passed && !c.Negative {
			seen := map[string]bool{}
			for _, a := range c.Audits {
				if a.Passed {
					continue
				}
				for _, f := range a.Failures {
					if !strings.Contains(f, "INV-") || seen[f] {
						continue
					}
					seen[f] = true
					invariantBreakers = append(invariantBreakers,
						fmt.Sprintf("`%s/%s`: %s", c.Topology, c.Design, f))
				}
			}
		}
		if c.Negative {
			fired := false
			for _, k := range c.Concurrency {
				if k.LostUpdates > 0 {
					fired = true
				}
			}
			for _, a := range c.Audits {
				for name, n := range a.DesignChecks {
					if n > 0 {
						fired = true
						_ = name
					}
				}
			}
			tag := fmt.Sprintf("`%s/%s`", c.Topology, c.Design)
			if fired {
				controlsFired = append(controlsFired, tag)
			} else {
				controlsMissed = append(controlsMissed, tag)
			}
		}
	}
	fmt.Fprintf(&b, "- Cells reported: %d; cells that failed: %d.\n", len(cells), len(failed))
	for _, f := range failed {
		fmt.Fprintf(&b, "  - failed: %s\n", f)
	}
	fmt.Fprintf(&b, "- Negative controls that fired: %d. Controls that did NOT fire: %d.\n",
		len(controlsFired), len(controlsMissed))
	for _, c := range controlsMissed {
		fmt.Fprintf(&b, "  - did not fire: %s — the associated correctness claim is **not demonstrated**\n", c)
	}
	if len(invariantBreakers) == 0 {
		b.WriteString("- No invariant violation was recorded by a design that is not a negative control.\n")
	} else {
		fmt.Fprintf(&b, "- Invariant violations recorded: %d.\n", len(invariantBreakers))
		for _, v := range invariantBreakers {
			fmt.Fprintf(&b, "  - %s\n", v)
		}
	}
	if best, worst, ok := fastestAndSlowestRead(cells); ok {
		fmt.Fprintf(&b, "- Highest and lowest valid read result, across all cells and reads: %.1f ops/s (`%s`) and %.1f ops/s (`%s`).\n",
			best.ops, best.label, worst.ops, worst.label)
	}
	b.WriteString("\nThis list ranks by number only. Which design is *better* is not a measurement and is not stated here.\n\n")

	// ---------------------------------------------------------------- analysis index
	b.WriteString("## Analyses of this data\n\n")
	if len(analyses) == 0 {
		fmt.Fprintf(&b, "No signed analysis exists yet for inputs digest `%s`.\n\n", digest)
	} else {
		b.WriteString("| Analysis | Analyst | Digest | Status |\n|---|---|---|---|\n")
		for _, a := range analyses {
			stale := ""
			if a.Digest != "" && a.Digest != digest {
				stale = " **stale — different digest**"
			}
			fmt.Fprintf(&b, "| `%s` | %s | `%s` | %s%s |\n",
				orNone(a.AnalysisID), a.Analyst, a.Digest, orNone(a.Status), stale)
		}
		b.WriteString("\n")
	}

	// ---------------------------------------------------------------- cells
	b.WriteString("## Cells\n\n")
	t := markdown.NewTable("Topology", "Design", "Family", "Kind", "Gate", "Load s", "Entries", "Serialized B")
	for _, c := range cells {
		gate := "pass"
		if !c.Gate.Passed {
			gate = "FAIL"
		}
		if c.Error != "" {
			gate = "aborted"
		}
		fmt.Fprintf(&b, "")
		t.Row(c.Topology, "`"+c.Design+"`", c.Family, string(catKind(c)), gate,
			fmt.Sprintf("%.2f", (c.Load.SchemaMS+c.Load.StagingMS+c.Load.LoadStmtMS+c.Load.IndexMS)/1000),
			fmt.Sprintf("%d", c.Dataset.TotalEntries),
			markdown.Bytes(c.Dataset.SerializedBytes))
	}
	t.Write(&b)
	b.WriteString("\n")

	// ---------------------------------------------------------------- dataset
	b.WriteString("## Dataset (identical for every design and engine)\n\n")
	d := cells[0].Dataset
	b.WriteString("| Property | Value |\n|---|---|\n")
	fmt.Fprintf(&b, "| Seed | %d |\n", d.Seed)
	fmt.Fprintf(&b, "| Product definitions | %d |\n", d.ProductDefs)
	fmt.Fprintf(&b, "| Environments | %d |\n", d.Environments)
	fmt.Fprintf(&b, "| Business units | %d |\n", d.BusinessUnits)
	fmt.Fprintf(&b, "| Installed products | %d |\n", d.Installations)
	fmt.Fprintf(&b, "| Entries per installed product (tier) | %d |\n", d.EntriesPerIP)
	fmt.Fprintf(&b, "| Cardinality mode | %s |\n", d.CardinalityMode)
	fmt.Fprintf(&b, "| Value regime | %s |\n", d.ValueRegime)
	fmt.Fprintf(&b, "| Key bytes | %d |\n", d.KeyBytes)
	fmt.Fprintf(&b, "| Value bytes | %d |\n", d.ValueBytes)
	fmt.Fprintf(&b, "| Total entries | %d |\n", d.TotalEntries)
	fmt.Fprintf(&b, "| Total serialized bytes | %s |\n", markdown.Bytes(d.SerializedBytes))
	fmt.Fprintf(&b, "| Mean serialized bytes per installation | %.0f |\n\n", d.MeanSerializedB)

	// ---------------------------------------------------------------- reads
	b.WriteString("## Reads (isolated)\n\n")
	b.WriteString("`ops/s` is the median across trials. A run with one trial has no spread and no error bar.\n\n")
	rt := markdown.NewTable("Topology", "Design", "Read", "ops/s", "p50 ms", "p99 ms", "max ms", "ops", "errors", "spread %")
	for _, c := range cells {
		for _, r := range c.Reads {
			rt.Row(c.Topology, "`"+c.DesignShort+"`", r.Name, markdown.Ops(r.OpsPerSec),
				markdown.MS(r.Latency.P50MS), markdown.MS(r.Latency.P99MS), markdown.MS(r.Latency.MaxMS),
				fmt.Sprintf("%d", r.Ops), fmt.Sprintf("%d", r.Errors), fmt.Sprintf("%.1f", r.SpreadPct))
		}
	}
	rt.Write(&b)
	b.WriteString("\n")

	// ---------------------------------------------------------------- writes
	b.WriteString("## Writes (isolated, one worker)\n\n")
	b.WriteString("`ops/s` here is the cost of one update. Concurrency is a different question and is measured separately.\n\n")
	wt := markdown.NewTable("Topology", "Design", "Write", "ops/s", "p50 ms", "p99 ms", "ops", "errors", "spread %")
	for _, c := range cells {
		for _, r := range c.Writes {
			wt.Row(c.Topology, "`"+c.DesignShort+"`", r.Name, markdown.Ops(r.OpsPerSec),
				markdown.MS(r.Latency.P50MS), markdown.MS(r.Latency.P99MS),
				fmt.Sprintf("%d", r.Ops), fmt.Sprintf("%d", r.Errors), fmt.Sprintf("%.1f", r.SpreadPct))
		}
	}
	wt.Write(&b)
	b.WriteString("\n")

	// ---------------------------------------------------------------- contention
	b.WriteString("## Contention (hot key, one installed product)\n\n")
	ct := markdown.NewTable("Topology", "Design", "Strategy", "Writers", "Acked", "Retries", "Conflicts",
		"ops/s", "p99 ms", "counter", "expected", "lost updates", "control")
	for _, c := range cells {
		for _, k := range c.Concurrency {
			ctrl := ""
			if c.Negative {
				ctrl = "negative control"
			}
			ct.Row(c.Topology, "`"+c.DesignShort+"`", strings.TrimPrefix(k.Name, "contention:"),
				fmt.Sprintf("%d", k.Writers), fmt.Sprintf("%d", k.Acknowledged),
				fmt.Sprintf("%d", k.Retries), fmt.Sprintf("%d", k.Conflicts),
				markdown.Ops(k.OpsPerSec), markdown.MS(k.Latency.P99MS),
				fmt.Sprintf("%d", k.FinalValue), fmt.Sprintf("%d", k.ExpectedValue),
				fmt.Sprintf("%d", k.LostUpdates), ctrl)
		}
	}
	ct.Write(&b)
	b.WriteString("\n")

	// ---------------------------------------------------------------- cadence
	if hasCadence(cells) {
		b.WriteString("## Cadence\n\n")
		b.WriteString("The offered rate is *calculated*: installed products / period. A compressed rate is a\ncompressed-time validation, never the temporal result it stands in for.\n\n")
		kt := markdown.NewTable("Topology", "Design", "Regime", "Fleet", "Calculated rate /s",
			"Offered", "Started", "Completed", "Dropped", "Queue max", "Client saturated", "lag p99 ms", "Compressed")
		for _, c := range cells {
			for _, k := range c.Cadence {
				kt.Row(c.Topology, "`"+c.DesignShort+"`", k.Regime, fmt.Sprintf("%d", k.Fleet),
					fmt.Sprintf("%.4f", k.OfferedRate),
					fmt.Sprintf("%d", k.Schedule.Offered), fmt.Sprintf("%d", k.Schedule.Started),
					fmt.Sprintf("%d", k.Schedule.Completed), fmt.Sprintf("%d", k.Schedule.Dropped),
					fmt.Sprintf("%d", k.Schedule.QueueDepthMax),
					fmt.Sprintf("%t", k.Schedule.ClientSaturated),
					markdown.MS(k.Schedule.SchedulingLagMS.P99MS),
					fmt.Sprintf("%t", k.Compressed))
			}
		}
		kt.Write(&b)
		b.WriteString("\n")
	}

	// ---------------------------------------------------------------- audits
	b.WriteString("## Audits\n\n")
	at := markdown.NewTable("Topology", "Design", "Phase", "Passed", "Checks", "Key/value mismatches", "Revision mismatches", "Design checks", "First failure")
	for _, c := range cells {
		for i, a := range c.Audits {
			dc := []string{}
			names := make([]string, 0, len(a.DesignChecks))
			for n := range a.DesignChecks {
				names = append(names, n)
			}
			sort.Strings(names)
			for _, n := range names {
				dc = append(dc, fmt.Sprintf("%s=%d", n, a.DesignChecks[n]))
			}
			first := ""
			if len(a.Failures) > 0 {
				first = firstLine(a.Failures[0])
			}
			at.Row(c.Topology, "`"+c.DesignShort+"`", fmt.Sprintf("%d", i+1), fmt.Sprintf("%t", a.Passed),
				fmt.Sprintf("%d", a.Checks), fmt.Sprintf("%d", a.KeyValueMismatches),
				fmt.Sprintf("%d", a.RevisionMismatch), strings.Join(dc, " "), first)
		}
	}
	at.Write(&b)
	b.WriteString("\n")

	// ---------------------------------------------------------------- storage
	if hasStorage(cells) {
		b.WriteString("## Storage\n\n")
		b.WriteString("PostgreSQL reports `pg_total_relation_size` per relation. The route to a YugabyteDB\n" +
			"size figure is a different query, so none is reported here rather than a number that\n" +
			"would be read as bytes on disk.\n\n")
		st := markdown.NewTable("Topology", "Design", "Relation", "Bytes")
		for _, c := range cells {
			names := make([]string, 0, len(c.Storage))
			for n := range c.Storage {
				names = append(names, n)
			}
			sort.Strings(names)
			for _, n := range names {
				st.Row(c.Topology, "`"+c.DesignShort+"`", n, markdown.Bytes(c.Storage[n]))
			}
		}
		st.Write(&b)
		b.WriteString("\n")
	}

	// ---------------------------------------------------------------- limitations
	b.WriteString("## Limitations stated with the numbers\n\n")
	b.WriteString("- The client and the database share the same eight heterogeneous cores; absolute\n")
	b.WriteString("  throughput is lower than a two-machine setup would give, and relative comparisons are\n")
	b.WriteString("  the reason the numbers are comparable at all.\n")
	b.WriteString("- There is no real network. Distributed designs look better here than they would across\n")
	b.WriteString("  availability zones, so `EXPLAIN (ANALYZE, DIST)` RPC counts are the portable signal.\n")
	b.WriteString("- The read and write loops are closed-loop per operation; deep tails are optimistic\n")
	b.WriteString("  floors (coordinated omission) and are compared between designs, never quoted as SLOs.\n")
	b.WriteString("- A single trial has no error bar. Cells whose own trials disagree by more than about a\n")
	b.WriteString("  fifth are flagged in the tables above.\n")
	if !strings.Contains(strings.Join(phasesOf(cells), ","), "explain") {
		b.WriteString("- Plan capture was not part of this run's phases.\n")
	}
	b.WriteString("\n")

	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		return err
	}
	return os.WriteFile(out, []byte(b.String()), 0o644)
}

// ---------------------------------------------------------------- helpers

type readPoint struct {
	ops   float64
	label string
}

func fastestAndSlowestRead(cells []Result) (readPoint, readPoint, bool) {
	var all []readPoint
	for _, c := range cells {
		if !c.Gate.Passed {
			continue
		}
		for _, r := range c.Reads {
			if r.OpsPerSec > 0 {
				all = append(all, readPoint{r.OpsPerSec, c.Topology + "/" + c.DesignShort + "/" + r.Name})
			}
		}
	}
	if len(all) == 0 {
		return readPoint{}, readPoint{}, false
	}
	sort.Slice(all, func(i, j int) bool { return all[i].ops < all[j].ops })
	return all[len(all)-1], all[0], true
}

func catKind(c Result) Kind {
	if c.Family == "document" {
		return Doc
	}
	return Rows
}

func hasCadence(cells []Result) bool {
	for _, c := range cells {
		if len(c.Cadence) > 0 {
			return true
		}
	}
	return false
}

func hasStorage(cells []Result) bool {
	for _, c := range cells {
		if len(c.Storage) > 0 {
			return true
		}
	}
	return false
}

func phasesOf(cells []Result) []string {
	var out []string
	for _, c := range cells {
		out = append(out, c.Options.Phases)
	}
	return out
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		s = s[:i]
	}
	if len(s) > 160 {
		s = s[:160] + "…"
	}
	return strings.ReplaceAll(s, "|", "\\|")
}
