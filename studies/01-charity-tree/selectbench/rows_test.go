package selectbench

import "testing"

// rowsFromOracle turns a canonical oracle value into the raw rows a database
// would return for the same contract, so the mapping can be checked both ways.
func rowsFromOracle(t *testing.T, group string, want any, fx Fixture) [][]any {
	t.Helper()
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
	case "sel_aggregate", "sel_aggregate_sel":
		v := want.([]AggRow)
		out := make([][]any, 0, len(v))
		for _, r := range v {
			out = append(out, []any{r.CustomerID, r.CustomerName, r.Total, r.Children})
		}
		return out
	case "sel_topn", "sel_topn_sel", "sel_topn_one":
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
	}
	t.Fatalf("no row builder for group %s", group)
	return nil
}

func nv(p *int) any {
	if p == nil {
		return nil
	}
	return *p
}

// The live path must accept exactly what the oracle expects, for every contract.
func TestVerifyRowsAcceptsOracleRows(t *testing.T) {
	c := load(t)
	fx := Teaching()
	a := args()
	for _, name := range c.Order {
		s := c.Statements[name]
		if s.Group == ControlGroup {
			continue
		}
		want, err := ExpectedAnswer(c, name, fx, a)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		rows := rowsFromOracle(t, s.Group, want, fx)
		if err := VerifyRows(c, name, fx, a, rows); err != nil {
			t.Errorf("%s (%s): %v", name, s.Group, err)
		}
	}
}

func TestVerifyRowsRejectsADifferentAnswer(t *testing.T) {
	c := load(t)
	fx := Teaching()
	a := args()
	full := rowsFromOracle(t, "sel_full_rows", FullRows(fx, a.Pred), fx)
	if len(full) < 2 {
		t.Fatal("fixture too small for this test")
	}
	// Drop a row: a subset is a different answer.
	if err := VerifyRows(c, "s01_join_where", fx, a, full[:len(full)-1]); err == nil {
		t.Error("a truncated result set was accepted")
	}
	// Wrong column shape (a membership result fed to a full-row statement).
	mem := rowsFromOracle(t, "sel_membership", Membership(fx, a.Pred), fx)
	if err := VerifyRows(c, "s01_join_where", fx, a, mem); err == nil {
		t.Error("a wrong column count was accepted")
	}
}

func TestVerifyRowsRefusesControls(t *testing.T) {
	c := load(t)
	fx := Teaching()
	if err := VerifyRows(c, "x01_optional_where", fx, args(), nil); err == nil {
		t.Error("a control was verified instead of being refused")
	}
}

// A NULL amount must survive the round trip as NULL, not as zero.
func TestCanonicalizePreservesNulls(t *testing.T) {
	got, err := Canonicalize("sel_pagination", [][]any{{13, nil}, {10, 50}})
	if err != nil {
		t.Fatal(err)
	}
	v := got.([]PageRow)
	if v[0].OrderID != 13 || v[1].OrderID != 10 {
		t.Fatalf("order wrong: %#v", v)
	}
	if v[0].Amount != nil {
		t.Errorf("order 13's NULL amount became %d", *v[0].Amount)
	}
}
