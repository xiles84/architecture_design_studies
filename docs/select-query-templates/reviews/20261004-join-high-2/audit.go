// Second independent HIGH review of the corrected join/CTE submission. It reads
// saved results and generated reports only; it starts no database.
//
//	podman run --rm --network none -v "<worktree>:/repo:ro" -w /repo \
//	  docker.io/library/golang:1.26-bookworm \
//	  go run ./docs/select-query-templates/reviews/20261004-join-high-2/audit.go
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const results = "studies/01-charity-tree/results"
const reports = "studies/01-charity-tree/reports"

type cell struct {
	Axis, Status, Plan string
	SamplesMS          []float64
}

type run struct {
	ID    string
	Cells []cell `json:"cells"`
	Error string `json:"error"`
}

func load(id string) (run, map[string]string) {
	var r run
	b, err := os.ReadFile(filepath.Join(results, id, "join.json"))
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(b, &r); err != nil {
		panic(err)
	}
	r.ID = id
	m := map[string]string{}
	f, err := os.Open(filepath.Join(results, id, "manifest.env"))
	if err == nil {
		s := bufio.NewScanner(f)
		for s.Scan() {
			if k, v, ok := strings.Cut(s.Text(), "="); ok {
				m[k] = v
			}
		}
		f.Close()
	}
	return r, m
}

