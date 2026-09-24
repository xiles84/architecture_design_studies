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

// A generated report holds measurements and no interpretation (methodology 11). Its
// TL;DR lists measured facts selected by fixed rules -- which cells failed, whether
// the controls fired, the wrong-read rate beside every throughput. It ranks by the
// number and says nothing about which scenario is better.
//
// One rule is absolute here: a relaxed cell's throughput NEVER appears without its
// wrong-read count and rate beside it. Never call a stale result fast.

func writeReport(resultsDir, out string) error {
	if resultsDir == "" {
		return fmt.Errorf("report: -results is required")
	}
	if out == "" {
		return fmt.Errorf("report: -report-out is required")
	}
	cells, err := filestore.ReadJSONTree[CellResult](resultsDir)
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
		if cells[i].Group != cells[j].Group {
			return cells[i].Group < cells[j].Group
		}
		return cells[i].Scenario < cells[j].Scenario
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
	fmt.Fprintf(&b, "# Study 05 — external cache throughput and consistency — run `%s`\n\n", runID)
	fmt.Fprintf(&b, "Generated %s from %d result file(s). This report contains measurements only; conclusions live in\nsigned analyses under `reports/analyses/`.\n\n", time.Now().UTC().Format(time.RFC3339), fileCount)

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
	if cells[0].BenchImage != "" {
		fmt.Fprintf(&b, "| Benchmark image | `%s` (`%s`) |\n", cells[0].BenchImage, cells[0].BenchImageID)
	}
	fmt.Fprintf(&b, "| Resource framing | `%s` |\n", cells[0].Options.ResourceFrame)
	fmt.Fprintf(&b, "| **Inputs digest** | `%s` |\n", digest)
	fmt.Fprintf(&b, "| Cells | %d |\n", len(cells))
	cells2 := uniqueScenarios(cells)
	fmt.Fprintf(&b, "| Distinct scenarios | %d |\n", len(cells2))
	fmt.Fprintf(&b, "| Hard TTL | %d s (probabilistic early expiry p = clamp(1 − remaining/%d s, 0, 1)) |\n", int(hardTTLSeconds), int(hardTTLSeconds))
	fmt.Fprintf(&b, "| Trialling | %d trial(s) per measurement; throughput is the median |\n\n", cells[0].Options.Trials)

	// ---------------------------------------------------------------- TL;DR
	b.WriteString("## TL;DR — measured facts, selected by fixed rules\n\n")
	var failed, controlsFired, controlsMissed []string
	var strictViolators, impossibleCells []string
	for _, c := range cells {
		tag := fmt.Sprintf("`%s/%s`", c.Topology, c.Scenario)
		if c.Error != "" {
			failed = append(failed, fmt.Sprintf("%s (%s)", tag, firstLine(c.Error)))
		}
		if c.Negative {
			fired := false
			for _, f := range c.Faults {
				if f.Reproduced {
					fired = true
				}
			}
			if fired {
				controlsFired = append(controlsFired, tag)
			} else {
				controlsMissed = append(controlsMissed, tag)
			}
		}
		if !c.Negative && c.StrictViolations > 0 {
			strictViolators = append(strictViolators, fmt.Sprintf("%s: %d stale-after-ack read(s)", tag, c.StrictViolations))
		}
		if c.ImpossibleValues > 0 {
			impossibleCells = append(impossibleCells, fmt.Sprintf("%s: %d impossible cache value(s)", tag, c.ImpossibleValues))
		}
	}
	fmt.Fprintf(&b, "- Cells reported: %d; cells that failed: %d.\n", len(cells), len(failed))
	for _, f := range failed {
		fmt.Fprintf(&b, "  - failed: %s\n", f)
	}
	fmt.Fprintf(&b, "- Negative controls that fired: %d. Controls that did NOT fire: %d.\n", len(controlsFired), len(controlsMissed))
	for _, c := range controlsMissed {
		fmt.Fprintf(&b, "  - did not fire: %s — the associated correctness claim is **not demonstrated**\n", c)
	}
	if len(strictViolators) == 0 {
		b.WriteString("- No strict cell recorded a stale-after-ack read.\n")
	} else {
		fmt.Fprintf(&b, "- Strict cells that violated their contract: %d.\n", len(strictViolators))
		for _, v := range strictViolators {
			fmt.Fprintf(&b, "  - %s\n", v)
		}
	}
	if len(impossibleCells) == 0 {
		b.WriteString("- No cell reported an impossible / uncommitted cache value.\n")
	} else {
		fmt.Fprintf(&b, "- Cells with impossible cache values: %d (**these cells are invalid under every freshness policy**).\n", len(impossibleCells))
		for _, v := range impossibleCells {
			fmt.Fprintf(&b, "  - %s\n", v)
		}
	}
	if hi, lo, ok := extremeThroughput(cells); ok {
		fmt.Fprintf(&b, "- Highest and lowest measured throughput across all cells: %s (`%s`) and %s (`%s`).\n",
			markdown.Ops(hi.v), hi.label, markdown.Ops(lo.v), lo.label)
	}
	if w, ok := highestWrongRate(cells); ok {
		fmt.Fprintf(&b, "- Highest wrong-read rate recorded: %.2f%% of all reads (`%s`), with %d wrong read(s) in %d.\n",
			w.rate, w.label, w.wrong, w.total)
	}
	b.WriteString("\nThis list ranks by number only. Which scenario is *better* is not a measurement and is not stated here.\n\n")

	// ---------------------------------------------------------------- analyses
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
			fmt.Fprintf(&b, "| `%s` | %s | `%s` | %s%s |\n", orNone(a.AnalysisID), a.Analyst, orNone(a.Digest), orNone(a.Status), stale)
		}
		b.WriteString("\n")
	}

	// ---------------------------------------------------------------- scenarios
	b.WriteString("## Scenarios\n\n")
	b.WriteString("Every scenario id retains all its dimensions: `<model>-<version>-<backend>-<strategy>-<freshness>-<writers>`.\n")
	b.WriteString("The strict-freshness policy column is the study's own account of HOW a strict read proves freshness.\n\n")
	t := markdown.NewTable("Topology", "Scenario", "Model", "Ver", "Backend", "Strategy", "Freshness", "Writers",
		"Strict policy", "Gate", "Strict viol.", "Impossible", "Control")
	for _, c := range cells {
		gate := "pass"
		if !c.Gate.Passed {
			gate = "FAIL"
		}
		if c.Error != "" {
			gate = "aborted"
		}
		ctrl := ""
		if c.Negative {
			fired := false
			for _, f := range c.Faults {
				if f.Reproduced {
					fired = true
				}
			}
			if fired {
				ctrl = "fired"
			} else {
				ctrl = "NOT fired"
			}
		}
		t.Row(c.Topology, "`"+c.Scenario+"`", c.ModelDim, c.VersionDim, c.BackendDim, c.StrategyDim,
			c.FreshnessDim, c.WritersDim, dash(c.StrictPolicy), gate,
			fmt.Sprintf("%d", c.StrictViolations), fmt.Sprintf("%d", c.ImpossibleValues), ctrl)
	}
	t.Write(&b)

	// ---------------------------------------------------------------- throughput
	b.WriteString("## Throughput\n\n")
	b.WriteString("A relaxed scenario's throughput is shown ONLY beside its wrong-read count and rate. A stale result is\nnot a fast result, and this table is built so it cannot be read as one.\n\n")
	tt := markdown.NewTable("Topology", "Scenario", "Phase", "Endpoint", "ops/s", "p50 ms", "p99 ms", "max ms", "ops", "errors", "spread %", "wrong reads", "wrong % of reads")
	for _, c := range cells {
		wr := totalWrong(c)
		for _, r := range c.Warm {
			tt.Row(c.Topology, "`"+c.ScenarioShort+"`", r.Name, "cacheable", markdown.Ops(r.OpsPerSec),
				markdown.MS(r.Latency.P50MS), markdown.MS(r.Latency.P99MS), markdown.MS(r.Latency.MaxMS),
				fmt.Sprintf("%d", r.Ops), fmt.Sprintf("%d", r.Errors), fmt.Sprintf("%.1f", r.SpreadPct),
				fmt.Sprintf("%d", wr.wrong), fmt.Sprintf("%.2f%%", wr.rate))
		}
		for _, r := range c.Mixed {
			tt.Row(c.Topology, "`"+c.ScenarioShort+"`", r.Name, "cacheable", markdown.Ops(r.OpsPerSec),
				markdown.MS(r.Latency.P50MS), markdown.MS(r.Latency.P99MS), markdown.MS(r.Latency.MaxMS),
				fmt.Sprintf("%d", r.Ops), fmt.Sprintf("%d", r.Errors), fmt.Sprintf("%.1f", r.SpreadPct),
				fmt.Sprintf("%d", wr.wrong), fmt.Sprintf("%.2f%%", wr.rate))
		}
		for _, r := range c.AppMix {
			tt.Row(c.Topology, "`"+c.ScenarioShort+"`", r.Name, "total app", markdown.Ops(r.OpsPerSec),
				markdown.MS(r.Latency.P50MS), markdown.MS(r.Latency.P99MS), markdown.MS(r.Latency.MaxMS),
				fmt.Sprintf("%d", r.Ops), fmt.Sprintf("%d", r.Errors), fmt.Sprintf("%.1f", r.SpreadPct),
				fmt.Sprintf("%d", wr.wrong), fmt.Sprintf("%.2f%%", wr.rate))
		}
	}
	tt.Write(&b)

	// ---------------------------------------------------------------- open loop
	if hasOpenLoop(cells) {
		b.WriteString("## Open-loop demand (fixed offered rate, no closed-loop backpressure)\n\n")
		b.WriteString("The offered rate is what the client scheduled; the delivered rate is what the database completed.\n")
		b.WriteString("`dropped` counts offered operations the generator could not attempt, and `client saturated` means the\n")
		b.WriteString("client — not the server — was the limit, which invalidates a claim about the server. The wrong-read\n")
		b.WriteString("count is the same phase's, so no delivered rate is shown without its correctness contract. Closed-loop\n")
		b.WriteString("results elsewhere in this report are floors and are never reused as open-loop.\n\n")
		ot := markdown.NewTable("Topology", "Scenario", "Phase", "offered ops/s", "delivered ops/s", "offered", "started", "completed", "rejected", "errors", "dropped", "client saturated", "p50 ms", "p99 ms", "max ms", "sched lag p99 ms", "wrong reads", "wrong % of reads")
		for _, c := range cells {
			for _, r := range c.OpenLoop {
				wrong, _, rate := wrongForPhase(c, r.Name)
				sat := "no"
				if r.ClientSaturated {
					sat = "YES — " + r.SaturationReason
				}
				ot.Row(c.Topology, "`"+c.ScenarioShort+"`", r.Name,
					fmt.Sprintf("%.0f", r.OfferedRate), fmt.Sprintf("%.0f", r.DeliveredRate),
					fmt.Sprintf("%d", r.Offered), fmt.Sprintf("%d", r.Started), fmt.Sprintf("%d", r.Completed),
					fmt.Sprintf("%d", r.Rejected), fmt.Sprintf("%d", r.Errors), fmt.Sprintf("%d", r.Dropped),
					sat, markdown.MS(r.Latency.P50MS), markdown.MS(r.Latency.P99MS), markdown.MS(r.Latency.MaxMS),
					markdown.MS(r.SchedulingLagMS.P99MS),
					fmt.Sprintf("%d", wrong), fmt.Sprintf("%.2f%%", rate))
			}
		}
		ot.Write(&b)
	}

	// ---------------------------------------------------------------- writes
	if hasWrites(cells) {
		b.WriteString("## Writes (isolated per operation, then the hotspot race)\n\n")
		wt := markdown.NewTable("Topology", "Scenario", "Write", "ops/s", "p50 ms", "p99 ms", "ops", "errors", "spread %")
		for _, c := range cells {
			for _, r := range c.Writes {
				wt.Row(c.Topology, "`"+c.ScenarioShort+"`", r.Name, markdown.Ops(r.OpsPerSec),
					markdown.MS(r.Latency.P50MS), markdown.MS(r.Latency.P99MS),
					fmt.Sprintf("%d", r.Ops), fmt.Sprintf("%d", r.Errors), fmt.Sprintf("%.1f", r.SpreadPct))
			}
		}
		wt.Write(&b)
	}

	// ---------------------------------------------------------------- wrong reads
	b.WriteString("## Wrong-read accounting (per phase)\n\n")
	b.WriteString("Freshness is judged by CONTENT HASH against an independent Go oracle, not by re-reading the database\non every hit. `ahead` counts reads that returned a later committed state than they required, which is\nlegitimate when a write commits mid-read. `concurrent/ambiguous` counts reads that overlapped a write and\nis kept separate from both fresh and wrong.\n\n")
	rt := markdown.NewTable("Topology", "Scenario", "Phase", "reads", "fresh", "wrong", "wrong % all", "hits",
		"wrong % hits", "keys", "max behind", "stale p50 ms", "stale p99 ms", "streak", "TTF p99 ms",
		"concurrent", "impossible", "causes")
	for _, c := range cells {
		for _, pw := range c.Wrong {
			s := pw.Summary
			rt.Row(c.Topology, "`"+c.ScenarioShort+"`", pw.Phase,
				fmt.Sprintf("%d", s.TotalReads), fmt.Sprintf("%d", s.FreshReads), fmt.Sprintf("%d", s.WrongReads),
				fmt.Sprintf("%.2f%%", s.WrongPctOfAllReads), fmt.Sprintf("%d", s.CacheHits),
				fmt.Sprintf("%.2f%%", s.WrongPctOfCacheHits), fmt.Sprintf("%d", s.UniqueKeysAffected),
				fmt.Sprintf("%d", s.MaxVersionsBehind), markdown.MS(s.StaleDuration.P50), markdown.MS(s.StaleDuration.P99),
				fmt.Sprintf("%d", s.ConsecutiveWrongMax), markdown.MS(s.TimeToFreshness.P99),
				fmt.Sprintf("%d", s.ConcurrentAmbiguous), fmt.Sprintf("%d", s.ImpossibleValues), causeString(s.ByCause))
		}
	}
	rt.Write(&b)

	// ---------------------------------------------------------------- expiry
	b.WriteString("## Probabilistic early expiration: observed vs specified\n\n")
	b.WriteString("The draw is `p = clamp(1 − remaining/300 s, 0, 1)`. The observed rate by age bucket is the empirical\ncheck on the implementation, alongside the unit tests that pin the same ages with a fake clock.\n\n")
	et := markdown.NewTable("Topology", "Scenario", "Age bucket", "hits", "early expiries", "observed %", "specified p (band)")
	for _, c := range cells {
		names := make([]string, 0, len(c.ExpiryByAge))
		for n := range c.ExpiryByAge {
			names = append(names, n)
		}
		sort.Strings(names)
		for _, n := range names {
			v := c.ExpiryByAge[n]
			rate := 0.0
			if v["hits"] > 0 {
				rate = 100 * float64(v["probabilistic_expiries"]) / float64(v["hits"])
			}
			et.Row(c.Topology, "`"+c.ScenarioShort+"`", n, fmt.Sprintf("%d", v["hits"]),
				fmt.Sprintf("%d", v["probabilistic_expiries"]), fmt.Sprintf("%.1f%%", rate), specifiedBand(n))
		}
	}
	if et.Len() == 0 {
		b.WriteString("*No cache hits were observed in this run, so there is no expiry evidence to report.*\n\n")
	} else {
		et.Write(&b)
	}

	// ---------------------------------------------------------------- lease
	b.WriteString("## Leases and stampede prevention\n\n")
	lt := markdown.NewTable("Topology", "Scenario", "acquired", "contended", "waits", "wait p99 ms", "timeouts",
		"fallbacks", "duplicate fills", "lease duration", "calibration")
	for _, c := range cells {
		if !strings.Contains(c.Scenario, "-none-") {
			l := c.Lease
			lt.Row(c.Topology, "`"+c.ScenarioShort+"`", fmt.Sprintf("%d", l.Acquisitions),
				fmt.Sprintf("%d", l.Contended), fmt.Sprintf("%d", l.Waits), markdown.MS(l.WaitP99MS),
				fmt.Sprintf("%d", l.Timeouts), fmt.Sprintf("%d", l.FallbackReads),
				fmt.Sprintf("%d", l.DuplicateFill), dash(l.LeaseDuration), dash(l.Calibration))
		}
	}
	lt.Write(&b)

	st := markdown.NewTable("Topology", "Scenario", "readers", "keys", "elapsed ms", "db loads", "loads/key", "fallbacks", "duplicate fills", "wrong", "impossible")
	nStampede := 0
	for _, c := range cells {
		if c.Stampede == nil {
			continue
		}
		nStampede++
		s := c.Stampede
		st.Row(c.Topology, "`"+c.ScenarioShort+"`", fmt.Sprintf("%d", s.Readers), fmt.Sprintf("%d", s.Keys),
			fmt.Sprintf("%.1f", s.ElapsedMS), fmt.Sprintf("%d", s.DatabaseLoads),
			fmt.Sprintf("%.2f", s.LoadsPerKey), fmt.Sprintf("%d", s.Fallbacks),
			fmt.Sprintf("%d", s.DuplicateFills), fmt.Sprintf("%d", s.WrongReads), fmt.Sprintf("%d", s.Impossible))
	}
	if nStampede > 0 {
		st.Write(&b)
	} else {
		b.WriteString("*The stampede phase did not run for any cell in this run.*\n\n")
	}

	// ---------------------------------------------------------------- churn
	ct := markdown.NewTable("Topology", "Scenario", "duration s", "TTL boundaries", "reads", "writes", "hard expiries", "early expiries", "wrong", "impossible")
	nChurn := 0
	for _, c := range cells {
		if c.Churn == nil {
			continue
		}
		nChurn++
		k := c.Churn
		ct.Row(c.Topology, "`"+c.ScenarioShort+"`", fmt.Sprintf("%.1f", k.DurationS),
			fmt.Sprintf("%.1f", k.TTLBoundaries), fmt.Sprintf("%d", k.Reads), fmt.Sprintf("%d", k.Writes),
			fmt.Sprintf("%d", k.HardExpired), fmt.Sprintf("%d", k.ProbExpired),
			fmt.Sprintf("%d", k.WrongReads), fmt.Sprintf("%d", k.Impossible))
	}
	if nChurn > 0 {
		b.WriteString("## Sustained churn across real TTL boundaries\n\n")
		b.WriteString("The primary hard TTL is 300 s and is never shortened for a run: these phases last long enough to\ncross real boundaries.\n\n")
		ct.Write(&b)
	}

	// ---------------------------------------------------------------- instances
	it := markdown.NewTable("Topology", "Scenario", "instances", "workers", "backend", "app ops/s", "spread %", "cacheable ops/s", "bypass reads", "wrong", "impossible", "note")
	nInst := 0
	for _, c := range cells {
		for _, ir := range c.Instances {
			nInst++
			it.Row(c.Topology, "`"+c.ScenarioShort+"`", fmt.Sprintf("%d", ir.Instances), fmt.Sprintf("%d", ir.Workers),
				ir.Backend, markdown.Ops(ir.AppOpsPerSec), fmt.Sprintf("%.1f", ir.AppSpreadPct),
				markdown.Ops(ir.CacheOpsPerSec), fmt.Sprintf("%d", ir.BypassReads),
				fmt.Sprintf("%d", ir.WrongReads), fmt.Sprintf("%d", ir.Impossible), ir.Note)
		}
	}
	if nInst > 0 {
		b.WriteString("## One versus three logical application instances (equal total workers)\n\n")
		it.Write(&b)
	}

	// ---------------------------------------------------------------- faults
	b.WriteString("## Fault injection and controls\n\n")
	ft := markdown.NewTable("Topology", "Scenario", "Fault", "reproduced", "wrong observed", "impossible observed", "correctness")
	nFaults := 0
	for _, c := range cells {
		for _, f := range c.Faults {
			nFaults++
			ft.Row(c.Topology, "`"+c.ScenarioShort+"`", f.Name, fmt.Sprintf("%t", f.Reproduced),
				fmt.Sprintf("%d", f.WrongReads), fmt.Sprintf("%d", f.Impossible), dash(f.Correctness))
		}
	}
	if nFaults > 0 {
		b.WriteString("A control's or a fault's speed is never shown as a result: a design that breaks an invariant quickly\nhas not been fast.\n\n")
		ft.Write(&b)
	} else {
		b.WriteString("*No fault phase ran in this run.*\n\n")
	}

	// ---------------------------------------------------------------- audits
	b.WriteString("## Ledger replay (audit)\n\n")
	at := markdown.NewTable("Topology", "Scenario", "Phase #", "passed", "checks", "content mismatches", "donation-set mismatches", "aggregate mismatches", "recent-slice mismatches", "design checks", "first failure")
	nAudits := 0
	for _, c := range cells {
		for i, a := range c.Audits {
			nAudits++
			names := make([]string, 0, len(a.DesignChecks))
			for n := range a.DesignChecks {
				names = append(names, n)
			}
			sort.Strings(names)
			var dc []string
			for _, n := range names {
				dc = append(dc, fmt.Sprintf("%s=%d", n, a.DesignChecks[n]))
			}
			first := ""
			if len(a.Failures) > 0 {
				first = firstLine(a.Failures[0])
			}
			at.Row(c.Topology, "`"+c.ScenarioShort+"`", fmt.Sprintf("%d", i+1), fmt.Sprintf("%t", a.Passed),
				fmt.Sprintf("%d", a.Checks), fmt.Sprintf("%d", a.KeyValueMismatches),
				fmt.Sprintf("%d", a.DonationSetMismatch), fmt.Sprintf("%d", a.AggregateMismatch),
				fmt.Sprintf("%d", a.RecentSliceMismatch), strings.Join(dc, " "), first)
		}
	}
	if nAudits > 0 {
		at.Write(&b)
	} else {
		b.WriteString("*No audit phase ran in this run.*\n\n")
	}

	// ---------------------------------------------------------------- cache stats
	b.WriteString("## Cache accounting\n\n")
	kt := markdown.NewTable("Topology", "Scenario", "backend", "capacity B", "resident B", "items", "evictions", "fills", "publishes", "fenced", "publish failures", "invalidations", "tombstone fences", "version validations", "bypass reads", "external writes", "unrecorded states confirmed", "write no-ops")
	for _, c := range cells {
		k := c.Cache
		kt.Row(c.Topology, "`"+c.ScenarioShort+"`", dash(k.Backend), markdown.Bytes(k.CapacityBytes),
			markdown.Bytes(k.ResidentBytes), fmt.Sprintf("%d", k.Items), fmt.Sprintf("%d", k.Evictions),
			fmt.Sprintf("%d", k.Fills), fmt.Sprintf("%d", k.Publishes), fmt.Sprintf("%d", k.PublishFenced),
			fmt.Sprintf("%d", k.PublishFailed), fmt.Sprintf("%d", k.Invalidations),
			fmt.Sprintf("%d", k.Tombstones), fmt.Sprintf("%d", k.ValidationQueries),
			fmt.Sprintf("%d", k.BypassReads), fmt.Sprintf("%d", k.ExternalWrites),
			fmt.Sprintf("%d", k.UnrecordedConfirmed), fmt.Sprintf("%d", k.WriteNoops))
	}
	kt.Write(&b)

	// ---------------------------------------------------------------- redis info
	if hasRedis(cells) {
		b.WriteString("## Redis accounting (from `INFO`)\n\n")
		b.WriteString("Redis's `allkeys-lru` is an APPROXIMATE LRU, and the memory backend's is exact; a difference in\nbehaviour between the two is partly a policy difference and is stated as such.\n\n")
		keys := []string{"redis_version", "maxmemory", "maxmemory_policy", "maxmemory_samples", "used_memory",
			"used_memory_peak", "keyspace_hits", "keyspace_misses", "expired_keys", "evicted_keys",
			"total_commands_processed", "total_net_input_bytes", "total_net_output_bytes", "used_cpu_sys", "used_cpu_user"}
		ri := markdown.NewTable(append([]string{"Topology", "Scenario"}, keys...)...)
		for _, c := range cells {
			if len(c.Redis) == 0 {
				continue
			}
			row := []string{c.Topology, "`" + c.ScenarioShort + "`"}
			for _, k := range keys {
				row = append(row, dash(c.Redis[k]))
			}
			ri.Row(row...)
		}
		ri.Write(&b)
	}

	// ---------------------------------------------------------------- client cpu
	b.WriteString("## Client and resource conditions\n\n")
	b.WriteString("CFS throttling is recorded, not assumed. A throttled client puts milliseconds into tails that belong\nto no scenario.\n\n")
	cpu := markdown.NewTable("Topology", "Scenario", "client periods", "client throttled", "client throttled ms", "framing", "workers")
	for _, c := range cells {
		periods, throttled, tms := int64(0), int64(0), int64(0)
		if c.ClientCPU != nil {
			periods = c.ClientCPU.Delta.NrPeriods
			throttled = c.ClientCPU.Delta.NrThrottled
			tms = c.ClientCPU.Delta.ThrottledUsec / 1000
		}
		cpu.Row(c.Topology, "`"+c.ScenarioShort+"`", fmt.Sprintf("%d", periods), fmt.Sprintf("%d", throttled),
			fmt.Sprintf("%d", tms), c.Options.ResourceFrame, fmt.Sprintf("%d", c.Options.Workers))
	}
	cpu.Write(&b)

	// ---------------------------------------------------------------- storage
	if hasStorage(cells) {
		b.WriteString("## Storage\n\n")
		b.WriteString("PostgreSQL reports `pg_total_relation_size` per relation. YugabyteDB's route to a size figure is a\ndifferent query, so none is reported there rather than a number read as bytes on disk.\n\n")
		sz := markdown.NewTable("Topology", "Scenario", "Relation", "Bytes")
		for _, c := range cells {
			names := make([]string, 0, len(c.Storage))
			for n := range c.Storage {
				names = append(names, n)
			}
			sort.Strings(names)
			for _, n := range names {
				sz.Row(c.Topology, "`"+c.ScenarioShort+"`", n, markdown.Bytes(c.Storage[n]))
			}
		}
		sz.Write(&b)
	}

	// ---------------------------------------------------------------- gates
	b.WriteString("## Correctness gates\n\n")
	gt := markdown.NewTable("Topology", "Scenario", "gate passed", "checks", "statements", "failures")
	for _, c := range cells {
		gt.Row(c.Topology, "`"+c.ScenarioShort+"`", fmt.Sprintf("%t", c.Gate.Passed),
			fmt.Sprintf("%d", c.Gate.Checks), strings.Join(c.Gate.Statements, " "),
			strings.Join(firstLines(c.Gate.Failures), " | "))
	}
	gt.Write(&b)

	// ---------------------------------------------------------------- ledger
	b.WriteString("## Ledger assertion: the requirement must never be ahead of the database\n\n")
	b.WriteString("After every writing phase the harness compares, for every donor, the content the database\n")
	b.WriteString("actually holds against the ledger's freshness requirement. A requirement BEHIND a committed\n")
	b.WriteString("but unacknowledged state is legitimate; a requirement AHEAD of the database is not, and fails\n")
	b.WriteString("the cell, because a bypass read that returned such a state would be an accounting defect\n")
	b.WriteString("rather than a design finding.\n\n")
	lgt := markdown.NewTable("Topology", "Scenario", "Phase", "donors", "mismatches", "passed", "first example")
	nLedger := 0
	for _, c := range cells {
		for _, lc := range c.Ledger {
			nLedger++
			ex := ""
			if len(lc.Examples) > 0 {
				ex = firstLine(lc.Examples[0])
			}
			lgt.Row(c.Topology, "`"+c.ScenarioShort+"`", lc.Phase, fmt.Sprintf("%d", lc.People),
				fmt.Sprintf("%d", lc.Mismatches), fmt.Sprintf("%t", lc.Passed), ex)
		}
	}
	if nLedger > 0 {
		lgt.Write(&b)
	} else {
		b.WriteString("*No ledger assertion ran in this run.*\n\n")
	}

	// ---------------------------------------------------------------- limitations
	b.WriteString("## Limitations stated with the numbers\n\n")
	b.WriteString("- The client, the cache and the database share the same eight heterogeneous cores of a laptop. Absolute\n")
	b.WriteString("  throughput is lower than a two-machine setup would give, and only relative comparisons within this\n")
	b.WriteString("  environment are meaningful.\n")
	b.WriteString("- There is no real network. A shared cache and a distributed database both look better here than they\n")
	b.WriteString("  would across availability zones; EXPLAIN (ANALYZE, DIST) RPC counts are the portable signal.\n")
	b.WriteString("- All containers run on one WSL2 machine. This is NOT evidence of data colocation or of network\n")
	b.WriteString("  behaviour; the placement pair reports engine evidence and physical colocation remains untested.\n")
	b.WriteString("- The workloads are closed-loop per operation within a phase (measure.Run), except the open-loop demand\n")
	b.WriteString("  phase, which schedules fixed arrival rates with measure.RunOpenLoop and reports the offered rate,\n")
	b.WriteString("  dropped arrivals, scheduling lag and whether the client saturated. Closed-loop results are therefore\n")
	b.WriteString("  optimistic floors and are compared between scenarios, never quoted as SLO figures.\n")
	b.WriteString("- A single trial has no error bar. Where trials were repeated the spread is shown per measurement, and\n")
	b.WriteString("  no conclusion may rest on a difference below the measured noise floor.\n")
	b.WriteString("- One DeepSeek HIGH agent planned, built, measured and analysed this study. There is no independent\n")
	b.WriteString("  model review of the implementation, the oracle or the conclusions.\n")
	b.WriteString("- A cache hit executes no SQL and therefore has no plan. The saved database operation is what a hit\n")
	b.WriteString("  avoids; the plans file records the statement it avoids.\n\n")

	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		return err
	}
	return os.WriteFile(out, []byte(b.String()), 0o644)
}

