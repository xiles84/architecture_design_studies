// Command probe asks the pinned engines the questions study 03's designs depend
// on, before any design SQL is trusted (HANDOFF §4). "Ask the engine what it is
// doing; do not recall it" -- LESSONS_LEARNED, study 02.
//
// Every answer is printed with the evidence that produced it, so the transcript
// saved under results/devchecks/ is the record.
package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"adsplatform/adapters/pgxdb"
	"adsplatform/core/inspect"
	"adsplatform/ports"
)

var (
	dsn    = flag.String("dsn", "", "connection string")
	engine = flag.String("engine", "postgres", "postgres | yugabyte")
	only   = flag.String("only", "", "run only this step, e.g. P13")
	p13Dur = flag.Duration("p13-duration", 60*time.Second, "P13: how long to run the hold-then-confirm loop")
	p13Wk  = flag.Int("p13-workers", 32, "P13: concurrent hold-then-confirm loops")
	p13Hot = flag.Int("p13-hot-blocks", 0, "P13: choose among only this many blocks (0 = any), so holds collide like buyers wanting the best seats")
	p13Idx = flag.Bool("p13-study-schema", false, "P13: add S1's partial expiry index and CHECK constraints")
)

func main() {
	flag.Parse()
	ctx := context.Background()
	db, err := pgxdb.Open(ctx, pgxdb.Config{DSN: *dsn, MaxConns: 80, StatementTimeoutMS: 60000, ApplicationName: "seat-probe"})
	if err != nil {
		fatal(err)
	}
	defer db.Close()
	for i := 0; i < 60; i++ {
		if _, err = db.Exec(ctx, "SELECT 1"); err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		fatal(err)
	}
	info := inspect.Describe(ctx, db, *engine)
	fmt.Printf("# Engine probe — %s\n\nversion: %s\n\n", *engine, firstLine(info.Version))

	steps := []struct {
		id, q string
		fn    func(context.Context, ports.DB) string
	}{
		{"P1", "now() inside an explicit transaction is the transaction start time", p1},
		{"P2", "FOR UPDATE NOWAIT accepted; SQLSTATE when the row is locked", p2},
		{"P3", "INSERT ... unnest ... ON CONFLICT DO UPDATE ... WHERE expired RETURNING", p3},
		{"P3b", "the same upsert waiting on an uncommitted valid claim, then re-checking", p3b},
		{"P4", "'infinity'::timestamptz stores, compares and indexes", p4},
		{"P5", "PRIMARY KEY ((a, b) HASH, c ASC) and a lookup plan", p5},
		{"P6", "pgx []int32 bound to = ANY($1::INT[])", p6},
		{"P7", "SERIALIZABLE read-then-update of one row by two sessions", p7},
		{"P8", "clock offset between this client and the database", p8},
		{"P9", "effective isolation for READ COMMITTED requests", func(context.Context, ports.DB) string {
			if info.EffectiveIsolation == "" {
				return "n/a (engine does not report it; PostgreSQL runs READ COMMITTED as requested)"
			}
			return "yb_effective_transaction_isolation_level = " + info.EffectiveIsolation
		}},
		{"P10", "RC conditional UPDATE waiting on a concurrent hold re-checks its WHERE (S1)", p10},
		{"P11", "RC unconditional UPDATE waiting on a concurrent hold overwrites it (S0)", p11},
		{"P12", "RC UPDATE with EXISTS on another table, waiting on a concurrent steal (E2)", p12},
		{"P13", "a hold committed by one transaction is visible to the next transaction's conditional UPDATE, under concurrency", p13},
	}
	for _, s := range steps {
		if *only != "" && s.id != *only {
			continue
		}
		if *only == "" && s.id == "P13" {
			continue // long-running; run with -only P13
		}
		fmt.Printf("## %s — %s\n\n", s.id, s.q)
		out := func() (res string) {
			defer func() {
				if r := recover(); r != nil {
					res = fmt.Sprintf("PANIC: %v", r)
				}
			}()
			return s.fn(ctx, db)
		}()
		fmt.Printf("%s\n\n", strings.TrimSpace(out))
	}
}

func fatal(err error) { fmt.Fprintln(os.Stderr, "fatal:", err); os.Exit(1) }

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}

