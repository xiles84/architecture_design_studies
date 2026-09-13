-- P4 — precreated-counter
--
-- P2 plus one denormalised column: event.seats_sold, maintained in the same
-- transaction as the sale. Availability stops being a COUNT over the event's
-- seats and becomes a single-row read, whatever the event's size.
--
-- The price is the point of the pair P2 -> P4. SKIP LOCKED let buyers of one
-- event work in parallel on different seat rows; the counter puts every one of
-- them back onto the SAME event row, in the same transaction, for as long as
-- that transaction lasts. The denormalisation reintroduces the hot row that
-- SKIP LOCKED existed to avoid.

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
    description  TEXT NOT NULL,
    -- Denormalised: always equal to the number of this event's sold tickets.
    -- The CHECK makes the database, not the application, the last line of
    -- defence against a counter that runs past the venue.
    seats_sold   INT NOT NULL DEFAULT 0,
    CONSTRAINT event_seats_sold_in_range CHECK (seats_sold BETWEEN 0 AND capacity)
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
