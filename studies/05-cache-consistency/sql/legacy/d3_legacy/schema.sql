-- Study 05 — LEGACY database model, and the flattened-FK reference cell.
--
-- This is Study 01's D3 "flattened-fk" schema, kept exactly as it is. It is the
-- study's principal source of truth, AND it is the "existing legacy database
-- model that cannot be modified for caching":
--
--   * no cache-specific version column,
--   * no cache outbox,
--   * no new trigger,
--   * no CDC configuration,
--   * no reliable cache-specific updated_at token.
--
-- The cache adapter may keep external metadata (the harness's oracle generation
-- counter is measurement metadata, not a schema change), but it must never
-- pretend that cache-side metadata can detect an unknown database writer. That
-- is why a STRICT legacy cache under external writers has to validate through
-- the authoritative database or bypass the cache entirely, and why the result of
-- such a cell says which of the two it did.
--
-- The denormalised donation.charity_id is D3's rolldown: the parent's key copied
-- onto the child so a charity-scoped question touches one table instead of two.
-- The reference pair that isolates it is reference/normalized_indexed, which is
-- the same schema without that column and with the join put back into the SQL.

CREATE TABLE charity (
    charity_id  BIGINT      PRIMARY KEY,
    name        TEXT        NOT NULL,
    country     TEXT        NOT NULL,
    founded_on  DATE        NOT NULL
);

CREATE TABLE person (
    person_id   BIGINT      PRIMARY KEY,
    charity_id  BIGINT      NOT NULL REFERENCES charity (charity_id),
    full_name   TEXT        NOT NULL,
    email       TEXT        NOT NULL,
    joined_at   TIMESTAMPTZ NOT NULL
);

CREATE TABLE donation (
    donation_id  BIGINT      PRIMARY KEY,
    person_id    BIGINT      NOT NULL REFERENCES person (person_id),
    charity_id   BIGINT      NOT NULL REFERENCES charity (charity_id),  -- rolldown
    amount_cents BIGINT      NOT NULL,
    currency     CHAR(3)     NOT NULL,
    donated_at   TIMESTAMPTZ NOT NULL,
    note         TEXT
);

-- Load staging, shared by every design. The dataset is written here with COPY and
-- then moved into `donation` by one design-specific statement, so no design pays a
-- per-row cost the others avoid (methodology 9/10: the load is timed but is never
-- the thing a study compares).
CREATE TABLE load_donation (
    donation_id  BIGINT      NOT NULL,
    person_id    BIGINT      NOT NULL,
    charity_id   BIGINT      NOT NULL,
    amount_cents BIGINT      NOT NULL,
    currency     CHAR(3)     NOT NULL,
    donated_at   TIMESTAMPTZ NOT NULL,
    note         TEXT
);
