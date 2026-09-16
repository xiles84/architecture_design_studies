package main

import (
	"testing"

	"adsplatform/core/catalog"
)

// TestReportCoverageMatchesSQL is AM-02.2's requirement: a declaration that
// can drift from the SQL is worth little, so this binds ReportCoverage to
// the actual catalogue. It fails if either direction disagrees --
// "answerable"/"partial" without a statement, or "unanswerable" with one.
func TestReportCoverageMatchesSQL(t *testing.T) {
	for _, d := range designs {
		src, err := readSQL(d.ID, "queries.sql")
		if err != nil {
			t.Fatalf("%s: %v", d.ID, err)
		}
		stmts, err := catalog.Parse(src)
		if err != nil {
			t.Fatalf("%s: %v", d.ID, err)
		}
		have := map[string]bool{}
		for _, s := range stmts {
			have[s.Name] = true
		}

		cov := ReportCoverage(d)
		for _, rname := range allReports {
			st, ok := cov[rname]
			if !ok {
				t.Errorf("%s: ReportCoverage has no entry for %s", d.ID, rname)
				continue
			}
			stmtName := rname + "_" + reportStmtSuffix[rname]
			exists := have[stmtName]
			switch st.Status {
			case "answerable", "partial":
				if !exists {
					t.Errorf("%s: %s declared %s but queries.sql has no %q statement",
						d.ID, rname, st.Status, stmtName)
				}
			case "unanswerable":
				if exists {
					t.Errorf("%s: %s declared unanswerable but queries.sql DOES have %q -- coverage is stale",
						d.ID, rname, stmtName)
				}
				if st.Note == "" {
					t.Errorf("%s: %s declared unanswerable with no reason", d.ID, rname)
				}
			default:
				t.Errorf("%s: %s has unknown status %q", d.ID, rname, st.Status)
			}
			if st.Status == "partial" && st.Note == "" {
				t.Errorf("%s: %s declared partial with no caveat", d.ID, rname)
			}
		}
	}
}

// TestReportWindowRegimesAreDistinctAndDerived checks the two window regimes
// against the dataset's own span, the same property study 01's equivalent
// test checks (RECENCY.md section 2): trailing ends at the data's own
// horizon; historical ends well before it; both regimes are the same width.
func TestReportWindowRegimesAreDistinctAndDerived(t *testing.T) {
	tSince, tUntil := reportWindowFor("trailing")
	hSince, hUntil := reportWindowFor("historical")

	if !tUntil.Equal(saleSpanEnd()) {
		t.Fatalf("trailing regime should end at the dataset's own sale span, got %s", tUntil)
	}
	if got := tUntil.Sub(tSince); got != reportWindow {
		t.Fatalf("trailing window width = %s, want %s", got, reportWindow)
	}
	if !hUntil.Before(saleSpanEnd()) {
		t.Fatal("historical regime must end strictly before the dataset's own sale span")
	}
	if got := saleSpanEnd().Sub(hUntil); got != historicalLead {
		t.Fatalf("historical lead = %s, want %s", got, historicalLead)
	}
	if got := hUntil.Sub(hSince); got != reportWindow {
		t.Fatalf("historical window width = %s, want %s", got, reportWindow)
	}
	if !hUntil.Before(tSince) {
		t.Fatal("the two regimes should not overlap at this dataset's scale")
	}
}
