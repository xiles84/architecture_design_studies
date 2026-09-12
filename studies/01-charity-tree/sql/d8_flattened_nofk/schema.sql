-- D8 "flattened-nofk": D3 with the FOREIGN KEY constraints removed.
--
-- Columns, indexes, queries and write statements are otherwise byte-identical to
-- D3, which makes the D3/D8 pair a controlled experiment on ONE question:
-- what does referential integrity enforcement cost?
--
-- What a foreign key actually does on the write path: every INSERT into donation
-- makes the engine prove that the parent rows exist, by taking a
-- KEY SHARE lock on one person row and one charity row. So a single-row insert
-- is really three row accesses plus two shared locks, and those locks are taken
-- on rows that other concurrent inserts are also touching.
--
-- On PostgreSQL that is a local buffer read. On YugabyteDB the parent may live on
-- a different node, so each check can become a cross-node RPC -- which is why
-- "drop the foreign keys" is folk wisdom in distributed SQL and why it deserves
-- an actual measurement rather than a shrug.
--
-- What you give up is not academic: nothing now stops a donation from pointing at
-- a person who does not exist, and nothing keeps donation.charity_id consistent
-- with person.charity_id. The measurement says what that risk buys you.

CREATE TABLE charity (
    charity_id  BIGINT      PRIMARY KEY,
    name        TEXT        NOT NULL,
    country     TEXT        NOT NULL,
    founded_on  DATE        NOT NULL
);

CREATE TABLE person (
    person_id   BIGINT      PRIMARY KEY,
    charity_id  BIGINT      NOT NULL,   -- no REFERENCES
    full_name   TEXT        NOT NULL,
    email       TEXT        NOT NULL,
    joined_at   TIMESTAMPTZ NOT NULL
);

CREATE TABLE donation (
    donation_id  BIGINT      PRIMARY KEY,
    person_id    BIGINT      NOT NULL,  -- no REFERENCES
    charity_id   BIGINT      NOT NULL,  -- no REFERENCES
    amount_cents BIGINT      NOT NULL,
    currency     CHAR(3)     NOT NULL,
    donated_at   TIMESTAMPTZ NOT NULL,
    note         TEXT
);
