package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ---------------------------------------------------------------------------
// Experiment E report — the same designs under two regimes.
//
// Every design in this study has a limitation, and every limitation has a regime
// in which it stops mattering: bounded history for the embedded design, many small
// charities for the hot rollup row, a small dataset for the unindexed design. A
// conclusion that says "D6 is bad" without saying "except when..." is incomplete
// in exactly the way that leads people to reject the right design for their
// domain.
//
// This report lays two runs side by side -- a baseline regime and an alternative
// one -- and answers two questions per design:
//
//   1. How much did the regime change THIS design?          (row-wise ratio)
//   2. Did the regime change which design is best?          (rank per question)
//
// The second question is the one that ends up in a conclusion.
// ---------------------------------------------------------------------------

func writeCompareReport(dirA, dirB, labelA, labelB, outPath string) error {
	a, err := loadResults(dirA)
	if err != nil {
		return fmt.Errorf("%s: %w", dirA, err)
	}
	bset, err := loadResults(dirB)
	if err != nil {
		return fmt.Errorf("%s: %w", dirB, err)
	}
	if labelA == "" {
		labelA = filepath.Base(dirA)
	}
	if labelB == "" {
		labelB = filepath.Base(dirB)
	}

	var b strings.Builder
	refA, refB := a.any(), bset.any()

	fmt.Fprintf(&b, "# Regime comparison — %s vs %s\n\n", labelA, labelB)
	fmt.Fprintf(&b, "| | %s | %s |\n|---|---|---|\n", labelA, labelB)
	fmt.Fprintf(&b, "| Run id | `%s` | `%s` |\n", refA.RunID, refB.RunID)
	fmt.Fprintf(&b, "| Environment | `%s` | `%s` |\n", refA.Environment, refB.Environment)
	fmt.Fprintf(&b, "| Scale | `%s` | `%s` |\n", refA.Scale, refB.Scale)
	fmt.Fprintf(&b, "| Profile | `%v` | `%v` |\n", refA.Dataset["profile"], refB.Dataset["profile"])
	fmt.Fprintf(&b, "| Charities / people / donations | %v / %v / %v | %v / %v / %v |\n",
		refA.Dataset["charities"], refA.Dataset["people"], refA.Dataset["donations"],
		refB.Dataset["charities"], refB.Dataset["people"], refB.Dataset["donations"])
	fmt.Fprintf(&b, "| Max donations per person | %v | %v |\n",
		refA.Dataset["max_donations_person"], refB.Dataset["max_donations_person"])

	if refA.Environment != refB.Environment {
		fmt.Fprintf(&b, "\n> ⚠️ **Different environments.** These runs are not comparable; see methodology rule 2.\n")
	}
	if fmt.Sprint(refA.Dataset["donations"]) != fmt.Sprint(refB.Dataset["donations"]) {
		fmt.Fprintf(&b, "\n> Donation counts differ slightly between the two datasets (the capped generator\n")
		fmt.Fprintf(&b, "> stops at the first donor that reaches the uncapped total, and extra charities\n")
		fmt.Fprintf(&b, "> consume extra random draws). The difference is a few rows in ~100 000, far\n")
		fmt.Fprintf(&b, "> below any effect reported here.\n")
	}

	for _, t := range topologyOrder {
		da, db := a.runs[t], bset.runs[t]
		if da == nil || db == nil {
			continue
		}
		var ds []string
		for _, d := range designOrder {
			if da[d] != nil && db[d] != nil {
				ds = append(ds, d)
			}
		}
		if len(ds) == 0 {
			continue
		}

		fmt.Fprintf(&b, "\n## %s\n\n", topologyLabel[t])

		// -- per-design summary ---------------------------------------------
		fmt.Fprintf(&b, "### How much the regime changed each design\n\n")
		fmt.Fprintf(&b, "Ratios are **%s ÷ %s**: above 1.0x the alternative regime helped this design.\n\n", labelB, labelA)
		fmt.Fprintf(&b, "| Design | Read score | ")
		for _, op := range writeOpOrder {
			fmt.Fprintf(&b, "%s | ", writeOpLabel[op])
		}
		fmt.Fprintf(&b, "\n|---|---:|")
		for range writeOpOrder {
			fmt.Fprintf(&b, "---:|")
		}
		fmt.Fprintf(&b, "\n")
		for _, d := range ds {
			ra, rb := da[d], db[d]
			fmt.Fprintf(&b, "| %s | %s | ", designShort[d], cmpRatio(geomeanReads(ra), geomeanReads(rb)))
			for _, op := range writeOpOrder {
				wa, wb := writeOf(ra, op), writeOf(rb, op)
				if wa == nil || wb == nil || wa.Skipped != "" || wb.Skipped != "" {
					fmt.Fprintf(&b, "— | ")
					continue
				}
				fmt.Fprintf(&b, "%s | ", cmpRatio(wa.OpsPerSec, wb.OpsPerSec))
			}
			fmt.Fprintf(&b, "\n")
		}

		// -- who wins each question, in each regime -------------------------
		fmt.Fprintf(&b, "\n### Which design is fastest, per question, in each regime\n\n")
		fmt.Fprintf(&b, "A regime that changes the **winner** is the result that belongs in a\n")
		fmt.Fprintf(&b, "conclusion. Rows where the winner changed are marked **→**.\n\n")
		fmt.Fprintf(&b, "| Question | Best in %s | Best in %s |\n|---|---|---|\n", labelA, labelB)
		for _, q := range sortedQueries(a, t) {
			wa, va := bestFor(ds, da, func(r *Run) float64 {
				if x := readOf(r, q); x != nil {
					return x.OpsPerSec
				}
				return 0
			})
			wb, vb := bestFor(ds, db, func(r *Run) float64 {
				if x := readOf(r, q); x != nil {
					return x.OpsPerSec
				}
				return 0
			})
			mark := ""
			if wa != wb {
				mark = " **→**"
			}
			fmt.Fprintf(&b, "| %s | %s (%s) | %s (%s)%s |\n",
				queryQuestion[q], designShort[wa], fmtOps(va), designShort[wb], fmtOps(vb), mark)
		}
		for _, op := range writeOpOrder {
			pick := func(r *Run) float64 {
				if w := writeOf(r, op); w != nil && w.Skipped == "" {
					return w.OpsPerSec
				}
				return 0
			}
			wa, va := bestFor(ds, da, pick)
			wb, vb := bestFor(ds, db, pick)
			if wa == "" && wb == "" {
				continue
			}
			mark := ""
			if wa != wb {
				mark = " **→**"
			}
			fmt.Fprintf(&b, "| write: %s | %s (%s) | %s (%s)%s |\n",
				writeOpLabel[op], designShort[wa], fmtOps(va), designShort[wb], fmtOps(vb), mark)
		}

		// -- full per-question ratios ---------------------------------------
		fmt.Fprintf(&b, "\n### Per-question change (%s ÷ %s)\n\n| Question | ", labelB, labelA)
		for _, d := range ds {
			fmt.Fprintf(&b, "%s | ", designShort[d])
		}
		fmt.Fprintf(&b, "\n|---|")
		for range ds {
			fmt.Fprintf(&b, "---:|")
		}
		fmt.Fprintf(&b, "\n")
		for _, q := range sortedQueries(a, t) {
			fmt.Fprintf(&b, "| %s | ", queryQuestion[q])
			for _, d := range ds {
				qa, qb := readOf(da[d], q), readOf(db[d], q)
				if qa == nil || qb == nil {
					fmt.Fprintf(&b, "— | ")
					continue
				}
				fmt.Fprintf(&b, "%s | ", cmpRatio(qa.OpsPerSec, qb.OpsPerSec))
			}
			fmt.Fprintf(&b, "\n")
		}
	}

	fmt.Fprintf(&b, "\n> Both runs are single-trial surveys unless their manifests say otherwise, so the\n")
	fmt.Fprintf(&b, "> ratios carry the same error bar as the generated reports they came from (the\n")
	fmt.Fprintf(&b, "> D4/D5 control put that at up to ~1.5x). Treat a change under 1.5x as noise and a\n")
	fmt.Fprintf(&b, "> change in winner between two designs within 1.5x of each other as a tie.\n")

	writeProvenance(&b, dirB, filepath.Join(filepath.Dir(outPath), "analyses"), refB.RunID)

	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(outPath, []byte(b.String()), 0o644)
}

func cmpRatio(a, b float64) string {
	if a <= 0 || b <= 0 {
		return "—"
	}
	r := b / a
	switch {
	case r >= 1.5:
		return fmt.Sprintf("**%.1fx**", r)
	case r <= 1/1.5:
		return fmt.Sprintf("**%.2fx**", r)
	default:
		return fmt.Sprintf("%.2fx", r)
	}
}

func bestFor(ds []string, runs map[string]*Run, pick func(*Run) float64) (string, float64) {
	best, bv := "", 0.0
	for _, d := range ds {
		if v := pick(runs[d]); v > bv {
			best, bv = d, v
		}
	}
	return best, bv
}
