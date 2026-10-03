package selectbench

import (
	"sort"
	"time"
)

// Args are the bound parameters for one evaluator call.
type Args struct {
	Pred   Pred
	N      int
	Limit  int
	Offset int
	Last   *PageRow
	// ParentMax restricts the aggregate/top-N sweep to parents with id <=
	// ParentMax (0 = every parent). It is the parent-selection / parent-count
	// axis of the aggregate-topn matrix.
	ParentMax int
	// SQL-side parameters, rendered into the statement by Bind.
	PredSQL          string
	LastAmount       int
	LastID           int
	Trials           int
	ValidateOnly     bool
	StatementTimeout time.Duration
}

// EvalFunc is a Go implementation of one catalogue statement's logical query
// over the fixture. The evaluators exist so the harness can check each
// template's *contract* without a database, and so the deliberately-wrong
// controls can be shown to produce a different answer. They are not a
// substitute for running the SQL: the DB-backed check runs under the benchmark
// lock in the measurement tasks.
type EvalFunc func(fx Fixture, a Args) any

// Evaluators maps every catalogue statement name to its Go mirror. A statement
// with no evaluator is a harness defect and is caught by TestEveryStatementHasEvaluator.
var Evaluators = map[string]EvalFunc{
	// full-row bag
	"s01_join_where":       func(fx Fixture, a Args) any { return FullRows(fx, a.Pred) },
	"s02_filtered_derived": func(fx Fixture, a Args) any { return FullRows(fx, a.Pred) },
	"s03_cte_filtered":     func(fx Fixture, a Args) any { return FullRows(fx, a.Pred) },
	"s04_join_on":          func(fx Fixture, a Args) any { return FullRows(fx, a.Pred) },
	"s05_comma_join":       func(fx Fixture, a Args) any { return FullRows(fx, a.Pred) },
	// membership (set of parent identities)
	"s06_exists_parent":        func(fx Fixture, a Args) any { return Membership(fx, a.Pred) },
	"s07_in_parent":            func(fx Fixture, a Args) any { return Membership(fx, a.Pred) },
	"s08_join_distinct_parent": func(fx Fixture, a Args) any { return Membership(fx, a.Pred) },
	// absence
	"s09_not_exists": func(fx Fixture, a Args) any { return Absence(fx, a.Pred) },
	"s10_left_anti":  func(fx Fixture, a Args) any { return Absence(fx, a.Pred) },
	// optional match (parents preserved)
	"s11_optional_on": func(fx Fixture, a Args) any { return OptionalMatch(fx, a.Pred) },
	// aggregates
	"s12_aggregate_join":   func(fx Fixture, a Args) any { return Aggregate(fx) },
	"s13_aggregate_preagg": func(fx Fixture, a Args) any { return Aggregate(fx) },
	// top-N per parent
	"s14_topn_rownumber": func(fx Fixture, a Args) any { return TopN(fx, a.N) },
	"s15_topn_lateral":   func(fx Fixture, a Args) any { return TopN(fx, a.N) },
	// aggregate / top-N sweep family (task select-aggregate-topn): the same
	// contracts restricted to a parent prefix, so one fixture serves the
	// parent-selection and parent-count axes.
	"a01_aggregate_join_sel":       func(fx Fixture, a Args) any { return AggregateUpTo(fx, a.ParentMax) },
	"a02_aggregate_preagg_sel":     func(fx Fixture, a Args) any { return AggregateUpTo(fx, a.ParentMax) },
	"a03_aggregate_correlated_sel": func(fx Fixture, a Args) any { return AggregateUpTo(fx, a.ParentMax) },
	"a04_topn_rownumber_sel":       func(fx Fixture, a Args) any { return TopNUpTo(fx, topNArg(a), a.ParentMax) },
	"a05_topn_lateral_sel":         func(fx Fixture, a Args) any { return TopNUpTo(fx, topNArg(a), a.ParentMax) },
	"a06_topn_distinct_on_sel":     func(fx Fixture, a Args) any { return TopNUpTo(fx, 1, a.ParentMax) },
	// pagination
	"s16_page_offset": func(fx Fixture, a Args) any { return Page(fx, a.Limit, a.Offset) },
	"s17_page_keyset": func(fx Fixture, a Args) any { return KeysetPage(fx, a.Limit, a.Last) },
	// deliberately wrong controls
	"x01_optional_where":         func(fx Fixture, a Args) any { return ControlOptionalWhere(fx, a.Pred) },
	"x02_not_in_nullable":        func(fx Fixture, a Args) any { return ControlNotNullIn(fx, a.Pred) },
	"x03_distinct_on_name":       func(fx Fixture, a Args) any { return ControlDistinctOnName(fx, a.Pred) },
	"x04_limit_before_join":      func(fx Fixture, a Args) any { return ControlLimitBeforeJoin(fx, a.Pred, a.Limit) },
	"x05_aggregate_double_count": func(fx Fixture, a Args) any { return DoubleCountAggregate(fx, a.ParentMax) },
	"x06_topn_rank_ties":         func(fx Fixture, a Args) any { return RankTiesTopN(fx, topNArg(a), a.ParentMax) },
}

