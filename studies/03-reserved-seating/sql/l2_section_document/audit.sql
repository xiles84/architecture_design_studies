-- L2 — section-document: how to recompute the truth from this design's tables.

-- name: a_tickets
-- params: none
SELECT event_id, seat_id, ticket_id, customer_id, COALESCE(hold_id, 0)::BIGINT AS hold_id
  FROM ticket;

-- name: a_inventory_sold
-- params: none
-- Sold claims in every document. Must agree with a_tickets seat for seat.
SELECT s.event_id, (c.key)::INT AS seat_id, COALESCE((c.value->>1)::BIGINT, 0) AS hold_id
  FROM event_section s
 CROSS JOIN LATERAL jsonb_each(s.claims) AS c
 WHERE (c.value->>0)::INT = 2;

-- name: a_valid_holds
-- params: none
SELECT (c.value->>1)::BIGINT AS hold_id, s.event_id, (c.key)::INT AS seat_id,
       to_timestamp((c.value->>3)::BIGINT / 1000.0) AS expires_at
  FROM event_section s
 CROSS JOIN LATERAL jsonb_each(s.claims) AS c
 WHERE (c.value->>0)::INT = 1
   AND (c.value->>3)::BIGINT > (EXTRACT(EPOCH FROM now()) * 1000)::BIGINT;

-- name: a_event_sold
-- params: event_id
SELECT COUNT(*)::BIGINT
  FROM event_section s
 CROSS JOIN LATERAL jsonb_each(s.claims) AS c
 WHERE s.event_id = $1
   AND (c.value->>0)::INT = 2;

-- name: a_inventory_drift
-- params: none
-- A claim for a seat that is not in that section of the event's venue, or a
-- section whose capacity disagrees with the venue.
SELECT s.event_id, s.section_no::BIGINT AS stored, (c.key)::BIGINT AS actual
  FROM event_section s
  JOIN event e ON e.event_id = s.event_id
 CROSS JOIN LATERAL jsonb_each(s.claims) AS c
  LEFT JOIN venue_seat vs ON vs.venue_id = e.venue_id AND vs.seat_id = (c.key)::INT
 WHERE vs.seat_id IS NULL OR vs.section_no <> s.section_no
UNION ALL
SELECT s.event_id, s.capacity::BIGINT, vsec.capacity::BIGINT
  FROM event_section s
  JOIN event e ON e.event_id = s.event_id
  JOIN venue_section vsec ON vsec.venue_id = e.venue_id AND vsec.section_no = s.section_no
 WHERE vsec.capacity <> s.capacity;

-- name: a_invalid_seats
-- params: none
SELECT t.event_id, t.seat_id
  FROM ticket t
  JOIN event e ON e.event_id = t.event_id
  LEFT JOIN venue_seat vs ON vs.venue_id = e.venue_id AND vs.seat_id = t.seat_id
 WHERE vs.seat_id IS NULL;
