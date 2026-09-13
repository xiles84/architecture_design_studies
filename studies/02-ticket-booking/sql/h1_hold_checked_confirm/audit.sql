-- Audit for the hold designs (H0 and H1 share it).

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
 WHERE seat_no IS NOT NULL
 GROUP BY event_id, seat_no
HAVING COUNT(*) > 1;

-- name: a_derived_drift
-- params: none
-- seats_taken must equal tickets plus holds still in status 'held'. Run after a
-- phase has finished and the sweeper has stopped, so no hold is mid-release.
-- stored = seats_taken, actual = tickets + live holds.
SELECT e.event_id,
       e.seats_taken::BIGINT AS stored,
       (COALESCE(s.sold, 0) + COALESCE(h.held, 0))::BIGINT AS actual
  FROM event e
  LEFT JOIN (SELECT event_id, COUNT(*) AS sold FROM ticket GROUP BY event_id) s
         ON s.event_id = e.event_id
  LEFT JOIN (SELECT event_id, COUNT(*) AS held FROM reservation
              WHERE status = 'held' GROUP BY event_id) h
         ON h.event_id = e.event_id
 WHERE e.seats_taken <> COALESCE(s.sold, 0) + COALESCE(h.held, 0);
