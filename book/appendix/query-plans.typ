#import "../lib/config.typ": *
#import "../lib/evidence.typ": claim-badge
#set heading(numbering: none)

// The plans, their SQL, their provenance and the node-tree checks are generated
// data. `selectbench/cmd/appendixsql` (Study 01) binds each quoted cell's statement
// with the harness's own catalogue and binder into query-plans.sql.json;
// `tools/evidence plans --write` resolves every exhibit in query-plans.spec.json
// against its claim's registry anchor, verifies the result file's sha256, the SQL
// against the cell and its catalogue sources against the producing commit, and
// each same/different node-tree statement, and writes query-plans.json. `go test` in
// tools/evidence fails when that file is stale. Only the reading notes are authored
// here, and they describe what the quoted plan shows, nothing the claim does not.
#let plan-data = json("query-plans.json")
#let exhibits = {
  let m = (:)
  for e in plan-data.exhibits { m.insert(e.id, e) }
  m
}
#let engine-name(t) = (
  "postgresql-1node": "PostgreSQL, single node",
  "yugabytedb-1node": "YugabyteDB, single node",
  "yugabytedb-3node-rf3": "YugabyteDB, three nodes (RF3)",
).at(t)

// The plan a node-tree statement compares against, named by its printed label when
// any exhibit quotes it (looked up by run and axis, so a cross-run comparison never
// borrows this exhibit's label), otherwise by its axis and run.
#let plan-hits(run, axis) = (plan-data.exhibits.filter(x => x.run_id == run)
  .map(x => x.plans.filter(p => p.axis == axis)).flatten())
#let plan-ref(run, axis, here) = {
  let hits = plan-hits(run, axis)
  let name = if hits.len() > 0 { ["#hits.first().label"] } else { box(raw(axis)) }
  if run == here { name } else [#name in run #box(raw(run))]
}

// A claim heading and its badge stay with the note that follows them.
#let plan-claim(title, id) = block(sticky: true, above: 1.4em, below: 0.7em)[
  #heading(level: 4, numbering: none, title)
  #claim-badge(id)
]

// The bound SQL the harness sent after `EXPLAIN (...)`. A rule on the left marks it
// as the query, set apart from the shaded plan beneath it.
#let sql-box(q) = block(
  width: 100%,
  inset: (x: 6pt, y: 5pt),
  radius: 2pt,
  stroke: (left: 2pt + palette.accent, rest: 0.5pt + palette.rule),
  above: 0.3em,
  below: 0.45em,
  breakable: false,
)[
  #set par(justify: false, leading: 0.5em)
  #set text(size: 7.4pt, hyphenate: false)
  #raw(q, lang: "sql")
]

// Statement SQL printed once: a later node-tree check that names the same text
// points back to it instead of repeating it.
#let shown-sql = state("plan-appendix-shown-sql", (:))
#let sql-once(q, where) = {
  context {
    let seen = shown-sql.get()
    if q in seen {
      text(size: 8pt, fill: palette.muted)[Same SQL as in the #seen.at(q) exhibit above.]
    } else {
      sql-box(q)
    }
  }
  shown-sql.update(s => {
    if q not in s { s.insert(q, where) }
    s
  })
}

#let plan-caption(label, detail) = block(above: 0.7em, below: 0.15em, width: 100%, sticky: true)[
  #set par(justify: false)
  #text(size: 9pt, weight: "bold")[#label]
  #h(0.4em) #text(size: 7.5pt, fill: palette.muted)[#detail]
]

