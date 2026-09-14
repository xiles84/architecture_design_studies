-- S2 — lock-then-update: the read catalogue.
--
-- Names and result shapes are identical in every design; verification compares
-- them with Go-computed truth. Expiry is judged at read time: a held row whose
-- hold_expires_at has passed is shown as available.

-- name: q01_section_map
-- params: event_id, section_no
-- The dynamic half of a seat map: the seats of one section that are NOT
-- available (1 = held, 2 = sold). The client merges it with the static layout.
SELECT seat_id,
       (CASE WHEN status = 'sold' THEN 2 ELSE 1 END)::SMALLINT AS state
  FROM event_seat
 WHERE event_id = $1
   AND section_no = $2
   AND (status = 'sold' OR (status = 'held' AND hold_expires_at > now()))
 ORDER BY seat_id;

-- name: q02_event_sections
-- params: event_id
-- Seats available per section: what a buyer looks at to pick a section. A count
-- over every seat row of the event, so it grows with the event.
SELECT section_no,
       COUNT(*) FILTER (WHERE status = 'available'
                           OR (status = 'held' AND hold_expires_at <= now()))::BIGINT AS available
  FROM event_seat
 WHERE event_id = $1
 GROUP BY section_no
 ORDER BY section_no;

-- name: q03_event_available
-- params: event_id
-- Seats available for the whole event.
SELECT COUNT(*)::BIGINT AS available
  FROM event_seat
 WHERE event_id = $1
   AND (status = 'available' OR (status = 'held' AND hold_expires_at <= now()));

-- name: q04_hold_seats
-- params: event_id, section_no, hold_id
-- The basket page: the seats of a hold that is still valid, and when it ends.
SELECT seat_id, hold_expires_at AS expires_at
  FROM event_seat
 WHERE event_id = $1
   AND section_no = $2
   AND hold_id = $3
   AND status = 'held'
   AND hold_expires_at > now()
 ORDER BY seat_id;

-- name: q05_customer_tickets
-- params: customer_id
-- "My tickets". Byte-identical in every design with the same ticket table: the
-- noise calibration set.
SELECT ticket_id, event_id, seat_id
  FROM ticket
 WHERE customer_id = $1
 ORDER BY ticket_id;

-- name: q06_ticket_by_id
-- params: ticket_id
-- Point lookup: the control.
SELECT ticket_id, event_id, seat_id, customer_id
  FROM ticket
 WHERE ticket_id = $1;
