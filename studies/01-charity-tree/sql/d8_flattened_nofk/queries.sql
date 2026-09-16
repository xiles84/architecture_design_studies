-- Query catalogue for D8 (flattened-nofk).
--
-- Byte-identical to D3. Foreign keys are a write-path and integrity feature; they
-- change no read plan here, so any read difference between D3 and D8 is noise.
-- Keeping the file rather than symlinking it keeps each design readable on its own.

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
-- Recency-window questions (study 01 v4, RECENCY.md). See d1's queries.sql
-- for the full explanation of the two window regimes and why a formulation
-- that only works in one is a trap.
--
-- This design carries the denormalised charity_id on donation, so q15
-- filters donation directly -- the same one-decision difference already
-- measured in q02/q05/q12.
-- ---------------------------------------------------------------------------

-- name: q13_donors_last_gift_window
-- params: since, until
SELECT t.person_id, p.full_name, t.last_at
  FROM (SELECT person_id, MAX(donated_at) AS last_at
          FROM donation
         GROUP BY person_id
        HAVING MAX(donated_at) >= $1 AND MAX(donated_at) < $2
         ORDER BY 2 DESC
         LIMIT 100) t
  JOIN person p ON p.person_id = t.person_id
 ORDER BY t.last_at DESC;

-- name: q14_donors_last_gift_window_count
-- params: since, until
SELECT COUNT(*) AS donor_count
  FROM (SELECT person_id
          FROM donation
         GROUP BY person_id
        HAVING MAX(donated_at) >= $1 AND MAX(donated_at) < $2) t;

-- name: q15_charity_donors_last_gift_window
-- params: charity_id, since, until
SELECT t.person_id, p.full_name, t.last_at
  FROM (SELECT person_id, MAX(donated_at) AS last_at
          FROM donation
         WHERE charity_id = $1
         GROUP BY person_id
        HAVING MAX(donated_at) >= $2 AND MAX(donated_at) < $3
         ORDER BY 2 DESC
         LIMIT 100) t
  JOIN person p ON p.person_id = t.person_id
 ORDER BY t.last_at DESC;

-- name: q16_lapsed_donors_count
-- params: since
SELECT COUNT(*) AS donor_count
  FROM (SELECT person_id
          FROM donation
         GROUP BY person_id
        HAVING MAX(donated_at) < $1) t;
