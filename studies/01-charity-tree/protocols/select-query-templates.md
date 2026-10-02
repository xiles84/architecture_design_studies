# Study protocol — Study 01 SELECT query templates

| Field | Value |
|---|---|
| Study | `studies/01-charity-tree/` (Study 01, charity tree), additive SELECT experiment |
| Protocol revision | 1 |
| Planner | deepseek-flash (DeepSeek), Deep Code CLI, effort not exposed; session capability HIGH, role planner |
| Starting main revision | `e4059822` (the claim commit for `task-20261002T131729Z-select-protocol`) |
| Status | ready (pre-measurement) |
| Required tag on integration | `study-01/select-v0-protocol` |
| Numbering | Study 01 is the charity tree; the next free study number is 10 (`studies/06..09` exist). This program adds no new study number: it is an additive SELECT experiment inside Study 01. |
| Owner inputs | [`docs/select-query-templates/REQUEST.md`](../../docs/select-query-templates/REQUEST.md), [`docs/select-query-templates/PLAN.md`](../../docs/select-query-templates/PLAN.md) |
| Method | [`docs/methodology.md`](../../docs/methodology.md), template [`docs/templates/STUDY_PROTOCOL.md`](../../docs/templates/STUDY_PROTOCOL.md) |

A protocol is written **before** any code that measures. Nothing in this file is a measured
number. Every figure below is an exact logical answer to a deterministic teaching fixture, or a
hypothesis to be tested. Performance numbers may enter the book only from registered cells of
the measuring tasks (3–6), each tagged and digested.

## 1. Question and invariant

A reader of the book has a required answer and a choice of readable SQL spellings for it. The
architectural question is a **decision**: *for a given answer contract, which SELECT template
should the reader reach for, and at what regime does its advantage reverse?* The answer must be
stated per **result contract** — bags of joined rows, sets of parent identities, aggregates,
parent-preserving optional matches, per-parent top-N, and ordered pagination are different
questions and comparing across them establishes nothing.

**Invariant (INV-SEL).** For every controlled pair, both statements must return *the same answer
under the pair's declared contract*: identical bag multiplicities for full-row contracts,
identical parent-identity sets for membership contracts, identical aggregate values and
empty-parent semantics for aggregate contracts, identical ordered output under the declared total
order for pagination/top-N contracts. A fast cell whose answer differs is a failed cell, and a
statement that answers a different question is **not** a competing spelling (the REQUEST's
EXISTS example is the model: it returns parent id 1 once and must never be ranked against a
full-child-row output).

The second invariant is provenance: **INV-PROV** — every timing cell names the pinned image
digest, dataset seed, index set, statistics state and prepared-plan policy that produced it, and
carries an EXPLAIN plan and its normalized signature beside the time.

## 2. Prediction (recorded before any measurement)

