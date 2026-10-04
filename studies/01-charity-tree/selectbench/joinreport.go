package selectbench

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// joinClass hashes every setting that makes two runs comparable: environment,
// dataset fingerprint, benchmark-client image and database-engine image. The
// replicated report keys on it in addition to engine/topology/axis, so runs with
// different fixtures, budgets or images are partitioned instead of pooled. A
// single class means the inputs were compatible; more than one is printed
// explicitly rather than averaged.
func joinClass(r JoinResult) string {
	p := r.Provenance
	var sb strings.Builder
	for _, k := range sortedKeys(p.Settings) {
		sb.WriteString(k)
		sb.WriteByte('=')
		sb.WriteString(p.Settings[k])
		sb.WriteByte(';')
	}
	// Key on the image DIGESTS, not their names: a client image is tagged with
	// its run id, so including the name would make every run its own class and
	// silently disable cross-run replication. Two runs are compatible when the
	// content they ran (client build, engine build, dataset, settings) matches.
	raw := strings.Join([]string{
		p.Environment, p.ClientImageDigest,
		p.EngineImageDigest, p.DatasetFingerprint, sb.String(),
	}, "|")
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])[:12]
}

func joinClassDescription(r JoinResult) string {
	p := r.Provenance
	return fmt.Sprintf("env=`%s`; dataset=`%s`; client=`%s`@`%s`; engine=`%s`@`%s`",
		orUnknown(p.Environment), orUnknown(p.DatasetFingerprint), orUnknown(p.ClientImage),
		orUnknown(p.ClientImageDigest), orUnknown(p.EngineImage), orUnknown(p.EngineImageDigest))
}

// Generated, measurements-only reports for the join/CTE family. Interpretation
// belongs in the signed analysis beside them.

func joinCellText(c Cell) string {
	if c.Axis == "" {
		return "missing"
	}
	if !c.TimingMeasured || c.Status != StatusValid {
		return string(c.Status)
	}
	return fmt.Sprintf("%.4f", c.MedianMS)
}

func joinRatioText(b, a Cell) string {
	if a.MedianMS <= 0 || b.MedianMS <= 0 || !a.TimingMeasured || !b.TimingMeasured {
		return "N/A"
	}
	return fmt.Sprintf("%.4f", b.MedianMS/a.MedianMS)
}

func sameJoinPlan(a, b Cell) string {
	if a.Plan == "" || b.Plan == "" {
		return "N/A"
	}
	return fmt.Sprintf("%v", SameStructure(a.Plan, b.Plan))
}

// joinSampleStats returns the min, median and max of a cell's recorded samples,
// so a reader can see the spread the median hides. n is 0 for an untimed cell.
func joinSampleStats(c Cell) (min, med, max float64, n int) {
	if !c.TimingMeasured || c.Status != StatusValid || len(c.SamplesMS) == 0 {
		return 0, 0, 0, 0
	}
	vals := append([]float64(nil), c.SamplesMS...)
	sort.Float64s(vals)
	m := vals[len(vals)/2]
	if len(vals)%2 == 0 {
		m = (vals[len(vals)/2-1] + m) / 2
	}
	return vals[0], m, vals[len(vals)-1], len(vals)
}

func fmtStat(v float64) string {
	if v == 0 {
		return "0"
	}
	return fmt.Sprintf("%.4f", v)
}

