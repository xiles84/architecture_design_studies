-- y2_noncolocated — the non-colocated half of the placement experiment. YugabyteDB
-- only.
--
-- The controlled pair's other half: same columns, same data, same operations, same
-- indexes and the same node count and resource budget as y1_colocated. The single
-- decision that differs is the child table's primary key, and therefore where a
-- donor's donations physically live:
--
--   PRIMARY KEY (donation_id)
--
-- YugabyteDB hashes the leading column, so donation_id is hashed and one donor's
-- donations scatter across every tablet in the cluster. "This donor's newest 20"
-- becomes a scatter-gather, and on a three-node cluster that is real cross-node RPC
-- traffic.
--
-- This is the shape a legacy application almost always has, because donation_id is
-- the natural identity key. That is why the pair is worth measuring rather than
-- asserting: the placement decision is invisible in the schema's columns and very
-- visible in the plan's RPC counts.
--
-- LIMITATION, stated with the numbers wherever they are used: all three nodes are
-- containers on one WSL machine, so this measures distributed machinery and not
-- network behaviour.
--
-- The queries/writes/audit are the flattened baseline's (see the design registry's
-- SQLDir).

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
