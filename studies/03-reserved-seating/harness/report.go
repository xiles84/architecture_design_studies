package main

// Adapted from studies/02-ticket-booking/harness/report.go at 3c0c3aa.

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
// One rule is applied throughout: a number produced while a design violated an
// invariant is never shown as a plain number. It is marked ❌, because a design
// that sells stolen seats quickly has not sold quickly.
// ---------------------------------------------------------------------------

var topologyOrder = []string{"pg-single", "yb-single", "yb-cluster3"}

var topologyLabel = map[string]string{
	"pg-single":   "PostgreSQL, 1 node",
	"yb-single":   "YugabyteDB, 1 node (RF=1)",
	"yb-cluster3": "YugabyteDB, 3 nodes (RF=3)",
}

var queryLabel = map[string]string{
	"q01_section_map":      "Seat map of a section",
	"q02_event_sections":   "Seats available per section",
	"q03_event_available":  "Seats available for an event",
	"q04_hold_seats":       "A hold's seats (basket page)",
	"q05_customer_tickets": "A customer's tickets",
	"q06_ticket_by_id":     "Ticket by id (control)",
	// Operational reports (REPORTS.md v2, AM-02). r01/r03/r05's window
	// regime lives in the Query name suffix ("@trailing"/"@historical"),
	// which writeReads' generic column-discovery already treats as a
	// distinct column, exactly as it does q01's tier suffix.
	"r01_holds_expiring_soon":                       "Holds expiring soon",
	"r02_seat_status_lookup":                        "Seat status lookup",
	"r03_section_sales_window@trailing":             "Section sales (trailing week)",
	"r03_section_sales_window@historical":           "Section sales (historical week)",
	"r04_event_recent_confirmations":                "Recent confirmations",
	"r05_customers_last_purchase_window@trailing":   "Customers' last purchase (trailing)",
	"r05_customers_last_purchase_window@historical": "Customers' last purchase (historical)",
}

type pair struct{ a, b, question string }

var pairs = []pair{
	{"s0_check_then_hold_rc", "s4_check_then_hold_serializable", "Identical SQL: what does SERIALIZABLE cost, and does it prevent theft?"},
	{"s4_check_then_hold_serializable", "s1_conditional_update", "Serializable read-then-write, or one conditional update?"},
	{"s1_conditional_update", "s2_lock_then_update", "Optimistic conditional update, or lock the seats first?"},
	{"s2_lock_then_update", "s3_lock_nowait", "Wait for a seat someone is taking, or fail fast and choose again?"},
	{"s1_conditional_update", "e1_sweeper_expiry", "Expiry judged at use time, or made to take effect by a sweeper?"},
	{"s1_conditional_update", "e2_cart_expiry", "Expiry copied onto every seat row, or stored once on the cart?"},
	{"e0_app_clock_expiry", "s1_conditional_update", "The application's clock, or the database's?"},
	{"k0_naive_confirm", "s1_conditional_update", "What does checking the hold at confirmation cost, and what does it prevent?"},
	{"s1_conditional_update", "k1_payment_window", "What does a guaranteed payment window cost, and what does it remove?"},
	{"s1_conditional_update", "l1_claim_rows", "Pre-created per-event seat rows, or claims created on hold?"},
	{"s1_conditional_update", "l2_section_document", "Seat rows, or an embedded section document?"},
	{"s1_conditional_update", "l3_section_sharded", "(YugabyteDB) all of an event's seats on one tablet, or spread by section?"},
	{"s1_conditional_update", "s1r_confirm_retry", "Does retrying a short confirmation once remove transient refusals, and what does it cost?"},
}

// noiseDesigns read q05/q06 through byte-identical SQL over identical tables.
var noiseDesigns = []string{"s0_check_then_hold_rc", "s1_conditional_update", "s2_lock_then_update", "s3_lock_nowait",
	"s4_check_then_hold_serializable", "e0_app_clock_expiry", "e1_sweeper_expiry", "k0_naive_confirm", "k1_payment_window", "l2_section_document", "s1r_confirm_retry"}

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

func isControl(id string) bool {
	d, err := designByID(id)
	return err == nil && d.NegativeControl != ""
}

// raceOf combines every trial of one tier: throughput is the median trial's;
// latencies are that trial's; violations are SUMMED over trials.
func raceOf(r *Run, tier int) *RaceResult {
	if r == nil {
		return nil
	}
	var trials []*RaceResult
	for i := range r.Races {
		if r.Races[i].Tier == tier {
			trials = append(trials, &r.Races[i])
		}
	}
	switch len(trials) {
	case 0:
		return nil
	case 1:
		return trials[0]
	}
	sort.Slice(trials, func(i, j int) bool { return trials[i].SeatsPerSec < trials[j].SeatsPerSec })
	agg := *trials[len(trials)/2]
	agg.Audit = &Audit{}
	agg.Events, agg.TimedOut, agg.UnderSold, agg.UnderSoldSeats, agg.Leaked = 0, 0, 0, 0, 0
	agg.DeferredHolds, agg.DeferredConfirmed, agg.DeferredRefused, agg.Errors = 0, 0, 0, 0
	agg.ConfirmRetries, agg.ConfirmRetrySuccesses = 0, 0
	var sps []float64
	for _, t := range trials {
		agg.Events += t.Events
		agg.TimedOut += t.TimedOut
		agg.UnderSold += t.UnderSold
		agg.UnderSoldSeats += t.UnderSoldSeats
		agg.Leaked += t.Leaked
		agg.DeferredHolds += t.DeferredHolds
		agg.DeferredConfirmed += t.DeferredConfirmed
		agg.DeferredRefused += t.DeferredRefused
		agg.ConfirmRetries += t.ConfirmRetries
		agg.ConfirmRetrySuccesses += t.ConfirmRetrySuccesses
		agg.Errors += t.Errors
		if t.Audit != nil {
			mergeAudit(agg.Audit, t.Audit)
		}
		sps = append(sps, t.SeatsPerSec)
	}
	agg.SeatsPerSec = measure.Median(sps)
	agg.TrialsSeatsPerSec = sps
	agg.SpreadPct = measure.SpreadPct(sps)
	return &agg
}

