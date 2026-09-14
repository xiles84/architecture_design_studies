-- E2 — cart-expiry: every statement the hold, checkout, sweeper, refund and
-- publishing flows issue.
--
-- Hold, one READ COMMITTED transaction:
--   w_insert_cart     -> the cart, with its expiry (written once)
--   w_hold_seats      -> S1's conditional update, expiry judged through the cart;
--                        fewer rows than seats: ROLLBACK (the cart goes too),
--                        re-read the seat map, choose again
--
-- Pre-payment check (autocommit):  w_check_hold
--
-- Confirm, one transaction (the harness first reads now() for its ledger):
--   w_confirm_cart    -> the cart must still be held and unexpired NOW (1 row)
--   w_confirm_seats   -> its seats must still point to it (N rows)
--   w_insert_tickets
--
-- Release by the buyer, one transaction: w_release_cart, w_release_hold
-- Sweeper, one transaction per batch:    w_expired_carts, w_release_cart_seats,
--                                        w_release_carts -- cart by cart
-- Refund, one transaction:               w_cancel_ticket, w_unsell_seat
-- Publish, one transaction:              w_publish_event, w_publish_seats

-- name: w_insert_cart
-- params: hold_id, customer_id, event_id, section_no, hold_ms
-- The only place the expiry is written.
INSERT INTO hold (hold_id, customer_id, event_id, section_no, status, expires_at, created_at)
VALUES ($1, $2, $3, $4, 'held', now() + make_interval(secs => $5::FLOAT8 / 1000.0), now())
RETURNING now() AS db_now, clock_timestamp() AS db_clock, expires_at;

-- name: w_hold_seats
-- params: hold_id, event_id, section_no, seat_ids
-- The design. A seat is free if available, or if the cart holding it has expired.
-- Correct only because a cart's expires_at is immutable (see schema.sql).
UPDATE event_seat es
   SET status = 'held',
       hold_id = $1
 WHERE es.event_id = $2
   AND es.section_no = $3
   AND es.seat_id = ANY ($4::INT[])
   AND (es.status = 'available'
        OR (es.status = 'held'
            AND EXISTS (SELECT 1 FROM hold h
                         WHERE h.hold_id = es.hold_id AND h.expires_at <= now())))
RETURNING es.seat_id, now() AS db_now, clock_timestamp() AS db_clock;

-- name: w_check_hold
-- params: event_id, section_no, hold_id, seat_ids
-- "Is my hold still valid?" before paying: the cart unexpired, every seat still its.
SELECT (SELECT COUNT(*)
          FROM event_seat
         WHERE event_id = $1
           AND section_no = $2
           AND seat_id = ANY ($4::INT[])
           AND hold_id = $3
           AND status = 'held') = cardinality($4::INT[])
   AND EXISTS (SELECT 1 FROM hold WHERE hold_id = $3 AND status = 'held' AND expires_at > now())
       AS valid;

-- name: w_confirm_cart
-- params: hold_id
-- Checked confirmation, on the one row that carries the expiry.
UPDATE hold
   SET status = 'confirmed'
 WHERE hold_id = $1
   AND status = 'held'
   AND expires_at > now()
RETURNING now() AS db_now, clock_timestamp() AS db_clock;

-- name: w_confirm_seats
-- params: event_id, section_no, hold_id, seat_ids, customer_id
-- The seats must still belong to the cart: a steal after the cart expired would
-- have repointed them.
UPDATE event_seat
   SET status = 'sold',
       customer_id = $5,
       sold_at = now()
 WHERE event_id = $1
   AND section_no = $2
   AND seat_id = ANY ($4::INT[])
   AND hold_id = $3
   AND status = 'held'
RETURNING seat_id, now() AS db_now, clock_timestamp() AS db_clock;

-- name: w_insert_tickets
-- params: ticket_ids, event_id, section_no, seat_ids, customer_id, hold_id, price_cents
-- One ticket per seat; ticket_ids and seat_ids are parallel arrays.
INSERT INTO ticket (ticket_id, event_id, section_no, seat_id, customer_id, hold_id, sold_at, price_cents)
SELECT t.ticket_id, $2::BIGINT, $3::INT, t.seat_id, $5::BIGINT, $6::BIGINT, now(), $7::BIGINT
  FROM unnest($1::BIGINT[], $4::INT[]) AS t (ticket_id, seat_id);

-- name: w_release_cart
-- params: hold_id
-- The buyer abandons the cart: one row.
UPDATE hold
   SET status = 'released'
 WHERE hold_id = $1
   AND status = 'held';

-- name: w_release_hold
-- params: event_id, section_no, hold_id, seat_ids
-- ...and gives back whatever seats still point to it.
UPDATE event_seat
   SET status = 'available', hold_id = NULL
 WHERE event_id = $1
   AND section_no = $2
   AND seat_id = ANY ($4::INT[])
   AND hold_id = $3
   AND status = 'held'
RETURNING seat_id, now() AS db_now, clock_timestamp() AS db_clock;

-- name: w_expired_carts
-- params: batch
-- The sweeper's batch, cart by cart, skipping carts a buyer is confirming. Only
-- tidying: expiry already took effect lazily.
SELECT hold_id, event_id, expires_at
  FROM hold
 WHERE status = 'held'
   AND expires_at <= now()
 ORDER BY expires_at
 LIMIT $1
   FOR UPDATE SKIP LOCKED;

-- name: w_release_cart_seats
-- params: hold_ids
-- The seats still pointing to those expired carts. A seat already stolen points to
-- its new cart and is left alone.
UPDATE event_seat
   SET status = 'available', hold_id = NULL
 WHERE hold_id = ANY ($1::BIGINT[])
   AND status = 'held'
RETURNING event_id, seat_id, now() AS db_now, clock_timestamp() AS db_clock;

-- name: w_release_carts
-- params: hold_ids
UPDATE hold
   SET status = 'released'
 WHERE hold_id = ANY ($1::BIGINT[])
   AND status = 'held'
   AND expires_at <= now();

-- name: w_cancel_ticket
-- params: ticket_id
-- A refund. No row means the ticket was already refunded.
DELETE FROM ticket
 WHERE ticket_id = $1
RETURNING event_id, section_no, seat_id, now() AS db_now, clock_timestamp() AS db_clock;

-- name: w_unsell_seat
-- params: event_id, section_no, seat_id
-- The refunded seat becomes available again.
UPDATE event_seat
   SET status = 'available', hold_id = NULL, customer_id = NULL, sold_at = NULL
 WHERE event_id = $1
   AND section_no = $2
   AND seat_id = $3
   AND status = 'sold';

-- name: w_publish_event
-- params: event_id, band_id, venue_id, name, starts_at, price_cents, description
INSERT INTO event (event_id, band_id, venue_id, name, starts_at, price_cents, description)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: w_publish_seats
-- params: event_id, venue_id
-- Pre-creation: one row per seat of the venue, at publishing time.
INSERT INTO event_seat (event_id, seat_id, section_no, status)
SELECT $1::BIGINT, seat_id, section_no, 'available'
  FROM venue_seat
 WHERE venue_id = $2::BIGINT;
