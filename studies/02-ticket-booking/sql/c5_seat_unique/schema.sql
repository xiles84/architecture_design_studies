-- C5 — seat-unique
--
-- No lock, no counter, no serializable transaction. The database's UNIQUE index
-- is the arbiter.
--
-- An event of capacity N has exactly N seat numbers, 1..N. Every ticket must
-- carry one (NOT NULL), no two tickets of an event may share one (UNIQUE), and
-- none may lie outside the venue (CHECK). With those three constraints an
-- overbooking is not prevented by the booking code at all -- it is
-- unrepresentable in the table.
--
-- The booking path is optimistic: read the highest seat number sold, try to
-- insert the next one with ON CONFLICT DO NOTHING, and if another buyer took it
-- first, read again. Every buyer of an event contends for the same next number,
-- so under a hot drop most attempts lose -- a retry storm instead of a queue.
--
-- The CHECK needs the capacity on the ticket row, because a CHECK constraint
-- cannot look at another table. The insert copies it from the event in the same
-- statement (INSERT ... SELECT FROM event), so the application never supplies
-- it. The known cost of that denormalisation: changing a venue's capacity later
-- would not revisit tickets already sold.
--
-- The design's structural weakness is cancellation. "Next seat = highest + 1"
-- never goes back to fill a hole, so after refunds the event reaches seat N and
-- reports sold out while seats sit empty. That is an UNDER-booking, invisible to
-- any overbooking check; the churn race measures it.

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

CREATE TABLE ticket (
    ticket_id       BIGINT PRIMARY KEY,
    event_id        BIGINT NOT NULL REFERENCES event (event_id),
    seat_no         INT NOT NULL,
    event_capacity  INT NOT NULL,
    customer_id     BIGINT NOT NULL,
    sold_at         TIMESTAMPTZ NOT NULL,
    price_cents     BIGINT NOT NULL,
    CONSTRAINT ticket_seat_inside_venue CHECK (seat_no BETWEEN 1 AND event_capacity)
);
