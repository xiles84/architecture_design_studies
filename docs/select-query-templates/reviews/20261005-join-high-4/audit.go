// Fourth independent HIGH review of the join/CTE submission (v4 topology runs).
// It reads saved results and generated reports only; it starts no database.
//
//	podman run --rm --network none -v "<worktree>:/repo:ro" -w /repo \
//	  docker.io/library/golang:1.26-bookworm \
//	  go run ./docs/select-query-templates/reviews/20261005-join-high-4/audit.go
package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const results = "studies/01-charity-tree/results"

type cell struct {
	Axis, Status string
	SamplesMS    []float64
}

type run struct {
	Cells    []cell `json:"cells"`
	Controls []struct {
		Detected bool `json:"detected"`
	} `json:"controls"`
	EndpointOps []int `json:"endpoint_ops"`
	Writer      *struct {
		Offered        int   `json:"offered"`
		Completed      int   `json:"completed"`
		Failed         int   `json:"failed"`
		EndpointOps    []int `json:"endpoint_ops"`
		LedgerMismatch int   `json:"ledger_mismatch"`
	} `json:"writer"`
	Placement []struct {
		Table       string `json:"table"`
		NumTablets  int    `json:"num_tablets"`
		IsColocated bool   `json:"is_colocated"`
	} `json:"placement"`
	Notes []string `json:"notes"`
}

type throttle struct {
	Node          string `json:"node"`
	NrPeriods     int    `json:"nr_periods"`
	NrThrottled   int    `json:"nr_throttled"`
	ThrottledUsec int    `json:"throttled_usec"`
}

func load(id string) (run, map[string]string, string) {
	var r run
	b, err := os.ReadFile(filepath.Join(results, id, "join.json"))
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(b, &r); err != nil {
		panic(err)
	}
	sum := sha256.Sum256(b)
	m := map[string]string{}
	if f, err := os.Open(filepath.Join(results, id, "manifest.env")); err == nil {
		s := bufio.NewScanner(f)
		for s.Scan() {
			if k, v, ok := strings.Cut(s.Text(), "="); ok {
				m[k] = v
			}
		}
		f.Close()
	}
	return r, m, hex.EncodeToString(sum[:])
}

func median(xs []float64) float64 {
	s := append([]float64(nil), xs...)
	sort.Float64s(s)
	n := len(s)
	if n%2 == 1 {
		return s[n/2]
	}
	return (s[n/2-1] + s[n/2]) / 2
}

func medians(r run) map[string]float64 {
	out := map[string]float64{}
	for _, c := range r.Cells {
		if c.Status == "valid" && len(c.SamplesMS) > 0 {
			out[c.Axis] = median(c.SamplesMS)
		}
	}
	return out
}

var leaderRE = regexp.MustCompile(`(yb-n[0-9]):[0-9]+`)

// leadersAfter returns relation -> distinct leader nodes from the post-run
// yb-admin listing (sections headed "## <relation>").
func leadersAfter(id string) map[string][]string {
	out := map[string][]string{}
	f, err := os.Open(filepath.Join(results, id, "resources", "placement.txt"))
	if err != nil {
		return out
	}
	defer f.Close()
	rel := ""
	seen := map[string]map[string]bool{}
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(line, "## ") {
			rel = strings.TrimPrefix(line, "## ")
			continue
		}
		if !regexp.MustCompile(`^[0-9a-f]{32}\s`).MatchString(line) {
			continue
		}
		if m := leaderRE.FindStringSubmatch(line); m != nil {
			if seen[rel] == nil {
				seen[rel] = map[string]bool{}
			}
			if !seen[rel][m[1]] {
				seen[rel][m[1]] = true
				out[rel] = append(out[rel], m[1])
			}
		}
	}
	for k := range out {
		sort.Strings(out[k])
	}
	return out
}

func leadersBefore(id string) []string {
	b, err := os.ReadFile(filepath.Join(results, id, "resources", "placement-leaders-before.txt"))
	if err != nil {
		return nil
	}
	seen := map[string]bool{}
	var out []string
	for _, m := range leaderRE.FindAllStringSubmatch(string(b), -1) {
		if !seen[m[1]] {
			seen[m[1]] = true
			out = append(out, m[1])
		}
	}
	sort.Strings(out)
	return out
}

