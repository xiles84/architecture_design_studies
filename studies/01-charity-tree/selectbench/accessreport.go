package selectbench

// Reports for the access-path matrix. A generated report holds measurements and
// no interpretation: verdicts on P10/P11 belong in the signed analysis.

import (
	"fmt"
	"sort"
	"strings"
)

func msOrNA(c Cell) string {
	if c.TimingMeasured && c.Status == StatusValid {
		return fmt.Sprintf("%.6f", c.MedianMS)
	}
	return "N/A"
}

func spreadOrNA(c Cell) string {
	if c.TimingMeasured && c.Status == StatusValid {
		return fmt.Sprintf("%.1f", c.SpreadPct)
	}
	return "N/A"
}

func writeAccessProvenance(b *strings.Builder, p Provenance) {
	b.WriteString("| Field | Value |\n|---|---|\n")
	fmt.Fprintf(b, "| repo commit | `%s` |\n", orUnknown(p.RepoCommit))
	fmt.Fprintf(b, "| benchmark image | `%s` |\n", orUnknown(p.EngineImage))
	fmt.Fprintf(b, "| image digest | `%s` |\n", orUnknown(p.ImageDigest))
	fmt.Fprintf(b, "| run tag | `%s` |\n", orUnknown(p.RunTag))
	fmt.Fprintf(b, "| environment | `%s` |\n", orUnknown(p.Environment))
	fmt.Fprintf(b, "| engine image digest | `%s` |\n", orUnknown(p.EngineImageDigest))
	fmt.Fprintf(b, "| engine version | %s |\n", orUnknown(p.EngineVersion))
	fmt.Fprintf(b, "| dataset seed | %d |\n", p.Seed)
	fmt.Fprintf(b, "| dataset fingerprint | `%s` |\n", orUnknown(p.DatasetFingerprint))
	for _, k := range sortedKeys(p.Settings) {
		fmt.Fprintf(b, "| setting: %s | %s |\n", k, p.Settings[k])
	}
	b.WriteString("\n")
}

