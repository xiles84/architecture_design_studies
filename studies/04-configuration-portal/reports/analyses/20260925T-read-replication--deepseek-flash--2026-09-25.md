---
analysis_id: 20260925T-read-replication--deepseek-flash--2026-09-25
run_id: 20260925T205751Z, 20260925T210116Z, 20260925T210441Z
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), HIGH-capability, Deep Code CLI; effort setting not exposed. Executor of task study-04/t1-read."
analyzed_at: 2026-09-25
inputs_digest: per-run 1b0251938b303689 (20260925T205751Z), bf263fa88ede79ed (20260925T210116Z), 0b7dfe9e95569ce6 (20260925T210441Z)
repo_commit: 9567bdb
status: current
headline: The read-phase replication of the n3/n4 rollup pair cannot resolve P2 — n4 read r01 1.5x faster in all three runs, but the two-design shuffle produced the same order for all three seeds, so design identity and position are confounded and the plans were not captured.
---

# Analysis — Study 04 T1 read phase: the rollup-maintenance pair, three runs — deepseek-flash

## TL;DR

- **P2 is unresolved, and the tranche shows why "unresolved" is the right word rather than "falsified".**
  `n4_rollup_app` read `r01_effective_config` about **1.52x** faster than `n3_rollup_trigger` in all three
  runs (n4 medians 25,176 / 25,176 / 24,292 ops/s against n3 16,607 / 16,953 / 15,936). That is outside
  the 30% threshold the protocol set, so on the numbers alone the prediction fails.
- **But the comparison is confounded, so the effect cannot be attributed to the design.** The two-design
  seeded shuffle produced the order `n4_rollup_app, n3_rollup_trigger` for **all three** seeds
  (101, 202, 303) — n4 ran first every time. Design identity and position moved together, which is
  precisely the confound the corpus warns about, and it cannot be separated from these runs.
- **The read SQL is byte-identical and the data is the same.** `diff n3_rollup_trigger/queries.sql
  n4_rollup_app/queries.sql` is empty; the schemas differ only in comments; `storage_bytes` and the load
  timings match exactly. So the usual explanations — a different query, a bigger table, more bloat — are
  ruled out from the committed artefacts.
- **The mechanism is not evidenced.** I ran the tranche with `--phases verify,read`, so no `explain`
  phase ran and no plans were captured. A 1.5x read difference with identical SQL and identical storage
  is exactly what a plan difference would look like, and this tranche cannot test that.
- The other four reads are identical within noise: the other four medians differ by 0.4–5.0%.

## What was measured

| Read | n3 median ops/s | n4 median ops/s | n4/n3 |
|---|---|---|---|
| `r01_effective_config` | 16,607 | 25,176 | **1.52x** |
| `r02_read_key` | 33,410 | 35,087 | 1.05x |
| `r03_revision_check` | 39,344 | 39,105 | 0.99x |
| `r04_list_installations` | 15,015 | 14,953 | 1.00x |
| `r05_search_key` | 24,453 | 24,256 | 0.99x |

Medians over the three runs; n3 per-run `r01` = 16,607 / 16,953 / 15,936, n4 = 25,591 / 25,176 / 24,292.
Every cell passed its correctness gate. All three runs are `pg-single`, `small`, three trials.

## Verdict on the protocol's predictions

| Id | Prediction | Verdict | Basis |
|---|---|---|---|
| P2 | no attributable read difference between n3 and n4 | **unresolved (confounded)** | a 1.52x gap on `r01` in 3/3 runs, but n4 held the same first position in all three, and no plans were captured |

P1, P3 and P4 belong to the write tranche and are unaffected.

## A finding about the method itself

**A seeded shuffle is not a varied order.** The runner's `designs_for` derives its order from
`ORDER_SEED`, and with a three-design list the three seeds did produce different orders (run 1:
`n4, x2, n3`; runs 2 and 3: `x2, n4, n3`). With the two-design list here, all three seeds produced the
same order. The realized order is recorded in every manifest — which is what makes this visible after
the fact — but "randomized order" in a replication record describes the *method*, not what happened.
A replication record should be read against the manifest's `design_order_<topology>` line, and a tranche
whose realized order never varies has not randomized anything.

**This also qualifies the write tranche.** The registered `v2-13` replication records
`order: randomized`, but its three runs realized only two distinct orders (seeds 202 and 303 coincided).
Its limits now say so.

## Where this analysis is weak

1. **No plans.** Without `explain`, the 1.52x `r01` gap has no mechanism attached. That is a defect in
   my phase selection, not a property of the designs.
2. **No position control.** The pair was never run in both orders. Doing so is the single experiment
   that would separate design from position here, and it is cheap.
3. **One scale, one engine, one host, three runs.** `run-level` replication at best.
4. **`r01` may simply be the first query of the first cell.** The measured within-cell position penalty
   from the write tranche (the first trial 1.5–5.6x slower than the cell median) is a different effect —
   but a first-cell effect of the same family has not been excluded, and could not be from this tranche.
5. **The observed gap is consistent, which is easy to over-read.** Three runs agreeing on a confounded
   comparison is still a confounded comparison.

## Next experiments

1. **The position control, first.** Run the pair as `n4, n3` and as `n3, n4` with the same seeds, with
   `explain` enabled. If the gap follows the order rather than the design, `r01` is a harness effect; if
   it follows the design, it is a real read cost of the trigger-maintained layout and needs a mechanism.
2. Register the writer-sweep crossing (4 levels, reversal between 1 and 4 writers) — blocked on a
   decision, because the v6 `crossover` lint requires five levels for a claimed reversal and this
   measurement has four.
3. Extend `designs_for` so a small design list can still realize distinct orders, or record realized
   order distinctness as a field the replication lint checks.
