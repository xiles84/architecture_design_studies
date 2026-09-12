package main

import (
	"fmt"
	"strings"
)

// Stmt is one named statement parsed out of a .sql catalogue file.
type Stmt struct {
	Name   string
	Params []string
	Doc    string
	SQL    string
}

// ParseCatalog reads the annotated-SQL format used throughout studies/*/sql:
//
//	-- name: q01_example
//	-- params: charity_id, person_id
//	-- free-form documentation lines
//	SELECT ...;
//
// A statement runs until the next "-- name:" line. Comment lines that appear
// after the statement body has started are kept as part of the SQL (they are
// legal SQL comments); only the leading comment block is treated as metadata.
func ParseCatalog(src string) ([]Stmt, error) {
	var out []Stmt
	var cur *Stmt
	bodyStarted := false

	flush := func() {
		if cur == nil {
			return
		}
		cur.SQL = strings.TrimSpace(cur.SQL)
		cur.SQL = strings.TrimSuffix(cur.SQL, ";")
		cur.Doc = strings.TrimSpace(cur.Doc)
		out = append(out, *cur)
		cur = nil
	}

	for _, line := range strings.Split(src, "\n") {
		line = strings.TrimRight(line, "\r")
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "-- name:") {
			flush()
			cur = &Stmt{Name: strings.TrimSpace(strings.TrimPrefix(trimmed, "-- name:"))}
			bodyStarted = false
			continue
		}
		if cur == nil {
			continue // preamble before the first -- name:
		}
		if !bodyStarted && strings.HasPrefix(trimmed, "-- params:") {
			raw := strings.TrimSpace(strings.TrimPrefix(trimmed, "-- params:"))
			if raw != "" && raw != "none" {
				for _, p := range strings.Split(raw, ",") {
					if p = strings.TrimSpace(p); p != "" {
						cur.Params = append(cur.Params, p)
					}
				}
			}
			continue
		}
		if !bodyStarted && strings.HasPrefix(trimmed, "--") {
			cur.Doc += strings.TrimSpace(strings.TrimPrefix(trimmed, "--")) + "\n"
			continue
		}
		if !bodyStarted && trimmed == "" {
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

// StmtMap indexes a catalogue by statement name.
func StmtMap(ss []Stmt) map[string]Stmt {
	m := make(map[string]Stmt, len(ss))
	for _, s := range ss {
		m[s.Name] = s
	}
	return m
}

// SplitDDL breaks a schema/index/trigger file into executable statements.
// It is deliberately simple but understands dollar-quoted bodies ($$ ... $$),
// which plain semicolon splitting would tear apart mid-function.
func SplitDDL(src string) []string {
	var out []string
	var sb strings.Builder
	inDollar := false

	for _, line := range strings.Split(src, "\n") {
		line = strings.TrimRight(line, "\r")
		if strings.Contains(line, "$$") {
			// Toggle once per occurrence so "AS $$" ... "$$ LANGUAGE" pairs up.
			for i := 0; i < strings.Count(line, "$$"); i++ {
				inDollar = !inDollar
			}
		}
		sb.WriteString(line + "\n")
		if !inDollar && strings.HasSuffix(strings.TrimSpace(line), ";") {
			if s := strings.TrimSpace(sb.String()); s != "" && !isAllComments(s) {
				out = append(out, s)
			}
			sb.Reset()
		}
	}
	if s := strings.TrimSpace(sb.String()); s != "" && !isAllComments(s) {
		out = append(out, s)
	}
	return out
}

func isAllComments(s string) bool {
	for _, line := range strings.Split(s, "\n") {
		t := strings.TrimSpace(line)
		if t != "" && !strings.HasPrefix(t, "--") {
			return false
		}
	}
	return true
}
