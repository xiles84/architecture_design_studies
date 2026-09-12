-- D1 "normalized-minimal": textbook 3NF. Primary keys and foreign keys only.
-- This is the design you get by drawing the ER diagram and stopping there.
-- PostgreSQL does NOT create an index for a foreign key's referencing column,
-- so donation.person_id is unindexed here. That is the point of the design.

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
    amount_cents BIGINT      NOT NULL,
    currency     CHAR(3)     NOT NULL,
    donated_at   TIMESTAMPTZ NOT NULL,
    note         TEXT
);
