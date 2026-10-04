package selectbench

import (
	"strings"
	"testing"
)

func TestLoadJoinCatalogue(t *testing.T) {
	cat, err := LoadJoin()
	if err != nil {
		t.Fatalf("LoadJoin: %v", err)
	}
	for _, name := range JoinStatements {
		if _, ok := cat.Statements[name]; !ok {
			t.Errorf("missing bag statement %s", name)
		}
	}
	if g := cat.Statements["j11_optional_on"].Group; g != GroupJoinOptional {
		t.Errorf("j11_optional_on group = %q, want %q", g, GroupJoinOptional)
	}
	if g := cat.Statements["x01_optional_where"].Group; g != ControlGroup {
		t.Errorf("x01_optional_where group = %q, want control", g)
	}
	reuse := 0
	for _, n := range cat.Order {
		if cat.Statements[n].Group == GroupJoinReuse {
			reuse++
		}
	}
	wantReuse := len(reusePolicies) * len(ReuseLevels()) * len(ReuseConsumerPercents())
	if reuse != wantReuse {
		t.Errorf("reuse statements = %d, want %d", reuse, wantReuse)
	}
	controls := JoinControlNames(cat)
	if len(controls) != 2 {
		t.Errorf("join controls = %v, want 2", controls)
	}
}

func TestEveryJoinStatementHasOracle(t *testing.T) {
	cat, err := LoadJoin()
	if err != nil {
		t.Fatalf("LoadJoin: %v", err)
	}
	jf := JoinEdge()
	a := JoinArgs{Threshold: 3, Refs: 2, ConsumerMax: 10}
	for _, name := range cat.Order {
		s := cat.Statements[name]
		if s.Group == ControlGroup {
			continue
		}
		if _, err := JoinExpectedAnswer(cat, name, jf, a); err != nil {
			t.Errorf("%s (%s): %v", name, s.Group, err)
		}
	}
}

func TestJoinBagEdgeCases(t *testing.T) {
	jf := JoinEdge()
	all := JoinBag(jf, 5)
	// Facts 101,102 (customer 1), 103 (customer 2, NULL amount), 105 (customer
	// 4) qualify; the orphan 104 must not appear.
	if len(all) != 4 {
		t.Fatalf("bag(threshold=5) = %d rows, want 4: %#v", len(all), all)
	}
	dup := 0
	for _, r := range all {
		if r.CustomerID == 1 {
			dup++
		}
	}
	if dup != 2 {
		t.Errorf("customer 1 rows = %d, want 2 (duplicate matches are part of the bag)", dup)
	}
	if all[2].Amount != nil {
		t.Errorf("expected the NULL-amount fact to keep its NULL, got %v", all[2].Amount)
	}
	if one := JoinBag(jf, 1); len(one) != 1 || one[0].FactID != 101 {
		t.Errorf("bag(threshold=1) = %#v, want just fact 101", one)
	}
	if zero := JoinBag(jf, 0); len(zero) != 0 {
		t.Errorf("bag(threshold=0) = %d rows, want 0", len(zero))
	}
}

func TestJoinOptionalPreservesUnmatchedParent(t *testing.T) {
	jf := JoinEdge()
	opt := JoinOptional(jf, 1)
	// Customer 3 has no facts at all; customers 2 and 4 have no qualifying
	// fact at threshold 1; each must still appear once with NULLs.
	got := map[int]int{}
	nullRows := map[int]bool{}
	for _, r := range opt {
		got[r.CustomerID]++
		if r.FactID == nil {
			nullRows[r.CustomerID] = true
		}
	}
	for _, id := range []int{2, 3, 4} {
		if got[id] != 1 || !nullRows[id] {
			t.Errorf("customer %d: rows=%d null=%v, want one preserved NULL row", id, got[id], nullRows[id])
		}
	}
	// The inner bag is a different answer: x01 (predicate in WHERE) must fail
	// detection against the optional oracle.
	if answersEqual(opt, JoinBag(jf, 1)) {
		t.Fatal("optional-match oracle equals the inner bag; the WHERE control could not fire")
	}
	if err := JoinDetectControl(mustCatalog(t), "x01_optional_where", jf, JoinArgs{Threshold: 1}, canonicalRows(t, GroupJoinOptional, JoinBag(jf, 1))); err != nil {
		t.Errorf("x01 control not detected: %v", err)
	}
}

