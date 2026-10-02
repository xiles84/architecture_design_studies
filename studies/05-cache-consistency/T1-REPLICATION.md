# T1 replication of Study 05's cache, churn and engine claims

**Task:** `task-20261001T133759Z-t5-study05-replication` · **branch** `study-05/t1-replication`
**Written before any run.** Eight claims are in scope: `v2-14`, `v2-16`, `v3-01`, `v3-02`, `v3-03`,
`v3-04`, `v3-05`, `v3-06`. Each is `single_run_directional` today.

## 1. Claims under test and their primaries

| claim | what it asserts (primary figure) | primary run |
|---|---|---|
| `v2-14-cache-throughput-gain` | a cache on warm cacheable reads buys **2.2-3.5x** throughput and an order of magnitude of p99 (15.8-27.1 ms → 0.8-2.0 ms) | `20260921T-survey3` |
| `v2-16-strict-freshness-read-cost` | every controlled strict/relaxed pair differs by **-14.2% to +12.3%**, inside the single-trial ~20% noise floor (read throughput only) | `20260921T-survey3` |
| `v3-01-churn-crosses-real-ttl` | a 1801 s churn against the fixed 300 s TTL crossed **6.0** boundaries, 0 wrong reads, hard-expiry counter 0 | `20260923T1100Z-churn1800-through` |
| `v3-02-churn-600-through-vs-aside` | at 601 s write-through and write-aside each crossed **2.0** boundaries with 0 wrong reads; no winner claimed | `20260923T1020Z-churn600-through` (+ `…1040Z-…aside`) |
| `v3-03-equal-total-framing-labelled` | a second, separately labelled equal-total framing exists (11k vs 8.0k ops/s); the two framings are never pooled | `20260923T1150Z-res2-equaltotal-pg` |
| `v3-04-yb-single-strict-relaxed` | on YugabyteDB 1-node (RF=1) relaxed read-through 596 vs strict 514 ops/s, one run, direction not repeated | `20260923T1152Z-res2-yb1` |
| `v3-05-yb-cluster3-strict-relaxed` | on 3-node RF=3 strict 942 vs relaxed 846, the **opposite** direction from 1-node | `20260923T1157Z-res2-yb3` |
| `v3-06-colocated-vs-noncolocated-locality` | y1-coloc faster than y2-noncol on warm (882/673) and cacheable (852/520) reads; placement **not** observed | `20260923T1215Z-res2-placement-yb3` |

The unit of work is a **fresh whole run**, not a reread of the primary reports.

## 2. Feasibility on the fixed machine (decided before running)

The environment is the fixed host `host-zenbook-ux5406sa`: 4 CPUs / 2 GiB podman machine, one benchmark
at a time (`ads-run-lock`). The images `postgres:17.11`, `redis:7.4.11-alpine`, `yugabyte:2025.2.6.0-b111`
and `localhost/cachebench:1` are all present locally.

| target | plan | basis |
|---|---|---|
| `v2-14`, `v2-16` | **measured**: pg-single + Redis, two whole runs in reversed pair order | fits; v2 PG+Redis arms already ran here |
| `v3-03` | **measured**: pg-single, `db-only` and `equal-total` framings, two whole runs each | fits |
| `v3-04` | **measured**: yb-single RF=1, two whole runs in reversed order | a single-node YugabyteDB arm runs cleanly here (study-03 feasibility, 2026-10-02) |
| `v3-02` | **measured**: pg-single, 600 s churn, through and aside, two whole runs in reversed order | fits; 600 s × 2 within budget |
| `v3-01` | **partially re-tested**: the 1800 s churn is the same mechanism as `v3-02` at a longer horizon. One whole run if the session budget allows; otherwise **deferred with the 600 s pair as the re-test** and the 1800 s figure left as the primary's. | fits but ~30 min/run |
| `v3-05` | **coverage gap**: needs `yb-cluster3` RF=3, which the 2 GiB machine cannot carry (standing blocker, refined 2026-10-02) | `yb-cluster3` unrunnable here |
| `v3-06` | **coverage gap**: the primary itself is on `yb-cluster3` and its placement was never observed. The placement mechanism (tablet/leader distribution) is unsupported by the harness, and one host is not colocation evidence. | task brief: record the placement mechanism as an unsupported coverage gap |

A target that cannot be measured is recorded as a coverage gap, not escalated (task brief).

## 3. Runs

All runs use the study's own runner, scale `small`, seed 42, the published budgets, one benchmark at a
time under `ads-run-lock`, from a clean committed tree, tagged `run/05-cache-consistency/<run-id>`.

**Pair order is a control read from the executed `design_order_<topology>` in each manifest, not from the
seed.** The seed only permutes the whole list, so each target's second run must show the pair genuinely
reversed; a run whose manifest does not reverse it is kept as a repeat and a new run with another
`--order-seed` is made (this is the t4 lesson, `LESSONS_LEARNED.md` 2026-10-02).

