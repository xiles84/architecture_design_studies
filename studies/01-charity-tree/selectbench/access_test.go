package selectbench

import (
	"context"
	"reflect"
	"strings"
	"testing"
)

func bigFixture() Fixture {
	// > 10,000 rows so every protocol page depth addresses a real page.
	return Generate(42, GenSpec{Customers: 40, OrdersPerParent: 300, NullAmountEvery: 23, NullFKEvery: 311})
}

func TestDefaultShapesAreDisjointAndNonEmpty(t *testing.T) {
	// Disjointness is a property of the arm thresholds and must hold on every
	// fixture. Arm-A non-emptiness (needed for the overlap control to fire) is
	// checked on the measurement fixture: the tiny teaching fixture has no row
	// under the selective threshold, so it cannot be the non-emptiness witness.
	big := bigFixture()
	for _, shape := range DefaultPredicateShapes() {
		for _, fx := range []Fixture{Teaching(), big} {
			if n := OverlapCount(fx, shape.A, shape.B); n != 0 {
				t.Errorf("shape %s: arms overlap on %d rows; UNION ALL would not be a valid spelling", shape.Name, n)
			}
		}
		if a := countTrue(big, shape.A); a == 0 {
			t.Errorf("shape %s: arm A matches nothing on the measurement fixture, so the overlap control cannot fire", shape.Name)
		}
	}
}

func countTrue(fx Fixture, pred Pred) int {
	if pred == nil {
		return 0
	}
	n := 0
	for _, o := range fx.Orders {
		if pred(o) {
			n++
		}
	}
	return n
}

func TestX05OverlapControlIsDetected(t *testing.T) {
	c := load(t)
	fx := bigFixture()
	var shape PredShape
	for _, s := range DefaultPredicateShapes() {
		if s.Name == "broad/indexable" {
			shape = s
		}
	}
	a := Args{PredA: shape.A, PredB: shape.B, PredASQL: shape.ASQL, PredBSQL: shape.BSQL}
	got := Evaluators["x05_union_all_overlap"](fx, a)
	bag, ok := got.([]PredAltRow)
	if !ok {
		t.Fatalf("x05 evaluator returned %T", got)
	}
	// x05 duplicates arm A only, so its bag is twice the arm-A matches; it must
	// differ from the two-arm OR oracle, which DetectControl proves.
	armA := PredicateAlt(fx, shape.A, nil)
	if len(bag) != 2*len(armA) {
		t.Fatalf("x05 did not duplicate the arm-A rows: got %d rows, arm A has %d", len(bag), len(armA))
	}
	if len(armA) == 0 {
		t.Fatal("fixture gives arm A no rows; the control cannot fire")
	}
	if err := DetectControl(c, "x05_union_all_overlap", fx, a, got, "sel_predicate_alt"); err != nil {
		t.Errorf("overlap control was not detected: %v", err)
	}
}

func TestPaginationKeysetPairsMatchOffsetAtEveryProtocolDepth(t *testing.T) {
	fx := bigFixture()
	const limit = 5
	for _, d := range DefaultPageDepths() {
		offset := Page(fx, limit, d)
		keyset := KeysetPage(fx, limit, cursorAt(fx, d-1))
		if !reflect.DeepEqual(offset, keyset) {
			t.Errorf("depth %d: keyset page differs from offset page:\n keyset %#v\n offset %#v", d, keyset, offset)
		}
	}
}

func TestBuildAccessSpecsHasDistinctAxesAndEveryOracle(t *testing.T) {
	c := load(t)
	fx := bigFixture()
	base := Args{Pred: AmountGreaterThan(100), PredSQL: "amount > 100"}
	specs := BuildAccessSpecs(fx, base, DefaultPredicateShapes(), DefaultPageDepths(), 5, 3)
	if len(specs) != 6*3+2*5 {
		t.Fatalf("specs = %d, want %d", len(specs), 6*3+2*5)
	}
	seen := map[string]bool{}
	for _, s := range specs {
		if seen[s.Axis] {
			t.Errorf("duplicate axis %q", s.Axis)
		}
		seen[s.Axis] = true
		if _, err := ExpectedAnswer(c, s.Statement, fx, s.Args); err != nil {
			t.Errorf("%s [%s]: %v", s.Statement, s.Axis, err)
		}
		if _, _, err := Bind(c.Statements[s.Statement].SQL, s.Args.ParamMap(c.Statements[s.Statement])); err != nil {
			t.Errorf("%s [%s] binding: %v", s.Statement, s.Axis, err)
		}
	}
	for _, wanted := range []string{
		"shape=selective/indexable;spelling=or",
		"shape=broad/expression;spelling=union_all",
		"shape=selective/cast;spelling=union",
		"shape=broad/cast;spelling=or",
		"depth=10000;spelling=keyset",
	} {
		if !seen[wanted] {
			t.Errorf("matrix is missing axis %q", wanted)
		}
	}
}

// matrixFakeQ returns the oracle rows for the exact bound statement/args, so the
// whole matrix can be exercised through the real cell loop without an engine.
type matrixFakeQ struct {
	rowsByBound map[string][][]any
}

