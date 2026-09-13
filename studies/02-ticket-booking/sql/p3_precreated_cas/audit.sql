-- Audit for the pre-created family. The harness compares these recomputations
-- with the capacity of every event and with its own count of the bookings and
-- cancellations it saw succeed.

-- name: a_sold_per_event
-- params: none
-- Every event, with the number of seats that are actually sold.
SELECT e.event_id,
       e.capacity,
       COALESCE(s.sold, 0)::BIGINT AS sold
  FROM event e
  LEFT JOIN (SELECT event_id, COUNT(*) AS sold
               FROM ticket WHERE status = 'sold'
              GROUP BY event_id) s ON s.event_id = e.event_id;

-- name: a_event_sold
-- params: event_id
SELECT COUNT(*)::BIGINT FROM ticket WHERE event_id = $1 AND status = 'sold';

-- name: a_duplicate_seats
-- params: none
-- A seat sold to two buyers. Structurally impossible here (one row per seat),
-- and checked anyway: an audit that only checks what a design could plausibly
-- get wrong is testing the auditor's assumptions, not the design.
SELECT event_id, seat_no, COUNT(*)::BIGINT AS n
  FROM ticket
 WHERE status = 'sold'
 GROUP BY event_id, seat_no
HAVING COUNT(*) > 1;

-- name: a_derived_drift
-- params: none
-- Pre-created inventory with no counter has no derived state that can drift.
-- Every event must still own exactly `capacity` seat rows.
SELECT e.event_id, e.capacity::BIGINT AS stored, COUNT(t.ticket_id)::BIGINT AS actual
  FROM event e
  LEFT JOIN ticket t ON t.event_id = e.event_id
 GROUP BY e.event_id, e.capacity
HAVING COUNT(t.ticket_id) <> e.capacity;
