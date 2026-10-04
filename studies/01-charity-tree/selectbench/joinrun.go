package selectbench

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

// The join/CTE cell runner.
//
// It mirrors the teaching runner's guarantees: a cell is only StatusValid when
// its answer matches the independent Go oracle, a rejected statement is
// unsupported, a budget overrun is a timeout, and no failing cell is ever timed
// or ranked. It differs in carrying the join Args and in recording the observed
// CTE materialization evidence the policy contract requires.

// JoinArgs are the bound parameters of one join/CTE cell.
type JoinArgs struct {
	Threshold        int
	Refs             int
	ConsumerMax      int
	Trials           int
	ValidateOnly     bool
	StatementTimeout time.Duration
}

// JoinSpec is one cell: a statement plus the exact Args it is bound with.
type JoinSpec struct {
	Statement string
	Axis      string
	Args      JoinArgs
}

// JoinParamMap renders the statement's declared parameters.
func JoinParamMap(s Statement, a JoinArgs) map[string]string {
	m := map[string]string{}
	for _, p := range s.Params {
		switch p {
		case "threshold":
			m[p] = strconv.Itoa(a.Threshold)
		case "consumer_max":
			m[p] = strconv.Itoa(a.ConsumerMax)
		default:
			m[p] = "1"
		}
	}
	return m
}

// planHasCTEScan reports whether the plan kept the CTE rather than inlining it.
// It is the observable evidence that a forced MATERIALIZED policy was honoured
// (or that a NOT MATERIALIZED / default CTE was inlined).
func planHasCTEScan(plan string) bool {
	return strings.Contains(strings.ToLower(plan), "cte scan on")
}

// applyPolicyEvidence records whether the forced CTE policy was observed. The
// evidence is the CTE Scan: PostgreSQL 17 does not print a Materialize node for
// an `AS MATERIALIZED` CTE (verified directly, 2026-10-04), it prints a `CTE e`
// subplan plus `CTE Scan on e`. An inlined CTE (default single reference, or
// NOT MATERIALIZED) is expanded into the main plan and prints neither. A forced
// MATERIALIZED cell that was inlined is unsupported: the harness must not
// silently report a default fold as a MATERIALIZED run.
func applyPolicyEvidence(cell *Cell, name, plan string) bool {
	policy := joinPolicyFor(name)
	if policy == "" {
		if strings.HasPrefix(name, "r_") {
			// reuse statements encode policy in the name
			switch {
			case strings.HasPrefix(name, "r_mat_"):
				policy = "materialized"
			case strings.HasPrefix(name, "r_notmat_"):
				policy = "not_materialized"
			default:
				policy = "default"
			}
		}
	}
	if policy == "" {
		return true
	}
	observed := planHasCTEScan(plan)
	switch policy {
	case "materialized":
		if !observed {
			cell.Status = StatusUnsupported
			cell.Note = "MATERIALIZED requested but the plan shows no CTE Scan (the CTE was inlined)"
			cell.PolicyNote = "policy=materialized; observed_cte_scan=false"
			return false
		}
		cell.PolicyNote = "policy=materialized; observed_cte_scan=true"
	case "not_materialized":
		cell.PolicyNote = "policy=not_materialized; observed_cte_scan=" + strconv.FormatBool(observed)
	default:
		cell.PolicyNote = "policy=default; observed_cte_scan=" + strconv.FormatBool(observed)
	}
	return true
}

// RunJoinSpecs executes one cell per spec, taking all of a cell's samples before
// the next cell starts.
func RunJoinSpecs(ctx context.Context, q Querier, c *Catalog, jf JoinFixture, specs []JoinSpec, runID, topology string) []Cell {
	cells := make([]Cell, 0, len(specs))
	for _, spec := range specs {
		cell, bound, ok := prepareJoinCell(ctx, q, c, jf, spec, runID, topology)
		trials := spec.Args.Trials
		if trials < 1 {
			trials = 1
		}
		for t := 0; ok && !spec.Args.ValidateOnly && t < trials && cell.Status == StatusValid; t++ {
			sampleJoinCell(ctx, q, c, jf, spec, bound, &cell)
		}
		finishCell(&cell)
		cells = append(cells, cell)
	}
	return cells
}

