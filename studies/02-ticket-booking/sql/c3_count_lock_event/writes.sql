-- Writes for C3 (count-lock-event): C1's statements, preceded by a parent lock.
--
-- Booking, one READ COMMITTED transaction:
--   w_lock_event     -> SELECT ... FOR UPDATE on the event row   <- the only addition
--   w_seats_left     -> capacity minus tickets that exist (fresh snapshot, after the lock)
--                       <= 0 -> sold out
--   w_insert_ticket  -> create the ticket

-- name: w_lock_event
-- params: event_id
SELECT event_id FROM event WHERE event_id = $1 FOR UPDATE;

-- name: w_seats_left
-- params: event_id
SELECT (e.capacity - (SELECT COUNT(*) FROM ticket t WHERE t.event_id = e.event_id))::BIGINT AS seats_left
  FROM event e
 WHERE e.event_id = $1;

-- name: w_insert_ticket
-- params: ticket_id, event_id, seat_no, customer_id, price_cents
INSERT INTO ticket (ticket_id, event_id, seat_no, customer_id, sold_at, price_cents)
VALUES ($1, $2, $3, $4, now(), $5);

-- name: w_cancel_ticket
-- params: ticket_id
-- Cancellation needs no lock: removing a ticket can only make room.
DELETE FROM ticket WHERE ticket_id = $1
RETURNING event_id, seat_no;

-- name: w_edit_event
-- params: description, event_id
-- Blocks on the booking lock, and bookings block on it: the organiser's edit
-- joins the queue.
UPDATE event SET description = $1 WHERE event_id = $2;

-- name: w_publish_event
-- params: event_id, band_id, name, venue, starts_at, capacity, price_cents, description
INSERT INTO event (event_id, band_id, name, venue, starts_at, capacity, price_cents, description)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
