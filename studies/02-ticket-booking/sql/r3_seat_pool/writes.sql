-- Writes for R3 (seat-pool).
--
-- Booking, one READ COMMITTED transaction:
--   w_pick_slot      -> lowest unsold seat nobody has locked (SKIP LOCKED)
--   w_consume_slot   -> delete it from the pool
--   w_insert_ticket  -> create the ticket with that seat number
-- and, only when w_pick_slot returns nothing:
--   w_any_slot       -> sold out, or all locked right now? (retry if the latter)
--
-- Cancellation, one transaction:
--   w_cancel_ticket  -> delete the ticket, returning its seat
--   w_restore_slot   -> put the seat back in the pool

-- name: w_pick_slot
-- params: event_id
SELECT seat_no
  FROM seat_slot
 WHERE event_id = $1
 ORDER BY seat_no
 LIMIT 1
   FOR UPDATE SKIP LOCKED;

-- name: w_any_slot
-- params: event_id
SELECT EXISTS (SELECT 1 FROM seat_slot WHERE event_id = $1);

-- name: w_consume_slot
-- params: event_id, seat_no
DELETE FROM seat_slot WHERE event_id = $1 AND seat_no = $2;

-- name: w_insert_ticket
-- params: ticket_id, event_id, seat_no, customer_id, price_cents
INSERT INTO ticket (ticket_id, event_id, seat_no, customer_id, sold_at, price_cents)
VALUES ($1, $2, $3, $4, now(), $5);

-- name: w_cancel_ticket
-- params: ticket_id
DELETE FROM ticket WHERE ticket_id = $1
RETURNING event_id, seat_no;

-- name: w_restore_slot
-- params: event_id, seat_no
INSERT INTO seat_slot (event_id, seat_no) VALUES ($1, $2);

-- name: w_edit_event
-- params: description, event_id
UPDATE event SET description = $1 WHERE event_id = $2;

-- name: w_publish_event
-- params: event_id, band_id, name, venue, starts_at, capacity, price_cents, description
INSERT INTO event (event_id, band_id, name, venue, starts_at, capacity, price_cents, description)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: w_publish_slots
-- params: event_id, capacity
-- Pre-creation again, but of a two-column token rather than a full ticket.
INSERT INTO seat_slot (event_id, seat_no)
SELECT $1::BIGINT, g FROM generate_series(1, $2::INT) AS g;
