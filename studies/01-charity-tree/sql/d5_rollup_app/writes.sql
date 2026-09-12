-- Write statements for D5 (rollup-app).
--
-- The same work D4's trigger does, but as separate statements the harness issues
-- inside one transaction. Experiment C combines these differently per strategy:
--
--   blind       w_insert_donation, w_rollup_person_ins, w_rollup_charity_ins
--   forupdate   w_lock_person, w_lock_charity, then the three above
--   optimistic  w_read_person_version, w_read_charity_version, insert, then the
--               _cas variants; a zero row count means someone else won, retry.

-- name: w_insert_donation
-- params: donation_id, person_id, charity_id, amount_cents, currency, donated_at, note
INSERT INTO donation (donation_id, person_id, charity_id, amount_cents, currency, donated_at, note)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: w_update_amount
-- params: amount_cents, donation_id
UPDATE donation SET amount_cents = $1 WHERE donation_id = $2;

-- name: w_delete_donation
-- params: donation_id
DELETE FROM donation WHERE donation_id = $1;

-- name: w_rollup_person_ins
-- params: amount_cents, donated_at, donation_id, person_id
UPDATE person SET
    donation_count      = donation_count + 1,
    total_donated_cents = total_donated_cents + $1,
    first_donation_at   = LEAST(first_donation_at, $2::TIMESTAMPTZ),
    last_donation_id    = CASE WHEN last_donation_at IS NULL OR $2::TIMESTAMPTZ >= last_donation_at
                               THEN $3 ELSE last_donation_id END,
    last_donation_at    = GREATEST(last_donation_at, $2::TIMESTAMPTZ),
    version             = version + 1
WHERE person_id = $4;

-- name: w_rollup_charity_ins
-- params: amount_cents, donated_at, donation_id, person_id, charity_id
UPDATE charity SET
    donation_count       = donation_count + 1,
    total_donated_cents  = total_donated_cents + $1,
    last_donation_id     = CASE WHEN last_donation_at IS NULL OR $2::TIMESTAMPTZ >= last_donation_at
                                THEN $3 ELSE last_donation_id END,
    last_donor_person_id = CASE WHEN last_donation_at IS NULL OR $2::TIMESTAMPTZ >= last_donation_at
                                THEN $4 ELSE last_donor_person_id END,
    last_donation_at     = GREATEST(last_donation_at, $2::TIMESTAMPTZ),
    version              = version + 1
WHERE charity_id = $5;

-- name: w_lock_person
-- params: person_id
SELECT version FROM person WHERE person_id = $1 FOR UPDATE;

-- name: w_lock_charity
-- params: charity_id
SELECT version FROM charity WHERE charity_id = $1 FOR UPDATE;

-- name: w_read_person_version
-- params: person_id
SELECT version FROM person WHERE person_id = $1;

-- name: w_read_charity_version
-- params: charity_id
SELECT version FROM charity WHERE charity_id = $1;

-- name: w_rollup_person_ins_cas
-- params: amount_cents, donated_at, donation_id, person_id, version
UPDATE person SET
    donation_count      = donation_count + 1,
    total_donated_cents = total_donated_cents + $1,
    first_donation_at   = LEAST(first_donation_at, $2::TIMESTAMPTZ),
    last_donation_id    = CASE WHEN last_donation_at IS NULL OR $2::TIMESTAMPTZ >= last_donation_at
                               THEN $3 ELSE last_donation_id END,
    last_donation_at    = GREATEST(last_donation_at, $2::TIMESTAMPTZ),
    version             = version + 1
WHERE person_id = $4 AND version = $5;

-- name: w_rollup_charity_ins_cas
-- params: amount_cents, donated_at, donation_id, person_id, charity_id, version
UPDATE charity SET
    donation_count       = donation_count + 1,
    total_donated_cents  = total_donated_cents + $1,
    last_donation_id     = CASE WHEN last_donation_at IS NULL OR $2::TIMESTAMPTZ >= last_donation_at
                                THEN $3 ELSE last_donation_id END,
    last_donor_person_id = CASE WHEN last_donation_at IS NULL OR $2::TIMESTAMPTZ >= last_donation_at
                                THEN $4 ELSE last_donor_person_id END,
    last_donation_at     = GREATEST(last_donation_at, $2::TIMESTAMPTZ),
    version              = version + 1
WHERE charity_id = $5 AND version = $6;

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
