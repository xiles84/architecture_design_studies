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
-- The counter must equal the seats actually sold. (Seat rows per event are
-- checked by P1-P3's audit; here the derived column is the thing that can drift.)
SELECT e.event_id, e.seats_sold::BIGINT AS stored, COALESCE(s.sold, 0)::BIGINT AS actual
  FROM event e
  LEFT JOIN (SELECT event_id, COUNT(*) AS sold
               FROM ticket WHERE status = 'sold'
              GROUP BY event_id) s ON s.event_id = e.event_id
 WHERE e.seats_sold <> COALESCE(s.sold, 0);
