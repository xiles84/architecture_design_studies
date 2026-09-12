package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ---------------------------------------------------------------------------
// EXPLAIN capture
// ---------------------------------------------------------------------------

// explainPrefix picks the richest EXPLAIN this engine understands.
//
//	PostgreSQL: BUFFERS attributes work to shared/local/temp block reads, which is
//	            the only honest way to compare access paths across designs --
//	            wall-clock alone confuses "less work" with "warmer cache".
//	YugabyteDB: DIST reports per-node RPC counts and the rows each RPC moved.
//	            On a distributed store the RPC count is the portable cost signal;
//	            milliseconds are an artefact of where the nodes happen to be.
func explainPrefix(engine string) string {
	if engine == "yugabyte" {
		return "EXPLAIN (ANALYZE, DIST, COSTS OFF, TIMING OFF)"
	}
	return "EXPLAIN (ANALYZE, BUFFERS, VERBOSE, COSTS ON)"
}

// ExplainAll captures a plan for every read query. Each query is executed twice
// and only the second plan is kept, so the recorded plan reflects a warm cache
// and a prepared statement rather than first-touch effects.
func ExplainAll(ctx context.Context, pool *pgxpool.Pool, d Design, ds *Dataset, engine string) (map[string]string, error) {
	src, err := readSQL(d.ID, "queries.sql")
	if err != nil {
		return nil, err
	}
	stmts, err := ParseCatalog(src)
	if err != nil {
		return nil, err
	}

	// Keys are fixed (so the plan is reproducible and identical across designs)
	// but chosen to be REPRESENTATIVE rather than convenient:
	//
	//   charity 1   the largest charity, 30% of all donors -- the most demanding
	//               charity-scoped case, not the easiest
	//   person      the busiest donor in that charity, so per-person plans show
	//               real work; person_id=1 might have three donations and make
	//               every design look equally good
	//   donation    one belonging to that busy donor
	//
	// The donation id matters more than it looks. In the embedded design, a
	// point lookup unnests the whole array of whichever person owns the id, so
	// its cost depends entirely on how many donations that person has. Explaining
	// with an id owned by a light donor would show a fast plan and contradict the
	// measured throughput, because random ids are disproportionately owned by
	// heavy donors -- they own more ids.
	t := computeTruth(ds, 1)
	donationID := t.donationID
	if idx := len(ds.DonationsByPerson); idx > 0 {
		for i := range ds.People {
			if ds.People[i].ID == t.personID && len(ds.DonationsByPerson[i]) > 0 {
				donationID = ds.Donations[ds.DonationsByPerson[i][len(ds.DonationsByPerson[i])/2]].ID
				break
			}
		}
	}
	vals := map[string]any{
		"charity_id":  int64(1),
		"person_id":   t.personID,
		"donation_id": donationID,
	}

	out := map[string]string{}
	for _, st := range stmts {
		args, err := bind(st.Params, vals)
		if err != nil {
			return nil, err
		}
		q := explainPrefix(engine) + " " + st.SQL
		var plan string
		for attempt := 0; attempt < 2; attempt++ {
			rows, err := pool.Query(ctx, q, args...)
			if err != nil {
				plan = "ERROR: " + err.Error()
				break
			}
			var sb strings.Builder
			for rows.Next() {
				var line string
				if err := rows.Scan(&line); err != nil {
					rows.Close()
					return nil, err
				}
				sb.WriteString(line + "\n")
			}
			rows.Close()
			if err := rows.Err(); err != nil {
				plan = "ERROR: " + err.Error()
				break
			}
			plan = sb.String()
		}
		out[st.Name] = strings.TrimRight(plan, "\n")
	}
	return out, nil
}

// ---------------------------------------------------------------------------
// Storage footprint
// ---------------------------------------------------------------------------

type RelSize struct {
	Name       string `json:"name"`
	Kind       string `json:"kind"` // table | index
	Table      string `json:"table,omitempty"`
	TotalBytes int64  `json:"total_bytes"`
	HeapBytes  int64  `json:"heap_bytes,omitempty"`
	IndexBytes int64  `json:"index_bytes,omitempty"`
}

type DBStats struct {
	RowCounts  map[string]int64 `json:"row_counts"`
	Relations  []RelSize        `json:"relations,omitempty"`
	TotalBytes int64            `json:"total_bytes"`
	Note       string           `json:"note,omitempty"`
}

// CollectStats records how much space a design costs. Storage is a first-class
// result here: a rollup or a denormalised key buys read speed with bytes, and a
// report that omits the bytes is only telling half the story.
func CollectStats(ctx context.Context, pool *pgxpool.Pool, d Design, engine string) (*DBStats, error) {
	s := &DBStats{RowCounts: map[string]int64{}}

	tables := []string{"charity", "person"}
	if !d.Embedded {
		tables = append(tables, "donation")
	}
	for _, t := range tables {
		var n int64
		if err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM "+t).Scan(&n); err != nil {
			return nil, fmt.Errorf("count %s: %w", t, err)
		}
		s.RowCounts[t] = n
	}

	if engine == "yugabyte" {
		// YugabyteDB stores data in DocDB, not in PostgreSQL heap files, so the
		// pg_*_size() family reports zero. Recording that fact is more useful
		// than recording a zero that looks like a measurement.
		s.Note = "size functions are not meaningful on YugabyteDB (storage lives in DocDB, " +
			"not PostgreSQL heap files); see the per-tablet sizes in the yb-admin output captured alongside this run"
		return s, nil
	}

	const q = `
SELECT c.relname,
       CASE c.relkind WHEN 'r' THEN 'table' WHEN 'i' THEN 'index' ELSE c.relkind::text END,
       COALESCE(t.relname, ''),
       pg_total_relation_size(c.oid),
       pg_relation_size(c.oid),
       CASE WHEN c.relkind = 'r' THEN pg_indexes_size(c.oid) ELSE 0 END
  FROM pg_class c
  JOIN pg_namespace n ON n.oid = c.relnamespace
  LEFT JOIN pg_index i ON i.indexrelid = c.oid
  LEFT JOIN pg_class t ON t.oid = i.indrelid
 WHERE n.nspname = 'public'
   AND c.relkind IN ('r', 'i')
 ORDER BY c.relkind, c.relname`

	rows, err := pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var r RelSize
		if err := rows.Scan(&r.Name, &r.Kind, &r.Table, &r.TotalBytes, &r.HeapBytes, &r.IndexBytes); err != nil {
			return nil, err
		}
		if r.Kind == "table" {
			s.TotalBytes += r.TotalBytes
		}
		s.Relations = append(s.Relations, r)
	}
	return s, rows.Err()
}

// EngineVersion records exactly what was measured. "PostgreSQL 17" is not a
// result; the full version banner is.
func EngineVersion(ctx context.Context, pool *pgxpool.Pool) string {
	var v string
	if err := pool.QueryRow(ctx, "SELECT version()").Scan(&v); err != nil {
		return "unknown: " + err.Error()
	}
	return strings.TrimSpace(v)
}
