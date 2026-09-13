package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"adsplatform/adapters/cgroup"
	"adsplatform/adapters/filestore"
	md "adsplatform/adapters/markdown"
	"adsplatform/core/measure"
	"adsplatform/core/provenance"
)

// ---------------------------------------------------------------------------
// Report generation: measurements only, no interpretation. Every number comes
// from exactly one result file; conclusions belong in a signed analysis.
//
// One rule is applied throughout: a number produced while a design violated the
// capacity invariant is never shown as a plain number. It is marked ❌ next to
// the violation, because a design that oversells quickly has not sold quickly.
// ---------------------------------------------------------------------------

var topologyOrder = []string{"pg-single", "yb-single", "yb-cluster3"}

var topologyLabel = map[string]string{
	"pg-single":   "PostgreSQL, 1 node",
	"yb-single":   "YugabyteDB, 1 node (RF=1)",
	"yb-cluster3": "YugabyteDB, 3 nodes (RF=3)",
}

var queryLabel = map[string]string{
	"q01_event_availability": "Seats left for an event",
	"q02_band_events":        "A band's events with availability",
	"q03_customer_tickets":   "A customer's tickets",
	"q04_ticket_by_id":       "Ticket by id (control)",
	"q05_band_tickets_sold":  "Tickets sold across a band",
}

type pair struct{ a, b, question string }

var pairs = []pair{
	{"p1_precreated_lock_first", "p2_precreated_skip_locked", "What does SKIP LOCKED buy over a plain FOR UPDATE on pre-created seats?"},
	{"p2_precreated_skip_locked", "p3_precreated_cas", "Locking read with SKIP LOCKED, or no lock and compare-and-set?"},
	{"p2_precreated_skip_locked", "p4_precreated_counter", "What does a denormalised availability counter cost the booking path?"},
	{"p4_precreated_counter", "c4_counter_guard", "Same counter serialising both: update a pre-created ticket, or insert a new one?"},
	{"c1_count_naive", "c2_count_serializable", "Identical SQL: what does SERIALIZABLE cost, and does it prevent overbooking?"},
	{"c1_count_naive", "c3_count_lock_event", "What does locking the parent event row cost, and does it prevent overbooking?"},
	{"c3_count_lock_event", "c4_counter_guard", "Inside the same per-event critical section: count the tickets, or keep a counter?"},
	{"c4_counter_guard", "c5_seat_unique", "A guarded counter row, or a unique constraint as the arbiter?"},
	{"c4_counter_guard", "r1_inventory_row", "What does moving the hot counter off the event row buy?"},
	{"r1_inventory_row", "r2_inventory_buckets", "What does sharding the counter into buckets buy?"},
	{"p2_precreated_skip_locked", "r3_seat_pool", "SKIP LOCKED on a wide pre-created ticket, or on a narrow slot plus a ticket insert?"},
	{"h0_hold_naive_confirm", "h1_hold_checked_confirm", "What does validating the hold at confirmation cost?"},
}

type report struct {
	dir   string
	runs  map[string]map[string]*Run
	topos []string
	ids   []string
}

func loadReport(dir string) (*report, error) {
	all, err := filestore.ReadJSONTree[Run](dir)
	if err != nil {
		return nil, err
	}
	rp := &report{dir: dir, runs: map[string]map[string]*Run{}}
	for i := range all {
		r := &all[i]
		if r.Study != studyID || r.DesignID == "" {
			continue
		}
		if rp.runs[r.Topology] == nil {
			rp.runs[r.Topology] = map[string]*Run{}
		}
		rp.runs[r.Topology][r.DesignID] = r
	}
	for _, t := range topologyOrder {
		if rp.runs[t] != nil {
			rp.topos = append(rp.topos, t)
		}
	}
	for _, d := range designs {
		for _, t := range rp.topos {
			if rp.runs[t][d.ID] != nil {
				rp.ids = append(rp.ids, d.ID)
				break
			}
		}
	}
	if len(rp.topos) == 0 {
		return nil, fmt.Errorf("no study %s results under %s", studyID, dir)
	}
	return rp, nil
}

func (rp *report) designsIn(t string) []string {
	var out []string
	for _, id := range rp.ids {
		if rp.runs[t][id] != nil {
			out = append(out, id)
		}
	}
	return out
}

func (rp *report) any() *Run {
	for _, t := range rp.topos {
		for _, id := range rp.designsIn(t) {
			return rp.runs[t][id]
		}
	}
	return nil
}

// raceOf combines every trial of one (mode, tier). Throughput is the median
// trial's; latencies are that median trial's (percentiles cannot be merged
// without the samples); violations and timeouts are SUMMED over trials, so an
// overbooking in any single trial stays visible.
func raceOf(r *Run, mode string, tier int) *RaceResult {
	if r == nil {
		return nil
	}
	var trials []*RaceResult
	for i := range r.Races {
		if r.Races[i].Mode == mode && r.Races[i].Tier == tier {
			trials = append(trials, &r.Races[i])
		}
	}
	switch len(trials) {
	case 0:
		return nil
	case 1:
		return trials[0]
	}
	sort.Slice(trials, func(i, j int) bool { return trials[i].SoldPerSec < trials[j].SoldPerSec })
	med := *trials[len(trials)/2]
	agg := med
	agg.Events, agg.Seats, agg.Sold, agg.Successes, agg.Retries, agg.Errors = 0, 0, 0, 0, 0, 0
	agg.Overbooked, agg.OverbookedSeats, agg.Underbooked, agg.UnderbookedSeats, agg.TimedOut = 0, 0, 0, 0, 0
	agg.Cancels, agg.SweepSold = 0, 0
	agg.Examples = nil
	var sps []float64
	for _, t := range trials {
		agg.Events += t.Events
		agg.Seats += t.Seats
		agg.Sold += t.Sold
		agg.Successes += t.Successes
		agg.Retries += t.Retries
		agg.Errors += t.Errors
		agg.Cancels += t.Cancels
		agg.SweepSold += t.SweepSold
		agg.Overbooked += t.Overbooked
		agg.OverbookedSeats += t.OverbookedSeats
		agg.Underbooked += t.Underbooked
		agg.UnderbookedSeats += t.UnderbookedSeats
		agg.TimedOut += t.TimedOut
		agg.Examples = append(agg.Examples, t.Examples...)
		if t.Audit != nil && (agg.Audit == nil || t.Audit.Violations() > agg.Audit.Violations()) {
			agg.Audit = t.Audit
		}
		sps = append(sps, t.SoldPerSec)
	}
	if agg.Successes > 0 {
		agg.AttemptsPerSuccess = float64(agg.Successes+agg.Retries) / float64(agg.Successes)
	}
	agg.SoldPerSec = measure.Median(sps)
	agg.TrialsSoldPerSec = sps
	agg.SpreadPct = measure.SpreadPct(sps)
	return &agg
}