// RenderAccessRun writes the per-whole-run measurements.
func RenderAccessRun(r AccessResult) (string, error) {
	var b strings.Builder
	b.WriteString("# SELECT access-path matrix (generated measurements)\n\n")
	b.WriteString("_Measurements only. Verdicts for P10/P11 belong in a signed analysis._\n\n")
	b.WriteString("**TL;DR:** within-run medians for OR / UNION ALL / UNION at three predicate shapes and ")
	b.WriteString("OFFSET / keyset at five page depths, each cell verified against the independent oracle. ")
	b.WriteString("A within-run difference is not a replicated result. Missing measurements are N/A.\n\n")

	fmt.Fprintf(&b, "| run | engine | topology | scale | page limit |\n|---|---|---|---|---|\n")
	fmt.Fprintf(&b, "| %s | %s | %s | %s | %d |\n\n", r.RunID, r.Engine, r.Topology, r.Scale, r.Limit)

	b.WriteString("## Provenance\n\n")
	writeAccessProvenance(&b, r.Provenance)

	b.WriteString("## Execution order\n\n")
	fmt.Fprintf(&b, "- default matrix: order seed `%d`, %d interleaved rounds\n", r.Order.Seed, r.Order.Rounds)
	fmt.Fprintf(&b, "- expression-index arm: order seed `%d`, %d interleaved rounds\n", r.IndexOrder.Seed, r.IndexOrder.Rounds)
	b.WriteString("- default-matrix sequence (one sample of every cell per round):\n")
	for i, axis := range r.Order.Sequence {
		fmt.Fprintf(&b, "  %d. %s\n", i+1, axis)
	}
	if len(r.IndexOrder.Sequence) > 0 {
		b.WriteString("- index-arm sequence:\n")
		for i, axis := range r.IndexOrder.Sequence {
			fmt.Fprintf(&b, "  %d. %s\n", i+1, axis)
		}
	}
	b.WriteString("\n")

	b.WriteString("## Disjunction disjointness proof\n\n")
	b.WriteString("UNION ALL is a competing spelling of OR only while the two arms cannot both match one row. ")
	b.WriteString("The count is computed by the independent Go oracle from the deterministic fixture.\n\n")
	b.WriteString("| shape | arm A | arm B | A matches | B matches | overlapping rows | disjoint |\n|---|---|---|---|---|---|---|\n")
	for _, o := range r.Overlap {
		fmt.Fprintf(&b, "| %s | `%s` | `%s` | %d | %d | %d | %t |\n",
			o.Shape, o.ASQL, o.BSQL, o.ArmAMatches, o.ArmBMatches, o.OverlapCount, o.Disjoint)
	}
	b.WriteString("\n")

	b.WriteString("## Predicate alternatives (P10)\n\n")
	b.WriteString("| shape | index arm | spelling | status | median ms | spread % | trials | rows | plan structure |\n|---|---|---|---|---|---|---|---|---|\n")
	for _, c := range sortCells(r.Cells) {
		if !strings.Contains(c.Axis, "spelling=") || strings.Contains(c.Axis, "depth=") || strings.Contains(c.Axis, "position=") {
			continue
		}
		fmt.Fprintf(&b, "| %s | %s | %s | %s | %s | %s | %d | %d | `%s` |\n",
			axisShape(c.Axis), axisIndexArm(c.Axis), axisSpelling(c.Axis), c.Status, msOrNA(c), spreadOrNA(c), c.Trials, c.ResultRows, c.structureSig())
	}
	b.WriteString("\n### P10 pairs\n\n")
	b.WriteString("Pairs compare the **default** arm only; the labelled `expression_index` arm is a package arm and is never pooled with the query-text pair.\n\n")
	writePairTable(&b, r.Pairs(), "shape=")

	b.WriteString("## Pagination (P11)\n\n")
	b.WriteString("| depth | spelling | status | median ms | spread % | trials | rows | plan structure |\n|---|---|---|---|---|---|---|---|\n")
	var depthRatios []string
	for _, d := range r.Depths {
		for _, spelling := range []string{"offset", "keyset"} {
			stmt := "s16_page_offset"
			if spelling == "keyset" {
				stmt = "s17_page_keyset"
			}
			if c, ok := r.Find(stmt, "depth="+itoa(d)+";spelling="+spelling); ok {
				fmt.Fprintf(&b, "| %d | %s | %s | %s | %s | %d | %d | `%s` |\n",
					d, spelling, c.Status, msOrNA(c), spreadOrNA(c), c.Trials, c.ResultRows, c.structureSig())
			}
		}
		if o, ok := r.Find("s16_page_offset", "depth="+itoa(d)+";spelling=offset"); ok {
			if k, ok2 := r.Find("s17_page_keyset", "depth="+itoa(d)+";spelling=keyset"); ok2 {
				if o.TimingMeasured && k.TimingMeasured && o.MedianMS > 0 {
					depthRatios = append(depthRatios, fmt.Sprintf("| %d | %.4f | %.6f | %.6f |", d, k.MedianMS/o.MedianMS, o.MedianMS, k.MedianMS))
				} else {
					depthRatios = append(depthRatios, fmt.Sprintf("| %d | N/A | %s | %s |", d, msOrNA(o), msOrNA(k)))
				}
			}
		}
	}
	b.WriteString("\nMedian ratio table (keyset / OFFSET; a ratio below 1 means keyset was faster in this run):\n\n")
	b.WriteString("| depth | ratio (keyset/offset) | offset median ms | keyset median ms |\n|---|---|---|---|\n")
	if len(depthRatios) == 0 {
		b.WriteString("| — | N/A | N/A | N/A |\n")
	} else {
		for _, row := range depthRatios {
			b.WriteString(row + "\n")
		}
	}
	b.WriteString("\n### P11 pairs\n\n")
	writePairTable(&b, r.Pairs(), "depth=")

	b.WriteString("## Repeated-design position control\n\n")
	if pp, ok := r.PositionPair(); ok {
		b.WriteString("The first matrix cell and a repeat of it at the end of the phase (`position=last`); a large ratio is a session position effect, not a spelling effect.\n\n")
		b.WriteString("| pair | first median ms | last median ms | ratio (last/first) | same plan |\n|---|---|---|---|---|\n")
		ratio := "N/A"
		if pp.RatioBA > 0 {
			ratio = fmt.Sprintf("%.4f", pp.RatioBA)
		}
		fmt.Fprintf(&b, "| %s | %s | %s | %s | %t |\n\n", pp.Label, msOrNA(pp.A), msOrNA(pp.B), ratio, pp.SamePlan)
	} else {
		b.WriteString("Not measured in this run.\n\n")
	}

	b.WriteString("## Expression-index arm: write and storage cost\n\n")
	ic := r.Index
	if ic.IndexName == "" {
		b.WriteString("Not attempted on this engine/topology.\n\n")
	} else {
		fmt.Fprintf(&b, "Index `%s` on expression `%s`; supported=%t.\n\n", ic.IndexName, ic.Expression, ic.Supported)
		b.WriteString("| measure | value |\n|---|---|\n")
		if ic.SizeBytes != nil {
			fmt.Fprintf(&b, "| index size (bytes) | %d |\n", *ic.SizeBytes)
		} else {
			b.WriteString("| index size (bytes) | N/A |\n")
		}
		if ic.WithIndexMS != nil {
			fmt.Fprintf(&b, "| insert batch with index (ms) | %.6f |\n", *ic.WithIndexMS)
		} else {
			b.WriteString("| insert batch with index (ms) | N/A |\n")
		}
		if ic.WithoutIndexMS != nil {
			fmt.Fprintf(&b, "| insert batch without index (ms) | %.6f |\n", *ic.WithoutIndexMS)
		} else {
			b.WriteString("| insert batch without index (ms) | N/A |\n")
		}
		fmt.Fprintf(&b, "| batch rows | %d |\n", ic.Rows)
		if ic.Note != "" {
			fmt.Fprintf(&b, "| note | %s |\n", ic.Note)
		}
		b.WriteString("\n")
	}

	b.WriteString("## Fixed-snapshot page equality\n\n")
	sc := r.Snapshot
	fmt.Fprintf(&b, "Isolation: `%s`; all depths equal: %t.\n\n", orUnknown(sc.Isolation), sc.AllEqual)
	if sc.Error != "" {
		fmt.Fprintf(&b, "Snapshot check error: %s\n\n", sc.Error)
	}
	b.WriteString("| depth | offset rows | keyset rows | equal |\n|---|---|---|---|\n")
	if len(sc.Depths) == 0 {
		b.WriteString("| — | N/A | N/A | N/A |\n")
	} else {
		for _, d := range sc.Depths {
			fmt.Fprintf(&b, "| %d | %d | %d | %t |\n", d.Depth, d.OffsetN, d.KeysetN, d.Equal)
		}
	}
	b.WriteString("\n")

	b.WriteString("## Live negative controls (per predicate shape)\n\n")
	b.WriteString("`x05_union_all_overlap` is this family's required control and must be detected in every shape. ")
	b.WriteString("`x01`-`x04` belong to other contracts; they are gated live on the teaching fixture, and here they are recorded as ")
	b.WriteString("`control_absent` when the measurement fixture cannot exercise them (for example `x01` needs an unmatched parent, ")
	b.WriteString("which a fixture where every parent has orders does not contain). An absent control is not evidence for the cells it cannot police.\n\n")
	if len(r.Controls) == 0 {
		b.WriteString("None recorded.\n\n")
	} else {
		for _, c := range r.Controls {
			fmt.Fprintf(&b, "- shape=%s `%s`: detected=%t, rows=%d%s\n", c.Shape, c.Statement, c.Detected, c.Rows, dash(c.Error))
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
			fmt.Fprintf(&b, "- `%s` [%s] — **%s**%s\n", c.Statement, c.Axis, c.Status, dash(c.Note))
		}
		b.WriteString("\nA cell that is invalid, unsupported or timed out is never ranked below.\n\n")
	}

	b.WriteString("## Captured plans and result sizes\n\n")
	for _, c := range sortCells(r.Cells) {
		fmt.Fprintf(&b, "### %s [%s]\n\nRows: %d; serialized decoded-row JSON bytes: %d (not wire bytes). Planning/execution: %s / %s ms.\n\n```text\n%s```\n\n",
			c.Statement, c.Axis, c.ResultRows, c.ResultJSONBytes, optionalMS(c.PlanningMS), optionalMS(c.ExecutionMS), c.Plan+"\n")
	}

	if len(r.Notes) > 0 {
		b.WriteString("## Limits and gaps recorded by the runner\n\n")
		for _, n := range r.Notes {
			fmt.Fprintf(&b, "- %s\n", n)
		}
		b.WriteString("\n")
	}
	if r.Error != "" {
		fmt.Fprintf(&b, "## Run error\n\n%s\n\n", r.Error)
	}
	return b.String(), nil
}

