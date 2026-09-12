-- D10: the D9 cache trigger, rewritten to stay correct under concurrent writes.
--
-- D9's trigger passed every isolated audit and then silently corrupted 3-5 donor
-- caches per 30-second window once reads and writes ran together
-- (results/20260913-d9-cache-race). The audit's recorded examples showed two
-- distinct bugs, and this file fixes each one with the smallest change that does:
--
-- BUG 1 -- ordering race (7 of 16 examples).
--   D9 PREPENDS the new donation, so the cache is ordered by COMMIT order. The
--   table is ordered by (donated_at, donation_id). Two inserts for the same donor
--   that commit in the opposite order to their timestamps -- adjacent ids such as
--   123936/123937 -- end up swapped.
--   FIX: merge instead of prepend. The new element and the existing elements are
--   re-sorted by (donated_at DESC, donation_id DESC) -- exactly the order the
--   table defines -- before truncating to 20. Commit order stops mattering.
--
-- BUG 2 -- lost update (9 of 16 examples).
--   On correction or removal, D9 REBUILDS the cache from the donation table in a
--   single UPDATE. That statement's snapshot is taken when the statement starts.
--   If a concurrent insert for the same donor commits while the UPDATE waits for
--   the person row lock, PostgreSQL re-checks the target row after the wait but
--   does NOT re-run the subquery against a fresh snapshot -- so the rebuild writes
--   a cache that is missing the donation that just committed.
--   FIX: take the person row lock FIRST, in its own statement, and only then run
--   the rebuild. Under READ COMMITTED each statement gets a new snapshot, so the
--   rebuild now starts after any competing writer on that donor has finished.
--
-- Neither fix touches the read path, so D9 and D10 read identically. What the
-- D9 -> D10 pair measures is the price of correctness on the write path.
--
-- Why the INSERT path needs no explicit lock: its SET expression reads
-- person.recent_donations from the row being updated. When that UPDATE waits on a
-- concurrent writer, PostgreSQL re-evaluates the expression against the newer row
-- version, so the merge always starts from the latest committed cache. Only the
-- rebuild, which reads a DIFFERENT table, is exposed to the stale snapshot.

CREATE OR REPLACE FUNCTION donation_recent_cache_rebuild(p_person_id BIGINT)
RETURNS VOID AS $$
BEGIN
    -- Lock first, so the rebuild below starts with a snapshot taken AFTER any
    -- concurrent writer on this donor has committed.
    PERFORM 1 FROM person WHERE person_id = p_person_id FOR UPDATE;

    UPDATE person SET recent_donations = COALESCE((
        SELECT jsonb_agg(jsonb_build_object(
                   'i', d.donation_id, 'a', d.amount_cents, 'c', d.currency,
                   't', d.donated_at,  'n', d.note)
                 ORDER BY d.donated_at DESC, d.donation_id DESC)
          FROM (SELECT * FROM donation
                 WHERE person_id = p_person_id
                 ORDER BY donated_at DESC, donation_id DESC
                 LIMIT 20) d
    ), '[]'::jsonb)
    WHERE person_id = p_person_id;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION donation_recent_cache() RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        -- Merge, then sort by the table's own order, then truncate. The sort keys
        -- match the audit exactly: donated_at DESC, donation_id DESC.
        UPDATE person SET recent_donations = (
            SELECT COALESCE(jsonb_agg(x.e ORDER BY x.at DESC, x.id DESC), '[]'::jsonb)
              FROM (
                    SELECT e, at, id
                      FROM (
                            SELECT jsonb_build_object(
                                       'i', NEW.donation_id, 'a', NEW.amount_cents,
                                       'c', NEW.currency,    't', NEW.donated_at,
                                       'n', NEW.note)             AS e,
                                   NEW.donated_at                 AS at,
                                   NEW.donation_id                AS id
                            UNION ALL
                            SELECT t.e,
                                   (t.e ->> 't')::TIMESTAMPTZ,
                                   (t.e ->> 'i')::BIGINT
                              FROM jsonb_array_elements(person.recent_donations) AS t(e)
                           ) merged
                     ORDER BY at DESC, id DESC
                     LIMIT 20
                   ) x
        )
        WHERE person_id = NEW.person_id;
        RETURN NEW;

    ELSIF TG_OP = 'DELETE' THEN
        PERFORM donation_recent_cache_rebuild(OLD.person_id);
        RETURN OLD;

    ELSE
        PERFORM donation_recent_cache_rebuild(NEW.person_id);
        RETURN NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER donation_recent_cache_trg
    AFTER INSERT OR UPDATE OR DELETE ON donation
    FOR EACH ROW EXECUTE FUNCTION donation_recent_cache();