func readOf(r *Run, q string, tier int) *ReadResult {
	if r == nil {
		return nil
	}
	for i := range r.Reads {
		if r.Reads[i].Query == q && r.Reads[i].Tier == tier {
			return &r.Reads[i]
		}
	}
	return nil
}

func writeOf(r *Run, op string, tier int) *WriteResult {
	if r == nil {
		return nil
	}
	for i := range r.Writes {
		if r.Writes[i].Op == op && r.Writes[i].Tier == tier {
			return &r.Writes[i]
		}
	}
	return nil
}

func (rp *report) raceTiers(mode string) []int {
	seen := map[int]bool{}
	var out []int
	for _, t := range rp.topos {
		for _, r := range rp.runs[t] {
			for _, rr := range r.Races {
				if rr.Mode == mode && !seen[rr.Tier] {
					seen[rr.Tier] = true
					out = append(out, rr.Tier)
				}
			}
		}
	}
	sort.Ints(out)
	return out
}

func tierLabel(t int) string {
	if t >= 1000 {
		return fmt.Sprintf("%dk seats", t/1000)
	}
	return fmt.Sprintf("%d seats", t)
}

// raceCell renders sold/s with every caveat that applies to it.
func raceCell(rr *RaceResult) string {
	if rr == nil {
		return "—"
	}
	v := md.Ops(rr.SoldPerSec)
	if len(rr.TrialsSoldPerSec) > 1 {
		v += fmt.Sprintf(" (%d trials, spread %.0f%%)", len(rr.TrialsSoldPerSec), rr.SpreadPct)
	}
	var marks []string
	if rr.Overbooked > 0 {
		marks = append(marks, fmt.Sprintf("❌ overbooked %d/%d (+%d)", rr.Overbooked, rr.Events, rr.OverbookedSeats))
	}
	if rr.Underbooked > 0 {
		marks = append(marks, fmt.Sprintf("⚠️ underbooked %d/%d (−%d)", rr.Underbooked, rr.Events, rr.UnderbookedSeats))
	}
	if rr.Audit != nil && rr.Audit.DriftEvents > 0 {
		marks = append(marks, fmt.Sprintf("⚠️ drift %d", rr.Audit.DriftEvents))
	}
	if rr.TimedOut > 0 {
		marks = append(marks, fmt.Sprintf("⏱ %d timed out at %.0f%% sold", rr.TimedOut, rr.TimedOutSoldPct))
	}
	if rr.EventsPlanned > rr.Events {
		marks = append(marks, fmt.Sprintf("⌛ only %d of %d events raced within the tier budget", rr.Events, rr.EventsPlanned))
	}
	if rr.Errors > 0 {
		marks = append(marks, fmt.Sprintf("%d errors", rr.Errors))
	}
	if len(marks) == 0 {
		return v
	}
	if rr.Overbooked > 0 {
		return fmt.Sprintf("~~%s~~ %s", v, strings.Join(marks, ", "))
	}
	return fmt.Sprintf("%s %s", v, strings.Join(marks, ", "))
}

func WriteReport(dir, outPath string) error {
	rp, err := loadReport(dir)
	if err != nil {
		return err
	}
	ref := rp.any()
	var b strings.Builder

	fmt.Fprintf(&b, "# Study 02 — results: avoiding overbooking (band → event → ticket)\n\n")
	fmt.Fprintf(&b, "Generated from `results/%s`. Every number comes from exactly one JSON file in that\n", filepath.Base(dir))
	fmt.Fprintf(&b, "directory. **This file contains measurements and no interpretation**; conclusions live in\n")
	fmt.Fprintf(&b, "signed analyses, indexed at the end.\n\n")

	rp.writeTLDR(&b)
	rp.writeSetup(&b, ref)
	rp.writeGate(&b)
	rp.writeRace(&b, "race")
	rp.writeRaceDetail(&b)
	rp.writeRace(&b, "churn")
	rp.writePrecreation(&b)
	rp.writeReads(&b)
	rp.writeHolds(&b)
	rp.writePairs(&b)
	rp.writeHazards(&b)

	fmt.Fprintf(&b, "\n---\n\nPlans for every read and write statement are in `results/%s/<topology>/plans/<design>.txt`\n", filepath.Base(dir))
	fmt.Fprintf(&b, "(`EXPLAIN (ANALYZE, BUFFERS)` on PostgreSQL, `EXPLAIN (ANALYZE, DIST)` on YugabyteDB; write\n")
	fmt.Fprintf(&b, "statements inside rolled-back transactions).\n")

	var versions []provenance.RepoVersion
	for _, t := range rp.topos {
		for _, id := range rp.designsIn(t) {
			versions = append(versions, rp.runs[t][id].Repo)
		}
	}
	provenance.WriteSection(&b, ref.RunID, versions, os.DirFS(dir), os.DirFS(filepath.Join(filepath.Dir(outPath), "analyses")))

	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(outPath, []byte(b.String()), 0o644)
}

