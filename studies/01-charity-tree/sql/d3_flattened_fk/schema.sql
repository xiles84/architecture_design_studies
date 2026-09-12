-- D3 "flattened-fk": D2, but the grandchild carries the grandparent's key.
--
-- donation.charity_id is redundant -- it is derivable by joining through person.
-- Denormalising it lets every charity-scoped question touch exactly one table.
-- The cost is a wider row, an extra index to maintain, and an invariant the
-- application must not break (a person's donations must all carry that person's
-- charity_id). Whether that trade is worth it is what the benchmark measures.

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
    charity_id   BIGINT      NOT NULL REFERENCES charity (charity_id),  -- denormalised
    amount_cents BIGINT      NOT NULL,
    currency     CHAR(3)     NOT NULL,
    donated_at   TIMESTAMPTZ NOT NULL,
    note         TEXT
);