// RunJoinSpecsInterleaved verifies and plans each cell once in the plan's order,
// then takes one sample of every valid cell per round, so within-run drift is
// spread across cells instead of landing on the late ones.
func RunJoinSpecsInterleaved(ctx context.Context, q Querier, c *Catalog, jf JoinFixture, specs []JoinSpec, plan OrderPlan, runID, topology string) []Cell {
	cells := make([]Cell, len(specs))
	bounds := make([]string, len(specs))
	ready := make([]bool, len(specs))
	for _, i := range plan.Order {
		cell, bound, ok := prepareJoinCell(ctx, q, c, jf, specs[i], runID, topology)
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
			sampleJoinCell(ctx, q, c, jf, specs[i], bounds[i], &cells[i])
		}
	}
	for i := range cells {
		finishCell(&cells[i])
	}
	return cells
}

func prepareJoinCell(ctx context.Context, q Querier, c *Catalog, jf JoinFixture, spec JoinSpec, runID, topology string) (Cell, string, bool) {
	a := spec.Args
	name := spec.Statement
	s, ok := c.Statements[name]
	if !ok {
		return Cell{RunID: runID, Topology: topology, Statement: name, Axis: spec.Axis,
			Status: StatusInvalid, Note: "unknown statement"}, "", false
	}
	params := JoinParamMap(s, a)
	cell := Cell{RunID: runID, Topology: topology, Statement: name, Group: s.Group, Axis: spec.Axis,
		Params: describeParams(params)}
	bound, _, err := Bind(s.SQL, params)
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
	if err := VerifyJoinRows(c, name, jf, a, rows); err != nil {
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
	if !applyPolicyEvidence(&cell, name, plan) {
		return cell, bound, false
	}
	cell.Status = StatusValid
	return cell, bound, true
}

func sampleJoinCell(ctx context.Context, q Querier, c *Catalog, jf JoinFixture, spec JoinSpec, bound string, cell *Cell) {
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
	if err := VerifyJoinRows(c, spec.Statement, jf, a, observed); err != nil {
		cell.Status, cell.Note = StatusInvalid, err.Error()
		return
	}
	cell.SamplesMS = append(cell.SamplesMS, elapsed)
}

// RunJoinControls executes deliberately-wrong controls and records whether each
// was DETECTED. A control cannot go through prepareJoinCell: its group has no
// oracle by design, so verification is replaced by detection against the
// contract it pretends to serve.
func RunJoinControls(ctx context.Context, q Querier, c *Catalog, jf JoinFixture, specs []JoinSpec) []JoinControlResult {
	var out []JoinControlResult
	for _, spec := range specs {
		group := joinControlContracts[spec.Statement]
		res := JoinControlResult{Statement: spec.Statement, Group: group}
		s, ok := c.Statements[spec.Statement]
		if !ok {
			res.Error = "unknown statement"
			out = append(out, res)
			continue
		}
		budget := spec.Args.StatementTimeout
		if budget <= 0 {
			budget = 2 * time.Minute
		}
		cellCtx, cancel := context.WithTimeout(ctx, budget)
		bound, _, err := Bind(s.SQL, JoinParamMap(s, spec.Args))
		var rows [][]any
		if err == nil {
			rows, err = q.Query(cellCtx, bound)
		}
		cancel()
		res.Rows = len(rows)
		if err == nil {
			err = JoinDetectControl(c, spec.Statement, jf, spec.Args, rows)
		}
		res.Detected = err == nil
		if err != nil {
			res.Error = err.Error()
		}
		out = append(out, res)
	}
	return out
}
