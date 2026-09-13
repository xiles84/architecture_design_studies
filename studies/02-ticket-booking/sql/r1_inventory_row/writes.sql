-- Writes for R1 (inventory-row). C4's flow, on a different row.
--
-- Booking, one READ COMMITTED transaction:
--   w_take_seat      -> guarded decrement on event_inventory; no row = sold out
--   w_insert_ticket  -> create the ticket

-- name: w_take_seat
-- params: event_id
UPDATE event_inventory
   SET remaining = remaining - 1
 WHERE event_id = $1
   AND remaining > 0
RETURNING remaining;

-- name: w_insert_ticket
-- params: ticket_id, event_id, seat_no, customer_id, price_cents
INSERT INTO ticket (ticket_id, event_id, seat_no, customer_id, sold_at, price_cents)
VALUES ($1, $2, $3, $4, now(), $5);

-- name: w_cancel_ticket
-- params: ticket_id
DELETE FROM ticket WHERE ticket_id = $1
RETURNING event_id, seat_no;

-- name: w_release_seat
-- params: event_id
UPDATE event_inventory SET remaining = remaining + 1 WHERE event_id = $1;

-- name: w_edit_event
-- params: description, event_id
-- A different row from every sale: no contention with booking.
UPDATE event SET description = $1 WHERE event_id = $2;

-- name: w_publish_event
-- params: event_id, band_id, name, venue, starts_at, capacity, price_cents, description
INSERT INTO event (event_id, band_id, name, venue, starts_at, capacity, price_cents, description)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: w_publish_inventory
-- params: event_id, capacity
INSERT INTO event_inventory (event_id, remaining) VALUES ($1, $2);
