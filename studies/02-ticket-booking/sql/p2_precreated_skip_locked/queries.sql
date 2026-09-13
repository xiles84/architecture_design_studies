-- Read catalogue for the pre-created family (P1, P2, P3 share it byte for byte).
--
-- Availability is a COUNT over the event's unsold rows. Its cost therefore
-- grows with event size -- q01 is benchmarked separately for every size tier,
-- which is where that shows.

-- name: q01_event_availability
-- params: event_id
-- How many seats are left for this event.
SELECT COUNT(*)::BIGINT AS seats_left
  FROM ticket
 WHERE event_id = $1
   AND status = 'available';

-- name: q02_band_events
-- params: band_id
-- A band's event list with live availability: one count per event.
SELECT e.event_id,
       e.capacity,
       (SELECT COUNT(*) FROM ticket t
         WHERE t.event_id = e.event_id AND t.status = 'available')::BIGINT AS seats_left
  FROM event e
 WHERE e.band_id = $1
 ORDER BY e.starts_at, e.event_id;

-- name: q03_customer_tickets
-- params: customer_id
-- "My tickets".
SELECT ticket_id, event_id, seat_no
  FROM ticket
 WHERE customer_id = $1
   AND status = 'sold'
 ORDER BY ticket_id;

-- name: q04_ticket_by_id
-- params: ticket_id
-- Point lookup: the control.
SELECT ticket_id, event_id, seat_no, customer_id
  FROM ticket
 WHERE ticket_id = $1;

-- name: q05_band_tickets_sold
-- params: band_id
-- Tickets sold across all of a band's events: aggregate up the whole tree.
SELECT COUNT(*)::BIGINT AS sold
  FROM ticket t
  JOIN event e ON e.event_id = t.event_id
 WHERE e.band_id = $1
   AND t.status = 'sold';
