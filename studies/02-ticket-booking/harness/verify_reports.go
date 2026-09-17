package main

import (
	"context"
	"fmt"
	"sort"
	"time"

	"adsplatform/core/catalog"
	"adsplatform/ports"
)

// verifyReports checks every report ReportCoverage marks answerable or
// partial, against truth computed independently in Go from the freshly
// loaded dataset (Verify always runs before any write, so ds.Sold is exactly
// the ground truth "sold" set, and reservation/sale_event start empty).
// An unanswerable report is not checked -- there is nothing to check.
func verifyReports(ctx context.Context, db ports.DB, d Design, ds *Dataset, q map[string]catalog.Stmt, k keySet, v *VerifyReport) {
	cov := ReportCoverage(d)
	ev := k.busyEvent
	if ev == nil && len(ds.Events) > 0 {
		ev = &ds.Events[0]
	}

	// r01: sales and revenue for one event, both regimes.
	if cov["r01"].Status != "unanswerable" && ev != nil {
		for _, regime := range []string{"trailing", "historical"} {
			since, until := reportWindowFor(regime)
			var wantN, wantRevenue int64
			for _, t := range ds.Sold {
				if t.EventID == ev.ID && !t.SoldAt.Before(since) && t.SoldAt.Before(until) {
					wantN++
					wantRevenue += t.PriceCents
				}
			}
			var gotN, gotRevenue int64
			st := q["r01_event_sales_window"]
			args, _ := catalog.Bind(st.Params, map[string]any{"event_id": ev.ID, "since": since, "until": until})
			err := db.QueryRow(ctx, st.SQL, args...).Scan(&gotN, &gotRevenue)
			ok := err == nil && gotN == wantN && gotRevenue == wantRevenue
			v.add(Check{Query: "r01_event_sales_window@" + regime, Key: fmt.Sprintf("event %d", ev.ID),
				OK: ok, Expect: fmt.Sprintf("%d tickets, %d cents", wantN, wantRevenue),
				Got: gotOrErrPair(gotN, gotRevenue, err)})
		}
	}

	// r02: last 50 sales for one event, newest first. Checked by VALUE (EH-02
	// AM-03.6): the multiset of returned sold_at must equal the truth's top-50
	// -- a row-count-only check would pass a query that returned the right
	// number of the WRONG 50 rows (e.g. oldest instead of newest).
	if cov["r02"].Status != "unanswerable" && ev != nil {
		var evTickets []*Ticket
		for i := range ds.Sold {
			if ds.Sold[i].EventID == ev.ID {
				evTickets = append(evTickets, &ds.Sold[i])
			}
		}
		sort.Slice(evTickets, func(i, j int) bool {
			if !evTickets[i].SoldAt.Equal(evTickets[j].SoldAt) {
				return evTickets[i].SoldAt.After(evTickets[j].SoldAt)
			}
			return evTickets[i].ID > evTickets[j].ID
		})
		wantN := min(50, len(evTickets))
		var want []time.Time
		for _, t := range evTickets[:wantN] {
			want = append(want, t.SoldAt)
		}
		st := q["r02_event_recent_buyers"]
		args, _ := catalog.Bind(st.Params, map[string]any{"event_id": ev.ID})
		got, err := scanTimeColumn(ctx, db, st.SQL, args, 3, 2)
		ok := err == nil && sameTimeMultiset(got, want)
		v.add(Check{Query: "r02_event_recent_buyers", Key: fmt.Sprintf("event %d (%d sold)", ev.ID, len(evTickets)),
			OK: ok, Expect: fmt.Sprintf("%d rows, sold_at multiset %s", wantN, abbreviateTimes(want)),
			Got: gotTimesOrErr(got, err)})
	}

	// r03: customers whose last purchase falls in the window, both regimes.
	// Checked by VALUE: the multiset of returned last_at must equal the
	// truth's top-100 by recency (EH-02 AM-03.6).
	if cov["r03"].Status != "unanswerable" {
		lastByCustomer := map[int64]time.Time{}
		purchasedInWindow := map[int64]bool{} // naive formulation, historical regime only
		for i := range ds.Sold {
			t := &ds.Sold[i]
			if cur, ok := lastByCustomer[t.CustomerID]; !ok || t.SoldAt.After(cur) {
				lastByCustomer[t.CustomerID] = t.SoldAt
			}
		}
		for _, regime := range []string{"trailing", "historical"} {
			since, until := reportWindowFor(regime)
			var want []time.Time
			for _, last := range lastByCustomer {
				if !last.Before(since) && last.Before(until) {
					want = append(want, last)
				}
			}
			sort.Slice(want, func(i, j int) bool { return want[i].After(want[j]) })
			if len(want) > 100 {
				want = want[:100]
			}
			st := q["r03_customers_last_purchase_window"]
			args, _ := catalog.Bind(st.Params, map[string]any{"since": since, "until": until})
			got, err := scanTimeColumn(ctx, db, st.SQL, args, 2, 1)
			ok := err == nil && sameTimeMultiset(got, want)
			c := Check{Query: "r03_customers_last_purchase_window@" + regime, Key: "top-100 by last purchase",
				OK: ok, Expect: fmt.Sprintf("%d rows, last_at multiset %s", len(want), abbreviateTimes(want)),
				Got: gotTimesOrErr(got, err)}

			// The recency trap (EH-02 AM-03.6): in the historical regime, also
			// compute the naive formulation (anyone who PURCHASED in the
			// window, not just whose LAST purchase falls in it). If the two
			// formulations' top-100 multisets agree, this dataset cannot tell
			// a correct r03 from a naive one apart at this size -- a gap in
			// the gate's own power, surfaced rather than hidden.
			if regime == "historical" {
				for i := range ds.Sold {
					t := &ds.Sold[i]
					if !t.SoldAt.Before(since) && t.SoldAt.Before(until) {
						purchasedInWindow[t.CustomerID] = true
					}
				}
				var naive []time.Time
				for cust := range purchasedInWindow {
					naive = append(naive, lastByCustomer[cust])
				}
				sort.Slice(naive, func(i, j int) bool { return naive[i].After(naive[j]) })
				if len(naive) > 100 {
					naive = naive[:100]
				}
				if sameTimeMultiset(naive, want) {
					c.Warning = "the naive formulation (purchased in the window) gives the same top-100 as the correct one (last purchase in the window) -- this dataset cannot distinguish them at this size"
				}
			}
			v.add(c)
		}
	}

	// r04: sell-through for the band with the most events.
	if cov["r04"].Status != "unanswerable" {
		type row struct{ event, cap, sold int64 }
		var want []row
		for i := range ds.Events {
			e := &ds.Events[i]
			if e.BandID == k.band {
				want = append(want, row{e.ID, int64(e.Capacity), int64(e.InitialSold)})
			}
		}
		sort.Slice(want, func(i, j int) bool { return want[i].event < want[j].event })
		st := q["r04_band_sellthrough"]
		args, _ := catalog.Bind(st.Params, map[string]any{"band_id": k.band})
		rows, err := db.Query(ctx, st.SQL, args...)
		var got []row
		if err == nil {
			for rows.Next() {
				var r row
				var pct float64
				if serr := rows.Scan(&r.event, &r.cap, &r.sold, &pct); serr == nil {
					got = append(got, r)
				}
			}
			rows.Close()
			err = rows.Err()
		}
		ok := err == nil && len(got) == len(want)
		if ok {
			for i := range want {
				if got[i] != want[i] {
					ok = false
					break
				}
			}
		}
		v.add(Check{Query: "r04_band_sellthrough", Key: fmt.Sprintf("band %d (%d events)", k.band, len(want)),
			OK: ok, Expect: fmt.Sprintf("%d rows", len(want)), Got: fmt.Sprintf("%d rows, err=%v", len(got), err)})
	}

	// r06: outstanding holds, checked only for Holds designs -- reservation is
	// empty on a fresh load, so the true answer is always zero here.
	if cov["r06"].Status != "unanswerable" && ev != nil {
		var got int64
		st := q["r06_outstanding_holds"]
		args, _ := catalog.Bind(st.Params, map[string]any{"event_id": ev.ID})
		err := db.QueryRow(ctx, st.SQL, args...).Scan(&got)
		v.add(Check{Query: "r06_outstanding_holds", Key: fmt.Sprintf("event %d, fresh load", ev.ID),
			OK: err == nil && got == 0, Expect: "0", Got: gotOrErr(got, err)})
	}

	// r05: refunds in a window, checked only where answerable (X1) -- zero
	// cancellations have happened on a fresh load, in either regime.
	if cov["r05"].Status == "answerable" {
		for _, regime := range []string{"trailing", "historical"} {
			since, until := reportWindowFor(regime)
			var gotN, gotCents int64
			st := q["r05_refunds_window"]
			args, _ := catalog.Bind(st.Params, map[string]any{"since": since, "until": until})
			err := db.QueryRow(ctx, st.SQL, args...).Scan(&gotN, &gotCents)
			v.add(Check{Query: "r05_refunds_window@" + regime, Key: "fresh load",
				OK: err == nil && gotN == 0 && gotCents == 0, Expect: "0 refunds, 0 cents", Got: gotOrErrPair(gotN, gotCents, err)})
		}
	}
}

