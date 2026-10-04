package selectbench

import (
	"fmt"
	"sort"
	"strings"
)

// RenderExistenceRun renders one whole membership/absence run. Measurements
// only; the contract-equality verdict and the P6/P7 calls belong in the signed
// analysis beside it.
func RenderExistenceRun(r ExistenceResult) (string, error) {
	if len(r.Cells) == 0 {
		return "", fmt.Errorf("no cells to render")
	}
	var b strings.Builder
	b.WriteString("# Study 01 SELECT membership & absence — generated measurements\n\n")
	b.WriteString("_Measurements only. Interpretation is in the signed analysis beside this report._\n\n")
	fmt.Fprintf(&b, "**TL;DR:** %d cells, %d valid, engine `%s`, topology `%s`, scale `%s`, seed %d.\n\n",
		len(r.Cells), r.Valid(), r.Engine, r.Topology, r.Scale, r.Seed)

	b.WriteString("## Provenance\n\n")
	p := r.Provenance
	fmt.Fprintf(&b, "- run id: `%s`\n- run tag: `%s`\n- environment: `%s`\n- repo commit: `%s`\n- engine version: `%s`\n- engine image: `%s` @ `%s`\n",
		r.RunID, p.RunTag, p.Environment, p.RepoCommit, p.EngineVersion, p.EngineImage, p.EngineImageDigest)
	fmt.Fprintf(&b, "- client CPU quota: `%s`, memory: `%s`\n", strings.TrimSpace(p.ClientCPUQuota), strings.TrimSpace(p.ClientMemoryLimit))
	b.WriteString("\n## Fanout levels (constant 10,000 child rows)\n\n| fanout | parents | children | dataset fingerprint |\n|---|---|---|---|\n")
	for _, l := range r.Levels {
		fmt.Fprintf(&b, "| %d | %d | %d | `%s` |\n", l.Fanout, l.Parents, l.Children, l.Fingerprint)
	}
	b.WriteString("\n")

	b.WriteString("## Execution order (per fanout phase)\n\n")
	for _, l := range r.Levels {
		fmt.Fprintf(&b, "- fanout %d: order seed `%d`, %d interleaved rounds\n", l.Fanout, l.Order.Seed, l.Order.Rounds)
		for i, axis := range l.Order.Sequence {
			fmt.Fprintf(&b, "  %d. %s\n", i+1, axis)
		}
	}
	b.WriteString("\n")

	b.WriteString("## P6 — membership (parent identity set)\n\nEach row is one fanout level; ms are this run's medians. `in/exists`, `jd/exists` and `jd/in` are the B spelling over A at the same level and predicate.\n\n")
	b.WriteString("| fanout | exists ms | in ms | in_notnull ms | join_distinct ms | in/exists | jd/exists | jd/in |\n|---|---|---|---|---|---|---|---|\n")
	for _, f := range r.Fanouts() {
		ex, _ := r.Find(ExistenceAxis("membership", "s06_exists_parent", f))
		in, _ := r.Find(ExistenceAxis("membership", "s07_in_parent", f))
		inn, _ := r.Find(ExistenceAxis("membership", "e12_in_notnull", f))
		jd, _ := r.Find(ExistenceAxis("membership", "s08_join_distinct_parent", f))
		fmt.Fprintf(&b, "| %d | %s | %s | %s | %s | %s | %s | %s |\n",
			f, exCell(ex), exCell(in), exCell(inn), exCell(jd), exRatio(in, ex), exRatio(jd, ex), exRatio(jd, in))
	}
	b.WriteString("\n")

	b.WriteString("## P7 — absence (parent identity set)\n\nAll three spellings use the NOT-NULL subquery contract; the nullable `NOT IN` case is the x02 control. Every valid answer here passed the oracle gate, so the sets are equal within contract.\n\n")
	b.WriteString("| fanout | not_exists ms | left_is_null ms | not_in_notnull ms |\n|---|---|---|---|\n")
	for _, f := range r.Fanouts() {
		ne, _ := r.Find(ExistenceAxis("absence", "s09_not_exists", f))
		la, _ := r.Find(ExistenceAxis("absence", "s10_left_anti", f))
		ni, _ := r.Find(ExistenceAxis("absence", "e13_not_in_notnull", f))
		fmt.Fprintf(&b, "| %d | %s | %s | %s |\n", f, exCell(ne), exCell(la), exCell(ni))
	}
	b.WriteString("\n")

	if pp, ok := r.PositionPair(); ok {
		fmt.Fprintf(&b, "## Repeated-design position control\n\n| first median ms | last median ms | ratio (last/first) | same plan |\n|---|---|---|---|\n| %.6f | %.6f | %s | %v |\n\n",
			pp.A.MedianMS, pp.B.MedianMS, exRatio(pp.B, pp.A), StructureSignature(pp.A.Plan) == StructureSignature(pp.B.Plan))
	}

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
	return b.String(), nil
}

// RenderExistenceReplicated summarises several whole runs, keyed on (engine,
// topology, axis) so environments and fanout levels are never pooled.
func RenderExistenceReplicated(runs []ExistenceResult) (string, error) {
	if len(runs) == 0 {
		return "", fmt.Errorf("no runs to summarise")
	}
	var b strings.Builder
	b.WriteString("# Study 01 SELECT membership & absence — replicated whole runs (generated measurements)\n\n")
	b.WriteString("_Measurements only. A stable sign across these runs is the evidence an analysis cites; it is assigned there._\n\n")
	fmt.Fprintf(&b, "**TL;DR:** %d whole runs summarised, each run's own median shown, never merged across engines or fanout levels.\n\n", len(runs))
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
	b.WriteString("\n## Membership pair ratios per run (per engine)\n\n| run | engine | topology | fanout | pair | ratio (B/A) | sign |\n|---|---|---|---|---|---|---|\n")
	for _, r := range runs {
		for _, p := range r.MembershipPairs() {
			if p.RatioBA <= 0 {
				continue
			}
			sign := "="
			if p.RatioBA > 1 {
				sign = ">1"
			} else if p.RatioBA < 1 {
				sign = "<1"
			}
			fmt.Fprintf(&b, "| %s | %s | %s | %d | %s | %.4f | %s |\n", r.RunID, r.Engine, r.Topology, p.Fanout, p.Label, p.RatioBA, sign)
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

func exCell(c Cell) string {
	if c.Axis == "" {
		return "missing"
	}
	if !c.TimingMeasured || c.Status != StatusValid {
		return string(c.Status)
	}
	return fmt.Sprintf("%.6f", c.MedianMS)
}

func exRatio(b, a Cell) string {
	if a.MedianMS <= 0 || b.MedianMS <= 0 || !a.TimingMeasured || !b.TimingMeasured {
		return "N/A"
	}
	return fmt.Sprintf("%.4f", b.MedianMS/a.MedianMS)
}
