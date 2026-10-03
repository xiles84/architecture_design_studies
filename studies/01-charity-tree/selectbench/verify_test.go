package selectbench

import (
	"testing"
)

func args() Args {
	return Args{Pred: AmountGreaterThan(100), N: 1, Limit: 2, Offset: 0}
}

func TestEveryValidStatementHasAnOracle(t *testing.T) {
	c := load(t)
	fx := Teaching()
	for _, name := range c.Order {
		s := c.Statements[name]
		if s.Group == ControlGroup {
			if _, err := ExpectedAnswer(c, name, fx, args()); err == nil {
				t.Errorf("%s is a control and should have no oracle", name)
			}
			continue
		}
		if _, err := ExpectedAnswer(c, name, fx, args()); err != nil {
			t.Errorf("%s (%s): %v", name, s.Group, err)
		}
	}
}

// Every valid statement's own answer must pass verification: that is the floor
// the DB mode enforces on a live cell.
func TestVerifyAcceptsTheTemplateAnswers(t *testing.T) {
	c := load(t)
	fx := Teaching()
	for _, name := range c.Order {
		if c.Statements[name].Group == ControlGroup {
			continue
		}
		got := Evaluators[name](fx, args())
		if err := Verify(c, name, fx, args(), got); err != nil {
			t.Errorf("%s: %v", name, err)
		}
	}
}

func TestVerifyRejectsAWrongAnswer(t *testing.T) {
	c := load(t)
	fx := Teaching()
	// A membership statement answered with the full-row bag answers a different
	// question and must be rejected.
	if err := Verify(c, "s06_exists_parent", fx, args(), FullRows(fx, args().Pred)); err == nil {
		t.Error("a membership cell with a full-row answer was accepted")
	}
	if err := Verify(c, "no_such_statement", fx, args(), nil); err == nil {
		t.Error("an unknown statement did not error")
	}
}

// Each control must differ from the oracle of the contract it pretends to serve.
func TestControlsAreDetectedAgainstTheirPretendedContract(t *testing.T) {
	c := load(t)
	fx := Teaching()
	a := args()
	cases := []struct {
		control   string
		pretended string
	}{
		{"x01_optional_where", "sel_optional_match"},
		{"x02_not_in_nullable", "sel_absence"},
		{"x03_distinct_on_name", "sel_membership"},
		{"x04_limit_before_join", "sel_full_rows"},
	}
	for _, tc := range cases {
		got := Evaluators[tc.control](fx, a)
		if err := DetectControl(c, tc.control, fx, a, got, tc.pretended); err != nil {
			t.Errorf("%v", err)
		}
	}
	// A mislabelled control (not in the control group) is a harness error.
	if err := DetectControl(c, "s01_join_where", fx, a, nil, "sel_full_rows"); err == nil {
		t.Error("a valid template was accepted as a control")
	}
}
