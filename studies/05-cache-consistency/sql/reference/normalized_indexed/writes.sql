-- normalized-indexed writes. The same five logical mutations as every design; the
-- only difference is that no copied charity_id has to be written or kept correct,
-- so an insert is narrower and reassignment is a one-column update.

-- name: w_load_donations
-- params: none
INSERT INTO donation (donation_id, person_id, amount_cents, currency, donated_at, note)
SELECT donation_id, person_id, amount_cents, currency, donated_at, note
  FROM load_donation;

-- name: w_donation_insert
-- params: donation_id, person_id, amount_cents, currency, donated_at, note
INSERT INTO donation (donation_id, person_id, amount_cents, currency, donated_at, note)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: w_donation_correct
-- params: donation_id, delta_cents, owner_id
UPDATE donation SET amount_cents = amount_cents + $2
 WHERE donation_id = $1 AND person_id = $3;

-- name: w_donation_delete
-- params: donation_id, owner_id
DELETE FROM donation WHERE donation_id = $1 AND person_id = $2;

-- name: w_person_update
-- params: person_id, full_name, email
UPDATE person SET full_name = $2, email = $3 WHERE person_id = $1;

-- name: w_donation_reassign
-- params: donation_id, new_person_id, owner_id
-- No copied key to keep in step: this is exactly the work the flattened model
-- makes the database do, and the reason reassignment appears in both cells.
UPDATE donation SET person_id = $2 WHERE donation_id = $1 AND person_id = $3;
