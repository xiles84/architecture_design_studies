-- Writes for P4 (precreated-counter): P2's booking path plus one statement.
--
-- Booking, one READ COMMITTED transaction:
--   w_pick_seat      -> FOR UPDATE SKIP LOCKED, as P2
--   w_sell_seat      -> as P2
--   w_count_sale     -> seats_sold + 1 on the event row   <- the only addition
-- with P2's exact fallback (w_any_available) when no unlocked seat is found.
--
-- Lock order is always ticket, then event -- in booking and in cancellation --
-- so the two paths cannot deadlock against each other.

-- name: w_pick_seat
-- params: event_id
SELECT ticket_id
  FROM ticket
 WHERE event_id = $1
   AND status = 'available'
 ORDER BY seat_no
 LIMIT 1
   FOR UPDATE SKIP LOCKED;

-- name: w_any_available
-- params: event_id
SELECT EXISTS (SELECT 1 FROM ticket WHERE event_id = $1 AND status = 'available');

-- name: w_sell_seat
-- params: customer_id, ticket_id
UPDATE ticket
   SET status = 'sold', customer_id = $1, sold_at = now()
 WHERE ticket_id = $2
   AND status = 'available';

-- name: w_count_sale
-- params: event_id
-- Every buyer of this event now waits here, one at a time, until the
-- transaction ahead of them commits.
UPDATE event SET seats_sold = seats_sold + 1 WHERE event_id = $1;

-- name: w_cancel_ticket
-- params: ticket_id
UPDATE ticket
   SET status = 'available', customer_id = NULL, sold_at = NULL
 WHERE ticket_id = $1
   AND status = 'sold'
RETURNING event_id, seat_no;

-- name: w_count_refund
-- params: event_id
UPDATE event SET seats_sold = seats_sold - 1 WHERE event_id = $1;

-- name: w_edit_event
-- params: description, event_id
-- Lands on the same row as w_count_sale: an organiser's edit queues behind the
-- sale in progress, and the sale behind the edit.
UPDATE event SET description = $1 WHERE event_id = $2;

-- name: w_publish_event
-- params: event_id, band_id, name, venue, starts_at, capacity, price_cents, description
INSERT INTO event (event_id, band_id, name, venue, starts_at, capacity, price_cents, description)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: w_publish_tickets
-- params: first_ticket_id, event_id, capacity, price_cents
INSERT INTO ticket (ticket_id, event_id, seat_no, status, price_cents)
SELECT $1::BIGINT + g - 1, $2::BIGINT, g, 'available', $4::BIGINT
  FROM generate_series(1, $3::INT) AS g;
