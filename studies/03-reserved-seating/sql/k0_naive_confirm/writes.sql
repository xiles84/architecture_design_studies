-- K0 — naive-confirm (NEGATIVE CONTROL): every statement the hold, checkout, sweeper, refund
-- and publishing flows issue.
--
-- Hold, one READ COMMITTED transaction:
--   w_hold_seats      -> fewer rows returned than seats requested = a seat was not
--                        free: ROLLBACK, re-read the seat map, choose again
--
-- Pre-payment check (autocommit, friendly, not a guarantee):
--   w_check_hold
--
-- Confirm, one transaction (the harness first reads now() for its ledger):
--   w_confirm_seats   -> NAIVE: marks the block sold without checking the hold
--                        (the negative control; S1 checks it here)
--   w_insert_tickets
--
-- Release by the buyer:        w_release_hold
-- Sweeper, one transaction:    w_expired_holds, then w_release_expired per section
-- Refund, one transaction:     w_cancel_ticket, then w_unsell_seat
-- Publish, one transaction:    w_publish_event, then w_publish_seats
--
-- Every statement that grants, extends, releases or sells returns now() AS db_now
-- (the transaction start: the clock the design's expiry logic uses) and
-- clock_timestamp() AS db_clock (when the row was actually written: the order of a
-- release and the grant that waited for it). The harness's ledger judges theft and
-- late sales by these database clocks.

-- name: w_hold_seats
-- params: hold_id, customer_id, hold_ms, event_id, section_no, seat_ids
-- The design. Takes every seat of the block that is available or whose hold has
-- expired; the caller rolls back unless every seat was taken.
UPDATE event_seat
   SET status = 'held',
       hold_id = $1,
       customer_id = $2,
       hold_expires_at = now() + make_interval(secs => $3::FLOAT8 / 1000.0)
 WHERE event_id = $4
   AND section_no = $5
   AND seat_id = ANY ($6::INT[])
   AND (status = 'available' OR (status = 'held' AND hold_expires_at <= now()))
RETURNING seat_id, now() AS db_now, clock_timestamp() AS db_clock, hold_expires_at AS expires_at;

-- name: w_check_hold
-- params: event_id, section_no, hold_id, seat_ids
-- "Is my hold still valid?" before paying. True when it runs; it can expire
-- during the payment that follows.
SELECT COUNT(*) = cardinality($4::INT[]) AS valid
  FROM event_seat
 WHERE event_id = $1
   AND section_no = $2
   AND seat_id = ANY ($4::INT[])
   AND hold_id = $3
   AND status = 'held'
   AND hold_expires_at > now();

-- name: w_confirm_seats
-- params: event_id, section_no, hold_id, seat_ids, customer_id
-- NAIVE: trusts the pre-payment check. If the hold expired during payment and
-- someone else holds or bought the seat, this overwrites them. It names the buyer
-- too, as code that "already checked" would.
UPDATE event_seat
   SET status = 'sold',
       hold_id = $3,
       customer_id = $5,
       sold_at = now()
 WHERE event_id = $1
   AND section_no = $2
   AND seat_id = ANY ($4::INT[])
RETURNING seat_id, now() AS db_now, clock_timestamp() AS db_clock;

-- name: w_insert_tickets
-- params: ticket_ids, event_id, section_no, seat_ids, customer_id, hold_id, price_cents
-- One ticket per seat; ticket_ids and seat_ids are parallel arrays.
INSERT INTO ticket (ticket_id, event_id, section_no, seat_id, customer_id, hold_id, sold_at, price_cents)
SELECT t.ticket_id, $2::BIGINT, $3::INT, t.seat_id, $5::BIGINT, $6::BIGINT, now(), $7::BIGINT
  FROM unnest($1::BIGINT[], $4::INT[]) AS t (ticket_id, seat_id);

-- name: w_release_hold
-- params: event_id, section_no, hold_id, seat_ids
-- The buyer removes the seats. Releases only what this hold still holds.
UPDATE event_seat
   SET status = 'available', hold_id = NULL, customer_id = NULL, hold_expires_at = NULL
 WHERE event_id = $1
   AND section_no = $2
   AND seat_id = ANY ($4::INT[])
   AND hold_id = $3
   AND status = 'held'
RETURNING seat_id, now() AS db_now, clock_timestamp() AS db_clock;

-- name: w_expired_holds
-- params: batch
-- The sweeper's batch: expired holds in expiry order, skipping rows a buyer is
-- touching. Here it only tidies up -- expiry already took effect lazily.
SELECT event_id, section_no, seat_id, hold_id, hold_expires_at AS expires_at
  FROM event_seat
 WHERE status = 'held'
   AND hold_expires_at <= now()
 ORDER BY hold_expires_at
 LIMIT $1
   FOR UPDATE SKIP LOCKED;

-- name: w_release_expired
-- params: event_id, section_no, seat_ids
-- Guarded again by expiry: never releases a hold that is still valid.
UPDATE event_seat
   SET status = 'available', hold_id = NULL, customer_id = NULL, hold_expires_at = NULL
 WHERE event_id = $1
   AND section_no = $2
   AND seat_id = ANY ($3::INT[])
   AND status = 'held'
   AND hold_expires_at <= now()
RETURNING seat_id, now() AS db_now, clock_timestamp() AS db_clock;

-- name: w_cancel_ticket
-- params: ticket_id
-- A refund. No row means the ticket was already refunded.
DELETE FROM ticket
 WHERE ticket_id = $1
RETURNING event_id, section_no, seat_id, now() AS db_now, clock_timestamp() AS db_clock;

-- name: w_unsell_seat
-- params: event_id, section_no, seat_id
-- The refunded seat becomes available again.
UPDATE event_seat
   SET status = 'available', hold_id = NULL, customer_id = NULL, hold_expires_at = NULL, sold_at = NULL
 WHERE event_id = $1
   AND section_no = $2
   AND seat_id = $3
   AND status = 'sold';

-- name: w_publish_event
-- params: event_id, band_id, venue_id, name, starts_at, price_cents, description
INSERT INTO event (event_id, band_id, venue_id, name, starts_at, price_cents, description)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: w_publish_seats
-- params: event_id, venue_id
-- Pre-creation: one row per seat of the venue, at publishing time.
INSERT INTO event_seat (event_id, seat_id, section_no, status)
SELECT $1::BIGINT, seat_id, section_no, 'available'
  FROM venue_seat
 WHERE venue_id = $2::BIGINT;