Written before the first run; not edited afterwards. Verdicts are recorded in the task-7 signed
analysis. `Recorded at` is the commit that introduces this file (the protocol commit; the hash is
recorded in the task's `RESULT.md`), which precedes all measurement code.

| Id | Comparison | Expected direction | Expected mechanism | Expected crossover / regime | Falsifier |
|---|---|---|---|---|---|
| P1 | the six inner-join spellings of one filtered fact, pairwise | **no resolved difference** in most regimes | the planner folds a one-use CTE and pushes the filter; the written order is not the physical order | a difference appears only where a CTE materializes or an index cannot be used | a pair differs beyond spread in the same direction in ≥3 runs at ≥3 trials with an identified plan difference |
| P2 | derived table vs plain join (same predicate) | no resolved difference | identical relational expression after subquery pull-up | — | a consistent, plan-explained difference |
| P3 | default CTE vs `MATERIALIZED` CTE, selective one-use predicate | MATERIALIZED is **slower** | the materialization boundary blocks pushing the selective predicate into the scan | reverses when the CTE is reused several times and the inner work is expensive | MATERIALIZED is never slower where the predicate is selective, or never faster under high reuse |
| P4 | `MATERIALIZED` vs `NOT MATERIALIZED` under repeated reuse of an expensive subquery | `MATERIALIZED` is **faster** above a reuse threshold | the expensive work is paid once instead of per reference | a bounded crossing vs reuse count | no crossing across ≥5 reuse levels, or the direction never reverses |
| P5 | LEFT JOIN with the child predicate in `ON` vs in `WHERE` | **different answers**, not a speed pair | the `WHERE` predicate on the child drops unmatched/left rows | n/a (answer contract, not performance) | the two return the same bag on a fixture with unmatched parents and NULL child keys |
| P6 | membership: EXISTS vs `IN` vs JOIN+`DISTINCT` on parent id | EXISTS/`IN` avoid row multiplication; `DISTINCT` costs de-duplication | the full join emits one row per qualifying child; membership is a set | the row-multiplication cost rises with child fanout | `DISTINCT` is never slower at high fanout, or EXISTS/`IN` differ in answer from the membership set |
| P7 | absence: `NOT EXISTS` vs LEFT JOIN `IS NULL` vs `NOT IN` | `NOT EXISTS`/anti-join are safe; nullable `NOT IN` can return **nothing** | `NOT IN` over a NULL yields UNKNOWN for every row | n/a unless the column is NOT NULL by contract | `NOT IN` over a nullable subquery returns the same set as `NOT EXISTS` |
| P8 | parent aggregate: post-join aggregate vs pre-aggregation vs correlated scalar | pre-aggregation avoids the join multiplier; correlated pays per parent | a join to two child collections multiplies rows and inflates SUM | correlated wins at small parent sets, loses at large | no reversal across ≥5 parent-selection levels, or totals disagree with the oracle |
| P9 | per-parent top-N: lateral/correlated indexed probe vs `ROW_NUMBER` vs `DISTINCT ON` (N=1) | the indexed probe wins at small parent sets and loses at large | one ordered index probe per parent vs one global sort | bounded crossing vs number of parents | no crossing across ≥5 levels, or `ROW_NUMBER` vs `RANK` tie behaviour differs from the declared contract |
| P10 | predicate alternatives: OR vs overlap-safe `UNION ALL` | OR can block an index; `UNION ALL` can use two index arms | a single index often cannot serve a disjunction across expressions | reverses when the predicate becomes overlap-unsafe or an index covers it | `UNION ALL` changes the answer (overlap) or OR is never worse on a non-indexable exact predicate |
| P11 | pagination: `OFFSET` vs keyset under a fixed snapshot | keyset cost is flat in page depth; OFFSET scans the skipped rows | OFFSET discards `N` rows each request | keyset advantage grows with page depth | OFFSET is flat across ≥5 page depths, or keyset returns a different page under the fixed snapshot |
| P12 | selectivity sweep across the six spellings | the spellings converge as the predicate becomes non-selective | the plan collapses to a scan regardless of spelling | a difference only at high selectivity | a spelling pair differs beyond spread at 50–90% selectivity |

Each prediction gets `confirmed`, `falsified` or `unresolved`. A non-appearing reversal (P4, P8,
P9, P11) is itself a result and is reported as one.

## 3. Designs, pairs and the three forms of evidence

Every comparison is `isolate`, `interaction` or `package`. The SELECT program is almost entirely
**isolate**: two spellings of one fact, with the same logical schema, indexes, data, settings and
answer contract, differing in exactly the SQL text (and, where declared, one physical factor).

| Comparison | Form | What changes | Attribution rule |
|---|---|---|---|
| spelling A vs spelling B, same predicate | isolate | SQL structure only | the delta may be attributed to the spelling |
| spelling × selectivity (or reuse, parent count, fanout) | interaction | spelling and a workload axis | neither alone may be attributed; report the crossing |
| spelling + new index (or + prepared statement, + rewritten schema) | package | syntax and a physical decision | **nothing may be attributed to the syntax**; the index is the changed factor and must be a separate isolate arm |

**Result contracts (answer groups).** Comparisons are valid only inside one contract.

1. **Full joined rows (bag).** Multiplicities matter. The contract is the bag of `(parent, child)`
   rows satisfying the predicate.
2. **Parent membership (set).** The contract is the set of parent identities with ≥1 qualifying
   child; each parent appears exactly once.
3. **Parent-preserving optional match.** The contract is the bag of parents, with qualifying child
   columns or NULLs; every parent appears at least once. (This is a different contract from
   full inner rows and from membership.)
4. **Parent aggregates.** The contract is one row per parent with declared SUM/COUNT semantics,
   including empty-parent and NULL-child rules.
5. **Per-parent top-N.** The contract is at most N child rows per parent under a declared total
   order, with declared tie behaviour.
6. **Ordered pagination.** The contract is a page of rows under a declared total order and a fixed
   snapshot, with keyset or OFFSET addressing.

## 4. Correctness gate and negative controls

The gate compares each cell's answer against an **independent Go oracle** computed from the
deterministic logical dataset, not one query against another. A failing cell reports no timing.
The oracle compares bag multiplicities, parent-identity sets, aggregate values and ordered
sequences as the contract requires.

| Control | Deliberately wrong how | Invariant it must violate | Must fire in | Read-path control? |
|---|---|---|---|---|
| `X-nullmod` | an inner-join filter written so a NULL child predicate drops a parent that the optional-match contract keeps | INV-SEL (optional-match contract) | every optional-match phase | yes |
| `X-mult` | plain JOIN + `DISTINCT` on parent **name** instead of parent id, on a fixture where two parents share a name | INV-SEL (membership contract) | membership phase | yes |
| `X-overlap` | `UNION ALL` where the two predicate arms overlap | INV-SEL (bag contract) | predicate-alternative phase | yes |
| `X-tie` | `ROW_NUMBER` without a total-order tiebreaker, reported as if deterministic | INV-SEL (top-N contract) | top-N phase | yes |
| `X-stale` | a plan/timing captured against a changed statistics snapshot while the answer is claimed for the declared one | INV-PROV | planning phase | no |

Rules carried from methodology 5 and 5a: every invariant has a control expected to break it; a
read-only SELECT phase that cannot exercise a control records `control_absent`, and then the other
designs' clean answers are **not** evidence of their correctness; invariants are two-sided (a
"never over/under" aggregate check is paired with the opposite direction).

