-- rollup-trigger writes: byte-identical to the flattened baseline's.
--
-- That is deliberate. The trigger fires on the child mutation, so the maintenance
-- cost is created by the schema, not by a different write path. Reading the same
-- statements on both cells means the write-cost difference in the results is the
-- trigger, not the SQL.

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
-- params: donation_id, delta_cents, owner_id
UPDATE donation SET amount_cents = amount_cents + $2
 WHERE donation_id = $1 AND person_id = $3;

-- name: w_donation_delete
-- params: donation_id, owner_id
DELETE FROM donation WHERE donation_id = $1 AND person_id = $2;

-- name: w_person_update
-- params: person_id, full_name, email
-- Touches only untouched columns, so the trigger does not fire and the stored
-- aggregate is not affected. A rollup design that put the aggregate on the same
-- row as frequently edited metadata would pay a very different price; that is a
-- separate design, not this one.
UPDATE person SET full_name = $2, email = $3 WHERE person_id = $1;

-- name: w_donation_reassign
-- params: donation_id, new_person_id, owner_id
UPDATE donation d
   SET person_id  = $2,
       charity_id = (SELECT p.charity_id FROM person p WHERE p.person_id = $2)
 WHERE d.donation_id = $1 AND d.person_id = $3;
