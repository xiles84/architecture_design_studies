-- Aggregate and per-parent top-N sweep family for the Study 01 SELECT experiment
-- (task select-aggregate-topn).
--
-- These statements ask the same logical questions as s12/s13/s14/s15, but
-- restricted to a prefix of parent ids (`:parent_max`, 0 = every parent), so a
-- single deterministic fixture serves the parent-selection and parent-count
-- axes without reloading. `:parent_max` is a VALUED parameter, not a predicate.
--
-- Groups:
--   sel_aggregate_sel  one row per selected parent: SUM(amount) and COUNT(*)
--   sel_topn_sel       at most :n child rows per selected parent, total order
--   sel_topn_one       DISTINCT ON: the N=1-only spelling of the same contract

-- name: a01_aggregate_join_sel
-- group: sel_aggregate_sel
-- params: parent_max
-- Post-join aggregate: LEFT JOIN preserves childless parents (INNER would drop
-- them and answer a different question), then GROUP BY. Correct here because
-- only one one-to-many relation is joined; adding a second would multiply the
-- totals, which is what the x05 control demonstrates.
SELECT c.id, c.name, COALESCE(SUM(o.amount), 0) AS total, COUNT(o.id) AS children
FROM customers c
LEFT JOIN orders o ON o.customer_id = c.id
WHERE :parent_max = 0 OR c.id <= :parent_max
GROUP BY c.id, c.name;

-- name: a02_aggregate_preagg_sel
-- group: sel_aggregate_sel
-- params: parent_max
-- Pre-aggregate the child side, then join: no row multiplication and no
-- per-parent re-aggregation of the big side.
SELECT c.id, c.name, COALESCE(a.total, 0) AS total, COALESCE(a.children, 0) AS children
FROM customers c
LEFT JOIN (SELECT customer_id, SUM(amount) AS total, COUNT(*) AS children
           FROM orders GROUP BY customer_id) a ON a.customer_id = c.id
WHERE :parent_max = 0 OR c.id <= :parent_max;

-- name: a03_aggregate_correlated_sel
-- group: sel_aggregate_sel
-- params: parent_max
-- Correlated scalar subqueries: one aggregate probe per selected parent.
SELECT c.id, c.name,
       COALESCE((SELECT SUM(o.amount) FROM orders o WHERE o.customer_id = c.id), 0) AS total,
       (SELECT COUNT(*) FROM orders o WHERE o.customer_id = c.id) AS children
FROM customers c
WHERE :parent_max = 0 OR c.id <= :parent_max;

-- name: a04_topn_rownumber_sel
-- group: sel_topn_sel
-- params: parent_max, n
-- One global sort with ROW_NUMBER and an explicit total-order tiebreaker
-- (amount DESC NULLS LAST, id ASC).
SELECT customer_id, id, amount FROM (
  SELECT o.customer_id, o.id, o.amount,
         ROW_NUMBER() OVER (PARTITION BY o.customer_id
                            ORDER BY o.amount DESC NULLS LAST, o.id ASC) AS rn
  FROM orders o
  WHERE o.customer_id IS NOT NULL
    AND (:parent_max = 0 OR o.customer_id <= :parent_max)) t
WHERE rn <= :n
ORDER BY customer_id ASC, amount DESC NULLS LAST, id ASC;

-- name: a05_topn_lateral_sel
-- group: sel_topn_sel
-- params: parent_max, n
-- One ordered index probe per selected parent through LATERAL ... LIMIT n,
-- under the same total order.
SELECT c.id, o.id, o.amount
FROM customers c
CROSS JOIN LATERAL (
  SELECT o.id, o.amount FROM orders o
  WHERE o.customer_id = c.id
  ORDER BY o.amount DESC NULLS LAST, o.id ASC
  LIMIT :n) o
WHERE :parent_max = 0 OR c.id <= :parent_max
ORDER BY c.id ASC, o.amount DESC NULLS LAST, o.id ASC;

-- name: a06_topn_distinct_on_sel
-- group: sel_topn_one
-- params: parent_max
-- DISTINCT ON is an N=1-only spelling (PostgreSQL/YugabyteDB extension). It
-- keeps one row per parent and is comparable to a04/a05 only at N=1.
SELECT DISTINCT ON (o.customer_id) o.customer_id, o.id, o.amount
FROM orders o
WHERE o.customer_id IS NOT NULL
  AND (:parent_max = 0 OR o.customer_id <= :parent_max)
ORDER BY o.customer_id ASC, o.amount DESC NULLS LAST, o.id ASC;

-- name: x05_aggregate_double_count
-- group: control_wrong_answer
-- params: parent_max
-- DELIBERATELY WRONG pre-aggregation grain: joining the child relation to
-- itself multiplies every total and count by the parent's child count. It must
-- be detected as a different answer, never shown as a faster aggregate.
SELECT c.id, c.name, COALESCE(SUM(o1.amount), 0) AS total, COUNT(o2.id) AS children
FROM customers c
JOIN orders o1 ON o1.customer_id = c.id
JOIN orders o2 ON o2.customer_id = c.id
WHERE :parent_max = 0 OR c.id <= :parent_max
GROUP BY c.id, c.name;

-- name: x06_topn_rank_ties
-- group: control_wrong_answer
-- params: parent_max, n
-- DELIBERATELY WRONG tie behaviour: RANK keeps every tied row at the boundary
-- rank (a separate contract), so a parent with a tie starting at rank 1 returns
-- more than :n rows. It must be detected wherever the fixture holds a tie.
SELECT customer_id, id, amount FROM (
  SELECT o.customer_id, o.id, o.amount,
         RANK() OVER (PARTITION BY o.customer_id ORDER BY o.amount DESC NULLS LAST) AS rnk
  FROM orders o
  WHERE o.customer_id IS NOT NULL
    AND (:parent_max = 0 OR o.customer_id <= :parent_max)) t
WHERE rnk <= :n
ORDER BY customer_id ASC, amount DESC NULLS LAST, id ASC;
