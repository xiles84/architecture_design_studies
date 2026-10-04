// Recheck immutable correction inputs without starting a database.
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type cell struct {
	Axis, Status, Plan string
	Trials int
	TimingMeasured bool
	MedianMS float64
	SamplesMS []float64
}
type order struct { Seed int64; Rounds int; Sequence []string }
type control struct { Statement string; Detected bool; Fanout int; Shape string }
type result struct {
	RunID string `json:"run_id"`
	Engine string
	Provenance struct { RunTag, RepoCommit, Environment string }
	Cells []cell
	Controls []control
	Order order
	IndexOrder order `json:"index_order"`
	Levels json.RawMessage
}

func main() {
	root := os.Args[1]
	files, err := filepath.Glob(filepath.Join(root, "studies/01-charity-tree/results/select-order-*/*.json")); must(err)
	sort.Strings(files)
	counts := map[string]int{}
	distinct := map[string]map[string]bool{}
	corpora := map[string][]string{}
	totalCells, totalSamples := 0, 0
	fmt.Println("# Independent correction-input check — GPT-6, 2026-10-04\n\nTL;DR: immutable inputs checked without a database run.\n\n| run | cells | samples | producing tag | normalized SHA-256 |\n|---|---:|---:|---|---|")
	for _, path := range files {
		family := strings.TrimSuffix(filepath.Base(path), ".json")
		if family != "access" && family != "aggregate" && family != "existence" { continue }
		raw, err := os.ReadFile(path); must(err)
		data := bytes.ReplaceAll(raw, []byte("\r\n"), []byte("\n"))
		var r result; must(json.Unmarshal(data, &r))
		check(r.Provenance.RunTag == "run/study-01/"+r.RunID, "missing producing tag")
		check(r.Provenance.RepoCommit == "ca2a38fbf9e53da00344578258138db9e01d3545", "wrong producing source")
		check(r.Provenance.Environment == "xiles-podman-select-order-correction-20261004", "wrong environment")
		byAxis := map[string]cell{}
		samples := 0
		for _, c := range r.Cells {
			check(c.Status == "valid" && c.TimingMeasured && c.Trials == 5 && len(c.SamplesMS) == 5 && c.Plan != "", "cell lacks valid timing/plan")
			check(!byAxis[c.Axis].TimingMeasured, "duplicate axis")
			byAxis[c.Axis] = c; samples += len(c.SamplesMS)
		}
		orders := []order{r.Order}
		if family == "existence" {
			var levels []struct { Fanout int; Order order }; must(json.Unmarshal(r.Levels, &levels))
			orders = nil; for _, l := range levels { orders = append(orders, l.Order) }
		} else if family == "access" { orders = append(orders, r.IndexOrder) }
		var sequence []string
		for _, o := range orders {
			check(o.Rounds == 5 && len(o.Sequence)>0, "order missing rounds/sequence")
			for _, axis := range o.Sequence {
				check(byAxis[axis].TimingMeasured, "order references absent cell")
				if strings.HasSuffix(axis, ";position=last") {
					check(strings.TrimSuffix(axis, ";position=last") == o.Sequence[0], "position control is not actual first cell")
					fmt.Fprintf(os.Stderr, "%s position last/first %.9f\n", r.RunID, byAxis[axis].MedianMS/byAxis[o.Sequence[0]].MedianMS)
				} else { sequence = append(sequence, axis) }
			}
		}
		key := family+"/"+r.Engine
		counts[key]++
		if distinct[key] == nil { distinct[key] = map[string]bool{} }; distinct[key][strings.Join(sequence, "|")] = true
		// Only required controls are counted: auxiliary controls may legitimately
		// not fire on a fixture lacking the case they exercise.
		required := map[string]int{}
		for _, c := range r.Controls {
			needed := (family=="access" && c.Statement=="x05_union_all_overlap") || (family=="aggregate" && (c.Statement=="x05_aggregate_double_count" || c.Statement=="x06_topn_rank_ties")) || (family=="existence" && (c.Statement=="x07_membership_multiplier" || c.Statement=="x02_not_in_nullable"))
			if needed { check(c.Detected, "required negative control did not fire"); required[c.Statement]++ }
		}
		fmt.Fprintf(os.Stderr, "%s controls %v\n", r.RunID, required)
		if family=="access" { check(required["x05_union_all_overlap"]==6, "missing overlap control") }
		if family=="aggregate" { check(required["x05_aggregate_double_count"]==1 && required["x06_topn_rank_ties"]==1, "missing aggregate control") }
		if family=="existence" { check(required["x07_membership_multiplier"]==3 && required["x02_not_in_nullable"]==3, "missing fanout control") }
		if family == "aggregate" {
			for _, n := range []int{1,10,100,1000,10000} {
				a := fmt.Sprintf("family=topn;stmt=distinct_on;parent_max=%d", n)
				b := fmt.Sprintf("family=topn;stmt=rownumber;parent_max=%d", n)
				fmt.Fprintf(os.Stderr, "%s level %d distinct_on/rownumber %.9f\n", r.RunID,n,byAxis[a].MedianMS/byAxis[b].MedianMS)
			}
		}
		corpora[family] = append(corpora[family], path)
		fmt.Printf("| %s | %d | %d | `%s` | `%x` |\n", r.RunID, len(r.Cells), samples, r.Provenance.RunTag, sha256.Sum256(data))
		totalCells+=len(r.Cells); totalSamples+=samples
	}
	check(len(counts)==6, "wrong family/engine coverage")
	for key, count := range counts { check(count==3 && len(distinct[key])==3, "need three distinct whole-run orders for "+key) }
	fmt.Printf("\n18 runs; %d valid cells; %d verified samples; 3 distinct recorded orders for each family/engine.\n", totalCells, totalSamples)
	for _, family := range []string{"access","aggregate","existence"} {
		h := sha256.New()
		for _, p := range corpora[family] {
			raw,err:=os.ReadFile(p);must(err); data:=bytes.ReplaceAll(raw,[]byte("\r\n"),[]byte("\n"))
			rel:=filepath.ToSlash(filepath.Join(filepath.Base(filepath.Dir(p)),filepath.Base(p)))
			fmt.Fprintf(h,"%s\n%d\n",rel,len(data));h.Write(data)
		}
		fmt.Printf("\n%s corpus: `%x`.\n",family,h.Sum(nil))
	}
}
func must(err error) { if err!=nil { panic(err) } }
func check(ok bool, why string) { if !ok { panic(why) } }
