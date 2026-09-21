-- embedded-locked: the concurrency-correct bounded-embedding trigger.
--
-- This is Study 01's D10 shape, restated for this study's element names. Both D9
-- bugs are fixed the same way D10 fixed them, and the fixes are the reason this
-- design can be compared to the others at all:
--
-- BUG 1 -- ordering race. Prepending the new donation orders the stored slice by
--   COMMIT order, not by the table's (donated_at DESC, donation_id DESC) order. Two
--   inserts for one donor that commit in the opposite order to their timestamps end
--   up swapped.
--   FIX: merge the new element with the existing ones, sort by the table's own
--   order, then truncate to 20. Commit order stops mattering.
--
-- BUG 2 -- lost update. Rebuilding the slice from the donation table in one UPDATE
--   takes the statement's snapshot when the statement starts. If a concurrent insert
--   for the same donor commits while the UPDATE waits for the parent row lock,
--   PostgreSQL re-checks the target row after the wait but does NOT re-run the
--   subquery against a fresh snapshot, so the rebuild stores a slice missing the
--   donation that just committed.
--   FIX: take the parent row lock FIRST, in its own statement, then rebuild. Under
--   READ COMMITTED each statement gets a new snapshot, so the rebuild starts after
--   any competing writer on that donor has committed.
--
-- The insert path needs no explicit lock: its SET expression reads
-- person.recent_donations from the row being updated, so a wait re-evaluates the
-- expression against the newer row version.
--
-- Attached AFTER the bulk load, with a backfill, so the loader does not pay a
-- per-row rebuild.

CREATE OR REPLACE FUNCTION portal_recent_rebuild(p_person_id BIGINT)
RETURNS VOID AS $$
BEGIN
    -- Lock first (BUG 2's fix): the rebuild below must start with a snapshot taken
    -- after any concurrent writer on this donor has committed.
    PERFORM 1 FROM person WHERE person_id = p_person_id FOR UPDATE;

    UPDATE person
       SET recent_donations = COALESCE((
            SELECT jsonb_agg(jsonb_build_object(
                       'id',   d.donation_id,
                       'amt',  d.amount_cents,
                       'cur',  d.currency,
                       'at',   d.donated_at,
                       'note', d.note)
                     ORDER BY d.donated_at DESC, d.donation_id DESC)
              FROM (SELECT * FROM donation
                     WHERE person_id = p_person_id
                     ORDER BY donated_at DESC, donation_id DESC
                     LIMIT 20) d
           ), '[]'::jsonb)
     WHERE person_id = p_person_id;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION portal_recent() RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        -- Merge, sort by the table's own order, truncate (BUG 1's fix).
        UPDATE person
           SET recent_donations = (
                SELECT COALESCE(jsonb_agg(x.e ORDER BY x.at DESC, x.id DESC), '[]'::jsonb)
                  FROM (
                        SELECT e, at, id
                          FROM (
                                SELECT jsonb_build_object(
                                           'id',   NEW.donation_id,
                                           'amt',  NEW.amount_cents,
                                           'cur',  NEW.currency,
                                           'at',   NEW.donated_at,
                                           'note', NEW.note)  AS e,
                                       NEW.donated_at        AS at,
                                       NEW.donation_id       AS id
                                UNION ALL
                                SELECT t.e,
                                       (t.e ->> 'at')::TIMESTAMPTZ,
                                       (t.e ->> 'id')::BIGINT
                                  FROM jsonb_array_elements(person.recent_donations) AS t(e)
                               ) merged
                         ORDER BY at DESC, id DESC
                         LIMIT 20
                       ) x
           )
         WHERE person_id = NEW.person_id;
        RETURN NEW;

    ELSIF TG_OP = 'DELETE' THEN
        PERFORM portal_recent_rebuild(OLD.person_id);
        RETURN OLD;

    ELSE
        PERFORM portal_recent_rebuild(NEW.person_id);
        -- Reassignment also changes the SOURCE donor's slice.
        IF NEW.person_id IS DISTINCT FROM OLD.person_id THEN
            PERFORM portal_recent_rebuild(OLD.person_id);
        END IF;
        RETURN NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER portal_recent_trg
    AFTER INSERT OR UPDATE OR DELETE ON donation
    FOR EACH ROW EXECUTE FUNCTION portal_recent();

-- Backfill the slice for every loaded donor.
UPDATE person p
   SET recent_donations = COALESCE((
        SELECT jsonb_agg(jsonb_build_object(
                   'id',   d.donation_id,
                   'amt',  d.amount_cents,
                   'cur',  d.currency,
                   'at',   d.donated_at,
                   'note', d.note)
                 ORDER BY d.donated_at DESC, d.donation_id DESC)
          FROM (SELECT * FROM donation
                 WHERE person_id = p.person_id
                 ORDER BY donated_at DESC, donation_id DESC
                 LIMIT 20) d
       ), '[]'::jsonb);
