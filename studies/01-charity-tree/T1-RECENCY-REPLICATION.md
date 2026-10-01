# T1: replicating Study 01's remaining claims — `v2-04` and `v2-05`

**Task:** `task-20261001T133759Z-t1-study01-replication` · **branch** `study-01/t1-remaining`

**Goal.** The two claims Study 01 still carries at `single_run_directional` and `repeated_controlled` were
measured once each, in one session each. This tranche re-runs their protocols as **whole runs** and asks
whether what they assert survives.

| claim | strength now | what it asserts | its primary |
|---|---|---|---|
| `v2-04-recency-maintained-index` | `single_run_directional` | the recency flag works, but an index on the rollup column answers the same question faster: `d20_recency_flag` against `d22_recency_rollup_idx` in the same run (README: 2.9x on PostgreSQL, 7.9x on YugabyteDB) | `20260916T090036Z-v3` (v4 recency protocol, 3 repeats) |
| `v2-05-d3-read-advantage-repeated` | `repeated_controlled` | `d3_flattened_fk`'s read advantage over `d2_normalized_indexed` survives five fresh loads in one session | `20260913T125342Z-v3` (5 trials, one session) |

**What is held constant.** Study 01's own v3/v4 protocols and configuration: scale `small`, `duration 3s`,
`warmup 1s`, the client budget (2 CPUs / 2 GiB) and the database budget on the same fixed host, the study's
own harness and its correctness gates including `AuditRecency`. **The factors varied are the run and the
design order** (the runner alternates AB/BA per trial), plus a **scoped design set**: this tranche runs five
designs on `pg-single` instead of the full matrix, and the manifest records the filter.

**Design.**

*Designs (5):* `d1_normalized_minimal` (the study's baseline), `d2_normalized_indexed`, `d3_flattened_fk`
(the `v2-05` pair), `d20_recency_flag`, `d22_recency_rollup_idx` (the `v2-04` pair).
*Experiments:* `reads` (the twelve core queries and the recency group, i.e. what `v2-05` turns on) and
`recency-reads` (q13–q16 in both window regimes, what `v2-04` turns on).
*Topology:* `pg-single` only.
*Runs:* two whole runs, each `--trials 2` so the runner's own AB/BA alternation gives both orders inside each
run — four fresh loads per configuration across the two runs, against the primaries' three (`v2-04`) and five
within one session (`v2-05`).

**Coverage gap, stated before running:** the primaries' `yb-single`/`yb-cluster3` arms, including the 7.9x
YugabyteDB figure `v2-04`'s README quotes, are **not** replicated. The podman machine has 2 GiB and the owner
has fixed the machine, so that arm is not measurable here. Whatever the PostgreSQL arm shows, the claims keep a
limit naming the YugabyteDB half as unreplicated.

## Required-comparison coverage (owner rule, 2026-09-14 — recorded before running)

| requirement | coverage |
|---|---|
| 1. cluster and client sizing, equal-total-budget control | Budgets are the study's published ones and unchanged from the primaries. **Not applicable:** one topology compared with its own past, no cluster-size comparison, so there is no equal-total-budget arm. The scoped design filter is recorded in the manifest. |
| 2. rollup, rolldown and embedding | **Covered:** `d22_recency_rollup_idx` is the rollup-index design this claim is about, and `d2`/`d3`/`d20` are the normalized, flattened and flag designs it is compared with. Rolldown has no Study 01 design and stays **outstanding**; embedding is not part of these two claims. |
| 3. colocated and non-colocated data | **Not applicable** (single-node PostgreSQL, no replication) and **outstanding** for the YugabyteDB arms this tranche excludes. |
| 4. optimistic and pessimistic concurrency | **Not applicable**: these are read claims with no concurrency arm; the recency maintenance path's locking is touched only by the maintenance experiment, which this tranche does not run. |

## Prediction

Recorded before measurement. Bands and orderings, not point figures — three replications in this corpus have
now moved point figures by 30–45% while their bands held.

| Id | Prediction | Falsifier |
|---|---|---|
| L1 | **The `v2-05` ordering holds in every trial and both runs**: `d3`'s core read score is above `d2`'s. | Any trial where `d2` ≥ `d3`. |
| L2 | **`v2-05`'s band holds**: `d3`/`d2` read-score ratio between **1.2x and 3x**, against the primary's reported advantage (its own limits call the insert side imprecise, not the read side). | A ratio outside that band in either run's trial medians. |
| L3 | **The `v2-04` ordering holds**: on the recency queries, `d22` is faster than `d20` in both runs. | `d20` at or above `d22` in either run. |
| L4 | **`v2-04`'s PostgreSQL speedup is real and smaller than a factor of ten**: `d22` over `d20` on the recency queries is **≥ 1.5x** and **< 10x**, against the primary's 2.9x. | Below 1.5x or above 10x. |
| L5 | **The flag does not beat the rollup index**, i.e. `d20` does not overtake `d22` even where its write path is cheaper; and `d1` remains the slowest of the five on the core queries. | `d20` ≥ `d22`, or `d1` not the slowest on the core read score. |
| L6 | **Correctness gates stay green**: `AuditRecency` and the study's other gates pass in every cell, no failed cells, both runs from a clean committed tree and tagged by the runner. | Any failed cell or a dirty tree. |

**What a pass decides.** `v2-05` can move to `run_replicated` (its five-fresh-loads evidence becomes four
fresh loads across two whole runs plus its existing within-run repeats). `v2-04` can move to `run_replicated`
for its **PostgreSQL** half, with its statement re-scoped to the replicated band and a limit naming the
YugabyteDB half as unreplicated. A failure leaves each claim's strength alone and adds a limit naming the
relation that did not reproduce.

**Known limits, stated in advance.** Five designs of the study's twenty-five; `pg-single` only; small scale;
two runs on one host; the read score is the study's blended twelve-query geomean, so these are claims about
that blend.
