-- Keep the same two parent-row locks as a full rollup, but maintain only sums.
-- The supported mutations are append, amount correction, donation deletion and
-- donor erasure. Moving a donation between donors/charities is outside this study.
CREATE OR REPLACE FUNCTION donation_sum_rollup() RETURNS TRIGGER AS $$
DECLARE
    pid BIGINT;
    cid BIGINT;
    delta BIGINT;
BEGIN
    IF TG_OP = 'INSERT' THEN
        pid := NEW.person_id;
        cid := NEW.charity_id;
        delta := NEW.amount_cents;
    ELSIF TG_OP = 'DELETE' THEN
        pid := OLD.person_id;
        cid := OLD.charity_id;
        delta := -OLD.amount_cents;
    ELSE
        IF NEW.person_id IS DISTINCT FROM OLD.person_id
           OR NEW.charity_id IS DISTINCT FROM OLD.charity_id THEN
            RAISE EXCEPTION 'D16 does not implement donation reassignment';
        END IF;
        pid := NEW.person_id;
        cid := NEW.charity_id;
        delta := NEW.amount_cents - OLD.amount_cents;
    END IF;
    UPDATE person SET total_donated_cents = total_donated_cents + delta
     WHERE person_id = pid;
    UPDATE charity SET total_donated_cents = total_donated_cents + delta
     WHERE charity_id = cid;
    IF TG_OP = 'DELETE' THEN
        RETURN OLD;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER donation_sum_rollup_trg
    AFTER INSERT OR UPDATE OR DELETE ON donation
    FOR EACH ROW EXECUTE FUNCTION donation_sum_rollup();
