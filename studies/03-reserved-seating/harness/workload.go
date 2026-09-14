package main

import (
	"context"
	"fmt"
	"maps"
	"math/rand"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"adsplatform/core/catalog"
	"adsplatform/core/inspect"
	"adsplatform/core/measure"
	"adsplatform/ports"
)

type Settings struct {
	Workers  int
	Duration time.Duration
	Warmup   time.Duration
	Trials   int
	Streams  int
	Tiers    []int
	WriteOps []string

	HoldTTL      time.Duration // real time: race and writes
	MaxConflicts int

	RaceBuyers     int
	RaceTimeout    time.Duration
	RaceTrials     int
	RaceTierBudget time.Duration
	RaceTiers      []int
	RaceDeferPct   int

	HumanMinute        time.Duration
	LifecycleBuyers    int
	LifecycleTiers     []int
	LifecycleEvents    int
	LifecycleTrials    int
	LcTTL              time.Duration // human time
	LcAbandonPct       int
	LcAbandonExplicit  int
	LcPaymentWindow    time.Duration // human time
	LcGuard            time.Duration // human time
	LcSweepEvery       time.Duration // human time
	LcOutageFrom       time.Duration // human time
	LcOutageTo         time.Duration // human time
	LcE0Skew           time.Duration // human time
	LcTimeout          time.Duration // human time
	LcSweepBatch       int
	LcProbeCapPerEvent int
}

// human converts a human-time duration into wall-clock time for the lifecycle.
func (s Settings) human(d time.Duration) time.Duration {
	return time.Duration(float64(d) * float64(s.HumanMinute) / float64(time.Minute))
}

// ---------------------------------------------------------------------------
// Reads: every question in isolation; seat map and availability per size tier.
// ---------------------------------------------------------------------------

type ReadResult struct {
	measure.Result
	Query string `json:"query"`
	Tier  int    `json:"tier,omitempty"`
}

func BenchmarkReads(ctx context.Context, db ports.DB, d Design, ds *Dataset, loadNow time.Time, s Settings) ([]ReadResult, error) {
	stmts, err := mustStmts(d.ID, "queries.sql")
	if err != nil {
		return nil, err
	}
	q := catalog.Map(stmts)
	budget := measure.Budget{Workers: s.Workers, Warmup: s.Warmup, Duration: s.Duration, Trials: s.Trials}
	run := func(name string, tier int, key func(r *rand.Rand) []any) ReadResult {
		st := q[name]
		label := name
		if tier > 0 {
			label = fmt.Sprintf("%s@%d", name, tier)
		}
		res := measure.Run(ctx, budget, label, func(ctx context.Context, r *rand.Rand) (measure.Outcome, error) {
			args := key(r)
			if len(st.Params) > len(args) {
				args = append(args, time.Now()) // E0: app_now, node 0
			}
			rows, err := db.Query(ctx, st.SQL, args...)
			if err != nil {
				return measure.Outcome{}, err
			}
			defer rows.Close()
			for rows.Next() {
				// Drain: the server has not necessarily done the work until the rows are fetched.
			}
			return measure.Outcome{}, rows.Err()
		})
		fmt.Printf("    %-34s %10.1f ops/s  p50=%7.3fms p99=%8.3fms  errors=%d%s\n",
			label, res.OpsPerSec, res.Latency.P50MS, res.Latency.P99MS, res.Errors, measure.SpreadNote(res))
		return ReadResult{Result: res, Query: name, Tier: tier}
	}

	var out []ReadResult
	for _, tier := range ds.tiersOf("catalogue") {
		evs := ds.eventsOf("catalogue", tier)
		v := ds.venue(evs[0].VenueID)
		out = append(out, run("q01_section_map", tier, func(r *rand.Rand) []any {
			return []any{evs[r.Intn(len(evs))].ID, int32(1 + r.Intn(len(v.Sections)))}
		}))
		out = append(out, run("q02_event_sections", tier, func(r *rand.Rand) []any { return []any{evs[r.Intn(len(evs))].ID} }))
		out = append(out, run("q03_event_available", tier, func(r *rand.Rand) []any { return []any{evs[r.Intn(len(evs))].ID} }))
	}
	var live []LoadedHold
	for _, h := range ds.Holds {
		if h.live() {
			live = append(live, h)
		}
	}
	if len(live) > 0 {
		out = append(out, run("q04_hold_seats", 0, func(r *rand.Rand) []any {
			h := live[r.Intn(len(live))]
			return []any{h.EventID, h.Section, h.ID}
		}))
	}
	customers := map[int64]bool{}
	for _, t := range ds.Sold {
		customers[t.CustomerID] = true
	}
	custList := make([]int64, 0, len(customers))
	for c := range customers {
		custList = append(custList, c)
	}
	sort.Slice(custList, func(i, j int) bool { return custList[i] < custList[j] })
	out = append(out, run("q05_customer_tickets", 0, func(r *rand.Rand) []any { return []any{custList[r.Intn(len(custList))]} }))
	out = append(out, run("q06_ticket_by_id", 0, func(r *rand.Rand) []any { return []any{ds.Sold[r.Intn(len(ds.Sold))].ID} }))
	return out, nil
}

