-- Audit for R3 (seat-pool).

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
SELECT event_id, seat_no, COUNT(*)::BIGINT AS n
  FROM ticket
 GROUP BY event_id, seat_no
HAVING COUNT(*) > 1;

-- name: a_derived_drift
-- params: none
-- The pool and the tickets must partition the venue: every seat is either a
-- slot or a ticket, never both, never neither. stored = slots + tickets,
-- actual = capacity; a seat present in both counts once extra.
SELECT e.event_id,
       (COALESCE(p.slots, 0) + COALESCE(s.sold, 0))::BIGINT AS stored,
       e.capacity::BIGINT AS actual
  FROM event e
  LEFT JOIN (SELECT event_id, COUNT(*) AS slots FROM seat_slot GROUP BY event_id) p
         ON p.event_id = e.event_id
  LEFT JOIN (SELECT event_id, COUNT(*) AS sold FROM ticket GROUP BY event_id) s
         ON s.event_id = e.event_id
  LEFT JOIN (SELECT x.event_id, COUNT(*) AS both_places
               FROM seat_slot x
               JOIN ticket t ON t.event_id = x.event_id AND t.seat_no = x.seat_no
              GROUP BY x.event_id) d
         ON d.event_id = e.event_id
 WHERE COALESCE(p.slots, 0) + COALESCE(s.sold, 0) <> e.capacity
    OR COALESCE(d.both_places, 0) > 0;
