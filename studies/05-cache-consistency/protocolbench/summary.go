package main

import (
	"adsplatform/core/provenance"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Trial medians are first reduced within each whole run, so many short trials
// cannot masquerade as many independently restarted runs.
func summarize(root, prefix, out string) error {
	rs, e := readResults(root)
	if e != nil {
		return e
	}
	type metrics struct{ Read, Write, Ack, ReadP99, WriteP99, DB, Source, Lock float64 }
	perRun := map[string]map[string]metrics{}
	positions := map[string]map[string]metrics{}
	digests := map[string]string{}
	controls := 0
	cells := 0
	maxThrottle := 0.
	evicted := 0
	var recovery, stale []float64
	for _, r := range rs {
		if !strings.HasPrefix(r.RunID, prefix+"-r") {
			continue
		}
		if !r.GatePassed {
			return fmt.Errorf("cannot summarize a failed cell %s", r.CellID)
		}
		cells++
		if perRun[r.RunID] == nil {
			perRun[r.RunID] = map[string]metrics{}
			positions[r.RunID] = map[string]metrics{}
			d, _, de := provenance.Digest(os.DirFS(filepath.Join(root, r.RunID)))
			if de != nil {
				return de
			}
			digests[r.RunID] = d
		}
		if r.ClientCPU.NrPeriods > 0 {
			v := float64(r.ClientCPU.NrThrottled) / float64(r.ClientCPU.NrPeriods)
			if v > maxThrottle {
				maxThrottle = v
			}
		}
		if infoField(r.RedisAfter, "evicted_keys") != "0" {
			evicted++
		}
		for _, c := range r.Controls {
			if !c.Passed {
				return fmt.Errorf("failed control")
			}
			controls++
			if c.RecoveryMS > 0 {
				recovery = append(recovery, c.RecoveryMS)
			}
			if c.StaleDurationMS > 0 {
				stale = append(stale, c.StaleDurationMS)
			}
		}
		phases := map[string][]trial{}
		for _, t := range r.Trials {
			phases[t.Phase] = append(phases[t.Phase], t)
		}
		for phase, ts := range phases {
			if len(ts) != 3 {
				return fmt.Errorf("expected three trials for %s", phase)
			}
			var a, b, c, d, f, g, h, i []float64
			for _, t := range ts {
				a = append(a, t.ReadOps)
				b = append(b, t.WriteOps)
				c = append(c, t.AckP50MS)
				d = append(d, t.ReadP99MS)
				f = append(f, t.WriteP99MS)
				g = append(g, t.DBP50MS)
				h = append(h, t.SourceQueriesPerRead)
				i = append(i, t.LockP99MS)
			}
			m := metrics{percentile(a, .5), percentile(b, .5), percentile(c, .5), percentile(d, .5), percentile(f, .5), percentile(g, .5), percentile(h, .5), percentile(i, .5)}
			key := strings.Join([]string{r.Arm, r.Protocol, r.Strategy, r.Kind, phase}, "/")
			if strings.Contains(r.CellID, "position-last") {
				positions[r.RunID][key] = m
			} else {
				perRun[r.RunID][key] = m
			}
		}
	}
	if len(perRun) != 3 || cells != 30 {
		return fmt.Errorf("expected three complete whole runs/30 cells, got %d/%d", len(perRun), cells)
	}
	var text strings.Builder
	text.WriteString("# Whole-run reduction of cache protocol measurements\n\nFixed rule: median of three trials within a whole run, then median/min/max across three restarted whole runs. Position-last replicas are excluded from the ordinary reduction and reported separately. No inferential significance or crossover is assigned.\n\n")
	runs := []string{}
	for r := range perRun {
		runs = append(runs, r)
	}
	sort.Strings(runs)
	for _, r := range runs {
		fmt.Fprintf(&text, "- %s: inputs digest `%s`.\n", r, digests[r])
	}
	fmt.Fprintf(&text, "\nCells/control gates: %d/%d. Nonzero Redis-eviction cells: %d. Maximum client throttled-period fraction over a complete cell: %.4f.\n\n", cells, controls, evicted, maxThrottle)
	keys := []string{}
	for k := range perRun[runs[0]] {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	text.WriteString("| Arm/protocol/strategy/deletion/phase | Read/s median [min,max] | Write/s median [min,max] | Ack p50 ms median [min,max] | Read p99 ms median [min,max] | Write p99 ms median [min,max] | Source queries/read median [min,max] |\n|---|---|---|---|---|---|---|\n")
	spread := func(xs []float64) string {
		sort.Float64s(xs)
		return fmt.Sprintf("%.3f [%.3f,%.3f]", percentile(xs, .5), xs[0], xs[len(xs)-1])
	}
	for _, k := range keys {
		var a, b, c, d, f, h []float64
		for _, r := range runs {
			m, ok := perRun[r][k]
			if !ok {
				return fmt.Errorf("missing cell %s in %s", k, r)
			}
			a = append(a, m.Read)
			b = append(b, m.Write)
			c = append(c, m.Ack)
			d = append(d, m.ReadP99)
			f = append(f, m.WriteP99)
			h = append(h, m.Source)
		}
		fmt.Fprintf(&text, "| %s | %s | %s | %s | %s | %s | %s |\n", k, spread(a), spread(b), spread(c), spread(d), spread(f), spread(h))
	}
	text.WriteString("\n## Same-candidate last/first position ratios\n\n| Whole run / cell / phase | Read/s ratio | Write/s ratio |\n|---|---:|---:|\n")
	for _, r := range runs {
		pk := []string{}
		for k := range positions[r] {
			pk = append(pk, k)
		}
		sort.Strings(pk)
		for _, k := range pk {
			last := positions[r][k]
			first, ok := perRun[r][k]
			if !ok {
				return fmt.Errorf("missing position baseline")
			}
			rr, wr := 0., 0.
			if first.Read > 0 {
				rr = last.Read / first.Read
			}
			if first.Write > 0 {
				wr = last.Write / first.Write
			}
			fmt.Fprintf(&text, "| %s / %s | %.4f | %.4f |\n", r, k, rr, wr)
		}
	}
	if len(recovery) > 0 {
		fmt.Fprintf(&text, "\nOrchestrated relay-resume recovery ms, median [min,max]: %s. Unsafe-barrier stale duration ms (includes the forced 30-ms pause), median [min,max]: %s.\n", spread(recovery), spread(stale))
	}
	return writeExclusive(out, []byte(text.String()))
}
