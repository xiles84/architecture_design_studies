-- D10 "embedded-hybrid-locked" (schema identical to D9): embed a BOUNDED slice of the children, keep the table.
--
-- D6 folds the whole donation table into the person row and discovers the
-- catch: appending one donation rewrites the entire document, so the cost of a
-- write grows with the history that is already there. The array is unbounded,
-- so the write cost is unbounded.
--
-- D9 is the shape production systems actually converge on. The donation table
-- stays exactly as it is in D3 -- it remains the source of truth, fully indexed,
-- able to answer every question. On top of it, person carries a cache of just
-- the most recent 20 donations, maintained by a trigger.
--
--   * the highest-QPS query in a system like this -- "show me this donor's
--     recent gifts" -- is answered from one row, with no second table, no index
--     descent and no sort, exactly as in D6
--   * every other question still has a real, indexed table to work against, so
--     none of D6's cross-person collapse applies
--   * write amplification is CAPPED at 20 elements instead of growing without
--     limit, so unlike D6 the insert cost does not depend on history length
--
-- What it costs: a second copy of recent data that can drift, a trigger on the
-- hot write path, and a cache size (20) that is now a design constant which has
-- to match what the application actually asks for. Ask for 21 and the cache
-- silently stops being usable.
--
-- The array is maintained NEWEST FIRST, which is the opposite of D6's
-- chronological order -- because here the useful end is the recent end.

CREATE TABLE charity (
    charity_id  BIGINT      PRIMARY KEY,
    name        TEXT        NOT NULL,
    country     TEXT        NOT NULL,
    founded_on  DATE        NOT NULL
);

CREATE TABLE person (
    person_id        BIGINT      PRIMARY KEY,
    charity_id       BIGINT      NOT NULL REFERENCES charity (charity_id),
    full_name        TEXT        NOT NULL,
    email            TEXT        NOT NULL,
    joined_at        TIMESTAMPTZ NOT NULL,
    -- bounded cache: the 20 most recent donations, newest first
    recent_donations JSONB       NOT NULL DEFAULT '[]'::jsonb
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
