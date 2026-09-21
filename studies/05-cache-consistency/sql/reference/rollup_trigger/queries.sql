-- rollup-trigger reads. Same column lists as every design; the portal read stops
-- aggregating and reads the two stored columns instead. That is the single
-- decision this cell isolates against the flattened baseline.

-- name: r_portal_person
-- params: person_id
-- p.donation_count / p.donation_total_cents replace the two correlated
-- subqueries. Whether that is actually cheaper is what the plan and the timing
-- say; on a warm cache the subqueries are two index descents, and this is a
-- wider parent row.
SELECT p.person_id,
       p.full_name,
       p.email,
       p.joined_at,
       p.charity_id,
       c.name    AS charity_name,
       c.country AS charity_country,
       p.donation_count       AS donation_count,
       p.donation_total_cents AS donation_total_cents
  FROM person p
  JOIN charity c ON c.charity_id = p.charity_id
 WHERE p.person_id = $1;

-- name: r_portal_recent
-- params: person_id
SELECT donation_id, amount_cents, currency, donated_at, note
  FROM donation
 WHERE person_id = $1
 ORDER BY donated_at DESC, donation_id DESC
 LIMIT 20;

-- name: r_charity_recent
-- params: charity_id
SELECT d.donation_id, d.amount_cents, d.donated_at, p.full_name
  FROM donation d
  JOIN person p ON p.person_id = d.person_id
 WHERE d.charity_id = $1
 ORDER BY d.donated_at DESC, d.donation_id DESC
 LIMIT 50;
