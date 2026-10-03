package selectbench

import (
	"strings"
	"testing"
)

func cell(stmt, group, plan string, ms, spread float64) Cell {
	return Cell{
		RunID: "r1", Topology: "pg-single", Statement: stmt, Group: group,
		Status: StatusValid, MedianMS: ms, SpreadPct: spread, Trials: 3, Plan: plan, TimingMeasured: true,
	}
}

// A generated report may hold measurements only; a caller that tries to embed a
// conclusion is refused.
func TestReportRefusesInterpretation(t *testing.T) {
	_, err := RenderComparison(ReportInput{
		Cells:          []Cell{cell("s01_join_where", "sel_full_rows", seqScan, 10, 5)},
		Interpretation: "the join template is fastest",
	})
	if err == nil {
		t.Fatal("a report accepted an interpretation: measurements and conclusions must stay separate")
	}
}

func TestReportNamesExcludedCellsAndDoesNotRankThem(t *testing.T) {
	out, err := RenderComparison(ReportInput{Cells: []Cell{
		cell("s01_join_where", "sel_full_rows", seqScan, 10, 5),
		{Statement: "s02_filtered_derived", Group: "sel_full_rows", Topology: "pg-single",
			Status: StatusTimeout, MedianMS: 1, Plan: seqScan, Note: "generator cap reached"},
	}})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "**timeout**") {
		t.Error("the timed-out cell was not named in the excluded section")
	}
	if strings.Contains(out, "`s01_join_where` vs `s02_filtered_derived`") {
		t.Error("a timed-out cell was ranked against a valid one")
	}
}

func TestReportDoesNotPoolControlsOrAcrossGroups(t *testing.T) {
	out, err := RenderComparison(ReportInput{Cells: []Cell{
		cell("s01_join_where", "sel_full_rows", seqScan, 10, 5),
		cell("x01_optional_where", ControlGroup, indexScan, 1, 5),
		cell("s06_exists_parent", "sel_membership", indexScan, 2, 5),
	}})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "a deliberately wrong control may never be pooled") {
		t.Error("the control was not refused as a comparison")
	}
	if !strings.Contains(out, "different semantic groups answer different questions") {
		t.Error("a cross-group pair was not refused")
	}
	if strings.Contains(out, "faster by") {
		t.Error("an incomparable pair was reported as a resolved difference")
	}
}

func TestReportDetectsSamePlan(t *testing.T) {
	out, err := RenderComparison(ReportInput{Cells: []Cell{
		cell("s01_join_where", "sel_full_rows", seqScan, 10, 5),
		cell("s02_filtered_derived", "sel_full_rows", seqScanResampled, 10.2, 5),
	}})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "same plan structure") {
		t.Errorf("two runs of the same plan were not recognised as the same structure:\n%s", out)
	}
}

func TestReportDoesNotResolveAWholeRunWinnerFromOneRun(t *testing.T) {
	out, err := RenderComparison(ReportInput{Cells: []Cell{
		cell("s01_join_where", "sel_full_rows", seqScan, 20, 3),
		cell("s02_filtered_derived", "sel_full_rows", indexScan, 5, 3),
	}})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "faster by") || !strings.Contains(out, "whole-run replication required") {
		t.Errorf("one run was treated as a replicated winner:\n%s", out)
	}
}

func TestReportCallsSmallDifferenceUnresolved(t *testing.T) {
	out, err := RenderComparison(ReportInput{Cells: []Cell{
		cell("s01_join_where", "sel_full_rows", seqScan, 10.0, 8),
		cell("s02_filtered_derived", "sel_full_rows", indexScan, 10.1, 8),
	}})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "no resolved difference") {
		t.Errorf("a difference inside the spread was not left unresolved:\n%s", out)
	}
}