// RenderJoinRun renders one whole join/CTE run. inputsDigest is the SHA-256 of
// the result JSON bytes, computed by the writer (it cannot live inside the file
// it hashes).
func RenderJoinRun(r JoinResult, inputsDigest string) (string, error) {
	// An early failure (feature check, connection, DDL) has no cells and no
	// feasibility evidence, but it still must produce a report: an error that
	// cannot be read back is not recorded evidence.
	if len(r.Cells) == 0 && r.Feasibility == nil && r.Error == "" {
		return "", fmt.Errorf("no cells and no feasibility evidence to render")
	}
	var b strings.Builder
	b.WriteString("# Study 01 SELECT join/filter and CTE materialization — generated measurements\n\n")
	b.WriteString("_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._\n\n")
	fmt.Fprintf(&b, "**TL;DR:** %d cells, %d valid, engine `%s`, topology `%s`, scale `%s`, seed %d, %d fact rows.\n\n",
		len(r.Cells), r.Valid(), r.Engine, r.Topology, r.Scale, r.Seed, r.FactCount)
	if r.Error != "" {
		fmt.Fprintf(&b, "**Run error:** %s\n\n", r.Error)
	}

	b.WriteString("## Provenance\n\n")
	p := r.Provenance
	fmt.Fprintf(&b, "- run id: `%s`\n- run tag: `%s`\n- environment: `%s`\n- repo commit: `%s`\n- engine version: `%s`\n",
		r.RunID, p.RunTag, p.Environment, p.RepoCommit, p.EngineVersion)
	fmt.Fprintf(&b, "- database engine image: `%s` @ `%s`\n", orUnknown(p.EngineImage), orUnknown(p.EngineImageDigest))
	fmt.Fprintf(&b, "- benchmark client image: `%s` @ `%s`\n", orUnknown(p.ClientImage), orUnknown(p.ClientImageDigest))
	fmt.Fprintf(&b, "- inputs digest: `%s`\n- dataset fingerprint: `%s`\n", orUnknown(inputsDigest), p.DatasetFingerprint)
	fmt.Fprintf(&b, "- client CPU quota: `%s`, memory: `%s`\n", strings.TrimSpace(p.ClientCPUQuota), strings.TrimSpace(p.ClientMemoryLimit))
	for _, k := range sortedKeys(p.Settings) {
		fmt.Fprintf(&b, "- setting `%s`: %s\n", k, p.Settings[k])
	}
	if len(r.Endpoints) > 0 {
		b.WriteString("\n### Endpoint distribution\n\nA cluster arm must rotate its reader across nodes; a single seed endpoint would not be balanced cluster access.\n\n| endpoint | operations | connections at end |\n|---|---|---|\n")
		for i, e := range r.Endpoints {
			ops, conns := 0, 0
			if i < len(r.EndpointOps) {
				ops = r.EndpointOps[i]
			}
			if i < len(r.EndpointConns) {
				conns = r.EndpointConns[i]
			}
			fmt.Fprintf(&b, "| %s | %d | %d |\n", e, ops, conns)
		}
	}
	b.WriteString("\n")

	fmt.Fprintf(&b, "## Fixture and selectivity mapping\n\n- arm: `%s`\n- fact rows: %d\n- parent rows: %d\n- rank thresholds: %v\n\nThe predicate is `rank <= threshold`, where `rank` is a deterministic shuffle of 1..%d; `rank` is a labelled proxy for the teaching selective predicate, not the same question. The `cte`/`cte-small` arm is the synthetic customers/facts microfixture; the `cte-domain` arm is Study 01's person/donation fixture with the person's charity projected and the views binding it to this catalogue.\n\n",
		orUnknown(r.Arm), r.FactCount, r.Customers, r.Thresholds, r.FactCount)

	if r.Feasibility != nil {
		renderJoinFeasibility(&b, *r.Feasibility)
	}

	b.WriteString("## P1 — six inner-join spellings and their CTE policies (median ms)\n\n")
	b.WriteString("| threshold | " + strings.Join(joinShortNames(), " | ") + " |\n")
	b.WriteString("|---|" + strings.Repeat("---|", len(JoinStatements)) + "\n")
	for _, t := range r.Thresholds {
		var cols []string
		for _, name := range JoinStatements {
			c, _ := r.Find(JoinBagAxis(name, t))
			cols = append(cols, joinCellText(c))
		}
		fmt.Fprintf(&b, "| %d (%s) | %s |\n", t, JoinThresholdPct(t, r.FactCount), strings.Join(cols, " | "))
	}
	b.WriteString("\n")

	b.WriteString("## P1 — bag spelling pairs against the canonical join+WHERE spelling\n\n")
	b.WriteString("| threshold | pair | A ms | B ms | B/A | same plan structure |\n|---|---|---|---|---|---|\n")
	for _, pr := range r.BagPairs() {
		fmt.Fprintf(&b, "| %d | %s | %.4f | %.4f | %s | %s |\n",
			pr.Threshold, pr.Label, pr.A.MedianMS, pr.B.MedianMS, joinRatioText(pr.B, pr.A), sameJoinPlan(pr.A, pr.B))
	}
	b.WriteString("\n")

	b.WriteString("## P3 — forced CTE policy against the default spelling\n\n")
	b.WriteString("| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |\n|---|---|---|---|---|---|---|\n")
	for _, pr := range r.PolicyPairs() {
		fmt.Fprintf(&b, "| %d | %s | %.4f | %.4f | %s | %s | %s |\n",
			pr.Threshold, pr.Label, pr.A.MedianMS, pr.B.MedianMS, joinRatioText(pr.B, pr.A), sameJoinPlan(pr.A, pr.B), pr.B.PolicyNote)
	}
	b.WriteString("\n")

	b.WriteString("## P4 — reuse ladder (median ms)\n\n")
	b.WriteString("| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |\n|---|---|---|---|---|---|---|---|\n")
	for _, policy := range []string{"default", "mat", "notmat"} {
		for _, refs := range ReuseLevels() {
			for _, pct := range ReuseConsumerPercents() {
				c, ok := r.Find(JoinReuseAxis(policy, refs, pct))
				if !ok {
					fmt.Fprintf(&b, "| %s | %d | %d%% | missing | | | | |\n", policy, refs, pct)
					continue
				}
				ms := joinCellText(c)
				refsSeen := "0"
				if c.Status == StatusValid {
					refsSeen = fmt.Sprint(c.ResultRows)
				}
				fmt.Fprintf(&b, "| %s | %d | %d%% | %s | %s | `%s` | %s | %s |\n",
					policy, refs, pct, ms, refsSeen, StructureSignature(c.Plan), c.PolicyNote, c.Status)
			}
		}
	}
	b.WriteString("\n")

	b.WriteString("## P4 — reuse policy pairs (B/A)\n\n| policy | refs | consumer | B/A | same plan structure |\n|---|---|---|---|---|\n")
	for _, pr := range r.ReusePairs() {
		fmt.Fprintf(&b, "| %s | %d | %d%% | %s | %s |\n", pr.Policy, pr.Refs, pr.Pct, joinRatioText(pr.B, pr.A), sameJoinPlan(pr.A, pr.B))
	}
	b.WriteString("\n")

	if pp, ok := r.JoinPositionPair(); ok {
		fmt.Fprintf(&b, "## Repeated-design position control\n\n| first median ms | last median ms | ratio (last/first) | same plan |\n|---|---|---|---|\n| %.6f | %.6f | %.4f | %s |\n\n",
			pp.A.MedianMS, pp.B.MedianMS, pp.RatioBA, sameJoinPlan(pp.A, pp.B))
	}

	b.WriteString("## Latency samples and derived throughput\n\n")
	b.WriteString("_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._\n\n")
	b.WriteString("| axis | status | min ms | median ms | max ms | samples | derived queries/s |\n|---|---|---|---|---|---|---|\n")
	statCells := append([]Cell(nil), r.Cells...)
	sort.Slice(statCells, func(i, j int) bool { return statCells[i].Axis < statCells[j].Axis })
	for _, c := range statCells {
		mn, med, mx, n := joinSampleStats(c)
		rate, mnS, medS, mxS := "N/A", "N/A", "N/A", "N/A"
		if n > 0 {
			mnS, medS, mxS = fmtStat(mn), fmtStat(med), fmtStat(mx)
			if med > 0 {
				rate = fmt.Sprintf("%.1f", 1000/med)
			}
		}
		fmt.Fprintf(&b, "| %s | %s | %s | %s | %s | %d | %s |\n",
			c.Axis, c.Status, mnS, medS, mxS, n, rate)
	}
	b.WriteString("\n")

	b.WriteString("## Matched blended workload\n\n")
	if r.Writer == nil {
		b.WriteString("No writer record: this run predates the phase field or is feasibility-only.\n\n")
	} else {
		wt := r.Writer
		fmt.Fprintf(&b, "- phase: `%s`; target offered rate: %d updates/s\n", orUnknown(wt.Phase), wt.TargetRate)
		fmt.Fprintf(&b, "- offered %d, committed %d, failed %d over %.3f s\n", wt.Offered, wt.Completed, wt.Failed, wt.WallSeconds)
		fmt.Fprintf(&b, "- offered rate %.3f/s, completed rate %.3f/s\n", wt.OfferedRate, wt.CompletedRate)
		if wt.Phase == "blended" {
			fmt.Fprintf(&b, "- reader operations overlapping a committed write: %d of %d\n", r.WriterOverlapSamples, r.WriterSamples)
			retries := 0
			for _, c := range r.Cells {
				retries += c.Retries
			}
			fmt.Fprintf(&b, "- reader retries of transient conflicts: %d\n", retries)
		}
		fmt.Fprintf(&b, "- per key: %s\n", joinPerKeySpread(wt.PerKey))
		fmt.Fprintf(&b, "- ledger: %d keys reconciled, %d mismatch", wt.LedgerKeys, wt.LedgerMismatch)
		if wt.LedgerDetail != "" {
			fmt.Fprintf(&b, " (%s)", wt.LedgerDetail)
		}
		b.WriteString("\n")
		if len(wt.LatencyMS) > 0 {
			fmt.Fprintf(&b, "- update latency: %s\n", joinLatencySummary(wt.LatencyMS))
		}
		if wt.Error != "" {
			fmt.Fprintf(&b, "- writer error: %s\n", wt.Error)
		}
		b.WriteString("\nA phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.\n\n")
	}

	b.WriteString("## Execution order\n\n")
	fmt.Fprintf(&b, "- order seed `%d`, %d interleaved rounds\n", r.Order.Seed, r.Order.Rounds)
	for i, axis := range r.Order.Sequence {
		fmt.Fprintf(&b, "  %d. %s\n", i+1, axis)
	}
	b.WriteString("\n")

	b.WriteString("## Live negative controls\n\n")
	if len(r.Controls) == 0 {
		b.WriteString("None recorded.\n\n")
	} else {
		b.WriteString("| control | pretended contract | detected | rows | note |\n|---|---|---|---|---|\n")
		for _, c := range r.Controls {
			fmt.Fprintf(&b, "| %s | %s | %v | %d | %s |\n", c.Statement, c.Group, c.Detected, c.Rows, c.Error)
		}
		b.WriteString("\n")
	}

	b.WriteString("## Excluded cells\n\n")
	var excluded []Cell
	for _, c := range r.Cells {
		if c.Status != StatusValid {
			excluded = append(excluded, c)
		}
	}
	if len(excluded) == 0 {
		b.WriteString("None: every cell is valid.\n\n")
	} else {
		for _, c := range excluded {
			fmt.Fprintf(&b, "- `%s` (%s) — **%s**%s\n", c.Axis, c.Statement, c.Status, dashJoin(c.Note))
		}
		b.WriteString("\nA cell that is invalid, unsupported or timed out is never ranked.\n\n")
	}

	b.WriteString("## Plan signatures\n\n| axis | statement | status | structure signature | policy note |\n|---|---|---|---|---|\n")
	cells := append([]Cell(nil), r.Cells...)
	sort.Slice(cells, func(i, j int) bool { return cells[i].Axis < cells[j].Axis })
	for _, c := range cells {
		fmt.Fprintf(&b, "| %s | %s | %s | `%s` | %s |\n", c.Axis, c.Statement, c.Status, StructureSignature(c.Plan), c.PolicyNote)
	}
	b.WriteString("\n")

	b.WriteString("## Recorded limits\n\n")
	if len(r.Notes) == 0 {
		b.WriteString("None recorded.\n")
	} else {
		for _, n := range r.Notes {
			fmt.Fprintf(&b, "- %s\n", n)
		}
	}
	return b.String(), nil
}