| batch | command (per run) | carries |
|---|---|---|
| A1/A2 | `--topologies pg-single --designs owned-opt-none-none-none-coord,owned-opt-redis-through-relaxed-coord,owned-opt-redis-through-strict-coord --phases verify,warm,mixed --trials 3 --tag --order-seed <s>` | `v2-14`, `v2-16` |
| B1/B2 | as A with `--designs owned-opt-redis-through-strict-coord --resource-frame db-only` | `v3-03` baseline |
| B3/B4 | as B1/B2 with `--resource-frame equal-total` | `v3-03` equal-total |
| C1/C2 | `--topologies yb-single --designs owned-opt-redis-through-strict-coord,owned-opt-redis-through-relaxed-coord --phases verify,warm,mixed --trials 3 --tag --order-seed <s>` | `v3-04` |
| D1/D2 | `--topologies pg-single --designs owned-opt-redis-through-strict-coord,owned-opt-redis-aside-strict-coord --phases verify,calibrate,warm,mixed,churn --churn-duration 10m --trials 3 --tag --order-seed <s>` | `v3-02` |

`v2-14`'s floor comparison is `owned-opt-none-none-none-coord` (no cache, but already paying the version
and outbox machinery) against `owned-opt-redis-through-relaxed-coord`. `v2-16`'s pair is
`…-through-relaxed-…` against `…-through-strict-…`.

## 4. Predictions with named falsifiers

| id | prediction | falsifier |
|---|---|---|
| P1 | `v2-14`: the cached warm read is **≥ 1.5x** the no-cache baseline in both orders, direction reproduced, magnitude allowed to move inside the primary's band and beyond | a ratio below 1.5x, or the ordering reversed in either order |
| P2 | `v2-16`: the strict/relaxed read difference is **inside the same run's spread** and shows no consistent direction across the two orders | a difference above the within-run spread with the **same sign in both orders** |
| P3 | `v3-03`: both resource framings are recorded (`resource_framing`, budgets) and reported separately, never pooled; the ratio is **not** claimed | either run omits its framing/budget record, or the analysis pools the two |
| P4 | `v3-04`: on `yb-single` the strict/relaxed difference is inside the run's spread, no consistent direction | a same-sign difference above spread in both orders |
| P5 | `v3-02`: at 600 s both strategies cross **≥ 2** real 300 s TTL boundaries with **0 wrong reads**, and neither is claimed to win | fewer than 2 crossings, or any wrong read, or a winner named |
| P6 | correctness and provenance: every cell's `verify`/`mixed` gate passes, wrong-read accounting present beside every throughput, all runs tagged from a clean tree | a failed cell not reported, a throughput without its wrong-read count, or a dirty-tree run |
| P7 | `v3-01`: the 1800 s horizon is either re-run and crosses ≥ 6.0 boundaries with 0 wrong reads, or is recorded as deferred with the 600 s pair as the re-test | a claim that the 1800 s figure replicated when it was not re-run |
| P8 | `v3-05`, `v3-06`: recorded as coverage gaps with the resource/placement reason, no promotion | promoting either without a run on `yb-cluster3` |

**What each pass decides.** A target whose direction and order both reproduce is promoted to
`run_replicated` with its statement re-scoped to the band actually seen; a target whose relation fails
keeps its strength and gains a limit naming what did not reproduce. `v2-14`/`v3-03`'s figures stay
package-level where their own limits say so. No target's scientific question is changed; if a promotion
would need that, this task escalates instead.

## 5. Required-comparison coverage (owner rule, 2026-09-14 — recorded before running)

| requirement | coverage |
|---|---|
| 1. cluster/client sizing, equal-total-budget control | **Covered for this task by B3/B4**: the `equal-total` arm holds totals constant and charges the cache a share, separately labelled from `db-only`. Host/client budgets are the published ones; the 2 GiB machine's ceiling is a standing limit on every cell. |
| 2. rollup, rolldown and embedding | **Not applicable to these claims** (cache freshness/backend/churn). Study 05's `ref-rollup-trigger`, `ref-embedded-locked` and `ref-normalized-indexed` cells exist as other claims, not targets here. The cache key's "newest 20 donations" is a bounded embedded payload but this task does not change its maintenance mechanism. |
| 3. colocated versus non-colocated data | **Outstanding — recorded as a coverage gap** (`v3-06`): the pair needs `yb-cluster3`, unrunnable at 2 GiB, and the harness has no verified tablet/leader placement record; one host is not colocation evidence. |
| 4. optimistic versus pessimistic concurrency | **Not applicable to the target claims** (they vary freshness, backend, strategy and duration, not the guard). The study's `owned-opt` vs `owned-pess` baselines exist; the task's `v3-05` gap is engine-topology, not arbitration. |

## 6. Limits stated in advance

Single host, one benchmark at a time; declared budgets still exceed the 2 GiB machine, so no container
limit is enforceable and the VM ceiling is the real one; Redis is a shared non-authoritative cache;
absolute throughput is day- and session-specific, so only within-run ratios and order agreement travel;
the 300 s TTL is fixed; churn volumes are a closed-loop floor; the primary's own figures are single-run and
are not silently re-labelled.
