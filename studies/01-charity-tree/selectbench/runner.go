package selectbench

import (
	"context"
	"encoding/json"
	"errors"
	"hash/fnv"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
)

// The SELECT runner.
//
// The database is reached through Querier, so the whole cell loop -- bind,
// explain, execute, verify, classify -- is unit-testable with a fake and needs
// no database to be checked. A live run supplies a pgx-backed Querier; the
// adapter is the only part that needs the engine, and it must set the search
// path to TeachingSchema so the catalogue's bare table names resolve.
//
// A cell is only StatusValid when its answer matches the contract oracle. A
// statement the engine rejects is unsupported; a statement that exceeds its
// budget is a timeout; a statement whose answer differs is invalid and is never
// timed. No path here reports a number for a cell that failed.

// Querier is the minimal database surface the runner needs.
type Querier interface {
	// Exec runs a statement for effect (DDL/load).
	Exec(ctx context.Context, sql string) error
	// Explain returns EXPLAIN (ANALYZE, BUFFERS, FORMAT TEXT) for a statement.
	Explain(ctx context.Context, sql string) (string, error)
	// Query runs a statement and returns its drained rows.
	Query(ctx context.Context, sql string) ([][]any, error)
}

// ErrUnsupported marks a statement the engine refused to plan or execute at all,
// as distinct from a load or connectivity failure.
var ErrUnsupported = errors.New("unsupported by engine")

// ParamMap renders the statement's declared parameters from Args.
func (a Args) ParamMap(s Statement) map[string]string {
	m := map[string]string{}
	for _, p := range s.Params {
		switch p {
		case "pred":
			m[p] = firstNonEmpty(a.PredSQL, "amount > 100")
		case "pred_a":
			m[p] = firstNonEmpty(a.PredASQL, "amount <= 100")
		case "pred_b":
			m[p] = firstNonEmpty(a.PredBSQL, "amount >= 400")
		case "n":
			m[p] = strconv.Itoa(topNArg(a))
		case "parent_max":
			m[p] = strconv.Itoa(a.ParentMax)
		case "limit":
			m[p] = strconv.Itoa(a.Limit)
		case "offset":
			m[p] = strconv.Itoa(a.Offset)
		case "last_amount":
			m[p] = strconv.Itoa(a.LastAmount)
			if a.Last != nil && a.Last.Amount != nil {
				m[p] = strconv.Itoa(*a.Last.Amount)
			}
		case "last_id":
			m[p] = strconv.Itoa(a.LastID)
			if a.Last != nil {
				m[p] = strconv.Itoa(a.Last.OrderID)
			}
		case "cursor_set":
			m[p] = "0"
			if a.Last != nil {
				m[p] = "1"
			}
		case "last_is_null":
			m[p] = "0"
			if a.Last != nil && a.Last.Amount == nil {
				m[p] = "1"
			}
		default:
			m[p] = "1"
		}
	}
	return m
}

// Prepare loads the fixture into the engine: schema, rows and statistics. The
// data is the same Fixture the oracles read, so a mismatch is a real defect and
// never a data difference.
func Prepare(ctx context.Context, q Querier, fx Fixture) error {
	if err := q.Exec(ctx, SchemaSQL()); err != nil {
		return err
	}
	if err := q.Exec(ctx, LoadSQL(fx)); err != nil {
		return err
	}
	return q.Exec(ctx, "ANALYZE "+TeachingSchema+".customers; ANALYZE "+TeachingSchema+".orders;")
}

// CellSpec is one cell to execute: a statement plus the exact Args it is bound
// with. The access-path and aggregate/top-N matrices bind different Args per
// cell (a parent prefix, a page offset, a cursor, a predicate shape), so the
// runner works from specs rather than one shared Args for every name.
type CellSpec struct {
	Statement string
	Axis      string
	Args      Args
}

// RunCells executes the named statements with one shared Args and returns one
// cell each. It never returns an error for an individual failing statement: a
// failure becomes the cell's status, so a run reports the failure instead of
// losing it.
func RunCells(ctx context.Context, q Querier, c *Catalog, fx Fixture, a Args, names []string, runID, topology string) []Cell {
	if len(names) == 0 {
		for _, n := range c.Order {
			s := c.Statements[n]
			if s.Group == ControlGroup {
				continue
			}
			names = append(names, n)
		}
	}
	specs := make([]CellSpec, 0, len(names))
	for _, name := range names {
		specs = append(specs, CellSpec{Statement: name, Args: a})
	}
	return RunSpecs(ctx, q, c, fx, specs, runID, topology)
}

