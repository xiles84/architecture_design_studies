---
analysis_id: 20260927T-v2-11-postgresql-half--deepseek-flash--2026-09-27
run_id: 20260927T202847Z
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), LOW/worker-capability session, Deep Code CLI; executing the study-03/v2-11-postgresql task."
analyzed_at: 2026-09-27
inputs_digest: b7bd31d3c19e14c6 (run 20260927T202847Z, pg-single); the YugabyteDB passes are 0e717c5c3a629cd6 and febd951548f154ca
repo_commit: bbe81338b91a7a5e8640f53dcff8c3ae15a3d5ca
status: current
headline: The PostgreSQL half of v2-11 reproduces - the same two designs differ by 0.83x there and by 24-28x on YugabyteDB - so the claim's engine-dependence is now replicated across both engines, while l3_section_sharded turns out to be a YugabyteDB-only design that the runner excludes from PostgreSQL by declaration.
---

# Analysis — study 03 `v2-11`, the PostgreSQL half — deepseek-flash

## TL;DR

- **The engine-dependence reproduces.** On `pg-single` the 100k-seat map reads **7,873 ops/s** for
  `l1_claim_rows` and **6,565** for `l2_section_document` — a multiple of **0.83x**, i.e. no attributable
  difference. The same pair on `yb-single` gives 24.4x and 27.8x. P1 and P4 are confirmed, and the contrast
  between engines is about **30x**, far beyond the 10x the protocol required.
- **`l3_section_sharded` was never supposed to run here.** I predicted it would fail on PostgreSQL, which
  would have blamed the design's load. Instead the run's manifest records
  **`yugabyte_only_designs: l3_section_sharded`** and the PostgreSQL report omits it: the runner excludes it
  **by declaration**, and the earlier passes' manifests said the same. P2 is **inapplicable**, and its
  falsifier framing in my own protocol was wrong — a defect in my pre-registration, not in the harness.
- **The correctness gate held**: both PostgreSQL cells passed 65/65 against Go-computed truth.
- **What this means for the claim**: the finding is no longer "a stored document is faster than rows" but
  "the same substitution buys a 25x multiple on one engine and nothing on the other" — measured twice on each
  engine now, and the second half had never been re-measured before this run.
- **Environment caveat, unchanged and still open**: the podman machine has 2 GiB and the declared container
  budgets cannot be enforced on it. This run is comparable with the others on this host and with no others.

## The comparison

| Question | `pg-single` l1 | `pg-single` l2 | multiple | `yb-single` multiple |
|---|---:|---:|---:|---:|
| Section map, 100k seats | 7,873 | 6,565 | **0.83x** | 24.4x / 27.8x |
| Section map, 10 seats | 36,601 | 26,265 | 0.72x | 0.4x |
| Event sections, 100k seats | 151 | 186 | 1.23x | 2.4x |
| Seats available for event, 100k | 158 | 193 | 1.22x | 1.3x |

Within-cell spreads are 1–22% here, against 207–340% on the YugabyteDB questions — this engine's
measurements are much tighter on this workload.

## Prediction verdicts

| Id | Prediction | Verdict | Basis |
|---|---|---|---|
| P1 | the PostgreSQL multiple falls between 0.6x and 1.7x | **confirmed** | 0.83x on the 100k map, 0.72–1.23x across the four |
| P2 | `l3`'s load failure is specific to the YugabyteDB topology | **inapplicable** | `l3` is declared `yugabyte_only_designs` and was never attempted on PostgreSQL |
| P3 | the correctness gate stays green | **confirmed** | 65/65 for both cells that ran |
| P4 | the engine contrast is at least 10x | **confirmed** | 0.83x against 24.4–27.8x, a factor of about 30 |

## What this changes

1. **`v2-11` gains the PostgreSQL half as replicated evidence** and a fourth whole run, with the engine
   contrast stated as the finding: the same design substitution is worth ~25x on one engine and nothing on
   the other.
2. **My own protocol carried a false premise**, and the manifest said so before I measured. The lesson is
   the one this corpus keeps relearning in other forms: **the artefacts already contain the answer if they
   are read first**. `yugabyte_only_designs` was in the baseline run's manifest as well; predicting a
   failure for a cell the runner is documented not to run wasted a prediction and, worse, framed a correct
   behaviour as a suspected defect.
3. **The `l3` question is now closed for PostgreSQL and unchanged for YugabyteDB**: the sharded design's
   load fails there in every pass, and its 2026-09-21 numbers (15.9 and 28.3 ops/s) have not been
   reproduced since.

## Where this analysis is weak

1. **Cross-engine absolutes are not comparable and none is quoted as if it were.** PostgreSQL is 200x faster
   on this question here; that is a statement about two containers on a 2 GiB machine, not a general ranking.
2. **One PostgreSQL run**, in one design order, with the pair measured in a single arrangement — the AB/BA
   control exists for the YugabyteDB pair, not for this one. The multiple is small enough that an order
   effect could matter more here than there.
3. **The environment ceiling** applies to every number above and is escalated, not resolved.
4. **The 3-node topology remains unmeasured**, so nothing here speaks to the claim's replication-factor
   dimension.

## Next experiments

1. **A second PostgreSQL pass in the reverse order**, which is cheap on this engine (this run took 15
   minutes) and would put the same AB/BA control on the engine-dependence as on the multiple.
2. **Re-run `l3` alone on `yb-single`** now that the engine is up and the sizing question is with the owner,
   to see whether its load completes with nothing else in the session.
3. **The sizing decision** — the one thing that would let any of these numbers be quoted as a resource
   envelope rather than a ceiling.
