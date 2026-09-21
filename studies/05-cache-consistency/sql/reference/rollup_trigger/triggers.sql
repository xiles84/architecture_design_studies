-- rollup-trigger: the trigger and the backfill.
--
-- Attached AFTER the bulk load (methodology, and study 01's convention): a
-- per-row trigger during the load would measure the loader, not the design.
--
-- The trigger RECOMPUTES the affected parent from the donation table rather than
-- incrementing a counter. That is deliberately the honest version: an
-- increment/decrement trigger is cheaper and is the one that silently drifts when
-- a write path forgets it, and this study compares read speed against maintenance
-- cost, not against a strawman. `a_rollup_drift` is what checks it.
--
-- The rebuild locks nothing explicitly. Under READ COMMITTED the UPDATE of the
-- parent takes its row lock; whether that is enough to keep the stored aggregate
-- correct under concurrent writes is exactly what the correctness gate and the
-- drift audit are for -- and study 01's D9 learned that lesson the expensive way.

CREATE OR REPLACE FUNCTION portal_rollup_rebuild(p_person_id BIGINT)
RETURNS VOID AS $$
BEGIN
    UPDATE person p
       SET donation_count       = a.n,
           donation_total_cents = a.total
      FROM (
            SELECT count(*)                       AS n,
                   COALESCE(sum(amount_cents), 0) AS total
              FROM donation
             WHERE person_id = p_person_id
           ) a
     WHERE p.person_id = p_person_id;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION portal_rollup() RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        PERFORM portal_rollup_rebuild(NEW.person_id);
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        PERFORM portal_rollup_rebuild(OLD.person_id);
        RETURN OLD;
    ELSE
        PERFORM portal_rollup_rebuild(NEW.person_id);
        -- Reassignment moves a donation between two parents, so both aggregates
        -- are now stale. Recomputing only the destination would leave the source
        -- claiming a donation it no longer has.
        IF NEW.person_id IS DISTINCT FROM OLD.person_id THEN
            PERFORM portal_rollup_rebuild(OLD.person_id);
        END IF;
        RETURN NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER portal_rollup_trg
    AFTER INSERT OR UPDATE OR DELETE ON donation
    FOR EACH ROW EXECUTE FUNCTION portal_rollup();

-- Backfill: every loaded parent gets its stored aggregate before any measurement.
UPDATE person p
   SET donation_count       = a.n,
       donation_total_cents = a.total
  FROM (
        SELECT person_id, count(*) AS n, COALESCE(sum(amount_cents), 0) AS total
          FROM donation
         GROUP BY person_id
       ) a
 WHERE p.person_id = a.person_id;

UPDATE person
   SET donation_count = 0, donation_total_cents = 0
 WHERE person_id NOT IN (SELECT DISTINCT person_id FROM donation);
