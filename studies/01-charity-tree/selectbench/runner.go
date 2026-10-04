package selectbench

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"strconv"
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

// RunSpecs executes one cell per spec, preserving order.
func RunSpecs(ctx context.Context, q Querier, c *Catalog, fx Fixture, specs []CellSpec, runID, topology string) []Cell {
	cells := make([]Cell, 0, len(specs))
	for _, spec := range specs {
		cells = append(cells, runOne(ctx, q, c, fx, spec, runID, topology))
	}
	return cells
}

func runOne(ctx context.Context, q Querier, c *Catalog, fx Fixture, spec CellSpec, runID, topology string) Cell {
	a := spec.Args
	name := spec.Statement
	s, ok := c.Statements[name]
	if !ok {
		return Cell{RunID: runID, Topology: topology, Statement: name, Axis: spec.Axis,
			Status: StatusInvalid, Note: "unknown statement"}
	}
	cell := Cell{RunID: runID, Topology: topology, Statement: name, Group: s.Group, Axis: spec.Axis,
		Params: describeParams(a.ParamMap(s))}

	bound, _, err := Bind(s.SQL, a.ParamMap(s))
	if err != nil {
		cell.Status, cell.Note = StatusInvalid, err.Error()
		return cell
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
		return cell
	}
	if err := VerifyRows(c, name, fx, a, rows); err != nil {
		cell.Status, cell.Note = StatusInvalid, err.Error()
		return cell
	}
	cell.ResultRows = len(rows)
	encoded, _ := json.Marshal(rows)
	cell.ResultJSONBytes = len(encoded)
	plan, err := q.Explain(cellCtx, bound)
	if err != nil {
		cell.Status, cell.Note = failureStatus(err), err.Error()
		return cell
	}
	cell.Plan = plan
	if plan == "" {
		cell.Status, cell.Note = StatusError, "missing EXPLAIN evidence"
		return cell
	}
	cell.PlanningMS, cell.ExecutionMS = planTimes(plan)
	cell.Status = StatusValid
	if a.ValidateOnly {
		return cell
	}
	trials := a.Trials
	if trials < 1 {
		trials = 1
	}
	for trial := 0; trial < trials; trial++ {
		start := time.Now()
		observed, err := q.Query(cellCtx, bound)
		elapsed := float64(time.Since(start).Nanoseconds()) / 1e6
		if err != nil {
			cell.Status, cell.Note = failureStatus(err), err.Error()
			break
		}
		if err := VerifyRows(c, name, fx, a, observed); err != nil {
			cell.Status, cell.Note = StatusInvalid, err.Error()
			break
		}
		cell.SamplesMS = append(cell.SamplesMS, elapsed)
	}
	if cell.Status == StatusValid {
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
	} else {
		// A late mismatch invalidates the whole cell, including earlier samples.
		cell.SamplesMS = nil
	}
	return cell
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