func code(err error) string {
	if err == nil {
		return "no error"
	}
	// pgx formats server errors as "ERROR: message (SQLSTATE XXXXX)".
	return err.Error()
}

func must(ctx context.Context, q ports.Queryer, sql string, args ...any) {
	if _, err := q.Exec(ctx, sql, args...); err != nil {
		panic(fmt.Sprintf("%s: %v", sql, err))
	}
}

func p1(ctx context.Context, db ports.DB) string {
	tx, err := db.Begin(ctx, ports.ReadCommitted)
	if err != nil {
		return code(err)
	}
	defer tx.Rollback(ctx)
	var a, b, ca, cb time.Time
	_ = tx.QueryRow(ctx, "SELECT now(), clock_timestamp()").Scan(&a, &ca)
	must(ctx, tx, "SELECT pg_sleep(0.3)")
	_ = tx.QueryRow(ctx, "SELECT now(), clock_timestamp()").Scan(&b, &cb)
	return fmt.Sprintf("now() first=%s second=%s equal=%v; clock_timestamp advanced %s",
		a.Format(time.RFC3339Nano), b.Format(time.RFC3339Nano), a.Equal(b), cb.Sub(ca))
}

func p2(ctx context.Context, db ports.DB) string {
	must(ctx, db, "DROP TABLE IF EXISTS probe_nowait")
	must(ctx, db, "CREATE TABLE probe_nowait (id INT PRIMARY KEY, v INT)")
	must(ctx, db, "INSERT INTO probe_nowait VALUES (1, 0)")
	defer db.Exec(ctx, "DROP TABLE IF EXISTS probe_nowait")
	a, _ := db.Begin(ctx, ports.ReadCommitted)
	defer a.Rollback(ctx)
	if _, err := a.Exec(ctx, "SELECT * FROM probe_nowait WHERE id = 1 FOR UPDATE"); err != nil {
		return "session A FOR UPDATE: " + code(err)
	}
	b, _ := db.Begin(ctx, ports.ReadCommitted)
	defer b.Rollback(ctx)
	t := time.Now()
	_, err := b.Exec(ctx, "SELECT * FROM probe_nowait WHERE id = 1 FOR UPDATE NOWAIT")
	return fmt.Sprintf("session B NOWAIT on the locked row after %s: %s (platform class %d)", time.Since(t).Round(time.Millisecond), code(err), db.Classify(err))
}

func p3(ctx context.Context, db ports.DB) string {
	must(ctx, db, "DROP TABLE IF EXISTS probe_claim")
	must(ctx, db, `CREATE TABLE probe_claim (event_id BIGINT, seat_id INT, hold_id BIGINT, expires_at TIMESTAMPTZ NOT NULL,
		PRIMARY KEY (event_id, seat_id))`)
	defer db.Exec(ctx, "DROP TABLE IF EXISTS probe_claim")
	must(ctx, db, "INSERT INTO probe_claim VALUES (1, 1, 10, now() - interval '1 minute'), (1, 2, 11, now() + interval '1 hour')")
	rows, err := db.Query(ctx, `INSERT INTO probe_claim AS c (event_id, seat_id, hold_id, expires_at)
		SELECT $1::BIGINT, s, $2::BIGINT, now() + interval '40 minutes' FROM unnest($3::INT[]) AS s
		ON CONFLICT (event_id, seat_id) DO UPDATE SET hold_id = EXCLUDED.hold_id, expires_at = EXCLUDED.expires_at
		WHERE c.expires_at <= now()
		RETURNING seat_id, now(), expires_at`, int64(1), int64(99), []int32{1, 2, 3})
	if err != nil {
		return code(err)
	}
	var got []int32
	for rows.Next() {
		var s int32
		var n, e time.Time
		_ = rows.Scan(&s, &n, &e)
		got = append(got, s)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return code(err)
	}
	sort.Slice(got, func(i, j int) bool { return got[i] < got[j] })
	return fmt.Sprintf("requested seats [1 expired, 2 valid, 3 new] -> returned %v (expected [1 3])", got)
}

