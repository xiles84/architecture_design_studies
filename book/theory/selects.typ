#import "../lib/config.typ": *
#import "../lib/teaching.typ": strategy, flow-picture
#set heading(numbering: none)
#pagebreak(weak: true)
#heading(level: 2, numbering: none, "Database SELECTs")
Two SQL expressions are comparable only if they return the same logical answer. Establish grain, duplicates, NULLs, ties and ordering first.

Strategies: #link(<theory-select-contract>)[The answer contract] · #link(<theory-select-join-filter>)[Joins and filter placement] · #link(<theory-select-cte>)[CTE folding and materialization] · #link(<theory-select-membership>)[Membership and absence] · #link(<theory-select-aggregate>)[Aggregation and per-parent top-N] · #link(<theory-select-pagination>)[Pagination and access paths].

#strategy("The answer contract", "theory-select-contract", [Define the answer before choosing the SQL spelling.],
  [A charity needs either one row per donor or one row per donation. Joining children changes row multiplicity unless the query restores the intended grain.],
  ([Name the output grain.], [Specify duplicates, NULLs, ties and order.], [Check answers against an independent oracle.],),
  ([A query has an explicit business question.], [Edge cases are part of verification.], [Alternative spellings preserve the contract.],),
  ([An inner join silently drops required parents.], [Deduplication removes meaningful repeated facts.], [A faster query answers a narrower question.],),
  ([Work required by the answer contract.], [None unless a maintained answer is introduced.], [Intermediate rows or maintained results.], [Validate all alternatives against the same oracle.], [Keep wrong-answer controls.],),
  flow-picture(([Business question], [Answer contract], [Oracle-checked result],), "The answer contract: the operation path.", [Business question → Answer contract → Oracle-checked result. Arrows show logical order, not a measured duration.]),
  ([Same logical answer across alternatives.], [Base facts and query results.], [Execution resources; not a new storage invariant.], [None for plain query spelling.], [The returned relation.], [Source snapshots or external copies.], [Unspecified grain, multiplicity, NULLs, ties or order.], [Compare edge-case outputs before timing.],), "case-charity-select-guide", measured: true, example: [
#table(columns: (1fr, 1.5fr), inset: 5pt, stroke: 0.4pt + palette.rule,
table.header([Business question], [Required grain]),
[Who donated?], [One row per qualifying donor.],
[Which donations arrived?], [One row per donation, including repeated donors.],
[Who has not donated?], [One row per donor with no matching donation.])
A join can multiply donor rows. DISTINCT repairs only the intended grain; it must not erase required donation rows.
])

#let query-example(title, body) = block(breakable: false)[
  #text(weight: "bold")[#title]
  #body
]

