# RESULT — `task-20261002T131729Z-select-join-cte`

LOW executor: model `deepseek-flash`, tool opencode/CLI, effort unknown, session
capability LOW, role executor. Goal `goal-20260922T025912Z-data-architecture-book`.
Branch `study-01/select-join-cte`. Producing commit `7e148174`. Task milestone
tag `study-01/select-v2-join-cte` (created at integration).

## What was asked

Execute the reviewed join/CTE measurement family: the six equivalent inner-join
spellings, the two CTE spellings under default / `MATERIALIZED` /
`NOT MATERIALIZED`, and a separate reuse contract, with plans and
no-resolved-difference outcomes retained; PostgreSQL 17 and YugabyteDB
single-node; the 3-node YB and placement subsets if they fit.

## What was done

Implemented and measured a self-contained join/CTE family (commits `7db4c737`,
`39c9eadc`, `f9e81517`, `94f2135f`), on a deterministic 200,000-row fact fixture
whose `rank` column is a seeded permutation of 1..N with exact thresholds
20/200/2,000/20,000/100,000/180,000. Every cell was verified against an
independent Go oracle before timing; a forced `MATERIALIZED` cell whose plan did
not keep the CTE (`CTE Scan`) was recorded unsupported rather than counted as a
default fold.

**Six whole runs, each 91/91 cells valid, 0 invalid / timeout / error:**

| run id | engine | topology | trials | inputs digest (SHA-256 of result JSON) |
|---|---|---|---|---|
| `select-join-pg1-20261004` | postgres 17.11 | single | 5 | `070df793e1354a33d7e9363877d83f9b77fc14647b9818a0631096c314f2cfac` |
| `select-join-pg2-20261004` | postgres 17.11 | single | 5 | `112577024534a9e91932f78d3ba8b9e530e0caea0aeefdf851737e578be7bab6` |
| `select-join-pg3-20261004` | postgres 17.11 | single | 5 | `05b2d7a99e38c77c84baf413a769f531179682921a2419bf73da89153b2d00ac` |
| `select-join-yb1-20261004` | yugabyte 2025.2.6 | single | 5 | `f650acde1ebb5541080360bca3e467863611e00322d3f580cefdd36ba90e0281` |
| `select-join-yb2-20261004` | yugabyte 2025.2.6 | single | 5 | `25e5f7bf01de5149bed501f089eebe900d7b462c2927eedaee302436b664d198` |
| `select-join-yb3-20261004` | yugabyte 2025.2.6 | single | 5 | `2d8394b9cc2a9a0f649367486f9d994fcfc1538e3e7d151bc75a1e8febe82e46` |

Each run has an annotated producing tag `run/study-01/<run-id>`. Seed 42 and the
fixture were held identical across runs; each run has its own seeded cell
permutation, interleaved samples and a first-cell-repeat position control.

Artifacts:

- per-run measurements only: `results/<run-id>/join.json`, `join.md`,
  `manifest.env`, `resources/`;
- replicated report (never pools engines): `reports/20261004T-select-join-replicated.md`;
- signed LOW analysis: `reports/analyses/20261004T-select-join-cte--deepseek-flash--2026-10-04.md`;
- per-engine feasibility evidence: `results/select-join-feas-{pg,pg2,pg3,yb}-20261004/`;
- premeasurement commitment: `protocols/select-join-cte-calibration-20261004.md`.

## Findings (measurements only; interpretation in the signed analysis)

- Six defaults (join+WHERE, predicate in ON, filtered derived table, filtered-
  input CTE, joined-output CTE, comma join) are plan-equivalent to join+WHERE on
  both engines. Within-run ratios 0.86–1.52x (PG) and 0.91–1.04x (YB); the 1.52x
  PG cell is same-plan and inside that run's 0.80 position-control drift.
- `MATERIALIZED` filtered-input CTE: 0.83–1.11x (PG), 0.97–1.08x (YB) — neutral.
- `MATERIALIZED` joined-output CTE: 6.8x at 10%, 27x at 1%, 88x at 0.01%
  (PG); 1.02–1.61x (YB). Absolute PG median ~37–48 ms at every threshold.
- Reuse at 16 references, `NOT MATERIALIZED` vs default: PG 0.74x at 1%
  consumer, 4.29x at 100%; YB 8.9x and 11.2x (no crossing, inlining always
  loses). Default equals `MATERIALIZED` in plan structure at references >= 2 on
  both engines.
- Rank index: PG uses `facts_rank_idx` for the 1% predicate; YB chooses a Seq
  Scan with `Storage Filter: (rank <= 2000)` (checked by hand).
- Position control (first/last within one run): PG 0.80, 0.90, 0.61; YB 0.86,
  1.35, 1.01.
- Negative controls: `x01_optional_where` detected on the edge fixture in both
  engines; `x02_join_limit_before` detected on the large fixture in both.

## Required-comparison coverage (AGENTS.md)

| Required comparison | Status here |
|---|---|
| Cluster/client sizing, per-node vs equal-total budget | **Partly covered.** Declared `DB_CPUS=2/DB_MEM=3g`, `CLIENT_CPUS=2/CLIENT_MEM=2g`; on the fixed 2 GiB podman machine the memory budgets cannot both be enforced, which the manifests record. Equal-total is single-node-not-applicable. Client is one sequential reader, recorded; no per-node throttling because there is one node. |
| Rollup / rolldown / embedding | **Not applicable.** A join-spelling experiment; it changes no derived state. Study 01's rollup/embedding designs are separate labelled arms and were not touched. |
| Colocated vs non-colocated data | **Explicit gap.** Not measured. The client DSN reaches one endpoint, the machine is a single host, and the 3-node arm is deferred; running containers on one laptop would not be placement evidence. |
| Optimistic vs pessimistic concurrency | **Not applicable.** Read-only isolated SELECT; no business invariant is mutated. |

## Gaps and failures

- 3-node YB RF=3 subset and the colocated/non-colocated pair: **not run** — the
  fixed 2 GiB machine carries one node, the client has no balanced multi-endpoint
  spread, and Study 07's topology work is deferred by the resource rule.
- No row-width sweep beyond the fixed fixture, no blended read/write phase, no
  between-consumer levels beyond 1%/100%, no p99 tail reporting, no real network,
  single host and shared cores.
- The widest threshold returns 179,911 rows, not 180,000: 89 facts have a NULL
  parent key and belong to no joined bag. The oracle agrees; the level is
  labelled 89.96%, not 90%.
- Two harness defects were found by the feasibility runs and fixed before the
  matrix (control detection inverted; wrong materialization heuristic). The
  failing feasibility outputs are retained as their evidence.

## Decision Log

See `docs/ai-work/tasks/2026/10/task-20261002T131729Z-select-join-cte/attempts/attempt-20261004T152035Z-f82109/PROGRESS.md`
(D1–D8). Notable: separate `selectjoin` schema/catalogue; generated reuse ladder;
unsupported policy cells are results, not failures; controls placed where they
can fire; PG materialization evidenced by `CTE Scan` (PostgreSQL prints no
`Materialize` node).

## Not claimed

No SELECT timing claim is promoted to the book by this task. A LOW signed
interpretation is provided and is separately attributed; book promotion waits
for the HIGH review `task-20261002T131729Z-select-synthesis` and publication
`task-20261004T142623Z-select-book-high`.
