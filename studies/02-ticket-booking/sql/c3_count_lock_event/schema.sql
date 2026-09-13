-- C3 — count-lock-event
--
-- Tables byte-identical to C1. The change is one statement at the start of the
-- booking transaction: lock the PARENT event row with SELECT ... FOR UPDATE,
-- then count and insert exactly as C1 does, at READ COMMITTED.
--
-- The lock turns "check, then act" into a critical section per event. The count
-- runs after the lock is granted, and at READ COMMITTED each statement takes a
-- fresh snapshot, so it sees every ticket the previous lock holder committed.
--
-- This is the pessimistic answer that keeps the naive design's data model. Its
-- cost is the count, now paid inside the critical section: every buyer of a
-- 100 000-seat event waits for the buyer ahead of them to count up to 100 000.

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
    description  TEXT NOT NULL
);

CREATE TABLE ticket (
    ticket_id    BIGINT PRIMARY KEY,
    event_id     BIGINT NOT NULL REFERENCES event (event_id),
    -- General admission: a booking made through this design has no seat number.
    -- Tickets present at load carry the seat they were sold as, so that every
    -- design starts from the same logical dataset.
    seat_no      INT,
    customer_id  BIGINT NOT NULL,
    sold_at      TIMESTAMPTZ NOT NULL,
    price_cents  BIGINT NOT NULL
);
