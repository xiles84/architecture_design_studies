-- normalized-indexed reads. Same column lists as every other design.
--
-- The charity feed is the discriminating read: without the copied key, the
-- predicate has to be answered through the person table, so the plan differs even
-- though the rows returned do not.

-- name: r_portal_person
-- params: person_id
SELECT p.person_id,
       p.full_name,
       p.email,
       p.joined_at,
       p.charity_id,
       c.name    AS charity_name,
       c.country AS charity_country,
       (SELECT count(*) FROM donation d WHERE d.person_id = p.person_id)                    AS donation_count,
       (SELECT COALESCE(sum(d.amount_cents), 0) FROM donation d WHERE d.person_id = p.person_id) AS donation_total_cents
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
-- The rolldown cost, in one statement: the charity predicate lives in `person`, so
-- this read reaches every donation of every person in the charity rather than the
-- charity's slice of a table already keyed by charity.
SELECT d.donation_id, d.amount_cents, d.donated_at, p.full_name
  FROM donation d
  JOIN person p ON p.person_id = d.person_id
 WHERE p.charity_id = $1
 ORDER BY d.donated_at DESC, d.donation_id DESC
 LIMIT 50;