## 5. Workload dimensions and levels

| Dimension | Levels | Why these levels |
|---|---|---|
| Selectivity (fraction of children qualifying) | 0.01%, 0.1%, 1%, 10%, 50%, 90% | locate where a spelling's selective access stops mattering |
| CTE reuse count (P4) | 1, 2, 4, 8, 16 | bound the materialization crossover |
| Parent selection (fraction of parents requested) | 0.1%, 1%, 10%, 100% | bound the per-parent-probe crossover (P8, P9) |
| Child fanout / skew | uniform; single hot parent; Zipfian | row-multiplication and probe cost |
| Parent-set size for correlated/lateral | 10, 100, 1 000, 10 000 | crossover in P9 |
| Page depth | 1, 10, 100, 1 000, 10 000 | OFFSET vs keyset |
| Row width | narrow, wide (blob column) | plan choice between index-only and heap access |

Rules: prefer levels chosen to locate a crossover; a sweep claims a crossing only with ≥5 levels
and a stated interval; change the shape, hold the volume (methodology 8a) — when a regime removes a
limitation, keep the total work comparable.

## 6. Regimes that remove each limitation

| Design | Its limitation | Regime that removes it | Designs re-run |
|---|---|---|---|
| materialized CTE | predicate cannot be pushed through the materialization boundary | CTE reused many times with expensive inner work | P4 pair at reuse 8/16 |
| correlated/lateral probe | per-parent probe cost dominates | small parent set, indexed child | P9 pair at 1% selection |
| OFFSET pagination | skipped rows rescanned | shallow page | P11 pair at depth 1 |
| `NOT IN` | NULL poisoning | column declared NOT NULL | P7 pair on a NOT NULL variant |
| index-only plan | heap fetches / visibility map | freshly vacuumed table | run each index arm once post-`VACUUM (ANALYZE)` |

## 7. Replication, order and variability

| Field | Plan |
|---|---|
| Whole runs per cell | ≥3 independent whole runs for any argued difference |
| Trials per cell within a run | ≥3 |
| Design order | seeded randomized per run; interleaved where several spellings share a phase |
| Position control | the same spelling first and last in a phase, on at least one run per task |
| Across-run variability | median plus across-run spread; a difference below measured variability is `unresolved` |
| Tier | `run-level` (same machine). No `hardware-independent` claim; one laptop establishes neither independent hardware nor real network behaviour |

Measurement is serial under `ads-run-lock`; results and reports are produced by the runner and
tagged `run/01-charity-tree/<run-id>`.

## 8. Coverage of the owner's required comparisons

- **Cluster and client sizing:** the protocol calculates current host/VM, client and system
  reserves, per-node and client budgets, readers/writers/connections and attainable load from the
  live environment (not from old reports). Topologies: PostgreSQL single-node, YugabyteDB
  single-node, YugabyteDB 3-node RF=3, with **distinct per-node and equal-total arms** where the
  topology is multi-node. A PostgreSQL replication/read-routing experiment is **not** implied by
  multiple PostgreSQL containers — that is an explicit gap.
