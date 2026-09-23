---
analysis_id: 20260923T1215Z-resources-engines--deepseek-flash--2026-09-23
run_id: 20260923T1215Z-res2-placement-yb3
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), HIGH-capability, Deep Code CLI; effort setting not exposed. Author of these runs."
analyzed_at: 2026-09-23
inputs_digest: a18817d592f767c5
repo_commit: 7ca602457d5336dde8f481093b1637a21f11a6b4
supersedes:
status: current
headline: The equal-total arm is labelled and gate-clean, strict/relaxed and colocated/non-colocated pairs now run on one host with three spread endpoints, but the runner cannot produce a tablet/leader distribution, so the placement evidence clause is a gap.
---

# Analysis — Study 05 v2 equal-total, engines and placement — deepseek-flash

> Deliverable of `task-20260922T200831Z-study05-v2-resources-engines`, from
> `studies/05-cache-consistency/HANDOFF-AMENDMENT-01-V2.md` (tag `study-05/v4-handoff`).

## TL;DR

- **The equal-total arm is measured and separately labelled.** `resource_framing: equal-total`
  versus `db-only` on the same design and topology: warm read 10,837 vs 8,027 ops/s, cacheable-99
  7,876 vs 8,130, app-99 7,945 vs 7,873 — the two arms are within ~1.03x on the mixed reads and are
  never pooled.
- **The strict/relaxed pair now runs per engine** with all three endpoints spread: on yb-single
  (RF=1) strict is 514 vs relaxed 596 app-99 ops/s; on yb-cluster3 (RF=3) strict is 942 vs relaxed
  846. The direction differs between engines and each is a single run, so neither is attributed.
- **The colocated/non-colocated pair is faster colocated**: y1 882/852/558 vs y2 673/557/387
  (warm/cacheable/app), ~1.3–1.5x. But the runner **could not produce a tablet/leader distribution**
  — `placement-evidence.txt` holds only the three server names — so the acceptance's "tablet/leader
  distribution recorded" is a **gap**, and one host cannot evidence colocation anyway.
- A real harness defect was found and fixed: every yb-cluster3 cell failed with
  `sslmode is invalid` because a comma-separated DSN list was handed to one pgx pool. `pgxdb.OpenSpread`
  now opens one pool per endpoint and spreads operations; `connection_nodes` is recorded in every
  result (1 or 3).

## 1. Measured cells (PostgreSQL 17.11 / YugabyteDB 2025.2.6 + Redis, scale small, 3 trials)

| Arm | Scenario | Endpoints | Gate | strict violations | impossible | warm ops/s | cacheable-99 | app-99 |
|---|---|---|---|---|---|---|---|---|
| db-only | owned-opt-redis-through-strict | 1 | pass | 0 | 0 | 8,027 | 8,130 | 7,873 |
| equal-total | owned-opt-redis-through-strict | 1 | pass | 0 | 0 | 10,837 | 7,876 | 7,945 |
| yb-single RF=1 | strict | 1 | pass | 0 | 0 | 9,100 | 6,208 | 514 |
| yb-single RF=1 | relaxed | 1 | pass | 0 | 0 | 6,353 | 8,071 | 596 |
| yb-cluster3 RF=3 | strict | 3 | pass | 0 | 0 | 9,189 | 10,406 | 942 |
| yb-cluster3 RF=3 | relaxed | 3 | pass | 0 | 0 | 8,179 | 8,975 | 846 |
| yb-cluster3 RF=3 | y1 colocated | 3 | pass | 0 | 0 | 882 | 852 | 558 |
| yb-cluster3 RF=3 | y2 non-colocated | 3 | pass | 0 | 0 | 673 | 557 | 387 |

All cells: `failed_cells: []`, `gate.passed = true`, `strict_contract_violations = 0`,
`impossible_cache_values = 0`, and the ledger assertions recorded no mismatch. Redis evictions were
0 at this scale (the cache fit the working set), so no eviction effect is claimed.

## 2. What the arms show

- **Equal-total vs db-only.** The equal-total frame charges a share of the whole machine to the
  Redis container; the mixed-read throughput is essentially unchanged, which is what a frame change
  should do here (the cache and the database are the same two containers either way). The value of
  the arm is that it is *labelled* — `resource_framing` is recorded in every result — so a later
  reader cannot pool it with the db-only numbers, which the v1 review asked for.
- **Strict vs relaxed, per engine.** Relaxed is faster on yb-single and slower on yb-cluster3.
  Both are single runs; the correct statement is that the pair is now measured on both engines with
  the endpoints spread, not that either freshness policy wins.
- **Colocated vs non-colocated.** y1's child key puts one donor's donations in one tablet; y2 hashes
  them across tablets. y1 is faster on all three reads (~1.3–1.5x). This is engine evidence that
  key choice changes locality, on one machine.

## 3. The harness defect this task found and fixed

Before the fix, **every** yb-cluster3 cell failed at connect: the runner hands the harness a
comma-separated DSN list and a single `pgxpool.ParseConfig` rejects it (`sslmode is invalid`). Two
consequences mattered: the cluster arm could not be measured at all, and a "fix" that used only the
first endpoint would have measured one node's SQL layer while calling it a cluster. `pgxdb.OpenSpread`
opens one pool per endpoint and round-robins operations (a transaction stays on the pool it began
on), and `connection_nodes` is recorded in every result. The pre-fix failed runs are kept, untagged,
in `results/`; the matrix above is one coherent commit.

## 4. Where this measurement is weak

1. **Placement distribution is a gap.** `placement-evidence.txt` records only the servers
   (`yb-n1|1`, `yb-n2|1`, `yb-n3|1`); the `yb-admin list_tablets` call produced nothing on this
   image. So the colocated/non-colocated speed difference is measured but *where the rows physically
   landed is not evidenced*.
2. **One host.** All containers run on one WSL2 machine, so this is not evidence of data colocation
   in the multi-host sense and the report says so.
3. **Single runs, one scale.** Each arm is one cell at scale small; the strict/relaxed direction is
   not repeated.
4. **Endpoint spread is round-robin by construction**, recorded as a count (1 or 3); per-endpoint
   operation counts are not instrumented.
