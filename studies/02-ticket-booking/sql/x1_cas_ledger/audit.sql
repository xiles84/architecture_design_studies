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