func p3b(ctx context.Context, db ports.DB) string {
	must(ctx, db, "DROP TABLE IF EXISTS probe_claim2")
	must(ctx, db, `CREATE TABLE probe_claim2 (event_id BIGINT, seat_id INT, hold_id BIGINT, expires_at TIMESTAMPTZ NOT NULL,
		PRIMARY KEY (event_id, seat_id))`)
	defer db.Exec(ctx, "DROP TABLE IF EXISTS probe_claim2")
	upsert := `INSERT INTO probe_claim2 AS c (event_id, seat_id, hold_id, expires_at)
		SELECT 1, s, $1::BIGINT, now() + interval '40 minutes' FROM unnest($2::INT[]) AS s
		ON CONFLICT (event_id, seat_id) DO UPDATE SET hold_id = EXCLUDED.hold_id, expires_at = EXCLUDED.expires_at
		WHERE c.expires_at <= now() RETURNING seat_id`
	a, _ := db.Begin(ctx, ports.ReadCommitted)
	defer a.Rollback(ctx)
	if _, err := a.Exec(ctx, upsert, int64(1), []int32{5}); err != nil {
		return "A: " + code(err)
	}
	var wg sync.WaitGroup
	var bRows int
	var bErr error
	var bWait time.Duration
	wg.Add(1)
	go func() {
		defer wg.Done()
		t := time.Now()
		b, err := db.Begin(ctx, ports.ReadCommitted)
		if err != nil {
			bErr = err
			return
		}
		defer b.Rollback(ctx)
		rows, err := b.Query(ctx, upsert, int64(2), []int32{5})
		if err != nil {
			bErr = err
			return
		}
		for rows.Next() {
			bRows++
		}
		rows.Close()
		bErr = rows.Err()
		bWait = time.Since(t)
		if bErr == nil {
			bErr = b.Commit(ctx)
		}
	}()
	time.Sleep(500 * time.Millisecond)
	cerr := a.Commit(ctx)
	wg.Wait()
	var holder int64
	_ = db.QueryRow(ctx, "SELECT hold_id FROM probe_claim2 WHERE seat_id = 5").Scan(&holder)
	return fmt.Sprintf("A commit: %s; B returned %d rows after %s, B error: %s; final holder %d (expected: B 0 rows or a retryable error, holder 1)",
		code(cerr), bRows, bWait.Round(time.Millisecond), code(bErr), holder)
}

func p4(ctx context.Context, db ports.DB) string {
	must(ctx, db, "DROP TABLE IF EXISTS probe_inf")
	must(ctx, db, "CREATE TABLE probe_inf (id INT PRIMARY KEY, expires_at TIMESTAMPTZ NOT NULL)")
	defer db.Exec(ctx, "DROP TABLE IF EXISTS probe_inf")
	if _, err := db.Exec(ctx, "CREATE INDEX probe_inf_exp ON probe_inf (expires_at ASC)"); err != nil {
		return "index: " + code(err)
	}
	must(ctx, db, "INSERT INTO probe_inf VALUES (1, 'infinity'), (2, now() - interval '1 minute'), (3, now() + interval '1 minute')")
	var expired, sold int
	var ok bool
	_ = db.QueryRow(ctx, "SELECT COUNT(*) FILTER (WHERE expires_at <= now()), COUNT(*) FILTER (WHERE expires_at = 'infinity'), bool_and(expires_at > now() OR id = 2) FROM probe_inf").Scan(&expired, &sold, &ok)
	var first int
	_ = db.QueryRow(ctx, "SELECT id FROM probe_inf WHERE expires_at <= now() ORDER BY expires_at LIMIT 1").Scan(&first)
	var back time.Time
	err := db.QueryRow(ctx, "SELECT expires_at FROM probe_inf WHERE id = 1").Scan(&back)
	return fmt.Sprintf("expired=%d (expected 1) infinity=%d (expected 1) others-in-future=%v; indexed expiry scan first id=%d (expected 2); scanning 'infinity' into time.Time: %s",
		expired, sold, ok, first, code(err))
}