func exists(id, name string) bool {
	_, err := os.Stat(filepath.Join(results, id, "resources", name))
	return err == nil
}

var arms = []string{"single", "equal", "tablet", "distributed", "colocated"}

func runID(arm, phase string, k int) string {
	return fmt.Sprintf("select-join-v4-topo-%s-%s%d-20261005", arm, phase, k)
}

// perRun returns the three per-run values of f for one arm and phase.
func perRun(arm, phase string, f func(map[string]float64) float64) []float64 {
	var out []float64
	for k := 1; k <= 3; k++ {
		r, _, _ := load(runID(arm, phase, k))
		out = append(out, f(medians(r)))
	}
	return out
}

func axis(stmt, t string) string { return "family=bag;stmt=" + stmt + ";threshold=" + t }

func get(stmt, t string) func(map[string]float64) float64 {
	return func(m map[string]float64) float64 { return m[axis(stmt, t)] }
}

func ratio(num, den string, t string) func(map[string]float64) float64 {
	return func(m map[string]float64) float64 { return m[axis(num, t)] / m[axis(den, t)] }
}

func fmtList(xs []float64, prec int) string {
	var s []string
	for _, x := range xs {
		s = append(s, fmt.Sprintf("%.*f", prec, x))
	}
	return strings.Join(s, ", ")
}

func minmax(xs []float64) (float64, float64) {
	lo, hi := xs[0], xs[0]
	for _, x := range xs {
		if x < lo {
			lo = x
		}
		if x > hi {
			hi = x
		}
	}
	return lo, hi
}

// separation says whether every run of a lies below (or above) every run of b.
func separation(a, b []float64) string {
	alo, ahi := minmax(a)
	blo, bhi := minmax(b)
	switch {
	case ahi < blo:
		return "a<b all runs"
	case alo > bhi:
		return "a>b all runs"
	default:
		return "overlap"
	}
}

