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
	fx := Teaching()
	for _, shape := range DefaultPredicateShapes() {
		if n := OverlapCount(fx, shape.A, shape.B); n != 0 {
			t.Errorf("shape %s: arms overlap on %d rows; UNION ALL would not be a valid spelling", shape.Name, n)
		}
		a, b := 0, 0
		for _, o := range fx.Orders {
			if shape.A(o) {
				a++
			}
			if shape.B(o) {
				b++
			}
		}
		if a == 0 {
			t.Errorf("shape %s: arm A matches nothing, so the overlap control cannot fire", shape.Name)
		}
		_ = b
	}
}

func TestX05OverlapControlIsDetected(t *testing.T) {
	c := load(t)
	fx := Teaching()
	shape := DefaultPredicateShapes()[0]
	a := Args{PredA: shape.A, PredB: shape.B, PredASQL: shape.ASQL, PredBSQL: shape.BSQL}
	got := Evaluators["x05_union_all_overlap"](fx, a)
	bag, ok := got.([]PredAltRow)
	if !ok {
		t.Fatalf("x05 evaluator returned %T", got)
	}
	if len(bag) != 2*len(PredicateAlt(fx, shape.A, shape.B)) {
		t.Fatalf("x05 did not duplicate the matching rows: got %d rows", len(bag))
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
	if len(specs) != 3*3+2*5 {
		t.Fatalf("specs = %d, want %d", len(specs), 3*3+2*5)
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
		"shape=indexable;spelling=or",
		"shape=expression;spelling=union_all",
		"shape=cast;spelling=union",
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

func TestAccessResultPairsCoverShapesAndDepths(t *testing.T) {
	axes := BuildAccessSpecs(bigFixture(), Args{}, DefaultPredicateShapes(), DefaultPageDepths(), 5, 0)
	r := AccessResult{Shapes: []string{"indexable", "expression", "cast"}, Depths: DefaultPageDepths()}
	for _, s := range axes {
		r.Cells = append(r.Cells, Cell{Statement: s.Statement, Axis: s.Axis,
			Status: StatusValid, TimingMeasured: true, MedianMS: 1, Trials: 1})
	}
	pairs := r.Pairs()
	if len(pairs) != 3*2+5 {
		t.Fatalf("pairs = %d, want %d", len(pairs), 3*2+5)
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