#strategy("Joins and filter placement", "theory-select-join-filter", [Equivalent inner-join spellings often share a plan; moving an outer-join filter can change the answer.],
  [A customer page needs each order above 100 together with its customer. Does writing the filter before the join reduce work?],
  ([Fix the returned columns and row multiplicities.], [Write the relationship and amount predicate explicitly.], [Inspect the actual plan; SQL layout does not prescribe execution order.],),
  ([A simple inner join preserves the required answer.], [A readable spelling makes the relationship clear.], [The optimizer can move the filter into an input scan.],),
  ([A child WHERE predicate discards required outer-join parents.], [A subquery changes duplicates, aggregation or limit semantics.], [Written filter order is treated as a timing guarantee.],),
  ([The chosen scan and join, not the position of WHERE.], [No extra maintained state.], [Plan-dependent join and intermediate-row memory.], [Preserve duplicates and required unmatched parents.], [Compare EXPLAIN ANALYZE and replicated timings with matched inputs.],),
  flow-picture(([Equivalent SQL forms], [Optimizer transformations], [Chosen scan and join],), "Query spelling and execution.", [Different SQL expressions can become the same scan and join. This is a conceptual optimization path, not proof of a particular engine plan.]),
  ([Return the same bag of customer-order rows.], [Customer and order base tables.], [Scan, join and client resources.], [None for these query forms.], [Transient joined rows.], [Results reflect the selected database snapshot.], [A changed answer or an unfavorable scan/join plan.], [Verify bags, sweep selectivity and compare plans across replicated runs.],), "case-charity-join-plan", measured: true, example: [
*Illustrative tables:* `customers(id, name)` and `orders(id, customer_id, amount)`.
Ana (id 1) has orders 11 and 14 worth 200 each. Ben (id 2) has an order worth 80 and one with NULL amount; Cara (id 3) has none.

All six forms return `(Ana, 11, 200)` and `(Ana, 14, 200)`. These are example values, not measurements. Without ORDER BY, row order is unspecified.

#query-example("1. Explicit join, WHERE filter", [
```sql
SELECT c.name, o.id AS order_id, o.amount
FROM customers c
JOIN orders o ON o.customer_id = c.id
WHERE o.amount > 100;
```
])
#query-example("2. Filter inside ON", [
```sql
SELECT c.name, o.id AS order_id, o.amount
FROM customers c JOIN orders o
  ON o.customer_id = c.id AND o.amount > 100;
```
])
#query-example("3. Filtered derived table", [
```sql
SELECT c.name, o.id AS order_id, o.amount
FROM customers c JOIN (
  SELECT id, customer_id, amount FROM orders
  WHERE amount > 100
) o ON o.customer_id = c.id;
```
])
#query-example("4. CTE naming the filtered orders", [
```sql
WITH expensive_orders AS (
  SELECT id, customer_id, amount FROM orders
  WHERE amount > 100
)
SELECT c.name, o.id AS order_id, o.amount
FROM customers c
JOIN expensive_orders o ON o.customer_id = c.id;
```
])
#query-example("5. CTE naming the joined rows", [
```sql
WITH customer_orders AS (
  SELECT c.name, o.id AS order_id, o.amount
  FROM customers c
  JOIN orders o ON o.customer_id = c.id
)
SELECT name, order_id, amount
FROM customer_orders WHERE amount > 100;
```
])
#query-example("6. Comma join, both predicates in WHERE", [
```sql
SELECT c.name, o.id AS order_id, o.amount
FROM customers c, orders o
WHERE o.customer_id = c.id AND o.amount > 100;
```
])
*Expected performance:* the direct join, ON filter and comma join often converge to the same plan. Simple derived tables and folded CTEs can converge too. A plan match is a hypothesis to test, not a guarantee of identical end-to-end timing.

*A different answer with LEFT JOIN.* The first expression keeps Ben and Cara with NULL order columns. The second drops them:
#query-example("Preserve every customer", [
```sql
SELECT c.name, o.id AS order_id, o.amount
FROM customers c LEFT JOIN orders o
  ON o.customer_id = c.id AND o.amount > 100;
```
])
#query-example("Keep only qualifying customer-order matches", [
```sql
SELECT c.name, o.id AS order_id, o.amount
FROM customers c LEFT JOIN orders o
  ON o.customer_id = c.id
WHERE o.amount > 100;
```
])
Filtering the child input before LEFT JOIN is equivalent to the first form. EXISTS instead returns Ana once when the question asks only who has a qualifying order; it cannot replace required order details.

*Documentation:* #link("https://www.postgresql.org/docs/17/queries-table-expressions.html")[PostgreSQL 17 joins] and #link("https://www.postgresql.org/docs/17/explicit-joins.html")[planner transformations]. Other engines and versions require their own plan checks.
])

