-- Legacy write catalogue.
--
-- Five logical mutations, because the study's question spans them:
--   insert, correction, deletion, person metadata, reassignment between people.
-- Reassignment is the one that touches TWO cache keys, and it is the reason the
-- owned model needs ordered locking.
--
-- Nothing here bumps a version or writes an outbox row: this model cannot be
-- changed for caching. That is the point of the pair.

-- name: w_load_donations
-- params: none
-- The only load statement. The staging table is the same for every design, so the
-- loader's cost does not favour one.
INSERT INTO donation (donation_id, person_id, charity_id, amount_cents, currency, donated_at, note)
SELECT donation_id, person_id, charity_id, amount_cents, currency, donated_at, note
  FROM load_donation;

-- name: w_donation_insert
-- params: donation_id, person_id, charity_id, amount_cents, currency, donated_at, note
INSERT INTO donation (donation_id, person_id, charity_id, amount_cents, currency, donated_at, note)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: w_donation_correct
-- params: donation_id, amount_cents
-- An amount correction is an UPDATE of the child row in place. It changes the
-- donated total and possibly the recent-20 slice's content, so every cache
-- strategy must invalidate or republish the donor's entry.
UPDATE donation SET amount_cents = $2 WHERE donation_id = $1;

-- name: w_donation_delete
-- params: donation_id
DELETE FROM donation WHERE donation_id = $1;

-- name: w_person_update
-- params: person_id, full_name, email
-- Touches no child row and no indexed column, but it changes the portal payload,
-- so a cache that only reacts to donation writes is wrong here.
UPDATE person SET full_name = $2, email = $3 WHERE person_id = $1;

-- name: w_donation_reassign
-- params: donation_id, new_person_id
-- Moving a donation between people changes TWO portal views, and on the legacy
-- model the copied charity_id must move with it or the rolldown invariant breaks.
-- `new_person_id` is validated by the harness against the dataset, and the new
-- charity is taken from the destination person so the copy cannot drift.
UPDATE donation d
   SET person_id  = $2,
       charity_id = (SELECT p.charity_id FROM person p WHERE p.person_id = $2)
 WHERE d.donation_id = $1;
