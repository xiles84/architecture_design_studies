package selectbench

import (
	"reflect"
	"strings"
	"testing"
)

// Pure tests for the aggregate / top-N sweep family. No database is needed: the
// evaluators and the contract oracles are checked against each other, and the
// two deliberately-wrong controls are shown to produce a different answer.

func sweepFixture() Fixture {
	return GuaranteedTie(Generate(42, GenSpec{Customers: 60, OrdersPerParent: 4, NullAmountEvery: 11, NullFKEvery: 199}))
}

func TestGuaranteedTieCreatesATieAtTheSmallestMultiChildParent(t *testing.T) {
	fx := GuaranteedTie(Teaching())
	counts := map[int]int{}
	top := map[int]int{}
	for _, o := range fx.Orders {
		if o.CustomerID == nil || o.Amount == nil {
			continue
		}
		counts[*o.CustomerID]++
		if top[*o.CustomerID] == 0 || *o.Amount > top[*o.CustomerID] {
			top[*o.CustomerID] = *o.Amount
		}
	}
	if counts[1] < 2 {
		t.Fatalf("fixture customer 1 has %d orders", counts[1])
	}
	// The tie must be AT THE TOP so RANK keeps an extra row at rank 1.
	atMax := 0
	for _, o := range fx.Orders {
		if o.CustomerID != nil && *o.CustomerID == 1 && o.Amount != nil && *o.Amount == top[1] {
			atMax++
		}
	}
	if atMax < 2 {
		t.Fatalf("GuaranteedTie left only %d orders at customer 1's maximum", atMax)
	}
}

func TestDoubleCountAggregateShowsTheWrongGrain(t *testing.T) {
	fx := sweepFixture()
	got := DoubleCountAggregate(fx, 0)
	want := AggregateUpTo(fx, 0)
	if reflect.DeepEqual(got, want) {
		t.Fatal("the wrong-grain control matched the aggregate oracle")
	}
	// The inflation must be the child count: a parent with two children has
	// twice the total.
	counts := map[int]int{}
	for _, o := range fx.Orders {
		if o.CustomerID != nil {
			counts[*o.CustomerID]++
		}
	}
	byID := map[int]AggRow{}
	for _, r := range want {
		byID[r.CustomerID] = r
	}
	for _, r := range got {
		w := byID[r.CustomerID]
		n := counts[r.CustomerID]
		if r.Total != w.Total*n || r.Children != w.Children*n {
			t.Fatalf("parent %d: double count %+v, expected base %+v x %d", r.CustomerID, r, w, n)
		}
	}
}

func TestRankTiesTopNKeepsBoundaryTies(t *testing.T) {
	fx := sweepFixture()
	got := RankTiesTopN(fx, 1, 0)
	want := TopNUpTo(fx, 1, 0)
	if reflect.DeepEqual(got, want) {
		t.Fatal("the RANK tie control matched the N=1 top-N oracle")
	}
	if len(got) <= len(want) {
		t.Fatalf("RANK did not keep extra tied rows: got %d, want more than %d", len(got), len(want))
	}
}

func TestAggregateTopNUptoRespectParentMax(t *testing.T) {
	fx := sweepFixture()
	for _, maxParent := range []int{1, 10, 60} {
		agg := AggregateUpTo(fx, maxParent)
		if len(agg) != maxParent {
			t.Errorf("AggregateUpTo(%d) returned %d parents", maxParent, len(agg))
		}
		for _, r := range agg {
			if r.CustomerID > maxParent {
				t.Errorf("AggregateUpTo(%d) included parent %d", maxParent, r.CustomerID)
			}
		}
		top := TopNUpTo(fx, 1, maxParent)
		for _, r := range top {
			if r.CustomerID > maxParent {
				t.Errorf("TopNUpTo(%d) included parent %d", maxParent, r.CustomerID)
			}
		}
		if got := len(AggregateUpTo(fx, 0)); got != len(fx.Customers) {
			t.Errorf("AggregateUpTo(0) returned %d rows for %d parents", got, len(fx.Customers))
		}
	}
}

func TestAggregateTopNStatementsAgreeWithTheirOracles(t *testing.T) {
	c := load(t)
	fx := sweepFixture()
	names := []string{
		"a01_aggregate_join_sel", "a02_aggregate_preagg_sel", "a03_aggregate_correlated_sel",
		"a04_topn_rownumber_sel", "a05_topn_lateral_sel", "a06_topn_distinct_on_sel",
	}
	for _, maxParent := range []int{0, 3, 60} {
		a := Args{N: 1, ParentMax: maxParent}
		for _, name := range names {
			if _, ok := Evaluators[name]; !ok {
				t.Fatalf("%s has no evaluator", name)
			}
			want, err := ExpectedAnswer(c, name, fx, a)
			if err != nil {
				t.Fatalf("%s: %v", name, err)
			}
			got := Evaluators[name](fx, a)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("%s at parent_max=%d disagrees with its oracle:\n got %#v\nwant %#v", name, maxParent, got, want)
			}
			if _, _, err := Bind(c.Statements[name].SQL, a.ParamMap(c.Statements[name])); err != nil {
				t.Errorf("%s binding: %v", name, err)
			}
		}
	}
}

