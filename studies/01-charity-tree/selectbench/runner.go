package selectbench

import (
	"context"
	"errors"
	"strconv"
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
		case "n":
			m[p] = strconv.Itoa(a.N)
		case "limit":
			m[p] = strconv.Itoa(a.Limit)
		case "offset":
			m[p] = strconv.Itoa(a.Offset)
		case "last_amount":
			m[p] = strconv.Itoa(a.LastAmount)
		case "last_id":
			m[p] = strconv.Itoa(a.LastID)
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

// RunCells executes the named statements and returns one cell each. It never
// returns an error for an individual failing statement: a failure becomes the
// cell's status, so a run reports the failure instead of losing it.
func RunCells(ctx context.Context, q Querier, c *Catalog, fx Fixture, a Args, names []string, runID, topology string) []Cell {
	if len(names) == 0 {
		for _, n := range c.Order {
			if c.Statements[n].Group != ControlGroup {
				names = append(names, n)
			}
		}
	}
	cells := make([]Cell, 0, len(names))
	for _, name := range names {
		s, ok := c.Statements[name]
		if !ok {
			cells = append(cells, Cell{RunID: runID, Topology: topology, Statement: name,
				Status: StatusInvalid, Note: "unknown statement"})
			continue
		}
		cell := Cell{RunID: runID, Topology: topology, Statement: name, Group: s.Group,
			Params: describeParams(a.ParamMap(s))}

		bound, _, err := Bind(s.SQL, a.ParamMap(s))
		if err != nil {
			cell.Status, cell.Note = StatusInvalid, err.Error()
			cells = append(cells, cell)
			continue
		}
		plan, err := q.Explain(ctx, bound)
		if err != nil {
			if errors.Is(err, ErrUnsupported) {
				cell.Status, cell.Note = StatusUnsupported, err.Error()
			} else {
				cell.Status, cell.Note = StatusTimeout, err.Error()
			}
			cells = append(cells, cell)
			continue
		}
		cell.Plan = plan

		rows, err := q.Query(ctx, bound)
		if err != nil {
			if errors.Is(err, ErrUnsupported) {
				cell.Status, cell.Note = StatusUnsupported, err.Error()
			} else {
				cell.Status, cell.Note = StatusTimeout, err.Error()
			}
			cells = append(cells, cell)
			continue
		}
		if err := VerifyRows(c, name, fx, a, rows); err != nil {
			cell.Status, cell.Note = StatusInvalid, err.Error()
			cells = append(cells, cell)
			continue
		}
		cell.Status, cell.Trials = StatusValid, 1
		cells = append(cells, cell)
	}
	return cells
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
