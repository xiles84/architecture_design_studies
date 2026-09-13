package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ArrivalResult struct {
	Offered          int64        `json:"offered"`
	Accepted         int64        `json:"accepted"`
	Completed        int64        `json:"completed"`
	Errors           int64        `json:"errors"`
	Dropped          int64        `json:"dropped"`
	ReadErrors       int64        `json:"read_errors"`
	FirstError       string       `json:"first_error,omitempty"`
	Retries          int64        `json:"retries"`
	WindowS          float64      `json:"window_s"`
	ElapsedS         float64      `json:"elapsed_s"`
	CompletedPerSec  float64      `json:"completed_per_sec_including_drain"`
	ReadsPerSec      float64      `json:"reads_per_sec"`
	Service          LatencyStats `json:"successful_service_latency"`
	Response         LatencyStats `json:"successful_scheduled_response_latency"`
	QueueDelay       LatencyStats `json:"successful_queue_delay"`
	SchedulerLag     LatencyStats `json:"scheduler_lag"`
	Reconciled       bool         `json:"reconciled"`
	WarmupErrors     int64        `json:"warmup_errors"`
	WarmupAudit      *RollupAudit `json:"warmup_audit,omitempty"`
	InitialDonations int64        `json:"initial_donations"`
	FinalDonations   int64        `json:"final_donations"`
	WarmupCompleted  int64        `json:"warmup_completed"`
}

// A bounded open-loop scheduler. Due times come from one fixed-rate timeline,
// not from worker completions; queue overflow is counted as rejected demand.
// Latency starts at the scheduled arrival and therefore includes scheduler and
// queue delay. Draining accepted work is measured separately from offering it.
func scheduleArrivals(ctx context.Context, duration time.Duration, rate float64, workers, capacity int, op func(context.Context, int64) error) *ArrivalResult {
	res := &ArrivalResult{WindowS: duration.Seconds()}
	if duration <= 0 {
		return res
	}
	type job struct {
		id  int64
		due time.Time
	}
	jobs := make(chan job, capacity)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var service, response, delay, lag []time.Duration
	var failures, completed atomic.Int64
	start := time.Now()
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				t0 := time.Now()
				err := op(ctx, j.id)
				t1 := time.Now()
				mu.Lock()
				if err != nil {
					failures.Add(1)
					if res.FirstError == "" {
						res.FirstError = err.Error()
					}
				} else {
					completed.Add(1)
					service = append(service, t1.Sub(t0))
					response = append(response, t1.Sub(j.due))
					delay = append(delay, t0.Sub(j.due))
				}
				mu.Unlock()
			}
		}()
	}
	count := int64(math.Ceil(duration.Seconds() * rate))
	for i := int64(0); i < count; i++ {
		due := start.Add(time.Duration(float64(time.Second) * float64(i) / rate))
		if wait := time.Until(due); wait > 0 {
			timer := time.NewTimer(wait)
			select {
			case <-timer.C:
			case <-ctx.Done():
				timer.Stop()
			}
		}
		if ctx.Err() != nil {
			break
		}
		res.Offered++
		lag = append(lag, time.Since(due))
		select {
		case jobs <- job{i, due}:
			res.Accepted++
		default:
			res.Dropped++
		}
	}
	// The configured observation window also matters at rates below one/sec.
	if wait := time.Until(start.Add(duration)); wait > 0 {
		timer := time.NewTimer(wait)
		select {
		case <-timer.C:
		case <-ctx.Done():
			timer.Stop()
		}
	}
	close(jobs)
	wg.Wait()
	res.ElapsedS = time.Since(start).Seconds()
	res.Completed = completed.Load()
	res.Errors = failures.Load()
	res.CompletedPerSec = float64(res.Completed) / res.ElapsedS
	res.Service = summarize(service)
	res.Response = summarize(response)
	res.QueueDelay = summarize(delay)
	res.SchedulerLag = summarize(lag)
	return res
}