func main() {
	fmt.Println("# Audit output — fourth HIGH join/CTE review (v4 topology runs)")
	fmt.Println()

	var ids []string
	entries, err := os.ReadDir(results)
	if err != nil {
		panic(err)
	}
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "select-join-v4-topo-") {
			ids = append(ids, e.Name())
		}
	}
	sort.Strings(ids)

	fmt.Println("## 1. Every v4 run: validity, samples, controls, build, digest")
	fmt.Println()
	fmt.Println("| run | topology/budget/phase/placement | valid/cells | 5-sample cells | controls fired | digest in join.md |")
	fmt.Println("|---|---|---|---|---|---|")
	builds, commits := map[string]int{}, map[string]int{}
	for _, id := range ids {
		r, m, sum := load(id)
		valid, five := 0, 0
		for _, c := range r.Cells {
			if c.Status == "valid" {
				valid++
			}
			if len(c.SamplesMS) == 5 {
				five++
			}
		}
		fired := 0
		for _, c := range r.Controls {
			if c.Detected {
				fired++
			}
		}
		rep, _ := os.ReadFile(filepath.Join(results, id, "join.md"))
		builds[m["client_image_id"]]++
		commits[m["repo_commit"]]++
		fmt.Printf("| %s | %s/%s/%s/%s | %d/%d | %d | %d/%d | %v |\n", strings.TrimSuffix(strings.TrimPrefix(id, "select-join-v4-topo-"), "-20261005"),
			m["topology"], m["node_budget"], m["phase"], m["placement"], valid, len(r.Cells), five, fired, len(r.Controls), strings.Contains(string(rep), sum))
	}
	fmt.Println()
	fmt.Printf("Runs: %d. Distinct client image ids: %d. Distinct repo commits: %d.\n\n", len(ids), len(builds), len(commits))

	fmt.Println("## 2. Placement evidence (cluster runs)")
	fmt.Println()
	fmt.Println("| run | tablets (yb_table_properties) | leaders after run (placement.txt) | donation leaders before | gate files | gate note |")
	fmt.Println("|---|---|---|---|---|---|")
	for _, id := range ids {
		_, m, _ := load(id)
		if m["topology"] != "cluster3" {
			continue
		}
		r, _, _ := load(id)
		var tabs []string
		for _, p := range r.Placement {
			if strings.HasSuffix(p.Table, "_pkey") {
				continue
			}
			c := ""
			if p.IsColocated {
				c = "c"
			}
			tabs = append(tabs, fmt.Sprintf("%s=%d%s", p.Table[strings.LastIndex(p.Table, ".")+1:], p.NumTablets, c))
		}
		la := leadersAfter(id)
		var rels []string
		for k := range la {
			rels = append(rels, k)
		}
		sort.Strings(rels)
		var ls []string
		for _, k := range rels {
			ls = append(ls, k+":"+strings.Join(la[k], "+"))
		}
		gate := ""
		if exists(id, "placement-prepared") {
			gate += "prepared "
		}
		if exists(id, "placement-gate.ok") {
			gate += "ok"
		}
		note := false
		for _, n := range r.Notes {
			if strings.Contains(n, "distributed placement gate") {
				note = true
			}
		}
		fmt.Printf("| %s | %s | %s | %s | %s | %v |\n", strings.TrimSuffix(strings.TrimPrefix(id, "select-join-v4-topo-"), "-20261005"),
			strings.Join(tabs, " "), strings.Join(ls, " "), strings.Join(leadersBefore(id), "+"), gate, note)
	}
	fmt.Println()

	for _, phase := range []string{"i", "b"} {
		fmt.Printf("## 3%s. join_where per-run medians (ms), phase %s\n\n", phase, map[string]string{"i": "isolated", "b": "blended"}[phase])
		fmt.Println("| arm | t=20 | t=2000 | t=180000 | joined MAT/default t=2000 | MAT t=2000 | MAT t=180000 |")
		fmt.Println("|---|---|---|---|---|---|---|")
		for _, a := range arms {
			fmt.Printf("| %s | %s | %s | %s | %s | %s | %s |\n", a,
				fmtList(perRun(a, phase, get("join_where", "20")), 2),
				fmtList(perRun(a, phase, get("join_where", "2000")), 2),
				fmtList(perRun(a, phase, get("join_where", "180000")), 1),
				fmtList(perRun(a, phase, ratio("cte_joined_outer_mat", "cte_joined_outer", "2000")), 2),
				fmtList(perRun(a, phase, get("cte_joined_outer_mat", "2000")), 1),
				fmtList(perRun(a, phase, get("cte_joined_outer_mat", "180000")), 1))
		}
		fmt.Println()
		fmt.Println("Across-run medians (median of the three run medians) and ratios to the single node:")
		fmt.Println()
		fmt.Println("| arm | t=20 | t=2000 | t=180000 | /single t=20 | /single t=2000 | /single t=180000 | MAT t=180000 /single | MAT t=2000 /single |")
		fmt.Println("|---|---|---|---|---|---|---|---|---|")
		base := map[string]float64{}
		for _, a := range arms {
			v20 := median(perRun(a, phase, get("join_where", "20")))
			v2k := median(perRun(a, phase, get("join_where", "2000")))
			v180 := median(perRun(a, phase, get("join_where", "180000")))
			m180 := median(perRun(a, phase, get("cte_joined_outer_mat", "180000")))
			m2k := median(perRun(a, phase, get("cte_joined_outer_mat", "2000")))
			if a == "single" {
				base["20"], base["2000"], base["180000"], base["m180"], base["m2k"] = v20, v2k, v180, m180, m2k
			}
			fmt.Printf("| %s | %.2f | %.2f | %.2f | %.2f | %.2f | %.2f | %.2f | %.2f |\n", a, v20, v2k, v180,
				v20/base["20"], v2k/base["2000"], v180/base["180000"], m180/base["m180"], m2k/base["m2k"])
		}
		fmt.Println()
		fmt.Println("Run-range separation of the analysis's pairs (a vs b):")
		fmt.Println()
		fmt.Println("| pair (a vs b) | t=20 | t=2000 | t=180000 | MAT/default t=2000 |")
		fmt.Println("|---|---|---|---|---|")
		for _, p := range [][2]string{{"distributed", "tablet"}, {"colocated", "distributed"}, {"distributed", "single"}, {"colocated", "single"}, {"tablet", "single"}, {"equal", "single"}} {
			fmt.Printf("| %s vs %s | %s | %s | %s | %s |\n", p[0], p[1],
				separation(perRun(p[0], phase, get("join_where", "20")), perRun(p[1], phase, get("join_where", "20"))),
				separation(perRun(p[0], phase, get("join_where", "2000")), perRun(p[1], phase, get("join_where", "2000"))),
				separation(perRun(p[0], phase, get("join_where", "180000")), perRun(p[1], phase, get("join_where", "180000"))),
				separation(perRun(p[0], phase, ratio("cte_joined_outer_mat", "cte_joined_outer", "2000")), perRun(p[1], phase, ratio("cte_joined_outer_mat", "cte_joined_outer", "2000"))))
		}
		fmt.Println()
	}

	fmt.Println("## 4. Blended phase: endpoints, writer and ledger")
	fmt.Println()
	fmt.Println("| run | reader endpoint ops | writer offered/completed/failed | writer endpoint ops | ledger mismatch |")
	fmt.Println("|---|---|---|---|---|")
	for _, id := range ids {
		r, m, _ := load(id)
		if m["phase"] != "blended" || r.Writer == nil {
			continue
		}
		fmt.Printf("| %s | %v | %d/%d/%d | %v | %d |\n", strings.TrimSuffix(strings.TrimPrefix(id, "select-join-v4-topo-"), "-20261005"),
			r.EndpointOps, r.Writer.Offered, r.Writer.Completed, r.Writer.Failed, r.Writer.EndpointOps, r.Writer.LedgerMismatch)
	}
	fmt.Println()

	fmt.Println("## 5. Node CPU throttling over each run (throttled periods / periods; throttled seconds)")
	fmt.Println()
	fmt.Println("| run | per node |")
	fmt.Println("|---|---|")
	for _, id := range ids {
		b, err := os.ReadFile(filepath.Join(results, id, "resources", "node-throttle.json"))
		if err != nil {
			fmt.Printf("| %s | not recorded |\n", id)
			continue
		}
		var ts []throttle
		_ = json.Unmarshal(b, &ts)
		var parts []string
		for _, t := range ts {
			parts = append(parts, fmt.Sprintf("%s %d/%d %.2fs", t.Node, t.NrThrottled, t.NrPeriods, float64(t.ThrottledUsec)/1e6))
		}
		fmt.Printf("| %s | %s |\n", strings.TrimSuffix(strings.TrimPrefix(id, "select-join-v4-topo-"), "-20261005"), strings.Join(parts, "; "))
	}
	fmt.Println()

	fmt.Println("## 6. P4 reuse ratios (NOT MATERIALIZED / MATERIALIZED) in the v3b domain runs, per run")
	fmt.Println()
	fmt.Println("| run | consumer | n=1 | n=2 | n=4 | n=8 | n=16 |")
	fmt.Println("|---|---|---|---|---|---|---|")
	var dom []string
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "select-join-v3b-domain-") {
			dom = append(dom, e.Name())
		}
	}
	sort.Strings(dom)
	for _, id := range dom {
		r, _, _ := load(id)
		m := medians(r)
		for _, c := range []string{"1%", "100%"} {
			var row []string
			for _, n := range []string{"1", "2", "4", "8", "16"} {
				num, ok1 := m["family=reuse;policy=notmat;refs="+n+";consumer="+c]
				den, ok2 := m["family=reuse;policy=mat;refs="+n+";consumer="+c]
				if ok1 && ok2 {
					row = append(row, fmt.Sprintf("%.3f", num/den))
				} else {
					row = append(row, "—")
				}
			}
			fmt.Printf("| %s | %s | %s |\n", strings.TrimSuffix(strings.TrimPrefix(id, "select-join-v3b-domain-"), "-20261005"), c, strings.Join(row, " | "))
		}
	}
}