// RunSpecs executes one cell per spec, preserving order, taking all of a
// cell's samples before the next cell starts. The order-corrected matrix uses
// RunSpecsInterleaved instead.
func RunSpecs(ctx context.Context, q Querier, c *Catalog, fx Fixture, specs []CellSpec, runID, topology string) []Cell {
	cells := make([]Cell, 0, len(specs))
	for _, spec := range specs {
		cell, bound, ok := prepareCell(ctx, q, c, fx, spec, runID, topology)
		trials := spec.Args.Trials
		if trials < 1 {
			trials = 1
		}
		for t := 0; ok && !spec.Args.ValidateOnly && t < trials && cell.Status == StatusValid; t++ {
			sampleCell(ctx, q, c, fx, spec, bound, &cell)
		}
		finishCell(&cell)
		cells = append(cells, cell)
	}
	return cells
}

// OrderPlan is a seeded execution order over n cells.
type OrderPlan struct {
	Seed  int64 `json:"seed"`
	Order []int `json:"order"`
}

// PlanOrder returns a deterministic Fisher-Yates permutation of 0..n-1 from the
// seed, so a run's cell order is reproducible and recorded rather than implied.
func PlanOrder(seed int64, n int) OrderPlan {
	order := make([]int, n)
	for i := range order {
		order[i] = i
	}
	r := rand.New(rand.NewSource(seed))
	r.Shuffle(n, func(i, j int) { order[i], order[j] = order[j], order[i] })
	return OrderPlan{Seed: seed, Order: order}
}

// Sequence reports the axis of each spec in plan order.
func (p OrderPlan) Sequence(specs []CellSpec) []string {
	out := make([]string, 0, len(p.Order))
	for _, i := range p.Order {
		out = append(out, specs[i].Axis)
	}
	return out
}

// OrderSeedFor derives a per-run order seed from the dataset seed and the run
// id, so every whole run gets its own recorded order.
func OrderSeedFor(seed int64, runID string) int64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(runID))
	return seed*1000003 + int64(h.Sum64()%1000000)
}

// RunOrder records how a phase was executed: the seed, the number of
// interleaved sample rounds, and the actual axis sequence.
type RunOrder struct {
	Seed     int64    `json:"seed"`
	Rounds   int      `json:"rounds"`
	Sequence []string `json:"sequence"`
}

// RunSpecsInterleaved verifies and plans each cell once in the plan's order,
// then takes one sample of every valid cell per round, so within-run drift is
// spread across cells instead of being concentrated in the late ones. Cells are
// returned in original spec order; the recorded sequence carries execution order.
func RunSpecsInterleaved(ctx context.Context, q Querier, c *Catalog, fx Fixture, specs []CellSpec, plan OrderPlan, runID, topology string) []Cell {
	cells := make([]Cell, len(specs))
	bounds := make([]string, len(specs))
	ready := make([]bool, len(specs))
	for _, i := range plan.Order {
		cell, bound, ok := prepareCell(ctx, q, c, fx, specs[i], runID, topology)
		cells[i], bounds[i], ready[i] = cell, bound, ok
	}
	rounds := 1
	for _, spec := range specs {
		if t := spec.Args.Trials; t > rounds {
			rounds = t
		}
	}
	for round := 0; round < rounds; round++ {
		for _, i := range plan.Order {
			if !ready[i] || specs[i].Args.ValidateOnly || cells[i].Status != StatusValid {
				continue
			}
			if t := specs[i].Args.Trials; t > 0 && round >= t {
				continue
			}
			sampleCell(ctx, q, c, fx, specs[i], bounds[i], &cells[i])
		}
	}
	for i := range cells {
		finishCell(&cells[i])
	}
	return cells
}

// positionPairCells finds the cell whose axis carries PositionSuffix and the
// cell it repeats, by stripping the suffix. It works for the actual randomized
// first cell, not a hard-coded one.
func positionPairCells(cells []Cell) (Cell, Cell, bool) {
	for _, last := range cells {
		if !strings.HasSuffix(last.Axis, PositionSuffix) {
			continue
		}
		base := strings.TrimSuffix(last.Axis, PositionSuffix)
		for _, first := range cells {
			if first.Axis == base {
				return first, last, true
			}
		}
	}
	return Cell{}, Cell{}, false
}

