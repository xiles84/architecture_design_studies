-- Writes for C5 (seat-unique). No explicit transaction: each statement autocommits.
--
-- Booking:
--   w_highest_seat  -> highest seat number sold for the event
--                      >= capacity -> sold out
--   w_claim_seat    -> insert seat highest+1; ON CONFLICT DO NOTHING
--                      0 rows -> another buyer took that number: retry from the top

-- name: w_highest_seat
-- params: event_id
SELECT COALESCE(MAX(seat_no), 0)::INT FROM ticket WHERE event_id = $1;

-- name: w_claim_seat
-- params: ticket_id, seat_no, customer_id, event_id
-- Capacity and price come from the event row inside the statement, so the
-- CHECK is enforced against the database's own copy of the capacity.
INSERT INTO ticket (ticket_id, event_id, seat_no, event_capacity, customer_id, sold_at, price_cents)
SELECT $1::BIGINT, e.event_id, $2::INT, e.capacity, $3::BIGINT, now(), e.price_cents
  FROM event e
 WHERE e.event_id = $4
ON CONFLICT (event_id, seat_no) DO NOTHING;

-- name: w_cancel_ticket
-- params: ticket_id
-- Leaves a hole in the seat numbers that booking will never revisit.
DELETE FROM ticket WHERE ticket_id = $1
RETURNING event_id, seat_no;

-- name: w_edit_event
-- params: description, event_id
UPDATE event SET description = $1 WHERE event_id = $2;

-- name: w_publish_event
-- params: event_id, band_id, name, venue, starts_at, capacity, price_cents, description
INSERT INTO event (event_id, band_id, name, venue, starts_at, capacity, price_cents, description)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