// writeTLDR is the report's summary: facts selected from the tables below by
// fixed rules, with no judgement about which design is better. Rankings are by
// the measured number only, and only among designs that did not break the
// invariant in that measurement. Interpretation belongs in a signed analysis,
// which has its own TL;DR.
func (rp *report) writeTLDR(b *strings.Builder) {
	fmt.Fprintf(b, "## TL;DR — measured facts\n\n")
	fmt.Fprintf(b, "*Selected from the tables below by fixed rules; no interpretation. Read a signed analysis for\n")
	fmt.Fprintf(b, "what these facts mean.*\n\n")

	cells, failed := 0, []string{}
	for _, tp := range rp.topos {
		for _, id := range rp.designsIn(tp) {
			cells++
			if e := rp.runs[tp][id].Error; e != "" {
				failed = append(failed, fmt.Sprintf("%s / %s (%s)", topologyLabel[tp], designShort(id), firstSentence(e)))
			}
		}
	}
	// A cell killed before it could write its JSON leaves no result file at all;
	// only the runner's manifest knows it existed.
	for _, fc := range rp.manifestFailures() {
		tp, id, _ := strings.Cut(fc, "/")
		if r := rp.runs[tp][id]; r != nil && r.Error != "" {
			continue // already listed with its error
		}
		if rp.runs[tp][id] == nil {
			failed = append(failed, fmt.Sprintf("%s / %s (no result file — see logs/)", topologyLabel[tp], designShort(id)))
		}
	}
	fmt.Fprintf(b, "- **Cells:** %d across %d topolog%s; %d failed%s.\n", cells, len(rp.topos), plural(len(rp.topos), "y", "ies"),
		len(failed), listSuffix(failed))

	// Invariant outcomes, per design, over every race/churn/holds/write audit.
	var broken, under []string
	controlsOK, controlsTotal := 0, 0
	for _, tp := range rp.topos {
		for _, d := range designs {
			r := rp.runs[tp][d.ID]
			if r == nil {
				continue
			}
			var over, underSeats int64
			for _, rr := range r.Races {
				over += rr.OverbookedSeats
				underSeats += rr.UnderbookedSeats
			}
			for _, h := range r.Holds {
				over += h.OverbookedSeats
			}
			if d.NegativeControl != "" {
				controlsTotal++
				if over > 0 {
					controlsOK++
				}
				continue
			}
			if over > 0 {
				broken = append(broken, fmt.Sprintf("%s on %s (+%d seats)", d.Short, topologyLabel[tp], over))
			}
			if underSeats > 0 {
				under = append(under, fmt.Sprintf("%s on %s (%d seats)", d.Short, topologyLabel[tp], underSeats))
			}
		}
	}
	fmt.Fprintf(b, "- **Negative controls:** %d of %d fired (overbooked, as designed to).\n", controlsOK, controlsTotal)
	if len(broken) == 0 {
		fmt.Fprintf(b, "- **Overbooking in designs meant to be correct:** none observed.\n")
	} else {
		fmt.Fprintf(b, "- ❌ **Overbooking in designs meant to be correct:** %s.\n", strings.Join(broken, "; "))
	}
	if len(under) == 0 {
		fmt.Fprintf(b, "- **Under-booking** (sold out with seats unsold): none observed.\n")
	} else {
		fmt.Fprintf(b, "- ⚠️ **Under-booking** (sold out with seats unsold): %s.\n", strings.Join(under, "; "))
	}

	// Fastest and slowest valid race result per tier, per topology.
	for _, tp := range rp.topos {
		var parts []string
		for _, tier := range rp.raceTiers("race") {
			type cand struct {
				id string
				v  float64
				to bool
			}
			var cs []cand
			for _, id := range rp.designsIn(tp) {
				if d, _ := designByID(id); d.NegativeControl != "" {
					continue
				}
				rr := raceOf(rp.runs[tp][id], "race", tier)
				if rr == nil || rr.Overbooked > 0 || rr.SoldPerSec <= 0 {
					continue
				}
				cs = append(cs, cand{id, rr.SoldPerSec, rr.TimedOut > 0})
			}
			if len(cs) == 0 {
				continue
			}
			sort.Slice(cs, func(i, j int) bool { return cs[i].v > cs[j].v })
			hi, lo := cs[0], cs[len(cs)-1]
			timedOut := 0
			for _, c := range cs {
				if c.to {
					timedOut++
				}
			}
			p := fmt.Sprintf("%s: highest %s %s/s, lowest %s %s/s", tierLabel(tier), designShort(hi.id), md.Ops(hi.v), designShort(lo.id), md.Ops(lo.v))
			if timedOut > 0 {
				p += fmt.Sprintf(" (%d of %d timed out)", timedOut, len(cs))
			}
			parts = append(parts, p)
		}
		if len(parts) > 0 {
			fmt.Fprintf(b, "- **Sell-out race, %s** (sold/s; designs meant to be correct, without overbooking):\n", topologyLabel[tp])
			for _, p := range parts {
				fmt.Fprintf(b, "  - %s\n", p)
			}
		}
	}

	// The pre-creation price, at the largest publish tier present.
	for _, tp := range rp.topos {
		var pre, cre []float64
		tier := 0
		for _, t := range allTiers {
			for _, id := range rp.designsIn(tp) {
				if writeOf(rp.runs[tp][id], "publish", t) != nil {
					tier = t
				}
			}
		}
		if tier == 0 {
			continue
		}
		for _, d := range designs {
			w := writeOf(rp.runs[tp][d.ID], "publish", tier)
			if w == nil || w.Latency.P50MS <= 0 {
				continue
			}
			if d.Precreated || d.SeatPool {
				pre = append(pre, w.Latency.P50MS)
			} else {
				cre = append(cre, w.Latency.P50MS)
			}
		}
		if len(pre) > 0 && len(cre) > 0 {
			sort.Float64s(pre)
			sort.Float64s(cre)
			fmt.Fprintf(b, "- **Publishing a %s event, %s** (p50): %s–%s ms with a row per seat (P1–P4, R3), %s–%s ms without.\n",
				tierLabel(tier), topologyLabel[tp], md.MS(pre[0]), md.MS(pre[len(pre)-1]), md.MS(cre[0]), md.MS(cre[len(cre)-1]))
		}
	}
	fmt.Fprintf(b, "\n")
}

// manifestFailures reads "topology/design" entries under failed_cells: in every
// pass of the run's manifest.
func (rp *report) manifestFailures() []string {
	b, err := os.ReadFile(filepath.Join(rp.dir, "manifest.yaml"))
	if err != nil {
		return nil
	}
	var out []string
	in := false
	for _, line := range strings.Split(string(b), "\n") {
		t := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(t, "failed_cells:"):
			in = true
		case in && strings.HasPrefix(t, "- "):
			out = append(out, strings.TrimSpace(strings.TrimPrefix(t, "- ")))
		default:
			in = false
		}
	}
	return out
}

