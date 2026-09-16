-- Query catalogue for D3 (flattened-fk).
--
-- Every charity-scoped query that needed a join through person in D1/D2 now
-- filters donation directly. The SQL gets shorter, which is itself part of the
-- finding: denormalising the grandparent key removes a join from six queries.

-- name: q01_last_donation_global
-- params: none
SELECT donation_id, person_id, amount_cents, currency, donated_at
  FROM donation
 ORDER BY donated_at DESC
 LIMIT 1;

-- name: q02_last_donation_charity
-- params: charity_id
-- No join: donation_charity_time_idx answers this with a one-row index scan.
SELECT donation_id, person_id, amount_cents, currency, donated_at
  FROM donation
 WHERE charity_id = $1
 ORDER BY donated_at DESC
 LIMIT 1;

-- name: q03_top_donor_charity
-- params: charity_id
-- Still a full aggregate over the charity's donations -- denormalising the key
-- removes the join but does NOT remove the GROUP BY. Only a rollup (D4/D5) can.
-- person is joined back for the name, but against a single already-chosen row.
SELECT p.person_id, p.full_name, t.total_cents
  FROM (SELECT person_id, SUM(amount_cents) AS total_cents
          FROM donation
         WHERE charity_id = $1
         GROUP BY person_id
         ORDER BY total_cents DESC
         LIMIT 1) t
  JOIN person p ON p.person_id = t.person_id;

-- name: q04_top_donors_leaderboard
-- params: charity_id
SELECT p.person_id, p.full_name, t.total_cents
  FROM (SELECT person_id, SUM(amount_cents) AS total_cents
          FROM donation
         WHERE charity_id = $1
         GROUP BY person_id
         ORDER BY total_cents DESC
         LIMIT 10) t
  JOIN person p ON p.person_id = t.person_id
 ORDER BY t.total_cents DESC;

-- name: q05_last_donor_charity
-- params: charity_id
SELECT d.person_id, p.full_name, d.donated_at
  FROM (SELECT person_id, donated_at
          FROM donation
         WHERE charity_id = $1
         ORDER BY donated_at DESC
         LIMIT 1) d
  JOIN person p ON p.person_id = d.person_id;

-- name: q06_person_first_last
-- params: person_id
SELECT MIN(donated_at) AS first_at, MAX(donated_at) AS last_at
  FROM donation
 WHERE person_id = $1;

-- name: q07_total_donated_global
-- params: none
SELECT COALESCE(SUM(amount_cents), 0) AS total_cents
  FROM donation;

-- name: q08_total_donated_charity
-- params: charity_id
-- donation_charity_amount_idx covers both the predicate and the summed column,
-- so this should be an index-only scan that never touches the heap.
SELECT COALESCE(SUM(amount_cents), 0) AS total_cents
  FROM donation
 WHERE charity_id = $1;

-- name: q09_person_recent_donations
-- params: person_id
SELECT donation_id, amount_cents, currency, donated_at, note
  FROM donation
 WHERE person_id = $1
 ORDER BY donated_at DESC
 LIMIT 20;

-- name: q10_donation_by_id
-- params: donation_id
SELECT donation_id, person_id, amount_cents, currency, donated_at, note
  FROM donation
 WHERE donation_id = $1;

-- name: q11_person_donation_count
-- params: person_id
SELECT COUNT(*) AS donation_count
  FROM donation
 WHERE person_id = $1;

-- name: q12_charity_recent_feed
-- params: charity_id
SELECT d.donation_id, d.amount_cents, d.donated_at, p.full_name
  FROM donation d
  JOIN person p ON p.person_id = d.person_id
 WHERE d.charity_id = $1
 ORDER BY d.donated_at DESC
 LIMIT 50;

-- ---------------------------------------------------------------------------
-- Recency-window questions (study 01 v4, RECENCY.md). See d3's queries.sql
-- for the full explanation of the two window regimes.
--
-- D19 is D3's exact schema, indexes and write path -- the only change is this
-- formulation: read the WINDOW first (an index range scan on donated_at),
-- then confirm each candidate donor gave nothing after the window closes.
-- Correct in both regimes -- a donor qualifies exactly when they gave inside
-- the window AND gave nothing after it. Costs what the window holds, unlike
-- D18's parent-driven probe (d18_recency_probe/queries.sql), which costs one
-- descent per donor whatever the window holds. Isolates SQL formulation
-- alone: see d18 -> d19 in RECENCY.md section 4.
-- ---------------------------------------------------------------------------

-- name: q13_donors_last_gift_window
-- params: since, until
SELECT t.person_id, p.full_name, t.last_at
  FROM (SELECT d.person_id, MAX(d.donated_at) AS last_at
          FROM donation d
         WHERE d.donated_at >= $1 AND d.donated_at < $2
         GROUP BY d.person_id) t
  JOIN person p ON p.person_id = t.person_id
 WHERE NOT EXISTS (SELECT 1 FROM donation d2
                     WHERE d2.person_id = t.person_id
                       AND d2.donated_at >= $2)
 ORDER BY t.last_at DESC
 LIMIT 100;

-- name: q14_donors_last_gift_window_count
-- params: since, until
SELECT COUNT(*) AS donor_count
  FROM (SELECT d.person_id, MAX(d.donated_at) AS last_at
          FROM donation d
         WHERE d.donated_at >= $1 AND d.donated_at < $2
         GROUP BY d.person_id) t
 WHERE NOT EXISTS (SELECT 1 FROM donation d2
                     WHERE d2.person_id = t.person_id
                       AND d2.donated_at >= $2);

-- name: q15_charity_donors_last_gift_window
-- params: charity_id, since, until
SELECT t.person_id, p.full_name, t.last_at
  FROM (SELECT d.person_id, MAX(d.donated_at) AS last_at
          FROM donation d
         WHERE d.charity_id = $1 AND d.donated_at >= $2 AND d.donated_at < $3
         GROUP BY d.person_id) t
  JOIN person p ON p.person_id = t.person_id
 WHERE NOT EXISTS (SELECT 1 FROM donation d2
                     WHERE d2.person_id = t.person_id
                       AND d2.donated_at >= $3)
 ORDER BY t.last_at DESC
 LIMIT 100;

-- name: q16_lapsed_donors_count
-- params: since
-- The window-first idea does not apply to "everyone whose last gift is
-- BEFORE since": there is no upper-bounded window to read first, so this
-- degenerates to the same statement D3 (formulation B) uses. Recorded here
-- rather than left out, per RECENCY.md's handoff instruction.
SELECT COUNT(*) AS donor_count
  FROM (SELECT person_id
          FROM donation
         GROUP BY person_id
        HAVING MAX(donated_at) < $1) t;
