-- y1_colocated — the colocated half of the placement experiment. YugabyteDB only.
--
-- The design is the flattened baseline: same columns, same reads, same writes, same
-- indexes. The SINGLE decision this cell makes is the physical placement of a
-- person's donations, expressed through the primary key:
--
--   PRIMARY KEY ((person_id) HASH, donated_at DESC, donation_id ASC)
--
--   * person_id HASH  -> every donation of one donor lives in ONE tablet, on one
--                        node. The tree's parent/child locality survives into the
--                        physical layout, so "this donor's newest 20" is a prefix
--                        read inside one tablet instead of a scatter-gather across
--                        every tablet in the cluster.
--   * donated_at DESC -> inside that tablet rows are stored newest first, so the
--                        recent-slice read needs no sort.
--   * donation_id ASC -> tiebreaker; keeps the key unique.
--
-- What it costs: donation_id is no longer the primary key, so a point lookup by
-- donation_id needs a secondary index, and in YugabyteDB a secondary index is a
-- distributed table of its own -- that lookup becomes two hops. That trade is what
-- the cell measures.
--
-- LIMITATION, stated with the numbers wherever they are used: all three nodes are
-- containers on one WSL machine. This measures the cost of the distributed
-- machinery, not the behaviour of a real network. The portable signal is the
-- EXPLAIN (ANALYZE, DIST) RPC/tablet counters, not milliseconds.
--
-- The queries/writes/audit are the flattened baseline's (see the design registry's
-- SQLDir): the point of a controlled pair is that only one thing changes.

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
    donation_id  BIGINT      NOT NULL,
    person_id    BIGINT      NOT NULL,
    charity_id   BIGINT      NOT NULL REFERENCES charity (charity_id),
    amount_cents BIGINT      NOT NULL,
    currency     CHAR(3)     NOT NULL,
    donated_at   TIMESTAMPTZ NOT NULL,
    note         TEXT,
    PRIMARY KEY ((person_id) HASH, donated_at DESC, donation_id ASC)
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
