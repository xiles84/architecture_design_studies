-- K1 — payment-window: the read catalogue.
--
-- Names and result shapes are identical in every design; verification compares
-- them with Go-computed truth. Expiry is judged at read time: a held row whose
-- hold_expires_at has passed is shown as available.

-- name: q01_section_map
-- params: event_id, section_no
-- The dynamic half of a seat map: the seats of one section that are NOT
-- available (1 = held, 2 = sold). The client merges it with the static layout.
SELECT seat_id,
       (CASE WHEN status = 'sold' THEN 2 ELSE 1 END)::SMALLINT AS state
  FROM event_seat
 WHERE event_id = $1
   AND section_no = $2
   AND (status = 'sold' OR (status = 'held' AND hold_expires_at > now()))
 ORDER BY seat_id;

-- name: q02_event_sections
-- params: event_id
-- Seats available per section: what a buyer looks at to pick a section. A count
-- over every seat row of the event, so it grows with the event.
SELECT section_no,
       COUNT(*) FILTER (WHERE status = 'available'
                           OR (status = 'held' AND hold_expires_at <= now()))::BIGINT AS available
  FROM event_seat
 WHERE event_id = $1
 GROUP BY section_no
 ORDER BY section_no;

-- name: q03_event_available
-- params: event_id
-- Seats available for the whole event.
SELECT COUNT(*)::BIGINT AS available
  FROM event_seat
 WHERE event_id = $1
   AND (status = 'available' OR (status = 'held' AND hold_expires_at <= now()));

-- name: q04_hold_seats
-- params: event_id, section_no, hold_id
-- The basket page: the seats of a hold that is still valid, and when it ends.
SELECT seat_id, hold_expires_at AS expires_at
  FROM event_seat
 WHERE event_id = $1
   AND section_no = $2
   AND hold_id = $3
   AND status = 'held'
   AND hold_expires_at > now()
 ORDER BY seat_id;

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
-- Operational reports (study 03 v2, REPORTS.md), r01-r02: this design keeps
-- one row per seat per event (event_seat), so both are a direct filter on it.
-- ---------------------------------------------------------------------------

-- name: r01_holds_expiring_soon
-- params: event_id, within
-- Which holds expire soon -- the operations desk's view, and what actually
-- sizes a sweeper. `within` is bound as an absolute cutoff instant (now +
-- the desk's lookahead), computed by the harness rather than in SQL, so the
-- comparison needs no engine-specific INTERVAL arithmetic.
SELECT seat_id, hold_id, customer_id, hold_expires_at
  FROM event_seat
 WHERE event_id = $1
   AND status = 'held'
   AND hold_expires_at <= $2
 ORDER BY hold_expires_at;

-- name: r02_seat_status_lookup
-- params: event_id, seat_id
-- Who holds or owns this seat right now -- the call the box office makes
-- while a customer is on the phone. `at` is when the current state started
-- (hold expiry if held, sale time if sold).
SELECT status, hold_id, customer_id, COALESCE(hold_expires_at, sold_at) AS at
  FROM event_seat
 WHERE event_id = $1 AND seat_id = $2;
