-- L1 — claim-rows: the read catalogue.
--
-- Names and result shapes are identical in every design; verification compares
-- them with Go-computed truth. A free seat has no claim, so every "how many left"
-- is the venue's capacity minus the valid claims -- it grows with what is claimed,
-- not with the size of the event.

-- name: q01_section_map
-- params: event_id, section_no
-- The dynamic half of a seat map: non-available seats (1 = held, 2 = sold).
SELECT seat_id,
       (CASE WHEN expires_at = 'infinity' THEN 2 ELSE 1 END)::SMALLINT AS state
  FROM seat_claim
 WHERE event_id = $1
   AND section_no = $2
   AND expires_at > now()
 ORDER BY seat_id;

-- name: q02_event_sections
-- params: event_id
-- Seats available per section: capacity minus valid claims.
SELECT s.section_no,
       (s.capacity - COALESCE(c.claimed, 0))::BIGINT AS available
  FROM event e
  JOIN venue_section s ON s.venue_id = e.venue_id
  LEFT JOIN (SELECT section_no, COUNT(*) AS claimed
               FROM seat_claim
              WHERE event_id = $1
                AND expires_at > now()
              GROUP BY section_no) c ON c.section_no = s.section_no
 WHERE e.event_id = $1
 ORDER BY s.section_no;

-- name: q03_event_available
-- params: event_id
-- Seats available for the whole event.
SELECT (v.capacity - (SELECT COUNT(*)
                        FROM seat_claim
                       WHERE event_id = $1
                         AND expires_at > now()))::BIGINT AS available
  FROM event e
  JOIN venue v ON v.venue_id = e.venue_id
 WHERE e.event_id = $1;

-- name: q04_hold_seats
-- params: event_id, section_no, hold_id
-- The basket page: the seats of a hold that is still valid, and when it ends.
SELECT seat_id, expires_at
  FROM seat_claim
 WHERE event_id = $1
   AND section_no = $2
   AND hold_id = $3
   AND expires_at > now()
   AND expires_at <> 'infinity'
 ORDER BY seat_id;

-- name: q05_customer_tickets
-- params: customer_id
-- "My tickets". Same SQL as every other design; the ticket table's extra unique
-- constraint is not on this path.
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
-- a row only for a seat that is HELD or SOLD (seat_claim); a free seat has
-- no row at all, so "available" is absence, not a status value.
-- ---------------------------------------------------------------------------

-- name: r01_holds_expiring_soon
-- params: event_id, within
-- `within` is bound as an absolute cutoff instant (see event_seat's own copy
-- of this comment for why).
SELECT seat_id, hold_id, customer_id, expires_at AS hold_expires_at
  FROM seat_claim
 WHERE event_id = $1
   AND expires_at < 'infinity'
   AND expires_at <= $2
 ORDER BY expires_at;

-- name: r02_seat_status_lookup
-- params: event_id, seat_id
-- Zero rows means the seat is available -- there is no row to return one.
SELECT CASE WHEN expires_at = 'infinity' THEN 'sold' ELSE 'held' END AS status,
       hold_id, customer_id,
       CASE WHEN expires_at = 'infinity' THEN claimed_at ELSE expires_at END AS at
  FROM seat_claim
 WHERE event_id = $1 AND seat_id = $2;
