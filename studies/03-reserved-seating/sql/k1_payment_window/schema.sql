-- K1 — payment-window: tables.
--
-- S1 plus a guaranteed checkout. The owner's requirement is that a buyer never
-- finds out, while paying, that the seats are gone. S1's checked confirmation
-- keeps the invariant but not quite that promise: a buyer who presses "pay" at
-- minute 39 and whose card takes two minutes is refused at minute 41, after
-- paying.
--
-- Here, starting checkout while the hold is valid extends it once:
--
--   UPDATE event_seat SET hold_expires_at = GREATEST(hold_expires_at, now() + 10 min),
--                         checkout_started_at = now()
--    WHERE ... hold_id = $hold AND status = 'held' AND hold_expires_at > now()
--      AND checkout_started_at IS NULL
--
-- A buyer who starts paying inside the 40 minutes is guaranteed the seats for as
-- long as a payment can take. A buyer too late to start is told so BEFORE paying.
-- "Once" matters: without checkout_started_at a buyer could extend a hold forever.
--
-- The price: every checkout is one more multi-row write, and an abandoned
-- checkout keeps its seats up to 10 minutes longer. The lifecycle experiment
-- measures how many boundary and late rejections this removes against S1.

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
    -- Set once, when the buyer starts paying; the payment window is granted then.
    checkout_started_at  TIMESTAMPTZ,
    PRIMARY KEY (event_id, seat_id),
    -- A half-applied hold, release or sale becomes an error instead of a quiet
    -- inconsistency. A held row may be past its expiry: that is lazy expiry.
    CONSTRAINT event_seat_available_is_empty CHECK (
        status <> 'available'
        OR (hold_id IS NULL AND customer_id IS NULL AND hold_expires_at IS NULL AND sold_at IS NULL
            AND checkout_started_at IS NULL)),
    CONSTRAINT event_seat_held_has_holder CHECK (
        status <> 'held'
        OR (hold_id IS NOT NULL AND customer_id IS NOT NULL AND hold_expires_at IS NOT NULL AND sold_at IS NULL)),
    CONSTRAINT event_seat_sold_has_buyer CHECK (
        status <> 'sold' OR (customer_id IS NOT NULL AND sold_at IS NOT NULL))
);
