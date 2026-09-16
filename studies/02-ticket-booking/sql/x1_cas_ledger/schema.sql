-- X1 — cas-ledger  (REPORTS.md v2)
--
-- P3's exact schema, plus one append-only table: sale_event. Every design in
-- this study either deletes a cancelled ticket or resets it to available, so
-- the sale itself leaves no trace once undone -- REPORTS.md section 2's
-- finding. X1 answers r01 (correctly, across refunds), r03 and r05 from a
-- ledger written in the SAME transaction as the sale and the cancellation
-- (see writes.sql), instead of from the ticket row that transaction undoes.
--
-- P3 is the base because the study's own analysis names it one of the two
-- fastest correct designs for a hot drop: the ledger's cost is measured where
-- it would hurt most, not where a slower design's noise would hide it.

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
    -- Edited by the organiser while a sale is running (the race's "editor").
    -- In designs that keep a counter on this row, those edits contend with
    -- bookings; here they do not.
    description  TEXT NOT NULL
);

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
    -- Edited by the organiser while a sale is running (the race's "editor").
    -- In designs that keep a counter on this row, those edits contend with
    -- bookings; here they do not.
    description  TEXT NOT NULL
);

CREATE TABLE ticket (
    ticket_id    BIGINT PRIMARY KEY,
    event_id     BIGINT NOT NULL REFERENCES event (event_id),
    seat_no      INT NOT NULL,
    status       TEXT NOT NULL CHECK (status IN ('available', 'sold')),
    customer_id  BIGINT,
    sold_at      TIMESTAMPTZ,
    price_cents  BIGINT NOT NULL,
    -- A sold ticket has a buyer and an unsold one does not. Cheap, and it turns
    -- a half-applied sale or cancellation into an error instead of a quiet
    -- inconsistency.
    CONSTRAINT ticket_sold_has_buyer CHECK ((status = 'sold') = (customer_id IS NOT NULL))
);

-- The ledger. Append-only: a sale never overwrites or deletes a prior entry,
-- which is exactly the property the base design lacks. Keyed by (event_id,
-- seat_no) rather than ticket_id, because that pair is what identifies a SEAT
-- across its sold/cancelled/resold lifetime; ticket_id is stable too in this
-- design (pre-created), but the pair is what the reconciliation audit joins on.
CREATE TABLE sale_event (
    -- Database-generated: unlike every other id in this study, nothing else
    -- ever needs to name a specific sale_event row, so there is no reason to
    -- thread a generator through the booking path for it.
    sale_event_id  BIGSERIAL   PRIMARY KEY,
    event_id       BIGINT      NOT NULL REFERENCES event (event_id),
    seat_no        INT         NOT NULL,
    customer_id    BIGINT      NOT NULL,
    kind           TEXT        NOT NULL CHECK (kind IN ('sold', 'cancelled')),
    at             TIMESTAMPTZ NOT NULL,
    price_cents    BIGINT      NOT NULL
);
