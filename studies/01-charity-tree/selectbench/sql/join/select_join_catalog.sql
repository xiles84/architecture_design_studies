-- Join/filter spellings and CTE materialization catalogue (task select-join-cte).
--
-- Loaded by LoadJoin from sql/join/ only, so the teaching-fixture families keep
-- the catalogue they were validated against. The schema is `selectjoin`:
-- customers(id, name) joined to facts(id, customer_id, rank, amount, payload).
--
-- `:threshold` is a rank predicate, not an amount predicate: `rank` is a
-- deterministic shuffle of 1..N, so `rank <= :threshold` selects an exact
-- fraction of the fact table (0.01%, 0.1%, 1%, 10%, 50%, 90% at the harness
-- thresholds). That mapping is labelled in every report; it is not the teaching
-- `amount > n` question.
--
-- Every statement in group sel_join_bag answers the SAME bag: the rows
-- (customer_id, customer_name, fact_id, amount) whose child rank is within the
-- threshold. They differ only in spelling, so a timing difference is
-- attributable to the spelling (or to the plan it selects), never the answer.
-- MATERIALIZED / NOT MATERIALIZED exist only for the two CTE spellings.

-- name: j01_join_where
-- group: sel_join_bag
-- params: threshold
-- Canonical template: join, then filter in WHERE.
SELECT c.id AS customer_id, c.name AS customer_name, f.id AS fact_id, f.amount
FROM customers c JOIN facts f ON f.customer_id = c.id
WHERE f.rank <= :threshold;

-- name: j02_join_on
-- group: sel_join_bag
-- params: threshold
-- Inner join, so a predicate in ON is logically the same bag as WHERE.
SELECT c.id AS customer_id, c.name AS customer_name, f.id AS fact_id, f.amount
FROM customers c JOIN facts f ON f.customer_id = c.id AND f.rank <= :threshold;

-- name: j03_filtered_derived
-- group: sel_join_bag
-- params: threshold
-- Filter first in a derived table, then join.
SELECT c.id AS customer_id, c.name AS customer_name, f.id AS fact_id, f.amount
FROM customers c
JOIN (SELECT * FROM facts WHERE rank <= :threshold) f ON f.customer_id = c.id;

-- name: j04_cte_filtered
-- group: sel_join_bag
-- params: threshold
-- The same filter named by a CTE; folding is a planner decision, not a promise.
-- This is the DEFAULT policy spelling.
WITH e AS (SELECT * FROM facts WHERE rank <= :threshold)
SELECT c.id AS customer_id, c.name AS customer_name, e.id AS fact_id, e.amount
FROM customers c JOIN e ON e.customer_id = c.id;

-- name: j04m_cte_filtered_mat
-- group: sel_join_bag
-- params: threshold
-- Forced MATERIALIZED filtered-input CTE.
WITH e AS MATERIALIZED (SELECT * FROM facts WHERE rank <= :threshold)
SELECT c.id AS customer_id, c.name AS customer_name, e.id AS fact_id, e.amount
FROM customers c JOIN e ON e.customer_id = c.id;

-- name: j04n_cte_filtered_notmat
-- group: sel_join_bag
-- params: threshold
-- Forced NOT MATERIALIZED filtered-input CTE; the planner may inline it.
WITH e AS NOT MATERIALIZED (SELECT * FROM facts WHERE rank <= :threshold)
SELECT c.id AS customer_id, c.name AS customer_name, e.id AS fact_id, e.amount
FROM customers c JOIN e ON e.customer_id = c.id;

-- name: j05_cte_joined_outer
-- group: sel_join_bag
-- params: threshold
-- Joined-output CTE with the predicate applied outside it: the boundary that
-- decides whether the selective filter can be pushed into the join. DEFAULT.
WITH e AS (
  SELECT c.id AS customer_id, c.name AS customer_name, f.id AS fact_id, f.amount, f.rank
  FROM customers c JOIN facts f ON f.customer_id = c.id
)
SELECT customer_id, customer_name, fact_id, amount FROM e WHERE rank <= :threshold;

-- name: j05m_cte_joined_outer_mat
-- group: sel_join_bag
-- params: threshold
-- Forced MATERIALIZED joined-output CTE; the outer predicate cannot be pushed
-- into a materialized CTE, so the join runs before the filter.
WITH e AS MATERIALIZED (
  SELECT c.id AS customer_id, c.name AS customer_name, f.id AS fact_id, f.amount, f.rank
  FROM customers c JOIN facts f ON f.customer_id = c.id
)
SELECT customer_id, customer_name, fact_id, amount FROM e WHERE rank <= :threshold;

-- name: j05n_cte_joined_outer_notmat
-- group: sel_join_bag
-- params: threshold
-- Forced NOT MATERIALIZED joined-output CTE; the planner may push the outer
-- predicate into the inlined query and filter before joining.
WITH e AS NOT MATERIALIZED (
  SELECT c.id AS customer_id, c.name AS customer_name, f.id AS fact_id, f.amount, f.rank
  FROM customers c JOIN facts f ON f.customer_id = c.id
)
SELECT customer_id, customer_name, fact_id, amount FROM e WHERE rank <= :threshold;

-- name: j06_comma_join
-- group: sel_join_bag
-- params: threshold
-- Legacy comma join with both predicates in WHERE.
SELECT c.id AS customer_id, c.name AS customer_name, f.id AS fact_id, f.amount
FROM customers c, facts f
WHERE c.id = f.customer_id AND f.rank <= :threshold;

-- name: j11_optional_on
-- group: sel_join_optional
-- params: threshold
-- Predicate in ON: every parent preserved, with qualifying children or NULLs.
-- The WHERE variant below is a DIFFERENT answer and is kept as a detected
-- control, never as a faster optional-match spelling.
SELECT c.id AS customer_id, c.name AS customer_name, f.id AS fact_id, f.amount
FROM customers c LEFT JOIN facts f ON f.customer_id = c.id AND f.rank <= :threshold;

-- name: x01_optional_where
-- group: control_wrong_answer
-- params: threshold
-- DELIBERATELY WRONG as an optional-match template: a child predicate in WHERE
-- drops every preserved parent with no qualifying child. It must be detected as
-- a different answer, never timed against sel_join_optional.
SELECT c.id AS customer_id, c.name AS customer_name, f.id AS fact_id, f.amount
FROM customers c LEFT JOIN facts f ON f.customer_id = c.id
WHERE f.rank <= :threshold;

-- name: x02_join_limit_before
-- group: control_wrong_answer
-- params: threshold
-- DELIBERATELY WRONG as a join-bag template: the limit is applied to the child
-- side before the multiplicative join, so the joined result is truncated
-- arbitrarily. The fixed limit 7 is mirrored by the Go control oracle.
SELECT c.id AS customer_id, c.name AS customer_name, f.id AS fact_id, f.amount
FROM customers c
JOIN (SELECT * FROM facts WHERE rank <= :threshold ORDER BY id LIMIT 7) f
  ON f.customer_id = c.id;
