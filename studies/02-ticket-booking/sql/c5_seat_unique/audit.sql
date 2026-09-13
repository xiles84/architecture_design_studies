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
-- C5's derived state is its seat-number frontier: MAX(seat_no) is what booking
-- believes has been sold. A frontier above the ticket count is a hole, which is
-- harmless until the frontier reaches capacity; from then on booking answers
-- "sold out" with seats empty. Only that harmful state is reported:
-- stored = frontier (= capacity), actual = tickets that exist.
SELECT t.event_id, MAX(t.seat_no)::BIGINT AS stored, COUNT(*)::BIGINT AS actual
  FROM ticket t
  JOIN event e ON e.event_id = t.event_id
 GROUP BY t.event_id, e.capacity
HAVING MAX(t.seat_no) >= e.capacity
   AND COUNT(*) < e.capacity;
