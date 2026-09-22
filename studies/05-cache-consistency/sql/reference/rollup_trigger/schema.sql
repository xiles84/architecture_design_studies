-- Reference cell: rollup-trigger. The rollup half of the information-placement
-- comparison.
--
-- Two aggregates the portal shows -- how many donations this donor made and how
-- much they gave -- stop being derived at read time and become stored columns on
-- the parent, maintained by a trigger on the child. One decision, isolated against
-- the flattened baseline: everything else (schema, indexes, the recent-20 read,
-- every mutation) is unchanged.
--
-- Study 01's D4 is this shape. What it buys and what it costs are both measured
-- here: the portal read stops aggregating, and every write starts recomputing a
-- parent's rollup inside the same transaction. Storing the aggregate also creates
-- a value that can drift, which is why a_rollup_drift exists.

CREATE TABLE charity (
    charity_id  BIGINT      PRIMARY KEY,
    name        TEXT        NOT NULL,
    country     TEXT        NOT NULL,
    founded_on  DATE        NOT NULL
);

CREATE TABLE person (
    person_id            BIGINT      PRIMARY KEY,
    charity_id           BIGINT      NOT NULL REFERENCES charity (charity_id),
    full_name            TEXT        NOT NULL,
    email                TEXT        NOT NULL,
    joined_at            TIMESTAMPTZ NOT NULL,
    donation_count       BIGINT      NOT NULL DEFAULT 0,
    donation_total_cents BIGINT      NOT NULL DEFAULT 0
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

CREATE TABLE load_donation (
    donation_id  BIGINT      NOT NULL,
    person_id    BIGINT      NOT NULL,
    charity_id   BIGINT      NOT NULL,
    amount_cents BIGINT      NOT NULL,
    currency     CHAR(3)     NOT NULL,
    donated_at   TIMESTAMPTZ NOT NULL,
    note         TEXT
);
