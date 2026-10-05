// SELECT synthesis evidence generator (HIGH, task-20261002T131729Z-select-synthesis).
//
// It reads saved, tagged results only and starts no database. Its default output
// is a generated, numbers-only report: every input run with its producing tag,
// commit and sha256(json), then per-run medians and ratios for every pair the
// synthesis analysis cites, with sign counts and plan-shape equality. It holds
// no interpretation; the verdicts live in the signed synthesis analysis. With
// -anchors it prints the registry anchors for each run set as JSON, so the
// evidence registry cites exactly the digests this report prints.
//
//	podman run --rm --network none -v "<worktree>:/repo:ro" -w /repo \
//	  docker.io/library/golang:1.26-bookworm \
//	  go run ./docs/select-query-templates/synthesis-20261005/evidence.go \
//	  > studies/01-charity-tree/reports/20261005T-select-synthesis-evidence.md
//
// Inputs: the randomized-order access/aggregate/existence runs
// (select-order-*-20261004), the v3b join/CTE domain runs and the v4
// placement-gated topology runs. The fixed-order primaries (select-access,
// select-agg, select-exist) are order-confounded and are not re-scored here.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	results     = "studies/01-charity-tree/results"
	reportPath  = "studies/01-charity-tree/reports/20261005T-select-synthesis-evidence.md"
	analysisRef = "studies/01-charity-tree/reports/analyses/20261005T-select-synthesis--claude-opus-5-5--2026-10-05.md"
)

type cell struct {
	Axis, Status, Plan string
	SamplesMS          []float64
}

type control struct {
	Statement  string
	Detected   bool
	Rows       int
	StatementL string `json:"statement"`
	DetectedL  bool   `json:"detected"`
}

type run struct {
	Cells    []cell    `json:"cells"`
	Controls []control `json:"controls"`
	Order    *struct {
		Seed     int64    `json:"seed"`
		Sequence []string `json:"sequence"`
	} `json:"order"`
	// Existence records one seeded order per fanout level instead of per run.
	Levels     json.RawMessage `json:"levels"`
	Provenance struct {
		RepoCommit string
		RunTag     string
	} `json:"provenance"`
	Snapshot *struct {
		Depths []struct {
			Equal bool `json:"equal"`
		} `json:"depths"`
	} `json:"snapshot_check"`
	Overlap []struct {
		OverlapCount int  `json:"overlap_count"`
		Disjoint     bool `json:"disjoint"`
	} `json:"overlap"`
	Placement []struct {
		Table      string `json:"table"`
		NumTablets int    `json:"num_tablets"`
	} `json:"placement"`
	EndpointOps []int `json:"endpoint_ops"`
}

