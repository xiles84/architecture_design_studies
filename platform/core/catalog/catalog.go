// Package catalog parses the annotated-SQL files every study keeps its designs
// in, and executes DDL scripts through the database port.
//
// SQL lives in readable files, embedded into each study's binary with go:embed,
// so the SQL that produced a result is provably the SQL sitting beside it.
package catalog

import (
	"context"
	"fmt"
	"strings"

	"adsplatform/ports"
)

// Stmt is one named statement parsed out of a .sql catalogue file.
type Stmt struct {
	Name   string
	Params []string
	Doc    string
	SQL    string
}

// Parse reads the catalogue format:
//
//	-- name: q01_example
//	-- params: event_id, customer_id
//	-- free-form documentation lines
//	SELECT ...;
//
// A statement runs until the next "-- name:" line. Only the leading comment
// block is metadata; comments inside the body stay part of the SQL.
func Parse(src string) ([]Stmt, error) {
	var out []Stmt
	var cur *Stmt
	bodyStarted := false

	flush := func() {
		if cur == nil {
			return
		}
		cur.SQL = strings.TrimSuffix(strings.TrimSpace(cur.SQL), ";")
		cur.Doc = strings.TrimSpace(cur.Doc)
		out = append(out, *cur)
		cur = nil
	}

	for _, line := range strings.Split(src, "\n") {
		line = strings.TrimRight(line, "\r")
		trimmed := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(trimmed, "-- name:"):
			flush()
			cur = &Stmt{Name: strings.TrimSpace(strings.TrimPrefix(trimmed, "-- name:"))}
			bodyStarted = false
			continue
		case cur == nil:
			continue
		case !bodyStarted && strings.HasPrefix(trimmed, "-- params:"):
			raw := strings.TrimSpace(strings.TrimPrefix(trimmed, "-- params:"))
			if raw != "" && raw != "none" {
				for _, p := range strings.Split(raw, ",") {
					if p = strings.TrimSpace(p); p != "" {
						cur.Params = append(cur.Params, p)
					}
				}
			}
			continue
		case !bodyStarted && strings.HasPrefix(trimmed, "--"):
			cur.Doc += strings.TrimSpace(strings.TrimPrefix(trimmed, "--")) + "\n"
			continue
		case !bodyStarted && trimmed == "":
			continue
		}
		bodyStarted = true
		cur.SQL += line + "\n"
	}
	flush()

	seen := map[string]bool{}
	for _, s := range out {
		if s.SQL == "" {
			return nil, fmt.Errorf("statement %q has no SQL body", s.Name)
		}
		if seen[s.Name] {
			return nil, fmt.Errorf("duplicate statement name %q", s.Name)
		}
		seen[s.Name] = true
	}
	return out, nil
}

// Map indexes a catalogue by statement name.
func Map(ss []Stmt) map[string]Stmt {
	m := make(map[string]Stmt, len(ss))
	for _, s := range ss {
		m[s.Name] = s
	}
	return m
}

// Bind turns a statement's "-- params:" names into positional arguments. This is
// what lets one driver execute statements whose signatures differ per design.
func Bind(params []string, vals map[string]any) ([]any, error) {
	out := make([]any, 0, len(params))
	for _, p := range params {
		v, ok := vals[p]
		if !ok {
			return nil, fmt.Errorf("no value bound for parameter %q", p)
		}
		out = append(out, v)
	}
	return out, nil
}

// SplitDDL breaks a schema/index/trigger file into executable statements. It
// understands dollar-quoted bodies ($$ ... $$), which naive semicolon splitting
// would tear apart mid-function.
func SplitDDL(src string) []string {
	var out []string
	var sb strings.Builder
	inDollar := false
	for _, line := range strings.Split(src, "\n") {
		line = strings.TrimRight(line, "\r")
		for i := 0; i < strings.Count(line, "$$"); i++ {
			inDollar = !inDollar
		}
		sb.WriteString(line + "\n")
		if !inDollar && strings.HasSuffix(strings.TrimSpace(line), ";") {
			if s := strings.TrimSpace(sb.String()); s != "" && !allComments(s) {
				out = append(out, s)
			}
			sb.Reset()
		}
	}
	if s := strings.TrimSpace(sb.String()); s != "" && !allComments(s) {
		out = append(out, s)
	}
	return out
}

func allComments(s string) bool {
	for _, line := range strings.Split(s, "\n") {
		if t := strings.TrimSpace(line); t != "" && !strings.HasPrefix(t, "--") {
			return false
		}
	}
	return true
}

// ExecScript runs every statement of a DDL file, naming the one that failed.
func ExecScript(ctx context.Context, q ports.Queryer, src string) error {
	for _, stmt := range SplitDDL(src) {
		if _, err := q.Exec(ctx, stmt); err != nil {
			head := strings.SplitN(strings.TrimSpace(stmt), "\n", 2)[0]
			return fmt.Errorf("DDL %q: %w", head, err)
		}
	}
	return nil
}
