-- Maintenance for donation.is_last_donation (RECENCY.md section 3).
--
-- The invariant is "exactly one TRUE row per donor, and it is their newest".
-- A MIN/MAX-shaped fact can be pushed forward for free (a new donation is
-- trivially newer than nothing) but not decremented: deleting the flagged
-- row means recomputing from the survivors, exactly the asymmetry D4's
-- rollup trigger already has for its own MIN/MAX columns.
--
-- Concurrency: two inserts for the SAME donor, racing at READ COMMITTED,
-- must not both conclude "I am the newest" and both set the flag. The lock
-- on the donor's own person row is what serialises them -- it is not a lock
-- study 01 needs for any other reason, borrowed here because donation has no
-- natural per-donor row to lock and person is guaranteed to exist first (the
-- FK). d21_recency_flag_unguarded is this file with step 1 removed, and is
-- expected to leave two flagged rows for one donor under concurrent inserts
-- (the negative control this design needs, methodology 5a).

CREATE OR REPLACE FUNCTION donation_last_flag_biu() RETURNS TRIGGER AS $$
BEGIN
    -- 2. Clear the flag on this donor's currently-flagged row, but only if
    --    the arriving donation is not older than it -- a backdated insert
    --    (RECENCY.md section 3, the maintenance write op) must not steal the
    --    flag from a donation that is genuinely newer.
    UPDATE donation
       SET is_last_donation = false
     WHERE person_id = NEW.person_id
       AND is_last_donation
       AND donated_at <= NEW.donated_at;

    -- 3. The arriving row is flagged only if nothing newer already exists.
    NEW.is_last_donation := NOT EXISTS (
        SELECT 1 FROM donation
         WHERE person_id = NEW.person_id
           AND donated_at > NEW.donated_at
    );

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER donation_last_flag_biu_trg
    BEFORE INSERT ON donation
    FOR EACH ROW
    EXECUTE FUNCTION donation_last_flag_biu();

-- A MAX cannot be decremented: if the deleted row carried the flag, promote
-- the donor's next-newest survivor. Ties broken by donation_id, matching the
-- gate's tie handling (RECENCY.md section 5).
CREATE OR REPLACE FUNCTION donation_last_flag_ad() RETURNS TRIGGER AS $$
DECLARE
    next_id BIGINT;
BEGIN
    IF NOT OLD.is_last_donation THEN
        RETURN OLD;
    END IF;

    PERFORM 1 FROM person WHERE person_id = OLD.person_id FOR UPDATE;

    SELECT donation_id INTO next_id
      FROM donation
     WHERE person_id = OLD.person_id
     ORDER BY donated_at DESC, donation_id DESC
     LIMIT 1;

    IF next_id IS NOT NULL THEN
        UPDATE donation SET is_last_donation = true WHERE donation_id = next_id;
    END IF;

    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER donation_last_flag_ad_trg
    AFTER DELETE ON donation
    FOR EACH ROW
    EXECUTE FUNCTION donation_last_flag_ad();

-- Deliberately no trigger on UPDATE of amount_cents: donated_at does not
-- move, so the flag does not move. If one fires here, the write benchmark
-- will show unexpected cost on w_update_amount and that is a harness bug.