// ---------------------------------------------------------------- helpers

type tp struct {
	v     float64
	label string
}

func extremeThroughput(cells []CellResult) (tp, tp, bool) {
	var all []tp
	for _, c := range cells {
		if !c.Gate.Passed || c.Error != "" || c.Negative {
			continue
		}
		for _, r := range c.Warm {
			all = append(all, tp{r.OpsPerSec, c.Topology + "/" + c.ScenarioShort + "/" + r.Name})
		}
		for _, r := range c.Mixed {
			all = append(all, tp{r.OpsPerSec, c.Topology + "/" + c.ScenarioShort + "/" + r.Name})
		}
	}
	if len(all) == 0 {
		return tp{}, tp{}, false
	}
	sort.Slice(all, func(i, j int) bool { return all[i].v < all[j].v })
	return all[len(all)-1], all[0], true
}

type wrongRate struct {
	rate         float64
	label        string
	wrong, total int64
}

func highestWrongRate(cells []CellResult) (wrongRate, bool) {
	var best wrongRate
	found := false
	for _, c := range cells {
		for _, pw := range c.Wrong {
			s := pw.Summary
			if s.TotalReads == 0 || s.WrongReads == 0 {
				continue
			}
			if !found || s.WrongPctOfAllReads > best.rate {
				best = wrongRate{s.WrongPctOfAllReads, c.Topology + "/" + c.ScenarioShort + "/" + pw.Phase, s.WrongReads, s.TotalReads}
				found = true
			}
		}
	}
	return best, found
}

