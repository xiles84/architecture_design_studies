package selectbench

import (
	"embed"
	"fmt"
	"sort"
	"strings"
)

// The join/CTE catalogue: the embedded base spellings in sql/join/, plus the
// generated reuse ladder. It is a separate catalogue from Load() so the
// teaching-fixture families keep exactly the catalogue they were validated on.

//go:embed sql/join/*.sql
var joinSQLFS embed.FS

// Semantic groups of the join/CTE family.
const (
	GroupJoinBag      = "sel_join_bag"
	GroupJoinOptional = "sel_join_optional"
	GroupJoinReuse    = "sel_join_reuse"
)

// ReuseLevels are the CTE reference counts. Five levels are the minimum a
// bounded crossing may be claimed from; the original {1,8} pair could not
// bound one, which is why HIGH approved the expansion.
func ReuseLevels() []int { return []int{1, 2, 4, 8, 16} }

// ReuseConsumerPercents are the per-consumer selection levels, as a percentage
// of the fact id space, applied OUTSIDE the common CTE.
func ReuseConsumerPercents() []int { return []int{1, 100} }

// reusePolicies maps a generated statement's policy id to its SQL keyword. The
// default spelling carries no keyword: folding is the planner's decision.
var reusePolicies = []struct{ id, keyword string }{
	{"default", ""},
	{"mat", "MATERIALIZED"},
	{"notmat", "NOT MATERIALIZED"},
}

// joinPolicyFor names the CTE materialization policy a base statement forces,
// empty for statements that have no CTE.
func joinPolicyFor(name string) string {
	switch name {
	case "j04_cte_filtered", "j05_cte_joined_outer":
		return "default"
	case "j04m_cte_filtered_mat", "j05m_cte_joined_outer_mat":
		return "materialized"
	case "j04n_cte_filtered_notmat", "j05n_cte_joined_outer_notmat":
		return "not_materialized"
	}
	return ""
}

// ReuseStatementName is the generated reuse statement's catalogue name.
func ReuseStatementName(policy string, refs, consumerPct int) string {
	return fmt.Sprintf("r_%s_n%d_c%d", policy, refs, consumerPct)
}

// buildReuseSQL builds the reuse ladder for one policy and reference count. The
// CTE computes the expensive token over the threshold-selected facts; each
// reference then consumes it under the same per-consumer filter, so materialize
// versus inline is the only difference between policies.
func buildReuseSQL(policy string, refs int) string {
	var kw string
	for _, p := range reusePolicies {
		if p.id == policy {
			kw = p.keyword
		}
	}
	head := "WITH e AS ("
	if kw != "" {
		head = "WITH e AS " + kw + " ("
	}
	var b strings.Builder
	b.WriteString(head + "\n")
	b.WriteString("  SELECT f.id AS consumer_id, " + ReuseExpr + " AS cksum\n")
	b.WriteString("  FROM facts f WHERE f.rank <= :threshold\n")
	b.WriteString(")\n")
	for r := 1; r <= refs; r++ {
		if r > 1 {
			b.WriteString("UNION ALL\n")
		}
		fmt.Fprintf(&b, "SELECT %d AS ref, COALESCE(sum(cksum), 0)::bigint AS checksum FROM e WHERE e.consumer_id <= :consumer_max\n", r)
	}
	b.WriteString("ORDER BY ref")
	return b.String()
}

func addReuseStatements(cat *Catalog) {
	for _, p := range reusePolicies {
		for _, refs := range ReuseLevels() {
			for _, pct := range ReuseConsumerPercents() {
				name := ReuseStatementName(p.id, refs, pct)
				if _, dup := cat.Statements[name]; dup {
					continue
				}
				cat.Statements[name] = Statement{
					Name:   name,
					Group:  GroupJoinReuse,
					Params: []string{"threshold", "consumer_max"},
					SQL:    buildReuseSQL(p.id, refs),
				}
				cat.Order = append(cat.Order, name)
			}
		}
	}
}

