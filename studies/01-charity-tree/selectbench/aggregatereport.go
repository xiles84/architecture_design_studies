package selectbench

import (
	"fmt"
	"sort"
	"strings"
)

// RenderAggRun renders one whole aggregate/top-N run: measurements only, with
// the provenance, per-cell plans and controls that produced them. Verdicts
// belong in a signed analysis, never here.
func RenderAggRun(r AggResult) (string, error) {
	if len(r.Cells) == 0 {
		return "", fmt.Errorf("no cells to render")
	}
	var b strings.Builder
	b.WriteString("# Study 01 SELECT aggregate & top-N — generated measurements\n\n")
	b.WriteString("_Measurements only. Interpretation is in the signed analysis beside this report._\n\n")
	fmt.Fprintf(&b, "**TL;DR:** %d cells, %d valid, engine `%s`, topology `%s`, scale `%s`, seed %d.\n\n",
		len(r.Cells), r.Valid(), r.Engine, r.Topology, r.Scale, r.Seed)

	b.WriteString("## Provenance\n\n")
	p := r.Provenance
	fmt.Fprintf(&b, "- run id: `%s`\n- run tag: `%s`\n- environment: `%s`\n- repo commit: `%s`\n- engine version: `%s`\n- engine image: `%s` @ `%s`\n",
		r.RunID, p.RunTag, p.Environment, p.RepoCommit, p.EngineVersion, p.EngineImage, p.EngineImageDigest)
	fmt.Fprintf(&b, "- dataset fingerprint: `%s`\n- client CPU quota: `%s`, memory: `%s`\n",
		p.DatasetFingerprint, strings.TrimSpace(p.ClientCPUQuota), strings.TrimSpace(p.ClientMemoryLimit))
	if s := throttledLine(p.ClientCPUStatBefore, p.ClientCPUStatAfter); s != "" {
		fmt.Fprintf(&b, "- client throttling: %s\n", s)
	}
	b.WriteString("\n")

	b.WriteString("## Execution order\n\n")
	fmt.Fprintf(&b, "- order seed `%d`; interleaved rounds: %d\n", r.Order.Seed, r.Order.Rounds)
	b.WriteString("- executed sequence (one sample of every cell per round):\n")
	for i, axis := range r.Order.Sequence {
		fmt.Fprintf(&b, "  %d. %s\n", i+1, axis)
	}
	b.WriteString("\n")

	b.WriteString("## P8 — parent aggregate (per selected parent)\n\n")
	b.WriteString("Each row is one parent-selection level; ms are that run's own medians. `pre/post` is pre-aggregate over post-join, `corr/*` correlated over the other spelling, all at the same level and arguments.\n\n")
	b.WriteString("| level (parents) | post_join ms | pre_agg ms | correlated ms | pre/post | corr/post | corr/pre |\n|---|---|---|---|---|---|---|\n")
	for _, level := range r.Levels {
		a, _ := r.Find(AggregateTopNAxis("aggregate", "a01_aggregate_join_sel", level))
		bb, _ := r.Find(AggregateTopNAxis("aggregate", "a02_aggregate_preagg_sel", level))
		c, _ := r.Find(AggregateTopNAxis("aggregate", "a03_aggregate_correlated_sel", level))
		fmt.Fprintf(&b, "| %d | %s | %s | %s | %s | %s | %s |\n",
			level, aggCell(a), aggCell(bb), aggCell(c), ratio(bb, a), ratio(c, a), ratio(c, bb))
	}
	b.WriteString("\n")

	b.WriteString("## P9 — per-parent top-N at N=1\n\n")
	b.WriteString("All three spellings are compared at N=1, because DISTINCT ON is an N=1-only spelling. `lat/rn` is LATERAL over ROW_NUMBER, `do/*` DISTINCT ON over the other spelling.\n\n")
	b.WriteString("| level (parents) | rownumber ms | lateral ms | distinct_on ms | lat/rn | do/rn | do/lat |\n|---|---|---|---|---|---|---|\n")
	for _, level := range r.Levels {
		a, _ := r.Find(AggregateTopNAxis("topn", "a04_topn_rownumber_sel", level))
		bb, _ := r.Find(AggregateTopNAxis("topn", "a05_topn_lateral_sel", level))
		c, _ := r.Find(AggregateTopNAxis("topn", "a06_topn_distinct_on_sel", level))
		fmt.Fprintf(&b, "| %d | %s | %s | %s | %s | %s | %s |\n",
			level, aggCell(a), aggCell(bb), aggCell(c), ratio(bb, a), ratio(c, a), ratio(c, bb))
	}
	b.WriteString("\n")

	b.WriteString("## Controls\n\n")
	if len(r.Controls) == 0 {
		b.WriteString("None recorded.\n\n")
	} else {
		b.WriteString("| control | pretended contract | detected | rows | note |\n|---|---|---|---|---|\n")
		for _, c := range r.Controls {
			fmt.Fprintf(&b, "| %s | %s | %v | %d | %s |\n", c.Statement, c.Group, c.Detected, c.Rows, c.Error)
		}
		b.WriteString("\n")
	}

	if pp, ok := r.PositionPair(); ok {
		ratio := "N/A"
		if pp.A.MedianMS > 0 {
			ratio = fmt.Sprintf("%.4f", pp.B.MedianMS/pp.A.MedianMS)
		}
		b.WriteString("## Repeated-design position control\n\n")
		b.WriteString("The first matrix cell repeated at the end of the phase; a large last/first ratio is a session position effect, not a spelling effect.\n\n")
		b.WriteString("| first median ms | last median ms | ratio (last/first) | same plan |\n|---|---|---|---|\n")
		fmt.Fprintf(&b, "| %.6f | %.6f | %s | %v |\n\n", pp.A.MedianMS, pp.B.MedianMS, ratio,
			StructureSignature(pp.A.Plan) == StructureSignature(pp.B.Plan))
	}

	b.WriteString("## Plans\n\n| axis | status | plan signature |\n|---|---|---|\n")
	cells := append([]Cell(nil), r.Cells...)
	sort.Slice(cells, func(i, j int) bool { return cells[i].Axis < cells[j].Axis })
	for _, c := range cells {
		fmt.Fprintf(&b, "| %s | %s | `%s` |\n", c.Axis, c.Status, StructureSignature(c.Plan))
	}
	b.WriteString("\n## Recorded limits\n\n")
	if len(r.Notes) == 0 {
		b.WriteString("None recorded.\n")
	} else {
		for _, n := range r.Notes {
			fmt.Fprintf(&b, "- %s\n", n)
		}
	}
	if r.Error != "" {
		fmt.Fprintf(&b, "\n**Run error:** %s\n", r.Error)
	}
	return b.String(), nil
}

