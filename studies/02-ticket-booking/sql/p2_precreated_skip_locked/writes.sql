-- Writes for P2 (precreated-skip-locked). Differs from P1 in w_pick_seat's
-- locking clause, plus the fallback that clause makes necessary.
--
-- Booking, one READ COMMITTED transaction:
--   w_pick_seat       -> lock the lowest seat NOBODY ELSE HAS LOCKED
--   w_sell_seat       -> mark it sold
-- and, only when w_pick_seat returns nothing:
--   w_any_available   -> is anything left at all? (non-locking)
--                        no  -> sold out
--                        yes -> every remaining seat is mid-sale by someone else;
--                               try again, because some of those sales may roll back
--
-- Without the fallback, the last few seats of a hot event would be reported as
-- sold out to every buyer who arrived while they were locked -- an under-sale
-- that no overbooking check would ever notice. The race's "underbooked" column
-- exists to catch exactly that.

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

-- name: w_cancel_ticket
-- params: ticket_id
UPDATE ticket
   SET status = 'available', customer_id = NULL, sold_at = NULL
 WHERE ticket_id = $1
   AND status = 'sold'
RETURNING event_id, seat_no;

-- name: w_edit_event
-- params: description, event_id
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
