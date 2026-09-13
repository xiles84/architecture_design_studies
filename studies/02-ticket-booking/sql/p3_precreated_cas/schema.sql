-- P3 — precreated-cas
--
-- Byte-identical tables to P1 and P2. The change is the concurrency control:
-- no locking read and no explicit transaction at all. A buyer reads a candidate
-- seat, then sells it with a compare-and-set UPDATE (`... AND status =
-- 'available'`). Zero rows affected means someone else sold it first: pick
-- another and try again. Optimistic concurrency, with the database's row-level
-- write check as the arbiter.
--
-- An optimistic design is only as good as its candidate choice. If every buyer
-- read "the lowest available seat" they would all collide on it, and P3 would be
-- a strawman. Each attempt therefore starts its search at a random seat, which
-- spreads buyers across the inventory; collisions then rise naturally as the
-- event fills and the free seats become few.

CREATE TABLE band (
    band_id  BIGINT PRIMARY KEY,
    name     TEXT NOT NULL,
    genre    TEXT NOT NULL,
    country  TEXT NOT NULL
);

CREATE TABLE event (
    event_id     BIGINT PRIMARY KEY,
    band_id      BIGINT NOT NULL REFERENCES band (band_id),
    name         TEXT NOT NULL,
    venue        TEXT NOT NULL,
    starts_at    TIMESTAMPTZ NOT NULL,
    capacity     INT NOT NULL CHECK (capacity > 0),
    price_cents  BIGINT NOT NULL,
    -- Edited by the organiser while a sale is running (the race's "editor").
    -- In designs that keep a counter on this row, those edits contend with
    -- bookings; here they do not.
    description  TEXT NOT NULL
);

CREATE TABLE ticket (
    ticket_id    BIGINT PRIMARY KEY,
    event_id     BIGINT NOT NULL REFERENCES event (event_id),
    seat_no      INT NOT NULL,
    status       TEXT NOT NULL CHECK (status IN ('available', 'sold')),
    customer_id  BIGINT,
    sold_at      TIMESTAMPTZ,
    price_cents  BIGINT NOT NULL,
    -- A sold ticket has a buyer and an unsold one does not. Cheap, and it turns
    -- a half-applied sale or cancellation into an error instead of a quiet
    -- inconsistency.
    CONSTRAINT ticket_sold_has_buyer CHECK ((status = 'sold') = (customer_id IS NOT NULL))
);
