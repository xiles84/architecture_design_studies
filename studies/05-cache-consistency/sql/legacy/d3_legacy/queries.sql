-- Legacy read catalogue.
--
-- Every design in Study 05 returns the SAME column list for these statements, so
-- one harness scan path serves them all and the canonical payload encoder is the
-- single place the content hash comes from. A design that returned a different
-- shape would be comparing its own JSON renderer, not its cache.

-- name: r_portal_person
-- params: person_id
-- The donor identity block plus the two aggregates the portal shows.
--
-- One statement, so a fill describes one committed instant: the person row, its
-- charity and the aggregates come from the same statement snapshot. Splitting it
-- into three round trips would give the filler three different database states to
-- stitch together, and the content hash would then be a hash of the harness's
-- stitching rather than of any state the database ever held.
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
-- The newest 20 donations in the deterministic order the cached representation
-- requires. `note` may be NULL; the payload encoder treats SQL NULL and JSON null
-- as the same value, which is why the embedded reference design (whose array
-- stores JSON null) still produces an identical content hash.
SELECT donation_id, amount_cents, currency, donated_at, note
  FROM donation
 WHERE person_id = $1
 ORDER BY donated_at DESC, donation_id DESC
 LIMIT 20;

-- name: r_charity_recent
-- params: charity_id
-- The blended workload's UNCACHEABLE operation: the charity's last 50 donations
-- with donor names. It is not cached in this study (the study's cache key is
-- donor-scoped), and it is here because it is the read that discriminates D3's
-- rolldown from the normalized reference: on the legacy model it filters on the
-- copied `d.charity_id` and never touches `person` for the predicate.
SELECT d.donation_id, d.amount_cents, d.donated_at, p.full_name
  FROM donation d
  JOIN person p ON p.person_id = d.person_id
 WHERE d.charity_id = $1
 ORDER BY d.donated_at DESC, d.donation_id DESC
 LIMIT 50;
