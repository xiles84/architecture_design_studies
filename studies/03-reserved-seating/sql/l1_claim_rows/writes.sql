-- L1 — claim-rows: every statement the hold, checkout, sweeper, refund and
-- publishing flows issue.
--
-- Hold, one READ COMMITTED transaction:
--   w_claim_seats     -> insert or steal a claim per seat; fewer rows than seats:
--                        ROLLBACK, re-read the seat map, choose again
--
-- Pre-payment check (autocommit):  w_check_hold
--
-- Confirm, one transaction (the harness first reads now() for its ledger):
--   w_confirm_seats   -> claims still this hold's and unexpired NOW become sold
--   w_insert_tickets  -> UNIQUE (event_id, seat_id): a second sale is refused here
--
-- Release by the buyer:        w_release_hold
-- Sweeper, one transaction:    w_expired_holds, then w_release_expired per section
-- Refund, one transaction:     w_cancel_ticket, then w_unsell_seat
-- Publish:                     w_publish_event -- nothing per seat
-- Load (harness, once):        w_load_sold_claims

-- name: w_claim_seats
-- params: event_id, section_no, hold_id, customer_id, hold_ms, seat_ids
-- The design. No row: inserted. An expired claim: stolen by the ON CONFLICT branch.
-- A valid claim (or a sale, expires_at = 'infinity'): the conflict branch updates
-- nothing and the seat is missing from RETURNING.
INSERT INTO seat_claim AS c (event_id, seat_id, section_no, hold_id, customer_id, expires_at, claimed_at)
SELECT $1::BIGINT, s, $2::INT, $3::BIGINT, $4::BIGINT,
       now() + make_interval(secs => $5::FLOAT8 / 1000.0), now()
  FROM unnest($6::INT[]) AS s
ON CONFLICT (event_id, seat_id) DO UPDATE
   SET hold_id = EXCLUDED.hold_id,
       customer_id = EXCLUDED.customer_id,
       expires_at = EXCLUDED.expires_at,
       claimed_at = EXCLUDED.claimed_at
 WHERE c.expires_at <= now()
RETURNING seat_id, now() AS db_now, clock_timestamp() AS db_clock, expires_at;

-- name: w_check_hold
-- params: event_id, section_no, hold_id, seat_ids
-- "Is my hold still valid?" before paying.
SELECT COUNT(*) = cardinality($4::INT[]) AS valid
  FROM seat_claim
 WHERE event_id = $1
   AND section_no = $2
   AND seat_id = ANY ($4::INT[])
   AND hold_id = $3
   AND expires_at > now()
   AND expires_at <> 'infinity';

-- name: w_confirm_seats
-- params: event_id, section_no, hold_id, seat_ids
-- Checked confirmation: a claim of this hold, unexpired now, becomes a sale.
UPDATE seat_claim
   SET expires_at = 'infinity'
 WHERE event_id = $1
   AND section_no = $2
   AND seat_id = ANY ($4::INT[])
   AND hold_id = $3
   AND expires_at > now()
   AND expires_at <> 'infinity'
RETURNING seat_id, now() AS db_now, clock_timestamp() AS db_clock;

-- name: w_insert_tickets
-- params: ticket_ids, event_id, section_no, seat_ids, customer_id, hold_id, price_cents
-- One ticket per seat; ticket_ids and seat_ids are parallel arrays.
INSERT INTO ticket (ticket_id, event_id, section_no, seat_id, customer_id, hold_id, sold_at, price_cents)
SELECT t.ticket_id, $2::BIGINT, $3::INT, t.seat_id, $5::BIGINT, $6::BIGINT, now(), $7::BIGINT
  FROM unnest($1::BIGINT[], $4::INT[]) AS t (ticket_id, seat_id);

-- name: w_release_hold
-- params: event_id, section_no, hold_id, seat_ids
-- The buyer removes the seats: the claims go.
DELETE FROM seat_claim
 WHERE event_id = $1
   AND section_no = $2
   AND seat_id = ANY ($4::INT[])
   AND hold_id = $3
   AND expires_at <> 'infinity'
RETURNING seat_id, now() AS db_now, clock_timestamp() AS db_clock;

-- name: w_expired_holds
-- params: batch
-- The sweeper's batch: expired claims in expiry order. Only tidying -- an expired
-- claim is already stealable.
SELECT event_id, section_no, seat_id, hold_id, expires_at
  FROM seat_claim
 WHERE expires_at <= now()
 ORDER BY expires_at
 LIMIT $1
   FOR UPDATE SKIP LOCKED;

-- name: w_release_expired
-- params: event_id, section_no, seat_ids
-- Guarded again by expiry: never deletes a claim that is still valid.
DELETE FROM seat_claim
 WHERE event_id = $1
   AND section_no = $2
   AND seat_id = ANY ($3::INT[])
   AND expires_at <= now()
RETURNING seat_id, now() AS db_now, clock_timestamp() AS db_clock;

-- name: w_cancel_ticket
-- params: ticket_id
-- A refund. No row means the ticket was already refunded.
DELETE FROM ticket
 WHERE ticket_id = $1
RETURNING event_id, section_no, seat_id, now() AS db_now, clock_timestamp() AS db_clock;

-- name: w_unsell_seat
-- params: event_id, section_no, seat_id
-- The refunded seat becomes free: its sold claim goes.
DELETE FROM seat_claim
 WHERE event_id = $1
   AND section_no = $2
   AND seat_id = $3
   AND expires_at = 'infinity';

-- name: w_publish_event
-- params: event_id, band_id, venue_id, name, starts_at, price_cents, description
-- Publishing is one row: the venue already defines the seats.
INSERT INTO event (event_id, band_id, venue_id, name, starts_at, price_cents, description)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: w_load_sold_claims
-- params: none
-- Run once by the loader, after the tickets are copied: every loaded sale gets its
-- sold claim. Done in SQL so the 'infinity' marker is written by the database.
INSERT INTO seat_claim (event_id, seat_id, section_no, hold_id, customer_id, expires_at, claimed_at)
SELECT event_id, seat_id, section_no, hold_id, customer_id, 'infinity', sold_at
  FROM ticket;