// LoadJoin parses the join/CTE catalogue and appends the generated reuse
// statements. It fails loudly on a duplicate or ungrouped statement, exactly as
// Load does for the teaching catalogue.
func LoadJoin() (*Catalog, error) {
	entries, err := joinSQLFS.ReadDir("sql/join")
	if err != nil {
		return nil, err
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
	cat := &Catalog{Statements: map[string]Statement{}}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}
		raw, err := joinSQLFS.ReadFile("sql/join/" + e.Name())
		if err != nil {
			return nil, err
		}
		stmts, err := parseSQL(string(raw))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", e.Name(), err)
		}
		for _, s := range stmts {
			if _, dup := cat.Statements[s.Name]; dup {
				return nil, fmt.Errorf("duplicate statement %q", s.Name)
			}
			cat.Statements[s.Name] = s
			cat.Order = append(cat.Order, s.Name)
		}
	}
	if len(cat.Order) == 0 {
		return nil, fmt.Errorf("empty join catalogue")
	}
	addReuseStatements(cat)
	return cat, nil
}

// --- canonicalization and verification -------------------------------------

func joinNeed(cols []any, n int, group string) error {
	if len(cols) != n {
		return fmt.Errorf("%s: expected %d columns, got %d", group, n, len(cols))
	}
	return nil
}

// JoinCanonicalize converts scanned rows into the canonical value of a join
// contract. It preserves multiplicities (the contracts are bags) and sorts each
// contract exactly as its oracle does.
func JoinCanonicalize(group string, rows [][]any) (any, error) {
	switch group {
	case GroupJoinBag:
		out := make([]JoinRow, 0, len(rows))
		for _, r := range rows {
			if err := joinNeed(r, 4, group); err != nil {
				return nil, err
			}
			cid, _, err := toInt(r[0])
			if err != nil {
				return nil, err
			}
			name, err := toString(r[1])
			if err != nil {
				return nil, err
			}
			fid, _, err := toInt(r[2])
			if err != nil {
				return nil, err
			}
			amt, anull, err := toInt(r[3])
			if err != nil {
				return nil, err
			}
			var amount *int
			if !anull {
				amount = ip(amt)
			}
			out = append(out, JoinRow{cid, name, fid, amount})
		}
		sort.SliceStable(out, func(i, j int) bool {
			if out[i].CustomerID != out[j].CustomerID {
				return out[i].CustomerID < out[j].CustomerID
			}
			return out[i].FactID < out[j].FactID
		})
		return out, nil

	case GroupJoinOptional:
		out := make([]JoinOptRow, 0, len(rows))
		for _, r := range rows {
			if err := joinNeed(r, 4, group); err != nil {
				return nil, err
			}
			cid, _, err := toInt(r[0])
			if err != nil {
				return nil, err
			}
			name, err := toString(r[1])
			if err != nil {
				return nil, err
			}
			fid, fnull, err := toInt(r[2])
			if err != nil {
				return nil, err
			}
			amt, anull, err := toInt(r[3])
			if err != nil {
				return nil, err
			}
			row := JoinOptRow{CustomerID: cid, CustomerName: name}
			if !fnull {
				row.FactID = ip(fid)
			}
			if !anull {
				row.Amount = ip(amt)
			}
			out = append(out, row)
		}
		sort.SliceStable(out, func(i, j int) bool {
			if out[i].CustomerID != out[j].CustomerID {
				return out[i].CustomerID < out[j].CustomerID
			}
			if out[i].FactID == nil || out[j].FactID == nil {
				return out[i].FactID == nil && out[j].FactID != nil
			}
			return *out[i].FactID < *out[j].FactID
		})
		return out, nil

	case GroupJoinReuse:
		out := make([]ReuseRow, 0, len(rows))
		for _, r := range rows {
			if err := joinNeed(r, 2, group); err != nil {
				return nil, err
			}
			ref, _, err := toInt(r[0])
			if err != nil {
				return nil, err
			}
			cksum, cnull, err := toInt(r[1])
			if err != nil {
				return nil, err
			}
			if cnull {
				return nil, fmt.Errorf("%s: NULL checksum where the contract requires an integer", group)
			}
			out = append(out, ReuseRow{Ref: ref, Checksum: int64(cksum)})
		}
		sort.SliceStable(out, func(i, j int) bool { return out[i].Ref < out[j].Ref })
		return out, nil

	case ControlGroup:
		return nil, fmt.Errorf("a deliberately wrong control has no canonical answer to verify against")
	}
	return nil, fmt.Errorf("unknown join group %q", group)
}

