package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"adsplatform/core/catalog"
	"adsplatform/ports"
)

// verifyReports checks every report ReportCoverage marks answerable or
// partial, against truth computed independently from the freshly loaded
// dataset (Verify always runs before any write). r06 is unanswerable
// everywhere and is not checked -- there is nothing to check.
func verifyReports(ctx context.Context, db ports.DB, d Design, ds *Dataset, q map[string]catalog.Stmt, loadNow time.Time, vr *VerifyReport) {
	cov := ReportCoverage(d)

	runReport := func(name, key string, args []any, cols int, want []string) {
		st, ok := q[name]
		if !ok {
			vr.add(Check{Query: name, Key: key, Got: "error: statement not found"})
			return
		}
		got, err := collect(ctx, db, st.SQL, args, cols)
		c := Check{Query: name, Key: key}
		if err != nil {
			c.Got = "error: " + err.Error()
		} else if len(got) == len(want) {
			c.OK = true
		} else {
			c.Expect = fmt.Sprintf("%d rows: %s", len(want), abbreviate(want))
			c.Got = fmt.Sprintf("%d rows: %s", len(got), abbreviate(got))
		}
		vr.add(c)
	}

	// r01: holds expiring within a generous cutoff (effectively "ever") on the
	// busiest catalogue event -- so the truth is exactly that event's total
	// live-held seat count, regardless of individual expiry timings.
	if cov["r01"].Status != "unanswerable" {
		if ev := mostSold(ds, 0); ev != nil {
			st := ds.State[ev.ID]
			wantHeld := 0
			for _, s := range st {
				if s == stLive {
					wantHeld++
				}
			}
			far := loadNow.Add(100 * 365 * 24 * time.Hour)
			runReport("r01_holds_expiring_soon", fmt.Sprintf("event %d, generous cutoff", ev.ID),
				[]any{ev.ID, far}, 4, make([]string, wantHeld))
		}
	}

	// r02: seat status lookup on one known-sold and one known-held seat.
	// Checked by STATUS PREFIX, not row count (any seat with a row returns
	// exactly one row regardless of status, so a count-only check would pass
	// a design that always answered "held").
	statusCheck := func(key string, args []any, wantStatus string) {
		st, ok := q["r02_seat_status_lookup"]
		if !ok {
			vr.add(Check{Query: "r02_seat_status_lookup", Key: key, Got: "error: statement not found"})
			return
		}
		got, err := collect(ctx, db, st.SQL, args, 4)
		c := Check{Query: "r02_seat_status_lookup", Key: key}
		switch {
		case err != nil:
			c.Got = "error: " + err.Error()
		case len(got) != 1:
			c.Expect, c.Got = fmt.Sprintf("1 row starting %q", wantStatus), fmt.Sprintf("%d rows", len(got))
		case !strings.HasPrefix(got[0], wantStatus+"|"):
			c.Expect, c.Got = wantStatus+"|...", got[0]
		default:
			c.OK = true
		}
		vr.add(c)
	}
	if cov["r02"].Status != "unanswerable" {
		if ev := mostSold(ds, 0); ev != nil {
			v := ds.venue(ev.VenueID)
			st := ds.State[ev.ID]
			if soldSec := sectionMost(v, st, stSold); soldSec != 0 {
				if seat := firstOfState(v, st, soldSec, stSold); seat != 0 {
					statusCheck(fmt.Sprintf("event %d seat %d (sold)", ev.ID, seat), []any{ev.ID, seat}, "sold")
				}
			}
			if liveSec := sectionMost(v, st, stLive); liveSec != 0 {
				if seat := firstOfState(v, st, liveSec, stLive); seat != 0 {
					statusCheck(fmt.Sprintf("event %d seat %d (held)", ev.ID, seat), []any{ev.ID, seat}, "held")
				}
			}
		}
	}

	// r03: confirmed sales per section, for a window covering the whole
	// dataset -- the truth is one row per section that has any sale.
	if cov["r03"].Status != "unanswerable" {
		if ev := mostSold(ds, 0); ev != nil {
			sections := map[int32]bool{}
			for i := range ds.Sold {
				if ds.Sold[i].EventID == ev.ID {
					sections[ds.Sold[i].Section] = true
				}
			}
			since := loadEpoch
			until := saleSpanEnd()
			runReport("r03_section_sales_window", fmt.Sprintf("event %d, whole span", ev.ID),
				[]any{ev.ID, since, until}, 3, make([]string, len(sections)))
		}
	}

	// r04: last 50 confirmations for the busiest event.
	if cov["r04"].Status != "unanswerable" {
		if ev := mostSold(ds, 0); ev != nil {
			n := 0
			for i := range ds.Sold {
				if ds.Sold[i].EventID == ev.ID {
					n++
				}
			}
			want := min(50, n)
			runReport("r04_event_recent_confirmations", fmt.Sprintf("event %d (%d sold)", ev.ID, n),
				[]any{ev.ID}, 4, make([]string, want))
		}
	}

	// r05: customers whose last purchase falls in the window, both regimes --
	// same shape as study 01/02's recency check.
	if cov["r05"].Status != "unanswerable" {
		lastByCustomer := map[int64]int64{}
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
			runReport("r05_customers_last_purchase_window", "regime "+regime,
				[]any{since, until}, 2, make([]string, want))
		}
	}
}

// firstOfState returns the first seat id in section secNo whose ground-truth
// state matches want, or 0 if none.
func firstOfState(v *Venue, st []uint8, secNo int32, want uint8) int32 {
	sec := v.section(secNo)
	for s := sec.FirstSeat; s < sec.FirstSeat+int32(sec.Capacity); s++ {
		if st[s-1] == want {
			return s
		}
	}
	return 0
}
