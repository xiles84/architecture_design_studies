package main

// Adapted from studies/02-ticket-booking/harness/load.go at 3c0c3aa.

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"adsplatform/core/catalog"
	"adsplatform/core/inspect"
	"adsplatform/ports"

	reservedseating "reservedseating"
)

func readSQL(designID, file string) (string, error) {
	b, err := reservedseating.SQL.ReadFile("sql/" + designID + "/" + file)
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

// LoadPhases records how long setup took. Pre-created seat rows are paid for here.
type LoadPhases struct {
	DropMS    float64          `json:"drop_ms"`
	SchemaMS  float64          `json:"schema_ms"`
	CopyMS    float64          `json:"copy_ms"`
	IndexMS   float64          `json:"index_ms"`
	AnalyzeMS float64          `json:"analyze_ms"`
	TotalMS   float64          `json:"total_ms"`
	Rows      map[string]int64 `json:"rows_loaded"`
	// LoadNow is the database clock every loaded hold's times derive from.
	LoadNow time.Time `json:"load_now"`
}

func msSince(t time.Time) float64 { return float64(time.Since(t).Microseconds()) / 1000 }

var studyTables = []string{"ticket", "seat_claim", "event_section", "event_seat", "hold", "event",
	"venue_seat", "venue_section", "venue", "band"}

// Load brings the database from anything to benchmark-ready for design d:
//
//	drop -> schema -> now() -> bulk COPY -> w_load_* -> indexes -> ANALYZE
//
// Every phase of every experiment starts from a fresh Load.
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

	// Every loaded hold's expiry is an offset from this instant, by the database's
	// own clock: a live hold is live on every engine, an expired one expired.
	if err := db.QueryRow(ctx, "SELECT now()").Scan(&ph.LoadNow); err != nil {
		return nil, fmt.Errorf("load now(): %w", err)
	}

	t = time.Now()
	if err := copyAll(ctx, db, d, ds, streams, ph.LoadNow, ph.Rows); err != nil {
		return nil, fmt.Errorf("copy: %w", err)
	}
	writes, err := mustStmts(d.ID, "writes.sql")
	if err != nil {
		return nil, err
	}
	for _, st := range writes {
		if strings.HasPrefix(st.Name, "w_load_") {
			n, err := db.Exec(ctx, st.SQL)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", st.Name, err)
			}
			ph.Rows[st.Name] = n
		}
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
	for _, tb := range studyTables {
		if _, ok := ph.Rows[tb]; !ok {
			continue
		}
		if _, err := db.Exec(ctx, "ANALYZE "+tb); err != nil {
			return nil, fmt.Errorf("analyze %s: %w", tb, err)
		}
	}
	ph.AnalyzeMS = msSince(t)
	ph.TotalMS = msSince(all)
	return ph, nil
}

// seatIndex maps a flat row number to (event, seat) over every event's venue.
type seatIndex struct {
	events []*Event
	starts []int // row number of each event's first seat
	total  int
}

func newSeatIndex(ds *Dataset) *seatIndex {
	si := &seatIndex{}
	for i := range ds.Events {
		e := &ds.Events[i]
		si.events = append(si.events, e)
		si.starts = append(si.starts, si.total)
		si.total += ds.venue(e.VenueID).Capacity
	}
	return si
}

func (si *seatIndex) at(i int) (*Event, int32) {
	k := sort.Search(len(si.starts), func(j int) bool { return si.starts[j] > i }) - 1
	return si.events[k], int32(i-si.starts[k]) + 1
}

type loadedIndex struct {
	sold  map[int64]map[int32]*Ticket
	holds map[int64]map[int32]*LoadedHold
}

func indexLoaded(ds *Dataset) loadedIndex {
	li := loadedIndex{sold: map[int64]map[int32]*Ticket{}, holds: map[int64]map[int32]*LoadedHold{}}
	for i := range ds.Sold {
		t := &ds.Sold[i]
		if li.sold[t.EventID] == nil {
			li.sold[t.EventID] = map[int32]*Ticket{}
		}
		li.sold[t.EventID][t.SeatID] = t
	}
	for i := range ds.Holds {
		h := &ds.Holds[i]
		if li.holds[h.EventID] == nil {
			li.holds[h.EventID] = map[int32]*LoadedHold{}
		}
		for _, s := range h.Seats {
			li.holds[h.EventID][s] = h
		}
	}
	return li
}

