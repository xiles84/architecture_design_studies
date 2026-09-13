-- Read catalogue for the hold designs (H0 and H1 share it). seats_taken counts
-- holds as well as sales, so "seats left" means "seats a new buyer can hold" --
-- which is what an event page should show. On quiescent data (no holds) it
-- equals capacity minus tickets, and that is what verification checks.

-- name: q01_event_availability
-- params: event_id
SELECT (capacity - seats_taken)::BIGINT AS seats_left
  FROM event
 WHERE event_id = $1;

-- name: q02_band_events
-- params: band_id
SELECT event_id, capacity, (capacity - seats_taken)::BIGINT AS seats_left
  FROM event
 WHERE band_id = $1
 ORDER BY starts_at, event_id;

-- name: q03_customer_tickets
-- params: customer_id
SELECT ticket_id, event_id, seat_no
  FROM ticket
 WHERE customer_id = $1
 ORDER BY ticket_id;

-- name: q04_ticket_by_id
-- params: ticket_id
SELECT ticket_id, event_id, seat_no, customer_id
  FROM ticket
 WHERE ticket_id = $1;

-- name: q05_band_tickets_sold
-- params: band_id
-- Counts tickets, not seats_taken: a hold is not a sale.
SELECT COUNT(*)::BIGINT AS sold
  FROM ticket t
  JOIN event e ON e.event_id = t.event_id
 WHERE e.band_id = $1;
