-- D20 "recency-flag": D3's exact schema, plus one column.
--
-- donation.charity_id is redundant -- it is derivable by joining through person.
-- Denormalising it lets every charity-scoped question touch exactly one table.
-- The cost is a wider row, an extra index to maintain, and an invariant the
-- application must not break (a person's donations must all carry that person's
-- charity_id). Whether that trade is worth it is what the benchmark measures.
--
-- donation.is_last_donation is the owner's proposal for "who last gave in a
-- period" (RECENCY.md): exactly one TRUE row per donor, maintained by a
-- trigger (triggers.sql) that locks the donor's person row first. What it
-- costs: a second UPDATE inside every insert's transaction, and partial-index
-- churn on the hot write path. What it buys: q13-q16 read one narrow index,
-- no aggregate, no second table. d21_recency_flag_unguarded is this same
-- schema with the trigger's lock removed -- the negative control.

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
    donation_id       BIGINT      PRIMARY KEY,
    person_id         BIGINT      NOT NULL REFERENCES person (person_id),
    charity_id        BIGINT      NOT NULL REFERENCES charity (charity_id),  -- denormalised
    amount_cents      BIGINT      NOT NULL,
    currency          CHAR(3)     NOT NULL,
    donated_at        TIMESTAMPTZ NOT NULL,
    note              TEXT,
    -- TRUE on exactly one donation per donor: their most recent. See the
    -- module comment above and triggers.sql for how that invariant holds.
    is_last_donation  BOOLEAN     NOT NULL DEFAULT false
);
