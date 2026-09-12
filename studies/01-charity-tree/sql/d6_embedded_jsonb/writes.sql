-- Write statements for D6 (embedded-jsonb).
--
-- Every one of these rewrites the whole person row. Under MVCC that allocates a
-- new row version the size of the entire donation history, so the cost of
-- recording one donation grows with how many that person has already made.

-- name: w_insert_donation
-- params: person_id, donation_id, amount_cents, currency, donated_at, note
UPDATE person
   SET donations = donations || jsonb_build_array(
           jsonb_build_object('i', $2::BIGINT,
                              'a', $3::BIGINT,
                              'c', $4::TEXT,
                              't', $5::TIMESTAMPTZ,
                              'n', $6::TEXT))
 WHERE person_id = $1;

-- name: w_update_amount
-- params: amount_cents, donation_id
-- No way to address one element in place: the array is rebuilt element by element.
UPDATE person p
   SET donations = (
        SELECT COALESCE(jsonb_agg(
                 CASE WHEN (t.e ->> 'i')::BIGINT = $2::BIGINT
                      THEN jsonb_set(t.e, '{a}', to_jsonb($1::BIGINT))
                      ELSE t.e END
                 ORDER BY t.ord), '[]'::jsonb)
          FROM jsonb_array_elements(p.donations) WITH ORDINALITY AS t(e, ord))
 WHERE p.donations @> jsonb_build_array(jsonb_build_object('i', $2::BIGINT));

-- name: w_delete_donation
-- params: donation_id
UPDATE person p
   SET donations = (
        SELECT COALESCE(jsonb_agg(t.e ORDER BY t.ord), '[]'::jsonb)
          FROM jsonb_array_elements(p.donations) WITH ORDINALITY AS t(e, ord)
         WHERE (t.e ->> 'i')::BIGINT <> $1::BIGINT)
 WHERE p.donations @> jsonb_build_array(jsonb_build_object('i', $1::BIGINT));

-- ---------------------------------------------------------------------------
-- Writes that ORIGINATE at the parent.
-- ---------------------------------------------------------------------------

-- name: w_update_person_profile
-- params: email, person_id
-- The same trivial edit as everywhere else -- but in this design the person row
-- also holds the entire donation history. Whether that matters is genuinely
-- uncertain and worth measuring rather than reasoning about: PostgreSQL does not
-- rewrite an unchanged TOASTed column, so a heavy donor's out-of-line array
-- survives untouched, while a median donor's ~2 kB array sits inline in the
-- tuple and is re-versioned in full. The GIN index on donations also makes a
-- HOT update less likely once pages fill up.
UPDATE person SET email = $1 WHERE person_id = $2;

-- name: w_delete_person_cascade
-- params: person_id
-- Erasure is where embedding should shine: the children are inside the parent,
-- so removing a donor and their entire history is a single row delete with no
-- child table to visit and no triggers to fire.
DELETE FROM person WHERE person_id = $1;
