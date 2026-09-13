-- H0 — hold-naive-confirm  (NEGATIVE CONTROL)   and   H1 — hold-checked-confirm
--
-- These tables are byte-identical in both directories; the two designs differ
-- in one statement of writes.sql (w_confirm_hold).
--
-- Real ticket sales are two steps with a person in between: the seat is HELD
-- while the buyer pays, and the hold expires if they wander off. The extra table
-- here is the reservation: a hold with an expiry.
--
-- The capacity counter on the event (seats_taken) counts held AND sold seats, so
-- a seat in someone's basket cannot be sold to anyone else. Taking a hold is
-- C4's guarded increment. A sweeper releases holds that have expired, giving the
-- seat back to the counter.
--
-- That creates the overbooking hazard this pair (H0 -> H1) exists to measure:
-- a buyer whose payment finishes AFTER their hold expired. Between the expiry and
-- the confirmation, the sweeper may have released the seat and another buyer
-- may have taken it. Confirming anyway sells the seat twice.
--
-- H0 confirms the way a first implementation usually does: it checked the hold
-- before taking payment, so after payment it simply marks the reservation
-- confirmed and issues the ticket. That check was true when it ran. It is the
-- negative control of the hold experiment -- expected to overbook whenever a
-- payment straddles an expiry, and reported as a failure if it does not.
--
-- H1 confirms with ONE conditional statement:
--
--   UPDATE reservation SET status = 'confirmed'
--    WHERE reservation_id = $1 AND status = 'held' AND expires_at > now()
--
-- Zero rows means the hold is gone and the buyer must be refunded, not ticketed.
-- The sweeper's release is the mirror image (status = 'held' AND expires_at <=
-- now()), and both take the reservation's row lock, so exactly one of them wins.
-- The expiry is judged by the database clock in both statements, never by the
-- application's.

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
    -- Held + sold. Never more than the venue.
    seats_taken  INT NOT NULL DEFAULT 0,
    CONSTRAINT event_seats_taken_in_range CHECK (seats_taken BETWEEN 0 AND capacity)
);

CREATE TABLE reservation (
    reservation_id  BIGINT PRIMARY KEY,
    event_id        BIGINT NOT NULL REFERENCES event (event_id),
    customer_id     BIGINT NOT NULL,
    status          TEXT NOT NULL CHECK (status IN ('held', 'confirmed', 'released')),
    expires_at      TIMESTAMPTZ NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL
);

CREATE TABLE ticket (
    ticket_id       BIGINT PRIMARY KEY,
    event_id        BIGINT NOT NULL REFERENCES event (event_id),
    seat_no         INT,
    customer_id     BIGINT NOT NULL,
    sold_at         TIMESTAMPTZ NOT NULL,
    price_cents     BIGINT NOT NULL,
    -- NULL for tickets present at load (sold before this system existed).
    reservation_id  BIGINT REFERENCES reservation (reservation_id)
);
