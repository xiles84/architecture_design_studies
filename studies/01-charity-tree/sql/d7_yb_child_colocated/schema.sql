-- D7 "yb-child-colocated": D3's columns, but the DONATION PRIMARY KEY is chosen
-- for data placement instead of for identity. YugabyteDB only.
--
-- YugabyteDB shards a table by hashing the leading HASH column(s) of its primary
-- key. With the default `donation_id BIGINT PRIMARY KEY`, donation_id is hashed,
-- so one person's 40 donations land on ~40 different tablets, and "the last 20
-- donations of person X" becomes a scatter-gather across every tablet in the
-- cluster. On a 3-node cluster that is real cross-node RPC traffic.
--
-- Here the key is ((person_id) HASH, donated_at DESC, donation_id ASC):
--   * person_id HASH  -> every donation of a person lives in ONE tablet, on one
--                        node. The parent/child locality of the tree is preserved
--                        in the physical layout.
--   * donated_at DESC -> inside that tablet the rows are already stored newest
--                        first, so "last N donations" is a prefix read with no
--                        sort and no extra index.
--   * donation_id ASC -> tiebreaker, keeps the key unique.
--
-- What it costs: donation_id is no longer the primary key, so a point lookup by
-- donation_id needs a secondary index -- which in YugabyteDB is itself a
-- distributed table, making that lookup a two-hop operation. Trading a fast
-- path for a fast path is the whole point; the benchmark says which way the
-- trade lands for this query mix.
--
-- The same locality idea is spelled `PARTITION BY HASH` + local indexes in
-- PostgreSQL, and `PARTITION KEY` in Cassandra-family stores. It is a general
-- tree-design lever, not a YugabyteDB quirk.

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
    person_id    BIGINT      NOT NULL REFERENCES person (person_id),
    charity_id   BIGINT      NOT NULL REFERENCES charity (charity_id),
    amount_cents BIGINT      NOT NULL,
    currency     CHAR(3)     NOT NULL,
    donated_at   TIMESTAMPTZ NOT NULL,
    note         TEXT,
    PRIMARY KEY ((person_id) HASH, donated_at DESC, donation_id ASC)
);
