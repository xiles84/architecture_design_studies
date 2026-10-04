-- Membership and absence spellings for the Study 01 SELECT experiment
-- (task select-existence).
--
-- The parent-identity membership (P6) and absence (P7) questions already have
-- their EXISTS / IN / JOIN+DISTINCT and NOT EXISTS / LEFT-IS-NULL spellings in
-- select_catalog.sql. These statements add the NOT-NULL subquery contract for
-- `IN`/`NOT IN` and the deliberately-wrong membership multiplier. The fanout
-- axis is a property of the loaded fixture, not of the SQL, so these statements
-- run unchanged at each fanout level.

-- name: e12_in_notnull
-- group: sel_membership
-- params: pred
-- `IN` over a subquery that guarantees a non-null child key. The positive
-- membership answer is the same whether or not the subquery filters NULLs;
-- this arm exists to pair with the NULL-contract case and to keep the parent
-- identity grain explicit.
SELECT c.id, c.name FROM customers c
WHERE c.id IN (SELECT o.customer_id FROM orders o
               WHERE o.customer_id IS NOT NULL AND (:pred));

-- name: e13_not_in_notnull
-- group: sel_absence
-- params: pred
-- `NOT IN` is safe only under a non-null contract: the subquery filters the
-- nullable child key. Without that filter one NULL makes every comparison
-- UNKNOWN and the answer is empty, which is what x02_not_in_nullable shows.
SELECT c.id, c.name FROM customers c
WHERE c.id NOT IN (SELECT o.customer_id FROM orders o
                   WHERE o.customer_id IS NOT NULL AND (:pred));

-- name: x07_membership_multiplier
-- group: control_wrong_answer
-- params: pred
-- DELIBERATELY WRONG membership grain: a plain child-row join emits one row per
-- qualifying child, not one row per parent, so it must be detected against the
-- membership set and never timed as a membership spelling.
SELECT c.id, c.name FROM customers c
JOIN orders o ON o.customer_id = c.id
WHERE (:pred);