func writePairTable(b *strings.Builder, pairs []PairDelta, prefix string) {
	b.WriteString("| pair | A median ms | B median ms | B/A | same plan | within-run verdict |\n|---|---|---|---|---|---|\n")
	any := false
	for _, p := range pairs {
		if !strings.HasPrefix(p.Label, prefix) {
			continue
		}
		any = true
		ratio := "N/A"
		if p.RatioBA > 0 {
			ratio = fmt.Sprintf("%.4f", p.RatioBA)
		}
		verdict := "incomplete"
		if p.Explanation != "" {
			verdict = p.Explanation
		}
		fmt.Fprintf(b, "| %s | %s | %s | %s | %t | %s |\n",
			p.Label, msOrNA(p.A), msOrNA(p.B), ratio, p.SamePlan, verdict)
	}
	if !any {
		b.WriteString("| — | N/A | N/A | N/A | false | no comparable pair recorded |\n")
	}
	b.WriteString("\n")
}

func axisShape(axis string) string {
	for _, part := range strings.Split(axis, ";") {
		if strings.HasPrefix(part, "shape=") {
			return strings.TrimPrefix(part, "shape=")
		}
	}
	return ""
}

func axisSpelling(axis string) string {
	for _, part := range strings.Split(axis, ";") {
		if strings.HasPrefix(part, "spelling=") {
			return strings.TrimPrefix(part, "spelling=")
		}
	}
	return ""
}

