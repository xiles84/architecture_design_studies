-- E1 — sweeper-expiry: tables.
--
-- Byte-identical tables to S1. The change is where expiry takes effect. In S1 an
-- expired hold is available the instant it expires, because every statement
-- compares hold_expires_at with now(). Here expiry is MATERIALISED: a hold
-- stays 'held' until the sweeper finds it expired and sets it back to
-- 'available'.
--
--   w_hold_seats  takes only status = 'available' seats
--   q01..q03      trust status: a held row is shown held, expired or not
--   the sweeper   is the only way an abandoned seat comes back
--
-- This is how many systems are built -- status columns are simpler to index and to
-- reason about than time comparisons, and a counter of available seats could be
-- maintained on top of it. The price: every abandoned seat is unavailable for the
-- gap between its expiry and the next sweep, and indefinitely if the sweeper
-- stops. The lifecycle experiment stops the sweeper for 40 human minutes in the
-- first event of each tier and probes what that does; the lazy designs are
-- expected to be unaffected.
--
-- Confirmation still checks hold_expires_at > now(): a hold the sweeper has not
-- reached yet is expired all the same.

-- ---------------------------------------------------------------------------
-- Common tables: identical in every design of this study.
-- ---------------------------------------------------------------------------

CREATE TABLE band (
    band_id  BIGINT PRIMARY KEY,
    name     TEXT NOT NULL,
    genre    TEXT NOT NULL,
    country  TEXT NOT NULL
);

CREATE TABLE venue (
    venue_id  BIGINT PRIMARY KEY,
    name      TEXT NOT NULL,
    capacity  INT NOT NULL CHECK (capacity > 0)
);

CREATE TABLE venue_section (
    venue_id    BIGINT NOT NULL REFERENCES venue (venue_id),
    section_no  INT NOT NULL,
    capacity    INT NOT NULL CHECK (capacity > 0),
    PRIMARY KEY (venue_id, section_no)
);

-- The physical seat map. Static: the client caches it, like a seat-map image.
CREATE TABLE venue_seat (
    venue_id      BIGINT NOT NULL REFERENCES venue (venue_id),
    seat_id       INT NOT NULL,
    section_no    INT NOT NULL,
    row_no        INT NOT NULL,
    seat_no       INT NOT NULL,
    quality_rank  INT NOT NULL,
    PRIMARY KEY (venue_id, seat_id)
);

CREATE TABLE event (
    event_id     BIGINT PRIMARY KEY,
    band_id      BIGINT NOT NULL REFERENCES band (band_id),
    venue_id     BIGINT NOT NULL REFERENCES venue (venue_id),
    name         TEXT NOT NULL,
    starts_at    TIMESTAMPTZ NOT NULL,
    price_cents  BIGINT NOT NULL,
    description  TEXT NOT NULL
);

-- A sold seat's purchase record. Deliberately WITHOUT a unique (event_id,
-- seat_id): in this design the seat row is the arbiter, and a second arbiter
-- here would stop the negative controls S0 and K0 from ever showing the double
-- sale they exist to show. (L1, which has no seat rows, makes it unique.)
CREATE TABLE ticket (
    ticket_id    BIGINT PRIMARY KEY,
    event_id     BIGINT NOT NULL REFERENCES event (event_id),
    section_no   INT NOT NULL,
    seat_id      INT NOT NULL,
    customer_id  BIGINT NOT NULL,
    -- NULL for tickets sold before this system existed (loaded).
    hold_id      BIGINT,
    sold_at      TIMESTAMPTZ NOT NULL,
    price_cents  BIGINT NOT NULL
);

-- ---------------------------------------------------------------------------
-- The inventory: one row per seat per event.
-- ---------------------------------------------------------------------------

CREATE TABLE event_seat (
    event_id         BIGINT NOT NULL REFERENCES event (event_id),
    seat_id          INT NOT NULL,
    -- Copied from venue_seat (immutable) so a section's seat map is one index
    -- range, without a join.
    section_no       INT NOT NULL,
    status           TEXT NOT NULL CHECK (status IN ('available', 'held', 'sold')),
    hold_id          BIGINT,
    customer_id      BIGINT,
    hold_expires_at  TIMESTAMPTZ,
    sold_at          TIMESTAMPTZ,
    PRIMARY KEY (event_id, seat_id),
    -- A half-applied hold, release or sale becomes an error instead of a quiet
    -- inconsistency. A held row may be past its expiry: that is lazy expiry.
    CONSTRAINT event_seat_available_is_empty CHECK (
        status <> 'available'
        OR (hold_id IS NULL AND customer_id IS NULL AND hold_expires_at IS NULL AND sold_at IS NULL)),
    CONSTRAINT event_seat_held_has_holder CHECK (
        status <> 'held'
        OR (hold_id IS NOT NULL AND customer_id IS NOT NULL AND hold_expires_at IS NOT NULL AND sold_at IS NULL)),
    CONSTRAINT event_seat_sold_has_buyer CHECK (
        status <> 'sold' OR (customer_id IS NOT NULL AND sold_at IS NOT NULL))
);
