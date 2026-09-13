-- Read catalogue for R1: availability comes from the inventory row.

-- name: q01_event_availability
-- params: event_id
SELECT remaining::BIGINT AS seats_left
  FROM event_inventory
 WHERE event_id = $1;

-- name: q02_band_events
-- params: band_id
SELECT e.event_id, e.capacity, i.remaining::BIGINT AS seats_left
  FROM event e
  JOIN event_inventory i ON i.event_id = e.event_id
 WHERE e.band_id = $1
 ORDER BY e.starts_at, e.event_id;

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
SELECT COALESCE(SUM(e.capacity - i.remaining), 0)::BIGINT AS sold
  FROM event e
  JOIN event_inventory i ON i.event_id = e.event_id
 WHERE e.band_id = $1;
