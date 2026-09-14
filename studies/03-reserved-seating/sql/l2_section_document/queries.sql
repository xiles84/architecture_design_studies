-- L2 — section-document: the read catalogue.
--
-- Names and result shapes are identical in every design; verification compares
-- them with Go-computed truth. Every seat question is answered by opening section
-- documents: claims are [state, hold_id, customer_id, expires_ms], a held claim is
-- valid while expires_ms is later than now() in epoch milliseconds.

-- name: q01_section_map
-- params: event_id, section_no
-- The dynamic half of a seat map, from ONE row: non-available seats.
SELECT (c.key)::INT AS seat_id,
       ((c.value->>0)::INT)::SMALLINT AS state
  FROM event_section s
 CROSS JOIN LATERAL jsonb_each(s.claims) AS c
 WHERE s.event_id = $1
   AND s.section_no = $2
   AND ((c.value->>0)::INT = 2
        OR (c.value->>3)::BIGINT > (EXTRACT(EPOCH FROM now()) * 1000)::BIGINT)
 ORDER BY 1;

-- name: q02_event_sections
-- params: event_id
-- Seats available per section: capacity minus valid claims, one document each.
SELECT s.section_no,
       (s.capacity - (SELECT COUNT(*)
                        FROM jsonb_each(s.claims) AS c
                       WHERE (c.value->>0)::INT = 2
                          OR (c.value->>3)::BIGINT > (EXTRACT(EPOCH FROM now()) * 1000)::BIGINT))::BIGINT AS available
  FROM event_section s
 WHERE s.event_id = $1
 ORDER BY s.section_no;

-- name: q03_event_available
-- params: event_id
-- Seats available for the whole event: every section document of the event.
SELECT COALESCE(SUM(s.capacity - (SELECT COUNT(*)
                                    FROM jsonb_each(s.claims) AS c
                                   WHERE (c.value->>0)::INT = 2
                                      OR (c.value->>3)::BIGINT > (EXTRACT(EPOCH FROM now()) * 1000)::BIGINT)), 0)::BIGINT AS available
  FROM event_section s
 WHERE s.event_id = $1;

-- name: q04_hold_seats
-- params: event_id, section_no, hold_id
-- The basket page: the seats of a hold that is still valid, and when it ends.
SELECT (c.key)::INT AS seat_id,
       to_timestamp((c.value->>3)::BIGINT / 1000.0) AS expires_at
  FROM event_section s
 CROSS JOIN LATERAL jsonb_each(s.claims) AS c
 WHERE s.event_id = $1
   AND s.section_no = $2
   AND (c.value->>0)::INT = 1
   AND (c.value->>1)::BIGINT = $3
   AND (c.value->>3)::BIGINT > (EXTRACT(EPOCH FROM now()) * 1000)::BIGINT
 ORDER BY 1;

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
