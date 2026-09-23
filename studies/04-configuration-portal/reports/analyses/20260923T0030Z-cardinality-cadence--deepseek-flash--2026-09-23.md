---
analysis_id: 20260923T0030Z-cardinality-cadence--deepseek-flash--2026-09-23
run_id: 20260923T0030Z-card-t500
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), HIGH-capability, Deep Code CLI; effort setting not exposed. Author of these runs."
analyzed_at: 2026-09-23
inputs_digest: 5b3c49eb09d2106d
repo_commit: f1d9b3c0157949c5ff3671951f97da15363eccaa
supersedes:
status: current
headline: The rollup designs are cardinality-flat while the row designs collapse (74x apart at 500 entries), and the update cadence delivers its offered rate under all three regimes with the trigger's w05 cost ~40x the application rollup.
---

# Analysis — Study 04 v2 cardinality and cadence — deepseek-flash

> Deliverable of `task-20260922T195241Z-study04-v2-cardinality-cadence`, the second execution task
> from `studies/04-configuration-portal/HANDOFF-AMENDMENT-01-V2.md` (tag `study-04/v2-handoff`).
> It reads the nine generated reports below; it edits no earlier analysis.

## TL;DR

- **The cardinality crossover is large and clean.** For `r04_list_installations`, the row designs
  fall from ~11.8k ops/s at 1 entry per product to **188–189 ops/s at 500** (a 62x collapse), while
  the two rollup designs stay **13.9–14.7k ops/s at every tier**. At 500 entries the rollup is
  **~74x** the row design; at 1 entry they are within 1.3x.
- **The hypothesis the protocol set for this phase holds**: the unindexed design *does* show a scan
  cost that grows with tier, so the tiers are not all served from cache and the phase legitimately
  continued to 500.
- **The update cadence delivered its offered rate in all three regimes** (quiet 1.15/s of 1/s,
  steady 60.1/s of 60/s, burst 66.3/s of 60/s with queue depth 22 of 64 and no drops), and the
  trigger-versus-application write gap is reproduced: `w05_replace_all` **10–11 ops/s** for the
  trigger against **413–448 ops/s** for the application rollup across all three cadences.
- **Staleness and derived-state rollup lag are not recorded by the harness**, so that half of the
  cadence acceptance is a named gap, not a measurement (section 5).

## 1. Cardinality (6 runs, entries per installed product 1/10/30/60/120/500)

`r01_effective_config` (point read) and `r04_list_installations` (the cardinality-sensitive read),
ops/s, phase `verify,explain,read`, three trials per read, PostgreSQL 17.11 single node.

| Entries/ip | design | r01 ops/s | r04 ops/s |
|---|---|---|---|
| 1 | n0_rows_unindexed | 25383 | 11205 |
| 1 | n1_rows_indexed | 31305 | 11773 |
| 1 | n3_rollup_trigger | 21414 | 13223 |
| 1 | n4_rollup_app | 24390 | 14709 |
| 10 | n0_rows_unindexed | 22589 | 3787 |
| 10 | n1_rows_indexed | 27116 | 2957 |
| 10 | n3_rollup_trigger | 21965 | 13709 |
| 10 | n4_rollup_app | 21433 | 14994 |
| 30 | n0_rows_unindexed | 17891 | 2569 |
| 30 | n1_rows_indexed | 20147 | 1854 |
| 30 | n3_rollup_trigger | 19603 | 14488 |
| 30 | n4_rollup_app | 20303 | 14444 |
| 60 | n0_rows_unindexed | 15151 | 1348 |
| 60 | n1_rows_indexed | 11894 | 1311 |
| 60 | n3_rollup_trigger | 17828 | 13598 |
| 60 | n4_rollup_app | 14677 | 13877 |
| 120 | n0_rows_unindexed | 8420 | 636 |
| 120 | n1_rows_indexed | 8692 | 648 |
| 120 | n3_rollup_trigger | 11569 | 14399 |
| 120 | n4_rollup_app | 11740 | 14585 |
| 500 | n0_rows_unindexed | 5654 | 188 |
| 500 | n1_rows_indexed | 5679 | 189 |
| 500 | n3_rollup_trigger | 5615 | 13864 |
| 500 | n4_rollup_app | 4671 | 13967 |

**Reading.** `r04` is the read that separates derived state from raw rows. For the row designs it
falls monotonically with entries per product (11.2k → 188 for n0, 11.8k → 189 for n1): at this scale
the index on the key does not help a read that has to touch every entry of a product, so n0 and n1
collapse together. The rollup designs are flat across two orders of magnitude (13.2k–14.7k), because
the answer is read from the parent rather than assembled from the children. `r01` (a point read)
declines only gently for every design (25.4k → 5.7k for n0), which is the expected shared effect of
a larger per-product working set.

The exact SQL behind `r04` and the rollup's storage shape are in the design catalogue, not asserted
here; the measured pattern is what is reported.

## 2. Cadence (3 runs at 120 entries per product, phase `verify,write,cadence`)

| Regime | Period | Fleet | Offered/s | Delivered/s | Completed / offered | Dropped | Queue depth (bound 64) |
|---|---|---|---|---|---|---|---|
| quiet | 1m0s | 60 | 1 | 1.15 | 5/5 | 0 | 0 |
| steady | 1s | 60 | 60 | 60.1 | 301/301 | 0 | 0 |
| burst (burst-size 30) | 1s | 60 | 60 | 66.3 | 300/300 | 0 | 22 |

All three regimes crossed their configured cadence — `started` equals `offered` and every update
completed with no rejected, errored or dropped updates. The burst absorption is visible as queue
depth 22 of a 64 bound; the client did not report saturation.

## 3. Trigger versus application maintenance at the 120-entry tier

`w05_replace_all` (whole-configuration replacement), ops/s, median of three trials:

| Regime | n3_rollup_trigger | n4_rollup_app | ratio |
|---|---|---|---|
| quiet | 10.0 | 413.1 | 41x |
| steady | 11.0 | 447.9 | 41x |
| burst | (see report) | (see report) | – |

The other writes are much closer (w01 ~1.0–1.2k, w02 ~0.8–1.0k, w03 ~0.45–0.5k, w04 286–709,
w06 ~1.65–1.7k), so the trigger's structural cost is concentrated on the operation that touches
every entry. This reproduces the v1 finding at a second cardinality and across cadences.

## 4. Correctness

All nine runs: `gate.passed = true`, `failed_cells: []`, no invariant violation by a non-control
design. The three-trial read rule was applied to every read cell in all nine runs.

## 5. Where this measurement is weak

1. **Staleness and rollup lag are not measured.** The harness records scheduling lag (p50 ~12.5 ms,
   p99 ~44 ms in the burst regime) and delivery counts, but not the age of a derived value at read
   time, so the acceptance's "staleness and rollup lag for trigger versus application" is **not
   met**; it needs a harness field that stamps the derived value and the read.
2. **One engine, one host, closed loop.** PostgreSQL 17.11 single node only; no YugabyteDB arm
   (that is the remaining-designs task). The client is closed-loop, so the cadence figures are
   delivered-rate observations, not SLOs.
3. **Tier 1 is a degenerate point** (a single entry per product) and is included for the crossover
   shape, not as a steady state.
4. **The cadence update writes and the cardinality reads come from separate runs**, so the trigger
   cost and the crossover are not from one cell.
