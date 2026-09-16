-- Synchronous upward consolidation of donation facts into person and charity.
--
-- Note on PostgreSQL semantics used below: LEAST() and GREATEST() ignore NULL
-- arguments (unlike most other SQL dialects), so they double as "initialise if
-- currently NULL". Also, inside a single UPDATE every column reference on the
-- right-hand side sees the pre-update value, so the CASE expressions comparing
-- against last_donation_at are reading the old value even though the same
-- statement is assigning a new one.
--
-- INSERT is cheap: pure increments.
-- DELETE and an UPDATE that moves donated_at are expensive: a MAX/MIN cannot be
-- decremented, so the extreme has to be recomputed from the surviving children.
-- That asymmetry is one of the headline findings of this study.

CREATE OR REPLACE FUNCTION donation_rollup_apply(
    p_person_id  BIGINT,
    p_charity_id BIGINT,
    p_recompute  BOOLEAN
) RETURNS VOID AS $$
BEGIN
    IF p_recompute THEN
        UPDATE person p SET
            first_donation_at = agg.min_at,
            last_donation_at  = agg.max_at,
            last_donation_id  = agg.last_id
        FROM (
            SELECT MIN(d.donated_at) AS min_at,
                   MAX(d.donated_at) AS max_at,
                   (SELECT d2.donation_id
                      FROM donation d2
                     WHERE d2.person_id = p_person_id
                     ORDER BY d2.donated_at DESC, d2.donation_id DESC
                     LIMIT 1) AS last_id
              FROM donation d
             WHERE d.person_id = p_person_id
        ) agg
        WHERE p.person_id = p_person_id;

        UPDATE charity c SET
            last_donation_at     = agg.max_at,
            last_donation_id     = agg.last_id,
            last_donor_person_id = agg.last_person
        FROM (
            SELECT MAX(d.donated_at) AS max_at,
                   (SELECT d2.donation_id
                      FROM donation d2
                     WHERE d2.charity_id = p_charity_id
                     ORDER BY d2.donated_at DESC, d2.donation_id DESC
                     LIMIT 1) AS last_id,
                   (SELECT d2.person_id
                      FROM donation d2
                     WHERE d2.charity_id = p_charity_id
                     ORDER BY d2.donated_at DESC, d2.donation_id DESC
                     LIMIT 1) AS last_person
              FROM donation d
             WHERE d.charity_id = p_charity_id
        ) agg
        WHERE c.charity_id = p_charity_id;
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION donation_rollup() RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE person SET
            donation_count      = donation_count + 1,
            total_donated_cents = total_donated_cents + NEW.amount_cents,
            first_donation_at   = LEAST(first_donation_at, NEW.donated_at),
            last_donation_id    = CASE WHEN last_donation_at IS NULL
                                         OR NEW.donated_at >= last_donation_at
                                       THEN NEW.donation_id ELSE last_donation_id END,
            last_donation_at    = GREATEST(last_donation_at, NEW.donated_at)
        WHERE person_id = NEW.person_id;

        UPDATE charity SET
            donation_count       = donation_count + 1,
            total_donated_cents  = total_donated_cents + NEW.amount_cents,
            last_donation_id     = CASE WHEN last_donation_at IS NULL
                                          OR NEW.donated_at >= last_donation_at
                                        THEN NEW.donation_id ELSE last_donation_id END,
            last_donor_person_id = CASE WHEN last_donation_at IS NULL
                                          OR NEW.donated_at >= last_donation_at
                                        THEN NEW.person_id ELSE last_donor_person_id END,
            last_donation_at     = GREATEST(last_donation_at, NEW.donated_at)
        WHERE charity_id = NEW.charity_id;

        RETURN NEW;

    ELSIF TG_OP = 'DELETE' THEN
        UPDATE person SET
            donation_count      = donation_count - 1,
            total_donated_cents = total_donated_cents - OLD.amount_cents
        WHERE person_id = OLD.person_id;

        UPDATE charity SET
            donation_count      = donation_count - 1,
            total_donated_cents = total_donated_cents - OLD.amount_cents
        WHERE charity_id = OLD.charity_id;

        -- The row is already gone from donation by the time an AFTER trigger runs,
        -- so the recompute below sees only survivors.
        PERFORM donation_rollup_apply(OLD.person_id, OLD.charity_id, TRUE);
        RETURN OLD;

    ELSE  -- UPDATE
        UPDATE person SET
            total_donated_cents = total_donated_cents
                                - OLD.amount_cents + NEW.amount_cents
        WHERE person_id = NEW.person_id;

        UPDATE charity SET
            total_donated_cents = total_donated_cents
                                - OLD.amount_cents + NEW.amount_cents
        WHERE charity_id = NEW.charity_id;

        IF NEW.donated_at IS DISTINCT FROM OLD.donated_at THEN
            PERFORM donation_rollup_apply(NEW.person_id, NEW.charity_id, TRUE);
        END IF;
        RETURN NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER donation_rollup_trg
    AFTER INSERT OR UPDATE OR DELETE ON donation
    FOR EACH ROW EXECUTE FUNCTION donation_rollup();