func copyAll(ctx context.Context, db ports.DB, d Design, ds *Dataset, streams int, loadNow time.Time, rows map[string]int64) error {
	cp := func(table string, cols []string, n int, row func(int) ([]any, error)) error {
		c, err := parallelCopy(ctx, db, table, cols, n, streams, row)
		if err != nil {
			return err
		}
		rows[table] += c
		return nil
	}
	if err := cp("band", []string{"band_id", "name", "genre", "country"}, len(ds.Bands), func(i int) ([]any, error) {
		b := ds.Bands[i]
		return []any{b.ID, b.Name, b.Genre, b.Country}, nil
	}); err != nil {
		return err
	}
	if err := cp("venue", []string{"venue_id", "name", "capacity"}, len(ds.Venues), func(i int) ([]any, error) {
		v := ds.Venues[i]
		return []any{v.ID, v.Name, int32(v.Capacity)}, nil
	}); err != nil {
		return err
	}
	type secRef struct {
		v *Venue
		s Section
	}
	var secs []secRef
	type vseat struct {
		v *Venue
		s Seat
	}
	var vseats []vseat
	for _, v := range ds.Venues {
		for _, s := range v.Sections {
			secs = append(secs, secRef{v, s})
		}
		for _, s := range v.Seats {
			vseats = append(vseats, vseat{v, s})
		}
	}
	if err := cp("venue_section", []string{"venue_id", "section_no", "capacity"}, len(secs), func(i int) ([]any, error) {
		return []any{secs[i].v.ID, secs[i].s.No, int32(secs[i].s.Capacity)}, nil
	}); err != nil {
		return err
	}
	if err := cp("venue_seat", []string{"venue_id", "seat_id", "section_no", "row_no", "seat_no", "quality_rank"}, len(vseats),
		func(i int) ([]any, error) {
			s := vseats[i].s
			return []any{vseats[i].v.ID, s.ID, s.Section, s.Row, s.No, s.Rank}, nil
		}); err != nil {
		return err
	}
	if err := cp("event", []string{"event_id", "band_id", "venue_id", "name", "starts_at", "price_cents", "description"}, len(ds.Events),
		func(i int) ([]any, error) {
			e := ds.Events[i]
			return []any{e.ID, e.BandID, e.VenueID, e.Name, e.StartsAt, e.PriceCents, e.Description}, nil
		}); err != nil {
		return err
	}
	if err := cp("ticket", []string{"ticket_id", "event_id", "section_no", "seat_id", "customer_id", "hold_id", "sold_at", "price_cents"}, len(ds.Sold),
		func(i int) ([]any, error) {
			t := ds.Sold[i]
			return []any{t.ID, t.EventID, t.Section, t.SeatID, t.CustomerID, nil, t.SoldAt, t.PriceCents}, nil
		}); err != nil {
		return err
	}

	li := indexLoaded(ds)
	switch d.Layout {
	case SeatRows:
		si := newSeatIndex(ds)
		return cp("event_seat", []string{"event_id", "seat_id", "section_no", "status", "hold_id", "customer_id", "hold_expires_at", "sold_at"},
			si.total, func(i int) ([]any, error) {
				e, seat := si.at(i)
				sec := ds.venue(e.VenueID).seat(seat).Section
				if t := li.sold[e.ID][seat]; t != nil {
					return []any{e.ID, seat, sec, "sold", nil, t.CustomerID, nil, t.SoldAt}, nil
				}
				if h := li.holds[e.ID][seat]; h != nil {
					return []any{e.ID, seat, sec, "held", h.ID, h.CustomerID, h.expiresAt(loadNow), nil}, nil
				}
				return []any{e.ID, seat, sec, "available", nil, nil, nil, nil}, nil
			})
	case CartRows:
		if err := cp("hold", []string{"hold_id", "customer_id", "event_id", "section_no", "status", "expires_at", "created_at"}, len(ds.Holds),
			func(i int) ([]any, error) {
				h := ds.Holds[i]
				return []any{h.ID, h.CustomerID, h.EventID, h.Section, "held", h.expiresAt(loadNow), h.createdAt(loadNow)}, nil
			}); err != nil {
			return err
		}
		si := newSeatIndex(ds)
		return cp("event_seat", []string{"event_id", "seat_id", "section_no", "status", "hold_id", "customer_id", "sold_at"},
			si.total, func(i int) ([]any, error) {
				e, seat := si.at(i)
				sec := ds.venue(e.VenueID).seat(seat).Section
				if t := li.sold[e.ID][seat]; t != nil {
					return []any{e.ID, seat, sec, "sold", nil, t.CustomerID, t.SoldAt}, nil
				}
				if h := li.holds[e.ID][seat]; h != nil {
					return []any{e.ID, seat, sec, "held", h.ID, nil, nil}, nil
				}
				return []any{e.ID, seat, sec, "available", nil, nil, nil}, nil
			})
	case ClaimRows:
		type claim struct {
			h    *LoadedHold
			seat int32
		}
		var claims []claim
		for i := range ds.Holds {
			for _, s := range ds.Holds[i].Seats {
				claims = append(claims, claim{&ds.Holds[i], s})
			}
		}
		// Sold claims are written by the design's w_load_sold_claims.
		return cp("seat_claim", []string{"event_id", "seat_id", "section_no", "hold_id", "customer_id", "expires_at", "claimed_at"}, len(claims),
			func(i int) ([]any, error) {
				c := claims[i]
				return []any{c.h.EventID, c.seat, c.h.Section, c.h.ID, c.h.CustomerID, c.h.expiresAt(loadNow), c.h.createdAt(loadNow)}, nil
			})
	case Documents:
		type docRef struct {
			e *Event
			s Section
		}
		var docs []docRef
		for i := range ds.Events {
			e := &ds.Events[i]
			for _, s := range ds.venue(e.VenueID).Sections {
				docs = append(docs, docRef{e, s})
			}
		}
		return cp("event_section", []string{"event_id", "section_no", "capacity", "version", "claims", "next_expiry"}, len(docs),
			func(i int) ([]any, error) {
				dr := docs[i]
				doc := Claims{}
				for seat := dr.s.FirstSeat; seat < dr.s.FirstSeat+int32(dr.s.Capacity); seat++ {
					if t := li.sold[dr.e.ID][seat]; t != nil {
						doc[seat] = ClaimVal{State: 2, Customer: t.CustomerID}
					} else if h := li.holds[dr.e.ID][seat]; h != nil {
						doc[seat] = ClaimVal{State: 1, Hold: h.ID, Customer: h.CustomerID, ExpiresMS: h.expiresAt(loadNow).UnixMilli()}
					}
				}
				text, next := doc.Encode()
				return []any{dr.e.ID, dr.s.No, int32(dr.s.Capacity), int64(0), text, next}, nil
			})
	}
	return fmt.Errorf("unknown layout %q", d.Layout)
}

