-- Write statements for D1/D2. Parameters are bound BY NAME by the harness
-- (the -- params: line lists them in $1..$n order), so designs with different
-- column sets can share one benchmark driver.

-- name: w_insert_donation
-- params: donation_id, person_id, amount_cents, currency, donated_at, note
INSERT INTO donation (donation_id, person_id, amount_cents, currency, donated_at, note)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: w_update_amount
-- params: amount_cents, donation_id
UPDATE donation SET amount_cents = $1 WHERE donation_id = $2;

-- name: w_delete_donation
-- params: donation_id
DELETE FROM donation WHERE donation_id = $1;

-- ---------------------------------------------------------------------------
-- Writes that ORIGINATE at the parent.
--
-- Everything above starts at donation; person and charity are only ever touched
-- as a side effect. But a real system also edits donors directly, and those two
-- operations are strongly design-discriminating -- so leaving them out would
-- systematically hide costs that fall on the embedded and rollup designs.
-- ---------------------------------------------------------------------------

-- name: w_update_person_profile
-- params: email, person_id
-- A donor changes their email. Touches no indexed column and no child data --
-- about as cheap an update as exists. The interest is in what each design has
-- made the person ROW into: in D6 that row also contains the donation history,
-- so this trivial edit has to rewrite (or at least re-version) a row whose size
-- depends on how much the donor has given.
UPDATE person SET email = $1 WHERE person_id = $2;

-- name: w_delete_person_cascade
-- params: person_id
-- Erasure: remove a donor and everything they gave. The logical outcome is
-- identical in every design, which is what makes the comparison fair -- but the
-- work is not. Here it is N child deletes plus one parent delete, and in the
-- rollup designs every one of those child deletes fires a trigger that
-- recomputes the parent's MIN/MAX from the survivors.
WITH removed AS (
    DELETE FROM donation WHERE person_id = $1 RETURNING 1
)
DELETE FROM person WHERE person_id = $1;