#strategy("CTE folding and materialization", "theory-select-cte", [A CTE names intermediate work; folding permits joint optimization, while materialization can reuse work and block predicate pushdown.],
  [The same customer-order query can scan few qualifying rows or create a large intermediate result. Reusing expensive work poses a different tradeoff.],
  ([Check the pinned engine's CTE syntax and default behavior.], [Keep one-use selective access separate from repeated expensive reuse.], [Verify folding or materialization in the actual plan.], [Measure saved computation against intermediate rows, memory and spills.],),
  ([Folding lets a selective outer filter reach an indexed base scan.], [Materialization avoids repeatedly evaluating an expensive input.], [The reused result fits the available execution memory.],),
  ([A one-use materialized result includes mostly discarded rows.], [Separate evaluation blocks a useful outer predicate.], [Forced folding repeats expensive work across many references.],),
  ([Filter pushdown versus shared or repeated computation.], [No persisted rollup or refresh job for a SELECT CTE.], [Transient result memory and possible temporary-file spills.], [Both policies must preserve the same answer and snapshot.], [Inspect engine support, CTE scans, loops and temp/RPC counters.],),
  flow-picture(([Named CTE input], [Fold or evaluate separately], [Filter and reuse],), "The CTE execution choice.", [Folding combines the CTE with its parent for optimization. Materialization evaluates a separate result that consumers can reuse. No branch is universally faster.]),
  ([Same values and multiplicities under either policy.], [Base facts and a transient CTE result.], [Repeated CPU work or intermediate memory/I/O.], [None; materialization here is per query, not a stored read model.], [The CTE output.], [Only the query snapshot; there is no later refresh lifecycle.], [Selective one-use access or repeated reuse with the wrong policy.], [Sweep selectivity and reuse separately with oracle checks and actual plan evidence.],), "case-charity-cte-reuse", measured: true, example: [
#table(columns: (1fr, 1.7fr), inset: 5pt, stroke: 0.4pt + palette.rule,
table.header([Format], [Expected behavior on PostgreSQL 17]),
[`WITH q AS (...)`], [A nonrecursive, side-effect-free, single-use CTE normally folds; multiple references normally materialize.],
[`WITH q AS MATERIALIZED (...)`], [Force separate CTE evaluation; outer predicates cannot be pushed inside it.],
[`WITH q AS NOT MATERIALIZED (...)`], [Permit eligible CTE folding; several references may repeat computation.])

*Selective one-use example:* add MATERIALIZED to form 5. The amount predicate stays outside the CTE, so it cannot reduce the joined rows inside. Add NOT MATERIALIZED to allow folding when the CTE is eligible.

*Repeated-use example:* several consumers need an expensive deterministic calculation over the same input. Materialization may compute it once; folding may compute it for each consumer. Selective consumers can still favor folding by avoiding most input rows.

These are predictions. Compare default, MATERIALIZED and NOT MATERIALIZED with identical answers and indexes. Record unsupported syntax or different defaults on the pinned YugabyteDB engine explicitly.

*Documentation:* #link("https://www.postgresql.org/docs/17/queries-with.html#QUERIES-WITH-CTE-MATERIALIZATION")[PostgreSQL 17 CTE folding and materialization]. This feature predates PostgreSQL 17; engine support is not established by a version-number guess.
])

#strategy("Membership and absence", "theory-select-membership", [Use a membership test when the question asks whether a related row exists.],
  [A charity asks for donors with a donation, not every donor-donation combination. Absence queries must also define NULL behavior.],
  ([Use existence semantics for a yes/no relationship.], [Check join duplication before adding DISTINCT.], [Test absence with nullable inputs.],),
  ([Only membership or absence is needed.], [Child fan-out is high.], [NULL handling is explicit.],),
  ([The output requires child details.], [NOT IN receives NULL without an explicit contract.], [An equivalent spelling is assumed faster everywhere.],),
  ([Plan-dependent probes or scans.], [No extra maintained state.], [Transient plan structures.], [Membership must preserve NULL and duplicate semantics.], [Inspect plans under measured fan-out regimes.],),
  flow-picture(([Parent candidate], [Related-row test], [Keep or reject parent],), "Membership and absence: the operation path.", [Parent candidate → Related-row test → Keep or reject parent. Arrows show logical order, not a measured duration.]),
  ([Return exactly the required parent membership.], [Parent and child rows.], [Read execution resources.], [None for equivalent SQL spellings.], [Membership answer.], [The selected snapshot can differ from later writes.], [Wrong NULL behavior or changed output grain.], [Sweep fan-out with matched answers and negative controls.],), "case-charity-membership", measured: true, example: [
*Illustrative SQL:* these generic table names show the contract, not a benchmark implementation.
```sql
SELECT p.id
FROM parent AS p
WHERE EXISTS (
  SELECT 1 FROM child AS c
  WHERE c.parent_id = p.id
);
```
Use NOT EXISTS for absence. A nullable NOT IN input can make a predicate unknown instead of true.
])

