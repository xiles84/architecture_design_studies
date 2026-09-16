package main

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"adsplatform/core/catalog"
	"adsplatform/core/inspect"
	"adsplatform/ports"

	ticketbooking "ticketbooking"
)

func readSQL(designID, file string) (string, error) {
	b, err := ticketbooking.SQL.ReadFile("sql/" + designID + "/" + file)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func mustStmts(designID, file string) ([]catalog.Stmt, error) {
	src, err := readSQL(designID, file)
	if err != nil {
		return nil, err
	}
	return catalog.Parse(src)
}

// LoadPhases records how long setup took. For this study load time is itself a
// result: pre-created designs write one row per seat before anything is sold.
type LoadPhases struct {
	DropMS    float64          `json:"drop_ms"`
	SchemaMS  float64          `json:"schema_ms"`
	CopyMS    float64          `json:"copy_ms"`
	IndexMS   float64          `json:"index_ms"`
	AnalyzeMS float64          `json:"analyze_ms"`
	TotalMS   float64          `json:"total_ms"`
	Rows      map[string]int64 `json:"rows_loaded"`
}

func msSince(t time.Time) float64 { return float64(time.Since(t).Microseconds()) / 1000 }

var studyTables = []string{"sale_event", "reservation", "seat_slot", "event_inventory_bucket", "event_inventory", "ticket", "event", "band"}

// Load brings the database from anything to benchmark-ready for design d:
//
//	drop -> schema -> bulk COPY -> indexes -> ANALYZE
//
// Indexes are built after the data lands, as a real migration would. Every
// phase of every experiment starts from a fresh Load: study 01 measured a 30x
// distortion from letting one write phase inherit another's bloat.
func Load(ctx context.Context, db ports.DB, d Design, ds *Dataset, streams int) (*LoadPhases, error) {
	return inspect.RetryOnTermination(ctx, db, 3, func() (*LoadPhases, error) {
		return load(ctx, db, d, ds, streams)
	})
}

func load(ctx context.Context, db ports.DB, d Design, ds *Dataset, streams int) (*LoadPhases, error) {
	ph := &LoadPhases{Rows: map[string]int64{}}
	all := time.Now()

	t := time.Now()
	for _, tb := range studyTables {
		if _, err := db.Exec(ctx, "DROP TABLE IF EXISTS "+tb+" CASCADE"); err != nil {
			return nil, fmt.Errorf("drop %s: %w", tb, err)
		}
	}
	ph.DropMS = msSince(t)

	t = time.Now()
	schema, err := readSQL(d.ID, "schema.sql")
	if err != nil {
		return nil, err
	}
	if err := catalog.ExecScript(ctx, db, schema); err != nil {
		return nil, fmt.Errorf("schema: %w", err)
	}
	ph.SchemaMS = msSince(t)

	t = time.Now()
	if err := copyAll(ctx, db, d, ds, streams, ph.Rows); err != nil {
		return nil, fmt.Errorf("copy: %w", err)
	}
	ph.CopyMS = msSince(t)

	t = time.Now()
	idx, err := readSQL(d.ID, "indexes.sql")
	if err != nil {
		return nil, err
	}
	if err := catalog.ExecScript(ctx, db, idx); err != nil {
		return nil, fmt.Errorf("indexes: %w", err)
	}
	ph.IndexMS = msSince(t)

	t = time.Now()
	for tb := range ph.Rows {
		if _, err := db.Exec(ctx, "ANALYZE "+tb); err != nil {
			return nil, fmt.Errorf("analyze %s: %w", tb, err)
		}
	}
	ph.AnalyzeMS = msSince(t)
	ph.TotalMS = msSince(all)
	return ph, nil
}

func copyAll(ctx context.Context, db ports.DB, d Design, ds *Dataset, streams int, rows map[string]int64) error {
	n, err := db.CopyFrom(ctx, "band", []string{"band_id", "name", "genre", "country"}, len(ds.Bands),
		func(i int) ([]any, error) {
			b := ds.Bands[i]
			return []any{b.ID, b.Name, b.Genre, b.Country}, nil
		})
	if err != nil {
		return fmt.Errorf("band: %w", err)
	}
	rows["band"] = n

	ecols := []string{"event_id", "band_id", "name", "venue", "starts_at", "capacity", "price_cents", "description"}
	if d.CounterColumn != "" {
		ecols = append(ecols, d.CounterColumn)
	}
	n, err = db.CopyFrom(ctx, "event", ecols, len(ds.Events), func(i int) ([]any, error) {
		e := ds.Events[i]
		row := []any{e.ID, e.BandID, e.Name, e.Venue, e.StartsAt, int32(e.Capacity), e.PriceCents, e.Description}
		if d.CounterColumn != "" {
			row = append(row, int32(e.InitialSold))
		}
		return row, nil
	})
	if err != nil {
		return fmt.Errorf("event: %w", err)
	}
	rows["event"] = n

	if d.Precreated {
		// One row per seat, sold or not. Seats are laid out per event, so a
		// row's index maps straight back to (event, seat).
		type seatRef struct {
			e    *Event
			seat int
		}
		var seats []seatRef
		for i := range ds.Events {
			e := &ds.Events[i]
			for s := 1; s <= e.Capacity; s++ {
				seats = append(seats, seatRef{e, s})
			}
		}
		sold := soldIndex(ds)
		n, err := parallelCopy(ctx, db, "ticket",
			[]string{"ticket_id", "event_id", "seat_no", "status", "customer_id", "sold_at", "price_cents"},
			len(seats), streams, func(i int) ([]any, error) {
				sr := seats[i]
				id := sr.e.FirstTicketID + int64(sr.seat) - 1
				if t, ok := sold[id]; ok {
					return []any{id, sr.e.ID, int32(sr.seat), "sold", t.CustomerID, t.SoldAt, sr.e.PriceCents}, nil
				}
				return []any{id, sr.e.ID, int32(sr.seat), "available", nil, nil, sr.e.PriceCents}, nil
			})
		if err != nil {
			return err
		}
		rows["ticket"] = n
	} else {
		cols := []string{"ticket_id", "event_id", "seat_no", "customer_id", "sold_at", "price_cents"}
		if d.TicketCapacity {
			cols = append(cols, "event_capacity")
		}
		capOf := map[int64]int{}
		for _, e := range ds.Events {
			capOf[e.ID] = e.Capacity
		}
		n, err := parallelCopy(ctx, db, "ticket", cols, len(ds.Sold), streams, func(i int) ([]any, error) {
			t := ds.Sold[i]
			row := []any{t.ID, t.EventID, int32(t.SeatNo), t.CustomerID, t.SoldAt, t.PriceCents}
			if d.TicketCapacity {
				row = append(row, int32(capOf[t.EventID]))
			}
			return row, nil
		})
		if err != nil {
			return err
		}
		rows["ticket"] = n
	}

	if d.InventoryRow {
		n, err := db.CopyFrom(ctx, "event_inventory", []string{"event_id", "remaining"}, len(ds.Events),
			func(i int) ([]any, error) {
				e := ds.Events[i]
				return []any{e.ID, int32(e.Capacity - e.InitialSold)}, nil
			})
		if err != nil {
			return fmt.Errorf("event_inventory: %w", err)
		}
		rows["event_inventory"] = n
	}

	if d.InventoryShard {
		type bucket struct {
			event            int64
			no, cap, remains int
		}
		var buckets []bucket
		for _, e := range ds.Events {
			k := min(8, e.Capacity)
			soldLeft := e.InitialSold
			for bno := 1; bno <= k; bno++ {
				c := e.Capacity / k
				if bno <= e.Capacity%k {
					c++
				}
				// Loaded sales drain buckets in order; which bucket a historical
				// sale came from is not observable, only the totals are.
				take := min(c, soldLeft)
				soldLeft -= take
				buckets = append(buckets, bucket{e.ID, bno, c, c - take})
			}
		}
		n, err := db.CopyFrom(ctx, "event_inventory_bucket",
			[]string{"event_id", "bucket_no", "bucket_capacity", "remaining"}, len(buckets),
			func(i int) ([]any, error) {
				b := buckets[i]
				return []any{b.event, int32(b.no), int32(b.cap), int32(b.remains)}, nil
			})
		if err != nil {
			return fmt.Errorf("event_inventory_bucket: %w", err)
		}
		rows["event_inventory_bucket"] = n
	}

	if d.SeatPool {
		type slot struct {
			event int64
			seat  int
		}
		var slots []slot
		for _, e := range ds.Events {
			for s := e.InitialSold + 1; s <= e.Capacity; s++ {
				slots = append(slots, slot{e.ID, s})
			}
		}
		n, err := parallelCopy(ctx, db, "seat_slot", []string{"event_id", "seat_no"}, len(slots), streams,
			func(i int) ([]any, error) { return []any{slots[i].event, int32(slots[i].seat)}, nil })
		if err != nil {
			return err
		}
		rows["seat_slot"] = n
	}

	if d.Holds {
		rows["reservation"] = 0
	}

	if d.Ledger {
		// Every ticket present at load was, logically, sold once. The ledger
		// must agree from the first row, or the reconciliation audit (r01/r03
		// would be short a history no sale ever recorded, not just from the
		// benchmark's own writes) would fail on a design that has not written
		// anything wrong yet.
		n, err := parallelCopy(ctx, db, "sale_event",
			[]string{"event_id", "seat_no", "customer_id", "kind", "at", "price_cents"},
			len(ds.Sold), streams, func(i int) ([]any, error) {
				t := ds.Sold[i]
				return []any{t.EventID, int32(t.SeatNo), t.CustomerID, "sold", t.SoldAt, t.PriceCents}, nil
			})
		if err != nil {
			return fmt.Errorf("sale_event: %w", err)
		}
		rows["sale_event"] = n
	}
	return nil
}

func soldIndex(ds *Dataset) map[int64]*Ticket {
	m := make(map[int64]*Ticket, len(ds.Sold))
	for i := range ds.Sold {
		m[ds.Sold[i].ID] = &ds.Sold[i]
	}
	return m
}

// parallelCopy splits one logical COPY into concurrent streams. On one
// PostgreSQL node it barely helps; on the 3-node cluster it is the difference
// between minutes and much longer, because each stream drives different tablet
// leaders.
func parallelCopy(ctx context.Context, db ports.DB, table string, cols []string, total, streams int, row func(int) ([]any, error)) (int64, error) {
	streams = max(streams, 1)
	if total == 0 {
		return 0, nil
	}
	chunk := (total + streams - 1) / streams
	var wg sync.WaitGroup
	errs := make([]error, streams)
	counts := make([]int64, streams)
	for w := 0; w < streams; w++ {
		start := w * chunk
		if start >= total {
			break
		}
		end := min(start+chunk, total)
		wg.Add(1)
		go func(w, start, end int) {
			defer wg.Done()
			counts[w], errs[w] = db.CopyFrom(ctx, table, cols, end-start, func(i int) ([]any, error) { return row(start + i) })
		}(w, start, end)
	}
	wg.Wait()
	var n int64
	for w := range errs {
		if errs[w] != nil {
			return 0, fmt.Errorf("%s: %w", table, errs[w])
		}
		n += counts[w]
	}
	return n, nil
}

// DBStats records storage. Pre-creation is paid in bytes as well as in time.
type RelSize struct {
	Name       string `json:"name"`
	Kind       string `json:"kind"`
	Table      string `json:"table,omitempty"`
	TotalBytes int64  `json:"total_bytes"`
}

type DBStats struct {
	Relations  []RelSize `json:"relations,omitempty"`
	TotalBytes int64     `json:"total_bytes"`
	Note       string    `json:"note,omitempty"`
}

func CollectStats(ctx context.Context, db ports.DB, engine string) (*DBStats, error) {
	s := &DBStats{}
	if engine == "yugabyte" {
		s.Note = "size functions are not meaningful on YugabyteDB (storage lives in DocDB, not PostgreSQL heap files)"
		return s, nil
	}
	rows, err := db.Query(ctx, `
SELECT c.relname,
       CASE c.relkind WHEN 'r' THEN 'table' ELSE 'index' END,
       COALESCE(t.relname, ''),
       pg_total_relation_size(c.oid)
  FROM pg_class c
  JOIN pg_namespace n ON n.oid = c.relnamespace
  LEFT JOIN pg_index i ON i.indexrelid = c.oid
  LEFT JOIN pg_class t ON t.oid = i.indrelid
 WHERE n.nspname = 'public' AND c.relkind IN ('r', 'i')
 ORDER BY c.relkind DESC, c.relname`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var r RelSize
		if err := rows.Scan(&r.Name, &r.Kind, &r.Table, &r.TotalBytes); err != nil {
			return nil, err
		}
		if r.Kind == "table" {
			s.TotalBytes += r.TotalBytes
		}
		s.Relations = append(s.Relations, r)
	}
	return s, rows.Err()
}

func indexCount(s *DBStats, table string) int {
	if s == nil {
		return -1
	}
	n := 0
	for _, r := range s.Relations {
		if r.Kind == "index" && strings.EqualFold(r.Table, table) {
			n++
		}
	}
	return n
}
