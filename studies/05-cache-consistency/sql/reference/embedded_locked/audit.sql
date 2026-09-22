-- embedded-locked audits: the baseline audits plus the invariant an embedded slice
-- creates.
--
-- a_embedded_drift compares the ORDERED ID SEQUENCE of the stored slice with the
-- newest 20 of the child table, not raw JSONB bytes. Study 01 learned that the
-- expensive way: the loader and PostgreSQL render timestamps differently in JSON,
-- and a byte comparison called a working cache broken. Sequence comparison asks the
-- question the invariant actually is -- does the slice hold exactly the newest 20,
-- in the right order -- and the Go oracle checks the values.

-- name: a_person_aggregate
-- params: person_id
SELECT count(*)                                  AS donation_count,
       COALESCE(sum(amount_cents), 0)            AS donation_total_cents
  FROM donation
 WHERE person_id = $1;

-- name: a_donation_ids
-- params: person_id
SELECT donation_id
  FROM donation
 WHERE person_id = $1
 ORDER BY donation_id;

-- name: a_recent_ids
-- params: person_id
SELECT donation_id
  FROM donation
 WHERE person_id = $1
 ORDER BY donated_at DESC, donation_id DESC
 LIMIT 20;

-- name: a_rolldown_drift
-- params: none
SELECT count(*)
  FROM donation d
  JOIN person p ON p.person_id = d.person_id
 WHERE d.charity_id <> p.charity_id;

-- name: a_embedded_drift
-- params: none
SELECT count(*)
  FROM person p
 WHERE COALESCE((
         SELECT jsonb_agg(x.id ORDER BY x.at DESC, x.id DESC)
           FROM (
                 SELECT (e ->> 'id')::BIGINT      AS id,
                        (e ->> 'at')::TIMESTAMPTZ AS at
                   FROM jsonb_array_elements(p.recent_donations) AS e
                ) x
       ), '[]'::jsonb)
   IS DISTINCT FROM
       COALESCE((
         SELECT jsonb_agg(d.donation_id ORDER BY d.donated_at DESC, d.donation_id DESC)
           FROM (SELECT * FROM donation
                  WHERE person_id = p.person_id
                  ORDER BY donated_at DESC, donation_id DESC
                  LIMIT 20) d
       ), '[]'::jsonb);
