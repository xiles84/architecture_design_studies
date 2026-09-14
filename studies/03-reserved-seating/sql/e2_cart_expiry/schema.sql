-- E2 — cart-expiry: tables.
--
-- S1 copies a hold's expiry onto every seat it covers: a block of four seats
-- stores the same timestamp four times. That is denormalisation, and it has the
-- usual price -- extending, releasing or listing a hold touches every seat row.
--
-- Here the expiry is stored ONCE, on a hold (cart) row. Seat rows only point to the
-- hold that has them:
--
--   hold        (hold_id, customer_id, event_id, section_no, status, expires_at)
--   event_seat  (event_id, seat_id, section_no, status, hold_id)
--
-- A hold is now one INSERT into hold plus S1's conditional UPDATE of the seats,
-- whose "has the previous hold expired?" test becomes a subquery on hold:
--
--   ... AND (status = 'available'
--            OR (status = 'held' AND EXISTS (SELECT 1 FROM hold h
--                 WHERE h.hold_id = event_seat.hold_id AND h.expires_at <= now())))
--
-- Releasing an abandoned cart is one row, and the sweeper works cart by cart. Every
-- seat-map read joins hold.
--
-- THIS DESIGN IS CORRECT ONLY BECAUSE A CART'S expires_at NEVER CHANGES AFTER IT IS
-- WRITTEN. A statement that waits for a row lock re-checks ITS OWN row; what a
-- subquery saw of another table can be stale (LESSONS_LEARNED, study 01's D9 cache).
-- If carts could be extended -- K1's payment window -- a steal could use the old
-- expiry and take a hold that was just extended. Making that safe would mean locking
-- the cart row in every steal. That combination (E2 + K1) is not measured.
-- Engine probe P12 checks the steal-vs-steal race on both engines.

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
-- The inventory: one row per seat per event, and one row per hold.
-- ---------------------------------------------------------------------------

CREATE TABLE hold (
    hold_id      BIGINT PRIMARY KEY,
    customer_id  BIGINT NOT NULL,
    event_id     BIGINT NOT NULL REFERENCES event (event_id),
    section_no   INT NOT NULL,
    -- 'held' until confirmed or released; an expired cart may still say 'held'
    -- until the sweeper reaches it -- expiry is judged by expires_at, lazily.
    status       TEXT NOT NULL CHECK (status IN ('held', 'confirmed', 'released')),
    -- Written once, never updated. See the header.
    expires_at   TIMESTAMPTZ NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL
);

CREATE TABLE event_seat (
    event_id    BIGINT NOT NULL REFERENCES event (event_id),
    seat_id     INT NOT NULL,
    section_no  INT NOT NULL,
    status      TEXT NOT NULL CHECK (status IN ('available', 'held', 'sold')),
    -- The hold that has this seat. NULL for available seats and loaded sales.
    hold_id     BIGINT REFERENCES hold (hold_id),
    customer_id BIGINT,
    sold_at     TIMESTAMPTZ,
    PRIMARY KEY (event_id, seat_id),
    CONSTRAINT event_seat_available_is_empty CHECK (
        status <> 'available' OR (hold_id IS NULL AND customer_id IS NULL AND sold_at IS NULL)),
    CONSTRAINT event_seat_held_has_hold CHECK (
        status <> 'held' OR (hold_id IS NOT NULL AND customer_id IS NULL AND sold_at IS NULL)),
    CONSTRAINT event_seat_sold_has_buyer CHECK (
        status <> 'sold' OR (customer_id IS NOT NULL AND sold_at IS NOT NULL))
);