// topNArg normalises the requested per-parent N to at least one. The sweep
// statements and their oracles share it, so an unset N can never make the SQL
// and the oracle disagree about how many rows a parent may contribute.
func topNArg(a Args) int {
	if a.N < 1 {
		return 1
	}
	return a.N
}

// KeysetPage returns the page strictly after the cursor under the declared order.
func KeysetPage(fx Fixture, limit int, last *PageRow) []PageRow {
	rows := append([]Order(nil), fx.Orders...)
	sort.Slice(rows, func(i, j int) bool { return orderLess(rows[i].Amount, rows[i].ID, rows[j].Amount, rows[j].ID) })
	var out []PageRow
	for _, o := range rows {
		if last != nil && !orderLess(last.Amount, last.OrderID, o.Amount, o.ID) {
			continue
		}
		out = append(out, PageRow{o.ID, o.Amount})
		if len(out) == limit {
			break
		}
	}
	return out
}

// --- deliberately wrong controls -----------------------------------------
//
// Each control must be *detected*: the harness asserts its answer differs from
// the contract oracle. None of them may ever be reported as a faster template.

// ControlOptionalWhere is the optional-match contract with the child predicate
// in WHERE instead of ON, which drops preserved parents. Its answer is the inner
// bag, not the optional-match bag.
func ControlOptionalWhere(fx Fixture, pred Pred) []FullRow { return FullRows(fx, pred) }

// ControlNotNullIn is `NOT IN` over a nullable child key. If any qualifying row
// has a NULL key, every comparison is UNKNOWN and the result is empty.
func ControlNotNullIn(fx Fixture, pred Pred) []int {
	for _, o := range fx.Orders {
		if pred(o) && o.CustomerID == nil {
			return []int{}
		}
	}
	return Absence(fx, pred)
}

// ControlDistinctOnName de-duplicates parents by name, silently merging two
// distinct identities that share a name.
func ControlDistinctOnName(fx Fixture, pred Pred) []string {
	names := byName(fx)
	seen := map[string]bool{}
	var out []string
	for _, id := range Membership(fx, pred) {
		n := names[id]
		if !seen[n] {
			seen[n] = true
			out = append(out, n)
		}
	}
	sort.Strings(out)
	return out
}

// ControlLimitBeforeJoin applies the limit to the child side before the join, so
// the joined result is truncated arbitrarily.
func ControlLimitBeforeJoin(fx Fixture, pred Pred, limit int) []FullRow {
	if limit <= 0 {
		return nil
	}
	names := byName(fx)
	kids := make([]Order, 0, len(fx.Orders))
	for _, o := range fx.Orders {
		if o.CustomerID != nil && pred(o) {
			kids = append(kids, o)
		}
	}
	sort.Slice(kids, func(i, j int) bool { return kids[i].ID < kids[j].ID })
	if len(kids) > limit {
		kids = kids[:limit]
	}
	var out []FullRow
	for _, o := range kids {
		out = append(out, FullRow{*o.CustomerID, names[*o.CustomerID], o.ID, o.Amount})
	}
	return out
}