// axisIndexArm names the physical arm: "default" is the query-text baseline,
// anything else is a package arm (currently only `expression_index`).
func axisIndexArm(axis string) string {
	for _, part := range strings.Split(axis, ";") {
		if strings.HasPrefix(part, "index=") {
			return strings.TrimPrefix(part, "index=")
		}
	}
	return "default"
}

// RenderAccessReplicated writes the cross-whole-run view. It makes no verdict:
// it reports each cell's across-run median/spread and each pair's per-run ratio
// and sign, which is what a signed analysis needs to call a direction stable.
func RenderAccessReplicated(runs []AccessResult) (string, error) {
	if len(runs) == 0 {
		return "", fmt.Errorf("no runs to summarise")
	}
	var b strings.Builder
	b.WriteString("# SELECT access-path matrix — replicated whole runs (generated measurements)\n\n")
	b.WriteString("_Measurements only. A stable sign across these runs is the evidence an analysis cites; it is assigned there._\n\n")
	fmt.Fprintf(&b, "**TL;DR:** %d whole runs summarised. Values are each run's own median; the across-run spread is shown, never merged into one number.\n\n", len(runs))

	b.WriteString("## Runs\n\n| run | engine | topology | scale | seed |\n|---|---|---|---|---|\n")
	for _, r := range runs {
		fmt.Fprintf(&b, "| %s | %s | %s | %s | %d |\n", r.RunID, r.Engine, r.Topology, r.Scale, r.Seed)
	}
	b.WriteString("\n")

	// Collect every (engine, topology, axis) cell in stable order. Engines are
	// never pooled: an across-run spread that mixed two engines would describe
	// no environment.
	type cellKey struct{ engine, topology, axis string }
	seen := map[cellKey]bool{}
	values := map[cellKey][]float64{}
	statuses := map[cellKey][]string{}
	var keys []cellKey
	for _, r := range runs {
		for _, c := range r.Cells {
			k := cellKey{r.Engine, r.Topology, c.Axis}
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
	b.WriteString("\n")

	b.WriteString("## Pair ratios per run (per engine)\n\n| engine | topology | pair | run ids | ratios (B/A) | signs |\n|---|---|---|---|---|---|\n")
	type pairKey struct{ engine, topology, label string }
	type agg struct {
		ids    []string
		ratios []string
		signs  []string
	}
	pairs := map[pairKey]*agg{}
	var pairOrder []pairKey
	pairSeen := map[pairKey]bool{}
	for _, r := range runs {
		pairsOfRun := r.Pairs()
		if pp, ok := r.PositionPair(); ok {
			pairsOfRun = append(pairsOfRun, pp)
		}
		for _, p := range pairsOfRun {
			if p.RatioBA <= 0 {
				continue
			}
			k := pairKey{r.Engine, r.Topology, p.Label}
			if pairs[k] == nil {
				pairs[k] = &agg{}
			}
			if !pairSeen[k] {
				pairSeen[k] = true
				pairOrder = append(pairOrder, k)
			}
			sign := "="
			if p.RatioBA > 1 {
				sign = ">1"
			} else if p.RatioBA < 1 {
				sign = "<1"
			}
			pairs[k].ratios = append(pairs[k].ratios, fmt.Sprintf("%.4f", p.RatioBA))
			pairs[k].signs = append(pairs[k].signs, sign)
		}
	}
	sort.Slice(pairOrder, func(i, j int) bool {
		if pairOrder[i].engine != pairOrder[j].engine {
			return pairOrder[i].engine < pairOrder[j].engine
		}
		if pairOrder[i].topology != pairOrder[j].topology {
			return pairOrder[i].topology < pairOrder[j].topology
		}
		return pairOrder[i].label < pairOrder[j].label
	})
	for _, k := range pairOrder {
		a := pairs[k]
		fmt.Fprintf(&b, "| %s | %s | %s | %s | %s | %s |\n", k.engine, k.topology, k.label, strings.Join(a.ids, ", "), strings.Join(a.ratios, ", "), strings.Join(a.signs, ", "))
	}
	b.WriteString("\n")

	b.WriteString("## Recorded limits\n\n")
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
