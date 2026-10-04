package selectbench

import (
	"fmt"
	"sort"
	"strings"
)

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

// RenderJoinRun renders one whole join/CTE run. inputsDigest is the SHA-256 of
// the result JSON bytes, computed by the writer (it cannot live inside the file
// it hashes).
func RenderJoinRun(r JoinResult, inputsDigest string) (string, error) {
	if len(r.Cells) == 0 {
		return "", fmt.Errorf("no cells to render")
	}
	var b strings.Builder
	b.WriteString("# Study 01 SELECT join/filter and CTE materialization — generated measurements\n\n")
	b.WriteString("_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._\n\n")
	fmt.Fprintf(&b, "**TL;DR:** %d cells, %d valid, engine `%s`, topology `%s`, scale `%s`, seed %d, %d fact rows.\n\n",
		len(r.Cells), r.Valid(), r.Engine, r.Topology, r.Scale, r.Seed, r.FactCount)

	b.WriteString("## Provenance\n\n")
	p := r.Provenance
	fmt.Fprintf(&b, "- run id: `%s`\n- run tag: `%s`\n- environment: `%s`\n- repo commit: `%s`\n- engine version: `%s`\n- engine image: `%s` @ `%s`\n",
		r.RunID, p.RunTag, p.Environment, p.RepoCommit, p.EngineVersion, p.EngineImage, p.EngineImageDigest)
	fmt.Fprintf(&b, "- inputs digest: `%s`\n- dataset fingerprint: `%s`\n", orUnknown(inputsDigest), p.DatasetFingerprint)
	fmt.Fprintf(&b, "- client CPU quota: `%s`, memory: `%s`\n", strings.TrimSpace(p.ClientCPUQuota), strings.TrimSpace(p.ClientMemoryLimit))
	for _, k := range sortedKeys(p.Settings) {
		fmt.Fprintf(&b, "- setting `%s`: %s\n", k, p.Settings[k])
	}
	b.WriteString("\n")

	fmt.Fprintf(&b, "## Fixture and selectivity mapping\n\n- fact rows: %d\n- customers: %d\n- rank thresholds: %v\n\nThe predicate is `rank <= threshold`, where `rank` is a deterministic shuffle of 1..%d; `rank` is only a labelled proxy for the teaching amount predicate.\n\n",
		r.FactCount, r.Customers, r.Thresholds, r.FactCount)

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
	fmt.Fprintf(&b, "**TL;DR:** %d whole runs summarised; each run's own median is shown and never merged across engines, topologies or thresholds.\n\n", len(runs))
	b.WriteString("## Runs\n\n| run | engine | topology | scale | fact rows | seed |\n|---|---|---|---|---|---|\n")
	for _, r := range runs {
		fmt.Fprintf(&b, "| %s | %s | %s | %s | %d | %d |\n", r.RunID, r.Engine, r.Topology, r.Scale, r.FactCount, r.Seed)
	}
	b.WriteString("\n")

	type key struct{ engine, topology, axis string }
	seen := map[key]bool{}
	values := map[key][]float64{}
	statuses := map[key][]string{}
	var keys []key
	for _, r := range runs {
		for _, c := range r.Cells {
			k := key{r.Engine, r.Topology, c.Axis}
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
		if keys[i].engine != keys[j].engine {
			return keys[i].engine < keys[j].engine
		}
		if keys[i].topology != keys[j].topology {
			return keys[i].topology < keys[j].topology
		}
		return keys[i].axis < keys[j].axis
	})
	b.WriteString("## Cells across runs (per engine and topology)\n\n| engine | topology | axis | run medians ms | across-run median | spread % | valid runs |\n|---|---|---|---|---|---|---|\n")
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
		fmt.Fprintf(&b, "| %s | %s | %s | %s | %s | %s | %d |\n",
			k.engine, k.topology, k.axis, strings.Join(shown, ", "), med, spread, len(medians))
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
