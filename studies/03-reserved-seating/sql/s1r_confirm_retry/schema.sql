-- S1r — confirm-retry: tables.
--
-- S1r is S1 with one decision added, in the application, not in this SQL: a
-- confirmation that matches fewer seats than the hold is rolled back and run once
-- more, at once, in a new transaction. Every file here is S1's but for the header.
-- The retry is safe by construction: w_confirm_seats is still guarded, so a second
-- attempt can never sell a hold that expired, was released or was taken. It exists
-- because YugabyteDB was seen refusing valid holds transiently under load (ER-01);
-- S1 -> S1r measures whether one retry removes that, and what it costs.
--
-- Reserved seating: every seat of every event exists as a row of event_seat,
-- created when the event is published, and the row itself carries the seat's
-- state -- available, held (by whom, until when) or sold.
--
-- What keeps two buyers off one seat is a single conditional UPDATE over the
-- whole block the buyer chose:
--
--   UPDATE event_seat SET status = 'held', hold_id = ..., hold_expires_at = now() + 40 min
--    WHERE event_id = ... AND section_no = ... AND seat_id = ANY (the block)
--      AND (status = 'available' OR (status = 'held' AND hold_expires_at <= now()))
--
-- If it returns fewer rows than seats requested, some seat was not free, and the
-- transaction rolls back: all or nothing. A second buyer reaching for the same
-- seat waits for the first buyer's row lock, then READ COMMITTED re-evaluates the
-- WHERE clause on the new row version and finds the seat held (engine probe P10).
--
-- Expiry is LAZY: an expired hold is simply a row whose hold_expires_at is in the
-- past. Nothing has to happen for the seat to become available again -- the next
-- buyer's UPDATE takes it, and every read judges expiry with now(). Two
-- consequences, both deliberate:
--
--   * correctness never depends on the sweeper; the sweeper only tidies rows up;
--   * no exact availability counter can be maintained anywhere, because an
--     expiry changes no row. Every "seats left" is a count over seat rows.
--
-- Every clock comparison uses the database's now(), which inside a transaction
-- is the transaction start time on both engines (engine probe P1).

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
