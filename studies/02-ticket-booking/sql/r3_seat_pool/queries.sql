-- Read catalogue for R3: availability counts the pool, sales count tickets.

-- name: q01_event_availability
-- params: event_id
SELECT COUNT(*)::BIGINT AS seats_left
  FROM seat_slot
 WHERE event_id = $1;

-- name: q02_band_events
-- params: band_id
SELECT e.event_id,
       e.capacity,
       (SELECT COUNT(*) FROM seat_slot s WHERE s.event_id = e.event_id)::BIGINT AS seats_left
  FROM event e
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
SELECT COUNT(*)::BIGINT AS sold
  FROM ticket t
  JOIN event e ON e.event_id = t.event_id
 WHERE e.band_id = $1;
