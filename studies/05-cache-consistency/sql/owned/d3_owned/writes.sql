-- Owned write catalogue.
--
-- The five logical mutations are the legacy model's, unchanged. What is added is
-- the cache-correctness machinery, and it is deliberately a small set of
-- statements so a reviewer can read exactly what a mutation costs:
--
--   * w_version_bump   -- one UPDATE of the parent row, in the mutation's own
--                         transaction, so the version moves atomically with the
--                         authoritative state;
--   * w_outbox_insert  -- the change event, same transaction;
--   * w_lock_person    -- the pessimistic strategy: lock the parent before
--                         reading anything;
--   * w_lock_persons2  -- ordered locking for reassignment; two parents, locked
--                         in a deterministic order so two concurrent
--                         reassignments in opposite directions cannot deadlock;
--   * w_version_cas    -- the optimistic strategy: bump only if the version has
--                         not moved. Zero rows affected IS the conflict signal.
--
-- The harness composes these inside one transaction. That is not a convenience:
-- the study's floor is that a cache value may be published only for a committed
-- database state, and composition inside the transaction is what makes "the
-- version moved", "the outbox has the event" and "the mutation committed" the
-- same event.

-- name: w_load_donations
-- params: none
INSERT INTO donation (donation_id, person_id, charity_id, amount_cents, currency, donated_at, note)
SELECT donation_id, person_id, charity_id, amount_cents, currency, donated_at, note
  FROM load_donation;

-- name: w_donation_insert
-- params: donation_id, person_id, charity_id, amount_cents, currency, donated_at, note
INSERT INTO donation (donation_id, person_id, charity_id, amount_cents, currency, donated_at, note)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: w_donation_correct
-- params: donation_id, amount_cents
UPDATE donation SET amount_cents = $2 WHERE donation_id = $1;

-- name: w_donation_delete
-- params: donation_id
DELETE FROM donation WHERE donation_id = $1;

-- name: w_person_update
-- params: person_id, full_name, email
UPDATE person SET full_name = $2, email = $3 WHERE person_id = $1;

-- name: w_donation_reassign
-- params: donation_id, new_person_id
UPDATE donation d
   SET person_id  = $2,
       charity_id = (SELECT p.charity_id FROM person p WHERE p.person_id = $2)
 WHERE d.donation_id = $1;

-- name: w_version_bump
-- params: person_id
UPDATE person SET cache_version = cache_version + 1 WHERE person_id = $1;

-- name: w_outbox_insert
-- params: person_id, kind, donation_id
INSERT INTO cache_outbox (person_id, kind, donation_id) VALUES ($1, $2, $3);

-- name: w_lock_person
-- params: person_id
-- The pessimistic strategy. FOR UPDATE, not FOR SHARE: the caller is about to
-- write. Blocking here is the design, not an accident, and the harness records
-- the wait.
SELECT cache_version FROM person WHERE person_id = $1 FOR UPDATE;

-- name: w_lock_persons2
-- params: person_a, person_b
-- Two people, one deterministic order. Without the ORDER BY, two concurrent
-- reassignments in opposite directions each hold the row the other needs and the
-- pair deadlocks -- which on PostgreSQL is a retryable error and on YugabyteDB a
-- transaction that must restart. Locking in id order removes the cycle instead of
-- retrying around it.
SELECT person_id, cache_version
  FROM person
 WHERE person_id IN ($1, $2)
 ORDER BY person_id
   FOR UPDATE;

-- name: w_version_cas
-- params: person_id, expected_version
-- The optimistic strategy. The guard and the bump are one statement, so there is
-- no window between checking the version and moving it. Zero rows affected means
-- someone else committed first; the caller retries within its deadline and counts
-- the conflict (methodology 6a).
UPDATE person
   SET cache_version = cache_version + 1
 WHERE person_id = $1
   AND cache_version = $2;

-- name: w_version_set
-- params: person_id, new_version
-- The UNSAFE form, used only by the negative control. The new version is computed
-- by the application from a value it read earlier and written back unguarded, so
-- two concurrent writers can both be acknowledged while one bump is lost. Every
-- candidate design uses w_version_cas or w_version_bump instead; this statement
-- exists so the control can be the wrong design on purpose.
UPDATE person SET cache_version = $2 WHERE person_id = $1;
