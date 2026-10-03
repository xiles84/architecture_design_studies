// Package selectbench is the additive SELECT harness for Study 01. It holds the
// embedded, named SQL catalogue, the independent Go oracles, the deterministic
// teaching/edge fixtures, and the deliberately-wrong read-path controls. It does
// not replace Study 01's historical harness: the existing `charitytree` main
// keeps its behaviour, and this package is compiled and tested separately.
package selectbench

import (
	"embed"
	"fmt"
	"sort"
	"strings"
)

//go:embed sql/*.sql
var sqlFS embed.FS

// ControlGroup is the semantic group of statements that are deliberately wrong.
// Nothing in this group may ever be compared against a valid template: its
// purpose is to be *detected* as a different answer.
const ControlGroup = "control_wrong_answer"

// Statement is one named SELECT from the embedded catalogue.
type Statement struct {
	Name   string
	Group  string
	Params []string
	SQL    string
}

// Catalog is the parsed, name-indexed catalogue in declaration order.
type Catalog struct {
	Order      []string
	Statements map[string]Statement
}

func (c *Catalog) Get(name string) (Statement, bool) {
	s, ok := c.Statements[name]
	return s, ok
}

// InGroup lists the statement names in one semantic group, in declaration order.
func (c *Catalog) InGroup(group string) []string {
	var out []string
	for _, n := range c.Order {
		if c.Statements[n].Group == group {
			out = append(out, n)
		}
	}
	return out
}

// Groups lists the distinct group ids in declaration order.
func (c *Catalog) Groups() []string {
	seen := map[string]bool{}
	var out []string
	for _, n := range c.Order {
		g := c.Statements[n].Group
		if !seen[g] {
			seen[g] = true
			out = append(out, g)
		}
	}
	return out
}

// Comparable reports whether two named statements may be pooled into one
// comparison. Statements must exist, must share a semantic group, and neither
// may be a deliberately-wrong control. This is the guard that stops a faster
// statement answering a different question from being read as a template win.
func (c *Catalog) Comparable(a, b string) (bool, error) {
	sa, ok := c.Statements[a]
	if !ok {
		return false, fmt.Errorf("unknown statement %q", a)
	}
	sb, ok := c.Statements[b]
	if !ok {
		return false, fmt.Errorf("unknown statement %q", b)
	}
	if sa.Group == ControlGroup || sb.Group == ControlGroup {
		return false, fmt.Errorf("%s and %s: a control may never be pooled with a template", a, b)
	}
	if sa.Group != sb.Group {
		return false, fmt.Errorf("%s (%s) and %s (%s): different semantic groups answer different questions",
			a, sa.Group, b, sb.Group)
	}
	return true, nil
}

// Load parses the embedded catalogue. It fails loudly on a missing group, a
// duplicate name, or a statement with no SQL body, so an unlabelled or
// ungrouped statement cannot reach a report.
func Load() (*Catalog, error) {
	entries, err := sqlFS.ReadDir("sql")
	if err != nil {
		return nil, err
	}
	// Sort by filename so declaration order never depends on the filesystem.
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
	cat := &Catalog{Statements: map[string]Statement{}}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}
		raw, err := sqlFS.ReadFile("sql/" + e.Name())
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
		return nil, fmt.Errorf("empty SELECT catalogue")
	}
	return cat, nil
}

func parseSQL(text string) ([]Statement, error) {
	var out []Statement
	var cur *Statement
	flush := func() error {
		if cur == nil {
			return nil
		}
		cur.SQL = strings.TrimSpace(cur.SQL)
		if cur.Name == "" {
			return fmt.Errorf("statement without a name")
		}
		if cur.Group == "" {
			return fmt.Errorf("statement %q has no group: an ungrouped statement cannot be compared safely", cur.Name)
		}
		if cur.SQL == "" {
			return fmt.Errorf("statement %q has no SQL body", cur.Name)
		}
		out = append(out, *cur)
		cur = nil
		return nil
	}
	for _, line := range strings.Split(text, "\n") {
		t := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(t, "-- name:"):
			if err := flush(); err != nil {
				return nil, err
			}
			cur = &Statement{Name: strings.TrimSpace(strings.TrimPrefix(t, "-- name:"))}
		case strings.HasPrefix(t, "-- group:"):
			if cur == nil {
				return nil, fmt.Errorf("group before any name")
			}
			cur.Group = strings.TrimSpace(strings.TrimPrefix(t, "-- group:"))
		case strings.HasPrefix(t, "-- params:"):
			if cur == nil {
				return nil, fmt.Errorf("params before any name")
			}
			for _, p := range strings.Split(strings.TrimPrefix(t, "-- params:"), ",") {
				if p = strings.TrimSpace(p); p != "" {
					cur.Params = append(cur.Params, p)
				}
			}
		case strings.HasPrefix(t, "--"), t == "":
			// comment or blank: keep the SQL readable, contribute nothing
		default:
			if cur == nil {
				return nil, fmt.Errorf("SQL before any -- name: %q", t)
			}
			cur.SQL += line + "\n"
		}
	}
	if err := flush(); err != nil {
		return nil, err
	}
	return out, nil
}
