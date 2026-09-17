package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"adsplatform/core/catalog"
	"adsplatform/core/measure"
	"adsplatform/ports"
)

// ---------------------------------------------------------------------------
// Operational reports (study 02 v2, REPORTS.md; AM-02, EH-02).
//
// r01-r06 are the back office's questions, added to every design's queries.sql
// without touching any existing statement, schema, index or write path
// (AM-02.1). Answerability is a mechanical fact about a design's SQL, not an
// opinion -- see ReportCoverage and its unit test in workload_test.go, which
// fails if the two disagree.
// ---------------------------------------------------------------------------

type ReportStatus struct {
	Status string `json:"status"` // "answerable" | "partial" | "unanswerable"
	Note   string `json:"note,omitempty"`
}

var allReports = []string{"r01", "r02", "r03", "r04", "r05", "r06"}

var reportName = map[string]string{
	"r01": "Sales and revenue in a window",
	"r02": "An event's recent buyers",
	"r03": "Customers whose last purchase falls in a window",
	"r04": "Sell-through per event, by band",
	"r05": "Refunds in a window",
	"r06": "Outstanding holds",
}

const partialRefundNote = "a ticket sold in the window and later refunded is not counted, because cancellation erases the sale"
const unanswerableRefundNote = "cancellation deletes or resets the sale; nothing records that it happened"
const unanswerableHoldsNote = "no reservation table in this design"

// partialLastPurchaseNote (EH-02 AM-03.6): r03 has the same caveat as r01,
// not the clean pass it was originally declared with. When a customer's most
// recent purchase is the one refunded, the ticket table dates them by an
// earlier purchase or drops them, because cancellation erases the sale that
// was their last one -- the same erasure r01 already accounts for.
const partialLastPurchaseNote = "a customer whose most recent purchase was refunded is dated by an earlier purchase, or drops out, because cancellation erases the sale"

// ReportCoverage states, per report, whether design d can answer it -- computed
// from the same flags that shape its schema, so the declaration cannot drift
// from what the design actually is without also changing what it does.
func ReportCoverage(d Design) map[string]ReportStatus {
	m := map[string]ReportStatus{
		"r02": {Status: "answerable"},
		"r04": {Status: "answerable"},
	}
	if d.Ledger {
		m["r01"] = ReportStatus{Status: "answerable"}
		m["r03"] = ReportStatus{Status: "answerable"}
		m["r05"] = ReportStatus{Status: "answerable"}
	} else {
		m["r01"] = ReportStatus{Status: "partial", Note: partialRefundNote}
		m["r03"] = ReportStatus{Status: "partial", Note: partialLastPurchaseNote}
		m["r05"] = ReportStatus{Status: "unanswerable", Note: unanswerableRefundNote}
	}
	if d.Holds {
		m["r06"] = ReportStatus{Status: "answerable"}
	} else {
		m["r06"] = ReportStatus{Status: "unanswerable", Note: unanswerableHoldsNote}
	}
	return m
}

// ---------------------------------------------------------------------------
// Window regimes (RECENCY.md section 2, reused here for r01/r03/r05: study
// 01's trailing/historical split, and the trap it exists to catch -- in the
// trailing regime "happened in the window" and "last happened in the window"
// coincide because nothing is newer; the historical regime is where a design
// that only checked the former would be wrong). Derived from the dataset's
// own load epoch, never from time.Now(), so two runs stay comparable.
// Sold-at timestamps span [loadEpoch, loadEpoch+90d).
// ---------------------------------------------------------------------------

const (
	reportWindow    = 7 * 24 * time.Hour
	historicalLead  = 30 * 24 * time.Hour
	saleSpanSeconds = 90 * 24 * 3600
)

func saleSpanEnd() time.Time { return loadEpoch.Add(saleSpanSeconds * time.Second) }

func reportWindowFor(regime string) (since, until time.Time) {
	until = saleSpanEnd()
	if regime == "historical" {
		until = until.Add(-historicalLead)
	}
	since = until.Add(-reportWindow)
	return since, until
}

// ---------------------------------------------------------------------------
// Multi-parameter read helper. BenchmarkReads' own run() binds exactly one
// positional key; the reports need up to three, so this is a sibling rather
// than a reshape of that call path (AM-02.4) -- the five original read
// questions keep executing exactly as they did before this file existed.
// ---------------------------------------------------------------------------