func plural(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}

func listSuffix(xs []string) string {
	if len(xs) == 0 {
		return ""
	}
	return ": " + strings.Join(xs, "; ")
}

func (rp *report) writeSetup(b *strings.Builder, ref *Run) {
	fmt.Fprintf(b, "## What was measured\n\n")
	t := md.NewTable("", "")
	t.Row("Run id", "`"+ref.RunID+"`")
	t.Row("Environment", fmt.Sprintf("[`%s`](../../../docs/environments/%s.md)", ref.Environment, ref.Environment))
	t.Row("Scale", fmt.Sprintf("`%s` — %v bands, %v events, %v seats, %v tickets sold at load",
		ref.Scale, ref.Dataset["bands"], ref.Dataset["events"], ref.Dataset["seats_total"], ref.Dataset["tickets_sold_at_load"]))
	t.Row("Events (count × seats)", eventsLine(ref.Dataset["events_by_kind_tier"]))
	t.Row("Seed", fmt.Sprintf("%v — identical data in every cell", ref.Dataset["seed"]))
	o := ref.Options
	t.Row("Reads / isolated writes", fmt.Sprintf("%v workers, %v measured + %v warmup, %v trial(s)", o["conns"], o["duration"], o["warmup"], o["trials"]))
	t.Row("Sell-out race", fmt.Sprintf("%v buyers per event, timeout %v, organiser edit every %v; churn race cancels %v%% of sales",
		o["race_buyers"], o["race_timeout"], o["editor_interval"], o["churn_pct"]))
	for _, tp := range rp.topos {
		for _, id := range rp.designsIn(tp) {
			if r := rp.runs[tp][id]; len(r.Holds) > 0 {
				o := r.Options
				t.Row("Holds — "+topologyLabel[tp], fmt.Sprintf("%v buyers, TTL %v, basket ≤ %v, payment %v, %v%% abandoned, sweeper every %v (time scale %v)",
					o["hold_buyers"], o["hold_ttl"], o["hold_think_max"], o["hold_pay"], o["hold_abandon_pct"], o["hold_sweep_every"], o["hold_time_scale"]))
				break
			}
		}
	}
	t.Write(b)

	fmt.Fprintf(b, "**Engines**\n\n")
	et := md.NewTable("Topology", "Version", "READ COMMITTED requests actually run as")
	for _, tp := range rp.topos {
		for _, id := range rp.designsIn(tp) {
			r := rp.runs[tp][id]
			eff := r.EngineInfo.EffectiveIsolation
			if eff == "" {
				eff = "read committed"
			}
			et.Row(topologyLabel[tp], "`"+firstLine(r.EngineInfo.Version)+"`", eff)
			break
		}
	}
	et.Write(b)
	fmt.Fprintf(b, "> Limits of this environment that bound every number below: the client and the\n")
	fmt.Fprintf(b, "> databases share one laptop's eight cores, and the \"3-node\" cluster has no real network\n")
	fmt.Fprintf(b, "> between nodes. The race is a closed loop (buyers wait for their answer), so its tails\n")
	fmt.Fprintf(b, "> are a floor, not an SLO figure.\n\n")
}

// eventsLine renders {"race_100": 20, ...} as "race: 100×20 · ...", by kind.
func eventsLine(v any) string {
	m, ok := v.(map[string]any)
	if !ok {
		return "—"
	}
	byKind := map[string][]string{}
	for _, kind := range []string{"catalogue", "race", "churn"} {
		for _, tier := range allTiers {
			if n, ok := m[fmt.Sprintf("%s_%d", kind, tier)]; ok {
				byKind[kind] = append(byKind[kind], fmt.Sprintf("%v × %s", n, tierLabel(tier)))
			}
		}
	}
	var parts []string
	for _, kind := range []string{"catalogue", "race", "churn"} {
		if len(byKind[kind]) > 0 {
			parts = append(parts, fmt.Sprintf("**%s**: %s", kind, strings.Join(byKind[kind], ", ")))
		}
	}
	return strings.Join(parts, " · ")
}

func (rp *report) writeGate(b *strings.Builder) {
	fmt.Fprintf(b, "## Correctness gate\n\n")
	fmt.Fprintf(b, "Every design must answer the read questions exactly as computed from the generated data and\n")
	fmt.Fprintf(b, "pass the overbooking audit on the freshly loaded tables before anything is timed. The audit\n")
	fmt.Fprintf(b, "then runs again after every phase that writes, on the state that phase left behind: events\n")
	fmt.Fprintf(b, "over capacity, seats sold twice, counters/buckets/pools that disagree with the tickets, and a\n")
	fmt.Fprintf(b, "reconciliation of the tickets that exist against the sales the harness saw commit.\n\n")

	t := md.NewTable("Topology", "Design", "Answers", "Audit on load", "After isolated writes", "After races", "After holds", "Cell")
	for _, tp := range rp.topos {
		for _, id := range rp.designsIn(tp) {
			r := rp.runs[tp][id]
			ans := "—"
			if r.Verify != nil {
				mark := "✅"
				if r.Verify.Failed > 0 {
					mark = "❌"
				}
				ans = fmt.Sprintf("%s %d/%d", mark, r.Verify.Passed, r.Verify.Passed+r.Verify.Failed)
			}
			cell := "✅"
			if r.Error != "" {
				cell = "❌ " + firstSentence(r.Error)
			}
			var wa, ra, ha []*Audit
			for i := range r.Writes {
				wa = append(wa, r.Writes[i].Audit)
			}
			for i := range r.Races {
				ra = append(ra, r.Races[i].Audit)
			}
			for i := range r.Holds {
				ha = append(ha, r.Holds[i].Audit)
			}
			t.Row(topologyLabel[tp], designShort(id), ans, auditCell(r.LoadAudit), auditsCell(wa), auditsCell(ra), auditsCell(ha), cell)
		}
	}
	t.Write(b)

	fmt.Fprintf(b, "### Negative controls\n\n")
	fmt.Fprintf(b, "Two designs are in the study **because they are wrong**. An audit that has never caught a\n")
	fmt.Fprintf(b, "wrong design has not been shown to work, so each control must be seen to fail. If a control\n")
	fmt.Fprintf(b, "did not fail, the absence of violations in the other designs of that experiment is not\n")
	fmt.Fprintf(b, "evidence of their correctness — only of too little contention.\n\n")
	nt := md.NewTable("Topology", "Control", "Expected failure", "Observed")
	for _, tp := range rp.topos {
		for _, d := range designs {
			if d.NegativeControl == "" || rp.runs[tp][d.ID] == nil {
				continue
			}
			r := rp.runs[tp][d.ID]
			var events, seats int64
			if d.ControlFires == "race" {
				for _, rr := range r.Races {
					events += int64(rr.Overbooked)
					seats += rr.OverbookedSeats
				}
			} else {
				for _, h := range r.Holds {
					events += int64(h.Overbooked)
					seats += h.OverbookedSeats
				}
			}
			obs := fmt.Sprintf("✅ fired: %d events overbooked by %d seats in total", events, seats)
			if events == 0 {
				obs = "⚠️ **did not fire** — this experiment did not generate enough contention to prove the audit"
			}
			nt.Row(topologyLabel[tp], d.Short, d.NegativeControl, obs)
		}
	}
	nt.Write(b)
}