func (q matrixFakeQ) Exec(context.Context, string) error { return nil }
func (q matrixFakeQ) Explain(context.Context, string) (string, error) {
	return "Limit\n  ->  Index Scan using orders_amount_idx on orders o\n", nil
}
func (q matrixFakeQ) Query(_ context.Context, sql string) ([][]any, error) {
	if rows, ok := q.rowsByBound[sql]; ok {
		return rows, nil
	}
	return nil, nil
}

func TestRunAccessMatrixMarksEveryCellValid(t *testing.T) {
	c := load(t)
	fx := bigFixture()
	base := Args{Pred: AmountGreaterThan(100), PredSQL: "amount > 100"}
	specs := BuildAccessSpecs(fx, base, DefaultPredicateShapes(), DefaultPageDepths(), 5, 3)

	q := matrixFakeQ{rowsByBound: map[string][][]any{}}
	for _, s := range specs {
		st := c.Statements[s.Statement]
		bound, _, err := Bind(st.SQL, s.Args.ParamMap(st))
		if err != nil {
			t.Fatal(err)
		}
		want, err := ExpectedAnswer(c, s.Statement, fx, s.Args)
		if err != nil {
			t.Fatal(err)
		}
		q.rowsByBound[bound] = rowsForGroup(st.Group, want, fx)
	}
	cells := RunSpecs(context.Background(), q, c, fx, specs, "run-test", "pg-single")
	if len(cells) != len(specs) {
		t.Fatalf("cells = %d, specs = %d", len(cells), len(specs))
	}
	for _, cell := range cells {
		if cell.Status != StatusValid {
			t.Errorf("%s [%s]: %s (%s)", cell.Statement, cell.Axis, cell.Status, cell.Note)
		}
		if !cell.TimingMeasured || len(cell.SamplesMS) != 3 {
			t.Errorf("%s [%s]: timing not measured over 3 trials: %+v", cell.Statement, cell.Axis, cell)
		}
	}
}

// rowsForGroup is rowsFromCanonical without the testing handle.
func rowsForGroup(group string, want any, fx Fixture) [][]any {
	return rowsFromCanonical(group, want, fx)
}

func TestPositionControlPairsFirstCellWithItsRepeat(t *testing.T) {
	shapes := DefaultPredicateShapes()
	base := Args{}
	specs := BuildAccessSpecs(bigFixture(), base, shapes, DefaultPageDepths(), 5, 0)
	specs = append(specs, PositionControlSpec(base, shapes[0], 0))

	names := make([]string, 0, len(shapes))
	for _, s := range shapes {
		names = append(names, s.Name)
	}
	r := AccessResult{Shapes: names, Depths: DefaultPageDepths()}
	for _, s := range specs {
		r.Cells = append(r.Cells, Cell{Statement: s.Statement, Axis: s.Axis,
			Status: StatusValid, TimingMeasured: true, MedianMS: 2, Trials: 1})
	}
	pp, ok := r.PositionPair()
	if !ok {
		t.Fatal("position pair not found")
	}
	if want := "shape=" + shapes[0].Name + ";spelling=or"; pp.A.Axis != want {
		t.Errorf("first axis = %q, want %q", pp.A.Axis, want)
	}
	if want := "shape=" + shapes[0].Name + ";spelling=or;position=last"; pp.B.Axis != want {
		t.Errorf("last axis = %q, want %q", pp.B.Axis, want)
	}
}

func TestAccessResultPairsCoverShapesAndDepths(t *testing.T) {
	shapes := DefaultPredicateShapes()
	axes := BuildAccessSpecs(bigFixture(), Args{}, shapes, DefaultPageDepths(), 5, 0)
	names := make([]string, 0, len(shapes))
	for _, s := range shapes {
		names = append(names, s.Name)
	}
	r := AccessResult{Shapes: names, Depths: DefaultPageDepths()}
	for _, s := range axes {
		r.Cells = append(r.Cells, Cell{Statement: s.Statement, Axis: s.Axis,
			Status: StatusValid, TimingMeasured: true, MedianMS: 1, Trials: 1})
	}
	pairs := r.Pairs()
	if len(pairs) != len(shapes)*2+5 {
		t.Fatalf("pairs = %d, want %d", len(pairs), len(shapes)*2+5)
	}
	for _, p := range pairs {
		if p.Explanation == "" {
			t.Errorf("pair %q has no explanation", p.Label)
		}
	}
}

func TestControlContractsCoverEveryControlStatement(t *testing.T) {
	c := load(t)
	controls := c.InGroup(ControlGroup)
	if len(controls) != len(controlContracts) {
		t.Fatalf("control statements = %d, controlContracts = %d", len(controls), len(controlContracts))
	}
	for _, name := range controls {
		if _, ok := controlContracts[name]; !ok {
			t.Errorf("control %s has no pretended contract", name)
		}
	}
	for name, group := range controlContracts {
		if _, ok := c.Statements[name]; !ok {
			t.Errorf("controlContracts names unknown statement %s", name)
		}
		if strings.TrimSpace(group) == "" {
			t.Errorf("control %s has an empty pretended contract", name)
		}
	}
}
