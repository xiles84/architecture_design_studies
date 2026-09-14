-- S4 — check-then-hold-serializable: tables.
--
-- S0 — check-then-hold-rc  (NEGATIVE CONTROL)   and   S4 — check-then-hold-serializable
--
-- These files are identical in both directories except for their first line; the
-- two designs differ only in the isolation level the harness uses.
--
-- The hold is written the way a first implementation usually is:
--
--   w_read_seats  -> are all the seats of the block free?
--   w_set_held    -> yes: UPDATE them to held, unconditionally
--
-- At READ COMMITTED (S0) that is wrong. Two buyers choose the same seat and both
-- read it free. The first updates and commits. The second's UPDATE waits for the
-- first's row lock, re-evaluates its WHERE clause on the new row version -- and
-- still matches, because the WHERE clause says nothing about status. It overwrites
-- the first hold silently (engine probe P11). Both buyers were told "the seats are
-- yours for 40 minutes"; the first finds out when paying. S0 is expected to show
-- theft and early rejections, and is reported as a failure if it does not.
--
-- At SERIALIZABLE (S4) the same SQL is correct: the two transactions read and then
-- write the same rows, the engine detects the conflict and one of them fails with
-- SQLSTATE 40001 (engine probe P7). The harness retries it, the retry reads the
-- seat held, and that buyer chooses again. Unlike study 02's C2, whose buyers all
-- read one event-wide count, the conflict set here is only the block -- buyers of
-- different seats do not conflict at all.
--
-- The confirmation, release, sweeper and reads are S1's.

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
