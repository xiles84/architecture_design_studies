-- Writes for X1 (cas-ledger). Byte-identical to P3 except w_sell_seat and
-- w_cancel_ticket: each becomes a data-modifying CTE that writes sale_event in
-- the SAME statement as the ticket UPDATE, which on PostgreSQL and YugabyteDB
-- both is one implicit transaction -- "the same transaction as the sale"
-- (REPORTS.md section 3) without touching Booker.bookCAS or Booker.Cancel in
-- harness/booking.go at all. w_sell_seat is still executed with Exec (not
-- QueryRow): the driver reports the FINAL statement's row count, so "0 rows"
-- still means "someone else sold it first, retry" exactly as it did in P3 --
-- the CAS loop needs no changes either. w_cancel_ticket's final SELECT
-- returns exactly the (event_id, seat_no) shape Booker.Cancel already scans.
--
-- Booking:
--   w_candidate_from(start_seat)  -> first available seat at or after a random seat
--   w_candidate_from(1)           -> only if nothing was found after start_seat (wrap around)
--                                     nothing here either -> sold out
--   w_sell_seat                   -> compare-and-set + ledger append; 0 rows = lost the race, retry

-- name: w_candidate_from
-- params: event_id, start_seat
SELECT ticket_id
  FROM ticket
 WHERE event_id = $1
   AND status = 'available'
   AND seat_no >= $2
 ORDER BY seat_no
 LIMIT 1;

-- name: w_sell_seat
-- params: customer_id, ticket_id
WITH sold AS (
    UPDATE ticket
       SET status = 'sold', customer_id = $1, sold_at = now()
     WHERE ticket_id = $2
       AND status = 'available'
    RETURNING event_id, seat_no, sold_at, price_cents
)
INSERT INTO sale_event (event_id, seat_no, customer_id, kind, at, price_cents)
SELECT event_id, seat_no, $1, 'sold', sold_at, price_cents FROM sold;

-- name: w_cancel_ticket
-- params: ticket_id
WITH cancelled AS (
    UPDATE ticket
       SET status = 'available', customer_id = NULL, sold_at = NULL
     WHERE ticket_id = $1
       AND status = 'sold'
    RETURNING event_id, seat_no, customer_id, price_cents
), logged AS (
    INSERT INTO sale_event (event_id, seat_no, customer_id, kind, at, price_cents)
    SELECT event_id, seat_no, customer_id, 'cancelled', now(), price_cents FROM cancelled
    RETURNING event_id, seat_no
)
SELECT event_id, seat_no FROM logged;

-- name: w_edit_event
-- params: description, event_id
UPDATE event SET description = $1 WHERE event_id = $2;

-- name: w_publish_event
-- params: event_id, band_id, name, venue, starts_at, capacity, price_cents, description
INSERT INTO event (event_id, band_id, name, venue, starts_at, capacity, price_cents, description)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: w_publish_tickets
-- params: first_ticket_id, event_id, capacity, price_cents
INSERT INTO ticket (ticket_id, event_id, seat_no, status, price_cents)
SELECT $1::BIGINT + g - 1, $2::BIGINT, g, 'available', $4::BIGINT
  FROM generate_series(1, $3::INT) AS g;