func auditCell(a *Audit) string {
	if a == nil {
		return "—"
	}
	if a.Violations() == 0 {
		return "✅"
	}
	var parts []string
	if a.OverbookedEvents > 0 {
		parts = append(parts, fmt.Sprintf("%d overbooked", a.OverbookedEvents))
	}
	if a.DuplicateSeats > 0 {
		parts = append(parts, fmt.Sprintf("%d dup seats", a.DuplicateSeats))
	}
	if a.DriftEvents > 0 {
		parts = append(parts, fmt.Sprintf("%d drift", a.DriftEvents))
	}
	if a.LedgerMismatches > 0 {
		parts = append(parts, fmt.Sprintf("%d ledger", a.LedgerMismatches))
	}
	return "❌ " + strings.Join(parts, ", ")
}

// auditsCell shows the worst of several audits. Race audits run after each tier
// on the whole database, so a violation carries forward into later tiers; the
// worst one is the honest summary.
func auditsCell(as []*Audit) string {
	var worst *Audit
	for _, a := range as {
		if a != nil && (worst == nil || a.Violations() > worst.Violations()) {
			worst = a
		}
	}
	if worst == nil {
		return "—"
	}
	return auditCell(worst)
}

func (rp *report) writeRace(b *strings.Builder, mode string) {
	tiers := rp.raceTiers(mode)
	if len(tiers) == 0 {
		return
	}
	if mode == "race" {
		fmt.Fprintf(b, "\n## The sell-out race — seats sold per second (higher is better)\n\n")
		fmt.Fprintf(b, "Every buyer arrives at an unsold event at the same instant and keeps buying until told\n")
		fmt.Fprintf(b, "\"sold out\". Throughput is successful sales over the time from the start to the last sale\n")
		fmt.Fprintf(b, "(summed over the tier's events); the time spent turning the rest of the crowd away is\n")
		fmt.Fprintf(b, "reported separately below as *sold out answered in*. Markers:\n")
		fmt.Fprintf(b, "❌ overbooked (events/total, extra seats) and struck through — **not a valid speed**;\n")
		fmt.Fprintf(b, "⚠️ underbooked (every buyer was sent away with seats unsold); ⏱ race hit its timeout.\n\n")
	} else {
		fmt.Fprintf(b, "\n## Churn race — selling out while buyers cancel\n\n")
		fmt.Fprintf(b, "As the race, but a share of successful buyers cancel immediately. After the crowd leaves,\n")
		fmt.Fprintf(b, "one last buyer sweeps up refunded seats; seats still unsold after that sweep are\n")
		fmt.Fprintf(b, "**underbooked** — the design said \"sold out\" while it had seats.\n\n")
	}
	for _, tp := range rp.topos {
		ids := rp.designsIn(tp)
		head := []string{"Design"}
		for _, t := range tiers {
			head = append(head, tierLabel(t)+">")
		}
		t := md.NewTable(head...)
		for _, id := range ids {
			row := []string{designShort(id)}
			for _, tier := range tiers {
				row = append(row, raceCell(raceOf(rp.runs[tp][id], mode, tier)))
			}
			t.Row(row...)
		}
		fmt.Fprintf(b, "### %s\n\n", topologyLabel[tp])
		t.Write(b)
		if mode == "churn" {
			ut := md.NewTable(append([]string{"Design"}, tierHeads(tiers, " — cancels / swept up / underbooked seats")...)...)
			for _, id := range ids {
				row := []string{designShort(id)}
				for _, tier := range tiers {
					rr := raceOf(rp.runs[tp][id], mode, tier)
					if rr == nil {
						row = append(row, "—")
						continue
					}
					row = append(row, fmt.Sprintf("%d / %d / %d", rr.Cancels, rr.SweepSold, rr.UnderbookedSeats))
				}
				ut.Row(row...)
			}
			ut.Write(b)
		}
	}
}

func tierHeads(tiers []int, suffix string) []string {
	var out []string
	for i, t := range tiers {
		h := tierLabel(t)
		if i == 0 {
			h += suffix
		}
		out = append(out, h+">")
	}
	return out
}

