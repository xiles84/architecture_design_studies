package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"adsplatform/core/measure"
	"adsplatform/ports"
)

// The ledger is the client's memory of what it was told. Every expectation the
// audit checks is derived from it and from the generated dataset -- never from
// the database -- so "the database agrees with itself" can never pass as
// correctness (methodology 5).
type ledger struct {
	mu sync.Mutex
	// values is tracked only where the final value is order-independent: a phase
	// whose writes race records counts instead, because two racing writers can
	// legitimately leave either value and asserting one of them would fail a
	// correct design.
	values map[int64]map[string]string
	rev    map[int64]int64

	acked     atomic.Int64
	retries   atomic.Int64
	conflicts atomic.Int64
	ambiguous atomic.Int64

	// incAcked counts acknowledged increments of the contention counter. Each
	// increment adds exactly one, so the counter's final value must be the value
	// it started at plus this number. Comparing the stored value against the
	// *last value written* instead would compare a number with itself and report
	// a lost update as a success -- which is what the first run of the lost-update
	// control did, with 583 acknowledged increments and a counter of 39.
	incAcked atomic.Int64
}

func newLedger(ds *Dataset) *ledger {
	l := &ledger{values: map[int64]map[string]string{}, rev: map[int64]int64{}}
	for _, ip := range ds.Installations {
		m := make(map[string]string, len(ds.Entries[ip.ID]))
		for _, e := range ds.Entries[ip.ID] {
			m[e.Key] = e.Value
		}
		l.values[ip.ID] = m
		l.rev[ip.ID] = 0
	}
	return l
}

func (l *ledger) snapshot() (map[int64]map[string]string, map[int64]int64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make(map[int64]map[string]string, len(l.values))
	for ip, m := range l.values {
		c := make(map[string]string, len(m))
		for k, v := range m {
			c[k] = v
		}
		out[ip] = c
	}
	rev := make(map[int64]int64, len(l.rev))
	for ip, r := range l.rev {
		rev[ip] = r
	}
	return out, rev
}

func (l *ledger) set(ip int64, key, value string) {
	l.mu.Lock()
	l.values[ip][key] = value
	l.mu.Unlock()
}

func (l *ledger) del(ip int64, key string) {
	l.mu.Lock()
	delete(l.values[ip], key)
	l.mu.Unlock()
}

func (l *ledger) bumped(ip int64) {
	l.mu.Lock()
	l.rev[ip]++
	l.mu.Unlock()
}

func (l *ledger) count(ip int64) int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.values[ip])
}

// ---------------------------------------------------------------- cell

// cell is one (topology, design) measurement: everything an operation needs.
type cell struct {
	design Design
	cat    *catalogue
	db     ports.DB
	ds     *Dataset
	led    *ledger
	opts   Options

	// opSeq numbers writes whose key names must be unique across workers.
	opSeq atomic.Int64
	// concIP and concKey are the hot spot the contention phase targets.
	concIP  int64
	concKey string
	// driftCount is the stale rollup value the control writes after commit.
	driftMu    sync.Mutex
	driftCount int
}

