// Command appendixsql writes the SQL behind every plan the book's query-plan
// appendix quotes. It starts no database.
//
// A result cell stores its statement name and the parameters it was bound with
// (`k=v,...`), not the SQL text. This command rebuilds the text with the
// harness's own code path: the embedded catalogue (selectbench.Load, or
// selectbench.LoadJoin, which also generates the CTE-reuse statements) and
// selectbench.Bind. The result is the statement the harness sent after
// `EXPLAIN (...)`, not a retyped copy.
//
// The SQL is only as old as this tree's catalogue code, so the output names the
// source files each catalogue is built from. `tools/evidence plans` checks those
// files against every exhibit's producing commit and refuses a difference.
//
//	go run ./selectbench/cmd/appendixsql -repo ../..        (from studies/01-charity-tree)
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"charitytree/selectbench"
)

const study = "studies/01-charity-tree"

type spec struct {
	Exhibits []struct {
		RunID string `json:"run_id"`
		File  string `json:"file"`
		Plans []struct {
			Axis string `json:"axis"`
		} `json:"plans"`
		Shapes []struct {
			Axis  string `json:"axis"`
			As    string `json:"as"`
			AsRun string `json:"as_run"`
		} `json:"shapes"`
	} `json:"exhibits"`
}

type cell struct {
	Axis      string `json:"Axis"`
	Statement string `json:"Statement"`
	Params    string `json:"Params"`
}

// Entry is one bound statement, keyed by the cell it was executed for.
type Entry struct {
	RunID     string `json:"run_id"`
	File      string `json:"file"`
	Axis      string `json:"axis"`
	Statement string `json:"statement"`
	Params    string `json:"params"`
	Catalogue string `json:"catalogue"`
	SQL       string `json:"sql"`
}

// Output is the generated file.
type Output struct {
	Schema      string              `json:"schema"`
	GeneratedBy string              `json:"generated_by"`
	Catalogues  map[string][]string `json:"catalogues"`
	Statements  []Entry             `json:"statements"`
}

// catalogueSources lists the repo files whose bytes determine a catalogue's
// bound SQL: the embedded .sql files, the parser, the binder and, for the join
// catalogue, the reuse-statement generator and the expression it embeds.
func catalogueSources(root, name string) ([]string, error) {
	sb := study + "/selectbench/"
	var glob string
	var code []string
	switch name {
	case "teaching":
		glob, code = sb+"sql/*.sql", []string{sb + "catalog.go", sb + "bind.go"}
	case "join":
		glob, code = sb+"sql/join/*.sql", []string{sb + "catalog.go", sb + "bind.go", sb + "joincatalog.go", sb + "joinfixture.go"}
	default:
		return nil, fmt.Errorf("unknown catalogue %q", name)
	}
	files, err := filepath.Glob(filepath.Join(root, filepath.FromSlash(glob)))
	if err != nil || len(files) == 0 {
		return nil, fmt.Errorf("catalogue %s: no files match %s", name, glob)
	}
	var out []string
	for _, f := range files {
		rel, err := filepath.Rel(root, f)
		if err != nil {
			return nil, err
		}
		out = append(out, filepath.ToSlash(rel))
	}
	out = append(out, code...)
	sort.Strings(out)
	return out, nil
}