func renderJoinFeasibility(b *strings.Builder, f JoinFeasibility) {
	b.WriteString("## Per-engine feasibility evidence\n\n")
	fmt.Fprintf(b, "- fact rows %d, customers %d, thresholds %v\n", f.FactCount, f.Customers, f.Thresholds)
	if f.ReuseExprSupported {
		b.WriteString("- expensive reuse token supported by the engine: yes\n")
	} else {
		fmt.Fprintf(b, "- expensive reuse token supported by the engine: **no**%s\n", dashJoin(f.ReuseExprError))
	}
	keys := make([]string, 0, len(f.ThresholdCounts))
	for k := range f.ThresholdCounts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	b.WriteString("\n| threshold | rows with rank <= threshold |\n|---|---|\n")
	for _, k := range keys {
		fmt.Fprintf(b, "| %s | %d |\n", k, f.ThresholdCounts[k])
	}
	b.WriteString("\n| statement | support | rank index observed | forced policy evidence |\n|---|---|---|---|\n")
	names := make([]string, 0, len(f.StatementSupport))
	for k := range f.StatementSupport {
		names = append(names, k)
	}
	sort.Strings(names)
	for _, n := range names {
		fmt.Fprintf(b, "| %s | %s | %v | %s |\n", n, f.StatementSupport[n], f.RankIndexUsed[n], f.PolicyEvidence[n])
	}
	for _, n := range f.Notes {
		fmt.Fprintf(b, "- gap: %s\n", n)
	}
	b.WriteString("\n")
}

