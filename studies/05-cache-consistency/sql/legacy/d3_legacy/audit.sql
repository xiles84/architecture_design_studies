-- Legacy audit catalogue.
--
-- The primary reconciliation is done in Go: the oracle recomputes what each key's
-- committed state must be from the operation ledger, and compares it with what the
-- read catalogue returns. These statements are the database-side cross-checks that
-- make the comparison independent of the read path -- a design whose reads and
-- whose audit share one bug would otherwise agree with itself, which is exactly
-- the failure methodology 5 exists to prevent.

-- name: a_person_aggregate
-- params: person_id
-- The count and total recomputed from the donation table, not from any stored
-- rollup. A design that maintains a rollup is checked against this.
SELECT count(*)                                  AS donation_count,
       COALESCE(sum(amount_cents), 0)            AS donation_total_cents
  FROM donation
 WHERE person_id = $1;

-- name: a_donation_ids
-- params: person_id
-- Every donation id the database holds for this person. The Go audit compares the
-- exact id set against the ledger, which is the source-write invariant: a write the
-- client saw acknowledged must be present, and a write it saw rejected must not be.
SELECT donation_id
  FROM donation
 WHERE person_id = $1
 ORDER BY donation_id;

-- name: a_recent_ids
-- params: person_id
-- The newest 20 ids in the study's deterministic order, computed independently of
-- the read catalogue. A cached slice that dropped or reordered an element shows up
-- here even when the aggregate looked right.
SELECT donation_id
  FROM donation
 WHERE person_id = $1
 ORDER BY donated_at DESC, donation_id DESC
 LIMIT 20;

-- name: a_rolldown_drift
-- params: none
-- Counts children whose copied charity_id disagrees with their parent's. This is
-- the invariant that D3's rolldown buys read speed with, and reassignment is the
-- mutation that can break it. It must be zero in every legacy/flattened cell.
SELECT count(*)
  FROM donation d
  JOIN person p ON p.person_id = d.person_id
 WHERE d.charity_id <> p.charity_id;