// BenchmarkReports runs r01-r06 for every report ReportCoverage marks
// answerable or partial. An unanswerable report is never timed and never
// appears as a zero -- it is simply absent from the returned slice, and the
// caller records the coverage table separately.
func BenchmarkReports(ctx context.Context, db ports.DB, d Design, ds *Dataset, s Settings) ([]ReadResult, error) {
	stmts, err := mustStmts(d.ID, "queries.sql")
	if err != nil {
		return nil, err
	}
	q := catalog.Map(stmts)
	cov := ReportCoverage(d)

	catalogueEvents := ds.eventsOf("catalogue", 0)
	bands := int64(len(ds.Bands))
	var holdEvents []*Event
	if d.Holds {
		for i := range ds.Events {
			holdEvents = append(holdEvents, &ds.Events[i])
		}
	}

	budget := measure.Budget{Workers: s.Workers, Warmup: s.Warmup, Duration: s.Duration, Trials: s.Trials}
	// run draws a FRESH key per execution via keyVals (never a single repeated
	// key -- the same rule BenchmarkReads' own comment states, and for the
	// same reason: a repeated key measures the buffer cache, not the design).
	run := func(name, label string, fixed map[string]any, keyVals func(r *rand.Rand) map[string]any) ReadResult {
		st, ok := q[name]
		if !ok {
			panic(fmt.Sprintf("%s: report %q declared answerable but has no statement", d.ID, name))
		}
		res := measure.Run(ctx, budget, label, func(ctx context.Context, r *rand.Rand) (measure.Outcome, error) {
			vals := make(map[string]any, len(fixed)+2)
			for k, v := range fixed {
				vals[k] = v
			}
			for k, v := range keyVals(r) {
				vals[k] = v
			}
			args, err := catalog.Bind(st.Params, vals)
			if err != nil {
				return measure.Outcome{}, err
			}
			rows, err := db.Query(ctx, st.SQL, args...)
			if err != nil {
				return measure.Outcome{}, err
			}
			defer rows.Close()
			for rows.Next() {
				// Drain: the server may not have done the work until fetched.
			}
			return measure.Outcome{}, rows.Err()
		})
		fmt.Printf("    %-40s %10.1f ops/s  p50=%7.3fms p99=%8.3fms  errors=%d%s\n",
			label, res.OpsPerSec, res.Latency.P50MS, res.Latency.P99MS, res.Errors, measure.SpreadNote(res))
		return ReadResult{Result: res, Query: label}
	}

	eventKey := func(r *rand.Rand) map[string]any {
		return map[string]any{"event_id": catalogueEvents[r.Intn(len(catalogueEvents))].ID}
	}

	var out []ReadResult
	if cov["r01"].Status != "unanswerable" {
		for _, regime := range []string{"trailing", "historical"} {
			since, until := reportWindowFor(regime)
			out = append(out, run("r01_event_sales_window", "r01_event_sales_window@"+regime,
				map[string]any{"since": since, "until": until}, eventKey))
		}
	}
	if cov["r02"].Status != "unanswerable" {
		out = append(out, run("r02_event_recent_buyers", "r02_event_recent_buyers", nil, eventKey))
	}
	if cov["r03"].Status != "unanswerable" {
		for _, regime := range []string{"trailing", "historical"} {
			since, until := reportWindowFor(regime)
			out = append(out, run("r03_customers_last_purchase_window", "r03_customers_last_purchase_window@"+regime,
				map[string]any{"since": since, "until": until}, func(r *rand.Rand) map[string]any { return map[string]any{} }))
		}
	}
	if cov["r04"].Status != "unanswerable" {
		out = append(out, run("r04_band_sellthrough", "r04_band_sellthrough", nil,
			func(r *rand.Rand) map[string]any { return map[string]any{"band_id": r.Int63n(bands) + 1} }))
	}
	if cov["r05"].Status != "unanswerable" {
		for _, regime := range []string{"trailing", "historical"} {
			since, until := reportWindowFor(regime)
			out = append(out, run("r05_refunds_window", "r05_refunds_window@"+regime,
				map[string]any{"since": since, "until": until}, func(r *rand.Rand) map[string]any { return map[string]any{} }))
		}
	}
	if cov["r06"].Status != "unanswerable" {
		out = append(out, run("r06_outstanding_holds", "r06_outstanding_holds", nil,
			func(r *rand.Rand) map[string]any {
				return map[string]any{"event_id": holdEvents[r.Intn(len(holdEvents))].ID}
			}))
	}
	return out, nil
}
