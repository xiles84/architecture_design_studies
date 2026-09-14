-- E2 — cart-expiry: the read catalogue.
--
-- Names and result shapes are identical in every design; verification compares
-- them with Go-computed truth. Expiry is judged at read time, as in S1 -- but it
-- lives on the cart, so every read that asks "held or free?" joins hold.

-- name: q01_section_map
-- params: event_id, section_no
-- The dynamic half of a seat map: non-available seats (1 = held, 2 = sold).
SELECT es.seat_id,
       (CASE WHEN es.status = 'sold' THEN 2 ELSE 1 END)::SMALLINT AS state
  FROM event_seat es
  LEFT JOIN hold h ON h.hold_id = es.hold_id
 WHERE es.event_id = $1
   AND es.section_no = $2
   AND (es.status = 'sold' OR (es.status = 'held' AND h.expires_at > now()))
 ORDER BY es.seat_id;

-- name: q02_event_sections
-- params: event_id
-- Seats available per section.
SELECT es.section_no,
       COUNT(*) FILTER (WHERE es.status = 'available'
                           OR (es.status = 'held' AND h.expires_at <= now()))::BIGINT AS available
  FROM event_seat es
  LEFT JOIN hold h ON h.hold_id = es.hold_id
 WHERE es.event_id = $1
 GROUP BY es.section_no
 ORDER BY es.section_no;

-- name: q03_event_available
-- params: event_id
-- Seats available for the whole event.
SELECT COUNT(*)::BIGINT AS available
  FROM event_seat es
  LEFT JOIN hold h ON h.hold_id = es.hold_id
 WHERE es.event_id = $1
   AND (es.status = 'available' OR (es.status = 'held' AND h.expires_at <= now()));

-- name: q04_hold_seats
-- params: event_id, section_no, hold_id
-- The basket page: the seats of a hold that is still valid, and when it ends.
SELECT es.seat_id, h.expires_at
  FROM event_seat es
  JOIN hold h ON h.hold_id = es.hold_id
 WHERE es.event_id = $1
   AND es.section_no = $2
   AND es.hold_id = $3
   AND es.status = 'held'
   AND h.expires_at > now()
 ORDER BY es.seat_id;

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
