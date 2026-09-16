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

-- ---------------------------------------------------------------------------
-- Operational reports (study 02 v2, REPORTS.md). The back office's questions,
-- not the buyer's -- and they land on the same tables the overbooking design
-- reshaped.
--
-- This design pre-creates a ticket row per SEAT, sold or not: every report
-- below must filter on status = 'sold', or it would count empty seats as sales.
-- ---------------------------------------------------------------------------

-- name: r01_event_sales_window
-- params: event_id, since, until
-- How the drop went: tickets sold and revenue inside a window. PARTIAL: a
-- ticket sold in this window and later refunded is NOT counted, because
-- cancellation resets the row to available (REPORTS.md section 2) and this
-- design keeps no record of the sale having happened.
SELECT COUNT(*) AS tickets_sold, COALESCE(SUM(price_cents), 0) AS revenue_cents
  FROM ticket
 WHERE event_id = $1 AND status = 'sold' AND sold_at >= $2 AND sold_at < $3;

-- name: r02_event_recent_buyers
-- params: event_id
-- The operations feed: the last 50 sales with the buyer, newest first.
SELECT ticket_id, customer_id, sold_at
  FROM ticket
 WHERE event_id = $1 AND status = 'sold'
 ORDER BY sold_at DESC
 LIMIT 50;

-- name: r03_customers_last_purchase_window
-- params: since, until
-- Customers whose LAST purchase falls in the window -- the recency question
-- RECENCY.md section 2 asks of donors, in this domain. The trap is the same:
-- in the trailing regime "bought in the window" and "last bought in the
-- window" happen to coincide (nothing is newer); the historical regime is
-- where a design that only checked "bought in the window" would be wrong.
SELECT t.customer_id, MAX(t.sold_at) AS last_at
  FROM ticket t
 WHERE t.status = 'sold'
 GROUP BY t.customer_id
HAVING MAX(t.sold_at) >= $1 AND MAX(t.sold_at) < $2
 ORDER BY 2 DESC
 LIMIT 100;

-- name: r04_band_sellthrough
-- params: band_id
-- The management report that decides the next tour: capacity, sold and the
-- percentage, per event.
SELECT e.event_id, e.capacity,
       COUNT(t.ticket_id) FILTER (WHERE t.status = 'sold') AS sold_count,
       ROUND(100.0 * COUNT(t.ticket_id) FILTER (WHERE t.status = 'sold') / e.capacity, 1) AS sold_pct
  FROM event e
  LEFT JOIN ticket t ON t.event_id = e.event_id
 WHERE e.band_id = $1
 GROUP BY e.event_id, e.capacity
 ORDER BY e.event_id;
