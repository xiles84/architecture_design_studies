---
analysis_id: 20260926T-v2-11-replication--deepseek-flash--2026-09-26
run_id: 20260926T234227Z
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), HIGH-capability, Deep Code CLI; effort setting not exposed. Executor of task study-03/v2-11-replication."
analyzed_at: 2026-09-26
inputs_digest: a2a6b655e14e7fa1 (run 20260926T234227Z); the baseline for comparison is run/03-reserved-seating/20260921T231328Z
repo_commit: 35abeff334f6100a9f79731b273a4880f4e9db43
status: current
headline: Study 03's first replication reproduces the document's advantage on the 100k seat map (27.8x against the original's 20.8x) while both absolutes fall 63-72 percent; the sharded variant's collapse could not be re-tested because its load failed three times.
---

# Analysis — study 03 `v2-11` replication: the multiple reproduces, the cell did not run — deepseek-flash

## TL;DR

- **The headline multiple reproduces.** On `yb-single` the stored section document read the 100k-seat map at
  **1,138 ops/s against the claim-rows design's 41**, a **27.8x** multiple where 2026-09-21 measured 20.8x
  (3.1k against 149). U1 is confirmed.
- **Both absolute figures fell far short of the original** — `l1_claim_rows` by **72%**, `l2_section_document`
  by **63%**. U3 is falsified. This is the third study in which the ratio reproduces and the numbers behind it
  do not.
- **The claim's second half could not be tested at all.** `l3_section_sharded` failed its *load* three times
  running — "terminating connection due to administrator command (SQLSTATE 57P01)" — before any measurement,
  so the sharded collapse was not re-observed. U2 is **not evaluable**, which is a coverage gap and not a
  result: the cell carries no evidence in either direction.
- **The correctness gate held where it ran**: L1 and L2 both passed 65/65 read questions against Go-computed
  truth. L3 never reached verification.
- **The tranche cannot raise the claim's strength.** The designs ran in one fixed order, so no position effect
  was controlled for, and the registry lint refuses `run_replicated` without randomized, interleaved or AB/BA
  order. The replication is recorded as a limit and here, not as a new strength, until an order-reversed run
  exists.

## The measurements

`yb-single`, `small`, phases `verify,read`, three trials per cell, medians of the three.

| Question | `l1_claim_rows` | `l2_section_document` | multiple |
|---|---:|---:|---:|
| Section map, 10 seats | 463 | 185 | 0.4x |
| Section map, 100k seats | **41** | **1,138** | **27.8x** |
| Event sections, 100k seats | 22 | 52 | 2.4x |
| Seats available for an event, 100k seats | 40 | 51 | 1.3x |
| `l3_section_sharded` | — | — | load failed, no data |

Against the 2026-09-21 baseline on the same topology: section map 100k — L1 149 → 41, L2 3.1k → 1,138.

## Prediction verdicts

| Id | Prediction | Verdict | Basis |
|---|---|---|---|
| U1 | the document is ≥4x faster on the 100k seat map | **confirmed** | 27.8x, against a 4x floor |
| U2 | the sharded variant records under 100 ops/s somewhere | **not evaluable** | the cell failed its load before measuring |
| U3 | every cell within ±40% of 2026-09-21 | **falsified** | −72% (L1) and −63% (L2) |
| U4 | the correctness gate stays green | **partial** | L1 and L2 65/65; L3 never ran |

## What this changes

1. **`v2-11` gains three limits**: the reproduced multiple, the un-run sharded half, and the fixed-order
   caveat that keeps its strength at `single_run_directional`.
2. **The pattern now has a third study.** Study 04 (rollup maintenance), study 02 (overbooking) and now
   study 03 (seat map) all show the same thing: **the ratio survives a different session and the absolute
   throughput does not.** Three studies, three workloads, three engines' worth of setups, thirteen days.
   That is now a measurement rule for this corpus rather than a study-specific curiosity.
3. **A new failure mode for the record.** Every earlier lost cell in this corpus was lost to a *correctness*
   failure — a design that answered wrongly. This one was lost to a *load* failure, three consecutive
   connection terminations, which no design property can explain. It is worth separating in the reports:
   "failed the gate" and "never got to the gate" are different kinds of absence.

## Where this analysis is weak

1. **The strongest half of the claim was not tested.** The engine-dependence — the same comparison showing
   no difference on PostgreSQL — is what makes `v2-11` interesting, and this tranche ran one engine.
2. **No order control.** One fixed design order means the 27.8x could carry a position effect of the size
   study 02 found on one of its designs (+58–76%). The multiple is large enough that a position effect is
   unlikely to reverse it, but that is an argument, not a control.
3. **The `l3` failure is unexplained beyond its symptom.** Connections were terminated by the server during
   `COPY` into `venue_seat`; whether that is a memory limit on this host, the topology script's
   configuration, or something else is not established, and it may recur on any three-design run.
4. **Wide within-cell spreads at the small tier** (182% and 196% on the 10-seat map) mean the small-tier
   numbers here are close to noise; only the 100k row is being interpreted.
5. **One session, one host, one scale**, and the 3-node topology was not measured.

## Next experiments

1. **Re-run with the design order reversed** (AB/BA), which is what would let the claim's strength move and
   is the only thing between this measurement and a registered replication.
2. **Diagnose the `l3` load failure** — it is a resource question about this host, and until it is
   understood every three-design run on `yb-single` carries the same risk.
3. **Add the PostgreSQL half** of the comparison in the same session, since engine-dependence is the claim.
