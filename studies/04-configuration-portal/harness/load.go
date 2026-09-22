package main

import (
	"context"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"adsplatform/core/catalog"
	"adsplatform/ports"
)

// A design's SQL is embedded in the binary, so the statements that produced a
// result are provably the statements that shipped beside it.
type catalogue struct {
	design Design
	stmts  map[string]catalog.Stmt
	order  []string
	schema string
	index  string
}

func loadCatalogue(fsys fs.FS, d Design) (*catalogue, error) {
	read := func(name string) (string, error) {
		b, err := fs.ReadFile(fsys, "sql/"+d.ID+"/"+name)
		if err != nil {
			return "", fmt.Errorf("design %s: %w", d.ID, err)
		}
		return string(b), nil
	}
	cat := &catalogue{design: d}
	var err error
	if cat.schema, err = read("schema.sql"); err != nil {
		return nil, err
	}
	if cat.index, err = read("indexes.sql"); err != nil {
		return nil, err
	}
	for _, f := range []string{"queries.sql", "writes.sql", "audit.sql"} {
		src, err := read(f)
		if err != nil {
			return nil, err
		}
		ss, err := catalog.Parse(src)
		if err != nil {
			return nil, fmt.Errorf("design %s %s: %w", d.ID, f, err)
		}
		for _, s := range ss {
			if _, dup := cat.stmts[s.Name]; dup && cat.stmts != nil {
				return nil, fmt.Errorf("design %s: statement %q declared twice", d.ID, s.Name)
			}
			if cat.stmts == nil {
				cat.stmts = map[string]catalog.Stmt{}
			}
			cat.stmts[s.Name] = s
			cat.order = append(cat.order, s.Name)
		}
	}
	// A statement that uses $1 but declares no "-- params:" binds no arguments,
	// and the first dev check found exactly that: two statements ran 2263 times
	// against a database that answered "there is no parameter $1" every time. The
	// catalogue is the place the mistake is visible, so it is caught here.
	for _, name := range cat.order {
		st := cat.stmts[name]
		if len(st.Params) == 0 && usesPlaceholder(st.SQL) {
			return nil, fmt.Errorf("design %s statement %s uses $n but declares no -- params", d.ID, name)
		}
	}

	// Every statement the registry says the harness will drive must exist. A
	// missing one is a harness bug: failing here keeps it from silently turning
	// into a design that does less than its name claims.
	var missing []string
	for _, need := range d.Needs {
		if _, ok := cat.stmts[need]; !ok {
			missing = append(missing, need)
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		return nil, fmt.Errorf("design %s is missing statements %s", d.ID, strings.Join(missing, ", "))
	}
	sort.Strings(cat.order)
	return cat, nil
}

// stmt returns a statement or panics: loadCatalogue already proved it exists.
func (c *catalogue) stmt(name string) catalog.Stmt {
	s, ok := c.stmts[name]
	if !ok {
		panic(fmt.Sprintf("design %s has no statement %q", c.design.ID, name))
	}
	return s
}

// args binds a statement's declared parameters. Using the catalogue's own
// parameter list is what lets one harness drive statements whose signatures
// differ between designs.
func (c *catalogue) args(name string, vals map[string]any) ([]any, error) {
	s := c.stmt(name)
	a, err := catalog.Bind(s.Params, vals)
	if err != nil {
		return nil, fmt.Errorf("design %s statement %s: %w", c.design.ID, name, err)
	}
	return a, nil
}

func (c *catalogue) exec(ctx context.Context, q ports.Queryer, name string, vals map[string]any) (int64, error) {
	a, err := c.args(name, vals)
	if err != nil {
		return 0, err
	}
	n, err := q.Exec(ctx, c.stmt(name).SQL, a...)
	if err != nil {
		return n, fmt.Errorf("%s: %w", name, err)
	}
	return n, nil
}

// ApplySchema creates the tables before the bulk load. Indexes (and, for the
// trigger design, triggers) are applied afterwards, because a per-row index
// maintained during a COPY measures the loader rather than the design.
func ApplySchema(ctx context.Context, db ports.DB, cat *catalogue) error {
	if err := catalog.ExecScript(ctx, db, cat.schema); err != nil {
		return fmt.Errorf("schema for %s: %w", cat.design.ID, err)
	}
	return nil
}

func ApplyIndexes(ctx context.Context, db ports.DB, cat *catalogue) error {
	if err := catalog.ExecScript(ctx, db, cat.index); err != nil {
		return fmt.Errorf("indexes for %s: %w", cat.design.ID, err)
	}
	return nil
}

// DropSchema removes the design's tables so a cell starts from nothing. Dropping
// is per-cell rather than per-run so a failed cell cannot poison the next one.
func DropSchema(ctx context.Context, db ports.DB) {
	for _, t := range []string{"load_entry", "config_entry", "config_document", "installed_product",
		"business_unit", "environment", "product_definition"} {
		_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS "+t+" CASCADE")
	}
	_, _ = db.Exec(ctx, "DROP FUNCTION IF EXISTS refresh_ip_rollup() CASCADE")
}

// BulkLoad writes the logical dataset. Reference tables go in directly; the
// configuration itself goes through the load_entry staging table and then one
// design-specific statement, so no design pays a per-row cost the others avoid.
type LoadInfo struct {
	Definitions   int64   `json:"definitions"`
	Environments  int64   `json:"environments"`
	BusinessUnits int64   `json:"business_units"`
	Installations int64   `json:"installations"`
	StagedEntries int64   `json:"staged_entries"`
	LoadedEntries int64   `json:"loaded_entries"`
	SchemaMS      float64 `json:"schema_ms"`
	StagingMS     float64 `json:"staging_ms"`
	LoadStmtMS    float64 `json:"load_statement_ms"`
	IndexMS       float64 `json:"index_ms"`
}

func BulkLoad(ctx context.Context, db ports.DB, cat *catalogue, ds *Dataset) (LoadInfo, error) {
	var li LoadInfo
	ms := func(done func() error) (float64, error) {
		t := nowMS()
		if err := done(); err != nil {
			return 0, err
		}
		return nowMS() - t, nil
	}

	li.Definitions, _ = db.CopyFrom(ctx, "product_definition", []string{"id", "name", "product_version"},
		len(ds.Definitions), func(i int) ([]any, error) {
			d := ds.Definitions[i]
			return []any{d.ID, d.Name, d.Version}, nil
		})
	li.Environments, _ = db.CopyFrom(ctx, "environment", []string{"id", "name"},
		len(ds.Environments), func(i int) ([]any, error) {
			e := ds.Environments[i]
			return []any{e.ID, e.Name}, nil
		})
	li.BusinessUnits, _ = db.CopyFrom(ctx, "business_unit", []string{"id", "name"},
		len(ds.BusinessUnits), func(i int) ([]any, error) {
			b := ds.BusinessUnits[i]
			return []any{b.ID, b.Name}, nil
		})
	n, err := db.CopyFrom(ctx, "installed_product",
		[]string{"id", "product_definition_id", "environment_id", "business_unit_id", "display_name"},
		len(ds.Installations), func(i int) ([]any, error) {
			ip := ds.Installations[i]
			var bu any
			if ip.BUID != 0 {
				bu = ip.BUID
			}
			return []any{ip.ID, ip.ProductDefID, ip.EnvID, bu, ip.DisplayName}, nil
		})
	if err != nil {
		return li, fmt.Errorf("installed_product: %w", err)
	}
	li.Installations = n

	// Stage every (installation, key, value) row.
	var stageErr error
	li.StagingMS, stageErr = ms(func() error {
		var total int
		for _, ip := range ds.Installations {
			total += len(ds.Entries[ip.ID])
		}
		idx, ipIdx := 0, 0
		cnt, err := db.CopyFrom(ctx, "load_entry", []string{"installed_product_id", "key", "value"}, total,
			func(i int) ([]any, error) {
				for idx >= len(ds.Entries[ds.Installations[ipIdx].ID]) {
					idx = 0
					ipIdx++
				}
				ip := ds.Installations[ipIdx]
				e := ds.Entries[ip.ID][idx]
				idx++
				return []any{ip.ID, e.Key, e.Value}, nil
			})
		li.StagedEntries = cnt
		return err
	})
	if stageErr != nil {
		return li, fmt.Errorf("staging load: %w", stageErr)
	}

	loadStmt := "w_load_entries"
	if cat.design.Kind == Doc {
		loadStmt = "w_load_documents"
	}
	li.LoadStmtMS, err = ms(func() error {
		n, err := cat.exec(ctx, db, loadStmt, nil)
		li.LoadedEntries = n
		return err
	})
	if err != nil {
		return li, fmt.Errorf("load statement: %w", err)
	}

	li.IndexMS, err = ms(func() error { return ApplyIndexes(ctx, db, cat) })
	if err != nil {
		return li, err
	}
	return li, nil
}

// usesPlaceholder reports whether a statement's SQL references a positional
// parameter. The leading comment block is stripped first: a "--" line mentioning
// $1 is documentation, not a binding.
func usesPlaceholder(sql string) bool {
	for _, line := range strings.Split(sql, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "--") {
			continue
		}
		for i := 0; i+1 < len(line); i++ {
			if line[i] == '$' && line[i+1] >= '0' && line[i+1] <= '9' {
				return true
			}
		}
	}
	return false
}
