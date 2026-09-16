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
	Workers   int
	Duration  time.Duration
	Warmup    time.Duration
	Trials    int
	Streams   int
	Tiers     []int
	WriteOps  []string
	RaceTiers []int

	RaceBuyers     int
	RaceTimeout    time.Duration
	RaceTrials     int
	RaceTierBudget time.Duration
	ChurnPct       int
	EditorInterval time.Duration

	HoldBuyers        int
	HoldTTL           time.Duration
	HoldThinkMax      time.Duration
	HoldPayMin        time.Duration
	HoldPayMax        time.Duration
	HoldAbandonPct    int
	HoldSweepEvery    time.Duration
	HoldEventsPerTier int
	HoldTiers         []int
}

// ---------------------------------------------------------------------------
// Reads: every question in isolation, availability once per size tier.
// ---------------------------------------------------------------------------

type ReadResult struct {
	measure.Result
	Query string `json:"query"`
	Tier  int    `json:"tier,omitempty"`
}

func BenchmarkReads(ctx context.Context, db ports.DB, d Design, ds *Dataset, s Settings) ([]ReadResult, error) {
	stmts, err := mustStmts(d.ID, "queries.sql")
	if err != nil {
		return nil, err
	}
	q := catalog.Map(stmts)

	// Key pools: fresh keys per execution, never a single repeated key.
	customers := map[int64]bool{}
	for _, t := range ds.Sold {
		customers[t.CustomerID] = true
	}
	custList := make([]int64, 0, len(customers))
	for c := range customers {
		custList = append(custList, c)
	}
	sort.Slice(custList, func(i, j int) bool { return custList[i] < custList[j] })

	budget := measure.Budget{Workers: s.Workers, Warmup: s.Warmup, Duration: s.Duration, Trials: s.Trials}
	run := func(name string, tier int, key func(r *rand.Rand) any) ReadResult {
		st := q[name]
		label := name
		if tier > 0 {
			label = fmt.Sprintf("%s@%d", name, tier)
		}
		res := measure.Run(ctx, budget, label, func(ctx context.Context, r *rand.Rand) (measure.Outcome, error) {
			rows, err := db.Query(ctx, st.SQL, key(r))
			if err != nil {
				return measure.Outcome{}, err
			}
			defer rows.Close()
			for rows.Next() {
				// Drain: the server has not necessarily done the work until the
				// rows are fetched.
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
		out = append(out, run("q01_event_availability", tier, func(r *rand.Rand) any { return evs[r.Intn(len(evs))].ID }))
	}
	bands := int64(len(ds.Bands))
	out = append(out, run("q02_band_events", 0, func(r *rand.Rand) any { return r.Int63n(bands) + 1 }))
	out = append(out, run("q03_customer_tickets", 0, func(r *rand.Rand) any { return custList[r.Intn(len(custList))] }))
	out = append(out, run("q04_ticket_by_id", 0, func(r *rand.Rand) any { return ds.Sold[r.Intn(len(ds.Sold))].ID }))
	out = append(out, run("q05_band_tickets_sold", 0, func(r *rand.Rand) any { return r.Int63n(bands) + 1 }))
	return out, nil
}

// ---------------------------------------------------------------------------
// Isolated writes, each on a fresh load.
// ---------------------------------------------------------------------------

type WriteResult struct {
	measure.Result
	Op    string `json:"op"`
	Tier  int    `json:"tier,omitempty"`
	Audit *Audit `json:"audit,omitempty"`
	// LedgerAudit: X1 only (REPORTS.md section 4 / AM-02.3).
	LedgerAudit *LedgerAudit `json:"ledger_audit,omitempty"`
	// AttemptsPerSuccess = (successes + retries) / successes: how much work one
	// sale took. 1.0 means no race was ever lost.
	AttemptsPerSuccess float64 `json:"attempts_per_success,omitempty"`
}

// BenchmarkBookSpread is demand spread across the catalogue, weighted by the
// seats each event has left -- the everyday case with little contention on any
// one event. A finite pool (seats), so it is bounded by count as well as time.
func BenchmarkBookSpread(ctx context.Context, bk *Booker, world *World, s Settings) WriteResult {
	var evs []*EventRef
	var cum []int64
	var total int64
	for _, e := range world.snapshot() {
		if e.Kind != "catalogue" {
			continue
		}
		left := int64(e.Capacity) - e.InitialSold
		if left <= 0 {
			continue
		}
		total += left
		evs = append(evs, e)
		cum = append(cum, total)
	}
	b := measure.Budget{Workers: s.Workers, Warmup: s.Warmup, Duration: s.Duration, Trials: s.Trials,
		OpsPerTrial: measure.FiniteBudget(total, s.Trials)}
	res := measure.Run(ctx, b, "book", func(ctx context.Context, r *rand.Rand) (measure.Outcome, error) {
		x := r.Int63n(total)
		i := sort.Search(len(cum), func(i int) bool { return cum[i] > x })
		br, err := bk.Book(ctx, evs[i], world.randomCustomer(r), r)
		return measure.Outcome{Rejected: br.SoldOut, Retries: br.Retries}, err
	})
	return finishWrite("book", 0, res)
}

// BenchmarkCancel refunds loaded tickets, each at most once.
func BenchmarkCancel(ctx context.Context, bk *Booker, ds *Dataset, s Settings) WriteResult {
	var pool []int64
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
		ok, retries, err := bk.Cancel(ctx, pool[i])
		if err == nil && !ok {
			err = fmt.Errorf("ticket %d was not sold: the cancel pool handed out a ticket twice", pool[i])
		}
		return measure.Outcome{Retries: retries}, err
	})
	return finishWrite("cancel", 0, res)
}

// BenchmarkPublish creates events of one size. Count-based: publishing a
// 100 000-seat event is too expensive to run for a fixed duration.
func BenchmarkPublish(ctx context.Context, bk *Booker, world *World, tier int, s Settings) WriteResult {
	perTrial := max(int64(2), int64(20000/tier))
	b := measure.Budget{Workers: min(4, s.Workers), WarmupOps: 1, OpsPerTrial: perTrial, Trials: s.Trials,
		SafetyLimit: 5 * time.Minute}
	res := measure.Run(ctx, b, fmt.Sprintf("publish@%d", tier), func(ctx context.Context, r *rand.Rand) (measure.Outcome, error) {
		_, err := bk.Publish(ctx, r.Int63n(world.bands)+1, tier, r)
		return measure.Outcome{}, err
	})
	return finishWrite("publish", tier, res)
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

// ---------------------------------------------------------------------------
// Plans
// ---------------------------------------------------------------------------

// ExplainAll captures plans for every read AND every write statement. For this
// study the write statements are where the designs live, so they are explained
// too -- each inside a transaction that is rolled back.
func ExplainAll(ctx context.Context, db ports.DB, d Design, ds *Dataset, world *World, engine string) (map[string]string, error) {
	k := chooseKeys(ds)
	if k.busyEvent == nil || k.soldTicket == nil {
		return nil, fmt.Errorf("dataset has no sold catalogue event to explain with")
	}
	ev := k.busyEvent
	vals := map[string]any{
		"event_id": ev.ID, "band_id": ev.BandID, "customer_id": k.customer,
		"ticket_id": k.soldTicket.ID, "start_seat": ev.Capacity / 2, "seat_no": ev.InitialSold + 1,
		"bucket_no": 1, "reservation_id": int64(1), "hold_ms": 250, "batch": 100,
		"description": "Explained.", "name": "explain show", "venue": "stadium",
		"starts_at": showEpoch, "capacity": ev.Capacity, "price_cents": ev.PriceCents,
		"first_ticket_id": world.nextTicket.Load() + 10_000_000,
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
	// Publishing a new event id must not collide with an existing one.
	vals["event_id"] = world.nextEvent.Load() + 1_000_000
	var publish, rest []catalog.Stmt
	for _, st := range writes {
		if strings.HasPrefix(st.Name, "w_publish_") {
			publish = append(publish, st)
		} else {
			rest = append(rest, st)
		}
	}
	// Publishing statements depend on each other (tickets reference the event),
	// so they are explained in order inside one rolled-back transaction, with
	// the busiest event's capacity -- the plan that shows what pre-creation costs.
	for name, plan := range inspect.ExplainSequence(ctx, db, engine, publish, vals) {
		out[name] = plan
	}
	vals["event_id"] = ev.ID
	for _, st := range rest {
		sv := vals
		if st.Name == "w_insert_ticket" || st.Name == "w_claim_seat" {
			// A statement that creates a ticket needs an id nobody holds, and no
			// reservation to point at.
			sv = maps.Clone(vals)
			sv["ticket_id"] = world.nextTicket.Load() + 20_000_000
			sv["reservation_id"] = nil
		}
		for name, plan := range inspect.Explain(ctx, db, engine, []catalog.Stmt{st}, sv, true) {
			out[name] = plan
		}
	}
	return out, nil
}

// atomicCursor hands out each index of a finite pool exactly once across workers.
type atomicCursor struct{ n atomic.Int64 }

func (c *atomicCursor) next() int64 { return c.n.Add(1) - 1 }