// hasOpenLoop reports whether any cell ran the open-loop demand phase.
func hasOpenLoop(cells []CellResult) bool {
	for _, c := range cells {
		if len(c.OpenLoop) > 0 {
			return true
		}
	}
	return false
}

// wrongForPhase returns the wrong-read account recorded for exactly one phase, so
// an open-loop rate is never shown without the correctness of its own phase.
func wrongForPhase(c CellResult, phase string) (wrong, total int64, rate float64) {
	for _, pw := range c.Wrong {
		if pw.Phase != phase {
			continue
		}
		wrong += pw.Summary.WrongReads
		total += pw.Summary.TotalReads
	}
	if total > 0 {
		rate = 100 * float64(wrong) / float64(total)
	}
	return wrong, total, rate
}

func totalWrong(c CellResult) struct {
	wrong, total int64
	rate         float64
} {
	var wrong, total int64
	for _, pw := range c.Wrong {
		wrong += pw.Summary.WrongReads
		total += pw.Summary.TotalReads
	}
	rate := 0.0
	if total > 0 {
		rate = 100 * float64(wrong) / float64(total)
	}
	return struct {
		wrong, total int64
		rate         float64
	}{wrong, total, rate}
}

func uniqueScenarios(cells []CellResult) []string {
	seen := map[string]bool{}
	var out []string
	for _, c := range cells {
		k := c.Topology + "/" + c.Scenario
		if !seen[k] {
			seen[k] = true
			out = append(out, k)
		}
	}
	return out
}

