-- Audit for the pre-created family. The harness compares these recomputations
-- with the capacity of every event and with its own count of the bookings and
-- cancellations it saw succeed.

-- name: a_sold_per_event
-- params: none
-- Every event, with the number of seats that are actually sold.
SELECT e.event_id,
       e.capacity,
       COALESCE(s.sold, 0)::BIGINT AS sold
  FROM event e
  LEFT JOIN (SELECT event_id, COUNT(*) AS sold
               FROM ticket WHERE status = 'sold'
              GROUP BY event_id) s ON s.event_id = e.event_id;

-- name: a_event_sold
-- params: event_id
SELECT COUNT(*)::BIGINT FROM ticket WHERE event_id = $1 AND status = 'sold';

-- name: a_duplicate_seats
-- params: none
-- A seat sold to two buyers. Structurally impossible here (one row per seat),
-- and checked anyway: an audit that only checks what a design could plausibly
-- get wrong is testing the auditor's assumptions, not the design.
SELECT event_id, seat_no, COUNT(*)::BIGINT AS n
  FROM ticket
 WHERE status = 'sold'
 GROUP BY event_id, seat_no
HAVING COUNT(*) > 1;

-- name: a_derived_drift
-- params: none
-- Pre-created inventory with no counter has no derived state that can drift.
-- Every event must still own exactly `capacity` seat rows.
SELECT e.event_id, e.capacity::BIGINT AS stored, COUNT(t.ticket_id)::BIGINT AS actual
  FROM event e
  LEFT JOIN ticket t ON t.event_id = e.event_id
 GROUP BY e.event_id, e.capacity
HAVING COUNT(t.ticket_id) <> e.capacity;

-- name: a_ledger_mismatches
-- params: none
-- The ledger reconciliation audit (REPORTS.md section 4): per seat, the
-- ledger's net (sold events minus cancelled events) must equal whether that
-- seat is CURRENTLY sold. A mismatch means sale_event was not written in the
-- same transaction as the ticket state change it claims to record -- exactly
-- what would happen if the ledger insert were ever split from the UPDATE.
-- (Reconciliation against what the harness itself saw commit is already
-- covered: RunAudit's existing ledger check compares ticket.status='sold'
-- counts to world.expectedSold(), and this check ties sale_event to that same
-- ticket state, so the composition covers both halves of the requirement.)
SELECT t.event_id, t.seat_no,
       COALESCE(l.net, 0)::BIGINT AS ledger_net,
       (t.status = 'sold')::INT::BIGINT AS actual_sold
  FROM ticket t
  LEFT JOIN (
        SELECT event_id, seat_no,
               COUNT(*) FILTER (WHERE kind = 'sold') - COUNT(*) FILTER (WHERE kind = 'cancelled') AS net
          FROM sale_event
         GROUP BY event_id, seat_no
  ) l ON l.event_id = t.event_id AND l.seat_no = t.seat_no
 WHERE COALESCE(l.net, 0) <> (t.status = 'sold')::INT;

-- name: a_ledger_attribution
-- params: none
-- Order-free attribution audit (EH-02 AM-04.3, deciding RR-ER-01). It never asks
-- "which record came first": neither `at` (now(), fixed at transaction START) nor
-- sale_event_id (YSQL hands out sequence blocks per connection) reconstructs commit
-- order on both engines, and an earlier ordered version of this check gave false
-- positives on PostgreSQL and up to 246 deterministic ones on YugabyteDB
-- (results/devchecks/am03-dc0{1,4}-*). Counting needs no order.
--   B  no buyer is refunded more times than they bought that seat -- what catches a
--      refund credited to the wrong (stale) buyer;
--   C  a seat that is sold now has its sale recorded verbatim (buyer, time, price),
--      and its current buyer has not refunded every sale they made of it.
-- ("Second sale with no refund between" is dropped: one row per seat plus the CAS
-- predicate make it structurally impossible, and a_duplicate_seats checks it.)
SELECT event_id, seat_no, 'buyer refunded more times than they bought this seat' AS problem
  FROM sale_event
 GROUP BY event_id, seat_no, customer_id
HAVING COUNT(*) FILTER (WHERE kind = 'cancelled') > COUNT(*) FILTER (WHERE kind = 'sold')
UNION ALL
SELECT t.event_id, t.seat_no, 'live sale not recorded verbatim in the ledger'
  FROM ticket t
 WHERE t.status = 'sold'
   AND NOT EXISTS (SELECT 1 FROM sale_event s
                    WHERE s.event_id = t.event_id AND s.seat_no = t.seat_no
                      AND s.kind = 'sold' AND s.customer_id = t.customer_id
                      AND s.at = t.sold_at AND s.price_cents = t.price_cents)
UNION ALL
SELECT t.event_id, t.seat_no, 'live buyer has refunded every sale they made of this seat'
  FROM ticket t
 WHERE t.status = 'sold'
   AND (SELECT COUNT(*) FILTER (WHERE s.kind = 'cancelled') - COUNT(*) FILTER (WHERE s.kind = 'sold')
          FROM sale_event s
         WHERE s.event_id = t.event_id AND s.seat_no = t.seat_no AND s.customer_id = t.customer_id) >= 0;
