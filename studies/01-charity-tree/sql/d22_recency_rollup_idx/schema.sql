-- D4 "rollup-trigger": D3 plus pre-computed aggregates on the parent rows,
-- maintained synchronously by database triggers.
--
-- Every question in the catalogue that asks "how much / how many / when last"
-- becomes a single-row primary-key read. The aggregate is consolidated upward
-- into the parent instead of being derived from the children on every read.
--
-- The bill arrives on the write path: one INSERT into donation fans out into
-- two UPDATEs of parent rows. The charity row in particular is a hot row --
-- every donation in the entire charity serialises on it. That contention is
-- not a bug in the design, it IS the design's cost, and experiment C measures it.

CREATE TABLE charity (
    charity_id           BIGINT      PRIMARY KEY,
    name                 TEXT        NOT NULL,
    country              TEXT        NOT NULL,
    founded_on           DATE        NOT NULL,
    -- consolidated from donation
    donation_count       BIGINT      NOT NULL DEFAULT 0,
    total_donated_cents  BIGINT      NOT NULL DEFAULT 0,
    last_donation_id     BIGINT,
    last_donation_at     TIMESTAMPTZ,
    last_donor_person_id BIGINT
);

CREATE TABLE person (
    person_id            BIGINT      PRIMARY KEY,
    charity_id           BIGINT      NOT NULL REFERENCES charity (charity_id),
    full_name            TEXT        NOT NULL,
    email                TEXT        NOT NULL,
    joined_at            TIMESTAMPTZ NOT NULL,
    -- consolidated from donation
    donation_count       BIGINT      NOT NULL DEFAULT 0,
    total_donated_cents  BIGINT      NOT NULL DEFAULT 0,
    first_donation_at    TIMESTAMPTZ,
    last_donation_at     TIMESTAMPTZ,
    last_donation_id     BIGINT
);

CREATE TABLE donation (
    donation_id  BIGINT      PRIMARY KEY,
    person_id    BIGINT      NOT NULL REFERENCES person (person_id),
    charity_id   BIGINT      NOT NULL REFERENCES charity (charity_id),
    amount_cents BIGINT      NOT NULL,
    currency     CHAR(3)     NOT NULL,
    donated_at   TIMESTAMPTZ NOT NULL,
    note         TEXT
);
