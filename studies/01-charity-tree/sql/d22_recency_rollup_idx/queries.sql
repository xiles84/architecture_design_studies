-- Query catalogue for D4 (rollup-trigger).
--
-- Seven of the twelve questions no longer touch the donation table at all: the
-- answer was consolidated into a parent row at write time. Those queries become
-- O(1) primary-key reads regardless of how many donations exist.
--
-- Queries that still need the child rows (q09, q10, q12) are unchanged from D3 --
-- a rollup cannot answer "show me the actual rows".

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

-- ---------------------------------------------------------------------------
-- Recency-window questions (study 01 v4, RECENCY.md). See d1's queries.sql
-- for the full explanation of the two window regimes.
--
-- person.last_donation_at is already maintained (D4's trigger / D5's
-- application path); a donor with no donations has it NULL, which the
-- comparisons below correctly exclude. This is the rollup answer, and it
-- has no global access path yet -- see D22/D23, which add exactly one index
-- to this same column and nothing else.
-- ---------------------------------------------------------------------------

-- name: q13_donors_last_gift_window
-- params: since, until
SELECT person_id, full_name, last_donation_at AS last_at
  FROM person
 WHERE last_donation_at >= $1 AND last_donation_at < $2
 ORDER BY last_donation_at DESC
 LIMIT 100;

-- name: q14_donors_last_gift_window_count
-- params: since, until
SELECT COUNT(*) AS donor_count
  FROM person
 WHERE last_donation_at >= $1 AND last_donation_at < $2;

-- name: q15_charity_donors_last_gift_window
-- params: charity_id, since, until
SELECT person_id, full_name, last_donation_at AS last_at
  FROM person
 WHERE charity_id = $1 AND last_donation_at >= $2 AND last_donation_at < $3
 ORDER BY last_donation_at DESC
 LIMIT 100;

-- name: q16_lapsed_donors_count
-- params: since
SELECT COUNT(*) AS donor_count
  FROM person
 WHERE last_donation_at < $1;