type loaded struct {
	id, file, sum, topology string
	r                       run
	med                     map[string]float64
	pln                     map[string]string
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

var (
	nodeLine  = regexp.MustCompile(`\((cost|actual)=`)
	digits    = regexp.MustCompile(`[0-9]+`)
	indexNode = regexp.MustCompile(`(Index (Only )?Scan|Bitmap Index Scan)`)
)

// planShape reduces an EXPLAIN text to its node tree (depth, node type,
// relation and index names; numbers removed). Equal shapes mean the same
// physical plan.
func planShape(p string) string {
	var out []string
	for _, line := range strings.Split(p, "\n") {
		if !nodeLine.MatchString(line) {
			continue
		}
		indent := len(line) - len(strings.TrimLeft(line, " "))
		node := strings.TrimPrefix(strings.TrimSpace(line), "->")
		for _, cut := range []string{"  (", " (cost=", " (actual"} {
			if i := strings.Index(node, cut); i >= 0 {
				node = node[:i]
				break
			}
		}
		out = append(out, fmt.Sprintf("%d:%s", indent/2, digits.ReplaceAllString(strings.TrimSpace(node), "#")))
	}
	return strings.Join(out, "|")
}

func load(dir, file, topology string) loaded {
	b, err := os.ReadFile(filepath.Join(results, dir, file))
	if err != nil {
		panic(err)
	}
	var r run
	if err := json.Unmarshal(b, &r); err != nil {
		panic(err)
	}
	sum := sha256.Sum256(b)
	l := loaded{id: dir, file: file, sum: hex.EncodeToString(sum[:]), topology: topology, r: r,
		med: map[string]float64{}, pln: map[string]string{}}
	for _, c := range r.Cells {
		if c.Status == "valid" && len(c.SamplesMS) > 0 {
			l.med[c.Axis] = median(c.SamplesMS)
			l.pln[c.Axis] = planShape(c.Plan)
		}
	}
	return l
}

type set struct {
	key, prefix, file, topology string
}

var sets = []set{
	{"access-pg", "select-order-access-postgres-", "access.json", "postgresql-1node"},
	{"access-yb", "select-order-access-yugabyte-", "access.json", "yugabytedb-1node"},
	{"agg-pg", "select-order-aggregate-postgres-", "aggregate.json", "postgresql-1node"},
	{"agg-yb", "select-order-aggregate-yugabyte-", "aggregate.json", "yugabytedb-1node"},
	{"exist-pg", "select-order-existence-postgres-", "existence.json", "postgresql-1node"},
	{"exist-yb", "select-order-existence-yugabyte-", "existence.json", "yugabytedb-1node"},
	{"dom-pg-i", "select-join-v3b-domain-pg-i", "join.json", "postgresql-1node"},
	{"dom-pg-b", "select-join-v3b-domain-pg-b", "join.json", "postgresql-1node"},
	{"dom-yb-i", "select-join-v3b-domain-yb-i", "join.json", "yugabytedb-1node"},
	{"dom-yb-b", "select-join-v3b-domain-yb-b", "join.json", "yugabytedb-1node"},
	{"topo-single-i", "select-join-v4-topo-single-i", "join.json", "yugabytedb-1node"},
	{"topo-single-b", "select-join-v4-topo-single-b", "join.json", "yugabytedb-1node"},
	{"topo-equal-i", "select-join-v4-topo-equal-i", "join.json", "yugabytedb-1node-equal-total"},
	{"topo-equal-b", "select-join-v4-topo-equal-b", "join.json", "yugabytedb-1node-equal-total"},
	{"topo-tablet-i", "select-join-v4-topo-tablet-i", "join.json", "yugabytedb-3node-rf3"},
	{"topo-tablet-b", "select-join-v4-topo-tablet-b", "join.json", "yugabytedb-3node-rf3"},
	{"topo-distributed-i", "select-join-v4-topo-distributed-i", "join.json", "yugabytedb-3node-rf3"},
	{"topo-distributed-b", "select-join-v4-topo-distributed-b", "join.json", "yugabytedb-3node-rf3"},
	{"topo-colocated-i", "select-join-v4-topo-colocated-i", "join.json", "yugabytedb-3node-rf3"},
	{"topo-colocated-b", "select-join-v4-topo-colocated-b", "join.json", "yugabytedb-3node-rf3"},
}

func dirs(prefix string) []string {
	es, err := os.ReadDir(results)
	if err != nil {
		panic(err)
	}
	var out []string
	for _, e := range es {
		if strings.HasPrefix(e.Name(), prefix) {
			out = append(out, e.Name())
		}
	}
	sort.Strings(out)
	return out
}

func sign(xs []float64) string {
	gt, lt := 0, 0
	for _, x := range xs {
		if x > 1 {
			gt++
		} else if x < 1 {
			lt++
		}
	}
	switch {
	case gt == len(xs):
		return ">1"
	case lt == len(xs):
		return "<1"
	}
	return "mixed"
}

func list(xs []float64, prec int) string {
	var s []string
	for _, x := range xs {
		s = append(s, fmt.Sprintf("%.*f", prec, x))
	}
	return strings.Join(s, ", ")
}

// ratios returns per-run med[b]/med[a].
func ratios(runs []loaded, a, b string) []float64 {
	var out []float64
	for _, l := range runs {
		va, oka := l.med[a]
		vb, okb := l.med[b]
		if !oka || !okb || va == 0 {
			out = append(out, math.NaN())
			continue
		}
		out = append(out, vb/va)
	}
	return out
}

func meds(runs []loaded, a string) []float64 {
	var out []float64
	for _, l := range runs {
		out = append(out, l.med[a])
	}
	return out
}

func samePlan(runs []loaded, a, b string) string {
	same := 0
	for _, l := range runs {
		if l.pln[a] != "" && l.pln[a] == l.pln[b] {
			same++
		}
	}
	return fmt.Sprintf("%d/%d", same, len(runs))
}

func minmax(xs []float64) (float64, float64) {
	lo, hi := math.Inf(1), math.Inf(-1)
	for _, x := range xs {
		if !math.IsNaN(x) {
			lo, hi = math.Min(lo, x), math.Max(hi, x)
		}
	}
	return lo, hi
}

func separation(a, b []float64) string {
	alo, ahi := minmax(a)
	blo, bhi := minmax(b)
	switch {
	case ahi < blo:
		return "a<b every run"
	case alo > bhi:
		return "a>b every run"
	}
	return "overlap"
}

func digestOf(runs []loaded) string {
	var lines []string
	for _, l := range runs {
		lines = append(lines, l.id+" "+l.sum)
	}
	sort.Strings(lines)
	h := sha256.Sum256([]byte(strings.Join(lines, "\n") + "\n"))
	return hex.EncodeToString(h[:])
}

func controls(runs []loaded, name string) (fired, seen int) {
	for _, l := range runs {
		for _, c := range l.r.Controls {
			st, det := c.Statement, c.Detected
			if st == "" {
				st, det = c.StatementL, c.DetectedL
			}
			if st == name {
				seen++
				if det {
					fired++
				}
			}
		}
	}
	return
}

func orderNote(l loaded) string {
	if l.r.Order != nil {
		return fmt.Sprintf("seed %d", l.r.Order.Seed)
	}
	var lv []struct {
		Order *struct {
			Seed int64 `json:"seed"`
		} `json:"order"`
	}
	if json.Unmarshal(l.r.Levels, &lv) == nil {
		var s []string
		for _, x := range lv {
			if x.Order != nil {
				s = append(s, fmt.Sprint(x.Order.Seed))
			}
		}
		if len(s) > 0 {
			return "per-level seeds " + strings.Join(s, "/")
		}
	}
	return "not recorded"
}

type anchor struct {
	RunID        string `json:"run_id"`
	RunTag       string `json:"run_tag"`
	RepoCommit   string `json:"repo_commit"`
	InputsDigest string `json:"inputs_digest"`
	Report       string `json:"report"`
	Analysis     string `json:"analysis"`
	Environment  string `json:"environment"`
	Topology     string `json:"topology"`
	Role         string `json:"role"`
}

func main() {
	anchorsMode := flag.Bool("anchors", false, "print registry anchors per run set as JSON")
	flag.Parse()

	fam := map[string][]loaded{}
	for _, s := range sets {
		for _, d := range dirs(s.prefix) {
			fam[s.key] = append(fam[s.key], load(d, s.file, s.topology))
		}
	}

	if *anchorsMode {
		out := map[string][]anchor{}
		for _, s := range sets {
			for _, l := range fam[s.key] {
				out[s.key] = append(out[s.key], anchor{
					RunID: l.id, RunTag: l.r.Provenance.RunTag, RepoCommit: l.r.Provenance.RepoCommit,
					InputsDigest: l.sum, Report: reportPath, Analysis: analysisRef,
					Environment: "host-zenbook-ux5406sa", Topology: l.topology, Role: "corroborating",
				})
			}
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "  ")
		if err := enc.Encode(out); err != nil {
			panic(err)
		}
		return
	}

	p := func(format string, args ...any) { fmt.Printf(format+"\n", args...) }
	p("# Study 01 SELECT synthesis — evidence tables (generated measurements)")
	p("")
	p("_Generated by `docs/select-query-templates/synthesis-20261005/evidence.go` from saved, tagged results. Measurements only: per-run medians of the five samples per cell, per-run ratios, sign counts and plan-shape equality. No interpretation; verdicts are in the signed synthesis analysis._")
	p("")
	p("**TL;DR:** 60 whole runs in 20 run sets; every cell valid; every applicable control fired. Ratios are reported per run and never pooled across engines, phases, topologies or levels.")
	p("")
	p("A ratio `B / A` below 1 means B was faster. \"Sign\" is `>1` or `<1` only when all runs agree, otherwise `mixed`. \"Separation\" compares run ranges: `a<b every run` means every run of a was below every run of b.")
	p("")

	p("## Inputs")
	p("")
	p("| run | set | topology | run tag | repo commit | order | sha256(json) |")
	p("|---|---|---|---|---|---|---|")
	var all []loaded
	for _, s := range sets {
		for _, l := range fam[s.key] {
			all = append(all, l)
			p("| %s | %s | %s | `%s` | `%s` | %s | `%s` |", l.id, s.key, l.topology, l.r.Provenance.RunTag, l.r.Provenance.RepoCommit, orderNote(l), l.sum)
		}
	}
	p("")
	p("Inputs digest (sha256 of the sorted `run sha256(json)` lines, LF-terminated): `%s` over %d runs.", digestOf(all), len(all))
	p("")

	p("## Validity, controls and physical access")
	p("")
	p("| set | runs | valid cells | cells whose plan contains an index node |")
	p("|---|---|---|---|")
	for _, s := range sets {
		valid, cells, idx := 0, 0, 0
		for _, l := range fam[s.key] {
			for _, c := range l.r.Cells {
				cells++
				if c.Status == "valid" {
					valid++
				}
				if indexNode.MatchString(c.Plan) {
					idx++
				}
			}
		}
		p("| %s | %d | %d/%d | %d/%d |", s.key, len(fam[s.key]), valid, cells, idx, cells)
	}
	p("")
	ctl := func(keys []string, names ...string) string {
		var rs []loaded
		for _, k := range keys {
			rs = append(rs, fam[k]...)
		}
		var parts []string
		for _, n := range names {
			f, s := controls(rs, n)
			parts = append(parts, fmt.Sprintf("`%s` fired %d/%d", n, f, s))
		}
		return strings.Join(parts, "; ")
	}
	p("- access: %s (x01/x02/x03 cannot fire on this fixture: control absent)", ctl([]string{"access-pg", "access-yb"}, "x05_union_all_overlap"))
	p("- aggregate: %s", ctl([]string{"agg-pg", "agg-yb"}, "x05_aggregate_double_count", "x06_topn_rank_ties"))
	p("- existence: %s", ctl([]string{"exist-pg", "exist-yb"}, "x02_not_in_nullable", "x07_membership_multiplier"))
	p("- join/CTE domain: %s", ctl([]string{"dom-pg-i", "dom-pg-b", "dom-yb-i", "dom-yb-b"}, "x01_optional_where", "x02_join_limit_before"))
	var nullRows []string
	for _, k := range []string{"exist-pg", "exist-yb"} {
		for _, l := range fam[k] {
			for _, c := range l.r.Controls {
				if c.Statement == "x02_not_in_nullable" {
					nullRows = append(nullRows, fmt.Sprint(c.Rows))
				}
			}
		}
	}
	p("- `x02_not_in_nullable` rows returned at every fanout level and run: %s", strings.Join(nullRows, ", "))
	snapOK, snapN, disj, disjN := 0, 0, 0, 0
	for _, k := range []string{"access-pg", "access-yb"} {
		for _, l := range fam[k] {
			if l.r.Snapshot != nil {
				for _, d := range l.r.Snapshot.Depths {
					snapN++
					if d.Equal {
						snapOK++
					}
				}
			}
			for _, o := range l.r.Overlap {
				disjN++
				if o.Disjoint && o.OverlapCount == 0 {
					disj++
				}
			}
		}
	}
	p("- fixed-snapshot page equality, OFFSET page = keyset page: %d/%d depth checks; UNION arms disjoint: %d/%d shapes", snapOK, snapN, disj, disjN)
	p("")

	p("## Position controls (one cell repeated first and last in a run; last / first)")
	p("")
	for _, s := range sets {
		var rs []string
		for _, l := range fam[s.key] {
			for ax, v := range l.med {
				if strings.HasSuffix(ax, ";position=last") {
					first := strings.TrimSuffix(ax, ";position=last")
					if f, ok := l.med[first]; ok && f > 0 {
						rs = append(rs, fmt.Sprintf("`%s` %.3f", first, v/f))
					}
				}
			}
		}
		sort.Strings(rs)
		if len(rs) > 0 {
			p("- %s: %s", s.key, strings.Join(rs, "; "))
		}
	}
	p("")

	pairTable := func(title string, keys []string, rows [][2]string) {
		p("### %s", title)
		p("")
		p("| set | B / A | per-run ratios | sign | plan shape equal |")
		p("|---|---|---|---|---|")
		for _, k := range keys {
			for _, r := range rows {
				rs := ratios(fam[k], r[0], r[1])
				p("| %s | `%s` / `%s` | %s | %s | %s |", k, r[1], r[0], list(rs, 3), sign(rs), samePlan(fam[k], r[0], r[1]))
			}
		}
		p("")
	}

	p("## Inner-join spellings and CTE policy (domain, P1/P2/P3/P12)")
	p("")
	domain := []string{"dom-pg-i", "dom-pg-b", "dom-yb-i", "dom-yb-b"}
	thresholds := []string{"20", "200", "2000", "20000", "100000", "180000"}
	var spellRows [][2]string
	for _, t := range thresholds {
		base := "family=bag;stmt=join_where;threshold=" + t
		for _, s := range []string{"join_on", "filtered_derived", "cte_filtered", "cte_joined_outer", "comma_join"} {
			spellRows = append(spellRows, [2]string{base, "family=bag;stmt=" + s + ";threshold=" + t})
		}
	}
	pairTable("Default spellings against `join_where`", domain, spellRows)
	var matRows [][2]string
	for _, t := range thresholds {
		matRows = append(matRows,
			[2]string{"family=bag;stmt=cte_filtered;threshold=" + t, "family=bag;stmt=cte_filtered_mat;threshold=" + t},
			[2]string{"family=bag;stmt=cte_joined_outer;threshold=" + t, "family=bag;stmt=cte_joined_outer_mat;threshold=" + t})
	}
	pairTable("Forced `MATERIALIZED` against the default CTE", domain, matRows)

	p("## CTE reuse (domain, P4)")
	p("")
	var reuseRows [][2]string
	for _, c := range []string{"1%", "100%"} {
		for _, n := range []string{"1", "2", "4", "8", "16"} {
			reuseRows = append(reuseRows, [2]string{"family=reuse;policy=mat;refs=" + n + ";consumer=" + c, "family=reuse;policy=notmat;refs=" + n + ";consumer=" + c})
		}
	}
	pairTable("`NOT MATERIALIZED` against `MATERIALIZED`", domain, reuseRows)
	var defRows [][2]string
	for _, c := range []string{"1%", "100%"} {
		for _, n := range []string{"1", "2", "16"} {
			defRows = append(defRows, [2]string{"family=reuse;policy=mat;refs=" + n + ";consumer=" + c, "family=reuse;policy=default;refs=" + n + ";consumer=" + c})
		}
	}
	pairTable("Engine default against `MATERIALIZED`", domain, defRows)

	p("## Membership and absence (P6/P7)")
	p("")
	var exRows [][2]string
	for _, f := range []string{"1", "10", "100"} {
		m := func(fm, s string) string { return "family=" + fm + ";stmt=" + s + ";fanout=" + f }
		exRows = append(exRows,
			[2]string{m("membership", "exists"), m("membership", "join_distinct")},
			[2]string{m("membership", "exists"), m("membership", "in")},
			[2]string{m("membership", "exists"), m("membership", "in_notnull")},
			[2]string{m("absence", "not_exists"), m("absence", "left_is_null")},
			[2]string{m("absence", "not_exists"), m("absence", "not_in_notnull")})
	}
	pairTable("Spelling pairs by child fanout (10,000 child rows; 10,000 / 1,000 / 100 parents)", []string{"exist-pg", "exist-yb"}, exRows)

	p("## Parent aggregates and per-parent top-1 (P8/P9)")
	p("")
	var agRows [][2]string
	for _, n := range []string{"1", "10", "100", "1000", "10000"} {
		m := func(fm, s string) string { return "family=" + fm + ";stmt=" + s + ";parent_max=" + n }
		agRows = append(agRows,
			[2]string{m("aggregate", "post_join"), m("aggregate", "correlated")},
			[2]string{m("aggregate", "post_join"), m("aggregate", "pre_agg")},
			[2]string{m("topn", "rownumber"), m("topn", "lateral")},
			[2]string{m("topn", "rownumber"), m("topn", "distinct_on")})
	}
	pairTable("Spelling pairs by number of parents requested", []string{"agg-pg", "agg-yb"}, agRows)
	p("Correlated scalar aggregate, per-run medians (ms):")
	p("")
	for _, k := range []string{"agg-pg", "agg-yb"} {
		for _, n := range []string{"1", "10000"} {
			ax := "family=aggregate;stmt=correlated;parent_max=" + n
			p("- %s `%s`: %s", k, ax, list(meds(fam[k], ax), 3))
		}
	}
	p("")

	p("## Predicate alternatives and pagination (P10/P11)")
	p("")
	var orRows [][2]string
	for _, s := range []string{"selective/indexable", "selective/expression", "selective/cast", "broad/indexable", "broad/expression", "broad/cast",
		"selective/expression;index=expression_index", "broad/expression;index=expression_index"} {
		orRows = append(orRows,
			[2]string{"shape=" + s + ";spelling=or", "shape=" + s + ";spelling=union_all"},
			[2]string{"shape=" + s + ";spelling=or", "shape=" + s + ";spelling=union"})
	}
	pairTable("`UNION ALL` and `UNION` against `OR`", []string{"access-pg", "access-yb"}, orRows)
	pairTable("Expression index present against absent (same SQL; the index arm runs as its own phase)", []string{"access-pg", "access-yb"}, [][2]string{
		{"shape=selective/expression;spelling=or", "shape=selective/expression;index=expression_index;spelling=or"},
		{"shape=broad/expression;spelling=or", "shape=broad/expression;index=expression_index;spelling=or"}})
	var pgRows [][2]string
	for _, d := range []string{"1", "10", "100", "1000", "10000"} {
		pgRows = append(pgRows, [2]string{"depth=" + d + ";spelling=offset", "depth=" + d + ";spelling=keyset"})
	}
	pairTable("Keyset against OFFSET by page depth (pages of 5 rows)", []string{"access-pg", "access-yb"}, pgRows)
	p("Per-run medians (ms) at the shallowest and deepest page:")
	p("")
	for _, k := range []string{"access-pg", "access-yb"} {
		for _, d := range []string{"1", "10000"} {
			for _, s := range []string{"offset", "keyset"} {
				ax := "depth=" + d + ";spelling=" + s
				p("- %s `%s`: %s", k, ax, list(meds(fam[k], ax), 3))
			}
		}
	}
	p("")

	p("## Placement on a 3-node YugabyteDB cluster (v4 topology runs)")
	p("")
	p("| set | run | tablets per relation | reader endpoint operations |")
	p("|---|---|---|---|")
	for _, k := range []string{"topo-tablet-i", "topo-distributed-i", "topo-colocated-i"} {
		for _, l := range fam[k] {
			var t []string
			for _, pl := range l.r.Placement {
				if !strings.HasSuffix(pl.Table, "_pkey") {
					t = append(t, fmt.Sprintf("%s=%d", pl.Table[strings.LastIndex(pl.Table, ".")+1:], pl.NumTablets))
				}
			}
			p("| %s | %s | %s | %v |", k, l.id, strings.Join(t, " "), l.r.EndpointOps)
		}
	}
	p("")
	p("| a vs b | phase | statement | a per-run medians (ms) | b per-run medians (ms) | ratio of across-run medians | separation |")
	p("|---|---|---|---|---|---|---|")
	for _, pr := range [][2]string{{"distributed", "tablet"}, {"colocated", "distributed"}, {"distributed", "single"}, {"distributed", "equal"}, {"tablet", "single"}, {"colocated", "single"}} {
		for _, ph := range []string{"i", "b"} {
			for _, t := range []string{"20", "2000", "180000"} {
				ax := "family=bag;stmt=join_where;threshold=" + t
				a, b := meds(fam["topo-"+pr[0]+"-"+ph], ax), meds(fam["topo-"+pr[1]+"-"+ph], ax)
				p("| %s vs %s | %s | `%s` | %s | %s | %.3f | %s |", pr[0], pr[1], ph, ax, list(a, 2), list(b, 2), median(a)/median(b), separation(a, b))
			}
		}
	}
	p("")
	p("Forced-`MATERIALIZED` penalty at 2,000 rows (`family=bag;stmt=cte_joined_outer_mat;threshold=2000` / `family=bag;stmt=cte_joined_outer;threshold=2000`), per run:")
	p("")
	for _, k := range []string{"topo-single-i", "topo-tablet-i", "topo-distributed-i", "topo-colocated-i", "topo-single-b", "topo-tablet-b", "topo-distributed-b", "topo-colocated-b"} {
		rs := ratios(fam[k], "family=bag;stmt=cte_joined_outer;threshold=2000", "family=bag;stmt=cte_joined_outer_mat;threshold=2000")
		p("- %s: %s", k, list(rs, 2))
	}
}