#let plan-exhibit(id) = {
  let e = exhibits.at(id)
  let where = engine-name(e.topology)
  block(above: 1.1em, below: 0.4em, width: 100%, sticky: true)[
    #set par(justify: false)
    #set text(size: 8pt, fill: palette.muted)
    #text(fill: palette.ink, weight: "bold")[#engine-name(e.topology)] · run #box(raw(e.run_id)) · tag #box(raw(e.run_tag)) \
    #raw(e.results_path) · sha256 #raw(e.inputs_digest)
  ]
  for p in e.plans {
    plan-caption(p.label, [SQL #raw(p.statement) · parameters #raw(p.params)])
    sql-box(p.sql)
    shown-sql.update(s => {
      if p.sql not in s { s.insert(p.sql, where) }
      s
    })
    block(above: 0.35em, below: 0.15em, width: 100%, sticky: true)[
      #set par(justify: false)
      #text(size: 7.5pt, fill: palette.muted)[Plan · #raw(p.axis) · #p.lines_kept of #p.lines_total lines kept]
    ]
    raw(p.text, block: true, lang: "explain")
  }
  // One line per compared plan and verdict, listing every statement checked against it,
  // then the SQL of each checked statement (and of a compared plan this appendix does
  // not quote), so every "same node tree" sentence can be read against its queries.
  let groups = ()
  for s in e.shapes {
    let key = s.as_run + "|" + s.at("as") + "|" + s.expect
    let i = groups.position(g => g.key == key)
    if i == none {
      groups.push((key: key, target: s.at("as"), as_run: s.as_run, expect: s.expect,
        target-statement: s.as_statement, target-sql: s.as_sql, items: (s,)))
    } else {
      groups.at(i).items.push(s)
    }
  }
  if groups.len() > 0 {
    block(above: 0.6em, below: 0.2em, width: 100%, sticky: true)[
      #set par(justify: false)
      #set text(size: 8.5pt)
      *Checked by the generator.*
      #groups.map(g => [
        #g.items.map(s => [#s.label #text(size: 7.5pt, fill: palette.muted)[(#raw(s.statement))]]).join(", "):
        #if g.expect == "same" [same node tree as] else [different node tree from]
        #plan-ref(g.as_run, g.target, e.run_id).
      ]).join(linebreak())
    ]
    for g in groups {
      for s in g.items {
        if s.label != "this plan" {
          plan-caption(s.label, [SQL #raw(s.statement)])
          sql-once(s.sql, where)
        }
      }
      if plan-hits(g.as_run, g.target).len() == 0 {
        plan-caption([Compared plan, not quoted here], [SQL #raw(g.target-statement)])
        sql-once(g.target-sql, where)
      }
    }
    v(0.6em)
  }
}

#pagebreak(weak: true)
#heading(level: 2, numbering: none, "Query plans behind the SELECT claims") <appendix-query-plans>

Every SELECT cell in the Charity study stored the text of one `EXPLAIN ANALYZE` execution in the run's
result file. It was captured after the cell's answer passed the correctness check and before its timed
samples.
This appendix quotes those plans for the claims in the Charity case's SELECT variations, from
#link(<case-charity-join-plan>)[join and filter] to
#link(<case-charity-placement>)[placement]. The times inside a plan come from that one instrumented
execution. The ratios in the claims come from the timed samples, so the two need not agree.

Each exhibit names the run, the tag of the code that produced it, the result file and its sha256.
The generator, `tools/evidence plans`, accepts a run only if it is one of the claim's registry
anchors and the file's digest equals the anchor's. It also recomputes every "same node tree" and
"different node tree" statement on these pages. A node tree is the plan without its numbers: node
types, relations, indexes and nesting. The run shown is one of three to twenty-four behind each
claim.

*Where the SQL comes from.* A result cell records the name of its catalogue statement and the
parameters it was bound with, not the SQL text. Each plan is therefore preceded by that statement
with its parameters bound, rebuilt by the study's own catalogue and binder: the code that ran it.
The generator refuses the SQL unless the statement and parameters match the cell. It also requires
every file the catalogue is built from to be byte-identical at the run's producing commit. Some
statements carry constant switches, such as `WHERE 1 = 0 OR c.id <= 1`. These let one catalogue
statement serve several arms of a run, and the plans show the planner folding the constant part
away. The membership, aggregate and access SQL reads the `customers` and `orders` tables of the
teaching fixture directly.

#block(breakable: false)[
  #set par(justify: false)
  #set text(size: 9pt)
  #table(columns: (auto, 1fr), inset: 6pt, stroke: 0.4pt + palette.rule,
    [Trimmed], [#plan-data.trim.description],
    [Kept attributes], [#plan-data.trim.kept_attributes.map(a => box(raw(a.trim()))).join(", ")],
    [Full text], [The untrimmed plan is the cell's `Plan` field in the named result file, under the axis printed beside each plan.],
  )
]

*Reading a plan.* A line starting with `->` is an input of the nearest line above it with less
indentation. `actual time=a..b` is milliseconds to the first and the last row, and `rows` is the row
count; both are per loop, and `loops` says how many times the node ran. `Rows Removed by Filter`
counts rows read and then discarded. YugabyteDB reports `Storage ... Read Requests` only when the plan
was captured with its `DIST` option: the join, CTE and placement runs. The access, aggregate and
membership runs used the plain options, so their YugabyteDB plans carry no request counts.

// ------------------------------------------------------------------ join
#heading(level: 3, numbering: none, "Join spellings and CTE fences")

The join, CTE and placement SQL reads two views. `customers` joins `person` to `charity`, and
`facts` reads `donation`. The planner expands both views, so the plans name the three tables.

#plan-claim("Filtered join spellings", "v6-sel-01-inner-join-spellings-same-plan")

The planner flattens every spelling of the filtered join into the same tree. On PostgreSQL the 1%
rank filter reaches the rank index (`Bitmap Index Scan on donation_rank_idx`) below both hash joins.
On YugabyteDB it reaches the range-sharded rank index (`Index Scan using donation_rank_idx`). The
generator compared the other five spellings with these plans: each has the same node tree, on both
engines.

#plan-exhibit("sel01-pg")
#plan-exhibit("sel01-yb")

#plan-claim("Forced MATERIALIZED around the join", "v6-sel-02-materialized-joined-cte-penalty")

The fence is the `CTE e` under a `CTE Scan`. Inside it the join reads all 200,000 donations
sequentially, and the filter outside then discards 198,000 of them (`Rows Removed by Filter`). Left
to the engine, the same CTE ran the `JOIN ... WHERE` plan above; the forced plan's node tree differs
from it.

#plan-exhibit("sel02-pg")
#plan-exhibit("sel02-yb")

#plan-claim("One CTE read twice", "v6-sel-03-cte-reuse-crossing")

MATERIALIZED runs the CTE body once (`CTE e`), and both references read the stored rows (`CTE Scan on
e` and `e_1`). NOT MATERIALIZED copies the body into each reference, so the index scan runs twice: 8
storage read requests on YugabyteDB against 4. At two references each engine's default has
MATERIALIZED's node tree.

#plan-exhibit("sel03-yb")
#plan-exhibit("sel03-pg")

#plan-claim("Inlining for a selective consumer on PostgreSQL", "v6-sel-04-cte-reuse-postgresql-selective-no-crossing")

With one reference and a 1% consumer, MATERIALIZED builds all 2,000 CTE rows and filters 1,983 of
them away afterwards. Inlined, both predicates reach one scan. The planner uses the primary-key index
for the consumer's `donation_id` range and tests the rank as a filter, with no stored intermediate
result.

#plan-exhibit("sel04-pg")

// ------------------------------------------------------------ membership
#heading(level: 3, numbering: none, "Membership and absence")

#plan-claim("JOIN + DISTINCT against EXISTS", "v6-sel-05-membership-join-distinct-penalty")

On PostgreSQL, EXISTS is a `Nested Loop Semi Join`: one index probe per parent, which stops at the
first qualifying order. JOIN + DISTINCT produces 7,240 joined rows and then removes
duplicates in a `HashAggregate`. On YugabyteDB, EXISTS was planned differently. It deduplicates the
qualifying orders by `customer_id` first and then hash-joins, so both YugabyteDB plans scan all 7,314
qualifying orders. No YugabyteDB plan in these runs used an index.

#plan-exhibit("sel05-pg")
#plan-exhibit("sel05-yb")

#plan-claim("IN against EXISTS", "v6-sel-06-membership-in-exists-no-difference")

PostgreSQL turned IN into the same semi-join as EXISTS: the two plans have the same node tree.

#plan-exhibit("sel06-pg")

#plan-claim("NOT IN over a nullable key", "v6-sel-07-nullable-not-in-empty")

No exhibit. The NOT IN over a nullable child key is a negative control: the harness records the rows
it returned (none, at every fanout), not its plan. NOT IN over a NOT NULL subquery is shown in the
next exhibit.

#plan-claim("Anti-join spellings on PostgreSQL", "v6-sel-08-anti-join-not-exists-postgresql")

NOT EXISTS becomes a `Nested Loop Anti Join` that probes the order index once for each of the 100
parents. LEFT JOIN ... IS NULL hash-joins all 7,314 qualifying orders and then discards the 7,240
joined rows. NOT IN over a NOT NULL subquery loads the qualifying orders into a hashed subplan. Both
alternatives read every qualifying order.

#plan-exhibit("sel08-pg")

// ------------------------------------------------------------- aggregates
#heading(level: 3, numbering: none, "Totals and the top row per parent")

#plan-claim("Correlated SUM on PostgreSQL", "v6-sel-09-aggregate-correlated-postgresql-crossover")

The correlated spelling runs two scalar subqueries per parent (`SubPlan 1` and `SubPlan 2`), each an
index probe. For one parent that is two probes. For 10,000 parents each subplan runs 10,000 times
(`loops=10000`). Join-then-GROUP-BY instead reads the 20,000 orders once and aggregates them in a
`HashAggregate`.

#plan-exhibit("sel09-pg")

#plan-claim("Correlated SUM on YugabyteDB", "v6-sel-10-aggregate-yugabytedb-no-reversal")

The structure is PostgreSQL's, with each aggregate split into a partial aggregate in storage and a
`Finalize Aggregate`. In this execution each probe took about 0.44 ms (`SubPlan 1`) and 0.22 ms
(`SubPlan 2`), against below 0.001 ms on PostgreSQL, and the 20,000 probes took 6.8 s in all. The plan
was captured without `DIST`, so it shows no storage requests, and the claim's per-request explanation
remains a hypothesis.

#plan-exhibit("sel10-yb")

#plan-claim("Pre-aggregating derived table", "v6-sel-18-aggregate-preagg-derived-table")

For one parent, join-then-GROUP-BY reads that parent's two orders through the index. The
pre-aggregating derived table first groups all 20,000 orders into 10,001 groups (`HashAggregate`) and
only then joins the one parent: the parent filter does not reach inside the derived table. At 10,000
parents both forms read every order once. YugabyteDB's plan at 10,000 parents also aggregates all
orders before the join.

#plan-exhibit("sel18-pg")
#plan-exhibit("sel18-yb")

#plan-claim("LATERAL probe against ROW_NUMBER()", "v6-sel-11-topn-lateral-probe-crossover")

LATERAL runs an index probe, a sort and a limit for each parent: once for one parent, and 10,000 times
for 10,000 parents (`loops=10000`, about 0.47 ms each on YugabyteDB). ROW_NUMBER() reads the 19,936
orders once, sorts them, and keeps the first row of each parent (`Run Condition`).

#plan-exhibit("sel11-yb")

#plan-claim("DISTINCT ON on PostgreSQL", "v6-sel-12-topn-distinct-on-postgresql")

Both spellings read the same 19,936 orders through `orders_customer_idx`, already in customer order,
and add an `Incremental Sort`. DISTINCT ON then keeps the first row per customer with `Unique`.
ROW_NUMBER() computes a window aggregate and filters it in a subquery.

#plan-exhibit("sel12-pg")

// ------------------------------------------------------------- access paths
#heading(level: 3, numbering: none, "Predicates and pagination")

#plan-claim("OR against UNION ALL", "v6-sel-13-or-not-beaten-by-union")

With selective arms, PostgreSQL serves OR with one `BitmapOr` over two index scans. UNION ALL runs
the two bitmap scans separately and appends them before the sort. With broad arms, OR reads the
primary-key index in id order and filters, while UNION ALL runs two index scans and sorts their
4,652 rows. On YugabyteDB the hash-sharded amount index cannot serve a range, so OR is one sequential
scan with the predicate pushed to storage (`Storage Filter`).

#plan-exhibit("sel13-pg")
#plan-exhibit("sel13-yb")

#plan-claim("Expression index", "v6-sel-14-expression-index-postgresql")

Without a matching index, the predicate on `amount + 0` is a filter over the primary-key scan, which
removes 11,749 rows. With the expression index present, the same SQL becomes a `BitmapOr` over
`orders_amount_expr_idx`. For the broad shape the index changes nothing: its node tree is the same
with or without the index.

#plan-exhibit("sel14-pg")

#plan-claim("Keyset against OFFSET", "v6-sel-15-keyset-offset-crossover")

At depth 10,000, OFFSET sorts all 12,000 orders and walks 10,005 sorted rows to return five. On
PostgreSQL the keyset predicate becomes a three-arm `BitmapOr`, including the `amount IS NULL` arm,
followed by a top-N heapsort over 2,000 rows. On YugabyteDB, keyset is a sequential scan with the
predicate as a `Storage Filter`, also followed by a top-N heapsort.

#plan-exhibit("sel15-pg")
#plan-exhibit("sel15-yb")

// ---------------------------------------------------------------- placement
#heading(level: 3, numbering: none, "Placement on three nodes")

#plan-claim("Splitting into three tablets", "v6-sel-16-placement-splitting-effect")

Placement did not change the plan: the single-tablet and the distributed runs have the same node
tree. The difference is inside the donation scan. In these two executions it made 176 storage read
requests in 80 ms on one tablet, and 59 requests in 17 ms across three tablets. These are single
instrumented executions; the claim's ratios come from the timed samples.

#plan-exhibit("sel16-tablet")
#plan-exhibit("sel16-distributed")

#plan-claim("Colocated database", "v6-sel-17-placement-colocated-vs-distributed")

The joins are the same as in the distributed run, but the colocated run read `person` and `charity`
through their primary-key indexes where the distributed run scanned them, so the node trees differ.
The donation table was scanned sequentially in every layout. In the colocated run that scan made 176
read requests, as on the single tablet.

#plan-exhibit("sel17-colocated")
