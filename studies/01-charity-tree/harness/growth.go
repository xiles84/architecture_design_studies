package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type GrowthPhase struct {
	Cycle      int               `json:"cycle"`
	Operation  string            `json:"operation"`
	Operations int               `json:"operations"`
	DurationS  float64           `json:"duration_s"`
	OpsPerSec  float64           `json:"ops_per_sec"`
	Verify     *VerifyReport     `json:"verify"`
	Audit      *RollupAudit      `json:"audit,omitempty"`
	Stats      *DBStats          `json:"stats"`
	Reads      []QueryResult     `json:"reads"`
	Snapshot   map[string]string `json:"snapshot"`
}

// Equal committed mutation counts age every design by the same logical work.
// Half of each inserted batch survives, so both total size and donor histories
// grow. Go mutates its reference dataset only after the SQL succeeds.
func benchmarkGrowth(ctx context.Context, pool *pgxpool.Pool, run *Run, d Design, ds *Dataset, opts BenchOpts) ([]GrowthPhase, error) {
	o := opts.Experiment
	if o.Cycles < 1 || o.Batch < 2 || d.AppRollup {
		return nil, fmt.Errorf("growth requires cycles>=1, batch>=2 and complete mutation maintenance")
	}
	src, err := readSQL(d.ID, "writes.sql")
	if err != nil {
		return nil, err
	}
	ss, err := ParseCatalog(src)
	if err != nil {
		return nil, err
	}
	w := StmtMap(ss)
	var out []GrowthPhase
	var nextID int64
	for _, v := range ds.Donations {
		if v.ID > nextID {
			nextID = v.ID
		}
	}
	exec := func(st Stmt, vals map[string]any) error {
		args, err := bind(st.Params, vals)
		if err != nil {
			return err
		}
		tag, err := pool.Exec(ctx, st.SQL, args...)
		if err != nil {
			return err
		}
		if tag.RowsAffected() != 1 {
			return fmt.Errorf("%s affected %d rows", st.Name, tag.RowsAffected())
		}
		return nil
	}
	for cycle := 1; cycle <= o.Cycles; cycle++ {
		batchStart := len(ds.Donations)
		for _, op := range []string{"insert", "update", "delete"} {
			phase := GrowthPhase{Cycle: cycle, Operation: op}
			start := time.Now()
			n := o.Batch
			if op == "delete" {
				n = o.Batch / 2
			}
			for i := 0; i < n; i++ {
				idx := batchStart + i
				switch op {
				case "insert":
					nextID++
					p := ds.People[(i+(cycle-1)*o.Batch)%len(ds.People)]
					v := Donation{ID: nextID, PersonID: p.ID, CharityID: p.CharityID, AmountCents: int64(1000 + i), Currency: "USD", DonatedAt: epochEnd.Add(time.Duration(nextID) * time.Millisecond)}
					if err := exec(w["w_insert_donation"], map[string]any{"donation_id": v.ID, "person_id": v.PersonID, "charity_id": v.CharityID, "amount_cents": v.AmountCents, "currency": v.Currency, "donated_at": v.DonatedAt, "note": nil}); err != nil {
						return out, err
					}
					ds.Donations = append(ds.Donations, v)
				case "update":
					v := &ds.Donations[idx]
					amount := v.AmountCents + 777
					if err := exec(w["w_update_amount"], map[string]any{"donation_id": v.ID, "amount_cents": amount}); err != nil {
						return out, err
					}
					v.AmountCents = amount
				case "delete":
					v := ds.Donations[idx]
					if err := exec(w["w_delete_donation"], map[string]any{"donation_id": v.ID}); err != nil {
						return out, err
					}
				}
				phase.Operations++
			}
			phase.DurationS = time.Since(start).Seconds()
			phase.OpsPerSec = float64(phase.Operations) / phase.DurationS
			if op == "delete" {
				ds.Donations = append(ds.Donations[:batchStart], ds.Donations[batchStart+n:]...)
			}
			ds.DonationsByPerson = make([][]int32, len(ds.People))
			for i, v := range ds.Donations {
				ds.DonationsByPerson[v.PersonID-1] = append(ds.DonationsByPerson[v.PersonID-1], int32(i))
			}
			phase.Verify, err = gate(ctx, pool, d, ds)
			if err != nil {
				return append(out, phase), err
			}
			if d.Rollups || d.SumOnly || d.RecentCache {
				phase.Audit, err = AuditRollups(ctx, pool, d)
				if err != nil {
					return append(out, phase), err
				}
				if auditBad(phase.Audit) > 0 {
					return append(out, phase), fmt.Errorf("growth invariant failed after %s", op)
				}
			}
			phase.Stats, err = CollectStats(ctx, pool, d, run.Engine)
			if err != nil {
				return out, err
			}
			phase.Snapshot = snapshot(ctx, pool, run.Engine)
			phase.Reads, err = BenchmarkReads(ctx, pool, d, NewBinder(ds), opts, []string{"q07_total_donated_global", "q08_total_donated_charity", "q09_person_recent_donations"})
			if err != nil {
				return out, err
			}
			if err := captureTargetPlans(ctx, pool, run, d, fmt.Sprintf("cycle%d-%s", cycle, op), ""); err != nil {
				return out, err
			}
			out = append(out, phase)
			fmt.Printf("  growth cycle=%d %s: %d operations; %.0f/s; bytes=%d\n", cycle, op, phase.Operations, phase.OpsPerSec, phase.Stats.TotalBytes)
		}
	}
	return out, nil
}
