-- Reference cell: embedded-locked. The embedding half of the
-- information-placement comparison, and the D9/D10 precedent applied here.
--
-- The portal's most-requested answer -- the newest 20 donations -- stops being a
-- second query against the child table and becomes a bounded slice stored on the
-- parent. One decision, isolated against the flattened baseline: same schema
-- otherwise, same indexes, same mutations, same reads for everything except the
-- recent slice.
--
-- Study 01's D9 is the naive version of this and its D10 is the concurrency-correct
-- one. Study 01 measured D9 silently corrupting 3-5 donor caches per 30-second
-- window once reads and writes overlapped, with two diagnosed causes: prepending by
-- commit order instead of the table's (donated_at, donation_id) order, and a
-- rebuild statement whose snapshot could not see a concurrently committed insert.
-- D10 fixed both by merging-and-sorting on insert and by locking the parent before
-- a rebuild.
--
-- THIS STUDY USES THE LOCKED, CORRECT VARIANT ONLY. The naive variant is not a
-- design here: it is already known to violate the invariant this study's whole
-- correctness floor rests on, and re-measuring a known-broken cache would add a
-- wrong number rather than a finding. The naive form is preserved as a precedent in
-- Study 01 (tags study-01/v1 and study-01/v3-mechanisms).

CREATE TABLE charity (
    charity_id  BIGINT      PRIMARY KEY,
    name        TEXT        NOT NULL,
    country     TEXT        NOT NULL,
    founded_on  DATE        NOT NULL
);

CREATE TABLE person (
    person_id          BIGINT      PRIMARY KEY,
    charity_id         BIGINT      NOT NULL REFERENCES charity (charity_id),
    full_name          TEXT        NOT NULL,
    email              TEXT        NOT NULL,
    joined_at          TIMESTAMPTZ NOT NULL,
    -- The bounded slice: the 20 most recent donations, newest first.
    recent_donations   JSONB       NOT NULL DEFAULT '[]'::jsonb
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
