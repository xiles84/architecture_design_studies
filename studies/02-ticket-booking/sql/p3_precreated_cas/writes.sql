-- Writes for P3 (precreated-cas). No explicit transaction on the booking path:
-- each statement autocommits.
--
-- Booking:
--   w_candidate_from(start_seat)  -> first available seat at or after a random seat
--   w_candidate_from(1)           -> only if nothing was found after start_seat (wrap around)
--                                     nothing here either -> sold out
--   w_sell_seat                   -> compare-and-set; 0 rows = lost the race, retry
--
-- The candidate read takes no lock, so "sold out" here is a statement about a
-- committed snapshot -- which is exactly what "sold out" means. The CAS is what
-- makes it impossible to sell a seat twice: two buyers can read the same
-- candidate, but only one UPDATE can move it from 'available' to 'sold'; the
-- other re-checks the condition on the new row version and matches nothing.

-- name: w_candidate_from
-- params: event_id, start_seat
SELECT ticket_id
  FROM ticket
 WHERE event_id = $1
   AND status = 'available'
   AND seat_no >= $2
 ORDER BY seat_no
 LIMIT 1;

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
