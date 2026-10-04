# Owner request: SELECT templates with examples and measured boundaries

Date: 2026-10-02. Session capability: HIGH (raw answer: HIGH (leader)).
Planner: GPT-6, Codex desktop; precise model variant and selected effort are not exposed.
Baseline local main: 1a178ce. Direct task request; no brainstorm was started.

## Verbatim implementation request

> I would like for my book to explore different ways to do the select and show examples and numbers when they thrive or die. so a reader can have an intuition of what template to use in a project.
> please create tasks for this.

The preceding question asked for an example of joining then filtering versus filtering then joining. The supplied answer illustrated customers' orders worth more than 100 in six full-row INNER JOIN forms and a seventh EXISTS form that answers a different question.

## Supplied teaching example, preserved as context

Example data (not a measurement): customers (1, Ana), (2, Ben); orders (10, 1, 50), (11, 1, 200), (12, 2, 80).
The six full-row forms return (Ana, order 11, amount 200):
1. INNER JOIN followed by WHERE amount > 100.
2. A filtered orders derived table joined to customers.
3. WITH expensive_orders naming the filtered input.
4. WITH customer_orders naming joined rows, followed by a filter.
5. INNER JOIN with amount > 100 in ON.
6. A comma join with the relationship and amount predicates in WHERE.
The EXISTS form returns the parent (id 1, Ana) once, even with several qualifying orders. It must not be ranked against full child-row output as an equivalent answer.

The supplied answer also distinguished written query structure from optimizer execution, CTE folding/materialization, and outer-join predicate placement. Preserve those distinctions, add exact runnable SQL/outputs, and measure rather than assume performance.

## Authorization and scope

Create durable executable repository tasks, not only a chat list. The task publication itself does not run experiments or rewrite the book. Downstream tasks cover protocol, implementation, measurements, evidence review and book publication with explicit dependencies and HIGH reviews. They use existing relational study data, not an unreviewed new datastore/domain study. Existing book/style approvals and immutable evidence remain authoritative.

## Documentation inputs to verify against pinned engines

- [PostgreSQL 17 table expressions and joins](https://www.postgresql.org/docs/17/queries-table-expressions.html).
- [PostgreSQL 17 WITH and materialization](https://www.postgresql.org/docs/17/queries-with.html).
- [PostgreSQL 17 subquery expressions](https://www.postgresql.org/docs/17/functions-subquery.html).

These document semantics/features, not benchmark results. Protocol authors must verify support on the pinned YugabyteDB image independently. The supplied answer's vbnet fences were display formatting; new examples use sql fences.

TASK COMPLETE
