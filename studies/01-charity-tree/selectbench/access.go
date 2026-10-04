package selectbench

// Access-path measurement matrix (task select-access-paths).
//
// The protocol fixes two bounded comparisons for this family:
//
//	P10  OR vs overlap-safe UNION ALL, at three predicate shapes
//	     (indexable, expression, cast), with the deliberately-wrong
//	     overlapping UNION ALL as the required negative control;
//	P11  OFFSET vs keyset pagination at page depths 1, 10, 100, 1k, 10k.
//
// This file holds only the pure part: the predicate shapes, the cell specs, and
// the result/report types. The database is reached through the Querier in the
// harness, so every cell here can be planned and reported without an engine.

import (
	"fmt"
	"sort"
	"strings"
)

// PredShape is one SQL spelling of the same logical disjunction, at one
// selectivity regime. The Go predicates A and B are the same logical conditions
// wherever the thresholds match; only ASQL/BSQL differ between the indexable,
// expression and cast spellings, which is what makes a timing difference
// attributable to the spelling rather than to a different answer.
type PredShape struct {
	Name string
	A    Pred
	B    Pred
	ASQL string
	BSQL string
}

// PredRegime is one selectivity regime: the two arm thresholds. The arms
// `amount <= Lo` and `amount >= Hi` stay disjoint by construction, so UNION ALL
// answers the same bag as OR and is a valid competing spelling at every regime.
type PredRegime struct {
	Name string
	Lo   int
	Hi   int
}

// DefaultPredicateRegimes cross the acceptance's selective/broad requirement.
// `selective` picks about 1% of rows per arm; `broad` about 20% per arm. The
// thresholds are workload levels, not results.
func DefaultPredicateRegimes() []PredRegime {
	return []PredRegime{
		{Name: "selective", Lo: 5, Hi: 496},
		{Name: "broad", Lo: 100, Hi: 400},
	}
}

// DefaultPredicateShapes are the bounded shapes from the protocol, crossed with
// the two selectivity regimes: {selective, broad} x {indexable, expression,
// cast}. The shape name keeps both coordinates so a report and a pair can never
// pool a selective cell with a broad one.
func DefaultPredicateShapes() []PredShape {
	spellings := []struct{ name, fmtA, fmtB string }{
		{"indexable", "amount <= %d", "amount >= %d"},
		{"expression", "(amount + 0) <= %d", "(amount + 0) >= %d"},
		{"cast", "(amount)::bigint <= %d", "(amount)::bigint >= %d"},
	}
	var out []PredShape
	for _, reg := range DefaultPredicateRegimes() {
		for _, sp := range spellings {
			out = append(out, PredShape{
				Name: reg.Name + "/" + sp.name,
				A:    AmountLE(reg.Lo), B: AmountGE(reg.Hi),
				ASQL: fmt.Sprintf(sp.fmtA, reg.Lo), BSQL: fmt.Sprintf(sp.fmtB, reg.Hi),
			})
		}
	}
	return out
}

// DefaultPageDepths are the five protocol levels for P11.
func DefaultPageDepths() []int { return []int{1, 10, 100, 1000, 10000} }

// ExpressionIndexName is the physical factor paired with the expression shape.
const ExpressionIndexName = "orders_amount_expr_idx"

// ExpressionIndexSQL creates the expression index that can serve the
// expression spelling. It is a *physical* arm: any gain it produces belongs to
// the index, never to the query text, which is why it is labelled separately.
func ExpressionIndexSQL() string {
	return "CREATE INDEX IF NOT EXISTS " + ExpressionIndexName +
		" ON " + TeachingSchema + ".orders ((amount + 0))"
}

// DropExpressionIndexSQL removes the expression index arm.
func DropExpressionIndexSQL() string {
	return "DROP INDEX IF EXISTS " + TeachingSchema + "." + ExpressionIndexName
}

// defaultArms mirrors the harness's default bound SQL for the two disjunction
// arms, so an Args with nil predicates still has an oracle consistent with the
// SQL that would actually run.
func defaultArms(a Args) (Pred, Pred) {
	a1, a2 := a.PredA, a.PredB
	if a1 == nil {
		a1 = AmountLE(100)
	}
	if a2 == nil {
		a2 = AmountGE(400)
	}
	return a1, a2
}

// shapeSpecs builds the three predicate-alternative cells for one shape.
func shapeSpecs(base Args, shape PredShape, trials int, indexArm string) []CellSpec {
	a := base
	a.PredA, a.PredB = shape.A, shape.B
	a.PredASQL, a.PredBSQL = shape.ASQL, shape.BSQL
	a.Trials = trials
	axis := "shape=" + shape.Name
	if indexArm != "" {
		axis += ";index=" + indexArm
	}
	return []CellSpec{
		{Statement: "s18_or_predicate", Axis: axis + ";spelling=or", Args: a},
		{Statement: "s19_union_all_predicate", Axis: axis + ";spelling=union_all", Args: a},
		{Statement: "s20_union_predicate", Axis: axis + ";spelling=union", Args: a},
	}
}

