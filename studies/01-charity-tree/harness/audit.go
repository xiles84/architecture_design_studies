package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ---------------------------------------------------------------------------
// Rollup audit
//
// The verification pass (verify.go) proves a design answers correctly on
// freshly loaded, quiescent data. This one asks the harder question:
//
//	after N concurrent writers hammered the same parent rows, do the
//	consolidated aggregates still agree with the children they summarise?
//
// That is the whole point of experiment C. "Optimistic vs pessimistic" is not a
// throughput question with a correctness footnote -- a strategy that is 3x
// faster and silently loses increments has not won anything. A lost update here
// means a charity's published total is quietly wrong, and nothing in the system
// would ever notice.
//
// The audit is a full recomputation from the donation table, so it is expensive
// and deliberately run only after a contended write benchmark.
// ---------------------------------------------------------------------------

type RollupAudit struct {
	Ran               bool  `json:"ran"`
	PersonRows        int64 `json:"person_rows_checked"`
	PersonMismatches  int64 `json:"person_mismatches"`
	CharityRows       int64 `json:"charity_rows_checked"`
	CharityMismatches int64 `json:"charity_mismatches"`
	// DriftCents is the signed error in the headline number: what the rollups
	// claim was donated minus what the donation rows actually add up to.
	DriftCents int64  `json:"drift_cents"`
	Note       string `json:"note,omitempty"`

	// D9 only: does the bounded embedded cache still hold exactly the newest
	// donations, in the right order?
	CacheChecked    int64 `json:"cache_rows_checked,omitempty"`
	CacheMismatches int64 `json:"cache_mismatches,omitempty"`
	// Up to five examples with a diagnosis, recorded only when the cache is wrong.
	CacheExamples []CacheMismatch `json:"cache_mismatch_examples,omitempty"`
}

// AuditRecentCache checks D9's embedded slice against the donation table.
//
// It compares the SEQUENCE OF DONATION IDS rather than the raw JSONB, on
// purpose. The loader writes timestamps as RFC 3339 with a "Z" suffix while
// PostgreSQL's own jsonb_build_object renders them with "+00:00"; both parse to
// the same instant and every query casts before use, but a byte comparison would
// call every trigger-touched row a mismatch and report a working cache as broken.
// Comparing ids tests what actually matters: the right children, in the right
// order, and none missing.
func AuditRecentCache(ctx context.Context, pool *pgxpool.Pool, d Design) (int64, int64, error) {
	if !d.RecentCache {
		return 0, 0, nil
	}
	const q = `
SELECT COUNT(*) AS checked,
       COUNT(*) FILTER (WHERE cached IS DISTINCT FROM actual) AS mismatches
  FROM (
    SELECT (SELECT COALESCE(ARRAY_AGG((t.e ->> 'i')::BIGINT ORDER BY t.ord), '{}')
              FROM jsonb_array_elements(p.recent_donations) WITH ORDINALITY AS t(e, ord)
           ) AS cached,
           (SELECT COALESCE(ARRAY_AGG(x.donation_id ORDER BY x.donated_at DESC, x.donation_id DESC), '{}')
              FROM (SELECT donation_id, donated_at
                      FROM donation
                     WHERE person_id = p.person_id
                     ORDER BY donated_at DESC, donation_id DESC
                     LIMIT $1) x
           ) AS actual
      FROM person p
  ) s`
	var checked, bad int64
	if err := pool.QueryRow(ctx, q, recentCacheSize).Scan(&checked, &bad); err != nil {
		return 0, 0, fmt.Errorf("recent-cache audit: %w", err)
	}
	return checked, bad, nil
}

// CacheMismatch is one person whose embedded cache disagrees with the table.
//
// Recording examples, not just a count, is what turns "the cache drifted" into a
// diagnosis. The two plausible failure modes leave different fingerprints:
//
//   - ORDERING race: the same ids appear in both lists, in a different order.
//     Two inserts for one donor committed in the opposite order to their
//     donated_at, and prepending put the later-committed one first.
//   - LOST UPDATE: an id is missing from (or extra in) the cache. A rebuild on
//     the update/delete path read its snapshot, a concurrent insert committed,
//     and the rebuild's UPDATE then overwrote the cache without it.
type CacheMismatch struct {
	PersonID  int64   `json:"person_id"`
	Cached    []int64 `json:"cached_ids"`
	Actual    []int64 `json:"actual_ids"`
	SameSet   bool    `json:"same_set"`
	Diagnosis string  `json:"diagnosis"`
}

