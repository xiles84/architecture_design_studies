-- R1 — inventory-row
--
-- C4 with the counter moved into its own table. The booking claims a seat with
-- the same guarded update, but on event_inventory instead of event:
--
--   UPDATE event_inventory SET remaining = remaining - 1
--    WHERE event_id = $1 AND remaining > 0
--
-- Correctness is identical to C4. What moves is WHICH ROW is hot. In C4 the
-- event row is both the catalogue entry (name, venue, description, read by
-- every page view and edited by the organiser) and the inventory lock every
-- buyer queues on. R1 splits those roles: the narrow inventory row takes the
-- contention; the event row goes back to being read and edited in peace.
--
-- The pair C4 -> R1 isolates that split. Whether it matters depends on the
-- engine's locking granularity: PostgreSQL locks whole rows, so an organiser's
-- edit and a sale conflict in C4; YugabyteDB detects conflicts per column for
-- updates that touch different columns, so it may not care.

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
    description  TEXT NOT NULL
);

-- The extra table: one narrow row per event, holding only what booking contends on.
CREATE TABLE event_inventory (
    event_id   BIGINT PRIMARY KEY REFERENCES event (event_id),
    remaining  INT NOT NULL,
    CONSTRAINT event_inventory_non_negative CHECK (remaining >= 0)
);

CREATE TABLE ticket (
    ticket_id    BIGINT PRIMARY KEY,
    event_id     BIGINT NOT NULL REFERENCES event (event_id),
    seat_no      INT,
    customer_id  BIGINT NOT NULL,
    sold_at      TIMESTAMPTZ NOT NULL,
    price_cents  BIGINT NOT NULL
);
