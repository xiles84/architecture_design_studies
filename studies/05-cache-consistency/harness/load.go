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

// A scenario's SQL is embedded in the binary, so the statements that produced a
// result are provably the statements that shipped beside it.

type catalogue struct {
	design Design
	stmts  map[string]catalog.Stmt
	order  []string

	schema   string
	indexes  string
	triggers string
}

func loadCatalogue(fsys fs.FS, d Design) (*catalogue, error) {
	own := d.OwnDir
	if own == "" {
		own = d.SQLDir
	}
	shared := d.SQLDir
	if shared == "" {
		shared = own
	}
	read := func(dir, name string) (string, error) {
		b, err := fs.ReadFile(fsys, "sql/"+dir+"/"+name)
		if err != nil {
			return "", fmt.Errorf("scenario %s: %w", d.ID, err)
		}
		return string(b), nil
	}
	cat := &catalogue{design: d}
	var err error
	if cat.schema, err = read(own, "schema.sql"); err != nil {
		return nil, err
	}
	if cat.indexes, err = read(own, "indexes.sql"); err != nil {
		return nil, err
	}
	// triggers.sql is optional: only the reference designs that maintain derived
	// data inside the database have one, and a missing file is not a missing
	// statement.
	if b, err := fs.ReadFile(fsys, "sql/"+own+"/triggers.sql"); err == nil {
		cat.triggers = string(b)
	}
	for _, f := range []string{"queries.sql", "writes.sql", "audit.sql"} {
		src, err := read(shared, f)
		if err != nil {
			return nil, err
		}
		ss, err := catalog.Parse(src)
		if err != nil {
			return nil, fmt.Errorf("scenario %s %s: %w", d.ID, f, err)
		}
		for _, s := range ss {
			if cat.stmts == nil {
				cat.stmts = map[string]catalog.Stmt{}
			}
			if _, dup := cat.stmts[s.Name]; dup {
				return nil, fmt.Errorf("scenario %s: statement %q declared twice", d.ID, s.Name)
			}
			cat.stmts[s.Name] = s
			cat.order = append(cat.order, s.Name)
		}
	}
	// A statement that uses $1 but declares no "-- params:" binds no arguments, and
	// study 04's first dev check found exactly that: two statements ran two thousand
	// times against a server answering "there is no parameter $1" while the read
	// side looked clean. The catalogue is where the mistake is visible.
	for _, name := range cat.order {
		st := cat.stmts[name]
		if len(st.Params) == 0 && usesPlaceholder(st.SQL) {
			return nil, fmt.Errorf("scenario %s statement %s uses $n but declares no -- params", d.ID, name)
		}
	}
	var missing []string
	for _, need := range d.Needs {
		if _, ok := cat.stmts[need]; !ok {
			missing = append(missing, need)
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		return nil, fmt.Errorf("scenario %s is missing statements %s", d.ID, strings.Join(missing, ", "))
	}
	sort.Strings(cat.order)
	return cat, nil
}

func (c *catalogue) stmt(name string) catalog.Stmt {
	s, ok := c.stmts[name]
	if !ok {
		panic(fmt.Sprintf("scenario %s has no statement %q", c.design.ID, name))
	}
	return s
}

func (c *catalogue) has(name string) bool { _, ok := c.stmts[name]; return ok }

func (c *catalogue) args(name string, vals map[string]any) ([]any, error) {
	s := c.stmt(name)
	a, err := catalog.Bind(s.Params, vals)
	if err != nil {
		return nil, fmt.Errorf("scenario %s statement %s: %w", c.design.ID, name, err)
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

// DropSchema removes the scenario's tables and functions. Per-cell rather than
// per-run, so a failed cell cannot poison the next one.
func DropSchema(ctx context.Context, db ports.DB) {
	for _, t := range []string{"load_donation", "cache_outbox", "donation", "person", "charity"} {
		_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS "+t+" CASCADE")
	}
	for _, fn := range []string{
		"portal_rollup()", "portal_rollup_rebuild(bigint)",
		"portal_recent()", "portal_recent_rebuild(bigint)",
	} {
		_, _ = db.Exec(ctx, "DROP FUNCTION IF EXISTS "+fn+" CASCADE")
	}
}

func ApplySchema(ctx context.Context, db ports.DB, cat *catalogue) error {
	if err := catalog.ExecScript(ctx, db, cat.schema); err != nil {
		return fmt.Errorf("schema for %s: %w", cat.design.ID, err)
	}
	return nil
}

// ApplyPostLoad attaches triggers and then indexes. Triggers first, because the
// trigger file's backfill has to see the loaded data and an index created before
// the backfill would be maintained by it for no reason; both AFTER the bulk load,
// because a per-row trigger or index during COPY measures the loader rather than
// the design.
func ApplyPostLoad(ctx context.Context, db ports.DB, cat *catalogue) error {
	if cat.triggers != "" {
		if err := catalog.ExecScript(ctx, db, cat.triggers); err != nil {
			return fmt.Errorf("triggers for %s: %w", cat.design.ID, err)
		}
	}
	if err := catalog.ExecScript(ctx, db, cat.indexes); err != nil {
		return fmt.Errorf("indexes for %s: %w", cat.design.ID, err)
	}
	return nil
}

type LoadInfo struct {
	Charities   int64   `json:"charities"`
	People      int64   `json:"people"`
	StagedRows  int64   `json:"staged_donations"`
	LoadedRows  int64   `json:"loaded_donations"`
	SchemaMS    float64 `json:"schema_ms"`
	BulkMS      float64 `json:"bulk_copy_ms"`
	LoadStmtMS  float64 `json:"load_statement_ms"`
	PostLoadMS  float64 `json:"post_load_ms"`
	HasTriggers bool    `json:"has_triggers"`
}

func BulkLoad(ctx context.Context, db ports.DB, cat *catalogue, ds *Dataset) (LoadInfo, error) {
	var li LoadInfo
	ms := func(done func() error) (float64, error) {
		t := nowMS()
		if err := done(); err != nil {
			return 0, err
		}
		return float64(nowMS() - t), nil
	}

	var err error
	li.SchemaMS, err = ms(func() error { return ApplySchema(ctx, db, cat) })
	if err != nil {
		return li, err
	}

	li.BulkMS, err = ms(func() error {
		n, err := db.CopyFrom(ctx, "charity", []string{"charity_id", "name", "country", "founded_on"},
			len(ds.Charities), func(i int) ([]any, error) {
				c := ds.Charities[i]
				return []any{c.ID, c.Name, c.Country, c.FoundedOn}, nil
			})
		if err != nil {
			return fmt.Errorf("charity: %w", err)
		}
		li.Charities = n
		n, err = db.CopyFrom(ctx, "person", []string{"person_id", "charity_id", "full_name", "email", "joined_at"},
			len(ds.People), func(i int) ([]any, error) {
				p := ds.People[i]
				return []any{p.ID, p.CharityID, p.FullName, p.Email, p.JoinedAt}, nil
			})
		if err != nil {
			return fmt.Errorf("person: %w", err)
		}
		li.People = n
		// Every donation row is staged, then moved by the design's own load
		// statement, so no design pays a per-row cost the others avoid.
		idx, personIdx := 0, 0
		cnt, err := db.CopyFrom(ctx, "load_donation",
			[]string{"donation_id", "person_id", "charity_id", "amount_cents", "currency", "donated_at", "note"},
			ds.TotalDonations, func(i int) ([]any, error) {
				for personIdx < len(ds.People) && idx >= len(ds.Donations[ds.People[personIdx].ID]) {
					idx = 0
					personIdx++
				}
				if personIdx >= len(ds.People) {
					return nil, fmt.Errorf("staging ran past the dataset")
				}
				p := ds.People[personIdx]
				d := ds.Donations[p.ID][idx]
				idx++
				var note any
				if d.Note != nil {
					note = *d.Note
				}
				return []any{d.ID, p.ID, d.CharityID, d.AmountCents, d.Currency, d.DonatedAt, note}, nil
			})
		if err != nil {
			return fmt.Errorf("load_donation: %w", err)
		}
		li.StagedRows = cnt
		return nil
	})
	if err != nil {
		return li, err
	}

	li.LoadStmtMS, err = ms(func() error {
		n, err := cat.exec(ctx, db, sLoadDonations, nil)
		li.LoadedRows = n
		return err
	})
	if err != nil {
		return li, fmt.Errorf("load statement: %w", err)
	}

	li.HasTriggers = cat.triggers != ""
	li.PostLoadMS, err = ms(func() error { return ApplyPostLoad(ctx, db, cat) })
	if err != nil {
		return li, err
	}
	return li, nil
}

// relationSizes records storage. On YugabyteDB pg_total_relation_size is not
// meaningful, so the study reports nothing there rather than a number a reader
// would take for bytes on disk.
func relationSizes(ctx context.Context, db ports.DB, engine string) map[string]int64 {
	if engine == "yugabyte" {
		return nil
	}
	out := map[string]int64{}
	for _, t := range []string{"charity", "person", "donation", "cache_outbox", "load_donation"} {
		var n int64
		err := db.QueryRow(ctx,
			"SELECT CASE WHEN to_regclass($1) IS NULL THEN 0 ELSE pg_total_relation_size(to_regclass($1)) END",
			t).Scan(&n)
		if err == nil && n > 0 {
			out[t] = n
		}
	}
	return out
}

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