func gotOrErrPair(a, b int64, err error) string {
	if err != nil {
		return "error: " + err.Error()
	}
	return fmt.Sprintf("%d, %d", a, b)
}

// countRows runs a query with any number of bound parameters and returns how
// many rows it produced, without caring about their column types -- what
// r03's check needs (row count against the independently computed truth),
// where scanStrings' fixed-int64-columns assumption does not fit (r03's
// second column is a timestamp, not an integer).
func countRows(ctx context.Context, db ports.DB, sql string, args []any) (int, error) {
	rows, err := db.Query(ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	n := 0
	for rows.Next() {
		n++
	}
	return n, rows.Err()
}

// scanTimeColumn runs a query and returns one TIMESTAMPTZ column's values,
// ignoring the rest -- what r02/r03's value checks need (EH-02 AM-03.6)
// without a fixed-column-type scanner for every other column in the row.
func scanTimeColumn(ctx context.Context, db ports.DB, sql string, args []any, cols, timeIdx int) ([]time.Time, error) {
	rows, err := db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []time.Time
	for rows.Next() {
		dest := make([]any, cols)
		var t time.Time
		for i := range dest {
			if i == timeIdx {
				dest[i] = &t
			} else {
				dest[i] = new(any)
			}
		}
		if err := rows.Scan(dest...); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

// sameTimeMultiset compares two multisets of timestamps at microsecond
// precision (EH-02 AM-03.6): PostgreSQL and Go both carry microsecond
// resolution, so anything finer is a false disagreement, not a real one.
func sameTimeMultiset(a, b []time.Time) bool {
	if len(a) != len(b) {
		return false
	}
	ua := make([]int64, len(a))
	ub := make([]int64, len(b))
	for i, t := range a {
		ua[i] = t.UnixMicro()
	}
	for i, t := range b {
		ub[i] = t.UnixMicro()
	}
	sort.Slice(ua, func(i, j int) bool { return ua[i] < ua[j] })
	sort.Slice(ub, func(i, j int) bool { return ub[i] < ub[j] })
	for i := range ua {
		if ua[i] != ub[i] {
			return false
		}
	}
	return true
}

func abbreviateTimes(ts []time.Time) string {
	if len(ts) == 0 {
		return "[]"
	}
	n := len(ts)
	if n > 3 {
		n = 3
	}
	s := "["
	for i := 0; i < n; i++ {
		if i > 0 {
			s += ", "
		}
		s += ts[i].UTC().Format(time.RFC3339Nano)
	}
	if len(ts) > n {
		s += fmt.Sprintf(", ... (%d more)", len(ts)-n)
	}
	return s + "]"
}

func gotTimesOrErr(ts []time.Time, err error) string {
	if err != nil {
		return "error: " + err.Error()
	}
	return fmt.Sprintf("%d rows, sold_at/last_at multiset %s", len(ts), abbreviateTimes(ts))
}

// ---------------------------------------------------------------------------
// Post-write report checks (EH-02 AM-03.5). Verify's report checks run once,
// on the freshly loaded dataset, before anything has been booked or
// cancelled -- which cannot see a report break under real writes, only under
// the load's own ground truth. This runs the same idea again after a write
// phase, against the harness's OWN counters (EventRef.booked/cancelled,
// already used by the overbooking audit), on whatever the phase actually
// did. Not a correctness gate: post-write violations here are reported like
// audit violations, not aborted on.
// ---------------------------------------------------------------------------

// mostActiveEvent returns the event this phase touched the most (booked plus
// cancelled), or nil if the phase touched nothing -- e.g. a write op whose
// pool was exhausted immediately, which is itself worth being unable to check.
func mostActiveEvent(world *World) *EventRef {
	var best *EventRef
	var bestActivity int64
	for _, e := range world.snapshot() {
		activity := e.booked.Load() + e.cancelled.Load()
		if activity > bestActivity {
			best, bestActivity = e, activity
		}
	}
	return best
}

// postWriteReportChecks checks r01 (for the phase's most active event) and
// r05 (summed over every event) against the harness's own counters. window
// is fixed at [loadEpoch, now+24h) -- wide enough to contain every possible
// sold_at/at this run could have produced, in either regime.
func postWriteReportChecks(ctx context.Context, db ports.DB, d Design, world *World, phase string) ([]Check, error) {
	cov := ReportCoverage(d)
	if cov["r01"].Status == "unanswerable" && cov["r05"].Status != "answerable" {
		return nil, nil
	}
	stmts, err := mustStmts(d.ID, "queries.sql")
	if err != nil {
		return nil, err
	}
	q := catalog.Map(stmts)
	since, until := loadEpoch, time.Now().Add(24*time.Hour)
	var checks []Check

	if cov["r01"].Status != "unanswerable" {
		if ev := mostActiveEvent(world); ev != nil {
			booked, cancelled := ev.booked.Load(), ev.cancelled.Load()
			want := ev.InitialSold + booked
			if cov["r01"].Status == "partial" {
				want -= cancelled
			}
			st := q["r01_event_sales_window"]
			args, _ := catalog.Bind(st.Params, map[string]any{"event_id": ev.ID, "since": since, "until": until})
			var gotN, gotRevenue int64
			err := db.QueryRow(ctx, st.SQL, args...).Scan(&gotN, &gotRevenue)
			checks = append(checks, Check{Query: "r01_event_sales_window@" + phase,
				Key: fmt.Sprintf("event %d (booked=%d, cancelled=%d)", ev.ID, booked, cancelled),
				OK:  err == nil && gotN == want, Expect: fmt.Sprintf("%d tickets", want), Got: gotOrErrPair(gotN, gotRevenue, err)})
		}
	}

	if cov["r05"].Status == "answerable" {
		var wantCancelled, tolerance int64
		for _, e := range world.snapshot() {
			wantCancelled += e.cancelled.Load()
			tolerance += e.ambiguous.Load()
		}
		st := q["r05_refunds_window"]
		args, _ := catalog.Bind(st.Params, map[string]any{"since": since, "until": until})
		var gotN, gotCents int64
		err := db.QueryRow(ctx, st.SQL, args...).Scan(&gotN, &gotCents)
		diff := gotN - wantCancelled
		if diff < 0 {
			diff = -diff
		}
		checks = append(checks, Check{Query: "r05_refunds_window@" + phase,
			Key: fmt.Sprintf("every event (harness cancelled=%d, tolerance=%d)", wantCancelled, tolerance),
			OK:  err == nil && diff <= tolerance, Expect: fmt.Sprintf("%d refunds (+/- %d)", wantCancelled, tolerance),
			Got: gotOrErrPair(gotN, gotCents, err)})
	}
	return checks, nil
}
