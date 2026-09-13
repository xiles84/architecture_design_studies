-- Read catalogue for R2: availability is the sum of an event's buckets.

-- name: q01_event_availability
-- params: event_id
SELECT COALESCE(SUM(remaining), 0)::BIGINT AS seats_left
  FROM event_inventory_bucket
 WHERE event_id = $1;

-- name: q02_band_events
-- params: band_id
SELECT e.event_id,
       e.capacity,
       (SELECT COALESCE(SUM(b.remaining), 0) FROM event_inventory_bucket b
         WHERE b.event_id = e.event_id)::BIGINT AS seats_left
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
SELECT COALESCE(SUM(b.bucket_capacity - b.remaining), 0)::BIGINT AS sold
  FROM event e
  JOIN event_inventory_bucket b ON b.event_id = e.event_id
 WHERE e.band_id = $1;