func median(xs []float64) float64 {
	s := append([]float64(nil), xs...)
	sort.Float64s(s)
	return s[len(s)/2]
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

func main() {
	fmt.Println("# Audit output — second HIGH join/CTE review")
	fmt.Println()

	// 1. Rank-index use per engine on the frozen domain matrix.
	fmt.Println("## 1. Rank index use in saved EXPLAIN ANALYZE plans (domain2 matrix)")
	fmt.Println()
	fmt.Println("| run | valid cells | plans naming donation_rank_idx | plans with Seq Scan on donation + Storage Filter on rank |")
	fmt.Println("|---|---|---|---|")
	for _, e := range []string{"pg", "yb"} {
		for k := 1; k <= 3; k++ {
			id := fmt.Sprintf("select-join-domain2-%s%d-20261004", e, k)
			r, _ := load(id)
			valid, idx, seq := 0, 0, 0
			for _, c := range r.Cells {
				if c.Status != "valid" {
					continue
				}
				valid++
				if strings.Contains(c.Plan, "donation_rank_idx") {
					idx++
				}
				if strings.Contains(c.Plan, "Seq Scan on donation") && strings.Contains(c.Plan, "Storage Filter: (rank") {
					seq++
				}
			}
			fmt.Printf("| %s | %d | %d | %d |\n", id, valid, idx, seq)
		}
	}
	fmt.Println()

	// 2. Topology arms: manifest budgets versus the generated report's classes.
	fmt.Println("## 2. Topology subset arms (manifest) and join_where medians (ms)")
	fmt.Println()
	fmt.Println("| arm | topology | node budget | phase | t=20 run medians | t=2000 run medians | t=180000 run medians |")
	fmt.Println("|---|---|---|---|---|---|---|")
	for _, arm := range []string{"single-i", "equal-i", "cluster-i", "single-bo", "equal-bo", "cluster-bo"} {
		var topo, budget, phase string
		per := map[string][]string{}
		for k := 1; k <= 3; k++ {
			r, m := load(fmt.Sprintf("select-join-topo-yb-%s%d-20261004", arm, k))
			topo, budget, phase = m["topology"], m["node_budget"], m["phase"]
			md := medians(r)
			for _, t := range []string{"20", "2000", "180000"} {
				per[t] = append(per[t], fmt.Sprintf("%.1f", md["family=bag;stmt=join_where;threshold="+t]))
			}
		}
		fmt.Printf("| %s | %s | %s | %s | %s | %s | %s |\n", arm, topo, budget, phase,
			strings.Join(per["20"], ", "), strings.Join(per["2000"], ", "), strings.Join(per["180000"], ", "))
	}
	fmt.Println()
	rep, err := os.ReadFile(filepath.Join(reports, "select-join-topology-subset-20261004.md"))
	if err != nil {
		panic(err)
	}
	fmt.Println("Generated report classes containing runs of more than one node budget:")
	fmt.Println()
	for _, line := range strings.Split(string(rep), "\n") {
		if !strings.HasPrefix(line, "| `") || !strings.Contains(line, "select-join-topo") {
			continue
		}
		cols := strings.Split(line, "|")
		if len(cols) < 3 {
			continue
		}
		budgets := map[string]bool{}
		for _, id := range strings.Split(cols[2], ",") {
			id = strings.TrimSpace(id)
			if id == "" {
				continue
			}
			_, m := load(id)
			budgets[m["node_budget"]] = true
		}
		if len(budgets) > 1 {
			var bs []string
			for b := range budgets {
				bs = append(bs, b)
			}
			sort.Strings(bs)
			fmt.Printf("- class %s pools node budgets %s\n", strings.TrimSpace(cols[1]), strings.Join(bs, " + "))
		}
	}
	fmt.Println()

	// 3. Recorded limits that contradict the run's own settings.
	fmt.Println("## 3. Generated per-run limits contradicting the run")
	fmt.Println()
	n, total := 0, 0
	for _, arm := range []string{"single-bo", "equal-bo", "cluster-bo"} {
		for k := 1; k <= 3; k++ {
			id := fmt.Sprintf("select-join-topo-yb-%s%d-20261004", arm, k)
			b, err := os.ReadFile(filepath.Join(results, id, "join.md"))
			if err != nil {
				panic(err)
			}
			total++
			s := string(b)
			if strings.Contains(s, "phase: `blended`") && strings.Contains(s, "no concurrent writers") {
				n++
			}
		}
	}
	fmt.Printf("- blended runs whose report says `no concurrent writers ... blended phase is a recorded gap`: %d of %d\n", n, total)
	fmt.Println()

	// 4. Contract pairs on the frozen domain matrix.
	fmt.Println("## 4. Contract pairs per run (domain2)")
	fmt.Println()
	fmt.Println("| run | NOT MAT / MAT c=100% n=1 | n=2 | n=16 | c=1% n=1 | c=1% n=16 | joined-output MAT / default t=20 | t=2000 | t=180000 |")
	fmt.Println("|---|---|---|---|---|---|---|---|---|")
	for _, e := range []string{"pg", "yb"} {
		for k := 1; k <= 3; k++ {
			id := fmt.Sprintf("select-join-domain2-%s%d-20261004", e, k)
			r, _ := load(id)
			m := medians(r)
			reuse := func(n, c string) float64 {
				return m["family=reuse;policy=notmat;refs="+n+";consumer="+c] / m["family=reuse;policy=mat;refs="+n+";consumer="+c]
			}
			p3 := func(t string) float64 {
				return m["family=bag;stmt=cte_joined_outer_mat;threshold="+t] / m["family=bag;stmt=cte_joined_outer;threshold="+t]
			}
			fmt.Printf("| %s | %.3f | %.3f | %.3f | %.3f | %.3f | %.2f | %.2f | %.2f |\n", id,
				reuse("1", "100%"), reuse("2", "100%"), reuse("16", "100%"), reuse("1", "1%"), reuse("16", "1%"),
				p3("20"), p3("2000"), p3("180000"))
		}
	}
	fmt.Println()

	// 5. Placement check run.
	fmt.Println("## 5. Placement feature-check runs")
	fmt.Println()
	for _, id := range []string{"select-join-topo-yb-colocated-20261004", "select-join-topo-yb-noncolocated-20261004"} {
		r, m := load(id)
		fmt.Printf("- %s: placement=%s cells=%d error=%q\n", id, m["placement"], len(r.Cells), r.Error)
	}
}