// ---------------------------------------------------------------------------
// Isolated writes, each on a fresh load.
// ---------------------------------------------------------------------------

type WriteResult struct {
	measure.Result
	Op                 string  `json:"op"`
	Tier               int     `json:"tier,omitempty"`
	Audit              *Audit  `json:"audit,omitempty"`
	AttemptsPerSuccess float64 `json:"attempts_per_success,omitempty"`
}

// BenchmarkHoldSpread is demand spread across the catalogue, weighted by the seats
// each event has left: choose a block, hold it, confirm at once.
func BenchmarkHoldSpread(ctx context.Context, sl *Seller, s Settings) WriteResult {
	ds := sl.world.ds
	var evs []*Event
	var cum []int64
	var total int64
	for _, e := range ds.eventsOf("catalogue", 0) {
		left := int64(truthAvailable(ds.State[e.ID]))
		if left <= 0 {
			continue
		}
		total += left
		evs = append(evs, e)
		cum = append(cum, total)
	}
	b := measure.Budget{Workers: s.Workers, Warmup: s.Warmup, Duration: s.Duration, Trials: s.Trials,
		OpsPerTrial: measure.FiniteBudget(int64(float64(total)/2.5), s.Trials)}
	res := measure.Run(ctx, b, "hold", func(ctx context.Context, r *rand.Rand) (measure.Outcome, error) {
		x := r.Int63n(total)
		i := sort.Search(len(cum), func(i int) bool { return cum[i] > x })
		buyer := &Buyer{seller: sl, node: r.Intn(2), r: r, maxConflicts: s.MaxConflicts}
		cust := sl.world.randomCustomer(r)
		a, err := buyer.Acquire(ctx, evs[i], cust, partySize(r), s.HoldTTL)
		if err != nil {
			return measure.Outcome{Retries: a.Conflicts + a.EngineRetries}, err
		}
		if !a.Granted {
			return measure.Outcome{Rejected: true, Retries: a.Conflicts + a.EngineRetries}, nil
		}
		cr, err := sl.Confirm(ctx, buyer.node, a.Block, a.Hold.HoldID, cust)
		if err == nil && !cr.Sold {
			err = fmt.Errorf("a %s hold was refused at its immediate confirmation (%s)", s.HoldTTL, cr.Class)
		}
		return measure.Outcome{Retries: a.Conflicts + a.EngineRetries + cr.Retries}, err
	})
	return finishWrite("hold", 0, res)
}