func causeString(m map[string]int64) string {
	if len(m) == 0 {
		return "—"
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%d", k, m[k]))
	}
	return strings.Join(parts, "; ")
}

// specifiedBand is the p the protocol's formula gives over each age bucket. The
// observed rate is compared with it; the unit tests pin the exact values at 0, 75,
// 150, 225 and 300 seconds.
func specifiedBand(bucket string) string {
	switch bucket {
	case "<0s":
		return "0"
	case "0-75s":
		return "0 to 0.25"
	case "75-150s":
		return "0.25 to 0.50"
	case "150-225s":
		return "0.50 to 0.75"
	case "225-300s":
		return "0.75 to 1.00"
	default:
		return "1"
	}
}

func hasWrites(cells []CellResult) bool {
	for _, c := range cells {
		if len(c.Writes) > 0 {
			return true
		}
	}
	return false
}

func hasRedis(cells []CellResult) bool {
	for _, c := range cells {
		if len(c.Redis) > 0 {
			return true
		}
	}
	return false
}

func hasStorage(cells []CellResult) bool {
	for _, c := range cells {
		if len(c.Storage) > 0 {
			return true
		}
	}
	return false
}

func firstLines(xs []string) []string {
	if len(xs) > 3 {
		xs = xs[:3]
	}
	out := make([]string, 0, len(xs))
	for _, x := range xs {
		out = append(out, firstLine(x))
	}
	return out
}

func dash(s string) string {
	if strings.TrimSpace(s) == "" {
		return "—"
	}
	return strings.ReplaceAll(s, "|", "\\|")
}
