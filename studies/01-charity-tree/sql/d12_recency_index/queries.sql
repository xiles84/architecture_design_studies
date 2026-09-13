-- d12_recency_index: controlled variant; see ../../ENHANCEMENTS.md.

-- name: q01_last_donation_global
-- params: none
-- "Information of the last donation." No charity_id on donation, no index on
-- donated_at: this has to read the whole table and keep a running top-1.
SELECT donation_id, person_id, amount_cents, currency, donated_at
  FROM donation
 ORDER BY donated_at DESC
 LIMIT 1;

-- name: q02_last_donation_charity
-- params: charity_id
-- Same question scoped to one charity. The charity is two levels up the tree,
-- so the scope predicate can only be applied after joining through person.
SELECT d.donation_id, d.person_id, d.amount_cents, d.currency, d.donated_at
  FROM donation d
  JOIN person p ON p.person_id = d.person_id
 WHERE p.charity_id = $1
 ORDER BY d.donated_at DESC
 LIMIT 1;

-- name: q03_top_donor_charity
-- params: charity_id
-- "Who is the person that donates more." Derived: aggregate every donation of
-- every person in the charity, then take the maximum.
SELECT p.person_id, p.full_name, SUM(d.amount_cents) AS total_cents
  FROM person p
  JOIN donation d ON d.person_id = p.person_id
 WHERE p.charity_id = $1
 GROUP BY p.person_id, p.full_name
 ORDER BY total_cents DESC
 LIMIT 1;

-- name: q04_top_donors_leaderboard
-- params: charity_id
-- The same aggregation, but returning a top-10 board. Included because a
-- leaderboard is the realistic version of the question and because LIMIT 10 vs
-- LIMIT 1 changes nothing about the cost -- the whole aggregate still runs.
SELECT p.person_id, p.full_name, SUM(d.amount_cents) AS total_cents
  FROM person p
  JOIN donation d ON d.person_id = p.person_id
 WHERE p.charity_id = $1
 GROUP BY p.person_id, p.full_name
 ORDER BY total_cents DESC
 LIMIT 10;

-- name: q05_last_donor_charity
-- params: charity_id
-- "Who was the last person that donated."
SELECT p.person_id, p.full_name, d.donated_at
  FROM donation d
  JOIN person p ON p.person_id = d.person_id
 WHERE p.charity_id = $1
 ORDER BY d.donated_at DESC
 LIMIT 1;

-- name: q06_person_first_last
-- params: person_id
-- "When was the first and last donation of a person."
SELECT MIN(donated_at) AS first_at, MAX(donated_at) AS last_at
  FROM donation
 WHERE person_id = $1;

-- name: q07_total_donated_global
-- params: none
-- "How much money was donated in total." Unavoidably a full aggregate here.
SELECT COALESCE(SUM(amount_cents), 0) AS total_cents
  FROM donation;

-- name: q08_total_donated_charity
-- params: charity_id
SELECT COALESCE(SUM(d.amount_cents), 0) AS total_cents
  FROM donation d
  JOIN person p ON p.person_id = d.person_id
 WHERE p.charity_id = $1;

-- name: q09_person_recent_donations
-- params: person_id
-- The bread-and-butter OLTP read: one donor's recent history. Typically the
-- highest-QPS query in a system like this, so its cost dominates real workloads
-- even though it is the least interesting question analytically.
SELECT donation_id, amount_cents, currency, donated_at, note
  FROM donation
 WHERE person_id = $1
 ORDER BY donated_at DESC
 LIMIT 20;

-- name: q10_donation_by_id
-- params: donation_id
-- Primary-key point lookup. The control: every design should be fast here, and
-- a design that is not has done something expensive to its identity access path.
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
-- "Recent activity" page: the last 50 donations to a charity with donor names.
SELECT d.donation_id, d.amount_cents, d.donated_at, p.full_name
  FROM donation d
  JOIN person p ON p.person_id = d.person_id
 WHERE p.charity_id = $1
 ORDER BY d.donated_at DESC
 LIMIT 50;
