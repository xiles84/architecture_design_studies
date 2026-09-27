---
analysis_id: 20260927T-close-out-tranche--deepseek-flash--2026-09-27
run_id: 20260927T205830Z (pg-single, pair reversed), 20260927T211303Z (yb-single, l3 alone)
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), LOW/worker-capability session, Deep Code CLI; executing the study-03/close-out task."
analyzed_at: 2026-09-27
inputs_digest: 2e2dfea128e77233 (pg pair reversed), 6885b9cdb94aedc3 (l3 alone); earlier passes 0e717c5c3a629cd6, febd951548f154ca, b7bd31d3c19e14c6
repo_commit: ae8684bf195a5f6bc7a7a7fcf204c672ba88d374
status: current
headline: The PostgreSQL half survives having its pair reversed (0.88x against 0.83x), giving the engine contrast an AB/BA control on both engines; and l3_section_sharded loads, verifies and reproduces its collapse when run alone, which makes its three earlier failures a session effect rather than a design property.
---

# Analysis — study 03 close-out tranche: both open items answered — deepseek-flash

## TL;DR

- **The PostgreSQL half has its AB/BA control.** With `l2` before `l1` the 100k-seat multiple is **0.88x**,
  against **0.83x** with the order reversed. P5, P6 and P7 are confirmed, and the engine contrast — roughly
  25x on YugabyteDB against essentially none on PostgreSQL — is now measured in **both orders on both
  engines**.
- **`l3_section_sharded` loads, verifies and reproduces when run alone.** D1 and D2 are confirmed: its load
  completed with nothing before it, it passed **65/65** against Go-computed truth, and its numbers land
  where 2026-09-21 put them (3,636 ops/s on the 10-seat map against 4.3k then; 23 against 28.3 on
  seats-available per section at 10 seats).
- **That makes the three earlier failures a session effect, not a design property.** The same design, the
  same topology, the same scale: it dies when other cells have loaded in the session and completes when none
  has, on a machine whose container budgets cannot be enforced.
- **It does not rewrite `v2-11`.** The pre-registration said Part 2 was a diagnostic and would not be folded
  into the claim's numbers, and that commitment is kept: the claim's sharded element stays as the 2026-09-21
  run left it, and registering the dedicated run as its replication is flagged as a decision rather than
  taken after seeing the result.
- **Correctness green everywhere**: 65/65 for both PostgreSQL cells and for `l3` alone.

## Part 1 — the PostgreSQL pair, both orders

| Order | `l1` 100k map | `l2` 100k map | multiple | vs YugabyteDB |
|---|---:|---:|---:|---:|
| `l1` then `l2` (2026-09-27, earlier) | 7,873 | 6,565 | 0.83x | 24.4–27.8x |
| `l2` then `l1` (this pass) | 7,144 | 6,317 | **0.88x** | 24.4–27.8x |

A 6% difference between orders; within-cell spreads 4–21%. Verdicts: **P5 confirmed**, **P6 confirmed**
(the contrast is ~28–32x, far above the 10x bar), **P7 confirmed**.

## Part 2 — `l3_section_sharded` alone

| Question | `l3` alone | what the 2026-09-21 matrix recorded |
|---|---:|---:|
| Section map, 10 seats | 3,636 | 4.3k |
| Section map, 100k seats | 1,702 | 2.2k |
| Seats available per section, 10 seats | 23 | 28.3 |
| Event sections, 100k seats | 13 | — |

Verification: **65/65**. Verdicts: **D1 confirmed** (it loads with an empty session around it),
**D2 confirmed** (its answers are right).

## Prediction verdicts

| Id | Prediction | Verdict | Basis |
|---|---|---|---|
| P5 | the PostgreSQL multiple stays in 0.6–1.7x when reversed | **confirmed** | 0.88x |
| P6 | the engine contrast stays ≥10x | **confirmed** | ~28–32x |
| P7 | PostgreSQL correctness stays green | **confirmed** | 65/65 and 65/65 |
| D1 | `l3` loads when nothing runs before it | **confirmed** | it loaded |
| D2 | its answers verify if it loads | **confirmed** | 65/65 |

## What this changes

1. **`v2-11` gains the reversed PostgreSQL pass** as a limit and a fifth whole run: the "no attributable
   difference on PostgreSQL" half now has the same order control as the multiple, so the engine contrast is
   the claim's most replicated element rather than its least.
2. **The `l3` question is answered but deliberately not registered.** Its failures are explained — session
   accumulation inside a machine that cannot honour its declared budgets — and the design is measurable by a
   dedicated run. Whether that run becomes the claim's sharded-element replication is a **decision**, because
   the protocol excluded it in advance and a worker may not quietly upgrade a diagnostic into evidence.
3. **A cheaper route for future matrices exists**: a design whose load is memory-heavy can be measured on
   its own rather than inside a matrix, and the difference between the two arrangements is itself a
   measurement of the environment rather than of the design.

## Where this analysis is weak

1. **Each part is a single run**, so neither new multiple has repetitions at its own level; the PostgreSQL
   control is one run per order, and `l3` was measured once.
2. **`l3`'s 2026-09-21 comparison is loose** — different run, same host, and the "matrix" values are quoted
   from a report rather than re-derived, so the agreement is order-of-magnitude, not a replication.
3. **The environment ceiling is unresolved** and applies to every number here; the l3 success under it is
   evidence about *this* machine's behaviour, not about the design's cost in general.
4. **Still unmeasured:** the 3-node topology, and any of this at a second scale.

## Next experiments

1. **Decide the `l3` registration** (owner or HIGH): if the dedicated run is accepted as the sharded
   element's replication, `v2-11`'s most interesting element gains evidence it has lacked since 2026-09-21.
2. **Repeat `l3` alone once** to give that decision a spread rather than a point.
3. **The sizing decision**, which would let any of these numbers be quoted as a resource envelope rather
   than a ceiling.
