package selectbench

import "fmt"

// Membership and absence sweep definitions for the Study 01 SELECT experiment
// (task select-existence).
//
// Fanout is a property of the loaded fixture, not of the SQL: every level loads
// the same logical volume (10,000 child rows) with {1, 10, 100} children per
// parent, so the parent count is {10,000, 1,000, 100}. The harness reloads the
// fixture per level and runs the same statements with the same predicate.

// FanoutLevels are the child-fanout levels, smallest fanout (most parents) first.
func FanoutLevels() []int { return []int{1, 10, 100} }

// TotalChildren is held constant across the fanout sweep.
const TotalChildren = 10000

// MembershipStatements: EXISTS, IN over the nullable child key, IN over the
// NOT-NULL contract, JOIN + DISTINCT. All answer the set of parent identities.
var MembershipStatements = []string{
	"s06_exists_parent",
	"s07_in_parent",
	"e12_in_notnull",
	"s08_join_distinct_parent",
}

// AbsenceStatements: NOT EXISTS, LEFT JOIN ... IS NULL, NOT IN under the
// NOT-NULL contract. All answer the set of parents with no qualifying child.
var AbsenceStatements = []string{
	"s09_not_exists",
	"s10_left_anti",
	"e13_not_in_notnull",
}

// ExistenceAxis names one cell. Fanout is in the coordinate so a report can
// never pool two fanout levels.
func ExistenceAxis(family, statement string, fanout int) string {
	return fmt.Sprintf("family=%s;stmt=%s;fanout=%d", family, shortExistenceStatement(statement), fanout)
}

func shortExistenceStatement(name string) string {
	switch name {
	case "s06_exists_parent":
		return "exists"
	case "s07_in_parent":
		return "in"
	case "e12_in_notnull":
		return "in_notnull"
	case "s08_join_distinct_parent":
		return "join_distinct"
	case "s09_not_exists":
		return "not_exists"
	case "s10_left_anti":
		return "left_is_null"
	case "e13_not_in_notnull":
		return "not_in_notnull"
	}
	return name
}

// BuildExistenceSpecs is the matrix for one fanout level.
func BuildExistenceSpecs(base Args, fanout, trials int) []CellSpec {
	var specs []CellSpec
	a := base
	a.Trials = trials
	for _, name := range MembershipStatements {
		specs = append(specs, CellSpec{Statement: name, Axis: ExistenceAxis("membership", name, fanout), Args: a})
	}
	for _, name := range AbsenceStatements {
		specs = append(specs, CellSpec{Statement: name, Axis: ExistenceAxis("absence", name, fanout), Args: a})
	}
	return specs
}

// ExistencePositionControlSpec repeats the first membership cell of the smallest
// fanout level within that level's phase.
func ExistencePositionControlSpec(base Args, trials int) CellSpec {
	fanout := FanoutLevels()[0]
	a := base
	a.Trials = trials
	return CellSpec{
		Statement: MembershipStatements[0],
		Axis:      ExistenceAxis("membership", MembershipStatements[0], fanout) + PositionSuffix,
		Args:      a,
	}
}

// ExistenceControlSpecs are the controls this family runs at every level.
func ExistenceControlSpecs(base Args, trials int) []ControlSpec {
	mk := func(name string) ControlSpec {
		a := base
		a.Trials = 1
		return ControlSpec{Statement: name, Args: a}
	}
	return []ControlSpec{mk("x02_not_in_nullable"), mk("x07_membership_multiplier")}
}

// FanoutLevel records one loaded fanout fixture.
type FanoutLevel struct {
	Fanout      int    `json:"fanout"`
	Parents     int    `json:"parents"`
	Children    int    `json:"children"`
	Fingerprint string `json:"fingerprint"`
}

// ExistenceResult is one whole membership/absence run.
type ExistenceResult struct {
	RunID      string          `json:"run_id"`
	Engine     string          `json:"engine"`
	Topology   string          `json:"topology"`
	Scale      string          `json:"scale"`
	Seed       int64           `json:"seed"`
	Levels     []FanoutLevel   `json:"levels"`
	Provenance Provenance      `json:"provenance"`
	Cells      []Cell          `json:"cells"`
	Controls   []ControlResult `json:"controls,omitempty"`
	Summary    map[string]int  `json:"summary"`
	Notes      []string        `json:"notes,omitempty"`
	Error      string          `json:"error,omitempty"`
}

// Valid counts cells whose answer matched the oracle and whose timing was measured.
func (r ExistenceResult) Valid() int { return r.Summary["valid"] }

// Fanouts returns the levels' fanouts in order.
func (r ExistenceResult) Fanouts() []int {
	out := make([]int, 0, len(r.Levels))
	for _, l := range r.Levels {
		out = append(out, l.Fanout)
	}
	return out
}

// Find returns the cell for an axis, if present.
func (r ExistenceResult) Find(axis string) (Cell, bool) {
	for _, c := range r.Cells {
		if c.Axis == axis {
			return c, true
		}
	}
	return Cell{}, false
}

// MembershipPair is one controlled comparison within one fanout level.
type MembershipPair struct {
	Fanout  int
	A, B    Cell
	RatioBA float64
	Label   string
}

// MembershipPairs compares each alternative membership spelling with EXISTS at
// the same fanout, and JOIN+DISTINCT with IN.
func (r ExistenceResult) MembershipPairs() []MembershipPair {
	var out []MembershipPair
	for _, fanout := range r.Fanouts() {
		exists, okE := r.Find(ExistenceAxis("membership", "s06_exists_parent", fanout))
		in, okI := r.Find(ExistenceAxis("membership", "s07_in_parent", fanout))
		inNN, okINN := r.Find(ExistenceAxis("membership", "e12_in_notnull", fanout))
		jd, okJD := r.Find(ExistenceAxis("membership", "s08_join_distinct_parent", fanout))
		if okE && okI {
			out = append(out, newMembershipPair(fanout, exists, in, "in vs exists"))
		}
		if okE && okINN {
			out = append(out, newMembershipPair(fanout, exists, inNN, "in_notnull vs exists"))
		}
		if okI && okJD {
			out = append(out, newMembershipPair(fanout, in, jd, "join_distinct vs in"))
		}
		if okE && okJD {
			out = append(out, newMembershipPair(fanout, exists, jd, "join_distinct vs exists"))
		}
	}
	return out
}

func newMembershipPair(fanout int, a, b Cell, label string) MembershipPair {
	p := MembershipPair{Fanout: fanout, A: a, B: b, Label: label}
	if a.MedianMS > 0 {
		p.RatioBA = b.MedianMS / a.MedianMS
	}
	return p
}

// PositionPair returns the first membership cell and its repeat, if present.
func (r ExistenceResult) PositionPair() (MembershipPair, bool) {
	fanout := FanoutLevels()[0]
	first, okFirst := r.Find(ExistenceAxis("membership", MembershipStatements[0], fanout))
	last, okLast := r.Find(ExistenceAxis("membership", MembershipStatements[0], fanout) + PositionSuffix)
	if !okFirst || !okLast {
		return MembershipPair{}, false
	}
	return newMembershipPair(fanout, first, last, "position control exists first vs last"), true
}
