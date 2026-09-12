-- D6 "embedded-jsonb": the child table is folded into the parent row.
--
-- The tree stops at two levels: charity -> person, and a person's donations live
-- inside that person's row as a JSONB array. This is the document-store shape,
-- expressed in PostgreSQL.
--
-- The promise: everything about one person is one row read. No join, no index
-- lookup into a second table, no fan-out.
--
-- The suspicion this study tests:
--   1. Any question that crosses persons (totals, global recency, leaderboards)
--      must unnest every array in the table.
--   2. Appending one donation rewrites the ENTIRE person row. Under MVCC that
--      means a new row version whose size grows with the donation history --
--      write cost is O(history), not O(1).
--   3. Rows above ~2 kB get TOASTed: the array is compressed and moved out of
--      line, so the "single row read" quietly becomes a second fetch.
--
-- Array element shape (kept short deliberately -- JSONB stores every key of
-- every element, so verbose key names multiply across millions of elements):
--   {"i": donation_id, "a": amount_cents, "c": currency, "t": donated_at, "n": note}
-- Elements are appended in donation order, so the array is ascending by "t".

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
    joined_at   TIMESTAMPTZ NOT NULL,
    donations   JSONB       NOT NULL DEFAULT '[]'::jsonb
);
