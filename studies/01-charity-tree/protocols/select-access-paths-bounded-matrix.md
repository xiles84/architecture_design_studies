# SELECT access paths — bounded measurement plan (LOW worker)

| Field | Value |
|---|---|
| Task | `task-20261002T131729Z-select-access-paths` |
| Session capability / role | LOW / executor (deepseek-flash, opencode CLI) |
| Protocol | [`select-query-templates.md`](select-query-templates.md) §12 row `select-access-paths` (revision 1, tag `study-01/select-v0-protocol`) |
| Harness | additive SELECT harness, tag `study-01/select-v1-harness` |
| Status | committed before any run; measurements not yet taken |

This plan **implements** the reviewed matrix; it does not amend the protocol. Any
change to a contract, a prediction or an accepted coverage would be a HIGH
decision and is not made here. The premeasurement commit hash is the commit that
introduces this file and is recorded in the task `RESULT.md`.

## 1. What is measured (P10, P11)

### 1.1 P10 — OR vs overlap-safe UNION / UNION ALL

One logical disjunction, three SQL spellings, three predicate shapes. The arms
are `amount <= 100` (A) and `amount >= 400` (B). They are **disjoint by
construction** (no value can be ≤100 and ≥400), so `UNION ALL` answers exactly
the OR bag and is a competing spelling; `UNION` de-duplicates and is the
answer-safe fallback when the arms might overlap. The independent Go oracle
(`PredicateAlt`) computes the bag from the fixture, and `OverlapCount` records
the disjointness proof per shape.

| Shape | Arm A SQL | Arm B SQL |
|---|---|---|
| `indexable` | `amount <= 100` | `amount >= 400` |
| `expression` | `(amount + 0) <= 100` | `(amount + 0) >= 400` |
| `cast` | `(amount)::bigint <= 100` | `(amount)::bigint >= 400` |

Statements (group `sel_predicate_alt`, bag of `(order id, amount)`): `s18_or_predicate`,
`s19_union_all_predicate`, `s20_union_predicate`.

**Required negative control:** `x05_union_all_overlap` is `UNION ALL` of one arm
with itself; every matching order appears twice, so it must fail the bag oracle.
It must fire in **every** predicate shape against the measured data, or the run
is refused. `x01`–`x04` continue to be gated live on the teaching fixture.

### 1.2 P11 — OFFSET vs keyset pagination

Same ordered page under the declared total order `amount DESC NULLS LAST, id ASC`,
page limit 5, at depths **1, 10, 100, 1,000, 10,000**. The keyset cursor is the
row at position depth-1, so both spellings must return the identical page. Each
cell is verified against the independent oracle; additionally, a single
`REPEATABLE READ` transaction runs both pages at every depth and requires them
to contain the same rows (the "fixed snapshot" part; there are no concurrent
writers).

### 1.3 Expression index — a labelled physical arm

`CREATE INDEX orders_amount_expr_idx ON selectbench.orders ((amount + 0))` can
serve the expression shape. The query text is unchanged, so any gain belongs to
the **index**, never to the syntax: the arm is labelled `index=expression_index`
and reported separately from the default arm. Its storage (`pg_relation_size`)
and insert cost (a 500-row batch inserted with and without the index) are
recorded. The default arm runs with the index dropped.

## 2. Data, resources and topology

- Scale `deep`: 200 customers × 60 orders = 12,000 rows, no hot parent, narrow
  rows. This is a worker-local dataset choice so every protocol page depth
  addresses a real page; it changes no contract.
- Engines: PostgreSQL 17 single-node and YugabyteDB single-node, each in its own
  whole run under `ads-run-lock`.
- Budgets are the study's published ones (database 2 CPU/3 GiB, client
  2 CPU/2 GiB); the runner records live guest memory/CPU and throttling counters.
- **Coverage gaps (recorded, not silently dropped):** no YB 3-node RF=3,
  no equal-total arm, no verified tablet/leader placement, no independent
  hardware or real network. Rollup/rolldown/embedding and optimistic/pessimistic
  concurrency are **not applicable** to these syntax-only isolate pairs.

## 3. Replication, order and variability

- ≥3 independent whole runs per engine, each a fresh build/load/tag; ≥3 trials
  per cell (this plan uses 5).
- Design order is interleaved by construction (cells run in a fixed matrix order
  per run); at least one run per engine records a first/last position control.
- Across-run medians and spreads are reported; a within-run difference inside the
  two cells' own spread is `no resolved difference`, and any difference is
  `within-run only` until the whole runs agree on its sign. Verdicts go in the
  signed analysis, not the generated report.

## 4. Outputs

Per whole run: `results/<run-id>/access.json`, `access.md`, `manifest.env`,
`resources/`. Cross-run: `access-replicated.md` generated from the run JSONs.
Every report names the run tag, producing commit, image digests, dataset
fingerprint, budgets and plan structure per cell.
