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

type Cell struct {
	Statement, Axis, Status, Plan string
	MedianMS                      float64
	TimingMeasured                bool
	SamplesMS                     []float64
}
type Provenance struct{ RepoCommit, RunTag, Environment string }
type Run struct {
	RunID      string     `json:"run_id"`
	Engine     string     `json:"engine"`
	Cells      []Cell     `json:"cells"`
	Provenance Provenance `json:"provenance"`
	Controls   []struct {
		Statement string
		Detected  bool
		Error     string
	} `json:"controls"`
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
func find(r Run, stmt, axis string) Cell {
	for _, c := range r.Cells {
		if c.Statement == stmt && c.Axis == axis {
			return c
		}
	}
	panic("missing cell " + r.RunID + " " + stmt + " " + axis)
}
func ratio(r Run, a, b, axisA, axisB string) float64 {
	return find(r, a, axisA).MedianMS / find(r, b, axisB).MedianMS
}
func lf(b []byte) []byte { return bytes.ReplaceAll(b, []byte("\r\n"), []byte("\n")) }
func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}
	results := filepath.Join(root, "studies/01-charity-tree/results")
	sets := map[string][]string{
		"access":    {"select-access-pg-1-20261003", "select-access-pg-2-20261003", "select-access-pg-3-20261003", "select-access-pg-pos-20261003", "select-access-yb-1-20261003", "select-access-yb-2-20261003", "select-access-yb-3-20261003"},
		"aggregate": {"select-agg-pg-1b-20261003", "select-agg-pg-2-20261003", "select-agg-pg-3-20261003", "select-agg-yb-1-20261003", "select-agg-yb-2-20261003", "select-agg-yb-3-20261003"},
		"existence": {"select-exist-pg-1b-20261004", "select-exist-pg-2-20261004", "select-exist-pg-3-20261004", "select-exist-yb-1-20261004", "select-exist-yb-2-20261004", "select-exist-yb-3-20261004"},
	}
	var paths []string
	for family, ids := range sets {
		for _, id := range ids {
			paths = append(paths, filepath.Join(results, id, family+".json"))
		}
	}
	sort.Strings(paths)
	hash := sha256.New()
	orders := map[string]map[string]int{}
	fmt.Println("# SELECT independent HIGH evidence check\n\nGenerated facts only. No database was started; timings below are read from saved result cells.\n\n## Input manifest\n\nSHA-256 hashes normalize CRLF to LF for checkout portability. Corpus digest frames sorted repository-relative paths and normalized byte lengths followed by file bytes.\n\n| run | JSON SHA-256 | producing commit | run tag | cells | CPU ids | MemTotal kB |\n|---|---|---|---|---|---|---|")
	var rows []string
	total, samples := 0, 0
	for _, path := range paths {
		raw, err := os.ReadFile(path)
		must(err)
		data := lf(raw)
		rel, err := filepath.Rel(root, path)
		must(err)
		rel = filepath.ToSlash(rel)
		fmt.Fprintf(hash, "%s\x00%d\x00", rel, len(data))
		hash.Write(data)
		var r Run
		must(json.Unmarshal(data, &r))
		if r.RunID != filepath.Base(filepath.Dir(path)) || len(r.Cells) == 0 || r.Provenance.RepoCommit == "" || r.Provenance.RunTag == "" {
			panic("missing run provenance: " + path)
		}
		var sequence []string
		for _, c := range r.Cells {
			if c.Status != "valid" || !c.TimingMeasured || c.MedianMS <= 0 || len(c.SamplesMS) < 3 || c.Plan == "" {
				panic("nonvalidated cell " + r.RunID + " " + c.Axis)
			}
			total++
			samples += len(c.SamplesMS)
			if !strings.Contains(c.Axis, "position=last") {
				sequence = append(sequence, c.Statement+"@"+c.Axis)
			}
		}
		family := strings.TrimSuffix(filepath.Base(path), ".json")
		required := map[string]int{}
		switch family {
		case "access":
			required["x05_union_all_overlap"] = 6
		case "aggregate":
			required["x05_aggregate_double_count"] = 1
			required["x06_topn_rank_ties"] = 1
		case "existence":
			required["x02_not_in_nullable"] = 3
			required["x07_membership_multiplier"] = 3
		}
		seen := map[string]int{}
		for _, control := range r.Controls {
			if required[control.Statement] > 0 {
				if !control.Detected || control.Error != "" {
					panic("required control not detected: " + r.RunID + " " + control.Statement)
				}
				seen[control.Statement]++
			}
		}
		for name, count := range required {
			if seen[name] != count {
				panic(fmt.Sprintf("control count %s %s: got %d want %d", r.RunID, name, seen[name], count))
			}
		}
		key := family + "/" + r.Engine
		if orders[key] == nil {
			orders[key] = map[string]int{}
		}
		orders[key][strings.Join(sequence, "|")]++
		cpus, err := os.ReadFile(filepath.Join(filepath.Dir(path), "resources/guest-cpus.txt"))
		must(err)
		mem, err := os.ReadFile(filepath.Join(filepath.Dir(path), "resources/guest-memory.txt"))
		must(err)
		memTotal := strings.Fields(strings.Split(string(mem), "\n")[0])[1]
		fmt.Printf("| %s | %x | %s | %s | %d | %s | %s |\n", r.RunID, sha256.Sum256(data), r.Provenance.RepoCommit, r.Provenance.RunTag, len(r.Cells), strings.TrimSpace(string(cpus)), memTotal)
		if family == "access" {
			for _, d := range []int{1, 10, 100, 1000, 10000} {
				rows = append(rows, fmt.Sprintf("| %s | keyset/OFFSET depth=%d | %.6f |", r.RunID, d, ratio(r, "s17_page_keyset", "s16_page_offset", fmt.Sprintf("depth=%d;spelling=keyset", d), fmt.Sprintf("depth=%d;spelling=offset", d))))
			}
		}
		if family == "aggregate" {
			for _, level := range []int{1, 10, 100, 1000, 10000} {
				for _, c := range r.Cells {
					if c.Axis == fmt.Sprintf("family=topn;stmt=distinct_on;parent_max=%d", level) {
						var other Cell
						for _, x := range r.Cells {
							if x.Axis == fmt.Sprintf("family=topn;stmt=rownumber;parent_max=%d", level) {
								other = x
							}
						}
						if other.MedianMS <= 0 {
							panic("missing rownumber")
						}
						rows = append(rows, fmt.Sprintf("| %s | DISTINCT ON/ROW_NUMBER parents=%d | %.6f |", r.RunID, level, c.MedianMS/other.MedianMS))
					}
				}
			}
		}
	}
	fmt.Printf("\nCorpus inputs_digest: %x\n\n19 runs; %d valid timed cells; %d saved drained-query samples.\n", hash.Sum(nil), total, samples)
	fmt.Println("\n## Ratios read directly from cells\n\nRatio below 1 means numerator has lower saved median; this is not a confidence test.\n\n| run | pair/level | ratio |\n|---|---|---|")
	for _, row := range rows {
		fmt.Println(row)
	}
	fmt.Println("\n## Recorded execution sequences (position repeats excluded)\n\n| family/engine | distinct sequences | runs |\n|---|---|---|")
	var keys []string
	for k := range orders {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		n := 0
		for _, v := range orders[k] {
			n += v
		}
		fmt.Printf("| %s | %d | %d |\n", k, len(orders[k]), n)
	}
}
