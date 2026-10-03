package selectbench

import (
	"reflect"
	"sort"
	"testing"
)

func load(t *testing.T) *Catalog {
	t.Helper()
	c, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	return c
}

func TestCatalogLoadsEveryStatementWithAGroup(t *testing.T) {
	c := load(t)
	if len(c.Order) < 15 {
		t.Fatalf("expected a populated catalogue, got %d statements", len(c.Order))
	}
	want := []string{"sel_full_rows", "sel_membership", "sel_absence", "sel_optional_match",
		"sel_aggregate", "sel_topn", "sel_pagination", ControlGroup}
	for _, g := range want {
		if len(c.InGroup(g)) == 0 {
			t.Errorf("semantic group %q is empty or missing", g)
		}
	}
	for _, n := range c.Order {
		s := c.Statements[n]
		if s.Group == "" {
			t.Errorf("%s: no group", n)
		}
		if s.SQL == "" {
			t.Errorf("%s: no SQL body", n)
		}
	}
}

func TestEveryStatementHasAnEvaluator(t *testing.T) {
	c := load(t)
	for _, n := range c.Order {
		if _, ok := Evaluators[n]; !ok {
			t.Errorf("%s has no Go evaluator: a template without a contract check cannot be measured", n)
		}
	}
}

// TestGroupGuardRefusesCrossGroupAndControls is the guard against pooling a
// faster statement that answers a different question.
func TestGroupGuardRefusesCrossGroupAndControls(t *testing.T) {
	c := load(t)
	cases := []struct {
		a, b string
		want bool
	}{
		{"s01_join_where", "s02_filtered_derived", true},     // same group
		{"s01_join_where", "s06_exists_parent", false},       // full rows vs membership
		{"s01_join_where", "x01_optional_where", false},      // template vs control
		{"x01_optional_where", "x02_not_in_nullable", false}, // control vs control
	}
	for _, tc := range cases {
		got, err := c.Comparable(tc.a, tc.b)
		if got != tc.want {
			t.Errorf("Comparable(%s,%s)=%v, want %v (err=%v)", tc.a, tc.b, got, tc.want, err)
		}
		if !tc.want && err == nil {
			t.Errorf("Comparable(%s,%s) was refused without an explaining error", tc.a, tc.b)
		}
	}
	if _, err := c.Comparable("no_such_statement", "s01_join_where"); err == nil {
		t.Error("an unknown statement name did not error")
	}
}

// TestGroupsAgreeWithOracles checks every valid template against the
// independent Go oracle for its declared contract.
func TestGroupsAgreeWithOracles(t *testing.T) {
	c := load(t)
	fx := Teaching()
	pred := AmountGreaterThan(100)
	a := Args{Pred: pred, N: 1, Limit: 3}

	checks := []struct {
		group string
		want  any
	}{
		{"sel_full_rows", FullRows(fx, pred)},
		{"sel_membership", Membership(fx, pred)},
		{"sel_absence", Absence(fx, pred)},
		{"sel_optional_match", OptionalMatch(fx, pred)},
		{"sel_aggregate", Aggregate(fx)},
		{"sel_topn", TopN(fx, 1)},
		{"sel_pagination", Page(fx, 3, 0)},
	}
	for _, ch := range checks {
		names := c.InGroup(ch.group)
		if len(names) == 0 {
			t.Errorf("group %s is empty", ch.group)
			continue
		}
		for _, n := range names {
			got := Evaluators[n](fx, a)
			if !reflect.DeepEqual(got, ch.want) {
				t.Errorf("%s (%s) disagrees with its oracle:\n got %#v\nwant %#v", n, ch.group, got, ch.want)
			}
		}
	}
}

