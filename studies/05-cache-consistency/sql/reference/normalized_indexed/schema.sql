-- Reference cell: normalized-indexed. The rolldown reference.
--
-- Identical to the flattened-FK model EXCEPT that the grandparent's key is not
-- copied onto the child: `donation` has no charity_id, and every charity-scoped
-- question joins through person. One decision, isolated, on the same data and the
-- same operations.
--
-- Study 01's D2 is this shape. It exists here because the study's blended workload
-- contains one charity-scoped read (r_charity_recent) that discriminates the two:
-- on the flattened model that read filters on the copied column and never touches
-- person; here it must join. What the copy costs on the write path -- a wider row
-- and one more index to maintain on every insert -- is visible in the write
-- numbers of the same cell.

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

-- No charity_id: the copied grandparent key is the one decision this cell removes.
CREATE TABLE donation (
    donation_id  BIGINT      PRIMARY KEY,
    person_id    BIGINT      NOT NULL REFERENCES person (person_id),
    amount_cents BIGINT      NOT NULL,
    currency     CHAR(3)     NOT NULL,
    donated_at   TIMESTAMPTZ NOT NULL,
    note         TEXT
);

-- The staging table keeps the same shape as every other design, so the loader is
-- held constant. The design's load statement simply ignores charity_id.
CREATE TABLE load_donation (
    donation_id  BIGINT      NOT NULL,
    person_id    BIGINT      NOT NULL,
    charity_id   BIGINT      NOT NULL,
    amount_cents BIGINT      NOT NULL,
    currency     CHAR(3)     NOT NULL,
    donated_at   TIMESTAMPTZ NOT NULL,
    note         TEXT
);