func benchmarkArrival(ctx context.Context, pool *pgxpool.Pool, d Design, ds *Dataset, b *Binder, opts BenchOpts) (*ArrivalResult, error) {
	o := opts.Experiment
	if o.ArrivalRate <= 0 || math.IsNaN(o.ArrivalRate) || math.IsInf(o.ArrivalRate, 0) || o.ArrivalRate > 1000000 || o.WriteWorkers < 1 || o.Queue < 0 || o.ReadWorkers < 0 || opts.Duration <= 0 {
		return nil, fmt.Errorf("invalid arrival controls")
	}
	if d.Embedded {
		return nil, fmt.Errorf("arrival reconciliation currently requires a donation table")
	}
	src, err := readSQL(d.ID, "writes.sql")
	if err != nil {
		return nil, err
	}
	ws, err := ParseCatalog(src)
	if err != nil {
		return nil, err
	}
	var retries atomic.Int64
	insert := b.insertFn(pool, d, StmtMap(ws), opts, &retries)
	qsrc, err := readSQL(d.ID, "queries.sql")
	if err != nil {
		return nil, err
	}
	qs, err := ParseCatalog(qsrc)
	if err != nil {
		return nil, err
	}
	queries := StmtMap(qs)
	rmix, err := readMixFor(opts.ReadMix)
	if err != nil {
		return nil, err
	}
	var baseline int64
	if err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM donation").Scan(&baseline); err != nil {
		return nil, err
	}
	phase := func(duration time.Duration) (*ArrivalResult, error) {
		var readers sync.WaitGroup
		var reads, readErrors atomic.Int64
		deadline := time.Now().Add(duration)
		for i := 0; i < o.ReadWorkers; i++ {
			readers.Add(1)
			go func(seed int64) {
				defer readers.Done()
				r := rand.New(rand.NewSource(seed))
				for time.Now().Before(deadline) && ctx.Err() == nil {
					st := queries[pickWeighted(rmix, r)]
					args, e := bind(st.Params, b.readVals(r))
					if e != nil {
						readErrors.Add(1)
						continue
					}
					rows, e := pool.Query(ctx, st.SQL, args...)
					if e != nil {
						readErrors.Add(1)
						continue
					}
					for rows.Next() {
					}
					rows.Close()
					if rows.Err() != nil {
						readErrors.Add(1)
					} else {
						reads.Add(1)
					}
				}
			}(int64(42 + i*7919))
		}
		r := scheduleArrivals(ctx, duration, o.ArrivalRate, o.WriteWorkers, o.Queue, func(ctx context.Context, id int64) error { return insert(ctx, requestRandom(id)) })
		readers.Wait()
		r.ReadErrors = readErrors.Load()
		if duration > 0 {
			r.ReadsPerSec = float64(reads.Load()) / duration.Seconds()
		}
		return r, nil
	}
	warm, err := phase(opts.Warmup)
	if err != nil {
		return nil, err
	}
	var warmAudit *RollupAudit
	if d.Rollups || d.SumOnly || d.RecentCache {
		warmAudit, err = AuditRollups(ctx, pool, d)
		if err != nil {
			return nil, err
		}
	}
	retries.Store(0)
	fmt.Println("arrival_measurement_started")
	r, err := phase(opts.Duration)
	if err != nil {
		return nil, err
	}
	r.Retries = retries.Load()
	r.WarmupErrors = warm.Errors + warm.ReadErrors
	r.WarmupAudit = warmAudit
	var after int64
	if err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM donation").Scan(&after); err != nil {
		return nil, err
	}
	r.Reconciled = after-baseline == warm.Completed+r.Completed
	r.InitialDonations, r.FinalDonations, r.WarmupCompleted = baseline, after, warm.Completed
	if !r.Reconciled {
		return r, fmt.Errorf("acknowledged inserts do not reconcile: initial=%d final=%d acknowledged=%d", baseline, after, warm.Completed+r.Completed)
	}
	return r, nil
}
