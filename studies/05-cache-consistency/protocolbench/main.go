// The entry adapter wires the port-based experiment to the PostgreSQL driver.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"adsplatform/adapters/cgroup"
	"adsplatform/adapters/pgxdb"
	"adsplatform/core/provenance"
)

type result struct {
	CellID                        string                 `json:"cell_id,omitempty"`
	RunID                         string                 `json:"run_id"`
	Environment                   string                 `json:"environment"`
	CreatedAt                     string                 `json:"created_at"`
	Repo                          provenance.RepoVersion `json:"repo"`
	Image                         string                 `json:"image"`
	Arm, Protocol, Strategy, Kind string
	Keys, Readers, Writers, Pool  int
	DurationS                     float64
	Seed                          int64
	Engine                        string
	Endpoint                      string
	Isolation                     string
	Controls                      []control
	GatePassed                    bool   `json:"gate_passed"`
	Error                         string `json:"error,omitempty"`
	FailureCounts                 struct{ WrongReads, ImpossibleReads, Errors int64 }
	Trials                        []trial `json:"trials,omitempty"`
	Plans                         map[string]string
	ClientCPU                     cgroup.CPUStat
	CgroupAvailable               bool
	RedisBefore, RedisAfter       string
	ObjectRows, SourceBytes       int64
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func run() error {
	cmd := flag.String("cmd", "run", "run, order, report, summary, or verify")
	dsn := flag.String("dsn", "postgres://bench:bench@cp-pg:5432/bench?sslmode=disable", "source")
	addr := flag.String("redis", "cp-redis:6379", "shared cache")
	arm := flag.String("arm", "A", "A publication or B deletion")
	protocol := flag.String("protocol", "sync", "version, sync, outbox")
	strategy := flag.String("strategy", "cas", "cas or lock")
	kind := flag.String("kind", "hard", "soft, hard, tombstone")
	keys := flag.Int("keys", 128, "logical keys")
	readers := flag.Int("readers", 8, "reader clients")
	writers := flag.Int("writers", 2, "writer clients")
	trials := flag.Int("trials", 3, "independent trials per phase")
	duration := flag.Duration("duration", 3*time.Second, "trial admission duration")
	seed := flag.Int64("seed", 42, "recorded seed")
	runID := flag.String("run-id", "", "whole-run identity")
	commit := flag.String("commit", "", "producing clean commit")
	tag := flag.String("tag", "", "immutable producing run tag")
	env := flag.String("environment", "xiles-podman-4cpu-2g-config-protocol-20261003", "environment")
	image := flag.String("image", "", "exact client image id")
	out := flag.String("out", "/results/result.json", "exclusive result destination")
	inputs := flag.String("inputs", "/results", "results directory")
	analyses := flag.String("analyses", "/analyses", "signed analyses directory")
	flag.Parse()
	if *cmd == "order" {
		cells := []string{"version-cas", "version-lock", "sync-cas", "sync-lock", "outbox-cas", "outbox-lock"}
		r := rand.New(rand.NewSource(*seed))
		r.Shuffle(len(cells), func(i, j int) { cells[i], cells[j] = cells[j], cells[i] })
		for _, c := range cells {
			fmt.Println(c)
		}
		fmt.Println(cells[0] + "-position-last")
		return nil
	}
	if *cmd == "report" {
		return writeReport(*inputs, *analyses, *out)
	}
	if *cmd == "summary" {
		return summarize(*inputs, *runID, *out)
	}
	if *cmd == "verify" {
		return verifyResults(*inputs)
	}
	if *cmd != "run" {
		return fmt.Errorf("unknown command %s", *cmd)
	}
	if *keys < 1 || *readers < 1 || *writers < 1 || *duration <= 0 || *trials < 1 {
		return fmt.Errorf("positive keys, clients, duration and trials required")
	}
	if !oneOf(*protocol, "version", "sync", "outbox") || !oneOf(*strategy, "cas", "lock") || !oneOf(*kind, "soft", "hard", "tombstone") || !oneOf(*arm, "A", "B") {
		return fmt.Errorf("invalid cell configuration")
	}
	if *arm == "B" && (*protocol != "sync" || *strategy != "cas") {
		return fmt.Errorf("Arm B holds synchronous publication and CAS constant")
	}
	if *commit == "" || *tag == "" || *runID == "" {
		return fmt.Errorf("clean producing commit, run id and annotated tag required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	pool := *readers + *writers + 1 + 3
	db, e := pgxdb.Open(ctx, pgxdb.Config{DSN: *dsn, MaxConns: pool, StatementTimeoutMS: 5000, ApplicationName: "cache-protocol-v5"})
	if e != nil {
		return e
	}
	defer db.Close()
	cache := newRedis(*addr)
	defer cache.close()
	b, e := newBench(db, cache, *protocol, *strategy, *kind, *keys)
	if e != nil {
		return e
	}
	defer b.stopRelay()
	res := result{RunID: *runID, Environment: *env, CreatedAt: time.Now().UTC().Format(time.RFC3339), Repo: provenance.RepoVersion{Commit: *commit, Tag: *tag}, Image: *image, Arm: *arm, Protocol: *protocol, Strategy: *strategy, Kind: *kind, Keys: *keys, Readers: *readers, Writers: *writers, Pool: pool, DurationS: duration.Seconds(), Seed: *seed, Endpoint: "cp-pg:5432 (all source requests; single node)", Isolation: "READ COMMITTED", Plans: map[string]string{}}
	before, ok := cgroup.Read()
	res.CgroupAvailable = ok
	if e = db.QueryRow(ctx, "SELECT version()").Scan(&res.Engine); e != nil {
		return e
	}
	info, e := cache.do(ctx, "INFO")
	if e != nil {
		return e
	}
	res.RedisBefore = fmt.Sprint(info)
	res.Controls, e = b.controls(ctx)
	if e == nil && *arm == "B" {
		var cs []control
		cs, e = b.deletionControls(ctx)
		res.Controls = append(res.Controls, cs...)
	}
	if e == nil {
		for _, n := range []string{"snapshot", "version", "pending"} {
			args := []any{int64(1)}
			if n == "pending" {
				args = nil
			}
			rows, qe := db.Query(ctx, "EXPLAIN (ANALYZE, BUFFERS, FORMAT TEXT) "+b.q(n), args...)
			if qe != nil {
				e = qe
				break
			}
			var lines []string
			for rows.Next() {
				var s string
				if qe = rows.Scan(&s); qe != nil {
					break
				}
				lines = append(lines, s)
			}
			if qe == nil {
				qe = rows.Err()
			}
			rows.Close()
			if qe != nil {
				e = qe
				break
			}
			res.Plans[n] = strings.Join(lines, "\n")
		}
	}
	phases := []string{"read", "write", "mixed"}
	if *arm == "B" {
		phases = []string{"delete-recreate"}
	}
	if e == nil {
	outer:
		for n := 1; n <= *trials; n++ {
			for _, phase := range phases {
				var t trial
				t, e = b.measure(ctx, phase, n, *readers, *writers, *duration, *seed+int64(n)*101)
				if e != nil {
					res.FailureCounts.WrongReads += t.WrongReads
					res.FailureCounts.ImpossibleReads += t.ImpossibleReads
					res.FailureCounts.Errors += t.Errors
					break outer
				}
				res.Trials = append(res.Trials, t)
			}
		}
	}
	b.stopRelay()
	after, available := cgroup.Read()
	res.ClientCPU = after.Sub(before)
	res.CgroupAvailable = ok && available
	if v, ie := cache.do(ctx, "INFO"); ie == nil {
		res.RedisAfter = fmt.Sprint(v)
	} else if e == nil {
		e = ie
	}
	if e == nil {
		e = db.QueryRow(ctx, "SELECT count(*) FROM protocol_object").Scan(&res.ObjectRows)
	}
	if e == nil {
		e = db.QueryRow(ctx, "SELECT sum(pg_total_relation_size(c::regclass))::bigint FROM unnest(ARRAY['protocol_generation','protocol_object','protocol_outbox']) c").Scan(&res.SourceBytes)
	}
	res.GatePassed = e == nil
	if e != nil {
		res.Error = e.Error()
		res.Trials = nil
		if res.FailureCounts.Errors == 0 {
			res.FailureCounts.Errors = 1
		}
	}
	if we := writeExclusiveJSON(*out, res); we != nil {
		return we
	}
	fmt.Printf("%s %s/%s/%s gate=%t trials=%d\n", *arm, *protocol, *strategy, *kind, res.GatePassed, len(res.Trials))
	return e
}
func oneOf(s string, xs ...string) bool {
	for _, x := range xs {
		if s == x {
			return true
		}
	}
	return false
}
func writeExclusiveJSON(p string, v any) error {
	bytes, e := json.MarshalIndent(v, "", "  ")
	if e != nil {
		return e
	}
	return writeExclusive(p, append(bytes, '\n'))
}
func writeExclusive(p string, bytes []byte) error {
	if e := os.MkdirAll(filepath.Dir(p), 0755); e != nil {
		return e
	}
	f, e := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if e != nil {
		return e
	}
	_, e = f.Write(bytes)
	ce := f.Close()
	if e != nil {
		return e
	}
	return ce
}
func readResults(root string) ([]result, error) {
	var all []result
	e := filepath.WalkDir(root, func(p string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if d.IsDir() || !strings.HasSuffix(p, ".json") || filepath.Base(filepath.Dir(p)) != "cells" {
			return nil
		}
		bytes, e := os.ReadFile(p)
		if e != nil {
			return e
		}
		var r result
		if e = json.Unmarshal(bytes, &r); e != nil {
			return e
		}
		if r.RunID != "" && r.Protocol != "" {
			r.CellID = strings.TrimSuffix(filepath.Base(p), ".json")
			all = append(all, r)
		}
		return nil
	})
	sort.Slice(all, func(i, j int) bool { return all[i].CreatedAt < all[j].CreatedAt })
	if e == nil && len(all) == 0 {
		e = fmt.Errorf("no protocol results")
	}
	return all, e
}
func verifyResults(root string) error {
	rs, e := readResults(root)
	if e != nil {
		return e
	}
	for _, r := range rs {
		if !r.GatePassed || r.Error != "" {
			return fmt.Errorf("failed cell %s %s %s: %s", r.RunID, r.Protocol, r.Strategy, r.Error)
		}
		for _, c := range r.Controls {
			if !c.Passed {
				return fmt.Errorf("failed negative/control gate %s", c.Name)
			}
		}
		want := 9
		if r.Arm == "B" {
			want = 3
		}
		if len(r.Trials) != want {
			return fmt.Errorf("cell has %d trials, want %d", len(r.Trials), want)
		}
		for _, t := range r.Trials {
			if t.Errors+t.WrongReads+t.ImpossibleReads != 0 {
				return fmt.Errorf("wrong answer or operation error in admitted timings")
			}
			if t.Phase == "read" {
				want := 0.
				if r.Protocol == "version" {
					want = 1
				}
				if t.SourceQueriesPerRead != want {
					return fmt.Errorf("warmed read source-query control: got %g want %g", t.SourceQueriesPerRead, want)
				}
			}
		}
	}
	fmt.Printf("verified %d cells: all control gates, phase/trial counts, clean read-only query contrasts and zero wrong/impossible reads\n", len(rs))
	return nil
}
func writeReport(root, analyses, out string) error {
	rs, e := readResults(root)
	if e != nil {
		return e
	}
	digest, n, e := provenance.Digest(os.DirFS(root))
	if e != nil {
		return e
	}
	var text strings.Builder
	fmt.Fprintf(&text, "# Cache publication and deletion measurements\n\n## TL;DR\n\n")
	passed, controls := 0, 0
	for _, r := range rs {
		if r.GatePassed {
			passed++
		}
		for _, c := range r.Controls {
			if c.Passed {
				controls++
			}
		}
	}
	fmt.Fprintf(&text, "Fixed selection rule: count every result and control; publish timing only for a cell whose complete correctness gate passed. %d/%d cells passed; %d controls passed. Inputs digest: `%s` (%d JSON files). This report contains measured facts; interpretation belongs in a signed analysis.\n\n", passed, len(rs), controls, digest, n)
	fmt.Fprintf(&text, "Environment: `%s`. Producing commit: `%s`; run tag: `%s`; client image: `%s`.\n\n", rs[0].Environment, rs[0].Repo.Commit, rs[0].Repo.Tag, rs[0].Image)
	text.WriteString("Single-node PostgreSQL and Redis share this laptop's Podman VM. Service budget: PostgreSQL 1 CPU/512 MiB plus Redis 1 CPU/256 MiB; client 1 CPU/512 MiB; planned reserve 1 CPU/768 MiB against machine configuration 4 CPU/2 GiB. The guest's actual resources and service cgroup counters are retained beside each cell. All source requests use one endpoint. Three-node/equal-total cluster and multi-host placement comparisons are resource/coverage gaps. Colocation, rollup, rolldown and embedding have no controlled contrast in this keyed projection. These are short, closed-loop measurements on 128 keys, with per-key client writer ordering, not saturation, donor-portal throughput or uncoordinated-writer evidence.\n\n")
	text.WriteString("## Every admitted trial\n\n| Arm / protocol / strategy / deletion | Phase / trial | Read/s | Write/s | Read p50/p99 ms | Write p50/p99 ms | DB p50 ms | Ack p50/p99 ms | Client lock p99 ms | Source queries/read | Relay queries / delivered |\n|---|---|---:|---:|---|---|---:|---|---:|---:|---|\n")
	for _, r := range rs {
		for _, t := range r.Trials {
			fmt.Fprintf(&text, "| %s / %s / %s / %s | %s / %d | %.1f | %.1f | %.3f / %.3f | %.3f / %.3f | %.3f | %.3f / %.3f | %.3f | %.6f | %d / %d |\n", r.Arm, r.Protocol, r.Strategy, r.Kind, t.Phase, t.Number, t.ReadOps, t.WriteOps, t.ReadP50MS, t.ReadP99MS, t.WriteP50MS, t.WriteP99MS, t.DBP50MS, t.AckP50MS, t.AckP99MS, t.LockP99MS, t.SourceQueriesPerRead, t.RelayQueries, t.DeliveredEvents)
		}
	}
	text.WriteString("\nA deletion/recreation pair counts as two acknowledged writes; write latencies refer to each individual write. Trial duration includes drained in-flight operations. Source queries/read counts snapshot and validation statements, excluding separately reported relay queries and post-trial audits. Fault-control timings below describe orchestrated pauses, not spontaneous outage distributions.\n\n## Correctness, resource counters and fault controls\n\n")
	for i, r := range rs {
		fmt.Fprintf(&text, "### Cell %d: %s / %s / %s / %s\n\nGate: `%t`. Error: `%s`. Wrong/impossible/errors: %d/%d/%d. Object rows after final trial: %d; source relation/index bytes: %d. Client cgroup periods/throttled: %d/%d; throttled microseconds: %d.\n\n", i+1, r.Arm, r.Protocol, r.Strategy, r.Kind, r.GatePassed, r.Error, r.FailureCounts.WrongReads, r.FailureCounts.ImpossibleReads, r.FailureCounts.Errors, r.ObjectRows, r.SourceBytes, r.ClientCPU.NrPeriods, r.ClientCPU.NrThrottled, r.ClientCPU.ThrottledUsec)
		for _, c := range r.Controls {
			fmt.Fprintf(&text, "- %s: pass=%t; %s; blocked/recovery/stale-duration ms=%.3f/%.3f/%.3f.\n", c.Name, c.Passed, c.Detail, c.BlockedMS, c.RecoveryMS, c.StaleDurationMS)
		}
		fmt.Fprintf(&text, "\nRedis evicted keys before/after: %s/%s. Cache policy: %s; maxmemory: %s bytes.\n\n", infoField(r.RedisBefore, "evicted_keys"), infoField(r.RedisAfter, "evicted_keys"), infoField(r.RedisAfter, "maxmemory_policy"), infoField(r.RedisAfter, "maxmemory"))
		names := []string{"snapshot", "version", "pending"}
		for _, name := range names {
			fmt.Fprintf(&text, "%s plan:\n\n```text\n%s\n```\n\n", name, r.Plans[name])
		}
	}
	text.WriteString("## Signed analyses index\n\n")
	as, ae := provenance.LoadAnalyses(os.DirFS(analyses), rs[0].RunID)
	if ae == nil {
		found := false
		for _, a := range as {
			if a.Digest == digest || a.Inputs[rs[0].RunID] == digest {
				fmt.Fprintf(&text, "- `%s`: %s (%s).\n", a.AnalysisID, a.Headline, a.File)
				found = true
			}
		}
		if !found {
			text.WriteString("No signed analysis currently indexes this digest.\n")
		}
	} else {
		text.WriteString("No signed analysis directory currently available.\n")
	}
	return writeExclusive(out, []byte(text.String()))
}
func infoField(s, k string) string {
	for _, line := range strings.Split(s, "\n") {
		if v, ok := strings.CutPrefix(strings.TrimSpace(line), k+":"); ok {
			return v
		}
	}
	return "unavailable"
}
