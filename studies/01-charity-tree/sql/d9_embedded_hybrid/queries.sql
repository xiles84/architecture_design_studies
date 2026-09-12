-- Query catalogue for D9 (embedded-hybrid).
--
-- Identical to D3 apart from q09, which is the point of the design: the
-- highest-QPS query in the catalogue is answered from the parent row instead of
-- from the child table. Everything else still uses the real, fully indexed
-- donation table, so none of D6's cross-person penalties apply here.

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
-- The design's whole reason to exist. Answered entirely from the parent row:
-- one primary-key read, no child table, no index descent, no sort. The cache is
-- already stored newest-first and already capped at 20, so array position IS
-- the answer order.
SELECT (e.value ->> 'i')::BIGINT      AS donation_id,
       (e.value ->> 'a')::BIGINT      AS amount_cents,
       e.value ->> 'c'                AS currency,
       (e.value ->> 't')::TIMESTAMPTZ AS donated_at,
       e.value ->> 'n'                AS note
  FROM person p
 CROSS JOIN LATERAL jsonb_array_elements(p.recent_donations) WITH ORDINALITY AS e(value, ord)
 WHERE p.person_id = $1
 ORDER BY e.ord;

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
