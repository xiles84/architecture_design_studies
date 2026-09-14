-- L1 — claim-rows: the read catalogue.
--
-- Names and result shapes are identical in every design; verification compares
-- them with Go-computed truth. A free seat has no claim, so every "how many left"
-- is the venue's capacity minus the valid claims -- it grows with what is claimed,
-- not with the size of the event.

-- name: q01_section_map
-- params: event_id, section_no
-- The dynamic half of a seat map: non-available seats (1 = held, 2 = sold).
SELECT seat_id,
       (CASE WHEN expires_at = 'infinity' THEN 2 ELSE 1 END)::SMALLINT AS state
  FROM seat_claim
 WHERE event_id = $1
   AND section_no = $2
   AND expires_at > now()
 ORDER BY seat_id;

-- name: q02_event_sections
-- params: event_id
-- Seats available per section: capacity minus valid claims.
SELECT s.section_no,
       (s.capacity - COALESCE(c.claimed, 0))::BIGINT AS available
  FROM event e
  JOIN venue_section s ON s.venue_id = e.venue_id
  LEFT JOIN (SELECT section_no, COUNT(*) AS claimed
               FROM seat_claim
              WHERE event_id = $1
                AND expires_at > now()
              GROUP BY section_no) c ON c.section_no = s.section_no
 WHERE e.event_id = $1
 ORDER BY s.section_no;

-- name: q03_event_available
-- params: event_id
-- Seats available for the whole event.
SELECT (v.capacity - (SELECT COUNT(*)
                        FROM seat_claim
                       WHERE event_id = $1
                         AND expires_at > now()))::BIGINT AS available
  FROM event e
  JOIN venue v ON v.venue_id = e.venue_id
 WHERE e.event_id = $1;

-- name: q04_hold_seats
-- params: event_id, section_no, hold_id
-- The basket page: the seats of a hold that is still valid, and when it ends.
SELECT seat_id, expires_at
  FROM seat_claim
 WHERE event_id = $1
   AND section_no = $2
   AND hold_id = $3
   AND expires_at > now()
   AND expires_at <> 'infinity'
 ORDER BY seat_id;

-- name: q05_customer_tickets
-- params: customer_id
-- "My tickets". Same SQL as every other design; the ticket table's extra unique
-- constraint is not on this path.
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
