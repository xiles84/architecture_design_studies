package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type experimentFile struct {
	Run  *Run
	Path string
}

func readExperiments(dir string) ([]experimentFile, error) {
	var out []experimentFile
	err := filepath.WalkDir(dir, func(path string, e os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if e.IsDir() || !strings.HasSuffix(path, ".json") {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		var r Run
		if err := json.Unmarshal(b, &r); err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		if r.Experiment == nil {
			return nil
		}
		out = append(out, experimentFile{&r, path})
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no experiment results in %s", dir)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out, nil
}

func experimentProblems(r *Run) string {
	var p []string
	if r.Error != "" {
		p = append(p, r.Error)
	}
	if r.Verify == nil || r.Verify.Failed > 0 {
		p = append(p, "initial correctness gate missing/failed")
	}
	if r.Audit != nil && auditBad(r.Audit) > 0 {
		p = append(p, fmt.Sprintf("post-write invariant: %d mismatches", auditBad(r.Audit)))
	}
	for _, q := range r.Reads {
		if q.Errors > 0 || q.Ops == 0 {
			p = append(p, fmt.Sprintf("%s: %d errors / %d operations", q.Query, q.Errors, q.Ops))
		}
	}
	for _, w := range r.Writes {
		if w.Errors > 0 || w.Exhausted || w.Skipped != "" || w.Ops == 0 {
			p = append(p, fmt.Sprintf("write %s: errors=%d exhausted=%t skipped=%s", w.Op, w.Errors, w.Exhausted, w.Skipped))
		}
	}
	if a := r.Experiment.Arrival; a != nil {
		if !a.Reconciled {
			p = append(p, "acknowledgement reconciliation failed")
		}
		if a.WarmupAudit != nil && auditBad(a.WarmupAudit) > 0 {
			p = append(p, fmt.Sprintf("warmup invariant: %d mismatches", auditBad(a.WarmupAudit)))
		}
		// Offered-load errors and rejections are outcomes, not successful capacity.
		if a.Errors+a.WarmupErrors+a.ReadErrors > 0 {
			p = append(p, fmt.Sprintf("arrival/read errors=%d", a.Errors+a.WarmupErrors+a.ReadErrors))
		}
	}
	for _, phase := range r.Experiment.Growth {
		if phase.Verify == nil || phase.Verify.Failed > 0 {
			p = append(p, "growth correctness failed")
		}
		for _, q := range phase.Reads {
			if q.Errors > 0 {
				p = append(p, "growth read error")
			}
		}
	}
	return strings.Join(p, "; ")
}

func experimentReadScore(r *Run) float64 {
	var sum float64
	n := 0
	for _, q := range r.Reads {
		if strings.Contains(q.Query, "@") {
			continue
		}
		if q.OpsPerSec <= 0 || q.Errors > 0 {
			return 0
		}
		sum += math.Log(q.OpsPerSec)
		n++
	}
	if n != 12 {
		return 0
	}
	return math.Exp(sum / float64(n))
}

func writeEnhancementReport(dir, outPath string) error {
	if outPath == "" {
		return fmt.Errorf("-report-out is required")
	}
	files, err := readExperiments(dir)
	if err != nil {
		return err
	}
	first := files[0].Run
	groups := map[string][]experimentFile{}
	failed := 0
	cacheBad := 0
	for _, f := range files {
		r := f.Run
		if r.Environment != first.Environment || r.RunID != first.RunID {
			return fmt.Errorf("report cannot mix environments or run ids")
		}
		if experimentProblems(r) != "" {
			failed++
		}
		if r.Audit != nil && r.Audit.CacheMismatches > 0 {
			cacheBad++
		}
		o := r.Experiment.Settings
		key := r.Topology + " / " + o.Condition + " / " + o.Mode + " / " + r.DesignID
		groups[key] = append(groups[key], f)
	}
	var b strings.Builder
	fmt.Fprintf(&b, "# Study 01 v3 — %s\n\n## TL;DR\n\n", first.RunID)
	fmt.Fprintf(&b, "- %d recorded cells; %d cells have errors, missing verification, or invariant failures.\n- %d cells recorded a post-write cache mismatch.\n- %d distinct topology/condition/mode/design groups are listed below; every trial is retained.\n\n", len(files), failed, cacheBad, len(groups))
	fmt.Fprintf(&b, "Environment: `%s`. Commit: `%s`. Run tag: `%s`.\n\n", first.Environment, first.Provenance["repo_commit"], first.Provenance["run_tag"])
	fmt.Fprintf(&b, "Measurements only. Each trial begins with a fresh load and a Go correctness gate.\nRead score is the geometric mean of 12 isolated query rates, excluding targeted reads.\nRates below are medians of trial rates; spread = (maximum − minimum) / median.\nFewer than five trials are explicitly counted. Failed cells are excluded from valid-rate summaries.\nClosed-loop query tails are not SLO measurements. Arrival response latency starts at the\nscheduled due time, includes queue delay, and covers successful requests only: errors and\nrejections are reported separately. Throughput includes draining accepted requests.\nAll nodes and the client share one laptop and its local container network.\n\n")
	fmt.Fprintf(&b, "## Completeness and failures\n\nSee [manifest](../results/%s/manifest.yaml) for planned cells, startup failures and interruptions.\n\n", first.RunID)
	for _, f := range files {
		if problem := experimentProblems(f.Run); problem != "" {
			rel, _ := filepath.Rel(filepath.Dir(outPath), f.Path)
			fmt.Fprintf(&b, "- [%s](%s): %s\n", filepath.Base(f.Path), filepath.ToSlash(rel), strings.ReplaceAll(problem, "\n", " "))
		}
	}
	if failed == 0 {
		b.WriteString("No recorded cell failed its gates or reported operation errors.\n")
	}
	keys := make([]string, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fs := groups[key]
		sort.Slice(fs, func(i, j int) bool { return fs[i].Run.Experiment.Settings.Trial < fs[j].Run.Experiment.Settings.Trial })
		r := fs[0].Run
		o := r.Experiment.Settings
		fmt.Fprintf(&b, "\n## %s\n\n", key)
		fmt.Fprintf(&b, "Preparation `%s`; scale `%s`; history multiplier %d; largest charity %d, smallest %d; database memory limit %d bytes.\n\n", o.Preparation, r.Scale, o.HistoryMultiplier, r.Experiment.Keys["largest"], r.Experiment.Keys["smallest"], o.DBMemoryBytes)
		fmt.Fprintf(&b, "| Trial | Order | Gate | Donors | Donations | Initial relation bytes | Source |\n|---|---:|---|---:|---:|---:|---|\n")
		metrics := map[string][]float64{}
		durations := map[string][]float64{}
		scores := []float64{}
		for _, f := range fs {
			r := f.Run
			valid := experimentProblems(r) == ""
			status := "pass"
			if !valid {
				status = "FAILED / diagnostic"
			}
			rel, _ := filepath.Rel(filepath.Dir(outPath), f.Path)
			size := int64(0)
			if r.Stats != nil {
				size = r.Stats.TotalBytes
			}
			fmt.Fprintf(&b, "| %d | %d | %s | %v | %v | %d | [JSON](%s) |\n", r.Experiment.Settings.Trial, r.Experiment.Settings.Order, status, r.Dataset["people"], r.Dataset["donations"], size, filepath.ToSlash(rel))
			if !valid {
				continue
			}
			if score := experimentReadScore(r); score > 0 {
				scores = append(scores, score)
			}
			for _, q := range r.Reads {
				metrics[q.Query] = append(metrics[q.Query], q.OpsPerSec)
				durations[q.Query] = append(durations[q.Query], q.DurationS)
			}
			for _, w := range r.Writes {
				k := "write:" + w.Op
				metrics[k] = append(metrics[k], w.OpsPerSec)
				durations[k] = append(durations[k], w.DurationS)
			}
		}
		if len(scores) > 0 {
			fmt.Fprintf(&b, "\n12-query score: median **%.2f**, spread **%.1f%%**, %d independent loads. Trial scores: %v.\n", median(scores), spreadPct(scores), len(scores), scores)
		}
		if len(metrics) > 0 {
			fmt.Fprintf(&b, "\n| Operation | Valid trials | Median ops/s | Spread %% | Min–max measured seconds | Trial ops/s |\n|---|---:|---:|---:|---|---|\n")
			mk := make([]string, 0, len(metrics))
			for k := range metrics {
				mk = append(mk, k)
			}
			sort.Strings(mk)
			for _, k := range mk {
				xs := metrics[k]
				ds := append([]float64(nil), durations[k]...)
				sort.Float64s(ds)
				fmt.Fprintf(&b, "| %s | %d | %.2f | %.1f | %.3f–%.3f | %s |\n", k, len(xs), median(xs), spreadPct(xs), ds[0], ds[len(ds)-1], ratesString(xs))
			}
		}
		for _, f := range fs {
			r := f.Run
			if a := r.Experiment.Arrival; a != nil {
				fmt.Fprintf(&b, "\nTrial %d: offered %d; accepted %d; completed %d; rejected %d; write errors %d; read errors %d; warmup errors %d; reconciled %t.\n", r.Experiment.Settings.Trial, a.Offered, a.Accepted, a.Completed, a.Dropped, a.Errors, a.ReadErrors, a.WarmupErrors, a.Reconciled)
				fmt.Fprintf(&b, "Offered %.0f/s, %d fixed readers, %d writers, hot probability %.2f, hot donors %d. %.2f completed/s including drain, %.2f reads/s. Window %.2fs, including drain %.2fs.\n", o.ArrivalRate, o.ReadWorkers, o.WriteWorkers, o.HotProbability, o.HotDonors, a.CompletedPerSec, a.ReadsPerSec, a.WindowS, a.ElapsedS)
				fmt.Fprintf(&b, "Successful-request p99: service %.3fms; scheduled response %.3fms; queue delay %.3fms. Scheduler lag p99 %.3fms; successful samples %d.\n", a.Service.P99MS, a.Response.P99MS, a.QueueDelay.P99MS, a.SchedulerLag.P99MS, a.Service.Count)
			}
			if len(r.Experiment.Growth) > 0 {
				fmt.Fprintf(&b, "\nTrial %d growth (sequential mutations, concurrent reads measured after each phase):\n\n| Cycle | Operation | Count | Seconds | Ops/s | Relation bytes | Verification failures |\n|---|---|---:|---:|---:|---:|---:|\n", r.Experiment.Settings.Trial)
				for _, p := range r.Experiment.Growth {
					bytes := int64(0)
					if p.Stats != nil {
						bytes = p.Stats.TotalBytes
					}
					bad := -1
					if p.Verify != nil {
						bad = p.Verify.Failed
					}
					fmt.Fprintf(&b, "| %d | %s | %d | %.3f | %.2f | %d | %d |\n", p.Cycle, p.Operation, p.Operations, p.DurationS, p.OpsPerSec, bytes, bad)
				}
			}
		}
	}
	writeProvenance(&b, dir, filepath.Join(filepath.Dir(outPath), "analyses"), first.RunID)
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return err
	}
	// A generated edition is archived before replacement; never silently erase it.
	if old, err := os.ReadFile(outPath); err == nil {
		if string(old) == b.String() {
			return nil
		}
		archive := filepath.Join(filepath.Dir(outPath), "outdated", filepath.Base(outPath)+".previous")
		for i := 1; ; i++ {
			p := fmt.Sprintf("%s-%d.md", archive, i)
			if _, err := os.Stat(p); os.IsNotExist(err) {
				if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
					return err
				}
				if err := os.WriteFile(p, old, 0644); err != nil {
					return err
				}
				break
			}
		}
	}
	return os.WriteFile(outPath, []byte(b.String()), 0644)
}

func ratesString(xs []float64) string {
	var ss []string
	for _, x := range xs {
		ss = append(ss, fmt.Sprintf("%.2f", x))
	}
	return strings.Join(ss, ", ")
}
