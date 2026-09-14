-- E2 — cart-expiry: how to recompute the truth from this design's tables.

-- name: a_tickets
-- params: none
SELECT event_id, seat_id, ticket_id, customer_id, COALESCE(hold_id, 0)::BIGINT AS hold_id
  FROM ticket;

-- name: a_inventory_sold
-- params: none
SELECT event_id, seat_id, COALESCE(hold_id, 0)::BIGINT AS hold_id
  FROM event_seat
 WHERE status = 'sold';

-- name: a_valid_holds
-- params: none
-- A seat is validly held when it points to a cart that is unexpired now.
SELECT es.hold_id, es.event_id, es.seat_id, h.expires_at
  FROM event_seat es
  JOIN hold h ON h.hold_id = es.hold_id
 WHERE es.status = 'held'
   AND h.expires_at > now();

-- name: a_event_sold
-- params: event_id
SELECT COUNT(*)::BIGINT FROM event_seat WHERE event_id = $1 AND status = 'sold';

-- name: a_inventory_drift
-- params: none
-- Every event must own exactly one seat row per venue seat.
SELECT e.event_id, v.capacity::BIGINT AS stored, COUNT(es.seat_id)::BIGINT AS actual
  FROM event e
  JOIN venue v ON v.venue_id = e.venue_id
  LEFT JOIN event_seat es ON es.event_id = e.event_id
 GROUP BY e.event_id, v.capacity
HAVING COUNT(es.seat_id) <> v.capacity;

-- name: a_invalid_seats
-- params: none
SELECT t.event_id, t.seat_id
  FROM ticket t
  JOIN event e ON e.event_id = t.event_id
  LEFT JOIN venue_seat vs ON vs.venue_id = e.venue_id AND vs.seat_id = t.seat_id
 WHERE vs.seat_id IS NULL;
