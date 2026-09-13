-- Writes for P1 (precreated-lock-first).
--
-- Booking, one READ COMMITTED transaction:
--   w_pick_seat  -> lock the lowest available seat (waits if another buyer holds it)
--   w_sell_seat  -> mark it sold
--
-- Why it cannot overbook: the row is locked before it is sold, and PostgreSQL
-- re-checks `status = 'available'` on the latest row version after waiting for
-- the lock. A seat sold by the previous lock holder no longer matches, and the
-- scan moves on to the next one. No row returned means no seat is left.
--
-- Why it convoys: every buyer asks for the SAME lowest seat. They queue on it,
-- then on the next, one at a time -- the event sells at the speed of one buyer
-- no matter how many seats remain.

-- name: w_pick_seat
-- params: event_id
SELECT ticket_id
  FROM ticket
 WHERE event_id = $1
   AND status = 'available'
 ORDER BY seat_no
 LIMIT 1
   FOR UPDATE;

-- name: w_sell_seat
-- params: customer_id, ticket_id
-- The status guard is redundant while the lock holds, and deliberately kept: if
-- it ever matches zero rows, the harness reports "a locked seat was no longer
-- available" as an error rather than selling it twice.
UPDATE ticket
   SET status = 'sold', customer_id = $1, sold_at = now()
 WHERE ticket_id = $2
   AND status = 'available';

-- name: w_cancel_ticket
-- params: ticket_id
-- A refund returns the seat to inventory. The row stays; only its state changes.
UPDATE ticket
   SET status = 'available', customer_id = NULL, sold_at = NULL
 WHERE ticket_id = $1
   AND status = 'sold'
RETURNING event_id, seat_no;

-- name: w_edit_event
-- params: description, event_id
-- The organiser edits the event page during the sale.
UPDATE event SET description = $1 WHERE event_id = $2;

-- name: w_publish_event
-- params: event_id, band_id, name, venue, starts_at, capacity, price_cents, description
INSERT INTO event (event_id, band_id, name, venue, starts_at, capacity, price_cents, description)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: w_publish_tickets
-- params: first_ticket_id, event_id, capacity, price_cents
-- The price of pre-creation, paid once per event: one row per seat, up front.
-- A 100 000-seat event is a 100 000-row insert before a single ticket is sold.
INSERT INTO ticket (ticket_id, event_id, seat_no, status, price_cents)
SELECT $1::BIGINT + g - 1, $2::BIGINT, g, 'available', $4::BIGINT
  FROM generate_series(1, $3::INT) AS g;
