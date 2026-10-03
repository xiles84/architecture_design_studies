package selectbench

import (
	"context"
	"errors"
	"strings"
	"testing"
)

// rowsFromCanonical is rowsFromOracle without the testing handle, so the runner
// tests can build the rows a database would return for an oracle answer.
func rowsFromCanonical(group string, want any, fx Fixture) [][]any {
	names := byName(fx)
	switch group {
	case "sel_full_rows":
		v := want.([]FullRow)
		out := make([][]any, 0, len(v))
		for _, r := range v {
			out = append(out, []any{r.CustomerID, r.CustomerName, r.OrderID, nv(r.Amount)})
		}
		return out
	case "sel_membership", "sel_absence":
		v := want.([]int)
		out := make([][]any, 0, len(v))
		for _, id := range v {
			out = append(out, []any{id, names[id]})
		}
		return out
	case "sel_optional_match":
		v := want.([]OptRow)
		out := make([][]any, 0, len(v))
		for _, r := range v {
			out = append(out, []any{r.CustomerID, r.CustomerName, nv(r.OrderID), nv(r.Amount)})
		}
		return out
	case "sel_aggregate":
		v := want.([]AggRow)
		out := make([][]any, 0, len(v))
		for _, r := range v {
			out = append(out, []any{r.CustomerID, r.CustomerName, r.Total, r.Children})
		}
		return out
	case "sel_topn":
		v := want.([]TopRow)
		out := make([][]any, 0, len(v))
		for _, r := range v {
			out = append(out, []any{r.CustomerID, r.OrderID, nv(r.Amount)})
		}
		return out
	case "sel_pagination":
		v := want.([]PageRow)
		out := make([][]any, 0, len(v))
		for _, r := range v {
			out = append(out, []any{r.OrderID, nv(r.Amount)})
		}
		return out
	case "sel_predicate_alt":
		v := want.([]PredAltRow)
		out := make([][]any, 0, len(v))
		for _, r := range v {
			out = append(out, []any{r.OrderID, nv(r.Amount)})
		}
		return out
	}
	return nil
}

type fakeQ struct {
	nameByBound  map[string]string // bound SQL -> statement name
	groupByName  map[string]string // statement name -> semantic group
	groupRows    map[string][][]any
	execs        []string
	override     map[string][][]any // statement name -> rows to return instead
	explainErr   map[string]error
	queryErr     map[string]error
	explainCalls int
}

func newFakeQ(t *testing.T, c *Catalog, fx Fixture, a Args) *fakeQ {
	t.Helper()
	f := &fakeQ{
		nameByBound: map[string]string{}, groupByName: map[string]string{}, groupRows: map[string][][]any{},
		override: map[string][][]any{}, explainErr: map[string]error{}, queryErr: map[string]error{},
	}
	for _, n := range c.Order {
		s := c.Statements[n]
		if s.Group == ControlGroup {
			continue
		}
		bound, _, err := Bind(s.SQL, a.ParamMap(s))
		if err != nil {
			t.Fatalf("%s: %v", n, err)
		}
		want, err := ExpectedAnswer(c, n, fx, a)
		if err != nil {
			t.Fatalf("%s: %v", n, err)
		}
		f.nameByBound[bound] = n
		f.groupByName[n] = s.Group
		if _, ok := f.groupRows[s.Group]; !ok {
			f.groupRows[s.Group] = rowsFromCanonical(s.Group, want, fx)
		}
	}
	return f
}

func (f *fakeQ) Exec(_ context.Context, sql string) error {
	f.execs = append(f.execs, sql)
	return nil
}

func (f *fakeQ) Explain(_ context.Context, sql string) (string, error) {
	f.explainCalls++
	name := f.nameByBound[sql]
	if err := f.explainErr[name]; err != nil {
		return "", err
	}
	return "Seq Scan on orders o  (cost=0.00..1.00 rows=1 width=8)\n  Filter: (amount > 100)\n", nil
}

func (f *fakeQ) Query(_ context.Context, sql string) ([][]any, error) {
	name := f.nameByBound[sql]
	if name == "" {
		return nil, errors.New("unexpected sql: " + sql)
	}
	if err := f.queryErr[name]; err != nil {
		return nil, err
	}
	if rows, ok := f.override[name]; ok {
		return rows, nil
	}
	return f.groupRows[f.groupByName[name]], nil
}

func TestPrepareLoadsSchemaFixtureAndStatistics(t *testing.T) {
	c := load(t)
	fx := Teaching()
	f := newFakeQ(t, c, fx, args())
	if err := Prepare(context.Background(), f, fx); err != nil {
		t.Fatal(err)
	}
	all := strings.Join(f.execs, "\n")
	for _, want := range []string{"CREATE SCHEMA IF NOT EXISTS selectbench",
		"INSERT INTO selectbench.customers", "INSERT INTO selectbench.orders",
		"ANALYZE selectbench.orders"} {
		if !strings.Contains(all, want) {
			t.Errorf("Prepare did not run %q", want)
		}
	}
}

func TestRunCellsMarksEveryTemplateValid(t *testing.T) {
	c := load(t)
	fx := Teaching()
	a := args()
	f := newFakeQ(t, c, fx, a)
	cells := RunCells(context.Background(), f, c, fx, a, nil, "run-1", "pg-single")
	if len(cells) == 0 {
		t.Fatal("no cells")
	}
	for _, cell := range cells {
		if cell.Status != StatusValid {
			t.Errorf("%s: status %s (%s)", cell.Statement, cell.Status, cell.Note)
		}
		if cell.Plan == "" {
			t.Errorf("%s: no plan captured", cell.Statement)
		}
		if cell.Group == "" {
			t.Errorf("%s: no group", cell.Statement)
		}
	}
}

func TestRunCellsMarksAWrongAnswerInvalid(t *testing.T) {
	c := load(t)
	fx := Teaching()
	a := args()
	f := newFakeQ(t, c, fx, a)
	// Drop a row from one statement's answer: a different answer, not a speed win.
	full := f.groupRows["sel_full_rows"]
	f.override["s01_join_where"] = full[:len(full)-1]

	cells := RunCells(context.Background(), f, c, fx, a, []string{"s01_join_where"}, "run-1", "pg-single")
	if len(cells) != 1 {
		t.Fatalf("cells = %d", len(cells))
	}
	if cells[0].Status != StatusInvalid {
		t.Errorf("a truncated answer was not marked invalid: %+v", cells[0])
	}
}

func TestRunCellsClassifiesUnsupportedAndTimeout(t *testing.T) {
	c := load(t)
	fx := Teaching()
	a := args()
	f := newFakeQ(t, c, fx, a)
	f.explainErr["s03_cte_filtered"] = ErrUnsupported
	f.queryErr["s14_topn_rownumber"] = context.DeadlineExceeded

	cells := RunCells(context.Background(), f, c, fx, a,
		[]string{"s03_cte_filtered", "s14_topn_rownumber"}, "run-1", "yb-single")
	got := map[string]Status{}
	for _, cell := range cells {
		got[cell.Statement] = cell.Status
	}
	if got["s03_cte_filtered"] != StatusUnsupported {
		t.Errorf("unsupported statement classified as %s", got["s03_cte_filtered"])
	}
	if got["s14_topn_rownumber"] != StatusTimeout {
		t.Errorf("timed-out statement classified as %s", got["s14_topn_rownumber"])
	}
}
