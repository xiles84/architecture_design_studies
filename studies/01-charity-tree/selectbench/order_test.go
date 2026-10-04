package selectbench

import (
	"context"
	"reflect"
	"testing"
)

// Tests for the order-correction machinery: a seeded per-run order and
// interleaved sample rounds.

func TestPlanOrderIsAPermutationAndDeterministic(t *testing.T) {
	for _, seed := range []int64{1, 2, 7, 42, 999983} {
		p := PlanOrder(seed, 12)
		if len(p.Order) != 12 {
			t.Fatalf("seed %d: order length %d", seed, len(p.Order))
		}
		seen := map[int]bool{}
		for _, i := range p.Order {
			if i < 0 || i >= 12 || seen[i] {
				t.Fatalf("seed %d: %v is not a permutation", seed, p.Order)
			}
			seen[i] = true
		}
		if !reflect.DeepEqual(p, PlanOrder(seed, 12)) {
			t.Fatalf("seed %d is not deterministic", seed)
		}
	}
	if reflect.DeepEqual(PlanOrder(1, 12).Order, PlanOrder(2, 12).Order) {
		t.Error("different seeds produced the same order")
	}
}

func TestOrderSeedDependsOnRunID(t *testing.T) {
	if OrderSeedFor(42, "select-x-1") == OrderSeedFor(42, "select-x-2") {
		t.Error("different run ids produced the same order seed")
	}
	if OrderSeedFor(42, "select-x-1") != OrderSeedFor(42, "select-x-1") {
		t.Error("order seed is not deterministic")
	}
}

func TestPlanSequenceReportsTheExecutedAxesInOrder(t *testing.T) {
	specs := []CellSpec{{Axis: "a"}, {Axis: "b"}, {Axis: "c"}, {Axis: "d"}}
	plan := PlanOrder(5, len(specs))
	seq := plan.Sequence(specs)
	if len(seq) != len(specs) {
		t.Fatalf("sequence length %d", len(seq))
	}
	for i, axis := range seq {
		if axis != specs[plan.Order[i]].Axis {
			t.Errorf("position %d: %q, want %q", i, axis, specs[plan.Order[i]].Axis)
		}
	}
}

func TestInterleavedMatchesSequentialStatusAndTrialCounts(t *testing.T) {
	c := load(t)
	fx := Teaching()
	a := Args{Pred: AmountGreaterThan(100), PredSQL: "amount > 100", N: 1, Limit: 2, Trials: 4}
	names := []string{"s01_join_where", "s06_exists_parent", "s09_not_exists", "s16_page_offset"}
	var specs []CellSpec
	for _, n := range names {
		specs = append(specs, CellSpec{Statement: n, Axis: "axis=" + n, Args: a})
	}
	seq := RunSpecs(context.Background(), newFakeQ(t, c, fx, a), c, fx, specs, "run", "single")
	plan := PlanOrder(11, len(specs))
	inter := RunSpecsInterleaved(context.Background(), newFakeQ(t, c, fx, a), c, fx, specs, plan, "run", "single")
	if len(inter) != len(specs) {
		t.Fatalf("interleaved returned %d cells, want %d", len(inter), len(specs))
	}
	byAxis := map[string]Cell{}
	for _, cell := range seq {
		byAxis[cell.Axis] = cell
	}
	for _, cell := range inter {
		want, ok := byAxis[cell.Axis]
		if !ok {
			t.Fatalf("unexpected axis %q", cell.Axis)
		}
		if cell.Status != want.Status || cell.Trials != want.Trials || cell.ResultRows != want.ResultRows {
			t.Errorf("%s: interleaved status=%s trials=%d rows=%d; sequential status=%s trials=%d rows=%d",
				cell.Axis, cell.Status, cell.Trials, cell.ResultRows, want.Status, want.Trials, want.ResultRows)
		}
	}
}

func TestInterleavedSkipsCellsWithNoTrials(t *testing.T) {
	c := load(t)
	fx := Teaching()
	a := Args{Pred: AmountGreaterThan(100), PredSQL: "amount > 100", N: 1, Limit: 2, Trials: 0}
	specs := []CellSpec{{Statement: "s01_join_where", Axis: "axis=a", Args: a}}
	plan := PlanOrder(3, 1)
	cells := RunSpecsInterleaved(context.Background(), newFakeQ(t, c, fx, a), c, fx, specs, plan, "run", "single")
	if len(cells) != 1 || cells[0].Trials != 1 || !cells[0].TimingMeasured {
		t.Fatalf("trials=0 did not fall back to one sample: %+v", cells[0])
	}
}
