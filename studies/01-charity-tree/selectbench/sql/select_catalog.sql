-- Named SELECT catalogue for the Study 01 additive SELECT experiment.
--
-- The `group:` line is a semantic-group id: two statements may be compared only
-- when their group ids match, because a faster statement that answers a different
-- question is not a competing template. The loader (catalog.go) refuses to pool
-- groups, and the comparison report must label any cross-group pair invalid.
--
-- Statements are written for PostgreSQL 17 / YSQL against the teaching fixture
-- `customers(id, name)` and `orders(id, customer_id, amount)`, which map onto
-- Study 01's charity/person/donation domain. `:pred` is bound by the harness.

-- name: s01_join_where
-- group: sel_full_rows
-- params: pred
-- The canonical template: join, then filter. Oracle: bag of (parent, child).
SELECT c.id, c.name, o.id, o.amount
FROM customers c JOIN orders o ON o.customer_id = c.id
WHERE :pred;

-- name: s02_filtered_derived
-- group: sel_full_rows
-- params: pred
-- Filter first in a derived table, then join.
SELECT c.id, c.name, o.id, o.amount
FROM customers c
JOIN (SELECT * FROM orders WHERE :pred) o ON o.customer_id = c.id;

-- name: s03_cte_filtered
-- group: sel_full_rows
-- params: pred
-- The same filter named by a CTE; folding is a planner decision, not a promise.
WITH e AS (SELECT * FROM orders WHERE :pred)
SELECT c.id, c.name, e.id, e.amount
FROM customers c JOIN e ON e.customer_id = c.id;

-- name: s04_join_on
-- group: sel_full_rows
-- params: pred
-- Predicate in ON rather than WHERE (inner join: logically equivalent).
SELECT c.id, c.name, o.id, o.amount
FROM customers c JOIN orders o ON o.customer_id = c.id AND (:pred);

-- name: s05_comma_join
-- group: sel_full_rows
-- params: pred
-- Comma join with both predicates in WHERE.
SELECT c.id, c.name, o.id, o.amount
FROM customers c, orders o
WHERE c.id = o.customer_id AND (:pred);

-- name: s06_exists_parent
-- group: sel_membership
-- params: pred
-- Parent membership: one row per parent with at least one qualifying child.
SELECT c.id, c.name FROM customers c
WHERE EXISTS (SELECT 1 FROM orders o WHERE o.customer_id = c.id AND (:pred));

-- name: s07_in_parent
-- group: sel_membership
-- params: pred
SELECT c.id, c.name FROM customers c
WHERE c.id IN (SELECT o.customer_id FROM orders o WHERE (:pred));

-- name: s08_join_distinct_parent
-- group: sel_membership
-- params: pred
-- De-duplicate on the parent IDENTITY, never on the name.
SELECT DISTINCT c.id, c.name FROM customers c
JOIN orders o ON o.customer_id = c.id WHERE (:pred);

-- name: s09_not_exists
-- group: sel_absence
-- params: pred
-- Parents with no qualifying child.
SELECT c.id, c.name FROM customers c
WHERE NOT EXISTS (SELECT 1 FROM orders o WHERE o.customer_id = c.id AND (:pred));

-- name: s10_left_anti
-- group: sel_absence
-- params: pred
SELECT c.id, c.name FROM customers c
LEFT JOIN orders o ON o.customer_id = c.id AND (:pred)
WHERE o.id IS NULL;

-- name: s11_optional_on
-- group: sel_optional_match
-- params: pred
-- Every parent preserved; qualifying children or NULLs. The WHERE variant below
-- is a DIFFERENT answer and is kept as a read-path control, not a template.
SELECT c.id, c.name, o.id, o.amount
FROM customers c LEFT JOIN orders o ON o.customer_id = c.id AND (:pred);

-- name: s12_aggregate_join
-- group: sel_aggregate
-- Post-join aggregate: correct only because one child collection is aggregated.
SELECT c.id, c.name, COALESCE(SUM(o.amount), 0) AS total, COUNT(o.id) AS children
FROM customers c LEFT JOIN orders o ON o.customer_id = c.id
GROUP BY c.id, c.name;

-- name: s13_aggregate_preagg
-- group: sel_aggregate
-- Pre-aggregate the child side, then join: no row multiplication.
SELECT c.id, c.name, COALESCE(a.total, 0) AS total, COALESCE(a.children, 0) AS children
FROM customers c
LEFT JOIN (SELECT customer_id, SUM(amount) AS total, COUNT(*) AS children
           FROM orders GROUP BY customer_id) a ON a.customer_id = c.id;

-- name: s14_topn_rownumber
-- group: sel_topn
-- params: n
-- Top-N per parent with an explicit total order (amount DESC, id ASC), so ties
-- are broken deterministically. The projection order is fixed across the
-- top-N statements: (customer_id, order_id, amount).
SELECT customer_id, id, amount FROM (
  SELECT o.*, ROW_NUMBER() OVER (PARTITION BY customer_id ORDER BY amount DESC NULLS LAST, id ASC) AS rn
  FROM orders o) t WHERE rn <= :n;

-- name: s15_topn_lateral
-- group: sel_topn
-- params: n
SELECT c.id, o.id, o.amount
FROM customers c
LEFT JOIN LATERAL (SELECT o.id, o.amount FROM orders o
                   WHERE o.customer_id = c.id
                   ORDER BY o.amount DESC NULLS LAST, o.id ASC LIMIT :n) o ON true;

-- name: s16_page_offset
-- group: sel_pagination
-- params: limit, offset
SELECT id, amount FROM orders ORDER BY amount DESC NULLS LAST, id ASC LIMIT :limit OFFSET :offset;

-- name: s17_page_keyset
-- group: sel_pagination
-- params: limit, last_amount, last_id
-- Keyset page after the last row of the previous page under the same order.
SELECT id, amount FROM orders
WHERE (amount, id) < (:last_amount, :last_id)
ORDER BY amount DESC NULLS LAST, id ASC LIMIT :limit;

-- name: x01_optional_where
-- group: control_wrong_answer
-- params: pred
-- DELIBERATELY WRONG as an optional-match template: the WHERE predicate drops
-- preserved parents. It must be detected as a different answer, never shown as
-- a faster optional-match alternative.
SELECT c.id, c.name, o.id, o.amount
FROM customers c LEFT JOIN orders o ON o.customer_id = c.id
WHERE (:pred);

-- name: x02_not_in_nullable
-- group: control_wrong_answer
-- params: pred
-- DELIBERATELY WRONG as an unconditional NOT EXISTS substitute: a NULL in the
-- subquery result makes every comparison UNKNOWN.
SELECT c.id, c.name FROM customers c
WHERE c.id NOT IN (SELECT o.customer_id FROM orders o WHERE (:pred));

-- name: x03_distinct_on_name
-- group: control_wrong_answer
-- params: pred
-- DELIBERATELY WRONG: de-duplicating on the name merges two distinct parents.
SELECT DISTINCT c.name FROM customers c
JOIN orders o ON o.customer_id = c.id WHERE (:pred);

-- name: x04_limit_before_join
-- group: control_wrong_answer
-- params: pred, limit
-- DELIBERATELY WRONG: the limit is applied before the multiplicative join, so
-- the joined result is truncated arbitrarily.
SELECT c.id, o.id FROM customers c
JOIN (SELECT * FROM orders WHERE (:pred) ORDER BY id LIMIT :limit) o
  ON o.customer_id = c.id;
