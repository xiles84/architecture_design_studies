// Third independent HIGH review of the join/CTE submission (v3b runs). It reads
// saved results and generated reports only; it starts no database.
//
//	podman run --rm --network none -v "<worktree>:/repo:ro" -w /repo \
//	  docker.io/library/golang:1.26-bookworm \
//	  go run ./docs/select-query-templates/reviews/20261005-join-high-3/audit.go
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
const reports = "studies/01-charity-tree/reports"

type cell struct {
	Axis, Status, Plan string
	SamplesMS          []float64
	PlanIndexScan      bool `json:"plan_index_scan"`
	PlanSeqScan        bool `json:"plan_seq_scan"`
	StorageRows        int  `json:"storage_rows_scanned"`
}

type run struct {
	Cells    []cell `json:"cells"`
	Controls []struct {
		Statement string `json:"statement"`
		Detected  bool   `json:"detected"`
	} `json:"controls"`
	Error       string `json:"error"`
	EndpointOps []int  `json:"endpoint_ops"`
	Writer      *struct {
		Offered, Completed, Failed int
		EndpointOps                []int `json:"endpoint_ops"`
	} `json:"writer"`
	WriterSamples int `json:"writer_samples"`
	WriterOverlap int `json:"writer_overlap_samples"`
	Placement     []struct {
		Table       string `json:"table"`
		NumTablets  int    `json:"num_tablets"`
		IsColocated bool   `json:"is_colocated"`
	} `json:"placement"`
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

var leaderRE = regexp.MustCompile(`^[0-9a-f]{32}\s.*\s(yb-n[0-9]):[0-9]+\s`)

// leaders returns table -> leader nodes from the saved yb-admin listing.
func leaders(id string) map[string][]string {
	out := map[string][]string{}
	f, err := os.Open(filepath.Join(results, id, "resources", "placement.txt"))
	if err != nil {
		return out
	}
	defer f.Close()
	table := ""
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(line, "## ") {
			table = strings.TrimPrefix(line, "## ")
			continue
		}
		if m := leaderRE.FindStringSubmatch(line); m != nil {
			out[table] = append(out[table], m[1])
		}
	}
	return out
}

