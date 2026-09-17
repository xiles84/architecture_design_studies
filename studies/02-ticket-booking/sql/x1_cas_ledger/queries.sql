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
-- Operational reports (study 02 v2, REPORTS.md). r02 and r04 are P3's own
-- formulation, unchanged (the ticket table is unchanged). r01, r03 and r05 are
-- what the ledger buys: answered from sale_event instead of from a ticket row
-- that a cancellation would otherwise have erased.
-- ---------------------------------------------------------------------------

-- name: r01_event_sales_window
-- params: event_id, since, until
-- How the drop went, CORRECTLY across refunds: a sale that happened in this
-- window and was later cancelled still counts as having happened then -- what
-- a finance report means by "sold in this period", and what P3 (see
-- p3_precreated_cas) cannot say.
SELECT COUNT(*) AS tickets_sold, COALESCE(SUM(price_cents), 0) AS revenue_cents
  FROM sale_event
 WHERE event_id = $1 AND kind = 'sold' AND at >= $2 AND at < $3;

-- name: r02_event_recent_buyers
-- params: event_id
SELECT ticket_id, customer_id, sold_at
  FROM ticket
 WHERE event_id = $1 AND status = 'sold'
 ORDER BY sold_at DESC
 LIMIT 50;

-- name: r03_customers_last_purchase_window
-- params: since, until
-- Answered from the ledger's 'sold' events, not the live ticket table
-- (EH-02 AM-03.6 corrects an earlier version of this comment: cancelling a
-- DIFFERENT, earlier ticket changes nothing on either formulation, since
-- neither one's last-purchase date depended on that ticket). What the ledger
-- buys is the case where the customer's MOST RECENT purchase is the one that
-- gets refunded: the ticket table then dates them by an earlier purchase, or
-- drops them if that was their only one, while the ledger still shows the
-- true last purchase -- which is why every non-ledger design's r03 is
-- `partial`, with exactly this caveat (REPORTS.md section 2).
SELECT customer_id, MAX(at) AS last_at
  FROM sale_event
 WHERE kind = 'sold'
 GROUP BY customer_id
HAVING MAX(at) >= $1 AND MAX(at) < $2
 ORDER BY 2 DESC
 LIMIT 100;

-- name: r04_band_sellthrough
-- params: band_id
SELECT e.event_id, e.capacity,
       COUNT(t.ticket_id) FILTER (WHERE t.status = 'sold') AS sold_count,
       ROUND(100.0 * COUNT(t.ticket_id) FILTER (WHERE t.status = 'sold') / e.capacity, 1) AS sold_pct
  FROM event e
  LEFT JOIN ticket t ON t.event_id = e.event_id
 WHERE e.band_id = $1
 GROUP BY e.event_id, e.capacity
 ORDER BY e.event_id;

-- name: r05_refunds_window
-- params: since, until
-- What the ledger buys that no other design in this study can answer at all:
-- cancellations and the money returned, inside a window.
SELECT COUNT(*) AS refunds, COALESCE(SUM(price_cents), 0) AS refunded_cents
  FROM sale_event
 WHERE kind = 'cancelled' AND at >= $1 AND at < $2;
