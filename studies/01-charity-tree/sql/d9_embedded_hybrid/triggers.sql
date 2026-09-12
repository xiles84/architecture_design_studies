-- Maintains person.recent_donations as the 20 most recent donations, newest first.
--
-- The trick that bounds the cost: the new element is given ordinal 0 and the
-- existing cached elements keep ordinals 1..n, so "keep ordinal < 20" prepends
-- the new donation and drops the oldest in a single pass. The array never grows
-- past 20 elements, which is the entire difference from D6 -- the write cost is
-- a constant, not a function of how much the donor has given before.
--
-- This assumes donations arrive in roughly chronological order, which is true of
-- an append-only stream. A backdated insert would land at the front of the cache
-- where it does not belong; the full table underneath still has the truth, and a
-- system that backdates would rebuild the cache from it rather than prepend.
--
-- UPDATE and DELETE rebuild the slice from the donation table instead of trying
-- to patch the array. They are rare next to inserts, and the honest recompute is
-- both simpler and impossible to get subtly wrong.

CREATE OR REPLACE FUNCTION donation_recent_cache_rebuild(p_person_id BIGINT)
RETURNS VOID AS $$
BEGIN
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
        UPDATE person SET recent_donations = (
            SELECT COALESCE(jsonb_agg(x.e ORDER BY x.ord), '[]'::jsonb)
              FROM (
                    SELECT jsonb_build_object(
                               'i', NEW.donation_id, 'a', NEW.amount_cents,
                               'c', NEW.currency,    't', NEW.donated_at,
                               'n', NEW.note) AS e,
                           0 AS ord
                    UNION ALL
                    SELECT t.e, t.ord
                      FROM jsonb_array_elements(person.recent_donations)
                           WITH ORDINALITY AS t(e, ord)
                   ) x
             WHERE x.ord < 20
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