// ---------------------------------------------------------------------------
// L2's section document.
// ---------------------------------------------------------------------------

// ClaimVal is one claimed seat: [state, hold_id, customer_id, expires_ms].
type ClaimVal struct {
	State     int64
	Hold      int64
	Customer  int64
	ExpiresMS int64
}

type Claims map[int32]ClaimVal

// Encode renders the document and its earliest held expiry (nil if none).
func (c Claims) Encode() (string, any) {
	m := make(map[string][4]int64, len(c))
	var next int64
	for s, v := range c {
		m[strconv.Itoa(int(s))] = [4]int64{v.State, v.Hold, v.Customer, v.ExpiresMS}
		if v.State == 1 && (next == 0 || v.ExpiresMS < next) {
			next = v.ExpiresMS
		}
	}
	b, _ := json.Marshal(m)
	if next == 0 {
		return string(b), nil
	}
	return string(b), time.UnixMilli(next).UTC()
}

func DecodeClaims(text string) (Claims, error) {
	var m map[string][4]int64
	if err := json.Unmarshal([]byte(text), &m); err != nil {
		return nil, fmt.Errorf("decode section document: %w", err)
	}
	c := make(Claims, len(m))
	for k, v := range m {
		s, err := strconv.Atoi(k)
		if err != nil {
			return nil, fmt.Errorf("decode section document: seat key %q", k)
		}
		c[int32(s)] = ClaimVal{State: v[0], Hold: v[1], Customer: v[2], ExpiresMS: v[3]}
	}
	return c, nil
}

// parallelCopy splits one logical COPY into concurrent streams.
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

// DBStats records storage: pre-created seat rows are paid in bytes as well as time.
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
