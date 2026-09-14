-- L1 — claim-rows: how to recompute the truth from this design's tables.

-- name: a_tickets
-- params: none
SELECT event_id, seat_id, ticket_id, customer_id, COALESCE(hold_id, 0)::BIGINT AS hold_id
  FROM ticket;

-- name: a_inventory_sold
-- params: none
-- Sold claims. Must agree with a_tickets seat for seat.
SELECT event_id, seat_id, COALESCE(hold_id, 0)::BIGINT AS hold_id
  FROM seat_claim
 WHERE expires_at = 'infinity';

-- name: a_valid_holds
-- params: none
SELECT hold_id, event_id, seat_id, expires_at
  FROM seat_claim
 WHERE expires_at > now()
   AND expires_at <> 'infinity';

-- name: a_event_sold
-- params: event_id
SELECT COUNT(*)::BIGINT FROM seat_claim WHERE event_id = $1 AND expires_at = 'infinity';

-- name: a_inventory_drift
-- params: none
-- Claims for seats outside their event's venue, or in the wrong section.
SELECT c.event_id, c.section_no::BIGINT AS stored, COALESCE(vs.section_no, -1)::BIGINT AS actual
  FROM seat_claim c
  JOIN event e ON e.event_id = c.event_id
  LEFT JOIN venue_seat vs ON vs.venue_id = e.venue_id AND vs.seat_id = c.seat_id
 WHERE vs.seat_id IS NULL OR vs.section_no <> c.section_no;

-- name: a_invalid_seats
-- params: none
SELECT t.event_id, t.seat_id
  FROM ticket t
  JOIN event e ON e.event_id = t.event_id
  LEFT JOIN venue_seat vs ON vs.venue_id = e.venue_id AND vs.seat_id = t.seat_id
 WHERE vs.seat_id IS NULL;