// jsonSet returns the document with one key set, preserving every other key.
// The document design changes its configuration in the application, not in SQL,
// which is exactly the mechanism under test: the whole document is rewritten on
// every publication.
func jsonSet(doc, key, value string) (string, error) {
	m := map[string]string{}
	if strings.TrimSpace(doc) != "" {
		if err := json.Unmarshal([]byte(doc), &m); err != nil {
			return "", fmt.Errorf("decode document: %w", err)
		}
	}
	m[key] = value
	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// pickRead chooses an installed product for a read: hot installations get most
// of the traffic, because a uniform draw would hide every cache and contention
// effect this study is about.
func (c *cell) pickRead(r *rand.Rand) Installation {
	return c.ds.PickInstallation(r, r.Intn(100) < 80)
}

// ---------------------------------------------------------------- reads

func (c *cell) readOp(name string) measure.Op {
	return func(ctx context.Context, r *rand.Rand) (measure.Outcome, error) {
		ip := c.pickRead(r)
		entries := c.ds.Entries[ip.ID]
		keyIdx := c.ds.PickKey(r, ip.ID)
		key := entries[keyIdx].Key
		vals := map[string]any{
			"installed_product_id": ip.ID,
			"key":                  key,
			"known_revision":       int64(0),
			"scope_kind":           "product_definition",
			"scope_id":             int64(ip.ProductDefID),
		}
		a, err := c.cat.args(name, vals)
		if err != nil {
			return measure.Outcome{}, err
		}
		rows, err := c.db.Query(ctx, c.cat.stmt(name).SQL, a...)
		if err != nil {
			return measure.Outcome{}, err
		}
		// Rows must be drained: a query whose rows are never fetched has not
		// necessarily been executed (methodology 7).
		n, err := drain(rows)
		if err != nil {
			return measure.Outcome{}, err
		}
		if n == 0 && name != sR03 {
			return measure.Outcome{Rejected: true}, nil
		}
		return measure.Outcome{}, nil
	}
}

// ---------------------------------------------------------------- writes

// publish wraps one publication: the change, the revision advance, and (for the
// application-maintained rollups) the rollup, all in one transaction so that
// "acknowledged" and "durable" are the same event (INV-3, INV-12).
func (c *cell) publish(ctx context.Context, ip int64, apply func(q ports.Queryer) error) error {
	tx, err := c.db.Begin(ctx, ports.ReadCommitted)
	if err != nil {
		return err
	}
	if err := apply(tx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	if _, err := c.cat.exec(ctx, tx, sRevBump, map[string]any{"installed_product_id": ip}); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	if c.design.Rollup == RollupApp {
		if _, err := c.cat.exec(ctx, tx, "w_recompute_rollup", map[string]any{"installed_product_id": ip}); err != nil {
			_ = tx.Rollback(ctx)
			return err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	c.led.bumped(ip)
	return nil
}

// publishDrift is the control's publication: the change and the revision commit
// first, and the rollup is written afterwards from a count read *inside* the
// publishing transaction. Two concurrent publications therefore each write a
// count that is already stale, and the stored rollup drifts (INV-9).
func (c *cell) publishDrift(ctx context.Context, ip int64, apply func(q ports.Queryer) error) error {
	tx, err := c.db.Begin(ctx, ports.ReadCommitted)
	if err != nil {
		return err
	}
	countSeen := c.led.count(ip)
	if err := apply(tx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	if _, err := c.cat.exec(ctx, tx, sRevBump, map[string]any{"installed_product_id": ip}); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	c.led.bumped(ip)

	// Second transaction: the rollup, computed from what was read before the
	// commit. This is the whole control.
	tx2, err := c.db.Begin(ctx, ports.ReadCommitted)
	if err != nil {
		return err
	}
	if _, err := tx2.Exec(ctx,
		"UPDATE installed_product SET config_count = $2 WHERE id = $1", ip, countSeen+1); err != nil {
		_ = tx2.Rollback(ctx)
		return err
	}
	if err := tx2.Commit(ctx); err != nil {
		return err
	}
	return nil
}

// uniqueKey produces a key name no other operation will use, so "add" and
// "delete" change the key set in an order-independent way and the audit can
// assert the final set exactly.
func (c *cell) uniqueKey(prefix string) string {
	return fmt.Sprintf("%s.%s_%06d", prefix, prefix, c.opSeq.Add(1))
}

func (c *cell) writeOp(name string) measure.Op {
	return func(ctx context.Context, r *rand.Rand) (measure.Outcome, error) {
		ip := c.ds.Installations[int(r.Int63())%len(c.ds.Installations)]
		entries := c.ds.Entries[ip.ID]
		key := entries[int(r.Int63())%len(entries)].Key

		switch name {
		case "w01_modify_key":
			val := fmt.Sprintf("modified-%d-%d", ip.ID, time.Now().UnixNano())
			err := c.publish(ctx, ip.ID, func(q ports.Queryer) error {
				_, e := c.cat.exec(ctx, q, name, map[string]any{"installed_product_id": ip.ID, "key": key, "value": val})
				return e
			})
			if err == nil {
				c.led.set(ip.ID, key, val)
			}
			return measure.Outcome{}, err

		case "w02_add_key":
			newKey := c.uniqueKey("added")
			val := fmt.Sprintf("added-%d", ip.ID)
			err := c.publish(ctx, ip.ID, func(q ports.Queryer) error {
				_, e := c.cat.exec(ctx, q, name, map[string]any{"installed_product_id": ip.ID, "key": newKey, "value": val})
				return e
			})
			if err == nil {
				c.led.set(ip.ID, newKey, val)
			}
			return measure.Outcome{}, err

		case "w03_delete_key":
			// Delete a key this phase added, so the delete has something to
			// remove without disturbing the dataset's own keys.
			newKey := c.uniqueKey("added")
			val := "transient"
			if err := c.publish(ctx, ip.ID, func(q ports.Queryer) error {
				_, e := c.cat.exec(ctx, q, "w02_add_key", map[string]any{"installed_product_id": ip.ID, "key": newKey, "value": val})
				return e
			}); err != nil {
				return measure.Outcome{}, err
			}
			err := c.publish(ctx, ip.ID, func(q ports.Queryer) error {
				_, e := c.cat.exec(ctx, q, name, map[string]any{"installed_product_id": ip.ID, "key": newKey})
				return e
			})
			return measure.Outcome{}, err

		case "w04_publish_batch":
			keys := []string{}
			vals := []string{}
			for i := 0; i < c.opts.Batch; i++ {
				keys = append(keys, c.uniqueKey("batch"))
				vals = append(vals, fmt.Sprintf("batch-%d", i))
			}
			err := c.publish(ctx, ip.ID, func(q ports.Queryer) error {
				_, e := c.cat.exec(ctx, q, sW04Up, map[string]any{
					"installed_product_id": ip.ID, "keys": keys, "values": vals})
				return e
			})
			if err == nil {
				for i, k := range keys {
					c.led.set(ip.ID, k, vals[i])
				}
			}
			return measure.Outcome{}, err

		case "w05_replace_all":
			keys := []string{}
			vals := []string{}
			for i := 0; i < c.opts.Tier; i++ {
				keys = append(keys, fmt.Sprintf("replaced.k_%06d", i))
				vals = append(vals, fmt.Sprintf("replaced-%d-%d", ip.ID, i))
			}
			err := c.publish(ctx, ip.ID, func(q ports.Queryer) error {
				if _, e := c.cat.exec(ctx, q, sW05Clear, map[string]any{"installed_product_id": ip.ID}); e != nil {
					return e
				}
				_, e := c.cat.exec(ctx, q, sW05Fill, map[string]any{
					"installed_product_id": ip.ID, "keys": keys, "values": vals})
				return e
			})
			if err == nil {
				c.led.mu.Lock()
				m := map[string]string{}
				for i, k := range keys {
					m[k] = vals[i]
				}
				c.led.values[ip.ID] = m
				c.led.mu.Unlock()
			}
			return measure.Outcome{}, err

		case "w06_update_metadata":
			tx, err := c.db.Begin(ctx, ports.ReadCommitted)
			if err != nil {
				return measure.Outcome{}, err
			}
			if _, err := c.cat.exec(ctx, tx, sW06, map[string]any{
				"installed_product_id": ip.ID,
				"display_name":         fmt.Sprintf("renamed-%d", ip.ID)}); err != nil {
				_ = tx.Rollback(ctx)
				return measure.Outcome{}, err
			}
			if err := tx.Commit(ctx); err != nil {
				return measure.Outcome{}, err
			}
			// Metadata is not a publication: the ledger's revision must not move.
			return measure.Outcome{}, nil

		case "wd1_write_doc":
			return c.docPublish(ctx, ip.ID, r)
		}
		return measure.Outcome{}, fmt.Errorf("unknown write op %q", name)
	}
}

// docPublish is the document design's publication: read the document, change it
// in the application, write it back guarded by the document's revision. A zero-row
// write means someone else published first, and the retry is counted (INV-13).
func (c *cell) docPublish(ctx context.Context, ip int64, r *rand.Rand) (measure.Outcome, error) {
	for attempt := 0; attempt <= c.opts.Retries; attempt++ {
		tx, err := c.db.Begin(ctx, ports.ReadCommitted)
		if err != nil {
			return measure.Outcome{}, err
		}
		a, _ := c.cat.args("wd1_read_doc", map[string]any{"installed_product_id": ip})
		var doc string
		var rev int64
		if err := tx.QueryRow(ctx, c.cat.stmt("wd1_read_doc").SQL, a...).Scan(&doc, &rev); err != nil {
			_ = tx.Rollback(ctx)
			return measure.Outcome{}, err
		}
		newKey := fmt.Sprintf("added.doc_%06d", c.opSeq.Add(1))
		newVal := fmt.Sprintf("doc-%d", ip)
		updated, err := jsonSet(doc, newKey, newVal)
		if err != nil {
			_ = tx.Rollback(ctx)
			return measure.Outcome{}, err
		}
		n, err := c.cat.exec(ctx, tx, "wd1_write_doc", map[string]any{
			"installed_product_id": ip, "doc": updated,
			"bytes": len(updated), "expected_revision": rev})
		if err != nil {
			_ = tx.Rollback(ctx)
			return measure.Outcome{}, err
		}
		if n == 0 {
			_ = tx.Rollback(ctx)
			c.led.conflicts.Add(1)
			continue
		}
		if _, err := c.cat.exec(ctx, tx, sRevBump, map[string]any{"installed_product_id": ip}); err != nil {
			_ = tx.Rollback(ctx)
			return measure.Outcome{}, err
		}
		if err := tx.Commit(ctx); err != nil {
			return measure.Outcome{}, err
		}
		c.led.set(ip, newKey, newVal)
		c.led.bumped(ip)
		return measure.Outcome{Retries: attempt}, nil
	}
	return measure.Outcome{}, fmt.Errorf("document publication gave up after %d attempts", c.opts.Retries)
}

// ---------------------------------------------------------------- contention

// incrementOp is the concurrency family's operation: a read-modify-write whose
// new value is derived from the value it read. That is what makes a lost update
// observable at all -- a design writing an explicit value can only be
// overwritten, which is correct behaviour rather than a lost update.
func (c *cell) incrementOp() measure.Op {
	return func(ctx context.Context, r *rand.Rand) (measure.Outcome, error) {
		ip, key := c.concIP, c.concKey
		switch c.design.Conc {
		case ConcOptimistic:
			for attempt := 0; attempt <= c.opts.Retries; attempt++ {
				tx, err := c.db.Begin(ctx, ports.ReadCommitted)
				if err != nil {
					return measure.Outcome{}, err
				}
				a, _ := c.cat.args("wc1_read_version", map[string]any{"installed_product_id": ip, "key": key})
				var value string
				var version int64
				if err := tx.QueryRow(ctx, c.cat.stmt("wc1_read_version").SQL, a...).Scan(&value, &version); err != nil {
					_ = tx.Rollback(ctx)
					return measure.Outcome{}, err
				}
				next := strconv.Itoa(atoi(value) + 1)
				n, err := c.cat.exec(ctx, tx, "wc1_update_guarded", map[string]any{
					"installed_product_id": ip, "key": key, "value": next, "expected_version": version})
				if err != nil {
					_ = tx.Rollback(ctx)
					return measure.Outcome{}, err
				}
				if n == 0 {
					// Someone published first. That is the signal, not a failure.
					_ = tx.Rollback(ctx)
					c.led.conflicts.Add(1)
					continue
				}
				if _, err := c.cat.exec(ctx, tx, sRevBump, map[string]any{"installed_product_id": ip}); err != nil {
					_ = tx.Rollback(ctx)
					return measure.Outcome{}, err
				}
				if err := tx.Commit(ctx); err != nil {
					return measure.Outcome{}, err
				}
				c.led.set(ip, key, next)
				c.led.bumped(ip)
				c.led.acked.Add(1)
				c.led.incAcked.Add(1)
				return measure.Outcome{Retries: attempt}, nil
			}
			c.led.retries.Add(1)
			return measure.Outcome{}, errors.New("optimistic retries exhausted")

		case ConcNone:
			// The drift control: publications race, but what breaks is the
			// derived rollup, which each of them writes from a count read before
			// the others committed.
			if c.design.Rollup != RollupDrift {
				return measure.Outcome{}, fmt.Errorf("design %s has no contention strategy", c.design.ID)
			}
			newKey := c.uniqueKey("drift")
			val := fmt.Sprintf("drift-%d", ip)
			err := c.publishDrift(ctx, ip, func(q ports.Queryer) error {
				_, e := c.cat.exec(ctx, q, sW02, map[string]any{
					"installed_product_id": ip, "key": newKey, "value": val})
				return e
			})
			if err == nil {
				c.led.set(ip, newKey, val)
				c.led.acked.Add(1)
			}
			return measure.Outcome{}, err

		case ConcPessimistic, ConcUnsafe:
			tx, err := c.db.Begin(ctx, ports.ReadCommitted)
			if err != nil {
				return measure.Outcome{}, err
			}
			if c.design.Conc == ConcPessimistic {
				a, _ := c.cat.args("wc2_lock_installed_product", map[string]any{"installed_product_id": ip})
				var rev int64
				if err := tx.QueryRow(ctx, c.cat.stmt("wc2_lock_installed_product").SQL, a...).Scan(&rev); err != nil {
					_ = tx.Rollback(ctx)
					return measure.Outcome{}, err
				}
			}
			read, write := "wc2_read_value", "wc2_write_value"
			if c.design.Conc == ConcUnsafe {
				// The control reads and writes without the lock above.
				read, write = "wx1_read_value", "wx1_write_value"
			}
			a, _ := c.cat.args(read, map[string]any{"installed_product_id": ip, "key": key})
			var value string
			var version int64
			if err := tx.QueryRow(ctx, c.cat.stmt(read).SQL, a...).Scan(&value, &version); err != nil {
				_ = tx.Rollback(ctx)
				return measure.Outcome{}, err
			}
			next := strconv.Itoa(atoi(value) + 1)
			if _, err := c.cat.exec(ctx, tx, write, map[string]any{
				"installed_product_id": ip, "key": key, "value": next}); err != nil {
				_ = tx.Rollback(ctx)
				return measure.Outcome{}, err
			}
			if _, err := c.cat.exec(ctx, tx, sRevBump, map[string]any{"installed_product_id": ip}); err != nil {
				_ = tx.Rollback(ctx)
				return measure.Outcome{}, err
			}
			if err := tx.Commit(ctx); err != nil {
				return measure.Outcome{}, err
			}
			c.led.set(ip, key, next)
			c.led.bumped(ip)
			c.led.acked.Add(1)
			c.led.incAcked.Add(1)
			return measure.Outcome{}, nil
		}
		return measure.Outcome{}, fmt.Errorf("design %s has no contention strategy", c.design.ID)
	}
}

// ---------------------------------------------------------------- small helpers

func atoi(s string) int {
	n, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return 0
	}
	return n
}

