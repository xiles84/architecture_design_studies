-- embedded-locked writes: byte-identical to the flattened baseline's.
--
-- The trigger fires on the child mutation, so the slice-maintenance cost is created
-- by the schema, not by a different write path. Reading the same statements on both
-- cells means the write-cost difference in the results is the embedding, not the
-- SQL.
--
-- The insert in the workload is a NEW donation with a fresh, monotonically
-- increasing id and a donated_at at "now", which is the shape a live portal
-- produces. A backdated insert would land in the merge and be sorted into place by
-- BUG 1's fix; the audits accept that, and the workload does not manufacture it.

-- name: w_load_donations
-- params: none
INSERT INTO donation (donation_id, person_id, charity_id, amount_cents, currency, donated_at, note)
SELECT donation_id, person_id, charity_id, amount_cents, currency, donated_at, note
  FROM load_donation;

-- name: w_donation_insert
-- params: donation_id, person_id, charity_id, amount_cents, currency, donated_at, note
INSERT INTO donation (donation_id, person_id, charity_id, amount_cents, currency, donated_at, note)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: w_donation_correct
-- params: donation_id, amount_cents
UPDATE donation SET amount_cents = $2 WHERE donation_id = $1;

-- name: w_donation_delete
-- params: donation_id
DELETE FROM donation WHERE donation_id = $1;

-- name: w_person_update
-- params: person_id, full_name, email
-- The embedded slice lives on this row, and this is the operation that makes
-- "embedded" interesting for a relational engine: an edit that touches nothing the
-- slice contains still rewrites (or re-versions) the row that contains it. Study 01
-- found PostgreSQL does not rewrite an unchanged TOASTed column, so a heavy donor's
-- out-of-line array survives while a median donor's ~2 kB array is re-versioned in
-- full. That behaviour is visible here, not assumed.
UPDATE person SET full_name = $2, email = $3 WHERE person_id = $1;

-- name: w_donation_reassign
-- params: donation_id, new_person_id
UPDATE donation d
   SET person_id  = $2,
       charity_id = (SELECT p.charity_id FROM person p WHERE p.person_id = $2)
 WHERE d.donation_id = $1;
