-- L2 — section-document: tables.
--
-- The seats of an event are embedded into their parent section, as one JSONB
-- document per (event, section) holding only the seats that are claimed:
--
--   event_section.claims = {"<seat_id>": [state, hold_id, customer_id, expires_ms], ...}
--   state 1 = held (valid while expires_ms is in the future), 2 = sold
--
-- A free seat is absent. A hold, a confirmation, a release and a sweep are all the
-- same operation: read the section's document and its version, change it in the
-- application, and write it back only if nobody else did in between:
--
--   UPDATE event_section SET claims = $new, version = version + 1
--    WHERE event_id = $e AND section_no = $s AND version = $read_version
--
-- Zero rows means another buyer changed the section first; the application reads it
-- again. This is optimistic concurrency on the parent row, the pattern of document
-- stores -- and study 01's embedding question (D6, D9) applied to seats:
--
--   * the seat map is a single-row read, the cheapest in the study;
--   * every write rewrites the whole section's document, and every writer of a
--     section contends on the same row: a hot section is serialised, whichever
--     seats its buyers chose;
--   * expiry is judged in the application against the database clock it read
--     (w_read_section returns now()), and lazily, as in S1.
--
-- A new hold's expiry is (read now() + TTL), so a hold is a few milliseconds
-- shorter than 40 minutes: never longer, and never early enough to be stolen.

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
-- The inventory: one document per event section.
-- ---------------------------------------------------------------------------

CREATE TABLE event_section (
    event_id     BIGINT NOT NULL REFERENCES event (event_id),
    section_no   INT NOT NULL,
    capacity     INT NOT NULL CHECK (capacity > 0),
    -- Incremented by every write: the compare-and-set token.
    version      BIGINT NOT NULL,
    -- Claimed seats only: {"<seat_id>": [state, hold_id, customer_id, expires_ms]}.
    claims       JSONB NOT NULL,
    -- The earliest expiry among held claims, NULL if none: lets the sweeper find
    -- sections worth tidying without opening every document.
    next_expiry  TIMESTAMPTZ,
    PRIMARY KEY (event_id, section_no)
);