func p5(ctx context.Context, db ports.DB) string {
	if *engine != "yugabyte" {
		return "n/a on PostgreSQL (YugabyteDB-only syntax)"
	}
	must(ctx, db, "DROP TABLE IF EXISTS probe_sharded")
	defer db.Exec(ctx, "DROP TABLE IF EXISTS probe_sharded")
	if _, err := db.Exec(ctx, "CREATE TABLE probe_sharded (event_id BIGINT, section_no INT, seat_id INT, v TEXT, PRIMARY KEY ((event_id, section_no) HASH, seat_id ASC))"); err != nil {
		return "create: " + code(err)
	}
	must(ctx, db, "INSERT INTO probe_sharded SELECT 1, s % 4, s, 'x' FROM generate_series(1, 400) AS s")
	rows, err := db.Query(ctx, "EXPLAIN (ANALYZE, DIST) SELECT * FROM probe_sharded WHERE event_id = 1 AND section_no = 2 AND seat_id = ANY (ARRAY[2, 6, 10])")
	if err != nil {
		return "explain: " + code(err)
	}
	var sb strings.Builder
	sb.WriteString("```\n")
	for rows.Next() {
		var l string
		_ = rows.Scan(&l)
		sb.WriteString(l + "\n")
	}
	rows.Close()
	sb.WriteString("```")
	return sb.String()
}

func p6(ctx context.Context, db ports.DB) string {
	var n int
	err := db.QueryRow(ctx, "SELECT COUNT(*) FROM generate_series(1, 10) AS g WHERE g = ANY ($1::INT[])", []int32{2, 4, 11}).Scan(&n)
	var m int
	err2 := db.QueryRow(ctx, "SELECT cardinality($1::INT[])", []int32{2, 4, 11}).Scan(&m)
	return fmt.Sprintf("matched %d (expected 2), error: %s; cardinality %d, error: %s", n, code(err), m, code(err2))
}

func p7(ctx context.Context, db ports.DB) string {
	must(ctx, db, "DROP TABLE IF EXISTS probe_ser")
	must(ctx, db, "CREATE TABLE probe_ser (id INT PRIMARY KEY, status TEXT)")
	must(ctx, db, "INSERT INTO probe_ser VALUES (1, 'available')")
	defer db.Exec(ctx, "DROP TABLE IF EXISTS probe_ser")
	var wg sync.WaitGroup
	results := make([]string, 2)
	start := make(chan struct{})
	readDone := make(chan struct{}, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			tx, err := db.Begin(ctx, ports.Serializable)
			if err != nil {
				results[i] = "begin: " + code(err)
				readDone <- struct{}{}
				return
			}
			defer tx.Rollback(ctx)
			var st string
			err = tx.QueryRow(ctx, "SELECT status FROM probe_ser WHERE id = 1").Scan(&st)
			readDone <- struct{}{}
			if err != nil {
				results[i] = "read: " + code(err)
				return
			}
			<-start
			if _, err := tx.Exec(ctx, "UPDATE probe_ser SET status = $1 WHERE id = 1", fmt.Sprintf("held-by-%d", i)); err != nil {
				results[i] = fmt.Sprintf("read %q, update: %s", st, code(err))
				return
			}
			if err := tx.Commit(ctx); err != nil {
				results[i] = fmt.Sprintf("read %q, commit: %s", st, code(err))
				return
			}
			results[i] = fmt.Sprintf("read %q, committed", st)
		}(i)
	}
	<-readDone
	<-readDone
	close(start)
	wg.Wait()
	var final string
	_ = db.QueryRow(ctx, "SELECT status FROM probe_ser WHERE id = 1").Scan(&final)
	return fmt.Sprintf("session 0: %s\nsession 1: %s\nfinal: %s (expected exactly one commit)", results[0], results[1], final)
}

func p8(ctx context.Context, db ports.DB) string {
	var offs []float64
	for i := 0; i < 20; i++ {
		t0 := time.Now()
		var n time.Time
		if err := db.QueryRow(ctx, "SELECT clock_timestamp()").Scan(&n); err != nil {
			return code(err)
		}
		t1 := time.Now()
		mid := t0.Add(t1.Sub(t0) / 2)
		offs = append(offs, float64(n.Sub(mid).Microseconds())/1000)
	}
	sort.Float64s(offs)
	return fmt.Sprintf("database minus client, ms, over 20 samples: median %.3f, min %.3f, max %.3f", offs[10], offs[0], offs[19])
}