// RenderJoinReplicated summarises several whole runs, keyed on (engine,
// topology, axis) so environments, thresholds and contracts are never pooled.
func RenderJoinReplicated(runs []JoinResult) (string, error) {
	if len(runs) == 0 {
		return "", fmt.Errorf("no runs to summarise")
	}
	var b strings.Builder
	b.WriteString("# Study 01 SELECT join/filter and CTE materialization — replicated whole runs (generated measurements)\n\n")
	b.WriteString("_Measurements only. A stable sign across runs is what an analysis may cite; the verdict is assigned there._\n\n")
	fmt.Fprintf(&b, "**TL;DR:** %d whole runs summarised; each run's own median is shown and never pooled across compatibility classes, engines, topologies or thresholds.\n\n", len(runs))
	// Compatibility classes: runs whose inputs differ (environment, fixture,
	// images, settings) are partitioned, never pooled.
	classes := map[string][]string{}
	var classOrder []string
	for _, r := range runs {
		cls := joinClass(r)
		if _, ok := classes[cls]; !ok {
			classOrder = append(classOrder, cls)
		}
		classes[cls] = append(classes[cls], r.RunID)
	}
	b.WriteString("## Compatibility classes\n\nRuns are grouped by environment, dataset, client and engine images, and settings; rows below are never pooled across classes.\n\n| class | runs | inputs |\n|---|---|---|\n")
	for _, cls := range classOrder {
		desc := ""
		for _, r := range runs {
			if joinClass(r) == cls {
				desc = joinClassDescription(r)
				break
			}
		}
		fmt.Fprintf(&b, "| `%s` | %s | %s |\n", cls, strings.Join(classes[cls], ", "), desc)
	}
	b.WriteString("\n")

	b.WriteString("## Runs\n\n| run | class | engine | topology | scale | fact rows | seed |\n|---|---|---|---|---|---|---|\n")
	for _, r := range runs {
		fmt.Fprintf(&b, "| %s | `%s` | %s | %s | %s | %d | %d |\n",
			r.RunID, joinClass(r), r.Engine, r.Topology, r.Scale, r.FactCount, r.Seed)
	}
	b.WriteString("\n")

	type key struct{ class, engine, topology, axis string }
	seen := map[key]bool{}
	values := map[key][]float64{}
	statuses := map[key][]string{}
	var keys []key
	for _, r := range runs {
		cls := joinClass(r)
		for _, c := range r.Cells {
			k := key{cls, r.Engine, r.Topology, c.Axis}
			if !seen[k] {
				seen[k] = true
				keys = append(keys, k)
			}
			if c.TimingMeasured && c.Status == StatusValid {
				values[k] = append(values[k], c.MedianMS)
			} else {
				statuses[k] = append(statuses[k], string(c.Status))
			}
		}
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].class != keys[j].class {
			return keys[i].class < keys[j].class
		}
		if keys[i].engine != keys[j].engine {
			return keys[i].engine < keys[j].engine
		}
		if keys[i].topology != keys[j].topology {
			return keys[i].topology < keys[j].topology
		}
		return keys[i].axis < keys[j].axis
	})
	b.WriteString("## Cells across runs (per compatibility class, engine and topology)\n\n| class | engine | topology | axis | run medians ms | across-run median | spread % | valid runs |\n|---|---|---|---|---|---|---|---|\n")
	for _, k := range keys {
		medians := append([]float64(nil), values[k]...)
		var shown []string
		for _, m := range medians {
			shown = append(shown, fmt.Sprintf("%.4f", m))
		}
		shown = append(shown, statuses[k]...)
		med, spread := "N/A", "N/A"
		if len(medians) > 0 {
			sorted := append([]float64(nil), medians...)
			sort.Float64s(sorted)
			m := sorted[len(sorted)/2]
			if len(sorted)%2 == 0 {
				m = (sorted[len(sorted)/2-1] + m) / 2
			}
			med = fmt.Sprintf("%.4f", m)
			if m > 0 {
				spread = fmt.Sprintf("%.1f", 100*(sorted[len(sorted)-1]-sorted[0])/m)
			}
		}
		fmt.Fprintf(&b, "| `%s` | %s | %s | %s | %s | %s | %s | %d |\n",
			k.class, k.engine, k.topology, k.axis, strings.Join(shown, ", "), med, spread, len(medians))
	}
	b.WriteString("\n## Pair ratios per run (per engine)\n\n| run | engine | topology | contract | pair | ratio (B/A) | sign |\n|---|---|---|---|---|---|---|\n")
	sign := func(v float64) string {
		switch {
		case v > 1:
			return ">1"
		case v < 1:
			return "<1"
		default:
			return "="
		}
	}
	for _, r := range runs {
		for _, pr := range r.BagPairs() {
			if pr.RatioBA <= 0 {
				continue
			}
			fmt.Fprintf(&b, "| %s | %s | %s | bag@%d | %s | %.4f | %s |\n", r.RunID, r.Engine, r.Topology, pr.Threshold, pr.Label, pr.RatioBA, sign(pr.RatioBA))
		}
		for _, pr := range r.ReusePairs() {
			if pr.RatioBA <= 0 {
				continue
			}
			fmt.Fprintf(&b, "| %s | %s | %s | reuse n=%d c=%d%% | %s | %.4f | %s |\n", r.RunID, r.Engine, r.Topology, pr.Refs, pr.Pct, pr.Label, pr.RatioBA, sign(pr.RatioBA))
		}
		if pp, ok := r.JoinPositionPair(); ok && pp.RatioBA > 0 {
			fmt.Fprintf(&b, "| %s | %s | %s | position | %s | %.4f | %s |\n", r.RunID, r.Engine, r.Topology, pp.Label, pp.RatioBA, sign(pp.RatioBA))
		}
	}
	b.WriteString("\n## Recorded limits\n\n")
	seenNote := map[string]bool{}
	any := false
	for _, r := range runs {
		for _, n := range r.Notes {
			if seenNote[n] {
				continue
			}
			seenNote[n] = true
			any = true
			fmt.Fprintf(&b, "- [%s] %s\n", r.RunID, n)
		}
	}
	if !any {
		b.WriteString("None recorded.\n")
	}
	return b.String(), nil
}

func joinShortNames() []string {
	out := make([]string, 0, len(JoinStatements))
	for _, n := range JoinStatements {
		out = append(out, shortJoinStatement(n))
	}
	return out
}

func dashJoin(s string) string {
	if strings.TrimSpace(s) == "" {
		return ""
	}
	return " — " + s
}
