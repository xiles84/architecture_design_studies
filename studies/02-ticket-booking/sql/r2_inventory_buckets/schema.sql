-- R2 — inventory-buckets
--
-- R1 with the single inventory row split into up to 8 bucket rows per event.
-- A 100 000-seat event gets 8 buckets of 12 500; a 10-seat gig gets 8 buckets
-- of one or two. The capacity invariant becomes "the buckets' remaining seats
-- add up to what is unsold", and each bucket individually can never go below
-- zero (CHECK), so no combination of bucket updates can oversell the event.
--
-- This is the sharded-counter technique. R1 serialises every buyer of an event
-- on one row; R2 lets up to 8 of them proceed at once, each on a different
-- bucket, picked at random with FOR UPDATE SKIP LOCKED.
--
-- What it costs:
--   * Availability is a SUM over 8 rows instead of one read.
--   * Buckets drain unevenly. Near sell-out, a buyer can find every bucket with
--     seats left locked by someone else. As in P2, "nothing unlocked" is not
--     "sold out", so the booking path carries an exact fallback check.
--   * A refund must find a bucket that has room to take the seat back.

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

CREATE TABLE event_inventory_bucket (
    event_id         BIGINT NOT NULL REFERENCES event (event_id),
    bucket_no        INT NOT NULL,
    bucket_capacity  INT NOT NULL,
    remaining        INT NOT NULL,
    PRIMARY KEY (event_id, bucket_no),
    CONSTRAINT bucket_remaining_in_range CHECK (remaining BETWEEN 0 AND bucket_capacity)
);

CREATE TABLE ticket (
    ticket_id    BIGINT PRIMARY KEY,
    event_id     BIGINT NOT NULL REFERENCES event (event_id),
    seat_no      INT,
    customer_id  BIGINT NOT NULL,
    sold_at      TIMESTAMPTZ NOT NULL,
    price_cents  BIGINT NOT NULL
);