// seatRace runs a hold by session A (uncommitted), then session B's statement on
// the same seat, commits A, and reports what B did.
func seatRace(ctx context.Context, db ports.DB, bSQL string) string {
	must(ctx, db, "DROP TABLE IF EXISTS probe_seat")
	must(ctx, db, `CREATE TABLE probe_seat (event_id BIGINT, seat_id INT, status TEXT NOT NULL, hold_id BIGINT,
		hold_expires_at TIMESTAMPTZ, PRIMARY KEY (event_id, seat_id))`)
	must(ctx, db, "INSERT INTO probe_seat SELECT 1, s, 'available', NULL, NULL FROM generate_series(1, 4) AS s")
	defer db.Exec(ctx, "DROP TABLE IF EXISTS probe_seat")
	a, _ := db.Begin(ctx, ports.ReadCommitted)
	defer a.Rollback(ctx)
	n, err := a.Exec(ctx, `UPDATE probe_seat SET status = 'held', hold_id = 1, hold_expires_at = now() + interval '40 minutes'
		WHERE event_id = 1 AND seat_id = ANY ($1::INT[])
		  AND (status = 'available' OR (status = 'held' AND hold_expires_at <= now()))`, []int32{2, 3})
	if err != nil || n != 2 {
		return fmt.Sprintf("A hold: %d rows, %s", n, code(err))
	}
	var wg sync.WaitGroup
	var bRows int64
	var bErr error
	var bWait time.Duration
	wg.Add(1)
	go func() {
		defer wg.Done()
		t := time.Now()
		b, err := db.Begin(ctx, ports.ReadCommitted)
		if err != nil {
			bErr = err
			return
		}
		defer b.Rollback(ctx)
		bRows, bErr = b.Exec(ctx, bSQL, []int32{3, 4})
		bWait = time.Since(t)
		if bErr == nil {
			bErr = b.Commit(ctx)
		}
	}()
	time.Sleep(500 * time.Millisecond)
	cerr := a.Commit(ctx)
	wg.Wait()
	rows, _ := db.Query(ctx, "SELECT seat_id, status, COALESCE(hold_id, 0) FROM probe_seat ORDER BY seat_id")
	var sb strings.Builder
	for rows.Next() {
		var s int
		var st string
		var h int64
		_ = rows.Scan(&s, &st, &h)
		fmt.Fprintf(&sb, " %d:%s/%d", s, st, h)
	}
	rows.Close()
	return fmt.Sprintf("A holds seats 2,3; B targets 3,4 while A is uncommitted. A commit: %s. B affected %d rows after %s, error: %s. Final:%s",
		code(cerr), bRows, bWait.Round(time.Millisecond), code(bErr), sb.String())
}

func p10(ctx context.Context, db ports.DB) string {
	return seatRace(ctx, db, `UPDATE probe_seat SET status = 'held', hold_id = 2, hold_expires_at = now() + interval '40 minutes'
		WHERE event_id = 1 AND seat_id = ANY ($1::INT[])
		  AND (status = 'available' OR (status = 'held' AND hold_expires_at <= now()))`) +
		"\n(expected: B affects 1 row -- seat 4 only -- so the harness rolls B back as a conflict; seat 3 stays with hold 1)"
}

func p11(ctx context.Context, db ports.DB) string {
	return seatRace(ctx, db, `UPDATE probe_seat SET status = 'held', hold_id = 2, hold_expires_at = now() + interval '40 minutes'
		WHERE event_id = 1 AND seat_id = ANY ($1::INT[])`) +
		"\n(expected for the S0 control: B affects 2 rows and seat 3 ends with hold 2 -- the first hold silently overwritten)"
}

