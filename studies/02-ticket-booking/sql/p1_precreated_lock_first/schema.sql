-- P1 — precreated-lock-first
--
-- Every seat of every event exists as a ticket row from the moment the event is
-- published. Selling a seat is an UPDATE from 'available' to 'sold'; nothing is
-- ever inserted on the booking path. Overbooking is impossible by construction
-- as long as one available row is never handed to two buyers -- and that is
-- exactly the property each pre-created design (P1, P2, P3) secures differently.
--
-- P1 secures it the textbook way: lock the lowest available seat with
-- SELECT ... FOR UPDATE, then sell it. Correct, and a convoy: every buyer of the
-- same event reaches for the same lowest seat, so they queue behind each other
-- one at a time, however many seats remain.

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
