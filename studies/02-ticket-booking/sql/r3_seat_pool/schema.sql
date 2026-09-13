-- R3 — seat-pool
--
-- A hybrid of the two families the study was asked to compare. The seats are
-- pre-created -- but as narrow tokens in an extra table, seat_slot, one row per
-- UNSOLD seat. The ticket is created at booking, carrying the seat it was cut
-- from. Booking moves a seat from the pool into a ticket:
--
--   take a slot nobody has locked (FOR UPDATE SKIP LOCKED), delete it,
--   insert the ticket with that seat number -- one transaction.
--
-- Overbooking is impossible because a slot can be deleted only once, and there
-- are only `capacity` of them; a second transaction that reaches for a deleted
-- slot re-checks it, finds it gone, and SKIP LOCKED means it never waited for it.
--
-- Against P2 (the same SKIP LOCKED pattern on pre-created ticket rows) the pair
-- P2 -> R3 asks: is it cheaper to UPDATE a wide pre-created ticket in place, or
-- to DELETE a two-column token and INSERT the ticket? The pool also shrinks as
-- the event sells, so "seats left" is a count over a table that gets smaller as
-- it gets hotter -- the opposite of P2's partial index, which accumulates dead
-- entries.

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

-- The extra table: unsold seats only. Primary key order (event, seat) keeps an
-- event's pool together and in seat order on both engines.
CREATE TABLE seat_slot (
    event_id  BIGINT NOT NULL REFERENCES event (event_id),
    seat_no   INT NOT NULL,
    PRIMARY KEY (event_id, seat_no)
);

CREATE TABLE ticket (
    ticket_id    BIGINT PRIMARY KEY,
    event_id     BIGINT NOT NULL REFERENCES event (event_id),
    seat_no      INT NOT NULL,
    customer_id  BIGINT NOT NULL,
    sold_at      TIMESTAMPTZ NOT NULL,
    price_cents  BIGINT NOT NULL
);
