-- embedded-locked reads. Same column lists as every design; the recent slice comes
-- from the parent's stored array instead of a second query against the child table.
-- That is the single decision this cell isolates.
--
-- The array is read through jsonb_array_elements in the STORED order, which the
-- trigger maintains as (donated_at DESC, donation_id DESC). The harness re-sorts
-- canonically anyway, so a slice that merely arrived in a different order would not
-- be called wrong -- element membership and values are what the audit compares.

-- name: r_portal_person
-- params: person_id
-- The aggregates are still derived from the child table: only the recent slice is
-- embedded, so the pair isolates embedding and not "store everything twice".
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
SELECT (e ->> 'id')::BIGINT          AS donation_id,
       (e ->> 'amt')::BIGINT         AS amount_cents,
       (e ->> 'cur')::CHAR(3)        AS currency,
       (e ->> 'at')::TIMESTAMPTZ     AS donated_at,
       (e ->> 'note')                AS note
  FROM person p,
       jsonb_array_elements(p.recent_donations) AS e
 WHERE p.person_id = $1;

-- name: r_charity_recent
-- params: charity_id
SELECT d.donation_id, d.amount_cents, d.donated_at, p.full_name
  FROM donation d
  JOIN person p ON p.person_id = d.person_id
 WHERE d.charity_id = $1
 ORDER BY d.donated_at DESC, d.donation_id DESC
 LIMIT 50;
