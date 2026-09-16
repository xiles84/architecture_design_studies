-- D5 "rollup-app": identical column layout to D4, but NO triggers.
--
-- The application maintains the same consolidated aggregates inside its own
-- transaction. Separating this from D4 isolates one question: given the same
-- amount of work, does it matter WHERE the rollup logic lives?
--
-- Holding the schema constant also lets experiment C swap concurrency-control
-- strategies (blind UPDATE / SELECT FOR UPDATE / optimistic version check)
-- without changing anything a query planner can see. The `version` columns exist
-- only for the optimistic strategy; the other strategies ignore them.

CREATE TABLE charity (
    charity_id           BIGINT      PRIMARY KEY,
    name                 TEXT        NOT NULL,
    country              TEXT        NOT NULL,
    founded_on           DATE        NOT NULL,
    donation_count       BIGINT      NOT NULL DEFAULT 0,
    total_donated_cents  BIGINT      NOT NULL DEFAULT 0,
    last_donation_id     BIGINT,
    last_donation_at     TIMESTAMPTZ,
    last_donor_person_id BIGINT,
    version              BIGINT      NOT NULL DEFAULT 0
);

CREATE TABLE person (
    person_id            BIGINT      PRIMARY KEY,
    charity_id           BIGINT      NOT NULL REFERENCES charity (charity_id),
    full_name            TEXT        NOT NULL,
    email                TEXT        NOT NULL,
    joined_at            TIMESTAMPTZ NOT NULL,
    donation_count       BIGINT      NOT NULL DEFAULT 0,
    total_donated_cents  BIGINT      NOT NULL DEFAULT 0,
    first_donation_at    TIMESTAMPTZ,
    last_donation_at     TIMESTAMPTZ,
    last_donation_id     BIGINT,
    version              BIGINT      NOT NULL DEFAULT 0
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