func (rp *report) writeRaceDetail(b *strings.Builder) {
	tiers := rp.raceTiers("race")
	if len(tiers) == 0 {
		return
	}
	fmt.Fprintf(b, "\n## The race in detail\n\n")
	fmt.Fprintf(b, "How each sale was achieved, not just how many. *Attempts per sale* is (sales + retries) ÷\n")
	fmt.Fprintf(b, "sales: 1.00 means no buyer ever lost a race and started again. *Sold out answered in* is the\n")
	fmt.Fprintf(b, "median latency of the final \"no\" each buyer received. *Organiser edit p99* is the latency of an\n")
	fmt.Fprintf(b, "unrelated update to the event page made while the crowd was buying.\n\n")
	for _, tp := range rp.topos {
		fmt.Fprintf(b, "### %s\n\n", topologyLabel[tp])
		for _, metric := range []struct {
			title string
			cell  func(*RaceResult) string
		}{
			{"Sale latency p99 (ms)", func(r *RaceResult) string { return md.MS(r.Latency.P99MS) }},
			{"Attempts per sale", func(r *RaceResult) string {
				if r.Successes == 0 {
					return "—"
				}
				return fmt.Sprintf("%.2f", r.AttemptsPerSuccess)
			}},
			{"Sold out answered in, p50 (ms)", func(r *RaceResult) string { return md.MS(r.RejectLatency.P50MS) }},
			{"Organiser edit p99 (ms)", func(r *RaceResult) string { return md.MS(r.EditorLatency.P99MS) }},
		} {
			head := []string{metric.title}
			for _, t := range tiers {
				head = append(head, tierLabel(t)+">")
			}
			t := md.NewTable(head...)
			for _, id := range rp.designsIn(tp) {
				row := []string{designShort(id)}
				for _, tier := range tiers {
					rr := raceOf(rp.runs[tp][id], "race", tier)
					if rr == nil {
						row = append(row, "—")
						continue
					}
					c := metric.cell(rr)
					if rr.Overbooked > 0 {
						c += " ❌"
					}
					row = append(row, c)
				}
				t.Row(row...)
			}
			t.Write(b)
		}
	}
}

func (rp *report) writePrecreation(b *strings.Builder) {
	fmt.Fprintf(b, "\n## Pre-created or created: the cost outside the race\n\n")
	fmt.Fprintf(b, "*Publish* creates a new event with its inventory in one transaction — one row, or one row\n")
	fmt.Fprintf(b, "per seat — and is shown as the median time to publish one event. *Book (spread)* is demand\n")
	fmt.Fprintf(b, "spread across the catalogue in proportion to seats left: little contention on any one event.\n")
	fmt.Fprintf(b, "*Cancel* refunds a sold ticket. Each operation ran on its own fresh load.\n\n")
	for _, tp := range rp.topos {
		head := []string{"Design"}
		for _, tier := range allTiers {
			head = append(head, "publish "+tierLabel(tier)+" (ms)>")
		}
		head = append(head, "book spread /s>", "cancel /s>", "load total>", "storage>")
		t := md.NewTable(head...)
		for _, id := range rp.designsIn(tp) {
			r := rp.runs[tp][id]
			row := []string{designShort(id)}
			for _, tier := range allTiers {
				w := writeOf(r, "publish", tier)
				if w == nil {
					row = append(row, "—")
					continue
				}
				row = append(row, md.MS(w.Latency.P50MS)+writeMark(w))
			}
			for _, op := range []string{"book", "cancel"} {
				w := writeOf(r, op, 0)
				if w == nil {
					row = append(row, "—")
					continue
				}
				row = append(row, md.Ops(w.OpsPerSec)+writeMark(w))
			}
			load, size := "—", "—"
			if r.Load != nil {
				load = fmt.Sprintf("%.1fs", r.Load.TotalMS/1000)
			}
			if r.Stats != nil && r.Stats.TotalBytes > 0 {
				size = md.Bytes(r.Stats.TotalBytes)
			}
			row = append(row, load, size)
			t.Row(row...)
		}
		fmt.Fprintf(b, "### %s\n\n", topologyLabel[tp])
		t.Write(b)
	}
	fmt.Fprintf(b, "Storage is PostgreSQL only: YugabyteDB keeps data in DocDB, where `pg_total_relation_size()`\n")
	fmt.Fprintf(b, "does not describe it.\n\n")
}

func writeMark(w *WriteResult) string {
	var m []string
	if w.Audit != nil && w.Audit.Violations() > 0 {
		m = append(m, "❌")
	}
	if w.Errors > 0 {
		m = append(m, fmt.Sprintf("(%d errors)", w.Errors))
	}
	if w.Exhausted {
		m = append(m, "(pool exhausted)")
	}
	if w.SafetyStopped {
		m = append(m, "(safety stop)")
	}
	if len(w.Trials) > 1 && w.SpreadPct > 20 {
		m = append(m, fmt.Sprintf("(±%.0f%%)", w.SpreadPct))
	}
	if len(m) == 0 {
		return ""
	}
	return " " + strings.Join(m, " ")
}

func (rp *report) writeReads(b *strings.Builder) {
	fmt.Fprintf(b, "\n## Reads (ops/s, higher is better)\n\n")
	fmt.Fprintf(b, "Each question benchmarked in isolation with a fresh key per execution. Availability is\n")
	fmt.Fprintf(b, "measured separately for every event size, because that is where counting designs pay.\n\n")
	for _, tp := range rp.topos {
		ids := rp.designsIn(tp)
		head := []string{"Question"}
		for _, id := range ids {
			head = append(head, designShort(id)+">")
		}
		t := md.NewTable(head...)
		for _, tier := range allTiers {
			row := []string{"Seats left — " + tierLabel(tier)}
			found := false
			for _, id := range ids {
				rr := readOf(rp.runs[tp][id], "q01_event_availability", tier)
				if rr == nil {
					row = append(row, "—")
					continue
				}
				found = true
				row = append(row, md.Ops(rr.OpsPerSec))
			}
			if found {
				t.Row(row...)
			}
		}
		for _, q := range []string{"q02_band_events", "q03_customer_tickets", "q04_ticket_by_id", "q05_band_tickets_sold"} {
			row := []string{queryLabel[q]}
			for _, id := range ids {
				rr := readOf(rp.runs[tp][id], q, 0)
				if rr == nil {
					row = append(row, "—")
					continue
				}
				row = append(row, md.Ops(rr.OpsPerSec))
			}
			t.Row(row...)
		}
		fmt.Fprintf(b, "### %s\n\n", topologyLabel[tp])
		t.Write(b)
	}
}

