package selectbench

import "fmt"

// Aggregate / top-N sweep definitions for the Study 01 SELECT experiment
// (task select-aggregate-topn).
//
// One deterministic fixture holds enough parents for the largest level, and
// every cell restricts the answer to a prefix of parent ids (`parent_max`). The
// five levels carry the protocol's {0.1%, 1%, 10%, 100%} selection plus 0.01%,
// so a crossing claim has the five levels the protocol requires to be bounded.
// On a 10,000-parent fixture they are 1, 10, 100, 1,000 and 10,000 parents --
// the same counts the top-N parent-count axis asks for.

// ParentLevels are the parent-count / parent-selection levels, smallest first.
func ParentLevels() []int { return []int{1, 10, 100, 1000, 10000} }

// Aggregate statements: post-join, pre-aggregate, correlated scalar.
var AggregateStatements = []string{
	"a01_aggregate_join_sel",
	"a02_aggregate_preagg_sel",
	"a03_aggregate_correlated_sel",
}

// Top-N statements: ROW_NUMBER, LATERAL, DISTINCT ON (N=1 only).
var TopNStatements = []string{
	"a04_topn_rownumber_sel",
	"a05_topn_lateral_sel",
	"a06_topn_distinct_on_sel",
}

// AggregateTopNAxis names one aggregate/top-N cell. The report and the
// replicated summary key on it, so it must be stable and engine-independent.
func AggregateTopNAxis(family, statement string, level int) string {
	return fmt.Sprintf("family=%s;stmt=%s;parent_max=%d", family, shortStatement(statement), level)
}

func shortStatement(name string) string {
	switch name {
	case "a01_aggregate_join_sel":
		return "post_join"
	case "a02_aggregate_preagg_sel":
		return "pre_agg"
	case "a03_aggregate_correlated_sel":
		return "correlated"
	case "a04_topn_rownumber_sel":
		return "rownumber"
	case "a05_topn_lateral_sel":
		return "lateral"
	case "a06_topn_distinct_on_sel":
		return "distinct_on"
	}
	return name
}

// BuildAggregateSpecs is the default matrix: every aggregate and top-N spelling
// at every parent level, with one shared timing budget. Controls (x05, x06) are
// run separately by the harness because they are expected to answer wrongly.
func BuildAggregateSpecs(base Args, levels []int, trials int) []CellSpec {
	var specs []CellSpec
	for _, level := range levels {
		for _, name := range AggregateStatements {
			a := base
			a.ParentMax = level
			a.N = 1
			a.Trials = trials
			specs = append(specs, CellSpec{Statement: name, Axis: AggregateTopNAxis("aggregate", name, level), Args: a})
		}
		for _, name := range TopNStatements {
			a := base
			a.ParentMax = level
			a.N = 1 // DISTINCT ON is N=1-only; the whole top-N axis is compared at N=1
			a.Trials = trials
			specs = append(specs, CellSpec{Statement: name, Axis: AggregateTopNAxis("topn", name, level), Args: a})
		}
	}
	return specs
}

// PositionSuffix marks the repeated-design position-control cell: the first
// matrix cell, re-run at the end of the phase under a distinct axis key.
const PositionSuffix = ";position=last"

// AggregatePositionControlSpec repeats the first matrix cell at the smallest parent
// level, so a within-run position/session effect can be read off beside the
// spelling pairs. It is an extra cell, never a replacement.
func AggregatePositionControlSpec(base Args, trials int) CellSpec {
	level := ParentLevels()[0]
	a := base
	a.ParentMax = level
	a.N = 1
	a.Trials = trials
	return CellSpec{
		Statement: AggregateStatements[0],
		Axis:      AggregateTopNAxis("aggregate", AggregateStatements[0], level) + PositionSuffix,
		Args:      a,
	}
}

// ControlSpec is a deliberately-wrong statement run under its own contract, so
// the run records whether the control fired. The aggregate family's controls
// are level-independent for detection; they are run at the widest level.
type ControlSpec struct {
	Statement string
	Args      Args
}

