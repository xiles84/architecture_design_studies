-- rollup-trigger audits: the baseline audits plus the invariant a stored aggregate
-- creates. a_rollup_drift is the whole reason storing an aggregate is a design
-- decision rather than an optimisation: a design that is faster because it keeps a
-- wrong number has not won.

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

-- name: a_rollup_drift
-- params: none
-- Counts parents whose two stored aggregates disagree with the child table they
-- summarise, in either direction: a count that is too low and a count that is too
-- high are both drift, and a design that only ever under-counts would pass a
-- one-sided check.
SELECT count(*)
  FROM person p
  LEFT JOIN (
        SELECT person_id, count(*) AS n, COALESCE(sum(amount_cents), 0) AS total
          FROM donation
         GROUP BY person_id
       ) a ON a.person_id = p.person_id
 WHERE p.donation_count       <> COALESCE(a.n, 0)
    OR p.donation_total_cents <> COALESCE(a.total, 0);
