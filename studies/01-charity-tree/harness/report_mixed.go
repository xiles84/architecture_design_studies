package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ---------------------------------------------------------------------------
// Experiment D report — what happens when reads and writes share a database.
//
// The headline is not throughput. It is the ratio of each split against the
// all-readers baseline: how much of a design's read advantage survives once
// writers are working on the same rows.
// ---------------------------------------------------------------------------

func splitLabel(m *MixedResult) string {
	return fmt.Sprintf("%d:%d", m.ReadWorkers, m.WriteWorkers)
}

func writeMixedReport(resultsDir, outPath string) error {
	rs, err := loadResults(resultsDir)
	if err != nil {
		return err
	}
	var b strings.Builder
	ref := rs.any()

	fmt.Fprintf(&b, "# Experiment D — reads and writes running concurrently\n\n")
	fmt.Fprintf(&b, "| | |\n|---|---|\n")
	fmt.Fprintf(&b, "| Run id | `%s` |\n", ref.RunID)
	fmt.Fprintf(&b, "| Environment | [`%s`](../../../docs/environments/%s.md) |\n", ref.Environment, ref.Environment)
	fmt.Fprintf(&b, "| Scale | `%s` — %v people, %v donations |\n", ref.Scale, ref.Dataset["people"], ref.Dataset["donations"])

	fmt.Fprintf(&b, "\n## Why this experiment exists\n\n")
	fmt.Fprintf(&b, "Every other measurement in this study runs reads alone, then writes alone. That\n")
	fmt.Fprintf(&b, "gives clean attribution, but it cannot see the costs that only exist when the two\n")
	fmt.Fprintf(&b, "overlap: MVCC bloat accumulating *while* readers scan, contention on rows writers\n")
	fmt.Fprintf(&b, "keep updating, autovacuum and compaction stealing CPU, and buffer-cache\n")
	fmt.Fprintf(&b, "competition.\n\n")
	fmt.Fprintf(&b, "All of those penalise designs that buy read speed with redundancy — which is most\n")
	fmt.Fprintf(&b, "of the designs here. **Isolated measurement therefore flatters exactly the designs\n")
	fmt.Fprintf(&b, "under scrutiny**, and this is where that bias gets checked instead of assumed away.\n\n")

	// mix weights, once
	if len(ref.Mixed) > 0 {
		fmt.Fprintf(&b, "### The workload\n\n")
		fmt.Fprintf(&b, "A donor portal plus a dashboard. Weights sum to 100.\n\n")
		fmt.Fprintf(&b, "| Read query | Weight | | Write op | Weight |\n|---|---:|---|---|---:|\n")
		rkeys := sortedByWeight(ref.Mixed[0].ReadMix)
		wkeys := sortedByWeight(ref.Mixed[0].WriteMix)
		for i := 0; i < len(rkeys); i++ {
			wcell := "| |"
			if i < len(wkeys) {
				wcell = fmt.Sprintf("| %s | %d |", wkeys[i], ref.Mixed[0].WriteMix[wkeys[i]])
			}
			fmt.Fprintf(&b, "| %s | %d | %s\n", queryQuestion[rkeys[i]], ref.Mixed[0].ReadMix[rkeys[i]], wcell)
		}
		fmt.Fprintf(&b, "\nSplits are **readers:writers**. The first is all-readers and is the baseline —\n")
		fmt.Fprintf(&b, "same query mix, same data, same process, so the only variable across splits is\n")
		fmt.Fprintf(&b, "the presence of writers. The dataset is reloaded before every split, so no split\n")
		fmt.Fprintf(&b, "inherits the bloat the previous one created.\n\n")
	}

	for _, t := range rs.topologies {
		fmt.Fprintf(&b, "## %s\n\n", topologyLabel[t])
		fmt.Fprintf(&b, "| Design | Split | Reads/s | **Read throughput kept** | Writes/s | Read p99 (ms) | Aggregates still correct |\n")
		fmt.Fprintf(&b, "|---|---|---:|---:|---:|---:|---|\n")
		for _, d := range rs.designsIn(t) {
			r := rs.runs[t][d]
			if len(r.Mixed) == 0 {
				continue
			}
			baseline := r.Mixed[0].ReadOpsPerSec
			for _, m := range r.Mixed {
				kept := "baseline"
				if m.WriteWorkers > 0 && baseline > 0 {
					kept = fmt.Sprintf("**%.0f%%**", m.ReadOpsPerSec/baseline*100)
				}
				// p99 of the dominant query, which is the one a user actually waits on.
				p99 := "—"
				if l, ok := m.PerQuery["q09_person_recent_donations"]; ok {
					p99 = fmtMS(l.P99MS)
				}
				audit := "—"
				if m.Audit != nil && m.Audit.Ran {
					bad := m.Audit.PersonMismatches + m.Audit.CharityMismatches + m.Audit.CacheMismatches
					if bad == 0 {
						audit = "yes"
					} else {
						audit = fmt.Sprintf("**NO — %d rows wrong**", bad)
					}
				}
				fmt.Fprintf(&b, "| %s | %s | %s | %s | %s | %s | %s |\n",
					designShort[d], splitLabel(m), fmtOps(m.ReadOpsPerSec), kept,
					fmtOps(m.WriteOpsPerSec), p99, audit)
			}
		}
		fmt.Fprintf(&b, "\n")
	}

	// Per-query degradation: which question suffers most when writers appear.
	fmt.Fprintf(&b, "## Which questions suffer most under write load\n\n")
	fmt.Fprintf(&b, "Percentage of the all-readers throughput retained at the heaviest write split.\n")
	fmt.Fprintf(&b, "A design can hold its aggregate throughput while one specific question collapses,\n")
	fmt.Fprintf(&b, "and the blended number would never show it.\n\n")
	for _, t := range rs.topologies {
		ds := rs.designsIn(t)
		if len(ds) == 0 {
			continue
		}
		fmt.Fprintf(&b, "### %s\n\n| Question | ", topologyLabel[t])
		for _, d := range ds {
			fmt.Fprintf(&b, "%s | ", designShort[d])
		}
		fmt.Fprintf(&b, "\n|---|")
		for range ds {
			fmt.Fprintf(&b, "---:|")
		}
		fmt.Fprintf(&b, "\n")

		var qnames []string
		seen := map[string]bool{}
		for _, d := range ds {
			for _, m := range rs.runs[t][d].Mixed {
				for q := range m.PerQueryOps {
					if !seen[q] {
						seen[q] = true
						qnames = append(qnames, q)
					}
				}
			}
		}
		sort.Strings(qnames)

		for _, q := range qnames {
			fmt.Fprintf(&b, "| %s | ", queryQuestion[q])
			for _, d := range ds {
				mx := rs.runs[t][d].Mixed
				if len(mx) < 2 {
					fmt.Fprintf(&b, "— | ")
					continue
				}
				base := mx[0].PerQueryOps[q]
				last := mx[len(mx)-1].PerQueryOps[q]
				if base <= 0 || last <= 0 {
					fmt.Fprintf(&b, "— | ")
					continue
				}
				fmt.Fprintf(&b, "%.0f%% | ", last/base*100)
			}
			fmt.Fprintf(&b, "\n")
		}
		fmt.Fprintf(&b, "\n")
	}

	fmt.Fprintf(&b, "> Readers and writers share the same eight host cores here, so some of the loss\n")
	fmt.Fprintf(&b, "> below is simply CPU being taken by the writers rather than interference in the\n")
	fmt.Fprintf(&b, "> database. The comparison that survives that caveat is **between designs at the\n")
	fmt.Fprintf(&b, "> same split** — they all lose the same CPU, so a design that keeps less of its\n")
	fmt.Fprintf(&b, "> throughput than another is losing it to something other than scheduling.\n")

	writeProvenance(&b, resultsDir, filepath.Join(filepath.Dir(outPath), "analyses"), ref.RunID)

	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(outPath, []byte(b.String()), 0o644)
}

func sortedByWeight(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if m[keys[i]] != m[keys[j]] {
			return m[keys[i]] > m[keys[j]]
		}
		return keys[i] < keys[j]
	})
	return keys
}