func TestJoinLimitBeforeControl(t *testing.T) {
	jf := JoinEdge()
	full := JoinBag(jf, 5)
	trunc := ControlJoinLimitBefore(jf, 5, 2)
	if answersEqual(full, trunc) {
		t.Fatal("limit-before-join control matched the bag oracle with limit 2")
	}
	if len(trunc) != 2 {
		t.Errorf("truncated control = %d rows, want 2", len(trunc))
	}
	// The SQL control uses a fixed limit of 7; on the 4-row edge fixture that
	// equals the bag, which is why the live gate runs it at 10% of the large
	// fixture instead. Here it must be seen to match, showing why.
	if !answersEqual(full, ControlJoinLimitBefore(jf, 5, 7)) {
		t.Fatal("edge fixture cannot distinguish a limit of 7; the gate must use the large fixture")
	}
	// Detection compares the observed answer with the contract oracle the
	// control pretends to serve, not with the control's own expected answer.
	cat := mustCatalog(t)
	rows := [][]any{}
	for _, r := range trunc {
		var amt any
		if r.Amount != nil {
			amt = *r.Amount
		}
		rows = append(rows, []any{r.CustomerID, r.CustomerName, r.FactID, amt})
	}
	if err := JoinDetectControl(cat, "x02_join_limit_before", jf, JoinArgs{Threshold: 5}, rows); err != nil {
		t.Errorf("x02 control not detected: %v", err)
	}
	// A control whose answer DID match the contracted bag must be reported as
	// not detected; otherwise the gate would pass on a control that never fires.
	if err := JoinDetectControl(cat, "x02_join_limit_before", jf, JoinArgs{Threshold: 5}, canonicalRows(t, GroupJoinBag, full)); err == nil {
		t.Error("control matched the bag oracle but was reported detected")
	}
}

func TestReuseAnswerDeterministic(t *testing.T) {
	jf := JoinEdge()
	a := ReuseAnswer(jf, 5, 4, 1000)
	if len(a) != 4 {
		t.Fatalf("reuse rows = %d, want 4", len(a))
	}
	for i, r := range a {
		if r.Ref != i+1 {
			t.Errorf("row %d ref = %d", i, r.Ref)
		}
		if r.Checksum != a[0].Checksum {
			t.Errorf("reference checksums differ: %d vs %d", r.Checksum, a[0].Checksum)
		}
	}
	narrow := ReuseAnswer(jf, 5, 1, 1)
	if len(narrow) != 1 {
		t.Fatalf("narrow reuse rows = %d, want 1", len(narrow))
	}
	if narrow[0].Checksum >= a[0].Checksum {
		t.Errorf("consumer filter did not reduce the consumed work: %d >= %d", narrow[0].Checksum, a[0].Checksum)
	}
}

func TestReusePolicySQLAndBinding(t *testing.T) {
	cat, err := LoadJoin()
	if err != nil {
		t.Fatalf("LoadJoin: %v", err)
	}
	for _, tc := range []struct {
		policy   string
		contains string
	}{{"mat", "AS MATERIALIZED ("}, {"notmat", "AS NOT MATERIALIZED ("}, {"default", "WITH e AS ("}} {
		name := ReuseStatementName(tc.policy, 4, 100)
		s, ok := cat.Statements[name]
		if !ok {
			t.Fatalf("missing %s", name)
		}
		if !strings.Contains(s.SQL, tc.contains) {
			t.Errorf("%s: SQL does not contain %q", name, tc.contains)
		}
		if got := strings.Count(s.SQL, "UNION ALL"); got != 3 {
			t.Errorf("%s: %d UNION ALL, want 3 for 4 references", name, got)
		}
		bound, used, err := Bind(s.SQL, JoinParamMap(s, JoinArgs{Threshold: 2000, ConsumerMax: 200000}))
		if err != nil {
			t.Fatalf("%s bind: %v", name, err)
		}
		if strings.Contains(bound, ":threshold") || strings.Contains(bound, ":consumer_max") {
			t.Errorf("%s: unbound placeholder remains: %s", name, bound)
		}
		if len(used) != 2 {
			t.Errorf("%s: used %v, want threshold+consumer_max", name, used)
		}
	}
}

func TestJoinMatrixShape(t *testing.T) {
	specs := BuildJoinSpecs(JoinArgs{}, 200000, 5)
	want := len(JoinStatements)*len(joinPercents) + len(reusePolicies)*len(ReuseLevels())*len(ReuseConsumerPercents())
	if len(specs) != want {
		t.Fatalf("matrix = %d cells, want %d", len(specs), want)
	}
	th := JoinThresholds(200000)
	if len(th) != 6 || th[0] != 20 || th[1] != 200 || th[2] != 2000 || th[3] != 20000 || th[4] != 100000 || th[5] != 180000 {
		t.Fatalf("thresholds on 200k = %v, want [20 200 2000 20000 100000 180000]", th)
	}
}

func mustCatalog(t *testing.T) *Catalog {
	t.Helper()
	cat, err := LoadJoin()
	if err != nil {
		t.Fatalf("LoadJoin: %v", err)
	}
	return cat
}

// canonicalRows turns an oracle answer into the scanned-row shape the live
// verifier sees, so the control-detection path is exercised without a database.
func canonicalRows(t *testing.T, group string, answer any) [][]any {
	t.Helper()
	switch group {
	case GroupJoinBag, GroupJoinOptional:
		rows := [][]any{}
		for _, r := range answer.([]JoinRow) {
			var amt any
			if r.Amount != nil {
				amt = *r.Amount
			}
			rows = append(rows, []any{r.CustomerID, r.CustomerName, r.FactID, amt})
		}
		return rows
	}
	t.Fatalf("no canonical row mapping for %s", group)
	return nil
}
