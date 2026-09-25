---
analysis_id: 20260925T-position-control--deepseek-flash--2026-09-25
run_id: 20260925T224415Z, 20260925T224738Z
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), HIGH-capability, Deep Code CLI; effort setting not exposed. Executor of task study-04/t1-followup."
analyzed_at: 2026-09-25
inputs_digest: per-run f72ad76160cdf980 (20260925T224415Z, n3 first), 7a6a175344a9c109 (20260925T224738Z, n4 first)
repo_commit: ecd9638
status: current
headline: The 1.52x read gap between n3 and n4 was a position effect, not an architectural difference — the design that runs first reads r01 about 1.55x faster, in both orders, with an identical plan.
---

# Analysis — Study 04 position control: the r01 read gap was positional — deepseek-flash

## TL;DR

- **The gap follows the order, not the design.** With `n3_rollup_trigger` first, `n4/n3` = **0.62**; with
  `n4_rollup_app` first, `n4/n3` = **1.50**. Each design is faster on `r01` when it holds first position.
- **Per design, the position effect is ~1.55x**: n3 read 20,810 ops/s first and 13,422 second (1.55x);
  n4 read 20,081 first and 12,809 second (1.57x). The effect is the same size for both designs, which is
  what a positional artefact looks like and not what an architectural difference looks like.
- **The plan is identical in all four combinations** — same node (`Sort`), same estimate
  (`cost=22.57..22.61 rows=18`), same actual (`rows=60`). So this is not a plan effect, and with
  byte-identical `queries.sql` and equal `storage_bytes` it is not a data effect either.
- **Verdict on P2: not falsified.** There is no attributable read difference between the designs. The
  earlier 1.52x was the confound this control was run to detect.
- **Boundary:** one run per order, one host, `small`, three trials. The effect is measured, not
  explained — the plan rules out the architecture-facing explanations, not every machine-level one.

## What was measured

| Run | Realized order | n3 ops/s | n4 ops/s | n4/n3 |
|---|---|---|---|---|
| 20260925T224415Z (seed 7) | `n3_rollup_trigger, n4_rollup_app` | 20,810 | 12,809 | **0.62** |
| 20260925T224738Z (seed 101) | `n4_rollup_app, n3_rollup_trigger` | 13,422 | 20,081 | **1.50** |

`pg-single`, `small`, `verify,explain,read`, three trials, `r01_effective_config`. Every cell passed its
correctness gate.

Both runs used the same runner invocation and the same two designs; only the seed differed, and the seed
is what the shuffle uses to decide which design is loaded and measured first. The realized order is
recorded in each manifest as `design_order_pg-single`.

## Why this resolves the earlier tranche

The preceding read tranche (`study-04/v2-t1-read`) found n4 about 1.52x faster in three runs and could
not attribute it, because all three seeds realized the same order — n4 first every time. This control
supplies the missing order: when n3 goes first, the same comparison reverses to 0.62. The magnitude is
also nearly identical in both directions (1.50 vs 1/0.62 = 1.61), which is the signature of a position
effect rather than a design effect.

## Where this analysis is weak

1. **One run per order.** The direction is clean and the effect is large, but this is a single run on each
   side, not a replicated control. Three runs per order would be the completed form.
2. **No mechanism.** Identical plans, identical SQL, identical storage and identical row counts leave
   machine-level explanations (cache state, CPU frequency, container warm-up) unruled-out and untested.
3. **One query of five.** `r02`–`r05` did not show the effect in the earlier tranche, so whatever it is,
   it is specific to `r01` for reasons not established here.
4. **The effect is large enough to matter for other studies.** Any single-order read comparison in this
   harness can carry a ~1.5x positional term; claims built on such comparisons need both orders.

## Next experiments

1. Three runs per order for this pair, to put a spread on the ~1.55x.
2. Re-check the earlier read tranche's other four queries under both orders — if `r01` is uniquely
   sensitive, the cause is likely in that query's access pattern.
3. Treat the unvaried-order check as a standing precondition: a read comparison reported from one
   realized order should say so.
