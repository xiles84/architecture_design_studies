-- L2 — section-document: every statement the hold, checkout, sweeper, refund and
-- publishing flows issue. The decisions live in the application: these statements
-- only read a document and write one back if its version has not moved.
--
-- Hold (autocommit):
--   w_read_section    -> version, claims, now()
--   (Go: every seat of the block absent or expired? no: choose again)
--   w_write_section   -> 0 rows = someone wrote the section first: read again
--
-- Pre-payment check: w_read_section, judged in Go.
--
-- Confirm, one transaction:
--   w_read_section (autocommit, before) -> the hold still valid by its now()?
--   w_write_section   -> claims of the hold become sold; 0 rows: ROLLBACK, read again
--   w_insert_tickets
--
-- Release by the buyer:     w_read_section, w_write_section
-- Sweeper, one transaction: w_expired_sections (locks the rows), then for each
--                           section w_read_section and w_write_section
-- Refund, one transaction:  w_cancel_ticket, then w_read_section and w_write_section
-- Publish, one transaction: w_publish_event, then w_publish_sections

-- name: w_read_section
-- params: event_id, section_no
-- The document as text, parsed in Go, with the version to write against and the
-- database clock every expiry decision is made with.
SELECT version, claims::TEXT AS claims, now() AS db_now, clock_timestamp() AS db_clock
  FROM event_section
 WHERE event_id = $1
   AND section_no = $2;

-- name: w_write_section
-- params: claims, next_expiry, event_id, section_no, version
-- Compare-and-set on the whole document.
UPDATE event_section
   SET claims = $1::JSONB,
       next_expiry = $2::TIMESTAMPTZ,
       version = version + 1
 WHERE event_id = $3
   AND section_no = $4
   AND version = $5
RETURNING now() AS db_now, clock_timestamp() AS db_clock;

-- name: w_insert_tickets
-- params: ticket_ids, event_id, section_no, seat_ids, customer_id, hold_id, price_cents
-- One ticket per seat; ticket_ids and seat_ids are parallel arrays.
INSERT INTO ticket (ticket_id, event_id, section_no, seat_id, customer_id, hold_id, sold_at, price_cents)
SELECT t.ticket_id, $2::BIGINT, $3::INT, t.seat_id, $5::BIGINT, $6::BIGINT, now(), $7::BIGINT
  FROM unnest($1::BIGINT[], $4::INT[]) AS t (ticket_id, seat_id);

-- name: w_expired_sections
-- params: batch
-- Sections holding at least one expired claim, locked for the sweeper and skipped
-- if a buyer is writing them.
SELECT event_id, section_no
  FROM event_section
 WHERE next_expiry <= now()
 ORDER BY next_expiry
 LIMIT $1
   FOR UPDATE SKIP LOCKED;

-- name: w_cancel_ticket
-- params: ticket_id
-- A refund. No row means the ticket was already refunded.
DELETE FROM ticket
 WHERE ticket_id = $1
RETURNING event_id, section_no, seat_id, now() AS db_now, clock_timestamp() AS db_clock;

-- name: w_publish_event
-- params: event_id, band_id, venue_id, name, starts_at, price_cents, description
INSERT INTO event (event_id, band_id, venue_id, name, starts_at, price_cents, description)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: w_publish_sections
-- params: event_id, venue_id
-- One empty document per section of the venue.
INSERT INTO event_section (event_id, section_no, capacity, version, claims, next_expiry)
SELECT $1::BIGINT, section_no, capacity, 0, '{}'::JSONB, NULL
  FROM venue_section
 WHERE venue_id = $2::BIGINT;