// BenchmarkRelease releases loaded live holds explicitly, each at most once.
func BenchmarkRelease(ctx context.Context, sl *Seller, s Settings) WriteResult {
	ds := sl.world.ds
	var pool []LoadedHold
	for _, h := range ds.Holds {
		if h.live() {
			pool = append(pool, h)
		}
	}
	rand.New(rand.NewSource(ds.Seed+13)).Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })
	var cursor atomicCursor
	b := measure.Budget{Workers: s.Workers, Warmup: s.Warmup, Duration: s.Duration, Trials: s.Trials,
		OpsPerTrial: measure.FiniteBudget(int64(len(pool)), s.Trials)}
	res := measure.Run(ctx, b, "release", func(ctx context.Context, r *rand.Rand) (measure.Outcome, error) {
		i := cursor.next()
		if i >= int64(len(pool)) {
			return measure.Outcome{}, measure.ErrExhausted
		}
		h := pool[i]
		n, err := sl.Release(ctx, 0, Block{Event: ds.event(h.EventID), Section: h.Section, Seats: h.Seats}, h.ID)
		if err == nil && n != len(h.Seats) {
			err = fmt.Errorf("released %d of %d seats of loaded hold %d", n, len(h.Seats), h.ID)
		}
		return measure.Outcome{}, err
	})
	return finishWrite("release", 0, res)
}

// BenchmarkCancel refunds loaded tickets, each at most once.
func BenchmarkCancel(ctx context.Context, sl *Seller, s Settings) WriteResult {
	ds := sl.world.ds
	pool := make([]int64, 0, len(ds.Sold))
	for _, t := range ds.Sold {
		pool = append(pool, t.ID)
	}
	rand.New(rand.NewSource(ds.Seed+11)).Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })
	var cursor atomicCursor
	b := measure.Budget{Workers: s.Workers, Warmup: s.Warmup, Duration: s.Duration, Trials: s.Trials,
		OpsPerTrial: measure.FiniteBudget(int64(len(pool)), s.Trials)}
	res := measure.Run(ctx, b, "cancel", func(ctx context.Context, r *rand.Rand) (measure.Outcome, error) {
		i := cursor.next()
		if i >= int64(len(pool)) {
			return measure.Outcome{}, measure.ErrExhausted
		}
		ok, err := sl.Cancel(ctx, pool[i])
		if err == nil && !ok {
			err = fmt.Errorf("ticket %d was not sold: the cancel pool handed out a ticket twice", pool[i])
		}
		return measure.Outcome{}, err
	})
	return finishWrite("cancel", 0, res)
}

// BenchmarkPublish creates events at a venue of one size. Count-based.
func BenchmarkPublish(ctx context.Context, sl *Seller, v *Venue, s Settings) WriteResult {
	perTrial := max(int64(2), int64(20000/v.Capacity))
	b := measure.Budget{Workers: min(4, s.Workers), WarmupOps: 1, OpsPerTrial: perTrial, Trials: s.Trials,
		SafetyLimit: 5 * time.Minute}
	bands := int64(len(sl.world.ds.Bands))
	res := measure.Run(ctx, b, fmt.Sprintf("publish@%d", v.Capacity), func(ctx context.Context, r *rand.Rand) (measure.Outcome, error) {
		_, err := sl.Publish(ctx, r.Int63n(bands)+1, v, r)
		return measure.Outcome{}, err
	})
	return finishWrite("publish", v.Capacity, res)
}

func finishWrite(op string, tier int, res measure.Result) WriteResult {
	w := WriteResult{Result: res, Op: op, Tier: tier}
	if res.Ops > 0 {
		w.AttemptsPerSuccess = float64(res.Ops+res.Retries) / float64(res.Ops)
	}
	flag := ""
	if res.Exhausted {
		flag = "  WARNING: pool exhausted"
	}
	if res.SafetyStopped {
		flag += "  WARNING: stopped by safety limit"
	}
	fmt.Printf("    %-34s %10.1f ops/s  p50=%7.3fms p99=%8.3fms  rejected=%d retries=%d errors=%d%s%s\n",
		res.Name, res.OpsPerSec, res.Latency.P50MS, res.Latency.P99MS, res.Rejected, res.Retries, res.Errors,
		measure.SpreadNote(res), flag)
	if res.FirstError != "" {
		fmt.Printf("      first error: %s\n", res.FirstError)
	}
	return w
}