- **Rollup, rolldown and embedding:** the SELECT families are syntax-only isolate pairs, so these
  mechanisms are **not applicable** to tasks 3–6; a justified non-applicability is recorded. The
  existing Study 01 designs `d4_rollup_trigger`, `d5_rollup_app`, `d6_embedded_jsonb`,
  `d9_embedded_hybrid`, `d10_embedded_hybrid_locked`, `d11_copied_key` remain the source for
  information-placement cost; if a measuring task varies one, it must use a **separate pair** with
  its own maintenance/storage/correctness arms and may not let the syntax arm claim the effect.
- **Colocated and non-colocated data:** the existing `d7_yb_child_colocated` / `d8_flattened_nofk`
  (and `y1`/`y2` reference) pairs are the placement evidence. Any 3-node arm must record verified
  tablet/leader distribution or be an explicit `coverage_gap`; single-host container proximity is
  not placement evidence.
- **Optimistic and pessimistic concurrency:** **non-applicable** — the SELECT program does not
  change arbitration. If a mixed read/write phase is added, the writer implementation is held
  fixed with one correct contract; varying it would require both strategies, contention, aborts
  and correctness controls and is a different study.
- **Any other decision this study varies:** SQL spelling (isolate), and — only where explicitly
  paired — one physical factor (index arm), prepared-statement policy, or statistics state.

## 9. Boundaries and non-applicability

This program will **not** establish: fastest-syntax rankings outside a stated regime; any claim
that written query structure equals execution order; behaviour on engines or versions other than
the pinned images; distributed-network effects (one host); persistent materialized-read-model
design (Study 08) or recursive hierarchy queries (Study 09); new concurrency strategies; a
PostgreSQL cluster/read-routing result from multiple containers; and any number for the
customers/orders fixture, which is a teaching analogue only. A crossover that never appears is a
reported outcome, not a failure.

## 10. Acceptance

- Every prediction carries a verdict; every control's firing is recorded per topology/phase.
- Every claim carries its evidence form and, if a package, its attribution status.
- Crossover/boundary results state dimension, levels and crossing interval with spread.
- Failed and timed-out cells are named and never reported as winners.
- Every result names environment, pinned images, run tag and inputs digest; every book number
  resolves to a registered claim in the active evidence package.
- Teaching SQL never carries a performance claim; the oracle never compares one query to another.

## 11. SQL example and contract atlas

Two fixtures. **Fixture A** is the REQUEST's supplied example, whose six inner-join forms must
return exactly one row. **Fixture B** adds the edge cases (NULL, duplicate child matches,
unmatched parent, ties, empty sets). No number below is a measurement.

Fixture A:
```
customers(id, name): (1,'Ana'), (2,'Ben')
orders(id, customer_id, amount): (10,1,50), (11,1,200), (12,2,80)
```
Fixture B:
```
customers(id, name): (1,'Ana'), (2,'Ben'), (3,'Cara')
orders(id, customer_id, amount): (10,1,50), (11,1,200), (12,2,80), (13,2,NULL), (14,1,200)
```

### 11.1 The seven REQUEST forms, exact outputs

Let `P = amount > 100`.

| # | Form | Fixture A output | Fixture B output |
|---|---|---|---|
| 1 | `FROM customers c JOIN orders o ON o.customer_id=c.id WHERE o.amount>100` | `(Ana,11,200)` | `(Ana,11,200),(Ana,14,200)` |
| 2 | `JOIN (SELECT * FROM orders WHERE amount>100) x ON x.customer_id=c.id` | `(Ana,11,200)` | `(Ana,11,200),(Ana,14,200)` |
| 3 | `WITH e AS (SELECT * FROM orders WHERE amount>100) SELECT … JOIN e …` | `(Ana,11,200)` | `(Ana,11,200),(Ana,14,200)` |
| 4 | `WITH co AS (SELECT … FROM customers JOIN orders …) SELECT … FROM co WHERE amount>100` | `(Ana,11,200)` | `(Ana,11,200),(Ana,14,200)` |
| 5 | `JOIN orders o ON o.customer_id=c.id AND o.amount>100` | `(Ana,11,200)` | `(Ana,11,200),(Ana,14,200)` |
| 6 | `FROM customers c, orders o WHERE c.id=o.customer_id AND o.amount>100` | `(Ana,11,200)` | `(Ana,11,200),(Ana,14,200)` |
| 7 | `WHERE EXISTS (SELECT 1 FROM orders o WHERE o.customer_id=c.id AND o.amount>100)` | `(1,Ana)` | `(1,Ana)` |

