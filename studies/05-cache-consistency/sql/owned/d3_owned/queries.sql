-- Owned read catalogue.
--
-- The three read statements are byte-identical to the legacy model's, so a
-- difference in a cached read is not explained by a different read. The model adds
-- exactly one statement: r_person_version, the authoritative token a strict cache
-- validates against when it cannot trust its own invalidation -- which is the
-- multi-instance local-cache case, and the legacy model has no equivalent.

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
SELECT d.donation_id, d.amount_cents, d.donated_at, p.full_name
  FROM donation d
  JOIN person p ON p.person_id = d.person_id
 WHERE d.charity_id = $1
 ORDER BY d.donated_at DESC, d.donation_id DESC
 LIMIT 50;

-- name: r_person_version
-- params: person_id
-- One indexed primary-key lookup, and the whole reason the owned model exists: it
-- lets a strict cache prove freshness against a single monotonic token instead of
-- re-reading the portal view. What that proof costs per read is measured, not
-- assumed -- it is a database round trip on every cache hit, and in a
-- single-instance deployment it is a cost the cache does not have to pay.
SELECT cache_version FROM person WHERE person_id = $1;