// atomicCursor hands out each index of a finite pool exactly once across workers.
type atomicCursor struct{ n atomic.Int64 }

func (c *atomicCursor) next() int64 { return c.n.Add(1) - 1 }

// ---------------------------------------------------------------------------
// Plans
// ---------------------------------------------------------------------------

// ExplainAll captures plans for every read AND write statement; writes run inside
// transactions that are rolled back. Keys: the largest catalogue event, its
// section with the most sold seats, one of its live holds, a sold ticket.
func ExplainAll(ctx context.Context, db ports.DB, d Design, ds *Dataset, w *World, engine string) (map[string]string, error) {
	var ev *Event
	for _, e := range ds.eventsOf("catalogue", 0) {
		if ev == nil || ds.venue(e.VenueID).Capacity > ds.venue(ev.VenueID).Capacity {
			ev = e
		}
	}
	if ev == nil || len(ds.Sold) == 0 {
		return nil, fmt.Errorf("dataset has no sold catalogue event to explain with")
	}
	v := ds.venue(ev.VenueID)
	st := ds.State[ev.ID]
	sec := sectionMost(v, st, stSold)
	var hold *LoadedHold
	for i := range ds.Holds {
		if h := &ds.Holds[i]; h.EventID == ev.ID && h.live() {
			hold = h
			break
		}
	}
	var ticket *Ticket
	for i := range ds.Sold {
		if ds.Sold[i].EventID == ev.ID {
			ticket = &ds.Sold[i]
			break
		}
	}
	seats := []int32{v.section(sec).FirstSeat, v.section(sec).FirstSeat + 1}
	holdID := int64(1)
	if hold != nil {
		holdID, seats, sec = hold.ID, hold.Seats, hold.Section
	}
	vals := map[string]any{
		"event_id": ev.ID, "section_no": sec, "seat_ids": seats, "hold_id": holdID, "hold_ids": []int64{holdID},
		"customer_id": int64(1), "hold_ms": float64(40 * 60 * 1000), "payment_window_ms": float64(10 * 60 * 1000),
		"ticket_ids": make([]int64, len(seats)), "price_cents": ev.PriceCents, "batch": 100, "app_now": time.Now(),
		"venue_id": ev.VenueID, "band_id": ev.BandID, "name": "explain show", "starts_at": showEpoch,
		"description": "Explained.", "seat_id": seats[0], "claims": "{}", "next_expiry": nil, "version": int64(0),
	}
	if ticket != nil {
		vals["ticket_id"] = ticket.ID
	} else {
		vals["ticket_id"] = int64(1)
	}
	for i := range seats {
		vals["ticket_ids"].([]int64)[i] = w.nextTicket.Load() + 20_000_000 + int64(i)
	}
	out := map[string]string{}
	reads, err := mustStmts(d.ID, "queries.sql")
	if err != nil {
		return nil, err
	}
	for name, plan := range inspect.Explain(ctx, db, engine, reads, vals, false) {
		out[name] = plan
	}
	writes, err := mustStmts(d.ID, "writes.sql")
	if err != nil {
		return nil, err
	}
	var publish, rest []catalog.Stmt
	for _, st := range writes {
		switch {
		case strings.HasPrefix(st.Name, "w_load_"):
			out[st.Name] = "NOT CAPTURED: run once by the loader, not by any flow"
		case strings.HasPrefix(st.Name, "w_publish_"):
			publish = append(publish, st)
		default:
			rest = append(rest, st)
		}
	}
	pv := maps.Clone(vals)
	pv["event_id"] = w.nextEvent.Load() + 5_000_000
	for name, plan := range inspect.ExplainSequence(ctx, db, engine, publish, pv) {
		out[name] = plan
	}
	for name, plan := range inspect.Explain(ctx, db, engine, rest, vals, true) {
		out[name] = plan
	}
	return out, nil
}