Forms 1–6 are one **bag** contract (rows 1 = 1; rows 2 = 2 on fixture B). Form 7 is a **membership**
contract (one parent identity, even when Ana has two qualifying orders). Form 7 must never be
ranked against forms 1–6.

### 11.2 Optional child matches (P5) — two different answers

Qualifying children `amount>100`, all parents preserved:
```
LEFT JOIN orders o ON o.customer_id=c.id AND o.amount>100
```
Fixture B → `(Ana,11,200),(Ana,14,200),(Ben,NULL,NULL),(Cara,NULL,NULL)` (4 rows).
Moving `o.amount>100` to `WHERE` → `(Ana,11,200),(Ana,14,200)` (2 rows): Ben and Cara are dropped.
The `WHERE` form answers the **inner** question, not the optional-match contract.

### 11.3 Membership and absence, exact outputs

| Contract | Statement | Fixture B |
|---|---|---|
| membership, ≥1 child `amount>100` | `WHERE EXISTS (…)` | `{Ana}` (1 row) |
| membership via join | `SELECT DISTINCT c.id, c.name … JOIN … WHERE amount>100` | `{Ana}` — note `DISTINCT` on `(id,name)`, **never** on `name` alone |
| absence | `WHERE NOT EXISTS (SELECT 1 FROM orders o WHERE o.customer_id=c.id AND o.amount>100)` | `{Ben, Cara}` |
| absence, left-anti | `LEFT JOIN orders o ON o.customer_id=c.id AND o.amount>100 WHERE o.id IS NULL` | `{Ben, Cara}` |
| absence, nullable `NOT IN` | `WHERE c.id NOT IN (SELECT customer_id FROM orders WHERE amount>100)` | `{Ben, Cara}` **here** — but if the subquery returns any NULL, every row becomes UNKNOWN and the result is **empty**: the trap. A NOT NULL contract removes it. |

`COUNT(*)` vs `COUNT(amount)` on fixture B all-orders: Ana `3` vs `3`; Ben `2` vs `1` (the NULL is
skipped by `COUNT(amount)`); Cara has no order rows, so an `INNER` aggregate emits nothing while a
`LEFT` aggregate emits `(Cara, 0, NULL)` — the empty-parent semantics depend on the contract.

### 11.4 Aggregates and per-parent top-N, exact outputs

Aggregate of **all** orders per customer (not filtered): Ana `SUM=450, COUNT(*)=3`; Ben
`SUM=80, COUNT(*)=2, COUNT(amount)=1`; Cara `0 rows` (inner) or `SUM=NULL, COUNT=0` (left).

Per-parent top-1 by `amount DESC` (fixture B, all orders), with tie between orders 11 and 14:
```
ROW_NUMBER() OVER (PARTITION BY customer_id ORDER BY amount DESC, id ASC) = 1   → Ana gets id 11
RANK()       OVER (PARTITION BY customer_id ORDER BY amount DESC)             = 1   → Ana gets id 11 and id 14
DISTINCT ON (customer_id) … ORDER BY customer_id, amount DESC, id ASC         → Ana gets id 11
```
`ROW_NUMBER`/`DISTINCT ON` need a total-order tiebreaker to be deterministic; `RANK` intentionally
returns both tied rows. The declared contract must pick one; the oracle enforces it.

### 11.5 Pagination, exact output shape

Under a fixed snapshot, ordered by `(amount DESC, id ASC)`:
```
page 1 OFFSET 0 LIMIT 2 → (200,11),(200,14)      page 2 OFFSET 2 LIMIT 2 → (80,12),(50,10)
keyset after last (amount,id)=(200,14), LIMIT 2 → (80,12),(50,10)
```
A NULL `amount` sorts last under `DESC` (default `NULLS FIRST` is for `ASC`; `DESC` puts NULLs
last), so the NULL-only order 13 is unorderable by amount and must be excluded or given a
declared position. A mutable offset shifts pages under concurrent writes; the snapshot/keyset
contract is what makes the page reproducible.

### 11.6 Mapping the teaching fixture onto Study 01

