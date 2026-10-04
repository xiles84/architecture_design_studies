package selectbench

import (
	"reflect"
	"strings"
	"testing"
)

// Pure tests for the membership/absence family (no database needed).

func existenceFixture() Fixture {
	return ExistenceShape(Generate(42, GenSpec{Customers: 40, OrdersPerParent: 5, NullAmountEvery: 11, NullFKEvery: 97}))
}

func TestMembershipWithMultiplicityDiffersFromTheSet(t *testing.T) {
	fx := existenceFixture()
	pred := AmountGreaterThan(100)
	set := Membership(fx, pred)
	mult := MembershipWithMultiplicity(fx, pred)
	if reflect.DeepEqual(set, mult) {
		t.Fatal("the multiplier control matched the membership set")
	}
	if len(mult) <= len(set) {
		t.Fatalf("multiplier returned %d rows for %d parents", len(mult), len(set))
	}
}

func TestExistenceStatementsAgreeWithTheirOracles(t *testing.T) {
	c := load(t)
	fx := existenceFixture()
	a := Args{Pred: AmountGreaterThan(100), PredSQL: "amount > 100"}
	for _, name := range []string{"e12_in_notnull", "e13_not_in_notnull", "x07_membership_multiplier"} {
		if _, ok := Evaluators[name]; !ok {
			t.Fatalf("%s has no evaluator", name)
		}
		if _, _, err := Bind(c.Statements[name].SQL, a.ParamMap(c.Statements[name])); err != nil {
			t.Errorf("%s binding: %v", name, err)
		}
	}
	for _, name := range []string{"e12_in_notnull", "e13_not_in_notnull"} {
		want, err := ExpectedAnswer(c, name, fx, a)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if got := Evaluators[name](fx, a); !reflect.DeepEqual(got, want) {
			t.Errorf("%s disagrees with its oracle:\n got %#v\nwant %#v", name, got, want)
		}
	}
}

func TestX07MembershipMultiplierIsDetected(t *testing.T) {
	c := load(t)
	fx := existenceFixture()
	a := Args{Pred: AmountGreaterThan(100)}
	got := Evaluators["x07_membership_multiplier"](fx, a)
	if err := DetectControl(c, "x07_membership_multiplier", fx, a, got, "sel_membership"); err != nil {
		t.Errorf("%v", err)
	}
}

func TestX02NullableNotInIsDetected(t *testing.T) {
	c := load(t)
	fx := existenceFixture()
	a := Args{Pred: AmountGreaterThan(100)}
	got := Evaluators["x02_not_in_nullable"](fx, a)
	if err := DetectControl(c, "x02_not_in_nullable", fx, a, got, "sel_absence"); err != nil {
		t.Errorf("%v", err)
	}
}

func TestExistenceSpecsCoverEveryFanoutAndSpelling(t *testing.T) {
	levels := FanoutLevels()
	specs := BuildExistenceSpecs(Args{}, levels[1], 5)
	if want := len(MembershipStatements) + len(AbsenceStatements); len(specs) != want {
		t.Fatalf("specs = %d, want %d", len(specs), want)
	}
	seen := map[string]bool{}
	for _, s := range specs {
		if seen[s.Axis] {
			t.Errorf("duplicate axis %q", s.Axis)
		}
		seen[s.Axis] = true
		if !strings.Contains(s.Axis, "fanout=10") {
			t.Errorf("%s lacks the fanout coordinate", s.Axis)
		}
		if s.Args.Trials != 5 {
			t.Errorf("%s trials = %d", s.Axis, s.Args.Trials)
		}
	}
	if len(ExistenceControlSpecs(Args{}, 3)) != 2 {
		t.Fatal("expected the two existence controls")
	}
}

func TestExistencePositionControlPairsTheFirstMembershipCell(t *testing.T) {
	fanout := FanoutLevels()[0]
	pc := ExistencePositionControlSpec(Args{}, 5)
	if pc.Statement != MembershipStatements[0] || !strings.HasSuffix(pc.Axis, PositionSuffix) {
		t.Fatalf("position spec = %+v", pc)
	}
	r := ExistenceResult{Levels: []FanoutLevel{{Fanout: fanout}}, Cells: []Cell{
		{Axis: ExistenceAxis("membership", MembershipStatements[0], fanout), Statement: MembershipStatements[0],
			Status: StatusValid, TimingMeasured: true, MedianMS: 1},
		{Axis: pc.Axis, Statement: MembershipStatements[0], Status: StatusValid, TimingMeasured: true, MedianMS: 3},
	}}
	pp, ok := r.PositionPair()
	if !ok || pp.RatioBA != 3 {
		t.Fatalf("position pair = %+v, ok=%v", pp, ok)
	}
}