// CacheMismatchExamples returns up to limit mismatching persons with a diagnosis.
func CacheMismatchExamples(ctx context.Context, pool *pgxpool.Pool, limit int) ([]CacheMismatch, error) {
	const q = `
SELECT person_id, cached, actual FROM (
    SELECT p.person_id,
           (SELECT COALESCE(ARRAY_AGG((t.e ->> 'i')::BIGINT ORDER BY t.ord), '{}')
              FROM jsonb_array_elements(p.recent_donations) WITH ORDINALITY AS t(e, ord)) AS cached,
           (SELECT COALESCE(ARRAY_AGG(x.donation_id ORDER BY x.donated_at DESC, x.donation_id DESC), '{}')
              FROM (SELECT donation_id, donated_at FROM donation
                     WHERE person_id = p.person_id
                     ORDER BY donated_at DESC, donation_id DESC LIMIT $1) x) AS actual
      FROM person p
) s
WHERE cached IS DISTINCT FROM actual
LIMIT $2`
	rows, err := pool.Query(ctx, q, recentCacheSize, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []CacheMismatch
	for rows.Next() {
		var m CacheMismatch
		if err := rows.Scan(&m.PersonID, &m.Cached, &m.Actual); err != nil {
			return nil, err
		}
		m.SameSet = sameSet(m.Cached, m.Actual)
		if m.SameSet {
			m.Diagnosis = "ordering race: same donations, different order"
		} else {
			m.Diagnosis = "lost or extra element: cache holds a different set of donations"
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func sameSet(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	seen := map[int64]int{}
	for _, v := range a {
		seen[v]++
	}
	for _, v := range b {
		seen[v]--
		if seen[v] < 0 {
			return false
		}
	}
	return true
}

// AuditRollups recomputes every aggregate from the donation table and counts the
// parent rows that disagree.
func AuditRollups(ctx context.Context, pool *pgxpool.Pool, d Design) (*RollupAudit, error) {
	a := &RollupAudit{}
	if d.SumOnly {
		a.Ran = true
		for _, table := range []string{"person", "charity"} {
			key := table + "_id"
			q := "SELECT COUNT(*), COUNT(*) FILTER (WHERE p.total_donated_cents IS DISTINCT FROM COALESCE(a.total,0)) FROM " + table + " p LEFT JOIN (SELECT " + key + ", SUM(amount_cents) total FROM donation GROUP BY " + key + ") a ON p." + key + "=a." + key
			var checked, bad int64
			if err := pool.QueryRow(ctx, q).Scan(&checked, &bad); err != nil {
				return nil, err
			}
			if table == "person" {
				a.PersonRows, a.PersonMismatches = checked, bad
			} else {
				a.CharityRows, a.CharityMismatches = checked, bad
			}
		}
		return a, nil
	}

	if d.RecentCache {
		checked, bad, err := AuditRecentCache(ctx, pool, d)
		if err != nil {
			return nil, err
		}
		a.Ran = true
		a.CacheChecked, a.CacheMismatches = checked, bad
		if bad > 0 {
			ex, err := CacheMismatchExamples(ctx, pool, 5)
			if err != nil {
				return nil, err
			}
			a.CacheExamples = ex
		}
	}

	if !d.Rollups {
		if !a.Ran {
			a.Note = "design carries no consolidated aggregates; nothing to audit"
		}
		return a, nil
	}
	a.Ran = true

	const personSQL = `
SELECT COUNT(*) AS checked,
       COUNT(*) FILTER (
           WHERE p.donation_count      IS DISTINCT FROM COALESCE(d.cnt, 0)
              OR p.total_donated_cents IS DISTINCT FROM COALESCE(d.total, 0)
       ) AS mismatches
  FROM person p
  LEFT JOIN (
        SELECT person_id, COUNT(*) AS cnt, SUM(amount_cents) AS total
          FROM donation
         GROUP BY person_id
  ) d ON d.person_id = p.person_id`

	if err := pool.QueryRow(ctx, personSQL).Scan(&a.PersonRows, &a.PersonMismatches); err != nil {
		return nil, fmt.Errorf("person audit: %w", err)
	}

	const charitySQL = `
SELECT COUNT(*) AS checked,
       COUNT(*) FILTER (
           WHERE c.donation_count      IS DISTINCT FROM COALESCE(d.cnt, 0)
              OR c.total_donated_cents IS DISTINCT FROM COALESCE(d.total, 0)
       ) AS mismatches,
       COALESCE(SUM(c.total_donated_cents) - SUM(COALESCE(d.total, 0)), 0) AS drift
  FROM charity c
  LEFT JOIN (
        SELECT charity_id, COUNT(*) AS cnt, SUM(amount_cents) AS total
          FROM donation
         GROUP BY charity_id
  ) d ON d.charity_id = c.charity_id`

	var drift any
	if err := pool.QueryRow(ctx, charitySQL).Scan(&a.CharityRows, &a.CharityMismatches, &drift); err != nil {
		return nil, fmt.Errorf("charity audit: %w", err)
	}
	a.DriftCents = toI64(drift)
	return a, nil
}