// parseParams inverts the harness's describeParams (sorted `k=v` joined by
// commas) and proves the inversion by rebuilding the string: a value holding a
// comma or a lost key fails here instead of binding a different statement.
func parseParams(s string) (map[string]string, error) {
	m := map[string]string{}
	if s == "" {
		return m, nil
	}
	for _, kv := range strings.Split(s, ",") {
		k, v, ok := strings.Cut(kv, "=")
		if !ok || k == "" {
			return nil, fmt.Errorf("parameter %q is not k=v", kv)
		}
		if _, dup := m[k]; dup {
			return nil, fmt.Errorf("parameter %q repeats", k)
		}
		m[k] = v
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var parts []string
	for _, k := range keys {
		parts = append(parts, k+"="+m[k])
	}
	if strings.Join(parts, ",") != s {
		return nil, fmt.Errorf("parameters %q do not round-trip; a value may contain a comma", s)
	}
	return m, nil
}

func main() {
	repo := flag.String("repo", ".", "repository root")
	specRel := flag.String("spec", "book/appendix/query-plans.spec.json", "appendix exhibit spec (repo-relative)")
	outRel := flag.String("out", "book/appendix/query-plans.sql.json", "output (repo-relative)")
	flag.Parse()
	if err := run(*repo, *specRel, *outRel); err != nil {
		fmt.Fprintln(os.Stderr, "appendixsql:", err)
		os.Exit(1)
	}
}

func run(repo, specRel, outRel string) error {
	root, err := filepath.Abs(repo)
	if err != nil {
		return err
	}
	raw, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(specRel)))
	if err != nil {
		return err
	}
	var sp spec
	if err := json.Unmarshal(raw, &sp); err != nil {
		return fmt.Errorf("%s: %w", specRel, err)
	}
	teaching, err := selectbench.Load()
	if err != nil {
		return err
	}
	join, err := selectbench.LoadJoin()
	if err != nil {
		return err
	}
	out := Output{Schema: "ads-plan-sql/1", GeneratedBy: study + "/selectbench/cmd/appendixsql", Catalogues: map[string][]string{}}
	for _, name := range []string{"teaching", "join"} {
		if out.Catalogues[name], err = catalogueSources(root, name); err != nil {
			return err
		}
	}

	// Every cell a plan or a node-tree check names, once.
	type key struct{ run, file, axis string }
	want := map[key]bool{}
	for _, e := range sp.Exhibits {
		for _, p := range e.Plans {
			want[key{e.RunID, e.File, p.Axis}] = true
		}
		for _, s := range e.Shapes {
			want[key{e.RunID, e.File, s.Axis}] = true
			other := e.RunID
			if s.AsRun != "" {
				other = s.AsRun
			}
			want[key{other, e.File, s.As}] = true
		}
	}
	runs := map[string]map[string]cell{}
	for k := range want {
		path := filepath.Join(root, filepath.FromSlash(study+"/results"), k.run, k.file)
		cells, ok := runs[path]
		if !ok {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			var r struct {
				Cells []cell `json:"cells"`
			}
			if err := json.Unmarshal(data, &r); err != nil {
				return fmt.Errorf("%s: %w", path, err)
			}
			cells = map[string]cell{}
			for _, c := range r.Cells {
				cells[c.Axis] = c
			}
			runs[path] = cells
		}
		c, ok := cells[k.axis]
		if !ok {
			return fmt.Errorf("%s/%s: no cell with axis %q", k.run, k.file, k.axis)
		}
		catName, cat := "teaching", teaching
		if k.file == "join.json" {
			catName, cat = "join", join
		}
		st, ok := cat.Get(c.Statement)
		if !ok {
			return fmt.Errorf("%s/%s %q: statement %s is not in the %s catalogue", k.run, k.file, k.axis, c.Statement, catName)
		}
		params, err := parseParams(c.Params)
		if err != nil {
			return fmt.Errorf("%s/%s %q: %w", k.run, k.file, k.axis, err)
		}
		bound, _, err := selectbench.Bind(st.SQL, params)
		if err != nil {
			return fmt.Errorf("%s/%s %q: bind %s: %w", k.run, k.file, k.axis, c.Statement, err)
		}
		out.Statements = append(out.Statements, Entry{RunID: k.run, File: k.file, Axis: k.axis,
			Statement: c.Statement, Params: c.Params, Catalogue: catName, SQL: bound})
	}
	sort.Slice(out.Statements, func(i, j int) bool {
		a, b := out.Statements[i], out.Statements[j]
		if a.RunID != b.RunID {
			return a.RunID < b.RunID
		}
		if a.File != b.File {
			return a.File < b.File
		}
		return a.Axis < b.Axis
	})
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(outRel)), buf.Bytes(), 0o644); err != nil {
		return err
	}
	fmt.Printf("appendixsql: wrote %s (%d statements)\n", outRel, len(out.Statements))
	return nil
}