#strategy("Aggregation and per-parent top-N", "theory-select-aggregate", [Choose between repeated probes and shared processing according to the requested parent set.],
  [A page requests a total or latest donation for each donor. Processing all donations may waste work when the page asks for few donors.],
  ([Define empty-parent and tie behavior.], [Compare a per-parent probe with a shared aggregate.], [Vary the requested parent set while holding facts fixed.],),
  ([A narrow request has a useful lookup path.], [Shared processing amortizes across many requested parents.], [The selected plan matches the request shape.],),
  ([Every parent causes an expensive repeated scan.], [Pre-aggregation processes mostly unrequested data.], [A top-one query silently discards tied answers.],),
  ([Per-parent probes or shared processing.], [No extra writes for read-time aggregation.], [Intermediate aggregates, sorting or ranking.], [Empty groups and ties need an explicit rule.], [Measure several regimes and query plans.],),
  flow-picture(([Requested parents], [Probe or process together], [Totals or top-N],), "Aggregation and per-parent top-N: the operation path.", [Requested parents → Probe or process together → Totals or top-N. Arrows show logical order, not a measured duration.]),
  ([Aggregate and top-N answers preserve empty groups and ties.], [Child facts plus transient processing.], [Query execution resources.], [None until a stored rollup is introduced.], [Read-time summaries and ranks.], [Maintained aggregates, if added separately.], [An unfavorable requested-parent count or access path.], [Sweep parent count with oracle checks and shuffled measurement order.],), "case-charity-aggregate", measured: true, example: [
#table(columns: (1fr, 1.5fr), inset: 5pt, stroke: 0.4pt + palette.rule,
table.header([Strategy], [Work it performs]),
[Correlated aggregate], [Compute the child answer separately for each requested parent.],
[Pre-aggregate and join], [Group children, then join those summaries to requested parents.],
[LATERAL / ordered probe], [Use a matching order and limit for each parent's requested top rows.])
COUNT of an empty set yields zero. A missing joined summary needs an explicit zero/default rule; top-N ties need their own contract.
])

#strategy("Pagination and access paths", "theory-select-pagination", [A cursor can continue from an ordered key; OFFSET skips through an ordered prefix.],
  [A donation feed needs the next page. Deep page numbers and a stable continuation cursor impose different navigation contracts.],
  ([Define a total order with a tie-breaker.], [Use a cursor carrying the last ordering keys.], [Compare equivalent pages at several depths.],),
  ([Navigation continues from the previous page.], [Ordering keys have a matching access path.], [Cursor semantics match concurrent change behavior.],),
  ([Users need arbitrary page-number jumps.], [Order lacks a stable tie-breaker.], [A shallow-page result is extrapolated to all depths.],),
  ([Continue after a key or scan and discard a prefix.], [Normal index maintenance.], [Cursor state plus access-path entries.], [Define behavior when rows appear, move or disappear.], [Version cursor formats and bound page sizes.],),
  flow-picture(([Ordered rows], [Last-key cursor], [Next matching page],), "Pagination and access paths: the operation path.", [Ordered rows → Last-key cursor → Next matching page. Arrows show logical order, not a measured duration.]),
  ([Return the correct next page under the chosen order.], [Rows, index and cursor keys.], [Query execution resources.], [Only ordinary index maintenance.], [Continuation cursor.], [A cursor can refer to rows changed since its creation.], [Unsupported jumps or an unfavorable engine/depth regime.], [Verify identical pages and sweep depth by engine.],), "case-charity-pagination", measured: true, example: [
*Illustrative ascending order:* the cursor holds both ordering keys.
```sql
WHERE (created_at, id) > (:last_time, :last_id)
ORDER BY created_at, id
LIMIT :page_size
```
The id resolves timestamp ties. OFFSET supports page-number navigation; this cursor supports continuation from a known key. A flat cost per page needs an index that matches the whole order.

*Access paths for a disjunction.* One OR over an indexed column can use that index once. Rewriting it as UNION ALL is safe only when the branches cannot overlap. A predicate written as an expression needs an index on that same expression.
])
