package main

import (
	"context"
	"fmt"
	"time"

	"adsplatform/core/measure"
	"adsplatform/ports"
)

// stalenessProbe measures the rollup lag: after a child change is acknowledged,
// how long the design's parent aggregate takes to reflect it.
//
// Study 04's n3 (trigger-maintained) and n4 (application-maintained) both update
// the aggregate inside the publishing transaction, so the honest expectation is
// a lag of about one read round-trip and no staleness window. The v2 review found
// the harness recorded no such number at all, and "there is no staleness window"
// is a claim that has to be shown rather than assumed. The drift control x2 is
// deliberately not probed here: it maintains the aggregate in a second
// transaction after commit, and its own audit is what measures that.
func (c *cell) stalenessProbe(ctx context.Context, samples int) (measure.LatencyStats, []string, error) {
	var lags []time.Duration
	var timeouts []string
	ips := c.ds.Installations
	if samples <= 0 || samples > len(ips) {
		samples = len(ips)
	}
	for i := 0; i < samples; i++ {
		ip := ips[i]
		before, ok, err := c.rollupCount(ctx, ip.ID)
		if err != nil {
			return measure.LatencyStats{}, nil, err
		}
		if !ok {
			continue
		}
		key := c.uniqueKey("lag")
		val := fmt.Sprintf("lag-%d", ip.ID)
		t0 := time.Now()
		if err := c.publish(ctx, ip.ID, func(q ports.Queryer) error {
			_, e := c.cat.exec(ctx, q, "w02_add_key", map[string]any{
				"installed_product_id": ip.ID, "key": key, "value": val})
			return e
		}); err != nil {
			return measure.LatencyStats{}, nil, err
		}
		// The ledger must carry the probe's key or the audit after this phase will
		// report the database as ahead of it.
		c.led.set(ip.ID, key, val)

		deadline := t0.Add(3 * time.Second)
		for {
			now, ok, err := c.rollupCount(ctx, ip.ID)
			if err != nil {
				return measure.LatencyStats{}, nil, err
			}
			if ok && now >= before+1 {
				lags = append(lags, time.Since(t0))
				break
			}
			if time.Now().After(deadline) {
				timeouts = append(timeouts, fmt.Sprintf(
					"installation %d: aggregate still %d after 3s, expected at least %d", ip.ID, now, before+1))
				break
			}
			time.Sleep(2 * time.Millisecond)
		}
	}
	return measure.Summarize(lags), timeouts, nil
}

// rollupCount reads the design's own aggregate for one installation through r04,
// the read every design must answer. The bool is false when the result carries no
// row for that installation.
func (c *cell) rollupCount(ctx context.Context, ipID int64) (int64, bool, error) {
	a, err := c.cat.args(sR04, map[string]any{})
	if err != nil {
		return 0, false, err
	}
	rows, err := c.db.Query(ctx, c.cat.stmt(sR04).SQL, a...)
	if err != nil {
		return 0, false, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			id              int64
			pdName, envName string
			buName          *string
			count, rev      int64
			lastModified    time.Time
		)
		if err := rows.Scan(&id, &pdName, &envName, &buName, &count, &rev, &lastModified); err != nil {
			return 0, false, err
		}
		if id == ipID {
			return count, true, nil
		}
	}
	return 0, false, rows.Err()
}