// prepareCell binds, verifies and plans one cell. The bound SQL and whether the
// cell is ready to time are returned with it.
func prepareCell(ctx context.Context, q Querier, c *Catalog, fx Fixture, spec CellSpec, runID, topology string) (Cell, string, bool) {
	a := spec.Args
	name := spec.Statement
	s, ok := c.Statements[name]
	if !ok {
		return Cell{RunID: runID, Topology: topology, Statement: name, Axis: spec.Axis,
			Status: StatusInvalid, Note: "unknown statement"}, "", false
	}
	cell := Cell{RunID: runID, Topology: topology, Statement: name, Group: s.Group, Axis: spec.Axis,
		Params: describeParams(a.ParamMap(s))}
	bound, _, err := Bind(s.SQL, a.ParamMap(s))
	if err != nil {
		cell.Status, cell.Note = StatusInvalid, err.Error()
		return cell, "", false
	}
	budget := a.StatementTimeout
	if budget <= 0 {
		budget = 2 * time.Minute
	}
	cellCtx, cancel := context.WithTimeout(ctx, budget)
	defer cancel()
	rows, err := q.Query(cellCtx, bound)
	if err != nil {
		cell.Status, cell.Note = failureStatus(err), err.Error()
		return cell, bound, false
	}
	if err := VerifyRows(c, name, fx, a, rows); err != nil {
		cell.Status, cell.Note = StatusInvalid, err.Error()
		return cell, bound, false
	}
	cell.ResultRows = len(rows)
	encoded, _ := json.Marshal(rows)
	cell.ResultJSONBytes = len(encoded)
	plan, err := q.Explain(cellCtx, bound)
	if err != nil {
		cell.Status, cell.Note = failureStatus(err), err.Error()
		return cell, bound, false
	}
	cell.Plan = plan
	if plan == "" {
		cell.Status, cell.Note = StatusError, "missing EXPLAIN evidence"
		return cell, bound, false
	}
	cell.PlanningMS, cell.ExecutionMS = planTimes(plan)
	cell.Status = StatusValid
	return cell, bound, true
}

// sampleCell takes one timed, verified sample and appends it to the cell.
func sampleCell(ctx context.Context, q Querier, c *Catalog, fx Fixture, spec CellSpec, bound string, cell *Cell) {
	a := spec.Args
	budget := a.StatementTimeout
	if budget <= 0 {
		budget = 2 * time.Minute
	}
	cellCtx, cancel := context.WithTimeout(ctx, budget)
	defer cancel()
	start := time.Now()
	observed, err := q.Query(cellCtx, bound)
	elapsed := float64(time.Since(start).Nanoseconds()) / 1e6
	if err != nil {
		cell.Status, cell.Note = failureStatus(err), err.Error()
		return
	}
	if err := VerifyRows(c, spec.Statement, fx, a, observed); err != nil {
		cell.Status, cell.Note = StatusInvalid, err.Error()
		return
	}
	cell.SamplesMS = append(cell.SamplesMS, elapsed)
}

// finishCell computes trials, median and spread from the accumulated samples and
// invalidates the whole cell, including earlier samples, if it failed.
func finishCell(cell *Cell) {
	if cell.Status != StatusValid || len(cell.SamplesMS) == 0 {
		cell.SamplesMS = nil
		cell.TimingMeasured = false
		return
	}
	cell.Trials = len(cell.SamplesMS)
	cell.TimingMeasured = true
	ordered := append([]float64(nil), cell.SamplesMS...)
	sort.Float64s(ordered)
	cell.MedianMS = ordered[len(ordered)/2]
	if len(ordered)%2 == 0 {
		cell.MedianMS = (ordered[len(ordered)/2-1] + cell.MedianMS) / 2
	}
	if cell.MedianMS > 0 {
		cell.SpreadPct = 100 * (ordered[len(ordered)-1] - ordered[0]) / cell.MedianMS
	}
}

func failureStatus(err error) Status {
	if errors.Is(err, ErrUnsupported) {
		return StatusUnsupported
	}
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return StatusTimeout
	}
	return StatusError
}

func describeParams(m map[string]string) string {
	if len(m) == 0 {
		return ""
	}
	keys := sortedKeys(m)
	out := ""
	for _, k := range keys {
		if out != "" {
			out += ","
		}
		out += k + "=" + m[k]
	}
	return out
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

// needsCursor reports whether a statement requires a bound pagination cursor.
// Such a statement is excluded from the default cell set: with no cursor it
// would compare an empty page against the first page and fail the gate for a
// reason that is the harness's, not the query's.
func needsCursor(s Statement) bool {
	for _, p := range s.Params {
		if p == "last_amount" || p == "last_id" {
			return true
		}
	}
	return false
}
