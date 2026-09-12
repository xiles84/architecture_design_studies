-- Query catalogue for D5 (rollup-app).
--
-- Byte-identical to D4 apart from this header. D4 and D5 store the same
-- aggregates and read them the same way; they differ only in who maintains them
-- on the write path. Reads are expected to match, and any measured difference
-- between them on the read side is noise -- which makes this pair a useful
-- built-in check on the harness's own measurement error.

-- name: q01_last_donation_global
-- params: none
-- There are only a handful of charities, so the global maximum is a tiny scan
-- over charity followed by one primary-key fetch of the donation itself.
SELECT d.donation_id, d.person_id, d.amount_cents, d.currency, d.donated_at
  FROM (SELECT last_donation_id
          FROM charity
         WHERE last_donation_id IS NOT NULL
         ORDER BY last_donation_at DESC
         LIMIT 1) c
  JOIN donation d ON d.donation_id = c.last_donation_id;

-- name: q02_last_donation_charity
-- params: charity_id
-- One primary-key read of charity, one primary-key read of donation.
SELECT d.donation_id, d.person_id, d.amount_cents, d.currency, d.donated_at
  FROM charity c
  JOIN donation d ON d.donation_id = c.last_donation_id
 WHERE c.charity_id = $1;

-- name: q03_top_donor_charity
-- params: charity_id
-- The GROUP BY is gone. person_top_donor_idx is ordered by total descending
-- within a charity, so this reads exactly one index entry.
SELECT person_id, full_name, total_donated_cents AS total_cents
  FROM person
 WHERE charity_id = $1
 ORDER BY total_donated_cents DESC
 LIMIT 1;

-- name: q04_top_donors_leaderboard
-- params: charity_id
SELECT person_id, full_name, total_donated_cents AS total_cents
  FROM person
 WHERE charity_id = $1
 ORDER BY total_donated_cents DESC
 LIMIT 10;

-- name: q05_last_donor_charity
-- params: charity_id
SELECT c.last_donor_person_id AS person_id, p.full_name, c.last_donation_at AS donated_at
  FROM charity c
  JOIN person p ON p.person_id = c.last_donor_person_id
 WHERE c.charity_id = $1;

-- name: q06_person_first_last
-- params: person_id
-- Single-row primary-key read; the MIN/MAX were computed at write time.
SELECT first_donation_at AS first_at, last_donation_at AS last_at
  FROM person
 WHERE person_id = $1;

-- name: q07_total_donated_global
-- params: none
-- Sums one row per charity instead of one row per donation.
SELECT COALESCE(SUM(total_donated_cents), 0) AS total_cents
  FROM charity;

-- name: q08_total_donated_charity
-- params: charity_id
SELECT total_donated_cents AS total_cents
  FROM charity
 WHERE charity_id = $1;

-- name: q09_person_recent_donations
-- params: person_id
-- Unchanged from D3: rollups store answers, not rows.
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
SELECT donation_count
  FROM person
 WHERE person_id = $1;

-- name: q12_charity_recent_feed
-- params: charity_id
SELECT d.donation_id, d.amount_cents, d.donated_at, p.full_name
  FROM donation d
  JOIN person p ON p.person_id = d.person_id
 WHERE d.charity_id = $1
 ORDER BY d.donated_at DESC
 LIMIT 50;
