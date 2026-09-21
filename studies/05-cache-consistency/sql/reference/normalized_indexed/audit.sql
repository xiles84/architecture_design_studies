-- normalized-indexed audits.
--
-- a_rolldown_drift does not exist here: there is no copied key to drift. Its
-- absence is the point -- the flattened cell pays for the invariant the normalized
-- cell does not have.

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
