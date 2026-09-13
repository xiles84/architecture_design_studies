-- Writes for C4 (counter-guard).
--
-- Booking, one READ COMMITTED transaction:
--   w_take_seat      -> guarded increment on the event row; no row = sold out
--   w_insert_ticket  -> create the ticket
--
-- Cancellation, one transaction:
--   w_cancel_ticket  -> delete the ticket
--   w_release_seat   -> decrement the counter

-- name: w_take_seat
-- params: event_id
UPDATE event
   SET seats_sold = seats_sold + 1
 WHERE event_id = $1
   AND seats_sold < capacity
RETURNING seats_sold;

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
UPDATE event SET seats_sold = seats_sold - 1 WHERE event_id = $1;

-- name: w_edit_event
-- params: description, event_id
-- Same row as the counter. See R1, which moves the counter off this row.
UPDATE event SET description = $1 WHERE event_id = $2;

-- name: w_publish_event
-- params: event_id, band_id, name, venue, starts_at, capacity, price_cents, description
INSERT INTO event (event_id, band_id, name, venue, starts_at, capacity, price_cents, description)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