// RenderAggReplicated summarises several whole runs, keyed on (engine,
// topology, axis) so two environments are never pooled into one spread.
func RenderAggReplicated(runs []AggResult) (string, error) {
	if len(runs) == 0 {
		return "", fmt.Errorf("no runs to summarise")
	}
	var b strings.Builder
	b.WriteString("# Study 01 SELECT aggregate & top-N — replicated whole runs (generated measurements)\n\n")
	b.WriteString("_Measurements only. A stable sign across these runs is the evidence an analysis cites; it is assigned there._\n\n")
	fmt.Fprintf(&b, "**TL;DR:** %d whole runs summarised, each run's own median shown; never merged across engines.\n\n", len(runs))

	b.WriteString("## Runs\n\n| run | engine | topology | scale | seed |\n|---|---|---|---|---|\n")
	for _, r := range runs {
		fmt.Fprintf(&b, "| %s | %s | %s | %s | %d |\n", r.RunID, r.Engine, r.Topology, r.Scale, r.Seed)
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
	b.WriteString("## Cells across runs (per engine)\n\n| engine | topology | axis | run medians ms | across-run median | spread % | valid runs |\n|---|---|---|---|---|---|---|\n")
	for _, k := range keys {
		medians := append([]float64(nil), values[k]...)
		var shown []string
		for _, m := range medians {
			shown = append(shown, fmt.Sprintf("%.6f", m))
		}
		shown = append(shown, statuses[k]...)
		med, spread := "N/A", "N/A"
		if len(medians) > 0 {
			sort.Float64s(medians)
			m := medians[len(medians)/2]
			if len(medians)%2 == 0 {
				m = (medians[len(medians)/2-1] + m) / 2
			}
			med = fmt.Sprintf("%.6f", m)
			if m > 0 {
				spread = fmt.Sprintf("%.1f", 100*(medians[len(medians)-1]-medians[0])/m)
			}
		}
		fmt.Fprintf(&b, "| %s | %s | %s | %s | %s | %s | %d |\n",
			k.engine, k.topology, k.axis, strings.Join(shown, ", "), med, spread, len(medians))
	}
	b.WriteString("\n## Pair ratios per run (per engine)\n\n| run | engine | topology | family | level | pair | ratio (B/A) | sign |\n|---|---|---|---|---|---|---|---|\n")
	for _, r := range runs {
		all := append(r.AggregatePairs(), r.TopNPairs()...)
		if pp, ok := r.PositionPair(); ok {
			all = append(all, pp)
		}
		for _, p := range all {
			if p.RatioBA <= 0 {
				continue
			}
			sign := "="
			if p.RatioBA > 1 {
				sign = ">1"
			} else if p.RatioBA < 1 {
				sign = "<1"
			}
			fmt.Fprintf(&b, "| %s | %s | %s | %d | %s | %.4f | %s |\n", r.Engine, r.Topology, p.Family, p.Level, p.Label, p.RatioBA, sign)
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

func aggCell(c Cell) string {
	if c.Axis == "" {
		return "missing"
	}
	if !c.TimingMeasured || c.Status != StatusValid {
		return string(c.Status)
	}
	return fmt.Sprintf("%.6f", c.MedianMS)
}

func ratio(b, a Cell) string {
	if a.MedianMS <= 0 || b.MedianMS <= 0 || !a.TimingMeasured || !b.TimingMeasured {
		return "N/A"
	}
	return fmt.Sprintf("%.4f", b.MedianMS/a.MedianMS)
}

func throttledLine(before, after string) string {
	get := func(s, field string) string {
		for _, line := range strings.Split(s, "\n") {
			if strings.HasPrefix(line, field) {
				return strings.TrimSpace(line)
			}
		}
		return ""
	}
	b1, a1 := get(before, "nr_throttled"), get(after, "nr_throttled")
	if b1 == "" && a1 == "" {
		return ""
	}
	return fmt.Sprintf("before %q, after %q", b1, a1)
}
