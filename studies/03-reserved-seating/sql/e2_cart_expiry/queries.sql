-- E2 — cart-expiry: the read catalogue.
--
-- Names and result shapes are identical in every design; verification compares
-- them with Go-computed truth. Expiry is judged at read time, as in S1 -- but it
-- lives on the cart, so every read that asks "held or free?" joins hold.

-- name: q01_section_map
-- params: event_id, section_no
-- The dynamic half of a seat map: non-available seats (1 = held, 2 = sold).
SELECT es.seat_id,
       (CASE WHEN es.status = 'sold' THEN 2 ELSE 1 END)::SMALLINT AS state
  FROM event_seat es
  LEFT JOIN hold h ON h.hold_id = es.hold_id
 WHERE es.event_id = $1
   AND es.section_no = $2
   AND (es.status = 'sold' OR (es.status = 'held' AND h.expires_at > now()))
 ORDER BY es.seat_id;

-- name: q02_event_sections
-- params: event_id
-- Seats available per section.
SELECT es.section_no,
       COUNT(*) FILTER (WHERE es.status = 'available'
                           OR (es.status = 'held' AND h.expires_at <= now()))::BIGINT AS available
  FROM event_seat es
  LEFT JOIN hold h ON h.hold_id = es.hold_id
 WHERE es.event_id = $1
 GROUP BY es.section_no
 ORDER BY es.section_no;

-- name: q03_event_available
-- params: event_id
-- Seats available for the whole event.
SELECT COUNT(*)::BIGINT AS available
  FROM event_seat es
  LEFT JOIN hold h ON h.hold_id = es.hold_id
 WHERE es.event_id = $1
   AND (es.status = 'available' OR (es.status = 'held' AND h.expires_at <= now()));

-- name: q04_hold_seats
-- params: event_id, section_no, hold_id
-- The basket page: the seats of a hold that is still valid, and when it ends.
SELECT es.seat_id, h.expires_at
  FROM event_seat es
  JOIN hold h ON h.hold_id = es.hold_id
 WHERE es.event_id = $1
   AND es.section_no = $2
   AND es.hold_id = $3
   AND es.status = 'held'
   AND h.expires_at > now()
 ORDER BY es.seat_id;

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

-- ---------------------------------------------------------------------------
-- Operational reports (study 03 v2, REPORTS.md). r03-r05 use the `ticket`
-- table, which is byte-identical across every design in this study (the
-- confirmed-sale record); this SQL is therefore the same for all of them.
-- ---------------------------------------------------------------------------

-- name: r03_section_sales_window
-- params: event_id, since, until
-- Confirmed sales per section inside a window -- selling pace, section by
-- section.
SELECT section_no, COUNT(*) AS confirmed, COALESCE(SUM(price_cents), 0) AS revenue_cents
  FROM ticket
 WHERE event_id = $1 AND sold_at >= $2 AND sold_at < $3
 GROUP BY section_no
 ORDER BY section_no;

-- name: r04_event_recent_confirmations
-- params: event_id
-- The operations feed: the last 50 confirmed sales, newest first.
SELECT ticket_id, seat_id, customer_id, sold_at
  FROM ticket
 WHERE event_id = $1
 ORDER BY sold_at DESC
 LIMIT 50;

-- name: r05_customers_last_purchase_window
-- params: since, until
-- Customers whose LAST purchase falls in the window -- the recency question
-- RECENCY.md section 2 asks of donors, in this domain. Same trap: in the
-- trailing regime "bought in the window" and "last bought in the window"
-- coincide (nothing is newer); the historical regime is where a design that
-- only checked the former would be wrong.
SELECT customer_id, MAX(sold_at) AS last_at
  FROM ticket
 GROUP BY customer_id
HAVING MAX(sold_at) >= $1 AND MAX(sold_at) < $2
 ORDER BY 2 DESC
 LIMIT 100;


-- ---------------------------------------------------------------------------
-- Operational reports (study 03 v2, REPORTS.md), r01-r02: this design's holds
-- are carts (one hold table row can cover several seats), so hold_expires_at
-- and the holding customer live on `hold`, not on event_seat -- unlike every
-- other design in this study, whose event_seat row carries its own expiry.
-- ---------------------------------------------------------------------------

-- name: r01_holds_expiring_soon
-- params: event_id, within
-- `within` is bound as an absolute cutoff instant (see event_seat's own copy
-- of this comment, in the designs that carry hold_expires_at directly).
SELECT es.seat_id, es.hold_id, h.customer_id, h.expires_at AS hold_expires_at
  FROM event_seat es
  JOIN hold h ON h.hold_id = es.hold_id
 WHERE es.event_id = $1
   AND es.status = 'held'
   AND h.expires_at <= $2
 ORDER BY h.expires_at;

-- name: r02_seat_status_lookup
-- params: event_id, seat_id
-- customer_id comes from the cart while held (event_seat's own customer_id
-- is NULL until sold, by this design's own CHECK constraint) and from
-- event_seat once sold; `at` follows the same split.
SELECT es.status,
       es.hold_id,
       COALESCE(es.customer_id, h.customer_id) AS customer_id,
       CASE WHEN es.status = 'sold' THEN es.sold_at ELSE h.expires_at END AS at
  FROM event_seat es
  LEFT JOIN hold h ON h.hold_id = es.hold_id
 WHERE es.event_id = $1 AND es.seat_id = $2;