// cursorAt is the last row of the page that ends at index `index`, i.e. the
// cursor that makes the keyset statement return Page(fx, limit, index+1). It is
// nil for a negative index (first page).
func cursorAt(fx Fixture, index int) *PageRow {
	if index < 0 {
		return nil
	}
	page := Page(fx, 1, index)
	if len(page) == 0 {
		return nil
	}
	row := page[0]
	return &row
}

func pageAxis(depth int, spelling string) string {
	return "depth=" + itoa(depth) + ";spelling=" + spelling
}

// PositionControlSpec repeats the first matrix cell (the first shape's OR
// spelling) at the very end of the run. The same spelling then appears both
// first and last in the phase, which is the protocol's repeated-design position
// control: the first occurrence is the implicit `position=first`, this is
// `position=last`.
func PositionControlSpec(base Args, shape PredShape, trials int) CellSpec {
	a := base
	a.PredA, a.PredB = shape.A, shape.B
	a.PredASQL, a.PredBSQL = shape.ASQL, shape.BSQL
	a.Trials = trials
	return CellSpec{Statement: "s18_or_predicate",
		Axis: "shape=" + shape.Name + ";spelling=or;position=last", Args: a}
}

// paginationSpecs builds the OFFSET/keyset pair at each depth. Both address the
// same ordered page: OFFSET skips `depth` rows, keyset resumes after the row at
// position depth-1.
func paginationSpecs(fx Fixture, base Args, depths []int, limit, trials int) []CellSpec {
	var out []CellSpec
	for _, d := range depths {
		off := base
		off.Limit, off.Offset, off.Trials, off.Last = limit, d, trials, nil
		out = append(out, CellSpec{Statement: "s16_page_offset", Axis: pageAxis(d, "offset"), Args: off})
		key := base
		key.Limit, key.Trials, key.Offset = limit, trials, 0
		key.Last = cursorAt(fx, d-1)
		out = append(out, CellSpec{Statement: "s17_page_keyset", Axis: pageAxis(d, "keyset"), Args: key})
	}
	return out
}

// BuildAccessSpecs is the full default-arm matrix: predicate-alternative cells
// for every shape, then the pagination pair at every depth. The expression
// index arm is added by the harness after it creates the index, so the default
// arm never silently contains a changed physical design.
func BuildAccessSpecs(fx Fixture, base Args, shapes []PredShape, depths []int, limit, trials int) []CellSpec {
	var out []CellSpec
	for _, shape := range shapes {
		out = append(out, shapeSpecs(base, shape, trials, "")...)
	}
	out = append(out, paginationSpecs(fx, base, depths, limit, trials)...)
	return out
}

// --- results ---------------------------------------------------------------

// OverlapRecord proves (or refuses) the disjointness a UNION ALL spelling needs.
type OverlapRecord struct {
	Shape        string `json:"shape"`
	ASQL         string `json:"arm_a"`
	BSQL         string `json:"arm_b"`
	ArmAMatches  int    `json:"arm_a_matches"`
	ArmBMatches  int    `json:"arm_b_matches"`
	OverlapCount int    `json:"overlap_count"`
	Disjoint     bool   `json:"disjoint"`
}

// AccessControl is a control result tagged with the predicate shape it ran in.
type AccessControl struct {
	Shape string `json:"shape"`
	ControlResult
}

// IndexCost is the write/storage cost of the expression-index arm.
type IndexCost struct {
	IndexName      string   `json:"index_name"`
	Expression     string   `json:"expression"`
	Supported      bool     `json:"supported"`
	SizeBytes      *int64   `json:"size_bytes,omitempty"`
	WithIndexMS    *float64 `json:"with_index_ms,omitempty"`
	WithoutIndexMS *float64 `json:"without_index_ms,omitempty"`
	Rows           int      `json:"rows"`
	Note           string   `json:"note,omitempty"`
}

// SnapshotCheck records the fixed-snapshot page-equality check: within one
// transaction, the OFFSET and keyset pages must contain the same rows at every
// depth.
type SnapshotCheck struct {
	Isolation string       `json:"isolation"`
	Depths    []SnapshotAt `json:"depths"`
	AllEqual  bool         `json:"all_equal"`
	Note      string       `json:"note,omitempty"`
	Error     string       `json:"error,omitempty"`
}

type SnapshotAt struct {
	Depth   int  `json:"depth"`
	Equal   bool `json:"equal"`
	OffsetN int  `json:"offset_rows"`
	KeysetN int  `json:"keyset_rows"`
}

// AccessResult is one whole access-path run for one engine/topology.
type AccessResult struct {
	RunID      string          `json:"run_id"`
	Engine     string          `json:"engine"`
	Topology   string          `json:"topology"`
	Scale      string          `json:"scale"`
	Seed       int64           `json:"seed"`
	Fixture    string          `json:"fixture_fingerprint"`
	Limit      int             `json:"page_limit"`
	Shapes     []string        `json:"shapes"`
	Depths     []int           `json:"page_depths"`
	Provenance Provenance      `json:"provenance"`
	Cells      []Cell          `json:"cells"`
	Summary    map[string]int  `json:"summary"`
	Controls   []AccessControl `json:"controls"`
	Overlap    []OverlapRecord `json:"overlap"`
	Index      IndexCost       `json:"index_cost"`
	Snapshot   SnapshotCheck   `json:"snapshot_check"`
	Order      RunOrder        `json:"order"`
	IndexOrder RunOrder        `json:"index_order"`
	Notes      []string        `json:"notes,omitempty"`
	Error      string          `json:"error,omitempty"`
}