func (rp *report) writeHolds(b *strings.Builder) {
	any := false
	for _, tp := range rp.topos {
		for _, r := range rp.runs[tp] {
			if len(r.Holds) > 0 {
				any = true
			}
		}
	}
	if !any {
		return
	}
	fmt.Fprintf(b, "\n## Holds with expiry (H0, H1)\n\n")
	fmt.Fprintf(b, "Buyers hold a seat, spend time in the basket, abandon some holds, and pay for the rest; a\n")
	fmt.Fprintf(b, "sweeper releases expired holds. Some payments finish after the hold expired. *Late confirms\n")
	fmt.Fprintf(b, "rejected* are those the design refused (and would refund); a design that accepts them can sell a\n")
	fmt.Fprintf(b, "released seat twice.\n\n")
	for _, tp := range rp.topos {
		t := md.NewTable("Design", "Event size", "Events>", "Confirmed/s>", "Holds>", "Abandoned>", "Expired at pre-check>",
			"Late confirms rejected>", "Swept>", "Hold p99 (ms)>", "Confirm p99 (ms)>", "Overbooked", "Errors>")
		rows := 0
		for _, id := range rp.designsIn(tp) {
			r := rp.runs[tp][id]
			for _, h := range r.Holds {
				ob := "✅ none"
				if h.Overbooked > 0 {
					ob = fmt.Sprintf("❌ %d events (+%d seats)", h.Overbooked, h.OverbookedSeats)
				}
				if h.TimedOut > 0 {
					ob += fmt.Sprintf(", ⏱ %d timed out", h.TimedOut)
				}
				t.Row(designShort(id), tierLabel(h.Tier), fmt.Sprint(h.Events), md.Ops(h.ConfirmedPerSec), fmt.Sprint(h.Holds),
					fmt.Sprint(h.Abandoned), fmt.Sprint(h.ExpiredAtCheck), fmt.Sprint(h.LateRejected), fmt.Sprint(h.SweptReleased),
					md.MS(h.HoldLatency.P99MS), md.MS(h.ConfirmLatency.P99MS), ob, fmt.Sprint(h.Errors))
				rows++
			}
		}
		if rows > 0 {
			fmt.Fprintf(b, "### %s\n\n", topologyLabel[tp])
			t.Write(b)
		}
	}
}

// noiseFloor derives this run's error bar from designs whose reads are
// byte-identical by construction (C1, C2 and C3 share queries.sql): every read
// difference between them is measurement noise.
func (rp *report) noiseFloor() (worst, typical float64, n int) {
	var ratios []float64
	controls := []string{"c1_count_naive", "c2_count_serializable", "c3_count_lock_event"}
	for _, tp := range rp.topos {
		for i := 0; i < len(controls); i++ {
			for j := i + 1; j < len(controls); j++ {
				a, bb := rp.runs[tp][controls[i]], rp.runs[tp][controls[j]]
				if a == nil || bb == nil {
					continue
				}
				for _, ra := range a.Reads {
					rb := readOf(bb, ra.Query, ra.Tier)
					if rb == nil || ra.OpsPerSec <= 0 || rb.OpsPerSec <= 0 {
						continue
					}
					ratios = append(ratios, math.Max(ra.OpsPerSec/rb.OpsPerSec, rb.OpsPerSec/ra.OpsPerSec))
				}
			}
		}
	}
	if len(ratios) == 0 {
		return 1.3, 1.3, 0
	}
	sort.Float64s(ratios)
	return ratios[len(ratios)-1], ratios[len(ratios)/2], len(ratios)
}

func (rp *report) writePairs(b *strings.Builder) {
	worst, typical, n := rp.noiseFloor()
	fmt.Fprintf(b, "\n## Controlled pairs — one decision at a time\n\n")
	if n > 0 {
		fmt.Fprintf(b, "**This run's error bar.** C1, C2 and C3 read through byte-identical SQL on identical data, so\n")
		fmt.Fprintf(b, "every read difference between them is noise: %d comparisons, median disagreement **%.2fx**, worst\n", n, typical)
		fmt.Fprintf(b, "**%.2fx**. Read rows inside the worst disagreement are omitted below. **Race and write rows are\n", worst)
		fmt.Fprintf(b, "always shown**: they are the price of each decision, and a hidden price reads as no price.\n\n")
	} else {
		fmt.Fprintf(b, "The C1/C2/C3 read control is not present in this run, so read differences are filtered at a\n")
		fmt.Fprintf(b, "conventional %.2fx rather than a measured error bar.\n\n", worst)
	}
	if worst > 1.4 {
		fmt.Fprintf(b, "> ⚠️ A worst-case control disagreement of %.2fx means this run supports order-of-magnitude\n", worst)
		fmt.Fprintf(b, "> read claims only. Arguments about smaller read differences need a repeated-trials run.\n\n")
	}
	for _, p := range pairs {
		t := md.NewTable("Topology", "Measurement", designShort(p.a)+">", designShort(p.b)+">", "Change")
		for _, tp := range rp.topos {
			ra, rb := rp.runs[tp][p.a], rp.runs[tp][p.b]
			if ra == nil || rb == nil {
				continue
			}
			for _, mode := range []string{"race", "churn"} {
				for _, tier := range rp.raceTiers(mode) {
					xa, xb := raceOf(ra, mode, tier), raceOf(rb, mode, tier)
					if xa == nil || xb == nil {
						continue
					}
					t.Row(topologyLabel[tp], fmt.Sprintf("%s, %s — sold/s", mode, tierLabel(tier)), raceCell(xa), raceCell(xb), changeCell(xa, xb))
				}
			}
			for _, op := range []string{"book", "cancel"} {
				wa, wb := writeOf(ra, op, 0), writeOf(rb, op, 0)
				if wa != nil && wb != nil {
					t.Row(topologyLabel[tp], op+" /s", md.Ops(wa.OpsPerSec)+writeMark(wa), md.Ops(wb.OpsPerSec)+writeMark(wb), md.Ratio(wa.OpsPerSec, wb.OpsPerSec))
				}
			}
			for _, tier := range allTiers {
				wa, wb := writeOf(ra, "publish", tier), writeOf(rb, "publish", tier)
				if wa != nil && wb != nil {
					t.Row(topologyLabel[tp], "publish "+tierLabel(tier)+" p50 ms", md.MS(wa.Latency.P50MS), md.MS(wb.Latency.P50MS),
						md.Ratio(wb.Latency.P50MS, wa.Latency.P50MS))
				}
			}
			for _, qa := range ra.Reads {
				qb := readOf(rb, qa.Query, qa.Tier)
				if qb == nil || qa.OpsPerSec <= 0 || qb.OpsPerSec <= 0 {
					continue
				}
				if r := qb.OpsPerSec / qa.OpsPerSec; r < worst && r > 1/worst {
					continue
				}
				label := queryLabel[qa.Query]
				if qa.Tier > 0 {
					label += " — " + tierLabel(qa.Tier)
				}
				t.Row(topologyLabel[tp], label, md.Ops(qa.OpsPerSec), md.Ops(qb.OpsPerSec), md.Ratio(qa.OpsPerSec, qb.OpsPerSec))
			}
			for i := range ra.Holds {
				if i < len(rb.Holds) && ra.Holds[i].Tier == rb.Holds[i].Tier {
					ha, hb := ra.Holds[i], rb.Holds[i]
					t.Row(topologyLabel[tp], fmt.Sprintf("holds, %s — confirmed/s (overbooked events)", tierLabel(ha.Tier)),
						fmt.Sprintf("%s (%d)", md.Ops(ha.ConfirmedPerSec), ha.Overbooked), fmt.Sprintf("%s (%d)", md.Ops(hb.ConfirmedPerSec), hb.Overbooked),
						md.Ratio(ha.ConfirmedPerSec, hb.ConfirmedPerSec))
				}
			}
		}
		if t.Len() == 0 {
			continue
		}
		fmt.Fprintf(b, "### %s → %s\n\n*%s*\n\n", designShort(p.a), designShort(p.b), p.question)
		t.Write(b)
	}
}

