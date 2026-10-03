package selectbench

import (
	"context"
	"math"
	"strings"
	"testing"
)

func TestOrderedPageCannotBeRepairedByCanonicalization(t *testing.T) {
	c, fx, a := load(t), Teaching(), args()
	rows := rowsFromOracle(t, "sel_pagination", Page(fx, a.Limit, 0), fx)
	rows[0], rows[1] = rows[1], rows[0]
	if err := VerifyRows(c, "s16_page_offset", fx, a, rows); err == nil {
		t.Fatal("wrong ORDER BY accepted")
	}
}

func TestNumericConversionRejectsLossyAnswers(t *testing.T) {
	for _, input := range []any{1.5, math.NaN(), math.Inf(1), "3 garbage", math.Exp2(63)} {
		if _, _, err := toInt(input); err == nil {
			t.Fatalf("lossy integer accepted: %v", input)
		}
	}
}

func TestOrderedTopNCannotBeRepairedByCanonicalization(t *testing.T) {
	c, fx, a := load(t), Teaching(), args()
	rows := rowsFromOracle(t, "sel_topn", TopN(fx, a.N), fx)
	rows[0], rows[1] = rows[1], rows[0]
	if err := VerifyRows(c, "s14_topn_rownumber", fx, a, rows); err == nil {
		t.Fatal("wrong top-N order accepted")
	}
}

func TestLiveControlsCannotPassOnMatchingAnswersOrWrongShapes(t *testing.T) {
	c, fx, a := load(t), Teaching(), args()
	for name, group := range controlContracts {
		want, _ := oracleFor(group, fx, a)
		rows := rowsFromOracle(t, group, want, fx)
		if err := DetectControlRows(c, name, fx, a, rows); err == nil {
			t.Fatalf("matching control accepted: %s", name)
		}
		if err := DetectControlRows(c, name, fx, a, [][]any{{"bad shape"}}); err == nil {
			t.Fatalf("decode failure counted as detection: %s", name)
		}
	}
}

func TestKeysetCursorBindsTieDirectionAndNullState(t *testing.T) {
	c, fx := load(t), Teaching()
	last := PageRow{OrderID: 11, Amount: ip(200)}
	a := Args{Pred: AmountGreaterThan(100), Limit: 2, Last: &last}
	s := c.Statements["s17_page_keyset"]
	bound, _, err := Bind(s.SQL, a.ParamMap(s))
	if err != nil || !strings.Contains(bound, "id > 11") {
		t.Fatalf("tie cursor: %s %v", bound, err)
	}
	want, _ := ExpectedAnswer(c, s.Name, fx, a)
	if !answersEqual(want, KeysetPage(fx, 2, &last)) {
		t.Fatal("cursor oracle compared with first page")
	}
	last.Amount = nil
	if a.ParamMap(s)["last_is_null"] != "1" {
		t.Fatal("NULL cursor not represented")
	}
}

func TestRunCellsMeasuresDrainedSamplesAndReportsErrors(t *testing.T) {
	c, fx, a := load(t), Teaching(), args()
	a.Trials = 3
	f := newFakeQ(t, c, fx, a)
	cells := RunCells(context.Background(), f, c, fx, a, []string{"s01_join_where"}, "run", "single")
	if !cells[0].TimingMeasured || len(cells[0].SamplesMS) != 3 || cells[0].ResultRows == 0 || cells[0].ResultJSONBytes == 0 {
		t.Fatalf("missing measurements: %+v", cells[0])
	}
	f.queryErr["s01_join_where"] = errorsForCheck{}
	cells = RunCells(context.Background(), f, c, fx, a, []string{"s01_join_where"}, "run", "single")
	if cells[0].Status != StatusError || cells[0].TimingMeasured || len(cells[0].SamplesMS) != 0 {
		t.Fatal("database error labelled timeout or timed")
	}
}

type errorsForCheck struct{}

func (errorsForCheck) Error() string { return "connection lost" }

type lateWrongQuerier struct {
	*fakeQ
	reads int
}

func (q *lateWrongQuerier) Query(ctx context.Context, sql string) ([][]any, error) {
	q.reads++
	rows, err := q.fakeQ.Query(ctx, sql)
	if q.reads == 4 && len(rows) > 0 {
		rows = rows[:len(rows)-1]
	}
	return rows, err
}

func TestLaterWrongAnswerDiscardsEarlierSamples(t *testing.T) {
	c, fx, a := load(t), Teaching(), args()
	a.Trials = 3
	q := &lateWrongQuerier{fakeQ: newFakeQ(t, c, fx, a)}
	cells := RunCells(context.Background(), q, c, fx, a, []string{"s01_join_where"}, "run", "single")
	if q.reads != 4 || cells[0].Status != StatusInvalid || cells[0].TimingMeasured || len(cells[0].SamplesMS) != 0 {
		t.Fatalf("late wrong answer retained timing: %+v", cells[0])
	}
}

func TestUntimedReportShowsNA(t *testing.T) {
	out, err := RenderComparison(ReportInput{Cells: []Cell{{Statement: "validation", Status: StatusValid}}})
	if err != nil || !strings.Contains(out, "N/A") {
		t.Fatalf("untimed zero reported as measurement: %s %v", out, err)
	}
}

func TestAliasesDoNotChangePhysicalStructure(t *testing.T) {
	if !SameStructure(seqScan, strings.ReplaceAll(seqScan, "orders o", "orders child")) {
		t.Fatal("alias renamed the physical plan")
	}
}

type missingPlanQuerier struct{ *fakeQ }

func (q missingPlanQuerier) Explain(context.Context, string) (string, error) { return "", nil }
func TestMissingPlanCannotProduceTiming(t *testing.T) {
	c, fx, a := load(t), Teaching(), args()
	cells := RunCells(context.Background(), missingPlanQuerier{newFakeQ(t, c, fx, a)}, c, fx, a, []string{"s01_join_where"}, "run", "single")
	if cells[0].Status != StatusError || cells[0].TimingMeasured {
		t.Fatal("missing plan admitted timing")
	}
}