// AggregateControlSpecs returns the control cells for the family.
func AggregateControlSpecs(base Args, trials int) []ControlSpec {
	mk := func(name string) ControlSpec {
		a := base
		a.ParentMax = 0
		a.N = 1
		a.Trials = trials
		return ControlSpec{Statement: name, Args: a}
	}
	return []ControlSpec{mk("x05_aggregate_double_count"), mk("x06_topn_rank_ties")}
}

// AggResult is one whole aggregate/top-N run.
type AggResult struct {
	RunID      string          `json:"run_id"`
	Engine     string          `json:"engine"`
	Topology   string          `json:"topology"`
	Scale      string          `json:"scale"`
	Seed       int64           `json:"seed"`
	Levels     []int           `json:"levels"`
	Provenance Provenance      `json:"provenance,omitempty"`
	Cells      []Cell          `json:"cells"`
	Controls   []ControlResult `json:"controls,omitempty"`
	Order      RunOrder        `json:"order"`
	Summary    map[string]int  `json:"summary"`
	Notes      []string        `json:"notes,omitempty"`
	Error      string          `json:"error,omitempty"`
}

// Valid counts cells whose answer matched the oracle and whose timing was
// measured.
func (r AggResult) Valid() int { return r.Summary["valid"] }

// Find returns the cell for an axis, if present.
func (r AggResult) Find(axis string) (Cell, bool) {
	for _, c := range r.Cells {
		if c.Axis == axis {
			return c, true
		}
	}
	return Cell{}, false
}

// Pair is one controlled comparison: A is the baseline spelling, B the other,
// and RatioBA is B's median divided by A's. Both cells share family and level.
type AggPair struct {
	Family  string
	Level   int
	A, B    Cell
	RatioBA float64
	Label   string
}

// AggregatePairs is a02/a01, a03/a01 and a03/a02 at each level.
func (r AggResult) AggregatePairs() []AggPair {
	var out []AggPair
	for _, level := range r.Levels {
		a, okA := r.Find(AggregateTopNAxis("aggregate", "a01_aggregate_join_sel", level))
		b, okB := r.Find(AggregateTopNAxis("aggregate", "a02_aggregate_preagg_sel", level))
		c, okC := r.Find(AggregateTopNAxis("aggregate", "a03_aggregate_correlated_sel", level))
		if okA && okB {
			out = append(out, newAggPair("aggregate", level, a, b, "pre_agg vs post_join"))
		}
		if okA && okC {
			out = append(out, newAggPair("aggregate", level, a, c, "correlated vs post_join"))
		}
		if okB && okC {
			out = append(out, newAggPair("aggregate", level, b, c, "correlated vs pre_agg"))
		}
	}
	return out
}

// TopNPairs is a05/a04, a06/a04 and a06/a05 at each level (all at N=1).
func (r AggResult) TopNPairs() []AggPair {
	var out []AggPair
	for _, level := range r.Levels {
		a, okA := r.Find(AggregateTopNAxis("topn", "a04_topn_rownumber_sel", level))
		b, okB := r.Find(AggregateTopNAxis("topn", "a05_topn_lateral_sel", level))
		c, okC := r.Find(AggregateTopNAxis("topn", "a06_topn_distinct_on_sel", level))
		if okA && okB {
			out = append(out, newAggPair("topn", level, a, b, "lateral vs rownumber"))
		}
		if okA && okC {
			out = append(out, newAggPair("topn", level, a, c, "distinct_on vs rownumber"))
		}
		if okB && okC {
			out = append(out, newAggPair("topn", level, b, c, "distinct_on vs lateral"))
		}
	}
	return out
}

// PositionPair returns the first matrix cell and its end-of-phase repeat, if the
// run carried the position control.
func (r AggResult) PositionPair() (AggPair, bool) {
	first, last, ok := positionPairCells(r.Cells)
	if !ok {
		return AggPair{}, false
	}
	return newAggPair("position", 0, first, last, "position control first vs last"), true
}

func newAggPair(family string, level int, a, b Cell, label string) AggPair {
	p := AggPair{Family: family, Level: level, A: a, B: b, Label: label}
	if a.MedianMS > 0 {
		p.RatioBA = b.MedianMS / a.MedianMS
	}
	return p
}
