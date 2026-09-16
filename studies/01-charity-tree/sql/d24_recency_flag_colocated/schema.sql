-- D24 "recency-flag-colocated": D20's flag, D7's placement. YugabyteDB only.
--
-- The primary-key rebuild is exactly D7's (see d7_yb_child_colocated/schema.sql
-- for the full reasoning): donation is keyed on ((person_id) HASH, donated_at
-- DESC, donation_id ASC), so one donor's donations share a tablet instead of
-- scattering across the cluster.
--
-- What is unchanged from D20: is_last_donation, maintained by the same trigger
-- (triggers.sql, byte-identical to D20's) and read by the same queries
-- (queries.sql, byte-identical to D20's). D20 -> D24 isolates data placement
-- alone (RECENCY.md section 4): whatever the flag's read/write cost is under
-- colocation, it is not explained by a changed statement.

CREATE TABLE charity (
    charity_id  BIGINT      PRIMARY KEY,
    name        TEXT        NOT NULL,
    country     TEXT        NOT NULL,
    founded_on  DATE        NOT NULL
);

CREATE TABLE person (
    person_id   BIGINT,
    charity_id  BIGINT      NOT NULL REFERENCES charity (charity_id),
    full_name   TEXT        NOT NULL,
    email       TEXT        NOT NULL,
    joined_at   TIMESTAMPTZ NOT NULL,
    PRIMARY KEY ((person_id) HASH)
);

CREATE TABLE donation (
    donation_id       BIGINT      NOT NULL,
    person_id         BIGINT      NOT NULL REFERENCES person (person_id),
    charity_id        BIGINT      NOT NULL REFERENCES charity (charity_id),
    amount_cents      BIGINT      NOT NULL,
    currency          CHAR(3)     NOT NULL,
    donated_at        TIMESTAMPTZ NOT NULL,
    note              TEXT,
    -- Same invariant and maintenance as D20 (see triggers.sql): TRUE on
    -- exactly one donation per donor, their most recent.
    is_last_donation  BOOLEAN     NOT NULL DEFAULT false,
    PRIMARY KEY ((person_id) HASH, donated_at DESC, donation_id ASC)
);