// JoinExpectedAnswer is the oracle answer for one join statement.
func JoinExpectedAnswer(c *Catalog, name string, jf JoinFixture, a JoinArgs) (any, error) {
	s, ok := c.Statements[name]
	if !ok {
		return nil, fmt.Errorf("unknown statement %q", name)
	}
	switch s.Group {
	case GroupJoinBag:
		return JoinBag(jf, a.Threshold), nil
	case GroupJoinOptional:
		return JoinOptional(jf, a.Threshold), nil
	case GroupJoinReuse:
		return ReuseAnswer(jf, a.Threshold, a.Refs, a.ConsumerMax), nil
	case ControlGroup:
		return nil, fmt.Errorf("a deliberately wrong control has no oracle: it must be seen to FAIL, not to match")
	}
	return nil, fmt.Errorf("unknown join group %q", s.Group)
}

// VerifyJoinRows canonicalizes and verifies one join cell's rows.
func VerifyJoinRows(c *Catalog, name string, jf JoinFixture, a JoinArgs, rows [][]any) error {
	s, ok := c.Statements[name]
	if !ok {
		return fmt.Errorf("unknown statement %q", name)
	}
	got, err := JoinCanonicalize(s.Group, rows)
	if err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	want, err := JoinExpectedAnswer(c, name, jf, a)
	if err != nil {
		return err
	}
	if !answersEqual(want, got) {
		return fmt.Errorf("%s (%s): observed answer does not match the contract oracle:\n got %#v\nwant %#v",
			name, s.Group, got, want)
	}
	return nil
}

// JoinControlContracts names the contract each join control pretends to serve.
var joinControlContracts = map[string]string{
	"x01_optional_where":    GroupJoinOptional,
	"x02_join_limit_before": GroupJoinBag,
}

// JoinControlNames lists the join controls in catalogue order.
func JoinControlNames(c *Catalog) []string {
	var out []string
	for _, n := range c.Order {
		if _, ok := joinControlContracts[n]; ok {
			out = append(out, n)
		}
	}
	return out
}

// JoinDetectControl requires the control's observed answer to DIFFER from the
// oracle of the contract it pretends to serve. nil means it fired.
func JoinDetectControl(c *Catalog, name string, jf JoinFixture, a JoinArgs, rows [][]any) error {
	group, ok := joinControlContracts[name]
	if !ok {
		return fmt.Errorf("unknown join control %s", name)
	}
	s, ok := c.Statements[name]
	if !ok || s.Group != ControlGroup {
		return fmt.Errorf("%s is not a deliberately wrong control", name)
	}
	got, err := JoinCanonicalize(group, rows)
	if err != nil {
		return err
	}
	var want any
	switch name {
	case "x01_optional_where":
		want, err = JoinOptional(jf, a.Threshold), error(nil)
	case "x02_join_limit_before":
		want = ControlJoinLimitBefore(jf, a.Threshold, 7)
	}
	if err != nil {
		return err
	}
	if answersEqual(want, got) {
		return fmt.Errorf("%s was NOT detected: its answer matched the %s oracle", name, group)
	}
	return nil
}
