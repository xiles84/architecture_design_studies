-- Read catalogue for R2: availability is the sum of an event's buckets.

-- name: q01_event_availability
-- params: event_id
SELECT COALESCE(SUM(remaining), 0)::BIGINT AS seats_left
  FROM event_inventory_bucket
 WHERE event_id = $1;

-- name: q02_band_events
-- params: band_id
SELECT e.event_id,
       e.capacity,
       (SELECT COALESCE(SUM(b.remaining), 0) FROM event_inventory_bucket b
         WHERE b.event_id = e.event_id)::BIGINT AS seats_left
  FROM event e
 WHERE e.band_id = $1
 ORDER BY e.starts_at, e.event_id;

-- name: q03_customer_tickets
-- params: customer_id
SELECT ticket_id, event_id, seat_no
  FROM ticket
 WHERE customer_id = $1
 ORDER BY ticket_id;

-- name: q04_ticket_by_id
-- params: ticket_id
SELECT ticket_id, event_id, seat_no, customer_id
  FROM ticket
 WHERE ticket_id = $1;

-- name: q05_band_tickets_sold
-- params: band_id
SELECT COALESCE(SUM(b.bucket_capacity - b.remaining), 0)::BIGINT AS sold
  FROM event e
  JOIN event_inventory_bucket b ON b.event_id = e.event_id
 WHERE e.band_id = $1;

-- ---------------------------------------------------------------------------
-- Operational reports (study 02 v2, REPORTS.md). The back office's questions,
-- not the buyer's -- and they land on the same tables the overbooking design
-- reshaped.
--
-- This design creates a ticket row only when a seat is sold, so every row in
-- `ticket` already represents a sale: no status filter is needed here.
-- ---------------------------------------------------------------------------

-- name: r01_event_sales_window
-- params: event_id, since, until
-- How the drop went: tickets sold and revenue inside a window. PARTIAL: a
-- ticket sold in this window and later refunded is NOT counted, because
-- cancellation erases the sale (REPORTS.md section 2) and this design keeps
-- no record of it.
SELECT COUNT(*) AS tickets_sold, COALESCE(SUM(price_cents), 0) AS revenue_cents
  FROM ticket
 WHERE event_id = $1 AND sold_at >= $2 AND sold_at < $3;

-- name: r02_event_recent_buyers
-- params: event_id
-- The operations feed: the last 50 sales with the buyer, newest first.
SELECT ticket_id, customer_id, sold_at
  FROM ticket
 WHERE event_id = $1
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
 GROUP BY t.customer_id
HAVING MAX(t.sold_at) >= $1 AND MAX(t.sold_at) < $2
 ORDER BY 2 DESC
 LIMIT 100;

-- name: r04_band_sellthrough
-- params: band_id
-- The management report that decides the next tour: capacity, sold and the
-- percentage, per event.
SELECT e.event_id, e.capacity, COUNT(t.ticket_id) AS sold_count,
       ROUND(100.0 * COUNT(t.ticket_id) / e.capacity, 1) AS sold_pct
  FROM event e
  LEFT JOIN ticket t ON t.event_id = e.event_id
 WHERE e.band_id = $1
 GROUP BY e.event_id, e.capacity
 ORDER BY e.event_id;
