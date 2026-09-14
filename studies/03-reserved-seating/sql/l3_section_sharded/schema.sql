-- L3 — section-sharded (YugabyteDB only): tables.
--
-- S1 with a different primary key on event_seat, and nothing else: the same
-- statements, the same indexes except the one this key makes redundant.
--
-- YugabyteDB hash-shards a table by the first column of its primary key. S1's key
-- (event_id, seat_id) therefore puts every seat of an event on one tablet: a hot
-- drop for one event is a hot tablet, whatever the number of nodes. Study 02 saw
-- three nodes buy nothing for one event's sale.
--
-- Here the hash covers (event_id, section_no), so the sections of an event spread
-- over the table's tablets, and a buyer's hold -- always within one section, with
-- section_no in every predicate -- still touches one tablet. Adding nodes can then
-- spread one event's load.
--
-- The syntax is YugabyteDB's; PostgreSQL has no hash-sharded keys, so this design
-- is never run there.

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
    -- YugabyteDB: hash over (event, section), range over the seat within it.
    PRIMARY KEY ((event_id, section_no) HASH, seat_id ASC),
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