func TestAggregateTopNControlsAreDetected(t *testing.T) {
	c := load(t)
	fx := sweepFixture()
	a := Args{N: 1}
	for name, group := range map[string]string{
		"x05_aggregate_double_count": "sel_aggregate_sel",
		"x06_topn_rank_ties":         "sel_topn_sel",
	} {
		got := Evaluators[name](fx, a)
		if err := DetectControl(c, name, fx, a, got, group); err != nil {
			t.Errorf("%v", err)
		}
	}
}

func TestBuildAggregateSpecsCoversEveryLevelAndSpelling(t *testing.T) {
	levels := ParentLevels()
	specs := BuildAggregateSpecs(Args{}, levels, 3)
	if want := (len(AggregateStatements) + len(TopNStatements)) * len(levels); len(specs) != want {
		t.Fatalf("specs = %d, want %d", len(specs), want)
	}
	seen := map[string]bool{}
	for _, s := range specs {
		if seen[s.Axis] {
			t.Errorf("duplicate axis %q", s.Axis)
		}
		seen[s.Axis] = true
		found := false
		for _, l := range levels {
			if s.Args.ParentMax == l {
				found = true
			}
		}
		if !found {
			t.Errorf("%s has parent_max %d, not a level", s.Axis, s.Args.ParentMax)
		}
		if s.Args.N != 1 || s.Args.Trials != 3 {
			t.Errorf("%s bound N=%d trials=%d", s.Axis, s.Args.N, s.Args.Trials)
		}
	}
	if len(AggregateControlSpecs(Args{}, 3)) != 2 {
		t.Fatal("expected the two aggregate controls")
	}
}

func TestAggregatePairsRecomputeBothDirections(t *testing.T) {
	r := AggResult{Levels: []int{10}}
	add := func(axis string, ms float64) {
		r.Cells = append(r.Cells, Cell{Axis: axis, Statement: "s", Status: StatusValid,
			TimingMeasured: true, MedianMS: ms})
	}
	add(AggregateTopNAxis("aggregate", "a01_aggregate_join_sel", 10), 4)
	add(AggregateTopNAxis("aggregate", "a02_aggregate_preagg_sel", 10), 1)
	add(AggregateTopNAxis("aggregate", "a03_aggregate_correlated_sel", 10), 2)
	pairs := r.AggregatePairs()
	if len(pairs) != 3 {
		t.Fatalf("aggregate pairs = %d, want 3", len(pairs))
	}
	byLabel := map[string]float64{}
	for _, p := range pairs {
		byLabel[p.Label] = p.RatioBA
		if p.Family != "aggregate" {
			t.Errorf("%s family = %q", p.Label, p.Family)
		}
	}
	if got := byLabel["pre_agg vs post_join"]; got != 0.25 {
		t.Errorf("pre/post = %v, want 0.25", got)
	}
	if got := byLabel["correlated vs pre_agg"]; got != 2 {
		t.Errorf("corr/pre = %v, want 2", got)
	}
}

func TestPositionControlSpecPairsWithTheFirstCell(t *testing.T) {
	levels := ParentLevels()
	pc := AggregatePositionControlSpec(Args{}, 3)
	if pc.Statement != AggregateStatements[0] || pc.Args.ParentMax != levels[0] || pc.Args.Trials != 3 {
		t.Fatalf("position spec = %+v", pc)
	}
	if !strings.HasSuffix(pc.Axis, PositionSuffix) {
		t.Fatalf("position axis %q lacks the suffix", pc.Axis)
	}
	r := AggResult{Levels: levels, Cells: []Cell{
		{Axis: AggregateTopNAxis("aggregate", AggregateStatements[0], levels[0]), Statement: AggregateStatements[0],
			Status: StatusValid, TimingMeasured: true, MedianMS: 1},
		{Axis: pc.Axis, Statement: AggregateStatements[0], Status: StatusValid, TimingMeasured: true, MedianMS: 2},
	}}
	pp, ok := r.PositionPair()
	if !ok {
		t.Fatal("position pair not found")
	}
	if pp.RatioBA != 2 {
		t.Fatalf("position ratio = %v, want 2", pp.RatioBA)
	}
	if _, ok := (AggResult{Levels: levels}).PositionPair(); ok {
		t.Fatal("position pair reported without the control cell")
	}
}
