-- P2 — precreated-skip-locked
--
-- Byte-identical tables to P1. The single change is in writes.sql: the seat is
-- picked with FOR UPDATE SKIP LOCKED, so a buyer passes over seats other buyers
-- are in the middle of taking instead of queueing behind them. This is the
-- job-queue pattern applied to inventory.
--
-- What SKIP LOCKED costs is an ambiguity: "no row" no longer means "sold out",
-- it can also mean "every remaining seat is locked right now by someone whose
-- transaction may yet roll back". Answering "sold out" on that would under-sell
-- the event. writes.sql therefore carries an exact, non-locking fallback check.

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
