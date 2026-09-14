-- L1 — claim-rows: tables.
--
-- No per-event seat rows at all. A venue's seats are defined once (venue_seat);
-- an event owns nothing but its event row until someone claims a seat. A claim is
-- a row created on hold:
--
--   seat_claim (event_id, seat_id) PRIMARY KEY -- at most one claim per seat
--
-- and the primary key, not a status check, keeps two buyers off one seat:
--
--   INSERT INTO seat_claim ... SELECT ... FROM unnest(block)
--   ON CONFLICT (event_id, seat_id) DO UPDATE SET <new holder>
--    WHERE seat_claim.expires_at <= now()
--
-- A free seat has no row and is inserted; an expired claim is stolen by the
-- ON CONFLICT branch; a valid claim makes the conflict branch update nothing, so
-- fewer rows come back than seats requested and the caller rolls back (engine
-- probes P3 and P3b). A sold seat keeps its claim with expires_at = 'infinity', so
-- it can never be stolen, and the ticket table's UNIQUE (event_id, seat_id) is a
-- second arbiter for sales.
--
-- This is study 02's question -- pre-create inventory or create it on demand? --
-- asked of marked seats. Publishing a 100 000-seat event is one row instead of
-- 100 000; a hold is an insert or upsert instead of an update; "seats left" is the
-- venue's capacity minus the valid claims, so it no longer grows with the event.
-- Expiry is lazy, as in S1.

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

-- A sold seat's purchase record. UNIQUE (event_id, seat_id) here, unlike every
-- other design: with no seat rows, the ticket table is where a second sale of a
-- seat would have to be stopped, and this design lets a unique index do it.
CREATE TABLE ticket (
    ticket_id    BIGINT PRIMARY KEY,
    event_id     BIGINT NOT NULL REFERENCES event (event_id),
    section_no   INT NOT NULL,
    seat_id      INT NOT NULL,
    customer_id  BIGINT NOT NULL,
    -- NULL for tickets sold before this system existed (loaded).
    hold_id      BIGINT,
    sold_at      TIMESTAMPTZ NOT NULL,
    price_cents  BIGINT NOT NULL,
    CONSTRAINT ticket_one_per_seat UNIQUE (event_id, seat_id)
);

-- ---------------------------------------------------------------------------
-- The inventory: a claim per seat that is held or sold. Free seats have no row.
-- ---------------------------------------------------------------------------

CREATE TABLE seat_claim (
    event_id     BIGINT NOT NULL REFERENCES event (event_id),
    seat_id      INT NOT NULL,
    section_no   INT NOT NULL,
    -- The hold that claimed the seat; NULL for tickets sold before this system (loaded).
    hold_id      BIGINT,
    customer_id  BIGINT NOT NULL,
    -- A hold's expiry; 'infinity' once the seat is sold.
    expires_at   TIMESTAMPTZ NOT NULL,
    claimed_at   TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (event_id, seat_id)
);
