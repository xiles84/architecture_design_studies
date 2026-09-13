-- Writes for R2 (inventory-buckets).
--
-- Booking, one READ COMMITTED transaction:
--   w_pick_bucket       -> a random bucket with seats left that nobody has locked
--   w_take_from_bucket  -> guarded decrement on that bucket
--   w_insert_ticket     -> create the ticket
-- and, only when w_pick_bucket returns nothing:
--   w_any_remaining     -> sold out, or just all locked right now? (retry if the latter)
--
-- Cancellation, one transaction:
--   w_cancel_ticket, then w_pick_refill_bucket + w_return_to_bucket
--   (a random bucket with room, skipping locked ones; retried if none is free)

-- name: w_pick_bucket
-- params: event_id
-- random() over at most 8 rows: spreads buyers across buckets instead of
-- sending them all to bucket 1.
SELECT bucket_no
  FROM event_inventory_bucket
 WHERE event_id = $1
   AND remaining > 0
 ORDER BY random()
 LIMIT 1
   FOR UPDATE SKIP LOCKED;

-- name: w_any_remaining
-- params: event_id
SELECT EXISTS (SELECT 1 FROM event_inventory_bucket WHERE event_id = $1 AND remaining > 0);

-- name: w_take_from_bucket
-- params: event_id, bucket_no
UPDATE event_inventory_bucket
   SET remaining = remaining - 1
 WHERE event_id = $1
   AND bucket_no = $2
   AND remaining > 0;

-- name: w_insert_ticket
-- params: ticket_id, event_id, seat_no, customer_id, price_cents
INSERT INTO ticket (ticket_id, event_id, seat_no, customer_id, sold_at, price_cents)
VALUES ($1, $2, $3, $4, now(), $5);

-- name: w_cancel_ticket
-- params: ticket_id
DELETE FROM ticket WHERE ticket_id = $1
RETURNING event_id, seat_no;

-- name: w_pick_refill_bucket
-- params: event_id
SELECT bucket_no
  FROM event_inventory_bucket
 WHERE event_id = $1
   AND remaining < bucket_capacity
 ORDER BY random()
 LIMIT 1
   FOR UPDATE SKIP LOCKED;

-- name: w_return_to_bucket
-- params: event_id, bucket_no
UPDATE event_inventory_bucket
   SET remaining = remaining + 1
 WHERE event_id = $1
   AND bucket_no = $2
   AND remaining < bucket_capacity;

-- name: w_edit_event
-- params: description, event_id
UPDATE event SET description = $1 WHERE event_id = $2;

-- name: w_publish_event
-- params: event_id, band_id, name, venue, starts_at, capacity, price_cents, description
INSERT INTO event (event_id, band_id, name, venue, starts_at, capacity, price_cents, description)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: w_publish_buckets
-- params: event_id, capacity
-- min(8, capacity) buckets; the remainder of capacity / buckets goes to the
-- first buckets, one seat each, so bucket sizes differ by at most one.
INSERT INTO event_inventory_bucket (event_id, bucket_no, bucket_capacity, remaining)
SELECT $1::BIGINT, b,
       ($2::INT / k) + CASE WHEN b <= $2::INT % k THEN 1 ELSE 0 END,
       ($2::INT / k) + CASE WHEN b <= $2::INT % k THEN 1 ELSE 0 END
  FROM (SELECT LEAST(8, $2::INT) AS k) s,
       generate_series(1, s.k) AS b;
