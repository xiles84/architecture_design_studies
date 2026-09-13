-- C1 — count-naive  (NEGATIVE CONTROL)   and   C2 — count-serializable
--
-- These files are byte-identical in c1_count_naive/ and c2_count_serializable/.
-- The two designs differ ONLY in the isolation level the harness requests for
-- the booking transaction (harness/designs.go): READ COMMITTED for C1,
-- SERIALIZABLE for C2.
--
-- No ticket row exists until a seat is sold. Booking is "check, then act":
--   count the event's tickets; if fewer than capacity, insert one.
--
-- At READ COMMITTED that is the classic race. Two buyers both count 99 of 100,
-- both insert, and the event holds 101 tickets. Nothing in the schema stops it:
-- the rows are different rows, so no lock and no constraint is ever contended.
-- C1 is in this study BECAUSE it is wrong. A correctness audit that has never
-- caught a wrong design has not been shown to work; C1 is the proof that the
-- overbooking audit can see an overbooking, and the report says so explicitly
-- whenever it fails to.
--
-- At SERIALIZABLE the same statements are safe: the engine detects that the two
-- transactions' reads and writes cannot be ordered serially and aborts one of
-- them (SQLSTATE 40001), which the harness retries. "Just use the right
-- isolation level" is a real strategy, and C2 prices it.

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
    ticket_id    BIGINT PRIMARY KEY,
    event_id     BIGINT NOT NULL REFERENCES event (event_id),
    -- General admission: a booking made through this design has no seat number.
    -- Tickets present at load carry the seat they were sold as, so that every
    -- design starts from the same logical dataset.
    seat_no      INT,
    customer_id  BIGINT NOT NULL,
    sold_at      TIMESTAMPTZ NOT NULL,
    price_cents  BIGINT NOT NULL
);
