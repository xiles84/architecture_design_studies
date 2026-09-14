-- S2 — lock-then-update: how to recompute the truth from this design's tables.
--
-- The harness compares these with its own ledger of what it saw commit (INV-1,
-- INV-5, INV-7) and with each other (derived-state drift). It runs them on
-- quiescent data, after a phase has stopped and the sweeper has finished.

-- name: a_tickets
-- params: none
-- Every ticket, for the per-seat ledger comparison.
SELECT event_id, seat_id, ticket_id, customer_id, COALESCE(hold_id, 0)::BIGINT AS hold_id
  FROM ticket;

-- name: a_inventory_sold
-- params: none
-- The inventory's own view of sold seats. Must agree with a_tickets seat for seat.
SELECT event_id, seat_id, COALESCE(hold_id, 0)::BIGINT AS hold_id
  FROM event_seat
 WHERE status = 'sold';

-- name: a_valid_holds
-- params: none
-- Holds valid now, judged by the database clock.
SELECT hold_id, event_id, seat_id, hold_expires_at AS expires_at
  FROM event_seat
 WHERE status = 'held'
   AND hold_expires_at > now();

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
-- Tickets for seats that do not exist in the event's venue.
SELECT t.event_id, t.seat_id
  FROM ticket t
  JOIN event e ON e.event_id = t.event_id
  LEFT JOIN venue_seat vs ON vs.venue_id = e.venue_id AND vs.seat_id = t.seat_id
 WHERE vs.seat_id IS NULL;
