-- Writes for H0 (hold-naive-confirm). H1's file differs in w_confirm_hold only.
--
-- Hold, one READ COMMITTED transaction:
--   w_take_seat       -> guarded increment of seats_taken; no row = sold out
--   w_insert_hold     -> reservation row, status 'held', expiring hold_ms from now
--
-- Buyer thinks, then (if they did not abandon the basket):
--   w_check_hold      -> "is my hold still valid?" -- a friendly pre-payment check,
--                        identical in H0 and H1. It is NOT a guarantee: the hold
--                        can expire during the payment that follows.
--
-- Payment happens (outside the database), then confirm, one transaction:
--   w_confirm_hold    -> the decision this pair isolates
--                        0 rows -> the hold is gone: refund, no ticket
--   w_insert_ticket   -> the ticket
--
-- Sweeper, one transaction per batch:
--   w_expired_holds   -> expired holds, locked, skipping any a confirm is touching
--   w_release_hold    -> guarded: only a hold still 'held' is released
--   w_release_seat    -> give the seat back (only if the release matched a row)
--
-- An immediate purchase (the survey's and the race's "book") is a hold followed
-- at once by its confirmation: two transactions, no think time.

-- name: w_take_seat
-- params: event_id
UPDATE event
   SET seats_taken = seats_taken + 1
 WHERE event_id = $1
   AND seats_taken < capacity
RETURNING seats_taken;

-- name: w_insert_hold
-- params: reservation_id, event_id, customer_id, hold_ms
INSERT INTO reservation (reservation_id, event_id, customer_id, status, expires_at, created_at)
VALUES ($1, $2, $3, 'held', now() + make_interval(secs => $4::FLOAT8 / 1000.0), now());

-- name: w_check_hold
-- params: reservation_id
SELECT (status = 'held' AND expires_at > now()) AS valid
  FROM reservation
 WHERE reservation_id = $1;

-- name: w_confirm_hold
-- params: reservation_id
-- NAIVE: trusts the pre-payment w_check_hold. If the hold expired during payment
-- and the sweeper already gave the seat to someone else, this still confirms it
-- and the ticket below is a second sale of the same seat.
UPDATE reservation
   SET status = 'confirmed'
 WHERE reservation_id = $1;

-- name: w_insert_ticket
-- params: ticket_id, event_id, seat_no, customer_id, price_cents, reservation_id
INSERT INTO ticket (ticket_id, event_id, seat_no, customer_id, sold_at, price_cents, reservation_id)
VALUES ($1, $2, $3, $4, now(), $5, $6);

-- name: w_expired_holds
-- params: batch
SELECT reservation_id, event_id
  FROM reservation
 WHERE status = 'held'
   AND expires_at <= now()
 ORDER BY expires_at
 LIMIT $1
   FOR UPDATE SKIP LOCKED;

-- name: w_release_hold
-- params: reservation_id
UPDATE reservation
   SET status = 'released'
 WHERE reservation_id = $1
   AND status = 'held';

-- name: w_release_seat
-- params: event_id
UPDATE event SET seats_taken = seats_taken - 1 WHERE event_id = $1;

-- name: w_cancel_ticket
-- params: ticket_id
DELETE FROM ticket WHERE ticket_id = $1
RETURNING event_id, seat_no;

-- name: w_edit_event
-- params: description, event_id
UPDATE event SET description = $1 WHERE event_id = $2;

-- name: w_publish_event
-- params: event_id, band_id, name, venue, starts_at, capacity, price_cents, description
INSERT INTO event (event_id, band_id, name, venue, starts_at, capacity, price_cents, description)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