`customers → charity`/`person`, `orders → donation`, `amount → donation.amount`,
`customer_id → donation.person_id`. The same seven forms and six contracts apply verbatim; Study
01's `d1_normalized_minimal`/`d2_normalized_indexed`/`d3_flattened_fk` provide the indexed
access paths and the `d4`/`d5`/`d6`/`d9`/`d10`/`d11` designs the information-placement arms.

## 12. Bounded matrices and LOW handoffs

Each measuring task runs under this protocol and produces run manifests, results, plans, generated
reports and a signed analysis section. Common rules: additive SELECT mode in the Study 01 harness
(image/output namespace separate from historical runs); independent Go oracle; drain result rows;
randomize keys; plan signature normalized (strip costs/rows/loops, keep node types and index
names); timeout per cell with the cell recorded as `timeout`; every DB start under `ads-run-lock`;
≥3 whole runs and ≥3 trials; seeded interleaving and a first/last position control.

| Task | Bounded matrix | Controls | Acceptance | Escalation |
|---|---|---|---|---|
| `select-harness` (2) | builds the oracle + catalogue + report renderer; no matrix | `X-nullmod`, `X-mult` unit tests | all six contracts reproduce fixture A/B exact outputs; oracle passes on both engines | a contract that cannot be expressed on the pinned engines |
| `select-join-cte` (3) | 6 spellings × 6 selectivity levels × {default, MATERIALIZED, NOT MATERIALIZED} at reuse ∈ {1,8}; PG-single + YB-single | `X-stale` | P1–P4 verdicts; plans for every cell | a spelling/feature unsupported on the pinned YB image |
| `select-existence` (4) | {EXISTS, IN, JOIN+DISTINCT} × fanout {1,10,100} ; {NOT EXISTS, LEFT-IS-NULL, NOT IN} × {NULL, NOT NULL} | `X-mult`, nullable-`NOT IN` fixture | P6–P7 verdicts; answers identical within contract | `IN` over a NULL contract cannot be made safe |
| `select-aggregate-topn` (5) | {post-join, pre-agg, correlated} × parent-selection {0.1%,1%,10%,100%}; {lateral, ROW_NUMBER, DISTINCT ON} × parent-count {10,100,1k,10k} | `X-tie` | P8–P9 verdicts or reported non-crossing; totals match oracle | a plan that cannot be forced/captured on YB |
| `select-access-paths` (6) | OR vs UNION ALL × {indexable, expression, cast} ; OFFSET vs keyset × page depths {1,10,100,1k,10k} | `X-overlap` | P10–P11 verdicts; page equality under snapshot | no snapshot mechanism available on an engine |
| `select-synthesis` (7, HIGH) | reads 3–6; issues corrections; signed analysis | — | every prediction verdict; mechanism companions | a result that changes the accepted matrix |
| `select-book` (8, LOW) | concept-first SELECT material + worked atlas | figure/token checks | book builds, pages inspected, claims resolve | Edition-3 pilot not integrated |

Topology rules: each measuring task includes **PG single-node and YB single-node** at minimum;
**YB 3-node RF=3** where the comparison is topology-relevant (per-node and equal-total arms
separated, balanced endpoints, measured throttling). Any 3-node arm that cannot be executed on
the fixed host is recorded as a `coverage_gap`, not silently dropped.

## 13. Feature support against the pinned engines

To be **verified in task 2/3 against the pinned images**, not recalled: CTE inlining and
`MATERIALIZED`/`NOT MATERIALIZED` support and defaults (PostgreSQL 17 changed the default for
one-use CTEs; YugabyteDB's YSQL must be checked independently); `DISTINCT ON`; `LATERAL`;
`NULLS FIRST/LAST`; `NOT IN` NULL semantics; available snapshot mechanisms (`REPEATABLE READ`,
`SET TRANSACTION SNAPSHOT`); and index-only-scan/visibility support. Until verified, each is
marked `unverified` and the corresponding matrix row is gated. Documentation inputs:
PostgreSQL 17 table expressions, `WITH` queries and subquery expressions (per REQUEST.md); the
pinned YugabyteDB image is the authority for YB behaviour.

## 14. Provenance and amendment

This protocol is revision 1, integrated under tag `study-01/select-v0-protocol`. Any narrowing of
the matrix, change to a contract, prediction or control requires a **dated amendment** in this
file reviewed by HIGH before the affected measurement runs. LOW workers may choose local,
reversible implementation details and log them; changes to meaning, scope or evidence gates do
not belong to a worker.