// TestContractsAreDistinctOnTheFixture shows that the contracts really do
// differ, which is why they may never be pooled into one ranking.
func TestContractsAreDistinctOnTheFixture(t *testing.T) {
	fx := Teaching()
	pred := AmountGreaterThan(100)
	if len(FullRows(fx, pred)) == len(Membership(fx, pred)) {
		t.Error("fixture does not exercise row multiplication: full rows and membership have equal size")
	}
	if len(Membership(fx, pred)) == len(OptionalMatch(fx, pred)) {
		t.Error("fixture does not exercise unmatched parents: membership and optional match have equal size")
	}
}

func TestPaginationKeysetMatchesOffsetPage(t *testing.T) {
	fx := Teaching()
	p1 := Page(fx, 3, 0)
	if len(p1) == 0 {
		t.Fatal("empty first page")
	}
	last := p1[len(p1)-1]
	offsetPage := Page(fx, 3, 3)
	keysetPage := KeysetPage(fx, 3, &last)
	if !reflect.DeepEqual(offsetPage, keysetPage) {
		t.Errorf("keyset page differs from the offset page under the fixed order:\n keyset %#v\n offset %#v", keysetPage, offsetPage)
	}
}

// TestControlsAreDetected is the read-path correctness control: each
// deliberately wrong fixture must produce an answer that differs from the
// contract oracle.
func TestControlsAreDetected(t *testing.T) {
	fx := Teaching()
	pred := AmountGreaterThan(100)
	a := Args{Pred: pred, N: 1, Limit: 2}

	if reflect.DeepEqual(Evaluators["x01_optional_where"](fx, a), OptionalMatch(fx, pred)) {
		t.Error("x01_optional_where was not detected: the WHERE variant matched the optional-match contract")
	}
	if reflect.DeepEqual(Evaluators["x02_not_in_nullable"](fx, a), Absence(fx, pred)) {
		t.Error("x02_not_in_nullable was not detected: nullable NOT IN matched NOT EXISTS")
	}
	// The membership contract is one row per parent IDENTITY; the control
	// de-duplicates by name and silently merges two identities.
	byIdentity := namesPerIdentity(fx, Membership(fx, pred))
	if reflect.DeepEqual(Evaluators["x03_distinct_on_name"](fx, a), byIdentity) {
		t.Error("x03_distinct_on_name was not detected: name-level DISTINCT matched parent identity")
	}
	if reflect.DeepEqual(Evaluators["x04_limit_before_join"](fx, a), FullRows(fx, pred)) {
		t.Error("x04_limit_before_join was not detected: limiting before the join matched the full bag")
	}
}

func namesPerIdentity(fx Fixture, ids []int) []string {
	names := byName(fx)
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		out = append(out, names[id])
	}
	sort.Strings(out)
	return out
}

// TestFixtureExercisesEveryRequiredEdgeCase guards the fixture itself: a test
// that passes on data lacking the edge case proves nothing.
func TestFixtureExercisesEveryRequiredEdgeCase(t *testing.T) {
	fx := Teaching()
	pred := AmountGreaterThan(100)

	var nullAmount, nullFK bool
	for _, o := range fx.Orders {
		if o.Amount == nil {
			nullAmount = true
		}
		if o.CustomerID == nil {
			nullFK = true
		}
	}
	if !nullAmount {
		t.Error("fixture has no NULL amount")
	}
	if !nullFK {
		t.Error("fixture has no NULL child foreign key (needed for the NOT IN control)")
	}

	// duplicate qualifying children for one parent, and a tie between them
	byParent := map[int]int{}
	for _, o := range fx.Orders {
		if o.CustomerID != nil && pred(o) {
			byParent[*o.CustomerID]++
		}
	}
	dup := false
	for _, n := range byParent {
		if n > 1 {
			dup = true
		}
	}
	if !dup {
		t.Error("fixture has no parent with duplicate qualifying children")
	}

	// unmatched parent
	if len(Absence(fx, pred)) == 0 {
		t.Error("fixture has no unmatched parent")
	}
	// all-miss predicate
	if n := len(Membership(fx, AmountGreaterThan(1000))); n != 0 {
		t.Errorf("expected an empty set for amount>1000, got %d", n)
	}
}
