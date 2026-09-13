-- Audit for designs that create tickets on booking. A ticket that exists is a
-- seat that is sold.

-- name: a_sold_per_event
-- params: none
SELECT e.event_id,
       e.capacity,
       COALESCE(s.sold, 0)::BIGINT AS sold
  FROM event e
  LEFT JOIN (SELECT event_id, COUNT(*) AS sold FROM ticket GROUP BY event_id) s
         ON s.event_id = e.event_id;

-- name: a_event_sold
-- params: event_id
SELECT COUNT(*)::BIGINT FROM ticket WHERE event_id = $1;

-- name: a_duplicate_seats
-- params: none
-- General-admission bookings carry no seat number, so only numbered tickets
-- can collide.
SELECT event_id, seat_no, COUNT(*)::BIGINT AS n
  FROM ticket
 WHERE seat_no IS NOT NULL
 GROUP BY event_id, seat_no
HAVING COUNT(*) > 1;

-- name: a_derived_drift
-- params: none
-- Seats the buckets say are sold must equal the tickets that exist, and the
-- buckets must add up to the event's capacity. stored = sold per buckets.
SELECT e.event_id, b.sold::BIGINT AS stored, COALESCE(s.sold, 0)::BIGINT AS actual
  FROM event e
  JOIN (SELECT event_id,
               SUM(bucket_capacity - remaining) AS sold,
               SUM(bucket_capacity) AS cap
          FROM event_inventory_bucket GROUP BY event_id) b ON b.event_id = e.event_id
  LEFT JOIN (SELECT event_id, COUNT(*) AS sold FROM ticket GROUP BY event_id) s
         ON s.event_id = e.event_id
 WHERE b.sold <> COALESCE(s.sold, 0)
    OR b.cap <> e.capacity;
