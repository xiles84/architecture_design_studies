package main

import (
	"context"
	"fmt"
	"sort"

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

	// r02: last 50 sales for one event, newest first.
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
		st := q["r02_event_recent_buyers"]
		args, _ := catalog.Bind(st.Params, map[string]any{"event_id": ev.ID})
		gotN, err := countRows(ctx, db, st.SQL, args)
		v.add(Check{Query: "r02_event_recent_buyers", Key: fmt.Sprintf("event %d (%d sold)", ev.ID, len(evTickets)),
			OK: err == nil && gotN == wantN, Expect: fmt.Sprintf("%d rows", wantN), Got: fmt.Sprintf("%d rows, err=%v", gotN, err)})
	}

	// r03: customers whose last purchase falls in the window, both regimes.
	if cov["r03"].Status != "unanswerable" {
		lastByCustomer := map[int64]int64{} // customer -> unix nanos of last sale
		for i := range ds.Sold {
			t := &ds.Sold[i]
			if cur, ok := lastByCustomer[t.CustomerID]; !ok || t.SoldAt.UnixNano() > cur {
				lastByCustomer[t.CustomerID] = t.SoldAt.UnixNano()
			}
		}
		for _, regime := range []string{"trailing", "historical"} {
			since, until := reportWindowFor(regime)
			want := 0
			for _, last := range lastByCustomer {
				if last >= since.UnixNano() && last < until.UnixNano() {
					want++
				}
			}
			if want > 100 {
				want = 100
			}
			st := q["r03_customers_last_purchase_window"]
			args, _ := catalog.Bind(st.Params, map[string]any{"since": since, "until": until})
			gotN, err := countRows(ctx, db, st.SQL, args)
			v.add(Check{Query: "r03_customers_last_purchase_window@" + regime, Key: "top-100 by last purchase",
				OK: err == nil && gotN == want, Expect: fmt.Sprintf("%d rows", want), Got: fmt.Sprintf("%d rows, err=%v", gotN, err)})
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