func main() {
	fmt.Println("# Audit output — third HIGH join/CTE review")
	fmt.Println()

	var ids []string
	entries, err := os.ReadDir(results)
	if err != nil {
		panic(err)
	}
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "select-join-v3b-") {
			ids = append(ids, e.Name())
		}
	}
	sort.Strings(ids)

	fmt.Println("## 1. Every v3b run: validity, samples, controls, build and digest")
	fmt.Println()
	fmt.Println("| run | arm settings | valid/cells | 5 samples | controls fired | report digest = sha256(join.json) |")
	fmt.Println("|---|---|---|---|---|---|")
	builds := map[string]int{}
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
		digestOK := strings.Contains(string(rep), sum)
		builds[m["client_image_id"]]++
		fmt.Printf("| %s | %s/%s/%s/%s/%s/%s | %d/%d | %d | %d/%d | %v |\n", id, m["engine"], m["topology"],
			m["node_budget"], m["phase"], m["placement"], m["index_arm"], valid, len(r.Cells), five, fired, len(r.Controls), digestOK)
	}
	fmt.Println()
	fmt.Printf("Client image ids across v3b runs: %d distinct\n\n", len(builds))

	fmt.Println("## 2. Index use in domain plans (cells with plan_index_scan / plan_seq_scan)")
	fmt.Println()
	fmt.Println("| run | index-scan cells | seq-scan cells | j01 t=20 storage rows |")
	fmt.Println("|---|---|---|---|")
	for _, id := range ids {
		if !strings.Contains(id, "-domain-") {
			continue
		}
		r, _, _ := load(id)
		idx, seq, j01 := 0, 0, -1
		for _, c := range r.Cells {
			if c.PlanIndexScan {
				idx++
			}
			if c.PlanSeqScan {
				seq++
			}
			if c.Axis == "family=bag;stmt=join_where;threshold=20" {
				j01 = c.StorageRows
			}
		}
		fmt.Printf("| %s | %d | %d | %d |\n", id, idx, seq, j01)
	}
	fmt.Println()

	fmt.Println("## 3. Contract pairs per domain run")
	fmt.Println()
	fmt.Println("| run | MAT/default joined t=20 | t=2000 | t=180000 | NOTMAT/MAT n=1 c=1% | n=1 c=100% | n=16 c=1% | n=16 c=100% |")
	fmt.Println("|---|---|---|---|---|---|---|---|")
	for _, id := range ids {
		if !strings.Contains(id, "-domain-") {
			continue
		}
		r, _, _ := load(id)
		m := medians(r)
		p3 := func(t string) float64 {
			return m["family=bag;stmt=cte_joined_outer_mat;threshold="+t] / m["family=bag;stmt=cte_joined_outer;threshold="+t]
		}
		p4 := func(n, c string) float64 {
			return m["family=reuse;policy=notmat;refs="+n+";consumer="+c] / m["family=reuse;policy=mat;refs="+n+";consumer="+c]
		}
		fmt.Printf("| %s | %.2f | %.2f | %.3f | %.3f | %.3f | %.3f | %.3f |\n", id, p3("20"), p3("2000"), p3("180000"),
			p4("1", "1%"), p4("1", "100%"), p4("16", "1%"), p4("16", "100%"))
	}
	fmt.Println()

	fmt.Println("## 4. Topology arms: join_where run medians (ms) and MAT/default at t=2000")
	fmt.Println()
	fmt.Println("| arm | t=20 | t=180000 | joined MAT/default t=2000 |")
	fmt.Println("|---|---|---|---|")
	for _, arm := range []string{"topo-single-i", "topo-equal-i", "topo-cluster-i", "coloc-i", "topo-single-b", "topo-equal-b", "topo-cluster-b", "coloc-b"} {
		var a, b, c []string
		for k := 1; k <= 3; k++ {
			r, _, _ := load(fmt.Sprintf("select-join-v3b-%s%d-20261005", arm, k))
			m := medians(r)
			a = append(a, fmt.Sprintf("%.2f", m["family=bag;stmt=join_where;threshold=20"]))
			b = append(b, fmt.Sprintf("%.1f", m["family=bag;stmt=join_where;threshold=180000"]))
			c = append(c, fmt.Sprintf("%.1f", m["family=bag;stmt=cte_joined_outer_mat;threshold=2000"]/m["family=bag;stmt=cte_joined_outer;threshold=2000"]))
		}
		fmt.Printf("| %s | %s | %s | %s |\n", arm, strings.Join(a, ", "), strings.Join(b, ", "), strings.Join(c, ", "))
	}
	fmt.Println()

	fmt.Println("## 5. Cluster placement: tablets and leader node per relation")
	fmt.Println()
	fmt.Println("| run | placement | relation tablets (yb_table_properties) | leader nodes (yb-admin) |")
	fmt.Println("|---|---|---|---|")
	for _, id := range ids {
		if !strings.Contains(id, "cluster") && !strings.Contains(id, "coloc") {
			continue
		}
		r, m, _ := load(id)
		var tabs []string
		for _, p := range r.Placement {
			tabs = append(tabs, fmt.Sprintf("%s=%d%s", strings.TrimPrefix(p.Table, "selectjoin_domain."), p.NumTablets,
				map[bool]string{true: "c", false: ""}[p.IsColocated]))
		}
		ls := leaders(id)
		nodes := map[string]bool{}
		var lt []string
		var keys []string
		for t := range ls {
			keys = append(keys, t)
		}
		sort.Strings(keys)
		for _, t := range keys {
			lt = append(lt, t+":"+strings.Join(ls[t], "+"))
			for _, n := range ls[t] {
				nodes[n] = true
			}
		}
		fmt.Printf("| %s | %s | %s | %s (distinct leader nodes: %d) |\n", id, m["placement"], strings.Join(tabs, " "), strings.Join(lt, " "), len(nodes))
	}
	fmt.Println()

	fmt.Println("## 6. Blended writer and endpoints")
	fmt.Println()
	fmt.Println("| run | reader endpoint ops | writer offered/completed/failed | writer endpoint ops | overlap/samples |")
	fmt.Println("|---|---|---|---|---|")
	for _, id := range ids {
		r, m, _ := load(id)
		if m["phase"] != "blended" || r.Writer == nil {
			continue
		}
		fmt.Printf("| %s | %v | %d/%d/%d | %v | %d/%d |\n", id, r.EndpointOps, r.Writer.Offered, r.Writer.Completed,
			r.Writer.Failed, r.Writer.EndpointOps, r.WriterOverlap, r.WriterSamples)
	}
	fmt.Println()

	fmt.Println("## 7. Generated per-run reports: stale limit text")
	fmt.Println()
	stale := 0
	for _, id := range ids {
		b, _ := os.ReadFile(filepath.Join(results, id, "join.md"))
		s := string(b)
		_, m, _ := load(id)
		if m["phase"] == "blended" && strings.Contains(s, "no concurrent writers") {
			stale++
		}
		if strings.Contains(s, "not parsed; raw EXPLAIN retained (coverage gap)") {
			stale++
		}
	}
	fmt.Printf("- stale limit phrases found: %d\n", stale)
	rep, _ := os.ReadFile(filepath.Join(reports, "select-join-v3b-topology-20261005.md"))
	n := strings.Count(string(rep), "\n| `")
	fmt.Printf("- topology report table rows starting with a class id: %d\n", n)
}
