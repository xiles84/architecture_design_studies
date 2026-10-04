// Independent HIGH review of saved join results; this starts no database.
package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Cell struct {
	Axis, Statement, Status, Plan string
	MedianMS, SpreadPct           float64
	Trials, ResultRows            int
	TimingMeasured                bool
	SamplesMS                     []float64
	PolicyNote                    string `json:"policy_note"`
}
type Result struct {
	RunID      string                                                                                    `json:"run_id"`
	Engine     string                                                                                    `json:"engine"`
	Cells      []Cell                                                                                    `json:"cells"`
	Summary    map[string]int                                                                            `json:"summary"`
	Provenance struct{ RepoCommit, RunTag, Environment, ClientCPUStatBefore, ClientCPUStatAfter string } `json:"provenance"`
	Order      struct {
		Seed     int64    `json:"seed"`
		Rounds   int      `json:"rounds"`
		Sequence []string `json:"sequence"`
	} `json:"order"`
	Controls []struct {
		Statement string `json:"statement"`
		Detected  bool   `json:"detected"`
	} `json:"controls"`
}

func must(ok bool, message string) {
	if !ok {
		panic(message)
	}
}
func read(path string) []byte {
	b, e := os.ReadFile(path)
	if e != nil {
		panic(e)
	}
	return []byte(strings.ReplaceAll(string(b), "\r\n", "\n"))
}
func median(a []float64) float64 {
	b := append([]float64(nil), a...)
	sort.Float64s(b)
	n := len(b)
	if n%2 == 0 {
		return (b[n/2-1] + b[n/2]) / 2
	}
	return b[n/2]
}
func close(a, b float64) bool { return math.Abs(a-b) < 1e-7 }
func find(r Result, axis string) Cell {
	for _, c := range r.Cells {
		if c.Axis == axis {
			return c
		}
	}
	panic("missing " + axis)
}
func bag(stmt string, t int) string { return fmt.Sprintf("family=bag;stmt=%s;threshold=%d", stmt, t) }
func reuse(policy string, n, p int) string {
	return fmt.Sprintf("family=reuse;policy=%s;refs=%d;consumer=%d%%", policy, n, p)
}

func main() {
	root := "studies/01-charity-tree/results"
	var runs []Result
	seeds := map[int64]bool{}
	fmt.Println("# Independent saved-input audit — HIGH GPT-6 / Codex desktop\n\nNo database starts or new timings. Hashes normalize checkout CRLF to committed LF.")
	fmt.Println("\n| run | input SHA-256 | cells | controls | samples >20% spread | last/first | actual guest memory kB |\n|---|---|---|---|---|---|---|")
	for _, id := range []string{"pg1", "pg2", "pg3", "yb1", "yb2", "yb3"} {
		dir := filepath.Join(root, "select-join-"+id+"-20261004")
		b := read(filepath.Join(dir, "join.json"))
		var r Result
		must(json.Unmarshal(b, &r) == nil, "JSON parse")
		digest := fmt.Sprintf("%x", sha256.Sum256(b))
		report := string(read(filepath.Join(dir, "join.md")))
		must(strings.Contains(report, digest), "original report digest mismatch: "+r.RunID)
		must(len(r.Cells) == 91 && r.Summary["valid"] == 91, "matrix size/validity")
		must(r.Order.Rounds == 5 && len(r.Order.Sequence) == 91, "order size/rounds")
		must(!seeds[r.Order.Seed], "duplicate seed")
		seeds[r.Order.Seed] = true
		axes := map[string]bool{}
		for _, a := range r.Order.Sequence {
			must(!axes[a], "duplicate order axis")
			axes[a] = true
		}
		must(r.Order.Sequence[90] == r.Order.Sequence[0]+";position=last", "not actual first repeated")
		noisy := 0
		for _, c := range r.Cells {
			must(c.Status == "valid" && c.TimingMeasured && c.Trials == 5 && len(c.SamplesMS) == 5, "invalid samples")
			must(close(median(c.SamplesMS), c.MedianMS), "median mismatch")
			must(axes[c.Axis], "cell missing from order")
			if c.SpreadPct > 20 {
				noisy++
			}
			if strings.Contains(c.Axis, "policy=mat;") || strings.Contains(c.Statement, "_mat") {
				must(strings.Contains(c.Plan, "CTE Scan on"), "missing forced materialization")
			}
		}
		must(len(r.Controls) == 2, "control count")
		for _, c := range r.Controls {
			must(c.Detected, "control failed")
		}
		first, last := find(r, r.Order.Sequence[0]), find(r, r.Order.Sequence[90])
		memory := strings.Fields(strings.Split(string(read(filepath.Join(dir, "resources/guest-memory.txt"))), "\n")[0])[1]
		fmt.Printf("| %s | `%s` | 91 | 2/2 | %d/91 | %.4f | %s |\n", r.RunID, digest, noisy, last.MedianMS/first.MedianMS, memory)
		must(strings.Contains(string(read(filepath.Join(dir, "manifest.env"))), "repo_commit="+r.Provenance.RepoCommit), "manifest commit")
		runs = append(runs, r)
	}
	fmt.Println("\n## Materialized joined-output / default joined-output (paired per-run ratios)\n\n| engine | threshold | ratios in run order |\n|---|---|---|")
	for _, engine := range []string{"postgres", "yugabyte"} {
		for _, t := range []int{20, 200, 2000, 20000, 100000, 180000} {
			var values []string
			for _, r := range runs {
				if r.Engine == engine {
					a, b := find(r, bag("cte_joined_outer", t)), find(r, bag("cte_joined_outer_mat", t))
					values = append(values, fmt.Sprintf("%.4f", b.MedianMS/a.MedianMS))
				}
			}
			fmt.Printf("| %s | %d | %s |\n", engine, t, strings.Join(values, ", "))
		}
	}
	fmt.Println("\n## NOT MATERIALIZED / MATERIALIZED (actual P4 pair, not /default)\n\n| engine | references | consumer id-space % | ratios in run order |\n|---|---|---|---|")
	for _, engine := range []string{"postgres", "yugabyte"} {
		for _, n := range []int{1, 2, 4, 8, 16} {
			for _, p := range []int{1, 100} {
				var values []string
				for _, r := range runs {
					if r.Engine == engine {
						a, b := find(r, reuse("mat", n, p)), find(r, reuse("notmat", n, p))
						values = append(values, fmt.Sprintf("%.4f", b.MedianMS/a.MedianMS))
					}
				}
				fmt.Printf("| %s | %d | %d | %s |\n", engine, n, p, strings.Join(values, ", "))
			}
		}
	}
	fmt.Println("\n## Broad default-spelling pair ratios to join+WHERE\n\n| engine | threshold | overall min | overall max |\n|---|---|---|---|")
	for _, engine := range []string{"postgres", "yugabyte"} {
		for _, t := range []int{100000, 180000} {
			lo, hi := math.Inf(1), 0.0
			for _, r := range runs {
				if r.Engine != engine {
					continue
				}
				a := find(r, bag("join_where", t))
				for _, s := range []string{"join_on", "filtered_derived", "cte_filtered", "cte_joined_outer", "comma_join"} {
					v := find(r, bag(s, t)).MedianMS / a.MedianMS
					lo = math.Min(lo, v)
					hi = math.Max(hi, v)
				}
			}
			fmt.Printf("| %s | %d | %.4f | %.4f |\n", engine, t, lo, hi)
		}
	}
	fmt.Println("\nAll saved-input arithmetic, original input digests, seeded orders, forced materialization and control checks passed. This does not accept absent workload/topology coverage or the author's interpretation.")
}
