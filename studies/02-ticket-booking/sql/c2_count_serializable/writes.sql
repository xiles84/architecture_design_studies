-- Writes for C1 (count-naive) and C2 (count-serializable) -- identical SQL.
--
-- Booking, one transaction (READ COMMITTED in C1, SERIALIZABLE in C2):
--   w_seats_left     -> capacity minus tickets that exist
--                       <= 0 -> sold out
--   w_insert_ticket  -> create the ticket
--
-- The gap between those two statements is one network round trip. That is the
-- whole race window in C1, and it is enough.

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
-- A refund deletes the ticket; the seat is available again because it is no
-- longer counted.
DELETE FROM ticket WHERE ticket_id = $1
RETURNING event_id, seat_no;

-- name: w_edit_event
-- params: description, event_id
UPDATE event SET description = $1 WHERE event_id = $2;

-- name: w_publish_event
-- params: event_id, band_id, name, venue, starts_at, capacity, price_cents, description
-- Publishing creates one row, whatever the capacity.
INSERT INTO event (event_id, band_id, name, venue, starts_at, capacity, price_cents, description)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
