-- Read catalogue for the counter designs (C4, H0, H1 share it; the hold designs
-- name the column seats_taken, see their own copy).

-- name: q01_event_availability
-- params: event_id
SELECT (capacity - seats_sold)::BIGINT AS seats_left
  FROM event
 WHERE event_id = $1;

-- name: q02_band_events
-- params: band_id
SELECT event_id, capacity, (capacity - seats_sold)::BIGINT AS seats_left
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
SELECT COALESCE(SUM(seats_sold), 0)::BIGINT AS sold
  FROM event
 WHERE band_id = $1;