// Valid returns the valid cells, in order.
func (r AccessResult) Valid() []Cell {
	var out []Cell
	for _, c := range r.Cells {
		if c.Status == StatusValid {
			out = append(out, c)
		}
	}
	return out
}

// Find returns the first cell matching both the statement and an axis
// substring, or false.
func (r AccessResult) Find(statement, axisContains string) (Cell, bool) {
	for _, c := range r.Cells {
		if c.Statement == statement && strings.Contains(c.Axis, axisContains) {
			return c, true
		}
	}
	return Cell{}, false
}

// PairDelta classifies two cells that should answer the same question. It is
// deliberately conservative: a difference inside the two cells' spread is "no
// resolved difference", and any difference is only within-run evidence until
// whole runs are replicated.
type PairDelta struct {
	Label       string
	A, B        Cell
	SamePlan    bool
	Resolved    bool
	RatioBA     float64 // median B / median A, 0 when A did not time
	Explanation string
}

// Pairs returns the comparable pairs this matrix exists to report: OR vs
// UNION ALL and OR vs UNION per shape, and OFFSET vs keyset per depth.
func (r AccessResult) Pairs() []PairDelta {
	var out []PairDelta
	for _, shape := range r.Shapes {
		if a, ok := r.Find("s18_or_predicate", "shape="+shape+";spelling=or"); ok {
			for _, other := range []struct{ stmt, spelling string }{
				{"s19_union_all_predicate", "union_all"},
				{"s20_union_predicate", "union"},
			} {
				if b, ok2 := r.Find(other.stmt, "shape="+shape+";spelling="+other.spelling); ok2 {
					out = append(out, classifyPair("shape="+shape+" or vs "+other.spelling, a, b))
				}
			}
		}
	}
	for _, d := range r.Depths {
		if o, ok := r.Find("s16_page_offset", "depth="+itoa(d)+";spelling=offset"); ok {
			if k, ok2 := r.Find("s17_page_keyset", "depth="+itoa(d)+";spelling=keyset"); ok2 {
				out = append(out, classifyPair("depth="+itoa(d)+" offset vs keyset", o, k))
			}
		}
	}
	return out
}

// PositionPair returns the repeated-design position control when the run
// contains a `position=last` repeat: the first shape's OR spelling measured as
// the first cell of the phase and again as the last. A large ratio here is a
// session position effect, not a spelling effect.
func (r AccessResult) PositionPair() (PairDelta, bool) {
	if len(r.Shapes) == 0 {
		return PairDelta{}, false
	}
	first, ok := r.Find("s18_or_predicate", "shape="+r.Shapes[0]+";spelling=or")
	if !ok {
		return PairDelta{}, false
	}
	last, ok := r.Find("s18_or_predicate", "shape="+r.Shapes[0]+";spelling=or;position=last")
	if !ok {
		return PairDelta{}, false
	}
	return classifyPair("position control "+r.Shapes[0]+" or first vs last", first, last), true
}

func classifyPair(label string, a, b Cell) PairDelta {
	pd := PairDelta{Label: label, A: a, B: b}
	if a.Status != StatusValid || b.Status != StatusValid || !a.TimingMeasured || !b.TimingMeasured {
		pd.Explanation = "pair incomplete: a cell is not a valid timing"
		return pd
	}
	pd.SamePlan = a.structureSig() == b.structureSig()
	if a.MedianMS > 0 {
		pd.RatioBA = b.MedianMS / a.MedianMS
	}
	diff := b.MedianMS - a.MedianMS
	tol := (a.SpreadPct + b.SpreadPct) / 100 * a.MedianMS
	if absf(diff) <= tol {
		pd.Explanation = "medians inside the two cells' own spread; no resolved difference within this run"
		return pd
	}
	pd.Resolved = true
	pd.Explanation = "within-run medians differ; whole-run replication required before a verdict"
	return pd
}

func absf(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}

// itoa is strconv.Itoa without importing strconv here.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

// sortCells orders cells by axis for stable report output.
func sortCells(cells []Cell) []Cell {
	out := append([]Cell(nil), cells...)
	sort.SliceStable(out, func(i, j int) bool { return out[i].Axis < out[j].Axis })
	return out
}

// IndexArmSpecs exposes the single-shape spec builder for the labelled
// expression-index arm the harness runs after it creates the index.
func IndexArmSpecs(base Args, shape PredShape, trials int, arm string) []CellSpec {
	return shapeSpecs(base, shape, trials, arm)
}

// CursorAt is the exported cursor position used by the harness snapshot check.
func CursorAt(fx Fixture, index int) *PageRow { return cursorAt(fx, index) }
