package selectbench

import (
	"fmt"
	"reflect"
)

// Per-cell verification.
//
// A timed cell is only valid when its observed answer equals the oracle for the
// statement's declared contract. This is the glue between a catalogue statement
// and the contract oracle: the DB mode calls Verify on every cell and aborts a
// cell whose answer differs, so a fast wrong answer can never be reported.

// answersEqual compares two contract answers. A nil slice and an empty slice are
// the same answer: "no rows" is "no rows", and a database that returns an empty
// result set must not be failed against an oracle that happened to leave its
// slice nil (or the reverse). Everything else is compared exactly.
func answersEqual(a, b any) bool {
	va, vb := reflect.ValueOf(a), reflect.ValueOf(b)
	if va.IsValid() && vb.IsValid() && va.Kind() == reflect.Slice && vb.Kind() == reflect.Slice &&
		va.Len() == 0 && vb.Len() == 0 {
		return true
	}
	return reflect.DeepEqual(a, b)
}

// oracleFor returns the canonical expected answer for a semantic group.
func oracleFor(group string, fx Fixture, a Args) (any, error) {
	switch group {
	case "sel_full_rows":
		return FullRows(fx, a.Pred), nil
	case "sel_membership":
		return Membership(fx, a.Pred), nil
	case "sel_absence":
		return Absence(fx, a.Pred), nil
	case "sel_optional_match":
		return OptionalMatch(fx, a.Pred), nil
	case "sel_aggregate":
		return Aggregate(fx), nil
	case "sel_topn":
		return TopN(fx, a.N), nil
	case "sel_pagination":
		return Page(fx, a.Limit, a.Offset), nil
	case "sel_predicate_alt":
		a1, a2 := defaultArms(a)
		return PredicateAlt(fx, a1, a2), nil
	case ControlGroup:
		return nil, fmt.Errorf("a deliberately wrong control has no oracle: it must be seen to FAIL, not to match")
	}
	return nil, fmt.Errorf("unknown semantic group %q", group)
}

// ExpectedAnswer is the oracle answer for one statement under its contract.
func ExpectedAnswer(c *Catalog, name string, fx Fixture, a Args) (any, error) {
	s, ok := c.Statements[name]
	if !ok {
		return nil, fmt.Errorf("unknown statement %q", name)
	}
	if needsCursor(s) {
		return KeysetPage(fx, a.Limit, a.Last), nil
	}
	return oracleFor(s.Group, fx, a)
}

// Verify compares an observed answer with the contract oracle. A non-nil error
// means the cell must not be timed or reported.
func Verify(c *Catalog, name string, fx Fixture, a Args, got any) error {
	want, err := ExpectedAnswer(c, name, fx, a)
	if err != nil {
		return err
	}
	if !answersEqual(want, got) {
		s := c.Statements[name]
		return fmt.Errorf("%s (%s): observed answer does not match the contract oracle:\n got %#v\nwant %#v",
			name, s.Group, got, want)
	}
	return nil
}

// DetectControl reports whether a control's observed answer differs from the
// oracle of the contract it pretends to serve. It returns nil only when the
// control was *detected* (i.e. the answers differ); if a control accidentally
// matches the oracle, that is the harness defect this function exists to catch.
func DetectControl(c *Catalog, control string, fx Fixture, a Args, got any, pretendedGroup string) error {
	if s, ok := c.Statements[control]; !ok {
		return fmt.Errorf("unknown statement %q", control)
	} else if s.Group != ControlGroup {
		return fmt.Errorf("%s is not a deliberately wrong control", control)
	}
	want, err := oracleFor(pretendedGroup, fx, a)
	if err != nil {
		return err
	}
	if answersEqual(want, got) {
		return fmt.Errorf("%s was NOT detected: its answer matched the %s oracle", control, pretendedGroup)
	}
	return nil
}