func p12(ctx context.Context, db ports.DB) string {
	must(ctx, db, "DROP TABLE IF EXISTS probe_es")
	must(ctx, db, "DROP TABLE IF EXISTS probe_cart")
	must(ctx, db, "CREATE TABLE probe_cart (hold_id BIGINT PRIMARY KEY, expires_at TIMESTAMPTZ NOT NULL)")
	must(ctx, db, "CREATE TABLE probe_es (event_id BIGINT, seat_id INT, status TEXT NOT NULL, hold_id BIGINT, PRIMARY KEY (event_id, seat_id))")
	defer db.Exec(ctx, "DROP TABLE IF EXISTS probe_es")
	defer db.Exec(ctx, "DROP TABLE IF EXISTS probe_cart")
	// Seat 1 is held by cart 10, which has expired. Two stealers race for it.
	must(ctx, db, "INSERT INTO probe_cart VALUES (10, now() - interval '1 minute')")
	must(ctx, db, "INSERT INTO probe_es VALUES (1, 1, 'held', 10)")
	steal := `UPDATE probe_es es SET status = 'held', hold_id = $1
		WHERE es.event_id = 1 AND es.seat_id = 1
		  AND (es.status = 'available' OR (es.status = 'held'
		       AND EXISTS (SELECT 1 FROM probe_cart h WHERE h.hold_id = es.hold_id AND h.expires_at <= now())))`
	a, _ := db.Begin(ctx, ports.ReadCommitted)
	defer a.Rollback(ctx)
	must(ctx, a, "INSERT INTO probe_cart VALUES (20, now() + interval '40 minutes')")
	na, err := a.Exec(ctx, steal, int64(20))
	if err != nil || na != 1 {
		return fmt.Sprintf("A steal: %d rows, %s", na, code(err))
	}
	var wg sync.WaitGroup
	var nb int64
	var bErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		b, err := db.Begin(ctx, ports.ReadCommitted)
		if err != nil {
			bErr = err
			return
		}
		defer b.Rollback(ctx)
		if _, err := b.Exec(ctx, "INSERT INTO probe_cart VALUES (30, now() + interval '40 minutes')"); err != nil {
			bErr = err
			return
		}
		nb, bErr = b.Exec(ctx, steal, int64(30))
		if bErr == nil {
			bErr = b.Commit(ctx)
		}
	}()
	time.Sleep(500 * time.Millisecond)
	cerr := a.Commit(ctx)
	wg.Wait()
	var holder int64
	_ = db.QueryRow(ctx, "SELECT hold_id FROM probe_es WHERE seat_id = 1").Scan(&holder)
	return fmt.Sprintf("A steals expired seat (cart 20, uncommitted); B tries the same steal (cart 30). A commit: %s. B affected %d rows, error: %s. Final holder %d (expected: B 0 rows or retryable error; holder 20)",
		code(cerr), nb, code(bErr), holder)
}