func mergeAudit(into, a *Audit) {
	into.Ledger.add(a.Ledger)
	into.DuplicateSeats += a.DuplicateSeats
	into.InventoryDrift += a.InventoryDrift
	into.InvalidSeats += a.InvalidSeats
	into.PartialHolds += a.PartialHolds
	into.TicketMismatches += a.TicketMismatches
	into.Events += a.Events
	into.Tickets += a.Tickets
}

// lifecycleOf sums every trial of one tier (throughput: the median trial's).
func lifecycleOf(r *Run, tier int) *LifecycleResult {
	if r == nil {
		return nil
	}
	var trials []*LifecycleResult
	for i := range r.Lifecycle {
		if r.Lifecycle[i].Tier == tier {
			trials = append(trials, &r.Lifecycle[i])
		}
	}
	switch len(trials) {
	case 0:
		return nil
	case 1:
		return trials[0]
	}
	sort.Slice(trials, func(i, j int) bool { return trials[i].ConfirmedSeatsPerSec < trials[j].ConfirmedSeatsPerSec })
	agg := *trials[len(trials)/2]
	agg.Audit = &Audit{}
	agg.RejectedLate, agg.RejectedBoundary, agg.RejectedEarly, agg.RejectedEarlyTransient = 0, 0, 0, 0
	agg.ConfirmRetries, agg.ConfirmRetrySuccesses = 0, 0
	agg.OutageProbed, agg.OutageUnavailable, agg.Leaked, agg.Errors = 0, 0, 0, 0
	agg.MonitorRetries, agg.MonitorTimeouts = 0, 0
	for _, t := range trials {
		agg.RejectedLate += t.RejectedLate
		agg.RejectedBoundary += t.RejectedBoundary
		agg.RejectedEarly += t.RejectedEarly
		agg.RejectedEarlyTransient += t.RejectedEarlyTransient
		agg.ConfirmRetries += t.ConfirmRetries
		agg.MonitorRetries += t.MonitorRetries
		agg.MonitorTimeouts += t.MonitorTimeouts
		agg.ConfirmRetrySuccesses += t.ConfirmRetrySuccesses
		agg.OutageProbed += t.OutageProbed
		agg.OutageUnavailable += t.OutageUnavailable
		agg.Leaked += t.Leaked
		agg.Errors += t.Errors
		if t.Audit != nil {
			mergeAudit(agg.Audit, t.Audit)
		}
	}
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

func (rp *report) tiers(kind string) []int {
	seen := map[int]bool{}
	var out []int
	for _, t := range rp.topos {
		for _, r := range rp.runs[t] {
			if kind == "race" {
				for _, x := range r.Races {
					if !seen[x.Tier] {
						seen[x.Tier] = true
						out = append(out, x.Tier)
					}
				}
			} else {
				for _, x := range r.Lifecycle {
					if !seen[x.Tier] {
						seen[x.Tier] = true
						out = append(out, x.Tier)
					}
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

// violationMarks lists every invariant an audit found broken.
func violationMarks(id string, a *Audit) []string {
	if a == nil {
		return nil
	}
	var m []string
	add := func(n int64, what string) {
		if n > 0 {
			m = append(m, fmt.Sprintf("%d %s", n, what))
		}
	}
	add(a.Ledger.Thefts, "thefts")
	add(a.Ledger.SoldSeatGrants, "grants on sold seats")
	add(a.Ledger.DoubleSales, "double sales")
	add(a.Ledger.SalesWithoutHold, "sales without the hold")
	add(a.Ledger.LateSales, "late sales")
	if a.Ledger.RejectedEarly > 0 {
		m = append(m, fmt.Sprintf("%d early rejections (transient: %s)",
			a.Ledger.RejectedEarly, transientCell(id, a, a.Ledger.RejectedEarlyTransient)))
	}
	add(a.DuplicateSeats, "seats with 2+ tickets")
	add(a.InventoryDrift, "inventory drift")
	add(a.InvalidSeats, "invalid seats")
	add(a.PartialHolds, "partial/unknown holds")
	add(a.TicketMismatches, "ticket mismatches")
	return m
}

// earlyCell is "early (transient)" from an audit, or — without one.
func earlyCell(id string, a *Audit) string {
	if a == nil {
		return "—"
	}
	return fmt.Sprintf("%d (%s)", a.Ledger.RejectedEarly, transientCell(id, a, a.Ledger.RejectedEarlyTransient))
}

// transientCell is the transient part of early rejections, or n/a for designs whose
// refusal is not a guarded statement (K0 has no guard; L2 checks in the application).
// transientCell reports the transient part of a design's early rejections, or says the
// class was not recorded. A class exists only where a diagnostic looked at the refusal
// (AM-03.3):
//
//   - K0 has no guarded statement to re-issue and never records one;
//   - L2 records one only from AM-03.2 onwards, so its refusals in earlier runs are
//     unclassified — visible as early rejections with nothing classified;
//   - every other design has re-issued its refusing statement since AM-01.2, so its
//     count is real, including a genuine zero.
//
// Keying only on the post-AM-03 "classified" counter would relabel every earlier run's
// recorded class as unknown, which is why the design decides.
func transientCell(id string, a *Audit, n int64) string {
	if a == nil {
		return "—"
	}
	if id == "k0_naive_confirm" {
		return "n/a"
	}
	if id == "l2_section_document" && a.Ledger.RejectedEarly > 0 && a.Ledger.RejectedEarlyClassified == 0 {
		return "not recorded"
	}
	return fmt.Sprint(n)
}

func raceCell(id string, rr *RaceResult) string {
	if rr == nil {
		return "—"
	}
	v := md.Ops(rr.SeatsPerSec)
	if len(rr.TrialsSeatsPerSec) > 1 {
		v += fmt.Sprintf(" (%d trials, spread %.0f%%)", len(rr.TrialsSeatsPerSec), rr.SpreadPct)
	}
	var marks []string
	if vm := violationMarks(id, rr.Audit); len(vm) > 0 {
		marks = append(marks, "❌ "+strings.Join(vm, ", "))
	}
	if rr.DeferredRefused > 0 {
		marks = append(marks, fmt.Sprintf("❌ %d of %d deferred confirmations refused", rr.DeferredRefused, rr.DeferredHolds))
	}
	if rr.UnderSold > 0 {
		marks = append(marks, fmt.Sprintf("⚠️ undersold %d/%d (−%d)", rr.UnderSold, rr.Events, rr.UnderSoldSeats))
	}
	if rr.Leaked > 0 {
		marks = append(marks, fmt.Sprintf("⚠️ leaked %d", rr.Leaked))
	}
	if rr.TimedOut > 0 {
		marks = append(marks, fmt.Sprintf("⏱ %d timed out at %.0f%% sold", rr.TimedOut, rr.TimedOutSoldPct))
	}
	if rr.EventsPlanned > rr.Events {
		marks = append(marks, fmt.Sprintf("⌛ %d of %d events raced within the tier budget", rr.Events, rr.EventsPlanned))
	}
	if rr.Errors > 0 {
		marks = append(marks, fmt.Sprintf("%d errors", rr.Errors))
	}
	if len(marks) == 0 {
		return v
	}
	if rr.Audit != nil && rr.Audit.Violations() > 0 {
		return fmt.Sprintf("~~%s~~ %s", v, strings.Join(marks, "; "))
	}
	return v + " " + strings.Join(marks, "; ")
}

func WriteReport(dir, outPath string) error {
	rp, err := loadReport(dir)
	if err != nil {
		return err
	}
	ref := rp.any()
	var b strings.Builder
	fmt.Fprintf(&b, "# Study 03 — results: reserved seating (venue → event → seat)\n\n")
	fmt.Fprintf(&b, "Generated from `results/%s`. Every number comes from exactly one JSON file in that\n", filepath.Base(dir))
	fmt.Fprintf(&b, "directory. **This file contains measurements and no interpretation**; conclusions live in\n")
	fmt.Fprintf(&b, "signed analyses, indexed at the end.\n\n")

	rp.writeTLDR(&b)
	rp.writeSetup(&b, ref)
	rp.writeGate(&b)
	rp.writeControls(&b)
	rp.writeGuarantee(&b)
	rp.writeRace(&b)
	rp.writeReportCoverage(&b)
	rp.writeReads(&b)
	rp.writeWrites(&b)
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

// controlFired reports whether a negative control broke an invariant in an
// experiment ("race" or "lifecycle") on a topology; ok=false means not measured.
func (rp *report) controlFired(tp, id, experiment string) (fired, ok bool, detail string) {
	r := rp.runs[tp][id]
	if r == nil {
		return false, false, ""
	}
	var agg Audit
	switch experiment {
	case "race":
		for _, x := range r.Races {
			ok = true
			if x.Audit != nil {
				mergeAudit(&agg, x.Audit)
			}
			if x.DeferredRefused > 0 {
				detail = fmt.Sprintf("%d deferred confirmations refused; ", x.DeferredRefused)
			}
		}
	case "lifecycle":
		for _, x := range r.Lifecycle {
			ok = true
			if x.Audit != nil {
				mergeAudit(&agg, x.Audit)
			}
		}
	}
	vm := violationMarks(id, &agg)
	return len(vm) > 0, ok, detail + strings.Join(vm, ", ")
}

// outageFired: E1's leak detector during the sweeper outage.
func (rp *report) outageFired(tp string) (fired, ok bool, detail string) {
	r := rp.runs[tp]["e1_sweeper_expiry"]
	if r == nil {
		return false, false, ""
	}
	var probed, unavailable int64
	for _, x := range r.Lifecycle {
		if x.OutageEvents > 0 {
			ok = true
			probed += x.OutageProbed
			unavailable += x.OutageUnavailable
		}
	}
	return unavailable > 0, ok, fmt.Sprintf("%d of %d probed seats unavailable after expiry", unavailable, probed)
}

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
	for _, fc := range rp.manifestFailures() {
		tp, id, _ := strings.Cut(fc, "/")
		if rp.runs[tp][id] == nil {
			failed = append(failed, fmt.Sprintf("%s / %s (no result file; see the manifest and logs)", topologyLabel[tp], designShort(id)))
		}
	}
	fmt.Fprintf(b, "- **Cells:** %d across %d topologies; %d failed%s.\n", cells, len(rp.topos), len(failed), listSuffix(failed))

	var fired, silent []string
	for _, tp := range rp.topos {
		for _, d := range designs {
			for _, exp := range d.ControlFires {
				f, ok, _ := rp.controlFired(tp, d.ID, exp)
				if !ok {
					continue
				}
				label := fmt.Sprintf("%s in the %s on %s", d.Short, exp, topologyLabel[tp])
				if f {
					fired = append(fired, label)
				} else {
					silent = append(silent, label)
				}
			}
		}
		if f, ok, _ := rp.outageFired(tp); ok {
			label := "E1 during the sweeper outage on " + topologyLabel[tp]
			if f {
				fired = append(fired, label)
			} else {
				silent = append(silent, label)
			}
		}
	}
	fmt.Fprintf(b, "- **Negative controls:** %d of %d fired as designed to.", len(fired), len(fired)+len(silent))
	if len(silent) > 0 {
		fmt.Fprintf(b, " ⚠️ **Did not fire:** %s — the absence of violations in the other designs of that experiment and topology is not evidence of correctness.", strings.Join(silent, "; "))
	}
	fmt.Fprintf(b, "\n")

	var broken []string
	for _, tp := range rp.topos {
		for _, id := range rp.designsIn(tp) {
			if isControl(id) {
				continue
			}
			var agg Audit
			r := rp.runs[tp][id]
			for _, x := range r.Races {
				if x.Audit != nil {
					mergeAudit(&agg, x.Audit)
				}
			}
			for _, x := range r.Lifecycle {
				if x.Audit != nil {
					mergeAudit(&agg, x.Audit)
				}
			}
			for _, x := range r.Writes {
				if x.Audit != nil && x.Op != "publish" {
					mergeAudit(&agg, x.Audit)
				}
			}
			if vm := violationMarks(id, &agg); len(vm) > 0 {
				broken = append(broken, fmt.Sprintf("%s on %s (%s)", designShort(id), topologyLabel[tp], strings.Join(vm, ", ")))
			}
		}
	}
	if len(broken) == 0 {
		fmt.Fprintf(b, "- **Invariant violations in designs meant to be correct:** none observed.\n")
	} else {
		fmt.Fprintf(b, "- ❌ **Invariant violations in designs meant to be correct:** %s.\n", strings.Join(broken, "; "))
	}

	for _, tp := range rp.topos {
		for _, id := range []string{"s1_conditional_update", "s1r_confirm_retry", "k1_payment_window"} {
			r := rp.runs[tp][id]
			if r == nil || len(r.Lifecycle) == 0 {
				continue
			}
			var late, boundary, early, transient, confirmed, retries, retrySold int64
			for _, x := range r.Lifecycle {
				late += x.RejectedLate
				boundary += x.RejectedBoundary
				early += x.RejectedEarly
				transient += x.RejectedEarlyTransient
				confirmed += x.ConfirmedHolds
				retries += x.ConfirmRetries
				retrySold += x.ConfirmRetrySuccesses
			}
			extra := ""
			if id == "s1r_confirm_retry" {
				extra = fmt.Sprintf("; confirmations retried %d, of which sold %d", retries, retrySold)
			}
			fmt.Fprintf(b, "- **Refused confirmations, %s, %s** (lifecycle, all tiers): %d late, %d boundary, %d early (%d transient); %d confirmed%s.\n",
				designShort(id), topologyLabel[tp], late, boundary, early, transient, confirmed, extra)
		}
		// AM-01: early rejections in the race, correct designs only, split by class.
		var raceEarly []string
		for _, id := range rp.designsIn(tp) {
			if d, err := designByID(id); err != nil || d.NegativeControl != "" {
				continue
			}
			var early, transient, retries, retrySold int64
			for _, x := range rp.runs[tp][id].Races {
				if x.Audit != nil {
					early += x.Audit.Ledger.RejectedEarly
					transient += x.Audit.Ledger.RejectedEarlyTransient
				}
				retries += x.ConfirmRetries
				retrySold += x.ConfirmRetrySuccesses
			}
			if early > 0 {
				raceEarly = append(raceEarly, fmt.Sprintf("%s %d (%d transient)", designShort(id), early, transient))
			}
			if id == "s1r_confirm_retry" && len(rp.runs[tp][id].Races) > 0 {
				fmt.Fprintf(b, "- **S1r confirmation retries, %s** (race, all tiers and trials): %d retried, %d of them sold.\n",
					topologyLabel[tp], retries, retrySold)
			}
		}
		if len(raceEarly) > 0 {
			fmt.Fprintf(b, "- **Early rejections in the race, correct designs, %s:** %s.\n", topologyLabel[tp], strings.Join(raceEarly, ", "))
		}
		if f, ok, detail := rp.outageFired(tp); ok {
			_ = f
			var lazy []string
			for _, id := range rp.designsIn(tp) {
				if id == "e1_sweeper_expiry" {
					continue
				}
				var u int64
				for _, x := range rp.runs[tp][id].Lifecycle {
					u += x.OutageUnavailable
				}
				if u > 0 {
					lazy = append(lazy, fmt.Sprintf("%s %d", designShort(id), u))
				}
			}
			others := "0 in every other design"
			if len(lazy) > 0 {
				others = "other designs: " + strings.Join(lazy, ", ")
			}
			fmt.Fprintf(b, "- **Sweeper outage, %s:** E1 %s; %s.\n", topologyLabel[tp], detail, others)
		}
	}

	for _, tp := range rp.topos {
		fmt.Fprintf(b, "- **Race, %s** (seats held/s; designs meant to be correct, without violations):\n", topologyLabel[tp])
		for _, tier := range rp.tiers("race") {
			type kv struct {
				id string
				v  float64
			}
			var xs []kv
			timedOut := 0
			for _, id := range rp.designsIn(tp) {
				if isControl(id) {
					continue
				}
				rr := raceOf(rp.runs[tp][id], tier)
				if rr == nil || rr.SeatsPerSec <= 0 || rr.Audit.Violations() > 0 || rr.DeferredRefused > 0 {
					continue
				}
				if rr.TimedOut > 0 {
					timedOut++
				}
				xs = append(xs, kv{id, rr.SeatsPerSec})
			}
			if len(xs) == 0 {
				continue
			}
			sort.Slice(xs, func(i, j int) bool { return xs[i].v > xs[j].v })
			note := ""
			if timedOut > 0 {
				note = fmt.Sprintf(" (%d of %d timed out)", timedOut, len(xs))
			}
			fmt.Fprintf(b, "  - %s: highest %s %s/s, lowest %s %s/s%s\n", tierLabel(tier),
				designShort(xs[0].id), md.Ops(xs[0].v), designShort(xs[len(xs)-1].id), md.Ops(xs[len(xs)-1].v), note)
		}
	}
	for _, tp := range rp.topos {
		var with, without []float64
		big := 0
		for _, id := range rp.designsIn(tp) {
			for _, w := range rp.runs[tp][id].Writes {
				if w.Op == "publish" && w.Tier > big {
					big = w.Tier
				}
			}
		}
		for _, id := range rp.designsIn(tp) {
			w := writeOf(rp.runs[tp][id], "publish", big)
			if w == nil || w.Latency.P50MS <= 0 {
				continue
			}
			if id == "l1_claim_rows" {
				without = append(without, w.Latency.P50MS)
			} else {
				with = append(with, w.Latency.P50MS)
			}
		}
		if len(with) > 0 && len(without) > 0 {
			sort.Float64s(with)
			fmt.Fprintf(b, "- **Publishing a %s event, %s** (p50): %s–%s ms with per-event inventory rows or documents, %s ms without (L1).\n",
				tierLabel(big), topologyLabel[tp], md.MS(with[0]), md.MS(with[len(with)-1]), md.MS(without[0]))
		}
	}
	fmt.Fprintf(b, "\n")
}

func (rp *report) manifestFailures() []string {
	raw, err := os.ReadFile(filepath.Join(rp.dir, "manifest.yaml"))
	if err != nil {
		return nil
	}
	var out []string
	in := false
	for _, line := range strings.Split(string(raw), "\n") {
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
	t.Row("Scale", fmt.Sprintf("`%s` — %v venues, %v events, %v tickets sold, %v live and %v expired holds at load",
		ref.Scale, ref.Dataset["venues"], ref.Dataset["events"], ref.Dataset["tickets_sold_at_load"],
		ref.Dataset["live_holds_at_load"], ref.Dataset["expired_holds_at_load"]))
	t.Row("Events (kind_seats: count)", eventsLine(ref.Dataset["events_by_kind_tier"]))
	t.Row("Seed", fmt.Sprintf("%v — identical data in every cell", ref.Dataset["seed"]))
	o := ref.Options
	t.Row("Reads / isolated writes", fmt.Sprintf("%v workers, %v measured + %v warmup, %v trial(s)", o["conns"], o["duration"], o["warmup"], o["trials"]))
	t.Row("Race", fmt.Sprintf("%v buyers per event, real hold TTL %v, %v%% of holds deferred until the crowd has gone, timeout %v, tier budget %v, give up after %v conflicts",
		o["race_buyers"], o["hold_ttl"], o["race_defer_pct"], o["race_timeout"], o["race_tier_budget"], o["max_conflicts"]))
	for _, tp := range rp.topos {
		for _, id := range rp.designsIn(tp) {
			r := rp.runs[tp][id]
			if len(r.Lifecycle) == 0 {
				continue
			}
			o := r.Options
			t.Row("Lifecycle — "+topologyLabel[tp], fmt.Sprintf("one human minute = %v; TTL %v, payment window (K1) %v, guard G %v, sweeper every %v, outage %v→%v in the first event of each tier, %v%% abandoned (%v%% of those explicitly), E0 skew %v; tiers %v, %v events per tier, %v buyers",
				o["human_minute"], o["lc_ttl"], o["lc_payment_window"], o["lc_guard"], o["lc_sweep_every"], o["lc_outage_from"],
				o["lc_outage_to"], o["lc_abandon_pct"], o["lc_abandon_explicit_pct"], o["lc_e0_skew"], o["lifecycle_tiers"],
				o["lifecycle_events_per_tier"], o["lifecycle_buyers"]))
			break
		}
	}
	t.Write(b)

	fmt.Fprintf(b, "**Engines**\n\n")
	et := md.NewTable("Topology", "Version", "READ COMMITTED runs as", "Client connections to", "DB clock − client clock (ms)>")
	for _, tp := range rp.topos {
		for _, id := range rp.designsIn(tp) {
			r := rp.runs[tp][id]
			eff := r.EngineInfo.EffectiveIsolation
			if eff == "" {
				eff = "read committed"
			}
			et.Row(topologyLabel[tp], "`"+firstLine(r.EngineInfo.Version)+"`", eff,
				fmt.Sprintf("%d node(s)", r.ConnectionNodes), fmt.Sprintf("%.2f", r.ClientDBClockOffsetMS))
			break
		}
	}
	et.Write(b)
	fmt.Fprintf(b, "> Limits of this environment that bound every number below: the client and the databases\n")
	fmt.Fprintf(b, "> share one laptop's eight cores; the \"3-node\" cluster has no real network between nodes; the\n")
	fmt.Fprintf(b, "> race and the lifecycle are closed loops, so tails are a floor; the lifecycle runs on\n")
	fmt.Fprintf(b, "> compressed time; E0's clock skew is injected.\n\n")
}

func (rp *report) writeGate(b *strings.Builder) {
	fmt.Fprintf(b, "## Correctness gate\n\n")
	fmt.Fprintf(b, "Six read questions checked against Go-computed truth on a load with sold seats, live holds and\n")
	fmt.Fprintf(b, "expired holds, plus an audit of the loaded state. A failure aborts the cell.\n\n")
	head := []string{"Design"}
	for _, tp := range rp.topos {
		head = append(head, topologyLabel[tp])
	}
	t := md.NewTable(head...)
	for _, id := range rp.ids {
		row := []string{designShort(id)}
		for _, tp := range rp.topos {
			r := rp.runs[tp][id]
			switch {
			case r == nil:
				row = append(row, "n/a")
			case r.Verify == nil:
				row = append(row, "not run: "+firstSentence(r.Error))
			case r.Verify.Failed > 0 || r.LoadAudit.Violations() > 0:
				row = append(row, fmt.Sprintf("❌ %d/%d; audit: %s", r.Verify.Passed, r.Verify.Passed+r.Verify.Failed, r.LoadAudit))
			default:
				cell := fmt.Sprintf("✅ %d/%d", r.Verify.Passed, r.Verify.Passed)
				if r.PreSweep > 0 {
					cell += fmt.Sprintf(" (after its sweeper released %d expired seats)", r.PreSweep)
				}
				row = append(row, cell)
			}
		}
		t.Row(row...)
	}
	t.Write(b)
}

func (rp *report) writeControls(b *strings.Builder) {
	fmt.Fprintf(b, "## Negative controls\n\n")
	fmt.Fprintf(b, "Deliberately wrong designs (and E1 with its sweeper stopped). Each must be seen to break the\n")
	fmt.Fprintf(b, "invariant it exists to test; otherwise the absence of violations elsewhere in that experiment\n")
	fmt.Fprintf(b, "is not evidence of correctness.\n\n")
	t := md.NewTable("Control", "Experiment", "Topology", "Fired?", "What was found")
	for _, d := range designs {
		for _, exp := range d.ControlFires {
			for _, tp := range rp.topos {
				f, ok, detail := rp.controlFired(tp, d.ID, exp)
				if !ok {
					continue
				}
				mark := "✅ fired"
				if !f {
					mark = "⚠️ **did not fire**"
				}
				t.Row(d.Short, exp, topologyLabel[tp], mark, detail)
			}
		}
	}
	for _, tp := range rp.topos {
		if f, ok, detail := rp.outageFired(tp); ok {
			mark := "✅ fired"
			if !f {
				mark = "⚠️ **did not fire**"
			}
			t.Row("E1 sweeper stopped", "lifecycle outage", topologyLabel[tp], mark, detail)
		}
	}
	if t.Len() == 0 {
		fmt.Fprintf(b, "No control was measured in this run.\n\n")
		return
	}
	t.Write(b)
}

func (rp *report) writeGuarantee(b *strings.Builder) {
	tiers := rp.tiers("lifecycle")
	if len(tiers) == 0 {
		return
	}
	fmt.Fprintf(b, "## The hold guarantee — lifecycle\n\n")
	fmt.Fprintf(b, "Compressed time. *Rejected* confirmations are classified by how long before the hold's expiry\n")
	fmt.Fprintf(b, "the confirming transaction started: **late** (after expiry), **boundary** (less than G before),\n")
	fmt.Fprintf(b, "**early** (G or more before — a violation). Violations are counted by the ledger and the audit.\n")
	fmt.Fprintf(b, "Release lag is in human minutes.\n\n")
	for _, tp := range rp.topos {
		t := md.NewTable("Design", "Tier", "Confirmed seats/s>", "Holds>", "Abandoned>", "Expired at check>", "Refused late / boundary / early (transient)>",
			"Outage: unavailable / probed>", "Leaked>", "Release lag p50/p99 (min)>", "Idle held seat-min>", "Monitor retries / events ended early>", "Violations")
		for _, id := range rp.designsIn(tp) {
			for _, tier := range tiers {
				x := lifecycleOf(rp.runs[tp][id], tier)
				if x == nil {
					continue
				}
				viol := "none"
				if vm := violationMarks(id, x.Audit); len(vm) > 0 {
					viol = "❌ " + strings.Join(vm, ", ")
				}
				rate := md.Ops(x.ConfirmedSeatsPerSec)
				if x.Audit.Violations() > 0 {
					rate = "~~" + rate + "~~"
				}
				outage := "—"
				if x.OutageEvents > 0 {
					outage = fmt.Sprintf("%d / %d", x.OutageUnavailable, x.OutageProbed)
				}
				expCheck := fmt.Sprint(x.ExpiredAtCheck)
				if x.ExpiredAtCheckout > 0 {
					expCheck += fmt.Sprintf(" (+%d at checkout)", x.ExpiredAtCheckout)
				}
				t.Row(designShort(id), tierLabel(tier), rate, fmt.Sprint(x.HoldsGranted),
					fmt.Sprintf("%d+%d", x.AbandonedSilent, x.AbandonedExplicit), expCheck,
					fmt.Sprintf("%d / %d / %d (%s)", x.RejectedLate, x.RejectedBoundary, x.RejectedEarly, transientCell(id, x.Audit, x.RejectedEarlyTransient)), outage,
					fmt.Sprint(x.Leaked), fmt.Sprintf("%.1f / %.1f", x.ReleaseLagHumanMin.P50MS, x.ReleaseLagHumanMin.P99MS),
					fmt.Sprintf("%.0f", x.IdleHeldSeatHumanMin),
					fmt.Sprintf("%d / %d", x.MonitorRetries, x.MonitorTimeouts), viol)
			}
		}
		if t.Len() > 0 {
			fmt.Fprintf(b, "### %s\n\n", topologyLabel[tp])
			t.Write(b)
		}
	}
}

func (rp *report) writeRace(b *strings.Builder) {
	tiers := rp.tiers("race")
	if len(tiers) == 0 {
		return
	}
	fmt.Fprintf(b, "## The race — a hot drop with seat choice\n\n")
	fmt.Fprintf(b, "Seats held per second, from the gate to the last hold granted. Conflicts and seat-map reads are\n")
	fmt.Fprintf(b, "per granted hold. Deferred confirmations kept a real 40-minute hold through the crowd.\n\n")
	for _, tp := range rp.topos {
		head := []string{"Design"}
		for _, tier := range tiers {
			head = append(head, tierLabel(tier)+">")
		}
		t := md.NewTable(head...)
		d := md.NewTable("Design", "Tier", "Conflicts/hold>", "Map reads/hold>", "Engine retries>", "Gave up>", "Hold p50/p99 ms>", "Confirm p99 ms>", "Deferred ok>", "Early rej. (transient)>", "Confirm retries run / sold>", "Final buyer sold>")
		for _, id := range rp.designsIn(tp) {
			row := []string{designShort(id)}
			for _, tier := range tiers {
				rr := raceOf(rp.runs[tp][id], tier)
				row = append(row, raceCell(id, rr))
				if rr != nil {
					d.Row(designShort(id), tierLabel(tier), fmt.Sprintf("%.2f", rr.ConflictsPerHold), fmt.Sprintf("%.2f", rr.MapReadsPerHold),
						fmt.Sprint(rr.EngineRetries), fmt.Sprint(rr.GaveUp), md.MS(rr.HoldLatency.P50MS)+" / "+md.MS(rr.HoldLatency.P99MS),
						md.MS(rr.ConfirmLatency.P99MS), fmt.Sprintf("%d/%d", rr.DeferredConfirmed, rr.DeferredHolds),
						earlyCell(id, rr.Audit),
						fmt.Sprintf("%d / %d", rr.ConfirmRetries, rr.ConfirmRetrySuccesses), fmt.Sprint(rr.SweepSold))
				}
			}
			t.Row(row...)
		}
		fmt.Fprintf(b, "### %s\n\n", topologyLabel[tp])
		t.Write(b)
		fmt.Fprintf(b, "<details><summary>Race detail — %s</summary>\n\n", topologyLabel[tp])
		d.Write(b)
		fmt.Fprintf(b, "</details>\n\n")
	}
}

// writeReportCoverage is the mandatory answerability table (REPORTS.md v2
// section 2 / AM-02.1-.2): mechanical, from ReportCoverage, printed BEFORE
// the throughput table so a reader sees why a dash in that table is a dash --
// r01-r05's actual ops/s (where answerable) is in ## Reads, which already
// discovers them generically; r06 never appears there because it is never
// benchmarked, unanswerable in every design in this study.
func (rp *report) writeReportCoverage(b *strings.Builder) {
	fmt.Fprintf(b, "## Operational reports: answerability (REPORTS.md v2)\n\n")
	fmt.Fprintf(b, "The back office's and operations desk's questions, not the buyer's. Added to every\n")
	fmt.Fprintf(b, "design's `queries.sql` with no change to any existing statement, schema, index or write\n")
	fmt.Fprintf(b, "path. ✅ answerable, ⚠️ partial (caveat below), ✗ unanswerable. Throughput for whatever\n")
	fmt.Fprintf(b, "each design can actually answer is in [Reads](#reads) below, under the same report id.\n\n")

	for _, tp := range rp.topos {
		ids := rp.designsIn(tp)
		head := []string{"Report"}
		for _, id := range ids {
			head = append(head, designShort(id)+">")
		}
		t := md.NewTable(head...)
		for _, rname := range allReports {
			row := []string{reportName[rname]}
			for _, id := range ids {
				d, err := designByID(id)
				if err != nil {
					row = append(row, "—")
					continue
				}
				switch ReportCoverage(d)[rname].Status {
				case "answerable":
					row = append(row, "✅")
				case "partial":
					row = append(row, "⚠️")
				default:
					row = append(row, "✗")
				}
			}
			t.Row(row...)
		}
		fmt.Fprintf(b, "### %s\n\n", topologyLabel[tp])
		t.Write(b)
		fmt.Fprintf(b, "\n")
		notes := map[string]bool{}
		for _, id := range ids {
			d, err := designByID(id)
			if err != nil {
				continue
			}
			for _, rname := range allReports {
				st := ReportCoverage(d)[rname]
				if st.Note != "" && !notes[rname+st.Note] {
					notes[rname+st.Note] = true
					fmt.Fprintf(b, "- **%s** (%s): %s\n", reportName[rname], designShort(id), st.Note)
				}
			}
		}
		fmt.Fprintf(b, "\n")
	}
}

func (rp *report) writeReads(b *strings.Builder) {
	fmt.Fprintf(b, "## Reads\n\n")
	fmt.Fprintf(b, "Each question in isolation, fresh key per execution, ops/s.\n\n")
	type col struct {
		q    string
		tier int
	}
	var cols []col
	seen := map[col]bool{}
	for _, tp := range rp.topos {
		for _, r := range rp.runs[tp] {
			for _, x := range r.Reads {
				c := col{x.Query, x.Tier}
				if !seen[c] {
					seen[c] = true
					cols = append(cols, c)
				}
			}
		}
	}
	sort.Slice(cols, func(i, j int) bool {
		if cols[i].q != cols[j].q {
			return cols[i].q < cols[j].q
		}
		return cols[i].tier < cols[j].tier
	})
	for _, tp := range rp.topos {
		head := []string{"Question"}
		ids := rp.designsIn(tp)
		for _, id := range ids {
			head = append(head, designShort(id)+">")
		}
		t := md.NewTable(head...)
		for _, c := range cols {
			label := queryLabel[c.q]
			if c.tier > 0 {
				label += " — " + tierLabel(c.tier)
			}
			row := []string{label}
			for _, id := range ids {
				x := readOf(rp.runs[tp][id], c.q, c.tier)
				if x == nil {
					row = append(row, "—")
				} else {
					row = append(row, md.Ops(x.OpsPerSec))
				}
			}
			t.Row(row...)
		}
		fmt.Fprintf(b, "### %s\n\n", topologyLabel[tp])
		t.Write(b)
	}
}

func (rp *report) writeWrites(b *strings.Builder) {
	fmt.Fprintf(b, "## Isolated writes, publishing and storage\n\n")
	for _, tp := range rp.topos {
		t := md.NewTable("Design", "Hold+confirm /s>", "Release /s>", "Refund /s>", "Publish p50 ms by size>", "Storage>", "Audit")
		for _, id := range rp.designsIn(tp) {
			r := rp.runs[tp][id]
			cell := func(op string) string {
				w := writeOf(r, op, 0)
				if w == nil {
					return "—"
				}
				v := md.Ops(w.OpsPerSec)
				if w.Audit != nil && w.Audit.Violations() > 0 {
					v = "~~" + v + "~~ ❌"
				}
				if w.Errors > 0 {
					v += fmt.Sprintf(" (%d errors)", w.Errors)
				}
				return v
			}
			var pubs []string
			for _, tier := range allTiers {
				if w := writeOf(r, "publish", tier); w != nil {
					pubs = append(pubs, fmt.Sprintf("%s: %s", tierLabel(tier), md.MS(w.Latency.P50MS)))
				}
			}
			storage := "—"
			if r.Stats != nil {
				if r.Stats.Note != "" {
					storage = "n/a on YugabyteDB"
				} else {
					storage = md.Bytes(r.Stats.TotalBytes)
				}
			}
			var marks []string
			for _, w := range r.Writes {
				if w.Audit != nil {
					marks = append(marks, violationMarks(id, w.Audit)...)
				}
			}
			audit := "consistent"
			if len(marks) > 0 {
				audit = "❌ " + strings.Join(marks, ", ")
			}
			t.Row(designShort(id), cell("hold"), cell("release"), cell("cancel"), strings.Join(pubs, " · "), storage, audit)
		}
		fmt.Fprintf(b, "### %s\n\n", topologyLabel[tp])
		t.Write(b)
	}
}

func (rp *report) noiseFloor() (worst, typical float64, n int) {
	var ratios []float64
	for _, tp := range rp.topos {
		for i := 0; i < len(noiseDesigns); i++ {
			for j := i + 1; j < len(noiseDesigns); j++ {
				a, c := rp.runs[tp][noiseDesigns[i]], rp.runs[tp][noiseDesigns[j]]
				if a == nil || c == nil {
					continue
				}
				for _, q := range []string{"q05_customer_tickets", "q06_ticket_by_id"} {
					ra, rc := readOf(a, q, 0), readOf(c, q, 0)
					if ra == nil || rc == nil || ra.OpsPerSec <= 0 || rc.OpsPerSec <= 0 {
						continue
					}
					ratios = append(ratios, math.Max(ra.OpsPerSec/rc.OpsPerSec, rc.OpsPerSec/ra.OpsPerSec))
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
		fmt.Fprintf(b, "**This run's error bar.** Ten designs read tickets through byte-identical SQL on identical data,\n")
		fmt.Fprintf(b, "so their `q05`/`q06` differences are noise: %d comparisons, median disagreement **%.2fx**, worst **%.2fx**.\n", n, typical, worst)
		fmt.Fprintf(b, "Read rows inside the worst disagreement are omitted below. **Race, lifecycle and write rows are always\n")
		fmt.Fprintf(b, "shown**, and so is every correctness difference.\n\n")
	}
	for _, p := range pairs {
		t := md.NewTable("Topology", "Measurement", designShort(p.a)+">", designShort(p.b)+">", "Change")
		for _, tp := range rp.topos {
			ra, rb := rp.runs[tp][p.a], rp.runs[tp][p.b]
			if ra == nil || rb == nil {
				continue
			}
			for _, tier := range rp.tiers("race") {
				xa, xb := raceOf(ra, tier), raceOf(rb, tier)
				if xa == nil || xb == nil {
					continue
				}
				change := md.Ratio(xa.SeatsPerSec, xb.SeatsPerSec)
				if xa.Audit.Violations() > 0 || xb.Audit.Violations() > 0 {
					change += " — ❌ a design with violations has no valid speed"
				}
				t.Row(topologyLabel[tp], "race, "+tierLabel(tier)+" — seats/s", raceCell(p.a, xa), raceCell(p.b, xb), change)
			}
			for _, tier := range rp.tiers("lifecycle") {
				xa, xb := lifecycleOf(ra, tier), lifecycleOf(rb, tier)
				if xa == nil || xb == nil {
					continue
				}
				t.Row(topologyLabel[tp], "lifecycle, "+tierLabel(tier)+" — refused late/boundary/early (transient); violations",
					fmt.Sprintf("%d/%d/%d (%d); %d", xa.RejectedLate, xa.RejectedBoundary, xa.RejectedEarly, xa.RejectedEarlyTransient, xa.Audit.Violations()),
					fmt.Sprintf("%d/%d/%d (%d); %d", xb.RejectedLate, xb.RejectedBoundary, xb.RejectedEarly, xb.RejectedEarlyTransient, xb.Audit.Violations()), "")
			}
			for _, op := range []string{"hold", "release", "cancel"} {
				wa, wb := writeOf(ra, op, 0), writeOf(rb, op, 0)
				if wa != nil && wb != nil {
					t.Row(topologyLabel[tp], op+" /s", md.Ops(wa.OpsPerSec), md.Ops(wb.OpsPerSec), md.Ratio(wa.OpsPerSec, wb.OpsPerSec))
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
		}
		if t.Len() == 0 {
			continue
		}
		fmt.Fprintf(b, "### %s → %s\n\n*%s*\n\n", designShort(p.a), designShort(p.b), p.question)
		t.Write(b)
	}
}

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
	fmt.Fprintf(b, "Cells show *throttled periods / periods (throttled time)*: for the **database** container(s) over\n")
	fmt.Fprintf(b, "the whole cell (summed across nodes), and for the **client** per phase.\n\n")
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

// dbCPU parses the runner's before/after cpu.stat captures, summed over nodes.
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

// eventsLine renders {"race_100": 20, ...} as "race_100: 20 · ...", sorted.
func eventsLine(v any) string {
	m, ok := v.(map[string]any)
	if !ok {
		return fmt.Sprint(v)
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, len(keys))
	for i, k := range keys {
		parts[i] = fmt.Sprintf("%s: %v", k, m[k])
	}
	return strings.Join(parts, " · ")
}
