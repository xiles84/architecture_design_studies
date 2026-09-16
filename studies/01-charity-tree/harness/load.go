package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	charitytree "charitytree"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func readSQL(designID, file string) (string, error) {
	b, err := charitytree.SQL.ReadFile("sql/" + designID + "/" + file)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func readSQLOptional(designID, file string) string {
	s, err := readSQL(designID, file)
	if err != nil {
		return ""
	}
	return s
}

// execDDL runs every statement of a DDL file, reporting which one failed.
func execDDL(ctx context.Context, pool *pgxpool.Pool, src string) error {
	for _, stmt := range SplitDDL(src) {
		if _, err := pool.Exec(ctx, stmt); err != nil {
			head := strings.SplitN(strings.TrimSpace(stmt), "\n", 2)[0]
			return fmt.Errorf("DDL %q: %w", head, err)
		}
	}
	return nil
}

// Drop removes everything the study creates, so a run always starts from a known
// state rather than from whatever the previous run left behind.
func Drop(ctx context.Context, pool *pgxpool.Pool) error {
	stmts := []string{
		`DROP TABLE IF EXISTS donation CASCADE`,
		`DROP TABLE IF EXISTS person CASCADE`,
		`DROP TABLE IF EXISTS charity CASCADE`,
		`DROP FUNCTION IF EXISTS donation_rollup() CASCADE`,
		`DROP FUNCTION IF EXISTS donation_rollup_apply(BIGINT, BIGINT, BOOLEAN) CASCADE`,
	}
	for _, s := range stmts {
		if _, err := pool.Exec(ctx, s); err != nil {
			return fmt.Errorf("drop: %s: %w", s, err)
		}
	}
	return nil
}

// LoadPhases records how long each stage of setup took. Load time is itself a
// design-relevant number: a schema that takes four times as long to populate is
// telling you something about its write path before any benchmark runs.
type LoadPhases struct {
	DropMS    float64 `json:"drop_ms"`
	SchemaMS  float64 `json:"schema_ms"`
	CopyMS    float64 `json:"copy_ms"`
	RollupMS  float64 `json:"rollup_ms"`
	IndexMS   float64 `json:"index_ms"`
	TriggerMS float64 `json:"trigger_ms"`
	AnalyzeMS float64 `json:"analyze_ms"`
	TotalMS   float64 `json:"total_ms"`
}

func msSince(t time.Time) float64 { return float64(time.Since(t).Microseconds()) / 1000.0 }

// SetupAndLoad brings a database from empty to benchmark-ready.
//
// Order matters and is the same for every design:
//
//	drop -> schema -> bulk COPY -> bulk rollup -> indexes -> triggers -> ANALYZE
//
// Indexes are built after the data lands because that is both faster and what a
// real migration does. Triggers are attached after the bulk rollup for the same
// reason: making the loader pay per-row trigger cost for data that a single
// aggregate query can compute would measure the loader, not the design. The
// per-row cost of those triggers is what the WRITE benchmark exists to measure.
func SetupAndLoad(ctx context.Context, pool *pgxpool.Pool, d Design, ds *Dataset, engine string, loadConns int) (*LoadPhases, error) {
	ph := &LoadPhases{}
	all := time.Now()

	t := time.Now()
	if err := Drop(ctx, pool); err != nil {
		return nil, err
	}
	ph.DropMS = msSince(t)

	t = time.Now()
	schema, err := readSQL(d.ID, "schema.sql")
	if err != nil {
		return nil, err
	}
	if err := execDDL(ctx, pool, schema); err != nil {
		return nil, fmt.Errorf("schema: %w", err)
	}
	ph.SchemaMS = msSince(t)

	t = time.Now()
	if err := copyData(ctx, pool, d, ds, loadConns); err != nil {
		return nil, fmt.Errorf("copy: %w", err)
	}
	ph.CopyMS = msSince(t)

	if d.Rollups {
		t = time.Now()
		if err := bulkRollup(ctx, pool); err != nil {
			return nil, fmt.Errorf("rollup: %w", err)
		}
		ph.RollupMS = msSince(t)
	}
	if d.SumOnly {
		t = time.Now()
		for _, table := range []string{"person", "charity"} {
			key := table + "_id"
			if _, err := pool.Exec(ctx, "UPDATE "+table+" p SET total_donated_cents = a.total FROM (SELECT "+key+", SUM(amount_cents) AS total FROM donation GROUP BY "+key+") a WHERE p."+key+" = a."+key); err != nil {
				return nil, fmt.Errorf("sum rollup: %w", err)
			}
		}
		ph.RollupMS = msSince(t)
	}

	t = time.Now()
	if idx := readSQLOptional(d.ID, "indexes.sql"); strings.TrimSpace(idx) != "" {
		if err := execDDL(ctx, pool, idx); err != nil {
			return nil, fmt.Errorf("indexes: %w", err)
		}
	}
	ph.IndexMS = msSince(t)

	if d.Triggers {
		t = time.Now()
		trg := readSQLOptional(d.ID, "triggers.sql")
		if err := execDDL(ctx, pool, trg); err != nil {
			return nil, fmt.Errorf("triggers: %w", err)
		}
		ph.TriggerMS = msSince(t)
	}

	t = time.Now()
	tables := []string{"charity", "person"}
	if !d.Embedded {
		tables = append(tables, "donation")
	}
	for _, tb := range tables {
		if _, err := pool.Exec(ctx, "ANALYZE "+tb); err != nil {
			return nil, fmt.Errorf("analyze %s: %w", tb, err)
		}
	}
	ph.AnalyzeMS = msSince(t)
	ph.TotalMS = msSince(all)
	return ph, nil
}

func copyData(ctx context.Context, pool *pgxpool.Pool, d Design, ds *Dataset, loadConns int) error {
	// charity
	cols := []string{"charity_id", "name", "country", "founded_on"}
	_, err := pool.CopyFrom(ctx, pgx.Identifier{"charity"}, cols,
		pgx.CopyFromSlice(len(ds.Charities), func(i int) ([]any, error) {
			c := ds.Charities[i]
			return []any{c.ID, c.Name, c.Country, c.FoundedOn}, nil
		}))
	if err != nil {
		return fmt.Errorf("charity: %w", err)
	}

	if d.Embedded {
		return copyPeopleEmbedded(ctx, pool, ds, loadConns)
	}

	pcols := []string{"person_id", "charity_id", "full_name", "email", "joined_at"}
	if d.RecentCache {
		pcols = append(pcols, "recent_donations")
	}
	if _, err := pool.CopyFrom(ctx, pgx.Identifier{"person"}, pcols,
		pgx.CopyFromSlice(len(ds.People), func(i int) ([]any, error) {
			p := ds.People[i]
			row := []any{p.ID, p.CharityID, p.FullName, p.Email, p.JoinedAt}
			if d.RecentCache {
				b, err := recentCacheJSON(ds, i)
				if err != nil {
					return nil, err
				}
				row = append(row, b)
			}
			return row, nil
		})); err != nil {
		return fmt.Errorf("person: %w", err)
	}

	dcols := []string{"donation_id", "person_id", "amount_cents", "currency", "donated_at", "note"}
	if d.CharityOnDonation {
		dcols = []string{"donation_id", "person_id", "charity_id", "amount_cents", "currency", "donated_at", "note"}
	}
	if d.LastFlag {
		dcols = append(dcols, "is_last_donation")
	}
	// lastFlag marks, for LastFlag designs, exactly the donations that are
	// their donor's newest. The loader already knows this from
	// DonationsByPerson (chronological, so the newest is the last index) --
	// computing it here means the bulk load sets the flag directly, the way
	// D4's rollups are loaded, instead of paying per-row trigger cost that
	// would measure the loader rather than the design (RECENCY.md section 3).
	var lastFlag []bool
	if d.LastFlag {
		lastFlag = make([]bool, len(ds.Donations))
		for _, idxs := range ds.DonationsByPerson {
			if len(idxs) > 0 {
				lastFlag[idxs[len(idxs)-1]] = true
			}
		}
	}
	row := func(i int) ([]any, error) {
		v := ds.Donations[i]
		var out []any
		if d.CharityOnDonation {
			out = []any{v.ID, v.PersonID, v.CharityID, v.AmountCents, v.Currency, v.DonatedAt, v.Note}
		} else {
			out = []any{v.ID, v.PersonID, v.AmountCents, v.Currency, v.DonatedAt, v.Note}
		}
		if d.LastFlag {
			out = append(out, lastFlag[i])
		}
		return out, nil
	}
	return parallelCopy(ctx, pool, "donation", dcols, len(ds.Donations), loadConns, row)
}

// jsonDonation is the on-disk element shape for D6. Keys are one character
// because JSONB stores every key of every element: with a million elements,
// verbose key names are megabytes of pure overhead.
type jsonDonation struct {
	I int64   `json:"i"`
	A int64   `json:"a"`
	C string  `json:"c"`
	T string  `json:"t"`
	N *string `json:"n"`
}

// recentCacheJSON builds D9's bounded slice for one person: the newest
// recentCacheSize donations, NEWEST FIRST.
//
// DonationsByPerson is in chronological order (the dataset is generated as an
// append-only stream), so this walks the tail backwards. Getting the direction
// wrong here would not fail the load -- it would quietly make q09 return the
// oldest donations while claiming to return the newest, which is exactly the
// class of bug the verification pass exists to catch.
func recentCacheJSON(ds *Dataset, personIdx int) (string, error) {
	idxs := ds.DonationsByPerson[personIdx]
	n := len(idxs)
	if n > recentCacheSize {
		n = recentCacheSize
	}
	arr := make([]jsonDonation, 0, n)
	for k := 0; k < n; k++ {
		v := ds.Donations[idxs[len(idxs)-1-k]]
		arr = append(arr, jsonDonation{
			I: v.ID, A: v.AmountCents, C: v.Currency,
			T: v.DonatedAt.UTC().Format(time.RFC3339Nano), N: v.Note,
		})
	}
	b, err := json.Marshal(arr)
	return string(b), err
}

func copyPeopleEmbedded(ctx context.Context, pool *pgxpool.Pool, ds *Dataset, loadConns int) error {
	cols := []string{"person_id", "charity_id", "full_name", "email", "joined_at", "donations"}
	row := func(i int) ([]any, error) {
		p := ds.People[i]
		idxs := ds.DonationsByPerson[i]
		arr := make([]jsonDonation, 0, len(idxs))
		for _, di := range idxs {
			v := ds.Donations[di]
			arr = append(arr, jsonDonation{
				I: v.ID, A: v.AmountCents, C: v.Currency,
				T: v.DonatedAt.UTC().Format(time.RFC3339Nano), N: v.Note,
			})
		}
		b, err := json.Marshal(arr)
		if err != nil {
			return nil, err
		}
		return []any{p.ID, p.CharityID, p.FullName, p.Email, p.JoinedAt, string(b)}, nil
	}
	return parallelCopy(ctx, pool, "person", cols, len(ds.People), loadConns, row)
}

// parallelCopy splits one logical COPY into n concurrent streams. On a single
// PostgreSQL node this barely helps; on a distributed cluster it is the
// difference between minutes and an hour, because each stream drives a different
// set of tablet leaders instead of serialising through one connection.
func parallelCopy(ctx context.Context, pool *pgxpool.Pool, table string, cols []string, total, conns int, row func(int) ([]any, error)) error {
	if conns < 1 {
		conns = 1
	}
	if total == 0 {
		return nil
	}
	chunk := (total + conns - 1) / conns
	var wg sync.WaitGroup
	errs := make([]error, conns)

	for w := 0; w < conns; w++ {
		start := w * chunk
		if start >= total {
			break
		}
		end := start + chunk
		if end > total {
			end = total
		}
		wg.Add(1)
		go func(w, start, end int) {
			defer wg.Done()
			i := start
			_, err := pool.CopyFrom(ctx, pgx.Identifier{table}, cols,
				pgx.CopyFromFunc(func() ([]any, error) {
					if i >= end {
						return nil, nil
					}
					v, err := row(i)
					i++
					return v, err
				}))
			errs[w] = err
		}(w, start, end)
	}
	wg.Wait()
	for _, e := range errs {
		if e != nil {
			return fmt.Errorf("%s: %w", table, e)
		}
	}
	return nil
}

// bulkRollup computes in two set-oriented statements what the trigger would
// otherwise compute one row at a time.
func bulkRollup(ctx context.Context, pool *pgxpool.Pool) error {
	const personSQL = `
UPDATE person p SET
    donation_count      = s.cnt,
    total_donated_cents = s.total,
    first_donation_at   = s.first_at,
    last_donation_at    = s.last_at,
    last_donation_id    = s.last_id
FROM (
    SELECT person_id,
           COUNT(*)                AS cnt,
           SUM(amount_cents)       AS total,
           MIN(donated_at)         AS first_at,
           MAX(donated_at)         AS last_at,
           (ARRAY_AGG(donation_id ORDER BY donated_at DESC, donation_id DESC))[1] AS last_id
      FROM donation
     GROUP BY person_id
) s
WHERE p.person_id = s.person_id`

	const charitySQL = `
UPDATE charity c SET
    donation_count       = s.cnt,
    total_donated_cents  = s.total,
    last_donation_at     = s.last_at,
    last_donation_id     = s.last_id,
    last_donor_person_id = s.last_person
FROM (
    SELECT charity_id,
           COUNT(*)          AS cnt,
           SUM(amount_cents) AS total,
           MAX(donated_at)   AS last_at,
           (ARRAY_AGG(donation_id ORDER BY donated_at DESC, donation_id DESC))[1] AS last_id,
           (ARRAY_AGG(person_id   ORDER BY donated_at DESC, donation_id DESC))[1] AS last_person
      FROM donation
     GROUP BY charity_id
) s
WHERE c.charity_id = s.charity_id`

	if _, err := pool.Exec(ctx, personSQL); err != nil {
		return fmt.Errorf("person rollup: %w", err)
	}
	if _, err := pool.Exec(ctx, charitySQL); err != nil {
		return fmt.Errorf("charity rollup: %w", err)
	}
	return nil
}

// SetupAndLoadRetry retries a load that failed because the SERVER killed the
// connection, and nothing else.
//
// Why this exists: on the three-node YugabyteDB cluster, the drop / create / COPY /
// index-backfill cycle that now runs before every isolated write operation
// occasionally has its session terminated mid-DDL ("terminating connection due to
// administrator command", SQLSTATE 57P01). The first time it happened it threw
// away a whole cell -- reads and inserts already measured -- over an event that
// says nothing about the design under test.
//
// The retry is deliberately narrow. A load is idempotent (it starts with a full
// drop), so repeating it is safe; but a genuine schema error, a failed
// verification or a timeout is a real result and must still fail the cell. Only
// server-initiated terminations and dropped connections are retried, the pool
// replaces the dead connection on its own, and every retry is logged so a cell
// that needed one can be identified afterwards.
func SetupAndLoadRetry(ctx context.Context, pool *pgxpool.Pool, d Design, ds *Dataset, engine string, loadConns int) (*LoadPhases, error) {
	const attempts = 3
	var last error
	for i := 1; i <= attempts; i++ {
		ph, err := SetupAndLoad(ctx, pool, d, ds, engine, loadConns)
		if err == nil {
			if i > 1 {
				fmt.Printf("  load succeeded on attempt %d after: %v\n", i, last)
			}
			return ph, nil
		}
		last = err
		if !isServerTermination(err) {
			return nil, err
		}
		fmt.Printf("  load attempt %d/%d lost its connection (%v); retrying\n", i, attempts, err)
		time.Sleep(time.Duration(i*10) * time.Second)
	}
	return nil, fmt.Errorf("load failed after %d attempts: %w", attempts, last)
}

func isServerTermination(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "57P01", "57P02", "57P03", "08006", "08003":
			return true
		}
	}
	msg := err.Error()
	return strings.Contains(msg, "conn closed") ||
		strings.Contains(msg, "unexpected EOF") ||
		strings.Contains(msg, "connection reset")
}