// p13 is a minimal reproduction, outside the harness, of what study 03's dev checks
// saw on YugabyteDB: a confirmation's conditional UPDATE matching none of the seats
// of a hold that had committed ~100 ms earlier, while a read right afterwards found
// the hold intact. Workers hold a random block with a conditional UPDATE, commit,
// then immediately confirm in a new transaction (SELECT now(); UPDATE ... WHERE
// hold_id = h AND status = 'held' AND hold_expires_at > now()). Readers scan
// per-section availability meanwhile, as buyers do. A confirmation matching fewer
// rows than the hold is re-checked with a plain read and counted as an anomaly
// when the read still finds the whole hold valid.
func p13(ctx context.Context, db ports.DB) string {
	must(ctx, db, "DROP TABLE IF EXISTS probe_p13")
	must(ctx, db, `CREATE TABLE probe_p13 (event_id BIGINT, seat_id INT, section_no INT, status TEXT NOT NULL,
		hold_id BIGINT, hold_expires_at TIMESTAMPTZ, PRIMARY KEY (event_id, seat_id))`)
	defer db.Exec(ctx, "DROP TABLE IF EXISTS probe_p13")
	must(ctx, db, "INSERT INTO probe_p13 SELECT 1, s, (s - 1) / 500 + 1, 'available', NULL, NULL FROM generate_series(1, 10000) AS s")
	must(ctx, db, "CREATE INDEX probe_p13_section ON probe_p13 (event_id, section_no, seat_id)")
	if *p13Idx {
		must(ctx, db, "CREATE INDEX probe_p13_expiry ON probe_p13 (hold_expires_at ASC) WHERE status = 'held'")
		must(ctx, db, "ALTER TABLE probe_p13 ADD CONSTRAINT p13_held CHECK (status <> 'held' OR (hold_id IS NOT NULL AND hold_expires_at IS NOT NULL))")
	}

	const hold = `UPDATE probe_p13 SET status = 'held', hold_id = $1, hold_expires_at = now() + interval '40 minutes'
		WHERE event_id = 1 AND section_no = $2 AND seat_id = ANY ($3::INT[])
		  AND (status = 'available' OR (status = 'held' AND hold_expires_at <= now()))
		RETURNING seat_id`
	const confirm = `UPDATE probe_p13 SET status = 'sold'
		WHERE event_id = 1 AND section_no = $1 AND seat_id = ANY ($3::INT[])
		  AND hold_id = $2 AND status = 'held' AND hold_expires_at > now()
		RETURNING seat_id`
	const check = `SELECT COUNT(*) FROM probe_p13 WHERE event_id = 1 AND section_no = $1 AND hold_id = $2
		AND status = 'held' AND hold_expires_at > now()`
	const reset = `UPDATE probe_p13 SET status = 'available', hold_id = NULL, hold_expires_at = NULL
		WHERE event_id = 1 AND section_no = $1 AND seat_id = ANY ($2::INT[])`

	count := func(rows ports.Rows, err error) (int, error) {
		if err != nil {
			return 0, err
		}
		defer rows.Close()
		n := 0
		for rows.Next() {
			n++
		}
		return n, rows.Err()
	}
	var (
		holds, confirms, anomalies, errs atomic.Int64
		nextHold                         atomic.Int64
		mu                               sync.Mutex
		examples                         []string
		wg                               sync.WaitGroup
	)
	nextHold.Store(1)
	deadline := time.Now().Add(*p13Dur)
	for w := 0; w < *p13Wk; w++ {
		wg.Add(1)
		go func(seed int64) {
			defer wg.Done()
			r := rand.New(rand.NewSource(seed))
			for time.Now().Before(deadline) {
				section := int32(1 + r.Intn(20))
				n := 2 + r.Intn(3)
				rowStart := int32((section-1)*500 + int32(r.Intn(20))*25 + 1)
				first := rowStart + int32(r.Intn(25-n+1))
				if *p13Hot > 0 {
					// The best blocks of section 1, row 1: every worker wants the same few.
					section, n = 1, 2
					first = int32(1 + r.Intn(min(*p13Hot, 24)))
				}
				seats := make([]int32, n)
				for i := range seats {
					seats[i] = first + int32(i)
				}
				h := nextHold.Add(1)
				tx, err := db.Begin(ctx, ports.ReadCommitted)
				if err != nil {
					errs.Add(1)
					continue
				}
				got, err := count(tx.Query(ctx, hold, h, section, seats))
				if err != nil || got != n {
					tx.Rollback(ctx)
					if err != nil {
						errs.Add(1)
					}
					continue
				}
				if err := tx.Commit(ctx); err != nil {
					errs.Add(1)
					continue
				}
				holds.Add(1)
				tx2, err := db.Begin(ctx, ports.ReadCommitted)
				if err != nil {
					errs.Add(1)
					continue
				}
				var now time.Time
				if err := tx2.QueryRow(ctx, "SELECT now()").Scan(&now); err != nil {
					tx2.Rollback(ctx)
					errs.Add(1)
					continue
				}
				matched, err := count(tx2.Query(ctx, confirm, section, h, seats))
				if err != nil {
					tx2.Rollback(ctx)
					errs.Add(1)
					continue
				}
				if matched != n {
					tx2.Rollback(ctx)
					var still int
					_ = db.QueryRow(ctx, check, section, h).Scan(&still)
					if still == n {
						anomalies.Add(1)
						mu.Lock()
						if len(examples) < 5 {
							examples = append(examples, fmt.Sprintf("hold %d section %d seats %v: confirm matched %d of %d; a read right after found %d valid", h, section, seats, matched, n, still))
						}
						mu.Unlock()
					}
				} else if err := tx2.Commit(ctx); err != nil {
					errs.Add(1)
					continue
				} else {
					confirms.Add(1)
				}
				_, _ = db.Exec(ctx, reset, section, seats)
			}
		}(int64(w) * 7919)
	}
	for w := 0; w < 4; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for time.Now().Before(deadline) {
				_, _ = count(db.Query(ctx, `SELECT section_no, COUNT(*) FILTER (WHERE status = 'available') FROM probe_p13 WHERE event_id = 1 GROUP BY section_no`))
			}
		}()
	}
	wg.Wait()
	return fmt.Sprintf("%d workers for %s: %d holds committed, %d confirmed, %d anomalies (confirm matched fewer rows than a hold a read then found intact), %d errors\n%s",
		*p13Wk, *p13Dur, holds.Load(), confirms.Load(), anomalies.Load(), errs.Load(), strings.Join(examples, "\n"))
}
