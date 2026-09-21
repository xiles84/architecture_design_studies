-- Study 05 — OWNED database model.
--
-- The database may be changed for cache correctness. Three additions, and nothing
-- else, because every addition is a cost the legacy model does not pay and that
-- cost is what the owned-vs-legacy comparison measures:
--
--   * person.cache_version -- a monotonic token on the authoritative parent. A
--     cache entry records the version of the state it holds, so a stale entry is
--     detectable without reading the cached content;
--   * cache_outbox -- a transactional change log. Every mutation that affects the
--     portal view appends an event in the SAME transaction as the mutation, so
--     "the version moved" and "the outbox has the event" cannot disagree with
--     "the mutation committed";
--   * atomic version changes -- each mutation statement's transaction bumps the
--     version, so no reader can see a committed mutation at the old version.
--
-- Ordered locking for reassignment lives in the write catalogue
-- (w_lock_persons2), because it is a statement, not a schema feature.
--
-- The schema is otherwise byte-identical to the legacy model, which is what makes
-- the pair a controlled comparison: same tables, same columns, same data, same
-- operations; the version token and the outbox are the only design difference.

CREATE TABLE charity (
    charity_id  BIGINT      PRIMARY KEY,
    name        TEXT        NOT NULL,
    country     TEXT        NOT NULL,
    founded_on  DATE        NOT NULL
);

CREATE TABLE person (
    person_id     BIGINT      PRIMARY KEY,
    charity_id    BIGINT      NOT NULL REFERENCES charity (charity_id),
    full_name     TEXT        NOT NULL,
    email         TEXT        NOT NULL,
    joined_at     TIMESTAMPTZ NOT NULL,
    cache_version BIGINT      NOT NULL DEFAULT 0
);

CREATE TABLE donation (
    donation_id  BIGINT      PRIMARY KEY,
    person_id    BIGINT      NOT NULL REFERENCES person (person_id),
    charity_id   BIGINT      NOT NULL REFERENCES charity (charity_id),  -- rolldown
    amount_cents BIGINT      NOT NULL,
    currency     CHAR(3)     NOT NULL,
    donated_at   TIMESTAMPTZ NOT NULL,
    note         TEXT
);

-- The change log. `person_id` is the cache key the event invalidates; `kind` names
-- the logical mutation so an operator can replay the log; `donation_id` is set for
-- donation-scoped events and NULL for a person metadata change. It is deliberately
-- not a queue with a consumer in this study: the study measures what having the
-- event costs and what it enables, not a full outbox-relay deployment.
CREATE TABLE cache_outbox (
    seq         BIGSERIAL   PRIMARY KEY,
    person_id   BIGINT      NOT NULL,
    kind        TEXT        NOT NULL,
    donation_id BIGINT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
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
