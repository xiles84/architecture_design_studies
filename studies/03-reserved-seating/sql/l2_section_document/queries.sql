-- L2 — section-document: the read catalogue.
--
-- Names and result shapes are identical in every design; verification compares
-- them with Go-computed truth. Every seat question is answered by opening section
-- documents: claims are [state, hold_id, customer_id, expires_ms], a held claim is
-- valid while expires_ms is later than now() in epoch milliseconds.

-- name: q01_section_map
-- params: event_id, section_no
-- The dynamic half of a seat map, from ONE row: non-available seats.
SELECT (c.key)::INT AS seat_id,
       ((c.value->>0)::INT)::SMALLINT AS state
  FROM event_section s
 CROSS JOIN LATERAL jsonb_each(s.claims) AS c
 WHERE s.event_id = $1
   AND s.section_no = $2
   AND ((c.value->>0)::INT = 2
        OR (c.value->>3)::BIGINT > (EXTRACT(EPOCH FROM now()) * 1000)::BIGINT)
 ORDER BY 1;

-- name: q02_event_sections
-- params: event_id
-- Seats available per section: capacity minus valid claims, one document each.
SELECT s.section_no,
       (s.capacity - (SELECT COUNT(*)
                        FROM jsonb_each(s.claims) AS c
                       WHERE (c.value->>0)::INT = 2
                          OR (c.value->>3)::BIGINT > (EXTRACT(EPOCH FROM now()) * 1000)::BIGINT))::BIGINT AS available
  FROM event_section s
 WHERE s.event_id = $1
 ORDER BY s.section_no;

-- name: q03_event_available
-- params: event_id
-- Seats available for the whole event: every section document of the event.
SELECT COALESCE(SUM(s.capacity - (SELECT COUNT(*)
                                    FROM jsonb_each(s.claims) AS c
                                   WHERE (c.value->>0)::INT = 2
                                      OR (c.value->>3)::BIGINT > (EXTRACT(EPOCH FROM now()) * 1000)::BIGINT)), 0)::BIGINT AS available
  FROM event_section s
 WHERE s.event_id = $1;

-- name: q04_hold_seats
-- params: event_id, section_no, hold_id
-- The basket page: the seats of a hold that is still valid, and when it ends.
SELECT (c.key)::INT AS seat_id,
       to_timestamp((c.value->>3)::BIGINT / 1000.0) AS expires_at
  FROM event_section s
 CROSS JOIN LATERAL jsonb_each(s.claims) AS c
 WHERE s.event_id = $1
   AND s.section_no = $2
   AND (c.value->>0)::INT = 1
   AND (c.value->>1)::BIGINT = $3
   AND (c.value->>3)::BIGINT > (EXTRACT(EPOCH FROM now()) * 1000)::BIGINT
 ORDER BY 1;

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
-- Operational reports (study 03 v2, REPORTS.md), r01-r02: this design has no
-- per-seat row at all -- claims live inside event_section.claims, one JSONB
-- document per section: {"<seat_id>": [state, hold_id, customer_id,
-- expires_ms], ...}, state 1 = held, 2 = sold. Both reports unnest it.
-- ---------------------------------------------------------------------------

-- name: r01_holds_expiring_soon
-- params: event_id, within
-- `within` is bound as an absolute cutoff instant (see event_seat's own copy
-- of this comment for why).
SELECT (kv.key)::INT AS seat_id,
       (kv.value ->> 1)::BIGINT AS hold_id,
       (kv.value ->> 2)::BIGINT AS customer_id,
       to_timestamp((kv.value ->> 3)::BIGINT / 1000.0) AS hold_expires_at
  FROM event_section es,
       LATERAL jsonb_each(es.claims) AS kv
 WHERE es.event_id = $1
   AND (kv.value ->> 0)::INT = 1
   AND to_timestamp((kv.value ->> 3)::BIGINT / 1000.0) <= $2
 ORDER BY 4;

-- name: r02_seat_status_lookup
-- params: event_id, seat_id
-- Zero rows means the seat is available -- the claims document has no key
-- for it. Reaches the seat's section through venue_seat, since a claim is
-- addressed by (event, section, seat-key) and this report is only given the
-- seat. `at` is NULL once sold: the document does not keep a sale timestamp,
-- only the (now stale) hold expiry -- an honest limitation of this design,
-- not a bug in the query.
SELECT CASE WHEN (es.claims -> ($2::text) ->> 0)::INT = 2 THEN 'sold' ELSE 'held' END AS status,
       (es.claims -> ($2::text) ->> 1)::BIGINT AS hold_id,
       (es.claims -> ($2::text) ->> 2)::BIGINT AS customer_id,
       CASE WHEN (es.claims -> ($2::text) ->> 0)::INT = 1
            THEN to_timestamp((es.claims -> ($2::text) ->> 3)::BIGINT / 1000.0)
       END AS at
  FROM event_section es
  JOIN event e ON e.event_id = es.event_id
  JOIN venue_seat vs ON vs.venue_id = e.venue_id AND vs.seat_id = $2
 WHERE es.event_id = $1
   AND es.section_no = vs.section_no
   AND es.claims ? ($2::text);
