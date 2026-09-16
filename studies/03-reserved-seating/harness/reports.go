package main

import (
	"time"
)

// ---------------------------------------------------------------------------
// Operational reports (study 03 v2, REPORTS.md; AM-02, EH-02).
//
// r01-r06 are the box office's and operations desk's questions, added to
// every design's queries.sql without touching any existing statement, schema,
// index or write path (AM-02.1). Answerability is a mechanical fact about a
// design's SQL, not an opinion -- see ReportCoverage and its unit test in
// reports_test.go, which fails if the two disagree.
// ---------------------------------------------------------------------------

type ReportStatus struct {
	Status string `json:"status"` // "answerable" | "partial" | "unanswerable"
	Note   string `json:"note,omitempty"`
}

var allReports = []string{"r01", "r02", "r03", "r04", "r05", "r06"}

var reportName = map[string]string{
	"r01": "Holds expiring soon",
	"r02": "Seat status lookup",
	"r03": "Confirmed sales per section in a window",
	"r04": "An event's recent confirmations",
	"r05": "Customers whose last purchase falls in a window",
	"r06": "Hold funnel in a window",
}

const (
	appClockNote    = "the answer depends on whose clock is asked -- the very ambiguity this design exists to expose"
	unanswerableR06 = "releasing an expired hold clears the row; no design in this study keeps a hold-event history"
)

// ReportCoverage states, per report, whether design d can answer it.
// r03-r05 read the `ticket` table, byte-identical across every design, so
// they are answerable everywhere. r01-r02 depend on d.Layout, the same flag
// the loader already uses to shape the inventory tables. r06 is
// unanswerable everywhere -- REPORTS.md's own finding, verified against all
// fourteen designs' write paths before this handoff was written (AM-02.0).
func ReportCoverage(d Design) map[string]ReportStatus {
	m := map[string]ReportStatus{
		"r01": {Status: "answerable"},
		"r02": {Status: "answerable"},
		"r03": {Status: "answerable"},
		"r04": {Status: "answerable"},
		"r05": {Status: "answerable"},
		"r06": {Status: "unanswerable", Note: unanswerableR06},
	}
	if d.AppClock {
		m["r01"] = ReportStatus{Status: "partial", Note: appClockNote}
	}
	return m
}

// ---------------------------------------------------------------------------
// Window regimes (RECENCY.md section 2, reused for r05 exactly as study 01
// and study 02 use it -- the trailing regime's cheap wrong answer happens to
// be right; the historical regime is where it would not be). Derived from
// the dataset's own load epoch, never from time.Now(), so two runs of r03/
// r05 stay comparable. r01's "within" is a genuinely live-clock question (an
// operations desk asks it against real time), so it alone uses an absolute
// cutoff computed from the wall clock when the benchmark runs.
// ---------------------------------------------------------------------------

const (
	reportWindow    = 7 * 24 * time.Hour
	historicalLead  = 30 * 24 * time.Hour
	saleSpanSeconds = 90 * 24 * 3600
	// holdsLookahead: r01 asks "expiring in the next N minutes"; ten minutes
	// is the scale a sweeper interval actually cares about.
	holdsLookahead = 10 * time.Minute
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
