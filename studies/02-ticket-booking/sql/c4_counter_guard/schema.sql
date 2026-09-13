-- C4 — counter-guard
--
-- C3 serialised bookings per event and then counted inside the critical section.
-- C4 keeps the serialisation and removes the count: the event row carries
-- seats_sold, and the booking claims a seat with ONE conditional update,
--
--   UPDATE event SET seats_sold = seats_sold + 1
--    WHERE event_id = $1 AND seats_sold < capacity
--
-- The row lock taken by that UPDATE is the critical section; the WHERE clause,
-- re-checked on the latest row version after any wait, is the capacity check.
-- Zero rows means sold out. The ticket is inserted in the same transaction.
--
-- The CHECK constraint below would stop an overbooking even if the guard were
-- deleted from the application's SQL. That is belt and braces on purpose: the
-- guard gives a clean "sold out" answer, the constraint makes the invariant the
-- database's job rather than every code path's.

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
    description  TEXT NOT NULL,
    seats_sold   INT NOT NULL DEFAULT 0,
    CONSTRAINT event_seats_sold_in_range CHECK (seats_sold BETWEEN 0 AND capacity)
);

CREATE TABLE ticket (
    ticket_id    BIGINT PRIMARY KEY,
    event_id     BIGINT NOT NULL REFERENCES event (event_id),
    seat_no      INT,
    customer_id  BIGINT NOT NULL,
    sold_at      TIMESTAMPTZ NOT NULL,
    price_cents  BIGINT NOT NULL
);