// writeHazards reports how often the benchmark client itself was CPU-throttled
// by its quota during each phase. A phase where the client spent a large share of
// its CFS periods throttled has tail latencies the client may have caused.
func (rp *report) writeHazards(b *strings.Builder) {
	var phases []string
	seen := map[string]bool{}
	for _, tp := range rp.topos {
		for _, r := range rp.runs[tp] {
			for ph := range r.ClientCPU {
				if !seen[ph] {
					seen[ph] = true
					phases = append(phases, ph)
				}
			}
		}
	}
	if len(phases) == 0 {
		return
	}
	sort.Strings(phases)
	fmt.Fprintf(b, "\n## Measurement hazard: CPU-quota throttling\n\n")
	fmt.Fprintf(b, "Every container runs under a CFS quota (`--cpus`). When a container exhausts a 100 ms period's\n")
	fmt.Fprintf(b, "allowance it is paused until the next period, which puts tens of milliseconds into whatever\n")
	fmt.Fprintf(b, "request it was serving or issuing — latency that belongs to no design. Cells show *throttled\n")
	fmt.Fprintf(b, "periods / periods (throttled time)*: for the **database** container(s) over the whole cell\n")
	fmt.Fprintf(b, "(summed across nodes), and for the **client** per phase.\n\n")
	for _, tp := range rp.topos {
		head := []string{"Design", "database, whole cell>"}
		for _, ph := range phases {
			head = append(head, "client "+ph+">")
		}
		t := md.NewTable(head...)
		for _, id := range rp.designsIn(tp) {
			r := rp.runs[tp][id]
			row := []string{designShort(id)}
			if c, ok := dbCPU(filepath.Join(rp.dir, tp, "logs", id+".dbcpu.txt")); ok && c.NrPeriods > 0 {
				row = append(row, fmt.Sprintf("%d/%d (%.0fs)", c.NrThrottled, c.NrPeriods, float64(c.ThrottledUsec)/1e6))
			} else {
				row = append(row, "—")
			}
			for _, ph := range phases {
				c, ok := r.ClientCPU[ph]
				if !ok || c.NrPeriods == 0 {
					row = append(row, "—")
					continue
				}
				row = append(row, fmt.Sprintf("%d/%d (%.1fs)", c.NrThrottled, c.NrPeriods, float64(c.ThrottledUsec)/1e6))
			}
			t.Row(row...)
		}
		fmt.Fprintf(b, "### %s\n\n", topologyLabel[tp])
		t.Write(b)
	}
}

// dbCPU parses the runner's before/after cpu.stat captures and returns the
// counters accumulated over the cell, summed over the database containers.
func dbCPU(path string) (cgroup.CPUStat, bool) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return cgroup.CPUStat{}, false
	}
	type key struct{ name, when string }
	stats := map[key]*cgroup.CPUStat{}
	var cur *cgroup.CPUStat
	for _, line := range strings.Split(string(raw), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "[") {
			name, rest, _ := strings.Cut(strings.TrimPrefix(line, "["), "] ")
			when, _, _ := strings.Cut(rest, " ")
			// A cell's first capture of each kind is the one that counts.
			k := key{name, when}
			if stats[k] != nil {
				cur = nil
				continue
			}
			cur = &cgroup.CPUStat{}
			stats[k] = cur
			continue
		}
		if cur == nil {
			continue
		}
		f := strings.Fields(line)
		if len(f) != 2 {
			continue
		}
		var v int64
		fmt.Sscan(f[1], &v)
		switch f[0] {
		case "usage_usec":
			cur.UsageUsec = v
		case "nr_periods":
			cur.NrPeriods = v
		case "nr_throttled":
			cur.NrThrottled = v
		case "throttled_usec":
			cur.ThrottledUsec = v
		}
	}
	var total cgroup.CPUStat
	found := false
	for k, before := range stats {
		if k.when != "before" {
			continue
		}
		after := stats[key{k.name, "after"}]
		if after == nil {
			continue
		}
		d := after.Sub(*before)
		total.UsageUsec += d.UsageUsec
		total.NrPeriods += d.NrPeriods
		total.NrThrottled += d.NrThrottled
		total.ThrottledUsec += d.ThrottledUsec
		found = true
	}
	return total, found
}

func changeCell(a, b *RaceResult) string {
	c := md.Ratio(a.SoldPerSec, b.SoldPerSec)
	if a.Overbooked > 0 || b.Overbooked > 0 {
		return c + " — ❌ not comparable: an overbooking design has no valid speed"
	}
	if a.TimedOut > 0 || b.TimedOut > 0 {
		return c + " — ⏱ a timed-out race measured only part of the sale"
	}
	return c
}

func firstSentence(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.IndexAny(s, ";\n"); i > 0 {
		s = s[:i]
	}
	if len(s) > 100 {
		s = s[:97] + "..."
	}
	return s
}
